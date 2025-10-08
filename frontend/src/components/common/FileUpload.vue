<script setup lang="ts">
import { ref } from 'vue'
import { NUpload, NButton, NSpace, NImage, NText, NIcon, type UploadFileInfo } from 'naive-ui'
import { uploadFile, deleteFile } from '@/service/api/common'
import { $t } from '@/utils'

interface Props {
  businessType: string // 业务类型
  accept?: string // 接受的文件类型
  maxSize?: number // 最大文件大小（MB）
  multiple?: boolean // 是否支持多选
  listType?: 'text' | 'image' | 'image-card' // 列表类型
  value?: string | string[] // 已上传文件URL
  disabled?: boolean
}

interface Emits {
  (e: 'update:value', value: string | string[]): void
  (e: 'success', fileInfo: Api.Common.UploadResponse): void
  (e: 'error', error: Error): void
}

const props = withDefaults(defineProps<Props>(), {
  accept: '*',
  maxSize: 100,
  multiple: false,
  listType: 'text',
  disabled: false
})

const emit = defineEmits<Emits>()

const fileList = ref<UploadFileInfo[]>([])
const loading = ref(false)

// 自定义上传
async function customRequest({ file, onFinish, onError }: any) {
  loading.value = true

  try {
    // 验证文件大小
    const fileSizeMB = file.file.size / 1024 / 1024
    if (fileSizeMB > props.maxSize) {
      const error = new Error(`文件大小不能超过${props.maxSize}MB`)
      onError()
      emit('error', error)
      window.$message?.error(error.message)
      return
    }

    // 调用上传API
    const response = await uploadFile(file.file, props.businessType)

    // 上传成功
    file.url = response.data.url
    file.id = response.data.id
    onFinish()

    emit('success', response.data)
    emit('update:value', props.multiple ? getUrls() : response.data.url)
    window.$message?.success('上传成功')
  }
  catch (error: any) {
    onError()
    emit('error', error)
    window.$message?.error('上传失败: ' + error.message)
  }
  finally {
    loading.value = false
  }
}

// 移除文件
async function handleRemove({ file }: { file: UploadFileInfo }) {
  if (file.id) {
    try {
      await deleteFile(file.id as string)
      window.$message?.success('删除成功')
      emit('update:value', props.multiple ? getUrls() : '')
    }
    catch (error: any) {
      window.$message?.error('删除失败: ' + error.message)
    }
  }
}

// 获取所有文件URL
function getUrls(): string[] {
  return fileList.value.filter(f => f.url).map(f => f.url as string)
}

// 预览文件
function handlePreview(file: UploadFileInfo) {
  if (file.url) {
    window.open(file.url, '_blank')
  }
}
</script>

<template>
  <NUpload
    v-model:file-list="fileList"
    :custom-request="customRequest"
    :accept="accept"
    :multiple="multiple"
    :list-type="listType"
    :disabled="disabled"
    @remove="handleRemove"
    @preview="handlePreview"
  >
    <NButton :loading="loading" :disabled="disabled">
      <template #icon>
        <nova-icon icon="carbon:upload" :size="18" />
      </template>
      点击上传
    </NButton>
  </NUpload>
</template>

<style scoped>
/* 根据需要添加样式 */
</style>
