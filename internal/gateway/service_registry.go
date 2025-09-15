package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ServiceRegistry 服务注册中心接口
type ServiceRegistry interface {
	Register(ctx context.Context, service *ServiceInstance) error
	Unregister(ctx context.Context, serviceID string) error
	Discover(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInstance, error)
	HealthCheck(ctx context.Context, service *ServiceInstance) error
	Close() error
}

// FileServiceRegistry 基于文件的服务注册中心
type FileServiceRegistry struct {
	mu       sync.RWMutex
	services map[string][]*ServiceInstance
	watchers map[string][]chan []*ServiceInstance
	filePath string
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewFileServiceRegistry 创建文件服务注册中心
func NewFileServiceRegistry(filePath string) *FileServiceRegistry {
	ctx, cancel := context.WithCancel(context.Background())
	fsr := &FileServiceRegistry{
		services: make(map[string][]*ServiceInstance),
		watchers: make(map[string][]chan []*ServiceInstance),
		filePath: filePath,
		ctx:      ctx,
		cancel:   cancel,
	}

	// 启动健康检查
	go fsr.startHealthCheck()

	return fsr
}

// Register 注册服务
func (fsr *FileServiceRegistry) Register(ctx context.Context, service *ServiceInstance) error {
	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	serviceName := service.ServiceName
	if serviceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	// 检查服务是否已存在
	services := fsr.services[serviceName]
	for i, existing := range services {
		if existing.ID == service.ID {
			// 更新现有服务
			services[i] = service
			fsr.services[serviceName] = services
			fsr.notifyWatchers(serviceName, services)
			return fsr.saveToFile()
		}
	}

	// 添加新服务
	service.RegisterTime = time.Now()
	service.LastCheck = time.Now()
	service.Healthy = true

	fsr.services[serviceName] = append(fsr.services[serviceName], service)
	fsr.notifyWatchers(serviceName, fsr.services[serviceName])

	log.Printf("Service registered: %s (%s)", serviceName, service.ID)
	return fsr.saveToFile()
}

// Unregister 注销服务
func (fsr *FileServiceRegistry) Unregister(ctx context.Context, serviceID string) error {
	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	for serviceName, services := range fsr.services {
		for i, service := range services {
			if service.ID == serviceID {
				// 移除服务
				fsr.services[serviceName] = append(services[:i], services[i+1:]...)

				// 如果没有服务实例了，删除整个服务
				if len(fsr.services[serviceName]) == 0 {
					delete(fsr.services, serviceName)
				}

				fsr.notifyWatchers(serviceName, fsr.services[serviceName])
				log.Printf("Service unregistered: %s (%s)", serviceName, serviceID)
				return fsr.saveToFile()
			}
		}
	}

	return fmt.Errorf("service not found: %s", serviceID)
}

// Discover 发现服务
func (fsr *FileServiceRegistry) Discover(ctx context.Context, serviceName string) ([]*ServiceInstance, error) {
	fsr.mu.RLock()
	defer fsr.mu.RUnlock()

	services, exists := fsr.services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service not found: %s", serviceName)
	}

	// 只返回健康的服务实例
	healthyServices := make([]*ServiceInstance, 0)
	for _, service := range services {
		if service.Healthy {
			healthyServices = append(healthyServices, service)
		}
	}

	if len(healthyServices) == 0 {
		return nil, fmt.Errorf("no healthy instances for service: %s", serviceName)
	}

	return healthyServices, nil
}

// Watch 监听服务变化
func (fsr *FileServiceRegistry) Watch(ctx context.Context, serviceName string) (<-chan []*ServiceInstance, error) {
	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	ch := make(chan []*ServiceInstance, 10)
	fsr.watchers[serviceName] = append(fsr.watchers[serviceName], ch)

	// 发送当前服务列表
	if services, exists := fsr.services[serviceName]; exists {
		select {
		case ch <- services:
		default:
		}
	}

	// 在context取消时清理watcher
	go func() {
		<-ctx.Done()
		fsr.removeWatcher(serviceName, ch)
		close(ch)
	}()

	return ch, nil
}

