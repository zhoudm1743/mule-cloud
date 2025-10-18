package dto

import "mule-cloud/internal/models"

// ProcedureReportRequest 工序上报请求
type ProcedureReportRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	BatchID       string `json:"batch_id"`                             // 可选，不分扎工序不需要
	BundleNo      string `json:"bundle_no"`                            // 可选
	ProcedureSeq  int    `json:"procedure_seq" binding:"required"`
	ProcedureName string `json:"procedure_name" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,gt=0"`
	Color         string `json:"color"`
	Size          string `json:"size"`
	Remark        string `json:"remark"`
}

// ProcedureReportResponse 工序上报响应
type ProcedureReportResponse struct {
	ReportID   string  `json:"report_id"`
	TotalPrice float64 `json:"total_price"`
	Message    string  `json:"message"`
}

// ReportListRequest 上报记录列表请求
type ReportListRequest struct {
	Page       int    `json:"page" form:"page"`
	PageSize   int    `json:"page_size" form:"page_size"`
	WorkerID   string `json:"worker_id" form:"worker_id"`
	ContractNo string `json:"contract_no" form:"contract_no"`
	StartDate  string `json:"start_date" form:"start_date"`
	EndDate    string `json:"end_date" form:"end_date"`
}

// ReportListResponse 上报记录列表响应
type ReportListResponse struct {
	Reports    []*models.ProcedureReport `json:"reports"`
	Total      int64                     `json:"total"`
	Statistics *ReportStatistics         `json:"statistics"`
}

// ReportStatistics 上报统计
type ReportStatistics struct {
	TotalQuantity int     `json:"total_quantity"`
	TotalAmount   float64 `json:"total_amount"`
}

// OrderProgressRequest 订单进度查询请求
type OrderProgressRequest struct {
	OrderID string `json:"order_id" uri:"order_id" binding:"required"`
}

// OrderProgressResponse 订单进度响应
type OrderProgressResponse struct {
	OrderID         string                             `json:"order_id"`
	ContractNo      string                             `json:"contract_no"`
	TotalQuantity   int                                `json:"total_quantity"`
	Procedures      []*models.OrderProcedureProgress   `json:"procedures"`
	OverallProgress float64                            `json:"overall_progress"`
}

// SalaryRequest 工资统计请求
type SalaryRequest struct {
	WorkerID  string `json:"worker_id" form:"worker_id"`
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
}

// SalaryResponse 工资统计响应
type SalaryResponse struct {
	WorkerID      string              `json:"worker_id"`
	WorkerName    string              `json:"worker_name"`
	WorkerNo      string              `json:"worker_no"`
	StartDate     string              `json:"start_date"`
	EndDate       string              `json:"end_date"`
	TotalQuantity int                 `json:"total_quantity"`
	TotalAmount   float64             `json:"total_amount"`
	Details       []map[string]interface{} `json:"details"`
}

