<script setup lang="ts">
import { fetchRegister } from '@/service'

const emit = defineEmits(['update:modelValue'])
function toLogin() {
  emit('update:modelValue', 'login')
}
const { t } = useI18n()

const rules = {
  phone: {
    required: true,
    trigger: 'blur',
    message: '请输入手机号',
  },
  nickname: {
    required: true,
    trigger: 'blur',
    message: '请输入昵称',
  },
  pwd: {
    required: true,
    trigger: 'blur',
    message: t('login.passwordRuleTip'),
  },
  rePwd: {
    required: true,
    trigger: 'blur',
    message: t('login.checkPasswordRuleTip'),
  },
}
const formValue = ref({
  phone: '',
  nickname: '',
  pwd: '',
  rePwd: '',
  email: '',
})

const isRead = ref(false)
const isLoading = ref(false)

async function handleRegister() {
  if (!isRead.value) {
    window.$message.warning('请阅读并同意用户协议')
    return
  }

  if (formValue.value.pwd !== formValue.value.rePwd) {
    window.$message.error('两次密码不一致')
    return
  }

  try {
    isLoading.value = true
    const { isSuccess, data } = await fetchRegister({
      phone: formValue.value.phone,
      password: formValue.value.pwd,
      nickname: formValue.value.nickname,
      email: formValue.value.email,
    })

    if (isSuccess) {
      window.$message.success(data.message || '注册成功！')
      toLogin()
    }
  }
  catch (e) {
    console.error('[Register Error]:', e)
  }
  finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div>
    <n-h2 depth="3" class="text-center">
      {{ $t('login.registerTitle') }}
    </n-h2>
    <n-form
      :rules="rules"
      :model="formValue"
      :show-label="false"
      size="large"
    >
      <n-form-item path="phone">
        <n-input
          v-model:value="formValue.phone"
          clearable
          placeholder="请输入手机号"
        />
      </n-form-item>
      <n-form-item path="nickname">
        <n-input
          v-model:value="formValue.nickname"
          clearable
          placeholder="请输入昵称"
        />
      </n-form-item>
      <n-form-item>
        <n-input
          v-model:value="formValue.email"
          clearable
          placeholder="请输入邮箱（可选）"
        />
      </n-form-item>
      <n-form-item path="pwd">
        <n-input
          v-model:value="formValue.pwd"
          type="password"
          :placeholder="$t('login.passwordPlaceholder')"
          clearable
          show-password-on="click"
        >
          <template #password-invisible-icon>
            <icon-park-outline-preview-close-one />
          </template>
          <template #password-visible-icon>
            <icon-park-outline-preview-open />
          </template>
        </n-input>
      </n-form-item>
      <n-form-item path="rePwd">
        <n-input
          v-model:value="formValue.rePwd"
          type="password"
          :placeholder="$t('login.checkPasswordPlaceholder')"
          clearable
          show-password-on="click"
        >
          <template #password-invisible-icon>
            <icon-park-outline-preview-close-one />
          </template>
          <template #password-visible-icon>
            <icon-park-outline-preview-open />
          </template>
        </n-input>
      </n-form-item>
      <n-form-item>
        <n-space
          vertical
          :size="20"
          class="w-full"
        >
          <n-checkbox v-model:checked="isRead">
            {{ $t('login.readAndAgree') }} <n-button
              type="primary"
              text
            >
              {{ $t('login.userAgreement') }}
            </n-button>
          </n-checkbox>
          <n-button
            block
            type="primary"
            :loading="isLoading"
            :disabled="isLoading"
            @click="handleRegister"
          >
            {{ $t('login.signUp') }}
          </n-button>
          <n-flex justify="center">
            <n-text>{{ $t('login.haveAccountText') }}</n-text>
            <n-button
              text
              type="primary"
              @click="toLogin"
            >
              {{ $t('login.signIn') }}
            </n-button>
          </n-flex>
        </n-space>
      </n-form-item>
    </n-form>
  </div>
</template>

<style scoped></style>
