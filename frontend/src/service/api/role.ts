import { request } from '../http'

// 获取租户下的所有角色（不分页）
export function fetchTenantRoles(tenantId: string) {
  return request.Get<Service.ResponseResult<Api.Role.RoleInfo[]>>('/admin/perms/roles/tenant', { params: { tenant_id: tenantId } })
}

// 获取单个角色
export function fetchRoleById(id: string) {
  return request.Get<Service.ResponseResult<Api.Role.RoleInfo>>(`/admin/perms/roles/${id}`)
}

// 分页查询角色
export function fetchRoleList(params: Api.Role.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Role.ListResponse>>('/admin/perms/roles', { params })
}

// 创建角色
export function createRole(data: Api.Role.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Role.RoleInfo>>('/admin/perms/roles', data)
}

// 更新角色
export function updateRole(id: string, data: Api.Role.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/perms/roles/${id}`, data)
}

// 删除角色
export function deleteRole(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/perms/roles/${id}`)
}

// 批量删除角色
export function batchDeleteRoles(ids: string[]) {
  return request.Post<Service.ResponseResult<any>>('/admin/perms/roles/batch-delete', { ids })
}

// 分配菜单权限给角色
export function assignRoleMenus(id: string, data: Api.Role.AssignMenusRequest) {
  return request.Post<Service.ResponseResult<any>>(`/admin/perms/roles/${id}/menus`, data)
}

// 获取角色的菜单权限
export function fetchRoleMenus(id: string) {
  return request.Get<Service.ResponseResult<string[]>>(`/admin/perms/roles/${id}/menus`)
}

