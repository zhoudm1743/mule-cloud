#!/usr/bin/env python3
"""
重置用户密码为 123456
"""

from pymongo import MongoClient
import hashlib

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

# MD5加密（与后端一致）
def hash_password(pwd):
    return hashlib.md5((pwd + "mule-zdm").encode()).hexdigest()

admin_coll = tenant_db['admin']

# 重置密码
new_password = hash_password("123456")
print(f"新密码（加密后）: {new_password}")

result = admin_coll.update_one(
    {'phone': '13838383388'},
    {'$set': {'password': new_password}}
)

if result.modified_count > 0:
    print("✅ 密码已重置为: 123456")
else:
    print("⚠️  密码未修改（可能已经是这个密码）")

# 验证
user = admin_coll.find_one({'phone': '13838383388'})
if user:
    print(f"\n📋 用户信息:")
    print(f"  phone: {user.get('phone')}")
    print(f"  nickname: {user.get('nickname')}")
    print(f"  password: {user.get('password')}")
    print(f"  roles: {user.get('roles')}")
    print(f"  status: {user.get('status')}")

client.close()

