import { router } from '@/router'
import { fetchLogin } from '@/service'
import { local } from '@/utils'
import { useRouteStore } from './router'
import { useTabStore } from './tab'

interface AuthStatus {
  userInfo: Api.Login.Info | null
  token: string
}
export const useAuthStore = defineStore('auth-store', {
  state: (): AuthStatus => {
    return {
      userInfo: local.get('userInfo'),
      token: local.get('accessToken') || '',
    }
  },
  getters: {
    /** 是否登录 */
    isLogin(state) {
      return Boolean(state.token)
    },
  },
  actions: {
    /* 登录退出，重置用户信息等 */
    async logout() {
      const route = unref(router.currentRoute)
      // 清除本地缓存
      this.clearAuthStorage()
      // 清空路由、菜单等数据
      const routeStore = useRouteStore()
      routeStore.resetRouteStore()
      // 清空标签栏数据
      const tabStore = useTabStore()
      tabStore.clearAllTabs()
      // 重置当前存储库
      this.$reset()
      // 重定向到登录页（如果当前不在登录页）
      if (route.name !== 'login') {
        // 如果当前页面需要认证，则保存重定向地址
        const redirect = route.meta.requiresAuth ? route.fullPath : undefined
        router.push({
          name: 'login',
          query: redirect ? { redirect } : undefined,
        })
      }
    },
    clearAuthStorage() {
      local.remove('accessToken')
      local.remove('refreshToken')
      local.remove('userInfo')
      local.remove('selected_tenant_id') // 清除系统管理员选择的租户上下文
    },

    /* 用户登录 */
    async login(phone: string, password: string, tenantCode?: string) {
      try {
        const loginData: any = { phone, password }
        if (tenantCode) {
          loginData.tenant_code = tenantCode
        }
        const { isSuccess, data } = await fetchLogin(loginData)
        if (!isSuccess)
          return

        // 处理登录信息
        await this.handleLoginInfo(data)
      }
      catch (e) {
        console.warn('[Login Error]:', e)
      }
    },

    /* 处理登录返回的数据 */
    async handleLoginInfo(data: Api.Login.Info) {
      // 将token和userInfo保存下来（后端返回的是 token，不是 accessToken）
      local.set('userInfo', data)
      local.set('accessToken', data.token) // 使用 data.token
      local.set('refreshToken', data.token) // 后端目前只返回一个token
      this.token = data.token
      this.userInfo = data

      // ✅ 系统管理员登录时，清除之前选择的租户上下文
      if (!data.tenant_id) {
        local.remove('selected_tenant_id')
      }

      // 添加路由和菜单
      const routeStore = useRouteStore()
      await routeStore.initAuthRoute()

      // 进行重定向跳转
      const route = unref(router.currentRoute)
      const query = route.query as { redirect: string }
      router.push({
        path: query.redirect || '/',
      })
    },
  },
})
