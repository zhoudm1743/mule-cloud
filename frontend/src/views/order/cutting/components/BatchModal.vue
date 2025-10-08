<script setup lang="ts">
import { computed, h, reactive, ref } from 'vue'
import { NButton, NDataTable, NFormItemGridItem, NGrid, NInput, NInputNumber, NModal, NSelect, NSpace, NText } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { bulkCreateCuttingBatch, fetchCuttingBatchList, fetchOrderDetail } from '@/service/api/order'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'BatchModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '裁剪制菲',
    edit: '编辑批次',
    view: '查看批次',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')
const currentTask = ref<Api.Order.CuttingTaskInfo | null>(null)
const orderData = ref<Api.Order.OrderInfo | null>(null)
const colorOptions = ref<any[]>([])
const orderSizes = ref<string[]>([]) // 订单的所有尺码
const existingBatches = ref<Api.Order.CuttingBatchInfo[]>([]) // 已有的批次列表

// 批次行数据类型
interface BatchRow {
  color: string
  layer_count: number
  size_quantities: Record<string, number> // 尺码 -> 每层数量
}

// 表单数据
const formDefault = () => ({
  task_id: '',
  bed_no: '',
  start_bundle_no: 1,
  batch_rows: [] as BatchRow[],
})

const formModel = reactive(formDefault())

const rules: any = {
  bed_no: { required: true, message: '请输入床号', trigger: 'blur' },
  start_bundle_no: { required: true, type: 'number', message: '请输入起始扎号', trigger: 'blur' },
}

const formRef = ref()
const loadingSubmit = ref(false)

// 初始化订单数据
async function initializeOrderData(task: Api.Order.CuttingTaskInfo) {
  try {
    const { data } = await fetchOrderDetail(task.order_id)
    
    if (data?.order) {
      orderData.value = data.order
      
      // 从订单明细中提取所有已配置的颜色（去重）
      const configuredColors = [...new Set((data.order.items || []).map(item => item.color))]
      colorOptions.value = configuredColors.map(color => ({
        label: color,
        value: color,
      }))
      
      // 提取所有尺码（去重并保持顺序）
      orderSizes.value = [...new Set((data.order.items || []).map(item => item.size))]
    }
  }
  catch (error) {
    orderData.value = null
    colorOptions.value = []
    orderSizes.value = []
    window.$message.warning('无法获取订单信息')
  }
}

// 加载已有的批次列表
async function loadExistingBatches(taskId: string) {
  try {
    const { data } = await fetchCuttingBatchList({
      task_id: taskId,
      page: 1,
      page_size: 9999, // 获取所有批次
    })
    existingBatches.value = data?.batches || []
  }
  catch (error) {
    existingBatches.value = []
  }
}

// 获取某个颜色+尺码在订单中的配置数量
function getOrderItemQuantity(color: string, size: string): number {
  if (!orderData.value)
    return 0
  const item = orderData.value.items.find(i => i.color === color && i.size === size)
  return item?.quantity || 0
}

// 获取某个颜色+尺码已经裁剪的数量（从已有批次中统计）
function getCutQuantity(color: string, size: string): number {
  let total = 0
  existingBatches.value.forEach((batch) => {
    if (batch.color === color) {
      const sizeDetail = batch.size_details.find(s => s.size === size)
      if (sizeDetail) {
        total += sizeDetail.quantity * batch.layer_count
      }
    }
  })
  return total
}

// 获取某个颜色+尺码还可以裁剪的数量
function getRemainingQuantity(color: string, size: string): number {
  const orderQty = getOrderItemQuantity(color, size)
  const cutQty = getCutQuantity(color, size)
  return Math.max(0, orderQty - cutQty)
}

// 添加批次行
function addBatchRow() {
  const sizeQuantities: Record<string, number> = {}
  orderSizes.value.forEach((size) => {
    sizeQuantities[size] = 0
  })
  
  formModel.batch_rows.push({
    color: '',
    layer_count: 1,
    size_quantities: sizeQuantities,
  })
}

// 删除批次行
function removeBatchRow(index: number) {
  formModel.batch_rows.splice(index, 1)
}

// 计算某行的合计扎数（非零尺码数量）
function getRowBundleCount(row: BatchRow): number {
  return Object.values(row.size_quantities).filter(qty => qty > 0).length
}

