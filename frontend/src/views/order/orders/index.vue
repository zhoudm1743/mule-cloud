<script setup lang="tsx">
import type { DataTableColumns } from 'naive-ui'
import CopyText from '@/components/custom/CopyText.vue'
import { useBoolean } from '@/hooks'
import { copyOrder, createCuttingTask, deleteOrder, fetchCuttingTaskByOrderId, fetchOrderList } from '@/service/api/order'
import { NAlert, NButton, NFormItem, NImage, NInput, NPopconfirm, NRadio, NRadioGroup, NSpace, NSwitch, NTag, NText } from 'naive-ui'
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

// 分页处理函数
function handlePageChange(p: number) {
  page.value = p
  fetchData()
}

function handlePageSizeChange(ps: number) {
  pageSize.value = ps
  page.value = 1
  fetchData()
}

// 搜索表单
const initialSearchForm = {
  contract_no: '',
  style_no: '',
  customer_id: '',
  status: undefined as number | undefined,
}
const searchForm = ref({ ...initialSearchForm })

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

// 复制订单相关
const showCopyModal = ref(false)
const copyOrderId = ref('')
const copyForm = ref({
  is_related: true,
  relation_type: 'add',
  relation_remark: ''
})

function openCopyModal(id: string) {
  copyOrderId.value = id
  copyForm.value = {
    is_related: true,
    relation_type: 'add',
    relation_remark: ''
  }
  showCopyModal.value = true
}

async function confirmCopy() {
  try {
    await copyOrder(copyOrderId.value, copyForm.value)
    window.$message.success('复制成功')
    showCopyModal.value = false
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '复制失败')
  }
}

// 重置搜索
function handleResetSearch() {
  searchForm.value = { ...initialSearchForm }
  page.value = 1
  fetchData()
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
            onClick={() => openCopyModal(row.id)}
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
  <NSpace vertical size="large">
    <NCard title="订单管理" :bordered="false" class="rounded-8px shadow-sm">
      <n-form :model="searchForm" label-placement="left" inline :show-feedback="false">
        <n-flex>
          <n-form-item label="合同号">
            <NInput
              v-model:value="searchForm.contract_no"
              placeholder="请输入合同号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-form-item label="款号">
            <NInput
              v-model:value="searchForm.style_no"
              placeholder="请输入款号"
              clearable
              class="w-200px"
              @keyup.enter="fetchData"
            />
          </n-form-item>
          <n-form-item label="订单状态">
            <NSelect
              v-model:value="searchForm.status"
              placeholder="请选择"
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
        </div>
        
        <NDataTable
          v-model:checked-row-keys="checkedRowKeys"
          :columns="columns"
          :data="tableData"
          :loading="loading"
          :scroll-x="2000"
          :row-key="(row: Api.Order.OrderInfo) => row.id"
        />
        <Pagination :count="total" :page="page" :page-size="pageSize" @change="handlePageChange" @update-page-size="handlePageSizeChange" />
      </NSpace>
    </NCard>

    <!-- 编辑弹窗 -->
    <TableModal ref="tableModalRef" @refresh="fetchData" />
    
    <!-- 复制订单弹窗 -->
    <NModal v-model:show="showCopyModal" preset="card" title="复制订单" class="w-500px">
      <NForm label-placement="left" :label-width="100">
        <NFormItem label="是否关联">
          <NSwitch v-model:value="copyForm.is_related">
            <template #checked>
              关联原订单
            </template>
            <template #unchecked>
              独立订单
            </template>
          </NSwitch>
        </NFormItem>
        
        <template v-if="copyForm.is_related">
          <NFormItem label="关联类型">
            <NRadioGroup v-model:value="copyForm.relation_type">
              <NSpace>
                <NRadio value="add">
                  追加订单
                  <NText depth="3" style="font-size: 12px; margin-left: 8px">
                    （客户增加数量）
                  </NText>
                </NRadio>
                <NRadio value="copy">
                  复制订单
                  <NText depth="3" style="font-size: 12px; margin-left: 8px">
                    （重复下单）
                  </NText>
                </NRadio>
              </NSpace>
            </NRadioGroup>
          </NFormItem>
          
          <NFormItem label="关联说明">
            <NInput
              v-model:value="copyForm.relation_remark"
              type="textarea"
              :rows="3"
              placeholder="选填，可说明追加原因或备注信息"
            />
          </NFormItem>
        </template>
        
        <NAlert type="info" style="margin-bottom: 16px">
          <template v-if="copyForm.is_related">
            <div v-if="copyForm.relation_type === 'add'">
              <strong>追加订单：</strong>新订单将关联原订单，合同号自动添加"-A"后缀，便于追踪和统计。
            </div>
            <div v-else>
              <strong>复制订单：</strong>新订单将关联原订单，合同号自动添加"-C"后缀，保留关联关系。
            </div>
          </template>
          <template v-else>
            <strong>独立订单：</strong>新订单与原订单无关联，合同号添加"-copy"后缀。
          </template>
        </NAlert>
      </NForm>
      
      <template #footer>
        <NSpace justify="end">
          <NButton @click="showCopyModal = false">
            取消
          </NButton>
          <NButton type="primary" @click="confirmCopy">
            确认复制
          </NButton>
        </NSpace>
      </template>
    </NModal>
  </NSpace>
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
