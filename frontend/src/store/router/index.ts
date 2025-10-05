import type { MenuOption } from 'naive-ui'
import { router } from '@/router'
import { staticRoutes } from '@/router/routes.static'
import { fetchUserRoutes } from '@/service'
import { useAuthStore } from '@/store/auth'
import { $t, local } from '@/utils'
import { createMenus, createRoutes, generateCacheRoutes } from './helper'

interface RoutesStatus {
  isInitAuthRoute: boolean
  menus: MenuOption[]
  rowRoutes: AppRoute.RowRoute[]
  activeMenu: string | null
  cacheRoutes: string[]
}
export const useRouteStore = defineStore('route-store', {
  state: (): RoutesStatus => {
    return {
      isInitAuthRoute: false,
      activeMenu: null,
      menus: [],
      rowRoutes: [],
      cacheRoutes: [],
    }
  },
  actions: {
    resetRouteStore() {
      this.resetRoutes()
      this.$reset()
    },
    resetRoutes() {
      if (router.hasRoute('appRoot'))
        router.removeRoute('appRoot')
    },
    // set the currently highlighted menu key
    setActiveMenu(key: string) {
      this.activeMenu = key
    },

    async initRouteInfo() {
      if (import.meta.env.VITE_ROUTE_LOAD_MODE === 'dynamic') {
        // Get user's route from JWT token (不传 id 参数，后端自动从 JWT 获取)
        const result = await fetchUserRoutes()

        if (!result.isSuccess || !result.data) {
          throw new Error('Failed to fetch user routes')
        }

        // 后端返回的是 { routes: [...] }，需要提取 routes 数组
        return result.data.routes
      }
      else {
        this.rowRoutes = staticRoutes
        return staticRoutes
      }
    },
    async initAuthRoute() {
      this.isInitAuthRoute = false

      try {
        // Initialize route information
        const rowRoutes = await this.initRouteInfo()
        if (!rowRoutes) {
          const error = new Error('Failed to get route information')
          // 不在这里显示错误消息，由调用者处理
          throw error
        }
        this.rowRoutes = rowRoutes

        // Generate actual route and insert
        const routes = createRoutes(rowRoutes)
        router.addRoute(routes)

        // Generate side menu
        this.menus = createMenus(rowRoutes)

        // Generate the route cache
        this.cacheRoutes = generateCacheRoutes(rowRoutes)

        this.isInitAuthRoute = true
      }
      catch (error) {
        // 重置状态并重新抛出错误，不显示消息
        this.isInitAuthRoute = false
        throw error
      }
    },
  },
})
