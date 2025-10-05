package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Password   string `json:"password" binding:"required"`
	TenantCode string `json:"tenant_code"` // 租户代码（可选，为空则查询系统库）
	IP         string `json:"ip"`          // 登录IP
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token           string              `json:"token"`
	UserID          string              `json:"user_id"`
	TenantID        string              `json:"tenant_id"` // 租户ID
	Phone           string              `json:"phone"`
	Nickname        string              `json:"nickname"`
	Avatar          string              `json:"avatar"`
	Role            []string            `json:"role"`
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty"` // 用户的菜单权限映射：{"admin": ["read", "create"], "finance": ["read"]}
	ExpiresAt       int64               `json:"expires_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID   string `json:"user_id"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

// RefreshTokenRequest 刷新Token请求
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// RefreshTokenResponse 刷新Token响应
type RefreshTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// TenantItem 租户列表项
type TenantItem struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

// GetTenantListResponse 获取租户列表响应
type GetTenantListResponse struct {
	Tenants []TenantItem `json:"tenants"`
	Total   int          `json:"total"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePasswordResponse 修改密码响应
type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetUserRoutesResponse 获取用户路由响应
type GetUserRoutesResponse struct {
	Routes []RouteItem `json:"routes"`
}

// RouteItem 路由项
type RouteItem struct {
	ID            string   `json:"id"`
	PID           *string  `json:"pid"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Title         string   `json:"title"`
	ComponentPath *string  `json:"componentPath"`
	Redirect      string   `json:"redirect,omitempty"`
	Icon          string   `json:"icon,omitempty"`
	RequiresAuth  bool     `json:"requiresAuth"`
	Roles         []string `json:"roles,omitempty"`
	KeepAlive     bool     `json:"keepAlive,omitempty"`
	Hide          bool     `json:"hide,omitempty"`
	Order         int      `json:"order,omitempty"`
	Href          string   `json:"href,omitempty"`
	ActiveMenu    string   `json:"activeMenu,omitempty"`
	WithoutTab    bool     `json:"withoutTab,omitempty"`
	PinTab        bool     `json:"pinTab,omitempty"`
	MenuType      string   `json:"menuType"`
	Status        int      `json:"status,omitempty"`
	CreatedAt     int64    `json:"created_at,omitempty"`
	UpdatedAt     int64    `json:"updated_at,omitempty"`
}

// GetProfileResponse 获取个人信息响应
type GetProfileResponse struct {
	UserID          string              `json:"user_id"`
	TenantID        string              `json:"tenant_id"` // 租户ID
	Phone           string              `json:"phone"`
	Nickname        string              `json:"nickname"`
	Avatar          string              `json:"avatar"`
	Email           string              `json:"email"`
	Role            []string            `json:"role"`
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty"` // 用户的菜单权限映射
	Status          int                 `json:"status"`
	CreatedAt       int64               `json:"created_at"`
	UpdatedAt       int64               `json:"updated_at"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// UpdateProfileResponse 更新个人信息响应
type UpdateProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
