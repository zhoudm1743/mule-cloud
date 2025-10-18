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
  fetchRoles()
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchRoles()
}

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
      page: page.value,
      page_size: pageSize.value,
      name: searchForm.value.name || undefined,
      code: searchForm.value.code || undefined,
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
      total.value = res.data.total || 0
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
  <NSpace vertical size="large">
    <NCard title="角色管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="角色名称">
            <n-input
              v-model:value="searchForm.name"
              placeholder="请输入角色名称"
              clearable
              class="w-200px"
              @keyup.enter="fetchRoles"
            />
          </n-form-item>
          <n-form-item label="角色代码">
            <n-input
              v-model:value="searchForm.code"
              placeholder="请输入角色代码"
              clearable
              class="w-200px"
              @keyup.enter="fetchRoles"
            />
          </n-form-item>
          <n-flex class="ml-auto">
            <NButton type="primary" @click="fetchRoles">
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
            新建角色
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
            确认删除所有选中角色？
          </NPopconfirm>
        </div>

        <n-data-table
          v-model:checked-row-keys="checkedRowKeys"
          :row-key="(row: Api.Role.RoleInfo) => row.id"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="1400"
        />
        <Pagination :count="total" @change="handlePageChange" />
      </NSpace>
    </NCard>

    <TableModal ref="tableModalRef" @refresh="fetchRoles" />
  </NSpace>
</template>
