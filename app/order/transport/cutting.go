package transport

import (
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/endpoint"
	"mule-cloud/app/order/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// ==================== 裁剪任务 Handlers ====================

// CreateCuttingTaskHandler 创建裁剪任务处理器
func CreateCuttingTaskHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingTaskCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateCuttingTaskEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListCuttingTasksHandler 裁剪任务列表处理器
func ListCuttingTasksHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingTaskListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListCuttingTasksEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetCuttingTaskHandler 获取裁剪任务详情处理器
func GetCuttingTaskHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "任务ID不能为空")
			return
		}

		ep := endpoint.GetCuttingTaskEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetCuttingTaskByOrderHandler 根据订单ID获取裁剪任务处理器
func GetCuttingTaskByOrderHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		if orderID == "" {
			response.Error(c, "订单ID不能为空")
			return
		}

		ep := endpoint.GetCuttingTaskByOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), orderID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ==================== 裁剪批次 Handlers ====================

// CreateCuttingBatchHandler 创建裁剪批次处理器
func CreateCuttingBatchHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingBatchCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateCuttingBatchEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// BulkCreateCuttingBatchHandler 批量创建裁剪批次处理器
func BulkCreateCuttingBatchHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingBatchBulkCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.BulkCreateCuttingBatchEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListCuttingBatchesHandler 裁剪批次列表处理器
func ListCuttingBatchesHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingBatchListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListCuttingBatchesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetCuttingBatchHandler 获取裁剪批次详情处理器
func GetCuttingBatchHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "批次ID不能为空")
			return
		}

		ep := endpoint.GetCuttingBatchEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteCuttingBatchHandler 删除裁剪批次处理器
func DeleteCuttingBatchHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "批次ID不能为空")
			return
		}

		ep := endpoint.DeleteCuttingBatchEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ClearTaskBatchesHandler 清空任务批次处理器
func ClearTaskBatchesHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		taskID := c.Param("taskId")
		if taskID == "" {
			response.Error(c, "任务ID不能为空")
			return
		}

		ep := endpoint.ClearTaskBatchesEndpoint(svc)
		resp, err := ep(c.Request.Context(), taskID)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// PrintCuttingBatchHandler 打印裁剪批次处理器
func PrintCuttingBatchHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "批次ID不能为空")
			return
		}

		ep := endpoint.PrintCuttingBatchEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// BatchPrintCuttingBatchesHandler 批量打印裁剪批次处理器
func BatchPrintCuttingBatchesHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BatchPrintRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.BatchPrintCuttingBatchesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ==================== 裁片监控 Handlers ====================

// ListCuttingPiecesHandler 裁片监控列表处理器
func ListCuttingPiecesHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CuttingPieceListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListCuttingPiecesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetCuttingPieceHandler 获取裁片监控详情处理器
func GetCuttingPieceHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "裁片ID不能为空")
			return
		}

		ep := endpoint.GetCuttingPieceEndpoint(svc)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateCuttingPieceProgressHandler 更新裁片进度处理器
func UpdateCuttingPieceProgressHandler(svc services.ICuttingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			response.Error(c, "裁片ID不能为空")
			return
		}

		var req dto.CuttingPieceProgressRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateCuttingPieceProgressEndpoint(svc)
		resp, err := ep(c.Request.Context(), map[string]interface{}{
			"id":       id,
			"progress": req.Progress,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
