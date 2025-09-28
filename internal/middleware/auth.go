package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zhoudm1743/mule-cloud/pkg/auth"
	"github.com/zhoudm1743/mule-cloud/pkg/logger"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(tokenManager *auth.TokenManager, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		token := extractToken(c)
		if token == "" {
			logger.Warn("Missing authorization token", "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing authorization token",
			})
			c.Abort()
			return
		}

		// 验证token
		claims, err := tokenManager.ValidateToken(token)
		if err != nil {
			logger.Warn("Invalid token", "error", err.Error(), "path", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid token",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("permissions", claims.Permissions)
		c.Set("claims", claims)

		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
func OptionalAuthMiddleware(tokenManager *auth.TokenManager, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			claims, err := tokenManager.ValidateToken(token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("roles", claims.Roles)
				c.Set("permissions", claims.Permissions)
				c.Set("claims", claims)
			}
		}
		c.Next()
	}
}

// extractToken 从请求中提取token
func extractToken(c *gin.Context) string {
	// 从Header中获取
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:6]) == "BEARER" {
		return bearerToken[7:]
	}

	// 从Query参数中获取
	token := c.Query("token")
	if token != "" {
		return token
	}

	return ""
}

// GetUserID 获取当前用户ID
func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get("user_id"); exists {
		return userID.(string)
	}
	return ""
}

// GetUsername 获取当前用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		return username.(string)
	}
	return ""
}

// GetUserRoles 获取当前用户角色
func GetUserRoles(c *gin.Context) []string {
	if roles, exists := c.Get("roles"); exists {
		return roles.([]string)
	}
	return []string{}
}

// GetUserPermissions 获取当前用户权限
func GetUserPermissions(c *gin.Context) []string {
	if permissions, exists := c.Get("permissions"); exists {
		return permissions.([]string)
	}
	return []string{}
}

// GetClaims 获取Claims
func GetClaims(c *gin.Context) *auth.CustomClaims {
	if claims, exists := c.Get("claims"); exists {
		return claims.(*auth.CustomClaims)
	}
	return nil
}

// HasRole 检查用户是否拥有指定角色
func HasRole(c *gin.Context, role string) bool {
	roles := GetUserRoles(c)
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasPermission 检查用户是否拥有指定权限
func HasPermission(c *gin.Context, permission string) bool {
	permissions := GetUserPermissions(c)
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// RequireRole 要求特定角色的中间件
func RequireRole(role string, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, role) {
			logger.Warn("Access denied - missing required role",
				"user_id", GetUserID(c),
				"required_role", role,
				"path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Access denied - insufficient role",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequirePermission 要求特定权限的中间件
func RequirePermission(permission string, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasPermission(c, permission) {
			logger.Warn("Access denied - missing required permission",
				"user_id", GetUserID(c),
				"required_permission", permission,
				"path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "Access denied - insufficient permission",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
