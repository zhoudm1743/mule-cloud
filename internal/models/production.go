package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductionPlan 生产计划模型
type ProductionPlan struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PlanNo  string             `bson:"plan_no" json:"plan_no" binding:"required"`
	OrderID primitive.ObjectID `bson:"order_id" json:"order_id" binding:"required"`
	StyleID primitive.ObjectID `bson:"style_id" json:"style_id"`

	// 计划信息
	PlannedQty  int           `bson:"planned_qty" json:"planned_qty"`
	ProcessFlow []ProcessStep `bson:"process_flow" json:"process_flow"`
	StartDate   time.Time     `bson:"start_date" json:"start_date"`
	EndDate     time.Time     `bson:"end_date" json:"end_date"`

	// 状态信息
	Status   PlanStatus `bson:"status" json:"status"`
	Priority int        `bson:"priority" json:"priority" default:"1"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// PlanStatus 计划状态
type PlanStatus int

const (
	PlanStatusDraft      PlanStatus = 0 // 草稿
	PlanStatusScheduled  PlanStatus = 1 // 已排产
	PlanStatusInProgress PlanStatus = 2 // 进行中
	PlanStatusCompleted  PlanStatus = 3 // 已完成
	PlanStatusCancelled  PlanStatus = 4 // 已取消
)

// CuttingTask 裁剪任务模型
type CuttingTask struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TaskNo  string             `bson:"task_no" json:"task_no" binding:"required"`
	OrderID primitive.ObjectID `bson:"order_id" json:"order_id"`
	StyleID primitive.ObjectID `bson:"style_id" json:"style_id"`
	PlanID  primitive.ObjectID `bson:"plan_id" json:"plan_id"`

	// 任务信息
	FabricType     string             `bson:"fabric_type" json:"fabric_type"`
	ColorID        primitive.ObjectID `bson:"color_id" json:"color_id"`
	TotalLayers    int                `bson:"total_layers" json:"total_layers"`
	PiecesPerLayer int                `bson:"pieces_per_layer" json:"pieces_per_layer"`
	TotalPieces    int                `bson:"total_pieces" json:"total_pieces"`

	// 时间信息
	ScheduledDate time.Time  `bson:"scheduled_date" json:"scheduled_date"`
	StartTime     *time.Time `bson:"start_time,omitempty" json:"start_time"`
	EndTime       *time.Time `bson:"end_time,omitempty" json:"end_time"`

	// 分配信息
	AssignedTo primitive.ObjectID `bson:"assigned_to,omitempty" json:"assigned_to"`

	// 状态信息
	Status       CuttingStatus `bson:"status" json:"status"`
	QualityGrade string        `bson:"quality_grade,omitempty" json:"quality_grade"`

	// 备注
	Remark string `bson:"remark,omitempty" json:"remark"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CuttingStatus 裁剪状态
type CuttingStatus int

const (
	CuttingStatusPending    CuttingStatus = 0 // 待裁剪
	CuttingStatusInProgress CuttingStatus = 1 // 裁剪中
	CuttingStatusCompleted  CuttingStatus = 2 // 已完成
	CuttingStatusDefective  CuttingStatus = 3 // 有缺陷
)

// ProductionProgress 生产进度模型
type ProductionProgress struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id"`
	StyleID   primitive.ObjectID `bson:"style_id" json:"style_id"`
	ProcessID primitive.ObjectID `bson:"process_id" json:"process_id"`

	// 进度信息
	PlannedQty     int     `bson:"planned_qty" json:"planned_qty"`
	CompletedQty   int     `bson:"completed_qty" json:"completed_qty"`
	DefectQty      int     `bson:"defect_qty" json:"defect_qty"`
	CompletionRate float64 `bson:"completion_rate" json:"completion_rate"`

	// 时间信息
	PlannedStart time.Time  `bson:"planned_start" json:"planned_start"`
	PlannedEnd   time.Time  `bson:"planned_end" json:"planned_end"`
	ActualStart  *time.Time `bson:"actual_start,omitempty" json:"actual_start"`
	ActualEnd    *time.Time `bson:"actual_end,omitempty" json:"actual_end"`

	// 状态信息
	Status ProgressStatus `bson:"status" json:"status"`

	// 审计信息
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// ProgressStatus 进度状态
type ProgressStatus int

const (
	ProgressStatusNotStarted ProgressStatus = 0 // 未开始
	ProgressStatusInProgress ProgressStatus = 1 // 进行中
	ProgressStatusCompleted  ProgressStatus = 2 // 已完成
	ProgressStatusBlocked    ProgressStatus = 3 // 阻塞
)

