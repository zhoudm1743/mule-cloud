<template>
	<view class="report-detail">
		<!-- 批次信息卡片 -->
		<view class="batch-info">
			<view class="info-header">
				<image v-if="batchData.order.style_image" :src="batchData.order.style_image" class="style-image" mode="aspectFill" />
				<view class="info-text">
					<view class="style-name">{{ batchData.order.style_name }}</view>
					<view class="contract-no">合同号：{{ batchData.batch.contract_no }}</view>
				</view>
			</view>
			<view class="info-grid">
				<view class="info-item">
					<text class="label">扎号</text>
					<text class="value">{{ batchData.batch.bundle_no || '-' }}</text>
				</view>
				<view class="info-item">
					<text class="label">颜色</text>
					<text class="value">{{ batchData.batch.color }}</text>
				</view>
				<view class="info-item">
					<text class="label">尺码</text>
					<text class="value">{{ batchData.batch.size || '-' }}</text>
				</view>
				<view class="info-item">
					<text class="label">数量</text>
					<text class="value highlight">{{ batchData.batch.quantity }}</text>
				</view>
			</view>
		</view>
		
		<!-- 工序选择 -->
		<view class="procedure-section">
			<view class="section-title">选择工序</view>
			<view class="procedure-list">
				<view 
					v-for="proc in batchData.order.procedures" 
					:key="proc.sequence"
					:class="['procedure-item', { 'active': selectedProcedure?.sequence === proc.sequence }]"
					@click="selectProcedure(proc)"
				>
					<view class="proc-info">
						<view class="proc-name">
							<text class="seq-badge">{{ proc.sequence }}</text>
							<text>{{ proc.procedure_name }}</text>
						</view>
						<view class="proc-price">{{ proc.unit_price }}元/件</view>
					</view>
					<view class="proc-progress" v-if="getProcedureProgress(proc.sequence)">
						<text>已上报：{{ getProcedureProgress(proc.sequence).reported_qty }}/{{ batchData.batch.quantity }}</text>
						<text v-if="getProcedureProgress(proc.sequence).is_completed" class="completed-badge">✓</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 上报信息 -->
		<view class="report-info" v-if="selectedProcedure">
			<view class="info-row">
				<text class="label">上报数量</text>
				<text class="value">{{ batchData.batch.quantity }} 件</text>
			</view>
			<view class="info-row">
				<text class="label">单价</text>
				<text class="value">{{ selectedProcedure.unit_price }} 元/件</text>
			</view>
			<view class="info-row highlight">
				<text class="label">预计工资</text>
				<text class="value">¥{{ calculateSalary() }}</text>
			</view>
		</view>
		
		<!-- 备注 -->
		<view class="remark-section">
			<textarea 
				v-model="remark" 
				placeholder="备注（可选）" 
				maxlength="200"
				class="remark-input"
			/>
		</view>
		
		<!-- 提交按钮 -->
		<view class="submit-section">
			<button 
				type="primary" 
				:disabled="!selectedProcedure || submitting"
				@click="submitReport"
				class="submit-btn"
			>
				{{ submitting ? '提交中...' : '提交上报' }}
			</button>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { submitProcedureReport } from '@/api/scan'

const batchData = ref({
	batch: {},
	order: { procedures: [] },
	batch_progress: []
})
const selectedProcedure = ref(null)
const remark = ref('')
const submitting = ref(false)

// 页面加载时获取参数
onMounted(() => {
	// 获取页面参数
	const pages = getCurrentPages()
	const currentPage = pages[pages.length - 1]
	const options = currentPage.options || {}
	
	if (options.data) {
		try {
			batchData.value = JSON.parse(decodeURIComponent(options.data))
			console.log('批次数据:', batchData.value)
		} catch (error) {
			console.error('解析数据失败:', error)
			uni.showToast({
				title: '数据错误',
				icon: 'none'
			})
			setTimeout(() => {
				uni.navigateBack()
			}, 1500)
		}
	}
})

// 选择工序
const selectProcedure = (proc) => {
	selectedProcedure.value = proc
}

// 获取工序进度
const getProcedureProgress = (seq) => {
	return batchData.value.batch_progress.find(p => p.procedure_seq === seq)
}

// 计算预计工资
const calculateSalary = () => {
	if (!selectedProcedure.value) return '0.00'
	return (batchData.value.batch.quantity * selectedProcedure.value.unit_price).toFixed(2)
}

// 提交上报
const submitReport = async () => {
	if (!selectedProcedure.value) {
		uni.showToast({ title: '请选择工序', icon: 'none' })
		return
	}
	
	if (submitting.value) return
	
	// 检查是否已完成上报
	const progress = getProcedureProgress(selectedProcedure.value.sequence)
	if (progress && progress.is_completed) {
		const confirmed = await new Promise(resolve => {
			uni.showModal({
				title: '提示',
				content: '该批次该工序已全部上报，确定要重复上报吗？',
				success: (res) => resolve(res.confirm)
			})
		})
		if (!confirmed) return
	}
	
	await doSubmit()
}

