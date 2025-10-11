# 小程序微服务

## 📋 概述

小程序微服务是专门为微信小程序提供的后端API服务，实现了微信登录、多租户管理、租户切换等核心功能。

## 🌟 核心特性

- ✅ **微信登录**：支持微信小程序登录，自动获取UnionID和OpenID
- ✅ **多租户支持**：一个微信用户可以关联多个租户（工厂）
- ✅ **租户切换**：用户可以在不同租户间自由切换
- ✅ **数据隔离**：每个租户的数据完全隔离在独立的数据库中
- ✅ **JWT认证**：使用JWT进行用户身份验证和授权
- ✅ **灵活绑定**：通过邀请码快速绑定新租户

## 🏗️ 架构设计

### 数据模型

```
系统库 (tenant_system)
  ├─ wechat_user        全局微信用户表
  ├─ user_tenant_map    用户-租户映射表
  └─ tenant             租户表

租户库 (tenant_xxx)
  ├─ member             租户成员表
  └─ ... 其他业务数据
```

### 核心流程

1. **首次登录**：微信授权 → 创建全局用户 → 绑定租户
2. **多租户登录**：微信授权 → 选择租户 → 生成JWT
3. **切换租户**：验证权限 → 生成新JWT → 切换数据库

## 📡 API接口

### 公开接口（无需认证）

#### 1. 微信登录
```http
POST /miniapp/wechat/login
Content-Type: application/json

{
  "code": "微信登录code",
  "encrypted_data": "加密的用户信息（可选）",
  "iv": "加密算法初始向量（可选）"
}
```

**响应示例 - 需要绑定租户**：
```json
{
  "code": 200,
  "data": {
    "need_bind_tenant": true,
    "user_info": {
      "id": "user_id",
      "union_id": "xxx",
      "open_id": "xxx",
      "nickname": "张三",
      "avatar": "头像URL"
    }
  }
}
```

**响应示例 - 需要选择租户**：
```json
{
  "code": 200,
  "data": {
    "need_select_tenant": true,
    "user_info": {...},
    "tenants": [
      {
        "tenant_id": "租户ID",
        "tenant_code": "ace",
        "tenant_name": "工厂A",
        "status": "active"
      }
    ]
  }
}
```

**响应示例 - 直接登录成功**：
```json
{
  "code": 200,
  "data": {
    "token": "JWT_TOKEN",
    "user_info": {...},
    "current_tenant": {
      "tenant_id": "租户ID",
      "tenant_code": "ace",
      "tenant_name": "工厂A",
      "roles": ["employee"]
    }
  }
}
```

#### 2. 绑定租户
```http
POST /miniapp/wechat/bind-tenant
Content-Type: application/json

{
  "user_id": "用户ID",
  "invite_code": "租户邀请码"
}
```

**响应示例**：
```json
{
  "code": 200,
  "data": {
    "success": true,
    "message": "绑定成功",
    "token": "JWT_TOKEN",
    "tenant_info": {
      "tenant_id": "租户ID",
      "tenant_code": "ace",
      "tenant_name": "工厂A"
    }
  }
}
```

#### 3. 选择租户
```http
POST /miniapp/wechat/select-tenant
Content-Type: application/json

{
  "user_id": "用户ID",
  "tenant_id": "租户ID"
}
```

### 认证接口（需要JWT Token）

#### 4. 切换租户
```http
POST /miniapp/wechat/switch-tenant
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "tenant_id": "目标租户ID"
}
```

**响应示例**：
```json
{
  "code": 200,
  "data": {
    "token": "NEW_JWT_TOKEN",
    "user_info": {...},
    "current_tenant": {...}
  }
}
```

#### 5. 获取用户信息
```http
GET /miniapp/user/info
Authorization: Bearer {JWT_TOKEN}
```

**响应示例**：
```json
{
  "code": 200,
  "data": {
    "user_info": {
      "id": "user_id",
      "nickname": "张三",
      "avatar": "头像URL",
      "phone": "13800138000"
    },
    "tenants": [
      {
        "tenant_id": "租户ID1",
        "tenant_name": "工厂A",
        "status": "active"
      },
      {
        "tenant_id": "租户ID2",
        "tenant_name": "工厂B",
        "status": "inactive"
      }
    ]
  }
}
```

#### 6. 更新用户信息
```http
PUT /miniapp/user/info
Authorization: Bearer {JWT_TOKEN}
Content-Type: application/json

{
  "nickname": "新昵称",
  "avatar": "新头像URL",
  "phone": "手机号"
}
```

