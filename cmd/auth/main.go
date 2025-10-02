package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	cachePkg "mule-cloud/core/cache"
	cfgPkg "mule-cloud/core/config"
	"mule-cloud/core/cousul"
	dbPkg "mule-cloud/core/database"
	jwtPkg "mule-cloud/core/jwt"
	loggerPkg "mule-cloud/core/logger"
	"mule-cloud/core/response"

	"mule-cloud/app/auth/services"
	"mule-cloud/app/auth/transport"
	"mule-cloud/app/gateway/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/auth.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 AuthService 启动中...",
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

	// 初始化JWT管理器
	jwtManager := jwtPkg.NewJWTManager(
		[]byte(cfg.JWT.SecretKey),
		time.Duration(cfg.JWT.ExpireTime)*time.Hour,
	)

	// 初始化认证服务
	authSvc := services.NewAuthService(jwtManager)

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())

	// 公开路由（不需要认证）
	public := r.Group("/auth")
	{
		public.POST("/login", transport.LoginHandler(authSvc))
		public.POST("/register", transport.RegisterHandler(authSvc))
		public.POST("/refresh", transport.RefreshTokenHandler(authSvc))
	}

	// 需要认证的路由
	protected := r.Group("/auth")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		protected.GET("/profile", transport.GetProfileHandler(authSvc))
		protected.PUT("/profile", transport.UpdateProfileHandler(authSvc))
		protected.POST("/password", transport.ChangePasswordHandler(authSvc))
		protected.GET("/getUserRoutes", transport.GetUserRoutesHandler(authSvc)) // 获取用户路由
	}

	// 健康检查（不需要认证）
	r.GET("/common/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Consul服务注册（如果启用）
	if cfg.Consul.Enabled {
		serviceConfig := &cousul.ServiceConfig{
			ServiceName:    cfg.Consul.ServiceName,
			ServiceAddress: cfg.Consul.ServiceIP, // 明确指定服务地址
			ServicePort:    cfg.Consul.ServicePort,
			Tags:           cfg.Consul.Tags,
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

		loggerPkg.Info("正在启动HTTP服务...",
			zap.String("address", fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)),
		)

		// 自动注册路由配置到网关
		routeConfig := &cousul.RouteConfig{
			Prefix:        "/auth",  // 微服务路径
			GatewayPrefix: "/admin", // 无网关前缀（公开访问）
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   false, // auth 服务本身不需要认证
			RequireRole:   []string{},
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
