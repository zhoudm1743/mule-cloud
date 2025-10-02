#!/bin/bash

# 初始化测试菜单数据
echo "========== 创建测试菜单 =========="

# 1. 仪表盘目录
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "dashboard",
    "path": "/dashboard",
    "title": "仪表盘",
    "requiresAuth": true,
    "icon": "icon-park-outline:analysis",
    "menuType": "dir",
    "order": 1
  }'
echo ""

# 2. 工作台页面
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "dashboard_workbench",
    "path": "/dashboard/workbench",
    "title": "工作台",
    "componentPath": "/dashboard/workbench/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:alarm",
    "pinTab": true,
    "menuType": "page",
    "order": 2
  }'
echo ""

# 3. 监控页
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "dashboard_monitor",
    "path": "/dashboard/monitor",
    "title": "监控页",
    "componentPath": "/dashboard/monitor/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:anchor",
    "menuType": "page",
    "order": 3
  }'
echo ""

# 4. 系统设置目录
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "setting",
    "path": "/setting",
    "title": "系统设置",
    "requiresAuth": true,
    "icon": "icon-park-outline:setting-two",
    "menuType": "dir",
    "order": 10
  }'
echo ""

# 5. 菜单管理
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "setting_menu",
    "path": "/setting/menu",
    "title": "菜单管理",
    "componentPath": "/setting/menu/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:application-menu",
    "menuType": "page",
    "order": 11
  }'
echo ""

# 6. 租户管理
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "setting_tenant",
    "path": "/setting/tenant",
    "title": "租户管理",
    "componentPath": "/setting/tenant/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:user",
    "menuType": "page",
    "order": 12
  }'
echo ""

# 7. 角色管理
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "setting_role",
    "path": "/setting/role",
    "title": "角色管理",
    "componentPath": "/setting/role/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:user",
    "menuType": "page",
    "order": 13
  }'
echo ""

# 8. 管理员管理
curl -X POST http://localhost:8080/admin/system/menus \
  -H "Content-Type: application/json" \
  -d '{
    "name": "setting_account",
    "path": "/setting/account",
    "title": "管理员管理",
    "componentPath": "/setting/account/index.vue",
    "requiresAuth": true,
    "icon": "icon-park-outline:user",
    "menuType": "page",
    "order": 14
  }'
echo ""

echo "========== 测试菜单创建完成 =========="

