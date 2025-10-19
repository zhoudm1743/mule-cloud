<template>
	<view class="skills-container">
		<!-- 加载中 -->
		<view v-if="loading" class="loading-section">
			<text class="loading-text">加载中...</text>
		</view>

		<!-- 技能证书内容 -->
		<view v-else class="skills-content">
			<!-- 提示信息 -->
			<view class="tip-card">
				<u-icon name="info-circle" size="36" color="#5EA3F2"></u-icon>
				<text class="tip-text">技能和证书由管理员维护，如有疑问请联系人事部门</text>
			</view>

			<!-- 技能列表 -->
			<view class="section-card">
				<view class="section-header">
					<text class="section-title">⭐ 技能列表</text>
				</view>

				<view v-if="profile.skills && profile.skills.length > 0" class="skills-list">
					<view v-for="(skill, index) in profile.skills" :key="index" class="skill-card">
						<view class="skill-header">
							<view class="skill-name-group">
								<text class="skill-name">{{ skill.name }}</text>
								<view class="skill-level-badge" :class="getSkillLevelClass(skill.level)">
									{{ getSkillLevelText(skill.level) }}
								</view>
							</view>
							<u-icon name="checkmark-circle-fill" size="40" color="#00b578"></u-icon>
						</view>

						<view v-if="skill.obtained_at" class="skill-info">
							<text class="info-label">获得时间：</text>
							<text class="info-value">{{ formatDate(skill.obtained_at) }}</text>
						</view>

						<view v-if="skill.process_ids && skill.process_ids.length > 0" class="skill-info">
							<text class="info-label">可操作工序：</text>
							<text class="info-value">{{ skill.process_ids.length }}个</text>
						</view>

						<view v-if="skill.remark" class="skill-remark">
							<text>{{ skill.remark }}</text>
						</view>
					</view>
				</view>

				<view v-else class="empty-section">
					<u-icon name="file-text" size="80" color="#e0e0e0"></u-icon>
					<text class="empty-text">暂无技能记录</text>
					<text class="empty-tip">完成技能培训后，管理员会为您添加技能</text>
				</view>
			</view>

			<!-- 证书列表 -->
			<view class="section-card">
				<view class="section-header">
					<text class="section-title">📄 证书列表</text>
				</view>

				<view v-if="profile.certificates && profile.certificates.length > 0" class="certificates-list">
					<view v-for="(cert, index) in profile.certificates" :key="index" class="cert-card">
						<view class="cert-header">
							<u-icon name="medal" size="40" color="#ff9800"></u-icon>
							<text class="cert-name">{{ cert.name }}</text>
						</view>

						<view class="cert-info-list">
							<view class="cert-info-item">
								<text class="label">证书编号</text>
								<text class="value">{{ cert.no || '未填写' }}</text>
							</view>
							<view class="cert-info-item">
								<text class="label">发证机关</text>
								<text class="value">{{ cert.issue_org || '未填写' }}</text>
							</view>
							<view class="cert-info-item">
								<text class="label">发证日期</text>
								<text class="value">{{ formatDate(cert.issued_at) }}</text>
							</view>
							<view class="cert-info-item">
								<text class="label">有效期</text>
								<text class="value" :class="{'expired': isExpired(cert.expired_at)}">
									{{ getExpiredText(cert.expired_at) }}
								</text>
							</view>
						</view>

						<view v-if="cert.file_url" class="cert-actions">
							<view class="view-cert-btn" @click="viewCertFile(cert.file_url)">
								<u-icon name="eye" size="28" color="#5EA3F2"></u-icon>
								<text>查看证书</text>
							</view>
						</view>
					</view>
				</view>

				<view v-else class="empty-section">
					<u-icon name="file-text" size="80" color="#e0e0e0"></u-icon>
					<text class="empty-text">暂无证书记录</text>
					<text class="empty-tip">获得职业资格证书后，管理员会为您录入</text>
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
		console.error('加载技能证书失败', error)
		uni.showToast({
			title: '加载失败',
			icon: 'none'
		})
	} finally {
		loading.value = false
	}
}

