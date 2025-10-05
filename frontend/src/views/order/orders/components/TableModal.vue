<script setup lang="ts">
import { computed, h, reactive, ref } from 'vue'
import { NCheckbox, NInput, NInputNumber, NSelect } from 'naive-ui'
import { createOrder, updateOrder, updateOrderProcedure, updateOrderStyle } from '@/service/api/order'
import { fetchAllStyles } from '@/service/api/order'
import { fetchAllColors, fetchAllCustomers, fetchAllOrderTypes, fetchAllProcedures, fetchAllSalesmans, fetchAllSizes } from '@/service/api/basic'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'OrderModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建订单',
    edit: '编辑订单',
    view: '查看订单',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')
const currentStep = ref(0) // 0, 1, 2 对应三个步骤

// 基础数据选项
const customerOptions = ref<any[]>([])
const styleOptions = ref<any[]>([])
const colorOptions = ref<any[]>([])
const sizeOptions = ref<any[]>([])
const orderTypeOptions = ref<any[]>([])
const salesmanOptions = ref<any[]>([])
const procedureOptions = ref<any[]>([])

// 表单数据
const formDefault = () => ({
  id: '',
  // 步骤1：基础信息
  contract_no: '',
  customer_id: '',
  delivery_date: undefined as string | undefined,
  order_type_id: '',
  salesman_id: '',
  remark: '',
  // 步骤2：款式数量
  style_id: '',
  colors: [] as string[],
  sizes: [] as string[],
  unit_price: 0,
  quantity: 0,
  items: [] as Api.Order.OrderItem[],
  // 步骤3：工序
  procedures: [] as Api.Order.OrderProcedure[],
  // 其他
  status: 0,
})

const formModel = reactive(formDefault())

const rules: any = {
  contract_no: { required: true, message: '请输入合同号', trigger: 'blur' },
  customer_id: { required: true, message: '请选择客户', trigger: 'change' },
}

const formRef = ref()
const loadingSubmit = ref(false)

// 加载基础数据
async function loadBaseData() {
  try {
    const [customers, styles, colors, sizes, orderTypes, salesmans, procedures] = await Promise.all([
      fetchAllCustomers(),
      fetchAllStyles(),
      fetchAllColors(),
      fetchAllSizes(),
      fetchAllOrderTypes(),
      fetchAllSalesmans(),
      fetchAllProcedures(),
    ])

    // 直接使用名称作为value，支持快速创建
    customerOptions.value = (customers.data?.customers || []).map(item => ({ label: item.value, value: item.value }))
    styleOptions.value = (styles.data?.styles || []).map(item => ({ label: `${item.style_no} - ${item.style_name}`, value: item.id, ...item }))
    colorOptions.value = (colors.data?.colors || []).map(item => ({ label: item.value, value: item.value }))
    sizeOptions.value = (sizes.data?.sizes || []).map(item => ({ label: item.value, value: item.value }))
    orderTypeOptions.value = (orderTypes.data?.order_types || []).map(item => ({ label: item.value, value: item.value }))
    salesmanOptions.value = (salesmans.data?.salesmans || []).map(item => ({ label: item.value, value: item.value }))
    procedureOptions.value = (procedures.data?.procedures || []).map(item => ({ label: item.value, value: item.value }))
  }
  catch (error: any) {
    window.$message.error('加载基础数据失败')
  }
}

function openModal(type: ModalType, data?: Api.Order.OrderInfo) {
  showModal()
  modalType.value = type
  currentStep.value = 0
  Object.assign(formModel, formDefault())
  loadBaseData()

  if (type === 'view' || type === 'edit') {
    if (data) {
      Object.assign(formModel, data)
      if (type === 'edit') {
        currentStep.value = 2 // 编辑时直接跳到最后一步
      }
    }
  }
}

