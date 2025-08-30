# 应用生命周期管理

本包提供了完整的应用生命周期管理功能，支持优雅启动、关闭和资源清理。

## 功能特性

### 🚀 服务管理
- **服务注册**: 统一注册和管理应用中的各种服务
- **按序启动**: 按照注册顺序依次启动服务
- **失败回滚**: 当某个服务启动失败时，自动回滚已启动的服务
- **优雅关闭**: 按相反顺序优雅关闭所有服务

### 📡 信号处理
- **信号监听**: 监听 SIGINT、SIGTERM、SIGHUP 等系统信号
- **优雅退出**: 接收到关闭信号时，优雅地停止所有服务
- **超时控制**: 每个服务的启动/停止都有超时控制

### ⚕️ 健康监控
- **定期检查**: 每30秒对所有服务进行健康检查
- **状态报告**: 提供服务状态查询接口
- **异常报警**: 当服务健康检查失败时记录警告日志

## 服务接口

每个服务需要实现 `Service` 接口：

```go
type Service interface {
    // Name 返回服务名称
    Name() string
    // Start 启动服务
    Start(ctx context.Context) error
    // Stop 停止服务
    Stop(ctx context.Context) error
    // HealthCheck 健康检查
    HealthCheck(ctx context.Context) error
}
```

## 使用示例

### 1. 创建生命周期管理器
```go
manager := lifecycle.NewManager()
```

### 2. 注册服务
```go
// 按启动顺序注册服务
manager.Register(lifecycle.NewLogService())
manager.Register(lifecycle.NewConfigService(configPath))
manager.Register(lifecycle.NewDatabaseService())
manager.Register(lifecycle.NewHTTPService())
```

### 3. 启动应用
```go
if err := manager.Start(); err != nil {
    log.Fatal("应用启动失败:", err)
}

// 等待应用关闭
manager.Wait()
```

## 内置服务适配器

### ConfigService
- **功能**: 管理配置文件加载和监听
- **启动**: 加载配置文件，开发环境下启动文件监听
- **健康检查**: 验证配置是否已正确初始化

### DatabaseService  
- **功能**: 管理数据库连接
- **启动**: 初始化数据库连接池
- **健康检查**: 验证数据库连接是否正常

### HTTPService
- **功能**: 管理HTTP服务器
- **启动**: 启动HTTP服务器
- **关闭**: 优雅关闭HTTP服务器
- **健康检查**: 检查HTTP服务器是否响应正常

### LogService
- **功能**: 管理日志系统
- **启动**: 初始化日志配置
- **健康检查**: 验证日志器是否可用

## 启动流程

```
📝 LogService    -> 初始化日志系统
📁 ConfigService -> 加载配置文件，启动文件监听
🗄️ DatabaseService -> 初始化数据库连接
🌐 HTTPService   -> 启动HTTP服务器
⚕️ HealthCheck   -> 开始定期健康检查
📡 SignalHandler -> 监听系统信号
```

## 关闭流程

```
📡 接收信号 (SIGINT/SIGTERM/SIGHUP)
🛑 停止HTTP服务器
🗄️ 关闭数据库连接  
📁 停止配置监听
📝 清理日志资源
✅ 应用完全关闭
```

## 错误处理

- **启动失败**: 自动回滚已启动的服务，确保资源正确清理
- **超时处理**: 每个操作都有30秒超时，避免无限阻塞
- **健康检查**: 定期检查服务状态，及时发现问题
- **日志记录**: 详细记录启动、关闭和错误信息

## 配置选项

目前所有超时时间都是硬编码的，未来版本将支持通过配置文件自定义：

- 服务启动超时: 30秒
- 服务停止超时: 30秒  
- 健康检查间隔: 30秒
- 健康检查超时: 10秒

## 扩展服务

要添加新的服务，只需实现 `Service` 接口并注册到管理器：

```go
type MyService struct{}

func (s *MyService) Name() string {
    return "MyService"
}

func (s *MyService) Start(ctx context.Context) error {
    // 启动逻辑
    return nil
}

func (s *MyService) Stop(ctx context.Context) error {
    // 停止逻辑
    return nil
}

func (s *MyService) HealthCheck(ctx context.Context) error {
    // 健康检查逻辑
    return nil
}

// 注册服务
manager.Register(&MyService{})
```

