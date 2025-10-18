<template>
  <NSpace vertical size="large" class="p-4">
    <NCard title="工作流设计器" :bordered="false" class="rounded-8px shadow-sm">
      <!-- 工具栏 -->
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="handleSave">
            <template #icon>
              <nova-icon icon="carbon:save" />
            </template>
            保存
          </NButton>
          <NButton type="success" @click="handlePublish">
            <template #icon>
              <nova-icon icon="carbon:checkmark" />
            </template>
            发布
          </NButton>
          <NButton @click="handlePreview">
            <template #icon>
              <nova-icon icon="carbon:view" />
            </template>
            预览
          </NButton>
          <NButton @click="handleBack">
            返回列表
          </NButton>
        </NSpace>
      </template>

      <!-- 主体区域 -->
      <NLayout has-sider style="height: calc(100vh - 220px)">
        <!-- 左侧组件面板 -->
        <NLayoutSider
          bordered
          width="200"
          :native-scrollbar="false"
        >
          <div class="p-4">
            <NText strong class="mb-4 block">
              节点类型
            </NText>
            <NSpace vertical>
              <div
                class="node-item"
                @mousedown="(e) => handleNodeDragStart('start', e)"
              >
                <nova-icon icon="carbon:play-filled" class="mr-2" />
                开始节点
              </div>
              <div
                class="node-item"
                @mousedown="(e) => handleNodeDragStart('normal', e)"
              >
                <nova-icon icon="carbon:checkbox" class="mr-2" />
                普通节点
              </div>
              <div
                class="node-item"
                @mousedown="(e) => handleNodeDragStart('end', e)"
              >
                <nova-icon icon="carbon:checkmark-filled" class="mr-2" />
                结束节点
              </div>
            </NSpace>
          </div>
        </NLayoutSider>

        <!-- 中间画布区域 -->
        <NLayoutContent
          class="canvas-container"
          content-style="background: #f5f5f5; position: relative;"
        >
          <div id="logic-flow-container" ref="containerRef" />
          
          <!-- 小地图（暂时禁用以避免销毁错误） -->
          <!-- <div id="logic-flow-minimap" class="minimap" /> -->
        </NLayoutContent>

        <!-- 右侧配置面板 -->
        <NLayoutSider
          bordered
          width="350"
          :native-scrollbar="false"
        >
          <div class="p-4">
            <NText strong class="mb-4 block">
              {{ selectedElement ? '属性配置' : '工作流信息' }}
            </NText>

            <!-- 工作流基本信息 -->
            <NForm v-if="!selectedElement" label-placement="top">
              <NFormItem label="工作流名称">
                <NInput v-model:value="workflowForm.name" placeholder="请输入工作流名称" />
              </NFormItem>
              <NFormItem label="工作流编码">
                <NInput v-model:value="workflowForm.code" placeholder="请输入工作流编码" />
              </NFormItem>
              <NFormItem label="描述">
                <NInput
                  v-model:value="workflowForm.description"
                  type="textarea"
                  placeholder="请输入描述"
                  :rows="3"
                />
              </NFormItem>
            </NForm>

            <!-- 节点配置 -->
            <NForm v-else-if="selectedElement.type === 'node'" label-placement="top">
              <NFormItem label="节点名称">
                <NInput v-model:value="selectedElement.properties.stateName" />
              </NFormItem>
              <NFormItem label="节点类型">
                <NSelect
                  v-model:value="selectedElement.properties.stateType"
                  :options="[
                    { label: '开始', value: 'start' },
                    { label: '普通', value: 'normal' },
                    { label: '结束', value: 'end' },
                  ]"
                />
              </NFormItem>
              <NFormItem label="节点颜色">
                <NColorPicker v-model:value="selectedElement.properties.stateColor" />
              </NFormItem>
              <NFormItem label="描述">
                <NInput
                  v-model:value="selectedElement.properties.description"
                  type="textarea"
                  :rows="3"
                />
              </NFormItem>
              <NSpace>
                <NButton type="primary" @click="handleUpdateNode">
                  更新
                </NButton>
                <NButton type="error" @click="handleDeleteNode">
                  删除
                </NButton>
              </NSpace>
            </NForm>

            <!-- 连线配置 -->
            <NForm v-else-if="selectedElement.type === 'edge'" label-placement="top">
              <NFormItem label="事件名称">
                <NInput v-model:value="selectedElement.properties.event" />
              </NFormItem>
              <NFormItem label="权限要求">
                <NInput v-model:value="selectedElement.properties.requireRole" placeholder="如：admin" />
              </NFormItem>
              <NFormItem label="描述">
                <NInput
                  v-model:value="selectedElement.properties.description"
                  type="textarea"
                  :rows="2"
                />
              </NFormItem>

              <NDivider />

              <NText strong class="mb-2 block">
                转换条件
              </NText>
              <NSpace vertical class="w-full mb-4">
                <NCard
                  v-for="(condition, index) in selectedElement.properties.conditions || []"
                  :key="index"
                  size="small"
                  :bordered="false"
                  class="bg-gray-50"
                >
                  <template #header>
                    条件 {{ index + 1 }}
                    <NButton
                      text
                      type="error"
                      size="small"
                      class="float-right"
                      @click="removeCondition(index)"
                    >
                      删除
                    </NButton>
                  </template>
                  <NSpace vertical class="w-full">
                    <NFormItem label="字段" size="small">
                      <NSelect
                        v-model:value="condition.field"
                        :options="availableFields"
                        placeholder="选择字段"
                        size="small"
                        @update:value="(value) => handleFieldChange(condition, value)"
                      />
                    </NFormItem>
                    <NFormItem label="操作符" size="small">
                      <NSelect
                        v-model:value="condition.operator"
                        :options="getOperatorsForField(condition.field)"
                        placeholder="选择操作符"
                        size="small"
                      />
                    </NFormItem>
                    <NFormItem label="比较值" size="small">
                      <NInputNumber
                        v-if="getFieldType(condition.field) === 'number'"
                        v-model:value="condition.value"
                        placeholder="输入数值"
                        size="small"
                        class="w-full"
                      />
                      <NSwitch
                        v-else-if="getFieldType(condition.field) === 'boolean'"
                        v-model:value="condition.value"
                        size="small"
                      />
                      <NInput
                        v-else
                        v-model:value="condition.value"
                        placeholder="输入值"
                        size="small"
                      />
                    </NFormItem>
                    <NFormItem v-if="getFieldDescription(condition.field)" label="说明" size="small">
                      <NText depth="3" style="font-size: 12px">
                        {{ getFieldDescription(condition.field) }}
                      </NText>
                    </NFormItem>
                  </NSpace>
                </NCard>
                <NButton dashed block @click="addCondition">
                  + 添加条件
                </NButton>
              </NSpace>

              <NSpace>
                <NButton type="primary" @click="handleUpdateEdge">
                  更新
                </NButton>
                <NButton type="error" @click="handleDeleteEdge">
                  删除
                </NButton>
              </NSpace>
            </NForm>
          </div>
        </NLayoutSider>
      </NLayout>
    </NCard>
  </NSpace>
