package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"mule-cloud/core/storage"
	"mule-cloud/internal/models"
	"mule-cloud/internal/repository"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileService 文件服务接口
type FileService interface {
	Upload(ctx context.Context, tenantCode, uploadBy, businessType string, file *multipart.FileHeader) (*models.FileInfo, error)
	Download(ctx context.Context, tenantCode, fileID string) (io.ReadCloser, *models.FileInfo, error)
	Delete(ctx context.Context, tenantCode, fileID string) error
	List(ctx context.Context, tenantCode string, page, pageSize int, businessType string) ([]*models.FileInfo, int64, error)
	GetPresignedURL(ctx context.Context, tenantCode, fileID string, expireSeconds int) (string, error)
}

// fileService 文件服务实现
type fileService struct {
	fileRepo repository.FileRepository
	storage  storage.Storage
}

// NewFileService 创建文件服务实例
func NewFileService(fileRepo repository.FileRepository, storage storage.Storage) FileService {
	return &fileService{
		fileRepo: fileRepo,
		storage:  storage,
	}
}

// Upload 上传文件
func (s *fileService) Upload(ctx context.Context, tenantCode, uploadBy, businessType string, fileHeader *multipart.FileHeader) (*models.FileInfo, error) {
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer file.Close()

	// 获取文件信息
	fileName := fileHeader.Filename
	fileSize := fileHeader.Size
	fileExt := strings.ToLower(filepath.Ext(fileName))
	contentType := fileHeader.Header.Get("Content-Type")

	// 验证文件大小（例如：限制100MB）
	const maxFileSize = 100 * 1024 * 1024 // 100MB
	if fileSize > maxFileSize {
		return nil, fmt.Errorf("文件大小超过限制(%dMB)", maxFileSize/(1024*1024))
	}

	// 验证文件类型（可根据需求配置）
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true, // 图片
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".ppt": true, ".pptx": true, // 文档
		".txt": true, ".csv": true, ".zip": true, ".rar": true, ".7z": true, // 其他
	}
	if !allowedExts[fileExt] {
		return nil, fmt.Errorf("不支持的文件类型: %s", fileExt)
	}

	// 生成存储键（路径）
	// 格式: tenant/business_type/yyyy/mm/dd/uuid.ext
	now := time.Now()
	storageKey := fmt.Sprintf("%s/%s/%s/%s",
		tenantCode,
		businessType,
		now.Format("2006/01/02"),
		uuid.New().String()+fileExt,
	)

	// 上传到存储
	url, err := s.storage.Upload(ctx, storageKey, file, fileSize, contentType)
	if err != nil {
		return nil, fmt.Errorf("上传文件到存储失败: %w", err)
	}

	// 创建文件记录
	fileInfo := &models.FileInfo{
		FileName:     fileName,
		FileSize:     fileSize,
		FileType:     contentType,
		FileExt:      fileExt,
		StorageKey:   storageKey,
		URL:          url,
		BusinessType: businessType,
		UploadBy:     uploadBy,
	}

	// 保存文件记录到数据库
	err = s.fileRepo.Create(ctx, tenantCode, fileInfo)
	if err != nil {
		// 如果数据库保存失败，尝试删除已上传的文件
		_ = s.storage.Delete(ctx, storageKey)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	return fileInfo, nil
}

// Download 下载文件
func (s *fileService) Download(ctx context.Context, tenantCode, fileID string) (io.ReadCloser, *models.FileInfo, error) {
	// 获取文件记录
	fileInfo, err := s.fileRepo.Get(ctx, tenantCode, fileID)
	if err != nil {
		return nil, nil, fmt.Errorf("获取文件记录失败: %w", err)
	}

	// 从存储下载文件
	reader, err := s.storage.Download(ctx, fileInfo.StorageKey)
	if err != nil {
		return nil, nil, fmt.Errorf("下载文件失败: %w", err)
	}

	return reader, fileInfo, nil
}

// Delete 删除文件
func (s *fileService) Delete(ctx context.Context, tenantCode, fileID string) error {
	// 获取文件记录
	fileInfo, err := s.fileRepo.Get(ctx, tenantCode, fileID)
	if err != nil {
		return fmt.Errorf("获取文件记录失败: %w", err)
	}

	// 从存储删除文件
	err = s.storage.Delete(ctx, fileInfo.StorageKey)
	if err != nil {
		return fmt.Errorf("删除存储文件失败: %w", err)
	}

	// 删除数据库记录
	err = s.fileRepo.Delete(ctx, tenantCode, fileID)
	if err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	return nil
}

// List 获取文件列表
func (s *fileService) List(ctx context.Context, tenantCode string, page, pageSize int, businessType string) ([]*models.FileInfo, int64, error) {
	return s.fileRepo.List(ctx, tenantCode, page, pageSize, businessType)
}

// GetPresignedURL 获取预签名URL
func (s *fileService) GetPresignedURL(ctx context.Context, tenantCode, fileID string, expireSeconds int) (string, error) {
	// 获取文件记录
	fileInfo, err := s.fileRepo.Get(ctx, tenantCode, fileID)
	if err != nil {
		return "", fmt.Errorf("获取文件记录失败: %w", err)
	}

	// 生成预签名URL
	url, err := s.storage.GetPresignedURL(ctx, fileInfo.StorageKey, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return url, nil
}
