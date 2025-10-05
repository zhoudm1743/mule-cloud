package middleware

import (
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/jwt"
	"mule-cloud/core/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorWithCode(c, 401, "未提供认证token")
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.ErrorWithCode(c, 401, "token格式错误，应为: Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 验证Token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			response.ErrorWithCode(c, 401, "token验证失败: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入Gin Context（保持向下兼容）
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("tenant_id", claims.TenantID)     // 保留 ID（兼容）
		c.Set("tenant_code", claims.TenantCode) // ✅ 新增：租户代码
		c.Set("roles", claims.Roles)
		c.Set("claims", claims)

		// ✅ 将租户信息存入标准Context（使用 TenantCode）
		ctx := c.Request.Context()
		ctx = tenantCtx.WithTenantCode(ctx, claims.TenantCode)
		ctx = tenantCtx.WithUserID(ctx, claims.UserID)
		ctx = tenantCtx.WithUsername(ctx, claims.Username)
		ctx = tenantCtx.WithRoles(ctx, claims.Roles)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsValue, exists := c.Get("claims")
		if !exists {
			response.ErrorWithCode(c, 401, "未认证")
			c.Abort()
			return
		}

		claims := claimsValue.(*jwt.Claims)
		if !claims.HasAnyRole(roles...) {
			response.ErrorWithCode(c, 403, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth 可选认证中间件（有token则验证，没有则跳过）
func OptionalAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := jwtManager.ValidateToken(parts[1])
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("tenant_id", claims.TenantID)     // 保留 ID（兼容）
				c.Set("tenant_code", claims.TenantCode) // ✅ 新增：租户代码
				c.Set("roles", claims.Roles)
				c.Set("claims", claims)

				// ✅ 将租户信息存入标准Context（使用 TenantCode）
				ctx := c.Request.Context()
				ctx = tenantCtx.WithTenantCode(ctx, claims.TenantCode)
				ctx = tenantCtx.WithUserID(ctx, claims.UserID)
				ctx = tenantCtx.WithUsername(ctx, claims.Username)
				ctx = tenantCtx.WithRoles(ctx, claims.Roles)
				c.Request = c.Request.WithContext(ctx)
			}
		}

		c.Next()
	}
}
