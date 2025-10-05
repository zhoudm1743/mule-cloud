package endpoint

import (
	"context"
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// GetColorEndpoint 获取颜色端点
func GetColorEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorListRequest)
		color, err := svc.Get(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return dto.ColorResponse{Color: color}, nil
	}
}

// GetAllColorsEndpoint 获取所有颜色端点（不分页）
func GetAllColorsEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorListRequest)
		colors, err := svc.GetAll(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.ColorListResponse{Colors: colors, Total: int64(len(colors))}, nil
	}
}

// ListColorsEndpoint 颜色列表端点（分页）
func ListColorsEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorListRequest)
		colors, total, err := svc.List(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.ColorListResponse{
			Colors: colors,
			Total:  total,
		}, nil
	}
}

// CreateColorEndpoint 创建颜色端点
func CreateColorEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorCreateRequest)
		color, err := svc.Create(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.ColorResponse{Color: color}, nil
	}
}

// UpdateColorEndpoint 更新颜色端点
func UpdateColorEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorUpdateRequest)
		color, err := svc.Update(ctx, req)
		if err != nil {
			return nil, err
		}
		return dto.ColorResponse{Color: color}, nil
	}
}

// DeleteColorEndpoint 删除颜色端点
func DeleteColorEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorListRequest)
		err := svc.Delete(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		return map[string]string{"message": "删除成功"}, nil
	}
}
