package models

// Admin 管理员模型
type Admin struct {
	ID        string   `json:"id" bson:"_id,omitempty"`
	Phone     string   `json:"phone" bson:"phone"`           // 手机号
	Password  string   `json:"-" bson:"password"`            // 密码（不返回给前端）
	Nickname  string   `json:"nickname" bson:"nickname"`     // 昵称
	Email     string   `json:"email" bson:"email"`           // 邮箱
	Avatar    string   `json:"avatar" bson:"avatar"`         // 头像
	Roles     []string `json:"role" bson:"role"`             // 角色ID数组（注意：数据库字段名为 role 单数）
	Status    int      `json:"status" bson:"status"`         // 状态：1-启用 0-禁用
	IsDeleted int      `json:"is_deleted" bson:"is_deleted"` // 是否删除：0-否 1-是
	Extend    Extend   `json:"extend" bson:"extend"`         // 扩展字段
	CreatedBy string   `json:"created_by" bson:"created_by"` // 创建人
	UpdatedBy string   `json:"updated_by" bson:"updated_by"` // 更新人
	CreatedAt int64    `json:"created_at" bson:"created_at"` // 创建时间
	UpdatedAt int64    `json:"updated_at" bson:"updated_at"` // 更新时间
	DeletedAt int64    `json:"deleted_at" bson:"deleted_at"` // 删除时间
}

type Extend struct {
	LastLoginAt int64  `json:"last_login_at" bson:"last_login_at"` // 最后登录时间
	LoginCount  int    `json:"login_count" bson:"login_count"`     // 登录次数
	LastLoginIP string `json:"last_login_ip" bson:"last_login_ip"` // 最后登录IP
}

// TableName 返回表名
func (Admin) TableName() string {
	return "admin"
}
