/// <reference path="../global.d.ts"/>

namespace Api {
  namespace Production {
    /* ========== 工序上报 ========== */

    /* 工序上报记录 */
    interface ProcedureReport {
      id: string
      order_id: string
      contract_no: string
      style_no: string
      style_name: string
      batch_id: string
      bundle_no: string
      color: string
      size: string
      quantity: number
      procedure_seq: number
      procedure_name: string
      unit_price: number
      total_price: number
      worker_id: string
      worker_name: string
      worker_no: string
      report_time: number
      remark: string
      created_at: number
      updated_at: number
    }

    /* 批次工序进度 */
    interface BatchProcedureProgress {
      id: string
      batch_id: string
      order_id: string
      contract_no: string
      bundle_no: string
      color: string
      size: string
      quantity: number
      procedure_seq: number
      procedure_name: string
      reported_qty: number
      is_completed: boolean
      completed_at: number
      created_at: number
      updated_at: number
    }

    /* 订单工序进度 */
    interface OrderProcedureProgress {
      id: string
      order_id: string
      contract_no: string
      total_qty: number  // 修复：字段名应与后端一致
      procedure_seq: number
      procedure_name: string
      reported_qty: number
      progress: number
      is_completed?: boolean  // 可选，后端可能不返回
      created_at: number
      updated_at: number
    }

    /* 上报记录列表请求 */
    interface ReportListRequest {
      page?: number
      page_size?: number
      contract_no?: string
      worker_id?: string
      start_date?: number
      end_date?: number
    }

    /* 上报记录列表响应 */
    interface ReportListResponse {
      reports: ProcedureReport[]
      total: number
    }

    /* 订单进度请求 */
    interface OrderProgressRequest {
      order_id: string
    }

    /* 订单进度响应 */
    interface OrderProgressResponse {
      order_id: string
      contract_no: string
      style_no: string
      total_quantity: number
      procedures: OrderProcedureProgress[]
    }

    /* 工资统计请求 */
    interface SalaryRequest {
      worker_id?: string
      start_date?: number
      end_date?: number
    }

    /* 工资统计响应 */
    interface SalaryResponse {
      worker_id: string
      worker_name: string
      total_pieces: number
      total_salary: number
      reports: ProcedureReport[]
    }

    /* ========== 质检和返工 ========== */

    /* 质检记录 */
    interface QualityInspection {
      id: string
      order_id: string
      contract_no: string
      style_no: string
      batch_id: string
      bundle_no: string
      color: string
      size: string
      quantity: number
      qualified_qty: number
      unqualified_qty: number
      inspection_time: number
      inspector_id: string
      inspector_name: string
      status: number // 1-合格 2-不合格 3-返工中
      defect_type: string
      defect_reason: string
      remark: string
      rework_id: string
      created_at: number
      updated_at: number
    }

    /* 返工记录 */
    interface ReworkRecord {
      id: string
      inspection_id: string
      order_id: string
      contract_no: string
      style_no: string
      batch_id: string
      bundle_no: string
      color: string
      size: string
      quantity: number
      defect_type: string
      defect_reason: string
      assigned_worker: string
      assigned_worker_name: string
      status: number // 1-待处理 2-处理中 3-已完成 4-已取消
      request_time: number
      start_time: number
      complete_time: number
      complete_images: string[]
      complete_remark: string
      created_at: number
      updated_at: number
    }
  }
}
