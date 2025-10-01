/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Menu {
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
      roles?: string[]
      keepAlive?: boolean
      hide?: boolean
      order?: number
      href?: string
      activeMenu?: string
      withoutTab?: boolean
      pinTab?: boolean
      menuType: 'dir' | 'page'
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

