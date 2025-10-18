package services

import (
	"context"
	"fmt"
	"time"

	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IWorkflowEngineService 工作流引擎服务接口
type IWorkflowEngineService interface {
	// InitOrderWorkflow 为订单初始化工作流实例
	InitOrderWorkflow(ctx context.Context, orderID string, workflowCode string) error

	// TransitionOrderState 执行订单状态转换
	TransitionOrderState(ctx context.Context, orderID string, event string, operator string, reason string, metadata map[string]interface{}) error

	// GetOrderWorkflowState 获取订单当前工作流状态
	GetOrderWorkflowState(ctx context.Context, orderID string) (*models.WorkflowInstance, error)

	// GetAvailableTransitions 获取订单当前可用的转换
	GetAvailableTransitions(ctx context.Context, orderID string) ([]models.WorkflowTransition, error)
}

type workflowEngineService struct {
	orderRepo            repository.OrderRepository
	workflowDefRepo      repository.IWorkflowDefinitionRepository
	workflowInstanceRepo repository.IWorkflowInstanceRepository
}

// NewWorkflowEngineService 创建工作流引擎服务
func NewWorkflowEngineService() IWorkflowEngineService {
	return &workflowEngineService{
		orderRepo:            repository.NewOrderRepository(),
		workflowDefRepo:      repository.NewWorkflowDefinitionRepository(),
		workflowInstanceRepo: repository.NewWorkflowInstanceRepository(),
	}
}

// InitOrderWorkflow 为订单初始化工作流实例
func (s *workflowEngineService) InitOrderWorkflow(ctx context.Context, orderID string, workflowCode string) error {
	// 获取订单
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return fmt.Errorf("订单不存在: %v", err)
	}

	// 获取激活的工作流定义
	definition, err := s.workflowDefRepo.GetActive(ctx, workflowCode)
	if err != nil {
		return fmt.Errorf("工作流定义不存在或未激活: %v", err)
	}

	// 检查是否已有工作流实例
	if order.WorkflowInstance != "" {
		return fmt.Errorf("订单已关联工作流实例")
	}

	// 找到开始状态
	var startState *models.WorkflowState
	for i := range definition.States {
		if definition.States[i].Type == "start" {
			startState = &definition.States[i]
			break
		}
	}
	if startState == nil {
		return fmt.Errorf("工作流定义没有开始状态")
	}

	// 创建工作流实例
	instance := &models.WorkflowInstance{
		ID:           primitive.NewObjectID().Hex(),
		WorkflowID:   definition.ID,
		EntityType:   "order",
		EntityID:     orderID,
		CurrentState: startState.Code,
		Variables:    make(map[string]interface{}),
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	err = s.workflowInstanceRepo.Create(ctx, instance)
	if err != nil {
		return fmt.Errorf("创建工作流实例失败: %v", err)
	}

	// 添加历史记录
	err = s.workflowInstanceRepo.AddHistory(ctx, instance.ID, models.WorkflowHistory{
		FromState: "",
		ToState:   startState.Code,
		Event:     "init",
		Operator:  "system",
		Reason:    "初始化工作流",
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("添加历史记录失败: %v", err)
	}

	// 更新订单
	err = s.orderRepo.Update(ctx, orderID, map[string]interface{}{
		"workflow_code":     workflowCode,
		"workflow_instance": instance.ID,
		"workflow_state":    startState.Code,
		"updated_at":        time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("更新订单失败: %v", err)
	}

	return nil
}

// TransitionOrderState 执行订单状态转换
func (s *workflowEngineService) TransitionOrderState(
	ctx context.Context,
	orderID string,
	event string,
	operator string,
	reason string,
	metadata map[string]interface{},
) error {
	// 获取订单
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return fmt.Errorf("订单不存在: %v", err)
	}

	// 检查是否已关联工作流
	if order.WorkflowInstance == "" {
		return fmt.Errorf("订单未关联工作流实例")
	}

	// 获取工作流实例
	instance, err := s.workflowInstanceRepo.Get(ctx, order.WorkflowInstance)
	if err != nil {
		return fmt.Errorf("工作流实例不存在: %v", err)
	}

	// 获取工作流定义
	definition, err := s.workflowDefRepo.Get(ctx, instance.WorkflowID)
	if err != nil {
		return fmt.Errorf("工作流定义不存在: %v", err)
	}

	// 查找匹配的转换规则
	var matchedTransition *models.WorkflowTransition
	for i := range definition.Transitions {
		trans := &definition.Transitions[i]
		if trans.FromState == instance.CurrentState && trans.Event == event {
			matchedTransition = trans
			break
		}
	}

	if matchedTransition == nil {
		return fmt.Errorf("没有找到匹配的状态转换规则: %s -> %s", instance.CurrentState, event)
	}

	// 检查转换条件
	if len(matchedTransition.Conditions) > 0 {
		// 合并订单数据和元数据
		contextData := make(map[string]interface{})
		contextData["order_id"] = orderID
		contextData["total_amount"] = order.TotalAmount
		contextData["quantity"] = order.Quantity
		contextData["progress"] = order.Progress
		contextData["status"] = order.Status
		for k, v := range metadata {
			contextData[k] = v
		}

		// 检查所有条件
		for _, condition := range matchedTransition.Conditions {
			if !s.checkCondition(condition, contextData) {
				return fmt.Errorf("转换条件不满足: %s", condition.Description)
			}
		}
	}

	// 执行转换
	oldState := instance.CurrentState
	newState := matchedTransition.ToState

	// 更新工作流实例
	instance.CurrentState = newState
	instance.UpdatedAt = time.Now().Unix()

	// 合并变量
	if metadata != nil {
		for k, v := range metadata {
			instance.Variables[k] = v
		}
	}

	err = s.workflowInstanceRepo.Update(ctx, instance.ID, map[string]interface{}{
		"current_state": newState,
		"variables":     instance.Variables,
		"updated_at":    instance.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("更新工作流实例失败: %v", err)
	}

	// 添加历史记录
	err = s.workflowInstanceRepo.AddHistory(ctx, instance.ID, models.WorkflowHistory{
		FromState: oldState,
		ToState:   newState,
		Event:     event,
		Operator:  operator,
		Reason:    reason,
		Timestamp: time.Now().Unix(),
		Metadata:  metadata,
	})
	if err != nil {
		return fmt.Errorf("添加历史记录失败: %v", err)
	}

	// 更新订单状态
	orderUpdate := map[string]interface{}{
		"workflow_state": newState,
		"updated_at":     time.Now().Unix(),
	}

	// 根据工作流状态映射订单状态（可选）
	statusMapping := map[string]int{
		"draft":      0, // 草稿
		"ordered":    1, // 已下单
		"production": 2, // 生产中
		"completed":  3, // 已完成
		"cancelled":  4, // 已取消
	}
	if mappedStatus, ok := statusMapping[newState]; ok {
		orderUpdate["status"] = mappedStatus
	}

	err = s.orderRepo.Update(ctx, orderID, orderUpdate)
	if err != nil {
		return fmt.Errorf("更新订单失败: %v", err)
	}

	// 执行转换动作
	s.executeActions(ctx, matchedTransition.Actions, orderID, metadata)

	return nil
}

// GetOrderWorkflowState 获取订单当前工作流状态
func (s *workflowEngineService) GetOrderWorkflowState(ctx context.Context, orderID string) (*models.WorkflowInstance, error) {
	// 获取订单
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %v", err)
	}

	if order.WorkflowInstance == "" {
		return nil, fmt.Errorf("订单未关联工作流实例")
	}

	// 获取工作流实例
	return s.workflowInstanceRepo.Get(ctx, order.WorkflowInstance)
}

// GetAvailableTransitions 获取订单当前可用的转换
func (s *workflowEngineService) GetAvailableTransitions(ctx context.Context, orderID string) ([]models.WorkflowTransition, error) {
	// 获取订单
	order, err := s.orderRepo.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("订单不存在: %v", err)
	}

	if order.WorkflowInstance == "" {
		return nil, fmt.Errorf("订单未关联工作流实例")
	}

	// 获取工作流实例
	instance, err := s.workflowInstanceRepo.Get(ctx, order.WorkflowInstance)
	if err != nil {
		return nil, fmt.Errorf("工作流实例不存在: %v", err)
	}

	// 获取工作流定义
	definition, err := s.workflowDefRepo.Get(ctx, instance.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("工作流定义不存在: %v", err)
	}

	// 查找当前状态可用的转换
	var availableTransitions []models.WorkflowTransition
	for _, trans := range definition.Transitions {
		if trans.FromState == instance.CurrentState {
			availableTransitions = append(availableTransitions, trans)
		}
	}

	return availableTransitions, nil
}

// checkCondition 检查单个条件
func (s *workflowEngineService) checkCondition(condition models.TransitionCondition, data map[string]interface{}) bool {
	if condition.Type != "field" {
		return true // 暂不支持其他类型
	}

	// 获取字段值
	fieldValue, ok := data[condition.Field]
	if !ok {
		return false
	}

	// 类型转换和比较
	return s.compareValues(fieldValue, condition.Operator, condition.Value)
}

// compareValues 比较值
func (s *workflowEngineService) compareValues(actual interface{}, operator string, expected interface{}) bool {
	switch operator {
	case "eq": // 等于
		return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
	case "ne": // 不等于
		return fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected)
	case "gt": // 大于
		return s.compareNumeric(actual, expected, ">")
	case "gte": // 大于等于
		return s.compareNumeric(actual, expected, ">=")
	case "lt": // 小于
		return s.compareNumeric(actual, expected, "<")
	case "lte": // 小于等于
		return s.compareNumeric(actual, expected, "<=")
	default:
		return false
	}
}

