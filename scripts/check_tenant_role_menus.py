# -*- coding: utf-8 -*-
from pymongo import MongoClient
import json

# 连接到 MongoDB
client = MongoClient('mongodb://localhost:27017/')

# 查询租户数据库
tenant_db = 'mule_68dda6cd04ba0d6c8dda4b7a'
db = client[tenant_db]

print(f'=== Database: {tenant_db} ===\n')

# 查询租户的默认角色
role = db.role.find_one({'code': 'tenant_admin', 'is_deleted': 0})

if role:
    print('Tenant Admin Role:')
    print(f"  ID: {role['_id']}")
    print(f"  Name: {role.get('name', 'N/A')}")
    print(f"  Code: {role.get('code', 'N/A')}")
    print('\nAssigned Menus:')
    
    menus = role.get('menus', [])
    if menus:
        for i, menu in enumerate(menus, 1):
            print(f"  {i}. {menu}")
    else:
        print('  (None)')
        
    print('\nMenu Permissions:')
    menu_permissions = role.get('menu_permissions', {})
    if menu_permissions:
        for menu_name, perms in menu_permissions.items():
            print(f"  {menu_name}: {perms}")
    else:
        print('  (None)')
else:
    print('ERROR: Tenant admin role not found')

client.close()

