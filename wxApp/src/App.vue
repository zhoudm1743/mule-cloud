<template>
	<view></view>
</template>

<script setup>
import { onLaunch, onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/store/modules/user'

onLaunch(() => {
	console.log('App Launch')
	
	// 检查登录状态
	checkLoginStatus()
})

onShow(() => {
	console.log('App Show')
})

// 检查登录状态
const checkLoginStatus = () => {
	const userStore = useUserStore()
	const token = uni.getStorageSync('token')
	
	console.log('检查登录状态, token:', token ? '存在' : '不存在')
	console.log('isLoggedIn:', userStore.isLoggedIn)
	console.log('hasTenant:', userStore.hasTenant)
	
	// 获取当前页面
	const pages = getCurrentPages()
	const currentPage = pages[pages.length - 1]
	const currentRoute = currentPage ? currentPage.route : ''
	
	// 如果不是登录相关页面，且未登录，跳转到登录页
	const authPages = ['pages/login/login', 'pages/bind-tenant/bind-tenant', 'pages/select-tenant/select-tenant']
	const isAuthPage = authPages.some(page => currentRoute.includes(page))
	
	if (!userStore.isLoggedIn && !isAuthPage) {
		console.log('未登录，跳转到登录页')
		uni.reLaunch({
			url: '/pages/login/login'
		})
	}
}
</script>

<style lang="scss">
@import 'uview-plus/index.scss';

page {
	background-color: #F0F8FF;
}
</style>
