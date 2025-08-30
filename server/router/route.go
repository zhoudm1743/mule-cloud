package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "mule-cloud/docs" // 导入生成的docs包
)

func Register(r *gin.Engine) {
	// 健康检查端点
	r.GET("/health", healthCheck)
	r.GET("/status", statusCheck)

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 注册API路由组
	api := r.Group("/api")
	regAdmin(api)
}

// healthCheck 健康检查端点
// @Summary 健康检查
// @Description 检查服务是否正常运行
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "服务正常"
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "服务运行正常",
	})
}

// statusCheck 状态检查端点
// @Summary 状态检查
// @Description 获取应用状态信息
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "应用状态信息"
// @Router /status [get]
func statusCheck(c *gin.Context) {
	// 返回基本的应用状态信息，避免循环导入
	c.JSON(http.StatusOK, gin.H{
		"status": "running",
		"name":   "mule-cloud",
		"health": "ok",
	})
}
