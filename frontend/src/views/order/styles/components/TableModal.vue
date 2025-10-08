<script setup lang="ts">
import { computed, h, reactive, ref } from 'vue'
import { NButton, NCheckbox, NInput, NInputNumber, NSelect, NSpace, NUpload, type UploadFileInfo, NRadioGroup, NRadio } from 'naive-ui'
import { createStyle, updateStyle } from '@/service/api/order'
import { fetchAllColors, fetchAllProcedures, fetchAllSizes } from '@/service/api/basic'
import { uploadFile, deleteFile } from '@/service/api/common'
import { useBoolean } from '@/hooks'

interface Emits {
  (e: 'refresh'): void
}

defineOptions({ name: 'StyleModal' })

const emit = defineEmits<Emits>()

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

type ModalType = 'add' | 'edit' | 'view'

const title = computed(() => {
  const titles: Record<ModalType, string> = {
    add: '新建款式',
    edit: '编辑款式',
    view: '查看款式',
  }
  return titles[modalType.value]
})

const modalType = ref<ModalType>('add')

// 基础数据选项
const colorOptions = ref<any[]>([])
const sizeOptions = ref<any[]>([])
const procedureOptions = ref<any[]>([])

const formDefault = (): Api.Order.CreateStyleRequest & { id?: string } => ({
  id: '',
  style_no: '',
  style_name: '',
  category: '',
  season: '',
  year: '',
  images: [],
  colors: [],
  sizes: [],
  unit_price: 0,
  remark: '',
  procedures: [],
  status: 1,
})

const formModel = reactive(formDefault())

const rules: any = {
  style_no: { required: true, message: '请输入款号', trigger: 'blur' },
  style_name: { required: true, message: '请输入款名', trigger: 'blur' },
}

const formRef = ref()
const loadingSubmit = ref(false)

// 图片上传相关
const fileList = ref<UploadFileInfo[]>([])
const uploadLoading = ref(false)

// 自定义图片上传
async function customUpload({ file, onFinish, onError }: any) {
  uploadLoading.value = true

  try {
    // 验证文件大小（10MB）
    const fileSizeMB = file.file.size / 1024 / 1024
    if (fileSizeMB > 10) {
      const error = new Error(`图片大小不能超过10MB`)
      onError()
      window.$message?.error(error.message)
      return
    }

    // 调用上传API
    const response = await uploadFile(file.file, 'style')

    // 上传成功，设置文件URL和ID
    const fileData = response.data  // alova 已经展开了响应，直接访问 data
    if (fileData) {
      file.url = fileData.url
      file.id = fileData.id
      
      // 更新formModel中的images数组
      if (!formModel.images) {
        formModel.images = []
      }
      formModel.images.push(fileData.url)
      
      onFinish()
      window.$message?.success('上传成功')
    } else {
      throw new Error('上传响应数据格式错误')
    }
  }
  catch (error: any) {
    onError()
    window.$message?.error('上传失败: ' + (error.message || '未知错误'))
  }
  finally {
    uploadLoading.value = false
  }
}

// 移除图片
async function handleRemoveImage({ file }: { file: UploadFileInfo }) {
  try {
    // 从图片数组中移除
    if (file.url && formModel.images) {
      const index = formModel.images.indexOf(file.url)
      if (index > -1) {
        formModel.images.splice(index, 1)
      }
    }
    
    // 如果有文件ID，调用删除API
    if (file.id && typeof file.id === 'string' && !file.id.startsWith('existing-')) {
      await deleteFile(file.id)
    }
    
    window.$message?.success('删除成功')
  }
  catch (error: any) {
    window.$message?.error('删除失败: ' + error.message)
  }
}

// 加载基础数据
async function loadBaseData() {
  try {
    const [colors, sizes, procedures] = await Promise.all([
      fetchAllColors(),
      fetchAllSizes(),
      fetchAllProcedures(),
    ])

    // 直接使用名称作为value，支持快速创建
    colorOptions.value = (colors.data?.colors || []).map((item: any) => ({ label: item.value, value: item.value }))
    sizeOptions.value = (sizes.data?.sizes || []).map((item: any) => ({ label: item.value, value: item.value }))
    procedureOptions.value = (procedures.data?.procedures || []).map((item: any) => ({ label: item.value, value: item.value }))
  }
  catch (error: any) {
    window.$message.error('加载基础数据失败: ' + error.message)
  }
}

function openModal(type: ModalType, data?: Api.Order.StyleInfo) {
  showModal()
  modalType.value = type
  Object.assign(formModel, formDefault())
  loadBaseData()

  // 重置图片列表
  fileList.value = []

  if (type === 'view' || type === 'edit') {
    if (data) {
      Object.assign(formModel, data)
      
      // 加载已有图片
      if (data.images && data.images.length > 0) {
        fileList.value = data.images.map((url, index) => ({
          id: `existing-${index}`,
          name: `图片${index + 1}`,
          status: 'finished',
          url,
        }))
      }
    }
  }
}

