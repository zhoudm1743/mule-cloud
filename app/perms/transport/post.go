package transport

import (
	"mule-cloud/app/perms/dto"
	"mule-cloud/app/perms/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建岗位
func CreatePostHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取创建人（待集成JWT后）
		perm := NewPermissionChecker(c)
		createdBy := perm.UserID
		if createdBy == "" {
			createdBy = "system"
		}

		post, err := postSvc.Create(c.Request.Context(), &req, createdBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "创建成功", post)
	}
}

// GetPostHandler 获取岗位详情
func GetPostHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		post, err := postSvc.GetByID(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, post)
	}
}

// ListPostsHandler 查询岗位列表
func ListPostsHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.ListPostRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		posts, total, err := postSvc.List(c.Request.Context(), &req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, map[string]interface{}{
			"posts": posts,
			"total": total,
			"page":  req.Page,
			"size":  req.PageSize,
		})
	}
}

// GetAllPostsHandler 获取所有岗位（不分页）
func GetAllPostsHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := postSvc.GetAll(c.Request.Context())
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, map[string]interface{}{
			"posts": posts,
			"total": len(posts),
		})
	}
}

// UpdatePostHandler 更新岗位
func UpdatePostHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req dto.UpdatePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		// TODO: 从上下文获取更新人（待集成JWT后）
		updatedBy := "system"

		err := postSvc.Update(c.Request.Context(), id, &req, updatedBy)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "更新成功", nil)
	}
}

// DeletePostHandler 删除岗位
func DeletePostHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		err := postSvc.Delete(c.Request.Context(), id)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "删除成功", nil)
	}
}

// BatchDeletePostsHandler 批量删除岗位
func BatchDeletePostsHandler(postSvc *services.PostService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BatchDeletePostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		err := postSvc.BatchDelete(c.Request.Context(), req.IDs)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.SuccessWithMsg(c, "批量删除成功", nil)
	}
}

