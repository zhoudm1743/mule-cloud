<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { deleteSize, fetchSizeList } from '@/service'
import { NButton, NPopconfirm, NSpace } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Basic.BasicInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const initialSearchForm = {
  value: '',
}
const searchForm = ref({ ...initialSearchForm })

// 分页处理
function handlePageChange(p: number) {
  page.value = p
  fetchData()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchData()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchData()
}

async function fetchData() {
  startLoading()
  try {
    const res = await fetchSizeList({
      page: page.value,
      page_size: pageSize.value,
      value: searchForm.value.value || undefined,
    })
    if (res.data) {
      tableData.value = res.data.sizes || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取尺码列表失败')
  }
  finally {
    endLoading()
  }
}

async function deleteHandler(id: string) {
  try {
    await deleteSize(id)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

async function batchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的尺码')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteSize(id)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量删除失败')
  }
}

const columns: DataTableColumns<Api.Basic.BasicInfo> = [
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
    title: '尺码名称',
    key: 'value',
    width: 200,
  },
  {
    title: '备注',
    key: 'remark',
    width: 300,
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render: (row) => {
      if (!row.created_at) return '-'
      return new Date(row.created_at * 1000).toLocaleString('zh-CN')
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 180,
    render: (row) => {
      if (!row.updated_at) return '-'
      return new Date(row.updated_at * 1000).toLocaleString('zh-CN')
    },
  },
  {
    title: '操作',
    align: 'center',
    key: 'actions',
    width: 200,
    fixed: 'right',
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
          <NPopconfirm onPositiveClick={() => deleteHandler(row.id)}>
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
  fetchData()
})
</script>

<template>
  <NSpace vertical size="large">
    <NCard title="尺码管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="尺码名称">
            <n-input
              v-model:value="searchForm.value"
              placeholder="请输入尺码名称"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchData">
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
            新建尺码
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
            确认删除选中的尺码？
          </NPopconfirm>
        </div>

        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1200"
          :row-key="(row: Api.Basic.BasicInfo) => row.id"
        />
        <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
      </NSpace>
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchData" />
  </NSpace>
</template>

