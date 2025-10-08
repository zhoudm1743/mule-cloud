package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetColorHandler 获取颜色处理器
func GetColorHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := binding.BindAll(c, &req); err != nil {
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

// GetAllColorsHandler 获取所有颜色处理器（不分页）
func GetAllColorsHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := binding.BindAll(c, &req); err != nil {
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

// ListColorsHandler 颜色列表处理器（分页）
func ListColorsHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListColorsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateColorHandler 创建颜色处理器
func CreateColorHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateColorEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateColorHandler 更新颜色处理器
func UpdateColorHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorUpdateRequest
		// 先绑定 JSON body（包含 required 字段）
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateColorEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteColorHandler 删除颜色处理器
func DeleteColorHandler(svc services.IColorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ColorListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteColorEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
