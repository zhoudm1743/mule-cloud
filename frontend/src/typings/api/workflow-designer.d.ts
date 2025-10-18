declare namespace Api {
  namespace WorkflowDesigner {
    /** 工作流状态定义 */
    interface WorkflowState {
      id: string
      name: string
      code: string
      type: 'start' | 'normal' | 'end'
      color: string
      description: string
      position?: StatePosition
      metadata?: Record<string, any>
    }

    /** 状态位置 */
    interface StatePosition {
      x: number
      y: number
    }

    /** 转换条件 */
    interface TransitionCondition {
      type: 'field' | 'script' | 'custom'
      field: string
      operator: 'eq' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'contains'
      value: any
      script: string
      description: string
      metadata?: Record<string, any>
    }

    /** 转换动作 */
    interface TransitionAction {
      type: 'update_field' | 'send_notification' | 'custom'
      field: string
      value: any
      script: string
      description: string
      metadata?: Record<string, any>
    }

    /** 工作流转换规则 */
    interface WorkflowTransition {
      id: string
      name: string
      from_state: string
      to_state: string
      event: string
      conditions: TransitionCondition[]
      actions: TransitionAction[]
      require_role: string
      description: string
      metadata?: Record<string, any>
    }

    /** 工作流定义 */
    interface WorkflowDefinition {
      id: string
      name: string
      code: string
      description: string
      states: WorkflowState[]
      transitions: WorkflowTransition[]
      version: number
      is_active: boolean
      metadata?: Record<string, any>
      created_at: number
      updated_at: number
      created_by: string
      updated_by: string
    }

    /** 创建/更新工作流请求 */
    interface WorkflowDefinitionRequest {
      name: string
      code: string
      description?: string
      states: WorkflowState[]
      transitions: WorkflowTransition[]
      metadata?: Record<string, any>
    }

    /** 工作流列表请求 */
    interface WorkflowListRequest {
      page: number
      page_size: number
    }

    /** 工作流列表响应 */
    interface WorkflowListResponse {
      workflows: WorkflowDefinition[]
      total: number
    }

    /** LogicFlow节点数据 */
    interface LogicFlowNode {
      id: string
      type: string
      x: number
      y: number
      text: {
        value: string
        x: number
        y: number
      }
      properties: {
        stateId: string
        stateName: string
        stateType: 'start' | 'normal' | 'end'
        stateColor: string
        description?: string
      }
    }

    /** LogicFlow边数据 */
    interface LogicFlowEdge {
      id: string
      type: string
      sourceNodeId: string
      targetNodeId: string
      startPoint: { x: number; y: number }
      endPoint: { x: number; y: number }
      text?: {
        value: string
        x: number
        y: number
      }
      properties: {
        transitionId: string
        event: string
        conditions?: TransitionCondition[]
        actions?: TransitionAction[]
        requireRole?: string
        description?: string
      }
    }

    /** LogicFlow图数据 */
    interface LogicFlowData {
      nodes: LogicFlowNode[]
      edges: LogicFlowEdge[]
    }

    /** 工作流模板 */
    interface WorkflowTemplate {
      id: string
      name: string
      code?: string
      description: string
      category: string
      icon: string
      preview: string
      states?: WorkflowTemplateState[]
      transitions?: WorkflowTemplateTransition[]
    }

    /** 模板状态定义 */
    interface WorkflowTemplateState {
      code: string
      name: string
      type: 'start' | 'normal' | 'end'
      color: string
      description: string
    }

    /** 模板转换定义 */
    interface WorkflowTemplateTransition {
      from: string
      to: string
      event: string
      event_label: string
      has_condition?: boolean
      condition_desc?: string
      require_role?: string
      role_desc?: string
      available_fields?: WorkflowConditionField[]
    }

    /** 可用的条件字段 */
    interface WorkflowConditionField {
      key: string
      label: string
      type: 'number' | 'string' | 'boolean'
      description: string
    }

    /** 模板列表响应 */
    interface TemplateListResponse {
      templates: WorkflowTemplate[]
    }

    /** 工作流实例 */
    interface WorkflowInstance {
      id: string
      workflow_id: string
      entity_type: string
      entity_id: string
      current_state: string
      history: WorkflowHistoryRecord[]
      variables: Record<string, any>
      created_at: number
      updated_at: number
    }

    /** 工作流历史记录 */
    interface WorkflowHistoryRecord {
      from_state: string
      to_state: string
      event: string
      operator: string
      reason: string
      timestamp: number
      metadata?: Record<string, any>
    }

    /** 执行转换请求 */
    interface ExecuteTransitionRequest {
      instance_id: string
      event: string
      reason?: string
      metadata?: Record<string, any>
    }
  }
}

