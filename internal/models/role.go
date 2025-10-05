package models

import "go.mongodb.org/mongo-driver/v2/bson"

// Role 角色模型
type Role struct {
	ID              string              `json:"id" bson:"_id,omitempty"`
	Name            string              `json:"name" bson:"name"`                                             // 角色名称
	Code            string              `json:"code" bson:"code"`                                             // 角色代码（唯一标识）
	Description     string              `json:"description" bson:"description"`                               // 角色描述
	Menus           []string            `json:"menus" bson:"menus"`                                           // 菜单名称数组（menu.name）
	MenuPermissions map[string][]string `json:"menu_permissions,omitempty" bson:"menu_permissions,omitempty"` // 菜单权限映射: {"admin": ["read", "create", "update"], "role": ["read"]}
	Status          int                 `json:"status" bson:"status"`                                         // 状态：1-启用 0-禁用
	IsDeleted       int                 `json:"is_deleted" bson:"is_deleted"`                                 // 是否删除：0-否 1-是
	CreatedBy       string              `json:"created_by" bson:"created_by"`                                 // 创建人
	UpdatedBy       string              `json:"updated_by" bson:"updated_by"`                                 // 更新人
	CreatedAt       int64               `json:"created_at" bson:"created_at"`                                 // 创建时间
	UpdatedAt       int64               `json:"updated_at" bson:"updated_at"`                                 // 更新时间
	DeletedAt       int64               `json:"deleted_at" bson:"deleted_at,omitempty"`                       // 删除时间
}

// TableName 返回表名
func (Role) TableName() string {
	return "role"
}

// BeforeCreate 创建前钩子
func (r *Role) BeforeCreate() {
	if r.ID == "" {
		r.ID = bson.NewObjectID().Hex()
	}
	if r.Menus == nil {
		r.Menus = []string{}
	}
}
