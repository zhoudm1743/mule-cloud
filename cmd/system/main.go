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

	"mule-cloud/app/system/services"
	"mule-cloud/app/system/transport"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/system.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 SystemService 启动中...",
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
	tenantSvc := services.NewTenantService()
	adminSvc := services.NewAdminService()
	menuSvc := services.NewMenuService()

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())

	// System路由组
	system := r.Group("/system")
	{
		// 租户路由
		tenant := system.Group("/tenants")
		{
			tenant.GET("/:id", transport.GetTenantHandler(tenantSvc))       // 获取单个租户
			tenant.GET("", transport.ListTenantsHandler(tenantSvc))         // 分页列表
			tenant.GET("/all", transport.GetAllTenantsHandler(tenantSvc))   // 获取所有（不分页）
			tenant.POST("", transport.CreateTenantHandler(tenantSvc))       // 创建租户
			tenant.PUT("/:id", transport.UpdateTenantHandler(tenantSvc))    // 更新租户
			tenant.DELETE("/:id", transport.DeleteTenantHandler(tenantSvc)) // 删除租户
		}

		// 管理员路由
		admin := system.Group("/admins")
		{
			admin.GET("/:id", transport.GetAdminHandler(adminSvc))       // 获取单个管理员
			admin.GET("", transport.ListAdminsHandler(adminSvc))         // 分页列表
			admin.GET("/all", transport.GetAllAdminsHandler(adminSvc))   // 获取所有（不分页）
			admin.POST("", transport.CreateAdminHandler(adminSvc))       // 创建管理员
			admin.PUT("/:id", transport.UpdateAdminHandler(adminSvc))    // 更新管理员
			admin.DELETE("/:id", transport.DeleteAdminHandler(adminSvc)) // 删除管理员
		}

		// 菜单路由（Nova-admin前端路由数据）
		menu := system.Group("/menus")
		{
			menu.GET("/all", transport.GetAllMenusHandler(menuSvc))                // 获取所有菜单（扁平结构）
			menu.GET("/:id", transport.GetMenuHandler(menuSvc))                    // 获取单个菜单
			menu.GET("", transport.ListMenusHandler(menuSvc))                      // 分页列表
			menu.POST("", transport.CreateMenuHandler(menuSvc))                    // 创建菜单
			menu.POST("/batch-delete", transport.BatchDeleteMenusHandler(menuSvc)) // 批量删除
			menu.PUT("/:id", transport.UpdateMenuHandler(menuSvc))                 // 更新菜单
			menu.DELETE("/:id", transport.DeleteMenuHandler(menuSvc))              // 删除菜单
		}
	}

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "UP",
			"service": cfg.Server.Name,
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
			Prefix:        "/system", // 微服务路径
			GatewayPrefix: "/admin",  // 网关前缀（后台接口必须通过 /admin 访问）
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   false,      // 需要认证
			RequireRole:   []string{}, // 需要 admin 角色
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
