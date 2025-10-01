package dto

// CreateMenuRequest 创建菜单请求
type CreateMenuRequest struct {
	PID           *string  `json:"pid"`                                        // 父级ID
	Name          string   `json:"name" binding:"required"`                    // 路由名称
	Path          string   `json:"path" binding:"required"`                    // 路由路径
	Title         string   `json:"title" binding:"required"`                   // 标题
	ComponentPath *string  `json:"componentPath"`                              // 组件路径
	Redirect      string   `json:"redirect"`                                   // 重定向
	Icon          string   `json:"icon"`                                       // 图标
	RequiresAuth  bool     `json:"requiresAuth"`                               // 需要认证
	Roles         []string `json:"roles"`                                      // 角色
	KeepAlive     bool     `json:"keepAlive"`                                  // 缓存
	Hide          bool     `json:"hide"`                                       // 隐藏
	Order         int      `json:"order"`                                      // 排序
	Href          string   `json:"href"`                                       // 外链
	ActiveMenu    string   `json:"activeMenu"`                                 // 高亮菜单
	WithoutTab    bool     `json:"withoutTab"`                                 // 不加Tab
	PinTab        bool     `json:"pinTab"`                                     // 固定Tab
	MenuType      string   `json:"menuType" binding:"required,oneof=dir page"` // 类型: dir/page
}

// UpdateMenuRequest 更新菜单请求
type UpdateMenuRequest struct {
	PID           *string  `json:"pid"`
	Name          string   `json:"name"`
	Path          string   `json:"path"`
	Title         string   `json:"title"`
	ComponentPath *string  `json:"componentPath"`
	Redirect      string   `json:"redirect"`
	Icon          string   `json:"icon"`
	RequiresAuth  *bool    `json:"requiresAuth"`
	Roles         []string `json:"roles"`
	KeepAlive     *bool    `json:"keepAlive"`
	Hide          *bool    `json:"hide"`
	Order         *int     `json:"order"`
	Href          string   `json:"href"`
	ActiveMenu    string   `json:"activeMenu"`
	WithoutTab    *bool    `json:"withoutTab"`
	PinTab        *bool    `json:"pinTab"`
	MenuType      string   `json:"menuType"`
}

// ListMenuRequest 查询菜单列表请求
type ListMenuRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name"`
	Title    string `form:"title"`
	MenuType string `form:"menuType"`
	Status   *int   `form:"status"`
}

// BatchDeleteMenuRequest 批量删除请求
type BatchDeleteMenuRequest struct {
	IDs []string `json:"ids" binding:"required,min=1"`
}
