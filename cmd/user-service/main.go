package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/zhoudm1743/mule-cloud/internal/handler"
	"github.com/zhoudm1743/mule-cloud/internal/middleware"
	"github.com/zhoudm1743/mule-cloud/internal/repository"
	"github.com/zhoudm1743/mule-cloud/internal/service"
	"github.com/zhoudm1743/mule-cloud/pkg/auth"
	"github.com/zhoudm1743/mule-cloud/pkg/cache"
	"github.com/zhoudm1743/mule-cloud/pkg/config"
	"github.com/zhoudm1743/mule-cloud/pkg/database"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("../../configs")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logConfig := logger.Config{
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		Output:     cfg.Log.Output,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
	}
	appLogger := logger.NewLogger(logConfig)

	// 初始化数据库
	mongoConfig := database.MongoConfig{
		URI:      cfg.Database.MongoDB.URI,
		Database: cfg.Database.MongoDB.Database,
		Username: cfg.Database.MongoDB.Username,
		Password: cfg.Database.MongoDB.Password,
	}

	mongoDB, err := database.NewMongoDB(mongoConfig)
	if err != nil {
		appLogger.Fatal("Failed to connect to MongoDB", "error", err)
	}
	defer mongoDB.Close(context.Background())

	// 创建索引
	err = mongoDB.CreateIndexes(context.Background())
	if err != nil {
		appLogger.Error("Failed to create indexes", "error", err)
	}

	// 初始化Redis缓存
	cacheConfig := cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	redisCache, err := cache.NewRedisCache(cacheConfig)
	if err != nil {
		appLogger.Fatal("Failed to connect to Redis", "error", err)
	}

	cacheManager := cache.NewCacheManager(redisCache)

	// 初始化认证服务
	jwtConfig := auth.JWTConfig{
		SecretKey:       cfg.JWT.SecretKey,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
	}
	authService := auth.NewAuthService(jwtConfig)
	tokenManager := auth.NewTokenManager(jwtConfig)

	// 初始化仓储层
	userRepo := repository.NewUserRepository(mongoDB.GetDatabase())

	// 初始化服务层
	userService := service.NewUserService(userRepo, authService, cacheManager, appLogger)

	// 初始化处理器层
	userHandler := handler.NewUserHandler(userService, appLogger)

	// 初始化Gin引擎
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

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
			"service":   "user-service",
			"timestamp": time.Now(),
		})
	})

	// API路由
	setupRoutes(router, userHandler, tokenManager, appLogger, redisCache)

	// 启动服务器
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	appLogger.Info("User service started",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"mode", cfg.Server.Mode)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down user service...")

	// 给服务5秒时间完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("User service stopped")
}

// setupRoutes 设置路由
func setupRoutes(
	router *gin.Engine,
	userHandler *handler.UserHandler,
	tokenManager *auth.TokenManager,
	logger logger.Logger,
	cache cache.Cache,
) {
	// API版本前缀
	v1 := router.Group("/api/v1")

	// 认证相关路由（无需认证）
	auth := v1.Group("/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.RefreshToken)
	}

	// 用户相关路由（需要认证）
	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(tokenManager, logger))
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.POST("/change-password", userHandler.ChangePassword)
		users.POST("/logout", userHandler.Logout)
	}

	// 管理员路由（需要认证和权限）
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenManager, logger))
	// admin.Use(middleware.RequireRole("admin", logger)) // 暂时注释，需要实现角色检查
	{
		admin.POST("/users", userHandler.CreateUser)
		admin.GET("/users", userHandler.ListUsers)
		admin.GET("/users/:id", userHandler.GetUser)
		admin.PUT("/users/:id", userHandler.UpdateUser)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
		admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
	}

	// 限流中间件示例
	rateLimit := v1.Group("/limited")
	rateLimit.Use(middleware.IPRateLimitMiddleware(cache, logger, 10)) // 每分钟10次请求
	{
		rateLimit.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "rate limited endpoint"})
		})
	}
}
