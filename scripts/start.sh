#!/bin/bash

# 信芙云服装生产管理系统启动脚本

set -e

echo "🚀 启动信芙云服装生产管理系统"
echo "================================"

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: 请先安装Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ 错误: 请先安装Docker Compose"
    exit 1
fi

# 检查端口是否被占用
check_port() {
    local port=$1
    local service=$2
    if netstat -tuln 2>/dev/null | grep ":$port " > /dev/null; then
        echo "⚠️  警告: 端口 $port 已被占用，可能会影响 $service 服务"
    fi
}

echo "🔍 检查端口占用情况..."
check_port 8080 "API网关"
check_port 8001 "用户服务"
check_port 27017 "MongoDB"
check_port 6379 "Redis"
check_port 8500 "Consul"
check_port 4222 "NATS"
check_port 9090 "Prometheus"
check_port 3000 "Grafana"

echo ""

# 构建镜像
echo "🔨 构建Docker镜像..."
docker-compose build

echo ""

# 启动服务
echo "🌟 启动所有服务..."
docker-compose up -d

echo ""

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "📊 检查服务状态..."
docker-compose ps

echo ""

# 测试服务连通性
echo "🧪 测试服务连通性..."

# 测试数据库连接
echo -n "MongoDB: "
if docker-compose exec -T mongodb mongosh --quiet --eval "db.runCommand('ping').ok" > /dev/null 2>&1; then
    echo "✅ 连接正常"
else
    echo "❌ 连接失败"
fi

# 测试Redis连接
echo -n "Redis: "
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ 连接正常"
else
    echo "❌ 连接失败"
fi

# 测试用户服务
echo -n "用户服务: "
if curl -s http://localhost:8001/health > /dev/null 2>&1; then
    echo "✅ 运行正常"
else
    echo "❌ 服务异常"
fi

# 测试API网关
echo -n "API网关: "
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ 运行正常"
else
    echo "❌ 服务异常"
fi

echo ""
echo "🎉 系统启动完成！"
echo ""
echo "📚 服务访问地址："
echo "  API网关:      http://localhost:8080"
echo "  用户服务:     http://localhost:8001"
echo "  Consul UI:    http://localhost:8500"
echo "  Prometheus:   http://localhost:9090"
echo "  Grafana:      http://localhost:3000"
echo ""
echo "🔑 默认登录账号："
echo "  管理员:       admin / password"
echo "  Grafana:      admin / admin123"
echo ""
echo "🧪 快速测试："
echo "  健康检查:     curl http://localhost:8080/health"
echo "  用户注册:     curl -X POST http://localhost:8001/api/v1/auth/register \\"
echo "                     -H 'Content-Type: application/json' \\"
echo "                     -d '{\"username\":\"test\",\"email\":\"test@example.com\",\"password\":\"123456\"}'"
echo "  用户登录:     curl -X POST http://localhost:8001/api/v1/auth/login \\"
echo "                     -H 'Content-Type: application/json' \\"
echo "                     -d '{\"username\":\"admin\",\"password\":\"password\"}'"
echo ""
echo "📖 查看日志："
echo "  所有服务:     docker-compose logs -f"
echo "  用户服务:     docker-compose logs -f user-service"
echo "  API网关:      docker-compose logs -f gateway"
echo ""
echo "⏹️  停止服务:"
echo "  docker-compose down"
echo ""
echo "💡 更多信息请查看 README.md 文件"
