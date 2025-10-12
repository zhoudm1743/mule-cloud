<template>
	<view class="scan-container">
		<!-- 扫码动画区域 -->
		<view class="scan-area">
			<view class="scan-frame">
				<view class="corner corner-tl"></view>
				<view class="corner corner-tr"></view>
				<view class="corner corner-bl"></view>
				<view class="corner corner-br"></view>
				<view class="scan-line"></view>
			</view>
			
			<view class="scan-tips">
				<text>将二维码放入框内，即可自动扫描</text>
			</view>
		</view>

		<!-- 扫码按钮 -->
		<view class="scan-actions">
			<u-button 
				type="primary" 
				:custom-style="{
					background: 'linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%)',
					border: 'none',
					borderRadius: '24rpx',
					height: '96rpx',
					fontSize: '32rpx',
					fontWeight: 'bold'
				}"
				@click="handleScan"
			>
				<view class="scan-btn-content">
					<u-icon name="scan" :size="40" color="#fff"></u-icon>
					<text class="scan-text">开始扫码</text>
				</view>
			</u-button>

			<u-button 
				type="info" 
				plain
				:custom-style="{
					borderColor: '#5EA3F2',
					color: '#5EA3F2',
					borderRadius: '24rpx',
					height: '96rpx',
					marginTop: '24rpx',
					fontSize: '32rpx'
				}"
				@click="handleManualInput"
			>
				手动输入
			</u-button>
		</view>

		<!-- 快速操作 -->
		<view class="quick-scan">
			<view class="quick-title">快速操作</view>
			<view class="quick-grid">
			<view class="quick-item" @click="handleQuickAction('order')">
				<view class="quick-icon" style="background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);">
					<u-icon name="order" :size="36" color="#fff"></u-icon>
				</view>
				<view class="quick-text">订单扫码</view>
			</view>
			<view class="quick-item" @click="handleQuickAction('process')">
				<view class="quick-icon" style="background: linear-gradient(135deg, #66BB6A 0%, #4CAF50 100%);">
					<u-icon name="clock" :size="36" color="#fff"></u-icon>
				</view>
				<view class="quick-text">工序上报</view>
			</view>
			<view class="quick-item" @click="handleQuickAction('quality')">
				<view class="quick-icon" style="background: linear-gradient(135deg, #FFA726 0%, #FF9800 100%);">
					<u-icon name="checkmark" :size="36" color="#fff"></u-icon>
				</view>
				<view class="quick-text">质检扫码</view>
			</view>
			<view class="quick-item" @click="handleQuickAction('finish')">
				<view class="quick-icon" style="background: linear-gradient(135deg, #AB47BC 0%, #9C27B0 100%);">
					<u-icon name="checkmark-circle-fill" :size="36" color="#fff"></u-icon>
				</view>
				<view class="quick-text">完工扫码</view>
			</view>
			</view>
		</view>

		<!-- 最近扫码记录 -->
		<view class="recent-scans">
			<view class="recent-header">
				<view class="recent-title">
					<u-icon name="file-text" :size="28" color="#333" style="margin-right: 8rpx;"></u-icon>
					<text>最近扫码</text>
				</view>
				<view class="recent-more" @click="handleViewAll">
					<text>查看全部</text>
					<u-icon name="arrow-right" :size="24" color="#5EA3F2" style="margin-left: 4rpx;"></u-icon>
				</view>
			</view>
			<view class="scan-list">
				<view 
					class="scan-record" 
					v-for="record in recentRecords" 
					:key="record.id"
					@click="handleRecordDetail(record)"
				>
			<view class="record-icon">
				<u-icon name="scan" :size="32" color="#5EA3F2"></u-icon>
			</view>
					<view class="record-info">
						<view class="record-title">{{ record.title }}</view>
						<view class="record-time">{{ record.time }}</view>
					</view>
					<view class="record-status" :class="record.statusClass">
						{{ record.status }}
					</view>
				</view>
			</view>
		</view>
		
		<!-- 底部导航栏 -->
		<TabBar :current="1" />
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import TabBar from '@/components/TabBar/TabBar.vue'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()

// 最近扫码记录（mock数据）
const recentRecords = ref([
	{ id: 1, title: '订单 #2024001', time: '2分钟前', status: '已处理', statusClass: 'success' },
	{ id: 2, title: '工序上报 #工001', time: '10分钟前', status: '已上报', statusClass: 'success' },
	{ id: 3, title: '质检 #QC2024', time: '1小时前', status: '合格', statusClass: 'success' }
])

onMounted(() => {
	// 页面加载时直接调用扫码
	handleScan()
})

// 扫码
const handleScan = () => {
	uni.scanCode({
		success: (res) => {
			console.log('扫码结果:', res)
			
			uni.showLoading({ title: '识别中...' })
			
			// TODO: 根据扫码结果调用后端API处理
			setTimeout(() => {
				uni.hideLoading()
				
				uni.showModal({
					title: '扫码成功',
					content: '扫描内容：' + res.result,
					confirmText: '继续扫码',
					cancelText: '返回首页',
					success: (modalRes) => {
						if (modalRes.confirm) {
							// 继续扫码
							handleScan()
						} else {
							// 返回首页
							uni.reLaunch({
								url: '/pages/index/index'
							})
						}
					}
				})
			}, 1000)
		},
		fail: (err) => {
			console.error('扫码失败:', err)
			uni.showToast({
				title: '扫码失败',
				icon: 'none'
			})
		}
	})
}

