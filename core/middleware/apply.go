package middleware

import (
	"mule-cloud/core/jwt"

	"github.com/gin-gonic/gin"
)

// MiddlewareConfig 中间件配置（可选）
type MiddlewareConfig struct {
	// SupportGateway 是否支持网关转发（默认: true）
	// true: 支持网关 X-User-ID 等 headers 和直接 JWT 验证
	// false: 只支持直接 JWT 验证
	SupportGateway bool

	// RequireRole 要求的角色（可选）
	// 例如: []string{"super"} 表示只有管理员可以访问
	RequireRole []string

	// Optional 是否可选认证（默认: false）
	// true: 有token则验证，没有则跳过
	// false: 必须提供token
	Optional bool
}

// Apply 应用标准中间件到路由组
//
// 这是唯一需要的中间件应用函数，适用于所有场景：
//
// 场景1: 标准业务服务（最常用）
//
//	middleware.Apply(group, jwtManager)
//
// 场景2: 只支持直接访问（不通过网关）
//
//	middleware.Apply(group, jwtManager, middleware.MiddlewareConfig{
//	    SupportGateway: false,
//	})
//
// 场景3: 要求管理员角色
//
//	middleware.Apply(group, jwtManager, middleware.MiddlewareConfig{
//	    RequireRole: []string{"super"},
//	})
//
// 场景4: 可选认证（公开接口）
//
//	middleware.Apply(group, jwtManager, middleware.MiddlewareConfig{
//	    Optional: true,
//	})
func Apply(group *gin.RouterGroup, jwtManager *jwt.JWTManager, config ...MiddlewareConfig) {
	// 默认配置
	cfg := MiddlewareConfig{
		SupportGateway: true,  // 默认支持网关
		RequireRole:    nil,   // 默认不要求特定角色
		Optional:       false, // 默认必须认证
	}

	// 如果提供了配置，使用提供的配置
	if len(config) > 0 {
		cfg = config[0]
	}

	// 1. 应用认证中间件
	if cfg.Optional {
		// 可选认证
		group.Use(OptionalAuth(jwtManager))
	} else if cfg.SupportGateway {
		// 支持网关或JWT认证（默认）
		group.Use(GatewayOrJWTAuth(jwtManager))
	} else {
		// 只支持JWT认证
		group.Use(JWTAuth(jwtManager))
	}

	// 2. 如果要求特定角色
	if len(cfg.RequireRole) > 0 {
		group.Use(RequireRole(cfg.RequireRole...))
	}

	// 3. 租户上下文切换（所有场景都需要）
	group.Use(TenantContextMiddleware())
}
