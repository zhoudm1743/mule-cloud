<script setup lang="ts">
import { ref } from 'vue'
import { NDrawer, NDrawerContent, NSpace, NCard, NDescriptions, NDescriptionsItem, NTag, NSpin, NCode } from 'naive-ui'
import { fetchOperationLogById } from '@/service'

const show = ref(false)
const loading = ref(false)
const logData = ref<Api.OperationLog.OperationLogInfo | null>(null)

// 打开抽屉
async function open(id: string) {
  show.value = true
  logData.value = null
  await fetchDetail(id)
}

// 获取详情数据
async function fetchDetail(id: string) {
  loading.value = true
  try {
    const res = await fetchOperationLogById(id)
    if (res.data?.log) {
      logData.value = res.data.log
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取操作日志详情失败')
    show.value = false
  }
  finally {
    loading.value = false
  }
}

// 获取HTTP方法的标签类型
function getMethodType(method: string): NaiveUI.ThemeColor {
  const methodMap: Record<string, NaiveUI.ThemeColor> = {
    'GET': 'info',
    'POST': 'success',
    'PUT': 'warning',
    'DELETE': 'error',
    'PATCH': 'warning',
  }
  return methodMap[method] || 'default'
}

// 获取状态码的标签类型
function getStatusType(code: number): NaiveUI.ThemeColor {
  if (code >= 200 && code < 300) return 'success'
  if (code >= 300 && code < 400) return 'info'
  if (code >= 400 && code < 500) return 'warning'
  if (code >= 500) return 'error'
  return 'default'
}

// 格式化日期时间
function formatDateTime(dateStr: string): string {
  try {
    const date = new Date(dateStr)
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch {
    return dateStr
  }
}

// 格式化耗时
function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}

// 格式化JSON
function formatJson(json: string): string {
  try {
    const obj = JSON.parse(json)
    return JSON.stringify(obj, null, 2)
  }
  catch {
    return json
  }
}

defineExpose({
  open,
})
</script>

<template>
  <NDrawer v-model:show="show" :width="800" placement="right">
    <NDrawerContent title="操作日志详情" closable>
      <NSpin :show="loading">
        <NSpace v-if="logData" vertical size="large">
          <!-- 基本信息 -->
          <NCard size="small" :bordered="false">
            <template #header>
              <span class="font-medium">基本信息</span>
            </template>
            <NDescriptions label-placement="left" :column="1">
              <NDescriptionsItem label="日志ID">
                {{ logData.id }}
              </NDescriptionsItem>
              <NDescriptionsItem label="操作用户">
                <div class="flex items-center gap-2">
                  <span>{{ logData.username }}</span>
                  <NTag size="small" type="info">ID: {{ logData.user_id }}</NTag>
                </div>
              </NDescriptionsItem>
              <NDescriptionsItem label="操作时间">
                {{ formatDateTime(logData.created_at) }}
              </NDescriptionsItem>
            </NDescriptions>
          </NCard>

          <!-- 请求信息 -->
          <NCard size="small" :bordered="false">
            <template #header>
              <span class="font-medium">请求信息</span>
            </template>
            <NDescriptions label-placement="left" :column="1">
              <NDescriptionsItem label="HTTP方法">
                <NTag :type="getMethodType(logData.method)">
                  {{ logData.method }}
                </NTag>
              </NDescriptionsItem>
              <NDescriptionsItem label="请求路径">
                <NCode :code="logData.path" language="text" />
              </NDescriptionsItem>
              <NDescriptionsItem label="资源名称">
                {{ logData.resource }}
              </NDescriptionsItem>
              <NDescriptionsItem label="操作类型">
                <NTag type="primary">
                  {{ logData.action }}
                </NTag>
              </NDescriptionsItem>
              <NDescriptionsItem label="客户端IP">
                {{ logData.ip }}
              </NDescriptionsItem>
              <NDescriptionsItem label="User Agent">
                <div class="break-all text-sm">{{ logData.user_agent }}</div>
              </NDescriptionsItem>
            </NDescriptions>
          </NCard>

          <!-- 响应信息 -->
          <NCard size="small" :bordered="false">
            <template #header>
              <span class="font-medium">响应信息</span>
            </template>
            <NDescriptions label-placement="left" :column="1">
              <NDescriptionsItem label="响应状态码">
                <NTag :type="getStatusType(logData.response_code)">
                  {{ logData.response_code }}
                </NTag>
              </NDescriptionsItem>
              <NDescriptionsItem label="请求耗时">
                <NTag type="info">
                  {{ formatDuration(logData.duration) }}
                </NTag>
              </NDescriptionsItem>
              <NDescriptionsItem v-if="logData.error" label="错误信息">
                <div class="text-red-500">{{ logData.error }}</div>
              </NDescriptionsItem>
            </NDescriptions>
          </NCard>

          <!-- 请求体 -->
          <NCard v-if="logData.request_body" size="small" :bordered="false">
            <template #header>
              <span class="font-medium">请求体</span>
            </template>
            <NCode 
              :code="formatJson(logData.request_body)" 
              language="json" 
              :show-line-numbers="true"
              class="max-h-400px overflow-auto"
            />
          </NCard>
        </NSpace>
      </NSpin>
    </NDrawerContent>
  </NDrawer>
</template>

<style scoped>
.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  font-size: 12px;
  color: var(--text-color-3);
}

.value {
  font-size: 14px;
  color: var(--text-color-1);
  font-weight: 500;
}
</style>

