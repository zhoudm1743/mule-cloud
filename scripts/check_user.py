#!/usr/bin/env python3
"""
检查用户在哪个数据库中
"""

from pymongo import MongoClient

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

phone = "17858361617"

print(f"查找手机号: {phone}")
print("="*60)

# 查看所有数据库
dbs = client.list_database_names()
mule_dbs = [db for db in dbs if db.startswith('mule')]

for db_name in sorted(mule_dbs):
    db = client[db_name]
    admin_coll = db['admin']
    
    user = admin_coll.find_one({'phone': phone})
    if user:
        print(f"\n✅ 找到用户！")
        print(f"数据库: {db_name}")
        print(f"ID: {user.get('_id')}")
        print(f"手机号: {user.get('phone')}")
        print(f"昵称: {user.get('nickname')}")
        print(f"角色: {user.get('role', user.get('roles', []))}")
        print(f"状态: {user.get('status')}")
        print(f"is_deleted: {user.get('is_deleted')}")
        if 'tenant_id' in user:
            print(f"⚠️  tenant_id: {user.get('tenant_id')} (旧字段，应该删除)")

print("\n" + "="*60)
print("查找完成")

client.close()

