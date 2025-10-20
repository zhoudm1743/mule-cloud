<script setup lang="ts">
import { ref, computed } from 'vue'
import { NCard, NSpace, NStatistic, NGrid, NGridItem, NButton, NSpin, NTag, NEmpty } from 'naive-ui'
import { fetchOperationLogStats } from '@/service'
import { useBoolean } from '@/hooks'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

// 统计数据
const statsData = ref<Api.OperationLog.StatsResponse | null>(null)

// 时间范围（默认最近7天）
const timeRange = ref<[number, number]>([
  Date.now() - 7 * 24 * 60 * 60 * 1000,
  Date.now(),
])

// 计算成功率
const successRate = computed(() => {
  if (!statsData.value || statsData.value.total === 0) return 0
  return ((statsData.value.success_num / statsData.value.total) * 100).toFixed(2)
})

// 计算失败率
const failRate = computed(() => {
  if (!statsData.value || statsData.value.total === 0) return 0
  return ((statsData.value.fail_num / statsData.value.total) * 100).toFixed(2)
})

// 格式化平均耗时
const avgTimeFormatted = computed(() => {
  if (!statsData.value) return '0ms'
  const avgTime = statsData.value.avg_time
  if (avgTime < 1000) return `${avgTime.toFixed(0)}ms`
  return `${(avgTime / 1000).toFixed(2)}s`
})

// 时间范围快捷选择
const timeRangeShortcuts = {
  '今天': () => {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    return [today.getTime(), Date.now()] as [number, number]
  },
  '最近3天': () => {
    return [Date.now() - 3 * 24 * 60 * 60 * 1000, Date.now()] as [number, number]
  },
  '最近7天': () => {
    return [Date.now() - 7 * 24 * 60 * 60 * 1000, Date.now()] as [number, number]
  },
  '最近30天': () => {
    return [Date.now() - 30 * 24 * 60 * 60 * 1000, Date.now()] as [number, number]
  },
}

