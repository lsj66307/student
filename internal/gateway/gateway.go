package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"student-management-system/internal/handler"
)

// GatewayConfig 网关配置
type GatewayConfig struct {
	Port         int                    `yaml:"port" json:"port"`
	Timeout      time.Duration          `yaml:"timeout" json:"timeout"`
	RateLimit    int                    `yaml:"rate_limit" json:"rate_limit"`
	Services     map[string]ServiceInfo `yaml:"services" json:"services"`
	LoadBalancer LoadBalancerConfig     `yaml:"load_balancer" json:"load_balancer"`
}

// ServiceInfo 服务信息
type ServiceInfo struct {
	URL    string `yaml:"url" json:"url"`
	Prefix string `yaml:"prefix" json:"prefix"`
	Weight int    `yaml:"weight" json:"weight"`
}

// LoadBalancerConfig 负载均衡配置
type LoadBalancerConfig struct {
	Strategy      string        `yaml:"strategy" json:"strategy"`
	HealthCheck   bool          `yaml:"health_check" json:"health_check"`
	CheckInterval time.Duration `yaml:"check_interval" json:"check_interval"`
	CheckTimeout  time.Duration `yaml:"check_timeout" json:"check_timeout"`
}

// Gateway 网关核心结构
type Gateway struct {
	config          *GatewayConfig
	configManager   *ConfigManager
	serviceRegistry ServiceRegistry
	loadBalancer    *LoadBalancer
	router          *gin.Engine
	services        sync.Map // map[string]*ServiceInstance
	ctx             context.Context
	cancel          context.CancelFunc
	startTime       time.Time
}

// NewGateway 创建网关实例
func NewGateway(config *GatewayConfig) *Gateway {
	ctx, cancel := context.WithCancel(context.Background())

	// 创建负载均衡器
	strategy := LoadBalancerStrategy(config.LoadBalancer.Strategy)
	if strategy == "" {
		strategy = RoundRobin
	}

	lb := NewLoadBalancer(strategy)



	// 创建配置管理器
	configManager, err := NewConfigManager("./config.yaml")
	if err != nil {
		log.Printf("Failed to create config manager: %v", err)
		configManager = nil
	}

	// 创建服务注册中心
	serviceRegistry := NewFileServiceRegistry("./services.json")

	gateway := &Gateway{
		config:          config,
		configManager:   configManager,
		serviceRegistry: serviceRegistry,
		loadBalancer:    lb,
		ctx:             ctx,
		cancel:          cancel,
		startTime:       time.Now(),
	}

	// 初始化服务
	gateway.initServices()

	// 启动健康检查
	if config.LoadBalancer.HealthCheck {
		gateway.startHealthCheck()
	}

	return gateway
}

// initServices 初始化服务实例
func (g *Gateway) initServices() {
	for serviceName, serviceInfo := range g.config.Services {
		instance := &ServiceInstance{
			ID:           fmt.Sprintf("%s-1", serviceName),
			URL:          serviceInfo.URL,
			Weight:       serviceInfo.Weight,
			Healthy:      true,
			Connections:  0,
			LastCheck:    time.Now(),
			ResponseTime: 0,
		}

		// 注册到负载均衡器
		g.loadBalancer.RegisterService(serviceName, instance)

		// 存储服务实例
		g.services.Store(serviceName, instance)
	}
}

// getPortFromURL 从URL中获取端口
func getPortFromURL(u *url.URL) int {
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			return 443
		}
		return 80
	}

	// 简单的端口转换
	if port == "8080" {
		return 8080
	}
	return 80
}

// SetupRouter 设置路由
func (g *Gateway) SetupRouter() *gin.Engine {
	if gin.Mode() == gin.TestMode {
		g.router = gin.New()
	} else {
		g.router = gin.Default()
	}

	// 添加全局中间件
	g.setupMiddlewares()

	// 设置路由
	g.setupRoutes()

	return g.router
}

// setupMiddlewares 设置中间件
func (g *Gateway) setupMiddlewares() {
	// 安全中间件
	g.router.Use(g.SecurityMiddleware())

	// CORS中间件
	g.router.Use(g.CORSMiddleware())

	// 日志中间件
	g.router.Use(g.LoggingMiddleware())

	// 限流中间件
	g.router.Use(g.RateLimitMiddleware())
}



