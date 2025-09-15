package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// ConfigManager 配置管理器
type ConfigManager struct {
	mu           sync.RWMutex
	config       *GatewayConfig
	configPath   string
	watchers     []ConfigWatcher
	lastModified time.Time
}

// ConfigWatcher 配置变更监听器
type ConfigWatcher interface {
	OnConfigChanged(config *GatewayConfig) error
}

// ServiceRegistryConfig 服务注册配置
type ServiceRegistryConfig struct {
	Type          string            `yaml:"type" json:"type"`                     // consul, etcd, file
	Endpoints     []string          `yaml:"endpoints" json:"endpoints"`           // 注册中心地址
	Prefix        string            `yaml:"prefix" json:"prefix"`                 // 服务前缀
	Timeout       time.Duration     `yaml:"timeout" json:"timeout"`               // 超时时间
	RetryInterval time.Duration     `yaml:"retry_interval" json:"retry_interval"` // 重试间隔
	HealthCheck   HealthCheckConfig `yaml:"health_check" json:"health_check"`     // 健康检查配置
	Metadata      map[string]string `yaml:"metadata" json:"metadata"`             // 元数据
}

// HealthCheckConfig 健康检查配置
type HealthCheckConfig struct {
	Enabled       bool          `yaml:"enabled" json:"enabled"`
	Interval      time.Duration `yaml:"interval" json:"interval"`
	Timeout       time.Duration `yaml:"timeout" json:"timeout"`
	Path          string        `yaml:"path" json:"path"`
	Method        string        `yaml:"method" json:"method"`
	ExpectedCode  int           `yaml:"expected_code" json:"expected_code"`
	MaxFailures   int           `yaml:"max_failures" json:"max_failures"`
	RetryInterval time.Duration `yaml:"retry_interval" json:"retry_interval"`
}

// EnhancedGatewayConfig 增强的网关配置
type EnhancedGatewayConfig struct {
	*GatewayConfig
	ServiceRegistry ServiceRegistryConfig `yaml:"service_registry" json:"service_registry"`
	Security        SecurityConfig        `yaml:"security" json:"security"`
	Monitoring      MonitoringConfig      `yaml:"monitoring" json:"monitoring"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	JWT       JWTConfig       `yaml:"jwt" json:"jwt"`
	CORS      CORSConfig      `yaml:"cors" json:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit" json:"rate_limit"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string        `yaml:"secret" json:"secret"`
	Expiration time.Duration `yaml:"expiration" json:"expiration"`
	Issuer     string        `yaml:"issuer" json:"issuer"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins" json:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods" json:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers" json:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers" json:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials" json:"allow_credentials"`
	MaxAge           int      `yaml:"max_age" json:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled   bool          `yaml:"enabled" json:"enabled"`
	RPS       int           `yaml:"rps" json:"rps"`
	Burst     int           `yaml:"burst" json:"burst"`
	Window    time.Duration `yaml:"window" json:"window"`
	KeyFunc   string        `yaml:"key_func" json:"key_func"` // ip, user, header
	Whitelist []string      `yaml:"whitelist" json:"whitelist"`
	Blacklist []string      `yaml:"blacklist" json:"blacklist"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled   bool   `yaml:"enabled" json:"enabled"`
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	Namespace string `yaml:"namespace" json:"namespace"`
	Subsystem string `yaml:"subsystem" json:"subsystem"`
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string) (*ConfigManager, error) {
	cm := &ConfigManager{
		configPath: configPath,
		watchers:   make([]ConfigWatcher, 0),
	}

	if err := cm.LoadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cm, nil
}

