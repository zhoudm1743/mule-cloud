package transport

import (
	"mule-cloud/app/order/dto"
	"mule-cloud/app/order/endpoint"
	"mule-cloud/app/order/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetOrderHandler 获取订单处理器
func GetOrderHandler(svc services.IOrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderListRequest
		if err := c.ShouldBindUri(&req); err != nil {
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
		if err := c.ShouldBind(&req); err != nil {
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
		if err := c.ShouldBindJSON(&req); err != nil {
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
		// 先绑定 URI 参数（ID）
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		// 再绑定 JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
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
		// 先绑定 URI 参数（ID）
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		// 再绑定 JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
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
		// 先绑定 JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		// 再绑定 URI 参数（ID）
		if err := c.ShouldBindUri(&req); err != nil {
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
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
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
		if err := c.ShouldBindUri(&req); err != nil {
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
