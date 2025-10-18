<template>
  <NSpace vertical size="large" class="p-4">
    <NCard title="工作流定义管理" :bordered="false" class="rounded-8px shadow-sm">
      <template #header-extra>
        <NSpace>
          <NButton type="primary" @click="handleCreate">
            <template #icon>
              <nova-icon icon="carbon:add" />
            </template>
            新建工作流
          </NButton>
          <NButton @click="handleTemplateLibrary">
            <template #icon>
              <nova-icon icon="carbon:template" />
            </template>
            模板库
          </NButton>
        </NSpace>
      </template>

      <!-- 工作流列表 -->
      <NSpace vertical size="large">
        <NDataTable
          :columns="columns"
          :data="workflowList"
          :loading="loading"
          :single-line="false"
        />
        <Pagination :count="total" @change="handlePageChange" />
      </NSpace>
    </NCard>

    <!-- 模板库模态框 -->
    <NModal
      v-model:show="showTemplateModal"
      preset="card"
      title="工作流模板库"
      style="width: 800px"
    >
      <NGrid :cols="2" :x-gap="16" :y-gap="16">
        <NGridItem v-for="template in templates" :key="template.id">
          <NCard
            :title="template.name"
            hoverable
            class="cursor-pointer"
            @click="handleApplyTemplate(template)"
          >
            <NSpace vertical>
              <NText depth="3">
                {{ template.description }}
              </NText>
              <NTag size="small">
                {{ template.category }}
              </NTag>
              <NText depth="3" class="text-sm">
                流程: {{ template.preview }}
              </NText>
            </NSpace>
          </NCard>
        </NGridItem>
      </NGrid>
    </NModal>
  </NSpace>
</template>

<script setup lang="tsx">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  NButton,
  NCard,
  NDataTable,
  NGrid,
  NGridItem,
  NModal,
  NPopconfirm,
  NSpace,
  NSwitch,
  NTag,
  NText,
} from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useBoolean } from '@/hooks'
import {
  activateWorkflowDefinition,
  deactivateWorkflowDefinition,
  deleteWorkflowDefinition,
  fetchWorkflowDefinitions,
  fetchWorkflowTemplates,
} from '@/service/api/workflow-designer'
import Pagination from '@/components/common/Pagination.vue'

defineOptions({ name: 'WorkflowDesignerList' })

const router = useRouter()
const { bool: loading, setTrue: startLoading, setFalse: endLoading } = useBoolean(false)

const workflowList = ref<Api.WorkflowDesigner.WorkflowDefinition[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showTemplateModal = ref(false)
const templates = ref<Api.WorkflowDesigner.WorkflowTemplate[]>([])

// 分页处理
function handlePageChange(p: number, ps: number) {
  page.value = p
  pageSize.value = ps
  fetchData()
}

// 表格列定义
const columns: DataTableColumns<Api.WorkflowDesigner.WorkflowDefinition> = [
  {
    title: '工作流名称',
    key: 'name',
    width: 200,
  },
  {
    title: '编码',
    key: 'code',
    width: 150,
  },
  {
    title: '描述',
    key: 'description',
    width: 250,
    ellipsis: {
      tooltip: true,
    },
  },
  {
    title: '版本',
    key: 'version',
    width: 80,
    render: row => <NTag size="small">v{row.version}</NTag>,
  },
  {
    title: '状态数',
    key: 'states',
    width: 80,
    render: row => row.states?.length || 0,
  },
  {
    title: '转换规则数',
    key: 'transitions',
    width: 100,
    render: row => row.transitions?.length || 0,
  },
  {
    title: '激活状态',
    key: 'is_active',
    width: 100,
    render: (row) => {
      return (
        <NSwitch
          value={row.is_active}
          onUpdateValue={(val: boolean) => handleToggleActive(row, val)}
        />
      )
    },
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 160,
    render: row => new Date(row.updated_at * 1000).toLocaleString('zh-CN'),
  },
  {
    title: '操作',
    key: 'actions',
    width: 250,
    fixed: 'right',
    render: (row) => {
      return (
        <NSpace>
          <NButton
            size="small"
            type="info"
            onClick={() => handleView(row)}
          >
            查看
          </NButton>
          <NButton
            size="small"
            type="primary"
            onClick={() => handleEdit(row)}
          >
            编辑
          </NButton>
          <NButton
            size="small"
            onClick={() => handleCopy(row)}
          >
            复制
          </NButton>
          <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
            {{
              default: () => '确定删除吗？',
              trigger: () => <NButton size="small" type="error">删除</NButton>,
            }}
          </NPopconfirm>
        </NSpace>
      )
    },
  },
]

onMounted(() => {
  fetchData()
  loadTemplates()
})

// 加载工作流列表
async function fetchData() {
  startLoading()
  try {
    const { data } = await fetchWorkflowDefinitions({
      page: page.value,
      page_size: pageSize.value,
    })
    if (data) {
      workflowList.value = data.workflows || []
      total.value = data.total || 0
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取工作流列表失败')
  }
  finally {
    endLoading()
  }
}

// 加载模板
async function loadTemplates() {
  try {
    const { data } = await fetchWorkflowTemplates()
    if (data) {
      templates.value = data.templates || []
    }
  }
  catch (error: any) {
    window.$message.error(error.message || '获取模板失败')
  }
}

// 新建工作流
function handleCreate() {
  router.push('/order/workflow/designer')
}

// 查看工作流
function handleView(row: Api.WorkflowDesigner.WorkflowDefinition) {
  router.push({
    path: '/order/workflow/designer',
    query: { id: row.id, mode: 'view' },
  })
}

// 编辑工作流
function handleEdit(row: Api.WorkflowDesigner.WorkflowDefinition) {
  router.push({
    path: '/order/workflow/designer',
    query: { id: row.id },
  })
}

// 复制工作流
function handleCopy(row: Api.WorkflowDesigner.WorkflowDefinition) {
  router.push({
    path: '/order/workflow/designer',
    query: { copyFrom: row.id },
  })
}

// 删除工作流
async function handleDelete(id: string) {
  try {
    await deleteWorkflowDefinition(id)
    window.$message.success('删除成功')
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '删除失败')
  }
}

// 切换激活状态
async function handleToggleActive(row: Api.WorkflowDesigner.WorkflowDefinition, active: boolean) {
  try {
    if (active) {
      await activateWorkflowDefinition(row.id)
      window.$message.success('激活成功')
    }
    else {
      await deactivateWorkflowDefinition(row.id)
      window.$message.success('停用成功')
    }
    fetchData()
  }
  catch (error: any) {
    window.$message.error(error.message || '操作失败')
  }
}

// 打开模板库
function handleTemplateLibrary() {
  showTemplateModal.value = true
}

// 应用模板
function handleApplyTemplate(template: Api.WorkflowDesigner.WorkflowTemplate) {
  router.push({
    path: '/order/workflow/designer',
    query: { template: template.id },
  })
}
</script>

<style scoped>
/* 样式 */
</style>

