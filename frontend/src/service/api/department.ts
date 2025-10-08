import { request } from '../http'

// 获取所有部门（不分页）
export function fetchAllDepartments() {
  return request.Get<Service.ResponseResult<{ departments: Api.Department.DepartmentInfo[], total: number }>>('/perms/departments/all')
}

// 获取单个部门
export function fetchDepartmentById(id: string) {
  return request.Get<Service.ResponseResult<Api.Department.DepartmentInfo>>(`/perms/departments/${id}`)
}

// 分页查询部门
export function fetchDepartmentList(params: Api.Department.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Department.ListResponse>>('/perms/departments', { params })
}

// 创建部门
export function createDepartment(data: Api.Department.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Department.DepartmentInfo>>('/perms/departments', data)
}

// 更新部门
export function updateDepartment(id: string, data: Api.Department.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/perms/departments/${id}`, data)
}

// 删除部门
export function deleteDepartment(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/perms/departments/${id}`)
}

// 批量删除部门
export function batchDeleteDepartments(ids: string[]) {
  return request.Post<Service.ResponseResult<any>>('/perms/departments/batch-delete', { ids })
}

