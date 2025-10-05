import { request } from '../http'

// 获取所有管理员（不分页）
export function fetchAllAdmins() {
  return request.Get<Service.ResponseResult<Api.Admin.AdminInfo[]>>('/perms/admins/all')
}

// 获取单个管理员
export function fetchAdminById(id: string) {
  return request.Get<Service.ResponseResult<Api.Admin.AdminInfo>>(`/perms/admins/${id}`)
}

// 分页查询管理员
export function fetchAdminList(params: Api.Admin.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Admin.ListResponse>>('/perms/admins', { params })
}

// 创建管理员
export function createAdmin(data: Api.Admin.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Admin.AdminInfo>>('/perms/admins', data)
}

// 更新管理员
export function updateAdmin(id: string, data: Api.Admin.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/perms/admins/${id}`, data)
}

// 删除管理员
export function deleteAdmin(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/perms/admins/${id}`)
}

// 分配角色给管理员
export function assignAdminRoles(id: string, data: Api.Admin.AssignRolesRequest) {
  return request.Post<Service.ResponseResult<any>>(`/perms/admins/${id}/roles`, data)
}

// 获取管理员的角色
export function fetchAdminRoles(id: string) {
  return request.Get<Service.ResponseResult<string[]>>(`/perms/admins/${id}/roles`)
}

// 移除管理员的角色
export function removeAdminRole(adminId: string, roleId: string) {
  return request.Delete<Service.ResponseResult<any>>(`/perms/admins/${adminId}/roles/${roleId}`)
}

