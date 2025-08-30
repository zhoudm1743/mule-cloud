# Swagger API 文档

## 简介

本项目已集成 Swagger API 文档，可以方便地查看和测试 API 接口。

## 访问方式

启动服务后，访问以下地址查看 API 文档：

```
http://localhost:8080/swagger/index.html
```

## API 接口概览

### 系统接口
- `GET /health` - 健康检查
- `GET /status` - 状态检查

### 认证接口
- `POST /api/admin/auth/login` - 管理员登录

### 管理员管理
- `GET /api/admin/system/admin` - 获取管理员列表
- `POST /api/admin/system/admin` - 创建管理员
- `PUT /api/admin/system/admin` - 更新管理员
- `DELETE /api/admin/system/admin` - 删除管理员

## 认证方式

需要认证的接口使用 JWT Token，在请求头中添加：

```
Authorization: Bearer <your-jwt-token>
```

## 获取 Token

通过登录接口获取 Token：

```bash
curl -X POST http://localhost:8080/api/admin/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "your-phone",
    "password": "your-password"
  }'
```

## 启动服务

1. **下载依赖包**：
   ```bash
   cd server
   go mod tidy
   go mod download
   ```

2. **启动服务**：
   ```bash
   go run main.go
   ```

3. **访问Swagger文档**：
   ```
   http://localhost:8080/swagger/index.html
   ```

## 重新生成文档

如果修改了 API 注释，可以重新生成文档：

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 在 server 目录下生成文档
swag init -g main.go --output ./docs
```

## 注意事项

1. 确保已安装所有依赖包：`go mod tidy`
2. 确保配置文件 `configs/app.yaml` 存在且配置正确
3. 如果遇到导入问题，请检查 go.mod 文件中的依赖项
4. 已修复 `swaggerFiles` 导入错误，现在使用 `"github.com/swaggo/files"`
