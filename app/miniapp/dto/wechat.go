package dto

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Code          string `json:"code" binding:"required"` // 微信登录code
	Nickname      string `json:"nickname"`                // 用户昵称（明文）
	Avatar        string `json:"avatar"`                  // 用户头像（明文）
	Gender        int    `json:"gender"`                  // 性别（明文）
	Country       string `json:"country"`                 // 国家（明文）
	Province      string `json:"province"`                // 省份（明文）
	City          string `json:"city"`                    // 城市（明文）
	EncryptedData string `json:"encrypted_data"`          // 加密的用户信息（备用）
	IV            string `json:"iv"`                      // 加密算法初始向量（备用）
	RawData       string `json:"raw_data"`                // 原始数据字符串（备用）
	Signature     string `json:"signature"`               // 签名（备用）
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	NeedSelectTenant bool             `json:"need_select_tenant"`       // 是否需要选择租户
	NeedBindTenant   bool             `json:"need_bind_tenant"`         // 是否需要绑定租户
	Token            string           `json:"token,omitempty"`          // JWT Token
	UserInfo         *WechatUserInfo  `json:"user_info"`                // 用户信息
	Tenants          []UserTenantInfo `json:"tenants,omitempty"`        // 用户关联的租户列表
	CurrentTenant    *UserTenantInfo  `json:"current_tenant,omitempty"` // 当前租户信息
}

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	ID       string `json:"id"`
	UnionID  string `json:"union_id"`
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone,omitempty"`
	Gender   int    `json:"gender"`
	Country  string `json:"country,omitempty"`
	Province string `json:"province,omitempty"`
	City     string `json:"city,omitempty"`
}

// UserTenantInfo 用户租户信息
type UserTenantInfo struct {
	TenantID   string   `json:"tenant_id"`
	TenantCode string   `json:"tenant_code"`
	TenantName string   `json:"tenant_name"`
	MemberID   string   `json:"member_id"`            // 成员ID
	Status     string   `json:"status"`               // active/inactive
	Roles      []string `json:"roles,omitempty"`      // 在该租户的角色
	JoinedAt   int64    `json:"joined_at"`            // 加入时间
	LeftAt     int64    `json:"left_at,omitempty"`    // 离职时间
	JobNumber  string   `json:"job_number,omitempty"` // 工号
	Department string   `json:"department,omitempty"` // 部门
	Position   string   `json:"position,omitempty"`   // 岗位
}

// BindTenantRequest 绑定租户请求
type BindTenantRequest struct {
	UserID     string `json:"user_id" binding:"required"`     // 用户ID
	InviteCode string `json:"invite_code" binding:"required"` // 租户邀请码
}

// BindTenantResponse 绑定租户响应
type BindTenantResponse struct {
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Token      string          `json:"token,omitempty"`       // 新的JWT Token
	TenantInfo *UserTenantInfo `json:"tenant_info,omitempty"` // 租户信息
}

// SelectTenantRequest 选择租户请求
type SelectTenantRequest struct {
	UserID   string `json:"user_id" binding:"required"`   // 用户ID
	TenantID string `json:"tenant_id" binding:"required"` // 租户ID
}

// SelectTenantResponse 选择租户响应
type SelectTenantResponse struct {
	Token         string          `json:"token"`          // 新的JWT Token
	UserInfo      *WechatUserInfo `json:"user_info"`      // 用户信息
	CurrentTenant *UserTenantInfo `json:"current_tenant"` // 当前租户信息
}

// SwitchTenantRequest 切换租户请求
type SwitchTenantRequest struct {
	TenantID string `json:"tenant_id" binding:"required"` // 要切换到的租户ID
}

// SwitchTenantResponse 切换租户响应
type SwitchTenantResponse struct {
	Token         string          `json:"token"`          // 新的JWT Token
	UserInfo      *WechatUserInfo `json:"user_info"`      // 用户信息
	CurrentTenant *UserTenantInfo `json:"current_tenant"` // 当前租户信息
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserID string `json:"user_id"` // 用户ID（从JWT中获取）
}

// GetUserInfoResponse 获取用户信息响应
type GetUserInfoResponse struct {
	UserInfo      *WechatUserInfo  `json:"user_info"`                // 用户信息
	Tenants       []UserTenantInfo `json:"tenants"`                  // 用户所有租户
	CurrentTenant *UserTenantInfo  `json:"current_tenant,omitempty"` // 当前租户
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone"`
	Gender     int    `json:"gender"`     // 性别：0-未知 1-男 2-女
	JobNumber  string `json:"job_number"` // 工号
	Department string `json:"department"` // 部门
	Position   string `json:"position"`   // 岗位
}

// UpdateUserInfoResponse 更新用户信息响应
type UpdateUserInfoResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetPhoneNumberRequest 获取手机号请求
type GetPhoneNumberRequest struct {
	Code string `json:"code" binding:"required"` // 微信返回的code
}

// GetPhoneNumberResponse 获取手机号响应
type GetPhoneNumberResponse struct {
	PhoneNumber string `json:"phone_number"` // 手机号
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}
