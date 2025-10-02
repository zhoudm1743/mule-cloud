package transport

import (
	"mule-cloud/app/auth/dto"
	"mule-cloud/app/auth/endpoint"
	"mule-cloud/app/auth/services"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// LoginHandler 登录处理器
func LoginHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeLoginEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// RegisterHandler 注册处理器
func RegisterHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeRegisterEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// RefreshTokenHandler 刷新Token处理器
func RefreshTokenHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.RefreshTokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeRefreshTokenEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// GetProfileHandler 获取个人信息处理器
func GetProfileHandler(svc services.IAuthService) gin.HandlerFunc {
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

// UpdateProfileHandler 更新个人信息处理器
func UpdateProfileHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.UpdateProfileRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUpdateProfileEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UpdateProfileRequest{
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

// ChangePasswordHandler 修改密码处理器
func ChangePasswordHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.ChangePasswordRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeChangePasswordEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.ChangePasswordRequest{
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

// GetUserRoutesHandler 获取用户路由处理器
func GetUserRoutesHandler(svc services.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 query 参数或 JWT 中间件获取用户ID
		userID := c.Query("id")
		if userID == "" {
			// 如果没有传 id，从 JWT 获取
			if id, exists := c.Get("user_id"); exists {
				userID = id.(string)
			} else {
				response.ErrorWithCode(c, 401, "未认证")
				return
			}
		}

		ep := endpoint.MakeGetUserRoutesEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.GetUserRoutesRequest{
			UserID: userID,
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}
