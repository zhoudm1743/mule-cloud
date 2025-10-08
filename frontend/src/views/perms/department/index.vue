<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { batchDeleteDepartments, deleteDepartment, fetchDepartmentList } from '@/service'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Department.DepartmentInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const initialSearchForm = {
  name: '',
  code: '',
}
const searchForm = ref({ ...initialSearchForm })

// 分页处理
function handlePageChange(p: number) {
  page.value = p
  fetchDepartments()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchDepartments()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchDepartments()
}

async function deleteData(id: string) {
  try {
    await deleteDepartment(id)
    window.$message.success('删除成功')
    fetchDepartments()
  }
  catch (e) {
    console.error('[Delete Department Error]:', e)
  }
}

async function handleBatchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的部门')
    return
  }
  try {
    await batchDeleteDepartments(checkedRowKeys.value)
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchDepartments()
  }
  catch (e) {
    console.error('[Batch Delete Error]:', e)
  }
}

const columns: DataTableColumns<Api.Department.DepartmentInfo> = [
  {
    type: 'selection',
    width: 30,
  },
  {
    title: 'ID',
    key: 'id',
    width: 200,
    render: row => <CopyText value={row.id} />,
  },
  {
    title: '部门编码',
    key: 'code',
    width: 150,
  },
  {
    title: '部门名称',
    key: 'name',
    width: 200,
  },
  {
    title: '父部门ID',
    key: 'parent_id',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => {
      return row.parent_id || '-'
    },
  },
  {
    title: '状态',
    align: 'center',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = row.status || 0
      const statusMap: Record<number, { type: NaiveUI.ThemeColor, text: string }> = {
        1: { type: 'success', text: '启用' },
        0: { type: 'error', text: '禁用' },
      }
      const statusInfo = statusMap[status] || statusMap[0]
      return <NTag type={statusInfo.type}>{statusInfo.text}</NTag>
    },
  },
  {
    title: '操作',
    align: 'center',
    key: 'actions',
    width: '12em',
    render: (row) => {
      return (
        <NSpace justify="center">
          <NButton
            text
            type="primary"
            onClick={() => tableModalRef.value.openModal('view', row)}
          >
            查看
          </NButton>
          <NButton
            text
            type="primary"
            onClick={() => tableModalRef.value.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteData(row.id)}>
            {{
              default: () => '确认删除？',
              trigger: () => <NButton text type="error">删除</NButton>,
            }}
          </NPopconfirm>
        </NSpace>
      )
    },
  },
]

onMounted(() => {
  fetchDepartments()
})

async function fetchDepartments() {
  startLoading()
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value,
      name: searchForm.value.name || undefined,
      code: searchForm.value.code || undefined,
    }
    
    const res = await fetchDepartmentList(params)
    if (res.data) {
      tableData.value = res.data.departments || []
      total.value = res.data.total || 0
    }
  }
  catch (e) {
    console.error('[Fetch Departments Error]:', e)
  }
  finally {
    endLoading()
  }
}
</script>

<template>
  <NSpace vertical size="large">
    <NCard title="部门管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="部门名称">
            <n-input
              v-model:value="searchForm.name"
              placeholder="请输入部门名称"
              clearable
              class="w-200px"
              @keyup.enter="fetchDepartments"
            />
          </n-form-item>
          <n-form-item label="部门编码">
            <n-input
              v-model:value="searchForm.code"
              placeholder="请输入部门编码"
              clearable
              class="w-200px"
              @keyup.enter="fetchDepartments"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchDepartments">
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
          <NButton type="primary" @click="tableModalRef.openModal('add')">
            <template #icon>
              <nova-icon icon="carbon:add" :size="18" />
            </template>
            新建部门
          </NButton>
          <NPopconfirm @positive-click="handleBatchDelete">
            <template #trigger>
              <NButton type="error">
                <template #icon>
                  <nova-icon icon="carbon:trash-can" :size="18" />
                </template>
                批量删除
              </NButton>
            </template>
            确认删除所有选中部门？
          </NPopconfirm>
        </div>

        <n-data-table
          v-model:checked-row-keys="checkedRowKeys"
          :row-key="(row: Api.Department.DepartmentInfo) => row.id"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1200"
        />
        <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
      </NSpace>
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchDepartments" />
  </NSpace>
</template>

