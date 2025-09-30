package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetColorHandler 获取颜色处理器
func GetColorHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetColorEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllColorsHandler 获取所有颜色处理器
func GetAllColorsHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		ep := endpoint.GetAllColorsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
