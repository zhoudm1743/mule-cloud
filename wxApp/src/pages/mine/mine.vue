<template>
	<view class="mine-container">
		<!-- 顶部用户信息区域 -->
		<view class="user-section">
			<view class="user-info" @click="handleEditProfile">
				<view class="avatar-wrapper">
					<image :src="memberProfile?.avatar || userInfo?.avatar || '/static/logo.png'" mode="aspectFill"></image>
				</view>
				<view class="user-detail">
					<view class="user-name">{{ memberProfile?.name || userInfo?.nickname || '微信用户' }}</view>
					<view class="user-phone">工号：{{ memberProfile?.job_number || '未设置' }}</view>
				</view>
				<u-icon name="arrow-right" :size="32" color="#fff"></u-icon>
			</view>
		</view>

		<!-- 当前企业 -->
		<view class="tenant-section" v-if="currentTenant" @click="showTenantSelector = true">
			<view class="tenant-content">
				<view class="tenant-icon">
					<u-icon name="home-fill" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="tenant-info">
					<view class="tenant-name">{{ currentTenant.tenant_name }}</view>
					<view class="tenant-code">{{ currentTenant.tenant_code }}</view>
				</view>
				<view class="tenant-badge" v-if="tenants.length > 1">{{ tenants.length }}</view>
				<u-icon name="arrow-right" :size="28" color="#999"></u-icon>
			</view>
		</view>

	<!-- 功能列表 -->
	<view class="menu-list">
		<!-- 个人档案 -->
		<view class="menu-item" @click="handleProfile">
			<view class="item-icon">
				<u-icon name="file-text-fill" :size="36" color="#5EA3F2"></u-icon>
			</view>
			<view class="item-text">个人档案</view>
			<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
		</view>

		<view class="divider"></view>

		<!-- 手机号
		<view class="menu-item" @click="handleBindPhone" v-if="!userInfo?.phone">
			<view class="item-icon">
				<u-icon name="phone-fill" :size="36" color="#5EA3F2"></u-icon>
			</view>
			<view class="item-text">绑定手机号</view>
			<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
		</view>
		
	<view class="menu-item" @click="handleChangePhone" v-else>
		<view class="item-icon">
			<u-icon name="phone-fill" :size="36" color="#5EA3F2"></u-icon>
		</view>
		<view class="item-content">
			<view class="item-text">手机号</view>
			<view class="item-desc">{{ userInfo.phone }}</view>
		</view>
		<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
	</view>

		<view class="divider"></view> -->

			<!-- <view class="menu-item" @click="handleClearCache">
				<view class="item-icon">
					<u-icon name="trash-fill" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="item-text">清除缓存</view>
				<view class="item-value">{{ cacheSize }}</view>
				<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="handleCheckUpdate">
				<view class="item-icon">
					<u-icon name="reload" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="item-text">检查更新</view>
				<view class="item-value">v1.0.0</view>
				<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
			</view> -->

			<view class="divider"></view>

			<view class="menu-item" @click="handleFeedback">
				<view class="item-icon">
					<u-icon name="chat-fill" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="item-text">意见反馈</view>
				<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="handleHelp">
				<view class="item-icon">
					<u-icon name="question-circle-fill" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="item-text">帮助中心</view>
				<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
			</view>

			<view class="menu-item" @click="handleAbout">
				<view class="item-icon">
					<u-icon name="info-circle-fill" :size="36" color="#5EA3F2"></u-icon>
				</view>
				<view class="item-text">关于我们</view>
				<u-icon name="arrow-right" :size="28" color="#ccc"></u-icon>
			</view>
		</view>

		<!-- 退出登录 -->
		<view class="logout-section">
			<view class="logout-btn" @click="handleLogout">
				<u-icon name="close-circle" :size="36" color="#f56c6c" style="margin-right: 12rpx;"></u-icon>
				<text>退出登录</text>
			</view>
		</view>

	<!-- 租户选择弹窗 -->
	<u-popup :show="showTenantSelector" @close="showTenantSelector = false" mode="bottom" :round="20">
		<view class="popup-container">
			<view class="popup-header">
				<text class="popup-title">选择企业</text>
				<view class="popup-close" @click="showTenantSelector = false">
					<u-icon name="close" :size="36" color="#999"></u-icon>
				</view>
			</view>
				<view class="popup-list">
					<view 
						class="popup-item"
						v-for="tenant in tenants"
						:key="tenant.tenant_id"
						:class="{ active: tenant.tenant_id === currentTenant?.tenant_id }"
						@click="handleSwitchTenant(tenant)"
					>
						<view class="popup-item-info">
							<view class="popup-item-name">{{ tenant.tenant_name }}</view>
							<view class="popup-item-code">{{ tenant.tenant_code }}</view>
						</view>
						<u-icon 
							v-if="tenant.tenant_id === currentTenant?.tenant_id"
							name="checkmark-circle-fill" 
							:size="48" 
							color="#5EA3F2"
						></u-icon>
					</view>
				</view>
			</view>
		</u-popup>

		
		<!-- 底部导航 -->
		<TabBar :current="2" />
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/modules/user'
import { switchTenant, updateUserInfo, bindPhone as bindPhoneAPI, unbindPhone as unbindPhoneAPI } from '@/api/auth'
import { getProfile } from '@/api/member'
import TabBar from '@/components/TabBar/TabBar.vue'

