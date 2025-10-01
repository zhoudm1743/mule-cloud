package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetOrderTypeHandler 获取订单类型处理器
func GetOrderTypeHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetOrderTypeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllOrderTypesHandler 获取所有订单类型处理器（不分页）
func GetAllOrderTypesHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllOrderTypesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListOrderTypesHandler 订单类型列表处理器（分页）
func ListOrderTypesHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListOrderTypesEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateOrderTypeHandler 创建订单类型处理器
func CreateOrderTypeHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateOrderTypeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateOrderTypeHandler 更新订单类型处理器
func UpdateOrderTypeHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeUpdateRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateOrderTypeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteOrderTypeHandler 删除订单类型处理器
func DeleteOrderTypeHandler(svc services.IOrderTypeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.OrderTypeListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteOrderTypeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
