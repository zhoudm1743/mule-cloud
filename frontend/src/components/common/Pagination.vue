<script setup lang="ts">
interface Props {
  count?: number
  page?: number
  pageSize?: number
  align?: 'left' | 'center' | 'right'
}
const props = withDefaults(defineProps<Props>(), {
  count: 0,
  page: 1,
  pageSize: 10,
  align: 'right',
})

const emit = defineEmits<{
  change: [page: number, pageSize: number] // 具名元组语法
}>()

const currentPage = ref(props.page)
const currentPageSize = ref(props.pageSize)
const displayOrder: Array<'pages' | 'size-picker' | 'quick-jumper'> = ['size-picker', 'pages']

// 监听 props 变化，同步内部状态
watch(() => props.page, (newPage) => {
  currentPage.value = newPage
})

watch(() => props.pageSize, (newPageSize) => {
  currentPageSize.value = newPageSize
})

function changePage() {
  emit('change', currentPage.value, currentPageSize.value)
}
</script>

<template>
  <div v-if="count > 0" :class="['pagination-wrapper', `pagination-${align}`]">
    <n-pagination
      v-model:page="currentPage"
      v-model:page-size="currentPageSize"
      :page-sizes="[10, 20, 30, 50]"
      :item-count="count"
      :display-order="displayOrder"
      show-size-picker
      @update-page="changePage"
      @update-page-size="changePage"
    />
  </div>
</template>

<style scoped>
.pagination-wrapper {
  padding: 10px 0;
  display: flex;
  width: 100%;
}

.pagination-left {
  justify-content: flex-start;
}

.pagination-center {
  justify-content: center;
}

.pagination-right {
  justify-content: flex-end;
}
</style>