const userStore = useUserStore()
const showTenantSelector = ref(false)
const cacheSize = ref('0KB')
const memberProfile = ref(null)

const userInfo = computed(() => userStore.userInfo || {})
const currentTenant = computed(() => userStore.currentTenant)
const tenants = computed(() => userStore.tenants || [])

onMounted(async () => {
	if (userStore.isLoggedIn) {
		await userStore.fetchUserInfo().catch(err => {
			console.error('获取用户信息失败', err)
		})
		// 获取员工档案信息
		await loadMemberProfile()
	}
	getCacheSize()
})

// 加载员工档案（静默加载）
const loadMemberProfile = async () => {
	try {
		const res = await getProfile(true)  // 静默加载
		if (res.code === 0) {
			memberProfile.value = res.data
		}
	} catch (error) {
		// 静默失败，不影响页面展示
		console.log('暂未获取到员工档案信息')
	}
}

const getCacheSize = () => {
	try {
		const res = uni.getStorageInfoSync()
		const size = res.currentSize || 0
		cacheSize.value = size < 1024 ? size + 'KB' : (size / 1024).toFixed(2) + 'MB'
	} catch (e) {
		console.error('获取缓存大小失败', e)
	}
}

const handleEditProfile = () => {
	uni.navigateTo({
		url: '/pages/edit-profile/edit-profile'
	})
}

const handleProfile = () => {
	uni.navigateTo({
		url: '/pages/member/profile/profile'
	})
}

const handleBindPhone = () => {
	uni.showModal({
		title: '绑定手机号',
		content: '此功能需要在小程序认证后使用微信授权获取',
		showCancel: false
	})
}

const handleChangePhone = () => {
	uni.showModal({
		title: '更换手机号',
		content: '请先解绑当前手机号，再绑定新手机号',
		showCancel: true,
		confirmText: '去解绑',
		success: (res) => {
			if (res.confirm) handleUnbindPhone()
		}
	})
}

const handleUnbindPhone = async () => {
	uni.showModal({
		title: '解绑手机号',
		content: '确定要解绑手机号吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '解绑中...' })
					await unbindPhoneAPI()
					await userStore.fetchUserInfo()
					uni.showToast({ title: '解绑成功', icon: 'success' })
				} catch (error) {
					uni.showToast({ title: error.message || '解绑失败', icon: 'none' })
				} finally {
					uni.hideLoading()
				}
			}
		}
	})
}

const handleSwitchTenant = async (tenant) => {
	if (tenant.tenant_id === currentTenant.value?.tenant_id) {
		showTenantSelector.value = false
		return
	}

	uni.showLoading({ title: '切换中...' })
	try {
		const res = await switchTenant(tenant.tenant_id)
		userStore.switchTenant({
			token: res.token,
			current_tenant: res.current_tenant
		})
		showTenantSelector.value = false
		uni.hideLoading()
		uni.showToast({ title: '切换成功', icon: 'success' })
		setTimeout(() => {
			uni.reLaunch({ url: '/pages/index/index' })
		}, 1000)
	} catch (error) {
		uni.hideLoading()
		uni.showToast({ title: '切换失败', icon: 'none' })
	}
}

const handleClearCache = () => {
	uni.showModal({
		title: '清除缓存',
		content: '确定要清除缓存吗？',
		success: (res) => {
			if (res.confirm) {
				try {
					const token = uni.getStorageSync('token')
					const userInfoData = uni.getStorageSync('userInfo')
					const currentTenantData = uni.getStorageSync('currentTenant')
					const tenantsData = uni.getStorageSync('tenants')
					
					uni.clearStorageSync()
					uni.setStorageSync('token', token)
					uni.setStorageSync('userInfo', userInfoData)
					uni.setStorageSync('currentTenant', currentTenantData)
					uni.setStorageSync('tenants', tenantsData)
					
					getCacheSize()
					uni.showToast({ title: '清除成功', icon: 'success' })
				} catch (e) {
					uni.showToast({ title: '清除失败', icon: 'none' })
				}
			}
		}
	})
}

const handleCheckUpdate = () => {
	uni.showToast({ title: '已是最新版本', icon: 'none' })
}

const handleFeedback = () => {
	uni.showToast({ title: '功能开发中', icon: 'none' })
}

const handleHelp = () => {
	uni.showToast({ title: '功能开发中', icon: 'none' })
}

const handleAbout = () => {
	uni.showModal({
		title: '关于我们',
		content: '智能工厂管理系统\n版本: v1.0.0\n\n服装生产全流程数字化管理',
		showCancel: false
	})
}