// 手动输入
const handleManualInput = () => {
	uni.showModal({
		title: '手动输入',
		editable: true,
		placeholderText: '请输入编号',
		success: (res) => {
			if (res.confirm && res.content) {
				console.log('手动输入:', res.content)
				// TODO: 处理手动输入的内容
			}
		}
	})
}

// 快速操作
const handleQuickAction = (type) => {
	console.log('快速操作:', type)
	// 根据类型直接调用扫码，后端根据类型处理
	handleScan()
}

// 查看全部记录
const handleViewAll = () => {
	uni.showToast({
		title: '功能开发中',
		icon: 'none'
	})
}

// 查看记录详情
const handleRecordDetail = (record) => {
	console.log('查看记录:', record)
	uni.showToast({
		title: '功能开发中',
		icon: 'none'
	})
}
</script>

<style lang="scss" scoped>
.scan-container {
	min-height: 100vh;
	background: linear-gradient(180deg, #F0F8FF 0%, #FFFFFF 40%);
	padding-bottom: 120rpx;
}

.scan-area {
	padding: 80rpx 0;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.scan-frame {
	width: 500rpx;
	height: 500rpx;
	position: relative;
	background: rgba(0, 0, 0, 0.5);
	border-radius: 20rpx;
	overflow: hidden;
	
	.corner {
		position: absolute;
		width: 60rpx;
		height: 60rpx;
		border: 6rpx solid #5EA3F2;
		
		&.corner-tl {
			top: 0;
			left: 0;
			border-right: none;
			border-bottom: none;
			border-top-left-radius: 20rpx;
		}
		
		&.corner-tr {
			top: 0;
			right: 0;
			border-left: none;
			border-bottom: none;
			border-top-right-radius: 20rpx;
		}
		
		&.corner-bl {
			bottom: 0;
			left: 0;
			border-right: none;
			border-top: none;
			border-bottom-left-radius: 20rpx;
		}
		
		&.corner-br {
			bottom: 0;
			right: 0;
			border-left: none;
			border-top: none;
			border-bottom-right-radius: 20rpx;
		}
	}
	
	.scan-line {
		position: absolute;
		top: 0;
		left: 0;
		right: 0;
		height: 4rpx;
		background: linear-gradient(90deg, transparent 0%, #5EA3F2 50%, transparent 100%);
		animation: scan 2s linear infinite;
	}
}

@keyframes scan {
	0% { transform: translateY(0); }
	100% { transform: translateY(500rpx); }
}

.scan-tips {
	margin-top: 40rpx;
	font-size: 28rpx;
	color: #666;
	text-align: center;
}

.scan-actions {
	padding: 0 40rpx;
	margin-top: 40rpx;
	
	.scan-btn-content {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 16rpx;
		
		.scan-text {
			font-size: 32rpx;
			font-weight: bold;
		}
	}
}

.quick-scan {
	margin: 60rpx 24rpx 40rpx;
	background: #fff;
	border-radius: 24rpx;
	padding: 40rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.05);
	
	.quick-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 32rpx;
	}
	
	.quick-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 32rpx;
	}
	
	.quick-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		
		.quick-icon {
			width: 96rpx;
			height: 96rpx;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			box-shadow: 0 4rpx 12rpx rgba(94, 163, 242, 0.3);
			margin-bottom: 16rpx;
		}
		
		.quick-text {
			font-size: 24rpx;
			color: #666;
			text-align: center;
		}
	}
}

.recent-scans {
	margin: 0 24rpx 40rpx;
	background: #fff;
	border-radius: 24rpx;
	padding: 40rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.05);
	
	.recent-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 32rpx;
		
		.recent-title {
			display: flex;
			align-items: center;
			font-size: 32rpx;
			font-weight: bold;
			color: #333;
		}
		
		.recent-more {
			display: flex;
			align-items: center;
			font-size: 26rpx;
			color: #5EA3F2;
		}
	}
	
	.scan-list {
		.scan-record {
			display: flex;
			align-items: center;
			padding: 24rpx 0;
			border-bottom: 1rpx solid #f5f5f5;
			
			&:last-child {
				border-bottom: none;
			}
			
			.record-icon {
				width: 80rpx;
				height: 80rpx;
				border-radius: 50%;
				background: #F0F8FF;
				display: flex;
				align-items: center;
				justify-content: center;
				margin-right: 24rpx;
			}
			
			.record-info {
				flex: 1;
				
				.record-title {
					font-size: 28rpx;
					color: #333;
					margin-bottom: 8rpx;
				}
				
				.record-time {
					font-size: 24rpx;
					color: #999;
				}
			}
			
			.record-status {
				font-size: 24rpx;
				padding: 8rpx 16rpx;
				border-radius: 8rpx;
				
				&.success {
					color: #4CAF50;
					background: rgba(76, 175, 80, 0.1);
				}
			}
		}
	}
}
</style>
