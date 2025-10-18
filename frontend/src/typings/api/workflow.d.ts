declare namespace Api {
  namespace Workflow {
    /** 状态定义 */
    interface StateDefinition {
      id: number
      name: string
      type: 'start' | 'normal' | 'end'
      color: string
    }

    /** 事件定义 */
    interface EventDefinition {
      name: string
      label: string
      requireCondition?: boolean
      requireRole?: string
    }

    /** 转换规则 */
    interface TransitionDefinition {
      from: number
      to: number
      event: string
      hasCondition?: boolean
      conditionDesc?: string
      requireRole?: string
      roleDesc?: string
    }

    /** 工作流定义 */
    interface WorkflowDefinition {
      states: StateDefinition[]
      events: EventDefinition[]
      transitions: TransitionDefinition[]
    }

    /** Mermaid 流程图 */
    interface MermaidDiagram {
      diagram: string
    }

    /** 转换规则响应 */
    interface TransitionRule {
      from: number
      from_name: string
      to: number
      to_name: string
      event: string
      has_condition?: boolean
      require_role?: string
    }

    /** 订单状态响应 */
    interface OrderStatus {
      order_id: string
      status: number
      status_name: string
    }

    /** 状态历史 */
    interface StateHistory {
      order_id: string
      from_state: number
      to_state: number
      event: string
      reason: string
      operator: string
      timestamp: number
      metadata?: Record<string, any>
    }

    /** 回滚记录 */
    interface RollbackRecord {
      order_id: string
      rollback_from: number
      rollback_to: number
      original_event: string
      reason: string
      operator: string
      timestamp: number
      metadata?: Record<string, any>
    }

    /** 状态转换请求 */
    interface TransitionRequest {
      order_id: string
      event: string
      reason?: string
      metadata?: Record<string, any>
    }

    /** 回滚请求 */
    interface RollbackRequest {
      order_id: string
      reason: string
    }
  }
}

