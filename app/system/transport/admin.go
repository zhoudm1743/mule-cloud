package transport

import (
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/endpoint"
	"mule-cloud/app/system/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetAdminHandler 获取管理员处理器
func GetAdminHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllAdminsHandler 获取所有管理员处理器（不分页）
func GetAllAdminsHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllAdminsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListAdminsHandler 管理员列表处理器（分页）
func ListAdminsHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListAdminsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateAdminHandler 创建管理员处理器
func CreateAdminHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateAdminHandler 更新管理员处理器
func UpdateAdminHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminUpdateRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteAdminHandler 删除管理员处理器
func DeleteAdminHandler(svc services.IAdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.AdminListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

