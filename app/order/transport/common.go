package transport

import (
	"mule-cloud/app/order/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
func HealthHandler(svc services.ICommonService) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := svc.Health()
		response.Success(c, map[string]string{"status": result})
	}
}
