package transport

import (
	"mule-cloud/app/test/endpoint"
	"mule-cloud/app/test/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// AdminGetHandler 获取管理员信息
func AdminGetHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 参数绑定
		var req endpoint.AdminRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 2. 获取用户信息（从网关传递）
		userID := c.GetHeader("X-User-ID")
		username := c.GetHeader("X-Username")

		// 3. 调用endpoint
		ep := endpoint.GetAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.AdminRequest{
			ID:     req.ID,
			UserID: userID,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		// 4. 添加额外信息到响应
		result := resp.(endpoint.AdminResponse)
		result.RequestedBy = username

		response.Success(c, result)
	}
}

// AdminDeleteHandler 删除管理员
func AdminDeleteHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.AdminRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 获取操作用户
		userID := c.GetHeader("X-User-ID")
		username := c.GetHeader("X-Username")

		// 调用endpoint
		ep := endpoint.DeleteAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.AdminRequest{
			ID:     req.ID,
			UserID: userID,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		result := resp.(endpoint.DeleteAdminResponse)
		result.DeletedBy = username

		response.Success(c, result)
	}
}

// AdminCreateHandler 创建管理员
func AdminCreateHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.CreateAdminRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		username := c.GetHeader("X-Username")

		ep := endpoint.CreateAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		result := resp.(endpoint.CreateAdminResponse)
		result.CreatedBy = username

		response.Success(c, result)
	}
}

// AdminUpdateHandler 更新管理员
func AdminUpdateHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req endpoint.UpdateAdminRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		username := c.GetHeader("X-Username")

		ep := endpoint.UpdateAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		result := resp.(endpoint.UpdateAdminResponse)
		result.UpdatedBy = username

		response.Success(c, result)
	}
}
