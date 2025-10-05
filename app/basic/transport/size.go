package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetSizeHandler 获取尺寸处理器
func GetSizeHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeGetRequest
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

// GetAllSizesHandler 获取所有尺寸处理器（不分页）
func GetAllSizesHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllSizesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListSizesHandler 尺寸列表处理器（分页）
func ListSizesHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListSizesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateSizeHandler 创建尺寸处理器
func CreateSizeHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateSizeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateSizeHandler 更新尺寸处理器
func UpdateSizeHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeUpdateRequest
		// 先绑定 JSON body（包含 required 字段）
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		// 再绑定 URI 参数（ID）
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateSizeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteSizeHandler 删除尺寸处理器
func DeleteSizeHandler(svc services.ISizeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SizeGetRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteSizeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
