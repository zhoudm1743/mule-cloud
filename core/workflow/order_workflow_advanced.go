package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mule-cloud/internal/repository"
)

// TransitionCondition 状态转换条件函数
type TransitionCondition func(ctx context.Context, orderID string, metadata map[string]interface{}) (bool, string)

// TransitionRule 增强的状态转换规则（支持条件）
type TransitionRule struct {
	From      OrderStatus
	Event     OrderEvent
	To        OrderStatus
	Condition TransitionCondition // 转换条件（可选）
	RequireRole string            // 需要的角色（可选）
}

// RollbackRecord 回滚记录
type RollbackRecord struct {
	OrderID       string                 `json:"order_id"`
	RollbackFrom  OrderStatus            `json:"rollback_from"`
	RollbackTo    OrderStatus            `json:"rollback_to"`
	OriginalEvent OrderEvent             `json:"original_event"`
	Reason        string                 `json:"reason"`
	Operator      string                 `json:"operator"`
	Timestamp     int64                  `json:"timestamp"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// 增强的状态转换规则表（支持条件）
var advancedTransitions = []TransitionRule{
	// 基础转换
	{From: StatusDraft, Event: EventSubmitOrder, To: StatusOrdered, Condition: nil},
	{From: StatusOrdered, Event: EventStartCutting, To: StatusProduction, Condition: nil},
	{From: StatusOrdered, Event: EventStartProduction, To: StatusProduction, Condition: nil},
	{From: StatusProduction, Event: EventUpdateProgress, To: StatusProduction, Condition: nil},
	
	// 完成订单 - 需要进度达到100%
	{
		From:  StatusProduction,
		Event: EventComplete,
		To:    StatusCompleted,
		Condition: func(ctx context.Context, orderID string, metadata map[string]interface{}) (bool, string) {
			// 检查进度是否达到100%
			if progress, ok := metadata["progress"].(float64); ok {
				if progress >= 1.0 {
					return true, ""
				}
				return false, fmt.Sprintf("进度不足：当前%.1f%%，需要100%%", progress*100)
			}
			
			// 如果没有传入进度，从数据库查询
			orderRepo := repository.NewOrderRepository()
			order, err := orderRepo.Get(ctx, orderID)
			if err != nil {
				return false, "无法获取订单信息"
			}
			
			if order.Progress >= 1.0 {
				return true, ""
			}
			return false, fmt.Sprintf("进度不足：当前%.1f%%，需要100%%", order.Progress*100)
		},
	},
	
	// 取消订单 - 需要管理员权限或特定角色
	{
		From:        StatusDraft,
		Event:       EventCancel,
		To:          StatusCancelled,
		RequireRole: "admin",
	},
	{
		From:        StatusOrdered,
		Event:       EventCancel,
		To:          StatusCancelled,
		RequireRole: "admin",
	},
	{
		From:        StatusProduction,
		Event:       EventCancel,
		To:          StatusCancelled,
		RequireRole: "admin",
	},
}

// CanTransitionWithCondition 检查是否可以进行状态转换（带条件检查）
func (w *OrderWorkflow) CanTransitionWithCondition(
	ctx context.Context,
	orderID string,
	currentStatus OrderStatus,
	event OrderEvent,
	userRole string,
	metadata map[string]interface{},
) (bool, string) {
	for _, rule := range advancedTransitions {
		if rule.From == currentStatus && rule.Event == event {
			// 检查角色权限
			if rule.RequireRole != "" && rule.RequireRole != userRole {
				return false, fmt.Sprintf("需要角色: %s", rule.RequireRole)
			}
			
			// 检查条件
			if rule.Condition != nil {
				canTransit, reason := rule.Condition(ctx, orderID, metadata)
				if !canTransit {
					return false, reason
				}
			}
			
			return true, ""
		}
	}
	
	return false, fmt.Sprintf("无效的状态转换: %d -> %s", currentStatus, event)
}

// TransitionToAdvanced 高级状态转换（支持条件和权限）
func (w *OrderWorkflow) TransitionToAdvanced(
	ctx context.Context,
	orderID string,
	event OrderEvent,
	operator string,
	userRole string,
	reason string,
	metadata map[string]interface{},
) error {
	// 获取当前状态
	currentStatus, err := w.GetCurrentStatus(ctx, orderID)
	if err != nil {
		return err
	}

	// 检查是否可以转换（带条件）
	canTransit, errMsg := w.CanTransitionWithCondition(ctx, orderID, currentStatus, event, userRole, metadata)
	if !canTransit {
		return fmt.Errorf("状态转换失败: %s", errMsg)
	}

	// 获取下一个状态
	var nextStatus OrderStatus
	for _, rule := range advancedTransitions {
		if rule.From == currentStatus && rule.Event == event {
			nextStatus = rule.To
			break
		}
	}

	// 记录状态历史（用于回滚）
	history := StateHistory{
		OrderID:   orderID,
		FromState: currentStatus,
		ToState:   nextStatus,
		Event:     event,
		Reason:    reason,
		Operator:  operator,
		Timestamp: time.Now().Unix(),
		Metadata:  metadata,
	}
	w.saveHistory(ctx, history)

	// 更新数据库
	err = w.orderRepo.Update(ctx, orderID, map[string]interface{}{
		"$set": map[string]interface{}{
			"status":     int(nextStatus),
			"updated_at": time.Now().Unix(),
		},
	})
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %v", err)
	}

	// 更新Redis缓存
	w.syncStatusToRedis(ctx, orderID, nextStatus)

	return nil
}

// RollbackLastTransition 回滚最后一次状态转换
func (w *OrderWorkflow) RollbackLastTransition(
	ctx context.Context,
	orderID string,
	operator string,
	reason string,
) error {
	// 获取最近的状态历史
	histories, err := w.GetHistory(ctx, orderID, 2)
	if err != nil || len(histories) < 1 {
		return fmt.Errorf("没有可回滚的历史记录")
	}

	lastHistory := histories[0]
	
	// 检查是否可以回滚（已完成和已取消的订单不允许回滚）
	if lastHistory.ToState == StatusCompleted || lastHistory.ToState == StatusCancelled {
		return fmt.Errorf("订单状态为 %s，不允许回滚", GetStatusName(lastHistory.ToState))
	}

	// 创建回滚记录
	rollback := RollbackRecord{
		OrderID:       orderID,
		RollbackFrom:  lastHistory.ToState,
		RollbackTo:    lastHistory.FromState,
		OriginalEvent: lastHistory.Event,
		Reason:        reason,
		Operator:      operator,
		Timestamp:     time.Now().Unix(),
		Metadata: map[string]interface{}{
			"original_transition": lastHistory,
		},
	}

	// 保存回滚记录
	w.saveRollbackRecord(ctx, rollback)

	// 更新数据库
	err = w.orderRepo.Update(ctx, orderID, map[string]interface{}{
		"$set": map[string]interface{}{
			"status":     int(lastHistory.FromState),
			"updated_at": time.Now().Unix(),
		},
	})
	if err != nil {
		return fmt.Errorf("回滚状态失败: %v", err)
	}

	// 更新Redis缓存
	w.syncStatusToRedis(ctx, orderID, lastHistory.FromState)

	// 记录回滚事件到历史
	rollbackHistory := StateHistory{
		OrderID:   orderID,
		FromState: lastHistory.ToState,
		ToState:   lastHistory.FromState,
		Event:     "rollback",
		Reason:    fmt.Sprintf("回滚: %s", reason),
		Operator:  operator,
		Timestamp: time.Now().Unix(),
		Metadata: map[string]interface{}{
			"is_rollback":        true,
			"rollback_from":      lastHistory.ToState,
			"original_event":     lastHistory.Event,
		},
	}
	w.saveHistory(ctx, rollbackHistory)

	return nil
}

// saveRollbackRecord 保存回滚记录到Redis
func (w *OrderWorkflow) saveRollbackRecord(ctx context.Context, rollback RollbackRecord) {
	rollbackKey := fmt.Sprintf("order:rollback:%s", rollback.OrderID)
	rollbackJSON, _ := json.Marshal(rollback)
	
	// 使用List保存回滚记录
	_ = w.redis.Client().LPush(ctx, rollbackKey, string(rollbackJSON)).Err()
	
	// 只保留最近50条回滚记录
	_ = w.redis.Client().LTrim(ctx, rollbackKey, 0, 49).Err()
	
	// 设置过期时间为90天
	_ = w.redis.Expire(ctx, rollbackKey, 90*24*time.Hour)
}

// GetRollbackHistory 获取回滚历史
func (w *OrderWorkflow) GetRollbackHistory(ctx context.Context, orderID string, limit int64) ([]RollbackRecord, error) {
	rollbackKey := fmt.Sprintf("order:rollback:%s", orderID)
	
	if limit <= 0 {
		limit = 10
	}
	
	results, err := w.redis.Client().LRange(ctx, rollbackKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	rollbacks := make([]RollbackRecord, 0, len(results))
	for _, result := range results {
		var rollback RollbackRecord
		if err := json.Unmarshal([]byte(result), &rollback); err == nil {
			rollbacks = append(rollbacks, rollback)
		}
	}

	return rollbacks, nil
}

// GetWorkflowDefinition 获取工作流定义（用于可视化）
func GetWorkflowDefinition() map[string]interface{} {
	// 状态列表
	states := []map[string]interface{}{
		{"id": 0, "name": "草稿", "type": "start", "color": "#909399"},
		{"id": 1, "name": "已下单", "type": "normal", "color": "#409EFF"},
		{"id": 2, "name": "生产中", "type": "normal", "color": "#E6A23C"},
		{"id": 3, "name": "已完成", "type": "end", "color": "#67C23A"},
		{"id": 4, "name": "已取消", "type": "end", "color": "#F56C6C"},
	}

	// 事件列表
	events := []map[string]interface{}{
		{"name": "submit_order", "label": "提交订单"},
		{"name": "start_cutting", "label": "开始裁剪"},
		{"name": "start_production", "label": "开始生产"},
		{"name": "update_progress", "label": "更新进度"},
		{"name": "complete", "label": "完成", "requireCondition": true},
		{"name": "cancel", "label": "取消", "requireRole": "admin"},
	}

	// 转换规则
	transitions := []map[string]interface{}{}
	for _, rule := range advancedTransitions {
		transition := map[string]interface{}{
			"from":  int(rule.From),
			"to":    int(rule.To),
			"event": string(rule.Event),
		}
		
		if rule.Condition != nil {
			transition["hasCondition"] = true
			transition["conditionDesc"] = "需要满足特定条件"
		}
		
		if rule.RequireRole != "" {
			transition["requireRole"] = rule.RequireRole
			transition["roleDesc"] = fmt.Sprintf("需要 %s 角色", rule.RequireRole)
		}
		
		transitions = append(transitions, transition)
	}

	return map[string]interface{}{
		"states":      states,
		"events":      events,
		"transitions": transitions,
	}
}

// GenerateMermaidDiagram 生成 Mermaid 流程图
func GenerateMermaidDiagram() string {
	diagram := "graph LR\n"
	diagram += "    Start([开始])\n"
	diagram += "    Draft[草稿]\n"
	diagram += "    Ordered[已下单]\n"
	diagram += "    Production[生产中]\n"
	diagram += "    Completed([已完成])\n"
	diagram += "    Cancelled([已取消])\n"
	diagram += "\n"
	diagram += "    Start --> Draft\n"
	diagram += "    Draft -->|提交订单| Ordered\n"
	diagram += "    Ordered -->|开始裁剪| Production\n"
	diagram += "    Ordered -->|开始生产| Production\n"
	diagram += "    Production -->|更新进度| Production\n"
	diagram += "    Production -->|完成<br/>进度>=100%| Completed\n"
	diagram += "    Draft -.->|取消<br/>需要admin| Cancelled\n"
	diagram += "    Ordered -.->|取消<br/>需要admin| Cancelled\n"
	diagram += "    Production -.->|取消<br/>需要admin| Cancelled\n"
	diagram += "\n"
	diagram += "    style Draft fill:#909399,color:#fff\n"
	diagram += "    style Ordered fill:#409EFF,color:#fff\n"
	diagram += "    style Production fill:#E6A23C,color:#fff\n"
	diagram += "    style Completed fill:#67C23A,color:#fff\n"
	diagram += "    style Cancelled fill:#F56C6C,color:#fff\n"
	
	return diagram
}

// GetTransitionRules 获取所有转换规则（用于前端展示）
func GetTransitionRules() []map[string]interface{} {
	rules := make([]map[string]interface{}, 0, len(advancedTransitions))
	
	for _, rule := range advancedTransitions {
		ruleMap := map[string]interface{}{
			"from":      int(rule.From),
			"from_name": GetStatusName(rule.From),
			"to":        int(rule.To),
			"to_name":   GetStatusName(rule.To),
			"event":     string(rule.Event),
		}
		
		if rule.Condition != nil {
			ruleMap["has_condition"] = true
		}
		
		if rule.RequireRole != "" {
			ruleMap["require_role"] = rule.RequireRole
		}
		
		rules = append(rules, ruleMap)
	}
	
	return rules
}

