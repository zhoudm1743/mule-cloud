package dto

import "mule-cloud/internal/models"

// OrderTypeRequest 订单类型请求
type OrderTypeListRequest struct {
	ID    string `uri:"id" query:"id"`
	Value string `query:"value"`

	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
}

type OrderTypeCreateRequest struct {
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type OrderTypeUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// OrderTypeResponse 订单类型响应
type OrderTypeResponse struct {
	OrderType *models.Basic `json:"order_type"`
}

// OrderTypeListResponse 订单类型列表响应
type OrderTypeListResponse struct {
	OrderTypes []models.Basic `json:"order_types"`
	Total      int64          `json:"total"`
}
