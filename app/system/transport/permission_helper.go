package transport

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// PermissionChecker 权限检查助手
type PermissionChecker struct {
	UserID   string
	TenantID string
	Roles    []interface{}
}

// NewPermissionChecker 从 Gin Context 创建权限检查器
func NewPermissionChecker(c *gin.Context) *PermissionChecker {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	roles, _ := c.Get("roles")

	userIDStr, _ := userID.(string)
	tenantIDStr, _ := tenantID.(string)
	rolesList, _ := roles.([]interface{})

	return &PermissionChecker{
		UserID:   userIDStr,
		TenantID: tenantIDStr,
		Roles:    rolesList,
	}
}

// IsSystemAdmin 是否为系统超管
func (p *PermissionChecker) IsSystemAdmin() bool {
	if p.Roles == nil {
		return false
	}
	for _, role := range p.Roles {
		if roleStr, ok := role.(string); ok && roleStr == "super" {
			// 系统超管：角色包含 "super" 且没有租户ID
			return p.TenantID == ""
		}
	}
	return false
}

// IsTenantAdmin 是否为租户超管
func (p *PermissionChecker) IsTenantAdmin() bool {
	if p.Roles == nil {
		return false
	}
	for _, role := range p.Roles {
		if roleStr, ok := role.(string); ok && roleStr == "tenant_admin" {
			// 租户超管：角色包含 "tenant_admin" 且有租户ID
			return p.TenantID != ""
		}
	}
	return false
}

// CanAccessTenant 是否可以访问指定租户的数据
func (p *PermissionChecker) CanAccessTenant(targetTenantID string) bool {
	// 系统超管：可以访问所有租户
	if p.IsSystemAdmin() {
		return true
	}

	// 租户超管：只能访问自己的租户
	if p.IsTenantAdmin() {
		return p.TenantID == targetTenantID
	}

	// 普通用户：只能访问自己的租户
	return p.TenantID == targetTenantID
}

// CheckTenantPermission 检查租户权限（用于创建、更新、删除操作）
func (p *PermissionChecker) CheckTenantPermission(targetTenantID string, operation string) error {
	// 系统超管：允许所有操作
	if p.IsSystemAdmin() {
		return nil
	}

	// 空租户ID的资源只有系统超管可以操作
	if targetTenantID == "" {
		return fmt.Errorf("无权限：只有系统超管可以%s系统级资源", operation)
	}

	// 租户超管和普通用户：只能操作自己租户的数据
	if p.TenantID != targetTenantID {
		return fmt.Errorf("无权限：不能%s其他租户的数据", operation)
	}

	return nil
}

// MustHavePermission 必须有权限，否则 panic（由 Gin 的 recovery 中间件捕获）
func (p *PermissionChecker) MustHavePermission(targetTenantID string, operation string) {
	if err := p.CheckTenantPermission(targetTenantID, operation); err != nil {
		panic(err)
	}
}

