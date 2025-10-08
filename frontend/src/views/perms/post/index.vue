<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { batchDeletePosts, deletePost, fetchPostList } from '@/service'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Post.PostInfo[]>([])
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
  fetchPosts()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchPosts()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchPosts()
}

async function deleteData(id: string) {
  try {
    await deletePost(id)
    window.$message.success('删除成功')
    fetchPosts()
  }
  catch (e) {
    console.error('[Delete Post Error]:', e)
  }
}

async function handleBatchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的岗位')
    return
  }
  try {
    await batchDeletePosts(checkedRowKeys.value)
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchPosts()
  }
  catch (e) {
    console.error('[Batch Delete Error]:', e)
  }
}

const columns: DataTableColumns<Api.Post.PostInfo> = [
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
    title: '岗位编码',
    key: 'code',
    width: 150,
  },
  {
    title: '岗位名称',
    key: 'name',
    width: 200,
  },
  {
    title: '父岗位ID',
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
  fetchPosts()
})

async function fetchPosts() {
  startLoading()
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value,
      name: searchForm.value.name || undefined,
      code: searchForm.value.code || undefined,
    }
    
    const res = await fetchPostList(params)
    if (res.data) {
      tableData.value = res.data.posts || []
      total.value = res.data.total || 0
    }
  }
  catch (e) {
    console.error('[Fetch Posts Error]:', e)
  }
  finally {
    endLoading()
  }
}
</script>

<template>
  <NSpace vertical size="large">
    <NCard title="岗位管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="岗位名称">
            <n-input
              v-model:value="searchForm.name"
              placeholder="请输入岗位名称"
              clearable
              class="w-200px"
              @keyup.enter="fetchPosts"
            />
          </n-form-item>
          <n-form-item label="岗位编码">
            <n-input
              v-model:value="searchForm.code"
              placeholder="请输入岗位编码"
              clearable
              class="w-200px"
              @keyup.enter="fetchPosts"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchPosts">
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
            新建岗位
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
            确认删除所有选中岗位？
          </NPopconfirm>
        </div>

        <n-data-table
          v-model:checked-row-keys="checkedRowKeys"
          :row-key="(row: Api.Post.PostInfo) => row.id"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1200"
        />
        <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
      </NSpace>
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchPosts" />
  </NSpace>
</template>

