package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 结构化日志器
type Logger struct {
	*logrus.Logger
}

// Fields 日志字段类型
type Fields map[string]interface{}

// 全局日志实例
var defaultLogger *Logger

// LogLevel 日志级别
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// Config 日志配置
type Config struct {
	Level      LogLevel `mapstructure:"level"`
	Format     string   `mapstructure:"format"`      // json, text
	Output     string   `mapstructure:"output"`      // stdout, file
	FilePath   string   `mapstructure:"file_path"`   // 日志文件路径
	MaxSize    int      `mapstructure:"max_size"`    // 文件最大大小(MB)
	MaxBackups int      `mapstructure:"max_backups"` // 保留备份数
	MaxAge     int      `mapstructure:"max_age"`     // 保留天数
	Compress   bool     `mapstructure:"compress"`    // 是否压缩
}

// Init 初始化日志器
func Init(config Config) error {
	logger := logrus.New()

	// 设置日志级别
	switch config.Level {
	case DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	case FatalLevel:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	if config.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// 设置输出
	if config.Output == "file" && config.FilePath != "" {
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	defaultLogger = &Logger{Logger: logger}
	return nil
}

// GetLogger 获取日志实例
func GetLogger() *Logger {
	if defaultLogger == nil {
		// 如果没有初始化，使用默认配置
		Init(Config{
			Level:  InfoLevel,
			Format: "text",
			Output: "stdout",
		})
	}
	return defaultLogger
}

// WithFields 添加字段
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
	return l.Logger.WithFields(logrus.Fields(fields))
}

// WithContext 添加上下文
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.Logger.WithContext(ctx)
}

// WithError 添加错误信息
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// Debug 调试日志
func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(args...)
}

// Debugf 格式化调试日志
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Info 信息日志
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

// Infof 格式化信息日志
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

// Warnf 格式化警告日志
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Error 错误日志
func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}

// Errorf 格式化错误日志
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// 全局便捷方法
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

func WithFields(fields Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return GetLogger().WithError(err)
}

func WithContext(ctx context.Context) *logrus.Entry {
	return GetLogger().WithContext(ctx)
}
