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
		if !tbl.Allow(key) {
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

	// 使用滑动窗口算法
	now := time.Now().Unix()
	windowStart := now - int64(rrl.window.Seconds())

	// 清理过期的记录
	rrl.client.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart, 10))

	// 获取当前窗口内的请求数
	count, err := rrl.client.ZCard(ctx, redisKey).Result()
	if err != nil {
		return false
	}

	if int(count) >= rrl.rate {
		return false
	}

	// 添加当前请求
	rrl.client.ZAdd(ctx, redisKey, &redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d", now),
	})

	// 设置过期时间
	rrl.client.Expire(ctx, redisKey, rrl.window)

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
		return &NoOpLimiter{}, nil
	}

	switch config.Type {
	case "memory":
		return NewTokenBucketLimiter(rate.Limit(config.Rate), config.Burst), nil
	case "redis":
		rdb := redis.NewClient(&redis.Options{
			Addr: config.RedisAddr,
			DB:   config.RedisDB,
		})
		return NewRedisRateLimiter(rdb, config.Rate, config.Window), nil
	default:
		return nil, fmt.Errorf("unsupported rate limiter type: %s", config.Type)
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
