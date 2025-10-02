<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { useAuthStore } from '@/store'
import { batchDeleteRoles, deleteRole, fetchRoleList } from '@/service'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)
const authStore = useAuthStore()

const tableModalRef = ref()
const tableData = ref<Api.Role.RoleInfo[]>([])
const checkedRowKeys = ref<string[]>([])

async function deleteData(id: string) {
  try {
    await deleteRole(id)
    window.$message.success('删除成功')
    fetchRoles()
  }
  catch (e) {
    console.error('[Delete Role Error]:', e)
  }
}

async function handleBatchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的角色')
    return
  }
  try {
    await batchDeleteRoles(checkedRowKeys.value)
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchRoles()
  }
  catch (e) {
    console.error('[Batch Delete Error]:', e)
  }
}

const columns: DataTableColumns<Api.Role.RoleInfo> = [
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
    title: '角色代码',
    key: 'code',
    width: 150,
  },
  {
    title: '角色名称',
    key: 'name',
    width: 200,
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '租户ID',
    key: 'tenant_id',
    width: 200,
    ellipsis: {
      tooltip: true,
    },
    render: (row) => {
      return row.tenant_id || '系统角色'
    },
  },
  {
    title: '菜单权限数',
    align: 'center',
    key: 'menus',
    width: 100,
    render: (row) => {
      return row.menus?.length || 0
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
    width: '18em',
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
  fetchRoles()
})

async function fetchRoles() {
  startLoading()
  try {
    // 构建查询参数
    const params: any = {
      page: 1,
      page_size: 100,
    }
    
    // 检查是否为系统超管（角色包含 'super' 且没有 tenant_id）
    const userInfo = authStore.userInfo
    const isSystemAdmin = userInfo?.role?.includes('super') && !userInfo?.tenant_id
    
    // 租户超管：只能查看本租户的角色
    if (!isSystemAdmin && userInfo?.tenant_id) {
      params.tenant_id = userInfo.tenant_id
    }
    // 系统超管：可以查看所有角色（不添加tenant_id过滤）
    
    const res = await fetchRoleList(params)
    if (res.data) {
      tableData.value = res.data.roles || []
    }
  }
  catch (e) {
    console.error('[Fetch Roles Error]:', e)
  }
  finally {
    endLoading()
  }
}
</script>

<template>
  <n-card title="角色管理">
    <template #header-extra>
      <n-flex>
        <NButton type="primary" @click="tableModalRef.openModal('add')">
          <template #icon>
            <icon-park-outline-add-one />
          </template>
          新建
        </NButton>
        <NButton type="primary" secondary @click="fetchRoles">
          <template #icon>
            <icon-park-outline-refresh />
          </template>
          刷新
        </NButton>
        <NPopconfirm @positive-click="handleBatchDelete">
          <template #trigger>
            <NButton type="error" secondary>
              <template #icon>
                <icon-park-outline-delete-five />
              </template>
              批量删除
            </NButton>
          </template>
          确认删除所有选中角色？
        </NPopconfirm>
      </n-flex>
    </template>

    <n-data-table
      v-model:checked-row-keys="checkedRowKeys"
      :row-key="(row: Api.Role.RoleInfo) => row.id"
      :columns="columns"
      :data="tableData"
      :loading="loading"
      size="small"
      :scroll-x="1400"
    />

    <TableModal ref="tableModalRef" @refresh="fetchRoles" />
  </n-card>
</template>
