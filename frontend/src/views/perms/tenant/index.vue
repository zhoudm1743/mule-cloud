<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { deleteTenant, fetchTenantList } from '@/service'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const tableModalRef = ref()
const tableData = ref<Api.Tenant.TenantInfo[]>([])
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
function handlePageChange(p: number, ps: number) {
  page.value = p
  pageSize.value = ps
  fetchTenants()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchTenants()
}

async function fetchTenants() {
  startLoading()
  try {
    const res = await fetchTenantList({
      page: page.value,
      page_size: pageSize.value,
      name: searchForm.value.name || undefined,
      code: searchForm.value.code || undefined,
    })
    if (res.data) {
      tableData.value = res.data.tenants || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取租户列表失败')
  }
  finally {
    endLoading()
  }
}

async function deleteTenantHandler(id: string) {
  try {
    await deleteTenant(id)
    window.$message.success('删除成功')
    fetchTenants()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

async function batchDeleteTenants() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的租户')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteTenant(id)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchTenants()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量删除失败')
  }
}

const columns: DataTableColumns<Api.Tenant.TenantInfo> = [
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
    title: '租户代码',
    key: 'code',
    width: 150,
  },
  {
    title: '租户名称',
    key: 'name',
    width: 200,
  },
  {
    title: '联系人',
    key: 'contact',
    width: 120,
  },
  {
    title: '联系电话',
    key: 'phone',
    width: 150,
  },
  {
    title: '联系邮箱',
    key: 'email',
    width: 200,
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
    title: '权限菜单数',
    align: 'center',
    key: 'menus',
    width: 100,
    render: (row) => {
      return row.menus?.length || 0
    },
  },
  {
    title: '操作',
    align: 'center',
    key: 'actions',
    width: '18em',
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
          <NButton
            text
            type="warning"
            onClick={() => tableModalRef.value.openModal('assignMenus', row)}
          >
            分配权限
          </NButton>
          <NPopconfirm onPositiveClick={() => deleteTenantHandler(row.id)}>
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
  fetchTenants()
})
</script>

<template>
  <NSpace vertical size="large">
    <NCard title="租户管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="租户名称">
            <n-input
              v-model:value="searchForm.name"
              placeholder="请输入租户名称"
              clearable
              class="w-200px"
              @keyup.enter="fetchTenants"
            />
          </n-form-item>
          <n-form-item label="租户代码">
            <n-input
              v-model:value="searchForm.code"
              placeholder="请输入租户代码"
              clearable
              class="w-200px"
              @keyup.enter="fetchTenants"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchTenants">
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
            新建租户
          </NButton>
          <NPopconfirm @positive-click="batchDeleteTenants">
            <template #trigger>
              <NButton type="error">
                <template #icon>
                  <nova-icon icon="carbon:trash-can" :size="18" />
                </template>
                批量删除
              </NButton>
            </template>
            确认删除选中的租户？
          </NPopconfirm>
        </div>

        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1400"
          :row-key="(row: Api.Tenant.TenantInfo) => row.id"
        />
        <Pagination :count="total" @change="handlePageChange" />
      </NSpace>
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchTenants" />
  </NSpace>
</template>
