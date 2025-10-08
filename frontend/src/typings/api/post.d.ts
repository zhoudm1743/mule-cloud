/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Post {
    /* 岗位信息 */
    interface PostInfo {
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

    /* 创建岗位请求 */
    interface CreateRequest {
      name: string
      code: string
      parent_id?: string
      status?: number
    }

    /* 更新岗位请求 */
    interface UpdateRequest {
      name?: string
      code?: string
      parent_id?: string
      status?: number
    }

    /* 查询岗位请求 */
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
      posts: PostInfo[]
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

