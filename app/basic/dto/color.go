package dto

import "mule-cloud/internal/models"

// ColorRequest 颜色请求
type ColorListRequest struct {
	ID    string `uri:"id" form:"id"`
	Value string `form:"value"`

	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

type ColorCreateRequest struct {
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

type ColorUpdateRequest struct {
	ID     string `uri:"id"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// ColorResponse 颜色响应
type ColorResponse struct {
	Color *models.Basic `json:"color"`
}

// ColorListResponse 颜色列表响应
type ColorListResponse struct {
	Colors []models.Basic `json:"colors"`
	Total  int64          `json:"total"`
}