</template>

<script setup lang="ts">
import { nextTick, onMounted, onUnmounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  NButton,
  NCard,
  NColorPicker,
  NDivider,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NLayout,
  NLayoutContent,
  NLayoutSider,
  NSelect,
  NSpace,
  NSwitch,
  NText,
} from 'naive-ui'
import LogicFlow from '@logicflow/core'
import { Menu, SelectionSelect, Snapshot } from '@logicflow/extension'
import '@logicflow/core/dist/index.css'
import '@logicflow/extension/dist/index.css'
import {
  createWorkflowDefinition,
  fetchWorkflowDefinition,
  fetchWorkflowTemplates,
  updateWorkflowDefinition,
} from '@/service/api/workflow-designer'

defineOptions({ name: 'WorkflowDesigner' })

const route = useRoute()
const router = useRouter()
const containerRef = ref()
let lf: LogicFlow | null = null

// 工作流表单
const workflowForm = ref({
  name: '',
  code: '',
  description: '',
})

// 选中的元素
const selectedElement = ref<any>(null)

// 可用的条件字段定义
const availableFields = [
  { label: '订单金额 (total_amount)', value: 'total_amount', type: 'number', description: '订单总金额，单位：元' },
  { label: '订单数量 (quantity)', value: 'quantity', type: 'number', description: '订单商品总数量' },
  { label: '完成进度 (progress)', value: 'progress', type: 'number', description: '订单完成进度，范围 0-1，1表示100%' },
  { label: '订单状态 (status)', value: 'status', type: 'number', description: '0-草稿 1-已下单 2-生产中 3-已完成 4-已取消' },
  { label: '单价 (unit_price)', value: 'unit_price', type: 'number', description: '商品单价' },
  { label: '合同号 (contract_no)', value: 'contract_no', type: 'string', description: '订单合同编号' },
  { label: '客户ID (customer_id)', value: 'customer_id', type: 'string', description: '客户唯一标识' },
  { label: '款式ID (style_id)', value: 'style_id', type: 'string', description: '款式唯一标识' },
  { label: '交货日期 (delivery_date)', value: 'delivery_date', type: 'string', description: '订单交货日期' },
]

