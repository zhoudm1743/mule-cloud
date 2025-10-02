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

async function fetchTenants() {
  startLoading()
  try {
    const res = await fetchTenantList({
      page: page.value,
      page_size: pageSize.value,
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
  <div class="h-full flex-col gap-16px overflow-hidden p-16px lt-sm:p-8px">
    <NCard title="租户管理" :bordered="false" class="h-full">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="tableModalRef.openModal('add')">
            <template #icon>
              <icon-park-outline-add-one />
            </template>
            新建
          </NButton>
          <NButton @click="fetchTenants()">
            <template #icon>
              <icon-park-outline-refresh />
            </template>
            刷新
          </NButton>
          <NPopconfirm @positive-click="batchDeleteTenants">
            <template #trigger>
              <NButton type="error">
                <template #icon>
                  <icon-park-outline-delete-five />
                </template>
                批量删除
              </NButton>
            </template>
            确认删除选中的租户？
          </NPopconfirm>
        </NSpace>
      </template>

      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="tableData"
        :loading="loading"
        :scroll-x="1400"
        :row-key="(row: Api.Tenant.TenantInfo) => row.id"
      />
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchTenants" />
  </div>
</template>
