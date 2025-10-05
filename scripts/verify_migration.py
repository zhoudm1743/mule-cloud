#!/usr/bin/env python3
"""
验证数据迁移结果
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
print("  数据迁移验证")
print("="*60)
print()

# 查看所有数据库
print("📚 所有数据库:")
dbs = client.list_database_names()
mule_dbs = [db for db in dbs if db.startswith('mule')]
for db in sorted(mule_dbs):
    print(f"  - {db}")
print()

# 验证系统数据库
print("📊 系统数据库 (mule):")
system_db = client['mule']
collections = ['admin', 'role', 'menu', 'basic', 'tenant']
for coll_name in collections:
    count = system_db[coll_name].count_documents({})
    print(f"  {coll_name}: {count} 条记录")
print()

# 验证租户数据库
tenant_dbs = [db for db in mule_dbs if db != 'mule']
for tenant_db_name in tenant_dbs:
    print(f"📊 租户数据库 ({tenant_db_name}):")
    tenant_db = client[tenant_db_name]
    collections = ['admin', 'role', 'menu', 'basic']
    for coll_name in collections:
        count = tenant_db[coll_name].count_documents({})
        if count > 0:
            print(f"  {coll_name}: {count} 条记录")
    print()

print("="*60)
print("  ✅ 验证完成！")
print("="*60)
print()
print("💡 下一步：启动应用服务测试")
print("  cd cmd/auth && go run main.go")
print("  cd cmd/system && go run main.go")
print()

client.close()

