package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetSalesmanHandler 获取业务员处理器
func GetSalesmanHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetSalesmanEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllSalesmansHandler 获取所有业务员处理器（不分页）
func GetAllSalesmansHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllSalesmansEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListSalesmansHandler 业务员列表处理器（分页）
func ListSalesmansHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListSalesmansEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateSalesmanHandler 创建业务员处理器
func CreateSalesmanHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanCreateRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateSalesmanEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateSalesmanHandler 更新业务员处理器
func UpdateSalesmanHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanUpdateRequest
		// 先绑定 JSON body（包含 required 字段）
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateSalesmanEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteSalesmanHandler 删除业务员处理器
func DeleteSalesmanHandler(svc services.ISalesmanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SalesmanListRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteSalesmanEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
