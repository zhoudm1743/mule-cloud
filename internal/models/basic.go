package models

// Basic 基础数据模型
type Basic struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	TenantID  string `json:"tenant_id" bson:"tenant_id"`   // 租户ID（空表示系统公共数据）
	Name      string `json:"name" bson:"name"`             // 名称
	Value     string `json:"value" bson:"value"`           // 值
	Remark    string `json:"remark" bson:"remark"`         // 备注
	IsCommon  bool   `json:"is_common" bson:"is_common"`   // 是否公共数据（其他租户可查询但不可修改删除）
	Status    int    `json:"status" bson:"status"`         // 状态：1-启用 0-禁用
	IsDeleted int    `json:"is_deleted" bson:"is_deleted"` // 是否删除：0-否 1-是
	CreatedBy string `json:"created_by" bson:"created_by"` // 创建人
	UpdatedBy string `json:"updated_by" bson:"updated_by"` // 更新人
	CreatedAt int64  `json:"created_at" bson:"created_at"` // 创建时间
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"` // 更新时间
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"` // 删除时间
}

// TableName 返回表名
func (Basic) TableName() string {
	return "basic"
}

// IsOwnedBy 检查是否属于指定租户（包括公共数据）
func (b *Basic) IsOwnedBy(tenantID string) bool {
	return b.TenantID == tenantID || b.IsCommon
}

// CanModifyBy 检查是否可以被指定租户修改（公共数据不能修改）
func (b *Basic) CanModifyBy(tenantID string) bool {
	return b.TenantID == tenantID && !b.IsCommon
}
