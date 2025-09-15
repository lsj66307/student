package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Gateway 网关接口
type GatewayInterface interface {
	// 启动网关
	Start() error
	// 停止网关
	Shutdown() error
	// 设置路由
	SetupRouter() *gin.Engine
	// 获取统计信息
	GetStats() map[string]interface{}
}

// LoadBalancer 负载均衡器接口
type LoadBalancerInterface interface {
	// 注册服务实例
	RegisterService(serviceName string, instance *ServiceInstance) error
	// 注销服务实例
	UnregisterService(serviceName string, instanceID string) error
	// 选择服务实例
	SelectInstance(serviceName string) (*ServiceInstance, error)
	// 获取所有服务实例
	GetInstances(serviceName string) []*ServiceInstance
	// 更新实例健康状态
	UpdateInstanceHealth(serviceName string, instanceID string, healthy bool)
}

// HealthChecker 健康检查器接口
type HealthCheckerInterface interface {
	// 开始健康检查
	Start(ctx context.Context)
	// 停止健康检查
	Stop()
	// 检查单个实例
	CheckInstance(instance *ServiceInstance) bool
}

// Middleware 中间件接口
type MiddlewareInterface interface {
	// 安全中间件
	SecurityMiddleware() gin.HandlerFunc
	// CORS中间件
	CORSMiddleware() gin.HandlerFunc
	// 日志中间件
	LoggingMiddleware() gin.HandlerFunc
	// 限流中间件
	RateLimitMiddleware() gin.HandlerFunc
	// 认证中间件
	AuthMiddleware() gin.HandlerFunc
}

// ServiceRegistry 服务注册表接口
type ServiceRegistryInterface interface {
	// 注册服务
	Register(serviceName string, instance *ServiceInstance) error
	// 注销服务
	Unregister(serviceName string, instanceID string) error
	// 发现服务
	Discover(serviceName string) ([]*ServiceInstance, error)
	// 监听服务变化
	Watch(serviceName string, callback func([]*ServiceInstance))
}

// Router 路由器接口
type RouterInterface interface {
	// 设置路由
	SetupRoutes(engine *gin.Engine)
	// 添加路由组
	AddRouteGroup(prefix string, handlers ...gin.HandlerFunc) *gin.RouterGroup
	// 代理请求
	ProxyRequest(c *gin.Context, targetURL string)
}

// ConfigManager 配置管理器接口
type ConfigManagerInterface interface {
	// 加载配置
	Load(configPath string) error
	// 获取配置
	Get(key string) interface{}
	// 设置配置
	Set(key string, value interface{})
	// 监听配置变化
	Watch(callback func(key string, value interface{}))
}



// RateLimiter 限流器接口
type RateLimiterInterface interface {
	// 检查是否允许请求
	Allow(key string) bool
	// 获取剩余配额
	Remaining(key string) int
	// 重置限流器
	Reset(key string)
}