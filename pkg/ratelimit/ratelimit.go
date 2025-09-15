package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
	"student-management-system/pkg/logger"
)

// RateLimiter 限流器接口
type RateLimiter interface {
	Allow(key string) bool
	Middleware() gin.HandlerFunc
}

// TokenBucketLimiter 基于令牌桶的限流器
type TokenBucketLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewTokenBucketLimiter 创建新的令牌桶限流器
func NewTokenBucketLimiter(r rate.Limit, b int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// getLimiter 获取或创建限流器
func (tbl *TokenBucketLimiter) getLimiter(key string) *rate.Limiter {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()

	limiter, exists := tbl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(tbl.rate, tbl.burst)
		tbl.limiters[key] = limiter
	}

	return limiter
}

// Allow 检查是否允许请求
func (tbl *TokenBucketLimiter) Allow(key string) bool {
	return tbl.getLimiter(key).Allow()
}

// Middleware 返回限流中间件
func (tbl *TokenBucketLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		logger.WithFields(logger.Fields{
			"client_ip": key,
		}).Debug("检查令牌桶限流")

		if !tbl.Allow(key) {
			logger.WithFields(logger.Fields{
				"client_ip": key,
			}).Warn("令牌桶限流触发")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		logger.WithFields(logger.Fields{
			"client_ip": key,
		}).Debug("令牌桶限流检查通过")
		c.Next()
	}
}

// RedisRateLimiter 基于Redis的分布式限流器
type RedisRateLimiter struct {
	client    *redis.Client
	rate      int
	window    time.Duration
	keyPrefix string
}

// NewRedisRateLimiter 创建新的Redis限流器
func NewRedisRateLimiter(client *redis.Client, rate int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		client:    client,
		rate:      rate,
		window:    window,
		keyPrefix: "rate_limit:",
	}
}

// Allow 检查是否允许请求
func (rrl *RedisRateLimiter) Allow(key string) bool {
	redisKey := rrl.keyPrefix + key
	ctx := rrl.client.Context()

	logger.WithFields(logger.Fields{
		"key":       key,
		"redis_key": redisKey,
		"rate":      rrl.rate,
		"window":    rrl.window,
	}).Debug("检查Redis限流")

	// 使用滑动窗口算法
	now := time.Now().Unix()
	windowStart := now - int64(rrl.window.Seconds())

	// 清理过期的记录
	err := rrl.client.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart, 10)).Err()
	if err != nil {
		logger.WithError(err).Error("清理过期限流记录失败")
		return false
	}

	// 获取当前窗口内的请求数
	count, err := rrl.client.ZCard(ctx, redisKey).Result()
	if err != nil {
		logger.WithError(err).Error("获取限流计数失败")
		return false
	}

	if int(count) >= rrl.rate {
		logger.WithFields(logger.Fields{
			"key":           key,
			"current_count": count,
			"rate_limit":    rrl.rate,
		}).Warn("Redis限流触发")
		return false
	}

	// 添加当前请求
	err = rrl.client.ZAdd(ctx, redisKey, &redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	}).Err()
	if err != nil {
		logger.WithError(err).Error("添加限流记录失败")
		return false
	}

	// 设置过期时间
	err = rrl.client.Expire(ctx, redisKey, rrl.window).Err()
	if err != nil {
		logger.WithError(err).Warn("设置限流记录过期时间失败")
	}

	logger.WithFields(logger.Fields{
		"key":           key,
		"current_count": count + 1,
		"rate_limit":    rrl.rate,
	}).Debug("Redis限流检查通过")

	return true
}

// Middleware 返回限流中间件
func (rrl *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !rrl.Allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Config 限流配置
type Config struct {
	Enabled   bool          `yaml:"enabled"`
	Type      string        `yaml:"type"` // "memory" or "redis"
	Rate      int           `yaml:"rate"`
	Burst     int           `yaml:"burst"`
	Window    time.Duration `yaml:"window"`
	RedisAddr string        `yaml:"redis_addr"`
	RedisDB   int           `yaml:"redis_db"`
}

// NewRateLimiter 根据配置创建限流器
func NewRateLimiter(config Config) (RateLimiter, error) {
	if !config.Enabled {
		logger.Info("限流功能已禁用")
		return &NoOpLimiter{}, nil
	}

	logger.WithFields(logger.Fields{
		"type":   config.Type,
		"rate":   config.Rate,
		"burst":  config.Burst,
		"window": config.Window,
	}).Info("初始化限流器")

	switch config.Type {
	case "memory":
		logger.Info("创建内存令牌桶限流器")
		return NewTokenBucketLimiter(rate.Limit(config.Rate), config.Burst), nil
	case "redis":
		logger.WithFields(logger.Fields{
			"redis_addr": config.RedisAddr,
			"redis_db":   config.RedisDB,
		}).Info("创建Redis限流器")
		rdb := redis.NewClient(&redis.Options{
			Addr: config.RedisAddr,
			DB:   config.RedisDB,
		})
		return NewRedisRateLimiter(rdb, config.Rate, config.Window), nil
	default:
		err := fmt.Errorf("unsupported rate limiter type: %s", config.Type)
		logger.WithError(err).Error("不支持的限流器类型")
		return nil, err
	}
}

// NoOpLimiter 空操作限流器（用于禁用限流）
type NoOpLimiter struct{}

func (nol *NoOpLimiter) Allow(key string) bool {
	return true
}

func (nol *NoOpLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
