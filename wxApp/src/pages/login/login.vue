<template>
	<view class="login-container">
		<!-- 装饰圆圈 -->
		<view class="circle circle-1"></view>
		<view class="circle circle-2"></view>
		
	<view class="logo-wrapper">
		<view class="logo-shadow">
			<view class="logo-icon">
				<view class="factory-icon">
					<view class="factory-roof"></view>
					<view class="factory-body">
						<view class="factory-window"></view>
						<view class="factory-window"></view>
					</view>
				</view>
			</view>
		</view>
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
				<view class="wechat-icon">
					<view class="wechat-bubble wechat-bubble-left"></view>
					<view class="wechat-bubble wechat-bubble-right"></view>
				</view>
				<text>微信一键登录</text>
			</button>

			<view class="agreement">
				<checkbox-group @change="onAgreementChange">
					<label>
						<checkbox :checked="agreed" color="#5EA3F2" />
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
			// 直接登录成功 - 保存token和用户信息
			userStore.setLoginInfo({
				token: res.token,
				user_info: res.user_info,
				current_tenant: res.current_tenant,
				tenants: res.tenants || [res.current_tenant]
			})

			uni.showToast({
				title: '登录成功',
				icon: 'success'
			})

		setTimeout(() => {
			uni.reLaunch({
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
	background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
	display: flex;
	flex-direction: column;
	padding: 40rpx;
	position: relative;
	overflow: hidden;
}

/* 装饰圆圈 */
.circle {
	position: absolute;
	border-radius: 50%;
	background: rgba(255, 255, 255, 0.1);
	
	&.circle-1 {
		width: 300rpx;
		height: 300rpx;
		top: -100rpx;
		left: -100rpx;
	}
	
	&.circle-2 {
		width: 200rpx;
		height: 200rpx;
		bottom: 100rpx;
		right: -50rpx;
	}
}

.logo-wrapper {
	text-align: center;
	margin-top: 100rpx;
	margin-bottom: 80rpx;
	z-index: 1;

	.logo-shadow {
		width: 200rpx;
		height: 200rpx;
		margin: 0 auto 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(255, 255, 255, 0.95);
		border-radius: 50%;
		box-shadow: 0 12rpx 48rpx rgba(0, 0, 0, 0.15);
		
		.logo-icon {
			width: 100%;
			height: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
		}
		
		.factory-icon {
			width: 100rpx;
			height: 100rpx;
			position: relative;
			
			.factory-roof {
				width: 0;
				height: 0;
				border-left: 50rpx solid transparent;
				border-right: 50rpx solid transparent;
				border-bottom: 30rpx solid #5EA3F2;
				position: absolute;
				top: 10rpx;
				left: 0;
			}
			
			.factory-body {
				width: 80rpx;
				height: 50rpx;
				background: #5EA3F2;
				position: absolute;
				bottom: 10rpx;
				left: 10rpx;
				border-radius: 4rpx;
				display: flex;
				justify-content: space-around;
				align-items: center;
				padding: 10rpx;
				box-sizing: border-box;
				
				.factory-window {
					width: 20rpx;
					height: 20rpx;
					background: rgba(255, 255, 255, 0.9);
					border-radius: 2rpx;
				}
			}
		}
	}

	.app-name {
		display: block;
		margin-top: 40rpx;
		font-size: 44rpx;
		font-weight: 600;
		color: #fff;
		letter-spacing: 2rpx;
	}

	.app-desc {
		display: block;
		margin-top: 20rpx;
		font-size: 28rpx;
		color: rgba(255, 255, 255, 0.95);
		letter-spacing: 1rpx;
	}
}

.login-card {
	background: #fff;
	border-radius: 24rpx;
	padding: 60rpx 40rpx;
	box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.1);

	.welcome-text {
		font-size: 48rpx;
		font-weight: 600;
		color: #333;
		margin-bottom: 16rpx;
		letter-spacing: 1rpx;
	}

	.tips-text {
		font-size: 28rpx;
		color: #999;
		margin-bottom: 60rpx;
		letter-spacing: 0.5rpx;
	}

	.login-btn {
		width: 100%;
		height: 100rpx;
		line-height: 100rpx;
		background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
		border-radius: 50rpx;
		border: none;
		font-size: 34rpx;
		font-weight: 500;
		color: #fff;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 8rpx 24rpx rgba(94, 163, 242, 0.4);
		letter-spacing: 2rpx;
		transition: all 0.3s ease;

		.wechat-icon {
			width: 40rpx;
			height: 40rpx;
			margin-right: 16rpx;
			position: relative;
			display: flex;
			align-items: center;
			justify-content: center;
			
			.wechat-bubble {
				position: absolute;
				background: #fff;
				border-radius: 50%;
				
				&-left {
					width: 24rpx;
					height: 24rpx;
					left: 0;
					top: 8rpx;
					
					&::before {
						content: '';
						position: absolute;
						width: 6rpx;
						height: 6rpx;
						background: #5EA3F2;
						border-radius: 50%;
						left: 6rpx;
						top: 6rpx;
					}
					
					&::after {
						content: '';
						position: absolute;
						width: 4rpx;
						height: 4rpx;
						background: #5EA3F2;
						border-radius: 50%;
						left: 14rpx;
						top: 8rpx;
					}
				}
				
				&-right {
					width: 20rpx;
					height: 20rpx;
					right: 0;
					top: 0;
					
					&::before {
						content: '';
						position: absolute;
						width: 5rpx;
						height: 5rpx;
						background: #5EA3F2;
						border-radius: 50%;
						left: 5rpx;
						top: 5rpx;
					}
					
					&::after {
						content: '';
						position: absolute;
						width: 3rpx;
						height: 3rpx;
						background: #5EA3F2;
						border-radius: 50%;
						left: 11rpx;
						top: 6rpx;
					}
				}
			}
		}
		
		&:active {
			transform: scale(0.98);
		}
		
		&::after {
			border: none;
		}
	}

	.agreement {
		margin-top: 40rpx;
		font-size: 24rpx;
		color: #999;
		line-height: 1.6;

		.agreement-text {
			margin-left: 8rpx;
		}

		.link {
			color: #5EA3F2;
			font-weight: 500;
		}
	}
}

.footer {
	margin-top: auto;
	text-align: center;
	padding: 40rpx 0;

	.copyright {
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.8);
	}
}
</style>
