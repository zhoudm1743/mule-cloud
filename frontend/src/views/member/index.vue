<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'
import axios from 'axios'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<any[]>([])
const checkedRowKeys = ref<string[]>([])
const total = ref(0)

// 搜索表单
const initialSearchForm = {
  name: '',
  job_number: '',
  department: '',
  status: undefined as string | undefined,
}
const searchForm = ref({ ...initialSearchForm })

// 状态选项
const statusOptions = [
  { label: '在职', value: 'active' },
  { label: '试用期', value: 'probation' },
  { label: '离职', value: 'inactive' },
]

// 获取员工列表
async function fetchData(page = 1, pageSize = 10) {
  startLoading()
  try {
    const res = await axios.get('/api/miniapp/member/list', {
      params: {
        page,
        page_size: pageSize,
        name: searchForm.value.name || undefined,
        job_number: searchForm.value.job_number || undefined,
        department: searchForm.value.department || undefined,
        status: searchForm.value.status,
      },
    })
    if (res.data?.data) {
      tableData.value = res.data.data.list || []
      total.value = res.data.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取员工列表失败')
  }
  finally {
    endLoading()
  }
}

// 分页处理
function handlePageChange(p: number, ps: number) {
  fetchData(p, ps)
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  fetchData()
}

// 删除员工
async function deleteHandler(id: string) {
  try {
    await axios.delete(`/api/miniapp/member/${id}`)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

// 批量删除
async function batchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的员工')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await axios.delete(`/api/miniapp/member/${id}`)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量删除失败')
  }
}

// 导出
async function handleExport() {
  try {
    window.$message.info('正在导出...')
    const res = await axios.get('/api/miniapp/member/export', {
      responseType: 'blob',
    })
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `员工数据_${new Date().toLocaleDateString()}.csv`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.$message.success('导出成功')
  }
  catch (error: any) {
    window.$message.error(error.message || '导出失败')
  }
}

// 导入
function handleImport() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.csv,.xlsx,.xls'
  input.onchange = async (e: any) => {
    const file = e.target.files[0]
    if (!file)
      return

    const formData = new FormData()
    formData.append('file', file)

    try {
      window.$message.info('正在导入...')
      const res = await axios.post('/api/miniapp/member/import', formData)
      window.$message.success(`导入成功：${res.data.data.success}条，失败：${res.data.data.failed}条`)
      fetchData()
    }
    catch (error: any) {
      window.$message.error(error.message || '导入失败')
    }
  }
  input.click()
}

const columns: DataTableColumns<any> = [
  {
    type: 'selection',
    width: 30,
  },
  {
    title: '工号',
    key: 'job_number',
    width: 120,
    render: row => <CopyText value={row.job_number} />,
  },
  {
    title: '姓名',
    key: 'name',
    width: 100,
  },
  {
    title: '性别',
    key: 'gender',
    width: 60,
    render: (row) => {
      const genderMap: Record<number, string> = { 0: '未知', 1: '男', 2: '女' }
      return genderMap[row.gender] || '未知'
    },
  },
  {
    title: '手机号',
    key: 'phone',
    width: 130,
    render: row => <CopyText value={row.phone} />,
  },
  {
    title: '部门',
    key: 'department',
    width: 120,
  },
  {
    title: '岗位',
    key: 'position',
    width: 120,
  },
  {
    title: '车间',
    key: 'workshop',
    width: 100,
  },
  {
    title: '班组',
    key: 'team',
    width: 100,
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const statusMap: Record<string, { label: string, type: any }> = {
        active: { label: '在职', type: 'success' },
        probation: { label: '试用期', type: 'warning' },
        inactive: { label: '离职', type: 'default' },
      }
      const status = statusMap[row.status] || { label: row.status, type: 'default' }
      return <NTag type={status.type}>{status.label}</NTag>
    },
  },
  {
    title: '入职日期',
    key: 'employed_at',
    width: 120,
    render: row => row.employed_at ? new Date(row.employed_at * 1000).toLocaleDateString() : '-',
  },
  {
    title: '工龄',
    key: 'work_years',
    width: 100,
    render: row => row.work_years ? `${row.work_years}年${row.work_months || 0}月` : '-',
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace>
          <NButton
            size={'small'}
            onClick={() => tableModalRef.value?.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteHandler(row.id)}>
            {{
              default: () => '确定删除吗？',
              trigger: () => <NButton size={'small'} type={'error'}>删除</NButton>,
            }}
          </NPopconfirm>
        </NSpace>
      )
    },
  },
]

onMounted(() => {
  fetchData()
})
</script>

<template>
  <NSpace vertical size="large">
    <NCard title="员工管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="姓名">
            <NInput
              v-model:value="searchForm.name"
              placeholder="请输入姓名"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-form-item label="工号">
            <NInput
              v-model:value="searchForm.job_number"
              placeholder="请输入工号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-form-item label="部门">
            <NInput
              v-model:value="searchForm.department"
              placeholder="请输入部门"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-form-item label="状态">
            <NSelect
              v-model:value="searchForm.status"
              placeholder="请选择"
              clearable
              class="w-120px"
              :options="statusOptions"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchData()">
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
        <div class="flex gap-4">
          <NButton type="primary" @click="tableModalRef?.openModal('add')">
            <template #icon>
              <nova-icon icon="carbon:add" :size="18" />
            </template>
            新增员工
          </NButton>
          <NPopconfirm @positive-click="batchDelete">
            <template #trigger>
              <NButton type="error">
                <template #icon>
                  <nova-icon icon="carbon:trash-can" :size="18" />
                </template>
                批量删除
              </NButton>
            </template>
            确定批量删除选中的员工吗？
          </NPopconfirm>
          <NButton @click="handleExport">
            <template #icon>
              <nova-icon icon="carbon:document-export" :size="18" />
            </template>
            导出
          </NButton>
          <NButton @click="handleImport">
            <template #icon>
              <nova-icon icon="carbon:document-import" :size="18" />
            </template>
            导入
          </NButton>
        </div>
        
        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1600"
          :row-key="(row: any) => row.id"
        />
        <Pagination :count="total" @change="handlePageChange" />
      </NSpace>
    </NCard>

    <!-- 编辑弹窗 -->
    <TableModal ref="tableModalRef" @refresh="fetchData" />
  </NSpace>
</template>

