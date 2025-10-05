#!/usr/bin/env python3
"""
检查菜单ID
"""

from pymongo import MongoClient

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

# 系统库
system_db = client['tenant_system']
menu_coll = system_db['menu']

# 查找特定ID
target_id = '68ddf8dfc09a1d346c845ced'
menu = menu_coll.find_one({'_id': target_id})

print(f"查找菜单 ID: {target_id}")
if menu:
    print(f"✅ 找到: {menu.get('title')} ({menu.get('name')})")
else:
    print(f"❌ 未找到")

print("\n" + "="*60)
print("📋 所有菜单的ID和名称:")
print("="*60)

all_menus = list(menu_coll.find({'is_deleted': 0}))
for m in all_menus:
    print(f"ID: {m['_id']}")
    print(f"   Name: {m.get('name')}")
    print(f"   Title: {m.get('title')}")
    print(f"   Path: {m.get('path')}")
    print()

client.close()

