#!/usr/bin/env python3
"""
检查特定租户数据库
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
tenant_db_name = f'tenant_{tenant_id}'

print(f"="*60)
print(f"📋 检查租户数据库: {tenant_db_name}")
print(f"="*60)

tenant_db = client[tenant_db_name]

# 检查集合
collections = tenant_db.list_collection_names()
print(f"\n集合列表: {collections}")

# 检查 admin 集合
if 'admin' in collections:
    admin_coll = tenant_db['admin']
    
    # 所有用户（包括删除的）
    all_users = list(admin_coll.find({}))
    print(f"\n所有用户记录数（包括已删除）: {len(all_users)}")
    
    for u in all_users:
        print(f"\n用户:")
        print(f"  _id: {u.get('_id')}")
        print(f"  phone: {u.get('phone')}")
        print(f"  nickname: {u.get('nickname')}")
        print(f"  is_deleted: {u.get('is_deleted')}")
        print(f"  status: {u.get('status')}")
        print(f"  roles: {u.get('roles')}")
    
    # 未删除的用户
    active_users = list(admin_coll.find({'is_deleted': 0}))
    print(f"\n未删除的用户数: {len(active_users)}")
    
    # 查找手机号 13838383388
    user_13838 = admin_coll.find_one({'phone': '13838383388'})
    if user_13838:
        print(f"\n✅ 找到用户 13838383388:")
        for key, value in user_13838.items():
            print(f"  {key}: {value}")
    else:
        print(f"\n❌ 未找到用户 13838383388")
else:
    print(f"\n❌ admin 集合不存在")

# 再查一下系统库中的租户信息
system_db = client['tenant_system']
tenant_coll = system_db['tenant']

tenant = tenant_coll.find_one({'_id': tenant_id})
if tenant:
    print(f"\n" + "="*60)
    print(f"📋 租户信息 (tenant_system.tenant):")
    print(f"="*60)
    for key, value in tenant.items():
        print(f"  {key}: {value}")
else:
    print(f"\n❌ 系统库中未找到租户 ID: {tenant_id}")

client.close()

