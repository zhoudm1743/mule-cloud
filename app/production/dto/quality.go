package dto

// InspectionRequest 质检提交请求
type InspectionRequest struct {
	OrderID        string   `json:"order_id" binding:"required"`
	BatchID        string   `json:"batch_id"`
	BundleNo       string   `json:"bundle_no"`
	ProcedureSeq   int      `json:"procedure_seq" binding:"required"`
	ProcedureName  string   `json:"procedure_name" binding:"required"`
	InspectedQty   int      `json:"inspected_qty" binding:"required,gt=0"`
	QualifiedQty   int      `json:"qualified_qty" binding:"required,gte=0"`
	UnqualifiedQty int      `json:"unqualified_qty" binding:"required,gte=0"`
	DefectTypes    []string `json:"defect_types"`
	DefectDesc     string   `json:"defect_desc"`
	Color          string   `json:"color"`
	Size           string   `json:"size"`
	Images         []string `json:"images"`
	Remark         string   `json:"remark"`
}

// InspectionResponse 质检提交响应
type InspectionResponse struct {
	InspectionID string  `json:"inspection_id"`
	QualityRate  float64 `json:"quality_rate"`
	NeedRework   bool    `json:"need_rework"`
	Message      string  `json:"message"`
}

// InspectionListRequest 质检列表请求
type InspectionListRequest struct {
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"page_size" form:"page_size"`
	InspectorID string `json:"inspector_id" form:"inspector_id"`
	ContractNo  string `json:"contract_no" form:"contract_no"`
	StartDate   int64  `json:"start_date" form:"start_date"`
	EndDate     int64  `json:"end_date" form:"end_date"`
}

// InspectionListResponse 质检列表响应
type InspectionListResponse struct {
	Inspections []*InspectionItem    `json:"inspections"`
	Total       int64                `json:"total"`
	Statistics  *InspectionStatistics `json:"statistics"`
}

// InspectionItem 质检记录项
type InspectionItem struct {
	ID              string   `json:"id"`
	OrderID         string   `json:"order_id"`
	ContractNo      string   `json:"contract_no"`
	StyleNo         string   `json:"style_no"`
	StyleName       string   `json:"style_name"`
	BundleNo        string   `json:"bundle_no"`
	Color           string   `json:"color"`
	Size            string   `json:"size"`
	ProcedureName   string   `json:"procedure_name"`
	InspectedQty    int      `json:"inspected_qty"`
	QualifiedQty    int      `json:"qualified_qty"`
	UnqualifiedQty  int      `json:"unqualified_qty"`
	QualityRate     float64  `json:"quality_rate"`
	DefectTypes     []string `json:"defect_types"`
	InspectorName   string   `json:"inspector_name"`
	InspectionTime  int64    `json:"inspection_time"`
	NeedRework      bool     `json:"need_rework"`
	ReworkID        string   `json:"rework_id"`
}

// InspectionStatistics 质检统计
type InspectionStatistics struct {
	TotalInspected   int     `json:"total_inspected"`
	TotalQualified   int     `json:"total_qualified"`
	TotalUnqualified int     `json:"total_unqualified"`
	QualityRate      float64 `json:"quality_rate"`
}