// 计算某行的合计件数
function getRowTotalPieces(row: BatchRow): number {
  return Object.values(row.size_quantities).reduce((sum, qty) => sum + (qty * row.layer_count), 0)
}

// 批次表格列定义
const batchColumns = computed<DataTableColumns<BatchRow>>(() => {
  const columns: DataTableColumns<BatchRow> = [
    {
      title: '颜色',
      key: 'color',
      width: 110,
      render: (row: BatchRow, index: number) => {
        return h(NSelect, {
          value: row.color,
          options: colorOptions.value,
          placeholder: '选择',
          filterable: true,
          size: 'small',
          'onUpdate:value': (value: string) => {
            formModel.batch_rows[index].color = value
          },
        })
      },
    },
    {
      title: '拉布层数',
      key: 'layer_count',
      width: 90,
      render: (row: BatchRow, index: number) => {
        return h(NInputNumber, {
          value: row.layer_count,
          min: 1,
          placeholder: '1',
          size: 'small',
          style: { width: '100%' },
          'onUpdate:value': (value: number | null) => {
            formModel.batch_rows[index].layer_count = value || 1
          },
        })
      },
    },
  ]

  // 动态添加尺码列
  orderSizes.value.forEach((size) => {
    columns.push({
      title: size,
      key: `size_${size}`,
      width: 120,
      render: (row: BatchRow, index: number) => {
        // 获取当前颜色+尺码的订单数量和已裁数量
        const orderQty = getOrderItemQuantity(row.color, size)
        const cutQty = getCutQuantity(row.color, size)
        const remaining = getRemainingQuantity(row.color, size)
        const maxPerLayer = row.layer_count > 0 ? Math.floor(remaining / row.layer_count) : 0
        
        // 判断是否超裁或接近超裁
        const isOvercut = cutQty > orderQty
        const isNearLimit = cutQty >= orderQty * 0.9
        
        return h('div', { class: 'flex flex-col gap-1' }, [
          // 输入框
          h(NInputNumber, {
            value: row.size_quantities[size] || 0,
            min: 0,
            max: maxPerLayer,
            placeholder: '0',
            size: 'small',
            style: { width: '100%' },
            'onUpdate:value': (value: number | null) => {
              formModel.batch_rows[index].size_quantities[size] = value || 0
            },
          }),
          // 进度信息（仅当选择了颜色时显示）
          row.color ? h('div', {
            class: 'text-xs text-center',
            style: {
              color: isOvercut ? '#f56c6c' : isNearLimit ? '#e6a23c' : '#909399',
            },
          }, `${cutQty}/${orderQty} 剩${remaining}`) : null,
        ])
      },
    })
  })

  // 添加合计列
  columns.push(
    {
      title: '合计扎数',
      key: 'bundle_count',
      width: 80,
      render: (row: BatchRow) => {
        const count = getRowBundleCount(row)
        return h(NText, { depth: 3, class: 'font-medium' }, { default: () => count })
      },
    },
    {
      title: '合计件数',
      key: 'total_pieces',
      width: 80,
      render: (row: BatchRow) => {
        const count = getRowTotalPieces(row)
        return h(NText, { type: 'info', class: 'font-medium' }, { default: () => count })
      },
    },
    {
      title: '操作',
      key: 'actions',
      width: 60,
      fixed: 'right',
      render: (row: BatchRow, index: number) => {
        return h(
          NButton,
          {
            size: 'small',
            type: 'error',
            secondary: true,
            onClick: () => removeBatchRow(index),
          },
          { default: () => '删除' },
        )
      },
    },
  )

  return columns
})

async function openModal(type: ModalType, task?: Api.Order.CuttingTaskInfo) {
  showModal()
  modalType.value = type
  currentTask.value = task || null

  if (type === 'add' && task) {
    Object.assign(formModel, formDefault())
    formModel.task_id = task.id
    
    // 并行加载订单数据和已有批次
    await Promise.all([
      initializeOrderData(task),
      loadExistingBatches(task.id),
    ])
    
    // 自动添加一行空数据，方便用户直接填写
    addBatchRow()
  }
}

