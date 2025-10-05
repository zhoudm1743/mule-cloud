package endpoint

import (
	"context"
	"mule-cloud/app/auth/dto"
	"mule-cloud/app/auth/services"

	"github.com/go-kit/kit/endpoint"
)

// MakeLoginEndpoint 创建登录端点
func MakeLoginEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.LoginRequest)
		return svc.Login(req)
	}
}

// MakeRegisterEndpoint 创建注册端点
func MakeRegisterEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.RegisterRequest)
		return svc.Register(req)
	}
}

// MakeRefreshTokenEndpoint 创建刷新Token端点
func MakeRefreshTokenEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(dto.RefreshTokenRequest)
		return svc.RefreshToken(req)
	}
}

// GetProfileEndpoint 获取个人信息端点
type GetProfileRequest struct {
	UserID string `json:"user_id"`
}

func MakeGetProfileEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetProfileRequest)
		return svc.GetProfile(ctx, req.UserID)
	}
}

// UpdateProfileEndpoint 更新个人信息端点
type UpdateProfileRequest struct {
	UserID string
	Data   dto.UpdateProfileRequest
}

func MakeUpdateProfileEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateProfileRequest)
		return svc.UpdateProfile(ctx, req.UserID, req.Data)
	}
}

// ChangePasswordEndpoint 修改密码端点
type ChangePasswordRequest struct {
	UserID string
	Data   dto.ChangePasswordRequest
}

func MakeChangePasswordEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangePasswordRequest)
		return svc.ChangePassword(ctx, req.UserID, req.Data)
	}
}

// GetUserRoutesEndpoint 获取用户路由端点
type GetUserRoutesRequest struct {
	UserID string `json:"user_id"`
}

func MakeGetUserRoutesEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRoutesRequest)
		return svc.GetUserRoutes(ctx, req.UserID)
	}
}

// MakeGetTenantListEndpoint 创建获取租户列表端点
func MakeGetTenantListEndpoint(svc services.IAuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return svc.GetTenantList()
	}
}
