package transport

import (
	"context"

	"github.com/gin-gonic/gin"
	"mule-cloud/app/order/endpoint"
	"mule-cloud/app/order/services"
	"mule-cloud/core/response"
)

// GetWorkflowTemplatesHandler 获取工作流模板
func GetWorkflowTemplatesHandler(s *services.WorkflowTemplateService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ep := endpoint.GetWorkflowTemplatesEndpoint(s)
		resp, err := ep(context.Background(), nil)
		if err != nil {
			response.Error(c, err.Error())
			return
		}
		response.Success(c, resp)
	}
}

