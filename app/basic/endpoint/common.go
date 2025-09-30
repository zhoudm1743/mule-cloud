package endpoint

import (
	"context"
	"mule-cloud/app/basic/services"

	"github.com/go-kit/kit/endpoint"
)

type CommonRequest struct{}

type CommonResponse struct {
	Status string `json:"status"`
}

func HealthEndpoint(svc services.ICommonService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := svc.Health()
		return resp, err
	}
}
