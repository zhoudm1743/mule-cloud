package transport

import (
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(svc services.ICommonService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ep := endpoint.HealthEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.CommonRequest{})
		if err != nil {
			response.Error(c, err.Error())
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
