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

	"mule-cloud/app/perms/services"
	"mule-cloud/app/perms/transport"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/perms.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 PermsService 启动中...",
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
	tenantSvc := services.NewTenantService()
	adminSvc := services.NewAdminService()
	menuSvc := services.NewMenuService()
	roleSvc := services.NewRoleService()

	// 初始化 JWT 管理器（用于直接访问时验证token）
	jwtManager := jwtPkg.NewJWTManager(nil, 0)

	// 初始化路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(response.RecoveryMiddleware())
	r.Use(response.UnifiedResponseMiddleware())
	r.Use(middleware.OperationLogMiddleware())

	// Perms路由组
	perms := r.Group("/perms")
	middleware.Apply(perms, jwtManager) // ✅ 一个函数搞定
	{
		// 租户路由
		tenant := perms.Group("/tenants")
		{
			tenant.GET("/:id", transport.GetTenantHandler(tenantSvc))                                          // 获取单个租户
			tenant.GET("", transport.ListTenantsHandler(tenantSvc))                                            // 分页列表
			tenant.GET("/all", transport.GetAllTenantsHandler(tenantSvc))                                      // 获取所有（不分页）
			tenant.POST("", transport.CreateTenantHandler(tenantSvc))                                          // 创建租户
			tenant.PUT("/:id", transport.UpdateTenantHandler(tenantSvc))                                       // 更新租户
			tenant.DELETE("/:id", transport.DeleteTenantHandler(tenantSvc))                                    // 删除租户
			tenant.POST("/:id/menus", transport.AssignTenantMenusHandler(tenantSvc.(*services.TenantService))) // 分配菜单权限（超管）
			tenant.GET("/:id/menus", transport.GetTenantMenusHandler(tenantSvc.(*services.TenantService)))     // 获取租户菜单权限
		}

		// 管理员路由
		admin := perms.Group("/admins")
		{
			admin.GET("/:id", transport.GetAdminHandler(adminSvc))                                                  // 获取单个管理员
			admin.GET("", transport.ListAdminsHandler(adminSvc))                                                    // 分页列表
			admin.GET("/all", transport.GetAllAdminsHandler(adminSvc))                                              // 获取所有（不分页）
			admin.POST("", transport.CreateAdminHandler(adminSvc))                                                  // 创建管理员
			admin.PUT("/:id", transport.UpdateAdminHandler(adminSvc))                                               // 更新管理员
			admin.DELETE("/:id", transport.DeleteAdminHandler(adminSvc))                                            // 删除管理员
			admin.POST("/:id/roles", transport.AssignAdminRolesHandler(adminSvc.(*services.AdminService)))          // 分配角色
			admin.GET("/:id/roles", transport.GetAdminRolesHandler(adminSvc.(*services.AdminService)))              // 获取管理员角色
			admin.DELETE("/:id/roles/:roleId", transport.RemoveAdminRoleHandler(adminSvc.(*services.AdminService))) // 移除角色
		}

		// 菜单路由（Nova-admin前端路由数据）
		menu := perms.Group("/menus")
		{
			menu.GET("/all", transport.GetAllMenusHandler(menuSvc))                // 获取所有菜单（扁平结构）
			menu.GET("/:id", transport.GetMenuHandler(menuSvc))                    // 获取单个菜单
			menu.GET("", transport.ListMenusHandler(menuSvc))                      // 分页列表
			menu.POST("", transport.CreateMenuHandler(menuSvc))                    // 创建菜单
			menu.POST("/batch-delete", transport.BatchDeleteMenusHandler(menuSvc)) // 批量删除
			menu.PUT("/:id", transport.UpdateMenuHandler(menuSvc))                 // 更新菜单
			menu.DELETE("/:id", transport.DeleteMenuHandler(menuSvc))              // 删除菜单
		}

		// 角色路由
		role := perms.Group("/roles")
		{
			role.GET("/:id", transport.GetRoleHandler(roleSvc))                    // 获取单个角色
			role.GET("", transport.ListRolesHandler(roleSvc))                      // 分页列表
			role.GET("/tenant", transport.GetTenantRolesHandler(roleSvc))          // 获取租户下的所有角色
			role.POST("", transport.CreateRoleHandler(roleSvc))                    // 创建角色
			role.PUT("/:id", transport.UpdateRoleHandler(roleSvc))                 // 更新角色
			role.DELETE("/:id", transport.DeleteRoleHandler(roleSvc))              // 删除角色
			role.POST("/batch-delete", transport.BatchDeleteRolesHandler(roleSvc)) // 批量删除
			role.POST("/:id/menus", transport.AssignMenusHandler(roleSvc))         // 分配菜单权限
			role.GET("/:id/menus", transport.GetRoleMenusHandler(roleSvc))         // 获取角色的菜单权限
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
			Prefix:        "/perms", // 微服务路径
			GatewayPrefix: "/admin", // 网关前缀（后台接口必须通过 /admin 访问）
			ServiceName:   cfg.Consul.ServiceName,
			RequireAuth:   true,       // 需要认证
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
