package transport

import (
	"mule-cloud/app/perms/dto"
	"mule-cloud/app/perms/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// CreateRoleHandler 创建角色
func CreateRoleHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 租户权限检查
		perm := NewPermissionChecker(c)
		// 数据库隔离后不需要租户权限检查

		// TODO: 从上下文获取创建人（待集成JWT后）
		createdBy := perm.UserID
		if createdBy == "" {
			createdBy = "system"
		}

		role, err := roleSvc.Create(c.Request.Context(), &req, createdBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "创建成功", role)
	}
}

// GetRoleHandler 获取角色详情
func GetRoleHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		role, err := roleSvc.GetByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, role)
	}
}

// ListRolesHandler 查询角色列表
func ListRolesHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ListRoleRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		roles, total, err := roleSvc.List(c.Request.Context(), &req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, map[string]interface{}{
			"roles": roles,
			"total": total,
			"page":  req.Page,
			"size":  req.PageSize,
		})
	}
}

// UpdateRoleHandler 更新角色
func UpdateRoleHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.UpdateRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 数据库隔离后不需要租户权限检查

		// TODO: 从上下文获取更新人（待集成JWT后）
		updatedBy := "system"

		err := roleSvc.Update(c.Request.Context(), id, &req, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "更新成功", nil)
	}
}

// DeleteRoleHandler 删除角色
func DeleteRoleHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		// 数据库隔离后不需要租户权限检查

		err := roleSvc.Delete(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "删除成功", nil)
	}
}

// BatchDeleteRolesHandler 批量删除角色
func BatchDeleteRolesHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BatchDeleteRoleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		err := roleSvc.BatchDelete(c.Request.Context(), req.IDs)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "批量删除成功", nil)
	}
}

// AssignMenusHandler 分配菜单权限（支持细粒度权限）
func AssignMenusHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		var req dto.AssignMenusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取更新人（待集成JWT后）
		updatedBy := "system"

		err := roleSvc.AssignMenus(c.Request.Context(), roleID, req.Menus, req.MenuPermissions, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "分配菜单权限成功", nil)
	}
}

// GetRoleMenusHandler 获取角色的菜单权限
func GetRoleMenusHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID := c.Param("id")

		menus, err := roleSvc.GetRoleMenus(c.Request.Context(), roleID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, menus)
	}
}

// GetTenantRolesHandler 获取租户下的所有角色
func GetTenantRolesHandler(roleSvc *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Query("tenant_id")
		if tenantID == "" {
			response.Error(c, "租户ID不能为空")
			return
		}

		roles, err := roleSvc.GetRolesByTenant(c.Request.Context(), tenantID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, roles)
	}
}
