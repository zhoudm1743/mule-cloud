<template>
	<view class="mine-container">
	<!-- 用户信息卡片 -->
	<view class="user-card">
		<view class="user-avatar">
			<image :src="userInfo?.avatar || '/static/logo.png'" mode="aspectFill"></image>
		</view>
		<view class="user-info">
			<view class="user-name">{{ userInfo?.nickname || '微信用户' }}</view>
			<view class="user-phone" @click="handleBindPhone">
				{{ userInfo?.phone || '点击绑定手机号' }}
			</view>
		</view>
	</view>

		<!-- 当前企业 -->
		<view class="current-tenant" v-if="currentTenant">
			<view class="tenant-header">
				<text class="tenant-label">当前企业</text>
				<text class="switch-btn" @click="showTenantSelector = true">切换 ›</text>
			</view>
			<view class="tenant-info">
				<view class="tenant-name">🏢 {{ currentTenant.tenant_name }}</view>
				<view class="tenant-code">企业代码：{{ currentTenant.tenant_code }}</view>
			</view>
		</view>

	<!-- 功能菜单 -->
	<view class="menu-list">
		<!-- 绑定手机号按钮 - 需要小程序认证后才能使用 -->
		<!-- #ifdef MP-WEIXIN -->
		<button 
			v-if="!userInfo?.phone"
			class="menu-item menu-button" 
			open-type="getPhoneNumber"
			@getphonenumber="handleGetPhoneNumber"
		>
			<view class="menu-icon">📱</view>
			<view class="menu-text">绑定手机号</view>
			<view class="menu-tip">(需小程序认证)</view>
			<view class="menu-arrow">›</view>
		</button>
		<!-- #endif -->

		<view class="menu-item" @click="handleUnbindPhone" v-if="userInfo?.phone">
			<view class="menu-icon">📱</view>
			<view class="menu-text">解绑手机号</view>
			<view class="menu-value">{{ userInfo.phone }}</view>
			<view class="menu-arrow">›</view>
		</view>

		<view class="menu-item" @click="handleMyTenants">
			<view class="menu-icon">🏢</view>
			<view class="menu-text">我的企业</view>
			<view class="menu-badge" v-if="tenants.length > 1">{{ tenants.length }}</view>
			<view class="menu-arrow">›</view>
		</view>

		<view class="menu-item" @click="handleAbout">
			<view class="menu-icon">ℹ️</view>
			<view class="menu-text">关于我们</view>
			<view class="menu-arrow">›</view>
		</view>

		<view class="menu-item" @click="handleLogout">
			<view class="menu-icon">🚪</view>
			<view class="menu-text logout-text">退出登录</view>
		</view>
	</view>

		<!-- 租户选择器 -->
		<uv-popup v-model="showTenantSelector" mode="bottom" :round="20">
			<view class="tenant-selector">
				<view class="selector-header">
					<text class="selector-title">选择企业</text>
					<text class="selector-close" @click="showTenantSelector = false">✕</text>
				</view>
				<view class="tenant-list">
					<view 
						class="tenant-item"
						v-for="tenant in tenants"
						:key="tenant.tenant_id"
						:class="{ active: tenant.tenant_id === currentTenant?.tenant_id }"
						@click="handleSwitchTenant(tenant)"
					>
						<view class="tenant-name">{{ tenant.tenant_name }}</view>
						<view class="tenant-status">
							{{ tenant.status === 'active' ? '✅ 在职' : '❌ 已离职' }}
						</view>
					</view>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/modules/user'
import { switchTenant } from '@/api/auth'

const userStore = useUserStore()
const showTenantSelector = ref(false)

const userInfo = computed(() => userStore.userInfo || {})
const currentTenant = computed(() => userStore.currentTenant)
const tenants = computed(() => userStore.tenants || [])

onMounted(() => {
	// 刷新用户信息
	if (userStore.isLoggedIn) {
		userStore.fetchUserInfo().catch(err => {
			console.error('获取用户信息失败', err)
		})
	}
})

