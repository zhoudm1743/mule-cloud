package dto

import "mule-cloud/internal/models"

// WorkflowDefinitionRequest 工作流定义请求
type WorkflowDefinitionRequest struct {
	Name        string                      `json:"name" binding:"required"`
	Code        string                      `json:"code" binding:"required"`
	Description string                      `json:"description"`
	States      []models.WorkflowState      `json:"states" binding:"required"`
	Transitions []models.WorkflowTransition `json:"transitions"`
	Metadata    map[string]interface{}      `json:"metadata,omitempty"`
}

// WorkflowListRequest 工作流列表请求
type WorkflowListRequest struct {
	Page     int64 `json:"page" form:"page"`
	PageSize int64 `json:"page_size" form:"page_size"`
}

// WorkflowListResponse 工作流列表响应
type WorkflowListResponse struct {
	Workflows []*models.WorkflowDefinition `json:"workflows"`
	Total     int64                        `json:"total"`
}

// ExecuteTransitionRequest 执行转换请求
type ExecuteTransitionRequest struct {
	InstanceID string                 `json:"instance_id" binding:"required"`
	Event      string                 `json:"event" binding:"required"`
	Reason     string                 `json:"reason"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// LogicFlowData LogicFlow图数据格式
type LogicFlowData struct {
	Nodes []LogicFlowNode `json:"nodes"`
	Edges []LogicFlowEdge `json:"edges"`
}

// LogicFlowNode LogicFlow节点
type LogicFlowNode struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	X          float64                `json:"x"`
	Y          float64                `json:"y"`
	Text       string                 `json:"text"`
	Properties map[string]interface{} `json:"properties"`
}

// LogicFlowEdge LogicFlow连线
type LogicFlowEdge struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	SourceNode string                 `json:"sourceNodeId"`
	TargetNode string                 `json:"targetNodeId"`
	Text       string                 `json:"text"`
	Properties map[string]interface{} `json:"properties"`
}

// ConvertToLogicFlowRequest 转换为LogicFlow格式请求
type ConvertToLogicFlowRequest struct {
	WorkflowID string `json:"workflow_id" binding:"required"`
}

// ConvertFromLogicFlowRequest 从LogicFlow格式转换请求
type ConvertFromLogicFlowRequest struct {
	Data LogicFlowData `json:"data" binding:"required"`
}

// WorkflowTemplateResponse 工作流模板响应
type WorkflowTemplateResponse struct {
	Templates []WorkflowTemplate `json:"templates"`
}

// WorkflowTemplate 工作流模板
type WorkflowTemplate struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Code        string                      `json:"code"`
	Description string                      `json:"description"`
	Category    string                      `json:"category"`
	Icon        string                      `json:"icon"`
	Preview     string                      `json:"preview"`
	States      []WorkflowTemplateState     `json:"states,omitempty"`
	Transitions []WorkflowTemplateTransition `json:"transitions,omitempty"`
}

// WorkflowTemplateState 模板状态定义
type WorkflowTemplateState struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        string `json:"type"` // start, normal, end
	Color       string `json:"color"`
	Description string `json:"description"`
}

// WorkflowTemplateTransition 模板转换定义
type WorkflowTemplateTransition struct {
	From            string                   `json:"from"`
	To              string                   `json:"to"`
	Event           string                   `json:"event"`
	EventLabel      string                   `json:"event_label"`
	HasCondition    bool                     `json:"has_condition"`
	ConditionDesc   string                   `json:"condition_desc,omitempty"`
	RequireRole     string                   `json:"require_role,omitempty"`
	RoleDesc        string                   `json:"role_desc,omitempty"`
	AvailableFields []WorkflowConditionField `json:"available_fields,omitempty"`
}

// WorkflowConditionField 可用的条件字段
type WorkflowConditionField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"` // number, string, boolean
	Description string `json:"description"`
}
