<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { assignTenantMenus, createTenant, fetchAllMenus, fetchTenantMenus, updateTenant } from '@/service'
import { useBoolean } from '@/hooks'
import { arrayToTree } from '@/utils'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'TenantModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view' | 'assignMenus'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建租户',
    edit: '编辑租户',
    view: '查看租户',
    assignMenus: '分配权限菜单',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

const formDefault = (): Api.Tenant.CreateRequest & { id?: string, menus?: string[] } => ({
  id: '',
  code: '',
  name: '',
  contact: '',
  phone: '',
  email: '',
  status: 1,
  menus: [],
})

const formModel = reactive(formDefault())

const rules: any = {
  code: { required: true, message: '请输入租户代码', trigger: 'blur' },
  name: { required: true, message: '请输入租户名称', trigger: 'blur' },
}

const formRef = ref()
const loadingSubmit = ref(false)
const allMenus = ref<Api.Menu.MenuItem[]>([])

async function getAllMenus() {
  try {
    const res = await fetchAllMenus()
    allMenus.value = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取菜单列表失败')
  }
}

async function getTenantMenus(id: string) {
  try {
    const res = await fetchTenantMenus(id)
    // 后端直接返回数组，不是 { menus: [...] } 格式
    formModel.menus = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取租户菜单权限失败')
  }
}

const menuTreeOptions = computed(() => {
  // 构建树形结构 - 使用 name 作为唯一标识
  // 需要先建立 id 到 name 的映射
  const idToNameMap = new Map<string, string>()
  allMenus.value.forEach(menu => {
    idToNameMap.set(menu.id, menu.name)
  })
  
  const menuItems = allMenus.value.map(menu => ({
    key: menu.name,        // 使用 name 作为 key（权限标识）
    label: menu.title || menu.name,
    id: menu.name,         // 使用 name 作为 id（用于树形结构）
    pid: menu.pid ? idToNameMap.get(menu.pid) : null, // 将 pid 也转换为 name
  }))
  return arrayToTree(menuItems)
})

function openModal(type: ModalType, data?: Api.Tenant.TenantInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())

  if (type === 'add') {
    // 新建模式
  }
  else if (type === 'view' || type === 'edit') {
    // 查看或编辑模式
    if (data) {
      Object.assign(formModel, data)
    }
  }
  else if (type === 'assignMenus') {
    // 分配菜单模式
    if (data) {
      formModel.id = data.id
      formModel.name = data.name
      getAllMenus()
      getTenantMenus(data.id)
    }
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  loadingSubmit.value = true

  try {
    if (modalType.value === 'add') {
      await createTenant(formModel as Api.Tenant.CreateRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      const { id, menus, ...updateData } = formModel
      if (!id) {
        window.$message.error('缺少租户ID')
        return
      }
      await updateTenant(id, updateData)
      window.$message.success('更新成功')
    }
    else if (modalType.value === 'assignMenus') {
      if (!formModel.id) {
        window.$message.error('缺少租户ID')
        return
      }
      
      // 自动添加父级菜单（确保菜单层级完整）
      const menusWithParents = new Set<string>(formModel.menus || [])
      
      // 为每个选中的菜单，递归添加其所有父级菜单
      const addParentMenus = (menuName: string) => {
        const menu = allMenus.value.find(m => m.name === menuName)
        if (menu && menu.pid) {
          // 找到父菜单
          const parentMenu = allMenus.value.find(m => m.id === menu.pid)
          if (parentMenu && parentMenu.name) {
            menusWithParents.add(parentMenu.name)
            // 递归添加父级的父级
            addParentMenus(parentMenu.name)
          }
        }
      }
      
      // 为所有选中的菜单添加父级
      ;(formModel.menus || []).forEach(menuName => {
        addParentMenus(menuName)
      })
      
      // 更新菜单列表
      const finalMenus = Array.from(menusWithParents)
      if (finalMenus.length !== (formModel.menus || []).length) {
        console.log('[租户菜单补全] 原始菜单:', formModel.menus)
        console.log('[租户菜单补全] 补全后菜单:', finalMenus)
      }
      
      await assignTenantMenus(formModel.id, { menus: finalMenus })
      window.$message.success('分配权限成功')
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
  <NModal v-model:show="visible" preset="card" :title="title" class="w-700px">
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      :label-width="100"
    >
      <NGrid v-if="modalType === 'assignMenus'" :cols="1" :x-gap="18">
        <NFormItemGridItem label="租户名称">
          <span>{{ formModel.name }}</span>
        </NFormItemGridItem>
        <NFormItemGridItem path="menus" label="权限菜单">
          <div class="h-400px overflow-y-auto">
            <NTree
              v-model:checked-keys="formModel.menus"
              :data="menuTreeOptions"
              checkable
              cascade
              key-field="key"
              label-field="label"
              children-field="children"
              :default-expand-all="true"
            />
          </div>
        </NFormItemGridItem>
      </NGrid>

      <NGrid v-else :cols="2" :x-gap="18">
        <NFormItemGridItem :span="2" path="code" label="租户代码">
          <NInput v-model:value="formModel.code" :disabled="modalType === 'view'" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="name" label="租户名称">
          <NInput v-model:value="formModel.name" :disabled="modalType === 'view'" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="contact" label="联系人">
          <NInput v-model:value="formModel.contact" :disabled="modalType === 'view'" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="phone" label="联系电话">
          <NInput v-model:value="formModel.phone" :disabled="modalType === 'view'" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="email" label="联系邮箱">
          <NInput v-model:value="formModel.email" :disabled="modalType === 'view'" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="status" label="状态">
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