// 获取字段类型
function getFieldType(fieldValue: string): string {
  const field = availableFields.find(f => f.value === fieldValue)
  return field?.type || 'string'
}

// 获取字段说明
function getFieldDescription(fieldValue: string): string {
  const field = availableFields.find(f => f.value === fieldValue)
  return field?.description || ''
}

// 根据字段类型获取可用的操作符
function getOperatorsForField(fieldValue: string) {
  const fieldType = getFieldType(fieldValue)
  
  if (fieldType === 'number') {
    return [
      { label: '等于 (=)', value: 'eq' },
      { label: '不等于 (≠)', value: 'ne' },
      { label: '大于 (>)', value: 'gt' },
      { label: '大于等于 (≥)', value: 'gte' },
      { label: '小于 (<)', value: 'lt' },
      { label: '小于等于 (≤)', value: 'lte' },
    ]
  }
  
  if (fieldType === 'boolean') {
    return [
      { label: '等于', value: 'eq' },
      { label: '不等于', value: 'ne' },
    ]
  }
  
  // string 类型
  return [
    { label: '等于', value: 'eq' },
    { label: '不等于', value: 'ne' },
    { label: '包含', value: 'contains' },
    { label: '开头是', value: 'starts_with' },
    { label: '结尾是', value: 'ends_with' },
  ]
}

// 字段变化时的处理
function handleFieldChange(condition: any, value: string) {
  // 清空之前的值
  condition.value = undefined
  condition.operator = 'eq'
  
  // 根据字段类型设置默认值
  const fieldType = getFieldType(value)
  if (fieldType === 'number') {
    condition.value = 0
  } else if (fieldType === 'boolean') {
    condition.value = false
  } else {
    condition.value = ''
  }
}

// 初始化LogicFlow
onMounted(async () => {
  // 先初始化 LogicFlow 实例
  initLogicFlow()
  
  // 等待 DOM 更新和 LogicFlow 初始化完成
  await nextTick()
  
  // 加载现有工作流（如果是编辑模式）
  const workflowId = route.query.id as string
  if (workflowId) {
    await loadWorkflow(workflowId)
  }
  
  // 加载模板（如果是从模板创建）
  const templateId = route.query.template as string
  if (templateId) {
    await loadTemplate(templateId)
  }
})

// 初始化LogicFlow实例
function initLogicFlow() {
  lf = new LogicFlow({
    container: document.getElementById('logic-flow-container')!,
    grid: {
      size: 10,
      visible: true,
    },
    background: {
      color: '#f7f9ff',
    },
    keyboard: {
      enabled: true,
    },
    style: {
      rect: {
        width: 100,
        height: 50,
        radius: 5,
        fill: '#409EFF',
        stroke: '#409EFF',
        strokeWidth: 2,
      },
      nodeText: {
        fontSize: 12,
        color: '#fff',
      },
      edgeText: {
        fontSize: 12,
        color: '#666',
        background: {
          fill: '#f5f5f5',
          padding: 4,
        },
      },
    },
    plugins: [Menu, SelectionSelect, Snapshot],
  })

  // 设置默认边类型
  lf.setDefaultEdgeType('polyline')

  // 绑定事件
  bindEvents()

  // 渲染数据
  if (workflowForm.value.name) {
    renderWorkflow()
  }
  else {
    // 初始示例数据
    lf.render({
      nodes: [],
      edges: [],
    })
  }
}

