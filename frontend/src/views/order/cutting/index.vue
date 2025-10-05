<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NInput, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { deleteCuttingBatch, fetchCuttingBatchList, fetchCuttingTaskList, printCuttingBatch } from '@/service/api/order'
import BatchModal from './components/BatchModal.vue'

defineOptions({ name: 'CuttingManagement' })

// 状态标签配置
const statusMap = {
  0: { label: '待裁剪', type: 'default' },
  1: { label: '裁剪中', type: 'info' },
  2: { label: '已完成', type: 'success' },
}

// 搜索参数
const searchParams = reactive({
  page: 1,
  page_size: 10,
  contract_no: '',
  style_no: '',
  status: undefined as number | undefined,
})

// 任务列表
const taskList = ref<Api.Order.CuttingTaskInfo[]>([])
const taskTotal = ref(0)
const loading = ref(false)

// 批次列表
const batchList = ref<Api.Order.CuttingBatchInfo[]>([])
const batchTotal = ref(0)
const batchLoading = ref(false)
const currentTaskId = ref('')

// 批次模态框
const batchModalRef = ref<InstanceType<typeof BatchModal> | null>(null)

// 获取任务列表
async function fetchTaskData() {
  loading.value = true
  try {
    const { data } = await fetchCuttingTaskList(searchParams)
    if (data) {
      taskList.value = data.tasks || []
      taskTotal.value = data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取任务列表失败')
  }
  finally {
    loading.value = false
  }
}

// 获取批次列表
async function fetchBatchData(taskId: string) {
  batchLoading.value = true
  currentTaskId.value = taskId
  try {
    const { data } = await fetchCuttingBatchList({
      task_id: taskId,
      page: 1,
      page_size: 100,
    })
    if (data) {
      batchList.value = data.batches || []
      batchTotal.value = data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取批次列表失败')
  }
  finally {
    batchLoading.value = false
  }
}

// 查看批次
function handleViewBatches(task: Api.Order.CuttingTaskInfo) {
  fetchBatchData(task.id)
}

// 添加批次
function handleAddBatch(task: Api.Order.CuttingTaskInfo) {
  currentTaskId.value = task.id
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(batchModalRef.value as any)?.openModal('add', task)
}

// 删除批次
async function handleDeleteBatch(batch: Api.Order.CuttingBatchInfo) {
  try {
    await deleteCuttingBatch(batch.id)
    window.$message.success('删除成功')
    fetchBatchData(currentTaskId.value)
    fetchTaskData() // 更新任务进度
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

// 打印批次
async function handlePrintBatch(batch: Api.Order.CuttingBatchInfo) {
  try {
    const { data } = await printCuttingBatch(batch.id)
    if (data) {
      window.$message.success(`打印成功，已打印 ${data.batch.print_count} 次`)
      // 这里应该触发实际的打印操作，可以使用浏览器打印API
      // 或者打开一个打印预览页面
      printQRCode(data.batch)
      fetchBatchData(currentTaskId.value)
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '打印失败')
  }
}

// 打印二维码
function printQRCode(batch: Api.Order.CuttingBatchInfo) {
  // 创建打印窗口
  const printWindow = window.open('', '_blank')
  if (!printWindow)
    return

  // 生成打印内容
  const printContent = `
    <html>
      <head>
        <title>裁剪批次 - ${batch.bundle_no}</title>
        <style>
          body { font-family: Arial, sans-serif; padding: 20px; }
          .header { text-align: center; margin-bottom: 20px; }
          .info { margin-bottom: 10px; }
          .qr-code { text-align: center; margin: 20px 0; }
          table { width: 100%; border-collapse: collapse; margin-top: 10px; }
          table, th, td { border: 1px solid black; }
          th, td { padding: 8px; text-align: left; }
          @media print {
            button { display: none; }
          }
        </style>
      </head>
      <body>
        <div class="header">
          <h2>裁剪批次标签</h2>
        </div>
        <div class="info">
          <p><strong>合同号：</strong>${batch.contract_no}</p>
          <p><strong>款号：</strong>${batch.style_no}</p>
          <p><strong>床号：</strong>${batch.bed_no}</p>
          <p><strong>扎号：</strong>${batch.bundle_no}</p>
          <p><strong>颜色：</strong>${batch.color}</p>
          <p><strong>拉布层数：</strong>${batch.layer_count}</p>
          <p><strong>总件数：</strong>${batch.total_pieces}</p>
        </div>
        <table>
          <thead>
            <tr>
              <th>尺码</th>
              <th>数量</th>
            </tr>
          </thead>
          <tbody>
            ${batch.size_details.map(sd => `
              <tr>
                <td>${sd.size}</td>
                <td>${sd.quantity}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
        <div class="qr-code">
          <p>二维码内容：${batch.qr_code || ''}</p>
          <p><small>扫码可查看详细信息</small></p>
        </div>
        <button onclick="window.print()">打印</button>
      </body>
    </html>
  `

  printWindow.document.write(printContent)
  printWindow.document.close()
}

// 任务表格列
const taskColumns = computed(() => [
  { title: '合同号', key: 'contract_no', width: 150 },
  { title: '款号', key: 'style_no', width: 120 },
  { title: '款名', key: 'style_name', width: 150 },
  { title: '客户', key: 'customer_name', width: 150 },
  {
    title: '进度',
    key: 'progress',
    width: 150,
    render: (row: Api.Order.CuttingTaskInfo) => `${row.cut_pieces} / ${row.total_pieces}`,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: Api.Order.CuttingTaskInfo) => {
      const config = statusMap[row.status as keyof typeof statusMap]
      return h(NTag, { type: config.type as any }, { default: () => config.label })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 250,
    render: (row: Api.Order.CuttingTaskInfo) => {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => handleViewBatches(row) }, { default: () => '查看批次' }),
          h(NButton, { size: 'small', type: 'primary', onClick: () => handleAddBatch(row) }, { default: () => '制菲' }),
        ],
      })
    },
  },
])

// 批次表格列
const batchColumns = computed(() => [
  { title: '床号', key: 'bed_no', width: 100 },
  { title: '扎号', key: 'bundle_no', width: 100 },
  { title: '颜色', key: 'color', width: 100 },
  { title: '拉布层数', key: 'layer_count', width: 100 },
  {
    title: '尺码明细',
    key: 'size_details',
    width: 200,
    render: (row: Api.Order.CuttingBatchInfo) => {
      return row.size_details.map(sd => `${sd.size}:${sd.quantity}`).join(', ')
    },
  },
  { title: '总件数', key: 'total_pieces', width: 100 },
  {
    title: '打印次数',
    key: 'print_count',
    width: 100,
    render: (row: Api.Order.CuttingBatchInfo) => row.print_count || 0,
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    render: (row: Api.Order.CuttingBatchInfo) => {
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', type: 'primary', onClick: () => handlePrintBatch(row) }, { default: () => '打印' }),
          h(NPopconfirm, { onPositiveClick: () => handleDeleteBatch(row) }, {
            default: () => '确定删除这个批次吗？',
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
          }),
        ],
      })
    },
  },
])

// 搜索
function handleSearch() {
  searchParams.page = 1
  fetchTaskData()
}

// 重置
function handleReset() {
  Object.assign(searchParams, {
    page: 1,
    page_size: 10,
    contract_no: '',
    style_no: '',
    status: undefined,
  })
  fetchTaskData()
}

// 刷新批次列表
function handleRefreshBatches() {
  if (currentTaskId.value) {
    fetchBatchData(currentTaskId.value)
    fetchTaskData() // 同时更新任务进度
  }
}

onMounted(() => {
  fetchTaskData()
})
</script>

<template>
  <div class="p-4">
    <NCard title="裁剪任务管理" :bordered="false" class="mb-4">
      <div class="mb-4 flex gap-4">
        <NInput
          v-model:value="searchParams.contract_no"
          placeholder="合同号"
          clearable
          class="w-200px"
        />
        <NInput
          v-model:value="searchParams.style_no"
          placeholder="款号"
          clearable
          class="w-200px"
        />
        <NButton type="primary" @click="handleSearch">
          搜索
        </NButton>
        <NButton @click="handleReset">
          重置
        </NButton>
      </div>

      <NDataTable
        :columns="taskColumns"
        :data="taskList"
        :loading="loading"
        :pagination="{
          page: searchParams.page,
          pageSize: searchParams.page_size,
          itemCount: taskTotal,
          showSizePicker: true,
          pageSizes: [10, 20, 50, 100],
          onChange: (page: number) => {
            searchParams.page = page
            fetchTaskData()
          },
          onUpdatePageSize: (pageSize: number) => {
            searchParams.page_size = pageSize
            searchParams.page = 1
            fetchTaskData()
          },
        }"
        :single-line="false"
      />
    </NCard>

    <NCard v-if="currentTaskId" title="裁剪批次列表" :bordered="false">
      <template #header-extra>
        <NButton @click="handleRefreshBatches">
          刷新
        </NButton>
      </template>

      <NDataTable
        :columns="batchColumns"
        :data="batchList"
        :loading="batchLoading"
        :single-line="false"
      />
    </NCard>

    <BatchModal ref="batchModalRef" @refresh="handleRefreshBatches" />
  </div>
</template>
