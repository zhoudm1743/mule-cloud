package transport

import (
	"mule-cloud/app/perms/dto"
	"mule-cloud/app/perms/endpoint"
	"mule-cloud/app/perms/services"
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

		// 数据库隔离后不需要租户权限检查

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

		// 数据库隔离后不需要租户权限检查

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

		// 数据库隔离后不需要租户权限检查

		ep := endpoint.DeleteAdminEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// AssignAdminRolesHandler 分配角色给管理员
func AssignAdminRolesHandler(adminSvc *services.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Param("id")

		var req dto.AssignRolesRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取操作人（待集成JWT后）
		updatedBy := "system"

		err := adminSvc.AssignRoles(c.Request.Context(), adminID, req.Roles, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "分配角色成功", nil)
	}
}

// GetAdminRolesHandler 获取管理员的角色
func GetAdminRolesHandler(adminSvc *services.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Param("id")

		roles, err := adminSvc.GetAdminRoles(c.Request.Context(), adminID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, roles)
	}
}

// RemoveAdminRoleHandler 移除管理员的某个角色
func RemoveAdminRoleHandler(adminSvc *services.AdminService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Param("id")
		roleID := c.Param("roleId")

		// TODO: 从上下文获取操作人（待集成JWT后）
		updatedBy := "system"

		err := adminSvc.RemoveRole(c.Request.Context(), adminID, roleID, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "移除角色成功", nil)
	}
}
