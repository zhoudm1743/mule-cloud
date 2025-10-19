import { request } from '../http'

/**
 * 上传文件
 */
export function uploadFile(file: File, businessType: string) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('business_type', businessType)

  // 注意：不要手动设置 Content-Type，让浏览器自动设置（包含 boundary）
  return request.Post<Service.ResponseResult<Api.Common.UploadResponse>>('/admin/common/files/upload', formData)
}

/**
 * 获取文件列表
 */
export function fetchFileList(params: Api.Common.FileListRequest) {
  return request.Get<Service.ResponseResult<Api.Common.FileListResponse>>('/admin/common/files', { params })
}

/**
 * 下载文件
 */
export function downloadFile(id: string) {
  return request.Get(`/admin/common/files/${id}`, {
    responseType: 'blob'
  })
}

/**
 * 删除文件
 */
export function deleteFile(id: string) {
  return request.Delete<Service.ResponseResult<any>>(`/admin/common/files/${id}`)
}

/**
 * 获取预签名URL
 */
export function getPresignedURL(id: string, expireSeconds?: number) {
  return request.Get<Service.ResponseResult<Api.Common.PresignedURLResponse>>(`/common/files/${id}/presigned`, {
    params: { expire_seconds: expireSeconds || 3600 }
  })
}
