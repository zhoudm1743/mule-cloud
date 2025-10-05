#!/usr/bin/env python3
"""
检查用户的角色信息
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

print("="*60)
print(f"📋 检查数据库: {tenant_db_name}")
print("="*60)

# 查找用户
admin_coll = tenant_db['admin']
user = admin_coll.find_one({'phone': '13838383388'})

if user:
    print(f"\n✅ 找到用户:")
    print(f"  _id: {user.get('_id')}")
    print(f"  phone: {user.get('phone')}")
    print(f"  nickname: {user.get('nickname')}")
    print(f"  roles: {user.get('roles')}")
    
    # 检查角色集合
    role_coll = tenant_db['role']
    print(f"\n📋 角色集合中的数据:")
    all_roles = list(role_coll.find({}))
    if not all_roles:
        print("  ❌ 角色集合为空！")
    else:
        for r in all_roles:
            print(f"  - {r.get('_id')}: {r.get('name')} ({r.get('code')})")
    
    # 查找用户需要的角色
    if user.get('roles'):
        print(f"\n🔍 查找用户的角色:")
        for role_id in user.get('roles'):
            role = role_coll.find_one({'_id': role_id})
            if role:
                print(f"  ✅ {role_id}: {role.get('name')}")
            else:
                print(f"  ❌ {role_id}: 不存在！")
else:
    print(f"\n❌ 未找到用户")

client.close()

