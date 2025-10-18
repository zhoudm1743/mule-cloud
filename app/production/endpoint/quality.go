package endpoint

import (
	"context"

	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/services"

	"github.com/go-kit/kit/endpoint"
)

// MakeSubmitInspectionEndpoint 提交质检
func MakeSubmitInspectionEndpoint(s services.IQualityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dto.InspectionRequest)
		return s.SubmitInspection(ctx, req)
	}
}

// MakeGetInspectionListEndpoint 获取质检列表
func MakeGetInspectionListEndpoint(s services.IQualityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*dto.InspectionListRequest)
		return s.GetInspectionList(ctx, req)
	}
}

// MakeGetInspectionEndpoint 获取质检详情
func MakeGetInspectionEndpoint(s services.IQualityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return s.GetInspection(ctx, id)
	}
}

// MakeDeleteInspectionEndpoint 删除质检记录
func MakeDeleteInspectionEndpoint(s services.IQualityService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		return nil, s.DeleteInspection(ctx, id)
	}
}

