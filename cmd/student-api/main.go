package main

import (
	"log"
	"student-management-system/internal/config"
	handler "student-management-system/internal/handler"
	repo "student-management-system/internal/repository"
	"student-management-system/pkg/cache"
	"student-management-system/pkg/errors"
	"student-management-system/pkg/logger"
	"student-management-system/pkg/ratelimit"
)

func main() {
	// 加载配置
	log.Println("正在加载配置...")
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志系统
	loggerConfig := logger.Config{
		Level:  logger.LogLevel(cfg.Log.Level),
		Format: cfg.Log.Format,
		Output: cfg.Log.Output,
	}
	if cfg.Log.Output == "file" {
		loggerConfig.FilePath = cfg.Log.FilePath
	}
	if err := logger.Init(loggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Starting student management system")

	// 初始化数据库连接
	logger.Info("正在初始化数据库连接...")
	err = repo.InitDB()
	if err != nil {
		logger.WithError(err).Fatal("数据库连接失败")
	}
	defer repo.CloseDB()
	logger.Info("Database connection established")

	// 初始化Redis连接
	logger.Info("正在初始化Redis连接...")
	err = repo.InitRedis(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Redis连接失败")
	}
	defer repo.CloseRedis()

	// 创建数据库表
	logger.Info("正在创建数据库表...")
	err = repo.CreateTables()
	if err != nil {
		logger.WithError(err).Fatal("创建数据库表失败")
	}

	// 初始化限流器
	logger.Info("正在初始化限流器...")
	rateLimitConfig := ratelimit.Config{
		Enabled:   cfg.RateLimit.Enabled,
		Type:      cfg.RateLimit.Type,
		Rate:      cfg.RateLimit.Rate,
		Burst:     cfg.RateLimit.Burst,
		Window:    cfg.RateLimit.Window,
		RedisAddr: cfg.RateLimit.RedisAddr,
		RedisDB:   cfg.RateLimit.RedisDB,
	}
	rateLimiter, err := ratelimit.NewRateLimiter(rateLimitConfig)
	if err != nil {
		logger.WithError(err).Fatal("创建限流器失败")
	}

	// 初始化缓存
	logger.Info("正在初始化缓存...")
	cacheConfig := cache.Config{
		Enabled:    cfg.Cache.Enabled,
		RedisAddr:  cfg.Cache.RedisAddr,
		RedisDB:    cfg.Cache.RedisDB,
		Prefix:     cfg.Cache.Prefix,
		DefaultTTL: cfg.Cache.DefaultTTL,
	}
	cacheInstance, err := cache.NewCache(cacheConfig)
	if err != nil {
		logger.WithError(err).Fatal("创建缓存失败")
	}
	cacheMiddleware := cache.NewCacheMiddleware(cacheInstance, cfg.Cache.DefaultTTL)

	// 设置路由
	logger.Info("正在设置路由...")
	router := handler.SetupRoutes(cfg)

	// 添加中间件
	router.Use(rateLimiter.Middleware())
	router.Use(errors.Recovery())
	router.Use(errors.LoggingMiddleware())
	router.Use(cacheMiddleware.Middleware())

	// 启动服务器
	logger.Info("学生管理系统启动中...")
	logger.Info("服务器运行在: http://localhost:3060")
	logger.Info("API文档: http://localhost:3060/")
	logger.Info("健康检查: http://localhost:3060/health")
	logger.Info("按 Ctrl+C 停止服务器")

	// 在端口3060启动服务器
	if err := router.Run(":3060"); err != nil {
		logger.WithError(err).Fatal("启动服务器失败")
	}
}
