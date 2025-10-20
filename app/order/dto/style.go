package dto

import "mule-cloud/internal/models"

// StyleListRequest 款式列表请求
type StyleListRequest struct {
	ID        string `uri:"id" form:"id"`
	StyleNo   string `form:"style_no"`   // 款号
	StyleName string `form:"style_name"` // 款名
	Category  string `form:"category"`   // 分类
	Season    string `form:"season"`     // 季节
	Year      string `form:"year"`       // 年份
	Status    int    `form:"status"`     // 状态

	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

// StyleCreateRequest 创建款式请求
type StyleCreateRequest struct {
	StyleNo    string                  `json:"style_no" binding"required"`   // 款号
	StyleName  string                  `json:"style_name" binding"required"` // 款名
	Category   string                  `json:"category"`                     // 分类
	Season     string                  `json:"season"`                       // 季节
	Year       string                  `json:"year"`                         // 年份
	Images     []string                `json:"images"`                       // 图片URL列表
	Colors     []string                `json:"colors"`                       // 颜色列表
	Sizes      []string                `json:"sizes"`                        // 尺码列表
	UnitPrice  float64                 `json:"unit_price"`                   // 单价
	Remark     string                  `json:"remark"`                       // 备注
	Procedures []models.StyleProcedure `json:"procedures"`                   // 工序清单
	Status     int                     `json:"status"`                       // 状态
}

// StyleUpdateRequest 更新款式请求
type StyleUpdateRequest struct {
	ID         string                  `uri:"id" binding"required"`
	StyleName  string                  `json:"style_name"` // 款名
	Category   string                  `json:"category"`   // 分类
	Season     string                  `json:"season"`     // 季节
	Year       string                  `json:"year"`       // 年份
	Images     []string                `json:"images"`     // 图片URL列表
	Colors     []string                `json:"colors"`     // 颜色列表
	Sizes      []string                `json:"sizes"`      // 尺码列表
	UnitPrice  float64                 `json:"unit_price"` // 单价
	Remark     string                  `json:"remark"`     // 备注
	Procedures []models.StyleProcedure `json:"procedures"` // 工序清单
	Status     int                     `json:"status"`     // 状态
}

// StyleResponse 款式响应
type StyleResponse struct {
	Style *models.Style `json:"style"`
}

// StyleListResponse 款式列表响应
type StyleListResponse struct {
	Styles []models.Style `json:"styles"`
	Total  int64          `json:"total"`
}