// 注册自定义节点
function registerCustomNodes() {
  if (!lf)
    return

  // 开始节点 - 使用内置的 rect 节点类型，通过样式自定义
  lf.setDefaultEdgeType('polyline')
}

// 绑定事件
function bindEvents() {
  if (!lf)
    return

  // 节点点击
  lf.on('node:click', ({ data }) => {
    selectedElement.value = {
      type: 'node',
      id: data.id,
      properties: { ...data.properties },
    }
  })

  // 边点击
  lf.on('edge:click', ({ data }) => {
    selectedElement.value = {
      type: 'edge',
      id: data.id,
      properties: { ...data.properties },
    }
  })

  // 画布点击
  lf.on('blank:click', () => {
    selectedElement.value = null
  })
}

// 节点拖拽开始
function handleNodeDragStart(type: string, event?: MouseEvent) {
  if (!lf)
    return

  // 防止选中文字
  if (event) {
    event.preventDefault()
  }

  const colorMap = {
    start: '#67C23A',
    normal: '#409EFF',
    end: '#F56C6C',
  }

  const nameMap = {
    start: '开始',
    normal: '新状态',
    end: '结束',
  }

  lf.dnd.startDrag({
    type: 'rect', // 使用内置的 rect 类型
    text: nameMap[type] || '状态',
    properties: {
      stateId: `state_${Date.now()}`,
      stateName: nameMap[type] || '新状态',
      stateType: type,
      stateColor: colorMap[type] || '#409EFF',
    },
  })
}

// 更新节点
function handleUpdateNode() {
  if (!lf || !selectedElement.value)
    return

  lf.updateNodeData({
    id: selectedElement.value.id,
    properties: selectedElement.value.properties,
    text: {
      value: selectedElement.value.properties.stateName,
    },
  })

  window.$message.success('节点更新成功')
}

// 删除节点
function handleDeleteNode() {
  if (!lf || !selectedElement.value)
    return

  lf.deleteNode(selectedElement.value.id)
  selectedElement.value = null
  window.$message.success('节点删除成功')
}

// 更新连线
function handleUpdateEdge() {
  if (!lf || !selectedElement.value)
    return

  lf.updateEdgeData({
    id: selectedElement.value.id,
    properties: selectedElement.value.properties,
    text: {
      value: selectedElement.value.properties.event || '',
    },
  })

  window.$message.success('连线更新成功')
}

// 删除连线
function handleDeleteEdge() {
  if (!lf || !selectedElement.value)
    return

  lf.deleteEdge(selectedElement.value.id)
  selectedElement.value = null
  window.$message.success('连线删除成功')
}

// 添加条件
function addCondition() {
  if (!selectedElement.value)
    return

  if (!selectedElement.value.properties.conditions) {
    selectedElement.value.properties.conditions = []
  }

  selectedElement.value.properties.conditions.push({
    type: 'field',
    field: '',
    operator: 'eq',
    value: '',
    description: '',
  })
}

// 删除条件
function removeCondition(index: number) {
  if (!selectedElement.value)
    return

  selectedElement.value.properties.conditions.splice(index, 1)
}

// 加载工作流
async function loadWorkflow(id: string) {
  try {
    console.log('🔍 开始加载工作流', id)
    const { data } = await fetchWorkflowDefinition(id)
    console.log('📦 获取到工作流数据', data)
    
    if (data) {
      workflowForm.value = {
        name: data.name,
        code: data.code,
        description: data.description,
      }
      
      console.log('✅ 工作流表单已填充', workflowForm.value)
      
      // 如果有状态数据，渲染到画布
      if (data.states && data.states.length > 0) {
        console.log('✅ 工作流包含状态定义，准备渲染', {
          statesCount: data.states.length,
          transitionsCount: data.transitions?.length || 0
        })
        renderWorkflowToCanvas(data)
      } else {
        console.warn('⚠️ 工作流不包含状态定义')
      }
    }
  }
  catch (error: any) {
    console.error('❌ 加载工作流失败', error)
    window.$message.error(error.message || '加载工作流失败')
  }
}

