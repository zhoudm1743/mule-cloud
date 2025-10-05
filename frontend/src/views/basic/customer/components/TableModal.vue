<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { createCustomer, updateCustomer } from '@/service'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'CustomerModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建客户',
    edit: '编辑客户',
    view: '查看客户',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

const formDefault = (): Api.Basic.CreateRequest & { id?: string } => ({
  id: '',
  value: '',
  remark: '',
})

const formModel = reactive(formDefault())

const rules: any = {
  value: { required: true, message: '请输入客户名称', trigger: 'blur' },
}

const formRef = ref()
const loadingSubmit = ref(false)

function openModal(type: ModalType, data?: Api.Basic.BasicInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())

  if (type === 'view' || type === 'edit') {
    if (data) {
      Object.assign(formModel, data)
    }
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  loadingSubmit.value = true

  try {
    if (modalType.value === 'add') {
      await createCustomer(formModel as Api.Basic.CreateRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      if (!formModel.id) {
        window.$message.error('缺少ID')
        return
      }
      // 只发送需要更新的字段
      const updateData: Api.Basic.UpdateRequest = {
        value: formModel.value,
        remark: formModel.remark,
      }
      await updateCustomer(formModel.id, updateData)
      window.$message.success('更新成功')
    }

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

defineExpose({ openModal })
</script>

<template>
  <NModal v-model:show="visible" preset="card" :title="title" class="w-600px">
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      :label-width="100"
    >
      <NGrid :cols="1" :x-gap="18">
        <NFormItemGridItem path="value" label="客户名称">
          <NInput v-model:value="formModel.value" :disabled="modalType === 'view'" placeholder="请输入客户名称" />
        </NFormItemGridItem>
        <NFormItemGridItem path="remark" label="备注">
          <NInput
            v-model:value="formModel.remark"
            :disabled="modalType === 'view'"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
          />
        </NFormItemGridItem>
      </NGrid>
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