// compareNumeric 数值比较
func (s *workflowEngineService) compareNumeric(actual interface{}, expected interface{}, operator string) bool {
	var actualNum, expectedNum float64

	switch v := actual.(type) {
	case int:
		actualNum = float64(v)
	case int64:
		actualNum = float64(v)
	case float64:
		actualNum = v
	case float32:
		actualNum = float64(v)
	default:
		return false
	}

	switch v := expected.(type) {
	case int:
		expectedNum = float64(v)
	case int64:
		expectedNum = float64(v)
	case float64:
		expectedNum = v
	case float32:
		expectedNum = float64(v)
	default:
		return false
	}

	switch operator {
	case ">":
		return actualNum > expectedNum
	case ">=":
		return actualNum >= expectedNum
	case "<":
		return actualNum < expectedNum
	case "<=":
		return actualNum <= expectedNum
	default:
		return false
	}
}

// executeActions 执行动作
func (s *workflowEngineService) executeActions(ctx context.Context, actions []models.TransitionAction, orderID string, metadata map[string]interface{}) {
	for _, action := range actions {
		switch action.Type {
		case "update_field":
			// 更新字段
			if action.Field != "" && action.Value != nil {
				_ = s.orderRepo.Update(ctx, orderID, map[string]interface{}{
					action.Field: action.Value,
					"updated_at": time.Now().Unix(),
				})
			}
		case "notify":
			// TODO: 发送通知
		case "webhook":
			// TODO: 调用webhook
		}
	}
}
