import { request } from '../http'

// 获取操作日志列表（分页）
export function fetchOperationLogList(params: Api.OperationLog.ListRequest) {
  return request.Get<Service.ResponseResult<Api.OperationLog.ListResponse>>('/admin/system/operation-logs', { params })
}

// 获取操作日志详情
export function fetchOperationLogById(id: string) {
  return request.Get<Service.ResponseResult<Api.OperationLog.DetailResponse>>(`/admin/system/operation-logs/${id}`)
}

// 获取操作日志统计
export function fetchOperationLogStats(params: Api.OperationLog.StatsRequest) {
  return request.Get<Service.ResponseResult<Api.OperationLog.StatsResponse>>('/admin/system/operation-logs/stats', { params })
}

