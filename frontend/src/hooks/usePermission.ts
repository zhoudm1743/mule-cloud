import { useAuthStore } from '@/store'
import { useRouteStore } from '@/store'
import { isArray, isString } from 'radash'

/** 权限判断 */
export function usePermission() {
  const authStore = useAuthStore()
  const routeStore = useRouteStore()

  /**
   * 检查角色权限（基于角色标识）
   * @param permission 角色标识，如 'super', 'admin'
   */
  function hasPermission(
    permission?: Entity.RoleType | Entity.RoleType[],
  ) {
    if (!permission)
      return true

    if (!authStore.userInfo)
      return false
    
    const { role } = authStore.userInfo
    
    // 防御性检查：role 可能为 null/undefined
    if (!role || !Array.isArray(role)) {
      return false
    }

    // 角色为super可直接通过
    let has = role.includes('super')
    if (!has) {
      if (isArray(permission))
        // 角色为数组, 判断是否有交集
        has = permission.some(i => role.includes(i))

      if (isString(permission))
        // 角色为字符串, 判断是否包含
        has = role.includes(permission)
    }
    return has
  }

  /**
   * 检查操作权限（细粒度权限）
   * @param menuName 菜单名称，如 'admin', 'role'
   * @param action 操作权限，如 'read', 'create', 'update', 'delete', 'pending', 'verify'
   * @returns 是否有权限
   * 
   * 使用示例：
   * hasAction('admin', 'create')  // 检查是否有创建管理员的权限
   * hasAction('finance', 'pending')  // 检查是否有挂账权限
   */
  function hasAction(menuName: string, action: string): boolean {
    // super（系统超管）和 tenant_admin（租户超管）角色拥有所有权限
    if (hasPermission('super') || hasPermission('tenant_admin')) {
      return true
    }

    if (!authStore.userInfo) {
      return false
    }

    const { menu_permissions } = authStore.userInfo
    
    // 如果没有 menu_permissions 数据，检查是否有菜单访问权限
    if (!menu_permissions) {
      // 降级方案：检查用户是否有该菜单的访问权限
      const menus = routeStore.menus || []
      const hasMenu = menus.some(m => m.name === menuName)
      
      if (!hasMenu) {
        return false
      }
      
      // 有菜单访问权限，假设有所有操作权限
      console.warn(`[usePermission] 缺少 menu_permissions 数据，假设有菜单就有权限: ${menuName}.${action}`)
      return true
    }

    // 检查是否有该菜单的权限
    const menuActions = menu_permissions[menuName]
    if (!menuActions) {
      return false
    }

    // 检查是否有指定的操作权限
    return menuActions.includes(action)
  }

  /**
   * 检查资源权限（组合格式）
   * @param resource 资源权限标识，格式：'menuName:action'，如 'admin:create', 'finance:pending'
   * @returns 是否有权限
   * 
   * 使用示例：
   * hasResource('admin:create')  // 检查是否有创建管理员的权限
   * hasResource('finance:pending')  // 检查是否有挂账权限
   */
  function hasResource(resource: string): boolean {
    const [menuName, action] = resource.split(':')
    if (!menuName || !action) {
      console.error(`[usePermission] 资源格式错误: ${resource}，正确格式: 'menuName:action'`)
      return false
    }
    return hasAction(menuName, action)
  }

  return {
    hasPermission,   // 角色权限检查
    hasAction,       // 操作权限检查（细粒度）
    hasResource,     // 资源权限检查（组合格式）
  }
}
