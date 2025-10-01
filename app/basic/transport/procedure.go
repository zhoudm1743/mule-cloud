package transport

import (
	"mule-cloud/app/basic/dto"
	"mule-cloud/app/basic/endpoint"
	"mule-cloud/app/basic/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// GetProcedureHandler 获取工序处理器
func GetProcedureHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetProcedureEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetAllProceduresHandler 获取所有工序处理器（不分页）
func GetAllProceduresHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.GetAllProceduresEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// ListProceduresHandler 工序列表处理器（分页）
func ListProceduresHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureListRequest
		if err := c.ShouldBind(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ListProceduresEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// CreateProcedureHandler 创建工序处理器
func CreateProcedureHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureCreateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.CreateProcedureEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateProcedureHandler 更新工序处理器
func UpdateProcedureHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureUpdateRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.UpdateProcedureEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// DeleteProcedureHandler 删除工序处理器
func DeleteProcedureHandler(svc services.IProcedureService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ProcedureListRequest
		if err := c.ShouldBindUri(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.DeleteProcedureEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
