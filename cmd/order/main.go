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

	"mule-cloud/app/order/services"
	"mule-cloud/app/order/transport"
	"mule-cloud/internal/repository"

	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/order.yaml", "配置文件路径")
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

	loggerPkg.Info("🚀 OrderService 启动中...",
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
	orderSvc := services.NewOrderService()
	styleSvc := services.NewStyleService()
	cuttingSvc := services.NewCuttingService(
		repository.NewCuttingTaskRepository(),
		repository.NewCuttingBatchRepository(),
		repository.NewCuttingPieceRepository(),
		repository.NewOrderRepository(),
	)
	commonSvc := services.NewCommonService()

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

	// Order路由组（需要认证）
	order := r.Group("/order")
	middleware.Apply(order, jwtManager) // ✅ 一个函数搞定
	{
		// 订单路由
		orders := order.Group("/orders")
		{
			orders.GET("/:id", transport.GetOrderHandler(orderSvc))                       // 获取单个订单
			orders.GET("", transport.ListOrdersHandler(orderSvc))                         // 分页列表
			orders.POST("", transport.CreateOrderHandler(orderSvc))                       // 创建订单（步骤1）
			orders.PUT("/:id/style", transport.UpdateOrderStyleHandler(orderSvc))         // 更新款式数量（步骤2）
			orders.PUT("/:id/procedure", transport.UpdateOrderProcedureHandler(orderSvc)) // 更新工序（步骤3）
			orders.PUT("/:id", transport.UpdateOrderHandler(orderSvc))                    // 更新订单
			orders.POST("/:id/copy", transport.CopyOrderHandler(orderSvc))                // 复制订单
			orders.DELETE("/:id", transport.DeleteOrderHandler(orderSvc))                 // 删除订单
		}

		// 款式路由
		styles := order.Group("/styles")
		{
			styles.GET("/:id", transport.GetStyleHandler(styleSvc))       // 获取单个款式
			styles.GET("", transport.ListStylesHandler(styleSvc))         // 分页列表
			styles.GET("/all", transport.GetAllStylesHandler(styleSvc))   // 获取所有（不分页）
			styles.POST("", transport.CreateStyleHandler(styleSvc))       // 创建款式
			styles.PUT("/:id", transport.UpdateStyleHandler(styleSvc))    // 更新款式
			styles.DELETE("/:id", transport.DeleteStyleHandler(styleSvc)) // 删除款式
		}

		// 裁剪路由
		cutting := order.Group("/cutting")
		{
			// 裁剪任务路由
			tasks := cutting.Group("/tasks")
			{
				tasks.POST("", transport.CreateCuttingTaskHandler(cuttingSvc))                    // 创建裁剪任务
				tasks.GET("", transport.ListCuttingTasksHandler(cuttingSvc))                      // 裁剪任务列表
				tasks.GET("/order/:order_id", transport.GetCuttingTaskByOrderHandler(cuttingSvc)) // 根据订单ID获取任务
				tasks.GET("/:id", transport.GetCuttingTaskHandler(cuttingSvc))                    // 获取裁剪任务详情
			}

			// 裁剪批次路由
			batches := cutting.Group("/batches")
			{
				batches.POST("", transport.CreateCuttingBatchHandler(cuttingSvc))          // 创建裁剪批次（制菲）
				batches.GET("", transport.ListCuttingBatchesHandler(cuttingSvc))           // 裁剪批次列表
				batches.GET("/:id", transport.GetCuttingBatchHandler(cuttingSvc))          // 获取裁剪批次详情
				batches.DELETE("/:id", transport.DeleteCuttingBatchHandler(cuttingSvc))    // 删除裁剪批次
				batches.POST("/:id/print", transport.PrintCuttingBatchHandler(cuttingSvc)) // 打印裁剪批次
			}

			// 裁片监控路由
			pieces := cutting.Group("/pieces")
			{
				pieces.GET("", transport.ListCuttingPiecesHandler(cuttingSvc))                       // 裁片列表
				pieces.GET("/:id", transport.GetCuttingPieceHandler(cuttingSvc))                     // 获取裁片详情
				pieces.PUT("/:id/progress", transport.UpdateCuttingPieceProgressHandler(cuttingSvc)) // 更新裁片进度
			}
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

		// 自动注册路由配置到网关
		routeConfig := &cousul.RouteConfig{
			Prefix:        "/order", // 微服务路径
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
