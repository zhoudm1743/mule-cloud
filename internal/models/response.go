package models

import "time"

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta 分页元数据
type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ListResponse 列表响应
type ListResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// StatsResponse 统计响应
type StatsResponse struct {
	TotalOrders       int64   `json:"total_orders"`
	ActiveOrders      int64   `json:"active_orders"`
	TotalWorkers      int64   `json:"total_workers"`
	ActiveWorkers     int64   `json:"active_workers"`
	TodayProduction   int64   `json:"today_production"`
	WeeklyProduction  int64   `json:"weekly_production"`
	MonthlyProduction int64   `json:"monthly_production"`
	AvgCompletionRate float64 `json:"avg_completion_rate"`
}

// OrderStatsResponse 订单统计响应
type OrderStatsResponse struct {
	TotalOrders      int64   `json:"total_orders"`
	PendingOrders    int64   `json:"pending_orders"`
	InProgressOrders int64   `json:"in_progress_orders"`
	CompletedOrders  int64   `json:"completed_orders"`
	CancelledOrders  int64   `json:"cancelled_orders"`
	TotalAmount      float64 `json:"total_amount"`
	AvgOrderValue    float64 `json:"avg_order_value"`
}

// ProductionStatsResponse 生产统计响应
type ProductionStatsResponse struct {
	TotalProduction   int64   `json:"total_production"`
	TodayProduction   int64   `json:"today_production"`
	WeeklyProduction  int64   `json:"weekly_production"`
	MonthlyProduction int64   `json:"monthly_production"`
	CompletionRate    float64 `json:"completion_rate"`
	DefectRate        float64 `json:"defect_rate"`
	AvgDailyOutput    float64 `json:"avg_daily_output"`
}

// WorkerStatsResponse 工人统计响应
type WorkerStatsResponse struct {
	TotalWorkers    int64   `json:"total_workers"`
	ActiveWorkers   int64   `json:"active_workers"`
	TodayAttendance int64   `json:"today_attendance"`
	AvgProductivity float64 `json:"avg_productivity"`
	TotalWorkHours  float64 `json:"total_work_hours"`
	AvgWorkHours    float64 `json:"avg_work_hours"`
}

// PayrollStatsResponse 工资统计响应
type PayrollStatsResponse struct {
	TotalPayroll    float64 `json:"total_payroll"`
	MonthlyPayroll  float64 `json:"monthly_payroll"`
	AvgWorkerSalary float64 `json:"avg_worker_salary"`
	OvertimePayroll float64 `json:"overtime_payroll"`
	BonusPayroll    float64 `json:"bonus_payroll"`
}

// DashboardResponse 仪表板响应
type DashboardResponse struct {
	Orders           OrderStatsResponse      `json:"orders"`
	Production       ProductionStatsResponse `json:"production"`
	Workers          WorkerStatsResponse     `json:"workers"`
	Payroll          PayrollStatsResponse    `json:"payroll"`
	RecentActivities []ActivityResponse      `json:"recent_activities"`
	Alerts           []AlertResponse         `json:"alerts"`
}

// ActivityResponse 活动响应
type ActivityResponse struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	UserName    string    `json:"user_name"`
	Timestamp   time.Time `json:"timestamp"`
}

// AlertResponse 告警响应
type AlertResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Level        string    `json:"level"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	Source       string    `json:"source"`
	Timestamp    time.Time `json:"timestamp"`
	Acknowledged bool      `json:"acknowledged"`
}

// OrderProgressResponse 订单进度响应
type OrderProgressResponse struct {
	OrderID        string                    `json:"order_id"`
	OrderNo        string                    `json:"order_no"`
	CustomerName   string                    `json:"customer_name"`
	StyleName      string                    `json:"style_name"`
	TotalQty       int                       `json:"total_qty"`
	CompletedQty   int                       `json:"completed_qty"`
	CompletionRate float64                   `json:"completion_rate"`
	Status         string                    `json:"status"`
	StartDate      time.Time                 `json:"start_date"`
	DeliveryDate   time.Time                 `json:"delivery_date"`
	Processes      []ProcessProgressResponse `json:"processes"`
}

// ProcessProgressResponse 工序进度响应
type ProcessProgressResponse struct {
	ProcessID      string     `json:"process_id"`
	ProcessName    string     `json:"process_name"`
	PlannedQty     int        `json:"planned_qty"`
	CompletedQty   int        `json:"completed_qty"`
	DefectQty      int        `json:"defect_qty"`
	CompletionRate float64    `json:"completion_rate"`
	Status         string     `json:"status"`
	PlannedStart   time.Time  `json:"planned_start"`
	PlannedEnd     time.Time  `json:"planned_end"`
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
}

// WorkerProductivityResponse 工人生产力响应
type WorkerProductivityResponse struct {
	WorkerID       string  `json:"worker_id"`
	WorkerName     string  `json:"worker_name"`
	TotalHours     float64 `json:"total_hours"`
	TotalOutput    int     `json:"total_output"`
	Productivity   float64 `json:"productivity"`
	QualityRate    float64 `json:"quality_rate"`
	AttendanceRate float64 `json:"attendance_rate"`
	Ranking        int     `json:"ranking"`
}

// PayrollSummaryResponse 工资汇总响应
type PayrollSummaryResponse struct {
	WorkerID       string  `json:"worker_id"`
	WorkerName     string  `json:"worker_name"`
	BaseSalary     float64 `json:"base_salary"`
	PieceSalary    float64 `json:"piece_salary"`
	OvertimeSalary float64 `json:"overtime_salary"`
	Bonus          float64 `json:"bonus"`
	Deduction      float64 `json:"deduction"`
	TotalSalary    float64 `json:"total_salary"`
	PayPeriod      string  `json:"pay_period"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}

// ValidationError 验证错误
type ValidationError struct {
	Field   string      `json:"field"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

// ValidationErrorResponse 验证错误响应
type ValidationErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

// PaginationRequest 分页请求
type PaginationRequest struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

// DateRangeRequest 日期范围请求
type DateRangeRequest struct {
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
}

// SortRequest 排序请求
type SortRequest struct {
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order,default=desc" binding:"oneof=asc desc"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Keyword string `form:"keyword"`
	Fields  string `form:"fields"`
}

// FilterRequest 过滤请求
type FilterRequest struct {
	Status    *int    `form:"status"`
	Category  string  `form:"category"`
	Type      string  `form:"type"`
	MinAmount float64 `form:"min_amount"`
	MaxAmount float64 `form:"max_amount"`
}

// ExportRequest 导出请求
type ExportRequest struct {
	Format    string    `form:"format,default=excel" binding:"oneof=excel csv pdf"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
	Fields    []string  `form:"fields"`
}

// ImportRequest 导入请求
type ImportRequest struct {
	FileURL   string `json:"file_url" binding:"required"`
	SheetName string `json:"sheet_name"`
	SkipRows  int    `json:"skip_rows"`
}

// BulkOperationRequest 批量操作请求
type BulkOperationRequest struct {
	IDs       []string    `json:"ids" binding:"required,min=1"`
	Operation string      `json:"operation" binding:"required"`
	Data      interface{} `json:"data"`
}

// BulkOperationResponse 批量操作响应
type BulkOperationResponse struct {
	SuccessCount int      `json:"success_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
}
