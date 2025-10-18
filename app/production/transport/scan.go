package transport

import (
	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/endpoint"
	"mule-cloud/app/production/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// ParseScanCodeHandler 扫码解析处理器
func ParseScanCodeHandler(svc services.IScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ScanCodeRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.ParseScanCodeEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
