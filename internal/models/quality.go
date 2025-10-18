package models

// QualityInspection 质检记录
type QualityInspection struct {
	ID              string   `json:"id" bson:"_id,omitempty"`
	OrderID         string   `json:"order_id" bson:"order_id"`                   // 订单ID
	ContractNo      string   `json:"contract_no" bson:"contract_no"`             // 合同号
	StyleNo         string   `json:"style_no" bson:"style_no"`                   // 款号
	StyleName       string   `json:"style_name" bson:"style_name"`               // 款名
	BatchID         string   `json:"batch_id" bson:"batch_id"`                   // 批次ID（可选）
	BundleNo        string   `json:"bundle_no" bson:"bundle_no"`                 // 扎号（可选）
	Color           string   `json:"color" bson:"color"`                         // 颜色
	Size            string   `json:"size" bson:"size"`                           // 尺码
	ProcedureSeq    int      `json:"procedure_seq" bson:"procedure_seq"`         // 质检的工序序号
	ProcedureName   string   `json:"procedure_name" bson:"procedure_name"`       // 工序名称
	InspectedQty    int      `json:"inspected_qty" bson:"inspected_qty"`         // 质检数量
	QualifiedQty    int      `json:"qualified_qty" bson:"qualified_qty"`         // 合格数量
	UnqualifiedQty  int      `json:"unqualified_qty" bson:"unqualified_qty"`     // 不合格数量
	QualityRate     float64  `json:"quality_rate" bson:"quality_rate"`           // 合格率（%）
	DefectTypes     []string `json:"defect_types" bson:"defect_types"`           // 缺陷类型
	DefectDesc      string   `json:"defect_desc" bson:"defect_desc"`             // 缺陷描述
	InspectorID     string   `json:"inspector_id" bson:"inspector_id"`           // 质检员ID
	InspectorName   string   `json:"inspector_name" bson:"inspector_name"`       // 质检员姓名
	InspectionTime  int64    `json:"inspection_time" bson:"inspection_time"`     // 质检时间
	Images          []string `json:"images" bson:"images"`                       // 质检照片
	Remark          string   `json:"remark" bson:"remark"`                       // 备注
	NeedRework      bool     `json:"need_rework" bson:"need_rework"`             // 是否需要返工
	ReworkID        string   `json:"rework_id" bson:"rework_id"`                 // 返工单ID
	IsDeleted       int      `json:"is_deleted" bson:"is_deleted"`               // 是否删除
	CreatedAt       int64    `json:"created_at" bson:"created_at"`               // 创建时间
	UpdatedAt       int64    `json:"updated_at" bson:"updated_at"`               // 更新时间
}

// TableName 返回表名
func (QualityInspection) TableName() string {
	return "quality_inspections"
}

// ReworkRecord 返工记录
type ReworkRecord struct {
	ID                 string   `json:"id" bson:"_id,omitempty"`
	OrderID            string   `json:"order_id" bson:"order_id"`                         // 订单ID
	ContractNo         string   `json:"contract_no" bson:"contract_no"`                   // 合同号
	StyleNo            string   `json:"style_no" bson:"style_no"`                         // 款号
	StyleName          string   `json:"style_name" bson:"style_name"`                     // 款名
	BatchID            string   `json:"batch_id" bson:"batch_id"`                         // 批次ID
	BundleNo           string   `json:"bundle_no" bson:"bundle_no"`                       // 扎号
	Color              string   `json:"color" bson:"color"`                               // 颜色
	Size               string   `json:"size" bson:"size"`                                 // 尺码
	InspectionID       string   `json:"inspection_id" bson:"inspection_id"`               // 关联的质检记录ID
	SourceProcedure    int      `json:"source_procedure" bson:"source_procedure"`         // 来源工序（哪个工序出的问题）
	SourceProcedureName string  `json:"source_procedure_name" bson:"source_procedure_name"` // 来源工序名称
	TargetProcedure    int      `json:"target_procedure" bson:"target_procedure"`         // 目标工序（返回到哪个工序）
	TargetProcedureName string  `json:"target_procedure_name" bson:"target_procedure_name"` // 目标工序名称
	ReworkQty          int      `json:"rework_qty" bson:"rework_qty"`                     // 返工数量
	ReworkReason       string   `json:"rework_reason" bson:"rework_reason"`               // 返工原因
	Status             int      `json:"status" bson:"status"`                             // 状态：0-待返工 1-返工中 2-已完成
	CreatedBy          string   `json:"created_by" bson:"created_by"`                     // 创建人（质检员）
	CreatedByName      string   `json:"created_by_name" bson:"created_by_name"`           // 创建人姓名
	AssignedWorker     string   `json:"assigned_worker" bson:"assigned_worker"`           // 返工工人
	AssignedWorkerName string   `json:"assigned_worker_name" bson:"assigned_worker_name"` // 返工工人姓名
	CompletedAt        int64    `json:"completed_at" bson:"completed_at"`                 // 完成时间
	Images             []string `json:"images" bson:"images"`                             // 返工照片
	Remark             string   `json:"remark" bson:"remark"`                             // 备注
	IsDeleted          int      `json:"is_deleted" bson:"is_deleted"`                     // 是否删除
	CreatedAt          int64    `json:"created_at" bson:"created_at"`                     // 创建时间
	UpdatedAt          int64    `json:"updated_at" bson:"updated_at"`                     // 更新时间
}

// TableName 返回表名
func (ReworkRecord) TableName() string {
	return "rework_records"
}