// setupRoutes 设置路由
func (g *Gateway) setupRoutes() {
	// 健康检查端点
	g.router.GET("/health", g.healthCheck)

	// 网关管理端点
	management := g.router.Group("/gateway")
	{
		management.GET("/stats", g.getStats)
		management.GET("/services", g.getServices)
		management.GET("/routes", g.getRoutes)
		management.POST("/services/:name/register", g.registerService)
		management.DELETE("/services/:name/unregister", g.unregisterService)
	}

	// API路由组
	api := g.router.Group("/api")
	{
		// 公开API（不需要认证）
		public := api.Group("/v1")
		{
			public.Any("/auth/*path", g.proxyHandler("student-api"))
		}

		// 需要认证的API
		protected := api.Group("/v1")
		protected.Use(g.AuthMiddleware())
		{
			protected.Any("/students/*path", g.proxyHandler("student-api"))
			protected.Any("/teachers/*path", g.proxyHandler("student-api"))
			protected.Any("/grades/*path", g.proxyHandler("student-api"))
		}
	}
}

// proxyHandler 通用代理处理器
func (g *Gateway) proxyHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		g.proxyRequest(c, serviceName)
	}
}

// proxyRequest 代理请求到后端服务
func (g *Gateway) proxyRequest(c *gin.Context, serviceName string) {
	// 获取服务实例
	instance, err := g.loadBalancer.GetInstance(serviceName)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, handler.Response{
			Code:    http.StatusServiceUnavailable,
			Message: "Service unavailable: " + err.Error(),
		})
		return
	}

	// 解析目标URL
	targetURL, err := url.Parse(instance.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler.Response{
			Code:    http.StatusInternalServerError,
			Message: "Invalid service URL",
		})
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ModifyResponse = g.modifyResponse
	proxy.ErrorHandler = g.proxyErrorHandler

	// 设置请求头
	g.setProxyHeaders(c, serviceName)

	// 记录请求
	g.logRequest(c, serviceName, targetURL.String())

	// 执行代理
	proxy.ServeHTTP(c.Writer, c.Request)

	// 释放连接
	defer g.loadBalancer.ReleaseConnection(instance)
}

// setProxyHeaders 设置代理请求头
func (g *Gateway) setProxyHeaders(c *gin.Context, serviceName string) {
	c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-Proto", "http")
	c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
	c.Request.Header.Set("X-Gateway-Service", serviceName)
	c.Request.Header.Set("X-Request-ID", generateRequestID())
	c.Request.Header.Set("X-Gateway-Version", "2.0.0")
}

// logRequest 记录请求日志
func (g *Gateway) logRequest(c *gin.Context, serviceName, targetURL string) {
	log.Printf("[PROXY] %s %s -> %s (%s)",
		c.Request.Method,
		c.Request.URL.Path,
		targetURL,
		serviceName)
}

// handleServiceError 处理服务错误
func (g *Gateway) handleServiceError(c *gin.Context, serviceName string, err error) {
	log.Printf("Service %s error: %v", serviceName, err)
	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error":   "Service Unavailable",
		"message": fmt.Sprintf("Service %s is currently unavailable", serviceName),
		"code":    "SERVICE_UNAVAILABLE",
	})
}

// proxyErrorHandler 代理错误处理器
func (g *Gateway) proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Proxy error: %v", err)
	w.WriteHeader(http.StatusBadGateway)
	w.Write([]byte(`{"error":"Bad Gateway","message":"Failed to proxy request"}`))
}

// modifyResponse 修改响应
func (g *Gateway) modifyResponse(resp *http.Response) error {
	resp.Header.Set("X-Gateway", "student-gateway")
	resp.Header.Set("X-Gateway-Version", "2.0.0")
	resp.Header.Set("X-Response-Time", time.Now().Format(time.RFC3339))
	return nil
}

