package models

// TenantMember 租户成员（存储在各租户库中）
// 这是对原有 Admin 模型的扩展，用于支持小程序用户
type TenantMember struct {
	ID string `json:"id" bson:"_id,omitempty"` // MongoDB ObjectID

	// 关联全局用户
	UnionID string `json:"union_id" bson:"union_id"` // 微信UnionID（关联全局用户）
	UserID  string `json:"user_id" bson:"user_id"`   // 全局用户ID（WechatUser._id）

	// 租户内信息（可被租户管理员修改）
	Name       string `json:"name" bson:"name"`               // 员工姓名（在本租户的称呼）
	JobNumber  string `json:"job_number" bson:"job_number"`   // 工号
	Department string `json:"department" bson:"department"`   // 部门
	Position   string `json:"position" bson:"position"`       // 岗位
	Phone      string `json:"phone" bson:"phone"`             // 联系电话（冗余）
	Avatar     string `json:"avatar" bson:"avatar"`           // 头像

	// 权限与角色
	Roles       []string `json:"role" bson:"role"`             // 角色ID数组（注意：字段名为role单数）
	Permissions []string `json:"permissions" bson:"permissions"` // 额外权限

	// 状态
	Status     string `json:"status" bson:"status"`             // 状态：active-在职 inactive-离职
	EmployedAt int64  `json:"employed_at" bson:"employed_at"`   // 入职时间
	LeftAt     int64  `json:"left_at" bson:"left_at,omitempty"` // 离职时间

	// 系统字段
	IsDeleted int    `json:"is_deleted" bson:"is_deleted"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdatedBy string `json:"updated_by" bson:"updated_by"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
}

// TableName 返回表名
func (TenantMember) TableName() string {
	return "member"
}

