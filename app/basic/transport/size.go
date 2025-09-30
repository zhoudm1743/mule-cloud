package transport

import (
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetSizeHandler 获取尺寸处理器
func GetSizeHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.SizeRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetSizeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllSizesHandler 获取所有尺寸处理器
func GetAllSizesHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ep := endpoint.GetAllSizesEndpoint(svc)
		resp, err := ep(c.Request.Context(), nil)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
