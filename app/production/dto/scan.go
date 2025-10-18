package dto

import "mule-cloud/internal/models"

// ScanCodeRequest 扫码解析请求
type ScanCodeRequest struct {
	QRCode string `json:"qr_code" binding:"required"`
}

// ScanCodeResponse 扫码解析响应
type ScanCodeResponse struct {
	Batch         *BatchInfo                       `json:"batch"`
	Order         *OrderInfo                       `json:"order"`
	BatchProgress []*models.BatchProcedureProgress `json:"batch_progress"`
}

// BatchInfo 批次信息
type BatchInfo struct {
	ID         string `json:"id"`
	BundleNo   string `json:"bundle_no"`
	Color      string `json:"color"`
	Size       string `json:"size"`
	Quantity   int    `json:"quantity"`
	BedNo      string `json:"bed_no"`
	ContractNo string `json:"contract_no"`
	StyleNo    string `json:"style_no"`
	TaskID     string `json:"task_id"`
	OrderID    string `json:"order_id"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	ID           string                     `json:"id"`
	ContractNo   string                     `json:"contract_no"`
	StyleNo      string                     `json:"style_no"`
	StyleName    string                     `json:"style_name"`
	StyleImage   string                     `json:"style_image"`
	CustomerName string                     `json:"customer_name"`
	Procedures   []models.OrderProcedure    `json:"procedures"`
}

