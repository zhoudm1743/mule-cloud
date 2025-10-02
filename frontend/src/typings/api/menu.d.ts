/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Menu {
    /* 权限定义 */
    interface Permission {
      action: string // 权限动作：read, create, update, delete, pending, verify...
      label: string // 显示名称：查看、创建、修改、删除、挂账、核销...
      description?: string // 描述
      is_basic: boolean // 是否基础权限（CRUD）
    }

    /* 菜单项（扁平结构 - 适配Nova-admin） */
    interface MenuItem {
      id: string
      pid: string | null
      name: string
      path: string
      title: string
      componentPath: string | null
      redirect?: string
      icon?: string
      requiresAuth: boolean
      roles?: string[] // 角色权限（系统级硬限制，一般不使用）
      keepAlive?: boolean
      hide?: boolean
      order?: number
      href?: string
      activeMenu?: string
      withoutTab?: boolean
      pinTab?: boolean
      menuType: 'dir' | 'page'
      available_permissions?: Permission[] // 该菜单支持的权限列表
      status?: number
      created_at?: number
      updated_at?: number
    }

    /* 创建菜单请求 */
    interface CreateRequest {
      pid?: string | null
      name: string
      path: string
      title: string
      componentPath?: string | null
      redirect?: string
      icon?: string
      requiresAuth?: boolean
      roles?: string[]
      keepAlive?: boolean
      hide?: boolean
      order?: number
      href?: string
      activeMenu?: string
      withoutTab?: boolean
      pinTab?: boolean
      menuType: 'dir' | 'page'
      available_permissions?: Permission[]
    }

    /* 更新菜单请求 */
    interface UpdateRequest {
      pid?: string | null
      name?: string
      path?: string
      title?: string
      componentPath?: string | null
      redirect?: string
      icon?: string
      requiresAuth?: boolean
      roles?: string[]
      keepAlive?: boolean
      hide?: boolean
      order?: number
      href?: string
      activeMenu?: string
      withoutTab?: boolean
      pinTab?: boolean
      menuType?: 'dir' | 'page'
      available_permissions?: Permission[]
    }

    /* 查询菜单请求 */
    interface ListRequest {
      page?: number
      pageSize?: number
      name?: string
      title?: string
      menuType?: string
      status?: number
    }

    /* 分页响应 */
    interface ListResponse {
      list: MenuItem[]
      total: number
      page: number
      pageSize: number
    }
  }
}

