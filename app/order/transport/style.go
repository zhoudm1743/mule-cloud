package transport

import (
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/endpoint"
	"mule-cloud/app/order/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetStyleHandler 获取款式处理器
func GetStyleHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetStyleEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllStylesHandler 获取所有款式处理器（不分页）
func GetAllStylesHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllStylesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListStylesHandler 款式列表处理器（分页）
func ListStylesHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListStylesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateStyleHandler 创建款式处理器
func CreateStyleHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateStyleEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateStyleHandler 更新款式处理器
func UpdateStyleHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleUpdateRequest
		// 先绑定 JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		// 再绑定 URI 参数（ID）
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateStyleEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteStyleHandler 删除款式处理器
func DeleteStyleHandler(svc services.IStyleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.StyleListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteStyleEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
