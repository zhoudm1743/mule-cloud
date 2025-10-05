#!/usr/bin/env python3
"""
创建 default 租户记录
"""

from pymongo import MongoClient
import time

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

# 系统库
system_db = client['tenant_system']
tenant_coll = system_db['tenant']

tenant_id = '68dda6cd04ba0d6c8dda4b7a'

# 检查是否已存在
existing = tenant_coll.find_one({'_id': tenant_id})

if existing:
    print(f"✅ 租户已存在: {existing.get('name')} ({existing.get('code')})")
else:
    # 创建租户记录
    now = int(time.time())
    tenant = {
        '_id': tenant_id,
        'code': 'default',
        'name': '默认租户',
        'contact': '管理员',
        'phone': '13838383388',
        'email': '',
        'menus': [],
        'status': 1,
        'is_deleted': 0,
        'created_at': now,
        'updated_at': now,
    }
    
    tenant_coll.insert_one(tenant)
    print(f"✅ 租户创建成功!")
    print(f"   ID: {tenant_id}")
    print(f"   Code: default")
    print(f"   Name: 默认租户")

# 验证
print("\n📋 所有租户:")
all_tenants = list(tenant_coll.find({'is_deleted': 0}))
for t in all_tenants:
    print(f"  - {t['_id']}: {t.get('name')} ({t.get('code')}) - status: {t.get('status')}")

client.close()

