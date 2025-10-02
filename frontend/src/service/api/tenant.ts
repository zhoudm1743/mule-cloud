import { request } from '../http'

// 获取所有租户（不分页）
export function fetchAllTenants() {
  return request.Get<Service.ResponseResult<{ tenants: Api.Tenant.TenantInfo[], total: number }>>('/system/tenants/all')
}

// 获取单个租户
export function fetchTenantById(id: string) {
  return request.Get<Service.ResponseResult<Api.Tenant.TenantInfo>>(`/system/tenants/${id}`)
}

// 分页查询租户
export function fetchTenantList(params: Api.Tenant.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Tenant.ListResponse>>('/system/tenants', { params })
}

// 创建租户
export function createTenant(data: Api.Tenant.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Tenant.TenantInfo>>('/system/tenants', data)
}

// 更新租户
export function updateTenant(id: string, data: Api.Tenant.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/system/tenants/${id}`, data)
}

// 删除租户
export function deleteTenant(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/system/tenants/${id}`)
}

// 分配菜单权限给租户（超管使用）
export function assignTenantMenus(id: string, data: Api.Tenant.AssignMenusRequest) {
  return request.Post<Service.ResponseResult<any>>(`/system/tenants/${id}/menus`, data)
}

// 获取租户的菜单权限
export function fetchTenantMenus(id: string) {
  return request.Get<Service.ResponseResult<string[]>>(`/system/tenants/${id}/menus`)
}

