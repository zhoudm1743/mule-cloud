package endpoint

import (
	"context"
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/services"

	"github.com/go-kit/kit/endpoint"
)

// GetStyleEndpoint 获取款式端点
func GetStyleEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleListRequest)
		style, err := svc.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.StyleResponse{Style: style}, nil
	}
}

// GetAllStylesEndpoint 获取所有款式端点（不分页）
func GetAllStylesEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleListRequest)
		styles, err := svc.GetAll(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.StyleListResponse{Styles: styles, Total: int64(len(styles))}, nil
	}
}

// ListStylesEndpoint 款式列表端点（分页）
func ListStylesEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleListRequest)
		styles, total, err := svc.List(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.StyleListResponse{
			Styles: styles,
			Total:  total,
		}, nil
	}
}

// CreateStyleEndpoint 创建款式端点
func CreateStyleEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleCreateRequest)
		style, err := svc.Create(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.StyleResponse{Style: style}, nil
	}
}

// UpdateStyleEndpoint 更新款式端点
func UpdateStyleEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleUpdateRequest)
		style, err := svc.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.StyleResponse{Style: style}, nil
	}
}

// DeleteStyleEndpoint 删除款式端点
func DeleteStyleEndpoint(svc services.IStyleService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.StyleListRequest)
		err := svc.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
