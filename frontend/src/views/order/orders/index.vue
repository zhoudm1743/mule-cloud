<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { copyOrder, createCuttingTask, deleteOrder, fetchCuttingTaskByOrderId, fetchOrderList } from '@/service/api/order'
import { NButton, NImage, NPopconfirm, NSpace, NTag } from 'naive-ui'
import { useRouter } from 'vue-router'
import TableModal from './components/TableModal.vue'

const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const router = useRouter()
const tableModalRef = ref()
const tableData = ref<Api.Order.OrderInfo[]>([])
const checkedRowKeys = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 搜索表单
const searchForm = ref({
  contract_no: '',
  style_no: '',
  customer_id: '',
  status: undefined as number | undefined,
})

// 订单状态映射
const statusMap = {
  0: { label: '草稿', type: 'default' },
  1: { label: '已下单', type: 'info' },
  2: { label: '生产中', type: 'warning' },
  3: { label: '已完成', type: 'success' },
  4: { label: '已取消', type: 'error' },
}

async function fetchData() {
  startLoading()
  try {
    const res = await fetchOrderList({
      page: page.value,
      page_size: pageSize.value,
      contract_no: searchForm.value.contract_no || undefined,
      style_no: searchForm.value.style_no || undefined,
      status: searchForm.value.status,
    })
    if (res.data) {
      tableData.value = res.data.orders || []
      total.value = res.data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取订单列表失败')
  }
  finally {
    endLoading()
  }
}

async function deleteHandler(id: string) {
  try {
    await deleteOrder(id)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

async function copyHandler(id: string) {
  try {
    await copyOrder(id)
    window.$message.success('复制成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '复制失败')
  }
}

async function batchDelete() {
  if (checkedRowKeys.value.length === 0) {
    window.$message.warning('请选择要删除的订单')
    return
  }
  try {
    for (const id of checkedRowKeys.value) {
      await deleteOrder(id)
    }
    window.$message.success('批量删除成功')
    checkedRowKeys.value = []
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '批量删除失败')
  }
}

// 创建裁剪任务并跳转
async function handleCreateCutting(order: Api.Order.OrderInfo) {
  try {
    // 先检查是否已经有裁剪任务
    try {
      const { data } = await fetchCuttingTaskByOrderId(order.id)
      if (data?.task) {
        window.$message.warning('该订单已有裁剪任务，即将跳转到裁剪管理页面')
        router.push('/order/cutting')
        return
      }
    }
    catch (error) {
      // 没有找到裁剪任务，继续创建
    }

    // 创建裁剪任务
    const { data } = await createCuttingTask({
      order_id: order.id,
    })

    if (data?.task) {
      window.$message.success('裁剪任务创建成功')
      // 跳转到裁剪管理页面
      router.push('/order/cutting')
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '创建裁剪任务失败')
  }
}

const columns: DataTableColumns<Api.Order.OrderInfo> = [
  {
    type: 'selection',
    width: 30,
  },
  {
    title: '图片',
    key: 'style_image',
    width: 80,
    render: (row) => {
      if (!row.style_image)
        return <div class="text-center">-</div>
      return <NImage width={60} height={60} src={row.style_image} />
    },
  },
  {
    title: '合同号',
    key: 'contract_no',
    width: 150,
    render: row => <CopyText value={row.contract_no} />,
  },
  {
    title: '客户',
    key: 'customer_name',
    width: 120,
  },
  {
    title: '款号',
    key: 'style_no',
    width: 120,
  },
  {
    title: '数量',
    key: 'quantity',
    width: 80,
  },
  {
    title: '单价',
    key: 'unit_price',
    width: 80,
  },
  {
    title: '总金额',
    key: 'total_amount',
    width: 100,
  },
  {
    title: '进度',
    key: 'progress',
    width: 80,
    render: (row) => {
      return `${row.progress.toFixed(0)}%`
    },
  },
  {
    title: '交货日期',
    key: 'delivery_date',
    width: 120,
  },
  {
    title: '订单类型',
    key: 'order_type_name',
    width: 100,
    render: row => row.order_type_name || '-',
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) => {
      const status = statusMap[row.status as keyof typeof statusMap] || statusMap[0]
      return <NTag type={status.type as any}>{status.label}</NTag>
    },
  },
  {
    title: '下单时间',
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
    width: 280,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace>
          <NButton
            size={'small'}
            onClick={() => tableModalRef.value?.openModal('view', row)}
          >
            详情
          </NButton>
          <NButton
            size={'small'}
            type={'primary'}
            onClick={() => tableModalRef.value?.openModal('edit', row)}
          >
            编辑
          </NButton>
          <NButton
            size={'small'}
            type={'info'}
            onClick={() => handleCreateCutting(row)}
          >
            裁剪
          </NButton>
          <NButton
            size={'small'}
            onClick={() => copyHandler(row.id)}
          >
            复制
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
    <NCard title="订单管理" :bordered="false" class="rounded-8px shadow-sm">
      <div class="flex-col">
        <!-- 搜索区域 -->
        <NSpace class="pb-12px" justify="space-between">
          <NSpace>
            <NInput
              v-model:value="searchForm.contract_no"
              placeholder="搜索合同号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
            <NInput
              v-model:value="searchForm.style_no"
              placeholder="搜索款号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
            <NSelect
              v-model:value="searchForm.status"
              placeholder="订单状态"
              clearable
              class="w-150px"
              :options="[
                { label: '草稿', value: 0 },
                { label: '已下单', value: 1 },
                { label: '生产中', value: 2 },
                { label: '已完成', value: 3 },
                { label: '已取消', value: 4 },
              ]"
            />
            <NButton type="primary" @click="fetchData">
              <template #icon>
                <nova-icon icon="carbon:search" :size="18" />
              </template>
              查询
            </NButton>
            <NButton @click="searchForm = { contract_no: '', style_no: '', customer_id: '', status: undefined }; fetchData()">
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
              新建订单
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
              确定批量删除选中的订单吗？
            </NPopconfirm>
          </NSpace>
        </NSpace>

        <!-- 表格 -->
        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="2000"
          :row-key="(row: Api.Order.OrderInfo) => row.id"
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
