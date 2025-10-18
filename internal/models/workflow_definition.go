package models

// WorkflowDefinition 工作流定义
type WorkflowDefinition struct {
	ID          string                 `bson:"_id,omitempty" json:"id"`
	Name        string                 `bson:"name" json:"name"`                             // 工作流名称
	Code        string                 `bson:"code" json:"code"`                             // 工作流唯一编码
	Description string                 `bson:"description" json:"description"`               // 描述
	States      []WorkflowState        `bson:"states" json:"states"`                         // 状态定义
	Transitions []WorkflowTransition   `bson:"transitions" json:"transitions"`               // 转换规则
	Version     int                    `bson:"version" json:"version"`                       // 版本号
	IsActive    bool                   `bson:"is_active" json:"is_active"`                   // 是否激活
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // 元数据
	CreatedAt   int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt   int64                  `bson:"updated_at" json:"updated_at"`
	CreatedBy   string                 `bson:"created_by" json:"created_by"`
	UpdatedBy   string                 `bson:"updated_by" json:"updated_by"`
}

// WorkflowState 工作流状态定义
type WorkflowState struct {
	ID          string                 `bson:"id" json:"id"`                                 // 状态ID
	Name        string                 `bson:"name" json:"name"`                             // 状态名称
	Code        string                 `bson:"code" json:"code"`                             // 状态编码
	Type        string                 `bson:"type" json:"type"`                             // 类型: start, normal, end
	Color       string                 `bson:"color" json:"color"`                           // 显示颜色
	Description string                 `bson:"description" json:"description"`               // 描述
	Position    *StatePosition         `bson:"position,omitempty" json:"position,omitempty"` // 可视化位置
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // 元数据
}

// StatePosition 状态在画布上的位置
type StatePosition struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

// WorkflowTransition 工作流转换规则
type WorkflowTransition struct {
	ID          string                 `bson:"id" json:"id"`                                 // 转换ID
	Name        string                 `bson:"name" json:"name"`                             // 转换名称
	FromState   string                 `bson:"from_state" json:"from_state"`                 // 起始状态ID
	ToState     string                 `bson:"to_state" json:"to_state"`                     // 目标状态ID
	Event       string                 `bson:"event" json:"event"`                           // 触发事件
	Conditions  []TransitionCondition  `bson:"conditions" json:"conditions"`                 // 转换条件
	Actions     []TransitionAction     `bson:"actions" json:"actions"`                       // 转换动作
	RequireRole string                 `bson:"require_role" json:"require_role"`             // 需要的角色
	Description string                 `bson:"description" json:"description"`               // 描述
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // 元数据
}

// TransitionCondition 转换条件
type TransitionCondition struct {
	Type        string                 `bson:"type" json:"type"`         // 条件类型: field, script, custom
	Field       string                 `bson:"field" json:"field"`       // 字段名
	Operator    string                 `bson:"operator" json:"operator"` // 操作符: eq, gt, gte, lt, lte, in, contains
	Value       interface{}            `bson:"value" json:"value"`       // 比较值
	Script      string                 `bson:"script" json:"script"`     // 脚本表达式
	Description string                 `bson:"description" json:"description"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// TransitionAction 转换动作
type TransitionAction struct {
	Type        string                 `bson:"type" json:"type"`     // 动作类型: update_field, send_notification, custom
	Field       string                 `bson:"field" json:"field"`   // 更新的字段
	Value       interface{}            `bson:"value" json:"value"`   // 更新的值
	Script      string                 `bson:"script" json:"script"` // 脚本
	Description string                 `bson:"description" json:"description"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

// WorkflowInstance 工作流实例（订单的工作流状态）
type WorkflowInstance struct {
	ID           string                 `bson:"_id,omitempty" json:"id"`
	WorkflowID   string                 `bson:"workflow_id" json:"workflow_id"`     // 工作流定义ID
	EntityType   string                 `bson:"entity_type" json:"entity_type"`     // 实体类型: order, task等
	EntityID     string                 `bson:"entity_id" json:"entity_id"`         // 实体ID（如订单ID）
	CurrentState string                 `bson:"current_state" json:"current_state"` // 当前状态ID
	History      []WorkflowHistory      `bson:"history" json:"history"`             // 历史记录
	Variables    map[string]interface{} `bson:"variables" json:"variables"`         // 工作流变量
	CreatedAt    int64                  `bson:"created_at" json:"created_at"`
	UpdatedAt    int64                  `bson:"updated_at" json:"updated_at"`
}

// WorkflowHistory 工作流历史记录
type WorkflowHistory struct {
	FromState string                 `bson:"from_state" json:"from_state"`
	ToState   string                 `bson:"to_state" json:"to_state"`
	Event     string                 `bson:"event" json:"event"`
	Operator  string                 `bson:"operator" json:"operator"`
	Reason    string                 `bson:"reason" json:"reason"`
	Timestamp int64                  `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
}
