<template>
	<view class="login-container">
		<view class="logo-wrapper">
			<image class="logo" src="/static/logo.png" mode="aspectFit"></image>
			<text class="app-name">智能工厂管理系统</text>
			<text class="app-desc">服装生产全流程数字化管理</text>
		</view>

		<view class="login-card">
			<view class="welcome-text">欢迎使用</view>
			<view class="tips-text">请使用微信授权登录</view>

			<button 
				class="login-btn" 
				type="primary" 
				@click="handleWechatLogin"
				:loading="loading"
			>
				<text class="icon">📱</text>
				<text>微信一键登录</text>
			</button>

			<view class="agreement">
				<checkbox-group @change="onAgreementChange">
					<label>
						<checkbox :checked="agreed" color="#1989fa" />
						<text class="agreement-text">
							我已阅读并同意
							<text class="link" @click.stop="showAgreement('user')">《用户协议》</text>
							和
							<text class="link" @click.stop="showAgreement('privacy')">《隐私政策》</text>
						</text>
					</label>
				</checkbox-group>
			</view>
		</view>

		<view class="footer">
			<text class="copyright">© 2025 智能工厂管理系统</text>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/store/modules/user'
import { bindTenant, selectTenant } from '@/api/auth'

const userStore = useUserStore()
const loading = ref(false)
const agreed = ref(false)

// 处理微信登录
const handleWechatLogin = async () => {
	if (!agreed.value) {
		uni.showToast({
			title: '请先阅读并同意用户协议',
			icon: 'none'
		})
		return
	}

	loading.value = true
	try {
		// 获取用户信息（新版API）
		let userProfile = null
		try {
			// #ifdef MP-WEIXIN
			const profileRes = await uni.getUserProfile({
				desc: '用于完善用户资料',
				lang: 'zh_CN'
			})
			userProfile = profileRes.userInfo
			console.log('获取用户信息成功', userProfile)
			// #endif
		} catch (e) {
			console.log('用户拒绝授权或获取失败，继续登录', e)
		}
		
		// 调用登录
		const res = await userStore.login(userProfile)

		// 根据响应处理不同情况
		if (res.need_bind_tenant) {
			// 需要绑定租户
			uni.navigateTo({
				url: '/pages/bind-tenant/bind-tenant?userId=' + res.user_info.id
			})
		} else if (res.need_select_tenant) {
			// 需要选择租户
			uni.navigateTo({
				url: '/pages/select-tenant/select-tenant?userId=' + res.user_info.id + '&tenants=' + JSON.stringify(res.tenants)
			})
		} else {
			// 直接登录成功
			userStore.setLoginInfo({
				token: res.token,
				user_info: res.user_info,
				current_tenant: res.current_tenant
			})

			uni.showToast({
				title: '登录成功',
				icon: 'success'
			})

			setTimeout(() => {
				uni.switchTab({
					url: '/pages/index/index'
				})
			}, 1000)
		}
	} catch (error) {
		console.error('登录失败', error)
		uni.showToast({
			title: error.message || '登录失败',
			icon: 'none'
		})
	} finally {
		loading.value = false
	}
}

// 获取用户信息（旧版，已废弃）
const handleGetUserInfo = (e) => {
	console.log('用户信息', e.detail)
}

// 同意协议
const onAgreementChange = (e) => {
	agreed.value = e.detail.value.length > 0
}

// 显示协议
const showAgreement = (type) => {
	uni.showModal({
		title: type === 'user' ? '用户协议' : '隐私政策',
		content: '这里是协议内容...',
		confirmText: '我知道了'
	})
}
</script>

<style lang="scss" scoped>
.login-container {
	min-height: 100vh;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	display: flex;
	flex-direction: column;
	padding: 40rpx;
}

.logo-wrapper {
	text-align: center;
	margin-top: 100rpx;
	margin-bottom: 80rpx;

	.logo {
		width: 160rpx;
		height: 160rpx;
		border-radius: 20rpx;
		box-shadow: 0 8rpx 16rpx rgba(0, 0, 0, 0.1);
	}

	.app-name {
		display: block;
		margin-top: 30rpx;
		font-size: 40rpx;
		font-weight: bold;
		color: #fff;
	}

	.app-desc {
		display: block;
		margin-top: 16rpx;
		font-size: 28rpx;
		color: rgba(255, 255, 255, 0.8);
	}
}

.login-card {
	background: #fff;
	border-radius: 24rpx;
	padding: 60rpx 40rpx;
	box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.1);

	.welcome-text {
		font-size: 48rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 16rpx;
	}

	.tips-text {
		font-size: 28rpx;
		color: #999;
		margin-bottom: 60rpx;
	}

	.login-btn {
		width: 100%;
		height: 96rpx;
		line-height: 96rpx;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 48rpx;
		border: none;
		font-size: 32rpx;
		color: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 8rpx 16rpx rgba(102, 126, 234, 0.4);

		.icon {
			font-size: 40rpx;
			margin-right: 16rpx;
		}
	}

	.agreement {
		margin-top: 40rpx;
		font-size: 24rpx;
		color: #999;

		.agreement-text {
			margin-left: 8rpx;
		}

		.link {
			color: #1989fa;
		}
	}
}

.footer {
	margin-top: auto;
	text-align: center;
	padding: 40rpx 0;

	.copyright {
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.6);
	}
}
</style>

