package models

// Tenant 租户模型
type Tenant struct {
	ID        string   `json:"id" bson:"_id,omitempty"`
	Code      string   `json:"code" bson:"code"`             // 租户代码
	Name      string   `json:"name" bson:"name"`             // 租户名称
	Contact   string   `json:"contact" bson:"contact"`       // 联系人
	Phone     string   `json:"phone" bson:"phone"`           // 联系电话
	Email     string   `json:"email" bson:"email"`           // 联系邮箱
	Menus     []string `json:"menus" bson:"menus"`           // 租户拥有的菜单权限（由超管分配）
	Status    int      `json:"status" bson:"status"`         // 状态：1-启用 0-禁用
	IsDeleted int      `json:"is_deleted" bson:"is_deleted"` // 是否删除：0-否 1-是
	CreatedBy string   `json:"created_by" bson:"created_by"` // 创建人
	UpdatedBy string   `json:"updated_by" bson:"updated_by"` // 更新人
	CreatedAt int64    `json:"created_at" bson:"created_at"` // 创建时间
	UpdatedAt int64    `json:"updated_at" bson:"updated_at"` // 更新时间
	DeletedAt int64    `json:"deleted_at" bson:"deleted_at"` // 删除时间
}

// TableName 返回表名
func (Tenant) TableName() string {
	return "tenant"
}

// HasMenu 检查租户是否拥有指定菜单权限
func (t *Tenant) HasMenu(menuID string) bool {
	for _, m := range t.Menus {
		if m == menuID {
			return true
		}
	}
	return false
}

// HasAllMenus 检查租户是否拥有所有指定的菜单权限
func (t *Tenant) HasAllMenus(menuIDs []string) bool {
	for _, menuID := range menuIDs {
		if !t.HasMenu(menuID) {
			return false
		}
	}
	return true
}
