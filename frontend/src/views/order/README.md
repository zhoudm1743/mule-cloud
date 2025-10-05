# 订单模块前端实现

## 概述

订单模块包含订单管理和款式管理两大功能，支持完整的订单生命周期管理和款式库管理。

## 目录结构

```
views/order/
├── orders/                  # 订单管理
│   ├── index.vue           # 订单列表页面
│   └── components/
│       └── TableModal.vue  # 订单编辑弹窗（三步骤）
├── styles/                  # 款式管理
│   ├── index.vue           # 款式列表页面
│   └── components/
│       └── TableModal.vue  # 款式编辑弹窗
└── README.md

service/api/
└── order.ts                 # 订单API服务

typings/api/
└── order.d.ts              # 订单类型定义
```

## 功能特性

### 订单管理 (orders/)

#### 列表功能
- **多条件搜索**：支持按合同号、款号、状态搜索
- **分页展示**：支持自定义每页显示数量
- **批量操作**：支持批量删除
- **状态标识**：不同颜色标签显示订单状态（草稿、已下单、生产中、已完成、已取消）

#### 订单创建（三步骤流程）

**步骤1：基础信息**
- 合同号（必填）
- 客户选择（必填，下拉选择）
- 交货日期（日期选择器）
- 订单类型（下拉选择）
- 业务员（下拉选择）
- 备注

**步骤2：款式数量**
- 款式选择（从款式库选择）
- 自动带出：颜色、尺码、单价、工序
- 单价设置
- 颜色尺码数量配置（动态表格）
  - 显示所有颜色+尺码组合
  - 每个组合可单独设置数量
  - 自动计算总数量

**步骤3：工序清单**
- 从款式自动复制工序
- 可编辑工价
- 可指定工人
- 可标记最终工序
- 可设置不分扎标记
- 完成后订单状态变为"已下单"

#### 订单操作
- **查看详情**：查看完整订单信息
- **编辑**：修改订单信息
- **复制**：快速复制订单创建新订单
- **删除**：软删除订单

### 款式管理 (styles/)

#### 列表功能
- **多条件搜索**：支持按款号、款名、状态搜索
- **分页展示**：支持自定义每页显示数量
- **批量操作**：支持批量删除
- **状态标识**：显示启用/禁用状态

#### 款式编辑（三标签页）

**基础信息标签**
- 款号（必填，创建后不可修改）
- 款名（必填）
- 分类
- 季节
- 年份
- 单价
- 状态（启用/禁用）
- 备注

**颜色尺码标签**
- 颜色多选
- 尺码多选

**工序清单标签**
- 添加工序
- 选择工序名称
- 设置工价
- 标记最终工序
- 设置不分扎
- 调整顺序（上移/下移）
- 删除工序

## API 接口

### 订单接口

```typescript
// 订单列表
fetchOrderList(params: Api.Order.ListRequest)

// 订单详情
fetchOrderDetail(id: string)

// 创建订单（步骤1）
createOrder(data: Api.Order.CreateRequest)

// 更新款式数量（步骤2）
updateOrderStyle(id: string, data: Api.Order.UpdateStyleRequest)

// 更新工序（步骤3）
updateOrderProcedure(id: string, data: Api.Order.UpdateProcedureRequest)

// 更新订单
updateOrder(id: string, data: Api.Order.UpdateRequest)

// 复制订单
copyOrder(id: string)

// 删除订单
deleteOrder(id: string)
```

### 款式接口

```typescript
// 款式列表
fetchStyleList(params: Api.Order.StyleListRequest)

// 所有款式（不分页）
fetchAllStyles(params?: Api.Order.StyleListRequest)

// 款式详情
fetchStyleDetail(id: string)

// 创建款式
createStyle(data: Api.Order.CreateStyleRequest)

// 更新款式
updateStyle(id: string, data: Api.Order.UpdateStyleRequest)

// 删除款式
deleteStyle(id: string)
```

## 数据类型

### 订单类型

```typescript
// 订单信息
interface OrderInfo {
  id: string
  contract_no: string          // 合同号
  style_id: string             // 款式ID
  style_no: string             // 款号
  style_name: string           // 款名
  style_image?: string         // 款式图片
  customer_id: string          // 客户ID
  customer_name: string        // 客户名称
  salesman_id?: string         // 业务员ID
  salesman_name?: string       // 业务员名称
  order_type_id?: string       // 订单类型ID
  order_type_name?: string     // 订单类型名称
  quantity: number             // 总数量
  unit_price: number           // 单价
  total_amount: number         // 总金额
  delivery_date?: string       // 交货日期
  progress: number             // 进度百分比
  status: number               // 状态
  remark?: string              // 备注
  colors: string[]             // 颜色列表
  sizes: string[]              // 尺码列表
  items: OrderItem[]           // 订单明细
  procedures: OrderProcedure[] // 工序列表
  created_at: number           // 创建时间
  updated_at: number           // 更新时间
}

// 订单明细
interface OrderItem {
  color_id: string
  color_name: string
  size_id: string
  size_name: string
  quantity: number
}

// 订单工序
interface OrderProcedure {
  sequence: number            // 顺序
  procedure_id: string        // 工序ID
  procedure_name: string      // 工序名称
  unit_price: number          // 工价
  assigned_worker?: string    // 指定工人
  is_slowest: boolean         // 是否最终工序
  no_bundle: boolean          // 不分扎上报
}
```