const handleLogout = () => {
	uni.showModal({
		title: '退出登录',
		content: '确定要退出登录吗？',
		success: (res) => {
			if (res.confirm) userStore.logout()
		}
	})
}
</script>

<style scoped lang="scss">
.mine-container {
	min-height: 100vh;
	background: #f5f7fa;
	padding-bottom: 120rpx;
}

// 顶部用户区域
.user-section {
	background: linear-gradient(180deg, #5EA3F2 0%, #4FC3F7 100%);
	padding: 40rpx 32rpx 60rpx;
}

.user-info {
	display: flex;
	align-items: center;
	padding: 32rpx;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 24rpx;
	backdrop-filter: blur(10rpx);
}

.avatar-wrapper {
	width: 120rpx;
	height: 120rpx;
	border-radius: 50%;
	overflow: hidden;
	margin-right: 24rpx;
	border: 4rpx solid rgba(255, 255, 255, 0.3);
	
	image {
		width: 100%;
		height: 100%;
	}
}

.user-detail {
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

// 企业卡片
.tenant-section {
	margin: -40rpx 32rpx 24rpx;
	background: #fff;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
}

.tenant-content {
	display: flex;
	align-items: center;
	padding: 32rpx;
}

.tenant-icon {
	width: 88rpx;
	height: 88rpx;
	background: #f0f8ff;
	border-radius: 16rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 24rpx;
}

.tenant-info {
	flex: 1;
	
	.tenant-name {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 8rpx;
	}
	
	.tenant-code {
		font-size: 24rpx;
		color: #999;
	}
}

.tenant-badge {
	min-width: 40rpx;
	height: 40rpx;
	line-height: 40rpx;
	padding: 0 12rpx;
	background: #f56c6c;
	color: #fff;
	font-size: 22rpx;
	border-radius: 20rpx;
	text-align: center;
	margin-right: 16rpx;
}

// 菜单列表
.menu-list {
	margin: 0 32rpx 24rpx;
	background: #fff;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
}

.menu-item {
	display: flex;
	align-items: center;
	padding: 32rpx;
	background: #fff;
	
	&:active {
		background: #f8f9fa;
	}
}


.item-icon {
	margin-right: 24rpx;
}

.item-text {
	flex: 1;
	font-size: 30rpx;
	color: #333;
}

.item-content {
	flex: 1;
	
	.item-text {
		font-size: 30rpx;
		color: #333;
		margin-bottom: 6rpx;
	}
	
	.item-desc {
		font-size: 24rpx;
		color: #999;
	}
}

.item-value {
	font-size: 26rpx;
	color: #999;
	margin-right: 16rpx;
}

.divider {
	height: 1rpx;
	background: #f0f0f0;
	margin: 0 32rpx;
}

// 退出登录
.logout-section {
	padding: 0 32rpx 32rpx;
}

.logout-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 32rpx;
	background: #fff;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
	
	text {
		font-size: 32rpx;
		color: #f56c6c;
		font-weight: bold;
	}
	
	&:active {
		opacity: 0.8;
	}
}

// 弹窗通用样式
.popup-container {
	padding: 32rpx 32rpx 64rpx;
}

.popup-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding-bottom: 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
	
	.popup-title {
		font-size: 36rpx;
		font-weight: bold;
		color: #333;
	}
}

.popup-list {
	max-height: 600rpx;
	overflow-y: auto;
}

.popup-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32rpx 0;
	border-bottom: 1rpx solid #f5f7fa;
	
	&:last-child {
		border-bottom: none;
	}
	
	&.active .popup-item-name {
		color: #5EA3F2;
		font-weight: bold;
	}
}

.popup-item-info {
	flex: 1;
	
	.popup-item-name {
		font-size: 32rpx;
		color: #333;
		margin-bottom: 8rpx;
	}
	
	.popup-item-code {
		font-size: 24rpx;
		color: #999;
	}
}

// 编辑表单
.edit-form {
	padding: 32rpx 0;
}

.form-row {
	display: flex;
	align-items: center;
	padding: 24rpx 0;
	
	.form-label {
		width: 120rpx;
		font-size: 30rpx;
		color: #666;
	}
	
	.form-value {
		flex: 1;
	}
	
	.form-input {
		flex: 1;
		padding: 20rpx 24rpx;
		background: #f5f7fa;
		border-radius: 12rpx;
		font-size: 30rpx;
	}
}

.avatar-edit {
	position: relative;
	width: 160rpx;
	height: 160rpx;
	border-radius: 16rpx;
	overflow: hidden;
	
	image {
		width: 100%;
		height: 100%;
	}
	
	.avatar-mask {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 60rpx;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
	}
}

.popup-footer {
	display: flex;
	gap: 20rpx;
	padding-top: 32rpx;
}

.footer-btn {
	flex: 1;
	height: 88rpx;
	border-radius: 44rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 32rpx;
	
	&.cancel-btn {
		background: #f5f7fa;
		color: #666;
	}
	
	&.confirm-btn {
		background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
		color: #fff;
		font-weight: bold;
	}
	
	&:active {
		opacity: 0.8;
	}
}
</style>

