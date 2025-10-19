<template>
	<view class="edit-profile-container">
		<!-- 头像 -->
		<view class="photos-section">
			<view class="photo-item" @click="handleChooseAvatar">
				<image :src="formData.avatar || '/static/logo.png'" mode="aspectFill" class="photo-img"></image>
				<view class="photo-tip">
					<u-icon v-if="!uploading.avatar" name="camera-fill" :size="28" color="#5EA3F2"></u-icon>
					<text>{{ uploading.avatar ? '上传中...' : '头像' }}</text>
				</view>
			</view>
		</view>

		<!-- 基本信息 -->
		<view class="form-section">
			<view class="section-title">基本信息</view>
			<view class="form-group">
				<view class="form-item">
					<view class="item-label">姓名</view>
					<input 
						class="item-input" 
						type="text"
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
					<view class="item-label">身份证号</view>
					<input 
						class="item-input" 
						type="idcard"
						v-model="formData.id_card_no" 
						placeholder="请输入身份证号"
						:disabled="!!originalIdCardNo"
						maxlength="18"
					/>
				</view>
				<view v-if="originalIdCardNo" class="form-tip">
					<text>身份证号首次填写后不可修改</text>
				</view>

				<view class="form-item" @click="showBirthdayPicker">
					<view class="item-label">出生日期</view>
					<view class="item-value">
						<text :class="{'placeholder': !formData.birthday}">
							{{ formData.birthday ? formatDate(formData.birthday) : '请选择出生日期' }}
						</text>
						<u-icon name="arrow-right" size="20" color="#999"></u-icon>
					</view>
				</view>

				<view class="form-item" @click="showNationPicker">
					<view class="item-label">民族</view>
					<view class="item-value">
						<text :class="{'placeholder': !formData.nation}">
							{{ formData.nation || '请选择民族' }}
						</text>
						<u-icon name="arrow-right" size="20" color="#999"></u-icon>
					</view>
				</view>

				<view class="form-item">
					<view class="item-label">籍贯</view>
					<input 
						class="item-input" 
						type="text"
						v-model="formData.native_place" 
						placeholder="如：广东深圳"
						maxlength="50"
					/>
				</view>

				<view class="form-item" @click="showMaritalStatusPicker">
					<view class="item-label">婚姻状况</view>
					<view class="item-value">
						<text :class="{'placeholder': !formData.marital_status}">
							{{ getMaritalStatusText(formData.marital_status) }}
						</text>
						<u-icon name="arrow-right" size="20" color="#999"></u-icon>
					</view>
				</view>

				<view class="form-item" @click="showPoliticalPicker">
					<view class="item-label">政治面貌</view>
					<view class="item-value">
						<text :class="{'placeholder': !formData.political}">
							{{ getPoliticalText(formData.political) }}
						</text>
						<u-icon name="arrow-right" size="20" color="#999"></u-icon>
					</view>
				</view>

				<view class="form-item" @click="showEducationPicker">
					<view class="item-label">学历</view>
					<view class="item-value">
						<text :class="{'placeholder': !formData.education}">
							{{ getEducationText(formData.education) }}
						</text>
						<u-icon name="arrow-right" size="20" color="#999"></u-icon>
					</view>
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

		<!-- 日期选择器 -->
		<u-datetime-picker
			:show="showDatePicker"
			v-model="selectedDate"
			mode="date"
			:max-date="maxDate"
			:min-date="minDate"
			@confirm="confirmBirthday"
			@cancel="showDatePicker = false"
		></u-datetime-picker>

		<!-- 选择器 -->
		<u-picker
			:show="showPicker"
			:columns="pickerColumns"
			@confirm="confirmPicker"
			@cancel="showPicker = false"
		></u-picker>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getProfile, updateBasicInfo, uploadPhoto } from '@/api/member'
import { uploadFile } from '@/api/auth'

const saving = ref(false)
const uploading = ref({
	avatar: false,
	photo: false
})

const originalIdCardNo = ref('') // 保存原始身份证号（用于判断是否可编辑）

const formData = ref({
	name: '',
	avatar: '',
	photo: '',
	gender: 0,
	id_card_no: '',
	birthday: 0,
	nation: '汉族',
	native_place: '',
	marital_status: 'single',
	political: 'masses',
	education: ''
})

// 日期选择器
const showDatePicker = ref(false)
const selectedDate = ref(new Date().getTime())
const maxDate = ref(new Date().getTime()) // 今天
const minDate = ref(new Date(1950, 0, 1).getTime()) // 1950年

// 选择器
const showPicker = ref(false)
const pickerColumns = ref([])
const pickerType = ref('') // nation / marital_status / political / education

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
				name: profile.name || '',
				avatar: profile.avatar || '',
				photo: profile.photo || '',
				gender: profile.gender || 0,
				id_card_no: profile.id_card_no || '',
				birthday: profile.birthday || 0,
				nation: profile.nation || '汉族',
				native_place: profile.native_place || '',
				marital_status: profile.marital_status || 'single',
				political: profile.political || 'masses',
				education: profile.education || ''
			}
			// 保存原始身份证号
			originalIdCardNo.value = profile.id_card_no || ''
		}
	} catch (error) {
		console.error('加载个人档案失败', error)
	}
}