### 款式类型

```typescript
// 款式信息
interface StyleInfo {
  id: string
  style_no: string             // 款号
  style_name: string           // 款名
  category?: string            // 分类
  season?: string              // 季节
  year?: string                // 年份
  images: string[]             // 图片列表
  colors: string[]             // 颜色列表
  sizes: string[]              // 尺码列表
  unit_price: number           // 单价
  remark?: string              // 备注
  procedures: StyleProcedure[] // 工序列表
  status: number               // 状态
  created_at: number
  updated_at: number
}

// 款式工序
interface StyleProcedure {
  sequence: number
  procedure_id: string
  procedure_name: string
  unit_price: number
  assigned_worker?: string
  is_slowest: boolean
  no_bundle: boolean
}
```

## 组件依赖

### UI 组件（Naive UI）
- `NCard` - 卡片容器
- `NDataTable` - 数据表格
- `NModal` - 弹窗
- `NForm` - 表单
- `NInput` - 输入框
- `NSelect` - 下拉选择
- `NDatePicker` - 日期选择器
- `NInputNumber` - 数字输入框
- `NButton` - 按钮
- `NSpace` - 间距布局
- `NGrid` - 网格布局
- `NSteps` - 步骤条
- `NTabs` - 标签页
- `NTag` - 标签
- `NCheckbox` - 复选框
- `NRadioGroup` - 单选组
- `NPopconfirm` - 气泡确认框
- `NImage` - 图片

### 自定义组件
- `CopyText` - 可复制文本组件

### Hooks
- `useBoolean` - 布尔值状态管理

## 使用说明

### 1. 创建订单

```vue
<script setup>
import { ref } from 'vue'

const tableModalRef = ref()

// 打开创建订单弹窗
function createOrder() {
  tableModalRef.value?.openModal('add')
}
</script>

<template>
  <NButton @click="createOrder">新建订单</NButton>
  <TableModal ref="tableModalRef" @refresh="fetchData" />
</template>
```

### 2. 编辑订单

```vue
<script setup>
const tableModalRef = ref()

// 打开编辑弹窗
function editOrder(order) {
  tableModalRef.value?.openModal('edit', order)
}
</script>
```

### 3. 查看订单详情

```vue
<script setup>
const tableModalRef = ref()

// 打开查看弹窗
function viewOrder(order) {
  tableModalRef.value?.openModal('view', order)
}
</script>
```

## 注意事项

1. **订单创建流程**
   - 必须按步骤完成，不能跳过
   - 步骤1创建后会生成订单ID
   - 步骤2、3基于订单ID更新数据
   - 只有完成步骤3后订单才会变为"已下单"状态

2. **款式关联**
   - 选择款式后会自动带出颜色、尺码、工序
   - 可以在订单中修改这些信息
   - 修改不会影响款式库中的原始数据

3. **数量计算**
   - 订单明细中的数量会自动汇总为总数量
   - 总金额 = 总数量 × 单价

4. **工序管理**
   - 工序可以调整顺序
   - 可以标记最终工序（用于生产调度）
   - 可以设置不分扎上报（特殊工序）

5. **状态流转**
   - 草稿(0) → 已下单(1) → 生产中(2) → 已完成(3)
   - 任何状态都可以取消(4)

## 路由配置

需要在路由配置中添加：

```typescript
{
  path: '/order',
  name: 'order',
  children: [
    {
      path: 'orders',
      name: 'order-orders',
      component: () => import('@/views/order/orders/index.vue'),
      meta: { title: '订单管理' }
    },
    {
      path: 'styles',
      name: 'order-styles',
      component: () => import('@/views/order/styles/index.vue'),
      meta: { title: '款式管理' }
    }
  ]
}
```

## 后续优化

- [ ] 添加订单详情页面（多标签页：订单信息、裁片监控、进度监控、订单白示、计件工资）
- [ ] 支持图片上传功能
- [ ] 支持批量导入订单
- [ ] 支持导出Excel
- [ ] 添加订单统计报表
- [ ] 添加进度条显示
- [ ] 支持订单打印
- [ ] 支持裁剪制菲打印
- [ ] 添加订单状态流转审核

## 相关文档

- [后端API文档](../../../docs/订单服务实现完成.md)
- [订单服务README](../../../app/order/README.md)
