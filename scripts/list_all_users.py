#!/usr/bin/env python3
"""
列出所有用户
"""

from pymongo import MongoClient

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

print("="*70)
print("  所有用户列表")
print("="*70)

# 查看所有数据库
dbs = client.list_database_names()
mule_dbs = [db for db in dbs if db.startswith('mule')]

for db_name in sorted(mule_dbs):
    db = client[db_name]
    admin_coll = db['admin']
    
    users = list(admin_coll.find({'is_deleted': 0}))
    
    if users:
        print(f"\n📊 数据库: {db_name}")
        print("-"*70)
        for user in users:
            print(f"  手机号: {user.get('phone')}")
            print(f"  昵称: {user.get('nickname')}")
            print(f"  角色: {user.get('role', user.get('roles', []))}")
            print(f"  状态: {'启用' if user.get('status') == 1 else '禁用'}")
            if 'tenant_id' in user:
                print(f"  ⚠️  tenant_id: {user.get('tenant_id')}")
            print()

print("="*70)
client.close()

