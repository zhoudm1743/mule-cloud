/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Basic {
    /* 基础数据信息 */
    interface BasicInfo {
      id: string
      value: string
      remark: string
      created_at: number
      updated_at: number
    }

    /* 创建请求 */
    interface CreateRequest {
      value: string
      remark?: string
    }

    /* 更新请求 */
    interface UpdateRequest {
      value?: string
      remark?: string
    }

    /* 查询请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      value?: string
    }

    /* 分页响应 */
    interface ListResponse<T = BasicInfo> {
      total: number
      page?: number
      page_size?: number
    }

    /* 客户 */
    namespace Customer {
      interface ListResponse extends Basic.ListResponse {
        customers: BasicInfo[]
      }
    }

    /* 颜色 */
    namespace Color {
      interface ListResponse extends Basic.ListResponse {
        colors: BasicInfo[]
      }
    }

    /* 业务员 */
    namespace Salesman {
      interface ListResponse extends Basic.ListResponse {
        salesmans: BasicInfo[]
      }
    }

    /* 尺码 */
    namespace Size {
      interface ListResponse extends Basic.ListResponse {
        sizes: BasicInfo[]
      }
    }

    /* 订单类型 */
    namespace OrderType {
      interface ListResponse extends Basic.ListResponse {
        order_types: BasicInfo[]
      }
    }

    /* 工序 */
    namespace Procedure {
      interface ListResponse extends Basic.ListResponse {
        procedures: BasicInfo[]
      }
    }
  }
}

