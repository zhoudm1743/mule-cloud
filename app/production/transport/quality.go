package transport

import (
	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/endpoint"
	"mule-cloud/app/production/services"

	"github.com/gin-gonic/gin"
)

// SubmitInspectionHandler 提交质检
func SubmitInspectionHandler(s services.IQualityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.InspectionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ep := endpoint.MakeSubmitInspectionEndpoint(s)
		resp, err := ep(c.Request.Context(), &req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// GetInspectionListHandler 获取质检列表
func GetInspectionListHandler(s services.IQualityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.InspectionListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ep := endpoint.MakeGetInspectionListEndpoint(s)
		resp, err := ep(c.Request.Context(), &req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// GetInspectionHandler 获取质检详情
func GetInspectionHandler(s services.IQualityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ep := endpoint.MakeGetInspectionEndpoint(s)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// DeleteInspectionHandler 删除质检记录
func DeleteInspectionHandler(s services.IQualityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ep := endpoint.MakeDeleteInspectionEndpoint(s)
		_, err := ep(c.Request.Context(), id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "删除成功"})
	}
}