const doSubmit = async () => {
	try {
		submitting.value = true
		uni.showLoading({ title: '提交中...' })
		
		const params = {
			order_id: batchData.value.order.id,
			batch_id: batchData.value.batch.id,
			bundle_no: batchData.value.batch.bundle_no,
			procedure_seq: selectedProcedure.value.sequence,
			procedure_name: selectedProcedure.value.procedure_name,
			quantity: batchData.value.batch.quantity,
			color: batchData.value.batch.color,
			size: batchData.value.batch.size,
			remark: remark.value
		}
		
		const res = await submitProcedureReport(params)
		
		uni.hideLoading()
		submitting.value = false
		
		// 显示上报成功
		uni.showModal({
			title: '上报成功',
			content: `工资：¥${res.data.total_price}`,
			showCancel: false,
			confirmText: '继续扫码',
			success: () => {
				// 返回扫码页，继续扫码
				uni.navigateBack()
			}
		})
	} catch (error) {
		uni.hideLoading()
		submitting.value = false
		
		console.error('上报失败:', error)
		uni.showModal({
			title: '上报失败',
			content: error.message || '网络错误，请重试',
			showCancel: false
		})
	}
}
</script>

<style lang="scss" scoped>
.report-detail {
	min-height: 100vh;
	background: #f5f5f5;
	padding-bottom: 32rpx;
}

.batch-info {
	background: #fff;
	margin: 24rpx;
	border-radius: 16rpx;
	padding: 32rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
	
	.info-header {
		display: flex;
		margin-bottom: 24rpx;
		
		.style-image {
			width: 120rpx;
			height: 120rpx;
			border-radius: 12rpx;
			margin-right: 24rpx;
			background: #f0f0f0;
		}
		
		.info-text {
			flex: 1;
			display: flex;
			flex-direction: column;
			justify-content: center;
			
			.style-name {
				font-size: 32rpx;
				font-weight: bold;
				color: #333;
				margin-bottom: 8rpx;
			}
			
			.contract-no {
				font-size: 26rpx;
				color: #999;
			}
		}
	}
	
	.info-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 24rpx;
		padding-top: 24rpx;
		border-top: 1rpx solid #f0f0f0;
		
		.info-item {
			display: flex;
			flex-direction: column;
			align-items: center;
			
			.label {
				font-size: 24rpx;
				color: #999;
				margin-bottom: 8rpx;
			}
			
			.value {
				font-size: 28rpx;
				color: #333;
				font-weight: 500;
				
				&.highlight {
					color: #5EA3F2;
					font-size: 32rpx;
					font-weight: bold;
				}
			}
		}
	}
}

.procedure-section {
	background: #fff;
	margin: 0 24rpx 24rpx;
	border-radius: 16rpx;
	padding: 32rpx;
	
	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 24rpx;
	}
	
	.procedure-list {
		.procedure-item {
			border: 2rpx solid #e0e0e0;
			border-radius: 12rpx;
			padding: 24rpx;
			margin-bottom: 16rpx;
			transition: all 0.3s;
			
			&:last-child {
				margin-bottom: 0;
			}
			
			&.active {
				border-color: #5EA3F2;
				background: rgba(94, 163, 242, 0.05);
			}
			
			.proc-info {
				display: flex;
				justify-content: space-between;
				align-items: center;
				margin-bottom: 8rpx;
				
				.proc-name {
					display: flex;
					align-items: center;
					font-size: 28rpx;
					color: #333;
					
					.seq-badge {
						display: inline-block;
						width: 40rpx;
						height: 40rpx;
						line-height: 40rpx;
						text-align: center;
						background: #5EA3F2;
						color: #fff;
						border-radius: 50%;
						font-size: 22rpx;
						margin-right: 12rpx;
					}
				}
				
				.proc-price {
					font-size: 28rpx;
					color: #ff6b6b;
					font-weight: bold;
				}
			}
			
			.proc-progress {
				font-size: 24rpx;
				color: #999;
				display: flex;
				align-items: center;
				gap: 8rpx;
				
				.completed-badge {
					color: #4CAF50;
					font-size: 28rpx;
				}
			}
		}
	}
}

.report-info {
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	margin: 0 24rpx 24rpx;
	border-radius: 16rpx;
	padding: 32rpx;
	color: #fff;
	
	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12rpx 0;
		font-size: 28rpx;
		
		&.highlight {
			padding-top: 20rpx;
			margin-top: 12rpx;
			border-top: 1rpx solid rgba(255, 255, 255, 0.3);
			font-size: 32rpx;
			font-weight: bold;
		}
	}
}

.remark-section {
	background: #fff;
	margin: 0 24rpx 24rpx;
	border-radius: 16rpx;
	padding: 24rpx;
	
	.remark-input {
		width: 100%;
		min-height: 150rpx;
		font-size: 28rpx;
		line-height: 1.6;
	}
}

.submit-section {
	padding: 0 24rpx;
	
	.submit-btn {
		width: 100%;
		height: 96rpx;
		background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
		border: none;
		border-radius: 48rpx;
		font-size: 32rpx;
		font-weight: bold;
		color: #fff;
		box-shadow: 0 8rpx 24rpx rgba(94, 163, 242, 0.4);
		
		&:disabled {
			opacity: 0.6;
		}
	}
}
</style>

