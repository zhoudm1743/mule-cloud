<script setup lang="ts">
import type {
  FormItemRule,
} from 'naive-ui'
import HelpInfo from '@/components/common/HelpInfo.vue'
import { Regex } from '@/constants'
import { useBoolean } from '@/hooks'
import { createMenu, updateMenu } from '@/service'

interface Props {
  modalName?: string
  allRoutes: Api.Menu.MenuItem[]
}

const {
  modalName = '',
  allRoutes,
} = defineProps<Props>()

const emit = defineEmits<{
  open: []
  close: []
  refresh: []
}>()

const { bool: modalVisible, setTrue: showModal, setFalse: hiddenModal } = useBoolean(false)
const { bool: submitLoading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const formDefault: Api.Menu.MenuItem = {
  name: '',
  path: '',
  id: '',
  pid: null,
  title: '',
  componentPath: null,
  requiresAuth: true,
  keepAlive: false,
  hide: false,
  withoutTab: true,
  pinTab: false,
  menuType: 'page',
  available_permissions: [],
}
const formModel = ref<Api.Menu.MenuItem>({ ...formDefault })

// 权限配置相关
const basicPermissions = [
  { action: 'read', label: '查看', description: '查看数据' },
  { action: 'create', label: '创建', description: '创建数据' },
  { action: 'update', label: '修改', description: '修改数据' },
  { action: 'delete', label: '删除', description: '删除数据' },
]

// 预定义业务动作列表（与后端 casbin_auth.go 保持一致）
const predefinedActions = [
  // 财务相关
  { value: 'pending', label: '挂账', category: '财务相关' },
  { value: 'verify', label: '核销', category: '财务相关' },
  { value: 'settle', label: '结算', category: '财务相关' },
  { value: 'reconcile', label: '对账', category: '财务相关' },
  { value: 'refund', label: '退款', category: '财务相关' },
  
  // 审批相关
  { value: 'approve', label: '批准', category: '审批相关' },
  { value: 'reject', label: '拒绝', category: '审批相关' },
  { value: 'audit', label: '审核', category: '审批相关' },
  { value: 'submit', label: '提交', category: '审批相关' },
  { value: 'withdraw', label: '撤回', category: '审批相关' },
  
  // 数据操作
  { value: 'export', label: '导出', category: '数据操作' },
  { value: 'import', label: '导入', category: '数据操作' },
  { value: 'sync', label: '同步', category: '数据操作' },
  { value: 'refresh', label: '刷新', category: '数据操作' },
  { value: 'calculate', label: '计算', category: '数据操作' },
  { value: 'generate', label: '生成', category: '数据操作' },
  
  // 状态变更
  { value: 'publish', label: '发布', category: '状态变更' },
  { value: 'cancel', label: '取消', category: '状态变更' },
  { value: 'close', label: '关闭', category: '状态变更' },
  { value: 'reopen', label: '重开', category: '状态变更' },
  { value: 'archive', label: '归档', category: '状态变更' },
  { value: 'restore', label: '恢复', category: '状态变更' },
  
  // 权限管理
  { value: 'assign', label: '分配', category: '权限管理' },
  { value: 'transfer', label: '转移', category: '权限管理' },
  { value: 'lock', label: '锁定', category: '权限管理' },
  { value: 'unlock', label: '解锁', category: '权限管理' },
  { value: 'enable', label: '启用', category: '权限管理' },
  { value: 'disable', label: '禁用', category: '权限管理' },
  
  // 其他
  { value: 'copy', label: '复制', category: '其他' },
  { value: 'move', label: '移动', category: '其他' },
  { value: 'merge', label: '合并', category: '其他' },
  { value: 'split', label: '拆分', category: '其他' },
  { value: 'convert', label: '转换', category: '其他' },
  { value: 'validate', label: '验证', category: '其他' },
  { value: 'notify', label: '通知', category: '其他' },
  { value: 'remind', label: '提醒', category: '其他' },
  { value: 'share', label: '分享', category: '其他' },
  { value: 'favorite', label: '收藏', category: '其他' },
  { value: 'star', label: '标星', category: '其他' },
  { value: 'pin', label: '置顶', category: '其他' },
  { value: 'unpin', label: '取消置顶', category: '其他' },
  { value: 'reset', label: '重置', category: '其他' },
  { value: 'retry', label: '重试', category: '其他' },
  { value: 'rollback', label: '回滚', category: '其他' },
  { value: 'upgrade', label: '升级', category: '其他' },
  { value: 'downgrade', label: '降级', category: '其他' },
]

// 根据 action 查找默认的 label
function getDefaultLabel(action: string): string {
  const found = predefinedActions.find(item => item.value === action)
  return found ? found.label : ''
}

// 扫描 views 目录下的所有 .vue 组件（排除 components 目录）
const viewComponents = computed(() => {
  const modules = import.meta.glob('/src/views/**/*.vue')
  const componentPaths: Array<{ label: string, value: string }> = []
  
  Object.keys(modules).forEach(path => {
    // 过滤掉 components 目录
    if (path.includes('/components/')) {
      return
    }
    
    // 转换路径格式：/src/views/system/admin/index.vue → /system/admin/index.vue
    const relativePath = path.replace('/src/views', '')
    
    // 生成显示标签（更友好的显示）
    const displayPath = relativePath
      .replace(/^\//, '')
      .replace(/\.vue$/, '')
    
    componentPaths.push({
      label: displayPath,
      value: relativePath,
    })
  })
  
  // 按路径排序
  return componentPaths.sort((a, b) => a.value.localeCompare(b.value))
})

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

async function openModal(type: ModalType = 'add', data?: Api.Menu.MenuItem) {
  emit('open')
  modalType.value = type
  showModal()
  const handlers = {
    async add() {
      formModel.value = { ...formDefault }
    },
    async view() {
      if (!data)
        return
      formModel.value = { ...data }
    },
    async edit() {
      if (!data)
        return
      formModel.value = { ...data }
    },
  }
  await handlers[type]()
}

function closeModal() {
  hiddenModal()
  endLoading()
  emit('close')
}

// 添加业务权限
function addBusinessPermission() {
  if (!formModel.value.available_permissions) {
    formModel.value.available_permissions = []
  }
  formModel.value.available_permissions.push({
    action: '',
    label: '',
    description: '',
    is_basic: false,
  })
}

// 当 action 改变时，自动填充默认的 label
function onActionChange(perm: any, newAction: string) {
  perm.action = newAction
  // 如果 label 为空，自动填充预定义的 label
  if (!perm.label) {
    perm.label = getDefaultLabel(newAction)
  }
}

// 删除权限
function removePermission(index: number) {
  formModel.value.available_permissions?.splice(index, 1)
}

// 检查基础权限是否被选中
function isBasicPermissionChecked(action: string) {
  return formModel.value.available_permissions?.some(p => p.is_basic && p.action === action) || false
}

// 切换基础权限
function toggleBasicPermission(bp: any, checked: boolean) {
  if (!formModel.value.available_permissions) {
    formModel.value.available_permissions = []
  }
  
  const index = formModel.value.available_permissions.findIndex(p => p.is_basic && p.action === bp.action)
  
  if (checked && index === -1) {
    // 添加
    formModel.value.available_permissions.push({
      action: bp.action,
      label: bp.label,
      description: bp.description,
      is_basic: true,
    })
  }
  else if (!checked && index !== -1) {
    // 删除
    formModel.value.available_permissions.splice(index, 1)
  }
}

defineExpose({
  openModal,
})

const formRef = ref()
async function submitModal() {
  const handlers = {
    async add() {
      try {
        const { isSuccess } = await createMenu({
          name: formModel.value.name,
          path: formModel.value.path,
          title: formModel.value.title,
          pid: formModel.value.pid,
          componentPath: formModel.value.componentPath,
          icon: formModel.value.icon,
          requiresAuth: formModel.value.requiresAuth,
          roles: formModel.value.roles,
          keepAlive: formModel.value.keepAlive,
          hide: formModel.value.hide,
          order: formModel.value.order,
          href: formModel.value.href,
          activeMenu: formModel.value.activeMenu,
          withoutTab: formModel.value.withoutTab,
          pinTab: formModel.value.pinTab,
          menuType: formModel.value.menuType,
          available_permissions: formModel.value.available_permissions,
        })
        if (isSuccess) {
          window.$message.success('创建成功')
          emit('refresh')
          return true
        }
      }
      catch (e) {
        console.error('[Create Menu Error]:', e)
        return false
      }
    },
    async edit() {
      try {
        const { isSuccess } = await updateMenu(formModel.value.id, {
          name: formModel.value.name,
          path: formModel.value.path,
          title: formModel.value.title,
          pid: formModel.value.pid,
          componentPath: formModel.value.componentPath,
          icon: formModel.value.icon,
          requiresAuth: formModel.value.requiresAuth,
          roles: formModel.value.roles,
          keepAlive: formModel.value.keepAlive,
          hide: formModel.value.hide,
          order: formModel.value.order,
          href: formModel.value.href,
          activeMenu: formModel.value.activeMenu,
          withoutTab: formModel.value.withoutTab,
          pinTab: formModel.value.pinTab,
          menuType: formModel.value.menuType,
          available_permissions: formModel.value.available_permissions,
        })
        if (isSuccess) {
          window.$message.success('更新成功')
          emit('refresh')
          return true
        }
      }
      catch (e) {
        console.error('[Update Menu Error]:', e)
        return false
      }
    },
    async view() {
      return true
    },
  }
  await formRef.value?.validate()
  startLoading()
  await handlers[modalType.value]() && closeModal()
  endLoading()
}

const dirTreeOptions = computed(() => {
  return filterDirectory(JSON.parse(JSON.stringify(allRoutes)))
})

function filterDirectory(node: any[]) {
  return node.filter((item) => {
    if (item.children) {
      const childDir = filterDirectory(item.children)
      if (childDir.length > 0)
        item.children = childDir
      else
        Reflect.deleteProperty(item, 'children')
    }

    return (item.menuType === 'dir')
  })
}

const rules = {
  name: {
    required: true,
    // message: '请输入菜单名称',
    validator(rule: FormItemRule, value: string) {
      if (!value)
        return new Error('请输入菜单名称')

      if (!new RegExp(Regex.RouteName).test(value))
        return new Error('菜单只能包含英文数字_!@#$%^&*~-')

      return true
    },
    trigger: 'blur',
  },
  path: {
    required: true,
    message: '请输入菜单路径',
    trigger: 'blur',
  },
  componentPath: {
    required: true,
    message: '请输入组件路径',
    trigger: 'blur',
  },
  title: {
    required: true,
    message: '请输入菜单标题',
    trigger: 'blur',
  },
}

// const options = ref()
// async function getRoleList() {
//   const { data } = await fetchRoleList()
//   options.value = data
// }
</script>

<template>
  <n-drawer
    v-model:show="modalVisible"
    :width="800"
    placement="right"
    :trap-focus="false"
    :block-scroll="true"
  >
    <n-drawer-content :title="modalTitle" closable>
      <n-form
        ref="formRef"
        :rules="rules"
        label-placement="left"
        :label-width="100"
        :model="formModel"
        :disabled="modalType === 'view'"
      >
      <n-grid :cols="2" :x-gap="18">
        <n-form-item-grid-item :span="2" path="pid">
          <template #label>
            父级目录
            <HelpInfo message="不填写则为顶层菜单" />
          </template>
          <n-tree-select
            v-model:value="formModel.pid" filterable clearable :options="dirTreeOptions" key-field="id"
            label-field="title" children-field="children" placeholder="请选择父级目录"
          />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="菜单名称" path="name">
          <n-input v-model:value="formModel.name" placeholder="Eg: system" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="标题" path="title">
          <n-input v-model:value="formModel.title" placeholder="Eg: My-System" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="2" label="路由路径" path="path">
          <n-input v-model:value="formModel.path" placeholder="Eg: /system/user" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="菜单类型" path="menuType">
          <n-radio-group v-model:value="formModel.menuType" name="radiogroup">
            <n-space>
              <n-radio value="dir">
                目录
              </n-radio>
              <n-radio value="page">
                页面
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="图标" path="icon">
          <icon-select v-model:value="formModel.icon" :disabled="modalType === 'view'" />
        </n-form-item-grid-item>
        <n-form-item-grid-item v-if="formModel.menuType === 'page'" :span="2" path="componentPath">
          <template #label>
            组件路径
            <HelpInfo message="从列表选择或手动输入组件路径" />
          </template>
          <n-select
            v-model:value="formModel.componentPath"
            :options="viewComponents"
            filterable
            tag
            placeholder="选择或输入组件路径"
            :render-label="(option: any) => option.label"
          >

          </n-select>
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" path="order">
          <template #label>
            菜单排序
            <HelpInfo message="数字越小，同级中越靠前" />
          </template>
          <n-input-number v-model:value="formModel.order" />
        </n-form-item-grid-item>
        <n-form-item-grid-item v-if="formModel.menuType === 'page'" :span="1" path="href">
          <template #label>
            外链页面
            <HelpInfo message="填写后，点击菜单将跳转到该地址，组件路径任意填写" />
          </template>
          <n-input v-model:value="formModel.href" placeholder="Eg: https://example.com" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="登录访问" path="requiresAuth">
          <n-switch v-model:value="formModel.requiresAuth" />
        </n-form-item-grid-item>
        <n-form-item-grid-item
          v-if="formModel.menuType === 'page'" :span="1" label="页面缓存"
          path="keepAlive"
        >
          <n-switch v-model:value="formModel.keepAlive" />
        </n-form-item-grid-item>
        <n-form-item-grid-item
          v-if="formModel.menuType === 'page'" :span="1" label="标签栏可见"
          path="withoutTab"
        >
          <n-switch v-model:value="formModel.withoutTab" />
        </n-form-item-grid-item>
        <n-form-item-grid-item v-if="formModel.menuType === 'page'" :span="1" label="常驻标签栏" path="pinTab">
          <n-switch v-model:value="formModel.pinTab" />
        </n-form-item-grid-item>
        <n-form-item-grid-item :span="1" label="侧边菜单隐藏" path="hide">
          <n-switch v-model:value="formModel.hide" />
        </n-form-item-grid-item>
        <n-form-item-grid-item
          v-if="formModel.menuType === 'page' && formModel.hide" :span="2"
          path="activeMenu"
        >
          <template #label>
            高亮菜单
            <HelpInfo message="当前路由不在左侧菜单显示，但需要高亮某个菜单" />
          </template>
          <n-input v-model:value="formModel.activeMenu" />
        </n-form-item-grid-item>
        
        <!-- 权限配置 -->
        <n-form-item-grid-item v-if="formModel.menuType === 'page'" :span="2" path="available_permissions">
          <template #label>
            权限配置
            <HelpInfo message="配置该菜单支持的操作权限（CRUD + 业务权限）" />
          </template>
          <n-card size="small">
            <div class="space-y-4">
              <!-- 基础权限 -->
              <div>
                <div class="mb-2 font-medium text-sm">
                  基础权限（CRUD）
                </div>
                <n-space>
                  <n-checkbox
                    v-for="bp in basicPermissions"
                    :key="bp.action"
                    :checked="isBasicPermissionChecked(bp.action)"
                    @update:checked="(checked) => toggleBasicPermission(bp, checked)"
                  >
                    {{ bp.label }}
                  </n-checkbox>
                </n-space>
              </div>
              
              <!-- 业务权限 -->
              <div>
                <div class="mb-2 font-medium text-sm flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <span>业务权限（如：挂账、核销等）</span>
                    <n-tooltip trigger="hover" placement="right">
                      <template #trigger>
                        <span class="text-blue-500 cursor-help">ℹ️</span>
                      </template>
                      <div class="text-sm max-w-xs">
                        <div class="font-medium mb-1">路由命名规范：</div>
                        <div class="text-gray-300 mb-2">路径最后一段必须是 action 名，支持任何 HTTP 方法</div>
                        <div class="space-y-1">
                          <div>• GET /finance/:id/pending（查询挂账列表）</div>
                          <div>• POST /finance/:id/pending（执行挂账）</div>
                          <div>• PUT /finance/:id/verify（更新核销状态）</div>
                          <div>• DELETE /finance/:id/cancel（取消操作）</div>
                        </div>
                      </div>
                    </n-tooltip>
                  </div>
                  <n-button size="tiny" @click="addBusinessPermission">
                    + 添加业务权限
                  </n-button>
                </div>
                <div
                  v-for="(perm, index) in formModel.available_permissions?.filter(p => !p.is_basic)"
                  :key="index"
                  class="mb-2 p-2 border rounded"
                >
                  <n-space align="center">
                    <n-tooltip trigger="hover">
                      <template #trigger>
                        <n-select
                          :value="perm.action"
                          :options="predefinedActions"
                          filterable
                          tag
                          placeholder="选择或输入权限动作"
                          style="width: 180px"
                          size="small"
                          @update:value="(val) => onActionChange(perm, val)"
                        >

                        </n-select>
                      </template>
                      <div class="text-sm">
                        <div class="font-medium mb-1">路由命名规范：</div>
                        <div class="text-gray-400 mb-1">路径最后一段必须是此 action 名</div>
                        <div>• GET /finance/pending（查询）</div>
                        <div>• POST /finance/:id/pending（创建/执行）</div>
                        <div>• PUT /workflow/:id/approve（更新）</div>
                        <div>• DELETE /order/:id/cancel（删除）</div>
                        <div class="mt-1 text-warning">⚠️ 支持任何 HTTP 方法，不限于 POST</div>
                      </div>
                    </n-tooltip>
                    <n-input
                      v-model:value="perm.label"
                      placeholder="显示名称 (自动填充)"
                      style="width: 120px"
                      size="small"
                    />
                    <n-input
                      v-model:value="perm.description"
                      placeholder="描述说明 (可选)"
                      style="width: 200px"
                      size="small"
                    />
                    <n-button
                      size="tiny"
                      type="error"
                      @click="removePermission(index + (formModel.available_permissions?.filter(p => p.is_basic).length || 0))"
                    >
                      删除
                    </n-button>
                  </n-space>
                </div>
                <div v-if="!formModel.available_permissions?.some(p => !p.is_basic)" class="text-gray-400 text-sm">
                  暂无业务权限，点击上方按钮添加
                </div>
              </div>
            </div>
          </n-card>
        </n-form-item-grid-item>
        
        <!-- <n-form-item-grid-item :span="2" path="roles">
          <template #label>
            访问角色
            <HelpInfo message="不填写则表示所有角色都可以访问" />
          </template>
          <n-select
            v-model:value="formModel.roles" multiple filterable
            label-field="role"
            value-field="id"
            :options="options"
          />
        </n-form-item-grid-item> -->
      </n-grid>
      </n-form>
      <template v-if="modalType !== 'view'" #footer>
        <n-space justify="end">
          <n-button @click="closeModal">
            取消
          </n-button>
          <n-button type="primary" :loading="submitLoading" @click="submitModal">
            提交
          </n-button>
        </n-space>
      </template>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped></style>
