import { request } from '../http'

interface Ilogin {
  phone: string // 后端使用 phone 字段，不是 userName
  password: string
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

export function fetchUserRoutes(params: { id: number }) {
  return request.Get<Service.ResponseResult<AppRoute.RowRoute[]>>('/auth/getUserRoutes', { params })
}

export function fetchUserProfile() {
  return request.Get<Service.ResponseResult<Api.Login.Info>>('/auth/profile')
}
