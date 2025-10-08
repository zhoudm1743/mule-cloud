# Common Service - 文件上传服务

## 概述

Common Service 提供通用的文件上传、下载、管理功能，支持多种存储后端（本地存储、MinIO、AWS S3、阿里云OSS）。

## 功能特性

- ✅ 多种存储后端支持
  - 本地文件系统存储
  - MinIO对象存储
  - AWS S3对象存储
  - 阿里云OSS对象存储

- ✅ 文件管理
  - 文件上传（支持大小和类型验证）
  - 文件下载
  - 文件删除
  - 文件列表查询（分页）
  - 预签名URL生成（临时访问）

- ✅ 多租户支持
  - 租户级别数据隔离
  - 文件按租户、业务类型、日期自动归类

- ✅ 安全性
  - 文件大小限制（默认100MB）
  - 文件类型白名单验证
  - 支持预签名URL实现临时访问

## API接口

### 上传文件
```
POST /admin/common/files/upload
Content-Type: multipart/form-data

参数:
- file: 上传的文件
- business_type: 业务类型（如：avatar, document, image等）

响应:
{
  "id": "文件ID",
  "file_name": "原始文件名",
  "file_size": 1024,
  "file_type": "image/jpeg",
  "file_ext": ".jpg",
  "url": "访问URL",
  "business_type": "avatar",
  "created_at": "2025-10-06 12:00:00"
}
```

### 获取文件列表
```
GET /admin/common/files?page=1&page_size=10&business_type=avatar

响应:
{
  "list": [...],
  "total": 100
}
```

### 下载文件
```
GET /admin/common/files/:id
```

### 删除文件
```
DELETE /admin/common/files/:id
```

### 获取预签名URL
```
GET /admin/common/files/:id/presigned?expire_seconds=3600

响应:
{
  "url": "预签名URL",
  "expire_at": "2025-10-06 13:00:00"
}
```

## 配置说明

### 本地存储（开发环境推荐）
```yaml
storage:
  type: "local"
  local:
    upload_dir: "./uploads"
    base_url: "http://localhost:8004/files"
```

### MinIO存储
```yaml
storage:
  type: "minio"
  minio:
    endpoint: "localhost:9000"
    access_key_id: "minioadmin"
    secret_access_key: "minioadmin"
    bucket_name: "mule-cloud"
    use_ssl: false
    region: "us-east-1"
```

### AWS S3存储
```yaml
storage:
  type: "s3"
  s3:
    endpoint: ""  # 留空使用AWS S3
    access_key_id: "your-access-key"
    secret_access_key: "your-secret-key"
    bucket_name: "mule-cloud"
    region: "us-east-1"
```

### 阿里云OSS存储
```yaml
storage:
  type: "oss"
  oss:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    access_key_id: "your-access-key"
    access_key_secret: "your-secret-key"
    bucket_name: "mule-cloud"
```

## 启动服务

```bash
go run cmd/common/main.go
```

或指定配置文件：
```bash
go run cmd/common/main.go --config config/common.yaml
```

## 文件组织结构

上传的文件按以下结构组织：
```
{tenant_code}/{business_type}/{yyyy}/{mm}/{dd}/{uuid}.{ext}
```

例如：
```
tenant001/avatar/2025/10/06/abc123-def456.jpg
tenant001/document/2025/10/06/xyz789-uvw012.pdf
```

## 支持的文件类型

### 图片
- .jpg, .jpeg, .png, .gif, .bmp, .webp

### 文档
- .pdf, .doc, .docx, .xls, .xlsx, .ppt, .pptx

### 其他
- .txt, .csv, .zip, .rar, .7z

可在 `app/common/services/file.go` 中的 `allowedExts` 变量配置。

## 注意事项

1. **文件大小限制**: 默认100MB，可在 `app/common/services/file.go` 中修改 `maxFileSize` 常量
2. **租户隔离**: 所有API调用需要在请求头中包含租户信息
3. **存储切换**: 切换存储类型后，历史文件需要手动迁移
4. **预签名URL**: 仅对象存储（MinIO/S3/OSS）支持，本地存储直接返回永久URL

## 前端使用

### 导入API
```typescript
import { uploadFile, fetchFileList, deleteFile } from '@/service/api/common'
```

### 上传文件
```typescript
const file = /* File对象 */
const response = await uploadFile(file, 'avatar')
console.log(response.data.url)
```

### 使用上传组件
```vue
<template>
  <FileUpload
    business-type="avatar"
    :max-size="5"
    list-type="image-card"
    @success="handleUploadSuccess"
  />
</template>

<script setup lang="ts">
import FileUpload from '@/components/common/FileUpload.vue'

function handleUploadSuccess(fileInfo: Api.Common.UploadResponse) {
  console.log('上传成功:', fileInfo)
}
</script>
```

## 扩展

### 添加新的存储后端

1. 在 `core/storage/` 目录创建新的存储实现文件
2. 实现 `Storage` 接口的所有方法
3. 在 `core/storage/factory.go` 的 `NewStorage` 函数中添加对应的case
4. 在 `core/config/config.go` 的 `StorageConfig` 中添加配置结构

### 添加新的文件类型

在 `app/common/services/file.go` 的 `Upload` 函数中修改 `allowedExts` map：
```go
allowedExts := map[string]bool{
    // 添加新类型
    ".mp4": true,
    ".avi": true,
    // ...
}
```

## 服务端口

- Common Service: 8004
- 健康检查: http://localhost:8004/health
- 静态文件（本地存储）: http://localhost:8004/files/*
