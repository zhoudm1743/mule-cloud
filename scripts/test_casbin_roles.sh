#!/bin/bash

# Casbin 角色权限测试脚本

BASE_URL="http://localhost:8080/admin/system"

echo "=== Casbin 角色权限系统测试 ==="
echo ""

# 1. 创建租户
echo "1. 创建测试租户..."
TENANT_RESPONSE=$(curl -s -X POST "$BASE_URL/tenants" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试公司",
    "code": "test_corp",
    "contact": "测试联系人",
    "phone": "13800138000",
    "email": "test@example.com"
  }')
echo $TENANT_RESPONSE | jq '.'
TENANT_ID=$(echo $TENANT_RESPONSE | jq -r '.data.id')
echo "租户ID: $TENANT_ID"
echo ""

# 2. 创建角色
echo "2. 创建测试角色..."
ROLE_RESPONSE=$(curl -s -X POST "$BASE_URL/roles" \
  -H "Content-Type: application/json" \
  -d "{
    \"tenant_id\": \"$TENANT_ID\",
    \"name\": \"系统管理员\",
    \"code\": \"sys_admin\",
    \"description\": \"系统管理员角色\",
    \"menus\": []
  }")
echo $ROLE_RESPONSE | jq '.'
ROLE_ID=$(echo $ROLE_RESPONSE | jq -r '.data.id')
echo "角色ID: $ROLE_ID"
echo ""

# 3. 创建菜单
echo "3. 创建测试菜单..."
MENU1_RESPONSE=$(curl -s -X POST "$BASE_URL/menus" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "dashboard",
    "path": "/dashboard",
    "title": "仪表盘",
    "icon": "icon-park-outline:analysis",
    "requiresAuth": true,
    "menuType": "dir"
  }')
echo $MENU1_RESPONSE | jq '.'
MENU1_ID=$(echo $MENU1_RESPONSE | jq -r '.data.id')
echo "菜单1 ID: $MENU1_ID"
echo ""

MENU2_RESPONSE=$(curl -s -X POST "$BASE_URL/menus" \
  -H "Content-Type: application/json" \
  -d "{
    \"pid\": \"$MENU1_ID\",
    \"name\": \"dashboard_workbench\",
    \"path\": \"/dashboard/workbench\",
    \"title\": \"工作台\",
    \"componentPath\": \"/dashboard/workbench/index.vue\",
    \"icon\": \"icon-park-outline:alarm\",
    \"requiresAuth\": true,
    \"menuType\": \"page\"
  }")
echo $MENU2_RESPONSE | jq '.'
MENU2_ID=$(echo $MENU2_RESPONSE | jq -r '.data.id')
echo "菜单2 ID: $MENU2_ID"
echo ""

# 4. 分配菜单给角色
echo "4. 分配菜单权限给角色..."
ASSIGN_RESPONSE=$(curl -s -X POST "$BASE_URL/roles/$ROLE_ID/menus" \
  -H "Content-Type: application/json" \
  -d "{
    \"menus\": [\"$MENU1_ID\", \"$MENU2_ID\"]
  }")
echo $ASSIGN_RESPONSE | jq '.'
echo ""

# 5. 获取角色的菜单
echo "5. 查询角色的菜单权限..."
curl -s -X GET "$BASE_URL/roles/$ROLE_ID/menus" | jq '.'
echo ""

# 6. 创建管理员
echo "6. 创建测试管理员..."
ADMIN_RESPONSE=$(curl -s -X POST "$BASE_URL/admins" \
  -H "Content-Type: application/json" \
  -d "{
    \"phone\": \"13900139000\",
    \"password\": \"123456\",
    \"nickname\": \"测试管理员\",
    \"email\": \"admin@test.com\",
    \"status\": 1
  }")
echo $ADMIN_RESPONSE | jq '.'
ADMIN_ID=$(echo $ADMIN_RESPONSE | jq -r '.data.id')
echo "管理员ID: $ADMIN_ID"
echo ""

# 7. 分配角色给管理员
echo "7. 分配角色给管理员..."
ASSIGN_ROLE_RESPONSE=$(curl -s -X POST "$BASE_URL/admins/$ADMIN_ID/roles" \
  -H "Content-Type: application/json" \
  -d "{
    \"roles\": [\"$ROLE_ID\"]
  }")
echo $ASSIGN_ROLE_RESPONSE | jq '.'
echo ""

# 8. 获取管理员的角色
echo "8. 查询管理员的角色..."
curl -s -X GET "$BASE_URL/admins/$ADMIN_ID/roles" | jq '.'
echo ""

# 9. 查询角色列表
echo "9. 查询角色列表..."
curl -s -X GET "$BASE_URL/roles?tenant_id=$TENANT_ID&page=1&page_size=10" | jq '.'
echo ""

# 10. 查询菜单列表
echo "10. 查询所有菜单..."
curl -s -X GET "$BASE_URL/menus/all" | jq '.'
echo ""

echo "=== 测试完成 ==="
echo ""
echo "创建的资源ID："
echo "  租户ID: $TENANT_ID"
echo "  角色ID: $ROLE_ID"
echo "  菜单1 ID: $MENU1_ID"
echo "  菜单2 ID: $MENU2_ID"
echo "  管理员ID: $ADMIN_ID"

