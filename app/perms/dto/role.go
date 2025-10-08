package dto

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	TenantID        string              `json:"tenant_id" binding"required"` // 租户ID
	Name            string              `json:"name" binding"required"`      // 角色名称
	Code            string              `json:"code" binding"required"`      // 角色代码
	Description     string              `json:"description"`                 // 角色描述
	Menus           []string            `json:"menus"`                       // 菜单名称数组（menu.name）
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty"`  // 菜单权限映射: {"admin": ["read", "create"], "role": ["read"]}
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	Name            string              `json:"name"`                       // 角色名称
	Description     string              `json:"description"`                // 角色描述
	Menus           []string            `json:"menus"`                      // 菜单名称数组（menu.name）
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty"` // 菜单权限映射
	Status          *int                `json:"status"`                     // 状态
}

// ListRoleRequest 查询角色列表请求
type ListRoleRequest struct {
	TenantID string `json:"tenant_id" form:"tenant_id"` // 租户ID（可选，用于筛选）
	Name     string `json:"name" form:"name"`           // 角色名称（模糊查询）
	Code     string `json:"code" form:"code"`           // 角色代码（模糊查询）
	Status   *int   `json:"status" form:"status"`       // 状态
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
}

// BatchDeleteRoleRequest 批量删除角色请求
type BatchDeleteRoleRequest struct {
	IDs []string `json:"ids" binding"required"` // 角色ID数组
}

// AssignMenusRequest 分配菜单权限请求
type AssignMenusRequest struct {
	Menus           []string            `json:"menus" binding"required"`    // 菜单名称数组（menu.name）
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty"` // 菜单权限映射: {"admin": ["read", "create"], "role": ["read"]}
}

// AssignRolesRequest 分配角色请求
type AssignRolesRequest struct {
	Roles []string `json:"roles" binding"required"` // 角色ID数组
}
