package repository

import (
	"context"
	"fmt"
	"student-management-system/internal/config"
	"student-management-system/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局Redis客户端
var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		logger.WithError(err).Error("Failed to load config for Redis")
		return fmt.Errorf("加载配置失败: %v", err)
	}

	logger.WithFields(map[string]interface{}{
		"host": cfg.Redis.Host,
		"port": cfg.Redis.Port,
		"db":   cfg.Redis.DB,
	}).Info("Connecting to Redis")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.WithError(err).Error("Redis connection test failed")
		// 允许程序在没有Redis的情况下运行
		RedisClient = nil
		logger.Warn("Redis connection failed, running without cache")
		return nil
	}

	logger.Info("Successfully connected to Redis")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient != nil {
		logger.Info("Closing Redis connection")
		return RedisClient.Close()
	}
	return nil
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}
