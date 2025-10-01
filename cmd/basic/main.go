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
	customerSvc := services.NewCustomerService()
	orderTypeSvc := services.NewOrderTypeService()
	procedureSvc := services.NewProcedureService()
	salesmanSvc := services.NewSalesmanService()
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
		color := basic.Group("/colors")
		{
			color.GET("/:id", transport.GetColorHandler(colorSvc))       // 获取单个颜色
			color.GET("", transport.ListColorsHandler(colorSvc))         // 分页列表
			color.GET("/all", transport.GetAllColorsHandler(colorSvc))   // 获取所有（不分页）
			color.POST("", transport.CreateColorHandler(colorSvc))       // 创建颜色
			color.PUT("/:id", transport.UpdateColorHandler(colorSvc))    // 更新颜色
			color.DELETE("/:id", transport.DeleteColorHandler(colorSvc)) // 删除颜色
		}

		// 尺寸路由
		size := basic.Group("/sizes")
		{
			size.GET("/:id", transport.GetSizeHandler(sizeSvc))       // 获取单个尺寸
			size.GET("", transport.ListSizesHandler(sizeSvc))         // 分页列表
			size.GET("/all", transport.GetAllSizesHandler(sizeSvc))   // 获取所有（不分页）
			size.POST("", transport.CreateSizeHandler(sizeSvc))       // 创建尺寸
			size.PUT("/:id", transport.UpdateSizeHandler(sizeSvc))    // 更新尺寸
			size.DELETE("/:id", transport.DeleteSizeHandler(sizeSvc)) // 删除尺寸
		}

		// 客户路由
		customer := basic.Group("/customers")
		{
			customer.GET("/:id", transport.GetCustomerHandler(customerSvc))       // 获取单个客户
			customer.GET("", transport.ListCustomersHandler(customerSvc))         // 分页列表
			customer.GET("/all", transport.GetAllCustomersHandler(customerSvc))   // 获取所有（不分页）
			customer.POST("", transport.CreateCustomerHandler(customerSvc))       // 创建客户
			customer.PUT("/:id", transport.UpdateCustomerHandler(customerSvc))    // 更新客户
			customer.DELETE("/:id", transport.DeleteCustomerHandler(customerSvc)) // 删除客户
		}

		// 订单类型路由
		orderType := basic.Group("/order_types")
		{
			orderType.GET("/:id", transport.GetOrderTypeHandler(orderTypeSvc))       // 获取单个订单类型
			orderType.GET("", transport.ListOrderTypesHandler(orderTypeSvc))         // 分页列表
			orderType.GET("/all", transport.GetAllOrderTypesHandler(orderTypeSvc))   // 获取所有（不分页）
			orderType.POST("", transport.CreateOrderTypeHandler(orderTypeSvc))       // 创建订单类型
			orderType.PUT("/:id", transport.UpdateOrderTypeHandler(orderTypeSvc))    // 更新订单类型
			orderType.DELETE("/:id", transport.DeleteOrderTypeHandler(orderTypeSvc)) // 删除订单类型
		}

		// 工序路由
		procedure := basic.Group("/procedures")
		{
			procedure.GET("/:id", transport.GetProcedureHandler(procedureSvc))       // 获取单个工序
			procedure.GET("", transport.ListProceduresHandler(procedureSvc))         // 分页列表
			procedure.GET("/all", transport.GetAllProceduresHandler(procedureSvc))   // 获取所有（不分页）
			procedure.POST("", transport.CreateProcedureHandler(procedureSvc))       // 创建工序
			procedure.PUT("/:id", transport.UpdateProcedureHandler(procedureSvc))    // 更新工序
			procedure.DELETE("/:id", transport.DeleteProcedureHandler(procedureSvc)) // 删除工序
		}

		// 业务员路由
		salesman := basic.Group("/salesmans")
		{
			salesman.GET("/:id", transport.GetSalesmanHandler(salesmanSvc))       // 获取单个业务员
			salesman.GET("", transport.ListSalesmansHandler(salesmanSvc))         // 分页列表
			salesman.GET("/all", transport.GetAllSalesmansHandler(salesmanSvc))   // 获取所有（不分页）
			salesman.POST("", transport.CreateSalesmanHandler(salesmanSvc))       // 创建业务员
			salesman.PUT("/:id", transport.UpdateSalesmanHandler(salesmanSvc))    // 更新业务员
			salesman.DELETE("/:id", transport.DeleteSalesmanHandler(salesmanSvc)) // 删除业务员
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
			Prefix:        "/basic", // 微服务路径
			GatewayPrefix: "/admin", // 网关前缀（后台接口必须通过 /admin 访问）
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
