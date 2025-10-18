package transport

import (
	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/endpoint"
	"mule-cloud/app/production/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// SubmitReportHandler 工序上报处理器
func SubmitReportHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureReportRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.SubmitReportEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetReportListHandler 上报记录列表处理器
func GetReportListHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ReportListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetReportListEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetReportByIDHandler 获取上报记录详情处理器
func GetReportByIDHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "上报记录ID不能为空")
			return
		}

		ep := endpoint.GetReportByIDEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteReportHandler 删除上报记录处理器
func DeleteReportHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "上报记录ID不能为空")
			return
		}

		ep := endpoint.DeleteReportEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetOrderProgressHandler 获取订单进度处理器
func GetOrderProgressHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		if orderID == "" {
			response.Error(c, "订单ID不能为空")
			return
		}

		ep := endpoint.GetOrderProgressEndpoint(svc)
		resp, err := ep(c.Request.Context(), orderID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetSalaryHandler 获取工资统计处理器
func GetSalaryHandler(svc services.IReportService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalaryRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetSalaryEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
