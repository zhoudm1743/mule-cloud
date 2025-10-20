package transport

import (
	"mule-cloud/app/system/dto"
	"mule-cloud/app/system/endpoint"
	"mule-cloud/app/system/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// ListOperationLogsHandler 操作日志列表处理器
func ListOperationLogsHandler(svc services.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OperationLogListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListOperationLogsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetOperationLogHandler 操作日志详情处理器
func GetOperationLogHandler(svc services.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OperationLogDetailRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetOperationLogEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// StatsOperationLogsHandler 操作日志统计处理器
func StatsOperationLogsHandler(svc services.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OperationLogStatsRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.StatsOperationLogsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
