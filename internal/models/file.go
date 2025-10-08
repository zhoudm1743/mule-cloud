package models

import (
	"time"
)

// FileInfo 文件信息
type FileInfo struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	TenantCode   string    `json:"tenant_code" bson:"tenant_code"`     // 租户代码
	FileName     string    `json:"file_name" bson:"file_name"`         // 原始文件名
	FileSize     int64     `json:"file_size" bson:"file_size"`         // 文件大小（字节）
	FileType     string    `json:"file_type" bson:"file_type"`         // 文件类型（MIME类型）
	FileExt      string    `json:"file_ext" bson:"file_ext"`           // 文件扩展名
	StorageType  string    `json:"storage_type" bson:"storage_type"`   // 存储类型（local/minio/s3/oss）
	StoragePath  string    `json:"storage_path" bson:"storage_path"`   // 存储路径
	StorageKey   string    `json:"storage_key" bson:"storage_key"`     // 存储键（用于对象存储）
	URL          string    `json:"url" bson:"url"`                     // 访问URL
	BusinessType string    `json:"business_type" bson:"business_type"` // 业务类型（avatar/document/image等）
	UploadBy     string    `json:"upload_by" bson:"upload_by"`         // 上传人ID
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

// TableName 返回集合名称
func (FileInfo) TableName() string {
	return "files"
}
