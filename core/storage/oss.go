package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSStorage 阿里云OSS存储实现
type OSSStorage struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
}

// NewOSSStorage 创建OSS存储实例
func NewOSSStorage(config OSSConfig) (*OSSStorage, error) {
	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	// 获取存储桶
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取OSS存储桶失败: %w", err)
	}

	return &OSSStorage{
		client:     client,
		bucket:     bucket,
		bucketName: config.BucketName,
	}, nil
}

// Upload 上传文件
func (s *OSSStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 上传文件
	err := s.bucket.PutObject(key, reader, oss.ContentType(contentType))
	if err != nil {
		return "", fmt.Errorf("上传文件到OSS失败: %w", err)
	}

	// 返回URL
	return s.GetURL(ctx, key)
}

// Download 下载文件
func (s *OSSStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	body, err := s.bucket.GetObject(key)
	if err != nil {
		return nil, fmt.Errorf("从OSS下载文件失败: %w", err)
	}

	return body, nil
}

// Delete 删除文件
func (s *OSSStorage) Delete(ctx context.Context, key string) error {
	err := s.bucket.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("从OSS删除文件失败: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *OSSStorage) Exists(ctx context.Context, key string) (bool, error) {
	isExist, err := s.bucket.IsObjectExist(key)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

// GetURL 获取文件访问URL
func (s *OSSStorage) GetURL(ctx context.Context, key string) (string, error) {
	// 构建OSS URL
	// 格式: https://bucket-name.endpoint/key
	endpoint := s.client.Config.Endpoint
	url := fmt.Sprintf("https://%s.%s/%s", s.bucketName, endpoint, key)
	return url, nil
}

// GetPresignedURL 获取预签名URL
func (s *OSSStorage) GetPresignedURL(ctx context.Context, key string, expireSeconds int) (string, error) {
	// 生成预签名URL
	signedURL, err := s.bucket.SignURL(key, oss.HTTPGet, int64(expireSeconds))
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return signedURL, nil
}
