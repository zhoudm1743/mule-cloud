<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { NCard, NSpace, NTag, NButton, NTimeline, NTimelineItem, NModal, NForm, NFormItem, NInput, NSelect, NSpin, NEmpty } from 'naive-ui'
import { fetchOrderWorkflowState, fetchOrderAvailableTransitions, executeOrderWorkflowTransition } from '@/service/api/order'

const props = defineProps<{
  orderId: string
}>()

const loading = ref(false)
const workflowInstance = ref<Api.Order.WorkflowInstance | null>(null)
const availableTransitions = ref<Api.Order.WorkflowTransition[]>([])

// 状态转换模态框
const showTransitionModal = ref(false)
const selectedTransition = ref<Api.Order.WorkflowTransition | null>(null)
const transitionForm = ref({
  reason: '',
  metadata: {} as Record<string, any>
})

// 状态类型映射
const stateTypeMap: Record<string, 'default' | 'info' | 'success' | 'warning' | 'error'> = {
  draft: 'default',
  ordered: 'info',
  production: 'warning',
  completed: 'success',
  cancelled: 'error'
}

// 状态名称映射
const stateNameMap: Record<string, string> = {
  draft: '草稿',
  ordered: '已下单',
  production: '生产中',
  completed: '已完成',
  cancelled: '已取消'
}

// 获取状态的显示名称
const getStateName = (state: string) => {
  return stateNameMap[state] || state
}

// 获取状态的类型
const getStateType = (state: string) => {
  return stateTypeMap[state] || 'default'
}

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp * 1000)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 加载工作流状态
async function loadWorkflowState() {
  if (!props.orderId) return
  
  loading.value = true
  try {
    // 获取工作流实例
    const stateRes = await fetchOrderWorkflowState(props.orderId)
    if (stateRes.data) {
      workflowInstance.value = stateRes.data.instance
    }

    // 获取可用转换
    const transitionsRes = await fetchOrderAvailableTransitions(props.orderId)
    if (transitionsRes.data) {
      availableTransitions.value = transitionsRes.data.transitions || []
    }
  } catch (error: any) {
    console.error('加载工作流状态失败:', error)
    // 如果是404或工作流不存在，不显示错误
    if (!error.message?.includes('404') && !error.message?.includes('未关联')) {
      window.$message.error(error.message || '加载工作流状态失败')
    }
  } finally {
    loading.value = false
  }
}

// 打开转换模态框
function openTransitionModal(transition: Api.Order.WorkflowTransition) {
  selectedTransition.value = transition
  transitionForm.value = {
    reason: '',
    metadata: {}
  }
  showTransitionModal.value = true
}

// 执行状态转换
async function handleTransition() {
  if (!selectedTransition.value) return false

  try {
    await executeOrderWorkflowTransition(props.orderId, {
      event: selectedTransition.value.event,
      reason: transitionForm.value.reason,
      metadata: transitionForm.value.metadata
    })
    window.$message.success('状态转换成功')
    // 重新加载工作流状态
    await loadWorkflowState()
    return true // 返回 true 关闭模态框
  } catch (error: any) {
    window.$message.error(error.message || '状态转换失败')
    return false // 返回 false 保持模态框打开
  }
}

// 历史记录排序（最新的在上面）
const sortedHistory = computed(() => {
  if (!workflowInstance.value?.history) return []
  return [...workflowInstance.value.history].sort((a, b) => b.timestamp - a.timestamp)
})

onMounted(() => {
  loadWorkflowState()
})

defineExpose({
  refresh: loadWorkflowState
})
</script>

