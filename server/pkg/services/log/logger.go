package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// 设置输出为标准输出
	Logger.SetOutput(os.Stdout)

	// 设置日志格式为带颜色的文本格式
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	Logger.SetLevel(logrus.InfoLevel)

	// 在开发环境下显示更多信息
	if os.Getenv("GIN_MODE") != "release" {
		Logger.SetLevel(logrus.DebugLevel)
	}
}

// GetLogger 返回全局logger实例
func GetLogger() *logrus.Logger {
	return Logger
}
