package transport

import (
	"mule-cloud/app/workflow/dto"
	"mule-cloud/app/workflow/services"
	corecontext "mule-cloud/core/context"
	"mule-cloud/core/response"
	"mule-cloud/core/workflow"

	"github.com/gin-gonic/gin"
)

// GetWorkflowDefinitionHandler 获取工作流定义
func GetWorkflowDefinitionHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		definition := svc.GetWorkflowDefinition()
		response.Success(c, definition)
	}
}

// GetMermaidDiagramHandler 获取 Mermaid 流程图
func GetMermaidDiagramHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		diagram := svc.GetMermaidDiagram()
		response.Success(c, dto.MermaidDiagramResponse{
			Diagram: diagram,
		})
	}
}

// GetTransitionRulesHandler 获取所有转换规则
func GetTransitionRulesHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		rules := svc.GetTransitionRules()
		response.Success(c, gin.H{
			"rules": rules,
		})
	}
}

// GetOrderStatusHandler 获取订单当前状态
func GetOrderStatusHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")

		status, err := svc.GetOrderStatus(c.Request.Context(), orderID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, dto.OrderStatusResponse{
			OrderID:    orderID,
			Status:     int(status),
			StatusName: workflow.GetStatusName(status),
		})
	}
}

// GetOrderHistoryHandler 获取订单状态历史
func GetOrderHistoryHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		limit := int64(20)

		if limitStr, ok := c.GetQuery("limit"); ok {
			var limitInt int
			if err := c.ShouldBindQuery(&struct{ Limit int }{limitInt}); err == nil && limitStr != "" {
				limit = int64(limitInt)
			}
		}

		history, err := svc.GetOrderHistory(c.Request.Context(), orderID, limit)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"order_id": orderID,
			"history":  history,
		})
	}
}

// GetRollbackHistoryHandler 获取回滚历史
func GetRollbackHistoryHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		limit := int64(10)

		rollbacks, err := svc.GetRollbackHistory(c.Request.Context(), orderID, limit)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"order_id":  orderID,
			"rollbacks": rollbacks,
		})
	}
}

// TransitionOrderHandler 执行状态转换
func TransitionOrderHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.TransitionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}

		// 获取当前用户
		operator := corecontext.GetUsername(c.Request.Context())
		// TODO: 从上下文获取用户角色
		userRole := "admin" // 临时使用admin角色

		err := svc.TransitionOrder(
			c.Request.Context(),
			req.OrderID,
			workflow.OrderEvent(req.Event),
			operator,
			userRole,
			req.Reason,
			req.Metadata,
		)

		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"message": "状态转换成功",
		})
	}
}

// RollbackOrderHandler 回滚订单状态
func RollbackOrderHandler(svc services.IWorkflowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RollbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.BadRequest(c, "参数错误: "+err.Error())
			return
		}

		// 获取当前用户
		operator := corecontext.GetUsername(c.Request.Context())

		err := svc.RollbackOrder(c.Request.Context(), req.OrderID, operator, req.Reason)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"message": "状态回滚成功",
		})
	}
}
