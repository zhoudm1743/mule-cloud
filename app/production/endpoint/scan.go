package endpoint

import (
	"context"

	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/services"

	"github.com/go-kit/kit/endpoint"
)

// ParseScanCodeEndpoint 扫码解析端点
func ParseScanCodeEndpoint(s services.IScanService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ScanCodeRequest)
		resp, err := s.ParseScanCode(ctx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
