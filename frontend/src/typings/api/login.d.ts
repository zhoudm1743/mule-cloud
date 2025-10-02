/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Login {
    /* 登录返回的用户信息 - 适配后端 auth 服务 */
    interface Info {
      /** JWT Token */
      token: string
      /** 用户ID */
      user_id: string
      /** 租户ID */
      tenant_id?: string
      /** 手机号 */
      phone: string
      /** 昵称 */
      nickname: string
      /** 头像 */
      avatar: string
      /** 用户角色 */
      role: string[]
      /** 菜单权限映射：{"admin": ["read", "create"], "finance": ["read", "pending"]} */
      menu_permissions?: Record<string, string[]>
      /** Token过期时间戳 */
      expires_at: number
    }

    /* 注册请求 */
    interface RegisterRequest {
      phone: string
      password: string
      nickname: string
      email?: string
    }

    /* 注册响应 */
    interface RegisterResponse {
      user_id: string
      phone: string
      nickname: string
      message: string
    }

    /* 刷新Token响应 */
    interface RefreshTokenResponse {
      token: string
      expires_at: number
    }

    /* 个人资料 */
    interface Profile {
      user_id: string
      tenant_id?: string
      phone: string
      nickname: string
      avatar: string
      email: string
      role: string[]
      menu_permissions?: Record<string, string[]>
      status: number
      created_at: number
      updated_at: number
    }
  }
}
