#!/usr/bin/env python3
"""
检查用户的所有字段
"""

from pymongo import MongoClient
import json

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
    print("📋 用户所有字段:")
    print(json.dumps({k: str(v) if not isinstance(v, (str, int, list, dict)) else v for k, v in user.items()}, indent=2, ensure_ascii=False))
else:
    print("❌ 用户不存在")

client.close()

