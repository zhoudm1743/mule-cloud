<template>
	<view class="edit-profile-container">
		<!-- 头像 -->
		<view class="avatar-section" @click="handleChooseAvatar">
			<image :src="formData.avatar || '/static/logo.png'" mode="aspectFill"></image>
			<view class="avatar-tip">
				<u-icon name="camera-fill" :size="32" color="#5EA3F2"></u-icon>
				<text>更换头像</text>
			</view>
		</view>

		<!-- 个人信息 -->
		<view class="form-section">
			<view class="section-title">个人信息</view>
			<view class="form-group">
				<view class="form-item">
					<view class="item-label">姓名</view>
					<input 
						class="item-input" 
						v-model="formData.name" 
						placeholder="请输入姓名"
						maxlength="20"
					/>
				</view>
				
				<view class="form-item">
					<view class="item-label">性别</view>
					<view class="item-value">
						<u-radio-group v-model="formData.gender" direction="row">
							<u-radio :name="1" label="男" activeColor="#5EA3F2"></u-radio>
							<u-radio :name="2" label="女" activeColor="#5EA3F2" style="margin-left: 60rpx;"></u-radio>
						</u-radio-group>
					</view>
				</view>

				<view class="form-item">
					<view class="item-label">联系电话</view>
					<view class="item-value">
						<text class="phone-text">{{ formData.phone || '未绑定' }}</text>
						<text class="bind-tip" @click="handleBindPhone" v-if="!formData.phone">去绑定</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 企业信息 -->
		<view class="form-section">
			<view class="section-title">企业信息</view>
			<view class="form-group">
				<view class="form-item">
					<view class="item-label">工号</view>
					<input 
						class="item-input" 
						v-model="formData.jobNumber" 
						placeholder="请输入工号"
						maxlength="20"
					/>
				</view>
				
				<view class="form-item">
					<view class="item-label">部门</view>
					<input 
						class="item-input" 
						v-model="formData.department" 
						placeholder="请输入部门"
						maxlength="50"
					/>
				</view>

				<view class="form-item">
					<view class="item-label">岗位</view>
					<input 
						class="item-input" 
						v-model="formData.position" 
						placeholder="请输入岗位"
						maxlength="50"
					/>
				</view>
			</view>
		</view>

		<!-- 保存按钮 -->
		<view class="save-section">
			<view class="save-btn" @click="handleSave">
				<text v-if="!saving">保存</text>
				<text v-else>保存中...</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/store/modules/user'
import { updateUserInfo } from '@/api/auth'

const userStore = useUserStore()
const saving = ref(false)

const formData = ref({
	name: '',
	avatar: '',
	gender: 0,
	phone: '',
	jobNumber: '',
	department: '',
	position: ''
})

onMounted(async () => {
	// 重新获取用户信息，确保拿到最新的成员详情
	try {
		await userStore.fetchUserInfo()
	} catch (error) {
		console.error('获取用户信息失败', error)
	}
	
	// 获取用户信息和租户成员信息
	const userInfo = userStore.userInfo || {}
	const tenants = userStore.tenants || []
	const currentTenant = userStore.currentTenant || {}
	
	// 从 tenants 数组中查找当前租户，获取完整的成员信息
	const currentTenantDetail = tenants.find(t => t.tenant_id === currentTenant.tenant_id) || currentTenant
	
	formData.value = {
		name: userInfo.nickname || '',
		avatar: userInfo.avatar || '',
		gender: userInfo.gender || 0,
		phone: userInfo.phone || '',
		jobNumber: currentTenantDetail.job_number || '',
		department: currentTenantDetail.department || '',
		position: currentTenantDetail.position || ''
	}
})

// 选择头像
const handleChooseAvatar = () => {
	uni.chooseImage({
		count: 1,
		sizeType: ['compressed'],
		sourceType: ['album', 'camera'],
		success: (res) => {
			formData.value.avatar = res.tempFilePaths[0]
			uni.showToast({
				title: '头像已选择',
				icon: 'none'
			})
		}
	})
}

// 绑定手机号
const handleBindPhone = () => {
	uni.showModal({
		title: '绑定手机号',
		content: '此功能需要在"我的"页面进行操作',
		showCancel: true,
		confirmText: '去绑定',
		success: (res) => {
			if (res.confirm) {
				uni.navigateBack()
			}
		}
	})
}

// 保存
const handleSave = async () => {
	if (!formData.value.name.trim()) {
		uni.showToast({
			title: '请输入姓名',
			icon: 'none'
		})
		return
	}

	saving.value = true
	try {
		await updateUserInfo({
			nickname: formData.value.name,
			avatar: formData.value.avatar,
			gender: formData.value.gender,
			job_number: formData.value.jobNumber,
			department: formData.value.department,
			position: formData.value.position
		})

		// 更新本地状态（全局用户信息）
		userStore.updateUserInfo({
			nickname: formData.value.name,
			avatar: formData.value.avatar,
			gender: formData.value.gender
		})

		// 重新获取用户完整信息（包括更新后的租户成员信息）
		await userStore.fetchUserInfo()

		uni.showToast({
			title: '保存成功',
			icon: 'success'
		})

		// 返回上一页
		setTimeout(() => {
			uni.navigateBack()
		}, 1000)
	} catch (error) {
		uni.showToast({
			title: error.message || '保存失败',
			icon: 'none'
		})
	} finally {
		saving.value = false
	}
}
</script>

<style scoped lang="scss">
.edit-profile-container {
	min-height: 100vh;
	background: #f5f7fa;
	padding-bottom: 120rpx;
}

// 头像区域
.avatar-section {
	position: relative;
	width: 200rpx;
	height: 200rpx;
	margin: 60rpx auto 40rpx;
	border-radius: 50%;
	overflow: hidden;
	box-shadow: 0 8rpx 24rpx rgba(94, 163, 242, 0.2);
	
	image {
		width: 100%;
		height: 100%;
	}
	
	.avatar-tip {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 80rpx;
		background: rgba(0, 0, 0, 0.6);
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 4rpx;
		
		text {
			font-size: 20rpx;
			color: #fff;
		}
	}
}

// 表单区域
.form-section {
	margin: 24rpx 32rpx;
}

.section-title {
	font-size: 26rpx;
	color: #999;
	padding: 0 8rpx 16rpx;
	font-weight: bold;
}

.form-group {
	background: #fff;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
}

.form-item {
	display: flex;
	align-items: center;
	padding: 32rpx;
	border-bottom: 1rpx solid #f5f7fa;
	
	&:last-child {
		border-bottom: none;
	}
}

.item-label {
	width: 140rpx;
	font-size: 30rpx;
	color: #333;
	font-weight: 500;
}

.item-input {
	flex: 1;
	text-align: right;
	font-size: 30rpx;
	color: #333;
}

.item-value {
	flex: 1;
	display: flex;
	align-items: center;
	justify-content: flex-end;
}

.phone-text {
	font-size: 30rpx;
	color: #666;
}

.bind-tip {
	margin-left: 16rpx;
	padding: 8rpx 24rpx;
	background: #e8f4ff;
	color: #5EA3F2;
	font-size: 24rpx;
	border-radius: 32rpx;
}

// 保存按钮
.save-section {
	padding: 40rpx 32rpx 32rpx;
}

.save-btn {
	padding: 32rpx;
	background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
	border-radius: 48rpx;
	box-shadow: 0 8rpx 16rpx rgba(94, 163, 242, 0.3);
	
	text {
		display: block;
		text-align: center;
		font-size: 32rpx;
		color: #fff;
		font-weight: bold;
	}
	
	&:active {
		opacity: 0.8;
	}
}
</style>
