#!/usr/bin/env python3
"""
查看租户详情
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

# 系统库
system_db = client['tenant_system']
tenant_coll = system_db['tenant']

print("="*60)
print("📋 所有租户的完整信息:")
print("="*60)

tenants = list(tenant_coll.find({}))
if not tenants:
    print("❌ 没有找到任何租户")
else:
    for t in tenants:
        print(f"\n租户 {t['_id']}:")
        # 打印所有字段
        for key, value in t.items():
            if key != '_id':
                print(f"  {key}: {value}")

# 查找 code 为 "default" 的租户
print("\n" + "="*60)
print("🔍 查找租户代码为 'default' 的租户:")
print("="*60)

default_tenant = tenant_coll.find_one({'code': 'default'})
if default_tenant:
    print(f"✅ 找到租户:")
    print(json.dumps({k: str(v) for k, v in default_tenant.items()}, indent=2, ensure_ascii=False))
else:
    print("❌ 未找到 code='default' 的租户")

client.close()

