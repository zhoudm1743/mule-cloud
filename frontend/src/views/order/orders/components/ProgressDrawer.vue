<script setup lang="ts">
import { computed, ref } from 'vue'
import { NDrawer, NDrawerContent, NSpace, NCard, NGrid, NGridItem, NStatistic, NTag, NText, NSpin } from 'naive-ui'
import { fetchOrderProgress } from '@/service/api/production'
import WorkflowStatus from './WorkflowStatus.vue'

const show = ref(false)
const loading = ref(false)
const orderInfo = ref<Api.Order.OrderInfo | null>(null)
const progressData = ref<Api.Production.OrderProcedureProgress[]>([])

// 打开抽屉
async function open(order: Api.Order.OrderInfo) {
  show.value = true
  orderInfo.value = order
  await fetchProgress(order.id)
}

// 获取进度数据
async function fetchProgress(orderId: string) {
  loading.value = true
  try {
    const res = await fetchOrderProgress(orderId)
    if (res.data) {
      progressData.value = res.data.procedures || []
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取订单进度失败')
  }
  finally {
    loading.value = false
  }
}

// 计算统计数据（按件数统计，而不是按工序数）
const progressStats = computed(() => {
  if (!progressData.value.length || !orderInfo.value) {
    return { total: 0, completed: 0, inProgress: 0, pending: 0 }
  }
  
  const totalQty = orderInfo.value.quantity || 0
  
  // 已完成：最后一道工序完成的件数
  const lastProcedure = progressData.value[progressData.value.length - 1]
  const completed = lastProcedure?.reported_qty || 0
  
  // 未开始：第一道工序还未开始的件数
  const firstProcedure = progressData.value[0]
  const pending = totalQty - (firstProcedure?.reported_qty || 0)
  
  // 进行中：已开始但未完成的件数
  const inProgress = totalQty - completed - pending
  
  return {
    total: totalQty,
    completed,
    inProgress,
    pending,
  }
})

// 根据进度获取颜色类名
function getProgressClass(progress: number) {
  if (progress >= 100) return 'progress-complete'
  if (progress >= 75) return 'progress-high'
  if (progress >= 50) return 'progress-medium'
  if (progress >= 25) return 'progress-low'
  if (progress > 0) return 'progress-very-low'
  return 'progress-none'
}

defineExpose({
  open,
})
</script>

<template>
  <NDrawer v-model:show="show" :width="800" placement="right">
    <NDrawerContent title="订单工序进度" closable>
      <NSpin :show="loading">
        <NSpace vertical size="large">
          <!-- 工作流状态 -->
          <WorkflowStatus v-if="orderInfo" :order-id="orderInfo.id" />

          <!-- 订单基本信息 -->
          <NCard v-if="orderInfo" title="订单信息" size="small" :bordered="false">
            <NGrid cols="2" x-gap="12" y-gap="12">
              <NGridItem>
                <div class="info-item">
                  <div class="label">合同号</div>
                  <div class="value">{{ orderInfo.contract_no }}</div>
                </div>
              </NGridItem>
              <NGridItem>
                <div class="info-item">
                  <div class="label">款号/款名</div>
                  <div class="value">{{ orderInfo.style_no }} / {{ orderInfo.style_name }}</div>
                </div>
              </NGridItem>
              <NGridItem>
                <div class="info-item">
                  <div class="label">客户</div>
                  <div class="value">{{ orderInfo.customer_name }}</div>
                </div>
              </NGridItem>
              <NGridItem>
                <div class="info-item">
                  <div class="label">订单数量</div>
                  <div class="value">{{ orderInfo.quantity }} 件</div>
                </div>
              </NGridItem>
            </NGrid>
          </NCard>

          <!-- 进度统计（按件数统计）-->
          <NGrid cols="4" x-gap="12">
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="订单数量" :value="progressStats.total">
                  <template #suffix>
                    <span style="font-size: 14px; margin-left: 4px;">件</span>
                  </template>
                </NStatistic>
              </NCard>
            </NGridItem>
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="已完成" :value="progressStats.completed">
                  <template #suffix>
                    <span style="font-size: 14px; margin: 0 4px;">件</span>
                    <NTag v-if="progressStats.total > 0" type="success" size="small">
                      {{ Math.round((progressStats.completed / progressStats.total) * 100) }}%
                    </NTag>
                  </template>
                </NStatistic>
              </NCard>
            </NGridItem>
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="进行中" :value="progressStats.inProgress">
                  <template #suffix>
                    <span style="font-size: 14px; margin: 0 4px;">件</span>
                    <NTag v-if="progressStats.total > 0" type="warning" size="small">
                      {{ Math.round((progressStats.inProgress / progressStats.total) * 100) }}%
                    </NTag>
                  </template>
                </NStatistic>
              </NCard>
            </NGridItem>
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="未开始" :value="progressStats.pending">
                  <template #suffix>
                    <span style="font-size: 14px; margin: 0 4px;">件</span>
                    <NTag v-if="progressStats.total > 0" type="default" size="small">
                      {{ Math.round((progressStats.pending / progressStats.total) * 100) }}%
                    </NTag>
                  </template>
                </NStatistic>
              </NCard>
            </NGridItem>
          </NGrid>

          <!-- 工序进度热力图 -->
          <NCard title="工序进度详情" size="small" :bordered="false">
            <div v-if="progressData.length === 0" class="text-center text-gray-400 py-8">
              暂无工序数据
            </div>
            <div v-else class="procedure-heatmap">
              <div
                v-for="procedure in progressData"
                :key="procedure.id"
                class="heatmap-row"
              >
                <div class="procedure-label">
                  <NText strong>{{ procedure.procedure_seq }}. {{ procedure.procedure_name }}</NText>
                </div>
                <div class="progress-bar-container">
                  <div
                    class="progress-bar"
                    :class="getProgressClass(procedure.progress)"
                    :style="{ width: `${Math.min(procedure.progress, 100)}%` }"
                  >
                    <span v-if="procedure.progress > 10" class="progress-text">
                      {{ Math.round(procedure.progress) }}%
                    </span>
                  </div>
                  <div class="progress-info">
                    <NText depth="3">
                      {{ procedure.reported_qty }} / {{ procedure.total_qty }} 件
                    </NText>
                    <NTag
                      v-if="procedure.is_completed"
                      type="success"
                      size="small"
                      :bordered="false"
                    >
                      已完成
                    </NTag>
                    <NTag
                      v-else-if="procedure.reported_qty > 0"
                      type="warning"
                      size="small"
                      :bordered="false"
                    >
                      进行中
                    </NTag>
                    <NTag
                      v-else
                      type="default"
                      size="small"
                      :bordered="false"
                    >
                      未开始
                    </NTag>
                  </div>
                </div>
              </div>
            </div>
          </NCard>
        </NSpace>
      </NSpin>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped lang="scss">
.info-item {
  .label {
    font-size: 12px;
    color: #999;
    margin-bottom: 4px;
  }
  .value {
    font-size: 14px;
    font-weight: 500;
    color: #333;
  }
}

.procedure-heatmap {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.heatmap-row {
  display: flex;
  align-items: center;
  gap: 16px;
}

.procedure-label {
  min-width: 100px;
  flex-shrink: 0;
}

.progress-bar-container {
  flex: 1;
  position: relative;
  height: 48px;
  background-color: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  align-items: center;
  padding: 0 12px;
}

.progress-bar {
  position: absolute;
  left: 0;
  top: 0;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding-right: 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
  
  .progress-text {
    color: white;
    font-weight: 600;
    font-size: 14px;
    text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
  }
}

.progress-info {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

// 热力图颜色方案（从冷色到暖色）
.progress-none {
  background: linear-gradient(90deg, #e5e7eb 0%, #d1d5db 100%);
}

.progress-very-low {
  background: linear-gradient(90deg, #dbeafe 0%, #bfdbfe 100%);
}

.progress-low {
  background: linear-gradient(90deg, #93c5fd 0%, #60a5fa 100%);
}

.progress-medium {
  background: linear-gradient(90deg, #fde047 0%, #facc15 100%);
}

.progress-high {
  background: linear-gradient(90deg, #fdba74 0%, #fb923c 100%);
}

.progress-complete {
  background: linear-gradient(90deg, #86efac 0%, #4ade80 100%);
}
</style>