// healthCheck 健康检查
func (g *Gateway) healthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(g.startTime).String(),
		"version":   "2.0.0",
		"services":  g.getServiceHealthStatus(),
	}

	c.JSON(http.StatusOK, status)
}

// getServiceHealthStatus 获取服务健康状态
func (g *Gateway) getServiceHealthStatus() map[string]interface{} {
	status := make(map[string]interface{})

	g.services.Range(func(key, value interface{}) bool {
		serviceName := key.(string)
		instance := value.(*ServiceInstance)

		status[serviceName] = gin.H{
			"healthy":    instance.Healthy,
			"last_check": instance.LastCheck.Unix(),
			"url":        instance.URL,
		}
		return true
	})

	return status
}

// getStats 获取网关统计信息
func (g *Gateway) getStats(c *gin.Context) {
	stats := gin.H{
		"gateway": gin.H{
			"version":    "2.0.0",
			"uptime":     time.Since(g.startTime).String(),
			"start_time": g.startTime.Unix(),
		},
		"config": gin.H{
			"port":       g.config.Port,
			"rate_limit": g.config.RateLimit,
			"timeout":    g.config.Timeout.String(),
		},
		"load_balancer": gin.H{
			"strategy":     g.config.LoadBalancer.Strategy,
			"health_check": g.config.LoadBalancer.HealthCheck,
		},
		"services": g.getServiceStats(),
	}

	c.JSON(http.StatusOK, stats)
}

// getServiceStats 获取服务统计信息
func (g *Gateway) getServiceStats() map[string]interface{} {
	stats := make(map[string]interface{})

	g.services.Range(func(key, value interface{}) bool {
		serviceName := key.(string)
		instance := value.(*ServiceInstance)

		stats[serviceName] = gin.H{
			"healthy":       instance.Healthy,
			"last_check":    instance.LastCheck.Unix(),
			"url":           instance.URL,
			"connections":   instance.Connections,
			"response_time": instance.ResponseTime.Milliseconds(),
		}
		return true
	})

	return stats
}

// GetServiceList 获取服务列表
func (g *Gateway) GetServiceList() []map[string]interface{} {
	var services []map[string]interface{}

	g.services.Range(func(key, value interface{}) bool {
		serviceName := key.(string)
		instance := value.(*ServiceInstance)

		services = append(services, map[string]interface{}{
			"name":    serviceName,
			"id":      instance.ID,
			"healthy": instance.Healthy,
			"url":     instance.URL,
		})
		return true
	})

	return services
}

// GetServiceStats 获取服务统计信息
func (g *Gateway) GetServiceStats() map[string]interface{} {
	stats := make(map[string]interface{})

	g.services.Range(func(key, value interface{}) bool {
		serviceName := key.(string)
		instance := value.(*ServiceInstance)

		healthyCount := 0
		if instance.Healthy {
			healthyCount = 1
		}

		stats[serviceName] = map[string]interface{}{
			"total_instances":   1,
			"healthy_instances": healthyCount,
			"connections":       instance.Connections,
			"response_time":     instance.ResponseTime.Milliseconds(),
		}
		return true
	})

	return stats
}

// getServices 获取服务列表
func (g *Gateway) getServices(c *gin.Context) {
	services := g.GetServiceList()
	c.JSON(http.StatusOK, handler.Response{
		Code:    http.StatusOK,
		Message: "Services retrieved successfully",
		Data: gin.H{
			"total":    len(services),
			"services": services,
		},
	})
}

// getRoutes 获取路由信息
func (g *Gateway) getRoutes(c *gin.Context) {
	routes := g.router.Routes()
	routeInfo := make([]gin.H, len(routes))

	for i, route := range routes {
		routeInfo[i] = gin.H{
			"method": route.Method,
			"path":   route.Path,
		}
	}

	c.JSON(http.StatusOK, handler.Response{
		Code:    http.StatusOK,
		Message: "Routes retrieved successfully",
		Data: gin.H{
			"total":  len(routes),
			"routes": routeInfo,
		},
	})
}

