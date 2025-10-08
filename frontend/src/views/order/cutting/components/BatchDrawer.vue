<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NCard, NDataTable, NDrawer, NDrawerContent, NInput, NPopconfirm, NSpace } from 'naive-ui'
import { useBoolean } from '@/hooks'
import { batchPrintCuttingBatches, deleteCuttingBatch, fetchCuttingBatchList, printCuttingBatch } from '@/service/api/order'
import PrintModal from './PrintModal.vue'
import Pagination from '@/components/common/Pagination.vue'

const { bool: visible, setTrue: openDrawer, setFalse: closeDrawer } = useBoolean(false)
const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const emit = defineEmits<{
  refresh: []
}>()

const printModalRef = ref<InstanceType<typeof PrintModal> | null>(null)
const tableData = ref<Api.Order.CuttingBatchInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentTaskId = ref('')

// 搜索表单
const searchForm = ref({
  bed_no: '',
  bundle_no: '',
  color: '',
})

// 分页处理函数
function handlePageChange(p: number) {
  page.value = p
  fetchData()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchData()
}

async function fetchData() {
  if (!currentTaskId.value)
    return

  startLoading()
  try {
    const res = await fetchCuttingBatchList({
      task_id: currentTaskId.value,
      page: page.value,
      page_size: pageSize.value,
      bed_no: searchForm.value.bed_no || undefined,
      bundle_no: searchForm.value.bundle_no || undefined,
      // color: searchForm.value.color || undefined,
    })
    if (res.data) {
      tableData.value = res.data.batches || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取批次列表失败')
  }
  finally {
    endLoading()
  }
}

// 打开抽屉
function open(taskId: string) {
  currentTaskId.value = taskId
  page.value = 1
  searchForm.value = {
    bed_no: '',
    bundle_no: '',
    color: '',
  }
  openDrawer()
  fetchData()
}

// 搜索
function handleSearch() {
  page.value = 1
  fetchData()
}

// 重置
function handleReset() {
  searchForm.value = {
    bed_no: '',
    bundle_no: '',
    color: '',
  }
  page.value = 1
  fetchData()
}

// 删除批次
async function deleteHandler(id: string) {
  try {
    await deleteCuttingBatch(id)
    window.$message.success('删除成功')
    fetchData()
    emit('refresh')
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

// 打印单个批次
async function handlePrintBatch(batch: Api.Order.CuttingBatchInfo) {
  try {
    const { data } = await printCuttingBatch(batch.id)
    if (data) {
      (printModalRef.value as any)?.openModal(data.batch)
      fetchData()
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '打印失败')
  }
}

// 批量打印（勾选的）
async function handleBatchPrint() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要打印的批次')
    return
  }

  try {
    // 使用批量打印API，一次性处理所有批次
    const { data } = await batchPrintCuttingBatches(checkedRowKeys.value)
    
    if (data?.batches && data.batches.length > 0) {
      (printModalRef.value as any)?.openModal(data.batches)
      window.$message.success(`已生成 ${data.count} 个批次的打印预览`)
    }
    else {
      window.$message.warning('没有可打印的批次')
    }

    checkedRowKeys.value = []
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量打印失败')
  }
}

// 打印所有（不分页，不需勾选）
async function handlePrintAll() {
  if (!currentTaskId.value) {
    window.$message.warning('无效的任务ID')
    return
  }

  try {
    // 先获取所有批次（不分页）
    const res = await fetchCuttingBatchList({
      task_id: currentTaskId.value,
      page: 1,
      page_size: 9999, // 获取所有批次
      bed_no: searchForm.value.bed_no || undefined,
      bundle_no: searchForm.value.bundle_no || undefined,
    })

    if (!res.data?.batches || res.data.batches.length === 0) {
      window.$message.warning('没有可打印的批次')
      return
    }

    const allIds = res.data.batches.map(batch => batch.id)
    
    // 批量打印
    const { data } = await batchPrintCuttingBatches(allIds)
    
    if (data?.batches && data.batches.length > 0) {
      (printModalRef.value as any)?.openModal(data.batches)
      window.$message.success(`已生成 ${data.count} 个批次的打印预览`)
    }

    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '打印失败')
  }
}

// 刷新
function handleRefresh() {
  fetchData()
  emit('refresh')
}

const columns: DataTableColumns<Api.Order.CuttingBatchInfo> = [
  {
    type: 'selection',
    width: 50,
  },
  { title: '床号', key: 'bed_no', width: 100 },
  { title: '扎号', key: 'bundle_no', width: 100 },
  { title: '颜色', key: 'color', width: 120 },
  { title: '拉布层数', key: 'layer_count', width: 100 },
  {
    title: '尺码明细',
    key: 'size_details',
    width: 200,
    render: row => row.size_details.map(sd => `${sd.size}:${sd.quantity}`).join(', '),
  },
  { title: '总件数', key: 'total_pieces', width: 100 },
  {
    title: '打印次数',
    key: 'print_count',
    width: 100,
    render: row => row.print_count || 0,
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace>
          <NButton size="small" type="primary" onClick={() => handlePrintBatch(row)}>
            打印
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteHandler(row.id)}>
            {{
              default: () => '确定删除这个批次吗？',
              trigger: () => (
                <NButton size="small" type="error">
                  删除
                </NButton>
              ),
            }}
          </NPopconfirm>
        </NSpace>
      )
    },
  },
]

defineExpose({ open })
</script>

<template>
  <NDrawer v-model:show="visible" :width="1200" placement="right">
    <NDrawerContent title="裁剪批次列表" closable>
      <NSpace vertical size="large">
        <!-- 搜索区域 -->
        <NCard :bordered="false" class="rounded-8px shadow-sm">
          <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
            <n-flex>
              <n-form-item label="床号">
                <NInput
                  v-model:value="searchForm.bed_no"
                  placeholder="请输入床号"
                  clearable
                  class="w-150px"
                  @keyup.enter="handleSearch"
                />
              </n-form-item>
              <n-form-item label="扎号">
                <NInput
                  v-model:value="searchForm.bundle_no"
                  placeholder="请输入扎号"
                  clearable
                  class="w-150px"
                  @keyup.enter="handleSearch"
                />
              </n-form-item>
              <n-form-item label="颜色">
                <NInput
                  v-model:value="searchForm.color"
                  placeholder="请输入颜色"
                  clearable
                  class="w-150px"
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

        <!-- 操作按钮和表格 -->
        <NCard :bordered="false" class="rounded-8px shadow-sm">
          <NSpace vertical size="large">
            <div class="flex gap-4">
              <NButton
                type="primary"
                @click="handlePrintAll"
              >
                <template #icon>
                  <nova-icon icon="carbon:printer" :size="18" />
                </template>
                打印所有
              </NButton>
              <NButton
                type="info"
                :disabled="checkedRowKeys.length === 0"
                @click="handleBatchPrint"
              >
                <template #icon>
                  <nova-icon icon="carbon:checkmark-outline" :size="18" />
                </template>
                打印已勾选 ({{ checkedRowKeys.length }})
              </NButton>
            </div>

            <NDataTable
              :columns="columns"
              :data="tableData"
              :loading="loading"
              :scroll-x="1200"
              :row-key="(row: Api.Order.CuttingBatchInfo) => row.id"
              v-model:checked-row-keys="checkedRowKeys"
            />
            <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
          </NSpace>
        </NCard>
      </NSpace>
    </NDrawerContent>
  </NDrawer>

  <PrintModal ref="printModalRef" />
</template>

