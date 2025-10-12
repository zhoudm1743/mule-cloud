<template>
	<view class="bind-tenant-container">
		<view class="header">
			<view class="title">绑定企业</view>
			<view class="desc">请输入企业邀请码加入企业</view>
		</view>

		<view class="form-card">
			<view class="form-item">
				<view class="label">企业邀请码</view>
				<input 
					class="input" 
					v-model="inviteCode" 
					placeholder="请输入企业邀请码"
					placeholder-class="placeholder"
					maxlength="20"
				/>
			</view>

			<view class="tips">
				<text>💡 邀请码可从企业管理员处获取</text>
			</view>

			<button 
				class="submit-btn" 
				type="primary"
				@click="handleBind"
				:loading="loading"
				:disabled="!inviteCode"
			>
				确认绑定
			</button>
		</view>

		<view class="help-text">
			<text>没有邀请码？</text>
			<text class="link" @click="contactAdmin">联系管理员获取</text>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { bindTenant } from '@/api/auth'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const userId = ref('')
const inviteCode = ref('')
const loading = ref(false)

// 获取页面参数
onMounted(() => {
	const pages = getCurrentPages()
	const currentPage = pages[pages.length - 1]
	const options = currentPage.options || {}
	userId.value = options.userId || ''
})

// 绑定租户
const handleBind = async () => {
	if (!inviteCode.value) {
		uni.showToast({
			title: '请输入邀请码',
			icon: 'none'
		})
		return
	}

	loading.value = true
	try {
		const res = await bindTenant(userId.value, inviteCode.value)

		// 保存登录信息
		userStore.setLoginInfo({
			token: res.token,
			user_info: userStore.userInfo,
			current_tenant: res.tenant_info,
			tenants: [res.tenant_info]
		})

		uni.showToast({
			title: '绑定成功',
			icon: 'success'
		})

	setTimeout(() => {
		uni.reLaunch({
			url: '/pages/index/index'
		})
	}, 1000)
	} catch (error) {
		console.error('绑定失败', error)
	} finally {
		loading.value = false
	}
}

// 联系管理员
const contactAdmin = () => {
	uni.showModal({
		title: '联系管理员',
		content: '请联系您的企业管理员获取邀请码',
		showCancel: false
	})
}
</script>

<style lang="scss" scoped>
.bind-tenant-container {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 40rpx;
}

.header {
	text-align: center;
	margin-bottom: 60rpx;

	.title {
		font-size: 48rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 16rpx;
	}

	.desc {
		font-size: 28rpx;
		color: #999;
	}
}

.form-card {
	background: #fff;
	border-radius: 16rpx;
	padding: 40rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

	.form-item {
		margin-bottom: 40rpx;

		.label {
			font-size: 28rpx;
			color: #666;
			margin-bottom: 16rpx;
		}

		.input {
			height: 88rpx;
			background: #f5f5f5;
			border-radius: 12rpx;
			padding: 0 24rpx;
			font-size: 32rpx;
			color: #333;
		}

		.placeholder {
			color: #ccc;
		}
	}

	.tips {
		padding: 24rpx;
		background: #f0f9ff;
		border-radius: 12rpx;
		margin-bottom: 40rpx;

		text {
			font-size: 26rpx;
			color: #1989fa;
		}
	}

	.submit-btn {
		width: 100%;
		height: 88rpx;
		background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
		border-radius: 12rpx;
		border: none;
		font-size: 32rpx;
		color: #fff;
	}
}

.help-text {
	text-align: center;
	margin-top: 40rpx;
	font-size: 28rpx;
	color: #999;

	.link {
		color: #1989fa;
		margin-left: 8rpx;
	}
}
</style>

