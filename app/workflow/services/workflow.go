package services

import (
	"context"
	"mule-cloud/core/workflow"
)

// IWorkflowService 工作流服务接口
type IWorkflowService interface {
	// 获取工作流定义
	GetWorkflowDefinition() map[string]interface{}
	
	// 获取 Mermaid 流程图
	GetMermaidDiagram() string
	
	// 获取订单当前状态
	GetOrderStatus(ctx context.Context, orderID string) (workflow.OrderStatus, error)
	
	// 获取订单状态历史
	GetOrderHistory(ctx context.Context, orderID string, limit int64) ([]workflow.StateHistory, error)
	
	// 获取订单回滚历史
	GetRollbackHistory(ctx context.Context, orderID string, limit int64) ([]workflow.RollbackRecord, error)
	
	// 执行状态转换
	TransitionOrder(ctx context.Context, orderID string, event workflow.OrderEvent, operator, userRole, reason string, metadata map[string]interface{}) error
	
	// 回滚状态
	RollbackOrder(ctx context.Context, orderID string, operator, reason string) error
	
	// 获取所有转换规则
	GetTransitionRules() []map[string]interface{}
}

type workflowService struct {
	workflow *workflow.OrderWorkflow
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService() IWorkflowService {
	return &workflowService{
		workflow: workflow.NewOrderWorkflow(),
	}
}

func (s *workflowService) GetWorkflowDefinition() map[string]interface{} {
	return workflow.GetWorkflowDefinition()
}

func (s *workflowService) GetMermaidDiagram() string {
	return workflow.GenerateMermaidDiagram()
}

func (s *workflowService) GetOrderStatus(ctx context.Context, orderID string) (workflow.OrderStatus, error) {
	return s.workflow.GetCurrentStatus(ctx, orderID)
}

func (s *workflowService) GetOrderHistory(ctx context.Context, orderID string, limit int64) ([]workflow.StateHistory, error) {
	return s.workflow.GetHistory(ctx, orderID, limit)
}

func (s *workflowService) GetRollbackHistory(ctx context.Context, orderID string, limit int64) ([]workflow.RollbackRecord, error) {
	return s.workflow.GetRollbackHistory(ctx, orderID, limit)
}

func (s *workflowService) TransitionOrder(
	ctx context.Context,
	orderID string,
	event workflow.OrderEvent,
	operator, userRole, reason string,
	metadata map[string]interface{},
) error {
	return s.workflow.TransitionToAdvanced(ctx, orderID, event, operator, userRole, reason, metadata)
}

func (s *workflowService) RollbackOrder(ctx context.Context, orderID string, operator, reason string) error {
	return s.workflow.RollbackLastTransition(ctx, orderID, operator, reason)
}

func (s *workflowService) GetTransitionRules() []map[string]interface{} {
	return workflow.GetTransitionRules()
}

