<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean, usePermission } from '@/hooks'
import { useAuthStore } from '@/store'
import { deleteAdmin, fetchAdminList } from '@/service'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)
const { hasAction, hasResource } = usePermission()
const authStore = useAuthStore()

const tableModalRef = ref()
const tableData = ref<Api.Admin.AdminInfo[]>([])
const checkedRowKeys = ref<string[]>([])

async function deleteData(id: string) {
  try {
    await deleteAdmin(id)
    window.$message.success('删除成功')
    fetchAdmins()
  }
  catch (e) {
    console.error('[Delete Admin Error]:', e)
  }
}

async function handleBatchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的管理员')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteAdmin(id)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchAdmins()
  }
  catch (e) {
    console.error('[Batch Delete Error]:', e)
  }
}

const columns: DataTableColumns<Api.Admin.AdminInfo> = [
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
    title: '昵称',
    key: 'nickname',
    width: 150,
  },
  {
    title: '手机号',
    key: 'phone',
    width: 150,
    render: row => <CopyText value={row.phone} />,
  },
  {
    title: '邮箱',
    key: 'email',
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
      if (!row.tenant_id) {
        return <NTag type="warning">系统管理员</NTag>
      }
      return row.tenant_id
    },
  },
  {
    title: '超级管理员',
    align: 'center',
    key: 'is_super',
    width: 100,
    render: (row) => {
      return row.is_super
        ? <NTag type="error">是</NTag>
        : <NTag type="default">否</NTag>
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
          {/* 查看按钮 - 需要 read 权限 */}
          {hasAction('admin', 'read') && (
            <NButton
              text
              type="primary"
              onClick={() => tableModalRef.value.openModal('view', row)}
            >
              查看
            </NButton>
          )}
          {/* 编辑按钮 - 需要 update 权限 */}
          {hasAction('admin', 'update') && (
            <NButton
              text
              type="primary"
              onClick={() => tableModalRef.value.openModal('edit', row)}
            >
              编辑
            </NButton>
          )}
          {/* 删除按钮 - 需要 delete 权限 */}
          {hasAction('admin', 'delete') && (
            <NPopconfirm onPositiveClick={() => deleteData(row.id)}>
              {{
                default: () => '确认删除',
                trigger: () => <NButton text type="error">删除</NButton>,
              }}
            </NPopconfirm>
          )}
        </NSpace>
      )
    },
  },
]

onMounted(() => {
  fetchAdmins()
})

async function fetchAdmins() {
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
    
    // 租户超管：只能查看本租户的管理员
    if (!isSystemAdmin && userInfo?.tenant_id) {
      params.tenant_id = userInfo.tenant_id
    }
    // 系统超管：可以查看所有管理员（不添加tenant_id过滤）
    
    const res = await fetchAdminList(params)
    if (res.data) {
      tableData.value = res.data.admins || []
    }
  }
  catch (e) {
    console.error('[Fetch Admins Error]:', e)
  }
  finally {
    endLoading()
  }
}
</script>

<template>
  <n-card title="管理员管理">
    <template #header-extra>
      <n-flex>
        <!-- 新建按钮 - 需要 create 权限 -->
        <NButton
          v-if="hasResource('admin:create')"
          type="primary"
          @click="tableModalRef.openModal('add')"
        >
          <template #icon>
            <icon-park-outline-add-one />
          </template>
          新建
        </NButton>
        <NButton type="primary" secondary @click="fetchAdmins">
          <template #icon>
            <icon-park-outline-refresh />
          </template>
          刷新
        </NButton>
        <!-- 批量删除按钮 - 需要 delete 权限 -->
        <NPopconfirm
          v-if="hasResource('admin:delete')"
          @positive-click="handleBatchDelete"
        >
          <template #trigger>
            <NButton type="error" secondary>
              <template #icon>
                <icon-park-outline-delete-five />
              </template>
              批量删除
            </NButton>
          </template>
          确认删除所有选中管理员？
        </NPopconfirm>
      </n-flex>
    </template>

    <n-data-table
      v-model:checked-row-keys="checkedRowKeys"
      :row-key="(row: Api.Admin.AdminInfo) => row.id"
      :columns="columns"
      :data="tableData"
      :loading="loading"
      size="small"
      :scroll-x="1400"
    />

    <TableModal ref="tableModalRef" @refresh="fetchAdmins" />
  </n-card>
</template>
