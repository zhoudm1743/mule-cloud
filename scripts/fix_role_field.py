#!/usr/bin/env python3
"""
修复 role 字段（Go 模型使用的是 role 而不是 roles）
"""

from pymongo import MongoClient

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

# 租户数据库
tenant_id = '68dda6cd04ba0d6c8dda4b7a'
tenant_db_name = f'mule_{tenant_id}'
tenant_db = client[tenant_db_name]

admin_coll = tenant_db['admin']

print("🔧 修复用户 role 字段...")

# 清空 role 字段，删除 roles 字段
result = admin_coll.update_one(
    {'phone': '13838383388'},
    {
        '$set': {'role': []},
        '$unset': {'roles': ''}
    }
)

if result.modified_count > 0:
    print("✅ role 字段已清空")
else:
    print("⚠️  未修改")

# 验证
user = admin_coll.find_one({'phone': '13838383388'})
if user:
    print(f"\n📋 修复后的字段:")
    print(f"  role: {user.get('role')}")
    print(f"  roles: {user.get('roles', '字段已删除')}")
    print(f"  status: {user.get('status')}")
    print(f"  is_deleted: {user.get('is_deleted')}")

client.close()

