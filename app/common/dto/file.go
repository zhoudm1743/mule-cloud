package dto

import "mule-cloud/internal/models"

// UploadRequest 上传请求（用于表单数据）
type UploadRequest struct {
	BusinessType string `form:"business_type" binding"required"` // 业务类型
}

// UploadResponse 上传响应
type UploadResponse struct {
	ID           string `json:"id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
	FileExt      string `json:"file_ext"`
	URL          string `json:"url"`
	BusinessType string `json:"business_type"`
	CreatedAt    string `json:"created_at"`
}

// FileListRequest 文件列表请求
type FileListRequest struct {
	Page         int    `form:"page" binding"required,min=1"`
	PageSize     int    `form:"page_size" binding"required,min=1,max=100"`
	BusinessType string `form:"business_type"` // 业务类型（可选）
}

// FileListResponse 文件列表响应
type FileListResponse struct {
	List  []FileItem `json:"list"`
	Total int64      `json:"total"`
}

// FileItem 文件项
type FileItem struct {
	ID           string `json:"id"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	FileType     string `json:"file_type"`
	FileExt      string `json:"file_ext"`
	URL          string `json:"url"`
	BusinessType string `json:"business_type"`
	UploadBy     string `json:"upload_by"`
	CreatedAt    string `json:"created_at"`
}

// PresignedURLRequest 预签名URL请求
type PresignedURLRequest struct {
	ExpireSeconds int `form:"expire_seconds" binding"min=60,max=3600"` // 过期时间（秒）
}

// PresignedURLResponse 预签名URL响应
type PresignedURLResponse struct {
	URL      string `json:"url"`
	ExpireAt string `json:"expire_at"`
}

// ModelToUploadResponse 模型转换为上传响应
func ModelToUploadResponse(file *models.FileInfo) *UploadResponse {
	return &UploadResponse{
		ID:           file.ID,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileExt:      file.FileExt,
		URL:          file.URL,
		BusinessType: file.BusinessType,
		CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ModelToFileItem 模型转换为文件项
func ModelToFileItem(file *models.FileInfo) FileItem {
	return FileItem{
		ID:           file.ID,
		FileName:     file.FileName,
		FileSize:     file.FileSize,
		FileType:     file.FileType,
		FileExt:      file.FileExt,
		URL:          file.URL,
		BusinessType: file.BusinessType,
		UploadBy:     file.UploadBy,
		CreatedAt:    file.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