// 添加工序
function addProcedure() {
  if (!formModel.procedures) {
    formModel.procedures = []
  }
  formModel.procedures.push({
    sequence: formModel.procedures.length + 1,
    procedure_name: '',
    unit_price: 0,
    assigned_worker: '',
    is_slowest: false,
    no_bundle: false,
  })
}

// 删除工序
function removeProcedure(index: number) {
  if (!formModel.procedures) return
  formModel.procedures.splice(index, 1)
  // 重新设置顺序
  formModel.procedures.forEach((proc, idx) => {
    proc.sequence = idx + 1
  })
}

// 工序上移
function moveProcedureUp(index: number) {
  if (!formModel.procedures || index === 0) return
  const temp = formModel.procedures[index]
  formModel.procedures[index] = formModel.procedures[index - 1]
  formModel.procedures[index - 1] = temp
  // 重新设置顺序
  formModel.procedures.forEach((proc, idx) => {
    proc.sequence = idx + 1
  })
}

// 工序下移
function moveProcedureDown(index: number) {
  if (!formModel.procedures || index === formModel.procedures.length - 1) return
  const temp = formModel.procedures[index]
  formModel.procedures[index] = formModel.procedures[index + 1]
  formModel.procedures[index + 1] = temp
  // 重新设置顺序
  formModel.procedures.forEach((proc, idx) => {
    proc.sequence = idx + 1
  })
}


