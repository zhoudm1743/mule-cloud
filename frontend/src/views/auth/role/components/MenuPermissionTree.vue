<script setup lang="ts">
import { ref } from 'vue'
import MenuTreeNode from './MenuTreeNode.vue'

interface Props {
  selectedMenus: string[]
  permissions: Record<string, Record<string, boolean>>
  treeData: any[]
  allMenus: Api.Menu.MenuItem[]
}

interface Emits {
  (e: 'update:selectedMenus', value: string[]): void
  (e: 'update:permissions', value: Record<string, Record<string, boolean>>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 展开状态
const expandedKeys = ref<Set<string>>(new Set())

// 初始化：展开所有节点
function initExpanded(nodes: any[]) {
  nodes.forEach(node => {
    if (node.children && node.children.length > 0) {
      expandedKeys.value.add(node.key)
      initExpanded(node.children)
    }
  })
}
initExpanded(props.treeData)

// 切换展开/折叠
function toggleExpand(key: string) {
  if (expandedKeys.value.has(key)) {
    expandedKeys.value.delete(key)
  } else {
    expandedKeys.value.add(key)
  }
}

// 切换菜单选中状态
function toggleMenu(menuName: string, node: any) {
  const newSelectedMenus = [...props.selectedMenus]
  const index = newSelectedMenus.indexOf(menuName)
  
  if (index > -1) {
    // 取消选中：移除自己和所有子节点
    const toRemove = [menuName]
    collectChildrenKeys(node, toRemove)
    emit('update:selectedMenus', newSelectedMenus.filter(m => !toRemove.includes(m)))
    
    // 清理权限数据
    const newPermissions = { ...props.permissions }
    toRemove.forEach(key => delete newPermissions[key])
    emit('update:permissions', newPermissions)
  } else {
    // 选中：添加自己（不级联）
    newSelectedMenus.push(menuName)
    emit('update:selectedMenus', newSelectedMenus)
    
    // 初始化权限数据
    initMenuPermission(menuName)
  }
}

// 收集所有子节点的key
function collectChildrenKeys(node: any, result: string[]) {
  if (node.children) {
    node.children.forEach((child: any) => {
      result.push(child.key)
      collectChildrenKeys(child, result)
    })
  }
}

// 初始化菜单权限
function initMenuPermission(menuName: string) {
  const menu = props.allMenus.find(m => m.name === menuName)
  if (!menu) return
  
  // 只有配置了 available_permissions 的菜单才初始化权限
  if (!menu.available_permissions || menu.available_permissions.length === 0) {
    return
  }
  
  const newPermissions = { ...props.permissions }
  newPermissions[menuName] = {}
  
  menu.available_permissions.forEach(perm => {
    newPermissions[menuName][perm.action] = true
  })
  
  emit('update:permissions', newPermissions)
}

// 更新权限
function updatePermission(menuName: string, action: string, enabled: boolean) {
  const newPermissions = { ...props.permissions }
  if (!newPermissions[menuName]) {
    newPermissions[menuName] = {}
  }
  newPermissions[menuName][action] = enabled
  emit('update:permissions', newPermissions)
}
</script>

<template>
  <div class="menu-permission-tree">
    <div v-if="!treeData || treeData.length === 0" class="text-gray-400 text-center py-8">
      暂无菜单数据
    </div>
    <div v-else>
      <MenuTreeNode
        v-for="node in treeData"
        :key="node.key"
        :node="node"
        :level="0"
        :selected-menus="selectedMenus"
        :permissions="permissions"
        :all-menus="allMenus"
        :expanded-keys="expandedKeys"
        @toggle-expand="toggleExpand"
        @toggle-menu="toggleMenu"
        @update-permission="updatePermission"
      />
    </div>
  </div>
</template>

<style scoped>
.menu-permission-tree {
  user-select: none;
}
</style>
