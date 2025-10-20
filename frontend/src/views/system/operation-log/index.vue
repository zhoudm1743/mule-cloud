<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import { useBoolean } from '@/hooks'
import { fetchOperationLogList } from '@/service'
import { NButton, NSpace, NTag } from 'naive-ui'
import Pagination from '@/components/common/Pagination.vue'
import DetailDrawer from './components/DetailDrawer.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const detailDrawerRef = ref()
const tableData = ref<Api.OperationLog.OperationLogInfo[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const initialSearchForm = {
  username: '',
  resource: '',
  method: '',
  action: '',
  response_code: undefined as number | undefined,
  start_time: undefined as number | undefined,
  end_time: undefined as number | undefined,
  dateRange: null as [number, number] | null,
}
const searchForm = ref({ ...initialSearchForm })

// 日期范围变化
function handleDateRangeChange(value: [number, number] | null) {
  if (value) {
    searchForm.value.start_time = Math.floor(value[0] / 1000)
    searchForm.value.end_time = Math.floor(value[1] / 1000)
  } else {
    searchForm.value.start_time = undefined
    searchForm.value.end_time = undefined
  }
}

// 分页处理
function handlePageChange(p: number, ps: number) {
  page.value = p
  pageSize.value = ps
  fetchLogs()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchLogs()
}

// 查看详情
function handleViewDetail(id: string) {
  detailDrawerRef.value.open(id)
}

// HTTP方法选项
const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
]

// 操作类型选项
const actionOptions = [
  { label: 'create', value: 'create' },
  { label: 'update', value: 'update' },
  { label: 'delete', value: 'delete' },
  { label: 'read', value: 'read' },
  { label: 'login', value: 'login' },
  { label: 'logout', value: 'logout' },
]

// 获取HTTP方法的标签类型
function getMethodType(method: string): NaiveUI.ThemeColor {
  const methodMap: Record<string, NaiveUI.ThemeColor> = {
    'GET': 'info',
    'POST': 'success',
    'PUT': 'warning',
    'DELETE': 'error',
    'PATCH': 'warning',
  }
  return methodMap[method] || 'default'
}

// 获取状态码的标签类型
function getStatusType(code: number): NaiveUI.ThemeColor {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 300 && code < 400) return 'info'
  if (code >= 400 && code < 500) return 'warning'
  if (code >= 500) return 'error'
  return 'default'
}

// 格式化日期时间
function formatDateTime(dateStr: string): string {
  try {
    const date = new Date(dateStr)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return dateStr
  }
}

// 格式化耗时
function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

const columns: DataTableColumns<Api.OperationLog.OperationLogInfo> = [
  {
    title: '操作时间',
    key: 'created_at',
    width: 180,
    render: row => formatDateTime(row.created_at),
  },
  {
    title: '操作用户',
    key: 'username',
    width: 120,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: 'HTTP方法',
    key: 'method',
    width: 100,
    align: 'center',
    render: row => <NTag type={getMethodType(row.method)} size="small">{row.method}</NTag>,
  },
  {
    title: '请求路径',
    key: 'path',
    width: 300,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '资源',
    key: 'resource',
    width: 150,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'action',
    width: 100,
    align: 'center',
    render: row => <NTag type="primary" size="small">{row.action}</NTag>,
  },
  {
    title: '状态码',
    key: 'response_code',
    width: 100,
    align: 'center',
    render: row => <NTag type={getStatusType(row.response_code)} size="small">{row.response_code}</NTag>,
  },
  {
    title: '耗时',
    key: 'duration',
    width: 100,
    align: 'center',
    render: row => <span>{formatDuration(row.duration)}</span>,
  },
  {
    title: '客户端IP',
    key: 'ip',
    width: 140,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    align: 'center',
    key: 'actions',
    width: 100,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace justify="center">
          <NButton
            text
            type="primary"
            onClick={() => handleViewDetail(row.id)}
          >
            查看详情
          </NButton>
        </NSpace>
      )
    },
  },
]

onMounted(() => {
  fetchLogs()
})

async function fetchLogs() {
  startLoading()
  try {
    const params: Api.OperationLog.ListRequest = {
      page: page.value,
      page_size: pageSize.value,
      username: searchForm.value.username || undefined,
      resource: searchForm.value.resource || undefined,
      method: searchForm.value.method || undefined,
      action: searchForm.value.action || undefined,
      response_code: searchForm.value.response_code,
      start_time: searchForm.value.start_time,
      end_time: searchForm.value.end_time,
    }
    
    const res = await fetchOperationLogList(params)
    if (res.data) {
      tableData.value = res.data.list || []
      total.value = res.data.total || 0
    }
  }
  catch (e) {
    console.error('[Fetch Operation Logs Error]:', e)
  }
  finally {
    endLoading()
  }
}
</script>

<template>
  <NSpace vertical size="large">
    <NCard :bordered="false" class="rounded-8px shadow-sm">
      <template #header>
        <span class="text-lg font-semibold">操作日志</span>
      </template>
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex :wrap="false">
          <n-form-item label="操作用户">
            <n-input
              v-model:value="searchForm.username"
              placeholder="请输入用户名"
              clearable
              class="w-150px"
              @keyup.enter="fetchLogs"
            />
          </n-form-item>
          <n-form-item label="资源名称">
            <n-input
              v-model:value="searchForm.resource"
              placeholder="请输入资源名称"
              clearable
              class="w-150px"
              @keyup.enter="fetchLogs"
            />
          </n-form-item>
          <n-form-item label="HTTP方法">
            <n-select
              v-model:value="searchForm.method"
              placeholder="请选择"
              clearable
              class="w-120px"
              :options="methodOptions"
            />
          </n-form-item>
          <n-form-item label="操作类型">
            <n-select
              v-model:value="searchForm.action"
              placeholder="请选择"
              clearable
              class="w-120px"
              :options="actionOptions"
            />
          </n-form-item>
          <n-form-item label="时间范围">
            <n-date-picker
              v-model:value="searchForm.dateRange"
              type="datetimerange"
              clearable
              class="w-350px"
              @update:value="handleDateRangeChange"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchLogs">
              <template #icon>
                <nova-icon icon="carbon:search" :size="18" />
              </template>
              搜索
            </NButton>
            <NButton strong secondary @click="handleResetSearch">
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
      <NSpace vertical size="large">
        <n-data-table
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1600"
        />
        <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" />
      </NSpace>
    </NCard>

    <!-- 详情抽屉 -->
    <DetailDrawer ref="detailDrawerRef" />
  </NSpace>
</template>

<style scoped>
:deep(.n-data-table) {
  font-size: 13px;
}
</style>

