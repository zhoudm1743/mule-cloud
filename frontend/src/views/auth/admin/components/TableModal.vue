<script setup lang="ts">
import { useBoolean } from '@/hooks'
import { createAdmin, fetchAllTenants, fetchTenantRoles, updateAdmin } from '@/service'
import { useAuthStore } from '@/store'

interface Props {
  modalName?: string
}

const {
  modalName = '',
} = defineProps<Props>()

const emit = defineEmits<{
  open: []
  close: []
  refresh: []
}>()

const { bool: modalVisible, setTrue: showModal, setFalse: hiddenModal } = useBoolean(false)

const { bool: submitLoading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)
const authStore = useAuthStore()

const formDefault: Api.Admin.CreateRequest = {
  phone: '',
  password: '',
  nickname: '',
  email: '',
  avatar: '',
  roles: [],
  tenant_id: '',
  status: 1,
}
const formModel = ref<Api.Admin.CreateRequest & { id?: string }>({ ...formDefault })

// 是否为系统超管（角色包含 'super' 且没有租户ID）
const isSystemAdmin = computed(() => {
  const userInfo = authStore.userInfo
  if (!userInfo || !userInfo.role || !Array.isArray(userInfo.role)) {
    return false
  }
  // 系统超管：角色包含 'super' 且 tenant_id 为空
  return userInfo.role.includes('super') && !userInfo.tenant_id
})

// 租户列表
const tenantOptions = ref<Api.Tenant.TenantInfo[]>([])
async function getTenantList() {
  if (!isSystemAdmin.value) return
  
  try {
    const res = await fetchAllTenants()
    tenantOptions.value = res.data?.tenants || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取租户列表失败')
  }
}

type ModalType = 'add' | 'view' | 'edit'
const modalType = shallowRef<ModalType>('add')
const modalTitle = computed(() => {
  const titleMap: Record<ModalType, string> = {
    add: '添加',
    view: '查看',
    edit: '编辑',
  }
  return `${titleMap[modalType.value]}${modalName}`
})

async function openModal(type: ModalType = 'add', data?: Api.Admin.AdminInfo) {
  emit('open')
  modalType.value = type
  showModal()
  
  const handlers = {
    async add() {
      formModel.value = { ...formDefault }
      
      // 系统超管：加载租户列表
      if (isSystemAdmin.value) {
        await getTenantList()
      }
      else {
        // 租户超管：自动设置为自己的租户
        formModel.value.tenant_id = authStore.userInfo?.tenant_id || ''
        // 立即加载角色列表
        if (formModel.value.tenant_id) {
          await getRoleList(formModel.value.tenant_id)
        }
      }
    },
    async view() {
      if (!data)
        return
      formModel.value = { ...data, password: '' }
      
      // 加载对应租户的角色列表
      if (data.tenant_id) {
        await getRoleList(data.tenant_id)
      }
    },
    async edit() {
      if (!data)
        return
      formModel.value = { ...data, password: '' }
      
      // 系统超管编辑时：加载租户列表
      if (isSystemAdmin.value) {
        await getTenantList()
      }
      
      // 加载对应租户的角色列表
      if (data.tenant_id) {
        await getRoleList(data.tenant_id)
      }
    },
  }
  await handlers[type]()
}

function closeModal() {
  hiddenModal()
  endLoading()
  emit('close')
}

defineExpose({
  openModal,
})

const formRef = ref()
async function submitModal() {
  const handlers = {
    async add() {
      try {
        await createAdmin(formModel.value)
        window.$message.success('创建成功')
        emit('refresh')
        return true
      }
      catch (error: any) {
        window.$message.error(error.message || '创建失败')
        return false
      }
    },
    async edit() {
      try {
        const { id, password, ...updateData } = formModel.value
        if (!id) {
          window.$message.error('缺少管理员ID')
          return false
        }
        // 如果密码不为空，才包含密码字段
        const submitData = password ? { ...updateData, password } : updateData
        await updateAdmin(id, submitData)
        window.$message.success('更新成功')
        emit('refresh')
        return true
      }
      catch (error: any) {
        window.$message.error(error.message || '更新失败')
        return false
      }
    },
    async view() {
      return true
    },
  }
  await formRef.value?.validate()
  startLoading()
  const success = await handlers[modalType.value]()
  if (success)
    closeModal()
  else
    endLoading()
}

