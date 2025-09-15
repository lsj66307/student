package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	App       AppConfig       `mapstructure:"app"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Log       LogConfig       `mapstructure:"log"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Server    ServerConfig    `mapstructure:"server"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rateLimit"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret    string        `mapstructure:"secret"`
	ExpiresIn time.Duration `mapstructure:"expires_in"`
	Issuer    string        `mapstructure:"issuer"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	ExposeHeaders    []string `mapstructure:"expose_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	Enabled   bool          `mapstructure:"enabled"`
	Type      string        `mapstructure:"type"`
	Rate      int           `mapstructure:"rate"`
	Burst     int           `mapstructure:"burst"`
	Window    time.Duration `mapstructure:"window"`
	RedisAddr string        `mapstructure:"redis_addr"`
	RedisDB   int           `mapstructure:"redis_db"`
}

// Load 加载配置文件
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	// 自动读取环境变量
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}

	// 手动处理环境变量替换
	if err := expandEnvVars(&config); err != nil {
		return nil, fmt.Errorf("环境变量解析失败: %v", err)
	}

	// 处理时间单位转换
	config.Database.ConnMaxLifetime *= time.Second
	config.JWT.ExpiresIn *= time.Second
	config.Server.ReadTimeout *= time.Second
	config.Server.WriteTimeout *= time.Second
	config.Server.IdleTimeout *= time.Second
	config.Server.ShutdownTimeout *= time.Second

	return &config, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "student-management-system")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.mode", "development")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "student_management")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parse_time", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 3600)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.file_path", "logs/app.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 5)
	viper.SetDefault("log.max_age", 30)
	viper.SetDefault("log.compress", true)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expires_in", 86400)
	viper.SetDefault("jwt.issuer", "student-management-system")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Server defaults
	viper.SetDefault("server.read_timeout", 30)
	viper.SetDefault("server.write_timeout", 30)
	viper.SetDefault("server.idle_timeout", 60)
	viper.SetDefault("server.shutdown_timeout", 10)

	// CORS defaults
	viper.SetDefault("cors.allow_origins", []string{"http://localhost:3000", "http://localhost:8080"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.expose_headers", []string{"Content-Length"})
	viper.SetDefault("cors.allow_credentials", true)
	viper.SetDefault("cors.max_age", 86400)

	// RateLimit defaults
	viper.SetDefault("rateLimit.enabled", true)
	viper.SetDefault("rateLimit.type", "memory")
	viper.SetDefault("rateLimit.rate", 100)
	viper.SetDefault("rateLimit.burst", 200)
	viper.SetDefault("rateLimit.window", 60)
	viper.SetDefault("rateLimit.redis_addr", "localhost:6379")
	viper.SetDefault("rateLimit.redis_db", 1)
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.Charset,
		c.ParseTime,
		c.Loc,
	)
}

// GetAddr 获取服务器监听地址
func (c *AppConfig) GetAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}

// IsDevelopment 判断是否为开发模式
func (c *AppConfig) IsDevelopment() bool {
	return c.Mode == "development"
}

// IsProduction 判断是否为生产模式
func (c *AppConfig) IsProduction() bool {
	return c.Mode == "production"
}

// IsTest 判断是否为测试模式
func (c *AppConfig) IsTest() bool {
	return c.Mode == "test"
}

// GetRedisAddr 获取Redis连接地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// expandEnvVars 展开环境变量
func expandEnvVars(config *Config) error {
	// 数据库配置
	if host := os.Getenv("DATABASE_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DATABASE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		} else {
			return fmt.Errorf("无效的数据库端口: %s", port)
		}
	}
	if username := os.Getenv("DATABASE_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("DATABASE_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if dbname := os.Getenv("DATABASE_DBNAME"); dbname != "" {
		config.Database.DBName = dbname
	}

	// JWT配置
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}

	// Redis配置
	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.Redis.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Redis.Port = p
		} else {
			return fmt.Errorf("无效的Redis端口: %s", port)
		}
	}
	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.Redis.Password = password
	}
	if db := os.Getenv("REDIS_DB"); db != "" {
		if d, err := strconv.Atoi(db); err == nil {
			config.Redis.DB = d
		} else {
			return fmt.Errorf("无效的Redis数据库编号: %s", db)
		}
	}

	return nil
}
