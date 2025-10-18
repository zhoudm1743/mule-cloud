import { request } from '../http'

/** 获取工作流定义 */
export function fetchWorkflowDefinition() {
  return request.Get<Service.ResponseResult<Api.Workflow.WorkflowDefinition>>('/order/workflow/definition')
}

/** 获取 Mermaid 流程图 */
export function fetchMermaidDiagram() {
  return request.Get<Service.ResponseResult<Api.Workflow.MermaidDiagram>>('/order/workflow/mermaid')
}

/** 获取转换规则 */
export function fetchTransitionRules() {
  return request.Get<Service.ResponseResult<{ rules: Api.Workflow.TransitionRule[] }>>('/order/workflow/rules')
}

/** 获取订单状态 */
export function fetchOrderStatus(orderId: string) {
  return request.Get<Service.ResponseResult<Api.Workflow.OrderStatus>>(`/order/workflow/orders/${orderId}/status`)
}

/** 获取状态历史 */
export function fetchOrderHistory(orderId: string, limit = 20) {
  return request.Get<Service.ResponseResult<{ order_id: string; history: Api.Workflow.StateHistory[] }>>(
    `/order/workflow/orders/${orderId}/history`,
    { params: { limit } },
  )
}

/** 获取回滚历史 */
export function fetchRollbackHistory(orderId: string, limit = 10) {
  return request.Get<Service.ResponseResult<{ order_id: string; rollbacks: Api.Workflow.RollbackRecord[] }>>(
    `/order/workflow/orders/${orderId}/rollbacks`,
    { params: { limit } },
  )
}

/** 执行状态转换 */
export function transitionOrder(data: Api.Workflow.TransitionRequest) {
  return request.Post<Service.ResponseResult<{ message: string }>>('/order/workflow/transition', data)
}

/** 回滚状态 */
export function rollbackOrder(data: Api.Workflow.RollbackRequest) {
  return request.Post<Service.ResponseResult<{ message: string }>>('/order/workflow/rollback', data)
}

