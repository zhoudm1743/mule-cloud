package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetOrderTypeEndpoint 获取订单类型端点
func GetOrderTypeEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeListRequest)
		orderType, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.OrderTypeResponse{OrderType: orderType}, nil
	}
}

// GetAllOrderTypesEndpoint 获取所有订单类型端点（不分页）
func GetAllOrderTypesEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeListRequest)
		orderTypes, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.OrderTypeListResponse{OrderTypes: orderTypes, Total: int64(len(orderTypes))}, nil
	}
}

// ListOrderTypesEndpoint 订单类型列表端点（分页）
func ListOrderTypesEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeListRequest)
		orderTypes, total, err := svc.List(req)
		if err != nil {
			return nil, err
		}
		return dto.OrderTypeListResponse{
			OrderTypes: orderTypes,
			Total:      total,
		}, nil
	}
}

// CreateOrderTypeEndpoint 创建订单类型端点
func CreateOrderTypeEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeCreateRequest)
		orderType, err := svc.Create(req)
		if err != nil {
			return nil, err
		}
		return dto.OrderTypeResponse{OrderType: orderType}, nil
	}
}

// UpdateOrderTypeEndpoint 更新订单类型端点
func UpdateOrderTypeEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeUpdateRequest)
		orderType, err := svc.Update(req)
		if err != nil {
			return nil, err
		}
		return dto.OrderTypeResponse{OrderType: orderType}, nil
	}
}

// DeleteOrderTypeEndpoint 删除订单类型端点
func DeleteOrderTypeEndpoint(svc services.IOrderTypeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OrderTypeListRequest)
		err := svc.Delete(req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
