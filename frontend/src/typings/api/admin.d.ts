/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Admin {
    /* 管理员信息 */
    interface AdminInfo {
      id: string
      tenant_id: string
      phone: string
      nickname: string
      email: string
      avatar: string
      roles: string[]
      is_super: boolean
      status: number
      created_at: number
      updated_at: number
    }

    /* 创建管理员请求 */
    interface CreateRequest {
      phone: string
      password: string
      nickname?: string
      email?: string
      avatar?: string
      tenant_id?: string  // 租户ID（空表示系统级用户）
      roles?: string[]
      status?: number
    }

    /* 更新管理员请求 */
    interface UpdateRequest {
      phone?: string
      password?: string
      nickname?: string
      email?: string
      avatar?: string
      tenant_id?: string  // 租户ID
      roles?: string[]
      status?: number
    }

    /* 查询管理员请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      phone?: string
      nickname?: string
      tenant_id?: string
      status?: number
    }

    /* 分页响应 */
    interface ListResponse {
      admins: AdminInfo[]
      total: number
      page: number
      size: number
    }

    /* 分配角色请求 */
    interface AssignRolesRequest {
      roles: string[]
    }
  }
}

