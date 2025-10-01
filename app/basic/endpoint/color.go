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
		color, err := svc.Get(req.ID)
		if err != nil {
			return nil, err
		}
		return dto.ColorResponse{Color: color}, nil
	}
}

func GetAllColorsEndpoint(svc services.IColorService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ColorListRequest)
		colors, err := svc.GetAll(req)
		if err != nil {
			return nil, err
		}
		return dto.ColorListResponse{Colors: colors}, nil
	}
}
