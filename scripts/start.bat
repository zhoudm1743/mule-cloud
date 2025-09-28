@echo off
chcp 65001 >nul
title 信芙云服装生产管理系统

echo 🚀 启动信芙云服装生产管理系统
echo ================================

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 请先安装Docker
    pause
    exit /b 1
)

REM 检查Docker Compose是否安装
docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: 请先安装Docker Compose
    pause
    exit /b 1
)

echo 🔍 检查端口占用情况...
netstat -an | findstr ":8080 " >nul && echo ⚠️  警告: 端口 8080 已被占用，可能会影响 API网关 服务
netstat -an | findstr ":8001 " >nul && echo ⚠️  警告: 端口 8001 已被占用，可能会影响 用户服务 服务
netstat -an | findstr ":27017 " >nul && echo ⚠️  警告: 端口 27017 已被占用，可能会影响 MongoDB 服务
netstat -an | findstr ":6379 " >nul && echo ⚠️  警告: 端口 6379 已被占用，可能会影响 Redis 服务

echo.

REM 构建镜像
echo 🔨 构建Docker镜像...
docker-compose build

echo.

REM 启动服务
echo 🌟 启动所有服务...
docker-compose up -d

echo.

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
echo 📊 检查服务状态...
docker-compose ps

echo.

REM 测试服务连通性
echo 🧪 测试服务连通性...

REM 测试用户服务
echo | set /p="用户服务: "
curl -s http://localhost:8001/health >nul 2>&1
if errorlevel 1 (
    echo ❌ 服务异常
) else (
    echo ✅ 运行正常
)

REM 测试API网关
echo | set /p="API网关: "
curl -s http://localhost:8080/health >nul 2>&1
if errorlevel 1 (
    echo ❌ 服务异常
) else (
    echo ✅ 运行正常
)

echo.
echo 🎉 系统启动完成！
echo.
echo 📚 服务访问地址：
echo   API网关:      http://localhost:8080
echo   用户服务:     http://localhost:8001
echo   Consul UI:    http://localhost:8500
echo   Prometheus:   http://localhost:9090
echo   Grafana:      http://localhost:3000
echo.
echo 🔑 默认登录账号：
echo   管理员:       admin / password
echo   Grafana:      admin / admin123
echo.
echo 🧪 快速测试：
echo   健康检查:     curl http://localhost:8080/health
echo   用户注册:     curl -X POST http://localhost:8001/api/v1/auth/register ^
echo                      -H "Content-Type: application/json" ^
echo                      -d "{\"username\":\"test\",\"email\":\"test@example.com\",\"password\":\"123456\"}"
echo   用户登录:     curl -X POST http://localhost:8001/api/v1/auth/login ^
echo                      -H "Content-Type: application/json" ^
echo                      -d "{\"username\":\"admin\",\"password\":\"password\"}"
echo.
echo 📖 查看日志：
echo   所有服务:     docker-compose logs -f
echo   用户服务:     docker-compose logs -f user-service
echo   API网关:      docker-compose logs -f gateway
echo.
echo ⏹️  停止服务:
echo   docker-compose down
echo.
echo 💡 更多信息请查看 README.md 文件
echo.
echo 按任意键关闭窗口...
pause >nul
