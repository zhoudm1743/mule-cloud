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

	"mule-cloud/app/production/services"
	"mule-cloud/app/production/transport"

	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/production.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 ProductionService 启动中...",
		zap.String("service", cfg.Server.Name),
		zap.Int("port", cfg.Server.Port),
	)

	// 初始化MongoDB DatabaseManager（如果启用）
	if cfg.MongoDB.Enabled {
		client, err := dbPkg.InitMongoDB(&cfg.MongoDB)
		if err != nil {
			loggerPkg.Fatal("初始化MongoDB失败", zap.Error(err))
		}
		dbPkg.InitDatabaseManager(client)
		loggerPkg.Info("✅ DatabaseManager初始化成功（支持多租户数据库隔离）")
	}

	// 初始化Redis（如果启用）
	if cfg.Redis.Enabled {
		if _, err := cachePkg.InitRedis(&cfg.Redis); err != nil {
			loggerPkg.Fatal("初始化Redis失败", zap.Error(err))
		}
		defer cachePkg.CloseRedis()
	}

	// 初始化服务
	scanSvc := services.NewScanService()
	reportSvc := services.NewReportService()
	qualitySvc := services.NewQualityService()
	reworkSvc := services.NewReworkService()

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())
	r.Use(middleware.OperationLogMiddleware())

	// 初始化 JWT 管理器（用于直接访问时验证token）
	jwtManager := jwtPkg.NewJWTManager(nil, 0)

	// Production路由组（需要认证）
	production := r.Group("/production")
	middleware.Apply(production, jwtManager) // ✅ 一个函数搞定
	{
		// 扫码路由
		production.POST("/scan/parse", transport.ParseScanCodeHandler(scanSvc)) // 扫码解析

		// 工序上报路由
		reports := production.Group("/reports")
		{
			reports.POST("", transport.SubmitReportHandler(reportSvc))       // 提交上报
			reports.GET("", transport.GetReportListHandler(reportSvc))       // 上报列表
			reports.GET("/:id", transport.GetReportByIDHandler(reportSvc))   // 上报详情
			reports.DELETE("/:id", transport.DeleteReportHandler(reportSvc)) // 删除上报记录
		}

		// 进度查询路由
		production.GET("/progress/:order_id", transport.GetOrderProgressHandler(reportSvc)) // 订单进度

		// 工资统计路由
		production.GET("/salary", transport.GetSalaryHandler(reportSvc)) // 工资统计

		// 质检路由
		inspections := production.Group("/inspections")
		{
			inspections.POST("", transport.SubmitInspectionHandler(qualitySvc))       // 提交质检
			inspections.GET("", transport.GetInspectionListHandler(qualitySvc))       // 质检列表
			inspections.GET("/:id", transport.GetInspectionHandler(qualitySvc))       // 质检详情
			inspections.DELETE("/:id", transport.DeleteInspectionHandler(qualitySvc)) // 删除质检记录
		}

		// 返工路由
		reworks := production.Group("/reworks")
		{
			reworks.POST("", transport.CreateReworkHandler(reworkSvc))               // 创建返工单
			reworks.GET("", transport.GetReworkListHandler(reworkSvc))               // 返工列表
			reworks.GET("/:id", transport.GetReworkHandler(reworkSvc))               // 返工详情
			reworks.PUT("/:id/complete", transport.CompleteReworkHandler(reworkSvc)) // 完成返工
			reworks.DELETE("/:id", transport.DeleteReworkHandler(reworkSvc))         // 删除返工记录
		}
	}

	// 健康检查端点（不需要认证）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "production-service",
		})
	})

	// Consul服务注册（如果启用）
	if cfg.Consul.Enabled {
		serviceConfig := &cousul.ServiceConfig{
			ServiceName:    cfg.Consul.ServiceName,
			ServiceAddress: cfg.Consul.ServiceIP, // 明确指定服务地址
			ServicePort:    cfg.Consul.ServicePort,
			Tags:           cfg.Consul.Tags,
			HealthCheck: &cousul.HealthCheck{
				HTTP:                           fmt.Sprintf("http://%s:%d/health", cfg.Consul.ServiceIP, cfg.Consul.ServicePort),
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

		// 自动注册路由配置到网关
		routeConfig := &cousul.RouteConfig{
			Prefix:        "/production", // 微服务路径
			GatewayPrefix: "/api",        // 网关前缀（小程序通过 /api 访问）
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   true,       // 需要认证
			RequireRole:   []string{}, // 角色要求
		}

		err = cousul.RegisterAndRun(r, serviceConfig, cfg.Consul.Address, routeConfig)
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
