package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"mule-cloud/core/cache"
	"mule-cloud/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// OrderStatus 订单状态
type OrderStatus int

const (
	StatusDraft      OrderStatus = 0 // 草稿
	StatusOrdered    OrderStatus = 1 // 已下单
	StatusProduction OrderStatus = 2 // 生产中
	StatusCompleted  OrderStatus = 3 // 已完成
	StatusCancelled  OrderStatus = 4 // 已取消
)

// OrderEvent 订单事件
type OrderEvent string

const (
	EventSubmitOrder     OrderEvent = "submit_order"     // 提交订单
	EventStartCutting    OrderEvent = "start_cutting"    // 开始裁剪
	EventStartProduction OrderEvent = "start_production" // 开始生产
	EventUpdateProgress  OrderEvent = "update_progress"  // 更新进度
	EventComplete        OrderEvent = "complete"         // 完成
	EventCancel          OrderEvent = "cancel"           // 取消
)

// StatusTransition 状态转换规则
type StatusTransition struct {
	From  OrderStatus
	Event OrderEvent
	To    OrderStatus
}

// 订单状态转换规则表
var orderTransitions = []StatusTransition{
	{StatusDraft, EventSubmitOrder, StatusOrdered},
	{StatusOrdered, EventStartCutting, StatusProduction},
	{StatusOrdered, EventStartProduction, StatusProduction},
	{StatusProduction, EventUpdateProgress, StatusProduction},
	{StatusProduction, EventComplete, StatusCompleted},
	{StatusDraft, EventCancel, StatusCancelled},
	{StatusOrdered, EventCancel, StatusCancelled},
	{StatusProduction, EventCancel, StatusCancelled},
}

