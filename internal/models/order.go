package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order 订单模型
type Order struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderNo       string             `bson:"order_no" json:"order_no" binding:"required"`
	CustomerID    primitive.ObjectID `bson:"customer_id" json:"customer_id" binding:"required"`
	StyleID       primitive.ObjectID `bson:"style_id" json:"style_id" binding:"required"`
	SalespersonID primitive.ObjectID `bson:"salesperson_id" json:"salesperson_id"`

	// 订单基本信息
	OrderType   string  `bson:"order_type" json:"order_type"`
	TotalQty    int     `bson:"total_qty" json:"total_qty"`
	UnitPrice   float64 `bson:"unit_price" json:"unit_price"`
	TotalAmount float64 `bson:"total_amount" json:"total_amount"`
	Currency    string  `bson:"currency" json:"currency" default:"CNY"`

	// 时间信息
	OrderDate    time.Time `bson:"order_date" json:"order_date"`
	DeliveryDate time.Time `bson:"delivery_date" json:"delivery_date"`

	// 状态信息
	Status   OrderStatus `bson:"status" json:"status"`
	Priority int         `bson:"priority" json:"priority" default:"1"`

	// 明细信息
	Items []OrderItem `bson:"items" json:"items"`

	// 备注信息
	Remark string `bson:"remark,omitempty" json:"remark"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Version   int                `bson:"version" json:"version"`
}

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusDraft      OrderStatus = 0 // 草稿
	OrderStatusConfirmed  OrderStatus = 1 // 已确认
	OrderStatusProduction OrderStatus = 2 // 生产中
	OrderStatusCompleted  OrderStatus = 3 // 已完成
	OrderStatusCancelled  OrderStatus = 4 // 已取消
)

// OrderItem 订单明细
type OrderItem struct {
	StyleID   primitive.ObjectID `bson:"style_id" json:"style_id"`
	ColorID   primitive.ObjectID `bson:"color_id" json:"color_id"`
	SizeID    primitive.ObjectID `bson:"size_id" json:"size_id"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	UnitPrice float64            `bson:"unit_price" json:"unit_price"`
	Amount    float64            `bson:"amount" json:"amount"`
	Remark    string             `bson:"remark,omitempty" json:"remark"`
}

// Style 款式模型
type Style struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StyleNo     string             `bson:"style_no" json:"style_no" binding:"required"`
	StyleName   string             `bson:"style_name" json:"style_name" binding:"required"`
	Category    string             `bson:"category" json:"category"`
	Season      string             `bson:"season" json:"season"`
	Year        int                `bson:"year" json:"year"`
	Description string             `bson:"description,omitempty" json:"description"`

	// 规格信息
	Fabric string   `bson:"fabric,omitempty" json:"fabric"`
	Colors []string `bson:"colors,omitempty" json:"colors"`
	Sizes  []string `bson:"sizes,omitempty" json:"sizes"`

	// 工艺信息
	ProcessFlow []ProcessStep `bson:"process_flow,omitempty" json:"process_flow"`

	// 成本信息
	MaterialCost   float64 `bson:"material_cost" json:"material_cost"`
	LaborCost      float64 `bson:"labor_cost" json:"labor_cost"`
	OverheadCost   float64 `bson:"overhead_cost" json:"overhead_cost"`
	TotalCost      float64 `bson:"total_cost" json:"total_cost"`
	SuggestedPrice float64 `bson:"suggested_price" json:"suggested_price"`

	// 图片信息
	Images []string `bson:"images,omitempty" json:"images"`

	// 状态信息
	Status StyleStatus `bson:"status" json:"status"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// StyleStatus 款式状态
type StyleStatus int

const (
	StyleStatusDraft    StyleStatus = 0 // 草稿
	StyleStatusActive   StyleStatus = 1 // 激活
	StyleStatusInactive StyleStatus = 2 // 停用
)

// ProcessStep 工序步骤
type ProcessStep struct {
	ProcessID     primitive.ObjectID `bson:"process_id" json:"process_id"`
	ProcessName   string             `bson:"process_name" json:"process_name"`
	Sequence      int                `bson:"sequence" json:"sequence"`
	Required      bool               `bson:"required" json:"required"`
	EstimatedTime float64            `bson:"estimated_time" json:"estimated_time"` // 预计工时（小时）
	StandardRate  float64            `bson:"standard_rate" json:"standard_rate"`   // 标准工价
}

// Customer 客户模型
type Customer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerNo  string             `bson:"customer_no" json:"customer_no" binding:"required"`
	CompanyName string             `bson:"company_name" json:"company_name" binding:"required"`
	ContactName string             `bson:"contact_name" json:"contact_name"`
	Phone       string             `bson:"phone" json:"phone"`
	Email       string             `bson:"email,omitempty" json:"email"`
	Address     string             `bson:"address,omitempty" json:"address"`

	// 业务信息
	CustomerType CustomerType `bson:"customer_type" json:"customer_type"`
	CreditLimit  float64      `bson:"credit_limit" json:"credit_limit"`
	PaymentTerm  string       `bson:"payment_term" json:"payment_term"`

	// 状态信息
	Status CustomerStatus `bson:"status" json:"status"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CustomerType 客户类型
