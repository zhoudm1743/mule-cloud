package logger

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// Config 日志配置
type Config struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Output     string `yaml:"output"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// logrusLogger logrus实现
type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

// NewLogger 创建新的日志记录器
func NewLogger(config Config) Logger {
	logger := logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 设置日志格式
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// 设置输出
	var output io.Writer
	switch config.Output {
	case "file":
		output = &lumberjack.Logger{
			Filename:   filepath.Join("logs", "app.log"),
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   true,
		}
	case "both":
		fileOutput := &lumberjack.Logger{
			Filename:   filepath.Join("logs", "app.log"),
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   true,
		}
		output = io.MultiWriter(os.Stdout, fileOutput)
	default:
		output = os.Stdout
	}
	logger.SetOutput(output)

	return &logrusLogger{
		logger: logger,
		entry:  logrus.NewEntry(logger),
	}
}

// Debug 调试日志
func (l *logrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugf 格式化调试日志
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info 信息日志
func (l *logrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infof 格式化信息日志
func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn 警告日志
func (l *logrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnf 格式化警告日志
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error 错误日志
func (l *logrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorf 格式化错误日志
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal 致命错误日志
func (l *logrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// WithField 添加字段
func (l *logrusLogger) WithField(key string, value interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithField(key, value),
	}
}

// WithFields 添加多个字段
func (l *logrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(logrus.Fields(fields)),
	}
}

// OperationLog 操作日志结构
type OperationLog struct {
	UserID     string      `json:"user_id"`
	Action     string      `json:"action"`
	Resource   string      `json:"resource"`
	ResourceID string      `json:"resource_id"`
	OldData    interface{} `json:"old_data,omitempty"`
	NewData    interface{} `json:"new_data,omitempty"`
	IP         string      `json:"ip"`
	UserAgent  string      `json:"user_agent"`
	Timestamp  time.Time   `json:"timestamp"`
	Duration   int64       `json:"duration_ms"`
	Error      string      `json:"error,omitempty"`
}

// LogOperation 记录操作日志
func LogOperation(logger Logger, opLog OperationLog) {
	logger.WithFields(map[string]interface{}{
		"type":        "operation",
		"user_id":     opLog.UserID,
		"action":      opLog.Action,
		"resource":    opLog.Resource,
		"resource_id": opLog.ResourceID,
		"ip":          opLog.IP,
		"user_agent":  opLog.UserAgent,
		"timestamp":   opLog.Timestamp,
		"duration_ms": opLog.Duration,
		"error":       opLog.Error,
	}).Info("Operation executed")
}
