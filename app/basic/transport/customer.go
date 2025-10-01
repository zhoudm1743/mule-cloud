package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetCustomerHandler 获取客户处理器
func GetCustomerHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetCustomerEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllCustomersHandler 获取所有客户处理器（不分页）
func GetAllCustomersHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllCustomersEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListCustomersHandler 客户列表处理器（分页）
func ListCustomersHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListCustomersEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateCustomerHandler 创建客户处理器
func CreateCustomerHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateCustomerEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateCustomerHandler 更新客户处理器
func UpdateCustomerHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerUpdateRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateCustomerEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteCustomerHandler 删除客户处理器
func DeleteCustomerHandler(svc services.ICustomerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CustomerListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteCustomerEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
