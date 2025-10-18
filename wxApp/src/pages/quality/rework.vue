<template>
	<view class="rework-page">
		<!-- 统计卡片 -->
		<view class="statistics">
			<view class="stat-item">
				<text class="stat-value">{{ statistics.pending }}</text>
				<text class="stat-label">待返工</text>
			</view>
			<view class="stat-divider"></view>
			<view class="stat-item">
				<text class="stat-value">{{ statistics.in_progress }}</text>
				<text class="stat-label">返工中</text>
			</view>
			<view class="stat-divider"></view>
			<view class="stat-item">
				<text class="stat-value">{{ statistics.completed }}</text>
				<text class="stat-label">已完成</text>
			</view>
		</view>

		<!-- 返工列表 -->
		<view class="rework-list">
			<view v-for="item in reworks" :key="item.id" class="rework-item" @click="handleDetail(item)">
				<view class="item-header">
					<text class="contract-no">{{ item.contract_no }}</text>
					<view class="status-tag" :class="'status-' + item.status">{{ item.status_text }}</view>
				</view>
				<view class="item-body">
					<view class="info-row">
						<text class="label">款名：</text>
						<text class="value">{{ item.style_name }}</text>
					</view>
					<view class="info-row">
						<text class="label">扎号：</text>
						<text class="value">{{ item.bundle_no || '-' }}</text>
					</view>
					<view class="info-row">
						<text class="label">颜色/尺码：</text>
						<text class="value">{{ item.color }} / {{ item.size || '-' }}</text>
					</view>
					<view class="info-row">
						<text class="label">返工工序：</text>
						<text class="value">{{ item.source_procedure_name }} → {{ item.target_procedure_name }}</text>
					</view>
					<view class="info-row">
						<text class="label">数量：</text>
						<text class="value highlight">{{ item.rework_qty }} 件</text>
					</view>
					<view class="info-row">
						<text class="label">原因：</text>
						<text class="value">{{ item.rework_reason }}</text>
					</view>
				</view>
				<view class="item-footer">
					<text class="time">{{ formatTime(item.created_at) }}</text>
					<button v-if="item.status === 0" size="mini" type="primary" @click.stop="handleAccept(item)">
						接受任务
					</button>
					<button v-if="item.status === 1" size="mini" type="primary" @click.stop="handleComplete(item)">
						完成返工
					</button>
				</view>
			</view>

			<!-- 加载更多 -->
			<view class="load-more" v-if="hasMore && !loading">
				<text @click="loadMore">加载更多</text>
			</view>
			<view class="load-more" v-if="loading">
				<text>加载中...</text>
			</view>
			<view class="load-more" v-if="!hasMore && reworks.length > 0">
				<text class="no-more">没有更多了</text>
			</view>

			<!-- 空状态 -->
			<view class="empty-state" v-if="!loading && reworks.length === 0">
				<u-icon name="file-text" size="100" color="#ccc" />
				<text>暂无返工记录</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getReworkList, completeRework } from '@/api/quality'

const reworks = ref([])
const statistics = ref({
	total: 0,
	pending: 0,
	in_progress: 0,
	completed: 0
})
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const hasMore = ref(true)

onMounted(() => {
	loadReworks()
})

// 加载返工列表
const loadReworks = async () => {
	if (loading.value) return

	try {
		loading.value = true
		const res = await getReworkList({
			page: page.value,
			page_size: pageSize.value
		})

		if (page.value === 1) {
			reworks.value = res.data.reworks
		} else {
			reworks.value.push(...res.data.reworks)
		}

		statistics.value = res.data.statistics
		hasMore.value = reworks.value.length < res.data.total
		loading.value = false
	} catch (error) {
		loading.value = false
		console.error('加载返工列表失败:', error)
	}
}

// 加载更多
const loadMore = () => {
	page.value++
	loadReworks()
}

// 查看详情
const handleDetail = (item) => {
	// TODO: 跳转到详情页
	console.log('查看详情:', item)
}

// 接受任务
const handleAccept = (item) => {
	uni.showModal({
		title: '接受返工任务',
		content: `确定接受该返工任务吗？`,
		success: async (res) => {
			if (res.confirm) {
				try {
					// TODO: 调用接受任务接口
					uni.showToast({ title: '已接受任务', icon: 'success' })
					setTimeout(() => {
						loadReworks()
					}, 1500)
				} catch (error) {
					uni.showToast({ title: '操作失败', icon: 'none' })
				}
			}
		}
	})
}

// 完成返工
const handleComplete = (item) => {
	uni.showModal({
		title: '完成返工',
		content: `确认已完成该返工任务吗？`,
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '提交中...' })
					await completeRework(item.id, {
						remark: '返工完成'
					})
					uni.hideLoading()
					uni.showToast({ title: '返工完成', icon: 'success' })
					setTimeout(() => {
						page.value = 1
						loadReworks()
					}, 1500)
				} catch (error) {
					uni.hideLoading()
					uni.showToast({ title: '操作失败', icon: 'none' })
				}
			}
		}
	})
}

// 格式化时间
const formatTime = (timestamp) => {
	const date = new Date(timestamp * 1000)
	return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
}
</script>

<style lang="scss" scoped>
.rework-page {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 32rpx;
}

.statistics {
	background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
	padding: 48rpx 24rpx;
	display: flex;
	justify-content: space-around;
	align-items: center;

	.stat-item {
		text-align: center;

		.stat-value {
			display: block;
			font-size: 48rpx;
			font-weight: bold;
			color: #fff;
			margin-bottom: 8rpx;
		}

		.stat-label {
			font-size: 24rpx;
			color: rgba(255, 255, 255, 0.8);
		}
	}

	.stat-divider {
		width: 1rpx;
		height: 80rpx;
		background: rgba(255, 255, 255, 0.3);
	}
}

.rework-list {
	padding: 24rpx;

	.rework-item {
		background: #fff;
		border-radius: 16rpx;
		padding: 32rpx;
		margin-bottom: 24rpx;
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

		.item-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 24rpx;

			.contract-no {
				font-size: 32rpx;
				font-weight: bold;
				color: #333;
			}

			.status-tag {
				padding: 8rpx 16rpx;
				border-radius: 8rpx;
				font-size: 24rpx;

				&.status-0 {
					background: #FFF3E0;
					color: #FF9800;
				}

				&.status-1 {
					background: #E3F2FD;
					color: #2196F3;
				}

				&.status-2 {
					background: #E8F5E9;
					color: #4CAF50;
				}
			}
		}

		.item-body {
			.info-row {
				display: flex;
				margin-bottom: 16rpx;
				font-size: 26rpx;

				.label {
					color: #999;
					min-width: 120rpx;
				}

				.value {
					flex: 1;
					color: #333;

					&.highlight {
						color: #5EA3F2;
						font-weight: 500;
					}
				}
			}
		}

		.item-footer {
			margin-top: 24rpx;
			padding-top: 24rpx;
			border-top: 1rpx solid #f0f0f0;
			display: flex;
			justify-content: space-between;
			align-items: center;

			.time {
				font-size: 24rpx;
				color: #999;
			}
		}
	}

	.load-more {
		text-align: center;
		padding: 24rpx 0;
		font-size: 26rpx;
		color: #999;

		.no-more {
			color: #ccc;
		}
	}

	.empty-state {
		padding: 200rpx 0;
		text-align: center;

		text {
			display: block;
			margin-top: 32rpx;
			font-size: 28rpx;
			color: #999;
		}
	}
}
</style>

