<script setup lang="ts">
import axios from 'axios'

interface Props {
  modelValue?: boolean
}

interface Emits {
  (e: 'update:modelValue', val: boolean): void
  (e: 'success'): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
})
const emit = defineEmits<Emits>()

const visible = computed({
  get() {
    return props.modelValue
  },
  set(val) {
    emit('update:modelValue', val)
  },
})

const modalType = ref<'add' | 'edit'>('add')
const loading = ref(false)
const formRef = ref()

const initialFormModel = {
  id: '',
  name: '',
  gender: 0,
  phone: '',
  id_card_no: '',
  job_number: '',
  department: '',
  position: '',
  workshop: '',
  team: '',
  employed_at: null as number | null,
  status: 'active',
}

const formModel = ref({ ...initialFormModel })

const rules = {
  name: { required: true, message: '请输入姓名', trigger: 'blur' },
  job_number: { required: true, message: '请输入工号', trigger: 'blur' },
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur' },
  ],
}

const genderOptions = [
  { label: '未知', value: 0 },
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]

const statusOptions = [
  { label: '在职', value: 'active' },
  { label: '试用期', value: 'probation' },
  { label: '离职', value: 'inactive' },
]

function openModal(type: 'add' | 'edit', data?: any) {
  modalType.value = type
  visible.value = true

  if (type === 'edit' && data) {
    formModel.value = {
      ...data,
      employed_at: data.employed_at ? data.employed_at * 1000 : null,
    }
  }
  else {
    formModel.value = { ...initialFormModel }
  }
}

async function handleSubmit() {
  await formRef.value?.validate()

  loading.value = true
  try {
    const submitData = {
      ...formModel.value,
      employed_at: formModel.value.employed_at ? Math.floor(formModel.value.employed_at / 1000) : 0,
    }

    if (modalType.value === 'add') {
      await axios.post('/api/miniapp/member', submitData)
      window.$message.success('新增成功')
    }
    else {
      await axios.put(`/api/miniapp/member/${formModel.value.id}`, submitData)
      window.$message.success('更新成功')
    }

    visible.value = false
    emit('success')
  }
  catch (error: any) {
    window.$message.error(error.response?.data?.message || error.message || '操作失败')
  }
  finally {
    loading.value = false
  }
}

defineExpose({ openModal })
</script>

<template>
  <n-modal
    v-model:show="visible"
    :title="modalType === 'add' ? '新增员工' : '编辑员工'"
    preset="card"
    class="w-800px"
  >
    <n-form
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      label-width="100"
    >
      <n-tabs type="line" animated>
        <n-tab-pane name="basic" tab="基本信息">
          <n-form-item label="姓名" path="name">
            <n-input v-model:value="formModel.name" placeholder="请输入姓名" />
          </n-form-item>
          <n-form-item label="性别" path="gender">
            <n-select v-model:value="formModel.gender" :options="genderOptions" />
          </n-form-item>
          <n-form-item label="手机号" path="phone">
            <n-input v-model:value="formModel.phone" placeholder="请输入手机号" />
          </n-form-item>
          <n-form-item label="身份证号" path="id_card_no">
            <n-input v-model:value="formModel.id_card_no" placeholder="请输入身份证号" />
          </n-form-item>
        </n-tab-pane>

        <n-tab-pane name="work" tab="工作信息">
          <n-form-item label="工号" path="job_number">
            <n-input v-model:value="formModel.job_number" placeholder="请输入工号" />
          </n-form-item>
          <n-form-item label="部门" path="department">
            <n-input v-model:value="formModel.department" placeholder="请输入部门" />
          </n-form-item>
          <n-form-item label="岗位" path="position">
            <n-input v-model:value="formModel.position" placeholder="请输入岗位" />
          </n-form-item>
          <n-form-item label="车间" path="workshop">
            <n-input v-model:value="formModel.workshop" placeholder="请输入车间" />
          </n-form-item>
          <n-form-item label="班组" path="team">
            <n-input v-model:value="formModel.team" placeholder="请输入班组" />
          </n-form-item>
          <n-form-item label="入职日期" path="employed_at">
            <n-date-picker
              v-model:value="formModel.employed_at"
              type="date"
              placeholder="请选择入职日期"
              class="w-full"
            />
          </n-form-item>
          <n-form-item label="状态" path="status">
            <n-select v-model:value="formModel.status" :options="statusOptions" />
          </n-form-item>
        </n-tab-pane>
      </n-tabs>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="visible = false">
          取消
        </n-button>
        <n-button type="primary" :loading="loading" @click="handleSubmit">
          确定
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

