<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { NButton, NDivider, NFormItemGridItem, NGrid, NInput, NInputNumber, NModal, NSpace } from 'naive-ui'
import { createCuttingBatch, fetchOrderDetail } from '@/service/api/order'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'BatchModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '制菲（创建裁剪批次）',
    edit: '编辑批次',
    view: '查看批次',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')
const currentTask = ref<Api.Order.CuttingTaskInfo | null>(null)

// 表单数据
const formDefault = () => ({
  task_id: '',
  bed_no: '',
  bundle_no: '',
  color: '',
  layer_count: 1,
  size_details: [] as Api.Order.SizeDetail[],
})

const formModel = reactive(formDefault())

const rules: any = {
  bed_no: { required: true, message: '请输入床号', trigger: 'blur' },
  bundle_no: { required: true, message: '请输入扎号', trigger: 'blur' },
  color: { required: true, message: '请输入颜色', trigger: 'blur' },
  layer_count: { required: true, type: 'number', message: '请输入拉布层数', trigger: 'blur' },
}

const formRef = ref()
const loadingSubmit = ref(false)

// 初始化尺码明细（从任务获取订单信息）
async function initializeSizeDetails(task: Api.Order.CuttingTaskInfo) {
  // 从订单API获取订单详情，然后提取尺码信息
  try {
    const { data } = await fetchOrderDetail(task.order_id)
    if (data?.order) {
      // 使用订单的尺码创建初始尺码明细
      formModel.size_details = (data.order.sizes || []).map(size => ({
        size: size,
        quantity: 0,
      }))
    }
  }
  catch (error) {
    // 如果获取失败，提供空数组让用户手动添加
    formModel.size_details = []
    window.$message.warning('无法获取订单尺码信息，请手动添加')
  }
}

// 添加尺码
function handleAddSize() {
  formModel.size_details.push({
    size: '',
    quantity: 0,
  })
}

// 删除尺码
function handleRemoveSize(index: number) {
  formModel.size_details.splice(index, 1)
}

// 计算总件数
const totalPieces = computed(() => {
  return formModel.size_details.reduce((sum, item) => {
    return sum + (item.quantity * formModel.layer_count)
  }, 0)
})

async function openModal(type: ModalType, task?: Api.Order.CuttingTaskInfo) {
  showModal()
  modalType.value = type
  currentTask.value = task || null

  if (type === 'add' && task) {
    Object.assign(formModel, formDefault())
    formModel.task_id = task.id
    await initializeSizeDetails(task)
  }
}

async function handleSubmit() {
  await formRef.value?.validate()

  // 验证尺码明细
  if (formModel.size_details.length === 0) {
    window.$message.error('请添加至少一个尺码')
    return
  }

  const hasInvalidSize = formModel.size_details.some(item => !item.size || item.quantity <= 0)
  if (hasInvalidSize) {
    window.$message.error('请填写完整的尺码信息')
    return
  }

  loadingSubmit.value = true
  try {
    await createCuttingBatch({
      task_id: formModel.task_id,
      bed_no: formModel.bed_no,
      bundle_no: formModel.bundle_no,
      color: formModel.color,
      layer_count: formModel.layer_count,
      size_details: formModel.size_details,
    })

    window.$message.success('制菲成功')
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

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="visible"
    :title="title"
    preset="card"
    class="w-800px"
    :mask-closable="false"
  >
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="100"
      :disabled="modalType === 'view'"
    >
      <NGrid :cols="2" :x-gap="18">
        <NFormItemGridItem path="bed_no" label="床号">
          <NInput v-model:value="formModel.bed_no" placeholder="请输入床号" />
        </NFormItemGridItem>
        <NFormItemGridItem path="bundle_no" label="扎号">
          <NInput v-model:value="formModel.bundle_no" placeholder="请输入扎号" />
        </NFormItemGridItem>
        <NFormItemGridItem path="color" label="颜色" :span="2">
          <NInput v-model:value="formModel.color" placeholder="请输入颜色" />
        </NFormItemGridItem>
        <NFormItemGridItem path="layer_count" label="拉布层数" :span="2">
          <NInputNumber
            v-model:value="formModel.layer_count"
            placeholder="拉布层数"
            :min="1"
            class="w-full"
          />
        </NFormItemGridItem>
      </NGrid>

      <NDivider>尺码明细</NDivider>

      <div class="mb-4">
        <NButton @click="handleAddSize">
          添加尺码
        </NButton>
      </div>

      <div v-for="(item, index) in formModel.size_details" :key="index" class="mb-4">
        <NGrid :cols="4" :x-gap="12">
          <NFormItemGridItem :path="`size_details[${index}].size`" label="尺码">
            <NInput v-model:value="item.size" placeholder="尺码" />
          </NFormItemGridItem>
          <NFormItemGridItem :path="`size_details[${index}].quantity`" label="每层数量">
            <NInputNumber v-model:value="item.quantity" placeholder="数量" :min="0" class="w-full" />
          </NFormItemGridItem>
          <NFormItemGridItem label="小计">
            <NInput :value="`${item.quantity * formModel.layer_count} 件`" readonly />
          </NFormItemGridItem>
          <NFormItemGridItem label="">
            <NButton type="error" @click="handleRemoveSize(index)">
              删除
            </NButton>
          </NFormItemGridItem>
        </NGrid>
      </div>

      <NDivider />

      <div class="text-lg font-bold">
        总件数：{{ totalPieces }} 件
      </div>
    </NForm>

    <template v-if="modalType !== 'view'" #footer>
      <NSpace justify="end">
        <NButton @click="hideModal()">
          取消
        </NButton>
        <NButton type="primary" :loading="loadingSubmit" @click="handleSubmit">
          确定
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
