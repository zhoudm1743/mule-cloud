<template>
	<view class="index-container">
		<!-- 顶部欢迎区域 -->
		<view class="welcome-section">
			<view class="welcome-text">
				<text class="greeting">你好，{{ userInfo.nickname || '用户' }}</text>
				<text class="date">{{ currentDate }}</text>
			</view>
			<view class="tenant-badge" v-if="currentTenant">
				<text>{{ currentTenant.tenant_name }}</text>
			</view>
		</view>

		<!-- 快捷功能 -->
		<view class="quick-actions">
			<view class="action-item" @click="handleAction('order')">
				<view class="action-icon">📋</view>
				<view class="action-text">订单管理</view>
			</view>
			<view class="action-item" @click="handleAction('production')">
				<view class="action-icon">🏭</view>
				<view class="action-text">生产进度</view>
			</view>
			<view class="action-item" @click="handleAction('quality')">
				<view class="action-icon">✅</view>
				<view class="action-text">质量检查</view>
			</view>
			<view class="action-item" @click="handleAction('report')">
				<view class="action-icon">📊</view>
				<view class="action-text">数据报表</view>
			</view>
		</view>

		<!-- 统计数据 -->
		<view class="stats-section">
			<view class="section-title">今日数据</view>
			<view class="stats-grid">
				<view class="stat-item">
					<view class="stat-value">{{ stats.todayOrders }}</view>
					<view class="stat-label">今日订单</view>
				</view>
				<view class="stat-item">
					<view class="stat-value">{{ stats.inProduction }}</view>
					<view class="stat-label">生产中</view>
				</view>
				<view class="stat-item">
					<view class="stat-value">{{ stats.completed }}</view>
					<view class="stat-label">已完成</view>
				</view>
				<view class="stat-item">
					<view class="stat-value">{{ stats.quality }}%</view>
					<view class="stat-label">合格率</view>
				</view>
			</view>
		</view>

		<!-- 最近订单 -->
		<view class="recent-orders">
			<view class="section-title">最近订单</view>
			<view class="order-list">
				<view 
					class="order-item"
					v-for="order in recentOrders"
					:key="order.id"
					@click="handleOrderDetail(order)"
				>
					<view class="order-info">
						<view class="order-name">{{ order.name }}</view>
						<view class="order-time">{{ order.time }}</view>
					</view>
					<view class="order-status" :class="order.status">
						{{ order.statusText }}
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()

const userInfo = computed(() => userStore.userInfo || {})
const currentTenant = computed(() => userStore.currentTenant)

// 当前日期
const currentDate = ref('')

// 统计数据（mock数据）
const stats = ref({
	todayOrders: 28,
	inProduction: 15,
	completed: 12,
	quality: 98.5
})

// 最近订单（mock数据）
const recentOrders = ref([
	{ id: 1, name: '春季新款连衣裙', time: '2小时前', status: 'processing', statusText: '生产中' },
	{ id: 2, name: '夏季T恤套装', time: '5小时前', status: 'completed', statusText: '已完成' },
	{ id: 3, name: '儿童校服定制', time: '1天前', status: 'pending', statusText: '待开始' }
])

onMounted(() => {
	// 检查登录状态
	if (!userStore.isLoggedIn) {
		uni.reLaunch({
			url: '/pages/login/login'
		})
		return
	}

	// 设置当前日期
	const now = new Date()
	const options = { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' }
	currentDate.value = now.toLocaleDateString('zh-CN', options)

	// 加载数据
	loadData()
})

// 加载数据
const loadData = async () => {
	// TODO: 调用实际的API加载数据
	console.log('加载数据...')
}

// 处理快捷操作
const handleAction = (type) => {
	uni.showToast({
		title: '功能开发中...',
		icon: 'none'
	})
}

// 查看订单详情
const handleOrderDetail = (order) => {
	uni.showToast({
		title: '订单详情开发中...',
		icon: 'none'
	})
}
</script>

<style lang="scss" scoped>
.index-container {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 40rpx;
}

.welcome-section {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	padding: 60rpx 40rpx;
	display: flex;
	justify-content: space-between;
	align-items: center;

	.welcome-text {
		.greeting {
			display: block;
			font-size: 40rpx;
			font-weight: bold;
			color: #fff;
			margin-bottom: 12rpx;
		}

		.date {
			display: block;
			font-size: 26rpx;
			color: rgba(255, 255, 255, 0.8);
		}
	}

	.tenant-badge {
		background: rgba(255, 255, 255, 0.2);
		padding: 12rpx 24rpx;
		border-radius: 20rpx;
		backdrop-filter: blur(10rpx);

		text {
			font-size: 24rpx;
			color: #fff;
		}
	}
}

.quick-actions {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 24rpx;
	padding: 32rpx 24rpx;
	background: #fff;
	margin: 24rpx;
	border-radius: 16rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

	.action-item {
		text-align: center;

		.action-icon {
			font-size: 60rpx;
			margin-bottom: 12rpx;
		}

		.action-text {
			font-size: 24rpx;
			color: #666;
		}
	}
}

.stats-section, .recent-orders {
	background: #fff;
	margin: 24rpx;
	padding: 32rpx;
	border-radius: 16rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.05);

	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 24rpx;
	}
}

.stats-grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 24rpx;

	.stat-item {
		text-align: center;

		.stat-value {
			font-size: 40rpx;
			font-weight: bold;
			color: #667eea;
			margin-bottom: 8rpx;
		}

		.stat-label {
			font-size: 24rpx;
			color: #999;
		}
	}
}

.order-list {
	.order-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 24rpx 0;
		border-bottom: 1rpx solid #f5f5f5;

		&:last-child {
			border-bottom: none;
		}

		.order-info {
			flex: 1;

			.order-name {
				font-size: 28rpx;
				color: #333;
				margin-bottom: 8rpx;
			}

			.order-time {
				font-size: 24rpx;
				color: #999;
			}
		}

		.order-status {
			font-size: 24rpx;
			padding: 8rpx 16rpx;
			border-radius: 8rpx;

			&.processing {
				background: #e6f7ff;
				color: #1890ff;
			}

			&.completed {
				background: #f6ffed;
				color: #52c41a;
			}

			&.pending {
				background: #fff7e6;
				color: #fa8c16;
			}
		}
	}
}
</style>