// 加载模板
async function loadTemplate(templateId: string) {
  try {
    console.log('🔍 开始加载模板', templateId)
    const { data } = await fetchWorkflowTemplates()
    console.log('📦 获取到模板数据', data)
    
    if (data?.templates) {
      const template = data.templates.find(t => t.id === templateId)
      console.log('🎯 找到模板', template)
      
      if (template) {
        // 应用模板数据到工作流表单
        workflowForm.value = {
          name: template.name,
          code: template.code || templateId,
          description: template.description,
        }
        
        console.log('✅ 工作流表单已填充', workflowForm.value)
        console.log('🔍 模板状态数据', { 
          hasStates: !!template.states, 
          statesCount: template.states?.length,
          states: template.states
        })
        
        // 如果模板包含状态和转换定义，自动渲染
        if (template.states && template.states.length > 0) {
          console.log('✅ 模板包含状态定义，准备渲染')
          renderTemplateToCanvas(template)
          window.$message.success(`已应用模板：${template.name}`)
        } else {
          console.warn('⚠️ 模板不包含状态定义')
          window.$message.success(`已应用模板：${template.name}，请在画布上拖拽节点创建工作流`)
        }
      } else {
        console.error('❌ 未找到指定的模板')
        window.$message.warning('未找到指定的模板')
      }
    } else {
      console.error('❌ 模板数据为空')
    }
  }
  catch (error: any) {
    console.error('❌ 加载模板失败', error)
    window.$message.error(error.message || '加载模板失败')
  }
}

// 将模板渲染到画布上
function renderTemplateToCanvas(template: Api.WorkflowDesigner.WorkflowTemplate) {
  console.log('🎨 开始渲染模板到画布', template)
  
  if (!lf) {
    console.error('❌ LogicFlow 实例未初始化')
    return
  }
  
  if (!template.states || template.states.length === 0) {
    console.warn('⚠️ 模板没有状态定义')
    return
  }
  
  const nodes: any[] = []
  const edges: any[] = []
  
  // 计算布局位置（垂直流式布局）
  const startX = 300
  const startY = 100
  const horizontalGap = 250
  const verticalGap = 150
  
  // 创建节点
  template.states.forEach((state, index) => {
    const x = startX + (index % 3) * horizontalGap
    const y = startY + Math.floor(index / 3) * verticalGap
    
    const node = {
      id: state.code,
      type: 'rect',
      x,
      y,
      text: state.name,
      properties: {
        state: {
          code: state.code,
          name: state.name,
          type: state.type,
          description: state.description,
        },
        color: state.color,
      },
    }
    
    console.log('📦 创建节点', node)
    nodes.push(node)
  })
  
  // 创建连线
  if (template.transitions) {
    template.transitions.forEach((transition, index) => {
      const conditions: any[] = []
      
      // 如果有条件，添加条件配置
      if (transition.has_condition && transition.available_fields) {
        transition.available_fields.forEach((field) => {
          conditions.push({
            field: field.key,
            operator: 'gte',
            value: field.type === 'number' ? 1.0 : '',
          })
        })
      }
      
      const edge = {
        id: `edge_${index}`,
        type: 'polyline',
        sourceNodeId: transition.from,
        targetNodeId: transition.to,
        text: transition.event_label,  // ✅ 显示中文标签
        properties: {
          event: transition.event,         // ✅ 英文事件代码
          eventLabel: transition.event_label,  // ✅ 中文标签（保存时优先使用）
          requireRole: transition.require_role || '',
          conditions,
        },
      }
      
      console.log('🔗 创建连线', edge)
      edges.push(edge)
    })
  }
  
  console.log('✅ 准备渲染', { nodes: nodes.length, edges: edges.length })
  
  // 渲染到画布
  try {
    lf.render({ nodes, edges })
    console.log('✅ 渲染成功')
    
    // 自动居中
    setTimeout(() => {
      lf.translateCenter()
      console.log('✅ 画布已居中')
    }, 100)
  } catch (error) {
    console.error('❌ 渲染失败', error)
  }
}

