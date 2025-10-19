import { request } from '../http'

// ==================== 客户 (Customer) ====================

// 分页查询客户
export function fetchCustomerList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Customer.ListResponse>>('/admin/basic/customers', { params })
}

// 获取所有客户（不分页）
export function fetchAllCustomers(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Customer.ListResponse>>('/admin/basic/customers/all', { params })
}

// 创建客户
export function createCustomer(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/customers', data)
}

// 更新客户
export function updateCustomer(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/customers/${id}`, data)
}

// 删除客户
export function deleteCustomer(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/customers/${id}`)
}

// ==================== 颜色 (Color) ====================

// 分页查询颜色
export function fetchColorList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Color.ListResponse>>('/admin/basic/colors', { params })
}

// 获取所有颜色（不分页）
export function fetchAllColors(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Color.ListResponse>>('/admin/basic/colors/all', { params })
}

// 创建颜色
export function createColor(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/colors', data)
}

// 更新颜色
export function updateColor(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/colors/${id}`, data)
}

// 删除颜色
export function deleteColor(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/colors/${id}`)
}

// ==================== 业务员 (Salesman) ====================

// 分页查询业务员
export function fetchSalesmanList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Salesman.ListResponse>>('/admin/basic/salesmans', { params })
}

// 获取所有业务员（不分页）
export function fetchAllSalesmans(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Salesman.ListResponse>>('/admin/basic/salesmans/all', { params })
}

// 创建业务员
export function createSalesman(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/salesmans', data)
}

// 更新业务员
export function updateSalesman(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/salesmans/${id}`, data)
}

// 删除业务员
export function deleteSalesman(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/salesmans/${id}`)
}

// ==================== 尺码 (Size) ====================

// 分页查询尺码
export function fetchSizeList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Size.ListResponse>>('/admin/basic/sizes', { params })
}

// 获取所有尺码（不分页）
export function fetchAllSizes(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Size.ListResponse>>('/admin/basic/sizes/all', { params })
}

// 创建尺码
export function createSize(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/sizes', data)
}

// 更新尺码
export function updateSize(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/sizes/${id}`, data)
}

// 删除尺码
export function deleteSize(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/sizes/${id}`)
}

// ==================== 订单类型 (OrderType) ====================

// 分页查询订单类型
export function fetchOrderTypeList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.OrderType.ListResponse>>('/admin/basic/order_types', { params })
}

// 获取所有订单类型（不分页）
export function fetchAllOrderTypes(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.OrderType.ListResponse>>('/admin/basic/order_types/all', { params })
}

// 创建订单类型
export function createOrderType(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/order_types', data)
}

// 更新订单类型
export function updateOrderType(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/order_types/${id}`, data)
}

// 删除订单类型
export function deleteOrderType(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/order_types/${id}`)
}

// ==================== 工序 (Procedure) ====================

// 分页查询工序
export function fetchProcedureList(params: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Procedure.ListResponse>>('/admin/basic/procedures', { params })
}

// 获取所有工序（不分页）
export function fetchAllProcedures(params?: Api.Basic.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Basic.Procedure.ListResponse>>('/admin/basic/procedures/all', { params })
}

// 创建工序
export function createProcedure(data: Api.Basic.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Basic.BasicInfo>>('/admin/basic/procedures', data)
}

// 更新工序
export function updateProcedure(id: string, data: Api.Basic.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/basic/procedures/${id}`, data)
}

// 删除工序
export function deleteProcedure(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/basic/procedures/${id}`)
}

