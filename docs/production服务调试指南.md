# Production 服务调试指南

## 一、修复的问题

### 1. 小程序动态导入错误 ✅
**问题**：`TypeError: r is not a function at scan.js:75`

**原因**：小程序不支持 ES6 动态 `import()`

**修复**：在文件顶部直接导入
```javascript
// ❌ 错误的方式（动态导入）
const { parseScanCode } = await import('@/api/scan')

// ✅ 正确的方式（静态导入）
import { parseScanCode } from '@/api/scan'
```

### 2. 添加租户上下文支持 ✅
在 `request.js` 中添加了 `X-Tenant-Code` header 支持：
```javascript
// 添加租户代码
if (currentTenant && currentTenant.tenant_code) {
    header['X-Tenant-Code'] = currentTenant.tenant_code;
}
```

### 3. 修正响应数据格式 ✅
修改了 `request.js` 的响应处理，返回完整的响应对象：
```javascript
// 之前：resolve(res.data.data) - 只返回 data 字段
// 现在：resolve(res.data) - 返回整个响应对象 { code, data, message }
```

## 二、服务配置

### Production 服务
- **服务名**：production-service
- **端口**：8008（已从 8007 改为 8008）
- **路由前缀**：`/production`
- **健康检查**：`http://localhost:8008/health`

### API 路径
小程序通过网关访问：
```
BASE_URL: https://dev.inzj.cn/api
扫码解析: POST /production/scan/parse
工序上报: POST /production/reports
上报列表: GET /production/reports
订单进度: GET /production/progress/:order_id
工资统计: GET /production/salary
```

## 三、测试步骤

### 1. 启动 Production 服务

```bash
# 方式1：使用启动脚本
.\start.ps1
# 选择 production 服务

# 方式2：单独启动
go run cmd/production/main.go -config=config/production.yaml
```

### 2. 验证服务状态

```bash
# 健康检查
curl http://localhost:8008/health

# 预期响应：
# {
#   "status": "ok",
#   "service": "production-service"
# }
```

### 3. 验证网关路由

```bash
# 通过网关访问健康检查
curl http://localhost:8080/admin/production/health

# 或者查看 Consul 服务注册
curl http://localhost:8500/v1/agent/services
```

### 4. 测试扫码解析 API

```bash
# 模拟扫码请求
curl -X POST http://localhost:8080/admin/production/scan/parse \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Tenant-Code: YOUR_TENANT_CODE" \
  -d '{
    "qr_code": "{\"order_id\":\"xxx\",\"bundle_no\":\"001\",\"color\":\"红色\"}"
  }'
```

## 四、小程序调试

### 1. 检查登录状态

在小程序调试控制台执行：
```javascript
console.log('Token:', uni.getStorageSync('token'))
console.log('Tenant:', uni.getStorageSync('currentTenant'))
```

确保：
- ✅ token 存在且有效
- ✅ currentTenant 对象包含 `tenant_code` 字段

### 2. 测试 API 请求

在 scan.vue 的 `onScanCode` 方法中添加调试日志：
```javascript
try {
    console.log('开始请求:', { qr_code: result })
    const res = await parseScanCode({ qr_code: result })
    console.log('请求成功:', res)
    // ...
} catch (error) {
    console.error('请求失败:', error)
    console.error('错误详情:', error.message)
}
```

### 3. 常见问题排查

#### 问题1：401 未认证
**原因**：token 无效或已过期
**解决**：重新登录

#### 问题2：网络请求失败
**原因**：
- 网关未启动
- production 服务未启动
- 网关未正确转发 production 路由

**解决**：
1. 确认所有服务都在运行
2. 检查 Consul 中 production 服务的注册状态
3. 检查网关路由配置

#### 问题3：二维码格式错误
**原因**：二维码内容不是有效的 JSON 或缺少必要字段

**示例有效二维码**：
```json
{
  "task_id": "批次任务ID",
  "order_id": "订单ID",
  "contract_no": "合同号",
  "style_no": "款号",
  "bed_no": "床号",
  "bundle_no": "扎号",
  "color": "颜色",
  "size": "尺码",
  "quantity": 50
}
```

## 五、完整的小程序请求流程

```
1. 用户扫码
   └─> scan.vue: onScanCode(e)
   
2. 调用 API
   └─> parseScanCode({ qr_code: result })
       └─> @/api/scan.js: parseScanCode()
           └─> @/utils/request.js: post()
   
3. 发送 HTTP 请求
   └─> URL: https://dev.inzj.cn/api/production/scan/parse
       Headers:
         - Authorization: Bearer {token}
         - X-Tenant-Code: {tenant_code}
       Body:
         - qr_code: {二维码内容}
   
4. 网关转发
   └─> Gateway (8080)
       └─> Production Service (8008)
   
5. 后端处理
   └─> cmd/production/main.go
       └─> transport/scan.go: ParseScanCodeHandler()
           └─> endpoint/scan.go: ParseScanCodeEndpoint()
               └─> services/scan.go: ParseScanCode()
   
6. 返回响应
   └─> {
         "code": 200,
         "data": {
           "batch": {...},
           "order": {...},
           "batch_progress": [...]
         },
         "message": "成功"
       }
   
7. 跳转页面
   └─> uni.navigateTo({
         url: '/pages/report/detail?data=...'
       })
```

## 六、注意事项

1. **端口已变更**：从 8007 改为 8008，确保配置文件和启动脚本都已更新
2. **租户隔离**：所有数据操作都基于租户代码隔离
3. **认证要求**：所有 API 都需要 JWT token 认证
4. **网关前缀**：小程序访问需要通过 `/admin/production/` 前缀
5. **响应格式**：统一使用 `{ code, data, message }` 格式

## 七、下一步

- [ ] 测试完整的扫码流程
- [ ] 测试工序上报功能
- [ ] 测试历史记录查询
- [ ] 验证租户数据隔离
- [ ] 性能测试和优化

