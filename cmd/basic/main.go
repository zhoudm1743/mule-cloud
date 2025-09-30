package main

import (
	"flag"
	"fmt"
	"log"
	cachePkg "mule-cloud/core/cache"
	cfgPkg "mule-cloud/core/config"
	"mule-cloud/core/cousul"
	dbPkg "mule-cloud/core/database"
	loggerPkg "mule-cloud/core/logger"
	"mule-cloud/core/response"

	"mule-cloud/app/basic/services"
	"mule-cloud/app/basic/transport"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/basic.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := cfgPkg.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	// 初始化日志系统
	if err := loggerPkg.InitLogger(&cfg.Log); err != nil {
		log.Fatalf("初始化日志系统失败: %v", err)
	}
	defer loggerPkg.Close()

	loggerPkg.Info("🚀 BasicService 启动中...",
		zap.String("service", cfg.Server.Name),
		zap.Int("port", cfg.Server.Port),
	)

	// 初始化MongoDB（如果启用）
	if cfg.MongoDB.Enabled {
		if _, err := dbPkg.InitMongoDB(&cfg.MongoDB); err != nil {
			loggerPkg.Fatal("初始化MongoDB失败", zap.Error(err))
		}
		defer dbPkg.CloseMongoDB()
	}

	// 初始化Redis（如果启用）
	if cfg.Redis.Enabled {
		if _, err := cachePkg.InitRedis(&cfg.Redis); err != nil {
			loggerPkg.Fatal("初始化Redis失败", zap.Error(err))
		}
		defer cachePkg.CloseRedis()
	}

	// 初始化服务
	colorSvc := services.NewColorService()
	sizeSvc := services.NewSizeService()
	commonSvc := services.NewCommonService()

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())

	// Basic路由组
	basic := r.Group("/basic")
	{
		// 颜色路由
		color := basic.Group("/color")
		{
			color.GET("/:id", transport.GetColorHandler(colorSvc))
			color.GET("", transport.GetAllColorsHandler(colorSvc))
		}

		// 尺寸路由
		size := basic.Group("/size")
		{
			size.GET("/:id", transport.GetSizeHandler(sizeSvc))
			size.GET("", transport.GetAllSizesHandler(sizeSvc))
		}
	}

	// Common路由组
	common := r.Group("/common")
	{
		common.GET("/health", transport.HealthHandler(commonSvc))
	}

	// Consul服务注册（如果启用）
	if cfg.Consul.Enabled {
		serviceConfig := &cousul.ServiceConfig{
			ServiceName: cfg.Consul.ServiceName,
			ServicePort: cfg.Consul.ServicePort,
			Tags:        cfg.Consul.Tags,
			HealthCheck: &cousul.HealthCheck{
				HTTP:                           fmt.Sprintf("http://%s:%d/common/health", cfg.Consul.ServiceIP, cfg.Consul.ServicePort),
				Interval:                       cfg.Consul.HealthCheckInterval,
				Timeout:                        cfg.Consul.HealthCheckTimeout,
				DeregisterCriticalServiceAfter: cfg.Consul.DeregisterAfter,
			},
		}

		loggerPkg.Info("准备注册到Consul",
			zap.String("service", serviceConfig.ServiceName),
			zap.Int("port", serviceConfig.ServicePort),
			zap.String("consul", cfg.Consul.Address),
		)

		err = cousul.RegisterAndRun(r, serviceConfig, cfg.Consul.Address)
		if err != nil {
			loggerPkg.Fatal("服务启动失败", zap.Error(err))
		}
	} else {
		// 不使用Consul，直接启动服务
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		loggerPkg.Info("服务启动（无Consul）",
			zap.String("address", addr),
		)
		if err := r.Run(addr); err != nil {
			loggerPkg.Fatal("服务启动失败", zap.Error(err))
		}
	}
}
