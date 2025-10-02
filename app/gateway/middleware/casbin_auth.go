package middleware

import (
	"fmt"
	"log"
	"mule-cloud/core/casbin"
	"mule-cloud/core/response"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// CasbinAuthMiddleware Casbin 鉴权中间件
func CasbinAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 从上下文获取用户信息（假设已经通过 JWT 中间件解析）
		// TODO: 集成 JWT 后从 token 中获取
		userID, exists := c.Get("user_id")
		if !exists {
			response.Error(c, "未授权：缺少用户信息")
			c.Abort()
			return
		}

		tenantID, _ := c.Get("tenant_id") // 租户ID
		roles, _ := c.Get("roles")        // 用户角色列表
		tenantIDStr, _ := tenantID.(string)

		// 检查角色权限
		if roleList, ok := roles.([]interface{}); ok {
			for _, role := range roleList {
				if roleStr, ok := role.(string); ok {
					// 系统级超管（tenant_id 为空）- 拥有所有权限
					if roleStr == "super" && tenantIDStr == "" {
						log.Printf("[Casbin] 系统超管访问: user=%v, %s %s", userID, method, path)
						c.Next()
						return
					}

					// 租户级超管（tenant_admin）- 拥有本租户所有权限
					if roleStr == "tenant_admin" && tenantIDStr != "" {
						log.Printf("[Casbin] 租户超管访问: tenant=%v, user=%v, %s %s", tenantIDStr, userID, method, path)
						c.Next()
						return
					}
				}
			}
		}

		// 智能解析资源和动作
		resource, action := parseResourceAndAction(path, method)

		// 检查权限（普通用户通过 Casbin）
		var allowed bool
		var err error

		if tenantIDStr != "" {
			// 租户用户
			userSub := fmt.Sprintf("tenant:%s:user:%s", tenantIDStr, userID)
			allowed, err = casbin.CheckPermission(userSub, resource, action)
			log.Printf("[Casbin] 租户用户权限检查: sub=%s, resource=%s, action=%s, allowed=%v", userSub, resource, action, allowed)
		} else {
			// 系统用户（非超管的系统用户）
			userSub := fmt.Sprintf("user:%s", userID)
			allowed, err = casbin.CheckPermission(userSub, resource, action)
			log.Printf("[Casbin] 系统用户权限检查: sub=%s, resource=%s, action=%s, allowed=%v", userSub, resource, action, allowed)
		}

		if err != nil {
			log.Printf("[Casbin] 权限检查失败: %v", err)
			response.Error(c, "权限检查失败")
			c.Abort()
			return
		}

		if !allowed {
			log.Printf("[Casbin] 权限拒绝: user=%v, resource=%s, action=%s", userID, resource, action)
			response.Error(c, "权限不足")
			c.Abort()
			return
		}

		log.Printf("[Casbin] 权限通过: user=%v, resource=%s, action=%s", userID, resource, action)
		c.Next()
	}
}

// parseResourceAndAction 智能解析资源路径和权限动作
//
// 规则：
// 1. 如果路径最后一段是业务动作词（非ID），则：
//   - resource = 去掉最后一段和所有ID参数，只保留资源名称
//   - action = 最后一段
//     例如：POST /finance/pending → resource="/finance", action="pending"
//     例如：POST /finance/123/pending → resource="/finance", action="pending"
//     例如：POST /finance/123/opt1/pending → resource="/finance", action="pending"
//
// 2. 否则使用 RESTful 风格：
//   - resource = 完整路径（保留ID用于精确匹配）
//   - action = HTTP方法映射（GET→read, POST→create...）
//     例如：POST /finance → resource="/finance", action="create"
//     例如：PUT /finance/123 → resource="/finance/123", action="update"
func parseResourceAndAction(path string, method string) (resource string, action string) {
	// 去掉前缀 /admin
	path = strings.TrimPrefix(path, "/admin")

	// 分割路径
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return path, getActionFromMethod(method)
	}

	// 获取最后一段
	lastSegment := parts[len(parts)-1]

	// 判断最后一段是否是业务动作（非ID形式）
	if isBusinessAction(lastSegment) {
		// 业务动作路径
		action = lastSegment

		// 去掉最后一段（action），保留中间部分
		resourceParts := parts[:len(parts)-1]

		// 清理掉所有看起来像 ID 的段，只保留资源名称
		cleanParts := []string{}
		for _, part := range resourceParts {
			if !isID(part) {
				cleanParts = append(cleanParts, part)
			}
		}

		// 构建资源路径
		if len(cleanParts) == 0 {
			resource = "/"
		} else {
			resource = "/" + strings.Join(cleanParts, "/")
		}

		log.Printf("[Casbin] 检测到业务动作: %s %s → resource=%s, action=%s", method, path, resource, action)
		return resource, action
	}

	// RESTful 风格路径
	resource = path
	action = getActionFromMethod(method)
	return resource, action
}

