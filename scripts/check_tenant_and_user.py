#!/usr/bin/env python3
"""
检查租户和用户数据
"""

from pymongo import MongoClient

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
admin_coll = system_db['admin']

print("="*60)
print("📋 系统数据库 (tenant_system) 中的租户:")
print("="*60)

tenants = list(tenant_coll.find({'is_deleted': 0}))
if not tenants:
    print("❌ 没有找到任何租户")
else:
    for t in tenants:
        print(f"\n租户ID: {t['_id']}")
        print(f"  Code: {t.get('code')}")
        print(f"  Name: {t.get('name')}")
        print(f"  Status: {t.get('status')}")
        
        # 检查该租户的数据库
        tenant_db_name = f"tenant_{t['_id']}"
        tenant_db = client[tenant_db_name]
        tenant_admin_coll = tenant_db['admin']
        
        users = list(tenant_admin_coll.find({'is_deleted': 0}))
        print(f"  租户数据库: {tenant_db_name}")
        print(f"  用户数量: {len(users)}")
        
        if users:
            print(f"  用户列表:")
            for u in users:
                print(f"    - {u.get('phone')} ({u.get('nickname')})")

print("\n" + "="*60)
print("📋 系统数据库中的管理员:")
print("="*60)

system_admins = list(admin_coll.find({'is_deleted': 0}))
if not system_admins:
    print("❌ 没有找到任何系统管理员")
else:
    for u in system_admins:
        print(f"  - {u.get('phone')} ({u.get('nickname')})")

client.close()

