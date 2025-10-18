package dto

// ReworkRequest 创建返工单请求
type ReworkRequest struct {
	InspectionID        string `json:"inspection_id"`
	OrderID             string `json:"order_id" binding:"required"`
	BatchID             string `json:"batch_id"`
	BundleNo            string `json:"bundle_no"`
	SourceProcedure     int    `json:"source_procedure" binding:"required"`
	TargetProcedure     int    `json:"target_procedure" binding:"required"`
	ReworkQty           int    `json:"rework_qty" binding:"required,gt=0"`
	ReworkReason        string `json:"rework_reason" binding:"required"`
	AssignedWorker      string `json:"assigned_worker"`
	Color               string `json:"color"`
	Size                string `json:"size"`
}

// ReworkResponse 创建返工单响应
type ReworkResponse struct {
	ReworkID string `json:"rework_id"`
	Message  string `json:"message"`
}

// ReworkListRequest 返工列表请求
type ReworkListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	Status     *int   `json:"status" form:"status"`
	WorkerID   string `json:"worker_id" form:"worker_id"`
	ContractNo string `json:"contract_no" form:"contract_no"`
}

// ReworkListResponse 返工列表响应
type ReworkListResponse struct {
	Reworks    []*ReworkItem      `json:"reworks"`
	Total      int64              `json:"total"`
	Statistics *ReworkStatistics  `json:"statistics"`
}

// ReworkItem 返工记录项
type ReworkItem struct {
	ID                  string `json:"id"`
	OrderID             string `json:"order_id"`
	ContractNo          string `json:"contract_no"`
	StyleNo             string `json:"style_no"`
	StyleName           string `json:"style_name"`
	BundleNo            string `json:"bundle_no"`
	Color               string `json:"color"`
	Size                string `json:"size"`
	SourceProcedureName string `json:"source_procedure_name"`
	TargetProcedureName string `json:"target_procedure_name"`
	ReworkQty           int    `json:"rework_qty"`
	ReworkReason        string `json:"rework_reason"`
	Status              int    `json:"status"`
	StatusText          string `json:"status_text"`
	CreatedByName       string `json:"created_by_name"`
	AssignedWorkerName  string `json:"assigned_worker_name"`
	CreatedAt           int64  `json:"created_at"`
	CompletedAt         int64  `json:"completed_at"`
}

// CompleteReworkRequest 完成返工请求
type CompleteReworkRequest struct {
	Images []string `json:"images"`
	Remark string   `json:"remark"`
}

// ReworkStatistics 返工统计
type ReworkStatistics struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
}

