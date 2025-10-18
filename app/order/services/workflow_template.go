package services

import (
	"context"
)

// WorkflowTemplateService 工作流模板服务
type WorkflowTemplateService struct{}

// NewWorkflowTemplateService 创建工作流模板服务
func NewWorkflowTemplateService() *WorkflowTemplateService {
	return &WorkflowTemplateService{}
}

// WorkflowTemplate 工作流模板
type WorkflowTemplate struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Code        string                   `json:"code"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Icon        string                   `json:"icon"`
	Preview     string                   `json:"preview"`
	States      []WorkflowTemplateState  `json:"states"`
	Transitions []WorkflowTemplateTransition `json:"transitions"`
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

// GetTemplates 获取所有模板
func (s *WorkflowTemplateService) GetTemplates(ctx context.Context) []WorkflowTemplate {
	return []WorkflowTemplate{
		s.getBasicOrderTemplate(),
	}
}

// getBasicOrderTemplate 获取基础订单流程模板（基于 order_workflow_advanced.go）
func (s *WorkflowTemplateService) getBasicOrderTemplate() WorkflowTemplate {
	return WorkflowTemplate{
		ID:          "basic_order_workflow",
		Name:        "基础订单流程",
		Code:        "basic_order",
		Description: "标准的订单处理流程，包含草稿、已下单、生产中、已完成、已取消五个状态，支持进度检查和权限控制",
		Category:    "订单管理",
		Icon:        "📦",
		Preview:     "草稿 → 已下单 → 生产中 → 已完成",
		States: []WorkflowTemplateState{
			{
				Code:        "draft",
				Name:        "草稿",
				Type:        "start",
				Color:       "#909399",
				Description: "订单初始状态，可以编辑订单信息",
			},
			{
				Code:        "ordered",
				Name:        "已下单",
				Type:        "normal",
				Color:       "#409EFF",
				Description: "订单已提交，等待开始生产",
			},
			{
				Code:        "production",
				Name:        "生产中",
				Type:        "normal",
				Color:       "#E6A23C",
				Description: "订单正在生产，可以更新进度",
			},
			{
				Code:        "completed",
				Name:        "已完成",
				Type:        "end",
				Color:       "#67C23A",
				Description: "订单已完成，进度达到100%",
			},
			{
				Code:        "cancelled",
				Name:        "已取消",
				Type:        "end",
				Color:       "#F56C6C",
				Description: "订单已取消，需要管理员权限",
			},
		},
		Transitions: []WorkflowTemplateTransition{
			{
				From:       "draft",
				To:         "ordered",
				Event:      "submit_order",
				EventLabel: "提交订单",
			},
			{
				From:       "ordered",
				To:         "production",
				Event:      "start_cutting",
				EventLabel: "开始裁剪",
			},
			{
				From:       "ordered",
				To:         "production",
				Event:      "start_production",
				EventLabel: "开始生产",
			},
			{
				From:       "production",
				To:         "production",
				Event:      "update_progress",
				EventLabel: "更新进度",
			},
			{
				From:          "production",
				To:            "completed",
				Event:         "complete",
				EventLabel:    "完成",
				HasCondition:  true,
				ConditionDesc: "进度必须达到100%",
				AvailableFields: []WorkflowConditionField{
					{
						Key:         "progress",
						Label:       "生产进度",
						Type:        "number",
						Description: "订单的完成进度（0-1之间的小数）",
					},
				},
			},
			{
				From:        "draft",
				To:          "cancelled",
				Event:       "cancel",
				EventLabel:  "取消",
				RequireRole: "admin",
				RoleDesc:    "需要管理员权限",
			},
			{
				From:        "ordered",
				To:          "cancelled",
				Event:       "cancel",
				EventLabel:  "取消",
				RequireRole: "admin",
				RoleDesc:    "需要管理员权限",
			},
			{
				From:        "production",
				To:          "cancelled",
				Event:       "cancel",
				EventLabel:  "取消",
				RequireRole: "admin",
				RoleDesc:    "需要管理员权限",
			},
		},
	}
}

