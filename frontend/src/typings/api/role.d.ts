/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Role {
    /* 角色信息 */
    interface RoleInfo {
      id: string
      tenant_id: string
      code: string
      name: string
      description: string
      menus: string[] // 菜单权限列表（存储菜单的 name 字段，如 'dashboard', 'admin'）
      menu_permissions?: Record<string, string[]> // 菜单权限映射：{"admin": ["read", "create", "update"], "role": ["read"]}
      status: number
      created_at: number
      updated_at: number
    }

    /* 创建角色请求 */
    interface CreateRequest {
      tenant_id?: string
      code: string
      name: string
      description?: string
      status?: number
    }

    /* 更新角色请求 */
    interface UpdateRequest {
      code?: string
      name?: string
      description?: string
      status?: number
    }

    /* 查询角色请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      tenant_id?: string
      code?: string
      name?: string
      status?: number
    }

    /* 分页响应 */
    interface ListResponse {
      roles: RoleInfo[]
      total: number
      page: number
      size: number
    }

    /* 分配菜单权限请求 */
    interface AssignMenusRequest {
      menus: string[] // 菜单名称数组（menu.name），如 ['dashboard', 'admin', 'role']
      menu_permissions?: Record<string, string[]> // 菜单权限映射（可选）：{"admin": ["read", "create"], "role": ["read"]}
    }

    /* 角色菜单权限响应 */
    interface MenusResponse {
      menus: string[]
    }
  }
}

