#!/usr/bin/env python3
"""
添加菜单到tenant_system数据库
"""

from pymongo import MongoClient
import time
from bson import ObjectId

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

# 检查是否已存在
existing = menu_coll.find_one({'name': 'dashboard'})

if existing:
    print(f"❌ 菜单 dashboard 已存在")
else:
    # 创建菜单
    now = int(time.time())
    menu = {
        '_id': str(ObjectId()),
        'name': 'dashboard',
        'path': '/dashboard',
        'title': '仪表盘',
        'pid': None,  # 顶级菜单
        'icon': 'icon-park-outline:analysis',
        'order': 1,
        'menu_type': 'page',
        'component_path': '/dashboard/workbench/index.vue',
        'keep_alive': True,
        'constant': False,
        'hide': False,
        'href': None,
        'multi_tab': False,
        'active_menu': None,
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    }
    
    menu_coll.insert_one(menu)
    print(f"✅ 菜单创建成功！")
    print(f"   name: dashboard")
    print(f"   title: 仪表盘")
    print(f"   path: /dashboard")
    print(f"   数据库: tenant_system.menu")

# 查看所有菜单
print("\n📋 当前所有菜单:")
all_menus = list(menu_coll.find({'is_deleted': 0}))
for m in all_menus:
    print(f"  - {m.get('name')}: {m.get('title')} ({m.get('path')})")

client.close()

