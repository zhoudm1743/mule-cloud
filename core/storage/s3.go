package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Storage AWS S3存储实现
type S3Storage struct {
	client     *s3.S3
	uploader   *s3manager.Uploader
	bucketName string
	region     string
}

// NewS3Storage 创建S3存储实例
func NewS3Storage(config S3Config) (*S3Storage, error) {
	// 创建AWS配置
	awsConfig := &aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	}

	// 如果指定了自定义端点（如MinIO兼容S3接口）
	if config.Endpoint != "" {
		awsConfig.Endpoint = aws.String(config.Endpoint)
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	// 创建会话
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, fmt.Errorf("创建AWS会话失败: %w", err)
	}

	// 创建S3客户端
	client := s3.New(sess)

	// 检查存储桶是否存在
	_, err = client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(config.BucketName),
	})
	if err != nil {
		// 存储桶不存在，尝试创建
		_, err = client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(config.BucketName),
		})
		if err != nil {
			return nil, fmt.Errorf("创建S3存储桶失败: %w", err)
		}
	}

	return &S3Storage{
		client:     client,
		uploader:   s3manager.NewUploader(sess),
		bucketName: config.BucketName,
		region:     config.Region,
	}, nil
}

// Upload 上传文件
func (s *S3Storage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 上传文件
	result, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("上传文件到S3失败: %w", err)
	}

	return result.Location, nil
}

// Download 下载文件
func (s *S3Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("从S3下载文件失败: %w", err)
	}

	return result.Body, nil
}

// Delete 删除文件
func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("从S3删除文件失败: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		// 检查是否是NotFound错误
		if _, ok := err.(s3.RequestFailure); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetURL 获取文件访问URL
func (s *S3Storage) GetURL(ctx context.Context, key string) (string, error) {
	// 构建S3 URL
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)
	return url, nil
}

// GetPresignedURL 获取预签名URL
func (s *S3Storage) GetPresignedURL(ctx context.Context, key string, expireSeconds int) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	// 生成预签名URL
	urlStr, err := req.Presign(time.Duration(expireSeconds) * time.Second)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return urlStr, nil
}
