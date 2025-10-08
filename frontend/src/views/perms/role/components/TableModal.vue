<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { assignRoleMenus, createRole, fetchAllMenus, fetchAllTenants, fetchRoleMenus, fetchTenantMenus, updateRole } from '@/service'
import { useBoolean } from '@/hooks'
import { useAuthStore } from '@/store'
import { arrayToTree } from '@/utils'
import MenuPermissionTree from './MenuPermissionTree.vue'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'RoleModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)
const authStore = useAuthStore()

// 是否为系统超管（角色包含 'super' 且没有租户ID）
const isSystemAdmin = computed(() => {
  const userInfo = authStore.userInfo
  if (!userInfo || !userInfo.role || !Array.isArray(userInfo.role)) {
    return false
  }
  // 系统超管：角色包含 'super' 且 tenant_id 为空
  return userInfo.role.includes('super') && !userInfo.tenant_id
})

type ModalType = 'add' | 'edit' | 'view' | 'assignMenus'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建角色',
    edit: '编辑角色',
    view: '查看角色',
    assignMenus: '分配菜单权限',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

const formDefault = (): Api.Role.CreateRequest & { id?: string, menus?: string[] } => ({
  id: '',
  tenant_id: '',
  name: '',
  code: '',
  description: '',
  status: 1,
  menus: [],
})

const formModel = reactive(formDefault())

const rules = computed(() => ({
  // 租户超管必须选择租户，系统超管可以不选（创建系统级角色）
  tenant_id: { 
    required: !isSystemAdmin.value, 
    message: '请选择租户', 
    trigger: 'change' 
  },
  name: { required: true, message: '请输入角色名称', trigger: 'blur' },
  code: { required: true, message: '请输入角色代码', trigger: 'blur' },
}))

const formRef = ref()
const loadingSubmit = ref(false)
const allTenants = ref<Api.Tenant.TenantInfo[]>([])
const allMenus = ref<Api.Menu.MenuItem[]>([])
const tenantMenus = ref<string[]>([])
const menuPermissions = ref<Record<string, Record<string, boolean>>>({}) // { "admin": { "read": true, "create": true } }

async function getAllTenants() {
  try {
    const res = await fetchAllTenants()
    // 后端返回的是 { tenants: [...], total: number } 格式
    allTenants.value = res.data?.tenants || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取租户列表失败')
  }
}

async function getAllMenus() {
  try {
    const res = await fetchAllMenus()
    allMenus.value = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取菜单列表失败')
  }
}

async function getTenantMenus(tenantId: string) {
  if (!tenantId)
    return
  try {
    const res = await fetchTenantMenus(tenantId)
    // 后端直接返回数组，不是 { menus: [...] } 格式
    tenantMenus.value = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取租户菜单权限失败')
  }
}

async function getRoleMenus(roleId: string) {
  try {
    const res = await fetchRoleMenus(roleId)
    // 后端直接返回数组，不是 { menus: [...] } 格式
    formModel.menus = res.data || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取角色菜单权限失败')
  }
}

// 初始化菜单权限
function initMenuPermissions(roleData?: Api.Role.RoleInfo) {
  menuPermissions.value = {}
  
  // 遍历所有可用菜单
  formModel.menus?.forEach(menuName => {
    const menu = allMenus.value.find(m => m.name === menuName)
    if (!menu) return
    
    // 只有配置了 available_permissions 的菜单才初始化权限
    if (!menu.available_permissions || menu.available_permissions.length === 0) {
      return
    }
    
    // 初始化该菜单的权限对象
    menuPermissions.value[menuName] = {}
    
    menu.available_permissions.forEach(perm => {
      // 如果角色有 menu_permissions 数据，使用它；否则默认全选
      if (roleData?.menu_permissions && roleData.menu_permissions[menuName]) {
        menuPermissions.value[menuName][perm.action] = roleData.menu_permissions[menuName].includes(perm.action)
      } else {
        // 默认全选
        menuPermissions.value[menuName][perm.action] = true
      }
    })
  })
}

