<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { createDepartment, fetchAllDepartments, updateDepartment } from '@/service'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'DepartmentModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建部门',
    edit: '编辑部门',
    view: '查看部门',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

const formDefault = (): Api.Department.CreateRequest & { id?: string } => ({
  id: '',
  name: '',
  code: '',
  parent_id: '',
  status: 1,
})

const formModel = reactive(formDefault())

const rules = computed(() => ({
  name: { required: true, message: '请输入部门名称', trigger: 'blur' },
  code: { required: true, message: '请输入部门编码', trigger: 'blur' },
}))

const formRef = ref()
const loadingSubmit = ref(false)
const allDepartments = ref<Api.Department.DepartmentInfo[]>([])

async function getAllDepartments() {
  try {
    const res = await fetchAllDepartments()
    allDepartments.value = res.data?.departments || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取部门列表失败')
  }
}

const departmentOptions = computed(() => {
  if (!Array.isArray(allDepartments.value)) {
    return []
  }
  return allDepartments.value
    .filter(dept => dept.id !== formModel.id) // 排除自己，避免循环引用
    .map(dept => ({
      label: dept.name,
      value: dept.id,
    }))
})

async function openModal(type: ModalType, data?: Api.Department.DepartmentInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())

  // 加载所有部门用于父部门选择
  await getAllDepartments()

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
      await createDepartment(formModel as Api.Department.CreateRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      const { id, ...updateData } = formModel
      if (!id) {
        window.$message.error('缺少部门ID')
        return
      }
      await updateDepartment(id, updateData)
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
        <NFormItemGridItem path="name" label="部门名称">
          <NInput v-model:value="formModel.name" :disabled="modalType === 'view'" placeholder="请输入部门名称" />
        </NFormItemGridItem>
        <NFormItemGridItem path="code" label="部门编码">
          <NInput v-model:value="formModel.code" :disabled="modalType === 'view'" placeholder="请输入部门编码" />
        </NFormItemGridItem>
        <NFormItemGridItem path="parent_id" label="父部门">
          <NSelect
            v-model:value="formModel.parent_id"
            :options="departmentOptions"
            :disabled="modalType === 'view'"
            placeholder="请选择父部门（为空表示顶级部门）"
            clearable
            filterable
          />
        </NFormItemGridItem>
        <NFormItemGridItem path="status" label="状态">
          <NRadioGroup v-model:value="formModel.status" :disabled="modalType === 'view'">
            <NRadio :value="1">
              启用
            </NRadio>
            <NRadio :value="0">
              禁用
            </NRadio>
          </NRadioGroup>
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

