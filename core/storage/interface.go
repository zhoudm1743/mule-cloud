package storage

import (
	"context"
	"io"
)

// Storage 存储接口
type Storage interface {
	// Upload 上传文件
	// ctx: 上下文
	// key: 存储键（文件路径或对象键）
	// reader: 文件读取器
	// size: 文件大小
	// contentType: 文件MIME类型
	// 返回: 访问URL, 错误
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) (string, error)

	// Download 下载文件
	// ctx: 上下文
	// key: 存储键
	// 返回: 文件读取器, 错误
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete 删除文件
	// ctx: 上下文
	// key: 存储键
	// 返回: 错误
	Delete(ctx context.Context, key string) error

	// Exists 检查文件是否存在
	// ctx: 上下文
	// key: 存储键
	// 返回: 是否存在, 错误
	Exists(ctx context.Context, key string) (bool, error)

	// GetURL 获取文件访问URL
	// ctx: 上下文
	// key: 存储键
	// 返回: URL, 错误
	GetURL(ctx context.Context, key string) (string, error)

	// GetPresignedURL 获取预签名URL（用于临时访问私有文件）
	// ctx: 上下文
	// key: 存储键
	// expireSeconds: 过期时间（秒）
	// 返回: URL, 错误
	GetPresignedURL(ctx context.Context, key string, expireSeconds int) (string, error)
}