// StateHistory 状态历史记录
type StateHistory struct {
	OrderID   string                 `json:"order_id"`
	FromState OrderStatus            `json:"from_state"`
	ToState   OrderStatus            `json:"to_state"`
	Event     OrderEvent             `json:"event"`
	Reason    string                 `json:"reason"`
	Operator  string                 `json:"operator"`
	Timestamp int64                  `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// OrderWorkflow 订单工作流
type OrderWorkflow struct {
	orderRepo repository.OrderRepository
	redis     *cache.RedisInstance
}

// NewOrderWorkflow 创建订单工作流
func NewOrderWorkflow() *OrderWorkflow {
	return &OrderWorkflow{
		orderRepo: repository.NewOrderRepository(),
		redis:     cache.Redis,
	}
}

// getOrderCacheKey 获取订单缓存键
func (w *OrderWorkflow) getOrderCacheKey(orderID string) string {
	return fmt.Sprintf("order:status:%s", orderID)
}

// getOrderHistoryKey 获取订单历史记录键
func (w *OrderWorkflow) getOrderHistoryKey(orderID string) string {
	return fmt.Sprintf("order:history:%s", orderID)
}

// GetCurrentStatus 获取当前订单状态（优先从Redis）
func (w *OrderWorkflow) GetCurrentStatus(ctx context.Context, orderID string) (OrderStatus, error) {
	// 先从Redis获取
	cacheKey := w.getOrderCacheKey(orderID)
	if statusStr, err := w.redis.Get(ctx, cacheKey); err == nil {
		var status OrderStatus
		if err := json.Unmarshal([]byte(statusStr), &status); err == nil {
			return status, nil
		}
	}

	// 从数据库获取
	order, err := w.orderRepo.Get(ctx, orderID)
	if err != nil {
		return 0, fmt.Errorf("获取订单失败: %v", err)
	}

	// 同步到Redis
	status := OrderStatus(order.Status)
	w.syncStatusToRedis(ctx, orderID, status)

	return status, nil
}

// syncStatusToRedis 同步状态到Redis
func (w *OrderWorkflow) syncStatusToRedis(ctx context.Context, orderID string, status OrderStatus) {
	cacheKey := w.getOrderCacheKey(orderID)
	statusJSON, _ := json.Marshal(status)
	_ = w.redis.Set(ctx, cacheKey, string(statusJSON), 24*time.Hour) // 缓存24小时
}

// CanTransition 检查是否可以进行状态转换
func (w *OrderWorkflow) CanTransition(currentStatus OrderStatus, event OrderEvent) bool {
	for _, transition := range orderTransitions {
		if transition.From == currentStatus && transition.Event == event {
			return true
		}
	}
	return false
}

// GetNextStatus 获取下一个状态
func (w *OrderWorkflow) GetNextStatus(currentStatus OrderStatus, event OrderEvent) (OrderStatus, error) {
	for _, transition := range orderTransitions {
		if transition.From == currentStatus && transition.Event == event {
			return transition.To, nil
		}
	}
	return currentStatus, fmt.Errorf("无效的状态转换: %d -> %s", currentStatus, event)
}

// TransitionTo 状态转换
func (w *OrderWorkflow) TransitionTo(
	ctx context.Context,
	orderID string,
	event OrderEvent,
	operator string,
	reason string,
	metadata map[string]interface{},
) error {
	// 获取当前状态
	currentStatus, err := w.GetCurrentStatus(ctx, orderID)
	if err != nil {
		return err
	}

	// 检查是否可以转换
	if !w.CanTransition(currentStatus, event) {
		return fmt.Errorf("不允许的状态转换: 状态[%d] 事件[%s]", currentStatus, event)
	}

	// 获取下一个状态
	nextStatus, err := w.GetNextStatus(currentStatus, event)
	if err != nil {
		return err
	}

	// 记录状态历史
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
	err = w.orderRepo.Update(ctx, orderID, bson.M{
		"$set": bson.M{
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

// saveHistory 保存状态历史到Redis
func (w *OrderWorkflow) saveHistory(ctx context.Context, history StateHistory) {
	historyKey := w.getOrderHistoryKey(history.OrderID)
	historyJSON, _ := json.Marshal(history)

	// 使用List保存历史记录，最新的在最前面
	_ = w.redis.Client().LPush(ctx, historyKey, string(historyJSON)).Err()

	// 只保留最近100条记录
	_ = w.redis.Client().LTrim(ctx, historyKey, 0, 99).Err()

	// 设置过期时间为90天
	_ = w.redis.Expire(ctx, historyKey, 90*24*time.Hour)
}

// GetHistory 获取订单状态历史
func (w *OrderWorkflow) GetHistory(ctx context.Context, orderID string, limit int64) ([]StateHistory, error) {
	historyKey := w.getOrderHistoryKey(orderID)

	if limit <= 0 {
		limit = 10
	}

	results, err := w.redis.Client().LRange(ctx, historyKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	histories := make([]StateHistory, 0, len(results))
	for _, result := range results {
		var history StateHistory
		if err := json.Unmarshal([]byte(result), &history); err == nil {
			histories = append(histories, history)
		}
	}

	return histories, nil
}

// UpdateProgress 更新订单进度（自动判断状态）
func (w *OrderWorkflow) UpdateProgress(
	ctx context.Context,
	orderID string,
	progress float64,
	operator string,
) error {
	// 获取当前状态
	currentStatus, err := w.GetCurrentStatus(ctx, orderID)
	if err != nil {
		return err
	}

	// 更新进度
	err = w.orderRepo.Update(ctx, orderID, bson.M{
		"$set": bson.M{
			"progress":   progress,
			"updated_at": time.Now().Unix(),
		},
	})
	if err != nil {
		return err
	}

	// 根据进度自动转换状态
	if progress >= 1.0 && currentStatus == StatusProduction {
		return w.TransitionTo(ctx, orderID, EventComplete, operator, "工序全部完成", map[string]interface{}{
			"progress": progress,
		})
	} else if progress > 0 && currentStatus == StatusOrdered {
		return w.TransitionTo(ctx, orderID, EventStartProduction, operator, "开始生产", map[string]interface{}{
			"progress": progress,
		})
	}

	return nil
}

// StartCutting 开始裁剪（状态转换）
func (w *OrderWorkflow) StartCutting(ctx context.Context, orderID, operator string) error {
	return w.TransitionTo(ctx, orderID, EventStartCutting, operator, "创建裁剪任务", nil)
}

// StartProduction 开始生产（状态转换）
func (w *OrderWorkflow) StartProduction(ctx context.Context, orderID, operator string, reason string) error {
	return w.TransitionTo(ctx, orderID, EventStartProduction, operator, reason, nil)
}

// CompleteOrder 完成订单
func (w *OrderWorkflow) CompleteOrder(ctx context.Context, orderID, operator, reason string) error {
	return w.TransitionTo(ctx, orderID, EventComplete, operator, reason, nil)
}

// CancelOrder 取消订单
func (w *OrderWorkflow) CancelOrder(ctx context.Context, orderID, operator, reason string) error {
	return w.TransitionTo(ctx, orderID, EventCancel, operator, reason, nil)
}

// InvalidateCache 使订单缓存失效
func (w *OrderWorkflow) InvalidateCache(ctx context.Context, orderID string) error {
	cacheKey := w.getOrderCacheKey(orderID)
	return w.redis.Del(ctx, cacheKey)
}

// GetStatusName 获取状态名称
func GetStatusName(status OrderStatus) string {
	switch status {
	case StatusDraft:
		return "草稿"
	case StatusOrdered:
		return "已下单"
	case StatusProduction:
		return "生产中"
	case StatusCompleted:
		return "已完成"
	case StatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}
