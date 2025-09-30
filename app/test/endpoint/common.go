package endpoint

import (
	"context"
	"mule-cloud/app/test/services"

	"github.com/go-kit/kit/endpoint"
)

type CommonRequest struct{}

type CommonResponse struct {
	Status string `json:"status"`
}

func HealthEndpoint(svc services.ICommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// req := request.(CommonRequest)
		resp, err := svc.Health()
		if err != nil {
			return nil, err
		}
		return CommonResponse{Status: resp}, nil
	}
}