## 🚀 快速开始

### 1. 配置文件

修改 `config/miniapp.yaml`：

```yaml
wechat:
  app_id: "你的小程序AppID"
  app_secret: "你的小程序AppSecret"
```

### 2. 启动服务

```bash
# 开发环境
go run cmd/miniapp/main.go

# 生产环境
go build -o miniapp cmd/miniapp/main.go
./miniapp
```

### 3. 小程序端调用示例

```javascript
// 1. 微信登录
async function wechatLogin() {
  // 获取微信登录code
  const { code } = await wx.login();
  
  // 调用后端登录接口
  const res = await wx.request({
    url: 'https://your-api.com/admin/miniapp/wechat/login',
    method: 'POST',
    data: { code }
  });
  
  if (res.data.data.need_bind_tenant) {
    // 需要绑定租户，跳转到输入邀请码页面
    wx.navigateTo({ url: '/pages/bind-tenant/index' });
  } else if (res.data.data.need_select_tenant) {
    // 需要选择租户
    wx.navigateTo({ 
      url: '/pages/select-tenant/index',
      data: { tenants: res.data.data.tenants }
    });
  } else {
    // 登录成功
    wx.setStorageSync('token', res.data.data.token);
    wx.switchTab({ url: '/pages/index/index' });
  }
}

// 2. 切换租户
async function switchTenant(tenantId) {
  const token = wx.getStorageSync('token');
  const res = await wx.request({
    url: 'https://your-api.com/admin/miniapp/wechat/switch-tenant',
    method: 'POST',
    header: {
      'Authorization': `Bearer ${token}`
    },
    data: { tenant_id: tenantId }
  });
  
  // 更新Token
  wx.setStorageSync('token', res.data.data.token);
  
  // 重新加载数据
  wx.reLaunch({ url: '/pages/index/index' });
}
```

## 🔒 安全说明

### JWT Token 结构

JWT Token包含以下信息：
- `user_id`：用户ID
- `username`：用户昵称
- `tenant_id`：当前租户ID
- `tenant_code`：当前租户代码
- `roles`：用户在当前租户的角色

### 租户隔离验证

- 每次请求自动验证用户是否有权访问指定租户
- 离职员工（status=inactive）只能查看历史数据，无法修改
- JWT中包含租户信息，后端自动切换对应租户数据库

## 📊 数据库索引建议

### wechat_user表
```javascript
db.wechat_user.createIndex({ "union_id": 1 }, { unique: true, sparse: true })
db.wechat_user.createIndex({ "open_id": 1 }, { unique: true })
db.wechat_user.createIndex({ "phone": 1 }, { unique: true, sparse: true })
db.wechat_user.createIndex({ "tenant_ids": 1 })
```

### user_tenant_map表
```javascript
db.user_tenant_map.createIndex({ "user_id": 1, "tenant_id": 1 }, { unique: true })
db.user_tenant_map.createIndex({ "union_id": 1, "status": 1 })
db.user_tenant_map.createIndex({ "tenant_id": 1, "status": 1 })
```

### member表（在各租户库中）
```javascript
db.member.createIndex({ "union_id": 1 }, { unique: true })
db.member.createIndex({ "user_id": 1 })
db.member.createIndex({ "status": 1, "is_deleted": 1 })
```

## 🐛 常见问题

### 1. UnionID为空？

需要在微信开放平台绑定小程序，否则只能获取到OpenID。如果没有UnionID，系统会自动使用OpenID作为唯一标识。

### 2. 如何生成邀请码？

当前版本邀请码使用租户的`code`字段。建议实现独立的邀请码系统：
- 支持设置有效期
- 支持限制使用次数
- 可以设置默认角色

### 3. 用户离职如何处理？

在租户管理后台将用户状态设置为`inactive`：
- 系统库：`user_tenant_map.status = 'inactive'`
- 租户库：`member.status = 'inactive'`

离职后用户仍可登录，但只能查看历史数据（只读模式）。

## 📚 相关文档

- [微信小程序多租户用户系统设计方案](../../docs/微信小程序多租户用户系统设计方案.md)
- [JWT和gRPC集成指南](../../docs/JWT和gRPC集成指南.md)
- [数据库级别租户隔离改造方案](../../docs/数据库级别租户隔离改造方案.md)

## 📝 更新日志

### v1.0.0 (2025-10-11)
- ✅ 实现微信登录功能
- ✅ 实现多租户管理
- ✅ 实现租户切换
- ✅ 支持一人多租户
- ✅ 完整的数据隔离方案

