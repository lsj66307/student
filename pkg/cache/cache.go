package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache 创建新的Redis缓存
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// buildKey 构建缓存键
func (rc *RedisCache) buildKey(key string) string {
	if rc.prefix == "" {
		return key
	}
	return rc.prefix + ":" + key
}

// Get 获取缓存值
func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, rc.buildKey(key)).Result()
}

// Set 设置缓存值
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, rc.buildKey(key), value, expiration).Err()
}

// Del 删除缓存
func (rc *RedisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := make([]string, len(keys))
	for i, key := range keys {
		redisKeys[i] = rc.buildKey(key)
	}

	return rc.client.Del(ctx, redisKeys...).Err()
}

// Exists 检查缓存是否存在
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.client.Exists(ctx, rc.buildKey(key)).Result()
	return result > 0, err
}

// CacheMiddleware 缓存中间件配置
type CacheMiddleware struct {
	cache      Cache
	defaultTTL time.Duration
	keyFunc    func(*gin.Context) string
}

// NewCacheMiddleware 创建缓存中间件
func NewCacheMiddleware(cache Cache, defaultTTL time.Duration) *CacheMiddleware {
	return &CacheMiddleware{
		cache:      cache,
		defaultTTL: defaultTTL,
		keyFunc:    defaultKeyFunc,
	}
}

// SetKeyFunc 设置缓存键生成函数
func (cm *CacheMiddleware) SetKeyFunc(keyFunc func(*gin.Context) string) {
	cm.keyFunc = keyFunc
}

// defaultKeyFunc 默认缓存键生成函数
func defaultKeyFunc(c *gin.Context) string {
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	method := c.Request.Method

	key := fmt.Sprintf("%s:%s", method, path)
	if query != "" {
		key += ":" + query
	}

	return key
}

// CacheResponse 缓存响应数据结构
type CacheResponse struct {
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
	Timestamp  time.Time         `json:"timestamp"`
}

// Middleware 返回缓存中间件
func (cm *CacheMiddleware) Middleware(ttl ...time.Duration) gin.HandlerFunc {
	cacheTTL := cm.defaultTTL
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}

	return func(c *gin.Context) {
		// 只缓存GET请求
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		cacheKey := cm.keyFunc(c)
		ctx := c.Request.Context()

		// 尝试从缓存获取
		cachedData, err := cm.cache.Get(ctx, cacheKey)
		if err == nil {
			var response CacheResponse
			if err := json.Unmarshal([]byte(cachedData), &response); err == nil {
				// 设置响应头
				for key, value := range response.Headers {
					c.Header(key, value)
				}
				c.Header("X-Cache", "HIT")
				c.Header("X-Cache-Time", response.Timestamp.Format(time.RFC3339))

				c.JSON(response.StatusCode, response.Body)
				c.Abort()
				return
			}
		}

		// 缓存未命中，继续处理请求
		c.Header("X-Cache", "MISS")

		// 使用自定义ResponseWriter来捕获响应
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
			headers:        make(map[string]string),
		}
		c.Writer = writer

		c.Next()

		// 只缓存成功的响应
		if writer.statusCode >= 200 && writer.statusCode < 300 {
			// 解析响应体
			var bodyData interface{}
			if len(writer.body) > 0 {
				json.Unmarshal(writer.body, &bodyData)
			}

			response := CacheResponse{
				StatusCode: writer.statusCode,
				Headers:    writer.headers,
				Body:       bodyData,
				Timestamp:  time.Now(),
			}

			// 序列化并缓存
			if data, err := json.Marshal(response); err == nil {
				cm.cache.Set(ctx, cacheKey, string(data), cacheTTL)
			}
		}
	}
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body       []byte
	headers    map[string]string
	statusCode int
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body = append(rw.body, data...)
	return rw.ResponseWriter.Write(data)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	// 复制响应头
	for key, values := range rw.ResponseWriter.Header() {
		if len(values) > 0 {
			rw.headers[key] = values[0]
		}
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

// InvalidatePattern 根据模式删除缓存
func (cm *CacheMiddleware) InvalidatePattern(ctx context.Context, pattern string) error {
	if rc, ok := cm.cache.(*RedisCache); ok {
		keys, err := rc.client.Keys(ctx, rc.buildKey(pattern)).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			// 移除前缀
			cleanKeys := make([]string, len(keys))
			for i, key := range keys {
				if rc.prefix != "" {
					cleanKeys[i] = strings.TrimPrefix(key, rc.prefix+":")
				} else {
					cleanKeys[i] = key
				}
			}
			return cm.cache.Del(ctx, cleanKeys...)
		}
	}
	return nil
}

// Config 缓存配置
type Config struct {
	Enabled    bool          `yaml:"enabled"`
	RedisAddr  string        `yaml:"redis_addr"`
	RedisDB    int           `yaml:"redis_db"`
	Prefix     string        `yaml:"prefix"`
	DefaultTTL time.Duration `yaml:"default_ttl"`
}

// NewCache 根据配置创建缓存
func NewCache(config Config) (Cache, error) {
	if !config.Enabled {
		return &NoOpCache{}, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		DB:   config.RedisDB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return NewRedisCache(rdb, config.Prefix), nil
}

// NoOpCache 空操作缓存（用于禁用缓存）
type NoOpCache struct{}

func (noc *NoOpCache) Get(ctx context.Context, key string) (string, error) {
	return "", redis.Nil
}

func (noc *NoOpCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}

func (noc *NoOpCache) Del(ctx context.Context, keys ...string) error {
	return nil
}

func (noc *NoOpCache) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}
