package dto

import "mule-cloud/internal/models"

// CuttingTaskCreateRequest 创建裁剪任务请求
type CuttingTaskCreateRequest struct {
	OrderID   string `json:"order_id" binding:"required"`
	CreatedBy string `json:"created_by"`
}

// CuttingTaskListRequest 裁剪任务列表请求
type CuttingTaskListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	ContractNo string `json:"contract_no" form:"contract_no"`
	StyleNo    string `json:"style_no" form:"style_no"`
	Status     *int   `json:"status" form:"status"`
}

// CuttingTaskListResponse 裁剪任务列表响应
type CuttingTaskListResponse struct {
	Tasks []*models.CuttingTask `json:"tasks"`
	Total int64                 `json:"total"`
}

// CuttingTaskResponse 裁剪任务响应
type CuttingTaskResponse struct {
	Task *models.CuttingTask `json:"task"`
}

// CuttingBatchCreateRequest 创建裁剪批次请求
type CuttingBatchCreateRequest struct {
	TaskID      string                `json:"task_id" binding:"required"`
	BedNo       string                `json:"bed_no" binding:"required"`
	BundleNo    string                `json:"bundle_no" binding:"required"`
	Color       string                `json:"color" binding:"required"`
	LayerCount  int                   `json:"layer_count" binding:"required,gt=0"`
	SizeDetails []models.SizeDetail   `json:"size_details" binding:"required,min=1"`
	CreatedBy   string                `json:"created_by"`
}

// CuttingBatchListRequest 裁剪批次列表请求
type CuttingBatchListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	TaskID     string `json:"task_id" form:"task_id"`
	ContractNo string `json:"contract_no" form:"contract_no"`
	BedNo      string `json:"bed_no" form:"bed_no"`
	BundleNo   string `json:"bundle_no" form:"bundle_no"`
}

// CuttingBatchListResponse 裁剪批次列表响应
type CuttingBatchListResponse struct {
	Batches []*models.CuttingBatch `json:"batches"`
	Total   int64                  `json:"total"`
}

// CuttingBatchResponse 裁剪批次响应
type CuttingBatchResponse struct {
	Batch *models.CuttingBatch `json:"batch"`
}

// CuttingPieceListRequest 裁片监控列表请求
type CuttingPieceListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	OrderID    string `json:"order_id" form:"order_id"`
	ContractNo string `json:"contract_no" form:"contract_no"`
	BedNo      string `json:"bed_no" form:"bed_no"`
	BundleNo   string `json:"bundle_no" form:"bundle_no"`
}

// CuttingPieceListResponse 裁片监控列表响应
type CuttingPieceListResponse struct {
	Pieces []*models.CuttingPiece `json:"pieces"`
	Total  int64                  `json:"total"`
}

// CuttingPieceResponse 裁片监控响应
type CuttingPieceResponse struct {
	Piece *models.CuttingPiece `json:"piece"`
}

// CuttingPieceProgressRequest 更新裁片进度请求
type CuttingPieceProgressRequest struct {
	Progress int `json:"progress" binding:"required,gte=0"`
}
