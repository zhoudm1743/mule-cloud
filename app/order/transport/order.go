package transport

import (
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/endpoint"
	"mule-cloud/app/order/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetOrderHandler 获取订单处理器
func GetOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListOrdersHandler 订单列表处理器（分页）
func ListOrdersHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListOrdersEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateOrderHandler 创建订单处理器
func CreateOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateOrderStyleHandler 更新订单款式处理器
func UpdateOrderStyleHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderStyleRequest
		// 使用统一的绑定方法，自动处理 URI 和 Body 参数
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateOrderStyleEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateOrderProcedureHandler 更新订单工序处理器
func UpdateOrderProcedureHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderProcedureRequest
		// 使用统一的绑定方法，自动处理 URI 和 Body 参数
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateOrderProcedureEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateOrderHandler 更新订单处理器
func UpdateOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderUpdateRequest
		// 使用统一的绑定方法，自动处理 URI 和 Body 参数
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CopyOrderHandler 复制订单处理器
func CopyOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderCopyRequest
		// 绑定URI参数（订单ID）
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// 绑定Body参数（关联信息）
		if err := binding.BindAll(c, &req); err != nil {
			// Body参数可选，不报错，使用默认值
		}

		ep := endpoint.CopyOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteOrderHandler 删除订单处理器
func DeleteOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteOrderEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// TransitionOrderWorkflowHandler 执行订单工作流状态转换
func TransitionOrderWorkflowHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderWorkflowTransitionRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.TransitionOrderWorkflowEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetOrderWorkflowStateHandler 获取订单工作流状态
func GetOrderWorkflowStateHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest // 复用，只需要ID
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetOrderWorkflowStateEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetOrderAvailableTransitionsHandler 获取订单可用的状态转换
func GetOrderAvailableTransitionsHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest // 复用，只需要ID
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetOrderAvailableTransitionsEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
