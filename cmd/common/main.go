package main

import (
	"flag"
	"fmt"
	"log"

	cfgPkg "mule-cloud/core/config"
	"mule-cloud/core/cousul"
	dbPkg "mule-cloud/core/database"
	jwtPkg "mule-cloud/core/jwt"
	loggerPkg "mule-cloud/core/logger"
	"mule-cloud/core/middleware"
	"mule-cloud/core/response"
	"mule-cloud/core/storage"

	"mule-cloud/app/common/services"
	"mule-cloud/app/common/transport"
	"mule-cloud/internal/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/common.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 CommonService 启动中...",
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

	// 初始化存储
	storageInstance, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		loggerPkg.Fatal("初始化存储失败", zap.Error(err))
	}
	loggerPkg.Info("✅ 存储初始化成功", zap.String("type", cfg.Storage.Type))

	// 获取MongoDB数据库实例
	db := dbPkg.GetDatabaseManager().GetDatabase(cfg.MongoDB.Database)

	// 初始化Repository
	fileRepo := repository.NewFileRepository(db)

	// 初始化Service
	fileService := services.NewFileService(fileRepo, storageInstance)

	// 初始化Transport
	fileTransport := transport.NewFileTransport(fileService, loggerPkg.Logger)

	// 初始化 JWT 管理器（用于直接访问时验证token）
	jwtManager := jwtPkg.NewJWTManager(nil, 0)

	// 创建Gin引擎
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger()) // 添加日志中间件，记录所有HTTP请求
	router.Use(gin.Recovery())
	router.Use(response.UnifiedResponseMiddleware())

	// 注册路由
	registerRoutes(router, fileTransport, jwtManager)

	// 注册到Consul（如果启用）
	if cfg.Consul.Enabled {
		serviceConfig := &cousul.ServiceConfig{
			ServiceName:    cfg.Consul.ServiceName,
			ServiceAddress: cfg.Consul.ServiceIP,
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
			Prefix:        "/common", // 微服务路径
			GatewayPrefix: "/admin",  // 网关前缀（后台接口必须通过 /admin 访问）
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   true,       // 文件上传需要认证（获取租户信息）
			RequireRole:   []string{}, // 无角色限制
		}

		err = cousul.RegisterAndRun(router, serviceConfig, cfg.Consul.Address, routeConfig)
		if err != nil {
			loggerPkg.Fatal("服务启动失败", zap.Error(err))
		}
	} else {
		// 不使用Consul，直接启动服务
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		loggerPkg.Info("🎉 CommonService 启动成功（无Consul）", zap.String("地址", addr))
		if err := router.Run(addr); err != nil {
			loggerPkg.Fatal("启动HTTP服务器失败", zap.Error(err))
		}
	}
}

// registerRoutes 注册路由
func registerRoutes(router *gin.Engine, fileTransport *transport.FileTransport, jwtManager *jwtPkg.JWTManager) {
	// 健康检查（无需认证）
	router.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	// 静态文件服务（用于本地存储访问，无需认证）
	router.Static("/files", "./uploads")

	// API路由组（Gateway会去掉/admin前缀，所以这里只需要/common）
	api := router.Group("/common")
	middleware.Apply(api, jwtManager) // 应用认证中间件
	{
		// 文件管理
		files := api.Group("/files")
		{
			files.POST("/upload", fileTransport.UploadHandler())                // 上传文件
			files.GET("", fileTransport.ListHandler())                          // 文件列表
			files.GET("/:id", fileTransport.DownloadHandler())                  // 下载文件
			files.DELETE("/:id", fileTransport.DeleteHandler())                 // 删除文件
			files.GET("/:id/presigned", fileTransport.GetPresignedURLHandler()) // 获取预签名URL
		}
	}
}
