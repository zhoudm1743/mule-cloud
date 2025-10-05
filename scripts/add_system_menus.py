#!/usr/bin/env python3
"""
添加系统管理菜单到tenant_system数据库
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

now = int(time.time())

# 1. 先创建系统管理父菜单（目录类型）
system_parent_id = str(ObjectId())
system_parent = {
    '_id': system_parent_id,
    'name': 'system',
    'path': '/system',
    'title': '系统管理',
    'pid': None,  # 顶级菜单
    'icon': 'icon-park-outline:setting',
    'order': 2,
    'menuType': 'directory',
    'componentPath': None,
    'keepAlive': True,
    'constant': False,
    'hide': False,
    'href': None,
    'multiTab': False,
    'activeMenu': None,
    'status': 1,
    'is_deleted': 0,
    'created_at': now,
    'updated_at': now,
    'created_by': 'system',
    'updated_by': 'system',
}

# 2. 子菜单
sub_menus = [
    {
        '_id': str(ObjectId()),
        'name': 'system_admin',
        'path': '/system/admin',
        'title': '管理员管理',
        'pid': system_parent_id,
        'icon': 'icon-park-outline:user',
        'order': 1,
        'menuType': 'page',
        'componentPath': '/system/admin/index.vue',
        'keepAlive': True,
        'constant': False,
        'hide': False,
        'href': None,
        'multiTab': False,
        'activeMenu': None,
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    },
    {
        '_id': str(ObjectId()),
        'name': 'system_role',
        'path': '/system/role',
        'title': '角色管理',
        'pid': system_parent_id,
        'icon': 'icon-park-outline:peoples',
        'order': 2,
        'menuType': 'page',
        'componentPath': '/system/role/index.vue',
        'keepAlive': True,
        'constant': False,
        'hide': False,
        'href': None,
        'multiTab': False,
        'activeMenu': None,
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    },
    {
        '_id': str(ObjectId()),
        'name': 'system_menu',
        'path': '/system/menu',
        'title': '菜单管理',
        'pid': system_parent_id,
        'icon': 'icon-park-outline:list',
        'order': 3,
        'menuType': 'page',
        'componentPath': '/system/menu/index.vue',
        'keepAlive': True,
        'constant': False,
        'hide': False,
        'href': None,
        'multiTab': False,
        'activeMenu': None,
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    },
    {
        '_id': str(ObjectId()),
        'name': 'system_tenant',
        'path': '/system/tenant',
        'title': '租户管理',
        'pid': system_parent_id,
        'icon': 'icon-park-outline:database',
        'order': 4,
        'menuType': 'page',
        'componentPath': '/system/tenant/index.vue',
        'keepAlive': True,
        'constant': False,
        'hide': False,
        'href': None,
        'multiTab': False,
        'activeMenu': None,
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    },
]

# 检查并插入
all_menus = [system_parent] + sub_menus

for menu in all_menus:
    existing = menu_coll.find_one({'name': menu['name']})
    if existing:
        print(f"⚠️  菜单 {menu['name']} 已存在，跳过")
    else:
        menu_coll.insert_one(menu)
        menu_type = menu['menuType']
        print(f"✅ 创建菜单: {menu['title']} ({menu['name']}) - {menu_type}")

print("\n" + "="*60)
print("📋 当前所有菜单:")
print("="*60)

all_menus = list(menu_coll.find({'is_deleted': 0}).sort('order', 1))
parent_menus = [m for m in all_menus if m.get('pid') is None]
child_menus = [m for m in all_menus if m.get('pid') is not None]

for parent in parent_menus:
    menu_type = parent.get('menuType', 'page')
    print(f"\n📁 {parent.get('title')} ({parent.get('path')}) - {menu_type}")
    
    # 找到所有子菜单
    children = [c for c in child_menus if c.get('pid') == parent['_id']]
    for child in sorted(children, key=lambda x: x.get('order', 0)):
        child_type = child.get('menuType', 'page')
        print(f"  └─ {child.get('title')} ({child.get('path')}) - {child_type}")

client.close()