<template>
  <NCard title="工作流状态" size="small" :bordered="false">
    <NSpin :show="loading">
      <div v-if="!workflowInstance" class="text-center py-8 text-gray-400">
        <NEmpty description="该订单未关联工作流" />
      </div>
      
      <NSpace v-else vertical size="large">
        <!-- 当前状态 -->
        <div>
          <div class="text-sm text-gray-500 mb-2">当前状态</div>
          <NTag :type="getStateType(workflowInstance.current_state)" size="large">
            {{ getStateName(workflowInstance.current_state) }}
          </NTag>
        </div>

        <!-- 可用转换 -->
        <div v-if="availableTransitions.length > 0">
          <div class="text-sm text-gray-500 mb-2">可执行操作</div>
          <NSpace>
            <NButton
              v-for="trans in availableTransitions"
              :key="trans.id"
              type="primary"
              @click="openTransitionModal(trans)"
            >
              {{ trans.name }}
            </NButton>
          </NSpace>
        </div>

        <!-- 历史记录 -->
        <div v-if="sortedHistory.length > 0">
          <div class="text-sm text-gray-500 mb-2">状态历史</div>
          <NTimeline>
            <NTimelineItem
              v-for="(item, index) in sortedHistory"
              :key="index"
              :type="getStateType(item.to_state)"
              :title="item.event === 'init' ? '初始化' : getStateName(item.to_state)"
              :time="formatTime(item.timestamp)"
            >
              <div class="text-sm">
                <div v-if="item.from_state">
                  <span class="text-gray-500">从</span>
                  <NTag :type="getStateType(item.from_state)" size="small" class="mx-1">
                    {{ getStateName(item.from_state) }}
                  </NTag>
                  <span class="text-gray-500">转换到</span>
                  <NTag :type="getStateType(item.to_state)" size="small" class="mx-1">
                    {{ getStateName(item.to_state) }}
                  </NTag>
                </div>
                <div v-else class="text-gray-500">
                  初始化为 
                  <NTag :type="getStateType(item.to_state)" size="small" class="mx-1">
                    {{ getStateName(item.to_state) }}
                  </NTag>
                </div>
                <div v-if="item.reason" class="text-gray-600 mt-1">
                  原因：{{ item.reason }}
                </div>
                <div v-if="item.operator" class="text-gray-400 mt-1">
                  操作人：{{ item.operator }}
                </div>
              </div>
            </NTimelineItem>
          </NTimeline>
        </div>
      </NSpace>
    </NSpin>

    <!-- 状态转换模态框 -->
    <NModal
      v-model:show="showTransitionModal"
      preset="dialog"
      :title="selectedTransition?.name || '状态转换'"
      positive-text="确认"
      negative-text="取消"
      :loading="loading"
      :trap-focus="true"
      :mask-closable="false"
      @positive-click="handleTransition"
    >
      <NForm :model="transitionForm" label-placement="left" label-width="80">
        <NFormItem label="转换说明">
          <div class="text-sm text-gray-600">
            <div>
              从 
              <NTag :type="getStateType(selectedTransition?.from_state || '')" size="small" class="mx-1">
                {{ getStateName(selectedTransition?.from_state || '') }}
              </NTag>
              转换到
              <NTag :type="getStateType(selectedTransition?.to_state || '')" size="small" class="mx-1">
                {{ getStateName(selectedTransition?.to_state || '') }}
              </NTag>
            </div>
            <div v-if="selectedTransition?.description" class="mt-2">
              {{ selectedTransition.description }}
            </div>
          </div>
        </NFormItem>

        <!-- 显示条件 -->
        <NFormItem v-if="selectedTransition?.conditions && selectedTransition.conditions.length > 0" label="转换条件">
          <div class="text-sm">
            <div v-for="(cond, idx) in selectedTransition.conditions" :key="idx" class="mb-1">
              <span class="text-gray-600">{{ cond.description || `${cond.field} ${cond.operator} ${cond.value}` }}</span>
            </div>
          </div>
        </NFormItem>

        <NFormItem label="转换原因">
          <NInput
            v-model:value="transitionForm.reason"
            type="textarea"
            placeholder="请输入转换原因（可选）"
            :rows="3"
          />
        </NFormItem>
      </NForm>
    </NModal>
  </NCard>
</template>

<style scoped lang="scss">
.mx-1 {
  margin-left: 0.25rem;
  margin-right: 0.25rem;
}

.mt-1 {
  margin-top: 0.25rem;
}

.mt-2 {
  margin-top: 0.5rem;
}

.mb-1 {
  margin-bottom: 0.25rem;
}

.mb-2 {
  margin-bottom: 0.5rem;
}

.py-8 {
  padding-top: 2rem;
  padding-bottom: 2rem;
}
</style>

