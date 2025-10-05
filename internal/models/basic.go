package models

// Basic 基础数据模型
type Basic struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Name      string `json:"name" bson:"name"`             // 名称
	Value     string `json:"value" bson:"value"`           // 值
	Remark    string `json:"remark" bson:"remark"`         // 备注
	IsCommon  bool   `json:"is_common" bson:"is_common"`   // 是否公共数据（保留字段，用于标记共享数据）
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
