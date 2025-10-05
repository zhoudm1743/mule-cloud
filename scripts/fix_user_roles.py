#!/usr/bin/env python3
"""
修复用户角色问题
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

# 方案：直接清空用户的角色数组
print("🔧 修复用户角色...")
result = admin_coll.update_one(
    {'phone': '13838383388'},
    {'$set': {'roles': []}}
)

if result.modified_count > 0:
    print("✅ 用户角色已清空")
else:
    print("⚠️  用户未修改（可能已经是空数组）")

# 验证
user = admin_coll.find_one({'phone': '13838383388'})
if user:
    print(f"\n验证用户信息:")
    print(f"  phone: {user.get('phone')}")
    print(f"  nickname: {user.get('nickname')}")
    print(f"  roles: {user.get('roles')}")
    print(f"  status: {user.get('status')}")

client.close()

