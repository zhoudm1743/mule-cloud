import { request } from '../http'

/** 获取工作流定义列表 */
export function fetchWorkflowDefinitions(params: Api.WorkflowDesigner.WorkflowListRequest) {
  return request.Get<Service.ResponseResult<Api.WorkflowDesigner.WorkflowListResponse>>(
    '/order/workflow/designer/definitions',
    { params },
  )
}

/** 创建工作流定义 */
export function createWorkflowDefinition(data: Api.WorkflowDesigner.WorkflowDefinitionRequest) {
  return request.Post<Service.ResponseResult<Api.WorkflowDesigner.WorkflowDefinition>>(
    '/order/workflow/designer/definitions',
    data,
  )
}

/** 获取工作流定义详情 */
export function fetchWorkflowDefinition(id: string) {
  return request.Get<Service.ResponseResult<Api.WorkflowDesigner.WorkflowDefinition>>(
    `/order/workflow/designer/definitions/${id}`,
  )
}

/** 更新工作流定义 */
export function updateWorkflowDefinition(id: string, data: Api.WorkflowDesigner.WorkflowDefinitionRequest) {
  return request.Put<Service.ResponseResult<{ message: string }>>(
    `/order/workflow/designer/definitions/${id}`,
    data,
  )
}

/** 删除工作流定义 */
export function deleteWorkflowDefinition(id: string) {
  return request.Delete<Service.ResponseResult<{ message: string }>>(
    `/order/workflow/designer/definitions/${id}`,
  )
}

/** 激活工作流定义 */
export function activateWorkflowDefinition(id: string) {
  return request.Post<Service.ResponseResult<{ message: string }>>(
    `/order/workflow/designer/definitions/${id}/activate`,
  )
}

/** 停用工作流定义 */
export function deactivateWorkflowDefinition(id: string) {
  return request.Post<Service.ResponseResult<{ message: string }>>(
    `/order/workflow/designer/definitions/${id}/deactivate`,
  )
}

/** 获取工作流模板列表 */
export function fetchWorkflowTemplates() {
  return request.Get<Service.ResponseResult<Api.WorkflowDesigner.TemplateListResponse>>(
    '/order/workflow/designer/templates',
  )
}

/** 获取工作流实例 */
export function fetchWorkflowInstance(entityType: string, entityId: string) {
  return request.Get<Service.ResponseResult<Api.WorkflowDesigner.WorkflowInstance>>(
    '/order/workflow/designer/instances',
    {
      params: {
        entity_type: entityType,
        entity_id: entityId,
      },
    },
  )
}

/** 执行工作流转换 */
export function executeWorkflowTransition(data: Api.WorkflowDesigner.ExecuteTransitionRequest) {
  return request.Post<Service.ResponseResult<{ message: string }>>(
    '/order/workflow/designer/execute',
    data,
  )
}

