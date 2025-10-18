package transport

import (
	"mule-cloud/app/production/dto"
	"mule-cloud/app/production/endpoint"
	"mule-cloud/app/production/services"

	"github.com/gin-gonic/gin"
)

// CreateReworkHandler 创建返工单
func CreateReworkHandler(s services.IReworkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ReworkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ep := endpoint.MakeCreateReworkEndpoint(s)
		resp, err := ep(c.Request.Context(), &req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// GetReworkListHandler 获取返工列表
func GetReworkListHandler(s services.IReworkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ReworkListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ep := endpoint.MakeGetReworkListEndpoint(s)
		resp, err := ep(c.Request.Context(), &req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// GetReworkHandler 获取返工详情
func GetReworkHandler(s services.IReworkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ep := endpoint.MakeGetReworkEndpoint(s)
		resp, err := ep(c.Request.Context(), id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, resp)
	}
}

// CompleteReworkHandler 完成返工
func CompleteReworkHandler(s services.IReworkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req dto.CompleteReworkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		ep := endpoint.MakeCompleteReworkEndpoint(s)
		_, err := ep(c.Request.Context(), map[string]interface{}{
			"id":  id,
			"req": &req,
		})
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "返工完成"})
	}
}

// DeleteReworkHandler 删除返工记录
func DeleteReworkHandler(s services.IReworkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ep := endpoint.MakeDeleteReworkEndpoint(s)
		_, err := ep(c.Request.Context(), id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "删除成功"})
	}
}

