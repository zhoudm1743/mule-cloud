<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { createPost, fetchAllPosts, updatePost } from '@/service'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'PostModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建岗位',
    edit: '编辑岗位',
    view: '查看岗位',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

const formDefault = (): Api.Post.CreateRequest & { id?: string } => ({
  id: '',
  name: '',
  code: '',
  parent_id: '',
  status: 1,
})

const formModel = reactive(formDefault())

const rules = computed(() => ({
  name: { required: true, message: '请输入岗位名称', trigger: 'blur' },
  code: { required: true, message: '请输入岗位编码', trigger: 'blur' },
}))

const formRef = ref()
const loadingSubmit = ref(false)
const allPosts = ref<Api.Post.PostInfo[]>([])

async function getAllPosts() {
  try {
    const res = await fetchAllPosts()
    allPosts.value = res.data?.posts || []
  }
  catch (error: any) {
    window.$message.error(error.message || '获取岗位列表失败')
  }
}

const postOptions = computed(() => {
  if (!Array.isArray(allPosts.value)) {
    return []
  }
  return allPosts.value
    .filter(post => post.id !== formModel.id) // 排除自己，避免循环引用
    .map(post => ({
      label: post.name,
      value: post.id,
    }))
})

async function openModal(type: ModalType, data?: Api.Post.PostInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())

  // 加载所有岗位用于父岗位选择
  await getAllPosts()

  if (type === 'view' || type === 'edit') {
    if (data) {
      Object.assign(formModel, data)
    }
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  loadingSubmit.value = true

  try {
    if (modalType.value === 'add') {
      await createPost(formModel as Api.Post.CreateRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      const { id, ...updateData } = formModel
      if (!id) {
        window.$message.error('缺少岗位ID')
        return
      }
      await updatePost(id, updateData)
      window.$message.success('更新成功')
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
  <NModal v-model:show="visible" preset="card" :title="title" class="w-600px">
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      :label-width="100"
    >
      <NGrid :cols="1" :x-gap="18">
        <NFormItemGridItem path="name" label="岗位名称">
          <NInput v-model:value="formModel.name" :disabled="modalType === 'view'" placeholder="请输入岗位名称" />
        </NFormItemGridItem>
        <NFormItemGridItem path="code" label="岗位编码">
          <NInput v-model:value="formModel.code" :disabled="modalType === 'view'" placeholder="请输入岗位编码" />
        </NFormItemGridItem>
        <NFormItemGridItem path="parent_id" label="父岗位">
          <NSelect
            v-model:value="formModel.parent_id"
            :options="postOptions"
            :disabled="modalType === 'view'"
            placeholder="请选择父岗位（为空表示顶级岗位）"
            clearable
            filterable
          />
        </NFormItemGridItem>
        <NFormItemGridItem path="status" label="状态">
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

