package models

// Order 订单模型
type Order struct {
	ID            string           `json:"id" bson:"_id,omitempty"`
	ContractNo    string           `json:"contract_no" bson:"contract_no"`         // 合同号
	StyleID       string           `json:"style_id" bson:"style_id"`               // 款式ID
	StyleNo       string           `json:"style_no" bson:"style_no"`               // 款号
	StyleName     string           `json:"style_name" bson:"style_name"`           // 款名
	StyleImage    string           `json:"style_image" bson:"style_image"`         // 款式图片
	CustomerID    string           `json:"customer_id" bson:"customer_id"`         // 客户ID
	CustomerName  string           `json:"customer_name" bson:"customer_name"`     // 客户名称
	SalesmanID    string           `json:"salesman_id" bson:"salesman_id"`         // 业务员ID
	SalesmanName  string           `json:"salesman_name" bson:"salesman_name"`     // 业务员名称
	OrderTypeID   string           `json:"order_type_id" bson:"order_type_id"`     // 订单类型ID
	OrderTypeName string           `json:"order_type_name" bson:"order_type_name"` // 订单类型名称
	Quantity      int              `json:"quantity" bson:"quantity"`               // 总数量
	UnitPrice     float64          `json:"unit_price" bson:"unit_price"`           // 单价
	TotalAmount   float64          `json:"total_amount" bson:"total_amount"`       // 总金额
	DeliveryDate  string           `json:"delivery_date" bson:"delivery_date"`     // 交货日期
	Progress      float64          `json:"progress" bson:"progress"`               // 进度百分比
	Status        int              `json:"status" bson:"status"`                   // 状态：0-草稿 1-已下单 2-生产中 3-已完成 4-已取消
	Remark        string           `json:"remark" bson:"remark"`                   // 备注
	Colors        []string         `json:"colors" bson:"colors"`                   // 颜色列表
	Sizes         []string         `json:"sizes" bson:"sizes"`                     // 尺码列表
	Items         []OrderItem      `json:"items" bson:"items"`                     // 订单明细（颜色+尺码组合）
	Procedures    []OrderProcedure `json:"procedures" bson:"procedures"`           // 工序清单
	IsDeleted     int              `json:"is_deleted" bson:"is_deleted"`           // 是否删除：0-否 1-是
	CreatedBy     string           `json:"created_by" bson:"created_by"`           // 创建人
	UpdatedBy     string           `json:"updated_by" bson:"updated_by"`           // 更新人
	CreatedAt     int64            `json:"created_at" bson:"created_at"`           // 创建时间
	UpdatedAt     int64            `json:"updated_at" bson:"updated_at"`           // 更新时间
	DeletedAt     int64            `json:"deleted_at" bson:"deleted_at"`           // 删除时间
}

// OrderItem 订单明细（颜色+尺码组合的数量）
type OrderItem struct {
	Color    string `json:"color" bson:"color"`       // 颜色名称
	Size     string `json:"size" bson:"size"`         // 尺码名称
	Quantity int    `json:"quantity" bson:"quantity"` // 数量
}

// OrderProcedure 订单工序
type OrderProcedure struct {
	Sequence       int     `json:"sequence" bson:"sequence"`               // 顺序
	ProcedureName  string  `json:"procedure_name" bson:"procedure_name"`   // 工序名称
	UnitPrice      float64 `json:"unit_price" bson:"unit_price"`           // 工价
	AssignedWorker string  `json:"assigned_worker" bson:"assigned_worker"` // 指定工人
	IsSlowest      bool    `json:"is_slowest" bson:"is_slowest"`           // 是否最终工序
	NoBundle       bool    `json:"no_bundle" bson:"no_bundle"`             // 不分扎上报
}

// TableName 返回表名
func (Order) TableName() string {
	return "orders"
}

// Style 款式模型
type Style struct {
	ID          string           `json:"id" bson:"_id,omitempty"`
	StyleNo     string           `json:"style_no" bson:"style_no"`       // 款号
	StyleName   string           `json:"style_name" bson:"style_name"`   // 款名
	Category    string           `json:"category" bson:"category"`       // 分类
	Season      string           `json:"season" bson:"season"`           // 季节
	Year        string           `json:"year" bson:"year"`               // 年份
	Description string           `json:"description" bson:"description"` // 描述
	Images      []string         `json:"images" bson:"images"`           // 图片URL列表
	Colors      []string         `json:"colors" bson:"colors"`           // 颜色列表
	Sizes       []string         `json:"sizes" bson:"sizes"`             // 尺码列表
	UnitPrice   float64          `json:"unit_price" bson:"unit_price"`   // 单价
	Remark      string           `json:"remark" bson:"remark"`           // 备注
	Procedures  []StyleProcedure `json:"procedures" bson:"procedures"`   // 工序清单
	Status      int              `json:"status" bson:"status"`           // 状态：1-启用 0-禁用
	IsDeleted   int              `json:"is_deleted" bson:"is_deleted"`   // 是否删除：0-否 1-是
	CreatedBy   string           `json:"created_by" bson:"created_by"`   // 创建人
	UpdatedBy   string           `json:"updated_by" bson:"updated_by"`   // 更新人
	CreatedAt   int64            `json:"created_at" bson:"created_at"`   // 创建时间
	UpdatedAt   int64            `json:"updated_at" bson:"updated_at"`   // 更新时间
	DeletedAt   int64            `json:"deleted_at" bson:"deleted_at"`   // 删除时间
}