// 绑定手机号 - 显示弹窗引导用户点击按钮
const handleBindPhone = () => {
	// 小程序需要通过 button open-type="getPhoneNumber" 来获取手机号
	// 这里不做任何处理，由模板中的 button 来触发
}

// 处理微信手机号授权回调
const handleGetPhoneNumber = async (e) => {
	console.log('手机号授权回调', e)
	
	// 检查是否有错误
	if (e.detail.errMsg && e.detail.errMsg !== 'getPhoneNumber:ok') {
		console.error('获取手机号失败', e.detail.errMsg)
		
		// 友好的错误提示
		let errorMsg = '获取手机号失败'
		if (e.detail.errMsg.includes('auth deny')) {
			errorMsg = '您拒绝了授权'
		} else if (e.detail.errMsg.includes('verify')) {
			errorMsg = '需要验证微信账号，请在微信中完成验证后重试'
		} else if (e.detail.errMsg.includes('fail')) {
			errorMsg = '获取失败，请确保小程序已通过微信认证'
		}
		
		uni.showModal({
			title: '提示',
			content: errorMsg + '\n\n开发环境可能无法使用此功能，需要小程序认证后在真机上测试',
			showCancel: false
		})
		return
	}
	
	if (e.detail.code) {
		try {
			uni.showLoading({ title: '绑定中...' })
			const { bindPhone } = await import('@/api/auth')
			const res = await bindPhone(e.detail.code)
			
			if (res.success) {
				uni.showToast({
					title: '绑定成功',
					icon: 'success'
				})
				// 刷新用户信息
				await userStore.fetchUserInfo()
			} else {
				uni.showToast({
					title: res.message || '绑定失败',
					icon: 'none'
				})
			}
		} catch (error) {
			console.error('绑定手机号失败', error)
			uni.showToast({
				title: error.message || '绑定失败',
				icon: 'none'
			})
		} finally {
			uni.hideLoading()
		}
	} else {
		uni.showToast({
			title: '取消授权',
			icon: 'none'
		})
	}
}

// 解绑手机号
const handleUnbindPhone = () => {
	uni.showModal({
		title: '解绑手机号',
		content: '确定要解绑手机号吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '解绑中...' })
					const { unbindPhone } = await import('@/api/auth')
					await unbindPhone()
					
					uni.showToast({
						title: '解绑成功',
						icon: 'success'
					})
					// 刷新用户信息
					await userStore.fetchUserInfo()
				} catch (error) {
					console.error('解绑失败', error)
					uni.showToast({
						title: error.message || '解绑失败',
						icon: 'none'
					})
				} finally {
					uni.hideLoading()
				}
			}
		}
	})
}

// 我的企业
const handleMyTenants = () => {
	uni.showToast({
		title: '功能开发中',
		icon: 'none'
	})
}

// 切换企业
const handleSwitchTenant = async (tenant) => {
	if (tenant.tenant_id === currentTenant.value?.tenant_id) {
		showTenantSelector.value = false
		return
	}

	uni.showLoading({ title: '切换中...' })

	try {
		const res = await switchTenant(tenant.tenant_id)

		// 更新状态
		userStore.switchTenant({
			token: res.token,
			current_tenant: res.current_tenant
		})

		showTenantSelector.value = false
		uni.hideLoading()
		uni.showToast({
			title: '切换成功',
			icon: 'success'
		})

		// 重新加载页面
		setTimeout(() => {
			uni.reLaunch({
				url: '/pages/index/index'
			})
		}, 1000)
	} catch (error) {
		uni.hideLoading()
		console.error('切换企业失败', error)
	}
}

// 关于我们
const handleAbout = () => {
	uni.showModal({
		title: '关于我们',
		content: '智能工厂管理系统 v1.0.0\n\n服装生产全流程数字化管理',
		showCancel: false
	})
}

// 退出登录
const handleLogout = () => {
	uni.showModal({
		title: '提示',
		content: '确定要退出登录吗？',
		success: (res) => {
			if (res.confirm) {
				userStore.logout()
			}
		}
	})
}
</script>