// 事件名称映射表（英文 -> 中文）
const eventLabelMap: Record<string, string> = {
  'submit_order': '提交订单',
  'start_cutting': '开始裁剪',
  'start_production': '开始生产',
  'update_progress': '更新进度',
  'complete': '完成',
  'cancel': '取消',
}

// 渲染工作流到画布（从已保存的工作流定义）
function renderWorkflowToCanvas(workflow: Api.WorkflowDesigner.WorkflowDefinition) {
  console.log('🎨 开始渲染工作流到画布', workflow)
  
  if (!lf) {
    console.error('❌ LogicFlow 实例未初始化')
    return
  }
  
  if (!workflow.states || workflow.states.length === 0) {
    console.warn('⚠️ 工作流没有状态定义')
    return
  }
  
  const nodes: any[] = []
  const edges: any[] = []
  
  // 创建节点
  workflow.states.forEach((state) => {
    const node = {
      id: state.code,
      type: 'rect',
      x: state.position?.x || 300,
      y: state.position?.y || 100,
      text: state.name,
      properties: {
        state: {
          code: state.code,
          name: state.name,
          type: state.type,
          description: state.description,
        },
        color: state.color,
      },
    }
    
    console.log('📦 创建节点', node)
    nodes.push(node)
  })
  
  // 创建连线
  if (workflow.transitions) {
    workflow.transitions.forEach((transition) => {
      // 如果 name 是英文事件代码，尝试从映射表中获取中文标签
      const displayText = eventLabelMap[transition.name] || transition.name || transition.event
      
      const edge = {
        id: transition.id,
        type: 'polyline',
        sourceNodeId: transition.from_state,
        targetNodeId: transition.to_state,
        text: displayText,  // ✅ 使用映射后的中文标签
        properties: {
          transitionId: transition.id,
          event: transition.event,
          eventLabel: displayText,  // ✅ 也更新 eventLabel
          requireRole: transition.require_role || '',
          conditions: transition.conditions || [],
          actions: transition.actions || [],
          description: transition.description || '',
        },
      }
      
      console.log('🔗 创建连线', {
        原始名称: transition.name,
        显示文本: displayText,
        事件代码: transition.event
      })
      edges.push(edge)
    })
  }
  
  console.log('✅ 准备渲染', { nodes: nodes.length, edges: edges.length })
  
  // 渲染到画布
  try {
    lf.render({ nodes, edges })
    console.log('✅ 渲染成功')
    
    // 自动居中
    setTimeout(() => {
      lf.translateCenter()
      console.log('✅ 画布已居中')
    }, 100)
  } catch (error) {
    console.error('❌ 渲染失败', error)
  }
}

