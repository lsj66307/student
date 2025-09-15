package repository

import (
	"context"
	"fmt"
	"log"
	"student-management-system/internal/config"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis(cfg *config.Config) error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
	})

	// 测试连接
	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败，将在无缓存模式下运行: %v", err)
		RedisClient = nil
		return nil // 不返回错误，允许程序继续运行
	}

	log.Printf("Redis连接成功: %s", pong)
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis连接已关闭")
	}
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}
