package main

import (
	"flag"
	"fmt"
	"log"
	"mule-cloud/app/gateway/middleware"
	cfgPkg "mule-cloud/core/config"
	hystrixPkg "mule-cloud/core/hystrix"
	jwtPkg "mule-cloud/core/jwt"
	"mule-cloud/core/response"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

// Gateway API网关结构（增强版）
type Gateway struct {
	consulClient *api.Client
	routes       map[string]*RouteConfig
	jwtManager   *jwtPkg.JWTManager
	rateLimiter  *middleware.RateLimiter
	config       *cfgPkg.Config
}

// RouteConfig 路由配置
type RouteConfig struct {
	ServiceName string   // Consul服务名
	RequireAuth bool     // 是否需要认证
	RequireRole []string // 需要的角色（为空则只需登录）
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

	// 构建路由配置
	routes := make(map[string]*RouteConfig)
	for prefix, routeCfg := range cfg.Gateway.Routes {
		routes[prefix] = &RouteConfig{
			ServiceName: routeCfg.ServiceName,
			RequireAuth: routeCfg.RequireAuth,
			RequireRole: routeCfg.RequireRole,
		}
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
		routes:       routes,
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
		path := c.Request.URL.Path

		// 1. 匹配路由前缀
		var routeConfig *RouteConfig
		var matchedPrefix string
		for prefix, config := range gw.routes {
			if strings.HasPrefix(path, prefix) {
				routeConfig = config
				matchedPrefix = prefix
				break
			}
		}

		if routeConfig == nil {
			c.JSON(404, gin.H{"code": 404, "msg": "路由不存在"})
			return
		}

		// 设置服务名称（供Hystrix中间件使用）
		c.Set("service_name", routeConfig.ServiceName)

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
		targetURL, err := gw.getServiceAddress(routeConfig.ServiceName)
		if err != nil {
			log.Printf("[网关错误] 服务不可用: %s, 错误: %v", routeConfig.ServiceName, err)
			c.JSON(503, gin.H{"code": 503, "msg": fmt.Sprintf("服务不可用: %s", routeConfig.ServiceName)})
			return
		}

		// 4. 构建反向代理
		target, _ := url.Parse(targetURL)
		proxy := httputil.NewSingleHostReverseProxy(target)

		// 5. 修改请求
		originalPath := c.Request.URL.Path
		c.Request.URL.Path = strings.TrimPrefix(originalPath, matchedPrefix)
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

		c.Request.Host = target.Host

		// 7. 记录日志
		log.Printf("[网关转发] %s %s → %s%s (服务: %s, 用户: %v)",
			c.Request.Method,
			originalPath,
			targetURL,
			c.Request.URL.Path,
			routeConfig.ServiceName,
			c.GetString("username"),
		)

		// 8. 执行代理转发
		proxy.ServeHTTP(c.Writer, c.Request)

		// 9. 记录响应时间
		duration := time.Since(startTime)
		log.Printf("[网关响应] %s %s 耗时: %v", c.Request.Method, originalPath, duration)
	}
}

// loginHandler 登录接口（示例）
func (gw *Gateway) loginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
			return
		}

		// 这里应该查询数据库验证用户名密码
		// 为了演示，简化处理
		if req.Username == "admin" && req.Password == "admin123" {
			token, err := gw.jwtManager.GenerateToken("1", "admin", []string{"admin", "user"})
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "生成token失败"})
				return
			}

			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "登录成功",
				"data": gin.H{
					"token":    token,
					"username": "admin",
					"roles":    []string{"admin", "user"},
				},
			})
			return
		}

		if req.Username == "user" && req.Password == "user123" {
			token, err := gw.jwtManager.GenerateToken("2", "user", []string{"user"})
			if err != nil {
				c.JSON(500, gin.H{"code": 500, "msg": "生成token失败"})
				return
			}

			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "登录成功",
				"data": gin.H{
					"token":    token,
					"username": "user",
					"roles":    []string{"user"},
				},
			})
			return
		}

		c.JSON(401, gin.H{"code": 401, "msg": "用户名或密码错误"})
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
			for _, routeConfig := range gw.routes {
				serviceSet[routeConfig.ServiceName] = true
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

	// 初始化Hystrix熔断器
	if cfg.Hystrix.Enabled {
		hystrixPkg.Init()
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
		public.POST("/login", gateway.loginHandler())
		public.GET("/health", gateway.healthHandler())
	}

	// 熔断器管理接口
	admin := r.Group("/gateway")
	{
		admin.GET("/hystrix/metrics", middleware.HystrixMetricsHandler())
		admin.GET("/hystrix/metrics/:service", middleware.HystrixMetricsHandler())
	}

	// 业务接口（需要认证 + 限流 + 熔断）
	api := r.Group("")
	if cfg.Gateway.RateLimit.Enabled {
		api.Use(gateway.rateLimiter.Middleware()) // 限流
	}
	api.Use(middleware.OptionalAuth(gateway.jwtManager)) // 可选认证（根据路由配置决定）
	if cfg.Hystrix.Enabled {
		api.Use(middleware.HystrixMiddleware()) // Hystrix熔断器
	}
	{
		api.Any("/test/*path", gateway.proxyHandler())
		api.Any("/basic/*path", gateway.proxyHandler())
		api.Any("/admin/*path", gateway.proxyHandler())
	}

	// 启动网关
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("========================================")
	log.Printf("🚀 %s 启动成功", cfg.Server.Name)
	log.Printf("📍 监听端口: %s", port)
	if cfg.Consul.Enabled {
		log.Printf("🔗 Consul地址: %s", cfg.Consul.Address)
	}
	log.Printf("🔐 JWT认证: 已启用")
	if cfg.Gateway.RateLimit.Enabled {
		log.Printf("⚡ 限流保护: 每秒%d请求", cfg.Gateway.RateLimit.Rate)
	}
	if cfg.Hystrix.Enabled {
		log.Printf("🔥 Hystrix熔断: 已启用")
	}
	log.Printf("📋 路由配置:")
	for prefix, routeConfig := range gateway.routes {
		authStr := "公开"
		if routeConfig.RequireAuth {
			if len(routeConfig.RequireRole) > 0 {
				authStr = fmt.Sprintf("需要角色: %v", routeConfig.RequireRole)
			} else {
				authStr = "需要登录"
			}
		}
		log.Printf("   %s/* → %s (%s)", prefix, routeConfig.ServiceName, authStr)
	}
	log.Printf("========================================")
	log.Printf("💡 测试命令:")
	log.Printf("   # 登录获取token")
	log.Printf("   curl -X POST http://localhost:8080/api/login -H \"Content-Type: application/json\" -d '{\"username\":\"admin\",\"password\":\"admin123\"}'")
	log.Printf("")
	log.Printf("   # 使用token访问需要认证的接口")
	log.Printf("   curl http://localhost:8080/test/admin/123 -H \"Authorization: Bearer {token}\"")
	log.Printf("")
	log.Printf("   # 访问公开接口（无需token）")
	log.Printf("   curl http://localhost:8080/basic/color/1")
	log.Printf("========================================")

	if err := r.Run(port); err != nil {
		log.Fatalf("网关启动失败: %v", err)
	}
}
