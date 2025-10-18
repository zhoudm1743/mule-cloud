package dto

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