async function handleSubmit() {
  await formRef.value?.validate()
  
  // 验证工序：至少有一个最终工序
  if (formModel.procedures && formModel.procedures.length > 0) {
    const hasFinalProcedure = formModel.procedures.some(proc => proc.is_slowest)
    if (!hasFinalProcedure) {
      window.$message.error('必须至少选择一个最终工序')
      return
    }
  }
  
  loadingSubmit.value = true

  try {
    if (modalType.value === 'add') {
      await createStyle(formModel as Api.Order.CreateStyleRequest)
      window.$message.success('创建成功')
    }
    else if (modalType.value === 'edit') {
      if (!formModel.id) {
        window.$message.error('缺少ID')
        return
      }
      const updateData = {
        style_name: formModel.style_name,
        category: formModel.category,
        season: formModel.season,
        year: formModel.year,
        images: formModel.images,
        colors: formModel.colors || [],
        sizes: formModel.sizes || [],
        unit_price: formModel.unit_price || 0,
        remark: formModel.remark,
        procedures: formModel.procedures || [],
        status: formModel.status,
      }
      await updateStyle(formModel.id, updateData as Api.Order.UpdateStyleRequest)
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

// 工序表格列定义
const procedureColumns = computed<any>(() => [
  { title: '顺序', key: 'sequence', width: 60 },
  {
    title: '工序名称',
    key: 'procedure_name',
    width: 200,
    render: (row: Api.Order.StyleProcedure, index: number) => {
      return h(NSelect, {
        value: formModel.procedures?.[index]?.procedure_name || '',
        disabled: modalType.value === 'view',
        options: procedureOptions.value,
        placeholder: '请选择或输入工序',
        filterable: true,
        tag: true,  // 支持快速创建
        'onUpdate:value': (value: string) => {
          if (formModel.procedures && formModel.procedures[index]) {
            formModel.procedures[index].procedure_name = value
          }
        },
      })
    },
  },
  {
    title: '工价',
    key: 'unit_price',
    width: 120,
    render: (row: Api.Order.StyleProcedure, index: number) => {
      return h(NInputNumber, {
        value: formModel.procedures?.[index]?.unit_price || 0,
        disabled: modalType.value === 'view',
        min: 0,
        precision: 2,
        class: 'w-full',
        'onUpdate:value': (value: number | null) => {
          if (formModel.procedures && formModel.procedures[index]) {
            formModel.procedures[index].unit_price = value || 0
          }
        },
      })
    },
  },
  {
    title: '最终工序',
    key: 'is_slowest',
    width: 80,
    align: 'center' as const,
    render: (row: Api.Order.StyleProcedure, index: number) => {
      return h(NCheckbox, {
        checked: formModel.procedures?.[index]?.is_slowest || false,
        disabled: modalType.value === 'view',
        'onUpdate:checked': (value: boolean) => {
          if (formModel.procedures && formModel.procedures[index]) {
            // 如果选中，取消其他工序的最终工序标记
            if (value) {
              formModel.procedures.forEach((proc, idx) => {
                proc.is_slowest = idx === index
              })
            } else {
              formModel.procedures[index].is_slowest = false
            }
          }
        },
      })
    },
  },
  {
    title: '不分扎',
    key: 'no_bundle',
    width: 80,
    align: 'center' as const,
    render: (row: Api.Order.StyleProcedure, index: number) => {
      return h(NCheckbox, {
        checked: formModel.procedures?.[index]?.no_bundle || false,
        disabled: modalType.value === 'view',
        'onUpdate:checked': (value: boolean) => {
          if (formModel.procedures && formModel.procedures[index]) {
            formModel.procedures[index].no_bundle = value
          }
        },
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render: (_row: Api.Order.StyleProcedure, index: number) => {
      if (modalType.value === 'view') return null
      return h(NSpace, {}, {
        default: () => [
          h(NButton, { size: 'small', onClick: () => moveProcedureUp(index), disabled: index === 0 }, { default: () => '上移' }),
          h(NButton, { size: 'small', onClick: () => moveProcedureDown(index), disabled: index === (formModel.procedures?.length || 0) - 1 }, { default: () => '下移' }),
          h(NButton, { size: 'small', type: 'error', onClick: () => removeProcedure(index) }, { default: () => '删除' }),
        ],
      })
    },
  },
])

defineExpose({ openModal })
</script>

<template>
  <NModal v-model:show="visible" preset="card" :title="title" class="w-900px">
    <NForm
      ref="formRef"
      :model="formModel"
      :rules="rules"
      label-placement="left"
      :label-width="100"
    >
      <NTabs type="line">
        <NTabPane name="basic" tab="基础信息">
          <NGrid :cols="2" :x-gap="18">
            <NFormItemGridItem path="style_no" label="款号" :span="2">
              <NInput v-model:value="formModel.style_no" :disabled="modalType === 'view' || modalType === 'edit'" placeholder="请输入款号" />
            </NFormItemGridItem>
            <NFormItemGridItem path="style_name" label="款名" :span="2">
              <NInput v-model:value="formModel.style_name" :disabled="modalType === 'view'" placeholder="请输入款名" />
            </NFormItemGridItem>
            <NFormItemGridItem path="images" label="款式图片" :span="2">
              <div class="w-full">
                <NUpload
                  v-model:file-list="fileList"
                  :custom-request="customUpload"
                  :disabled="modalType === 'view'"
                  accept="image/*"
                  list-type="image-card"
                  :max="5"
                  @remove="handleRemoveImage"
                >
                  <NButton v-if="modalType !== 'view'" :loading="uploadLoading" :disabled="uploadLoading">
                    <template #icon>
                      <nova-icon icon="carbon:upload" :size="18" />
                    </template>
                    点击上传
                  </NButton>
                </NUpload>
                <div class="text-gray-400 text-sm mt-2">支持jpg、png格式，最多5张，每张不超过10MB</div>
              </div>
            </NFormItemGridItem>
            <NFormItemGridItem path="category" label="分类">
              <NInput v-model:value="formModel.category" :disabled="modalType === 'view'" placeholder="请输入分类" />
            </NFormItemGridItem>
            <NFormItemGridItem path="season" label="季节">
              <NInput v-model:value="formModel.season" :disabled="modalType === 'view'" placeholder="请输入季节" />
            </NFormItemGridItem>
            <NFormItemGridItem path="year" label="年份">
              <NInput v-model:value="formModel.year" :disabled="modalType === 'view'" placeholder="请输入年份" />
            </NFormItemGridItem>
            <NFormItemGridItem path="unit_price" label="单价">
              <NInputNumber v-model:value="formModel.unit_price" :disabled="modalType === 'view'" placeholder="单价" class="w-full" :precision="2" />
            </NFormItemGridItem>
            <NFormItemGridItem path="status" label="状态" :span="2">
              <NRadioGroup v-model:value="formModel.status" :disabled="modalType === 'view'">
                <NRadio :value="1">
                  启用
                </NRadio>
                <NRadio :value="0">
                  禁用
                </NRadio>
              </NRadioGroup>
            </NFormItemGridItem>
            <NFormItemGridItem path="remark" label="备注" :span="2">
              <NInput
                v-model:value="formModel.remark"
                :disabled="modalType === 'view'"
                type="textarea"
                :rows="3"
                placeholder="请输入备注"
              />
            </NFormItemGridItem>
          </NGrid>
        </NTabPane>

        <NTabPane name="spec" tab="颜色尺码">
          <NGrid :cols="1" :x-gap="18">
            <NFormItemGridItem path="colors" label="颜色">
              <NSelect
                v-model:value="formModel.colors"
                :disabled="modalType === 'view'"
                :options="colorOptions"
                placeholder="请选择颜色，可直接输入新颜色"
                multiple
                filterable
                tag
              />
            </NFormItemGridItem>
            <NFormItemGridItem path="sizes" label="尺码">
              <NSelect
                v-model:value="formModel.sizes"
                :disabled="modalType === 'view'"
                :options="sizeOptions"
                placeholder="请选择尺码，可直接输入新尺码"
                multiple
                filterable
                tag
              />
            </NFormItemGridItem>
          </NGrid>
        </NTabPane>

        <NTabPane name="procedure" tab="工序清单">
          <NSpace vertical class="w-full">
            <NButton v-if="modalType !== 'view'" type="primary" @click="addProcedure">
              <template #icon>
                <nova-icon icon="carbon:add" :size="18" />
              </template>
              添加工序
            </NButton>

            <NDataTable
              :columns="procedureColumns"
              :data="formModel.procedures || []"
              size="small"
              max-height="400px"
            />
          </NSpace>
        </NTabPane>
      </NTabs>
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