// WorkReport 工作上报模型
type WorkReport struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkerID  primitive.ObjectID `bson:"worker_id" json:"worker_id" binding:"required"`
	OrderID   primitive.ObjectID `bson:"order_id" json:"order_id" binding:"required"`
	ProcessID primitive.ObjectID `bson:"process_id" json:"process_id" binding:"required"`

	// 工作信息
	Date      time.Time `bson:"date" json:"date" binding:"required"`
	StartTime time.Time `bson:"start_time" json:"start_time" binding:"required"`
	EndTime   time.Time `bson:"end_time" json:"end_time" binding:"required"`
	WorkHours float64   `bson:"work_hours" json:"work_hours"`
	Quantity  int       `bson:"quantity" json:"quantity" binding:"required,min=1"`

	// 工价信息
	UnitPrice float64 `bson:"unit_price" json:"unit_price"`
	Amount    float64 `bson:"amount" json:"amount"`

	// 质量信息
	QualityGrade string   `bson:"quality_grade" json:"quality_grade"`
	DefectQty    int      `bson:"defect_qty" json:"defect_qty"`
	DefectTypes  []string `bson:"defect_types,omitempty" json:"defect_types"`

	// 状态信息
	Status     ReportStatus       `bson:"status" json:"status"`
	ReviewedBy primitive.ObjectID `bson:"reviewed_by,omitempty" json:"reviewed_by"`
	ReviewedAt *time.Time         `bson:"reviewed_at,omitempty" json:"reviewed_at"`

	// 备注
	Remark string `bson:"remark,omitempty" json:"remark"`

	// 审计信息
	ReportedBy primitive.ObjectID `bson:"reported_by" json:"reported_by"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

// ReportStatus 上报状态
type ReportStatus int

const (
	ReportStatusPending  ReportStatus = 0 // 待审核
	ReportStatusApproved ReportStatus = 1 // 已审核
	ReportStatusRejected ReportStatus = 2 // 已驳回
)

// Process定义已移动到master_data.go中

// ProcessStatus 工序状态
type ProcessStatus int

const (
	ProcessStatusActive   ProcessStatus = 1 // 激活
	ProcessStatusInactive ProcessStatus = 0 // 停用
)

// QualityStandard 质量标准
type QualityStandard struct {
	Item        string `bson:"item" json:"item"`
	Standard    string `bson:"standard" json:"standard"`
	CheckMethod string `bson:"check_method,omitempty" json:"check_method"`
}

// Size和Color定义已移动到master_data.go中

// SizeStatus 尺码状态
type SizeStatus int

const (
	SizeStatusActive   SizeStatus = 1 // 激活
	SizeStatusInactive SizeStatus = 0 // 停用
)

// ColorStatus 颜色状态
type ColorStatus int

const (
	ColorStatusActive   ColorStatus = 1 // 激活
	ColorStatusInactive ColorStatus = 0 // 停用
)

// CreateProductionPlanRequest 创建生产计划请求
type CreateProductionPlanRequest struct {
	OrderID     string        `json:"order_id" binding:"required"`
	PlannedQty  int           `json:"planned_qty" binding:"required,min=1"`
	ProcessFlow []ProcessStep `json:"process_flow" binding:"required"`
	StartDate   time.Time     `json:"start_date" binding:"required"`
	EndDate     time.Time     `json:"end_date" binding:"required"`
	Priority    int           `json:"priority"`
}

// CreateCuttingTaskRequest 创建裁剪任务请求
type CreateCuttingTaskRequest struct {
	OrderID        string    `json:"order_id" binding:"required"`
	PlanID         string    `json:"plan_id"`
	FabricType     string    `json:"fabric_type" binding:"required"`
	ColorID        string    `json:"color_id" binding:"required"`
	TotalLayers    int       `json:"total_layers" binding:"required,min=1"`
	PiecesPerLayer int       `json:"pieces_per_layer" binding:"required,min=1"`
	ScheduledDate  time.Time `json:"scheduled_date" binding:"required"`
	AssignedTo     string    `json:"assigned_to"`
}

// CreateWorkReportRequest 创建工作上报请求
type CreateWorkReportRequest struct {
	WorkerID     string    `json:"worker_id" binding:"required"`
	OrderID      string    `json:"order_id" binding:"required"`
	ProcessID    string    `json:"process_id" binding:"required"`
	Date         time.Time `json:"date" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required,min=1"`
	QualityGrade string    `json:"quality_grade"`
	DefectQty    int       `json:"defect_qty"`
	DefectTypes  []string  `json:"defect_types"`
	Remark       string    `json:"remark"`
}

// UpdateWorkReportRequest 更新工作上报请求
type UpdateWorkReportRequest struct {
	Quantity     int      `json:"quantity,omitempty"`
	QualityGrade string   `json:"quality_grade,omitempty"`
	DefectQty    int      `json:"defect_qty,omitempty"`
	DefectTypes  []string `json:"defect_types,omitempty"`
	Remark       string   `json:"remark,omitempty"`
}

// BatchWorkReportRequest 批量工作上报请求
type BatchWorkReportRequest struct {
	WorkerID string                    `json:"worker_id" binding:"required"`
	Date     time.Time                 `json:"date" binding:"required"`
	Reports  []SingleWorkReportRequest `json:"reports" binding:"required,min=1"`
}

// SingleWorkReportRequest 单个工作上报请求
type SingleWorkReportRequest struct {
	OrderID      string    `json:"order_id" binding:"required"`
	ProcessID    string    `json:"process_id" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
	Quantity     int       `json:"quantity" binding:"required,min=1"`
	QualityGrade string    `json:"quality_grade"`
	DefectQty    int       `json:"defect_qty"`
	Remark       string    `json:"remark"`
}

// WorkReportListRequest 工作上报列表请求
type WorkReportListRequest struct {
	Page      int       `form:"page,default=1" binding:"min=1"`
	PageSize  int       `form:"page_size,default=10" binding:"min=1,max=100"`
	WorkerID  string    `form:"worker_id"`
	OrderID   string    `form:"order_id"`
	ProcessID string    `form:"process_id"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	Status    *int      `form:"status"`
}

// ProductionProgressRequest 生产进度查询请求
type ProductionProgressRequest struct {
	OrderID   string `form:"order_id"`
	StyleID   string `form:"style_id"`
	ProcessID string `form:"process_id"`
}
