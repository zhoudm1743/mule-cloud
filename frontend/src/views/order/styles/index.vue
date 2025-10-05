<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { deleteStyle, fetchStyleList } from '@/service/api/order'
import { NButton, NImage, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Order.StyleInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const searchForm = ref({
  style_no: '',
  style_name: '',
  status: undefined as number | undefined,
})

async function fetchData() {
  startLoading()
  try {
    const res = await fetchStyleList({
      page: page.value,
      page_size: pageSize.value,
      style_no: searchForm.value.style_no || undefined,
      style_name: searchForm.value.style_name || undefined,
      status: searchForm.value.status,
    })
    if (res.data) {
      tableData.value = res.data.styles || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取款式列表失败')
  }
  finally {
    endLoading()
  }
}

async function deleteHandler(id: string) {
  try {
    await deleteStyle(id)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

async function batchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的款式')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteStyle(id)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量删除失败')
  }
}

const columns: DataTableColumns<Api.Order.StyleInfo> = [
  {
    type: 'selection',
    width: 30,
  },
  {
    title: '图片',
    key: 'images',
    width: 80,
    render: (row) => {
      if (!row.images || row.images.length === 0)
        return <div class="text-center">-</div>
      return <NImage width={60} height={60} src={row.images[0]} />
    },
  },
  {
    title: 'ID',
    key: 'id',
    width: 180,
    render: row => <CopyText value={row.id} />,
  },
  {
    title: '款号',
    key: 'style_no',
    width: 150,
  },
  {
    title: '款名',
    key: 'style_name',
    width: 150,
  },
  {
    title: '单价',
    key: 'unit_price',
    width: 100,
  },
  {
    title: '颜色',
    key: 'colors',
    width: 200,
    render: (row) => {
      if (!row.colors || row.colors.length === 0) return '-'
      return row.colors.join(', ')
    },
  },
  {
    title: '尺码',
    key: 'sizes',
    width: 200,
    render: (row) => {
      if (!row.sizes || row.sizes.length === 0) return '-'
      return row.sizes.join(', ')
    },
  },
  {
    title: '工序数',
    key: 'procedures',
    width: 80,
    render: (row) => {
      return row.procedures?.length || 0
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      return row.status === 1
        ? <NTag type="success">启用</NTag>
        : <NTag type="default">禁用</NTag>
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 160,
    render: (row) => {
      if (!row.created_at) return '-'
      return new Date(row.created_at * 1000).toLocaleString('zh-CN')
    },
  },
  {
    title: '备注',
    key: 'remark',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace>
          <NButton
            size={'small'}
            onClick={() => tableModalRef.value?.openModal('view', row)}
          >
            查看
          </NButton>
          <NButton
            size={'small'}
            type={'primary'}
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
  <div>
    <NCard title="款式管理" :bordered="false" class="rounded-8px shadow-sm">
      <div class="flex-col">
        <!-- 搜索区域 -->
        <NSpace class="pb-12px" justify="space-between">
          <NSpace>
            <NInput
              v-model:value="searchForm.style_no"
              placeholder="搜索款号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
            <NInput
              v-model:value="searchForm.style_name"
              placeholder="搜索款名"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
            <NSelect
              v-model:value="searchForm.status"
              placeholder="状态"
              clearable
              class="w-120px"
              :options="[
                { label: '启用', value: 1 },
                { label: '禁用', value: 0 },
              ]"
            />
            <NButton type="primary" @click="fetchData">
              <template #icon>
                <nova-icon icon="carbon:search" :size="18" />
              </template>
              查询
            </NButton>
            <NButton @click="searchForm = { style_no: '', style_name: '', status: undefined }; fetchData()">
              <template #icon>
                <nova-icon icon="carbon:reset" :size="18" />
              </template>
              重置
            </NButton>
          </NSpace>
          <NSpace>
            <NButton type="primary" @click="tableModalRef?.openModal('add')">
              <template #icon>
                <nova-icon icon="carbon:add" :size="18" />
              </template>
              新建款式
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
              确定批量删除选中的款式吗？
            </NPopconfirm>
          </NSpace>
        </NSpace>

        <!-- 表格 -->
        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1800"
          :row-key="(row: Api.Order.StyleInfo) => row.id"
          :pagination="{
            page,
            pageSize,
            pageCount: Math.ceil(total / pageSize),
            showSizePicker: true,
            pageSizes: [10, 20, 50, 100],
            onChange: (p: number) => { page = p; fetchData() },
            onUpdatePageSize: (ps: number) => { pageSize = ps; page = 1; fetchData() },
          }"
          class="flex-1-hidden"
        />
      </div>
    </NCard>

    <!-- 编辑弹窗 -->
    <TableModal ref="tableModalRef" @refresh="fetchData" />
  </div>
</template>

<style scoped>
.flex-col {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.flex-1-hidden {
  flex: 1;
  overflow: hidden;
}
</style>
