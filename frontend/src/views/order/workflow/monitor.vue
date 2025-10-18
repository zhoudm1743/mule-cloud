<template>
  <NSpace vertical size="large" class="p-4">
    <NCard title="订单工作流管理" :bordered="false" class="rounded-8px shadow-sm">
      <NTabs v-model:value="activeTab" type="line">
        <!-- 工作流可视化 -->
        <NTabPane name="visualization" tab="工作流可视化">
          <NSpace vertical size="large">
            <NAlert type="info" closable>
              订单工作流定义了订单从创建到完成的整个生命周期，包括状态、事件和转换规则。
            </NAlert>

            <!-- Mermaid 流程图 -->
            <NCard title="流程图" :bordered="false">
              <NSpin :show="loading">
                <div v-if="mermaidCode" class="mermaid-container">
                  <pre class="mermaid">{{ mermaidCode }}</pre>
                </div>
              </NSpin>
            </NCard>

            <!-- 状态说明 -->
            <NCard title="状态说明" :bordered="false">
              <NSpace>
                <NTag
                  v-for="state in states"
                  :key="state.id"
                  :type="getStateTagType(state.type)"
                  :bordered="false"
                  size="large"
                >
                  {{ state.name }}
                </NTag>
              </NSpace>
            </NCard>

            <!-- 事件说明 -->
            <NCard title="事件说明" :bordered="false">
              <NDataTable
                :columns="eventColumns"
                :data="events"
                :pagination="false"
                :single-line="false"
              />
            </NCard>
          </NSpace>
        </NTabPane>

        <!-- 转换规则 -->
        <NTabPane name="rules" tab="转换规则">
          <NSpace vertical size="large">
            <NAlert type="warning" closable>
              转换规则定义了订单状态之间的转换条件和权限要求。
            </NAlert>

            <NDataTable
              :columns="ruleColumns"
              :data="rules"
              :pagination="false"
              :single-line="false"
            />
          </NSpace>
        </NTabPane>

        <!-- 订单状态查询 -->
        <NTabPane name="query" tab="订单状态查询">
          <NSpace vertical size="large">
            <NForm :model="queryForm" label-placement="left" inline>
              <NFormItem label="订单ID">
                <NInput v-model:value="queryForm.order_id" placeholder="请输入订单ID" />
              </NFormItem>
              <NFormItem>
                <NButton type="primary" @click="handleQueryOrder">
                  查询
                </NButton>
              </NFormItem>
            </NForm>

            <!-- 当前状态 -->
            <NCard v-if="currentOrder" title="当前状态" :bordered="false">
              <NDescriptions :column="3" label-placement="left">
                <NDescriptionsItem label="订单ID">
                  {{ currentOrder.order_id }}
                </NDescriptionsItem>
                <NDescriptionsItem label="状态">
                  <NTag :type="getStatusTagType(currentOrder.status)" :bordered="false">
                    {{ currentOrder.status_name }}
                  </NTag>
                </NDescriptionsItem>
              </NDescriptions>
            </NCard>

            <!-- 状态历史 -->
            <NCard v-if="orderHistory.length > 0" title="状态历史" :bordered="false">
              <NTimeline>
                <NTimelineItem
                  v-for="(item, index) in orderHistory"
                  :key="index"
                  :type="getHistoryType(item)"
                  :title="`${getStateName(item.from_state)} → ${getStateName(item.to_state)}`"
                  :time="formatTimestamp(item.timestamp)"
                >
                  <NText depth="3">
                    事件: {{ item.event }}
                  </NText>
                  <br>
                  <NText depth="3">
                    操作人: {{ item.operator }}
                  </NText>
                  <br>
                  <NText v-if="item.reason" depth="3">
                    原因: {{ item.reason }}
                  </NText>
                </NTimelineItem>
              </NTimeline>
            </NCard>

            <!-- 回滚历史 -->
            <NCard v-if="rollbackHistory.length > 0" title="回滚历史" :bordered="false">
              <NTimeline>
                <NTimelineItem
                  v-for="(item, index) in rollbackHistory"
                  :key="index"
                  type="error"
                  :title="`回滚: ${getStateName(item.rollback_from)} → ${getStateName(item.rollback_to)}`"
                  :time="formatTimestamp(item.timestamp)"
                >
                  <NText depth="3">
                    原始事件: {{ item.original_event }}
                  </NText>
                  <br>
                  <NText depth="3">
                    回滚原因: {{ item.reason }}
                  </NText>
                  <br>
                  <NText depth="3">
                    操作人: {{ item.operator }}
                  </NText>
                </NTimelineItem>
              </NTimeline>
            </NCard>

            <!-- 操作按钮 -->
            <NSpace v-if="currentOrder">
              <NButton type="error" @click="openRollbackModal">
                回滚状态
              </NButton>
            </NSpace>
          </NSpace>
        </NTabPane>
      </NTabs>
    </NCard>

    <!-- 回滚模态框 -->
    <NModal
      v-model:show="showRollbackModal"
      preset="dialog"
      title="回滚订单状态"
      positive-text="确定"
      negative-text="取消"
      @positive-click="handleRollback"
    >
      <NForm label-placement="left" label-width="100px" class="mt-4">
        <NFormItem label="订单ID">
          <NText>{{ currentOrder?.order_id }}</NText>
        </NFormItem>
        <NFormItem label="当前状态">
          <NTag :type="getStatusTagType(currentOrder?.status || 0)" :bordered="false">
            {{ currentOrder?.status_name }}
          </NTag>
        </NFormItem>
        <NFormItem label="回滚原因" required>
          <NInput
            v-model:value="rollbackForm.reason"
            type="textarea"
            placeholder="请输入回滚原因"
            :rows="3"
          />
        </NFormItem>
      </NForm>
    </NModal>
  </NSpace>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NDescriptions,
  NDescriptionsItem,
  NForm,
  NFormItem,
  NInput,
  NModal,
  NSpace,
  NSpin,
  NTabPane,
  NTabs,
  NTag,
  NText,
  NTimeline,
  NTimelineItem,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import {
  fetchMermaidDiagram,
  fetchOrderHistory,
  fetchOrderStatus,
  fetchRollbackHistory,
  fetchTransitionRules,
  fetchWorkflowDefinition,
  rollbackOrder,
} from '@/service/api/workflow'
import { useBoolean } from '@/hooks'