async function handleSubmit() {
  await formRef.value?.validate()

  // 验证至少有一行数据
  if (formModel.batch_rows.length === 0) {
    window.$message.error('请至少添加一行数据')
    return
  }

  // 验证每行数据
  for (let i = 0; i < formModel.batch_rows.length; i++) {
    const row = formModel.batch_rows[i]
    
    if (!row.color) {
      window.$message.error(`第 ${i + 1} 行：请选择颜色`)
      return
    }

    if (!row.layer_count || row.layer_count <= 0) {
      window.$message.error(`第 ${i + 1} 行：请输入拉布层数`)
      return
    }

    // 检查是否至少有一个尺码有数量
    const hasQuantity = Object.values(row.size_quantities).some(qty => qty > 0)
    if (!hasQuantity) {
      window.$message.error(`第 ${i + 1} 行：请至少填写一个尺码的数量`)
    return
    }

    // 验证数量约束
    for (const [size, qty] of Object.entries(row.size_quantities)) {
      if (qty > 0) {
        const totalCut = qty * row.layer_count
        const remaining = getRemainingQuantity(row.color, size)
        if (totalCut > remaining) {
          window.$message.error(`第 ${i + 1} 行：${row.color} - ${size} 裁剪数量超过剩余数量`)
          return
        }
      }
    }
  }

  loadingSubmit.value = true
  try {
    // 构造批量创建请求数据
    let bundleNo = formModel.start_bundle_no
    const batches: Api.Order.BatchItem[] = formModel.batch_rows.map((row) => {
      // 将 size_quantities 转换为 size_details 数组
      const sizeDetails = Object.entries(row.size_quantities)
        .filter(([_, qty]) => qty > 0)
        .map(([size, qty]) => ({ size, quantity: qty }))

      const batch: Api.Order.BatchItem = {
        bundle_no: String(bundleNo).padStart(2, '0'), // 补零，确保至少两位数
        color: row.color,
        layer_count: row.layer_count,
        size_details: sizeDetails,
      }

      bundleNo++
      return batch
    })

    // 批量创建批次（一次性提交）
    const { data } = await bulkCreateCuttingBatch({
      task_id: formModel.task_id,
      bed_no: formModel.bed_no,
      batches,
    })

    window.$message.success(`成功创建 ${data?.count || formModel.batch_rows.length} 个批次`)
    hideModal()
    emit('refresh')
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
  finally {
    loadingSubmit.value = false
  }
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="visible"
    :title="title"
    preset="card"
    class="w-1100px"
    :mask-closable="false"
  >
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="100"
      :disabled="modalType === 'view'"
    >
      <!-- 表单区域 -->
      <div>
        <!-- 基本信息 -->
        <NGrid :cols="2" :x-gap="18" class="mb-3">
          <NFormItemGridItem path="bed_no" label="设置床号">
          <NInput v-model:value="formModel.bed_no" placeholder="请输入床号" />
        </NFormItemGridItem>
          <NFormItemGridItem path="start_bundle_no" label="起始扎号">
          <NInputNumber
              v-model:value="formModel.start_bundle_no"
            :min="1"
              placeholder="1"
            class="w-full"
          />
        </NFormItemGridItem>
      </NGrid>

        <!-- 操作提示 -->
        <div class="mb-3 px-3 py-2 bg-blue-50 rounded text-sm text-gray-600">
          <span class="font-medium text-gray-700">操作提示：</span>
          点击下方"添加一行"按钮，选择颜色和拉布层数，填写各尺码数量
        </div>

        <!-- 批次数据表格 -->
        <NDataTable
          :columns="batchColumns"
          :data="formModel.batch_rows"
          :bordered="true"
          :single-line="false"
          size="small"
          :scroll-x="900"
          max-height="400"
          class="mb-3"
        />

        <!-- 添加按钮 -->
        <NButton dashed block @click="addBatchRow">
          + 添加一行
        </NButton>
      </div>
    </NForm>

    <template v-if="modalType !== 'view'" #footer>
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-500">
          <span v-if="formModel.batch_rows.length > 0">
            已添加 <span class="font-bold text-primary">{{ formModel.batch_rows.length }}</span> 行，将创建 <span class="font-bold text-primary">{{ formModel.batch_rows.length }}</span> 个批次
          </span>
          <span v-else class="text-orange-500">请至少添加一行数据</span>
        </div>
        <NSpace>
          <NButton @click="hideModal">
            取消
          </NButton>
          <NButton type="primary" :loading="loadingSubmit" :disabled="formModel.batch_rows.length === 0" @click="handleSubmit">
            保存
          </NButton>
        </NSpace>
      </div>
    </template>
  </NModal>
</template>
