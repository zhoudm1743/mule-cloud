package services

import (
	"context"
	"fmt"

	corecontext "mule-cloud/core/context"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
)

// IWorkflowDesignerService 工作流设计器服务接口
type IWorkflowDesignerService interface {
	// 工作流定义管理
	CreateDefinition(ctx context.Context, workflow *models.WorkflowDefinition) (*models.WorkflowDefinition, error)
	UpdateDefinition(ctx context.Context, id string, workflow *models.WorkflowDefinition) error
	GetDefinition(ctx context.Context, id string) (*models.WorkflowDefinition, error)
	ListDefinitions(ctx context.Context, page, pageSize int64) ([]*models.WorkflowDefinition, int64, error)
	DeleteDefinition(ctx context.Context, id string) error
	ActivateDefinition(ctx context.Context, id string) error
	DeactivateDefinition(ctx context.Context, id string) error

	// 工作流实例管理
	GetInstance(ctx context.Context, entityType, entityID string) (*models.WorkflowInstance, error)
	ExecuteTransition(ctx context.Context, instanceID, event string, operator, reason string, metadata map[string]interface{}) error
}

type workflowDesignerService struct {
	defRepo  repository.IWorkflowDefinitionRepository
	instRepo repository.IWorkflowInstanceRepository
}

// NewWorkflowDesignerService 创建工作流设计器服务
func NewWorkflowDesignerService() IWorkflowDesignerService {
	return &workflowDesignerService{
		defRepo:  repository.NewWorkflowDefinitionRepository(),
		instRepo: repository.NewWorkflowInstanceRepository(),
	}
}

func (s *workflowDesignerService) CreateDefinition(ctx context.Context, workflow *models.WorkflowDefinition) (*models.WorkflowDefinition, error) {
	// 设置创建人
	workflow.CreatedBy = corecontext.GetUsername(ctx)
	workflow.UpdatedBy = corecontext.GetUsername(ctx)

	// 验证工作流定义
	if err := s.validateDefinition(workflow); err != nil {
		return nil, err
	}

	// 创建
	if err := s.defRepo.Create(ctx, workflow); err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *workflowDesignerService) UpdateDefinition(ctx context.Context, id string, workflow *models.WorkflowDefinition) error {
	// 验证工作流定义
	if err := s.validateDefinition(workflow); err != nil {
		return err
	}

	// 更新
	workflow.UpdatedBy = corecontext.GetUsername(ctx)
	return s.defRepo.Update(ctx, id, workflow)
}

func (s *workflowDesignerService) GetDefinition(ctx context.Context, id string) (*models.WorkflowDefinition, error) {
	return s.defRepo.Get(ctx, id)
}

func (s *workflowDesignerService) ListDefinitions(ctx context.Context, page, pageSize int64) ([]*models.WorkflowDefinition, int64, error) {
	return s.defRepo.List(ctx, page, pageSize)
}

func (s *workflowDesignerService) DeleteDefinition(ctx context.Context, id string) error {
	return s.defRepo.Delete(ctx, id)
}

func (s *workflowDesignerService) ActivateDefinition(ctx context.Context, id string) error {
	return s.defRepo.Activate(ctx, id)
}

func (s *workflowDesignerService) DeactivateDefinition(ctx context.Context, id string) error {
	return s.defRepo.Deactivate(ctx, id)
}

func (s *workflowDesignerService) GetInstance(ctx context.Context, entityType, entityID string) (*models.WorkflowInstance, error) {
	return s.instRepo.GetByEntity(ctx, entityType, entityID)
}

func (s *workflowDesignerService) ExecuteTransition(ctx context.Context, instanceID, event string, operator, reason string, metadata map[string]interface{}) error {
	// 获取实例
	instance, err := s.instRepo.Get(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("获取工作流实例失败: %v", err)
	}

	// 获取工作流定义
	workflow, err := s.defRepo.Get(ctx, instance.WorkflowID)
	if err != nil {
		return fmt.Errorf("获取工作流定义失败: %v", err)
	}

	// 查找匹配的转换规则
	var transition *models.WorkflowTransition
	for i := range workflow.Transitions {
		t := &workflow.Transitions[i]
		if t.FromState == instance.CurrentState && t.Event == event {
			transition = t
			break
		}
	}

	if transition == nil {
		return fmt.Errorf("未找到有效的转换规则: 从状态 %s 通过事件 %s", instance.CurrentState, event)
	}

	// 检查条件
	if !s.checkConditions(transition.Conditions, instance, metadata) {
		return fmt.Errorf("转换条件不满足")
	}

	// 执行转换
	history := models.WorkflowHistory{
		FromState: instance.CurrentState,
		ToState:   transition.ToState,
		Event:     event,
		Operator:  operator,
		Reason:    reason,
		Timestamp: 0, // Will be set in repository
		Metadata:  metadata,
	}

	// 更新实例状态
	err = s.instRepo.Update(ctx, instanceID, map[string]interface{}{
		"current_state": transition.ToState,
	})
	if err != nil {
		return fmt.Errorf("更新工作流实例失败: %v", err)
	}

	// 添加历史记录
	err = s.instRepo.AddHistory(ctx, instanceID, history)
	if err != nil {
		return fmt.Errorf("添加历史记录失败: %v", err)
	}

	// 执行动作
	s.executeActions(transition.Actions, instance, metadata)

	return nil
}

