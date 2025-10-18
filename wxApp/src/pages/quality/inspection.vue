<template>
	<view class="inspection-page">
		<!-- 批次信息卡片 -->
		<view class="batch-info" v-if="batchData.batch">
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

		<!-- 质检表单 -->
		<view class="form-section">
			<view class="section-title">质检信息</view>
			
			<!-- 工序选择 -->
			<view class="form-item">
				<view class="label">选择工序</view>
			<picker mode="selector" :range="procedures" range-key="procedure_name" @change="onProcedureChange">
				<view class="picker">
					{{ selectedProcedure ? selectedProcedure.procedure_name : '请选择工序' }}
				</view>
			</picker>
			</view>

			<!-- 质检数量 -->
			<view class="form-item">
				<view class="label">质检数量</view>
				<input type="number" v-model.number="form.inspected_qty" placeholder="请输入质检数量" />
			</view>

			<!-- 合格数量 -->
			<view class="form-item">
				<view class="label">合格数量</view>
				<input type="number" v-model.number="form.qualified_qty" placeholder="请输入合格数量" />
			</view>

			<!-- 不合格数量 -->
			<view class="form-item">
				<view class="label">不合格数量</view>
				<input type="number" v-model.number="form.unqualified_qty" placeholder="请输入不合格数量" class="error-input" />
			</view>

			<!-- 合格率 -->
			<view class="quality-rate" v-if="form.inspected_qty > 0">
				<text>合格率：</text>
				<text class="rate" :class="{ 'low-rate': qualityRate < 95 }">{{ qualityRate.toFixed(2) }}%</text>
			</view>

			<!-- 缺陷类型 -->
			<view class="form-item" v-if="form.unqualified_qty > 0">
				<view class="label">缺陷类型</view>
				<view class="defect-tags">
					<view v-for="(defect, index) in defectOptions" :key="index" 
						:class="['defect-tag', { 'active': form.defect_types.includes(defect) }]"
						@click="toggleDefect(defect)">
						{{ defect }}
					</view>
				</view>
			</view>

			<!-- 缺陷描述 -->
			<view class="form-item" v-if="form.unqualified_qty > 0">
				<view class="label">缺陷描述</view>
				<textarea v-model="form.defect_desc" placeholder="请描述具体缺陷情况" maxlength="200" />
			</view>

			<!-- 备注 -->
			<view class="form-item">
				<view class="label">备注</view>
				<textarea v-model="form.remark" placeholder="备注（可选）" maxlength="200" />
			</view>
		</view>

		<!-- 提交按钮 -->
		<view class="submit-section">
			<button type="primary" @click="submit" :disabled="submitting">
				{{ submitting ? '提交中...' : '提交质检' }}
			</button>
		</view>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { submitInspection } from '@/api/quality'

const batchData = ref({
	batch: {},
	order: { procedures: [] }
})
const selectedProcedure = ref(null)
const procedures = computed(() => batchData.value.order.procedures || [])

const form = ref({
	inspected_qty: 0,
	qualified_qty: 0,
	unqualified_qty: 0,
	defect_types: [],
	defect_desc: '',
	remark: ''
})

const submitting = ref(false)

// 缺陷类型选项
const defectOptions = ['线头', '污渍', '破损', '色差', '尺寸不符', '其他']

// 计算合格率
const qualityRate = computed(() => {
	if (form.value.inspected_qty <= 0) return 0
	return (form.value.qualified_qty / form.value.inspected_qty) * 100
})

// 页面加载时获取参数
onMounted(() => {
	const pages = getCurrentPages()
	const currentPage = pages[pages.length - 1]
	const options = currentPage.options || {}
	
	if (options.data) {
		try {
			batchData.value = JSON.parse(decodeURIComponent(options.data))
		} catch (error) {
			console.error('解析数据失败:', error)
		}
	}
})

// 选择工序
const onProcedureChange = (e) => {
	selectedProcedure.value = procedures.value[e.detail.value]
}

// 切换缺陷类型
const toggleDefect = (defect) => {
	const index = form.value.defect_types.indexOf(defect)
	if (index > -1) {
		form.value.defect_types.splice(index, 1)
	} else {
		form.value.defect_types.push(defect)
	}
}

