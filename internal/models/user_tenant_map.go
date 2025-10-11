package models

// UserTenantMap 用户-租户关系映射表（存储在系统库）
type UserTenantMap struct {
	ID         string `json:"id" bson:"_id,omitempty"`        // MongoDB ObjectID
	UserID     string `json:"user_id" bson:"user_id"`         // 微信用户ID（WechatUser._id）
	UnionID    string `json:"union_id" bson:"union_id"`       // 冗余字段，便于查询
	TenantID   string `json:"tenant_id" bson:"tenant_id"`     // 租户ID
	TenantCode string `json:"tenant_code" bson:"tenant_code"` // 租户代码
	MemberID   string `json:"member_id" bson:"member_id"`     // 在租户库中的成员ID

	// 关系状态
	Status   string `json:"status" bson:"status"`         // 状态：active-在职 inactive-离职 pending-待审核
	JoinedAt int64  `json:"joined_at" bson:"joined_at"`   // 加入时间
	LeftAt   int64  `json:"left_at" bson:"left_at"`       // 离职时间（如果已离职）

	// 系统字段
	IsDeleted int   `json:"is_deleted" bson:"is_deleted"` // 软删除
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}

// TableName 返回表名
func (UserTenantMap) TableName() string {
	return "user_tenant_map"
}

