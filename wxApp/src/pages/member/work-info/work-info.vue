<template>
	<view class="work-info-container">
		<!-- 加载中 -->
		<view v-if="loading" class="loading-section">
			<text class="loading-text">加载中...</text>
		</view>

		<!-- 工作信息内容 -->
		<view v-else class="work-info-content">
			<!-- 提示信息 -->
			<view class="tip-card">
				<u-icon name="info-circle" size="36" color="#5EA3F2"></u-icon>
				<text class="tip-text">以下信息为只读，如有错误请联系人事部门修改</text>
			</view>

			<!-- 基本信息 -->
			<view class="section-card">
				<view class="section-title">基本信息</view>
				<view class="info-list">
					<view class="info-item">
						<text class="label">工号</text>
						<text class="value">{{ profile.job_number || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="label">部门</text>
						<text class="value">{{ profile.department || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="label">岗位</text>
						<text class="value">{{ profile.position || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="label">车间</text>
						<text class="value">{{ profile.workshop || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="label">班组</text>
						<text class="value">{{ profile.team || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="label">班组长</text>
						<text class="value">{{ profile.team_leader || '未设置' }}</text>
					</view>
				</view>
			</view>

			<!-- 工作时间 -->
			<view class="section-card">
				<view class="section-title">工作时间</view>
				<view class="info-list">
					<view class="info-item">
						<text class="label">入职日期</text>
						<text class="value highlight">{{ formatDate(profile.employed_at) }}</text>
					</view>
					<view class="info-item" v-if="profile.regular_at">
						<text class="label">转正日期</text>
						<text class="value">{{ formatDate(profile.regular_at) }}</text>
					</view>
					<view class="info-item">
						<text class="label">工龄</text>
						<text class="value highlight">{{ profile.work_years }}年{{ profile.work_months }}个月</text>
					</view>
				</view>
			</view>

			<!-- 合同信息 -->
			<view class="section-card">
				<view class="section-title">合同信息</view>
				<view class="info-list">
					<view class="info-item">
						<text class="label">合同类型</text>
						<text class="value">{{ getContractTypeText(profile.contract_type) }}</text>
					</view>
					<view class="info-item" v-if="profile.contract_start_at">
						<text class="label">合同开始</text>
						<text class="value">{{ formatDate(profile.contract_start_at) }}</text>
					</view>
					<view class="info-item" v-if="profile.contract_end_at">
						<text class="label">合同结束</text>
						<text class="value">{{ formatDate(profile.contract_end_at) }}</text>
					</view>
				</view>
			</view>

			<!-- 状态 -->
			<view class="section-card">
				<view class="section-title">状态</view>
				<view class="info-list">
					<view class="info-item">
						<text class="label">当前状态</text>
						<view class="status-badge" :class="getStatusClass(profile.status)">
							{{ getStatusText(profile.status) }}
						</view>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProfile } from '@/api/member'

const loading = ref(true)
const profile = ref({})

onMounted(async () => {
	await loadProfile()
})

const loadProfile = async () => {
	loading.value = true
	try {
		const res = await getProfile()
		if (res.code === 0) {
			profile.value = res.data || {}
		} else {
			uni.showToast({
				title: res.msg || '加载失败',
				icon: 'none'
			})
		}
	} catch (error) {
		console.error('加载工作信息失败', error)
		uni.showToast({
			title: '加载失败',
			icon: 'none'
		})
	} finally {
		loading.value = false
	}
}

const formatDate = (timestamp) => {
	if (!timestamp || timestamp === 0) return '未设置'
	const date = new Date(timestamp * 1000)
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const getContractTypeText = (type) => {
	const map = {
		fulltime: '全职',
		parttime: '兼职',
		intern: '实习',
		dispatch: '劳务派遣'
	}
	return map[type] || '未设置'
}

const getStatusText = (status) => {
	const map = {
		active: '在职',
		probation: '试用期',
		inactive: '离职',
		suspended: '停职'
	}
	return map[status] || status
}

const getStatusClass = (status) => {
	const map = {
		active: 'active',
		probation: 'probation',
		inactive: 'inactive',
		suspended: 'suspended'
	}
	return map[status] || ''
}
</script>

<style scoped lang="scss">
.work-info-container {
	min-height: 100vh;
	background: #f5f7fa;
	padding-bottom: 40rpx;
}

.loading-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	height: 80vh;
	gap: 20rpx;
}

.loading-text {
	font-size: 28rpx;
	color: #999;
}

.work-info-content {
	padding: 0 24rpx;
}

.tip-card {
	display: flex;
	align-items: center;
	gap: 16rpx;
	background: #e8f4ff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin: 24rpx 0;

	.tip-text {
		flex: 1;
		font-size: 26rpx;
		color: #5EA3F2;
		line-height: 1.6;
	}
}

.section-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
}

.section-title {
	font-size: 32rpx;
	color: #333;
	font-weight: bold;
	margin-bottom: 24rpx;
}

.info-list {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.info-item {
	display: flex;
	align-items: center;
	justify-content: space-between;

	.label {
		font-size: 28rpx;
		color: #666;
	}

	.value {
		font-size: 28rpx;
		color: #333;
		font-weight: 500;

		&.highlight {
			color: #5EA3F2;
			font-weight: bold;
		}
	}
}

.status-badge {
	padding: 8rpx 24rpx;
	border-radius: 32rpx;
	font-size: 26rpx;
	font-weight: 500;

	&.active {
		background: #e8f8f5;
		color: #00b578;
	}

	&.probation {
		background: #fff3e0;
		color: #ff9800;
	}

	&.inactive {
		background: #f5f5f5;
		color: #999;
	}

	&.suspended {
		background: #ffebee;
		color: #f44336;
	}
}
</style>

