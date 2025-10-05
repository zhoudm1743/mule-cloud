import { request } from '../http'

interface Ilogin {
  phone: string // 后端使用 phone 字段，不是 userName
  password: string
  tenant_code?: string // 租户代码（可选）
}

export function fetchLogin(data: Ilogin) {
  const methodInstance = request.Post<Service.ResponseResult<Api.Login.Info>>('/auth/login', data)
  methodInstance.meta = {
    authRole: null,
  }
  return methodInstance
}

export function fetchRegister(data: Api.Login.RegisterRequest) {
  const methodInstance = request.Post<Service.ResponseResult<Api.Login.RegisterResponse>>('/auth/register', data)
  methodInstance.meta = {
    authRole: null,
  }
  return methodInstance
}

export function fetchUpdateToken(data: { token: string }) {
  const method = request.Post<Service.ResponseResult<Api.Login.RefreshTokenResponse>>('/auth/refresh', data)
  method.meta = {
    authRole: 'refreshToken',
  }
  return method
}

// 获取用户路由返回结构
interface GetUserRoutesResponse {
  routes: AppRoute.RowRoute[]
}

export function fetchUserRoutes(params?: { id?: string }) {
  return request.Get<Service.ResponseResult<GetUserRoutesResponse>>('/auth/getUserRoutes', { params })
}

export function fetchUserProfile() {
  return request.Get<Service.ResponseResult<Api.Login.Info>>('/auth/profile')
}

// 租户列表项
interface TenantItem {
  code: string
  name: string
  status: number
}

// 获取租户列表响应
interface GetTenantListResponse {
  tenants: TenantItem[]
  total: number
}

// 获取租户列表（用于登录页面选择）
export function fetchLoginTenantList() {
  const methodInstance = request.Get<Service.ResponseResult<GetTenantListResponse>>('/auth/tenants')
  methodInstance.meta = {
    authRole: null, // 无需认证
  }
  return methodInstance
}
