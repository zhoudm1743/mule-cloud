package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/internal/middleware"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
	UserService string
	// 其他服务地址可以在这里添加
	OrderService        string
	ProductionService   string
	TimesheetService    string
	PayrollService      string
	ReportService       string
	MasterDataService   string
	NotificationService string
	FileService         string
}

// LoadServiceRegistry 加载服务注册表
func LoadServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		UserService:         getEnv("USER_SERVICE_URL", "http://localhost:8001"),
		OrderService:        getEnv("ORDER_SERVICE_URL", "http://localhost:8003"),
		ProductionService:   getEnv("PRODUCTION_SERVICE_URL", "http://localhost:8004"),
		TimesheetService:    getEnv("TIMESHEET_SERVICE_URL", "http://localhost:8005"),
		PayrollService:      getEnv("PAYROLL_SERVICE_URL", "http://localhost:8006"),
		ReportService:       getEnv("REPORT_SERVICE_URL", "http://localhost:8007"),
		MasterDataService:   getEnv("MASTER_DATA_SERVICE_URL", "http://localhost:8002"),
		NotificationService: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8008"),
		FileService:         getEnv("FILE_SERVICE_URL", "http://localhost:8009"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 初始化日志
	logConfig := logger.Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}
	appLogger := logger.NewLogger(logConfig)

	// 加载服务注册表
	services := LoadServiceRegistry()

	// 初始化Gin引擎
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 添加中间件
	router.Use(middleware.RecoveryMiddleware(appLogger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(appLogger))
	router.Use(middleware.SecurityMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "api-gateway",
			"timestamp": time.Now(),
			"services": gin.H{
				"user-service": services.UserService,
			},
		})
	})

	// API路由
	setupRoutes(router, services, appLogger)

	// 启动服务器
	port := getEnv("GATEWAY_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start gateway", "error", err)
		}
	}()

	appLogger.Info("API Gateway started", "port", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down API Gateway...")

	// 给服务5秒时间完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Gateway forced to shutdown", "error", err)
	}

	appLogger.Info("API Gateway stopped")
}

// setupRoutes 设置路由
func setupRoutes(router *gin.Engine, services *ServiceRegistry, logger logger.Logger) {
	api := router.Group("/api")

	// 用户服务路由
	userGroup := api.Group("/v1")
	{
		// 认证相关
		userGroup.Any("/auth/*path", createProxy(services.UserService, logger))

		// 用户相关
		userGroup.Any("/users/*path", createProxy(services.UserService, logger))

		// 管理员相关
		userGroup.Any("/admin/users/*path", createProxy(services.UserService, logger))
	}

	// 其他服务的路由可以在这里添加
	// 订单服务
	orderGroup := api.Group("/v1")
	{
		orderGroup.Any("/orders/*path", createProxy(services.OrderService, logger))
		orderGroup.Any("/customers/*path", createProxy(services.OrderService, logger))
		orderGroup.Any("/styles/*path", createProxy(services.OrderService, logger))
		orderGroup.Any("/salespersons/*path", createProxy(services.OrderService, logger))
	}

	// 生产服务
	productionGroup := api.Group("/v1")
	{
		productionGroup.Any("/production/*path", createProxy(services.ProductionService, logger))
		productionGroup.Any("/cutting/*path", createProxy(services.ProductionService, logger))
	}

	// 工时服务
	timesheetGroup := api.Group("/v1")
	{
		timesheetGroup.Any("/work-reports/*path", createProxy(services.TimesheetService, logger))
		timesheetGroup.Any("/timesheets/*path", createProxy(services.TimesheetService, logger))
	}

	// 工资服务
	payrollGroup := api.Group("/v1")
	{
		payrollGroup.Any("/payroll/*path", createProxy(services.PayrollService, logger))
		payrollGroup.Any("/salary/*path", createProxy(services.PayrollService, logger))
	}

	// 报表服务
	reportGroup := api.Group("/v1")
	{
		reportGroup.Any("/reports/*path", createProxy(services.ReportService, logger))
		reportGroup.Any("/statistics/*path", createProxy(services.ReportService, logger))
		reportGroup.Any("/dashboard/*path", createProxy(services.ReportService, logger))
	}

	// 基础数据服务
	masterDataGroup := api.Group("/v1")
	{
		masterDataGroup.Any("/processes/*path", createProxy(services.MasterDataService, logger))
		masterDataGroup.Any("/sizes/*path", createProxy(services.MasterDataService, logger))
		masterDataGroup.Any("/colors/*path", createProxy(services.MasterDataService, logger))
		masterDataGroup.Any("/workers/*path", createProxy(services.MasterDataService, logger))
	}

	// 通知服务
	notificationGroup := api.Group("/v1")
	{
		notificationGroup.Any("/notifications/*path", createProxy(services.NotificationService, logger))
	}

	// 文件服务
	fileGroup := api.Group("/v1")
	{
		fileGroup.Any("/files/*path", createProxy(services.FileService, logger))
		fileGroup.Any("/upload/*path", createProxy(services.FileService, logger))
	}
}

// createProxy 创建反向代理
func createProxy(targetURL string, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		target, err := url.Parse(targetURL)
		if err != nil {
			logger.Error("Invalid target URL", "url", targetURL, "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "Service unavailable",
			})
			return
		}

		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(target)

		// 自定义Director以修改请求
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host

			// 保持原始路径
			originalPath := c.Param("path")
			if originalPath != "" {
				// 获取路由组前缀
				routePrefix := getRoutePrefix(c.FullPath())
				// 构建完整路径
				req.URL.Path = "/api/v1" + routePrefix + originalPath
			}

			// 设置X-Forwarded-*头
			req.Header.Set("X-Forwarded-For", c.ClientIP())
			req.Header.Set("X-Forwarded-Proto", "http")
			req.Header.Set("X-Forwarded-Host", c.Request.Host)

			// 传递请求ID
			if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
				req.Header.Set("X-Request-ID", requestID)
			}
		}

		// 自定义错误处理
		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			logger.Error("Proxy error", "url", targetURL, "error", err)
			rw.WriteHeader(http.StatusBadGateway)
			rw.Write([]byte(`{"code":502,"message":"Service unavailable"}`))
		}

		// 执行代理
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// getRoutePrefix 获取路由前缀
func getRoutePrefix(fullPath string) string {
	// 从完整路径中提取路由前缀
	parts := strings.Split(fullPath, "/")
	if len(parts) >= 4 {
		return "/" + parts[3] // /api/v1/xxx/*path -> /xxx
	}
	return ""
}
