#!/usr/bin/env python3
"""
列出所有数据库
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
print("📋 所有数据库:")
print("="*60)

db_names = client.list_database_names()
for db_name in db_names:
    if 'tenant' in db_name or 'mule' in db_name:
        db = client[db_name]
        collections = db.list_collection_names()
        print(f"\n数据库: {db_name}")
        print(f"  集合: {', '.join(collections)}")
        
        # 如果有 admin 集合，显示用户数
        if 'admin' in collections:
            admin_count = db['admin'].count_documents({'is_deleted': 0})
            print(f"  用户数: {admin_count}")
            if admin_count > 0:
                users = list(db['admin'].find({'is_deleted': 0}, {'phone': 1, 'nickname': 1}))
                for u in users:
                    print(f"    - {u.get('phone', 'N/A')} ({u.get('nickname', 'N/A')})")
        
        # 如果有 tenant 集合，显示租户数
        if 'tenant' in collections:
            tenant_count = db['tenant'].count_documents({})
            print(f"  租户数: {tenant_count}")
            if tenant_count > 0:
                tenants = list(db['tenant'].find({}, {'code': 1, 'name': 1, 'status': 1}))
                for t in tenants:
                    print(f"    - {t.get('code', 'N/A')} ({t.get('name', 'N/A')}) - status:{t.get('status', 'N/A')}")

client.close()

