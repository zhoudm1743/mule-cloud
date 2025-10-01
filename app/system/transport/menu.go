package transport

import (
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/services"
	"mule-cloud/core/response"
	"mule-cloud/internal/models"

	"github.com/gin-gonic/gin"
)

// GetAllMenusHandler 获取所有菜单（Nova-admin路由数据）
func GetAllMenusHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := menuSvc.GetAll(c.Request.Context())
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, menus)
	}
}

// GetMenuHandler 获取单个菜单
func GetMenuHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		menu, err := menuSvc.GetByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, menu)
	}
}

// CreateMenuHandler 创建菜单
func CreateMenuHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateMenuRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		menu := &models.Menu{
			PID:           req.PID,
			Name:          req.Name,
			Path:          req.Path,
			Title:         req.Title,
			ComponentPath: req.ComponentPath,
			Redirect:      req.Redirect,
			Icon:          req.Icon,
			RequiresAuth:  req.RequiresAuth,
			Roles:         req.Roles,
			KeepAlive:     req.KeepAlive,
			Hide:          req.Hide,
			Order:         req.Order,
			Href:          req.Href,
			ActiveMenu:    req.ActiveMenu,
			WithoutTab:    req.WithoutTab,
			PinTab:        req.PinTab,
			MenuType:      req.MenuType,
		}

		err := menuSvc.Create(c.Request.Context(), menu)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, menu)
	}
}

// UpdateMenuHandler 更新菜单
func UpdateMenuHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.UpdateMenuRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 构建更新字段
		updates := make(map[string]interface{})
		if req.PID != nil {
			updates["pid"] = req.PID
		}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Path != "" {
			updates["path"] = req.Path
		}
		if req.Title != "" {
			updates["title"] = req.Title
		}
		if req.ComponentPath != nil {
			updates["componentPath"] = req.ComponentPath
		}
		if req.Redirect != "" {
			updates["redirect"] = req.Redirect
		}
		if req.Icon != "" {
			updates["icon"] = req.Icon
		}
		if req.RequiresAuth != nil {
			updates["requiresAuth"] = *req.RequiresAuth
		}
		if req.Roles != nil {
			updates["roles"] = req.Roles
		}
		if req.KeepAlive != nil {
			updates["keepAlive"] = *req.KeepAlive
		}
		if req.Hide != nil {
			updates["hide"] = *req.Hide
		}
		if req.Order != nil {
			updates["order"] = *req.Order
		}
		if req.Href != "" {
			updates["href"] = req.Href
		}
		if req.ActiveMenu != "" {
			updates["activeMenu"] = req.ActiveMenu
		}
		if req.WithoutTab != nil {
			updates["withoutTab"] = *req.WithoutTab
		}
		if req.PinTab != nil {
			updates["pinTab"] = *req.PinTab
		}
		if req.MenuType != "" {
			updates["menuType"] = req.MenuType
		}

		err := menuSvc.Update(c.Request.Context(), id, updates)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "更新成功", nil)
	}
}

// DeleteMenuHandler 删除菜单
func DeleteMenuHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := menuSvc.Delete(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "删除成功", nil)
	}
}

// ListMenusHandler 分页查询菜单
func ListMenusHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ListMenuRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 默认值
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}

		// 构建过滤条件
		filters := make(map[string]interface{})
		if req.Name != "" {
			filters["name"] = req.Name
		}
		if req.Title != "" {
			filters["title"] = req.Title
		}
		if req.MenuType != "" {
			filters["menuType"] = req.MenuType
		}
		if req.Status != nil {
			filters["status"] = *req.Status
		}

		menus, total, err := menuSvc.List(c.Request.Context(), req.Page, req.PageSize, filters)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"list":     menus,
			"total":    total,
			"page":     req.Page,
			"pageSize": req.PageSize,
		})
	}
}

// BatchDeleteMenusHandler 批量删除菜单
func BatchDeleteMenusHandler(menuSvc *services.MenuService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BatchDeleteMenuRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		err := menuSvc.BatchDelete(c.Request.Context(), req.IDs)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "批量删除成功", nil)
	}
}
