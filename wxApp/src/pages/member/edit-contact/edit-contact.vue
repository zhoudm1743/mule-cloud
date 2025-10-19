<template>
	<view class="edit-contact-container">
		<view class="form-section">
			<view class="section-title">联系信息</view>
			<view class="form-group">
				<view class="form-item">
					<view class="item-label">手机号</view>
					<input 
						class="item-input" 
						type="number"
						v-model="formData.phone" 
						placeholder="请输入手机号"
						maxlength="11"
					/>
				</view>
				
				<view class="form-item">
					<view class="item-label">邮箱</view>
					<input 
						class="item-input" 
						type="text"
						v-model="formData.email" 
						placeholder="请输入邮箱"
					/>
				</view>

				<view class="form-item">
					<view class="item-label">家庭住址</view>
					<textarea 
						class="item-textarea" 
						v-model="formData.address" 
						placeholder="请输入家庭住址"
						maxlength="200"
					/>
				</view>
			</view>
		</view>

		<view class="form-section">
			<view class="section-title">紧急联系人</view>
			<view class="form-group">
				<view class="form-item">
					<view class="item-label">姓名</view>
					<input 
						class="item-input" 
						type="text"
						v-model="formData.emergency_contact" 
						placeholder="请输入紧急联系人姓名"
						maxlength="20"
					/>
				</view>
				
				<view class="form-item">
					<view class="item-label">电话</view>
					<input 
						class="item-input" 
						type="number"
						v-model="formData.emergency_phone" 
						placeholder="请输入紧急联系电话"
						maxlength="11"
					/>
				</view>

				<view class="form-item">
					<view class="item-label">关系</view>
					<input 
						class="item-input" 
						type="text"
						v-model="formData.emergency_relation" 
						placeholder="如：父亲、配偶等"
						maxlength="20"
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
import { getProfile, updateContactInfo } from '@/api/member'

const saving = ref(false)
const formData = ref({
	phone: '',
	email: '',
	address: '',
	emergency_contact: '',
	emergency_phone: '',
	emergency_relation: ''
})

onMounted(async () => {
	await loadProfile()
})

// 加载个人档案
const loadProfile = async () => {
	try {
		const res = await getProfile()
		if (res.code === 0) {
			const profile = res.data || {}
			formData.value = {
				phone: profile.phone || '',
				email: profile.email || '',
				address: profile.address || '',
				emergency_contact: profile.emergency_contact || '',
				emergency_phone: profile.emergency_phone || '',
				emergency_relation: profile.emergency_relation || ''
			}
		}
	} catch (error) {
		console.error('加载个人档案失败', error)
	}
}

// 保存
const handleSave = async () => {
	saving.value = true
	try {
		const res = await updateContactInfo(formData.value)
		
		if (res.code === 0) {
			uni.showToast({
				title: '保存成功',
				icon: 'success'
			})

			// 返回上一页
			setTimeout(() => {
				uni.navigateBack()
			}, 1000)
		} else {
			uni.showToast({
				title: res.msg || '保存失败',
				icon: 'none'
			})
		}
	} catch (error) {
		uni.showToast({
			title: '保存失败',
			icon: 'none'
		})
	} finally {
		saving.value = false
	}
}
</script>

<style scoped lang="scss">
.edit-contact-container {
	min-height: 100vh;
	background: #f5f7fa;
	padding-bottom: 120rpx;
}

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

.item-textarea {
	flex: 1;
	font-size: 30rpx;
	color: #333;
	min-height: 120rpx;
	padding: 16rpx 0;
	line-height: 1.6;
}

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
