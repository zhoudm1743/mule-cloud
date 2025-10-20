/// <reference path="../global.d.ts"/>

namespace Api {
  namespace OperationLog {
    /* 操作日志信息 */
    interface OperationLogInfo {
      id: string
      user_id: string
      username: string
      method: string // HTTP方法
      path: string // 请求路径
      resource: string // 资源名称
      action: string // 操作类型
      request_body: string // 请求体
      response_code: number // 响应状态码
      duration: number // 耗时（毫秒）
      ip: string // 客户端IP
      user_agent: string // 用户代理
      error?: string // 错误信息
      created_at: string // 创建时间
    }

    /* 查询操作日志请求 */
    interface ListRequest {
      page: number
      page_size: number
      user_id?: string
      username?: string
      method?: string
      resource?: string
      action?: string
      response_code?: number
      start_time?: number
      end_time?: number
    }

    /* 分页响应 */
    interface ListResponse {
      list: OperationLogInfo[]
      total: number
      page: number
      page_size: number
    }

    /* 统计请求 */
    interface StatsRequest {
      start_time: number
      end_time: number
      group_by?: string
    }

    /* 用户统计 */
    interface UserStats {
      user_id: string
      username: string
      count: number
    }

    /* 操作统计 */
    interface ActionStats {
      action: string
      count: number
    }

    /* 统计响应 */
    interface StatsResponse {
      total: number
      success_num: number
      fail_num: number
      avg_time: number
      stats: Record<string, any>
      top_users: UserStats[]
      top_actions: ActionStats[]
    }

    /* 详情响应 */
    interface DetailResponse {
      log: OperationLogInfo
    }
  }
}

