import { request } from '../http'

// 获取所有菜单（扁平结构）
export function fetchAllMenus() {
  return request.Get<Service.ResponseResult<Api.Menu.MenuItem[]>>('/admin/perms/menus/all')
}

// 获取单个菜单
export function fetchMenuById(id: string) {
  return request.Get<Service.ResponseResult<Api.Menu.MenuItem>>(`/admin/perms/menus/${id}`)
}

// 分页查询菜单
export function fetchMenuList(params: Api.Menu.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Menu.ListResponse>>('/admin/perms/menus', { params })
}

// 创建菜单
export function createMenu(data: Api.Menu.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Menu.MenuItem>>('/admin/perms/menus', data)
}

// 更新菜单
export function updateMenu(id: string, data: Api.Menu.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/perms/menus/${id}`, data)
}

// 删除菜单
export function deleteMenu(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/perms/menus/${id}`)
}

// 批量删除菜单
export function batchDeleteMenus(ids: string[]) {
  return request.Post<Service.ResponseResult<any>>('/admin/perms/menus/batch-delete', { ids })
}

