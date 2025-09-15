package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"student-management-system/internal/config"
	handler "student-management-system/internal/handler"
	repo "student-management-system/internal/repository"
	"student-management-system/pkg/errors"
	"student-management-system/pkg/logger"
	"student-management-system/pkg/ratelimit"
	"syscall"
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
		log.Fatalf("初始化日志系统失败: %v", err)
	}

	logger.Info("正在启动学生管理系统")

	// 初始化数据库连接
	logger.Info("正在初始化数据库连接...")
	err = repo.InitDB()
	if err != nil {
		logger.WithError(err).Fatal("数据库连接失败")
	}
	defer repo.CloseDB()
	logger.Info("数据库连接已建立")

	// 初始化Redis连接
	logger.Info("正在初始化Redis连接...")
	err = repo.InitRedis()
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

	// 设置路由
	logger.Info("正在设置路由...")
	router := handler.SetupRoutes(cfg)

	// 添加中间件
	router.Use(rateLimiter.Middleware())
	router.Use(errors.Recovery())
	router.Use(errors.LoggingMiddleware())

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:         ":3060",
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 启动服务器
	logger.Info("学生管理系统启动中...")
	logger.Info("服务器运行在: http://localhost:3060")
	logger.Info("API文档: http://localhost:3060/")
	logger.Info("健康检查: http://localhost:3060/health")
	logger.Info("按 Ctrl+C 停止服务器")

	// 在goroutine中启动服务器
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("启动服务器失败")
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("正在关闭服务器...")

	// 优雅关闭服务器，等待现有连接完成
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("服务器强制关闭")
	}

	logger.Info("服务器已关闭")
}
