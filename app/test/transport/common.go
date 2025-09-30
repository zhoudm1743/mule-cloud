package transport

import (
	"mule-cloud/app/test/services"
	"mule-cloud/core/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(svc services.ICommonService) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := svc.Health()
		if err != nil {
			response.Error(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": resp,
		})
	}
}
