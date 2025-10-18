<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { NDrawer, NDrawerContent, NSpace, NCard, NGrid, NGridItem, NStatistic, NTag, NTimeline, NTimelineItem, NProgress, NText, NSpin } from 'naive-ui'
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

// 计算统计数据
const progressStats = computed(() => {
  if (!progressData.value.length) {
    return { total: 0, completed: 0, inProgress: 0, pending: 0 }
  }
  
  const completed = progressData.value.filter(p => p.is_completed).length
  const inProgress = progressData.value.filter(p => !p.is_completed && p.reported_qty > 0).length
  const pending = progressData.value.filter(p => p.reported_qty === 0).length
  
  return {
    total: progressData.value.length,
    completed,
    inProgress,
    pending,
  }
})

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

          <!-- 进度统计 -->
          <NGrid cols="4" x-gap="12">
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="总工序数" :value="progressStats.total" />
              </NCard>
            </NGridItem>
            <NGridItem>
              <NCard size="small" :bordered="false">
                <NStatistic label="已完成" :value="progressStats.completed">
                  <template #suffix>
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
                    <NTag v-if="progressStats.total > 0" type="default" size="small">
                      {{ Math.round((progressStats.pending / progressStats.total) * 100) }}%
                    </NTag>
                  </template>
                </NStatistic>
              </NCard>
            </NGridItem>
          </NGrid>

          <!-- 工序进度时间轴 -->
          <NCard title="工序进度详情" size="small" :bordered="false">
            <div v-if="progressData.length === 0" class="text-center text-gray-400 py-8">
              暂无工序数据
            </div>
            <NTimeline v-else>
              <NTimelineItem
                v-for="procedure in progressData"
                :key="procedure.id"
                :type="procedure.is_completed ? 'success' : (procedure.reported_qty > 0 ? 'warning' : 'default')"
                :title="`${procedure.procedure_seq}. ${procedure.procedure_name}`"
              >
                <template #icon>
                  <div class="timeline-icon">
                    {{ procedure.procedure_seq }}
                  </div>
                </template>
                <div class="procedure-detail">
                  <NSpace vertical size="small">
                    <NSpace align="center">
                      <NText strong>完成进度：</NText>
                      <NProgress
                        type="line"
                        :percentage="Math.round(procedure.progress)"
                        :status="procedure.is_completed ? 'success' : 'default'"
                        style="width: 300px"
                      />
                      <NText>{{ procedure.reported_qty }} / {{ procedure.total_quantity }} 件</NText>
                    </NSpace>
                    <NSpace>
                      <NTag v-if="procedure.is_completed" type="success" size="small">
                        ✓ 已完成
                      </NTag>
                      <NTag v-else-if="procedure.reported_qty > 0" type="warning" size="small">
                        进行中
                      </NTag>
                      <NTag v-else type="default" size="small">
                        未开始
                      </NTag>
                    </NSpace>
                  </NSpace>
                </div>
              </NTimelineItem>
            </NTimeline>
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

.timeline-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #5ea3f2;
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.procedure-detail {
  padding: 8px 0;
}
</style>

