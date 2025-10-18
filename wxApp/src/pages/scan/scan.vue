<template>
	<view class="scan-container">
		<!-- 自定义相机扫码界面 -->
		<view class="camera-container" v-if="showCamera">
			<!-- 相机组件 -->
			<camera 
				class="camera" 
				device-position="back" 
				:flash="flashOn ? 'torch' : 'off'"
				mode="scanCode"
				@scancode="onScanCode"
				@error="onCameraError"
			>
				<!-- 扫码框 -->
				<cover-view class="camera-cover">
					<!-- 中间扫码框 -->
					<cover-view class="scan-area">
						<cover-view :class="['scan-frame', { 'scan-success': scanSuccess }]">
						</cover-view>
						<cover-view class="scan-tips">
							<cover-text>将二维码放入框内扫描</cover-text>
						</cover-view>
					</cover-view>
					
					<!-- 底部按钮 -->
					<cover-view class="camera-footer">
						<cover-view class="footer-btn" @tap="closeScan">
							<cover-text class="footer-text">X</cover-text>
						</cover-view>
					</cover-view>
				</cover-view>
			</camera>
		</view>
		
		<!-- 引导页面 -->
		<view class="guide-page" v-else>
			<!-- 扫码提示区域 -->
			<view class="scan-header">
				<view class="scan-icon-wrapper">
					<u-icon name="scan" :size="120" color="#5EA3F2"></u-icon>
				</view>
				<view class="scan-title">扫码识别</view>
				<view class="scan-desc">点击下方按钮开始扫码或手动输入编号</view>
			</view>

			<!-- 扫码按钮 -->
			<view class="scan-actions">
				<u-button 
					type="primary" 
					:custom-style="{
						background: 'linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%)',
						border: 'none',
						borderRadius: '48rpx',
						height: '120rpx',
						fontSize: '36rpx',
						fontWeight: 'bold',
						boxShadow: '0 8rpx 24rpx rgba(94, 163, 242, 0.4)'
					}"
					@click="openScan"
				>
					<view class="scan-btn-content">
						<u-icon name="scan" :size="48" color="#fff"></u-icon>
						<text class="scan-text">开始扫码</text>
					</view>
				</u-button>

				<u-button 
					type="info" 
					plain
					:custom-style="{
						borderColor: '#5EA3F2',
						color: '#5EA3F2',
						borderRadius: '48rpx',
						height: '120rpx',
						marginTop: '32rpx',
						fontSize: '36rpx'
					}"
					@click="handleManualInput"
				>
					<view class="scan-btn-content">
						<u-icon name="edit-pen" :size="44" color="#5EA3F2"></u-icon>
						<text>手动输入</text>
					</view>
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
							<view class="record-subtitle">{{ record.subtitle }}</view>
							<view class="record-time">{{ record.time }}</view>
						</view>
						<view class="record-status" :class="record.statusClass">
							{{ record.status }}
						</view>
					</view>
					
					<!-- 空状态 -->
					<view class="empty-state" v-if="recentRecords.length === 0">
						<text class="empty-text">暂无扫码记录</text>
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
import { parseScanCode, getReportList } from '@/api/scan'

const userStore = useUserStore()
const showCamera = ref(false)
const flashOn = ref(false)
const scanSuccess = ref(false)
const isScanning = ref(false) // 防止重复扫描

// 最近扫码记录
const recentRecords = ref([])

onMounted(() => {
	loadRecentRecords()
})

// 加载最近的上报记录
const loadRecentRecords = async () => {
	try {
		const res = await getReportList({ page: 1, page_size: 5 })
		if (res.data && res.data.reports) {
			recentRecords.value = res.data.reports.map(record => ({
				id: record.id,
				title: `${record.procedure_name} - ${record.contract_no}`,
				subtitle: `${record.color} / ${record.size || '-'} / ${record.quantity}件`,
				time: formatTimeAgo(record.report_time),
				status: '已上报',
				statusClass: 'success',
				rawData: record
			}))
		}
	} catch (error) {
		console.error('加载最近记录失败:', error)
		// 失败时不显示错误，保持空列表
	}
}

// 格式化时间为"xx前"
const formatTimeAgo = (timestamp) => {
	const now = Math.floor(Date.now() / 1000)
	const diff = now - timestamp
	
	if (diff < 60) return '刚刚'
	if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
	if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
	if (diff < 2592000) return `${Math.floor(diff / 86400)}天前`
	
	// 超过30天显示日期
	const date = new Date(timestamp * 1000)
	return `${date.getMonth() + 1}月${date.getDate()}日`
}

// 打开扫码界面
const openScan = () => {
	scanSuccess.value = false
	isScanning.value = false // 重置扫描状态
	// 请求相机权限
	uni.authorize({
		scope: 'scope.camera',
		success: () => {
			showCamera.value = true
		},
		fail: () => {
			uni.showModal({
				title: '需要相机权限',
				content: '请在设置中开启相机权限',
				confirmText: '去设置',
				success: (res) => {
					if (res.confirm) {
						uni.openSetting()
					}
				}
			})
		}
	})
}

// 关闭扫码界面
const closeScan = () => {
	showCamera.value = false
	flashOn.value = false
	scanSuccess.value = false
	isScanning.value = false
}

