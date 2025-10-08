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
		currentTenantCode := tenantCtx.GetTenantCode(ctx)

		log.Printf("[租户上下文] 当前 tenantCode: '%s'", currentTenantCode)

		// ✅ 只有系统管理员（tenantCode为空或"system"）可以切换租户上下文
		if currentTenantCode == "" || currentTenantCode == "system" {
			// ✅ 检查是否请求切换到特定租户（使用 tenant_code）
			contextTenantCode := c.GetHeader("X-Tenant-Context")
			log.Printf("[租户上下文] X-Tenant-Context header: '%s'", contextTenantCode)

			if contextTenantCode != "" {
				// 验证是否为系统管理员
				rolesValue, exists := c.Get("roles")
				log.Printf("[租户上下文] roles exist: %v, roles: %v", exists, rolesValue)

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

						log.Printf("[租户上下文] isSuperAdmin: %v", isSuperAdmin)

						if isSuperAdmin {
							// ✅ 允许切换到指定租户上下文（使用 tenant_code）
							ctx = tenantCtx.WithTenantCode(ctx, contextTenantCode)
							c.Request = c.Request.WithContext(ctx)

							log.Printf("[租户上下文切换] ✅ 系统管理员切换到租户: %s", contextTenantCode)

							// ✅ 验证切换后的 tenantCode
							newTenantCode := tenantCtx.GetTenantCode(c.Request.Context())
							log.Printf("[租户上下文切换] ✅ 切换后验证 tenantCode: '%s'", newTenantCode)
						} else {
							log.Printf("[租户上下文切换] ❌ 非系统管理员尝试切换租户，已拒绝")
						}
					}
				}
			} else {
				log.Printf("[租户上下文] 没有 X-Tenant-Context header，不切换")
			}
		} else {
			log.Printf("[租户上下文] 当前不是系统管理员（tenantCode=%s），不允许切换", currentTenantCode)
		}

		c.Next()
	}
}
