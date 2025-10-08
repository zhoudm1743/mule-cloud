package dto

// CreateDepartmentRequest 创建部门请求
type CreateDepartmentRequest struct {
	Name     string `json:"name" binding:"required"`     // 部门名称
	Code     string `json:"code" binding:"required"`     // 部门编码
	ParentID string `json:"parent_id"`                   // 父部门ID
	Status   int    `json:"status"`                      // 状态
}

// UpdateDepartmentRequest 更新部门请求
type UpdateDepartmentRequest struct {
	Name     string `json:"name"`      // 部门名称
	Code     string `json:"code"`      // 部门编码
	ParentID string `json:"parent_id"` // 父部门ID
	Status   *int   `json:"status"`    // 状态
}

// ListDepartmentRequest 查询部门列表请求
type ListDepartmentRequest struct {
	Name     string `json:"name" form:"name"`           // 部门名称（模糊查询）
	Code     string `json:"code" form:"code"`           // 部门编码（模糊查询）
	ParentID string `json:"parent_id" form:"parent_id"` // 父部门ID
	Status   *int   `json:"status" form:"status"`       // 状态
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
}

// BatchDeleteDepartmentRequest 批量删除部门请求
type BatchDeleteDepartmentRequest struct {
	IDs []string `json:"ids" binding:"required"` // 部门ID数组
}