// 切换闪光灯
const toggleFlash = () => {
	flashOn.value = !flashOn.value
}

// 扫码成功回调
const onScanCode = async (e) => {
	console.log('扫码结果:', e.detail)
	const result = e.detail.result
	
	if (!result) return
	
	// 防止重复扫描
	if (isScanning.value) {
		console.log('正在处理中，忽略重复扫描')
		return
	}
	isScanning.value = true
	
	// 扫描成功，边框高亮
	scanSuccess.value = true
	
	// 延迟一下让用户看到边框高亮效果
	setTimeout(async () => {
		try {
			uni.showLoading({ title: '识别中...' })
			
			// 调用后端解析扫码内容
			const res = await parseScanCode({ qr_code: result })
			
			uni.hideLoading()
			
			// 关闭相机
			showCamera.value = false
			scanSuccess.value = false
			isScanning.value = false
			
			// 跳转到上报详情页
			uni.navigateTo({
				url: `/pages/report/detail?data=${encodeURIComponent(JSON.stringify(res.data))}`
			})
		} catch (error) {
			uni.hideLoading()
			showCamera.value = false
			scanSuccess.value = false
			isScanning.value = false
			
			console.error('扫码识别失败:', error)
			uni.showModal({
				title: '识别失败',
				content: error.message || '无法识别该二维码，请重试',
				showCancel: false,
				confirmText: '重新扫码',
				success: (modalRes) => {
					if (modalRes.confirm) {
						// 重新扫码
						setTimeout(() => {
							openScan()
						}, 300)
					}
				}
			})
		}
	}, 300)
}

// 相机错误回调
const onCameraError = (e) => {
	console.error('相机错误:', e.detail)
	showCamera.value = false
	uni.showToast({
		title: '相机启动失败',
		icon: 'none'
	})
}

// 从相册选择二维码
const chooseFromAlbum = () => {
	uni.chooseImage({
		count: 1,
		sourceType: ['album'],
		success: (res) => {
			// 关闭相机
			showCamera.value = false
			
			// TODO: 识别图片中的二维码
			uni.showToast({
				title: '功能开发中',
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
	openScan()
}

// 查看全部记录
const handleViewAll = () => {
	uni.navigateTo({
		url: '/pages/report/history'
	})
}

// 查看记录详情
const handleRecordDetail = (record) => {
	console.log('查看记录:', record)
	// 跳转到上报记录历史页面
	uni.navigateTo({
		url: '/pages/report/history'
	})
}
</script>

<style lang="scss" scoped>
.scan-container {
	min-height: 100vh;
	background: linear-gradient(180deg, #F0F8FF 0%, #FFFFFF 40%);
	padding-bottom: 120rpx;
}

// 相机容器
.camera-container {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	z-index: 9999;
	background: #000;
}

.camera {
	width: 100%;
	height: 100%;
	background: rgba(0, 0, 0, 0.3);
}

// 相机遮罩层
.camera-cover {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	display: flex;
	flex-direction: column;
}

// 中间扫码区域
.scan-area {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
}

.scan-frame {
	width: 500rpx;
	height: 500rpx;
	border: 4rpx solid #0f0;
	opacity: 0.3;
	transition: opacity 0.3s;
	
	&.scan-success {
		opacity: 1;
	}
}

.scan-tips {
	margin-top: 40rpx;
	
	cover-text {
		font-size: 28rpx;
		color: #fff;
		text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.8);
	}
}

// 底部按钮
.camera-footer {
	display: flex;
	align-items: center;
	justify-content: space-around;
	padding: 60rpx 120rpx;
	padding-bottom: calc(60rpx + env(safe-area-inset-bottom));
}

.footer-btn {
	width: 120rpx;
	height: 120rpx;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	
	.footer-text {
		width: 60rpx;
		height: 60rpx;
		font-size: 32rpx;
		color: #fff;
		text-align: center;
		line-height: 60rpx;
	}
}

// 引导页面
.guide-page {
	min-height: 100vh;
}

.scan-header {
	padding: 100rpx 40rpx 60rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	text-align: center;
	
	.scan-icon-wrapper {
		width: 200rpx;
		height: 200rpx;
		background: linear-gradient(135deg, rgba(94, 163, 242, 0.1) 0%, rgba(79, 195, 247, 0.1) 100%);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 40rpx;
		box-shadow: 0 8rpx 32rpx rgba(94, 163, 242, 0.15);
	}
	
	.scan-title {
		font-size: 44rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 20rpx;
		letter-spacing: 2rpx;
	}
	
	.scan-desc {
		font-size: 28rpx;
		color: #999;
		line-height: 1.6;
		max-width: 500rpx;
	}
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
			font-size: 36rpx;
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
					margin-bottom: 6rpx;
					font-weight: 500;
				}
				
				.record-subtitle {
					font-size: 24rpx;
					color: #666;
					margin-bottom: 6rpx;
				}
				
				.record-time {
					font-size: 22rpx;
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
		
		.empty-state {
			padding: 80rpx 0;
			text-align: center;
			
			.empty-text {
				font-size: 26rpx;
				color: #999;
			}
		}
	}
}
</style>
