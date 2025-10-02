<script setup lang="ts">
interface Props {
  node: any
  level: number
  selectedMenus: string[]
  permissions: Record<string, Record<string, boolean>>
  allMenus: Api.Menu.MenuItem[]
  expandedKeys: Set<string>
}

interface Emits {
  (e: 'toggleExpand', key: string): void
  (e: 'toggleMenu', menuName: string, node: any): void
  (e: 'updatePermission', menuName: string, action: string, enabled: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 判断是否有子节点
function hasChildren(): boolean {
  return props.node.children && props.node.children.length > 0
}

// 是否展开
function isExpanded(): boolean {
  return props.expandedKeys.has(props.node.key)
}

// 是否选中
function isChecked(): boolean {
  return props.selectedMenus.includes(props.node.key)
}

// 获取菜单权限配置
function getMenuPermissions() {
  const menu = props.allMenus.find(m => m.name === props.node.key)
  if (!menu) return { basic: [], business: [] }
  
  const basic: any[] = []
  const business: any[] = []
  
  // 只有明确配置了 available_permissions 的菜单才返回权限
  // 没有配置的菜单不显示权限配置区
  if (menu.available_permissions && menu.available_permissions.length > 0) {
    menu.available_permissions.forEach(perm => {
      if (perm.is_basic) {
        basic.push(perm)
      } else {
        business.push(perm)
      }
    })
  }
  
  return { basic, business }
}

// 判断是否应该显示权限配置区
function shouldShowPermissions(): boolean {
  const menu = props.allMenus.find(m => m.name === props.node.key)
  // 只有配置了 available_permissions 的菜单才显示权限配置
  return !!(menu?.available_permissions && menu.available_permissions.length > 0)
}
</script>

<template>
  <div class="menu-tree-node">
    <!-- 节点本身 -->
    <div 
      class="flex items-start py-2 px-2 rounded transition-colors cursor-pointer"
      :class="isChecked() ? 'bg-blue-50' : ''"
      :style="{ paddingLeft: `${level * 20}px` }"
      @mouseenter="$event.currentTarget.classList.add('bg-gray-100')"
      @mouseleave="$event.currentTarget.classList.remove('bg-gray-100')"
    >
      <!-- 展开/折叠图标 -->
      <div 
        class="w-20px h-20px flex items-center justify-center flex-shrink-0"
        @click.stop="hasChildren() && emit('toggleExpand', node.key)"
      >
        <svg 
          v-if="hasChildren()"
          class="text-gray-500 transition-transform"
          :class="{ 'transform rotate-90': isExpanded() }"
          style="width: 16px; height: 16px;"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M8.59 16.59L13.17 12L8.59 7.41L10 6l6 6l-6 6l-1.41-1.41z" fill="currentColor"/>
        </svg>
        <span v-else class="text-gray-300 text-xs">•</span>
      </div>
      
      <!-- 菜单复选框 -->
      <div class="flex-shrink-0 mr-2" @click.stop>
        <NCheckbox 
          :checked="isChecked()" 
          size="small"
          @update:checked="() => emit('toggleMenu', node.key, node)"
        />
      </div>
      
      <!-- 菜单标题和权限 -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2 mb-1">
          <span 
            class="text-sm"
            :class="isChecked() ? 'font-medium text-blue-700' : 'text-gray-700'"
          >
            {{ node.label }}
          </span>
          <span v-if="!hasChildren()" class="text-xs text-gray-400">({{ node.key }})</span>
        </div>
        
        <!-- 权限配置区 - 只在选中、没有子节点、且配置了权限时显示 -->
        <div 
          v-if="isChecked() && !hasChildren() && shouldShowPermissions()" 
          class="mt-2 space-y-2"
          @click.stop
        >
          <!-- 基础权限 -->
          <div 
            v-if="getMenuPermissions().basic.length > 0" 
            class="flex items-center gap-2 flex-wrap"
          >
            <NTag size="small" type="info" :bordered="false">基础</NTag>
            <NCheckbox
              v-for="perm in getMenuPermissions().basic"
              :key="perm.action"
              :checked="permissions[node.key]?.[perm.action] ?? true"
              size="small"
              @update:checked="(val: boolean) => emit('updatePermission', node.key, perm.action, val)"
            >
              {{ perm.label }}
            </NCheckbox>
          </div>
          
          <!-- 业务权限 -->
          <div 
            v-if="getMenuPermissions().business.length > 0" 
            class="flex items-center gap-2 flex-wrap"
          >
            <NTag size="small" type="warning" :bordered="false">业务</NTag>
            <template v-for="perm in getMenuPermissions().business" :key="perm.action">
              <NTooltip v-if="perm.description" trigger="hover">
                <template #trigger>
                  <NCheckbox
                    :checked="permissions[node.key]?.[perm.action] ?? true"
                    size="small"
                    @update:checked="(val: boolean) => emit('updatePermission', node.key, perm.action, val)"
                  >
                    {{ perm.label }}
                  </NCheckbox>
                </template>
                {{ perm.description }}
              </NTooltip>
              <NCheckbox
                v-else
                :checked="permissions[node.key]?.[perm.action] ?? true"
                size="small"
                @update:checked="(val: boolean) => emit('updatePermission', node.key, perm.action, val)"
              >
                {{ perm.label }}
              </NCheckbox>
            </template>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 子节点（递归） -->
    <div v-if="hasChildren() && isExpanded()" class="menu-tree-children">
      <MenuTreeNode 
        v-for="child in node.children" 
        :key="child.key" 
        :node="child" 
        :level="level + 1"
        :selected-menus="selectedMenus"
        :permissions="permissions"
        :all-menus="allMenus"
        :expanded-keys="expandedKeys"
        @toggle-expand="(key) => emit('toggleExpand', key)"
        @toggle-menu="(menuName, node) => emit('toggleMenu', menuName, node)"
        @update-permission="(menuName, action, enabled) => emit('updatePermission', menuName, action, enabled)"
      />
    </div>
  </div>
</template>

<style scoped>
.menu-tree-node {
  user-select: none;
}

.menu-tree-children {
  border-left: 1px solid #e5e7eb;
  margin-left: 0.5rem;
}
</style>

