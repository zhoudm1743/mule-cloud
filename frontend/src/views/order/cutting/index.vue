<script setup lang="tsx">
import { computed, h, onMounted, reactive, ref } from 'vue'
import { NButton, NCard, NDataTable, NInput, NInputNumber, NModal, NForm, NFormItem, NSpace, NTabs, NTabPane, NTag, NText, NProgress } from 'naive-ui'
import { fetchCuttingTaskList, fetchCuttingPieceList, updateCuttingPieceProgress } from '@/service/api/order'
import { useBoolean } from '@/hooks'
import BatchModal from './components/BatchModal.vue'
import BatchDrawer from './components/BatchDrawer.vue'

defineOptions({ name: 'CuttingManagement' })

// Tab 切换
const activeTab = ref('tasks')

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
const { bool: taskLoading, setTrue: startTaskLoading, setFalse: endTaskLoading } = useBoolean(false)

// 裁片列表
const pieceList = ref<Api.Order.CuttingPieceInfo[]>([])
const pieceTotal = ref(0)
const { bool: pieceLoading, setTrue: startPieceLoading, setFalse: endPieceLoading } = useBoolean(false)

// 裁片搜索参数
const initialPieceSearchParams = {
  page: 1,
  page_size: 20,
  contract_no: '',
  bed_no: '',
  bundle_no: '',
}
const pieceSearchParams = reactive({ ...initialPieceSearchParams })

// 进度更新模态框
const showProgressModal = ref(false)
const currentPiece = ref<Api.Order.CuttingPieceInfo | null>(null)
const progressForm = ref({
  progress: 0,
})

// 当前任务ID
const currentTaskId = ref('')

// 批次模态框
const batchModalRef = ref<InstanceType<typeof BatchModal> | null>(null)
// 批次抽屉
const batchDrawerRef = ref<InstanceType<typeof BatchDrawer> | null>(null)

// 获取任务列表
async function fetchTaskData() {
  startTaskLoading()
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
    endTaskLoading()
  }
}

