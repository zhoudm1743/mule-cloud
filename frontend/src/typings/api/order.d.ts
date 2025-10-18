/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Order {
    /* 订单明细 */
    interface OrderItem {
      color: string
      size: string
      quantity: number
    }

    /* 订单工序 */
    interface OrderProcedure {
      sequence: number
      procedure_name: string
      unit_price: number
      assigned_worker?: string
      is_slowest: boolean
      no_bundle: boolean
    }

    /* 订单信息 */
    interface OrderInfo {
      id: string
      contract_no: string
      style_id: string
      style_no: string
      style_name: string
      style_image?: string
      customer_id: string
      customer_name: string
      salesman_id?: string
      salesman_name?: string
      order_type_id?: string
      order_type_name?: string
      quantity: number
      unit_price: number
      total_amount: number
      delivery_date?: string
      progress: number
      status: number // 0-草稿 1-已下单 2-生产中 3-已完成 4-已取消
      remark?: string
      colors: string[]
      sizes: string[]
      items: OrderItem[]
      procedures: OrderProcedure[]
      // 工作流相关
      workflow_code?: string
      workflow_instance?: string
      workflow_state?: string
      created_at: number
      updated_at: number
    }

    /* 创建订单请求（步骤1）*/
    interface CreateRequest {
      contract_no: string
      customer_id: string
      delivery_date?: string
      order_type_id?: string
      salesman_id?: string
      remark?: string
    }

    /* 更新款式数量请求（步骤2）*/
    interface UpdateStyleRequest {
      style_id: string
      colors: string[]
      sizes: string[]
      unit_price: number
      quantity: number
      items: OrderItem[]
    }

    /* 更新工序请求（步骤3）*/
    interface UpdateProcedureRequest {
      procedures: OrderProcedure[]
    }

    /* 更新订单请求 */
    interface UpdateRequest {
      contract_no?: string
      style_id?: string
      customer_id?: string
      salesman_id?: string
      order_type_id?: string
      colors?: string[]
      sizes?: string[]
      unit_price?: number
      quantity?: number
      delivery_date?: string
      status?: number
      remark?: string
      items?: OrderItem[]
      procedures?: OrderProcedure[]
    }

    /* 查询请求 */
    interface ListRequest {
      page?: number
      page_size?: number
      contract_no?: string
      style_no?: string
      customer_id?: string
      salesman_id?: string
      order_type_id?: string
      status?: number
      start_date?: string
      end_date?: string
      order_start?: string
      order_end?: string
      remark?: string
    }

    /* 订单列表响应 */
    interface ListResponse {
      orders: OrderInfo[]
      total: number
    }

    /* 订单详情响应 */
    interface OrderResponse {
      order: OrderInfo
    }

    /* 款式工序 */
    interface StyleProcedure {
      sequence: number
      procedure_name: string
      unit_price: number
      assigned_worker?: string
      is_slowest: boolean
      no_bundle: boolean
    }

    /* 款式信息 */
    interface StyleInfo {
      id: string
      style_no: string
      style_name: string
      category?: string
      season?: string
      year?: string
      images: string[]
      colors: string[]
      sizes: string[]
      unit_price: number
      remark?: string
      procedures: StyleProcedure[]
      status: number
      created_at: number
      updated_at: number
    }

    /* 创建款式请求 */
    interface CreateStyleRequest {
      style_no: string
      style_name: string
      category?: string
      season?: string
      year?: string
      images?: string[]
      colors?: string[]
      sizes?: string[]
      unit_price?: number
      remark?: string
      procedures?: StyleProcedure[]
      status?: number
    }

    /* 更新款式请求 */
    interface UpdateStyleRequest {
      style_name?: string
      category?: string
      season?: string
      year?: string
      images?: string[]
      colors?: string[]
      sizes?: string[]
      unit_price?: number
      remark?: string
      procedures?: StyleProcedure[]
      status?: number
    }

    /* 款式查询请求 */
    interface StyleListRequest {
      page?: number
      page_size?: number
      style_no?: string
      style_name?: string
      category?: string
      season?: string
      year?: string
      status?: number
    }

    /* 款式列表响应 */
    interface StyleListResponse {
      styles: StyleInfo[]
      total: number
    }

    /* 款式详情响应 */
    interface StyleResponse {
      style: StyleInfo
    }

    /* ========== 裁剪相关 ========== */

    /* 尺码明细 */
    interface SizeDetail {
      size: string
      quantity: number
    }

    /* 裁剪任务信息 */
    interface CuttingTaskInfo {
      id: string
      order_id: string
      contract_no: string
      style_no: string
      style_name: string
      customer_name: string
      total_pieces: number
      cut_pieces: number
      status: number // 0-待裁剪 1-裁剪中 2-已完成
      batches: CuttingBatchInfo[]
      created_by?: string
      updated_by?: string
      created_at: number
      updated_at: number
    }

    /* 裁剪批次信息 */
    interface CuttingBatchInfo {
      id: string
      task_id: string
      order_id: string
      contract_no: string
      style_no: string
      bed_no: string
      bundle_no: string
      color: string
      layer_count: number
      size_details: SizeDetail[]
      total_pieces: number
      qr_code?: string
      print_count?: number
      created_by?: string
      created_at?: number
      printed_at?: number
    }

    /* 裁片监控信息 */
    interface CuttingPieceInfo {
      id: string
      order_id: string
      contract_no: string
      style_no: string
      bed_no: string
      bundle_no: string
      color: string
      size: string
      quantity: number
      progress: number
      total_process: number
      created_at: number
    }

    /* 创建裁剪任务请求 */
    interface CreateCuttingTaskRequest {
      order_id: string
      created_by?: string
    }

    /* 裁剪任务列表请求 */
    interface CuttingTaskListRequest {
      page?: number
      page_size?: number
      contract_no?: string
      style_no?: string
      status?: number
    }

    /* 裁剪任务列表响应 */
    interface CuttingTaskListResponse {
      tasks: CuttingTaskInfo[]
      total: number
    }

    /* 裁剪任务详情响应 */
    interface CuttingTaskResponse {
      task: CuttingTaskInfo
    }

    /* 创建裁剪批次请求 */
    interface CreateCuttingBatchRequest {
      task_id: string
      bed_no: string
      bundle_no: string
      color: string
      layer_count: number
      size_details: SizeDetail[]
      created_by?: string
    }

    /* 批次项（批量创建时使用） */
    interface BatchItem {
      bundle_no: string
      color: string
      layer_count: number
      size_details: SizeDetail[]
    }

    /* 批量创建裁剪批次请求 */
    interface BulkCreateCuttingBatchRequest {
      task_id: string
      bed_no: string
      batches: BatchItem[]
      created_by?: string
    }

    /* 批量创建裁剪批次响应 */
    interface BulkCreateCuttingBatchResponse {
      batches: CuttingBatchInfo[]
      count: number
    }

    /* 裁剪批次列表请求 */
    interface CuttingBatchListRequest {
      page?: number
      page_size?: number
      task_id?: string
      contract_no?: string
      bed_no?: string
      bundle_no?: string
    }

    /* 裁剪批次列表响应 */
    interface CuttingBatchListResponse {
      batches: CuttingBatchInfo[]
      total: number
    }

    /* 裁剪批次详情响应 */
    interface CuttingBatchResponse {
      batch: CuttingBatchInfo
    }

    /* 批量打印响应 */
    interface BatchPrintResponse {
      batches: CuttingBatchInfo[]
      count: number
    }

    /* 裁片监控列表请求 */
    interface CuttingPieceListRequest {
      page?: number
      page_size?: number
      order_id?: string
      contract_no?: string
      bed_no?: string
      bundle_no?: string
    }

    /* 裁片监控列表响应 */
    interface CuttingPieceListResponse {
      pieces: CuttingPieceInfo[]
      total: number
    }

    /* 裁片监控详情响应 */
    interface CuttingPieceResponse {
      piece: CuttingPieceInfo
    }

    /* 更新裁片进度请求 */
    interface UpdateCuttingPieceProgressRequest {
      progress: number
    }

    /* ========== 工作流相关 ========== */

    /* 工作流历史记录 */
    interface WorkflowHistory {
      from_state: string
      to_state: string
      event: string
      operator: string
      reason: string
      timestamp: number
      metadata?: Record<string, any>
    }

    /* 工作流实例 */
    interface WorkflowInstance {
      id: string
      workflow_id: string
      entity_type: string
      entity_id: string
      current_state: string
      history: WorkflowHistory[]
      variables: Record<string, any>
      created_at: number
      updated_at: number
    }

    /* 转换条件 */
    interface TransitionCondition {
      type: string
      field: string
      operator: string
      value: any
      description?: string
    }

    /* 转换动作 */
    interface TransitionAction {
      type: string
      field?: string
      value?: any
      description?: string
    }

    /* 工作流转换规则 */
    interface WorkflowTransition {
      id: string
      name: string
      from_state: string
      to_state: string
      event: string
      conditions: TransitionCondition[]
      actions: TransitionAction[]
      require_role?: string
      description?: string
    }

    /* 获取工作流状态响应 */
    interface WorkflowStateResponse {
      instance: WorkflowInstance
    }

    /* 获取可用转换响应 */
    interface WorkflowTransitionsResponse {
      transitions: WorkflowTransition[]
    }

    /* 执行工作流转换请求 */
    interface WorkflowTransitionRequest {
      event: string
      operator?: string
      reason?: string
      metadata?: Record<string, any>
    }
  }
}
