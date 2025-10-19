import { request } from '../http'

// 获取所有路由信息
export function fetchAllRoutes() {
  return request.Get<Service.ResponseResult<AppRoute.RowRoute[]>>('/admin/getUserRoutes')
}

// 获取所有用户信息
export function fetchUserPage() {
  return request.Get<Service.ResponseResult<Entity.User[]>>('/admin/userPage')
}
// 获取所有角色列表（已废弃，使用 role.ts 中的 fetchRoleList）
// export function fetchRoleList() {
//   return request.Get<Service.ResponseResult<Entity.Role[]>>('/admin/role/list')
// }

/**
 * 请求获取字典列表
 *
 * @param code - 字典编码，用于筛选特定的字典列表
 * @returns 返回的字典列表数据
 */
export function fetchDictList(code?: string) {
  const params = { code }
  return request.Get<Service.ResponseResult<Entity.Dict[]>>('/admin/dict/list', { params })
}