// 监听菜单选择变化
watch(() => formModel.menus, (newMenus) => {
  if (modalType.value !== 'assignMenus') return
  
  // 为新增的菜单添加默认权限
  newMenus?.forEach(menuName => {
    if (!menuPermissions.value[menuName]) {
      const menu = allMenus.value.find(m => m.name === menuName)
      
      // 只有配置了 available_permissions 的菜单才初始化权限
      if (menu?.available_permissions && menu.available_permissions.length > 0) {
        menuPermissions.value[menuName] = {}
        menu.available_permissions.forEach(perm => {
          menuPermissions.value[menuName][perm.action] = true // 默认全选
        })
      }
    }
  })
  
  // 删除未选中菜单的权限
  Object.keys(menuPermissions.value).forEach(menuName => {
    if (!newMenus?.includes(menuName)) {
      delete menuPermissions.value[menuName]
    }
  })
})

const tenantOptions = computed(() => {
  if (!Array.isArray(allTenants.value)) {
    return []
  }
  return allTenants.value.map(tenant => ({
    label: tenant.name,
    value: tenant.id,
  }))
})

const menuTreeOptions = computed(() => {
  // 如果是分配菜单模式，只显示租户拥有的菜单
  let availableMenus = allMenus.value
  
  if (modalType.value === 'assignMenus') {
    // 系统超管：可以看到所有菜单
    if (isSystemAdmin.value) {
      availableMenus = allMenus.value
    }
    // 租户超管：只能看到自己租户的菜单
    else if (tenantMenus.value.length > 0) {
      // 严格过滤：只保留租户拥有的菜单
      availableMenus = allMenus.value.filter(menu => tenantMenus.value.includes(menu.name))
    }
    // 如果租户菜单为空（异常情况），不显示任何菜单
    else {
      availableMenus = []
    }
  }

  // 如果没有可用菜单，直接返回空数组
  if (availableMenus.length === 0) {
    return []
  }

  // 构建树形结构 - 使用 name 作为唯一标识
  // 创建一个包含所有菜单的 name 映射（用于查找父节点）
  const allMenuNameMap = new Map<string, string>()
  allMenus.value.forEach(menu => {
    allMenuNameMap.set(menu.id, menu.name)
  })
  
  // 创建可用菜单的 name 集合（用于过滤父节点）
  const availableMenuNames = new Set(availableMenus.map(m => m.name))
  
  // 调试日志
  if (modalType.value === 'assignMenus' && !isSystemAdmin.value) {
    console.log('[菜单树] 可用菜单列表:', Array.from(availableMenuNames))
    console.log('[菜单树] 租户菜单列表:', tenantMenus.value)
  }
  
  const menuItems = availableMenus.map(menu => {
    // 获取父节点的 name
    let parentName = menu.pid ? allMenuNameMap.get(menu.pid) : null
    
    // 如果父节点不在可用菜单中，设置为 null（将此节点提升为根节点）
    if (parentName && !availableMenuNames.has(parentName)) {
      parentName = null
    }
    
    return {
      key: menu.name,        // 使用 name 作为 key（权限标识）
      label: menu.title || menu.name,
      id: menu.name,         // 使用 name 作为 id（用于树形结构）
      pid: parentName,       // 父节点的 name（如果父节点不可用则为 null）
      menuData: menu,        // 保存完整菜单数据，用于权限渲染
    }
  })
  
  return arrayToTree(menuItems)
})

async function onTenantChange(tenantId: string) {
  formModel.menus = []
  await getTenantMenus(tenantId)
}

