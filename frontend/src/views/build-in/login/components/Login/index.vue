<script setup lang="ts">
import type { FormInst, SelectOption } from 'naive-ui'
import { useAuthStore } from '@/store'
import { local } from '@/utils'
import { fetchLoginTenantList } from '@/service'

const emit = defineEmits(['update:modelValue'])

const authStore = useAuthStore()

function toOtherForm(type: any) {
  emit('update:modelValue', type)
}

const { t } = useI18n()
const rules = computed(() => {
  return {
    phone: {
      required: true,
      trigger: 'blur',
      message: t('login.accountRuleTip'), // 手机号验证
    },
    pwd: {
      required: true,
      trigger: 'blur',
      message: t('login.passwordRuleTip'),
    },
  }
})
const formValue = ref({
  phone: '13800138000', // 默认手机号
  pwd: '123456',
  tenantCode: '', // 租户代码（可选，空字符串表示系统管理员）
})
const isRemember = ref(false)
const isLoading = ref(false)

// 租户列表
const tenantOptions = ref<SelectOption[]>([
  {
    label: '系统管理员',
    value: '',
  },
])
const isLoadingTenants = ref(false)

// 获取租户列表
async function loadTenantList() {
  try {
    isLoadingTenants.value = true
    const { data, error } = await fetchLoginTenantList()
    if (!error && data) {
      const tenants = data.tenants || []
      tenantOptions.value = [
        {
          label: '系统管理员',
          value: '',
        },
        ...tenants.map(tenant => ({
          label: `${tenant.name} (${tenant.code})`,
          value: tenant.code,
        })),
      ]
    }
  }
  catch (err) {
    console.error('加载租户列表失败:', err)
  }
  finally {
    isLoadingTenants.value = false
  }
}

const formRef = ref<FormInst | null>(null)
function handleLogin() {
  formRef.value?.validate(async (errors) => {
    if (errors)
      return

    isLoading.value = true
    const { phone, pwd, tenantCode } = formValue.value

    if (isRemember.value)
      local.set('loginAccount', { phone, pwd, tenantCode })
    else local.remove('loginAccount')

    await authStore.login(phone, pwd, tenantCode)
    isLoading.value = false
  })
}
onMounted(() => {
  checkUserAccount()
  loadTenantList()
})
function checkUserAccount() {
  const loginAccount = local.get('loginAccount')
  if (!loginAccount)
    return

  formValue.value = loginAccount
  isRemember.value = true
}
</script>

<template>
  <div>
    <n-h2 depth="3" class="text-center">
      {{ $t('login.signInTitle') }}
    </n-h2>
    <n-form ref="formRef" :rules="rules" :model="formValue" :show-label="false" size="large">
      <n-form-item path="tenantCode">
        <n-select
          v-model:value="formValue.tenantCode"
          :options="tenantOptions"
          :loading="isLoadingTenants"
          placeholder="选择租户"
          clearable
        >
          <template #prefix>
            <n-icon><icon-park-outline-building-four /></n-icon>
          </template>
        </n-select>
      </n-form-item>
      <n-form-item path="phone">
        <n-input v-model:value="formValue.phone" clearable placeholder="请输入手机号">
          <template #prefix>
            <n-icon><icon-park-outline-user /></n-icon>
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="pwd">
        <n-input v-model:value="formValue.pwd" type="password" :placeholder="$t('login.passwordPlaceholder')" clearable show-password-on="click">
          <template #prefix>
            <n-icon><icon-park-outline-lock /></n-icon>
          </template>
          <template #password-invisible-icon>
            <icon-park-outline-preview-close-one />
          </template>
          <template #password-visible-icon>
            <icon-park-outline-preview-open />
          </template>
        </n-input>
      </n-form-item>
      <n-space vertical :size="20">
        <div class="flex-y-center justify-between">
          <n-checkbox v-model:checked="isRemember">
            {{ $t('login.rememberMe') }}
          </n-checkbox>
          <n-button type="primary" text @click="toOtherForm('resetPwd')">
            {{ $t('login.forgotPassword') }}
          </n-button>
        </div>
        <n-button block type="primary" size="large" :loading="isLoading" :disabled="isLoading" @click="handleLogin">
          {{ $t('login.signIn') }}
        </n-button>
        <n-flex>
          <n-text>{{ $t('login.noAccountText') }}</n-text>
          <n-button type="primary" text @click="toOtherForm('register')">
            {{ $t('login.signUp') }}
          </n-button>
        </n-flex>
      </n-space>
    </n-form>
    <n-divider>
      <span op-80>{{ $t('login.or') }}</span>
    </n-divider>
    <n-space justify="center">
      <n-button circle>
        <template #icon>
          <n-icon><icon-park-outline-wechat /></n-icon>
        </template>
      </n-button>
      <n-button circle>
        <template #icon>
          <n-icon><icon-park-outline-tencent-qq /></n-icon>
        </template>
      </n-button>
      <n-button circle>
        <template #icon>
          <n-icon><icon-park-outline-github-one /></n-icon>
        </template>
      </n-button>
    </n-space>
  </div>
</template>

<style scoped></style>
