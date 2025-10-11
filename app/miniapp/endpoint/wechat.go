package endpoint

import (
	"context"
	"mule-cloud/app/miniapp/dto"
	"mule-cloud/app/miniapp/services"

	"github.com/go-kit/kit/endpoint"
)

// MakeWechatLoginEndpoint 创建微信登录端点
func MakeWechatLoginEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.WechatLoginRequest)
		return svc.WechatLogin(req)
	}
}

// MakeBindTenantEndpoint 创建绑定租户端点
func MakeBindTenantEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.BindTenantRequest)
		return svc.BindTenant(req)
	}
}

// MakeSelectTenantEndpoint 创建选择租户端点
func MakeSelectTenantEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.SelectTenantRequest)
		return svc.SelectTenant(req)
	}
}

// SwitchTenantRequest 切换租户请求（包含userID）
type SwitchTenantRequest struct {
	UserID string
	Data   dto.SwitchTenantRequest
}

// MakeSwitchTenantEndpoint 创建切换租户端点
func MakeSwitchTenantEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SwitchTenantRequest)
		return svc.SwitchTenant(req.UserID, req.Data)
	}
}

// GetUserInfoRequest 获取用户信息请求
type GetUserInfoRequest struct {
	UserID string
}

// MakeGetUserInfoEndpoint 创建获取用户信息端点
func MakeGetUserInfoEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserInfoRequest)
		return svc.GetUserInfo(req.UserID)
	}
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	UserID string
	Data   dto.UpdateUserInfoRequest
}

// MakeUpdateUserInfoEndpoint 创建更新用户信息端点
func MakeUpdateUserInfoEndpoint(svc services.IWechatService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateUserInfoRequest)
		return svc.UpdateUserInfo(req.UserID, req.Data)
	}
}

