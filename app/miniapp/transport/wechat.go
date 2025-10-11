package transport

import (
	"mule-cloud/app/miniapp/dto"
	"mule-cloud/app/miniapp/endpoint"
	"mule-cloud/app/miniapp/services"
	"mule-cloud/core/binding"
	"mule-cloud/core/response"

	"github.com/gin-gonic/gin"
)

// WechatLoginHandler 微信登录处理器
func WechatLoginHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.WechatLoginRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeWechatLoginEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// BindTenantHandler 绑定租户处理器
func BindTenantHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.BindTenantRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeBindTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// SelectTenantHandler 选择租户处理器
func SelectTenantHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.SelectTenantRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeSelectTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), req)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// SwitchTenantHandler 切换租户处理器（需要JWT认证）
func SwitchTenantHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.SwitchTenantRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeSwitchTenantEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.SwitchTenantRequest{
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

// GetUserInfoHandler 获取用户信息处理器（需要JWT认证）
func GetUserInfoHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		ep := endpoint.MakeGetUserInfoEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.GetUserInfoRequest{
			UserID: userID.(string),
		})
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UpdateUserInfoHandler 更新用户信息处理器（需要JWT认证）
func UpdateUserInfoHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.UpdateUserInfoRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		ep := endpoint.MakeUpdateUserInfoEndpoint(svc)
		resp, err := ep(c.Request.Context(), endpoint.UpdateUserInfoRequest{
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

// GetPhoneNumberHandler 获取微信手机号处理器（需要JWT认证）
func GetPhoneNumberHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		var req dto.GetPhoneNumberRequest
		if err := binding.BindAll(c, &req); err != nil {
			response.Error(c, "参数错误: "+err.Error())
			return
		}

		resp, err := svc.GetPhoneNumber(userID.(string), req.Code)
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, resp)
	}
}

// UnbindPhoneHandler 解绑手机号处理器（需要JWT认证）
func UnbindPhoneHandler(svc services.IWechatService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中间件获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			return
		}

		err := svc.UnbindPhone(userID.(string))
		if err != nil {
			response.Error(c, err.Error())
			return
		}

		response.Success(c, gin.H{
			"success": true,
			"message": "解绑成功",
		})
	}
}
