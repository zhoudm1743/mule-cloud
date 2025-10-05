package middleware

import (
	"log"
	tenantCtx "mule-cloud/core/context"

	"github.com/gin-gonic/gin"
)

// TenantContextMiddleware 租户上下文切换中间件
// 允许系统管理员通过 X-Tenant-Context header 切换到指定租户的上下文
func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		currentTenantID := tenantCtx.GetTenantID(ctx)

		// 只有系统管理员（tenantID为空）可以切换租户上下文
		if currentTenantID == "" {
			// 检查是否请求切换到特定租户
			contextTenantID := c.GetHeader("X-Tenant-Context")

			if contextTenantID != "" {
				// 验证是否为系统管理员
				rolesValue, exists := c.Get("roles")
				if exists {
					if roles, ok := rolesValue.([]string); ok {
						// 检查是否有 super 角色
						isSuperAdmin := false
						for _, role := range roles {
							if role == "super" {
								isSuperAdmin = true
								break
							}
						}

						if isSuperAdmin {
							// 允许切换到指定租户上下文
							ctx = tenantCtx.WithTenantID(ctx, contextTenantID)
							c.Request = c.Request.WithContext(ctx)

							log.Printf("[租户上下文切换] 系统管理员切换到租户: %s", contextTenantID)
						} else {
							log.Printf("[租户上下文切换] 非系统管理员尝试切换租户，已拒绝")
						}
					}
				}
			}
		}

		c.Next()
	}
}
