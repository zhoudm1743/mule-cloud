package endpoint

import (
	"context"

	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/services"

	"github.com/go-kit/kit/endpoint"
)

// SubmitReportEndpoint 工序上报端点
func SubmitReportEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ProcedureReportRequest)
		resp, err := s.SubmitReport(ctx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// GetReportListEndpoint 上报记录列表端点
func GetReportListEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.ReportListRequest)
		resp, err := s.GetReportList(ctx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// GetReportByIDEndpoint 获取上报记录详情端点
func GetReportByIDEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		report, err := s.GetReportByID(ctx, id)
		if err != nil {
			return nil, err
		}
		return report, nil
	}
}

// DeleteReportEndpoint 删除上报记录端点
func DeleteReportEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id := request.(string)
		err := s.DeleteReport(ctx, id)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{"message": "删除成功"}, nil
	}
}

// GetOrderProgressEndpoint 获取订单进度端点
func GetOrderProgressEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		orderID := request.(string)
		resp, err := s.GetOrderProgress(ctx, orderID)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// GetSalaryEndpoint 获取工资统计端点
func GetSalaryEndpoint(s services.IReportService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SalaryRequest)
		resp, err := s.GetSalary(ctx, &req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
