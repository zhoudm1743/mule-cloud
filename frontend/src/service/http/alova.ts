import { local } from '@/utils'
import { createAlova } from 'alova'
import { createServerTokenAuthentication } from 'alova/client'
import adapterFetch from 'alova/fetch'
import VueHook from 'alova/vue'
import type { VueHookType } from 'alova/vue'
import {
  DEFAULT_ALOVA_OPTIONS,
  DEFAULT_BACKEND_OPTIONS,
} from './config'
import {
  handleBusinessError,
  handleRefreshToken,
  handleResponseError,
  handleServiceResult,
} from './handle'

const { onAuthRequired, onResponseRefreshToken } = createServerTokenAuthentication<VueHookType>({
  // 服务端判定token过期
  refreshTokenOnSuccess: {
    // 当服务端返回401时，表示token过期
    isExpired: async (response, method) => {
      const res = await response.clone().json()

      const isExpired = method.meta && method.meta.isExpired
      const authRole = method.meta && method.meta.authRole
      
      // 排除刷新token请求、登录请求和已标记为过期的请求
      if (authRole === 'refreshToken' || authRole === null || isExpired) {
        return false
      }
      
      return response.status === 401 || res.code === 401
    },

    // 当token过期时触发，在此函数中触发刷新token
    handler: async (_response, method) => {
      // 此处采取限制，防止过期请求无限循环重发
      if (!method.meta)
        method.meta = { isExpired: true }
      else
        method.meta.isExpired = true

      await handleRefreshToken()
    },
  },
  // 添加token到请求头
  assignToken: (method) => {
    method.config.headers.Authorization = `Bearer ${local.get('accessToken')}`
  },
})

// docs path of alova.js https://alova.js.org/
export function createAlovaInstance(
  alovaConfig: Service.AlovaConfig,
  backendConfig?: Service.BackendConfig,
) {
  const _backendConfig = { ...DEFAULT_BACKEND_OPTIONS, ...backendConfig }
  const _alovaConfig = { ...DEFAULT_ALOVA_OPTIONS, ...alovaConfig }

  return createAlova({
    statesHook: VueHook,
    requestAdapter: adapterFetch(),
    cacheFor: null,
    baseURL: _alovaConfig.baseURL,
    timeout: _alovaConfig.timeout,

    beforeRequest: onAuthRequired((method) => {
      if (method.meta?.isFormPost) {
        method.config.headers['Content-Type'] = 'application/x-www-form-urlencoded'
        method.data = new URLSearchParams(method.data as URLSearchParams).toString()
      }
      
      // ✅ 添加租户上下文 header（如果已选择租户）
      // 后端的 TenantContextMiddleware 会验证用户是否有权限切换租户
      // 这样可以避免前端时序问题（userInfo 可能还未从 localStorage 恢复）
      const selectedTenantCode = local.get('selected_tenant_code')
      if (selectedTenantCode) {
        method.config.headers['X-Tenant-Context'] = selectedTenantCode
      }
      
      alovaConfig.beforeRequest?.(method)
    }),
    responded: onResponseRefreshToken({
      // 请求成功的拦截器
      onSuccess: async (response, method) => {
        const { status } = response

        if (status === 200) {
          // 返回blob数据
          if (method.meta?.isBlob)
            return response.blob()

          // 返回json数据
          const apiData = await response.json()
          // 请求成功
          if (apiData[_backendConfig.codeKey] === _backendConfig.successCode)
            return handleServiceResult(apiData)

          // 业务请求失败
          const errorResult = handleBusinessError(apiData, _backendConfig)
          return handleServiceResult(errorResult, false)
        }
        // 接口请求失败
        const errorResult = handleResponseError(response)
        return handleServiceResult(errorResult, false)
      },
      onError: (error, method) => {
        // 如果是登录、注册或刷新token请求失败，不显示警告（会由业务层处理）
        const authRole = method.meta && method.meta.authRole
        if (authRole === 'refreshToken' || authRole === null) {
          return
        }
        const tip = `[${method.type}] - [${method.url}] - ${error.message}`
        window.$message?.warning(tip)
      },

      onComplete: async (_method) => {
        // 处理请求完成逻辑
      },
    }),
  })
}
