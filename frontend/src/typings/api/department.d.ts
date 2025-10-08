/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Department {
    /* 部门信息 */
    interface DepartmentInfo {
      id: string
      name: string
      code: string
      parent_id: string
      status: number
      created_at: number
      updated_at: number
      created_by: string
      updated_by: string
    }

    /* 创建部门请求 */
    interface CreateRequest {
      name: string
      code: string
      parent_id?: string
      status?: number
    }

    /* 更新部门请求 */
    interface UpdateRequest {
      name?: string
      code?: string
      parent_id?: string
      status?: number
    }

    /* 查询部门请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      name?: string
      code?: string
      parent_id?: string
      status?: number
    }

    /* 分页响应 */
    interface ListResponse {
      departments: DepartmentInfo[]
      total: number
      page: number
      size: number
    }

    /* 批量删除请求 */
    interface BatchDeleteRequest {
      ids: string[]
    }
  }
}

