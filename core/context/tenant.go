package context

import (
	"context"
)

// Context Key 定义
type contextKey string

const (
	// TenantIDKey 实际上存储的是租户代码（code）而不是ID
	// 为了向后兼容保留这个名字，但实际存储的是 tenant_code
	TenantIDKey contextKey = "tenant_id"
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	RolesKey    contextKey = "roles"
)

// WithTenantID 设置租户ID到Context
// 注意：虽然名字是 ID，但实际应该传入租户 code
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// WithTenantCode 设置租户代码到Context（别名，语义更清晰）
// 推荐使用这个函数而不是 WithTenantID
func WithTenantCode(ctx context.Context, tenantCode string) context.Context {
	return WithTenantID(ctx, tenantCode)
}

// GetTenantID 从Context获取租户ID
// 注意：实际返回的是租户 code 而不是 ID
func GetTenantID(ctx context.Context) string {
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

// GetTenantCode 从Context获取租户代码（别名，语义更清晰）
// 推荐使用这个函数而不是 GetTenantID
func GetTenantCode(ctx context.Context) string {
	return GetTenantID(ctx)
}

// WithUserID 设置用户ID到Context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从Context获取用户ID
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// WithUsername 设置用户名到Context
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, UsernameKey, username)
}

// GetUsername 从Context获取用户名
func GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(UsernameKey).(string); ok {
		return username
	}
	return ""
}

// WithRoles 设置角色到Context
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, RolesKey, roles)
}

// GetRoles 从Context获取角色
func GetRoles(ctx context.Context) []string {
	if roles, ok := ctx.Value(RolesKey).([]string); ok {
		return roles
	}
	return []string{}
}

// WithUserInfo 一次性设置所有用户信息到Context
func WithUserInfo(ctx context.Context, tenantID, userID, username string, roles []string) context.Context {
	ctx = WithTenantID(ctx, tenantID)
	ctx = WithUserID(ctx, userID)
	ctx = WithUsername(ctx, username)
	ctx = WithRoles(ctx, roles)
	return ctx
}

// GetUserInfo 一次性获取所有用户信息
func GetUserInfo(ctx context.Context) (tenantID, userID, username string, roles []string) {
	tenantID = GetTenantID(ctx)
	userID = GetUserID(ctx)
	username = GetUsername(ctx)
	roles = GetRoles(ctx)
	return
}