// 保存工作流
async function handleSave() {
  if (!lf)
    return

  if (!workflowForm.value.name || !workflowForm.value.code) {
    window.$message.warning('请填写工作流名称和编码')
    return
  }

  try {
    const graphData = lf.getGraphData()
    console.log('📊 获取画布数据', graphData)

    // 转换为工作流定义格式
    const states: Api.WorkflowDesigner.WorkflowState[] = graphData.nodes.map(node => {
      // 兼容两种格式：直接属性 或 state 对象
      const stateInfo = node.properties.state || node.properties
      const stateCode = node.id || stateInfo.code || node.properties.stateId
      
      // LogicFlow 的 text 可能是对象 {value: string} 或字符串
      let stateName = ''
      if (typeof node.text === 'object' && node.text !== null) {
        stateName = node.text.value || ''
      } else if (typeof node.text === 'string') {
        stateName = node.text
      } else {
        stateName = stateInfo.name || node.properties.stateName || ''
      }
      
      const stateType = stateInfo.type || node.properties.stateType || 'normal'
      const stateColor = node.properties.color || stateInfo.color || node.properties.stateColor || '#409EFF'
      const stateDesc = stateInfo.description || node.properties.description || ''
      
      console.log('📦 转换节点', { 
        原始节点: node,
        text类型: typeof node.text,
        text值: node.text,
        提取的状态: { code: stateCode, name: stateName, type: stateType } 
      })
      
      return {
        id: stateCode,
        name: stateName,
        code: stateCode,
        type: stateType,
        color: stateColor,
        description: stateDesc,
        position: { x: node.x!, y: node.y! },
      }
    })

    const transitions: Api.WorkflowDesigner.WorkflowTransition[] = graphData.edges.map(edge => {
      // LogicFlow 的 text 可能是对象 {value: string} 或字符串
      let edgeText = ''
      if (typeof edge.text === 'object' && edge.text !== null) {
        edgeText = edge.text.value || ''
      } else if (typeof edge.text === 'string') {
        edgeText = edge.text
      }
      
      // 优先使用 eventLabel（中文标签），然后是 text，最后才是 event（英文事件名）
      const displayName = edge.properties.eventLabel || edgeText || edge.properties.event || ''
      const eventCode = edge.properties.event || edge.properties.eventLabel || edgeText || ''
      
      console.log('🔗 转换连线', {
        原始连线: edge,
        text类型: typeof edge.text,
        text值: edge.text,
        显示名称: displayName,
        事件代码: eventCode
      })
      
      return {
        id: edge.properties.transitionId || edge.id || `trans_${Date.now()}_${Math.random()}`,
        name: displayName,        // ✅ 使用中文显示名称
        from_state: edge.sourceNodeId,
        to_state: edge.targetNodeId,
        event: eventCode,         // ✅ 使用英文事件代码
        conditions: edge.properties.conditions || [],
        actions: edge.properties.actions || [],
        require_role: edge.properties.requireRole || '',
        description: edge.properties.description || '',
      }
    })

    // 验证数据
    const invalidStates = states.filter(s => !s.name || !s.code)
    if (invalidStates.length > 0) {
      console.error('❌ 无效的状态', invalidStates)
      window.$message.error('存在无效的状态节点，请确保所有节点都有名称')
      return
    }
    
    const payload = {
      ...workflowForm.value,
      states,
      transitions,
    }
    
    console.log('💾 准备保存工作流', payload)

    const workflowId = route.query.id as string
    if (workflowId) {
      // 更新
      await updateWorkflowDefinition(workflowId, payload)
      window.$message.success('保存成功')
    }
    else {
      // 创建
      await createWorkflowDefinition(payload)
      window.$message.success('创建成功')
      router.push('/order/workflow/designer/list')
    }
  }
  catch (error: any) {
    console.error('❌ 保存失败', error)
    window.$message.error(error.message || '保存失败')
  }
}

// 发布工作流
async function handlePublish() {
  await handleSave()
  // TODO: 调用激活API
  window.$message.success('发布成功')
}

// 预览
function handlePreview() {
  if (!lf)
    return

  const graphData = lf.getGraphData()
  console.log('工作流数据：', graphData)
  window.$message.info('请查看控制台')
}

// 返回列表
function handleBack() {
  router.push('/order/workflow/designer/list')
}

// 清理
onUnmounted(() => {
  if (lf) {
    try {
      // 清理事件监听
      lf.off('node:click')
      lf.off('edge:click')
      lf.off('blank:click')
      
      // 销毁实例
      lf.destroy()
    } catch (error) {
      // 忽略销毁错误（MiniMap 插件的已知问题）
      console.warn('LogicFlow destroy warning:', error)
    } finally {
      lf = null
    }
  }
})
</script>

<style scoped>
#logic-flow-container {
  width: 100%;
  height: 100%;
}

.minimap {
  position: absolute;
  bottom: 20px;
  right: 20px;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.node-item {
  padding: 12px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  cursor: move;
  display: flex;
  align-items: center;
  transition: all 0.3s;
  background: white;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.node-item:hover {
  border-color: #409EFF;
  background: #ecf5ff;
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.canvas-container {
  position: relative;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

#logic-flow-container {
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}
</style>

