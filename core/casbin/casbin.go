package casbin

import (
	"fmt"
	"log"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
)

var (
	// Enforcer 全局 Casbin enforcer 实例
	Enforcer *casbin.Enforcer
	once     sync.Once
)

// Config Casbin 配置
type Config struct {
	MongoURI     string // MongoDB 连接URI
	DatabaseName string // 数据库名称
	ModelPath    string // 模型文件路径（可选，默认使用内置模型）
}

// InitCasbin 初始化 Casbin
func InitCasbin(cfg *Config) (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		// 创建 MongoDB 适配器
		adapter, adapterErr := mongodbadapter.NewAdapter(cfg.MongoURI)
		if adapterErr != nil {
			err = fmt.Errorf("创建Casbin MongoDB适配器失败: %w", adapterErr)
			return
		}

		// 加载模型
		var m model.Model
		if cfg.ModelPath != "" {
			// 从文件加载模型
			m, err = model.NewModelFromFile(cfg.ModelPath)
		} else {
			// 使用内置模型
			m, err = model.NewModelFromString(getDefaultModel())
		}

		if err != nil {
			err = fmt.Errorf("加载Casbin模型失败: %w", err)
			return
		}

		// 创建 Enforcer
		Enforcer, err = casbin.NewEnforcer(m, adapter)
		if err != nil {
			err = fmt.Errorf("创建Casbin Enforcer失败: %w", err)
			return
		}

		// 加载策略
		if loadErr := Enforcer.LoadPolicy(); loadErr != nil {
			log.Printf("⚠️  加载Casbin策略失败: %v", loadErr)
		}

		log.Println("✅ Casbin 初始化成功")
	})

	return Enforcer, err
}

// getDefaultModel 获取默认 RBAC 模型
func getDefaultModel() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
}

