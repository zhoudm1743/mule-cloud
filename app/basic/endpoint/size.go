package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetSizeEndpoint 获取尺寸端点
func GetSizeEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeGetRequest)
		size, err := svc.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.SizeResponse{Size: size}, nil
	}
}

// GetAllSizesEndpoint 获取所有尺寸端点（不分页）
func GetAllSizesEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeListRequest)
		sizes, err := svc.GetAll(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.SizeListResponse{
			Sizes: sizes,
			Total: int64(len(sizes)),
		}, nil
	}
}

// ListSizesEndpoint 尺寸列表端点（分页）
func ListSizesEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeListRequest)
		sizes, total, err := svc.List(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.SizeListResponse{
			Sizes: sizes,
			Total: total,
		}, nil
	}
}

// CreateSizeEndpoint 创建尺寸端点
func CreateSizeEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeCreateRequest)
		size, err := svc.Create(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.SizeResponse{Size: size}, nil
	}
}

// UpdateSizeEndpoint 更新尺寸端点
func UpdateSizeEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeUpdateRequest)
		size, err := svc.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.SizeResponse{Size: size}, nil
	}
}

// DeleteSizeEndpoint 删除尺寸端点
func DeleteSizeEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SizeGetRequest)
		err := svc.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
