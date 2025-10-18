package models

// ProcedureReport 工序上报记录
type ProcedureReport struct {
	ID            string  `json:"id" bson:"_id,omitempty"`
	OrderID       string  `json:"order_id" bson:"order_id"`             // 订单ID
	ContractNo    string  `json:"contract_no" bson:"contract_no"`       // 合同号
	StyleNo       string  `json:"style_no" bson:"style_no"`             // 款号
	StyleName     string  `json:"style_name" bson:"style_name"`         // 款名
	BatchID       string  `json:"batch_id" bson:"batch_id"`             // 批次ID（可选，不分扎工序不需要）
	BundleNo      string  `json:"bundle_no" bson:"bundle_no"`           // 扎号（可选）
	Color         string  `json:"color" bson:"color"`                   // 颜色
	Size          string  `json:"size" bson:"size"`                     // 尺码（可选）
	Quantity      int     `json:"quantity" bson:"quantity"`             // 上报数量
	ProcedureSeq  int     `json:"procedure_seq" bson:"procedure_seq"`   // 工序序号
	ProcedureName string  `json:"procedure_name" bson:"procedure_name"` // 工序名称
	UnitPrice     float64 `json:"unit_price" bson:"unit_price"`         // 工价
	TotalPrice    float64 `json:"total_price" bson:"total_price"`       // 总工资 = 数量 * 工价
	WorkerID      string  `json:"worker_id" bson:"worker_id"`           // 工人ID（从登录信息获取）
	WorkerName    string  `json:"worker_name" bson:"worker_name"`       // 工人姓名
	WorkerNo      string  `json:"worker_no" bson:"worker_no"`           // 工号
	ReportTime    int64   `json:"report_time" bson:"report_time"`       // 上报时间
	Remark        string  `json:"remark" bson:"remark"`                 // 备注
	IsDeleted     int     `json:"is_deleted" bson:"is_deleted"`         // 是否删除：0-否 1-是
	CreatedAt     int64   `json:"created_at" bson:"created_at"`         // 创建时间
	UpdatedAt     int64   `json:"updated_at" bson:"updated_at"`         // 更新时间
}

// TableName 返回表名
func (ProcedureReport) TableName() string {
	return "procedure_reports"
}

// BatchProcedureProgress 批次工序进度
type BatchProcedureProgress struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	BatchID       string `json:"batch_id" bson:"batch_id"`             // 批次ID
	BundleNo      string `json:"bundle_no" bson:"bundle_no"`           // 扎号
	OrderID       string `json:"order_id" bson:"order_id"`             // 订单ID
	ProcedureSeq  int    `json:"procedure_seq" bson:"procedure_seq"`   // 工序序号
	ProcedureName string `json:"procedure_name" bson:"procedure_name"` // 工序名称
	Quantity      int    `json:"quantity" bson:"quantity"`             // 批次总数量
	ReportedQty   int    `json:"reported_qty" bson:"reported_qty"`     // 已上报数量
	IsCompleted   bool   `json:"is_completed" bson:"is_completed"`     // 是否完成
	CompletedAt   int64  `json:"completed_at" bson:"completed_at"`     // 完成时间
	CreatedAt     int64  `json:"created_at" bson:"created_at"`         // 创建时间
	UpdatedAt     int64  `json:"updated_at" bson:"updated_at"`         // 更新时间
}

// TableName 返回表名
func (BatchProcedureProgress) TableName() string {
	return "batch_procedure_progress"
}

// OrderProcedureProgress 订单工序进度汇总
type OrderProcedureProgress struct {
	ID            string  `json:"id" bson:"_id,omitempty"`
	OrderID       string  `json:"order_id" bson:"order_id"`             // 订单ID
	ContractNo    string  `json:"contract_no" bson:"contract_no"`       // 合同号
	ProcedureSeq  int     `json:"procedure_seq" bson:"procedure_seq"`   // 工序序号
	ProcedureName string  `json:"procedure_name" bson:"procedure_name"` // 工序名称
	TotalQty      int     `json:"total_qty" bson:"total_qty"`           // 订单总数量
	ReportedQty   int     `json:"reported_qty" bson:"reported_qty"`     // 已上报数量
	Progress      float64 `json:"progress" bson:"progress"`             // 进度百分比
	CreatedAt     int64   `json:"created_at" bson:"created_at"`         // 创建时间
	UpdatedAt     int64   `json:"updated_at" bson:"updated_at"`         // 更新时间
}

// TableName 返回表名
func (OrderProcedureProgress) TableName() string {
	return "order_procedure_progress"
}