// 获取统计数据
async function fetchStats() {
  startLoading()
  try {
    const params: Api.OperationLog.StatsRequest = {
      start_time: Math.floor(timeRange.value[0] / 1000),
      end_time: Math.floor(timeRange.value[1] / 1000),
    }
    
    const res = await fetchOperationLogStats(params)
    if (res.data) {
      statsData.value = res.data
    }
  }
  catch (e) {
    console.error('[Fetch Stats Error]:', e)
    window.$message.error('获取统计数据失败')
  }
  finally {
    endLoading()
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <NSpace vertical size="large">
    <!-- 筛选条件 -->
    <NCard title="统计筛选" :bordered="false" class="rounded-8px shadow-sm">
      <NSpace>
        <n-date-picker
          v-model:value="timeRange"
          type="datetimerange"
          clearable
          class="w-400px"
          :shortcuts="timeRangeShortcuts"
        />
        <NButton type="primary" @click="fetchStats">
          <template #icon>
            <nova-icon icon="carbon:analytics" :size="18" />
          </template>
          查询统计
        </NButton>
      </NSpace>
    </NCard>

    <NSpin :show="loading">
      <template v-if="statsData">
        <!-- 总览统计 -->
        <NCard title="总览" :bordered="false" class="rounded-8px shadow-sm">
          <NGrid cols="2 600:4" responsive="screen" :x-gap="24" :y-gap="16">
            <NGridItem>
              <NStatistic label="总操作数" :value="statsData.total">
                <template #prefix>
                  <nova-icon icon="carbon:data-1" :size="20" class="text-primary" />
                </template>
              </NStatistic>
            </NGridItem>
            <NGridItem>
              <NStatistic label="成功操作" :value="statsData.success_num">
                <template #prefix>
                  <nova-icon icon="carbon:checkmark-filled" :size="20" class="text-success" />
                </template>
                <template #suffix>
                  <NTag type="success" size="small" class="ml-2">
                    {{ successRate }}%
                  </NTag>
                </template>
              </NStatistic>
            </NGridItem>
            <NGridItem>
              <NStatistic label="失败操作" :value="statsData.fail_num">
                <template #prefix>
                  <nova-icon icon="carbon:warning-filled" :size="20" class="text-error" />
                </template>
                <template #suffix>
                  <NTag type="error" size="small" class="ml-2">
                    {{ failRate }}%
                  </NTag>
                </template>
              </NStatistic>
            </NGridItem>
            <NGridItem>
              <NStatistic label="平均耗时">
                <template #prefix>
                  <nova-icon icon="carbon:time" :size="20" class="text-warning" />
                </template>
                <template #default>
                  {{ avgTimeFormatted }}
                </template>
              </NStatistic>
            </NGridItem>
          </NGrid>
        </NCard>

        <!-- TOP操作用户 -->
        <NCard title="TOP 10 操作用户" :bordered="false" class="rounded-8px shadow-sm">
          <template v-if="statsData.top_users && statsData.top_users.length > 0">
            <NSpace vertical size="small">
              <div
                v-for="(user, index) in statsData.top_users"
                :key="user.user_id"
                class="flex items-center justify-between p-3 rounded hover:bg-gray-50 dark:hover:bg-gray-800 transition"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="flex items-center justify-center w-8 h-8 rounded-full font-bold"
                    :class="{
                      'bg-yellow-100 text-yellow-700': index === 0,
                      'bg-gray-200 text-gray-700': index === 1,
                      'bg-orange-100 text-orange-700': index === 2,
                      'bg-blue-50 text-blue-600': index > 2,
                    }"
                  >
                    {{ index + 1 }}
                  </div>
                  <div>
                    <div class="font-medium">{{ user.username }}</div>
                    <div class="text-xs text-gray-500">ID: {{ user.user_id }}</div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="primary" size="small">
                    {{ user.count }} 次操作
                  </NTag>
                </div>
              </div>
            </NSpace>
          </template>
          <NEmpty v-else description="暂无数据" />
        </NCard>

        <!-- TOP操作类型 -->
        <NCard title="TOP 10 操作类型" :bordered="false" class="rounded-8px shadow-sm">
          <template v-if="statsData.top_actions && statsData.top_actions.length > 0">
            <NSpace vertical size="small">
              <div
                v-for="(action, index) in statsData.top_actions"
                :key="action.action"
                class="flex items-center justify-between p-3 rounded hover:bg-gray-50 dark:hover:bg-gray-800 transition"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="flex items-center justify-center w-8 h-8 rounded-full font-bold"
                    :class="{
                      'bg-yellow-100 text-yellow-700': index === 0,
                      'bg-gray-200 text-gray-700': index === 1,
                      'bg-orange-100 text-orange-700': index === 2,
                      'bg-blue-50 text-blue-600': index > 2,
                    }"
                  >
                    {{ index + 1 }}
                  </div>
                  <div class="font-medium">{{ action.action }}</div>
                </div>
                <div class="flex items-center gap-2">
                  <NTag type="info" size="small">
                    {{ action.count }} 次
                  </NTag>
                  <div class="w-40 bg-gray-200 dark:bg-gray-700 rounded-full h-2">
                    <div
                      class="bg-blue-500 h-2 rounded-full"
                      :style="{ width: `${(action.count / statsData!.total) * 100}%` }"
                    ></div>
                  </div>
                </div>
              </div>
            </NSpace>
          </template>
          <NEmpty v-else description="暂无数据" />
        </NCard>
      </template>

      <template v-else>
        <NCard :bordered="false" class="rounded-8px shadow-sm">
          <NEmpty description="请选择时间范围并查询统计数据" />
        </NCard>
      </template>
    </NSpin>
  </NSpace>
</template>

<style scoped>
:deep(.n-statistic) {
  .n-statistic-value {
    font-size: 28px;
    font-weight: 600;
  }
}
</style>

