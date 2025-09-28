package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username" binding:"required,min=3,max=50"`
	Email       string             `bson:"email" json:"email" binding:"required,email"`
	Password    string             `bson:"password" json:"-"`
	Phone       string             `bson:"phone,omitempty" json:"phone"`
	RealName    string             `bson:"real_name,omitempty" json:"real_name"`
	Avatar      string             `bson:"avatar,omitempty" json:"avatar"`
	Status      UserStatus         `bson:"status" json:"status"`
	LastLoginAt *time.Time         `bson:"last_login_at,omitempty" json:"last_login_at"`
	RoleIDs     []string           `bson:"role_ids,omitempty" json:"role_ids"`
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"created_by"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	Version     int                `bson:"version" json:"version"`
}

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusInactive UserStatus = 0 // 未激活
	UserStatusActive   UserStatus = 1 // 激活
	UserStatusBlocked  UserStatus = 2 // 被阻止
	UserStatusDeleted  UserStatus = 3 // 已删除
)

// Role 角色模型
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" binding:"required,min=2,max=50"`
	Code        string             `bson:"code" json:"code" binding:"required,min=2,max=50"`
	Description string             `bson:"description,omitempty" json:"description"`
	Permissions []string           `bson:"permissions,omitempty" json:"permissions"`
	Status      RoleStatus         `bson:"status" json:"status"`
	CreatedBy   primitive.ObjectID `bson:"created_by,omitempty" json:"created_by"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// RoleStatus 角色状态
type RoleStatus int

const (
	RoleStatusInactive RoleStatus = 0 // 未激活
	RoleStatusActive   RoleStatus = 1 // 激活
)

// Permission 权限模型
type Permission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Code        string             `bson:"code" json:"code" binding:"required"`
	Module      string             `bson:"module" json:"module"`
	Action      string             `bson:"action" json:"action"`
	Resource    string             `bson:"resource" json:"resource"`
	Description string             `bson:"description,omitempty" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Worker 工人模型
type Worker struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	WorkerNo   string             `bson:"worker_no" json:"worker_no" binding:"required"`
	Name       string             `bson:"name" json:"name" binding:"required"`
	Phone      string             `bson:"phone,omitempty" json:"phone"`
	IDCard     string             `bson:"id_card,omitempty" json:"id_card"`
	Department string             `bson:"department,omitempty" json:"department"`
	Position   string             `bson:"position,omitempty" json:"position"`
	SkillLevel SkillLevel         `bson:"skill_level" json:"skill_level"`
	HourlyRate float64            `bson:"hourly_rate" json:"hourly_rate"`
	PieceRates []PieceRate        `bson:"piece_rates,omitempty" json:"piece_rates"`
	Status     WorkerStatus       `bson:"status" json:"status"`
	HireDate   time.Time          `bson:"hire_date" json:"hire_date"`
	CreatedBy  primitive.ObjectID `bson:"created_by,omitempty" json:"created_by"`
	UpdatedBy  primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// SkillLevel 技能等级
type SkillLevel int

const (
	SkillLevelJunior       SkillLevel = 1 // 初级
	SkillLevelIntermediate SkillLevel = 2 // 中级
	SkillLevelSenior       SkillLevel = 3 // 高级
	SkillLevelExpert       SkillLevel = 4 // 专家
)

// WorkerStatus 工人状态
type WorkerStatus int

const (
	WorkerStatusActive    WorkerStatus = 1 // 在职
	WorkerStatusInactive  WorkerStatus = 0 // 离职
	WorkerStatusSuspended WorkerStatus = 2 // 停职
)

// PieceRate 计件工价
type PieceRate struct {
	ProcessID   primitive.ObjectID `bson:"process_id" json:"process_id"`
	ProcessName string             `bson:"process_name" json:"process_name"`
	Rate        float64            `bson:"rate" json:"rate"`
	Unit        string             `bson:"unit" json:"unit"`
	EffectiveAt time.Time          `bson:"effective_at" json:"effective_at"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int64     `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	User         *UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       primitive.ObjectID `json:"id"`
	Username string             `json:"username"`
	Email    string             `json:"email"`
	RealName string             `json:"real_name"`
	Avatar   string             `json:"avatar"`
	Roles    []string           `json:"roles"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	RealName string `json:"real_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status"`
	RoleID   string `form:"role_id"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	RealName string   `json:"real_name,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	RoleIDs  []string `json:"role_ids,omitempty"`
	Status   int      `json:"status"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Email    string   `json:"email,omitempty" binding:"omitempty,email"`
	RealName string   `json:"real_name,omitempty"`
	Phone    string   `json:"phone,omitempty"`
	RoleIDs  []string `json:"role_ids,omitempty"`
	Status   int      `json:"status"`
}
