package endpoint

import (
	"context"
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/services"

	"github.com/go-kit/kit/endpoint"
)

// GetOrderEndpoint 获取订单端点
func GetOrderEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderListRequest)
		order, err := svc.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// ListOrdersEndpoint 订单列表端点（分页）
func ListOrdersEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderListRequest)
		orders, total, err := svc.List(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.OrderListResponse{
			Orders: orders,
			Total:  total,
		}, nil
	}
}

// CreateOrderEndpoint 创建订单端点
func CreateOrderEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderCreateRequest)
		order, err := svc.Create(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// UpdateOrderStyleEndpoint 更新订单款式端点
func UpdateOrderStyleEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderStyleRequest)
		order, err := svc.UpdateStyle(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// UpdateOrderProcedureEndpoint 更新订单工序端点
func UpdateOrderProcedureEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderProcedureRequest)
		order, err := svc.UpdateProcedure(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// UpdateOrderEndpoint 更新订单端点
func UpdateOrderEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderUpdateRequest)
		order, err := svc.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// CopyOrderEndpoint 复制订单端点
func CopyOrderEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderCopyRequest)
		order, err := svc.Copy(ctx, req.ID, req.IsRelated, req.RelationType, req.RelationRemark)
		if err != nil {
			return nil, err
		}
		return dto.OrderResponse{Order: order}, nil
	}
}

// DeleteOrderEndpoint 删除订单端点
func DeleteOrderEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderListRequest)
		err := svc.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}

// TransitionOrderWorkflowEndpoint 执行订单工作流状态转换端点
func TransitionOrderWorkflowEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderWorkflowTransitionRequest)
		err := svc.TransitionWorkflowState(ctx, req)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "状态转换成功"}, nil
	}
}

// GetOrderWorkflowStateEndpoint 获取订单工作流状态端点
func GetOrderWorkflowStateEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderListRequest)
		instance, err := svc.GetWorkflowState(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.OrderWorkflowStateResponse{Instance: instance}, nil
	}
}

// GetOrderAvailableTransitionsEndpoint 获取订单可用状态转换端点
func GetOrderAvailableTransitionsEndpoint(svc services.IOrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderListRequest)
		transitions, err := svc.GetAvailableTransitions(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.OrderWorkflowTransitionsResponse{Transitions: transitions}, nil
	}
}
