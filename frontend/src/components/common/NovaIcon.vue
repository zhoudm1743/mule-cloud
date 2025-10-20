<script setup lang="ts">
import { Icon } from '@iconify/vue'

interface iconPorps {
  /* 图标名称 */
  icon?: string
  /* 图标颜色 */
  color?: string
  /* 图标大小 */
  size?: number
  /* 图标深度 */
  depth?: 1 | 2 | 3 | 4 | 5
}
const { size = 18, icon } = defineProps<iconPorps>()

// 在组件外部缓存 SVG 导入结果，避免每次都重新创建
const svgModules = import.meta.glob<string>('@/assets/svg-icons/*.svg', {
  query: '?raw',
  import: 'default',
  eager: true,
})

const isLocal = computed(() => {
  return icon && icon.startsWith('local:')
})

// 使用 computed 缓存本地图标的获取结果
const localIconContent = computed(() => {
  if (!icon || !isLocal.value) return ''
  const svgName = icon.replace('local:', '')
  return svgModules[`/src/assets/svg-icons/${svgName}.svg`] || ''
})
</script>

<template>
  <n-icon
    v-if="icon"
    :size="size"
    :depth="depth"
    :color="color"
  >
    <template v-if="isLocal">
      <i v-html="localIconContent" />
    </template>
    <template v-else>
      <Icon :icon="icon" />
    </template>
  </n-icon>
</template>
