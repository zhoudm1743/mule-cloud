#!/usr/bin/env python3
"""
检查 mule_ 开头的数据库
"""

from pymongo import MongoClient

# 连接配置
host = 'localhost'
port = 27015
username = 'root'
password = 'bgg8384495'

connection_string = f"mongodb://{username}:{password}@{host}:{port}/"
client = MongoClient(connection_string)

print("="*60)
print("📋 所有以 mule_ 开头的数据库:")
print("="*60)

db_names = client.list_database_names()
mule_dbs = [db for db in db_names if db.startswith('mule_')]

if not mule_dbs:
    print("❌ 没有找到以 mule_ 开头的数据库")
else:
    for db_name in mule_dbs:
        print(f"\n数据库: {db_name}")
        db = client[db_name]
        collections = db.list_collection_names()
        print(f"  集合: {collections}")
        
        if 'admin' in collections:
            admin_coll = db['admin']
            users = list(admin_coll.find({}))
            print(f"  用户数: {len(users)}")
            for u in users:
                print(f"    - {u.get('phone', 'N/A')} ({u.get('nickname', 'N/A')}) - is_deleted: {u.get('is_deleted', 'N/A')}")

# 提取租户ID
tenant_id = '68dda6cd04ba0d6c8dda4b7a'
mule_db_name = f'mule_{tenant_id}'

print(f"\n" + "="*60)
print(f"🔍 检查数据库: {mule_db_name}")
print(f"="*60)

if mule_db_name in db_names:
    db = client[mule_db_name]
    collections = db.list_collection_names()
    print(f"集合: {collections}")
    
    if 'admin' in collections:
        admin_coll = db['admin']
        user = admin_coll.find_one({'phone': '13838383388'})
        if user:
            print(f"\n✅ 找到用户 13838383388:")
            for key, value in user.items():
                print(f"  {key}: {value}")
        else:
            print(f"\n❌ 未找到用户 13838383388")
            all_users = list(admin_coll.find({}))
            print(f"\n所有用户:")
            for u in all_users:
                print(f"  - {u.get('phone')} ({u.get('nickname')})")
else:
    print(f"❌ 数据库 {mule_db_name} 不存在")

client.close()

