package dto

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string   `json:"token"`
	UserID    string   `json:"user_id"`
	Phone     string   `json:"phone"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Role      []string `json:"role"`
	ExpiresAt int64    `json:"expires_at"`
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

// GetProfileResponse 获取个人信息响应
type GetProfileResponse struct {
	UserID    string   `json:"user_id"`
	Phone     string   `json:"phone"`
	Nickname  string   `json:"nickname"`
	Avatar    string   `json:"avatar"`
	Email     string   `json:"email"`
	Role      []string `json:"role"`
	Status    int      `json:"status"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
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