// 提交质检
const submit = async () => {
	if (submitting.value) return

	// 校验
	if (!selectedProcedure.value) {
		uni.showToast({ title: '请选择工序', icon: 'none' })
		return
	}
	if (!form.value.inspected_qty || form.value.inspected_qty <= 0) {
		uni.showToast({ title: '请输入质检数量', icon: 'none' })
		return
	}
	if (form.value.qualified_qty + form.value.unqualified_qty !== form.value.inspected_qty) {
		uni.showToast({ title: '合格数量与不合格数量之和必须等于质检数量', icon: 'none' })
		return
	}

	try {
		submitting.value = true
		uni.showLoading({ title: '提交中...' })

		const params = {
			order_id: batchData.value.order.id,
			batch_id: batchData.value.batch.id,
			bundle_no: batchData.value.batch.bundle_no,
			procedure_seq: selectedProcedure.value.sequence,
			procedure_name: selectedProcedure.value.procedure_name,
			inspected_qty: form.value.inspected_qty,
			qualified_qty: form.value.qualified_qty,
			unqualified_qty: form.value.unqualified_qty,
			defect_types: form.value.defect_types,
			defect_desc: form.value.defect_desc,
			color: batchData.value.batch.color,
			size: batchData.value.batch.size,
			remark: form.value.remark
		}

		const res = await submitInspection(params)

		uni.hideLoading()
		submitting.value = false

		// 显示结果
		uni.showModal({
			title: '质检提交成功',
			content: `合格率: ${res.data.quality_rate.toFixed(2)}%${res.data.need_rework ? '\n\n检测到不合格品，需要返工' : ''}`,
			showCancel: false,
			success: () => {
				uni.navigateBack()
			}
		})
	} catch (error) {
		uni.hideLoading()
		submitting.value = false
		console.error('质检提交失败:', error)
		uni.showModal({
			title: '提交失败',
			content: error.message || '网络错误，请重试',
			showCancel: false
		})
	}
}
</script>

<style lang="scss" scoped>
.inspection-page {
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
				margin-bottom: 12rpx;
			}

			.contract-no {
				font-size: 26rpx;
				color: #666;
			}
		}
	}

	.info-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 24rpx;

		.info-item {
			text-align: center;

			.label {
				display: block;
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
					font-weight: bold;
				}
			}
		}
	}
}

.form-section {
	background: #fff;
	margin: 24rpx;
	border-radius: 16rpx;
	padding: 32rpx;

	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 32rpx;
	}

	.form-item {
		margin-bottom: 32rpx;

		.label {
			font-size: 28rpx;
			color: #666;
			margin-bottom: 16rpx;
		}

		input, textarea {
			width: 100%;
			padding: 24rpx;
			border: 1rpx solid #ddd;
			border-radius: 8rpx;
			font-size: 28rpx;
		}

		.error-input {
			border-color: #ff4444;
		}

		.picker {
			padding: 24rpx;
			border: 1rpx solid #ddd;
			border-radius: 8rpx;
			font-size: 28rpx;
		}

		textarea {
			height: 150rpx;
		}
	}

	.quality-rate {
		padding: 24rpx;
		background: #f0f0f0;
		border-radius: 8rpx;
		text-align: center;
		margin-bottom: 32rpx;
		font-size: 28rpx;

		.rate {
			font-size: 36rpx;
			font-weight: bold;
			color: #4CAF50;

			&.low-rate {
				color: #ff4444;
			}
		}
	}

	.defect-tags {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;

		.defect-tag {
			padding: 12rpx 24rpx;
			border: 1rpx solid #ddd;
			border-radius: 24rpx;
			font-size: 24rpx;
			color: #666;

			&.active {
				background: #5EA3F2;
				border-color: #5EA3F2;
				color: #fff;
			}
		}
	}
}

.submit-section {
	padding: 0 24rpx;

	button {
		border-radius: 48rpx;
		height: 96rpx;
		font-size: 32rpx;
	}
}
</style>

