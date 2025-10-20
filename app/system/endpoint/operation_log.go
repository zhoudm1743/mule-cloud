package endpoint

import (
	"context"
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/services"
)

// ListOperationLogsEndpoint 列出操作日志的 Endpoint
func ListOperationLogsEndpoint(svc services.OperationLogService) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OperationLogListRequest)
		return svc.List(ctx, req)
	}
}

// GetOperationLogEndpoint 获取操作日志详情的 Endpoint
func GetOperationLogEndpoint(svc services.OperationLogService) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OperationLogDetailRequest)
		return svc.GetByID(ctx, req.ID)
	}
}

// StatsOperationLogsEndpoint 操作日志统计的 Endpoint
func StatsOperationLogsEndpoint(svc services.OperationLogService) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.OperationLogStatsRequest)
		return svc.Stats(ctx, req)
	}
}
