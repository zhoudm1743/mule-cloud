# -*- coding: utf-8 -*-
"""
修复租户权限：移除租户管理员不应该拥有的权限
- 移除 system_menu (菜单管理)
- 移除 system_tenant (租户管理)
"""
from pymongo import MongoClient
from datetime import datetime

# 连接到 MongoDB
client = MongoClient('mongodb://localhost:27017/')

# 租户数据库
tenant_db = 'mule_68dda6cd04ba0d6c8dda4b7a'
db = client[tenant_db]

print(f'=== Fixing permissions for database: {tenant_db} ===\n')

# 查找租户管理员角色
role = db.role.find_one({'code': 'tenant_admin', 'is_deleted': 0})

if not role:
    print('ERROR: Tenant admin role not found!')
    client.close()
    exit(1)

print(f'Found role: {role["name"]}')
print(f'  Current menus: {role.get("menus", [])}')
print()

# 正确的菜单列表（移除系统级别资源）
correct_menus = [
    'dashboard',      # 仪表盘
    'system',         # 系统管理（父菜单）
    'system_admin',   # 管理员管理
    'system_role',    # 角色管理
]

# 正确的菜单权限（移除系统级别资源）
correct_menu_permissions = {
    'system_admin': ['read', 'create', 'update', 'delete'],
    'system_role': ['read', 'create', 'update', 'delete', 'menus'],
}

# 检查是否需要更新
needs_update = False
current_menus = role.get('menus', [])
current_permissions = role.get('menu_permissions', {})

if 'system_menu' in current_menus or 'system_tenant' in current_menus:
    needs_update = True
    print('ISSUE: Found system-level menus that should be removed:')
    if 'system_menu' in current_menus:
        print('  - system_menu (Menu Management)')
    if 'system_tenant' in current_menus:
        print('  - system_tenant (Tenant Management)')
    print()

if 'system_menu' in current_permissions or 'system_tenant' in current_permissions:
    needs_update = True
    print('ISSUE: Found system-level permissions that should be removed')
    print()

if needs_update:
    # 更新角色
    result = db.role.update_one(
        {'_id': role['_id']},
        {
            '$set': {
                'menus': correct_menus,
                'menu_permissions': correct_menu_permissions,
                'updated_at': datetime.now()
            }
        }
    )
    
    if result.modified_count > 0:
        print('SUCCESS: Role permissions updated!')
        print(f'  New menus: {correct_menus}')
        print(f'  New permissions: {list(correct_menu_permissions.keys())}')
    else:
        print('WARNING: No changes made')
else:
    print('OK: Role permissions are already correct!')
    print('  No update needed')

client.close()

