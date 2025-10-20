package dto

import "mule-cloud/internal/models"

// ProcedureRequest 工序请求
type ProcedureListRequest struct {
	ID    string `uri:"id" form:"id"`
	Value string `form:"value"`

	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type ProcedureCreateRequest struct {
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type ProcedureUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// ProcedureResponse 工序响应
type ProcedureResponse struct {
	Procedure *models.Basic `json:"procedure"`
}

// ProcedureListResponse 工序列表响应
type ProcedureListResponse struct {
	Procedures []models.Basic `json:"procedures"`
	Total      int64          `json:"total"`
}