<style lang="scss" scoped>
.mine-container {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 40rpx;
}

.user-card {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	padding: 60rpx 40rpx;
	display: flex;
	align-items: center;

	.user-avatar {
		width: 120rpx;
		height: 120rpx;
		border-radius: 60rpx;
		overflow: hidden;
		border: 4rpx solid rgba(255, 255, 255, 0.3);
		margin-right: 32rpx;

		image {
			width: 100%;
			height: 100%;
		}
	}

	.user-info {
		flex: 1;

		.user-name {
			font-size: 36rpx;
			font-weight: bold;
			color: #fff;
			margin-bottom: 12rpx;
		}

		.user-phone {
			font-size: 26rpx;
			color: rgba(255, 255, 255, 0.8);
		}
	}
}

.current-tenant {
	background: #fff;
	margin: 24rpx;
	padding: 32rpx;
	border-radius: 16rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

	.tenant-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;

		.tenant-label {
			font-size: 26rpx;
			color: #999;
		}

		.switch-btn {
			font-size: 28rpx;
			color: #1989fa;
		}
	}

	.tenant-info {
		.tenant-name {
			font-size: 32rpx;
			font-weight: bold;
			color: #333;
			margin-bottom: 12rpx;
		}

		.tenant-code {
			font-size: 26rpx;
			color: #999;
		}
	}
}

.menu-list {
	background: #fff;
	margin: 24rpx;
	border-radius: 16rpx;
	overflow: hidden;

	.menu-item {
		display: flex;
		align-items: center;
		padding: 32rpx;
		border-bottom: 1rpx solid #f5f5f5;
		transition: all 0.3s;

		&:last-child {
			border-bottom: none;
		}

		&:active {
			background: #f5f5f5;
		}

		.menu-icon {
			font-size: 40rpx;
			margin-right: 24rpx;
		}

		.menu-text {
			flex: 1;
			font-size: 30rpx;
			color: #333;

			&.logout-text {
				color: #ff4d4f;
			}
		}

		.menu-tip {
			font-size: 22rpx;
			color: #999;
			margin-right: 8rpx;
		}

		.menu-value {
			font-size: 26rpx;
			color: #999;
			margin-right: 16rpx;
		}

		.menu-badge {
			background: #ff4d4f;
			color: #fff;
			font-size: 22rpx;
			padding: 4rpx 12rpx;
			border-radius: 20rpx;
			margin-right: 16rpx;
		}

		.menu-arrow {
			font-size: 36rpx;
			color: #ccc;
		}
	}

	// 按钮样式重置（用于绑定手机号按钮）
	.menu-button {
		padding: 0;
		margin: 0;
		border: none;
		background: transparent;
		line-height: normal;
		text-align: left;
		border-radius: 0;
		width: 100%;
		display: flex;
		align-items: center;
		justify-content: flex-start;

		&::after {
			border: none;
		}
	}
}

.tenant-selector {
	background: #fff;
	border-radius: 20rpx 20rpx 0 0;
	padding-bottom: env(safe-area-inset-bottom);

	.selector-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 32rpx;
		border-bottom: 1rpx solid #f5f5f5;

		.selector-title {
			font-size: 32rpx;
			font-weight: bold;
			color: #333;
		}

		.selector-close {
			font-size: 40rpx;
			color: #999;
		}
	}

	.tenant-list {
		max-height: 600rpx;
		overflow-y: auto;

		.tenant-item {
			padding: 32rpx;
			border-bottom: 1rpx solid #f5f5f5;
			display: flex;
			justify-content: space-between;
			align-items: center;
			transition: all 0.3s;

			&:active {
				background: #f5f5f5;
			}

			&.active {
				background: #e6f7ff;
			}

			.tenant-name {
				font-size: 30rpx;
				color: #333;
			}

			.tenant-status {
				font-size: 24rpx;
				color: #52c41a;
			}
		}
	}
}
</style>

