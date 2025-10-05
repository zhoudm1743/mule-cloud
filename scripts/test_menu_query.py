#!/usr/bin/env python3
"""
测试菜单查询
"""

from pymongo import MongoClient
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

# 测试ID
test_id = '68ddf8dfc09a1d346c845ced'

print(f"测试菜单ID: {test_id}")
print("="*60)

# 方式1: 使用字符串查询
print("\n1️⃣ 使用字符串查询:")
menu1 = menu_coll.find_one({'_id': test_id, 'is_deleted': 0})
if menu1:
    print(f"   ✅ 找到: {menu1.get('title')}")
    print(f"   _id 类型: {type(menu1['_id'])}")
    print(f"   _id 值: {menu1['_id']}")
else:
    print(f"   ❌ 未找到")

# 方式2: 使用ObjectId查询
print("\n2️⃣ 使用ObjectId查询:")
try:
    obj_id = ObjectId(test_id)
    menu2 = menu_coll.find_one({'_id': obj_id, 'is_deleted': 0})
    if menu2:
        print(f"   ✅ 找到: {menu2.get('title')}")
        print(f"   _id 类型: {type(menu2['_id'])}")
        print(f"   _id 值: {menu2['_id']}")
    else:
        print(f"   ❌ 未找到")
except Exception as e:
    print(f"   ❌ 错误: {e}")

# 查看所有菜单的_id类型
print("\n3️⃣ 所有菜单的_id类型:")
all_menus = list(menu_coll.find({'is_deleted': 0}))
for m in all_menus:
    print(f"   {m.get('name'):20} | _id类型: {type(m['_id']).__name__:15} | _id值: {m['_id']}")

client.close()

