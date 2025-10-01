package models

// Menu 菜单/路由模型 - 适配Nova-admin前端路由结构（扁平结构）
type Menu struct {
	// ========== Nova-admin 路由字段（扁平结构，所有字段在同一层级）==========
	ID            string   `json:"id" bson:"_id,omitempty"`
	PID           *string  `bson:"pid" json:"pid"`                                   // 父级路由ID，顶级为null
	Name          string   `bson:"name" json:"name"`                                 // 路由名称(唯一标识)
	Path          string   `bson:"path" json:"path"`                                 // 路由路径
	Title         string   `bson:"title" json:"title"`                               // 页面标题
	ComponentPath *string  `bson:"componentPath" json:"componentPath"`               // 组件路径（目录为null）
	Redirect      string   `bson:"redirect,omitempty" json:"redirect,omitempty"`     // 重定向
	Icon          string   `bson:"icon,omitempty" json:"icon,omitempty"`             // 图标
	RequiresAuth  bool     `bson:"requiresAuth" json:"requiresAuth"`                 // 是否需要认证
	Roles         []string `bson:"roles,omitempty" json:"roles,omitempty"`           // 角色权限
	KeepAlive     bool     `bson:"keepAlive,omitempty" json:"keepAlive,omitempty"`   // 页面缓存
	Hide          bool     `bson:"hide,omitempty" json:"hide,omitempty"`             // 是否隐藏
	Order         int      `bson:"order,omitempty" json:"order,omitempty"`           // 排序
	Href          string   `bson:"href,omitempty" json:"href,omitempty"`             // 外链
	ActiveMenu    string   `bson:"activeMenu,omitempty" json:"activeMenu,omitempty"` // 高亮菜单
	WithoutTab    bool     `bson:"withoutTab,omitempty" json:"withoutTab,omitempty"` // 不加入Tab
	PinTab        bool     `bson:"pinTab,omitempty" json:"pinTab,omitempty"`         // 固定Tab
	MenuType      string   `bson:"menuType" json:"menuType"`                         // 菜单类型: dir/page

	// ========== 通用字段 ==========
	Status    int    `json:"status" bson:"status"`         // 状态: 1-启用 0-禁用
	IsDeleted int    `json:"is_deleted" bson:"is_deleted"` // 是否删除: 0-否 1-是
	CreatedBy string `json:"created_by" bson:"created_by"` // 创建人
	UpdatedBy string `json:"updated_by" bson:"updated_by"` // 更新人
	CreatedAt int64  `json:"created_at" bson:"created_at"` // 创建时间
	UpdatedAt int64  `json:"updated_at" bson:"updated_at"` // 更新时间
	DeletedAt int64  `json:"deleted_at" bson:"deleted_at"` // 删除时间
}

// TableName 集合名称
func (Menu) TableName() string {
	return "menus"
}