const formatDate = (timestamp) => {
	if (!timestamp || timestamp === 0) return '未填写'
	const date = new Date(timestamp * 1000)
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const getSkillLevelText = (level) => {
	const map = {
		beginner: '初级',
		intermediate: '中级',
		advanced: '高级',
		expert: '专家'
	}
	return map[level] || level
}

const getSkillLevelClass = (level) => {
	return level
}

const getExpiredText = (expiredAt) => {
	if (!expiredAt || expiredAt === 0) return '长期有效'
	const date = new Date(expiredAt * 1000)
	const now = new Date()
	if (date < now) return '已过期'
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const isExpired = (expiredAt) => {
	if (!expiredAt || expiredAt === 0) return false
	const date = new Date(expiredAt * 1000)
	return date < new Date()
}

const viewCertFile = (url) => {
	// 预览证书文件
	uni.downloadFile({
		url: url,
		success: (res) => {
			if (res.statusCode === 200) {
				uni.openDocument({
					filePath: res.tempFilePath,
					showMenu: true,
					success: () => {
						console.log('打开文档成功')
					}
				})
			}
		}
	})
}
</script>

<style scoped lang="scss">
.skills-container {
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

.skills-content {
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

.section-header {
	margin-bottom: 24rpx;
}

.section-title {
	font-size: 32rpx;
	color: #333;
	font-weight: bold;
}

// 技能列表
.skills-list {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.skill-card {
	background: #f8f9fa;
	border-radius: 16rpx;
	padding: 24rpx;
}

.skill-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 16rpx;
}

.skill-name-group {
	display: flex;
	align-items: center;
	gap: 12rpx;
}

.skill-name {
	font-size: 30rpx;
	color: #333;
	font-weight: bold;
}

.skill-level-badge {
	padding: 4rpx 16rpx;
	border-radius: 24rpx;
	font-size: 24rpx;
	font-weight: 500;

	&.beginner {
		background: #e3f2fd;
		color: #2196f3;
	}

	&.intermediate {
		background: #e8f5e9;
		color: #4caf50;
	}

	&.advanced {
		background: #fff3e0;
		color: #ff9800;
	}

	&.expert {
		background: #fce4ec;
		color: #e91e63;
	}
}

.skill-info {
	font-size: 26rpx;
	color: #666;
	margin-bottom: 8rpx;

	.info-label {
		color: #999;
	}

	.info-value {
		color: #333;
	}
}

.skill-remark {
	margin-top: 16rpx;
	padding: 16rpx;
	background: #fff;
	border-radius: 12rpx;

	text {
		font-size: 26rpx;
		color: #666;
		line-height: 1.6;
	}
}

// 证书列表
.certificates-list {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.cert-card {
	background: #f8f9fa;
	border-radius: 16rpx;
	padding: 24rpx;
}

.cert-header {
	display: flex;
	align-items: center;
	gap: 16rpx;
	margin-bottom: 20rpx;

	.cert-name {
		font-size: 30rpx;
		color: #333;
		font-weight: bold;
	}
}

.cert-info-list {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
}

.cert-info-item {
	display: flex;
	justify-content: space-between;
	font-size: 26rpx;

	.label {
		color: #999;
	}

	.value {
		color: #333;

		&.expired {
			color: #f44336;
		}
	}
}

.cert-actions {
	margin-top: 20rpx;
	padding-top: 20rpx;
	border-top: 1rpx solid #e0e0e0;
}

.view-cert-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 8rpx;
	padding: 16rpx;
	background: #e8f4ff;
	border-radius: 12rpx;

	text {
		font-size: 28rpx;
		color: #5EA3F2;
		font-weight: 500;
	}

	&:active {
		opacity: 0.8;
	}
}

// 空状态
.empty-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 80rpx 40rpx;

	.empty-text {
		font-size: 28rpx;
		color: #999;
		margin-top: 24rpx;
	}

	.empty-tip {
		font-size: 24rpx;
		color: #ccc;
		margin-top: 12rpx;
		text-align: center;
	}
}
</style>

