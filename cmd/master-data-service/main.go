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
	"github.com/zhoudm1743/mule-cloud/pkg/config"
	"github.com/zhoudm1743/mule-cloud/pkg/database"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title 信芙云基础数据服务 API
// @version 1.0
// @description 服装生产管理系统基础数据服务，提供工序、尺码、颜色、客户、业务员等基础数据管理功能
// @host localhost:8002
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	appLogger.Info("Starting Master Data Service...")

	// 连接MongoDB
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
	defer func() {
		if err := mongoDB.Close(context.Background()); err != nil {
			appLogger.Error("Failed to disconnect from MongoDB", "error", err)
		}
	}()

	// 连接Redis (暂时注释掉，基础数据服务暂不需要Redis)
	// redisCache, err := cache.NewRedisCache(cfg.Cache.Redis)
	// if err != nil {
	//	appLogger.Fatal("Failed to connect to Redis", "error", err)
	// }

	// 初始化JWT认证服务
	jwtConfig := auth.JWTConfig{
		SecretKey:       cfg.JWT.SecretKey,
		AccessTokenTTL:  cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		Issuer:          cfg.JWT.Issuer,
	}

	tokenManager := auth.NewTokenManager(jwtConfig)

	// 初始化Repository层
	processRepo := repository.NewProcessRepository(mongoDB.GetDatabase(), appLogger)
	sizeRepo := repository.NewSizeRepository(mongoDB.GetDatabase(), appLogger)
	colorRepo := repository.NewColorRepository(mongoDB.GetDatabase(), appLogger)
	customerRepo := repository.NewCustomerRepository(mongoDB.GetDatabase(), appLogger)
	salespersonRepo := repository.NewSalespersonRepository(mongoDB.GetDatabase(), appLogger)

	// 初始化Service层
	processService := service.NewProcessService(processRepo, appLogger)
	sizeService := service.NewSizeService(sizeRepo, appLogger)
	colorService := service.NewColorService(colorRepo, appLogger)
	customerService := service.NewCustomerService(customerRepo, appLogger)
	salespersonService := service.NewSalespersonService(salespersonRepo, appLogger)

	// 初始化Handler层
	masterDataHandler := handler.NewMasterDataHandler(
		processService,
		sizeService,
		colorService,
		customerService,
		salespersonService,
		appLogger,
	)

	// 创建索引
	if err := createIndexes(context.Background(), mongoDB.GetDatabase(), appLogger); err != nil {
		appLogger.Error("Failed to create indexes", "error", err)
	}

	// 初始化Gin路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 全局中间件
	router.Use(middleware.LoggingMiddleware(appLogger, middleware.LoggingConfig{}))
	router.Use(middleware.RecoveryMiddleware(appLogger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.SecurityMiddleware())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":   "master-data-service",
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	// API版本前缀
	v1 := router.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(tokenManager, appLogger))

	// 工序路由
	processes := v1.Group("/processes")
	{
		processes.POST("", masterDataHandler.CreateProcess)
		processes.GET("", masterDataHandler.ListProcesses)
		processes.GET("/active", masterDataHandler.GetActiveProcesses)
		processes.GET("/:id", masterDataHandler.GetProcess)
		processes.PUT("/:id", masterDataHandler.UpdateProcess)
		processes.DELETE("/:id", masterDataHandler.DeleteProcess)
	}

	// 尺码路由
	sizes := v1.Group("/sizes")
	{
		sizes.POST("", masterDataHandler.CreateSize)
		sizes.GET("", masterDataHandler.ListSizes)
		sizes.GET("/active", masterDataHandler.GetActiveSizes)
		sizes.GET("/:id", masterDataHandler.GetSize)
		sizes.PUT("/:id", masterDataHandler.UpdateSize)
		sizes.DELETE("/:id", masterDataHandler.DeleteSize)
	}

	// 颜色路由
	colors := v1.Group("/colors")
	{
		colors.POST("", masterDataHandler.CreateColor)
		colors.GET("", masterDataHandler.ListColors)
		colors.GET("/:id", func(c *gin.Context) {
			// TODO: 实现获取单个颜色
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		colors.PUT("/:id", func(c *gin.Context) {
			// TODO: 实现更新颜色
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		colors.DELETE("/:id", func(c *gin.Context) {
			// TODO: 实现删除颜色
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
	}

	// 客户路由
	customers := v1.Group("/customers")
	{
		customers.POST("", masterDataHandler.CreateCustomer)
		customers.GET("", masterDataHandler.ListCustomers)
		customers.GET("/:id", func(c *gin.Context) {
			// TODO: 实现获取单个客户
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		customers.PUT("/:id", func(c *gin.Context) {
			// TODO: 实现更新客户
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		customers.DELETE("/:id", func(c *gin.Context) {
			// TODO: 实现删除客户
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
	}

	// 业务员路由
	salespersons := v1.Group("/salespersons")
	{
		salespersons.POST("", masterDataHandler.CreateSalesperson)
		salespersons.GET("", masterDataHandler.ListSalespersons)
		salespersons.GET("/:id", func(c *gin.Context) {
			// TODO: 实现获取单个业务员
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		salespersons.PUT("/:id", func(c *gin.Context) {
			// TODO: 实现更新业务员
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
		salespersons.DELETE("/:id", func(c *gin.Context) {
			// TODO: 实现删除业务员
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented yet"})
		})
	}

	// 启动HTTP服务器
	port := ":8002"
	server := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start master data service", "error", err)
		}
	}()

	appLogger.Info("Master Data Service started", "port", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Master Data Service...")

	// 给服务5秒时间完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Master Data Service forced to shutdown", "error", err)
	}

	appLogger.Info("Master Data Service stopped")
}

// createIndexes 创建数据库索引
func createIndexes(ctx context.Context, db *mongo.Database, logger logger.Logger) error {
	logger.Info("Creating database indexes...")

	// 工序索引
	processCollection := db.Collection("processes")
	if _, err := processCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return fmt.Errorf("failed to create process code index: %w", err)
	}

	if _, err := processCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "category", Value: 1}, {Key: "is_active", Value: 1}},
	}); err != nil {
		return fmt.Errorf("failed to create process category index: %w", err)
	}

	// 尺码索引
	sizeCollection := db.Collection("sizes")
	if _, err := sizeCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return fmt.Errorf("failed to create size code index: %w", err)
	}

	// 颜色索引
	colorCollection := db.Collection("colors")
	if _, err := colorCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return fmt.Errorf("failed to create color code index: %w", err)
	}

	// 客户索引
	customerCollection := db.Collection("customers")
	if _, err := customerCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return fmt.Errorf("failed to create customer code index: %w", err)
	}

	if _, err := customerCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "customer_type", Value: 1}, {Key: "region", Value: 1}},
	}); err != nil {
		return fmt.Errorf("failed to create customer type index: %w", err)
	}

	// 业务员索引
	salespersonCollection := db.Collection("salespersons")
	if _, err := salespersonCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return fmt.Errorf("failed to create salesperson code index: %w", err)
	}

	if _, err := salespersonCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "department", Value: 1}, {Key: "status", Value: 1}},
	}); err != nil {
		return fmt.Errorf("failed to create salesperson department index: %w", err)
	}

	logger.Info("Database indexes created successfully")
	return nil
}
