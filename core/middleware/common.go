package middleware

import (
	"mule-cloud/core/jwt"

	"github.com/gin-gonic/gin"
)

// ApplyCommonMiddlewares 应用标准中间件到路由组
// 适用于需要认证和租户隔离的业务服务
//
// 使用示例：
//
//	protected := r.Group("/api")
//	middleware.ApplyCommonMiddlewares(protected, jwtManager)
func ApplyCommonMiddlewares(group *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	group.Use(JWTAuth(jwtManager))       // JWT 认证
	group.Use(TenantContextMiddleware()) // 租户上下文切换（支持系统管理员切换租户）
}

// ApplyOptionalAuthMiddlewares 应用可选认证中间件
// 适用于公开 + 需要认证的混合路由
//
// 使用示例：
//
//	public := r.Group("/public")
//	middleware.ApplyOptionalAuthMiddlewares(public, jwtManager)
func ApplyOptionalAuthMiddlewares(group *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	group.Use(OptionalAuth(jwtManager))  // 可选认证
	group.Use(TenantContextMiddleware()) // 租户上下文切换
}

// ApplyAdminMiddlewares 应用管理员中间件
// 要求用户必须有 super 角色
//
// 使用示例：
//
//	admin := r.Group("/admin")
//	middleware.ApplyAdminMiddlewares(admin, jwtManager)
func ApplyAdminMiddlewares(group *gin.RouterGroup, jwtManager *jwt.JWTManager) {
	group.Use(JWTAuth(jwtManager))       // JWT 认证
	group.Use(RequireRole("super"))      // 要求 super 角色
	group.Use(TenantContextMiddleware()) // 租户上下文切换
}
