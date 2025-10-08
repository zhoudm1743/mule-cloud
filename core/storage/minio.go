package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOStorage MinIO存储实现
type MinIOStorage struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOStorage 创建MinIO存储实例
func NewMinIOStorage(config MinIOConfig) (*MinIOStorage, error) {
	// 初始化MinIO客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: config.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("初始化MinIO客户端失败: %w", err)
	}

	// 检查存储桶是否存在，不存在则创建
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("检查存储桶失败: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{Region: config.Region})
		if err != nil {
			return nil, fmt.Errorf("创建存储桶失败: %w", err)
		}
	}

	return &MinIOStorage{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

// Upload 上传文件
func (s *MinIOStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 上传文件
	_, err := s.client.PutObject(ctx, s.bucketName, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("上传文件到MinIO失败: %w", err)
	}

	// 返回URL
	return s.GetURL(ctx, key)
}

// Download 下载文件
func (s *MinIOStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("从MinIO下载文件失败: %w", err)
	}

	return object, nil
}

// Delete 删除文件
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("从MinIO删除文件失败: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *MinIOStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetURL 获取文件访问URL
func (s *MinIOStorage) GetURL(ctx context.Context, key string) (string, error) {
	// 对于公开访问的存储桶，可以直接返回URL
	// 格式: http(s)://endpoint/bucket/key
	endpoint := s.client.EndpointURL().String()
	url := fmt.Sprintf("%s/%s/%s", endpoint, s.bucketName, key)
	return url, nil
}

// GetPresignedURL 获取预签名URL
func (s *MinIOStorage) GetPresignedURL(ctx context.Context, key string, expireSeconds int) (string, error) {
	// 生成预签名URL
	expires := time.Duration(expireSeconds) * time.Second
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucketName, key, expires, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return presignedURL.String(), nil
}
