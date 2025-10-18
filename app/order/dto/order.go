package dto

import "mule-cloud/internal/models"

// OrderListRequest 订单列表请求
type OrderListRequest struct {
	ID          string `uri:"id" query:"id"`
	ContractNo  string `query:"contract_no"`   // 合同号
	StyleNo     string `query:"style_no"`      // 款号
	CustomerID  string `query:"customer_id"`   // 客户ID
	SalesmanID  string `query:"salesman_id"`   // 业务员ID
	OrderTypeID string `query:"order_type_id"` // 订单类型ID
	Status      int    `query:"status"`        // 状态
	StartDate   string `query:"start_date"`    // 开始日期（交货日期范围）
	EndDate     string `query:"end_date"`      // 结束日期
	OrderStart  string `query:"order_start"`   // 下单时间范围开始
	OrderEnd    string `query:"order_end"`     // 下单时间范围结束
	Remark      string `query:"remark"`        // 备注

	Page     int64 `query:"page"`
	PageSize int64 `query:"page_size"`
}

// OrderCreateRequest 创建订单请求（步骤1：基础信息）
type OrderCreateRequest struct {
	ContractNo   string `json:"contract_no" binding:"required"` // 合同号
	CustomerID   string `json:"customer_id" binding:"required"` // 客户ID
	DeliveryDate string `json:"delivery_date"`                  // 交货日期
	OrderTypeID  string `json:"order_type_id"`                  // 订单类型ID
	SalesmanID   string `json:"salesman_id"`                    // 业务员ID
	Remark       string `json:"remark"`                         // 备注
}

// OrderStyleRequest 订单款式数量（步骤2）
type OrderStyleRequest struct {
	ID        string             `uri:"id" binding:"required"`
	StyleID   string             `json:"style_id" binding:"required"`   // 款式ID
	Colors    []string           `json:"colors"`                        // 选择的颜色列表
	Sizes     []string           `json:"sizes"`                         // 选择的尺码列表
	UnitPrice float64            `json:"unit_price" binding:"required"` // 单价
	Quantity  int                `json:"quantity" binding:"required"`   // 总数量
	Items     []models.OrderItem `json:"items" binding:"required"`      // 订单明细（颜色+尺码组合）
}

// OrderProcedureRequest 订单工序（步骤3）
type OrderProcedureRequest struct {
	ID         string                  `uri:"id" binding:"required"`
	Procedures []models.OrderProcedure `json:"procedures"` // 工序清单
}

// OrderUpdateRequest 更新订单请求
type OrderUpdateRequest struct {
	ID           string                  `uri:"id" binding:"required"`
	ContractNo   string                  `json:"contract_no"`   // 合同号
	StyleID      string                  `json:"style_id"`      // 款式ID
	CustomerID   string                  `json:"customer_id"`   // 客户ID
	SalesmanID   string                  `json:"salesman_id"`   // 业务员ID
	OrderTypeID  string                  `json:"order_type_id"` // 订单类型ID
	Colors       []string                `json:"colors"`        // 颜色列表
	Sizes        []string                `json:"sizes"`         // 尺码列表
	UnitPrice    float64                 `json:"unit_price"`    // 单价
	Quantity     int                     `json:"quantity"`      // 数量
	DeliveryDate string                  `json:"delivery_date"` // 交货日期
	Status       int                     `json:"status"`        // 状态
	Remark       string                  `json:"remark"`        // 备注
	Items        []models.OrderItem      `json:"items"`         // 订单明细
	Procedures   []models.OrderProcedure `json:"procedures"`    // 工序清单
}

// OrderCopyRequest 复制订单请求
type OrderCopyRequest struct {
	ID             string `uri:"id" binding:"required"`
	IsRelated      bool   `json:"is_related"`      // 是否关联原订单
	RelationType   string `json:"relation_type"`   // 关联类型：copy-复制 add-追加
	RelationRemark string `json:"relation_remark"` // 关联说明
}

// CuttingBatchRequest 裁剪制菲请求
type CuttingBatchRequest struct {
	OrderID      string `json:"order_id" binding:"required"`       // 订单ID
	BedNo        string `json:"bed_no" binding:"required"`         // 床号
	BundleNo     string `json:"bundle_no" binding:"required"`      // 扎号
	ColorID      string `json:"color_id" binding:"required"`       // 颜色ID
	LayerCount   int    `json:"layer_count" binding:"required"`    // 拉布层数
	PiecesPerSet int    `json:"pieces_per_set" binding:"required"` // 每扎件数
}

// OrderResponse 订单响应
type OrderResponse struct {
	Order *models.Order `json:"order"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Orders []models.Order `json:"orders"`
	Total  int64          `json:"total"`
}

// OrderWorkflowTransitionRequest 订单工作流状态转换请求
type OrderWorkflowTransitionRequest struct {
	ID       string                 `uri:"id" binding:"required"`     // 订单ID
	Event    string                 `json:"event" binding:"required"` // 事件名称
	Operator string                 `json:"operator"`                 // 操作人
	Reason   string                 `json:"reason"`                   // 转换原因
	Metadata map[string]interface{} `json:"metadata"`                 // 元数据
}

// OrderWorkflowStateResponse 订单工作流状态响应
type OrderWorkflowStateResponse struct {
	Instance *models.WorkflowInstance `json:"instance"`
}

// OrderWorkflowTransitionsResponse 订单可用转换响应
type OrderWorkflowTransitionsResponse struct {
	Transitions []models.WorkflowTransition `json:"transitions"`
}