// validateDefinition 验证工作流定义
func (s *workflowDesignerService) validateDefinition(workflow *models.WorkflowDefinition) error {
	if workflow.Name == "" {
		return fmt.Errorf("工作流名称不能为空")
	}
	if workflow.Code == "" {
		return fmt.Errorf("工作流编码不能为空")
	}
	if len(workflow.States) == 0 {
		return fmt.Errorf("至少需要定义一个状态")
	}

	// 验证状态
	stateMap := make(map[string]bool)
	hasStart := false
	hasEnd := false
	for _, state := range workflow.States {
		if state.ID == "" || state.Name == "" {
			return fmt.Errorf("状态ID和名称不能为空")
		}
		if stateMap[state.ID] {
			return fmt.Errorf("状态ID重复: %s", state.ID)
		}
		stateMap[state.ID] = true

		if state.Type == "start" {
			hasStart = true
		}
		if state.Type == "end" {
			hasEnd = true
		}
	}

	if !hasStart {
		return fmt.Errorf("至少需要一个起始状态")
	}
	if !hasEnd {
		return fmt.Errorf("至少需要一个结束状态")
	}

	// 验证转换规则
	for _, trans := range workflow.Transitions {
		if !stateMap[trans.FromState] {
			return fmt.Errorf("转换规则引用了不存在的起始状态: %s", trans.FromState)
		}
		if !stateMap[trans.ToState] {
			return fmt.Errorf("转换规则引用了不存在的目标状态: %s", trans.ToState)
		}
	}

	return nil
}

// checkConditions 检查转换条件
func (s *workflowDesignerService) checkConditions(conditions []models.TransitionCondition, instance *models.WorkflowInstance, metadata map[string]interface{}) bool {
	if len(conditions) == 0 {
		return true // 无条件，直接通过
	}

	for _, condition := range conditions {
		if !s.checkCondition(condition, instance, metadata) {
			return false
		}
	}

	return true
}

// checkCondition 检查单个条件
func (s *workflowDesignerService) checkCondition(condition models.TransitionCondition, instance *models.WorkflowInstance, metadata map[string]interface{}) bool {
	switch condition.Type {
	case "field":
		// 从metadata或variables中获取字段值
		var fieldValue interface{}
		if val, ok := metadata[condition.Field]; ok {
			fieldValue = val
		} else if val, ok := instance.Variables[condition.Field]; ok {
			fieldValue = val
		} else {
			return false
		}

		// 比较值
		return s.compareValues(fieldValue, condition.Operator, condition.Value)

	case "script":
		// TODO: 实现脚本执行（可以集成JavaScript引擎）
		return true

	default:
		return true
	}
}

// compareValues 比较值
func (s *workflowDesignerService) compareValues(fieldValue interface{}, operator string, compareValue interface{}) bool {
	switch operator {
	case "eq":
		return fieldValue == compareValue
	case "gt":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := compareValue.(float64); ok {
				return fv > cv
			}
		}
	case "gte":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := compareValue.(float64); ok {
				return fv >= cv
			}
		}
	case "lt":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := compareValue.(float64); ok {
				return fv < cv
			}
		}
	case "lte":
		if fv, ok := fieldValue.(float64); ok {
			if cv, ok := compareValue.(float64); ok {
				return fv <= cv
			}
		}
	}
	return false
}

// executeActions 执行转换动作
func (s *workflowDesignerService) executeActions(actions []models.TransitionAction, instance *models.WorkflowInstance, metadata map[string]interface{}) {
	for _, action := range actions {
		switch action.Type {
		case "update_field":
			// 更新工作流变量
			if instance.Variables == nil {
				instance.Variables = make(map[string]interface{})
			}
			instance.Variables[action.Field] = action.Value

		case "send_notification":
			// TODO: 发送通知

		default:
			// 自定义动作
		}
	}
}
