package http

import (
	"net/http"
	"time"

	"mule-cloud/pkg/services/config"
	"mule-cloud/pkg/services/log"
	"mule-cloud/router"

	"github.com/gin-gonic/gin"
)

func Run() *http.Server {
	// 获取配置
	serverConfig := config.GetServer()

	// 根据环境设置gin模式和自定义输出
	if config.GetApp().IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
		// 设置gin的调试输出到我们的自定义writer
		gin.DefaultWriter = log.NewGinDebugWriter()
	}

	// 创建gin引擎，保留gin的路由注册日志但使用自定义请求日志
	r := gin.New()

	// 使用自定义的logrus中间件
	r.Use(log.GinLogger())
	r.Use(log.GinRecovery())

	// 注册路由
	router.Register(r)

	return &http.Server{
		Addr:           serverConfig.GetServerAddr(),
		Handler:        r,
		ReadTimeout:    time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(serverConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: serverConfig.MaxHeaderBytes,
	}
}