defineOptions({ name: 'OrderWorkflow' })

const activeTab = ref('visualization')
const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

// 工作流定义
const states = ref<Api.Workflow.StateDefinition[]>([])
const events = ref<Api.Workflow.EventDefinition[]>([])
const transitions = ref<Api.Workflow.TransitionDefinition[]>([])
const mermaidCode = ref('')

// 转换规则
const rules = ref<Api.Workflow.TransitionRule[]>([])

// 订单查询
const queryForm = ref({
  order_id: '',
})
const currentOrder = ref<Api.Workflow.OrderStatus | null>(null)
const orderHistory = ref<Api.Workflow.StateHistory[]>([])
const rollbackHistory = ref<Api.Workflow.RollbackRecord[]>([])

// 回滚
const showRollbackModal = ref(false)
const rollbackForm = ref({
  reason: '',
})

// 事件表格列
const eventColumns: DataTableColumns<Api.Workflow.EventDefinition> = [
  { title: '事件名称', key: 'label' },
  { title: '事件标识', key: 'name' },
  {
    title: '特殊要求',
    key: 'requirements',
    render: (row) => {
      const requirements = []
      if (row.requireCondition)
        requirements.push('需要满足条件')
      if (row.requireRole)
        requirements.push(`需要角色: ${row.requireRole}`)
      return requirements.length > 0 ? requirements.join(', ') : '-'
    },
  },
]

// 转换规则表格列
const ruleColumns: DataTableColumns<Api.Workflow.TransitionRule> = [
  { title: '起始状态', key: 'from_name' },
  { title: '目标状态', key: 'to_name' },
  { title: '触发事件', key: 'event' },
  {
    title: '条件',
    key: 'has_condition',
    render: row => (row.has_condition ? '需要满足特定条件' : '-'),
  },
  {
    title: '权限要求',
    key: 'require_role',
    render: row => (row.require_role ? row.require_role : '-'),
  },
]