// CheckPermission 检查权限
func CheckPermission(sub, obj, act string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	return Enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加权限策略
func AddPolicy(sub, obj, act string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	added, err := Enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return false, err
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return added, nil
}

// RemovePolicy 移除权限策略
func RemovePolicy(sub, obj, act string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	removed, err := Enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return false, err
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return removed, nil
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(user, role string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	added, err := Enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		return false, err
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return added, nil
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(user, role string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	removed, err := Enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		return false, err
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return removed, nil
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(user string) ([]string, error) {
	if Enforcer == nil {
		return nil, fmt.Errorf("Casbin未初始化")
	}

	return Enforcer.GetRolesForUser(user)
}

// GetUsersForRole 获取拥有某角色的所有用户
func GetUsersForRole(role string) ([]string, error) {
	if Enforcer == nil {
		return nil, fmt.Errorf("Casbin未初始化")
	}

	return Enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser 获取用户的所有权限
func GetPermissionsForUser(user string) ([][]string, error) {
	if Enforcer == nil {
		return nil, fmt.Errorf("Casbin未初始化")
	}

	perms := Enforcer.GetPermissionsForUser(user)
	return perms, nil
}

// DeleteRole 删除角色及其所有相关策略
func DeleteRole(role string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	// 删除角色的所有权限
	removed, err := Enforcer.RemoveFilteredPolicy(0, role)
	if err != nil {
		return false, err
	}

	// 删除角色的所有分组关系
	if _, err := Enforcer.RemoveFilteredGroupingPolicy(1, role); err != nil {
		log.Printf("⚠️  删除角色分组关系失败: %v", err)
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return removed, nil
}

// DeleteUser 删除用户及其所有相关策略
func DeleteUser(user string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	// 删除用户的所有权限
	removed1, err := Enforcer.RemoveFilteredPolicy(0, user)
	if err != nil {
		return false, err
	}

	// 删除用户的所有角色
	removed2, err := Enforcer.RemoveFilteredGroupingPolicy(0, user)
	if err != nil {
		return false, err
	}

	// 持久化
	if saveErr := Enforcer.SavePolicy(); saveErr != nil {
		log.Printf("⚠️  保存策略失败: %v", saveErr)
	}

	return removed1 || removed2, nil
}

// SyncRoleMenus 同步角色的菜单权限（废弃，请使用 SyncRoleMenusWithPermissions）
// 先删除角色的所有权限，再批量添加新权限
func SyncRoleMenus(tenantID, roleID string, menuPaths []string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}

	roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)

	// 删除角色的所有现有权限
	if _, err := Enforcer.RemoveFilteredPolicy(0, roleSub); err != nil {
		return fmt.Errorf("删除角色权限失败: %w", err)
	}

	// 批量添加新权限（默认：read, create, update, delete）
	for _, menuPath := range menuPaths {
		// 添加基础 CRUD 权限
		actions := []string{"read", "create", "update", "delete"}
		for _, action := range actions {
			if _, err := Enforcer.AddPolicy(roleSub, menuPath, action); err != nil {
				log.Printf("⚠️  添加权限失败 %s -> %s: %v", roleSub, menuPath, err)
			}
		}
	}

	// 持久化
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}

// SyncRoleMenusWithPermissions 同步角色的菜单权限（支持细粒度权限）
// menuPermissions: {"admin": ["read", "create", "update"], "role": ["read"]}
// menuPathMap: {"admin": "/system/admin", "role": "/system/role"} - 菜单名到路径的映射
func SyncRoleMenusWithPermissions(tenantID, roleID string, menuPermissions map[string][]string, menuPathMap map[string]string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}

	roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)

	// 删除角色的所有现有权限
	if _, err := Enforcer.RemoveFilteredPolicy(0, roleSub); err != nil {
		return fmt.Errorf("删除角色权限失败: %w", err)
	}

	// 批量添加新权限
	for menuName, actions := range menuPermissions {
		// 从映射中获取路径
		menuPath := menuPathMap[menuName]
		if menuPath == "" {
			// 如果映射中没有，使用降级方案
			menuPath = getMenuPathFromNameFallback(menuName)
			log.Printf("⚠️  菜单 %s 没有提供路径映射，使用降级路径: %s", menuName, menuPath)
		}

		for _, action := range actions {
			if _, err := Enforcer.AddPolicy(roleSub, menuPath, action); err != nil {
				log.Printf("⚠️  添加权限失败 %s -> %s:%s: %v", roleSub, menuPath, action, err)
			}
		}
	}

	// 持久化
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	log.Printf("[Casbin] 角色权限同步成功: %s, %d个菜单", roleSub, len(menuPermissions))
	return nil
}

// getMenuPathFromNameFallback 降级方案：使用静态映射
func getMenuPathFromNameFallback(menuName string) string {
	mapping := map[string]string{
		"dashboard": "/dashboard",
		"perms":     "/perms",
		"admin":     "/perms/admin",
		"role":      "/perms/role",
		"tenant":    "/perms/tenant",
		"menu":      "/perms/menu",
	}

	if path, ok := mapping[menuName]; ok {
		return path
	}

	// 默认返回 /菜单名
	return "/" + menuName
}

// SyncUserRoles 同步用户的角色
// 先删除用户的所有角色，再批量添加新角色
func SyncUserRoles(tenantID, userID string, roleIDs []string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}

	userSub := fmt.Sprintf("tenant:%s:user:%s", tenantID, userID)

	// 删除用户的所有现有角色
	if _, err := Enforcer.RemoveFilteredGroupingPolicy(0, userSub); err != nil {
		return fmt.Errorf("删除用户角色失败: %w", err)
	}

	// 批量添加新角色
	for _, roleID := range roleIDs {
		roleSub := fmt.Sprintf("tenant:%s:role:%s", tenantID, roleID)
		if _, err := Enforcer.AddGroupingPolicy(userSub, roleSub); err != nil {
			log.Printf("⚠️  添加用户角色失败 %s -> %s: %v", userSub, roleSub, err)
		}
	}

	// 持久化
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}

// CheckUserPermission 检查用户是否有权限访问某个资源
func CheckUserPermission(tenantID, userID, resource, action string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	userSub := fmt.Sprintf("tenant:%s:user:%s", tenantID, userID)
	return Enforcer.Enforce(userSub, resource, action)
}

// CheckSuperAdmin 检查是否是超级管理员
func CheckSuperAdmin(userID string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	superSub := fmt.Sprintf("super:user:%s", userID)
	// 超级管理员拥有所有权限
	return Enforcer.Enforce(superSub, "*", "*")
}

// AddSuperAdmin 添加超级管理员
func AddSuperAdmin(userID string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}

	superSub := fmt.Sprintf("super:user:%s", userID)

	// 添加到超级管理员组
	if _, err := Enforcer.AddGroupingPolicy(superSub, "super:admin"); err != nil {
		return fmt.Errorf("添加超级管理员失败: %w", err)
	}

	// 超级管理员拥有所有权限
	if _, err := Enforcer.AddPolicy("super:admin", "*", "*"); err != nil {
		return fmt.Errorf("添加超级管理员权限失败: %w", err)
	}

	// 持久化
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	return nil
}

// SyncTenantMenus 同步租户的菜单权限（超管分配菜单给租户）
// 这定义了租户的权限边界，租户内的角色不能超出这个范围
func SyncTenantMenus(tenantID string, menuPaths []string) error {
	if Enforcer == nil {
		return fmt.Errorf("Casbin未初始化")
	}

	tenantSub := fmt.Sprintf("tenant:%s", tenantID)

	// 删除租户的所有现有权限
	if _, err := Enforcer.RemoveFilteredPolicy(0, tenantSub); err != nil {
		return fmt.Errorf("删除租户权限失败: %w", err)
	}

	// 批量添加新权限
	for _, menuPath := range menuPaths {
		// 添加读权限
		if _, err := Enforcer.AddPolicy(tenantSub, menuPath, "read"); err != nil {
			log.Printf("⚠️  添加权限失败 %s -> %s: %v", tenantSub, menuPath, err)
		}
		// 添加写权限
		if _, err := Enforcer.AddPolicy(tenantSub, menuPath, "write"); err != nil {
			log.Printf("⚠️  添加权限失败 %s -> %s: %v", tenantSub, menuPath, err)
		}
	}

	// 持久化
	if err := Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	log.Printf("[Casbin] 租户权限同步成功: %s, %d个菜单", tenantID, len(menuPaths))
	return nil
}

// CheckTenantPermission 检查租户是否有权限访问某个资源
func CheckTenantPermission(tenantID, resource, action string) (bool, error) {
	if Enforcer == nil {
		return false, fmt.Errorf("Casbin未初始化")
	}

	tenantSub := fmt.Sprintf("tenant:%s", tenantID)
	return Enforcer.Enforce(tenantSub, resource, action)
}
