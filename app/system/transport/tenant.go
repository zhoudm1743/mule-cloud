package transport

import (
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/endpoint"
	"mule-cloud/app/system/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetTenantHandler 获取租户处理器
func GetTenantHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllTenantsHandler 获取所有租户处理器（不分页）
func GetAllTenantsHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllTenantsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListTenantsHandler 租户列表处理器（分页）
func ListTenantsHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListTenantsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateTenantHandler 创建租户处理器
func CreateTenantHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateTenantHandler 更新租户处理器
func UpdateTenantHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantUpdateRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteTenantHandler 删除租户处理器
func DeleteTenantHandler(svc services.ITenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TenantListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// AssignTenantMenusHandler 分配菜单权限给租户（超管使用）
func AssignTenantMenusHandler(tenantSvc *services.TenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		var req dto.AssignTenantMenusRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取操作人（待集成JWT后）
		updatedBy := "system"

		err := tenantSvc.AssignMenus(c.Request.Context(), tenantID, req.Menus, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "分配菜单权限成功", nil)
	}
}

// GetTenantMenusHandler 获取租户的菜单权限
func GetTenantMenusHandler(tenantSvc *services.TenantService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.Param("id")

		menus, err := tenantSvc.GetTenantMenus(c.Request.Context(), tenantID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, menus)
	}
}