type CustomerType int

const (
	CustomerTypeRegular   CustomerType = 1 // 普通客户
	CustomerTypeVIP       CustomerType = 2 // VIP客户
	CustomerTypeWholesale CustomerType = 3 // 批发客户
)

// CustomerStatus 客户状态
type CustomerStatus int

const (
	CustomerStatusActive   CustomerStatus = 1 // 激活
	CustomerStatusInactive CustomerStatus = 0 // 停用
)

// Salesperson 业务员模型
type Salesperson struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name   string             `bson:"name" json:"name" binding:"required"`
	Phone  string             `bson:"phone" json:"phone"`
	Email  string             `bson:"email,omitempty" json:"email"`
	Region string             `bson:"region,omitempty" json:"region"`

	// 业绩信息
	SalesTarget float64 `bson:"sales_target" json:"sales_target"`
	Commission  float64 `bson:"commission" json:"commission"`

	// 状态信息
	Status SalespersonStatus `bson:"status" json:"status"`

	// 审计信息
	CreatedBy primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// SalespersonStatus 业务员状态
type SalespersonStatus int

const (
	SalespersonStatusActive   SalespersonStatus = 1 // 激活
	SalespersonStatusInactive SalespersonStatus = 0 // 停用
)

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	Page          int       `form:"page,default=1" binding:"min=1"`
	PageSize      int       `form:"page_size,default=10" binding:"min=1,max=100"`
	OrderNo       string    `form:"order_no"`
	CustomerID    string    `form:"customer_id"`
	Status        *int      `form:"status"`
	StartDate     time.Time `form:"start_date"`
	EndDate       time.Time `form:"end_date"`
	SalespersonID string    `form:"salesperson_id"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	CustomerID    string      `json:"customer_id" binding:"required"`
	StyleID       string      `json:"style_id" binding:"required"`
	SalespersonID string      `json:"salesperson_id"`
	OrderType     string      `json:"order_type"`
	DeliveryDate  time.Time   `json:"delivery_date" binding:"required"`
	Items         []OrderItem `json:"items" binding:"required,min=1"`
	Remark        string      `json:"remark"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	CustomerID    string      `json:"customer_id"`
	StyleID       string      `json:"style_id"`
	SalespersonID string      `json:"salesperson_id"`
	OrderType     string      `json:"order_type"`
	DeliveryDate  time.Time   `json:"delivery_date"`
	Items         []OrderItem `json:"items"`
	Status        int         `json:"status"`
	Remark        string      `json:"remark"`
}

// CopyOrderRequest 复制订单请求
type CopyOrderRequest struct {
	OrderID      string    `json:"order_id" binding:"required"`
	DeliveryDate time.Time `json:"delivery_date" binding:"required"`
}

// StyleListRequest 款式列表请求
type StyleListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Category string `form:"category"`
	Season   string `form:"season"`
	Year     *int   `form:"year"`
	Status   *int   `form:"status"`
}

// CreateStyleRequest 创建款式请求
type CreateStyleRequest struct {
	StyleNo        string        `json:"style_no" binding:"required"`
	StyleName      string        `json:"style_name" binding:"required"`
	Category       string        `json:"category"`
	Season         string        `json:"season"`
	Year           int           `json:"year"`
	Description    string        `json:"description"`
	Fabric         string        `json:"fabric"`
	Colors         []string      `json:"colors"`
	Sizes          []string      `json:"sizes"`
	ProcessFlow    []ProcessStep `json:"process_flow"`
	MaterialCost   float64       `json:"material_cost"`
	LaborCost      float64       `json:"labor_cost"`
	OverheadCost   float64       `json:"overhead_cost"`
	SuggestedPrice float64       `json:"suggested_price"`
}

// CustomerListRequest 客户列表请求
type CustomerListRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Type     *int   `form:"type"`
	Status   *int   `form:"status"`
}

// CreateCustomerRequest 创建客户请求
type CreateCustomerRequest struct {
	CompanyName  string  `json:"company_name" binding:"required"`
	ContactName  string  `json:"contact_name"`
	Phone        string  `json:"phone"`
	Email        string  `json:"email"`
	Address      string  `json:"address"`
	CustomerType int     `json:"customer_type"`
	CreditLimit  float64 `json:"credit_limit"`
	PaymentTerm  string  `json:"payment_term"`
}
