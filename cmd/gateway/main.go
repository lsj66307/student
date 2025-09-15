package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"student-management-system/internal/gateway"
)

func main() {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 加载配置
	cfg, err := loadConfig("config/gateway.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建网关实例
	gw := gateway.NewGateway(cfg)

	// 设置路由
	router := gw.SetupRouter()

	// 启动服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  60 * time.Second,
	}

	// 优雅关闭
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gateway server...")
		gw.Shutdown()
		server.Close()
	}()

	log.Printf("Gateway server starting on port %d", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig 加载配置文件
func loadConfig(filename string) (*gateway.GatewayConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// 如果配置文件不存在，使用默认配置
		log.Printf("Config file not found, using default config")
		return &gateway.GatewayConfig{
			Port:      8090,
			Timeout:   15 * time.Second,
			RateLimit: 100,
			Services: map[string]gateway.ServiceInfo{
				"student-api": {
					URL:    "http://localhost:8080",
					Prefix: "/api/v1",
					Weight: 1,
				},
			},
			LoadBalancer: gateway.LoadBalancerConfig{
				Strategy:      "round_robin",
				HealthCheck:   true,
				CheckInterval: 30 * time.Second,
				CheckTimeout:  5 * time.Second,
			},
		}, nil
	}

	var cfg gateway.GatewayConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
