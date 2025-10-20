package dto

import "mule-cloud/internal/models"

// SizeRequest 尺寸请求
type SizeGetRequest struct {
	ID string `uri:"id" binding"required"`
}

type SizeListRequest struct {
	Value    string `form:"value"`
	ID       string `uri:"id" form:"id"`
	Page     int64  `form:"page"`
	PageSize int64  `form:"page_size"`
}

type SizeCreateRequest struct {
	Value  string `json:"value" binding"required"`
	Remark string `json:"remark"`
}

type SizeUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value" binding"required"`
	Remark string `json:"remark"`
}

// SizeResponse 尺寸响应
type SizeResponse struct {
	Size *models.Basic `json:"size"`
}

// SizeListResponse 尺寸列表响应
type SizeListResponse struct {
	Sizes []*models.Basic `json:"sizes"`
	Total int64           `json:"total"`
}