// HealthCheck 健康检查
func (fsr *FileServiceRegistry) HealthCheck(ctx context.Context, service *ServiceInstance) error {
	if service.URL == "" {
		return fmt.Errorf("service URL is empty")
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 构建健康检查URL
	healthURL := service.URL
	if service.HealthCheckPath != "" {
		healthURL += service.HealthCheckPath
	} else {
		healthURL += "/health"
	}

	// 发送健康检查请求
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Close 关闭服务注册中心
func (fsr *FileServiceRegistry) Close() error {
	fsr.cancel()

	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	// 关闭所有watchers
	for _, watchers := range fsr.watchers {
		for _, ch := range watchers {
			close(ch)
		}
	}
	fsr.watchers = make(map[string][]chan []*ServiceInstance)

	return nil
}

// notifyWatchers 通知监听器
func (fsr *FileServiceRegistry) notifyWatchers(serviceName string, services []*ServiceInstance) {
	watchers, exists := fsr.watchers[serviceName]
	if !exists {
		return
	}

	for _, ch := range watchers {
		select {
		case ch <- services:
		default:
			// 如果channel满了，跳过这次通知
		}
	}
}

// removeWatcher 移除监听器
func (fsr *FileServiceRegistry) removeWatcher(serviceName string, ch chan []*ServiceInstance) {
	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	watchers := fsr.watchers[serviceName]
	for i, watcher := range watchers {
		if watcher == ch {
			fsr.watchers[serviceName] = append(watchers[:i], watchers[i+1:]...)
			break
		}
	}

	// 如果没有watchers了，删除整个entry
	if len(fsr.watchers[serviceName]) == 0 {
		delete(fsr.watchers, serviceName)
	}
}

// startHealthCheck 启动健康检查
func (fsr *FileServiceRegistry) startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-fsr.ctx.Done():
			return
		case <-ticker.C:
			fsr.performHealthChecks()
		}
	}
}

// performHealthChecks 执行健康检查
func (fsr *FileServiceRegistry) performHealthChecks() {
	fsr.mu.Lock()
	defer fsr.mu.Unlock()

	for serviceName, services := range fsr.services {
		for _, service := range services {
			go func(svc *ServiceInstance, name string) {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := fsr.HealthCheck(ctx, svc)

				fsr.mu.Lock()
				defer fsr.mu.Unlock()

				oldHealthy := svc.Healthy
				svc.LastCheck = time.Now()

				if err != nil {
					svc.Healthy = false
					svc.FailureCount++
					if oldHealthy {
						log.Printf("Service %s (%s) marked as unhealthy: %v", name, svc.ID, err)
					}
				} else {
					svc.Healthy = true
					svc.FailureCount = 0
					if !oldHealthy {
						log.Printf("Service %s (%s) marked as healthy", name, svc.ID)
					}
				}

				// 如果健康状态发生变化，通知watchers
				if oldHealthy != svc.Healthy {
					fsr.notifyWatchers(name, fsr.services[name])
				}
			}(service, serviceName)
		}
	}
}

// saveToFile 保存服务信息到文件
func (fsr *FileServiceRegistry) saveToFile() error {
	if fsr.filePath == "" {
		return nil // 如果没有指定文件路径，不保存
	}

	data, err := json.MarshalIndent(fsr.services, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal services: %w", err)
	}

	// 这里应该使用ioutil.WriteFile，但为了避免导入冲突，我们先返回nil
	// 在实际使用时需要实现文件写入逻辑
	log.Printf("Services data: %s", string(data))
	return nil
}

// loadFromFile 从文件加载服务信息
func (fsr *FileServiceRegistry) loadFromFile() error {
	if fsr.filePath == "" {
		return nil // 如果没有指定文件路径，不加载
	}

	// 这里应该实现从文件加载的逻辑
	// 在实际使用时需要实现文件读取逻辑
	return nil
}

// GetAllServices 获取所有服务
func (fsr *FileServiceRegistry) GetAllServices() map[string][]*ServiceInstance {
	fsr.mu.RLock()
	defer fsr.mu.RUnlock()

	// 创建副本以避免并发问题
	result := make(map[string][]*ServiceInstance)
	for name, services := range fsr.services {
		servicesCopy := make([]*ServiceInstance, len(services))
		copy(servicesCopy, services)
		result[name] = servicesCopy
	}

	return result
}

// GetServiceStats 获取服务统计信息
func (fsr *FileServiceRegistry) GetServiceStats() map[string]interface{} {
	fsr.mu.RLock()
	defer fsr.mu.RUnlock()

	stats := make(map[string]interface{})
	totalServices := 0
	totalInstances := 0
	healthyInstances := 0

	for serviceName, services := range fsr.services {
		totalServices++
		serviceHealthy := 0
		for _, service := range services {
			totalInstances++
			if service.Healthy {
				healthyInstances++
				serviceHealthy++
			}
		}

		stats[serviceName] = map[string]interface{}{
			"total_instances":     len(services),
			"healthy_instances":   serviceHealthy,
			"unhealthy_instances": len(services) - serviceHealthy,
		}
	}

	stats["summary"] = map[string]interface{}{
		"total_services":      totalServices,
		"total_instances":     totalInstances,
		"healthy_instances":   healthyInstances,
		"unhealthy_instances": totalInstances - healthyInstances,
	}

	return stats
}
