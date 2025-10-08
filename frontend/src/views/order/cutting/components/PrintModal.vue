<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NModal, NQrCode, NSpace } from 'naive-ui'
import { useBoolean } from '@/hooks'

defineOptions({ name: 'PrintModal' })

const { bool: visible, setTrue: showModal, setFalse: hideModal } = useBoolean(false)

const batchList = ref<Api.Order.CuttingBatchInfo[]>([])

// 为批次（扎号）生成二维码数据
const getQRCodeValue = (batch: Api.Order.CuttingBatchInfo) => {
  if (!batch)
    return ''
  
  // 二维码包含整个批次（扎号）的完整信息
  const qrData = {
    order_id: batch.order_id,
    contract_no: batch.contract_no,
    style_no: batch.style_no,
    bed_no: batch.bed_no,
    bundle_no: batch.bundle_no,
    color: batch.color,
    size: batch.size_details[0]?.size, // 单个尺码
    quantity: batch.total_pieces, // 总件数
  }
  
  return JSON.stringify(qrData)
}

// 打开模态框（支持单个或多个批次）
function openModal(batches: Api.Order.CuttingBatchInfo | Api.Order.CuttingBatchInfo[]) {
  if (Array.isArray(batches)) {
    batchList.value = batches
  }
  else {
    batchList.value = [batches]
  }
  showModal()
}

// 打印
function handlePrint() {
  window.print()
}

defineExpose({
  openModal,
})
</script>

<template>
  <NModal
    v-model:show="visible"
    title="打印预览"
    preset="card"
    class="w-1400px"
    :mask-closable="false"
  >
    <div class="print-container">
      <!-- 网格布局：每行3个标签 -->
      <div class="labels-grid">
        <div v-for="batch in batchList" :key="batch.id" class="label-card">
          <!-- 二维码 -->
          <div class="qr-section">
            <NQrCode :value="getQRCodeValue(batch)" :size="100" />
          </div>
          
          <!-- 信息区域 -->
          <div class="info-section">
            <div class="info-item">
              <span class="label">订单号</span>
              <span class="value">{{ batch.contract_no }}</span>
            </div>
            <div class="info-item">
              <span class="label">款号</span>
              <span class="value">{{ batch.style_no }}</span>
            </div>
            <div class="info-item">
              <span class="label">颜色</span>
              <span class="value">{{ batch.color }}</span>
            </div>
            <div class="info-item">
              <span class="label">床号</span>
              <span class="value">{{ batch.bed_no }}</span>
            </div>
            <div class="info-item highlight">
              <span class="label">扎号</span>
              <span class="value">{{ batch.bundle_no }}</span>
            </div>
            <div class="divider" />
            <div class="info-item size">
              <span class="label">尺码</span>
              <span class="value">{{ batch.size_details[0]?.size || '-' }}</span>
            </div>
            <div class="info-item total">
              <span class="label">数量</span>
              <span class="value">{{ batch.total_pieces }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <NSpace justify="space-between">
        <span class="text-gray-500">共 {{ batchList.length }} 个批次</span>
        <NSpace>
          <NButton @click="hideModal">
            取消
          </NButton>
          <NButton type="primary" @click="handlePrint">
            打印全部
          </NButton>
        </NSpace>
      </NSpace>
    </template>
  </NModal>
</template>

<style scoped>
.print-container {
  padding: 10px;
  max-height: 600px;
  overflow-y: auto;
}

/* 网格布局：每行3列 */
.labels-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

/* 单个标签卡片 */
.label-card {
  border: 2px solid #333;
  border-radius: 4px;
  padding: 12px;
  background: white;
  display: flex;
  flex-direction: column;
  align-items: center;
  page-break-inside: avoid;
  min-height: 280px;
}

.qr-section {
  margin-bottom: 10px;
  display: flex;
  justify-content: center;
}

.info-section {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  padding: 3px 0;
}

.info-item .label {
  color: #666;
  font-weight: 500;
}

.info-item .value {
  color: #333;
  font-weight: 600;
}

.info-item.highlight {
  background: #f0f7ff;
  padding: 6px 8px;
  border-radius: 4px;
  margin: 4px 0;
}

.info-item.highlight .value {
  color: #1890ff;
  font-size: 16px;
  font-weight: 700;
}

.divider {
  height: 1px;
  background: #e0e0e0;
  margin: 6px 0;
}

.info-item.size .value {
  color: #1890ff;
  font-size: 18px;
  font-weight: 700;
}

.info-item.total {
  background: #fff7e6;
  padding: 6px 8px;
  border-radius: 4px;
  border: 1px solid #ffd591;
}

.info-item.total .label {
  color: #d46b08;
}

.info-item.total .value {
  color: #d46b08;
  font-size: 18px;
  font-weight: 700;
}

/* 打印样式 */
@media print {
  @page {
    size: A4;
    margin: 10mm;
  }

  .print-container {
    max-height: none !important;
    overflow: visible !important;
    padding: 0;
  }

  /* 打印时每行3列 */
  .labels-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
  }

  .label-card {
    border: 1.5px solid #000;
    padding: 8px;
    min-height: auto;
    page-break-inside: avoid;
  }

  /* 每9个标签后分页（3行×3列=9个） */
  .label-card:nth-child(9n) {
    page-break-after: always;
  }

  .qr-section {
    margin-bottom: 6px;
  }

  .info-item {
    font-size: 10px;
    padding: 2px 0;
  }

  .info-item.highlight {
    padding: 4px 6px;
  }

  .info-item.highlight .value {
    font-size: 13px;
  }

  .info-item.size .value {
    font-size: 14px;
  }

  .info-item.total {
    padding: 4px 6px;
  }

  .info-item.total .value {
    font-size: 14px;
  }

  .divider {
    margin: 4px 0;
  }

  /* 隐藏模态框装饰 */
  :deep(.n-card__header),
  :deep(.n-card__footer),
  button {
    display: none !important;
  }
}

/* 屏幕较小时改为2列 */
@media screen and (max-width: 1200px) {
  .labels-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