// 选择头像（同时赋值给avatar和photo）
const handleChooseAvatar = () => {
	if (uploading.value.avatar) return
	
	uni.chooseImage({
		count: 1,
		sizeType: ['compressed'],
		sourceType: ['album', 'camera'],
		success: async (res) => {
			const tempFilePath = res.tempFilePaths[0]
			
			try {
				uploading.value.avatar = true
				
				// 先本地预览
				formData.value.avatar = tempFilePath
				formData.value.photo = tempFilePath
				
				// 上传到服务器（使用 auth.js 的 uploadFile）
				const uploadResult = await uploadFile(tempFilePath, 'avatar')
				const photoUrl = uploadResult.url || uploadResult.file_url
				
				// 同时更新avatar和photo两个字段
				await uploadPhoto({ type: 'avatar', url: photoUrl })
				await uploadPhoto({ type: 'photo', url: photoUrl })
				
				// 更新为服务器URL
				formData.value.avatar = photoUrl
				formData.value.photo = photoUrl
				
				uni.showToast({
					title: '头像已更换',
					icon: 'success'
				})
			} catch (error) {
				console.error('头像上传失败', error)
				uni.showToast({
					title: error.message || '头像上传失败',
					icon: 'none'
				})
			} finally {
				uploading.value.avatar = false
			}
		}
	})
}

// 显示出生日期选择器
const showBirthdayPicker = () => {
	if (formData.value.birthday) {
		selectedDate.value = formData.value.birthday * 1000
	}
	showDatePicker.value = true
}

// 确认出生日期
const confirmBirthday = (e) => {
	formData.value.birthday = Math.floor(e.value / 1000)
	showDatePicker.value = false
}

// 显示民族选择器
const showNationPicker = () => {
	pickerType.value = 'nation'
	pickerColumns.value = [['汉族', '回族', '维吾尔族', '壮族', '满族', '苗族', '彝族', '土家族', '藏族', '蒙古族', '其他']]
	showPicker.value = true
}

// 显示婚姻状况选择器
const showMaritalStatusPicker = () => {
	pickerType.value = 'marital_status'
	pickerColumns.value = [[{text: '未婚', value: 'single'}, {text: '已婚', value: 'married'}, {text: '离异', value: 'divorced'}]]
	showPicker.value = true
}

// 显示政治面貌选择器
const showPoliticalPicker = () => {
	pickerType.value = 'political'
	pickerColumns.value = [[{text: '群众', value: 'masses'}, {text: '团员', value: 'league'}, {text: '党员', value: 'party'}]]
	showPicker.value = true
}

// 显示学历选择器
const showEducationPicker = () => {
	pickerType.value = 'education'
	pickerColumns.value = [[
		{text: '小学', value: 'primary'},
		{text: '初中', value: 'middle'},
		{text: '高中', value: 'high'},
		{text: '大专', value: 'college'},
		{text: '本科', value: 'bachelor'},
		{text: '硕士', value: 'master'},
		{text: '博士', value: 'doctor'}
	]]
	showPicker.value = true
}

// 确认选择器
const confirmPicker = (e) => {
	const value = e.value[0]
	if (pickerType.value === 'nation') {
		formData.value.nation = value
	} else if (pickerType.value === 'marital_status') {
		formData.value.marital_status = value.value || value
	} else if (pickerType.value === 'political') {
		formData.value.political = value.value || value
	} else if (pickerType.value === 'education') {
		formData.value.education = value.value || value
	}
	showPicker.value = false
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
		await updateBasicInfo({
			name: formData.value.name,
			gender: formData.value.gender,
			id_card_no: formData.value.id_card_no,
			birthday: formData.value.birthday,
			nation: formData.value.nation,
			native_place: formData.value.native_place,
			marital_status: formData.value.marital_status,
			political: formData.value.political,
			education: formData.value.education
		})

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

// 格式化日期
const formatDate = (timestamp) => {
	if (!timestamp || timestamp === 0) return ''
	const date = new Date(timestamp * 1000)
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// 文本转换函数
const getMaritalStatusText = (status) => {
	const map = {single: '未婚', married: '已婚', divorced: '离异'}
	return map[status] || '请选择'
}

const getPoliticalText = (political) => {
	const map = {party: '党员', league: '团员', masses: '群众'}
	return map[political] || '请选择'
}

const getEducationText = (education) => {
	const map = {
		primary: '小学', middle: '初中', high: '高中',
		college: '大专', bachelor: '本科', master: '硕士', doctor: '博士'
	}
	return map[education] || '请选择'
}
</script>

<style scoped lang="scss">
.edit-profile-container {
	min-height: 100vh;
	background: #f5f7fa;
	padding-bottom: 120rpx;
}

// 头像和证件照区域
.photos-section {
	display: flex;
	justify-content: center;
	padding: 60rpx 0 40rpx;
}

.photo-item {
	position: relative;
	width: 160rpx;
	height: 160rpx;
	border-radius: 16rpx;
	overflow: hidden;
	box-shadow: 0 8rpx 24rpx rgba(94, 163, 242, 0.2);
	cursor: pointer;
	
	.photo-img {
		width: 100%;
		height: 100%;
	}
	
	.photo-tip {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 60rpx;
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
	gap: 8rpx;

	text {
		font-size: 30rpx;
		color: #333;

		&.placeholder {
			color: #ccc;
		}
	}
}

.form-tip {
	padding: 16rpx 32rpx 0;
	
	text {
		font-size: 24rpx;
		color: #ff9800;
	}
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
