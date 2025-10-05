<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { deleteProcedure, fetchProcedureList } from '@/service'
import { NButton, NPopconfirm, NSpace } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Basic.BasicInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

async function fetchData() {
  startLoading()
  try {
    const res = await fetchProcedureList({
      page: page.value,
      page_size: pageSize.value,
    })
    if (res.data) {
      tableData.value = res.data.procedures || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取工序列表失败')
  }
  finally {
    endLoading()
  }
}

async function deleteHandler(id: string) {
  try {
    await deleteProcedure(id)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

async function batchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的工序')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteProcedure(id)
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
    title: '工序名称',
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
  <div class="h-full flex-col gap-16px overflow-hidden p-16px lt-sm:p-8px">
    <NCard title="工序管理" :bordered="false" class="h-full">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="tableModalRef.openModal('add')">
            <template #icon>
              <icon-park-outline-add-one />
            </template>
            新建
          </NButton>
          <NButton @click="fetchData()">
            <template #icon>
              <icon-park-outline-refresh />
            </template>
            刷新
          </NButton>
          <NPopconfirm @positive-click="batchDelete">
            <template #trigger>
              <NButton type="error">
                <template #icon>
                  <icon-park-outline-delete-five />
                </template>
                批量删除
              </NButton>
            </template>
            确认删除选中的工序？
          </NPopconfirm>
        </NSpace>
      </template>

      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :scroll-x="1200"
        :row-key="(row: Api.Basic.BasicInfo) => row.id"
      />
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchData" />
  </div>
</template>