onMounted(() => {
  loadWorkflowDefinition()
  loadTransitionRules()
  initMermaid()
})

// 加载工作流定义
async function loadWorkflowDefinition() {
  startLoading()
  try {
    const { data } = await fetchWorkflowDefinition()
    if (data) {
      states.value = data.states
      events.value = data.events
      transitions.value = data.transitions
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取工作流定义失败')
  }
  finally {
    endLoading()
  }
}

// 加载转换规则
async function loadTransitionRules() {
  try {
    const { data } = await fetchTransitionRules()
    if (data) {
      rules.value = data.rules
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取转换规则失败')
  }
}

// 初始化 Mermaid
async function initMermaid() {
  try {
    const { data } = await fetchMermaidDiagram()
    if (data) {
      mermaidCode.value = data.diagram

      // 动态加载 Mermaid 库
      if (typeof window !== 'undefined' && !(window as any).mermaid) {
        const script = document.createElement('script')
        script.src = 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js'
        script.onload = () => {
          (window as any).mermaid.initialize({ startOnLoad: true, theme: 'default' })
          setTimeout(() => {
            (window as any).mermaid.contentLoaded()
          }, 100)
        }
        document.head.appendChild(script)
      }
      else {
        setTimeout(() => {
          (window as any).mermaid.contentLoaded()
        }, 100)
      }
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取流程图失败')
  }
}

// 查询订单
async function handleQueryOrder() {
  if (!queryForm.value.order_id) {
    window.$message.warning('请输入订单ID')
    return
  }

  try {
    // 获取当前状态
    const statusRes = await fetchOrderStatus(queryForm.value.order_id)
    if (statusRes.data) {
      currentOrder.value = statusRes.data
    }

    // 获取状态历史
    const historyRes = await fetchOrderHistory(queryForm.value.order_id)
    if (historyRes.data) {
      orderHistory.value = historyRes.data.history
    }

    // 获取回滚历史
    const rollbackRes = await fetchRollbackHistory(queryForm.value.order_id)
    if (rollbackRes.data) {
      rollbackHistory.value = rollbackRes.data.rollbacks
    }

    window.$message.success('查询成功')
  }
  catch (error: any) {
    window.$message.error(error.message || '查询失败')
  }
}

// 打开回滚模态框
function openRollbackModal() {
  rollbackForm.value.reason = ''
  showRollbackModal.value = true
}

// 执行回滚
async function handleRollback() {
  if (!rollbackForm.value.reason) {
    window.$message.warning('请输入回滚原因')
    return
  }

  if (!currentOrder.value) {
    window.$message.warning('请先查询订单')
    return
  }

  try {
    await rollbackOrder({
      order_id: currentOrder.value.order_id,
      reason: rollbackForm.value.reason,
    })

    window.$message.success('回滚成功')
    showRollbackModal.value = false

    // 重新查询
    handleQueryOrder()
  }
  catch (error: any) {
    window.$message.error(error.message || '回滚失败')
  }
}

// 获取状态标签类型
function getStateTagType(type: string) {
  switch (type) {
    case 'start':
      return 'info'
    case 'end':
      return 'success'
    default:
      return 'warning'
  }
}

// 获取状态标签类型（根据状态值）
function getStatusTagType(status: number) {
  switch (status) {
    case 0:
      return 'default'
    case 1:
      return 'info'
    case 2:
      return 'warning'
    case 3:
      return 'success'
    case 4:
      return 'error'
    default:
      return 'default'
  }
}

// 获取历史项类型
function getHistoryType(item: Api.Workflow.StateHistory) {
  if (item.metadata?.is_rollback)
    return 'error'
  if (item.to_state === 3)
    return 'success'
  if (item.to_state === 4)
    return 'error'
  return 'info'
}

// 获取状态名称
function getStateName(status: number) {
  const stateMap: Record<number, string> = {
    0: '草稿',
    1: '已下单',
    2: '生产中',
    3: '已完成',
    4: '已取消',
  }
  return stateMap[status] || '未知'
}

// 格式化时间戳
function formatTimestamp(timestamp: number) {
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped>
.mermaid-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
}

.mermaid {
  font-size: 14px;
}
</style>

