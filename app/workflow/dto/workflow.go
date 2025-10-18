package dto

// TransitionRequest 状态转换请求
type TransitionRequest struct {
	OrderID  string                 `json:"order_id" binding:"required"`
	Event    string                 `json:"event" binding:"required"`
	Reason   string                 `json:"reason"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// RollbackRequest 回滚请求
type RollbackRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Reason  string `json:"reason" binding:"required"`
}

// OrderStatusResponse 订单状态响应
type OrderStatusResponse struct {
	OrderID    string `json:"order_id"`
	Status     int    `json:"status"`
	StatusName string `json:"status_name"`
}

// WorkflowDefinitionResponse 工作流定义响应
type WorkflowDefinitionResponse struct {
	States      []map[string]interface{} `json:"states"`
	Events      []map[string]interface{} `json:"events"`
	Transitions []map[string]interface{} `json:"transitions"`
}

// MermaidDiagramResponse Mermaid 流程图响应
type MermaidDiagramResponse struct {
	Diagram string `json:"diagram"`
}

