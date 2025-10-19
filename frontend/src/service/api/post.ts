import { request } from '../http'

// 获取所有岗位（不分页）
export function fetchAllPosts() {
  return request.Get<Service.ResponseResult<{ posts: Api.Post.PostInfo[], total: number }>>('/admin/perms/posts/all')
}

// 获取单个岗位
export function fetchPostById(id: string) {
  return request.Get<Service.ResponseResult<Api.Post.PostInfo>>(`/admin/perms/posts/${id}`)
}

// 分页查询岗位
export function fetchPostList(params: Api.Post.ListRequest) {
  return request.Get<Service.ResponseResult<Api.Post.ListResponse>>('/admin/perms/posts', { params })
}

// 创建岗位
export function createPost(data: Api.Post.CreateRequest) {
  return request.Post<Service.ResponseResult<Api.Post.PostInfo>>('/admin/perms/posts', data)
}

// 更新岗位
export function updatePost(id: string, data: Api.Post.UpdateRequest) {
  return request.Put<Service.ResponseResult<any>>(`/admin/perms/posts/${id}`, data)
}

// 删除岗位
export function deletePost(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/perms/posts/${id}`)
}

// 批量删除岗位
export function batchDeletePosts(ids: string[]) {
  return request.Post<Service.ResponseResult<any>>('/admin/perms/posts/batch-delete', { ids })
}

