package dto

// CreatePostRequest 创建岗位请求
type CreatePostRequest struct {
	Name     string `json:"name" binding:"required"` // 岗位名称
	Code     string `json:"code" binding:"required"` // 岗位编码
	ParentID string `json:"parent_id"`               // 父岗位ID
	Status   int    `json:"status"`                  // 状态
}

// UpdatePostRequest 更新岗位请求
type UpdatePostRequest struct {
	Name     string `json:"name"`      // 岗位名称
	Code     string `json:"code"`      // 岗位编码
	ParentID string `json:"parent_id"` // 父岗位ID
	Status   *int   `json:"status"`    // 状态
}

// ListPostRequest 查询岗位列表请求
type ListPostRequest struct {
	Name     string `json:"name" form:"name"`           // 岗位名称（模糊查询）
	Code     string `json:"code" form:"code"`           // 岗位编码（模糊查询）
	ParentID string `json:"parent_id" form:"parent_id"` // 父岗位ID
	Status   *int   `json:"status" form:"status"`       // 状态
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
}

// BatchDeletePostRequest 批量删除岗位请求
type BatchDeletePostRequest struct {
	IDs []string `json:"ids" binding:"required"` // 岗位ID数组
}

