package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	uploadDir string // 上传目录
	baseURL   string // 访问基础URL
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(config LocalConfig) (*LocalStorage, error) {
	// 确保上传目录存在
	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	return &LocalStorage{
		uploadDir: config.UploadDir,
		baseURL:   config.BaseURL,
	}, nil
}

// Upload 上传文件
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error) {
	// 构建完整路径
	fullPath := filepath.Join(s.uploadDir, key)

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入文件
	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回访问URL
	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}

// Download 下载文件
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.uploadDir, key)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	return file, nil
}

// Delete 删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.uploadDir, key)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.uploadDir, key)

	_, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetURL 获取文件访问URL
func (s *LocalStorage) GetURL(ctx context.Context, key string) (string, error) {
	url := fmt.Sprintf("%s/%s", s.baseURL, key)
	return url, nil
}

// GetPresignedURL 获取预签名URL（本地存储不需要预签名）
func (s *LocalStorage) GetPresignedURL(ctx context.Context, key string, expireSeconds int) (string, error) {
	return s.GetURL(ctx, key)
}
