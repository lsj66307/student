package gateway

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// LoadBalancerStrategy 负载均衡策略
type LoadBalancerStrategy string

const (
	RoundRobin LoadBalancerStrategy = "round_robin"
	Random     LoadBalancerStrategy = "random"
	Weighted   LoadBalancerStrategy = "weighted"
	LeastConn  LoadBalancerStrategy = "least_conn"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID              string        `json:"id"`
	ServiceName     string        `json:"service_name"`
	URL             string        `json:"url"`
	Weight          int           `json:"weight"`
	Healthy         bool          `json:"healthy"`
	Connections     int           `json:"connections"`
	LastCheck       time.Time     `json:"last_check"`
	ResponseTime    time.Duration `json:"response_time"`
	RegisterTime    time.Time     `json:"register_time"`
	FailureCount    int           `json:"failure_count"`
	HealthCheckPath string        `json:"health_check_path"`
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	strategy LoadBalancerStrategy
	services map[string][]*ServiceInstance
	counters map[string]int
	mu       sync.RWMutex
	rand     *rand.Rand
}

// NewLoadBalancer 创建负载均衡器
func NewLoadBalancer(strategy LoadBalancerStrategy) *LoadBalancer {
	return &LoadBalancer{
		strategy: strategy,
		services: make(map[string][]*ServiceInstance),
		counters: make(map[string]int),
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// RegisterService 注册服务实例
func (lb *LoadBalancer) RegisterService(serviceName string, instance *ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if _, exists := lb.services[serviceName]; !exists {
		lb.services[serviceName] = make([]*ServiceInstance, 0)
		lb.counters[serviceName] = 0
	}

	lb.services[serviceName] = append(lb.services[serviceName], instance)
}

// UnregisterService 注销服务实例
func (lb *LoadBalancer) UnregisterService(serviceName, instanceID string) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	instances, exists := lb.services[serviceName]
	if !exists {
		return
	}

	for i, instance := range instances {
		if instance.ID == instanceID {
			lb.services[serviceName] = append(instances[:i], instances[i+1:]...)
			break
		}
	}
}

// GetInstance 获取服务实例
func (lb *LoadBalancer) GetInstance(serviceName string) (*ServiceInstance, error) {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances, exists := lb.services[serviceName]
	if !exists || len(instances) == 0 {
		return nil, errors.New("no available instances")
	}

	// 过滤健康的实例
	healthyInstances := make([]*ServiceInstance, 0)
	for _, instance := range instances {
		if instance.Healthy {
			healthyInstances = append(healthyInstances, instance)
		}
	}

	if len(healthyInstances) == 0 {
		return nil, errors.New("no healthy instances")
	}

	var selectedInstance *ServiceInstance
	switch lb.strategy {
	case RoundRobin:
		selectedInstance = lb.roundRobin(serviceName, healthyInstances)
	case Random:
		selectedInstance = lb.random(healthyInstances)
	case Weighted:
		selectedInstance = lb.weighted(healthyInstances)
	case LeastConn:
		selectedInstance = lb.leastConnections(healthyInstances)
	default:
		selectedInstance = healthyInstances[0]
	}

	// 增加连接计数
	if selectedInstance != nil {
		selectedInstance.Connections++
	}

	return selectedInstance, nil
}

// GetInstances 获取所有服务实例
func (lb *LoadBalancer) GetInstances(serviceName string) []*ServiceInstance {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	instances, exists := lb.services[serviceName]
	if !exists {
		return nil
	}

	// 返回副本以避免并发修改
	result := make([]*ServiceInstance, len(instances))
	copy(result, instances)
	return result
}

// GetStrategy 获取负载均衡策略
func (lb *LoadBalancer) GetStrategy() LoadBalancerStrategy {
	return lb.strategy
}

// SetStrategy 设置负载均衡策略
func (lb *LoadBalancer) SetStrategy(strategy LoadBalancerStrategy) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.strategy = strategy
}

// roundRobin 轮询选择
func (lb *LoadBalancer) roundRobin(serviceName string, instances []*ServiceInstance) *ServiceInstance {
	count := lb.counters[serviceName]
	instance := instances[count%len(instances)]
	lb.counters[serviceName] = (count + 1) % len(instances)
	return instance
}

// random 随机选择
func (lb *LoadBalancer) random(instances []*ServiceInstance) *ServiceInstance {
	return instances[lb.rand.Intn(len(instances))]
}

// weighted 加权选择
func (lb *LoadBalancer) weighted(instances []*ServiceInstance) *ServiceInstance {
	totalWeight := 0
	for _, instance := range instances {
		totalWeight += instance.Weight
	}

	if totalWeight == 0 {
		return instances[lb.rand.Intn(len(instances))]
	}

	randomWeight := lb.rand.Intn(totalWeight)
	currentWeight := 0

	for _, instance := range instances {
		currentWeight += instance.Weight
		if randomWeight < currentWeight {
			return instance
		}
	}

	return instances[len(instances)-1]
}

// leastConnections 最少连接数选择
func (lb *LoadBalancer) leastConnections(instances []*ServiceInstance) *ServiceInstance {
	if len(instances) == 0 {
		return nil
	}

	minConnections := instances[0].Connections
	selected := instances[0]

	for _, instance := range instances[1:] {
		if instance.Connections < minConnections {
			minConnections = instance.Connections
			selected = instance
		}
	}

	return selected
}

// ReleaseConnection 释放连接
func (lb *LoadBalancer) ReleaseConnection(instance *ServiceInstance) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if instance.Connections > 0 {
		instance.Connections--
	}
}

// UpdateInstanceHealth 更新实例健康状态
func (lb *LoadBalancer) UpdateInstanceHealth(serviceName, instanceID string, healthy bool, responseTime time.Duration) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	instances, exists := lb.services[serviceName]
	if !exists {
		return
	}

	for _, instance := range instances {
		if instance.ID == instanceID {
			instance.Healthy = healthy
			instance.LastCheck = time.Now()
			instance.ResponseTime = responseTime
			break
		}
	}
}

// GetStats 获取负载均衡器统计信息
func (lb *LoadBalancer) GetStats() map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["strategy"] = lb.strategy
	stats["services"] = make(map[string]interface{})

	for serviceName, instances := range lb.services {
		healthyCount := 0
		totalConnections := 0
		for _, instance := range instances {
			if instance.Healthy {
				healthyCount++
			}
			totalConnections += instance.Connections
		}

		stats["services"].(map[string]interface{})[serviceName] = map[string]interface{}{
			"total_instances":   len(instances),
			"healthy_instances": healthyCount,
			"total_connections": totalConnections,
		}
	}

	return stats
}