// LoadConfig 加载配置
func (cm *ConfigManager) LoadConfig() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 检查文件是否存在
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", cm.configPath)
	}

	// 获取文件修改时间
	fileInfo, err := os.Stat(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// 如果文件没有修改，直接返回
	if !cm.lastModified.IsZero() && fileInfo.ModTime().Equal(cm.lastModified) {
		return nil
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var config GatewayConfig
	ext := filepath.Ext(cm.configPath)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	// 验证配置
	if err := cm.validateConfig(&config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 设置默认值
	cm.setDefaults(&config)

	// 更新配置
	oldConfig := cm.config
	cm.config = &config
	cm.lastModified = fileInfo.ModTime()

	// 通知监听器
	if oldConfig != nil {
		cm.notifyWatchers(&config)
	}

	log.Printf("Config loaded successfully from %s", cm.configPath)
	return nil
}

// validateConfig 验证配置
func (cm *ConfigManager) validateConfig(config *GatewayConfig) error {
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d", config.Port)
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if config.RateLimit < 0 {
		return fmt.Errorf("rate limit cannot be negative")
	}

	// 验证服务配置
	for name, service := range config.Services {
		if service.URL == "" {
			return fmt.Errorf("service %s: URL cannot be empty", name)
		}
		if service.Weight < 0 {
			return fmt.Errorf("service %s: weight cannot be negative", name)
		}
	}

	return nil
}

// setDefaults 设置默认值
func (cm *ConfigManager) setDefaults(config *GatewayConfig) {
	if config.Port == 0 {
		config.Port = 8080
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if config.RateLimit == 0 {
		config.RateLimit = 1000
	}

	if config.LoadBalancer.Strategy == "" {
		config.LoadBalancer.Strategy = "round_robin"
	}

	if config.LoadBalancer.CheckInterval == 0 {
		config.LoadBalancer.CheckInterval = 30 * time.Second
	}

	if config.LoadBalancer.CheckTimeout == 0 {
		config.LoadBalancer.CheckTimeout = 5 * time.Second
	}

	// 设置服务默认权重
	for name, service := range config.Services {
		if service.Weight == 0 {
			service.Weight = 1
			config.Services[name] = service
		}
	}
}

// GetConfig 获取当前配置
func (cm *ConfigManager) GetConfig() *GatewayConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// UpdateConfig 更新配置
func (cm *ConfigManager) UpdateConfig(config *GatewayConfig) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 验证配置
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	// 设置默认值
	cm.setDefaults(config)

	// 更新配置
	cm.config = config

	// 保存到文件
	if err := cm.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 通知监听器
	cm.notifyWatchers(config)

	return nil
}

// saveConfig 保存配置到文件
func (cm *ConfigManager) saveConfig() error {
	ext := filepath.Ext(cm.configPath)
	var data []byte
	var err error

	switch ext {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(cm.config)
	case ".json":
		data, err = json.MarshalIndent(cm.config, "", "  ")
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return ioutil.WriteFile(cm.configPath, data, 0644)
}

// AddWatcher 添加配置监听器
func (cm *ConfigManager) AddWatcher(watcher ConfigWatcher) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

// RemoveWatcher 移除配置监听器
func (cm *ConfigManager) RemoveWatcher(watcher ConfigWatcher) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, w := range cm.watchers {
		if w == watcher {
			cm.watchers = append(cm.watchers[:i], cm.watchers[i+1:]...)
			break
		}
	}
}

// notifyWatchers 通知所有监听器
func (cm *ConfigManager) notifyWatchers(config *GatewayConfig) {
	for _, watcher := range cm.watchers {
		go func(w ConfigWatcher) {
			if err := w.OnConfigChanged(config); err != nil {
				log.Printf("Config watcher error: %v", err)
			}
		}(watcher)
	}
}

// WatchConfig 监控配置文件变化
func (cm *ConfigManager) WatchConfig() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := cm.LoadConfig(); err != nil {
			log.Printf("Failed to reload config: %v", err)
		}
	}
}

// GetServiceConfig 获取服务配置
func (cm *ConfigManager) GetServiceConfig(serviceName string) (ServiceInfo, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	service, exists := cm.config.Services[serviceName]
	return service, exists
}

// AddService 添加服务配置
func (cm *ConfigManager) AddService(serviceName string, service ServiceInfo) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.config.Services == nil {
		cm.config.Services = make(map[string]ServiceInfo)
	}

	cm.config.Services[serviceName] = service

	// 保存配置
	if err := cm.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 通知监听器
	cm.notifyWatchers(cm.config)

	return nil
}

// RemoveService 移除服务配置
func (cm *ConfigManager) RemoveService(serviceName string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.config.Services == nil {
		return fmt.Errorf("no services configured")
	}

	delete(cm.config.Services, serviceName)

	// 保存配置
	if err := cm.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// 通知监听器
	cm.notifyWatchers(cm.config)

	return nil
}
