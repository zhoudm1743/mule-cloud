<template>
	<view class="report-history">
		<!-- 统计卡片 -->
		<view class="statistics">
			<view class="stat-item">
				<text class="stat-value">{{ statistics.total_quantity }}</text>
				<text class="stat-label">总件数</text>
			</view>
			<view class="stat-divider"></view>
			<view class="stat-item">
				<text class="stat-value highlight">¥{{ statistics.total_amount }}</text>
				<text class="stat-label">总工资</text>
			</view>
		</view>
		
		<!-- 筛选栏 -->
		<view class="filter-bar">
			<picker mode="date" :value="filterDate" @change="onDateChange">
				<view class="filter-item">
					<u-icon name="calendar" size="20" color="#5EA3F2" />
					<text>{{ filterDate || '选择日期' }}</text>
				</view>
			</picker>
			<view class="filter-item" @click="clearFilter" v-if="filterDate">
				<text>清除</text>
			</view>
		</view>
		
		<!-- 记录列表 -->
		<view class="record-list">
			<view 
				v-for="record in records" 
				:key="record.id"
				class="record-item"
				@click="handleRecordDetail(record)"
			>
				<view class="record-header">
					<text class="contract-no">{{ record.contract_no }}</text>
					<text class="amount">¥{{ record.total_price }}</text>
				</view>
				<view class="record-body">
					<view class="record-info">
						<text>扎号：{{ record.bundle_no || '-' }}</text>
					</view>
					<view class="record-info">
						<text>{{ record.color }} / {{ record.size || '-' }}</text>
					</view>
					<view class="record-info">
						<text>{{ record.procedure_name }}</text>
						<text>{{ record.quantity }}件 × {{ record.unit_price }}元</text>
					</view>
				</view>
				<view class="record-footer">
					<text class="time">{{ formatTime(record.report_time) }}</text>
				</view>
			</view>
			
			<!-- 加载更多 -->
			<view class="load-more" v-if="hasMore && !loading">
				<text @click="loadMore">加载更多</text>
			</view>
			<view class="load-more" v-if="loading">
				<text>加载中...</text>
			</view>
			<view class="load-more" v-if="!hasMore && records.length > 0">
				<text class="no-more">没有更多了</text>
			</view>
			
			<!-- 空状态 -->
			<view class="empty-state" v-if="!loading && records.length === 0">
				<u-icon name="file-text" size="100" color="#ccc" />
				<text>暂无上报记录</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getReportList } from '@/api/scan'

const statistics = ref({
	total_quantity: 0,
	total_amount: 0
})
const filterDate = ref('')
const records = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = 20
const hasMore = ref(true)

onMounted(() => {
	loadRecords()
})

// 加载记录
const loadRecords = async (isLoadMore = false) => {
	if (loading.value) return
	
	try {
		loading.value = true
		
		const params = {
			page: page.value,
			page_size: pageSize
		}
		
		if (filterDate.value) {
			params.start_date = filterDate.value
			params.end_date = filterDate.value
		}
		
		const res = await getReportList(params)
		
		if (isLoadMore) {
			records.value = [...records.value, ...res.data.reports]
		} else {
			records.value = res.data.reports
		}
		
		// 更新统计
		if (res.data.statistics) {
			statistics.value = res.data.statistics
		}
		
		// 判断是否还有更多
		hasMore.value = records.value.length < res.data.total
		
		loading.value = false
	} catch (error) {
		loading.value = false
		console.error('加载记录失败:', error)
		uni.showToast({
			title: error.message || '加载失败',
			icon: 'none'
		})
	}
}

// 日期筛选
const onDateChange = (e) => {
	filterDate.value = e.detail.value
	page.value = 1
	loadRecords()
}

// 清除筛选
const clearFilter = () => {
	filterDate.value = ''
	page.value = 1
	loadRecords()
}

// 加载更多
const loadMore = () => {
	page.value++
	loadRecords(true)
}

// 格式化时间
const formatTime = (timestamp) => {
	const date = new Date(timestamp * 1000)
	const now = new Date()
	const diff = now.getTime() - date.getTime()
	
	// 1分钟内
	if (diff < 60 * 1000) {
		return '刚刚'
	}
	
	// 1小时内
	if (diff < 60 * 60 * 1000) {
		return `${Math.floor(diff / (60 * 1000))}分钟前`
	}
	
	// 今天
	if (date.toDateString() === now.toDateString()) {
		return `今天 ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
	}
	
	// 昨天
	const yesterday = new Date(now)
	yesterday.setDate(yesterday.getDate() - 1)
	if (date.toDateString() === yesterday.toDateString()) {
		return `昨天 ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
	}
	
	// 其他
	return `${date.getMonth() + 1}月${date.getDate()}日 ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

// 查看详情
const handleRecordDetail = (record) => {
	console.log('查看记录:', record)
	// TODO: 跳转到详情页或显示详情弹窗
	uni.showToast({
		title: '功能开发中',
		icon: 'none'
	})
}
</script>

<style lang="scss" scoped>
.report-history {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 32rpx;
}

.statistics {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	margin: 24rpx;
	border-radius: 16rpx;
	padding: 48rpx 32rpx;
	display: flex;
	align-items: center;
	justify-content: space-around;
	box-shadow: 0 8rpx 24rpx rgba(102, 126, 234, 0.3);
	
	.stat-item {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		
		.stat-value {
			font-size: 48rpx;
			font-weight: bold;
			color: #fff;
			margin-bottom: 8rpx;
			
			&.highlight {
				color: #ffd700;
			}
		}
		
		.stat-label {
			font-size: 24rpx;
			color: rgba(255, 255, 255, 0.8);
		}
	}
	
	.stat-divider {
		width: 2rpx;
		height: 60rpx;
		background: rgba(255, 255, 255, 0.3);
	}
}

.filter-bar {
	margin: 24rpx;
	display: flex;
	gap: 16rpx;
	
	.filter-item {
		flex: 1;
		background: #fff;
		border-radius: 12rpx;
		padding: 20rpx 24rpx;
		display: flex;
		align-items: center;
		gap: 12rpx;
		font-size: 28rpx;
		color: #333;
		
		&:last-child {
			flex: 0 0 auto;
			color: #5EA3F2;
		}
	}
}

.record-list {
	margin: 0 24rpx;
	
	.record-item {
		background: #fff;
		border-radius: 12rpx;
		padding: 24rpx;
		margin-bottom: 16rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
		
		.record-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 16rpx;
			
			.contract-no {
				font-size: 30rpx;
				font-weight: bold;
				color: #333;
			}
			
			.amount {
				font-size: 32rpx;
				font-weight: bold;
				color: #ff6b6b;
			}
		}
		
		.record-body {
			.record-info {
				display: flex;
				justify-content: space-between;
				font-size: 26rpx;
				color: #666;
				margin-bottom: 8rpx;
				
				&:last-child {
					margin-bottom: 0;
				}
			}
		}
		
		.record-footer {
			margin-top: 12rpx;
			padding-top: 12rpx;
			border-top: 1rpx solid #f0f0f0;
			
			.time {
				font-size: 24rpx;
				color: #999;
			}
		}
	}
	
	.load-more {
		text-align: center;
		padding: 32rpx 0;
		font-size: 28rpx;
		color: #999;
		
		.no-more {
			color: #ccc;
		}
	}
	
	.empty-state {
		text-align: center;
		padding: 120rpx 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 24rpx;
		
		text {
			font-size: 28rpx;
			color: #999;
		}
	}
}
</style>

