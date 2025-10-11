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

	"mule-cloud/app/miniapp/services"
	"mule-cloud/app/miniapp/transport"
	"mule-cloud/core/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/miniapp.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 MiniappService 启动中...",
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

	// 初始化JWT管理器
	jwtManager := jwtPkg.NewJWTManager(
		[]byte(cfg.JWT.SecretKey),
		time.Duration(cfg.JWT.ExpireTime)*time.Hour,
	)

	// 初始化微信服务
	wechatSvc := services.NewWechatService(
		jwtManager,
		cfg.Wechat.AppID,
		cfg.Wechat.AppSecret,
	)

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())
	r.Use(middleware.OperationLogMiddleware())

	// 公开路由（不需要认证）
	public := r.Group("/miniapp")
	{
		public.POST("/wechat/login", transport.WechatLoginHandler(wechatSvc))          // 微信登录
		public.POST("/wechat/bind-tenant", transport.BindTenantHandler(wechatSvc))     // 绑定租户
		public.POST("/wechat/select-tenant", transport.SelectTenantHandler(wechatSvc)) // 选择租户
	}

	// 需要认证的路由
	protected := r.Group("/miniapp")
	middleware.Apply(protected, jwtManager) // 应用JWT认证中间件
	{
		protected.POST("/wechat/switch-tenant", transport.SwitchTenantHandler(wechatSvc)) // 切换租户
		protected.GET("/user/info", transport.GetUserInfoHandler(wechatSvc))              // 获取用户信息
		protected.PUT("/user/info", transport.UpdateUserInfoHandler(wechatSvc))           // 更新用户信息
		protected.POST("/wechat/phone", transport.GetPhoneNumberHandler(wechatSvc))       // 绑定手机号
		protected.DELETE("/wechat/phone", transport.UnbindPhoneHandler(wechatSvc))        // 解绑手机号
	}

	// 健康检查（不需要认证）
	r.GET("/common/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Consul服务注册（如果启用）
	if cfg.Consul.Enabled {
		serviceConfig := &cousul.ServiceConfig{
			ServiceName:    cfg.Consul.ServiceName,
			ServiceAddress: cfg.Consul.ServiceIP,
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
			Prefix:        "/miniapp",
			GatewayPrefix: "/api",
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   false, // 小程序服务有部分公开接口
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
