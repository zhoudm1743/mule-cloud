import { request } from '../http'

// ==================== 订单 (Order) ====================

// 分页查询订单
export function fetchOrderList(params: Api.Order.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.ListResponse>>('/order/orders', { params })
}

// 获取订单详情
export function fetchOrderDetail(id: string) {
  return request.Get<Service.ResponseResult<Api.Order.OrderResponse>>(`/order/orders/${id}`)
}

// 创建订单（步骤1：基础信息）
export function createOrder(data: Api.Order.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Order.OrderResponse>>('/order/orders', data)
}

// 更新订单款式数量（步骤2）
export function updateOrderStyle(id: string, data: Api.Order.UpdateStyleRequest) {
  return request.Put<Service.ResponseResult<Api.Order.OrderResponse>>(`/order/orders/${id}/style`, data)
}

// 更新订单工序（步骤3）
export function updateOrderProcedure(id: string, data: Api.Order.UpdateProcedureRequest) {
  return request.Put<Service.ResponseResult<Api.Order.OrderResponse>>(`/order/orders/${id}/procedure`, data)
}

// 更新订单
export function updateOrder(id: string, data: Api.Order.UpdateRequest) {
  return request.Put<Service.ResponseResult<Api.Order.OrderResponse>>(`/order/orders/${id}`, data)
}

// 复制订单
export function copyOrder(id: string, data?: {
  is_related?: boolean
  relation_type?: string
  relation_remark?: string
}) {
  return request.Post<Service.ResponseResult<Api.Order.OrderResponse>>(`/order/orders/${id}/copy`, data || {})
}

// 删除订单
export function deleteOrder(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/order/orders/${id}`)
}

// ==================== 款式 (Style) ====================

// 分页查询款式
export function fetchStyleList(params: Api.Order.StyleListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.StyleListResponse>>('/order/styles', { params })
}

// 获取所有款式（不分页）
export function fetchAllStyles(params?: Api.Order.StyleListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.StyleListResponse>>('/order/styles/all', { params })
}

// 获取款式详情
export function fetchStyleDetail(id: string) {
  return request.Get<Service.ResponseResult<Api.Order.StyleResponse>>(`/order/styles/${id}`)
}

// 创建款式
export function createStyle(data: Api.Order.CreateStyleRequest) {
  return request.Post<Service.ResponseResult<Api.Order.StyleResponse>>('/order/styles', data)
}

// 更新款式
export function updateStyle(id: string, data: Api.Order.UpdateStyleRequest) {
  return request.Put<Service.ResponseResult<Api.Order.StyleResponse>>(`/order/styles/${id}`, data)
}

// 删除款式
export function deleteStyle(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/order/styles/${id}`)
}

// ==================== 裁剪任务 (Cutting Task) ====================

// 创建裁剪任务
export function createCuttingTask(data: Api.Order.CreateCuttingTaskRequest) {
  return request.Post<Service.ResponseResult<Api.Order.CuttingTaskResponse>>('/order/cutting/tasks', data)
}

// 获取裁剪任务列表
export function fetchCuttingTaskList(params: Api.Order.CuttingTaskListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingTaskListResponse>>('/order/cutting/tasks', { params })
}

// 获取裁剪任务详情
export function fetchCuttingTaskDetail(id: string) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingTaskResponse>>(`/order/cutting/tasks/${id}`)
}

// 根据订单ID获取裁剪任务
export function fetchCuttingTaskByOrderId(orderId: string) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingTaskResponse>>(`/order/cutting/tasks/order/${orderId}`)
}

// ==================== 裁剪批次 (Cutting Batch) ====================

// 创建裁剪批次
export function createCuttingBatch(data: Api.Order.CreateCuttingBatchRequest) {
  return request.Post<Service.ResponseResult<Api.Order.CuttingBatchResponse>>('/order/cutting/batches', data)
}

// 批量创建裁剪批次（制菲）
export function bulkCreateCuttingBatch(data: Api.Order.BulkCreateCuttingBatchRequest) {
  return request.Post<Service.ResponseResult<Api.Order.BulkCreateCuttingBatchResponse>>('/order/cutting/batches/bulk', data)
}

// 获取裁剪批次列表
export function fetchCuttingBatchList(params: Api.Order.CuttingBatchListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingBatchListResponse>>('/order/cutting/batches', { params })
}

// 获取裁剪批次详情
export function fetchCuttingBatchDetail(id: string) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingBatchResponse>>(`/order/cutting/batches/${id}`)
}

// 删除裁剪批次
export function deleteCuttingBatch(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/order/cutting/batches/${id}`)
}

// 打印裁剪批次
export function printCuttingBatch(id: string) {
  return request.Post<Service.ResponseResult<Api.Order.CuttingBatchResponse>>(`/order/cutting/batches/${id}/print`, {})
}

// 批量打印裁剪批次
export function batchPrintCuttingBatches(ids: string[]) {
  return request.Post<Service.ResponseResult<Api.Order.BatchPrintResponse>>('/order/cutting/batches/batch-print', { ids })
}

// ==================== 裁片监控 (Cutting Piece) ====================

// 获取裁片监控列表
export function fetchCuttingPieceList(params: Api.Order.CuttingPieceListRequest) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingPieceListResponse>>('/order/cutting/pieces', { params })
}

// 获取裁片监控详情
export function fetchCuttingPieceDetail(id: string) {
  return request.Get<Service.ResponseResult<Api.Order.CuttingPieceResponse>>(`/order/cutting/pieces/${id}`)
}

// 更新裁片进度
export function updateCuttingPieceProgress(id: string, data: Api.Order.UpdateCuttingPieceProgressRequest) {
  return request.Put<Service.ResponseResult<any>>(`/order/cutting/pieces/${id}/progress`, data)
}

// ==================== 订单工作流 (Order Workflow) ====================

// 获取订单工作流状态
export function fetchOrderWorkflowState(orderId: string) {
  return request.Get<Service.ResponseResult<Api.Order.WorkflowStateResponse>>(`/order/orders/${orderId}/workflow/state`)
}

// 获取订单可用的状态转换
export function fetchOrderAvailableTransitions(orderId: string) {
  return request.Get<Service.ResponseResult<Api.Order.WorkflowTransitionsResponse>>(`/order/orders/${orderId}/workflow/transitions`)
}

// 执行订单工作流状态转换
export function executeOrderWorkflowTransition(orderId: string, data: Api.Order.WorkflowTransitionRequest) {
  return request.Post<Service.ResponseResult<any>>(`/order/orders/${orderId}/workflow/transition`, data)
}
