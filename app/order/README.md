# 订单服务 (Order Service)

订单管理微服务，支持订单和款式的完整生命周期管理。

## 功能特性

### 订单管理
- **订单列表**：分页查询、多条件搜索（合同号、款号、客户、业务员、订单类型、日期范围等）
- **订单创建**：三步骤创建流程
  - 步骤1：基础信息（合同号、客户、交货日期、订单类型、业务员、备注）
  - 步骤2：款式数量（选择款式、颜色尺码组合、单价、数量）
  - 步骤3：工序清单（配置工序、工价、指定工人）
- **订单编辑**：支持修改订单各项信息
- **订单复制**：快速复制现有订单创建新订单
- **订单详情**：查看完整订单信息
- **订单删除**：软删除订单

### 款式管理
- **款式库**：款式列表查询、分页、搜索
- **款式创建**：添加新款式（款号、款名、颜色、尺码、单价、工序清单、图片）
- **款式编辑**：修改款式信息
- **款式删除**：软删除款式

### 数据模型

#### 订单 (Order)
```go
{
  "id": "订单ID",
  "contract_no": "合同号",
  "style_id": "款式ID",
  "style_no": "款号",
  "style_name": "款名",
  "style_image": "款式图片",
  "customer_id": "客户ID",
  "customer_name": "客户名称",
  "salesman_id": "业务员ID",
  "salesman_name": "业务员名称",
  "order_type_id": "订单类型ID",
  "order_type_name": "订单类型名称",
  "quantity": 450,
  "unit_price": 2.00,
  "total_amount": 900.00,
  "delivery_date": "2022-10-09",
  "progress": 0.0,
  "status": 1,
  "remark": "备注",
  "colors": ["颜色ID列表"],
  "sizes": ["尺码ID列表"],
  "items": [
    {
      "color_id": "颜色ID",
      "color_name": "青蓝",
      "size_id": "尺码ID",
      "size_name": "L",
      "quantity": 100
    }
  ],
  "procedures": [
    {
      "sequence": 1,
      "procedure_id": "工序ID",
      "procedure_name": "质检",
      "unit_price": 0.3,
      "assigned_worker": "指定工人",
      "is_slowest": false,
      "no_bundle": true
    }
  ]
}
```

#### 款式 (Style)
```go
{
  "id": "款式ID",
  "style_no": "0031214",
  "style_name": "测试",
  "category": "分类",
  "season": "季节",
  "year": "2022",
  "images": ["图片URL列表"],
  "colors": ["颜色列表"],
  "sizes": ["尺码列表"],
  "unit_price": 2.00,
  "remark": "备注",
  "procedures": [
    {
      "sequence": 1,
      "procedure_id": "工序ID",
      "procedure_name": "质检",
      "unit_price": 0.3,
      "assigned_worker": "",
      "is_slowest": false,
      "no_bundle": true
    }
  ],
  "status": 1
}
```

## API 接口

### 订单接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /order/orders/:id | 获取单个订单 |
| GET | /order/orders | 订单列表（分页） |
| POST | /order/orders | 创建订单（步骤1：基础信息） |
| PUT | /order/orders/:id/style | 更新订单款式（步骤2：款式数量） |
| PUT | /order/orders/:id/procedure | 更新订单工序（步骤3：工序清单） |
| PUT | /order/orders/:id | 更新订单 |
| POST | /order/orders/:id/copy | 复制订单 |
| DELETE | /order/orders/:id | 删除订单 |

### 款式接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /order/styles/:id | 获取单个款式 |
| GET | /order/styles | 款式列表（分页） |
| GET | /order/styles/all | 获取所有款式（不分页） |
| POST | /order/styles | 创建款式 |
| PUT | /order/styles/:id | 更新款式 |
| DELETE | /order/styles/:id | 删除款式 |

## 启动服务

### 配置文件

配置文件位于 `config/order.yaml`，需要配置：
- MongoDB 连接信息
- Redis 连接信息
- Consul 服务注册信息
- JWT 密钥

### 启动命令

```bash
# 编译
go build -o order.exe cmd/order/main.go

# 运行
./order.exe -config=config/order.yaml
```

### 通过 Consul 访问

服务注册到 Consul 后，可通过网关访问：
- 网关地址：`http://gateway:8080/admin/order/*`
- 直接访问：`http://localhost:8084/order/*`

## 技术栈

- **框架**: Gin + go-kit
- **数据库**: MongoDB（支持多租户数据库隔离）
- **缓存**: Redis
- **服务发现**: Consul
- **认证**: JWT
- **日志**: Zap
- **中间件**: 
  - 租户上下文
  - JWT认证
  - 操作日志
  - 统一响应

## 租户隔离

服务支持完整的租户数据库级别隔离：
- 每个租户拥有独立的数据库
- 通过请求头 `X-Tenant-ID` 或 `X-Tenant-Code` 识别租户
- 自动切换到对应租户的数据库进行操作

## 开发说明

### 目录结构

```
app/order/
├── dto/           # 数据传输对象
│   ├── order.go
│   └── style.go
├── endpoint/      # 端点层（go-kit）
│   ├── order.go
│   └── style.go
├── services/      # 业务逻辑层
│   ├── order.go
│   ├── style.go
│   └── common.go
├── transport/     # HTTP传输层
│   ├── order.go
│   ├── style.go
│   └── common.go
└── README.md

cmd/order/         # 启动文件
└── main.go

internal/
├── models/        # 数据模型
│   └── order.go
└── repository/    # 数据访问层
    ├── order.go
    └── style.go

config/
└── order.yaml     # 配置文件
```

### 添加新功能

1. 在 `internal/models/` 定义数据模型
2. 在 `internal/repository/` 实现数据访问
3. 在 `app/order/dto/` 定义 DTO
4. 在 `app/order/services/` 实现业务逻辑
5. 在 `app/order/endpoint/` 定义端点
6. 在 `app/order/transport/` 实现 HTTP 处理器
7. 在 `cmd/order/main.go` 注册路由

## 未来扩展

- [ ] 裁剪制菲打印功能
- [ ] 订单进度监控
- [ ] 裁片监控
- [ ] 计件工资统计
- [ ] 订单白示（甘特图）
- [ ] 批量导入导出
- [ ] 订单状态流转
- [ ] 消息通知

## 相关文档

- [快速开始](../../docs/快速开始.md)
- [租户隔离架构](../../docs/数据库级别租户隔离-最终完成.md)
- [中间件使用指南](../../docs/中间件极简使用指南.md)
- [操作日志](../../docs/操作日志中间件使用指南.md)