// 步骤1：创建基础信息
async function handleStep1() {
  await formRef.value?.validate()
  loadingSubmit.value = true

  try {
    const res = await createOrder({
      contract_no: formModel.contract_no,
      customer_id: formModel.customer_id,
      delivery_date: formModel.delivery_date || undefined,
      order_type_id: formModel.order_type_id || undefined,
      salesman_id: formModel.salesman_id || undefined,
      remark: formModel.remark || undefined,
    })

    if (res.data?.order) {
      formModel.id = res.data.order.id
      window.$message.success('基础信息保存成功')
      currentStep.value = 1
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
  finally {
    loadingSubmit.value = false
  }
}

// 当选择客户时，自动填充客户名称
function onCustomerChange(customerId: string) {
  const customer = customerOptions.value.find(c => c.value === customerId)
  if (customer) {
    // 客户名称会在后端保存时自动填充
  }
}

// 当选择款式时，自动加载款式的颜色、尺码和工序
function onStyleChange(styleId: string) {
  const style = styleOptions.value.find(s => s.value === styleId)
  if (style) {
    // 自动填充款式信息
    formModel.colors = style.colors || []
    formModel.sizes = style.sizes || []
    formModel.unit_price = style.unit_price || 0
    
    // 复制工序清单
    formModel.procedures = (style.procedures || []).map((proc: any, index: number) => ({
      sequence: index + 1,
      procedure_name: proc.procedure_name || '',
      unit_price: proc.unit_price || 0,
      assigned_worker: proc.assigned_worker || '',
      is_slowest: proc.is_slowest || false,
      no_bundle: proc.no_bundle || false,
    }))

    // 生成订单明细（所有颜色和尺码的组合）
    generateOrderItems()
  }
}

// 生成订单明细
function generateOrderItems() {
  const items: Api.Order.OrderItem[] = []
  
  if (formModel.colors.length === 0 || formModel.sizes.length === 0) {
    formModel.items = []
    return
  }
  
  formModel.colors.forEach((color) => {
    formModel.sizes.forEach((size) => {
      // 查找是否已存在该组合，保留原有数量
      const existingItem = formModel.items.find(
        item => item.color === color && item.size === size
      )
      
      items.push({
        color: color,
        size: size,
        quantity: existingItem?.quantity || 0,
      })
    })
  })
  formModel.items = items
  
  // 重新计算总数量
  formModel.quantity = items.reduce((sum, item) => sum + item.quantity, 0)
}

// 当颜色或尺码改变时，重新生成明细
function onColorsChange() {
  generateOrderItems()
}

function onSizesChange() {
  generateOrderItems()
}

// 步骤2：保存款式数量
async function handleStep2() {
  if (!formModel.style_id) {
    window.$message.error('请选择款式')
    return
  }
  if (formModel.items.length === 0) {
    window.$message.error('请配置颜色尺码组合')
    return
  }

  // 计算总数量
  const totalQuantity = formModel.items.reduce((sum, item) => sum + item.quantity, 0)
  formModel.quantity = totalQuantity

  loadingSubmit.value = true
  try {
    await updateOrderStyle(formModel.id, {
      style_id: formModel.style_id,
      colors: formModel.colors,
      sizes: formModel.sizes,
      unit_price: formModel.unit_price,
      quantity: formModel.quantity,
      items: formModel.items,
    })

    window.$message.success('款式数量保存成功')
    currentStep.value = 2
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
  finally {
    loadingSubmit.value = false
  }
}

// 步骤3：保存工序
async function handleStep3() {
  // 验证工序：至少有一个最终工序
  if (formModel.procedures && formModel.procedures.length > 0) {
    const hasFinalProcedure = formModel.procedures.some(proc => proc.is_slowest)
    if (!hasFinalProcedure) {
      window.$message.error('必须至少选择一个最终工序')
      return
    }
  }
  
  loadingSubmit.value = true
  try {
    await updateOrderProcedure(formModel.id, {
      procedures: formModel.procedures,
    })

    window.$message.success('订单创建成功')
    hideModal()
    emit('refresh')
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
  finally {
    loadingSubmit.value = false
  }
}

// 编辑模式的保存
async function handleEditSubmit() {
  await formRef.value?.validate()
  
  // 验证工序：至少有一个最终工序
  if (formModel.procedures && formModel.procedures.length > 0) {
    const hasFinalProcedure = formModel.procedures.some(proc => proc.is_slowest)
    if (!hasFinalProcedure) {
      window.$message.error('必须至少选择一个最终工序')
      return
    }
  }
  
  loadingSubmit.value = true

  try {
    await updateOrder(formModel.id, {
      contract_no: formModel.contract_no,
      customer_id: formModel.customer_id,
      delivery_date: formModel.delivery_date || undefined,
      order_type_id: formModel.order_type_id || undefined,
      salesman_id: formModel.salesman_id || undefined,
      remark: formModel.remark || undefined,
      style_id: formModel.style_id || undefined,
      colors: formModel.colors,
      sizes: formModel.sizes,
      unit_price: formModel.unit_price,
      quantity: formModel.quantity,
      items: formModel.items,
      procedures: formModel.procedures,
      status: formModel.status,
    })

    window.$message.success('更新成功')
    hideModal()
    emit('refresh')
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
  finally {
    loadingSubmit.value = false
  }
}

// 订单明细表格列定义
const itemColumns = computed(() => [
  { title: '颜色', key: 'color' },
  { title: '尺码', key: 'size' },
  {
    title: '数量',
    key: 'quantity',
    render: (row: Api.Order.OrderItem, index: number) => {
      return h(NInputNumber, {
        value: formModel.items[index].quantity,
        disabled: modalType.value === 'view',
        min: 0,
        class: 'w-full',
        'onUpdate:value': (value: number | null) => {
          formModel.items[index].quantity = value || 0
          // 自动计算总数量
          formModel.quantity = formModel.items.reduce((sum, item) => sum + item.quantity, 0)
        },
      })
    },
  },
])

// 工序表格列定义
const procedureColumns = computed(() => [
  { title: '顺序', key: 'sequence', width: 60 },
  { title: '工序名称', key: 'procedure_name' },
  {
    title: '工价',
    key: 'unit_price',
    width: 120,
    render: (row: Api.Order.OrderProcedure, index: number) => {
      return h(NInputNumber, {
        value: formModel.procedures[index].unit_price,
        disabled: modalType.value === 'view',
        min: 0,
        precision: 2,
        class: 'w-full',
        'onUpdate:value': (value: number | null) => {
          formModel.procedures[index].unit_price = value || 0
        },
      })
    },
  },
  {
    title: '指定工人',
    key: 'assigned_worker',
    width: 120,
    render: (row: Api.Order.OrderProcedure, index: number) => {
      return h(NInput, {
        value: formModel.procedures[index].assigned_worker,
        disabled: modalType.value === 'view',
        placeholder: '指定工人',
        'onUpdate:value': (value: string) => {
          formModel.procedures[index].assigned_worker = value
        },
      })
    },
  },
  {
    title: '最终工序',
    key: 'is_slowest',
    width: 80,
    render: (row: Api.Order.OrderProcedure, index: number) => {
      return h(NCheckbox, {
        checked: formModel.procedures[index].is_slowest,
        disabled: modalType.value === 'view',
        'onUpdate:checked': (value: boolean) => {
          // 如果选中，取消其他工序的最终工序标记
          if (value) {
            formModel.procedures.forEach((proc, idx) => {
              proc.is_slowest = idx === index
            })
          } else {
            formModel.procedures[index].is_slowest = false
          }
        },
      })
    },
  },
  {
    title: '不分扎',
    key: 'no_bundle',
    width: 80,
    render: (row: Api.Order.OrderProcedure, index: number) => {
      return h(NCheckbox, {
        checked: formModel.procedures[index].no_bundle,
        disabled: modalType.value === 'view',
        'onUpdate:checked': (value: boolean) => {
          formModel.procedures[index].no_bundle = value
        },
      })
    },
  },
])

defineExpose({ openModal })
</script>

<template>
  <NModal v-model:show="visible" preset="card" :title="title" class="w-1000px">
    <div v-if="modalType === 'add'" class="mb-4">
      <NSteps :current="currentStep">
        <NStep title="基础信息" description="填写订单基本信息" />
        <NStep title="款式数量" description="选择款式并配置数量" />
        <NStep title="工序清单" description="配置工序信息" />
      </NSteps>
    </div>

    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      :label-width="100"
    >
      <!-- 步骤1：基础信息 -->
      <div v-show="modalType !== 'add' || currentStep === 0">
        <NGrid :cols="2" :x-gap="18">
          <NFormItemGridItem path="contract_no" label="合同号" :span="2">
            <NInput v-model:value="formModel.contract_no" :disabled="modalType === 'view'" placeholder="请输入合同号" />
          </NFormItemGridItem>
          <NFormItemGridItem path="customer_id" label="客户" :span="2">
            <NSelect
              v-model:value="formModel.customer_id"
              :disabled="modalType === 'view'"
              :options="customerOptions"
              placeholder="请选择客户"
              filterable
              @update:value="onCustomerChange"
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="delivery_date" label="交货日期">
            <NDatePicker
              v-model:formatted-value="formModel.delivery_date"
              :disabled="modalType === 'view'"
              type="date"
              value-format="yyyy-MM-dd"
              placeholder="请选择交货日期"
              class="w-full"
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="order_type_id" label="订单类型">
            <NSelect
              v-model:value="formModel.order_type_id"
              :disabled="modalType === 'view'"
              :options="orderTypeOptions"
              placeholder="请选择订单类型"
              clearable
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="salesman_id" label="业务员" :span="2">
            <NSelect
              v-model:value="formModel.salesman_id"
              :disabled="modalType === 'view'"
              :options="salesmanOptions"
              placeholder="请选择业务员"
              clearable
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="remark" label="备注" :span="2">
            <NInput
              v-model:value="formModel.remark"
              :disabled="modalType === 'view'"
              type="textarea"
              :rows="3"
              placeholder="请输入备注"
            />
          </NFormItemGridItem>
        </NGrid>
      </div>

      <!-- 步骤2：款式数量 -->
      <div v-show="modalType !== 'add' || currentStep === 1">
        <NGrid :cols="2" :x-gap="18">
          <NFormItemGridItem path="style_id" label="款式" :span="2">
            <NSelect
              v-model:value="formModel.style_id"
              :disabled="modalType === 'view'"
              :options="styleOptions"
              placeholder="请选择款式"
              filterable
              @update:value="onStyleChange"
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="colors" label="颜色" :span="2">
            <NSelect
              v-model:value="formModel.colors"
              :disabled="modalType === 'view'"
              :options="colorOptions"
              placeholder="请选择颜色，可直接输入新颜色"
              multiple
              filterable
              tag
              @update:value="onColorsChange"
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="sizes" label="尺码" :span="2">
            <NSelect
              v-model:value="formModel.sizes"
              :disabled="modalType === 'view'"
              :options="sizeOptions"
              placeholder="请选择尺码，可直接输入新尺码"
              multiple
              filterable
              tag
              @update:value="onSizesChange"
            />
          </NFormItemGridItem>
          <NFormItemGridItem path="unit_price" label="单价">
            <NInputNumber v-model:value="formModel.unit_price" :disabled="modalType === 'view'" placeholder="单价" class="w-full" :precision="2" />
          </NFormItemGridItem>
          <NFormItemGridItem path="quantity" label="总数量">
            <NInputNumber v-model:value="formModel.quantity" disabled placeholder="总数量" class="w-full" />
          </NFormItemGridItem>
        </NGrid>

        <div v-if="formModel.items.length > 0" class="mt-4">
          <NDivider>颜色尺码数量配置</NDivider>
          <NDataTable
            :columns="itemColumns"
            :data="formModel.items"
            size="small"
            max-height="300px"
          />
        </div>
      </div>

      <!-- 步骤3：工序清单 -->
      <div v-show="modalType !== 'add' || currentStep === 2">
        <NDivider>工序清单</NDivider>
        <NDataTable
          :columns="procedureColumns"
          :data="formModel.procedures"
          size="small"
          max-height="300px"
        />
      </div>
    </NForm>

    <template v-if="modalType !== 'view'" #footer>
      <NSpace justify="end">
        <NButton v-if="modalType === 'add' && currentStep > 0" @click="currentStep--">
          上一步
        </NButton>
        <NButton @click="hideModal()">
          取消
        </NButton>
        <NButton
          v-if="modalType === 'add' && currentStep === 0"
          type="primary"
          :loading="loadingSubmit"
          @click="handleStep1"
        >
          下一步
        </NButton>
        <NButton
          v-if="modalType === 'add' && currentStep === 1"
          type="primary"
          :loading="loadingSubmit"
          @click="handleStep2"
        >
          下一步
        </NButton>
        <NButton
          v-if="modalType === 'add' && currentStep === 2"
          type="primary"
          :loading="loadingSubmit"
          @click="handleStep3"
        >
          完成
        </NButton>
        <NButton
          v-if="modalType === 'edit'"
          type="primary"
          :loading="loadingSubmit"
          @click="handleEditSubmit"
        >
          确定
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
