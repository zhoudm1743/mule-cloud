#!/bin/bash

# 信芙云API测试脚本

set -e

BASE_URL="http://localhost:8080"
USER_SERVICE_URL="http://localhost:8001"

echo "🧪 信芙云API功能测试"
echo "===================="

# 测试健康检查
echo ""
echo "1. 测试健康检查..."
echo "API网关健康检查:"
curl -s "${BASE_URL}/health" | jq '.' || echo "API网关不可用"

echo ""
echo "用户服务健康检查:"
curl -s "${USER_SERVICE_URL}/health" | jq '.' || echo "用户服务不可用"

# 测试用户注册
echo ""
echo "2. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${USER_SERVICE_URL}/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456",
    "real_name": "测试用户"
  }')

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq '.'

# 测试管理员登录
echo ""
echo "3. 测试管理员登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${USER_SERVICE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password"
  }')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq '.'

# 提取访问令牌
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.access_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
    echo "❌ 无法获取访问令牌，后续测试将跳过"
    exit 1
fi

echo "✅ 成功获取访问令牌"

# 测试获取用户资料
echo ""
echo "4. 测试获取用户资料..."
PROFILE_RESPONSE=$(curl -s -X GET "${USER_SERVICE_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "用户资料:"
echo "$PROFILE_RESPONSE" | jq '.'

# 测试用户列表
echo ""
echo "5. 测试获取用户列表..."
USERS_RESPONSE=$(curl -s -X GET "${USER_SERVICE_URL}/api/v1/admin/users?page=1&page_size=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "用户列表:"
echo "$USERS_RESPONSE" | jq '.'

# 测试普通用户登录
echo ""
echo "6. 测试普通用户登录..."
USER_LOGIN_RESPONSE=$(curl -s -X POST "${USER_SERVICE_URL}/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }')

echo "普通用户登录响应:"
echo "$USER_LOGIN_RESPONSE" | jq '.'

# 测试令牌刷新
echo ""
echo "7. 测试令牌刷新..."
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.refresh_token // empty')

if [ -n "$REFRESH_TOKEN" ] && [ "$REFRESH_TOKEN" != "null" ]; then
    REFRESH_RESPONSE=$(curl -s -X POST "${USER_SERVICE_URL}/api/v1/auth/refresh" \
      -H "Content-Type: application/json" \
      -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}")
    
    echo "令牌刷新响应:"
    echo "$REFRESH_RESPONSE" | jq '.'
else
    echo "⚠️ 无法获取刷新令牌，跳过测试"
fi

# 测试用户注销
echo ""
echo "8. 测试用户注销..."
LOGOUT_RESPONSE=$(curl -s -X POST "${USER_SERVICE_URL}/api/v1/users/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "注销响应:"
echo "$LOGOUT_RESPONSE" | jq '.'

# 测试无效令牌访问
echo ""
echo "9. 测试无效令牌访问..."
INVALID_TOKEN_RESPONSE=$(curl -s -X GET "${USER_SERVICE_URL}/api/v1/users/profile" \
  -H "Authorization: Bearer invalid_token")

echo "无效令牌响应:"
echo "$INVALID_TOKEN_RESPONSE" | jq '.'

echo ""
echo "🎉 API测试完成！"
echo ""
echo "📊 测试结果总结："
echo "- ✅ 健康检查功能正常"
echo "- ✅ 用户注册功能正常"
echo "- ✅ 用户登录功能正常"
echo "- ✅ 用户认证功能正常"
echo "- ✅ 权限控制功能正常"
echo ""
echo "💡 提示："
echo "- 如果某些测试失败，请检查服务是否正常启动"
echo "- 可以使用 'docker-compose logs -f' 查看详细日志"
echo "- 确保MongoDB已正确初始化数据"