// StyleProcedure 款式工序
type StyleProcedure struct {
	Sequence       int     `json:"sequence" bson:"sequence"`               // 顺序
	ProcedureName  string  `json:"procedure_name" bson:"procedure_name"`   // 工序名称
	UnitPrice      float64 `json:"unit_price" bson:"unit_price"`           // 工价
	AssignedWorker string  `json:"assigned_worker" bson:"assigned_worker"` // 指定工人
	IsSlowest      bool    `json:"is_slowest" bson:"is_slowest"`           // 是否最终工序
	NoBundle       bool    `json:"no_bundle" bson:"no_bundle"`             // 不分扎上报
}

// TableName 返回表名
func (Style) TableName() string {
	return "styles"
}

// CuttingTask 裁剪任务
type CuttingTask struct {
	ID           string         `json:"id" bson:"_id,omitempty"`
	OrderID      string         `json:"order_id" bson:"order_id"`           // 订单ID
	ContractNo   string         `json:"contract_no" bson:"contract_no"`     // 合同号
	StyleNo      string         `json:"style_no" bson:"style_no"`           // 款号
	StyleName    string         `json:"style_name" bson:"style_name"`       // 款名
	CustomerName string         `json:"customer_name" bson:"customer_name"` // 客户名称
	TotalPieces  int            `json:"total_pieces" bson:"total_pieces"`   // 总件数
	CutPieces    int            `json:"cut_pieces" bson:"cut_pieces"`       // 已裁件数
	Status       int            `json:"status" bson:"status"`               // 状态：0-待裁剪 1-裁剪中 2-已完成
	Batches      []CuttingBatch `json:"batches" bson:"batches"`             // 裁剪批次列表
	IsDeleted    int            `json:"is_deleted" bson:"is_deleted"`       // 是否删除：0-否 1-是
	CreatedBy    string         `json:"created_by" bson:"created_by"`       // 创建人
	UpdatedBy    string         `json:"updated_by" bson:"updated_by"`       // 更新人
	CreatedAt    int64          `json:"created_at" bson:"created_at"`       // 创建时间
	UpdatedAt    int64          `json:"updated_at" bson:"updated_at"`       // 更新时间
	DeletedAt    int64          `json:"deleted_at" bson:"deleted_at"`       // 删除时间
}

// CuttingBatch 裁剪批次（制菲）
type CuttingBatch struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	TaskID      string       `json:"task_id" bson:"task_id"`           // 裁剪任务ID
	OrderID     string       `json:"order_id" bson:"order_id"`         // 订单ID
	ContractNo  string       `json:"contract_no" bson:"contract_no"`   // 合同号
	StyleNo     string       `json:"style_no" bson:"style_no"`         // 款号
	BedNo       string       `json:"bed_no" bson:"bed_no"`             // 床号
	BundleNo    string       `json:"bundle_no" bson:"bundle_no"`       // 扎号
	Color       string       `json:"color" bson:"color"`               // 颜色名称
	LayerCount  int          `json:"layer_count" bson:"layer_count"`   // 拉布层数
	SizeDetails []SizeDetail `json:"size_details" bson:"size_details"` // 尺码明细
	TotalPieces int          `json:"total_pieces" bson:"total_pieces"` // 总件数
	QRCode      string       `json:"qr_code" bson:"qr_code"`           // 二维码内容
	PrintCount  int          `json:"print_count" bson:"print_count"`   // 打印次数
	IsDeleted   int          `json:"is_deleted" bson:"is_deleted"`     // 是否删除：0-否 1-是
	CreatedBy   string       `json:"created_by" bson:"created_by"`     // 创建人
	CreatedAt   int64        `json:"created_at" bson:"created_at"`     // 创建时间
	PrintedAt   int64        `json:"printed_at" bson:"printed_at"`     // 最后打印时间
}

// SizeDetail 尺码明细
type SizeDetail struct {
	Size     string `json:"size" bson:"size"`         // 尺码名称
	Quantity int    `json:"quantity" bson:"quantity"` // 数量
}

// CuttingPiece 裁片监控
type CuttingPiece struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	OrderID      string `json:"order_id" bson:"order_id"`           // 订单ID
	ContractNo   string `json:"contract_no" bson:"contract_no"`     // 合同号
	StyleNo      string `json:"style_no" bson:"style_no"`           // 款号
	BedNo        string `json:"bed_no" bson:"bed_no"`               // 床号
	BundleNo     string `json:"bundle_no" bson:"bundle_no"`         // 扎号
	Color        string `json:"color" bson:"color"`                 // 颜色名称
	Size         string `json:"size" bson:"size"`                   // 尺码名称
	Quantity     int    `json:"quantity" bson:"quantity"`           // 数量
	Progress     int    `json:"progress" bson:"progress"`           // 进度（已完成工序数）
	TotalProcess int    `json:"total_process" bson:"total_process"` // 总工序数
	CreatedAt    int64  `json:"created_at" bson:"created_at"`       // 创建时间
}

// TableName 返回表名
func (CuttingTask) TableName() string {
	return "cutting_tasks"
}

// TableName 返回表名
func (CuttingBatch) TableName() string {
	return "cutting_batches"
}

// TableName 返回表名
func (CuttingPiece) TableName() string {
	return "cutting_pieces"
}
