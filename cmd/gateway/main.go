package main

import (
	"flag"
	"fmt"
	"log"
	"mule-cloud/app/gateway/middleware"
	cfgPkg "mule-cloud/core/config"
	hystrixPkg "mule-cloud/core/hystrix"
	jwtPkg "mule-cloud/core/jwt"
	loggerPkg "mule-cloud/core/logger"
	"mule-cloud/core/response"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

// Gateway API网关结构（增强版）
type Gateway struct {
	consulClient *api.Client
	routeManager *middleware.DynamicRouteManager // 动态路由管理器
	jwtManager   *jwtPkg.JWTManager
	rateLimiter  *middleware.RateLimiter
	config       *cfgPkg.Config
}

// NewGateway 创建增强版网关实例
func NewGateway(cfg *cfgPkg.Config) (*Gateway, error) {
	// 连接Consul
	var client *api.Client
	if cfg.Consul.Enabled {
		config := api.DefaultConfig()
		config.Address = cfg.Consul.Address
		config.Scheme = cfg.Consul.Scheme

		var err error
		client, err = api.NewClient(config)
		if err != nil {
			return nil, fmt.Errorf("连接Consul失败: %v", err)
		}
	}

	// 创建动态路由管理器
	var routeManager *middleware.DynamicRouteManager
	if cfg.Consul.Enabled && client != nil {
		routeManager = middleware.NewDynamicRouteManager(client)
		log.Println("✅ 启用动态路由管理器 (基于Consul KV)")

		// 从配置文件迁移路由到Consul（如果Consul中没有配置）
		if len(routeManager.GetAllRoutes()) == 0 && len(cfg.Gateway.Routes) > 0 {
			log.Println("🔄 检测到Consul中无路由配置，正在从配置文件迁移...")
			for prefix, routeCfg := range cfg.Gateway.Routes {
				config := &middleware.RouteConfig{
					ServiceName:   routeCfg.ServiceName,
					GatewayPrefix: "", // 默认无前缀，保持兼容
					RequireAuth:   routeCfg.RequireAuth,
					RequireRole:   routeCfg.RequireRole,
				}
				if err := routeManager.AddRoute(prefix, config); err != nil {
					log.Printf("⚠️  迁移路由配置失败 (%s): %v", prefix, err)
				}
			}
		}

		// 从配置文件迁移Hystrix配置到Consul（如果Consul中没有配置）
		if len(routeManager.GetAllHystrixConfigs()) == 0 && len(cfg.Hystrix.Command) > 0 {
			log.Println("🔄 检测到Consul中无Hystrix配置，正在从配置文件迁移...")
			for serviceName, cmdCfg := range cfg.Hystrix.Command {
				config := &middleware.DynamicHystrixConfig{
					Timeout:                cmdCfg.Timeout,
					MaxConcurrentRequests:  cmdCfg.MaxConcurrentRequests,
					RequestVolumeThreshold: cmdCfg.RequestVolumeThreshold,
					SleepWindow:            cmdCfg.SleepWindow,
					ErrorPercentThreshold:  cmdCfg.ErrorPercentThreshold,
				}
				if err := routeManager.AddHystrixConfig(serviceName, config); err != nil {
					log.Printf("⚠️  迁移Hystrix配置失败 (%s): %v", serviceName, err)
				}
			}
		}
	} else {
		log.Println("⚠️  Consul未启用，动态路由功能将不可用")
	}

	// JWT管理器
	jwtSecret := []byte(cfg.JWT.SecretKey)
	expireTime := time.Duration(cfg.JWT.ExpireTime) * time.Hour

	// 限流器
	var rateLimiter *middleware.RateLimiter
	if cfg.Gateway.RateLimit.Enabled {
		rateLimiter = middleware.NewRateLimiter(cfg.Gateway.RateLimit.Rate)
	}

	return &Gateway{
		consulClient: client,
		routeManager: routeManager,
		jwtManager:   jwtPkg.NewJWTManager(jwtSecret, expireTime),
		rateLimiter:  rateLimiter,
		config:       cfg,
	}, nil
}

// getServiceAddress 从Consul获取健康的服务地址
func (gw *Gateway) getServiceAddress(serviceName string) (string, error) {
	services, _, err := gw.consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("查询服务失败: %v", err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("未找到可用的服务实例: %s", serviceName)
	}

	// 简单负载均衡：返回第一个健康实例
	service := services[0].Service
	return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

// proxyHandler 反向代理处理器（增强版）
func (gw *Gateway) proxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		originalPath := c.Request.URL.Path

		// 1. 匹配路由配置（从动态路由管理器获取）
		var routeConfig *middleware.RouteConfig
		var serviceName string

		if gw.routeManager != nil {
			// 使用动态路由
			routes := gw.routeManager.GetAllRoutes()
			for routePrefix, config := range routes {
				// 构建完整的匹配路径（网关前缀 + 路由前缀）
				fullPrefix := config.GatewayPrefix + routePrefix
				if fullPrefix == "" {
					fullPrefix = routePrefix
				}
				if strings.HasPrefix(originalPath, fullPrefix) {
					routeConfig = config
					serviceName = config.ServiceName
					break
				}
			}
		} else {
			// 降级到静态路由（从配置文件）
			for prefix, staticCfg := range gw.config.Gateway.Routes {
				if strings.HasPrefix(originalPath, prefix) {
					routeConfig = &middleware.RouteConfig{
						ServiceName: staticCfg.ServiceName,
						RequireAuth: staticCfg.RequireAuth,
						RequireRole: staticCfg.RequireRole,
					}
					serviceName = staticCfg.ServiceName
					break
				}
			}
		}

		if routeConfig == nil {
			c.JSON(404, gin.H{"code": 404, "msg": "路由不存在"})
			return
		}

		// 设置服务名称（供Hystrix中间件使用）
		c.Set("service_name", serviceName)

		// 2. 认证检查（如果需要）
		if routeConfig.RequireAuth {
			claimsValue, exists := c.Get("claims")
			if !exists {
				c.JSON(401, gin.H{"code": 401, "msg": "需要认证"})
				return
			}

			// 角色检查
			if len(routeConfig.RequireRole) > 0 {
				claims := claimsValue.(*jwtPkg.Claims)
				if !claims.HasAnyRole(routeConfig.RequireRole...) {
					c.JSON(403, gin.H{"code": 403, "msg": "权限不足"})
					return
				}
			}
		}

		// 3. 从Consul获取服务地址
		targetURL, err := gw.getServiceAddress(serviceName)
		if err != nil {
			log.Printf("[网关错误] 服务不可用: %s, 错误: %v", serviceName, err)
			c.JSON(503, gin.H{"code": 503, "msg": fmt.Sprintf("服务不可用: %s", serviceName)})
			return
		}

		// 4. 构建反向代理
		target, _ := url.Parse(targetURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// 5. 修改请求路径（去掉网关配置的前缀）
		targetPath := originalPath
		if routeConfig.GatewayPrefix != "" && strings.HasPrefix(originalPath, routeConfig.GatewayPrefix) {
			// 去掉网关前缀，保留路由前缀
			targetPath = strings.TrimPrefix(originalPath, routeConfig.GatewayPrefix)
			if targetPath == "" {
				targetPath = "/"
			}
		}
		c.Request.URL.Path = targetPath
		c.Request.URL.Host = target.Host
		c.Request.URL.Scheme = target.Scheme

		// 6. 设置转发头（包括用户信息）
		c.Request.Header.Set("X-Forwarded-Host", c.Request.Host)
		c.Request.Header.Set("X-Real-IP", c.ClientIP())
		c.Request.Header.Set("X-Gateway", "mule-cloud-gateway-")

		// 传递用户信息到后端服务
		if userID, exists := c.Get("user_id"); exists {
			c.Request.Header.Set("X-User-ID", userID.(string))
		}
		if username, exists := c.Get("username"); exists {
			c.Request.Header.Set("X-Username", username.(string))
		}
		if tenantID, exists := c.Get("tenant_id"); exists {
			c.Request.Header.Set("X-Tenant-ID", tenantID.(string))
		}
		// ✅ 新增：传递租户代码（用于数据库连接）
		if tenantCode, exists := c.Get("tenant_code"); exists {
			c.Request.Header.Set("X-Tenant-Code", tenantCode.(string))
		}
		if rolesValue, exists := c.Get("roles"); exists {
			if roles, ok := rolesValue.([]string); ok && len(roles) > 0 {
				c.Request.Header.Set("X-Roles", strings.Join(roles, ","))
			}
		}
		
		// ✅ 重要：转发前端发送的 X-Tenant-Context header（用于超管切换租户）
		// 这个 header 是前端直接发送的，不在 JWT token 中，需要单独转发
		if contextTenant := c.GetHeader("X-Tenant-Context"); contextTenant != "" {
			c.Request.Header.Set("X-Tenant-Context", contextTenant)
			log.Printf("[网关转发] 转发租户上下文: %s", contextTenant)
		}

		c.Request.Host = target.Host

		// 7. 记录日志
		log.Printf("[网关转发] %s %s → %s%s (服务: %s, 前缀: %s, 用户: %v)",
			c.Request.Method,
			originalPath,
			targetURL,
			c.Request.URL.Path,
			serviceName,
			routeConfig.GatewayPrefix,
			c.GetString("username"),
		)

		// 8. 执行代理转发
		proxy.ServeHTTP(c.Writer, c.Request)

		// 9. 记录响应时间
		duration := time.Since(startTime)
		log.Printf("[网关响应] %s %s 耗时: %v", c.Request.Method, originalPath, duration)
	}
}

// healthHandler 健康检查
func (gw *Gateway) healthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		healthStatus := gin.H{
			"status":  "healthy",
			"gateway": gw.config.Server.Name,
		}

		// 检查Consul连接
		if gw.config.Consul.Enabled && gw.consulClient != nil {
			_, err := gw.consulClient.Agent().Self()
			if err != nil {
				c.JSON(503, gin.H{"status": "unhealthy", "error": "Consul连接失败"})
				return
			}
			healthStatus["consul"] = gw.config.Consul.Address

			// 检查服务状态
			services := make(map[string]string)
			serviceSet := make(map[string]bool)

			// 从动态路由管理器获取服务列表
			if gw.routeManager != nil {
				routes := gw.routeManager.GetAllRoutes()
				for _, routeConfig := range routes {
					serviceSet[routeConfig.ServiceName] = true
				}
			} else {
				// 降级到静态配置
				for _, routeConfig := range gw.config.Gateway.Routes {
					serviceSet[routeConfig.ServiceName] = true
				}
			}

			for svcName := range serviceSet {
				addr, err := gw.getServiceAddress(svcName)
				if err != nil {
					services[svcName] = "不可用"
				} else {
					services[svcName] = addr
				}
			}
			healthStatus["services"] = services
		}

		c.JSON(200, healthStatus)
	}
}

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/gateway.yaml", "配置文件路径")
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

	// 初始化Hystrix熔断器
	if cfg.Hystrix.Enabled {
		// 从配置文件读取服务级别配置
		commands := make(map[string]hystrixPkg.Config)
		for serviceName, cmdCfg := range cfg.Hystrix.Command {
			commands[serviceName] = hystrixPkg.Config{
				Timeout:                cmdCfg.Timeout,
				MaxConcurrentRequests:  cmdCfg.MaxConcurrentRequests,
				RequestVolumeThreshold: cmdCfg.RequestVolumeThreshold,
				SleepWindow:            cmdCfg.SleepWindow,
				ErrorPercentThreshold:  cmdCfg.ErrorPercentThreshold,
			}
		}
		hystrixPkg.InitWithConfig(commands)
	}

	// 创建网关实例
	gateway, err := NewGateway(cfg)
	if err != nil {
		log.Fatalf("创建网关失败: %v", err)
	}

	// 创建Gin路由
	gin.SetMode(cfg.Server.Mode)
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())                         // 日志
	r.Use(response.RecoveryMiddleware())        // 统一错误恢复
	r.Use(response.UnifiedResponseMiddleware()) // 统一响应
	r.Use(middleware.CORS())                    // 跨域

	// 公开接口（无需认证）
	public := r.Group("/api")
	{
		public.GET("/health", gateway.healthHandler())
	}

	// 网关管理接口（动态路由和熔断器管理）
	admin := r.Group("/gateway")
	{
		// Hystrix 指标监控
		admin.GET("/hystrix/metrics", middleware.HystrixMetricsHandler())
		admin.GET("/hystrix/metrics/:service", middleware.HystrixMetricsHandler())

		// 动态路由管理 API（需要动态路由管理器）
		if gateway.routeManager != nil {
			adminHandlers := middleware.NewAdminHandlers(gateway.routeManager)

			// 路由配置管理
			adminAPI := admin.Group("/admin")
			{
				// 路由管理
				adminAPI.GET("/routes", adminHandlers.ListRoutes)
				adminAPI.GET("/routes/*prefix", adminHandlers.GetRoute)
				adminAPI.POST("/routes", adminHandlers.AddRoute)
				adminAPI.PUT("/routes/*prefix", adminHandlers.UpdateRoute)
				adminAPI.DELETE("/routes/*prefix", adminHandlers.DeleteRoute)

				// Hystrix 配置管理
				adminAPI.GET("/hystrix", adminHandlers.ListHystrixConfigs)
				adminAPI.GET("/hystrix/:service", adminHandlers.GetHystrixConfig)
				adminAPI.POST("/hystrix", adminHandlers.AddHystrixConfig)
				adminAPI.PUT("/hystrix/:service", adminHandlers.UpdateHystrixConfig)
				adminAPI.DELETE("/hystrix/:service", adminHandlers.DeleteHystrixConfig)

				// 配置重载
				adminAPI.POST("/reload", adminHandlers.ReloadConfig)
			}
		}
	}

	// 业务接口（动态路由）
	// 使用 NoRoute 作为兜底，根据路由配置决定是否需要认证
	var handlers []gin.HandlerFunc
	if cfg.Gateway.RateLimit.Enabled {
		handlers = append(handlers, gateway.rateLimiter.Middleware())
	}
	handlers = append(handlers, middleware.OptionalAuth(gateway.jwtManager))
	if cfg.Hystrix.Enabled {
		handlers = append(handlers, middleware.HystrixMiddleware())
	}
	handlers = append(handlers, gateway.proxyHandler())
	r.NoRoute(handlers...)
	// 启动网关
	port := fmt.Sprintf(":%d", cfg.Server.Port)

	loggerPkg.Info("🚀 Gateway 启动中...",
		zap.String("service", cfg.Server.Name),
		zap.Int("port", cfg.Server.Port),
	)
	if err := r.Run(port); err != nil {
		log.Fatalf("网关启动失败: %v", err)
	}
}
