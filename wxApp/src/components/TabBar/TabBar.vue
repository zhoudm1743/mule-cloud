<template>
	<u-tabbar
		:value="currentIndex"
		@change="handleChange"
		:placeholder="true"
		:safeAreaInsetBottom="true"
		activeColor="#5EA3F2"
		inactiveColor="#999999"
		:border="true"
	>
		<u-tabbar-item text="首页">
			<template #active-icon>
				<u-icon name="home-fill" :size="40"></u-icon>
			</template>
			<template #inactive-icon>
				<u-icon name="home" :size="40"></u-icon>
			</template>
		</u-tabbar-item>
		
		<u-tabbar-item text="扫码">
			<template #active-icon>
				<u-icon name="scan" :size="42" color="#5EA3F2"></u-icon>
			</template>
			<template #inactive-icon>
				<u-icon name="scan" :size="42"></u-icon>
			</template>
		</u-tabbar-item>
		
		<u-tabbar-item text="我的">
			<template #active-icon>
				<u-icon name="account-fill" :size="40"></u-icon>
			</template>
			<template #inactive-icon>
				<u-icon name="account" :size="40"></u-icon>
			</template>
		</u-tabbar-item>
	</u-tabbar>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
	current: {
		type: Number,
		default: 0
	}
})

const currentIndex = ref(props.current)

const tabList = [
	'/pages/index/index',
	'/pages/scan/scan',
	'/pages/mine/mine'
]

watch(() => props.current, (newVal) => {
	currentIndex.value = newVal
})

const handleChange = (index) => {
	if (index === currentIndex.value) return
	
	// Use redirectTo instead of switchTab (no system tabBar configured)
	uni.redirectTo({
		url: tabList[index]
	})
}
</script>

<style scoped lang="scss">
// Styles handled by u-tabbar component internally
</style>

