#!/bin/bash

# 认证服务 API 测试脚本
# 使用方法: bash scripts/test_auth_api.sh

BASE_URL="http://localhost:8002"
echo "🚀 开始测试认证服务 API..."
echo "基础URL: $BASE_URL"
echo "================================"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
PASS=0
FAIL=0

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local token=$5
    
    echo ""
    echo -e "${YELLOW}测试: $name${NC}"
    echo "请求: $method $endpoint"
    
    if [ -z "$token" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $token" \
            -d "$data")
    fi
    
    echo "响应: $response"
    
    # 检查是否成功
    if echo "$response" | grep -q '"code":0'; then
        echo -e "${GREEN}✅ 测试通过${NC}"
        ((PASS++))
        echo "$response"
    else
        echo -e "${RED}❌ 测试失败${NC}"
        ((FAIL++))
    fi
    
    echo "--------------------------------"
}

# 1. 测试健康检查
test_api "健康检查" "GET" "/health" "" ""

# 2. 测试用户注册
PHONE="138$(date +%s | tail -c 9)"  # 生成唯一手机号
echo "生成测试手机号: $PHONE"
test_api "用户注册" "POST" "/auth/register" \
    "{\"phone\":\"$PHONE\",\"password\":\"123456\",\"nickname\":\"测试用户\",\"email\":\"test@example.com\"}" \
    ""

# 3. 测试重复注册（应该失败）
test_api "重复注册（预期失败）" "POST" "/auth/register" \
    "{\"phone\":\"$PHONE\",\"password\":\"123456\",\"nickname\":\"测试用户\",\"email\":\"test@example.com\"}" \
    ""

# 4. 测试登录（使用预置的测试账号）
echo ""
echo -e "${YELLOW}使用预置测试账号登录: 13800138000${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"phone":"13800138000","password":"123456"}')

echo "登录响应: $LOGIN_RESPONSE"

if echo "$LOGIN_RESPONSE" | grep -q '"code":0'; then
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
    echo -e "${GREEN}✅ 登录成功${NC}"
    echo "Token: ${TOKEN:0:50}..."
    ((PASS++))
else
    echo -e "${RED}❌ 登录失败，请先运行初始化脚本${NC}"
    echo "运行: mongosh mongodb://root:bgg8384495@localhost:27015/mule?authSource=admin < scripts/init_auth_users.js"
    ((FAIL++))
    TOKEN=""
fi

# 如果登录成功，继续测试需要认证的接口
if [ ! -z "$TOKEN" ]; then
    # 5. 测试获取个人信息
    test_api "获取个人信息" "GET" "/auth/profile" "" "$TOKEN"
    
    # 6. 测试更新个人信息
    test_api "更新个人信息" "PUT" "/auth/profile" \
        "{\"nickname\":\"更新后的昵称\",\"email\":\"newemail@example.com\"}" \
        "$TOKEN"
    
    # 7. 测试刷新Token
    test_api "刷新Token" "POST" "/auth/refresh" \
        "{\"token\":\"$TOKEN\"}" \
        ""
    
    # 8. 测试修改密码（暂时注释，避免影响后续测试）
    # test_api "修改密码" "POST" "/auth/password" \
    #     "{\"old_password\":\"123456\",\"new_password\":\"654321\"}" \
    #     "$TOKEN"
    
    # 9. 测试无效Token（应该失败）
    test_api "无效Token访问（预期失败）" "GET" "/auth/profile" "" "invalid_token_here"
fi

# 10. 测试错误的登录密码
test_api "错误密码登录（预期失败）" "POST" "/auth/login" \
    "{\"phone\":\"13800138000\",\"password\":\"wrongpassword\"}" \
    ""

# 11. 测试不存在的用户
test_api "不存在的用户登录（预期失败）" "POST" "/auth/login" \
    "{\"phone\":\"99999999999\",\"password\":\"123456\"}" \
    ""

# 输出测试统计
echo ""
echo "================================"
echo -e "${GREEN}测试完成！${NC}"
echo -e "通过: ${GREEN}$PASS${NC}"
echo -e "失败: ${RED}$FAIL${NC}"
echo "================================"

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${RED}⚠️  有 $FAIL 个测试失败${NC}"
    exit 1
fi

