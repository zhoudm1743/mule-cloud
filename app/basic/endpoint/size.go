package endpoint

import (
	"context"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

// SizeRequest 尺寸请求
type SizeRequest struct {
	ID string `uri:"id" binding:"required"`
}

// SizeResponse 尺寸响应
type SizeResponse struct {
	Size *services.Size `json:"size"`
}

// SizeListResponse 尺寸列表响应
type SizeListResponse struct {
	Sizes []*services.Size `json:"sizes"`
	Total int              `json:"total"`
}

// GetSizeEndpoint 获取尺寸端点
func GetSizeEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SizeRequest)
		size, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return SizeResponse{Size: size}, nil
	}
}

// GetAllSizesEndpoint 获取所有尺寸端点
func GetAllSizesEndpoint(svc services.ISizeService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sizes, err := svc.GetAll()
		if err != nil {
			return nil, err
		}
		return SizeListResponse{
			Sizes: sizes,
			Total: len(sizes),
		}, nil
	}
}
