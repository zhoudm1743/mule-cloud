import { request } from '../http'

// ==================== 工序上报 ====================

// 获取上报记录列表
export function fetchReportList(params: Api.Production.ReportListRequest) {
  return request.Get<Service.ResponseResult<Api.Production.ReportListResponse>>('/api/production/reports', { params })
}

// 获取上报记录详情
export function fetchReportDetail(id: string) {
  return request.Get<Service.ResponseResult<{ report: Api.Production.ProcedureReport }>>(`/api/production/reports/${id}`)
}

// 删除上报记录
export function deleteReport(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/api/production/reports/${id}`)
}

// 获取订单进度
export function fetchOrderProgress(orderId: string) {
  return request.Get<Service.ResponseResult<Api.Production.OrderProgressResponse>>(`/api/production/progress/${orderId}`)
}

// 获取工资统计
export function fetchSalary(params: Api.Production.SalaryRequest) {
  return request.Get<Service.ResponseResult<Api.Production.SalaryResponse>>('/api/production/salary', { params })
}
