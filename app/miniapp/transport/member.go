package transport

import (
	"mule-cloud/app/miniapp/dto"
	"mule-cloud/app/miniapp/endpoint"
	"mule-cloud/app/miniapp/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// ========== 员工档案相关 Handler ==========

// GetProfileHandler 获取个人档案处理器（需要JWT认证）
func GetProfileHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		ep := endpoint.MakeGetProfileEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.GetProfileRequest{
			UserID: userID.(string),
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateBasicInfoHandler 更新基本信息处理器（需要JWT认证）
func UpdateBasicInfoHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.UpdateBasicInfoRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUpdateBasicInfoEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UpdateBasicInfoRequest{
			UserID: userID.(string),
			Data:   req,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateContactInfoHandler 更新联系信息处理器（需要JWT认证）
func UpdateContactInfoHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.UpdateContactInfoRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUpdateContactInfoEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UpdateContactInfoRequest{
			UserID: userID.(string),
			Data:   req,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UploadPhotoHandler 上传照片处理器（需要JWT认证）
func UploadPhotoHandler(svc services.IMemberService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.UploadPhotoRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUploadPhotoEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UploadPhotoRequest{
			UserID: userID.(string),
			Data:   req,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
