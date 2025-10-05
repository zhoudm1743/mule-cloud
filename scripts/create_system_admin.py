#!/usr/bin/env python3
"""
创建系统管理员
"""

from pymongo import MongoClient
import hashlib
import time
from bson import ObjectId

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

# MD5加密（与Go后端一致）
def md5_password(pwd):
    return hashlib.md5((pwd + "mule-zdm").encode()).hexdigest()

# 系统库
system_db = client['mule']
admin_coll = system_db['admin']

# 检查是否已存在
phone = "17858361617"
existing = admin_coll.find_one({'phone': phone})

if existing:
    print(f"❌ 用户 {phone} 已存在")
else:
    # 创建系统管理员
    now = int(time.time())
    admin = {
        '_id': str(ObjectId()),
        'phone': phone,
        'password': md5_password('123456'),
        'nickname': '系统管理员',
        'email': '',
        'avatar': '',
        'role': ['super'],  # 超级管理员角色
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    }
    
    admin_coll.insert_one(admin)
    print(f"✅ 系统管理员创建成功！")
    print(f"   手机号: {phone}")
    print(f"   密码: 123456")
    print(f"   角色: super")
    print(f"   数据库: mule (系统库)")

# 同时创建系统角色
role_coll = system_db['role']
super_role = role_coll.find_one({'code': 'super'})

if not super_role:
    role = {
        '_id': str(ObjectId()),
        'code': 'super',
        'name': '超级管理员',
        'description': '系统超级管理员',
        'menus': [],
        'menu_permissions': {},
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
        'created_by': 'system',
        'updated_by': 'system',
    }
    role_coll.insert_one(role)
    print(f"✅ 超级管理员角色创建成功！")

client.close()

