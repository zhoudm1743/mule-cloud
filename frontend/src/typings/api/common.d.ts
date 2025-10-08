declare namespace Api {
  namespace Common {
    // 文件信息
    interface FileInfo {
      id: string
      file_name: string
      file_size: number
      file_type: string
      file_ext: string
      url: string
      business_type: string
      upload_by: string
      created_at: string
    }

    // 上传响应
    interface UploadResponse {
      id: string
      file_name: string
      file_size: number
      file_type: string
      file_ext: string
      url: string
      business_type: string
      created_at: string
    }

    // 文件列表请求
    interface FileListRequest {
      page: number
      page_size: number
      business_type?: string
    }

    // 文件列表响应
    interface FileListResponse {
      list: FileInfo[]
      total: number
    }

    // 预签名URL请求
    interface PresignedURLRequest {
      expire_seconds?: number
    }

    // 预签名URL响应
    interface PresignedURLResponse {
      url: string
      expire_at: string
    }
  }
}