const rules = {
  phone: {
    required: true,
    message: '请输入手机号',
    trigger: 'blur',
  },
  password: {
    required: modalType.value === 'add',
    message: '请输入密码',
    trigger: 'blur',
  },
}

const roleOptions = ref<Api.Role.RoleInfo[]>([])
async function getRoleList(tenantId?: string) {
  if (!tenantId) {
    roleOptions.value = []
    return
  }
  
  try {
    const res = await fetchTenantRoles(tenantId)
    roleOptions.value = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取角色列表失败')
  }
}

// 监听租户变化，动态加载角色
watch(() => formModel.value.tenant_id, async (newTenantId) => {
  if (modalType.value === 'add' || modalType.value === 'edit') {
    if (newTenantId) {
      await getRoleList(newTenantId)
      // 清空已选角色（因为不同租户的角色不同）
      formModel.value.roles = []
    }
    else {
      roleOptions.value = []
      formModel.value.roles = []
    }
  }
})
</script>

<template>
  <NModal
    v-model:show="modalVisible"
    :mask-closable="false"
    preset="card"
    :title="modalTitle"
    class="w-700px"
    :segmented="{
      content: true,
      action: true,
    }"
  >
    <NForm ref="formRef" :rules="rules" label-placement="left" :model="formModel" :label-width="100" :disabled="modalType === 'view'">
      <NGrid :cols="2" :x-gap="18">
        <!-- 租户选择（仅系统超管可见） -->
        <NFormItemGridItem 
          v-if="isSystemAdmin" 
          :span="2" 
          label="所属租户" 
          path="tenant_id"
        >
          <NSelect
            v-model:value="formModel.tenant_id"
            :options="tenantOptions"
            label-field="name"
            value-field="id"
            placeholder="请选择租户（为空表示系统级用户）"
            clearable
            filterable
          />
        </NFormItemGridItem>
        
        <!-- 租户信息显示（非系统超管只读） -->
        <NFormItemGridItem 
          v-else-if="formModel.tenant_id"
          :span="2" 
          label="所属租户"
        >
          <NInput :value="formModel.tenant_id" readonly placeholder="本租户" />
        </NFormItemGridItem>
        
        <NFormItemGridItem :span="1" label="手机号" path="phone">
          <NInput v-model:value="formModel.phone" placeholder="请输入手机号" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="1" label="密码" path="password">
          <NInput v-model:value="formModel.password" type="password" placeholder="请输入密码" show-password-on="click" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="1" label="昵称" path="nickname">
          <NInput v-model:value="formModel.nickname" placeholder="请输入昵称" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="1" label="邮箱" path="email">
          <NInput v-model:value="formModel.email" placeholder="请输入邮箱" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" label="角色" path="roles">
          <NSelect
            v-model:value="formModel.roles" multiple filterable
            label-field="name"
            value-field="id"
            :options="roleOptions"
            placeholder="请选择角色"
          />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" label="头像URL" path="avatar">
          <NInput v-model:value="formModel.avatar" placeholder="请输入头像URL" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="1" label="状态" path="status">
          <NSwitch
            v-model:value="formModel.status"
            :checked-value="1" :unchecked-value="0"
          >
            <template #checked>
              启用
            </template>
            <template #unchecked>
              禁用
            </template>
          </NSwitch>
        </NFormItemGridItem>
      </NGrid>
    </NForm>
    <template #action>
      <NSpace justify="center">
        <NButton @click="closeModal">
          取消
        </NButton>
        <NButton type="primary" :loading="submitLoading" @click="submitModal">
          提交
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>
