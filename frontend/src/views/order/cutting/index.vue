<script setup lang="ts">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NInput, NSpace, NTag, NText } from 'naive-ui'
import { fetchCuttingTaskList } from '@/service/api/order'
import BatchModal from './components/BatchModal.vue'
import BatchDrawer from './components/BatchDrawer.vue'

defineOptions({ name: 'CuttingManagement' })

// 状态标签配置
const statusMap = {
  0: { label: '待裁剪', type: 'default' },
  1: { label: '裁剪中', type: 'info' },
  2: { label: '已完成', type: 'success' },
}

// 搜索参数
const initialSearchParams = {
  page: 1,
  page_size: 10,
  contract_no: '',
  style_no: '',
  status: undefined as number | undefined,
}
const searchParams = reactive({ ...initialSearchParams })

// 任务列表
const taskList = ref<Api.Order.CuttingTaskInfo[]>([])
const taskTotal = ref(0)
const loading = ref(false)

// 当前任务ID
const currentTaskId = ref('')

// 批次模态框
const batchModalRef = ref<InstanceType<typeof BatchModal> | null>(null)
// 批次抽屉
const batchDrawerRef = ref<InstanceType<typeof BatchDrawer> | null>(null)

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

// 查看批次
function handleViewBatches(task: Api.Order.CuttingTaskInfo) {
  ;(batchDrawerRef.value as any)?.open(task.id)
}

// 添加批次
function handleAddBatch(task: Api.Order.CuttingTaskInfo) {
  currentTaskId.value = task.id
  ;(batchModalRef.value as any)?.openModal('add', task)
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
    width: 200,
    render: (row: Api.Order.CuttingTaskInfo) => {
      const isOvercut = row.cut_pieces > row.total_pieces
      const progress = `${row.cut_pieces} / ${row.total_pieces}`
      
      if (isOvercut) {
        return h(NSpace, { align: 'center' }, {
          default: () => [
            h(NText, { type: 'error', strong: true }, { default: () => progress }),
            h(NText, { type: 'error', depth: 3, style: { fontSize: '12px' } }, { default: () => '(超量)' }),
          ],
        })
      }
      
      return progress
    },
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


// 搜索
function handleSearch() {
  searchParams.page = 1
  fetchTaskData()
}

// 重置
function handleReset() {
  Object.assign(searchParams, { ...initialSearchParams })
  fetchTaskData()
}

// 刷新任务和批次列表
function handleRefreshBatches() {
  fetchTaskData()
}

// 任务分页
function handlePageChange(page: number) {
  searchParams.page = page
  fetchTaskData()
}

// 任务分页大小
function handlePageSizeChange(pageSize: number) {
  searchParams.page_size = pageSize
  searchParams.page = 1
  fetchTaskData()
}

// 任务分页配置
const taskPagination = computed(() => ({
  page: searchParams.page,
  pageSize: searchParams.page_size,
  pageCount: Math.ceil(taskTotal.value / searchParams.page_size),
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
}))

onMounted(() => {
  fetchTaskData()
})
</script>

<template>
  <NSpace vertical size="large" class="p-4">
    <NCard title="裁剪任务管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchParams" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="合同号">
            <NInput
              v-model:value="searchParams.contract_no"
              placeholder="请输入合同号"
              clearable
              class="w-200px"
              @keyup.enter="handleSearch"
            />
          </n-form-item>
          <n-form-item label="款号">
            <NInput
              v-model:value="searchParams.style_no"
              placeholder="请输入款号"
              clearable
              class="w-200px"
              @keyup.enter="handleSearch"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="handleSearch">
              <template #icon>
                <nova-icon icon="carbon:search" :size="18" />
              </template>
              搜索
            </NButton>
            <NButton strong secondary @click="handleReset">
              <template #icon>
                <nova-icon icon="carbon:reset" :size="18" />
              </template>
              重置
            </NButton>
          </n-flex>
        </n-flex>
      </n-form>
    </NCard>

    <NCard :bordered="false" class="rounded-8px shadow-sm">
      <NDataTable
        :columns="taskColumns"
        :data="taskList"
        :loading="loading"
        :single-line="false"
      />
      <Pagination :count="taskTotal" :page="searchParams.page" :page-size="searchParams.page_size" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
    </NCard>

    <BatchModal ref="batchModalRef" @refresh="handleRefreshBatches" />
    <BatchDrawer ref="batchDrawerRef" @refresh="handleRefreshBatches" />
  </NSpace>
</template>