// isBusinessAction 判断是否是业务动作
// 规则：在预定义的业务动作列表中
func isBusinessAction(segment string) bool {
	// 常见业务动作列表（可扩展）
	businessActions := map[string]bool{
		// 财务相关
		"pending":   true, // 挂账
		"verify":    true, // 核销
		"settle":    true, // 结算
		"reconcile": true, // 对账
		"refund":    true, // 退款

		// 审批相关
		"approve":  true, // 批准
		"reject":   true, // 拒绝
		"audit":    true, // 审核
		"submit":   true, // 提交
		"withdraw": true, // 撤回

		// 数据操作
		"export":    true, // 导出
		"import":    true, // 导入
		"sync":      true, // 同步
		"refresh":   true, // 刷新
		"calculate": true, // 计算
		"generate":  true, // 生成

		// 状态变更
		"publish": true, // 发布
		"cancel":  true, // 取消
		"close":   true, // 关闭
		"reopen":  true, // 重开
		"archive": true, // 归档
		"restore": true, // 恢复

		// 权限管理
		"assign":   true, // 分配
		"transfer": true, // 转移
		"lock":     true, // 锁定
		"unlock":   true, // 解锁
		"enable":   true, // 启用
		"disable":  true, // 禁用

		// 其他
		"copy":      true, // 复制
		"move":      true, // 移动
		"merge":     true, // 合并
		"split":     true, // 拆分
		"convert":   true, // 转换
		"validate":  true, // 验证
		"notify":    true, // 通知
		"remind":    true, // 提醒
		"share":     true, // 分享
		"favorite":  true, // 收藏
		"star":      true, // 标星
		"pin":       true, // 置顶
		"unpin":     true, // 取消置顶
		"reset":     true, // 重置
		"retry":     true, // 重试
		"rollback":  true, // 回滚
		"upgrade":   true, // 升级
		"downgrade": true, // 降级
	}

	return businessActions[segment]
}

// isID 判断是否是 ID 参数
func isID(segment string) bool {
	// MongoDB ObjectID (24位十六进制)
	if matched, _ := regexp.MatchString(`^[\da-fA-F]{24}$`, segment); matched {
		return true
	}
	// 纯数字
	if matched, _ := regexp.MatchString(`^[0-9]+$`, segment); matched {
		return true
	}
	// UUID
	if matched, _ := regexp.MatchString(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, segment); matched {
		return true
	}

	return false
}

// getActionFromMethod 根据 HTTP 方法映射到权限动作（细粒度CRUD）
func getActionFromMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET", "HEAD":
		return "read" // 查询
	case "POST":
		return "create" // 创建
	case "PUT", "PATCH":
		return "update" // 修改
	case "DELETE":
		return "delete" // 删除
	case "OPTIONS":
		return "*" // 预检请求放行
	default:
		return "read"
	}
}

// CasbinAuthConfig 可选配置的鉴权中间件
type CasbinAuthConfig struct {
	SkipPaths []string // 跳过鉴权的路径
}

// CasbinAuthMiddlewareWithConfig 带配置的 Casbin 鉴权中间件
func CasbinAuthMiddlewareWithConfig(config CasbinAuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否在跳过列表中
		for _, skipPath := range config.SkipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// 执行正常鉴权流程
		CasbinAuthMiddleware()(c)
	}
}
