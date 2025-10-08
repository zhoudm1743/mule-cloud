package transport

import (
	"mule-cloud/app/common/dto"
	"mule-cloud/app/common/services"
	"mule-cloud/core/context"
	"mule-cloud/core/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FileTransport 文件传输层
type FileTransport struct {
	fileService services.FileService
	logger      *zap.Logger
}

// NewFileTransport 创建文件传输层实例
func NewFileTransport(fileService services.FileService, logger *zap.Logger) *FileTransport {
	return &FileTransport{
		fileService: fileService,
		logger:      logger,
	}
}

// UploadHandler 上传文件
func (t *FileTransport) UploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		tenantCode := context.GetTenantCode(c.Request.Context())
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "租户信息缺失"})
			return
		}

		// 获取用户信息
		userID, _ := c.Get("user_id")
		uploadBy := ""
		if userID != nil {
			uploadBy = userID.(string)
		}

		// 解析表单参数
		var req dto.UploadRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
			return
		}

		// 获取上传的文件
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "获取上传文件失败: " + err.Error()})
			return
		}

		// 调用服务上传文件
		fileInfo, err := t.fileService.Upload(c.Request.Context(), tenantCode, uploadBy, req.BusinessType, fileHeader)
		if err != nil {
			t.logger.Error("上传文件失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "上传文件失败: " + err.Error()})
			return
		}

		response.Success(c, dto.ModelToUploadResponse(fileInfo))
	}
}

// DownloadHandler 下载文件
func (t *FileTransport) DownloadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		tenantCode := context.GetTenantCode(c.Request.Context())
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "租户信息缺失"})
			return
		}

		// 获取文件ID
		fileID := c.Param("id")
		if fileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
			return
		}

		// 调用服务下载文件
		reader, fileInfo, err := t.fileService.Download(c.Request.Context(), tenantCode, fileID)
		if err != nil {
			t.logger.Error("下载文件失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "下载文件失败: " + err.Error()})
			return
		}
		defer reader.Close()

		// 设置响应头
		c.Header("Content-Type", fileInfo.FileType)
		c.Header("Content-Disposition", "attachment; filename="+fileInfo.FileName)
		c.Header("Content-Length", strconv.FormatInt(fileInfo.FileSize, 10))

		// 流式响应文件内容
		c.DataFromReader(http.StatusOK, fileInfo.FileSize, fileInfo.FileType, reader, nil)
	}
}

// DeleteHandler 删除文件
func (t *FileTransport) DeleteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		tenantCode := context.GetTenantCode(c.Request.Context())
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "租户信息缺失"})
			return
		}

		// 获取文件ID
		fileID := c.Param("id")
		if fileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
			return
		}

		// 调用服务删除文件
		err := t.fileService.Delete(c.Request.Context(), tenantCode, fileID)
		if err != nil {
			t.logger.Error("删除文件失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败: " + err.Error()})
			return
		}

		response.Success(c, gin.H{"message": "删除成功"})
	}
}

// ListHandler 获取文件列表
func (t *FileTransport) ListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		tenantCode := context.GetTenantCode(c.Request.Context())
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "租户信息缺失"})
			return
		}

		// 解析查询参数
		var req dto.FileListRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
			return
		}

		// 调用服务获取列表
		files, total, err := t.fileService.List(c.Request.Context(), tenantCode, req.Page, req.PageSize, req.BusinessType)
		if err != nil {
			t.logger.Error("获取文件列表失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件列表失败: " + err.Error()})
			return
		}

		// 转换响应
		list := make([]dto.FileItem, 0, len(files))
		for _, file := range files {
			list = append(list, dto.ModelToFileItem(file))
		}

		response.Success(c, dto.FileListResponse{
			List:  list,
			Total: total,
		})
	}
}

// GetPresignedURLHandler 获取预签名URL
func (t *FileTransport) GetPresignedURLHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户信息
		tenantCode := context.GetTenantCode(c.Request.Context())
		if tenantCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "租户信息缺失"})
			return
		}

		// 获取文件ID
		fileID := c.Param("id")
		if fileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
			return
		}

		// 解析查询参数
		var req dto.PresignedURLRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
			return
		}

		// 默认过期时间为3600秒（1小时）
		if req.ExpireSeconds == 0 {
			req.ExpireSeconds = 3600
		}

		// 调用服务获取预签名URL
		url, err := t.fileService.GetPresignedURL(c.Request.Context(), tenantCode, fileID, req.ExpireSeconds)
		if err != nil {
			t.logger.Error("获取预签名URL失败", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取预签名URL失败: " + err.Error()})
			return
		}

		expireAt := time.Now().Add(time.Duration(req.ExpireSeconds) * time.Second)
		response.Success(c, dto.PresignedURLResponse{
			URL:      url,
			ExpireAt: expireAt.Format("2006-01-02 15:04:05"),
		})
	}
}
