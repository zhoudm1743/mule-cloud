<template>
	<view class="profile-container">
		<!-- 加载中 -->
		<view v-if="loading" class="loading-section">
			<text class="loading-text">加载中...</text>
		</view>

		<!-- 个人档案内容 -->
		<view v-else class="profile-content">
		<!-- 头部信息 -->
		<view class="header-section">
			<image :src="profile.avatar || '/static/logo.png'" mode="aspectFill" class="avatar"></image>
			<view class="header-info">
				<text class="name">{{ profile.name || '未设置' }}</text>
				<text class="position">{{ profile.department }} · {{ profile.position }}</text>
				<text class="job-number">工号：{{ profile.job_number || '未设置' }}</text>
			</view>
		</view>

			<!-- 基本信息 -->
			<view class="section-card" @click="navigateToEdit('basic')">
				<view class="section-header">
					<text class="section-title">基本信息</text>
					<u-icon name="arrow-right" size="20" color="#999"></u-icon>
				</view>
				<view class="info-grid">
					<view class="info-item">
						<text class="info-label">姓名</text>
						<text class="info-value">{{ profile.name || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">性别</text>
						<text class="info-value">{{ getGenderText(profile.gender) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">身份证号</text>
						<text class="info-value">{{ profile.id_card_no || '未填写' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">出生日期</text>
						<text class="info-value">{{ formatDate(profile.birthday) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">年龄</text>
						<text class="info-value">{{ profile.age }}岁</text>
					</view>
					<view class="info-item">
						<text class="info-label">民族</text>
						<text class="info-value">{{ profile.nation || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">籍贯</text>
						<text class="info-value">{{ profile.native_place || '未设置' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">婚姻状况</text>
						<text class="info-value">{{ getMaritalStatusText(profile.marital_status) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">政治面貌</text>
						<text class="info-value">{{ getPoliticalText(profile.political) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">学历</text>
						<text class="info-value">{{ getEducationText(profile.education) }}</text>
					</view>
				</view>
			</view>

			<!-- 联系信息 -->
			<view class="section-card" @click="navigateToEdit('contact')">
				<view class="section-header">
					<text class="section-title">联系信息</text>
					<u-icon name="arrow-right" size="20" color="#999"></u-icon>
				</view>
				<view class="info-list">
					<view class="info-row">
						<text class="info-label">手机号</text>
						<text class="info-value">{{ profile.phone || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">邮箱</text>
						<text class="info-value">{{ profile.email || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">家庭住址</text>
						<text class="info-value">{{ profile.address || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">紧急联系人</text>
						<text class="info-value">{{ profile.emergency_contact || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">紧急电话</text>
						<text class="info-value">{{ profile.emergency_phone || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">关系</text>
						<text class="info-value">{{ profile.emergency_relation || '未设置' }}</text>
					</view>
				</view>
			</view>

			<!-- 工作信息 -->
			<view class="section-card" @click="navigateToWork">
				<view class="section-header">
					<text class="section-title">工作信息</text>
					<u-icon name="arrow-right" size="20" color="#999"></u-icon>
				</view>
				<view class="info-list">
					<view class="info-row">
						<text class="info-label">工号</text>
						<text class="info-value">{{ profile.job_number || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">部门</text>
						<text class="info-value">{{ profile.department || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">岗位</text>
						<text class="info-value">{{ profile.position || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">车间</text>
						<text class="info-value">{{ profile.workshop || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">班组</text>
						<text class="info-value">{{ profile.team || '未设置' }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">入职日期</text>
						<text class="info-value">{{ formatDate(profile.employed_at) }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">工龄</text>
						<text class="info-value">{{ profile.work_years }}年{{ profile.work_months }}个月</text>
					</view>
				</view>
			</view>

			<!-- 技能证书 -->
			<view class="section-card" @click="navigateToSkills">
				<view class="section-header">
					<text class="section-title">技能证书</text>
					<u-icon name="arrow-right" size="20" color="#999"></u-icon>
				</view>
				<view v-if="profile.skills && profile.skills.length > 0" class="skill-list">
					<view v-for="(skill, index) in profile.skills" :key="index" class="skill-item">
						<text class="skill-name">{{ skill.name }}</text>
						<text class="skill-level">{{ getSkillLevelText(skill.level) }}</text>
					</view>
				</view>
				<view v-else class="empty-text">
					<text>暂无技能记录</text>
				</view>
			</view>

			<!-- 薪资信息 -->
			<view class="section-card">
				<view class="section-header">
					<text class="section-title">薪资信息</text>
				</view>
				<view class="info-list">
					<view class="info-row">
						<text class="info-label">薪资类型</text>
						<text class="info-value">{{ getSalaryTypeText(profile.salary_type) }}</text>
					</view>
					<text class="tip-text">具体金额请联系人事部门</text>
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

// 加载个人档案
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
		console.error('加载个人档案失败', error)
		uni.showToast({
			title: '加载失败',
			icon: 'none'
		})
	} finally {
		loading.value = false
	}
}

// 导航到编辑页面
const navigateToEdit = (type) => {
	if (type === 'basic') {
		uni.navigateTo({ url: '/pages/edit-profile/edit-profile' })
	} else if (type === 'contact') {
		uni.navigateTo({ url: '/pages/member/edit-contact/edit-contact' })
	}
}

// 导航到工作信息
const navigateToWork = () => {
	uni.navigateTo({ url: '/pages/member/work-info/work-info' })
}

// 导航到技能证书
const navigateToSkills = () => {
	uni.navigateTo({ url: '/pages/member/skills/skills' })
}

// 格式化日期
const formatDate = (timestamp) => {
	if (!timestamp || timestamp === 0) return '未设置'
	const date = new Date(timestamp * 1000)
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// 性别文本
const getGenderText = (gender) => {
	const map = { 0: '未知', 1: '男', 2: '女' }
	return map[gender] || '未知'
}

// 婚姻状况文本
const getMaritalStatusText = (status) => {
	const map = {
		single: '未婚',
		married: '已婚',
		divorced: '离异'
	}
	return map[status] || '未设置'
}

// 政治面貌文本
const getPoliticalText = (political) => {
	const map = {
		party: '党员',
		league: '团员',
		masses: '群众'
	}
	return map[political] || '未设置'
}

// 学历文本
const getEducationText = (education) => {
	const map = {
		primary: '小学',
		middle: '初中',
		high: '高中',
		college: '大专',
		bachelor: '本科',
		master: '硕士',
		doctor: '博士'
	}
	return map[education] || '未设置'
}

// 技能等级文本
const getSkillLevelText = (level) => {
	const map = {
		beginner: '初级',
		intermediate: '中级',
		advanced: '高级',
		expert: '专家'
	}
	return map[level] || level
}

// 薪资类型文本
const getSalaryTypeText = (type) => {
	const map = {
		hourly: '计时',
		piece: '计件',
		monthly: '月薪',
		mixed: '混合'
	}
	return map[type] || '未设置'
}
</script>

<style scoped lang="scss">
.profile-container {
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

.profile-content {
	padding: 0 24rpx;
}

// 头部信息
.header-section {
	background: linear-gradient(135deg, #5EA3F2 0%, #4FC3F7 100%);
	border-radius: 24rpx;
	padding: 40rpx;
	margin: 24rpx 0;
	display: flex;
	align-items: center;
	box-shadow: 0 8rpx 24rpx rgba(94, 163, 242, 0.3);
}

.avatar {
	width: 120rpx;
	height: 120rpx;
	border-radius: 16rpx;
	border: 4rpx solid rgba(255, 255, 255, 0.3);
	margin-right: 24rpx;
}

.header-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: 8rpx;

	.name {
		font-size: 36rpx;
		color: #fff;
		font-weight: bold;
	}

	.position {
		font-size: 28rpx;
		color: rgba(255, 255, 255, 0.9);
	}

	.job-number {
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.8);
	}
}

// 卡片区域
.section-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
}

.section-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 24rpx;
}

.section-title {
	font-size: 32rpx;
	color: #333;
	font-weight: bold;
}

// 网格布局（基本信息）
.info-grid {
	display: grid;
	grid-template-columns: repeat(2, 1fr);
	gap: 24rpx;
}

.info-item {
	display: flex;
	flex-direction: column;
	gap: 8rpx;
}

.info-label {
	font-size: 24rpx;
	color: #999;
}

.info-value {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
}

// 列表布局（联系信息、工作信息）
.info-list {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.info-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

// 技能列表
.skill-list {
	display: flex;
	flex-wrap: wrap;
	gap: 16rpx;
}

.skill-item {
	display: flex;
	align-items: center;
	gap: 8rpx;
	padding: 12rpx 24rpx;
	background: #e8f4ff;
	border-radius: 32rpx;

	.skill-name {
		font-size: 28rpx;
		color: #5EA3F2;
		font-weight: 500;
	}

	.skill-level {
		font-size: 24rpx;
		color: #5EA3F2;
	}
}

.empty-text {
	text-align: center;
	padding: 40rpx 0;

	text {
		font-size: 28rpx;
		color: #999;
	}
}

.tip-text {
	font-size: 24rpx;
	color: #999;
	margin-top: 16rpx;
	font-style: italic;
}
</style>

