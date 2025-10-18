package endpoint

import (
	"context"

	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/services"

	"github.com/go-kit/kit/endpoint"
)

// MakeCreateReworkEndpoint 创建返工单
func MakeCreateReworkEndpoint(s services.IReworkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dto.ReworkRequest)
		return s.CreateRework(ctx, req)
	}
}

// MakeGetReworkListEndpoint 获取返工列表
func MakeGetReworkListEndpoint(s services.IReworkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dto.ReworkListRequest)
		return s.GetReworkList(ctx, req)
	}
}

// MakeGetReworkEndpoint 获取返工详情
func MakeGetReworkEndpoint(s services.IReworkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return s.GetRework(ctx, id)
	}
}

// MakeCompleteReworkEndpoint 完成返工
func MakeCompleteReworkEndpoint(s services.IReworkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		params := request.(map[string]interface{})
		id := params["id"].(string)
		req := params["req"].(*dto.CompleteReworkRequest)
		return nil, s.CompleteRework(ctx, id, req)
	}
}

// MakeDeleteReworkEndpoint 删除返工记录
func MakeDeleteReworkEndpoint(s services.IReworkService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return nil, s.DeleteRework(ctx, id)
	}
}
