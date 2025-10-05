# -*- coding: utf-8 -*-
from pymongo import MongoClient
from datetime import datetime
import hashlib

# 连接到 MongoDB
client = MongoClient('mongodb://localhost:27017/')

# 租户数据库
tenant_db = 'mule_68dda6cd04ba0d6c8dda4b7a'
db = client[tenant_db]

# 密码加密函数（与Go代码一致）
def hash_password(password):
    return hashlib.md5(password.encode()).hexdigest()

# 查询租户管理员角色
role = db.role.find_one({'code': 'tenant_admin', 'is_deleted': 0})

if not role:
    print('ERROR: Tenant admin role not found!')
    client.close()
    exit(1)

print(f'Found role: {role["name"]} (ID: {role["_id"]})')

# 检查用户是否已存在
existing_user = db.admin.find_one({'phone': '13838383388'})
if existing_user:
    print('User 13838383388 already exists!')
    print(f'  ID: {existing_user["_id"]}')
    print(f'  Nickname: {existing_user.get("nickname", "N/A")}')
    print(f'  Roles: {existing_user.get("roles", [])}')
else:
    # 创建新用户
    new_user = {
        'phone': '13838383388',
        'password': hash_password('123456'),
        'nickname': '租户管理员',
        'avatar': '',
        'email': '',
        'roles': [str(role['_id'])],
        'status': 1,
        'is_deleted': 0,
        'created_at': datetime.now(),
        'updated_at': datetime.now()
    }
    
    result = db.admin.insert_one(new_user)
    print('User created successfully!')
    print(f'  ID: {result.inserted_id}')
    print(f'  Phone: 13838383388')
    print(f'  Password: 123456')
    print(f'  Nickname: 租户管理员')
    print(f'  Role: {role["name"]}')

client.close()

