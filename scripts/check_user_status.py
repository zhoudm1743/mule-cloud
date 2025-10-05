#!/usr/bin/env python3
"""
检查并修复用户状态
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
user = admin_coll.find_one({'phone': '13838383388'})

if user:
    print("📋 用户完整信息:")
    print(f"  _id: {user.get('_id')}")
    print(f"  phone: {user.get('phone')}")
    print(f"  nickname: {user.get('nickname')}")
    print(f"  password: {user.get('password', 'N/A')[:20]}...")  # 只显示前20位
    print(f"  roles: {user.get('roles')}")
    print(f"  status: {user.get('status')}")
    print(f"  is_deleted: {user.get('is_deleted')}")
    
    # 确保状态正确
    if user.get('status') != 1:
        print(f"\n⚠️  用户状态异常，修复为 1（启用）")
        admin_coll.update_one(
            {'phone': '13838383388'},
            {'$set': {'status': 1}}
        )
        print("✅ 状态已修复")
    else:
        print(f"\n✅ 用户状态正常")
    
    if user.get('is_deleted') != 0:
        print(f"\n⚠️  用户已删除，修复为 0（未删除）")
        admin_coll.update_one(
            {'phone': '13838383388'},
            {'$set': {'is_deleted': 0}}
        )
        print("✅ 删除标记已修复")
    else:
        print(f"✅ 用户未删除")
else:
    print("❌ 用户不存在")

client.close()

