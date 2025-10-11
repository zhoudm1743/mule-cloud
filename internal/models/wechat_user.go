package models

// WechatUser 全局微信用户（存储在系统库）
type WechatUser struct {
	ID        string   `json:"id" bson:"_id,omitempty"`          // MongoDB ObjectID
	UnionID   string   `json:"union_id" bson:"union_id"`         // 微信UnionID（全局唯一）
	OpenID    string   `json:"open_id" bson:"open_id"`           // 微信OpenID（小程序唯一）
	Phone     string   `json:"phone" bson:"phone"`               // 手机号（可能为空）
	Nickname  string   `json:"nickname" bson:"nickname"`         // 微信昵称
	Avatar    string   `json:"avatar" bson:"avatar"`             // 微信头像
	Gender    int      `json:"gender" bson:"gender"`             // 性别：0-未知 1-男 2-女
	Country   string   `json:"country" bson:"country"`           // 国家
	Province  string   `json:"province" bson:"province"`         // 省份
	City      string   `json:"city" bson:"city"`                 // 城市
	Language  string   `json:"language" bson:"language"`         // 语言

	// 多租户关联（仅用于快速查询）
	TenantIDs []string `json:"tenant_ids" bson:"tenant_ids"` // 用户关联的租户ID列表

	// 系统字段
	Status      int   `json:"status" bson:"status"`                         // 状态：1-正常 0-禁用
	IsDeleted   int   `json:"is_deleted" bson:"is_deleted"`                 // 是否删除
	CreatedAt   int64 `json:"created_at" bson:"created_at"`                 // 创建时间
	UpdatedAt   int64 `json:"updated_at" bson:"updated_at"`                 // 更新时间
	LastLoginAt int64 `json:"last_login_at" bson:"last_login_at,omitempty"` // 最后登录时间
}

// TableName 返回表名
func (WechatUser) TableName() string {
	return "wechat_user"
}

