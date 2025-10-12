<template>
	<view class="select-tenant-container">
		<view class="header">
			<view class="title">选择企业</view>
			<view class="desc">您在以下企业中有账号，请选择</view>
		</view>

		<view class="tenant-list">
			<view 
				class="tenant-item"
				v-for="tenant in tenants"
				:key="tenant.tenant_id"
				@click="handleSelect(tenant)"
			>
				<view class="tenant-info">
					<view class="tenant-name">{{ tenant.tenant_name }}</view>
					<view class="tenant-code">企业代码：{{ tenant.tenant_code }}</view>
					<view class="tenant-status" :class="tenant.status">
						{{ tenant.status === 'active' ? '✅ 在职' : '❌ 已离职' }}
					</view>
				</view>
				<view class="arrow">›</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { selectTenant } from '@/api/auth'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const userId = ref('')
const tenants = ref([])

// 获取页面参数
onMounted(() => {
	const pages = getCurrentPages()
	const currentPage = pages[pages.length - 1]
	const options = currentPage.options || {}
	
	userId.value = options.userId || ''
	if (options.tenants) {
		try {
			tenants.value = JSON.parse(decodeURIComponent(options.tenants))
		} catch (e) {
			console.error('解析租户列表失败', e)
		}
	}
})

// 选择租户
const handleSelect = async (tenant) => {
	uni.showLoading({
		title: '登录中...'
	})

	try {
		const res = await selectTenant(userId.value, tenant.tenant_id)

		// 保存登录信息
		userStore.setLoginInfo({
			token: res.token,
			user_info: res.user_info,
			current_tenant: res.current_tenant,
			tenants: tenants.value
		})

		uni.hideLoading()
		uni.showToast({
			title: '登录成功',
			icon: 'success'
		})

	setTimeout(() => {
		uni.reLaunch({
			url: '/pages/index/index'
		})
	}, 1000)
	} catch (error) {
		uni.hideLoading()
		console.error('选择租户失败', error)
	}
}
</script>

<style lang="scss" scoped>
.select-tenant-container {
	min-height: 100vh;
	background: #f5f5f5;
	padding: 40rpx;
}

.header {
	text-align: center;
	margin-bottom: 40rpx;

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

.tenant-list {
	.tenant-item {
		background: #fff;
		border-radius: 16rpx;
		padding: 32rpx;
		margin-bottom: 24rpx;
		display: flex;
		align-items: center;
		justify-content: space-between;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);
		transition: all 0.3s;

		&:active {
			transform: scale(0.98);
			opacity: 0.9;
		}

		.tenant-info {
			flex: 1;

			.tenant-name {
				font-size: 32rpx;
				font-weight: bold;
				color: #333;
				margin-bottom: 12rpx;
			}

			.tenant-code {
				font-size: 26rpx;
				color: #999;
				margin-bottom: 12rpx;
			}

			.tenant-status {
				font-size: 24rpx;
				padding: 4rpx 12rpx;
				border-radius: 8rpx;
				display: inline-block;

				&.active {
					background: #e7f9ef;
					color: #52c41a;
				}

				&.inactive {
					background: #fff1f0;
					color: #ff4d4f;
				}
			}
		}

		.arrow {
			font-size: 48rpx;
			color: #ccc;
			margin-left: 16rpx;
		}
	}
}
</style>

