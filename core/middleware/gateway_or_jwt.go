package middleware

import (
	tenantCtx "mule-cloud/core/context"
	"mule-cloud/core/jwt"
	"mule-cloud/core/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// GatewayOrJWTAuth 支持网关转发和直接JWT验证的认证中间件
//
// 适用场景：
// 1. 通过网关访问：使用网关转发的 X-User-ID, X-Username, X-Tenant-ID, X-Roles headers
// 2. 直接访问服务：验证 Authorization header 中的 JWT token
//
// 使用示例：
//
//	protected := r.Group("/api")
//	middleware.ApplyGatewayOrJWTMiddlewares(protected, jwtManager)
func GatewayOrJWTAuth(jwtManager *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID, username, tenantID, tenantCode string
		var roles []string

		// 优先使用网关传递的用户信息headers
		xUserID := c.GetHeader("X-User-ID")
		xUsername := c.GetHeader("X-Username")
		xTenantID := c.GetHeader("X-Tenant-ID")
		xTenantCode := c.GetHeader("X-Tenant-Code") // ✅ 新增：租户代码
		xRoles := c.GetHeader("X-Roles")

		if xUserID != "" || xUsername != "" {
			// 场景1: 使用网关传递的信息（网关已验证过JWT）
			userID = xUserID
			username = xUsername
			tenantID = xTenantID
			tenantCode = xTenantCode // ✅ 新增
			if xRoles != "" {
				roles = strings.Split(xRoles, ",")
			}
		} else {
			// 场景2: 直接访问服务（非网关），需要验证JWT
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				response.ErrorWithCode(c, 401, "未提供认证token")
				c.Abort()
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.ErrorWithCode(c, 401, "token格式错误")
				c.Abort()
				return
			}

			claims, err := jwtManager.ValidateToken(parts[1])
			if err != nil {
				response.ErrorWithCode(c, 401, "token验证失败: "+err.Error())
				c.Abort()
				return
			}

			userID = claims.UserID
			username = claims.Username
			tenantID = claims.TenantID
			tenantCode = claims.TenantCode // ✅ 新增
			roles = claims.Roles
		}

		// 将用户信息存入Gin Context（向下兼容）
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("tenant_id", tenantID)     // 保留 ID（兼容）
		c.Set("tenant_code", tenantCode) // ✅ 新增：租户代码
		c.Set("roles", roles)

		// ✅ 将租户信息存入标准Context（使用 TenantCode 进行数据库连接）
		ctx := c.Request.Context()
		ctx = tenantCtx.WithTenantCode(ctx, tenantCode)
		ctx = tenantCtx.WithUserID(ctx, userID)
		ctx = tenantCtx.WithUsername(ctx, username)
		ctx = tenantCtx.WithRoles(ctx, roles)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// ApplyGatewayOrJWTMiddlewares 应用支持网关和直接访问的标准中间件
// 适用于通过网关或直接访问的业务服务
//
// 使用示例：
//
//	protected := r.Group("/api")
//	middleware.ApplyGatewayOrJWTMiddlewares(protected, jwtManager)
func ApplyGatewayOrJWTMiddlewares(group *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	group.Use(GatewayOrJWTAuth(jwtManager)) // 网关或JWT认证
	group.Use(TenantContextMiddleware())    // 租户上下文切换
}