// registerService 动态注册服务
func (g *Gateway) registerService(c *gin.Context) {
	serviceName := c.Param("name")

	var req struct {
		URL    string `json:"url" binding:"required"`
		Weight int    `json:"weight"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, handler.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	_, err := url.Parse(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, handler.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid URL",
		})
		return
	}

	instance := &ServiceInstance{
		ID:              fmt.Sprintf("%s-%d", serviceName, time.Now().Unix()),
		ServiceName:     serviceName,
		URL:             req.URL,
		Weight:          req.Weight,
		Healthy:         true,
		Connections:     0,
		LastCheck:       time.Now(),
		ResponseTime:    0,
		RegisterTime:    time.Now(),
		FailureCount:    0,
		HealthCheckPath: "/health",
	}

	// 注册到负载均衡器
	g.loadBalancer.RegisterService(serviceName, instance)

	// 注册到服务注册中心
	if g.serviceRegistry != nil {
		if err := g.serviceRegistry.Register(g.ctx, instance); err != nil {
			log.Printf("Failed to register service to registry: %v", err)
		}
	}

	// 存储到本地缓存
	g.services.Store(serviceName, instance)

	c.JSON(http.StatusOK, handler.Response{
		Code:    http.StatusOK,
		Message: "Service registered successfully",
		Data: gin.H{
			"service":  serviceName,
			"instance": instance,
		},
	})
}

// RegisterService 注册新服务
func (g *Gateway) RegisterService(serviceName, serviceURL string, weight int) error {
	instance := &ServiceInstance{
		ID:              fmt.Sprintf("%s-%d", serviceName, time.Now().Unix()),
		ServiceName:     serviceName,
		URL:             serviceURL,
		Weight:          weight,
		Healthy:         true,
		Connections:     0,
		LastCheck:       time.Now(),
		ResponseTime:    0,
		RegisterTime:    time.Now(),
		FailureCount:    0,
		HealthCheckPath: "/health",
	}

	// 注册到负载均衡器
	g.loadBalancer.RegisterService(serviceName, instance)

	// 注册到服务注册中心
	if g.serviceRegistry != nil {
		if err := g.serviceRegistry.Register(g.ctx, instance); err != nil {
			log.Printf("Failed to register service to registry: %v", err)
		}
	}

	// 存储到本地缓存
	g.services.Store(serviceName, instance)

	return nil
}

// unregisterService 注销服务
func (g *Gateway) unregisterService(c *gin.Context) {
	serviceName := c.Param("name")

	g.services.Delete(serviceName)

	c.JSON(http.StatusOK, handler.Response{
		Code:    http.StatusOK,
		Message: "Service unregistered successfully",
		Data: gin.H{
			"service": serviceName,
		},
	})
}

// startHealthCheck 启动健康检查
func (g *Gateway) startHealthCheck() {
	interval := g.config.LoadBalancer.CheckInterval
	if interval == 0 {
		interval = 30 * time.Second
	}

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				g.performHealthChecks()
			case <-g.ctx.Done():
				return
			}
		}
	}()
}

// performHealthChecks 执行健康检查
func (g *Gateway) performHealthChecks() {
	g.services.Range(func(key, value interface{}) bool {
		serviceName := key.(string)
		instance := value.(*ServiceInstance)

		go g.checkServiceHealth(serviceName, instance)
		return true
	})
}

// checkServiceHealth 检查单个服务健康状态
func (g *Gateway) checkServiceHealth(serviceName string, instance *ServiceInstance) {
	timeout := g.config.LoadBalancer.CheckTimeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	client := &http.Client{Timeout: timeout}
	healthURL := fmt.Sprintf("%s/health", instance.URL)

	start := time.Now()
	resp, err := client.Get(healthURL)
	responseTime := time.Since(start)

	if err != nil {
		instance.Healthy = false
		log.Printf("Health check failed for %s: %v", serviceName, err)
		return
	}
	defer resp.Body.Close()

	instance.Healthy = resp.StatusCode == http.StatusOK
	instance.LastCheck = time.Now()
	instance.ResponseTime = responseTime

	log.Printf("Health check for %s: healthy=%v", serviceName, instance.Healthy)
}

// Shutdown 优雅关闭
func (g *Gateway) Shutdown() {
	log.Println("Shutting down gateway...")
	g.cancel()
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
