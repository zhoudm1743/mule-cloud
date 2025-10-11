/**
 * 用户状态管理
 */
import { defineStore } from 'pinia'
import { wechatLogin, getUserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', {
	state: () => ({
		token: uni.getStorageSync('token') || '',
		userInfo: uni.getStorageSync('userInfo') || null,
		tenants: uni.getStorageSync('tenants') || [],
		currentTenant: uni.getStorageSync('currentTenant') || null
	}),

	getters: {
		// 是否已登录
		isLoggedIn: (state) => !!state.token,
		
		// 是否有租户
		hasTenant: (state) => !!state.currentTenant,
		
		// 当前租户名称
		currentTenantName: (state) => state.currentTenant?.tenant_name || ''
	},

	actions: {
		/**
		 * 微信登录
		 * @param {Object} userProfile 用户信息（可选）
		 */
		async login(userProfile = null) {
			try {
				// 获取微信登录code
				const { code } = await uni.login({
					provider: 'weixin'
				})

				// 准备登录参数
				const loginParams = { code }
				
				// 如果有用户信息，添加到请求中
				if (userProfile) {
					loginParams.nickname = userProfile.nickName
					loginParams.avatar = userProfile.avatarUrl
					loginParams.gender = userProfile.gender
					loginParams.country = userProfile.country
					loginParams.province = userProfile.province
					loginParams.city = userProfile.city
				}

				// 调用后端登录接口
				const res = await wechatLogin(loginParams)

				// 保存用户信息
				this.userInfo = res.user_info
				uni.setStorageSync('userInfo', res.user_info)

				return res
			} catch (error) {
				console.error('登录失败', error)
				throw error
			}
		},

		/**
		 * 设置Token和租户信息
		 */
		setLoginInfo({ token, user_info, current_tenant, tenants }) {
			this.token = token
			this.userInfo = user_info
			this.currentTenant = current_tenant
			this.tenants = tenants || []

			// 持久化存储
			uni.setStorageSync('token', token)
			uni.setStorageSync('userInfo', user_info)
			uni.setStorageSync('currentTenant', current_tenant)
			uni.setStorageSync('tenants', tenants || [])
		},

		/**
		 * 切换租户
		 */
		switchTenant({ token, current_tenant }) {
			this.token = token
			this.currentTenant = current_tenant

			// 持久化存储
			uni.setStorageSync('token', token)
			uni.setStorageSync('currentTenant', current_tenant)
		},

		/**
		 * 更新用户信息
		 */
		updateUserInfo(userInfo) {
			this.userInfo = { ...this.userInfo, ...userInfo }
			uni.setStorageSync('userInfo', this.userInfo)
		},

		/**
		 * 获取用户完整信息
		 */
		async fetchUserInfo() {
			try {
				const res = await getUserInfo()
				this.userInfo = res.user_info
				this.tenants = res.tenants || []
				this.currentTenant = res.current_tenant

				uni.setStorageSync('userInfo', res.user_info)
				uni.setStorageSync('tenants', res.tenants)
				uni.setStorageSync('currentTenant', res.current_tenant)

				return res
			} catch (error) {
				console.error('获取用户信息失败', error)
				throw error
			}
		},

		/**
		 * 退出登录
		 */
		logout() {
			this.token = ''
			this.userInfo = null
			this.currentTenant = null
			this.tenants = []

			// 清除存储
			uni.removeStorageSync('token')
			uni.removeStorageSync('userInfo')
			uni.removeStorageSync('currentTenant')
			uni.removeStorageSync('tenants')

			// 跳转到登录页
			uni.reLaunch({
				url: '/pages/login/login'
			})
		}
	},

	// 持久化
	persist: true
})

