#!/usr/bin/env python3
"""
列出所有菜单
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

print("="*60)
print("📋 tenant_system 数据库中的所有菜单:")
print("="*60)

all_menus = list(menu_coll.find({'is_deleted': 0}).sort('order', 1))
parent_menus = [m for m in all_menus if m.get('pid') is None]
child_menus = [m for m in all_menus if m.get('pid') is not None]

print(f"\n总计: {len(all_menus)} 个菜单 ({len(parent_menus)} 个父菜单, {len(child_menus)} 个子菜单)\n")

for parent in parent_menus:
    menu_type = parent.get('menu_type', 'page')
    order = parent.get('order', 0)
    print(f"📁 [{order}] {parent.get('title')} ({parent.get('name')}) - {menu_type}")
    print(f"   Path: {parent.get('path')}")
    
    # 找到所有子菜单
    children = [c for c in child_menus if c.get('pid') == parent['_id']]
    if children:
        for child in sorted(children, key=lambda x: x.get('order', 0)):
            child_type = child.get('menu_type', 'page')
            child_order = child.get('order', 0)
            component = child.get('component_path', 'N/A')
            print(f"   └─ [{child_order}] {child.get('title')} ({child.get('name')}) - {child_type}")
            print(f"      Path: {child.get('path')}")
            print(f"      Component: {component}")
    print()

client.close()

