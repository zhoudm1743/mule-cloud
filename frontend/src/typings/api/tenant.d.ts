/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Tenant {
    /* 租户信息 */
    interface TenantInfo {
      id: string
      code: string
      name: string
      contact: string
      phone: string
      email: string
      menus: string[] // 菜单权限列表（存储菜单的 name 字段，如 'dashboard', 'admin'）
      status: number
      created_at: number
      updated_at: number
    }

    /* 创建租户请求 */
    interface CreateRequest {
      code: string
      name: string
      contact?: string
      phone?: string
      email?: string
      status?: number
    }

    /* 更新租户请求 */
    interface UpdateRequest {
      code?: string
      name?: string
      contact?: string
      phone?: string
      email?: string
      status?: number
    }

    /* 查询租户请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      code?: string
      name?: string
      status?: number
    }

    /* 分页响应 */
    interface ListResponse {
      list: TenantInfo[]
      total: number
      page: number
      page_size: number
    }

    /* 分配菜单权限请求 */
    interface AssignMenusRequest {
      menus: string[] // 菜单名称数组（menu.name），如 ['dashboard', 'admin', 'role']
    }

    /* 租户菜单权限响应 */
    interface MenusResponse {
      menus: string[]
    }
  }
}