// 获取裁片列表
async function fetchPieceData() {
  startPieceLoading()
  try {
    const { data } = await fetchCuttingPieceList(pieceSearchParams)
    if (data) {
      pieceList.value = data.pieces || []
      pieceTotal.value = data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取裁片列表失败')
  }
  finally {
    endPieceLoading()
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
function handlePageChange(p: number, ps: number) {
  searchParams.page = p
  searchParams.page_size = ps
  fetchTaskData()
}

// 裁片表格列
const pieceColumns = computed(() => [
  { title: '合同号', key: 'contract_no', width: 140, fixed: 'left' as const },
  { title: '款号', key: 'style_no', width: 100 },
  { title: '床号', key: 'bed_no', width: 100 },
  { title: '扎号', key: 'bundle_no', width: 100 },
  { title: '颜色', key: 'color', width: 100 },
  { title: '尺码', key: 'size', width: 80 },
  { title: '数量', key: 'quantity', width: 100 },
  {
    title: '进度',
    key: 'progress',
    width: 250,
    render: (row: Api.Order.CuttingPieceInfo) => {
      const progressPercent = row.total_process > 0
        ? Math.round((row.progress / row.total_process) * 100)
        : 0
      
      return h(NSpace, { vertical: true, size: 4 }, {
        default: () => [
          h(NProgress, {
            type: 'line',
            percentage: progressPercent,
            status: progressPercent === 100 ? 'success' : 'default',
            showIndicator: false,
          }),
          h(NText, { depth: 3, style: { fontSize: '12px' } }, {
            default: () => `${row.progress} / ${row.total_process} 道工序 (${progressPercent}%)`,
          }),
        ],
      })
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row: Api.Order.CuttingPieceInfo) => {
      const progressPercent = row.total_process > 0
        ? Math.round((row.progress / row.total_process) * 100)
        : 0
      
      if (progressPercent === 100) {
        return h(NTag, { type: 'success' }, { default: () => '已完成' })
      }
      else if (progressPercent > 0) {
        return h(NTag, { type: 'warning' }, { default: () => '进行中' })
      }
      else {
        return h(NTag, { type: 'default' }, { default: () => '未开始' })
      }
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: (row: Api.Order.CuttingPieceInfo) => {
      return new Date(row.created_at * 1000).toLocaleString('zh-CN')
    },
  },
  // {
  //   title: '操作',
  //   key: 'actions',
  //   width: 120,
  //   fixed: 'right' as const,
  //   render: (row: Api.Order.CuttingPieceInfo) => {
  //     return h(NSpace, {}, {
  //       default: () => [
  //         h(
  //           NButton,
  //           {
  //             size: 'small',
  //             type: 'primary',
  //             onClick: () => openProgressModal(row),
  //           },
  //           { default: () => '更新进度' },
  //         ),
  //       ],
  //     })
  //   },
  // },
])

// 裁片搜索
function handlePieceSearch() {
  pieceSearchParams.page = 1
  fetchPieceData()
}

// 裁片重置
function handlePieceReset() {
  Object.assign(pieceSearchParams, { ...initialPieceSearchParams })
  fetchPieceData()
}

// 裁片分页
function handlePiecePageChange(p: number, ps: number) {
  pieceSearchParams.page = p
  pieceSearchParams.page_size = ps
  fetchPieceData()
}

// 打开进度更新模态框
function openProgressModal(piece: Api.Order.CuttingPieceInfo) {
  currentPiece.value = piece
  progressForm.value.progress = piece.progress
  showProgressModal.value = true
}

// 更新进度
async function handleUpdateProgress() {
  if (!currentPiece.value) return
  
  try {
    await updateCuttingPieceProgress(currentPiece.value.id, {
      progress: progressForm.value.progress,
    })
    window.$message.success('更新成功')
    showProgressModal.value = false
    fetchPieceData()
  }
  catch (error: any) {
    window.$message.error(error.message || '更新失败')
  }
}

// Tab 切换时加载数据
function handleTabChange(value: string) {
  activeTab.value = value
  if (value === 'tasks' && taskList.value.length === 0) {
    fetchTaskData()
  }
  else if (value === 'pieces' && pieceList.value.length === 0) {
    fetchPieceData()
  }
}

onMounted(() => {
  fetchTaskData()
})
</script>

<template>
  <NSpace vertical size="large" class="p-4">
    <NCard title="裁剪管理" :bordered="false" class="rounded-8px shadow-sm">
      <NTabs v-model:value="activeTab" type="line" @update:value="handleTabChange">
        <!-- 裁剪任务 Tab -->
        <NTabPane name="tasks" tab="裁剪任务">
          <NSpace vertical size="large">
            <!-- 搜索栏 -->
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

            <!-- 任务表格 -->
            <NDataTable
              :columns="taskColumns"
              :data="taskList"
              :loading="taskLoading"
              :single-line="false"
            />
            <Pagination 
              :count="taskTotal" 
              @change="handlePageChange" 
            />
          </NSpace>
        </NTabPane>

        <!-- 裁片监控 Tab -->
        <NTabPane name="pieces" tab="裁片监控">
          <NSpace vertical size="large">
            <!-- 搜索栏 -->
            <n-form :model="pieceSearchParams" label-placement="left" inline :show-feedback="false">
              <n-flex>
                <n-form-item label="合同号">
                  <NInput
                    v-model:value="pieceSearchParams.contract_no"
                    placeholder="请输入合同号"
                    clearable
                    class="w-200px"
                    @keyup.enter="handlePieceSearch"
                  />
                </n-form-item>
                <n-form-item label="床号">
                  <NInput
                    v-model:value="pieceSearchParams.bed_no"
                    placeholder="请输入床号"
                    clearable
                    class="w-150px"
                    @keyup.enter="handlePieceSearch"
                  />
                </n-form-item>
                <n-form-item label="扎号">
                  <NInput
                    v-model:value="pieceSearchParams.bundle_no"
                    placeholder="请输入扎号"
                    clearable
                    class="w-150px"
                    @keyup.enter="handlePieceSearch"
                  />
                </n-form-item>
                <n-flex class="ml-auto">
                  <NButton type="primary" @click="handlePieceSearch">
                    <template #icon>
                      <nova-icon icon="carbon:search" :size="18" />
                    </template>
                    搜索
                  </NButton>
                  <NButton strong secondary @click="handlePieceReset">
                    <template #icon>
                      <nova-icon icon="carbon:reset" :size="18" />
                    </template>
                    重置
                  </NButton>
                </n-flex>
              </n-flex>
            </n-form>

            <!-- 裁片表格 -->
            <NDataTable
              :columns="pieceColumns"
              :data="pieceList"
              :loading="pieceLoading"
              :scroll-x="1500"
              :single-line="false"
            />
            <Pagination 
              :count="pieceTotal" 
              @change="handlePiecePageChange" 
            />
          </NSpace>
        </NTabPane>
      </NTabs>
    </NCard>

    <!-- 批次模态框 -->
    <BatchModal ref="batchModalRef" @refresh="handleRefreshBatches" />
    <!-- 批次抽屉 -->
    <BatchDrawer ref="batchDrawerRef" @refresh="handleRefreshBatches" />

    <!-- 进度更新模态框 -->
    <NModal
      v-model:show="showProgressModal"
      preset="dialog"
      title="更新裁片进度"
      positive-text="确定"
      negative-text="取消"
      @positive-click="handleUpdateProgress"
    >
      <NForm v-if="currentPiece" label-placement="left" label-width="100px" class="mt-4">
        <NFormItem label="合同号">
          <NText>{{ currentPiece.contract_no }}</NText>
        </NFormItem>
        <NFormItem label="扎号">
          <NText>{{ currentPiece.bundle_no }}</NText>
        </NFormItem>
        <NFormItem label="颜色/尺码">
          <NText>{{ currentPiece.color }} / {{ currentPiece.size }}</NText>
        </NFormItem>
        <NFormItem label="总工序数">
          <NText>{{ currentPiece.total_process }} 道</NText>
        </NFormItem>
        <NFormItem label="当前进度" required>
          <NInputNumber
            v-model:value="progressForm.progress"
            :min="0"
            :max="currentPiece.total_process"
            :step="1"
            placeholder="请输入完成的工序数"
            style="width: 100%"
          />
        </NFormItem>
        <NFormItem>
          <NText depth="3" style="font-size: 12px">
            提示：输入已完成的工序数量（0 ~ {{ currentPiece.total_process }}）
          </NText>
        </NFormItem>
      </NForm>
    </NModal>
  </NSpace>
</template>
