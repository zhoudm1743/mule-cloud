<template>
	<view>？</view>
</template>

<script setup>
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'

onLaunch(() => {
	console.log('App Launch')
	
	// 检查更新（仅微信小程序）
	// #ifdef MP-WEIXIN
	const updateManager = uni.getUpdateManager()
	updateManager.onCheckForUpdate((res) => {
		console.log('是否有新版本:', res.hasUpdate)
	})
	updateManager.onUpdateReady(() => {
		uni.showModal({
			title: '更新提示',
			content: '新版本已经准备好，是否重启应用？',
			success: (res) => {
				if (res.confirm) {
					updateManager.applyUpdate()
				}
			}
		})
	})
	updateManager.onUpdateFailed(() => {
		uni.showModal({
			title: '更新失败',
			content: '新版本下载失败，请删除小程序后重新搜索打开',
			showCancel: false
		})
	})
	// #endif
})

onShow(() => {
	console.log('App Show')
})

onHide(() => {
	console.log('App Hide')
})
</script>

<style lang="scss">
@import '@/uni.scss';

/* 全局样式 */
page {
	background-color: #f5f5f5;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif;
}

/* 通用按钮样式 */
button {
	&::after {
		border: none;
	}
}

/* 重置 uv-ui 样式 */
.uv-popup__content {
	background-color: transparent !important;
}
</style>
