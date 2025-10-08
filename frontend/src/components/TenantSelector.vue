<script setup lang="ts">
import { fetchAllTenants } from '@/service'
import { useAuthStore } from '@/store'
import { local } from '@/utils'
import { computed, onMounted, ref } from 'vue'
import type { SelectOption } from 'naive-ui'

const authStore = useAuthStore()

// 租户选项
const tenantOptions = ref<SelectOption[]>([
  {
    label: '所有租户（系统视图）',
    value: '',
  },
])

// 当前选择的租户（使用 Code 而不是 ID）
const selectedTenantId = ref<string>('')

// 是否为系统管理员
const isSystemAdmin = computed(() => {
  const userTenantId = authStore.userInfo?.tenant_id
  return !userTenantId || userTenantId === ''
})

// 加载租户列表
async function loadTenants() {
  try {
    const { data } = await fetchAllTenants()
    if (data && data.tenants) {
      tenantOptions.value = [
        {
          label: '所有租户（系统视图）',
          value: '',
        },
        ...data.tenants.map(tenant => ({
          label: `${tenant.name} (${tenant.code})`,
          value: tenant.code, // ✅ 使用 code 而不是 id
        })),
      ]
    }
  }
  catch (error) {
    console.error('加载租户列表失败:', error)
  }
}

// 切换租户
function onTenantChange(value: string) {
  // ✅ 保存选择的租户 Code（而不是 ID）
  local.set('selected_tenant_code', value)
  
  // 刷新当前页面
  window.location.reload()
}

// 恢复上次选择的租户
function restoreSelection() {
  const savedTenantCode = local.get('selected_tenant_code')
  if (savedTenantCode) {
    selectedTenantId.value = savedTenantCode
  }
}

onMounted(() => {
  if (isSystemAdmin.value) {
    loadTenants()
    restoreSelection()
  }
})

defineExpose({
  isSystemAdmin,
  selectedTenantId,
})
</script>

<template>
  <div v-if="isSystemAdmin" class="tenant-selector">
    <n-space align="center">
      <n-text>管理租户：</n-text>
      <n-select
        v-model:value="selectedTenantId"
        :options="tenantOptions"
        placeholder="选择租户"
        clearable
        style="width: 250px"
        @update:value="onTenantChange"
      >
        <template #prefix>
          <n-icon><icon-park-outline-building-four /></n-icon>
        </template>
      </n-select>
      <n-text v-if="selectedTenantId" type="warning" depth="3">
        当前正在管理该租户的数据
      </n-text>
    </n-space>
  </div>
</template>

<style scoped>
.tenant-selector {
  padding: 12px 16px;
  background-color: var(--n-color-target);
  border-bottom: 1px solid var(--n-border-color);
}
</style>
