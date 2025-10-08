package storage

import (
	"fmt"
	"mule-cloud/core/config"
)

// NewStorage 根据配置创建存储实例
func NewStorage(cfg *config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "local":
		return NewLocalStorage(LocalConfig{
			UploadDir: cfg.Local.UploadDir,
			BaseURL:   cfg.Local.BaseURL,
		})
	case "minio":
		return NewMinIOStorage(MinIOConfig{
			Endpoint:        cfg.MinIO.Endpoint,
			AccessKeyID:     cfg.MinIO.AccessKeyID,
			SecretAccessKey: cfg.MinIO.SecretAccessKey,
			BucketName:      cfg.MinIO.BucketName,
			UseSSL:          cfg.MinIO.UseSSL,
			Region:          cfg.MinIO.Region,
		})
	case "s3":
		return NewS3Storage(S3Config{
			Endpoint:        cfg.S3.Endpoint,
			AccessKeyID:     cfg.S3.AccessKeyID,
			SecretAccessKey: cfg.S3.SecretAccessKey,
			BucketName:      cfg.S3.BucketName,
			Region:          cfg.S3.Region,
		})
	case "oss":
		return NewOSSStorage(OSSConfig{
			Endpoint:        cfg.OSS.Endpoint,
			AccessKeyID:     cfg.OSS.AccessKeyID,
			AccessKeySecret: cfg.OSS.AccessKeySecret,
			BucketName:      cfg.OSS.BucketName,
		})
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", cfg.Type)
	}
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	UploadDir string // 上传目录
	BaseURL   string // 访问基础URL
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint        string // 端点
	AccessKeyID     string // 访问密钥ID
	SecretAccessKey string // 访问密钥
	BucketName      string // 存储桶名称
	UseSSL          bool   // 是否使用SSL
	Region          string // 区域
}

// S3Config S3配置
type S3Config struct {
	Endpoint        string // 端点（可选，默认AWS S3）
	AccessKeyID     string // 访问密钥ID
	SecretAccessKey string // 访问密钥
	BucketName      string // 存储桶名称
	Region          string // 区域
}

// OSSConfig 阿里云OSS配置
type OSSConfig struct {
	Endpoint        string // 端点
	AccessKeyID     string // 访问密钥ID
	AccessKeySecret string // 访问密钥
	BucketName      string // 存储桶名称
}