async function openModal(type: ModalType, data?: Api.Role.RoleInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())

  if (type === 'add') {
    // 系统超管：可以选择任意租户
    if (isSystemAdmin.value) {
      await getAllTenants()
    }
    else {
      // 租户超管：自动设置为自己的租户
      formModel.tenant_id = authStore.userInfo?.tenant_id || ''
    }
  }
  else if (type === 'view' || type === 'edit') {
    if (data) {
      Object.assign(formModel, data)
    }
    if (type === 'edit') {
      // 系统超管编辑时：加载租户列表
      if (isSystemAdmin.value) {
        await getAllTenants()
      }
    }
  }
  else if (type === 'assignMenus') {
    if (data) {
      // 数据库物理隔离模式：每个租户只能访问自己数据库中的数据，无需额外检查
      
      formModel.id = data.id
      formModel.name = data.name
      formModel.tenant_id = data.tenant_id
      await getAllMenus()
      
      // 系统超管：加载角色所属租户的菜单
      // 租户超管：只加载自己租户的菜单
      const tenantIdToLoad = isSystemAdmin.value 
        ? data.tenant_id 
        : (authStore.userInfo?.tenant_id || '')
        
      if (tenantIdToLoad) {
        await getTenantMenus(tenantIdToLoad)

      }
      
      await getRoleMenus(data.id)
      // 初始化权限数据
      initMenuPermissions(data)
    }
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  loadingSubmit.value = true

  try {
    if (modalType.value === 'add') {
      await createRole(formModel as Api.Role.CreateRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      const { id, menus, ...updateData } = formModel
      if (!id) {
        window.$message.error('缺少角色ID')
        return
      }
      await updateRole(id, updateData)
      window.$message.success('更新成功')
    }
    else if (modalType.value === 'assignMenus') {
      if (!formModel.id) {
        window.$message.error('缺少角色ID')
        return
      }
      
      // 自动添加父级菜单（仅在租户拥有的菜单范围内）
      const menusWithParents = new Set<string>(formModel.menus || [])
      
      // 获取租户菜单范围（系统超管使用所有菜单，租户超管使用租户菜单）
      const availableMenus = isSystemAdmin.value 
        ? allMenus.value.map(m => m.name)
        : tenantMenus.value
      
      // 为每个选中的菜单，递归添加其所有父级菜单（但仅限租户拥有的）
      const addParentMenus = (menuName: string) => {
        const menu = allMenus.value.find(m => m.name === menuName)
        if (menu && menu.pid) {
          // 找到父菜单
          const parentMenu = allMenus.value.find(m => m.id === menu.pid)
          if (parentMenu && parentMenu.name) {
            // 只有当父菜单在可用范围内时才添加
            if (availableMenus.includes(parentMenu.name)) {
              menusWithParents.add(parentMenu.name)
              // 递归添加父级的父级
              addParentMenus(parentMenu.name)
            }
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
        console.log('[菜单补全] 原始菜单:', formModel.menus)
        console.log('[菜单补全] 补全后菜单:', finalMenus)
        console.log('[菜单补全] 可用菜单范围:', availableMenus)
      }
      formModel.menus = finalMenus
      
      // 转换 menuPermissions 为后端格式
      const menu_permissions: Record<string, string[]> = {}
      Object.entries(menuPermissions.value).forEach(([menuName, perms]) => {
        menu_permissions[menuName] = Object.entries(perms)
          .filter(([_, enabled]) => enabled)
          .map(([action, _]) => action)
      })
      
      await assignRoleMenus(formModel.id, { 
        menus: formModel.menus || [],
        menu_permissions,
      })
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
        <NFormItemGridItem label="角色名称">
          <span class="font-medium">{{ formModel.name }}</span>
        </NFormItemGridItem>
        <NFormItemGridItem :span="1" label="菜单权限配置">
          <div class="h-500px overflow-y-auto pr-2">
            <MenuPermissionTree
              v-model:selected-menus="formModel.menus!"
              v-model:permissions="menuPermissions"
              :tree-data="menuTreeOptions"
              :all-menus="allMenus"
            />
          </div>
        </NFormItemGridItem>
      </NGrid>

      <NGrid v-else :cols="2" :x-gap="18">
        <!-- 系统超管：可以选择租户（可为空创建系统级角色） -->
        <NFormItemGridItem 
          v-if="isSystemAdmin" 
          :span="2" 
          path="tenant_id" 
          label="所属租户"
        >
          <NSelect
            v-model:value="formModel.tenant_id"
            :options="tenantOptions"
            :disabled="modalType === 'view'"
            placeholder="请选择租户（为空表示系统级角色）"
            clearable
            filterable
            @update:value="onTenantChange"
          />
        </NFormItemGridItem>
        
        <!-- 租户超管：只显示当前租户（不可修改） -->
        <NFormItemGridItem 
          v-else 
          :span="2" 
          label="所属租户"
        >
          <NInput 
            :value="formModel.tenant_id" 
            readonly 
            placeholder="本租户" 
          />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="code" label="角色代码">
          <NInput v-model:value="formModel.code" :disabled="modalType === 'view'" placeholder="如：admin" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="name" label="角色名称">
          <NInput v-model:value="formModel.name" :disabled="modalType === 'view'" placeholder="如：系统管理员" />
        </NFormItemGridItem>
        <NFormItemGridItem :span="2" path="description" label="角色描述">
          <NInput
            v-model:value="formModel.description"
            type="textarea"
            :disabled="modalType === 'view'"
            placeholder="角色描述"
          />
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
