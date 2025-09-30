package dto

import "mule-cloud/models"

// ColorRequest 颜色请求
type ColorListRequest struct {
	ID    string `uri:"id"`
	Value string `query:"value"`

	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
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
	Colors []*models.Basic `json:"colors"`
	Total  int             `json:"total"`
}
