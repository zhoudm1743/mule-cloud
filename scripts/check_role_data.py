#!/usr/bin/env python3
"""检查 role 数据在哪个数据库"""

from pymongo import MongoClient

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

print("\n" + "="*60)
print("🔍 检查 role 数据分布")
print("="*60)

# 检查系统库
print("\n📦 系统库 (tenant_system):")
system_db = client['tenant_system']
if 'role' in system_db.list_collection_names():
    role_count = system_db['role'].count_documents({})
    print(f"  ✓ role 集合存在，共 {role_count} 条记录")
    if role_count > 0:
        roles = list(system_db['role'].find({}, {'_id': 1, 'name': 1, 'code': 1}).limit(5))
        for role in roles:
            print(f"    - {role.get('name')} ({role.get('code')})")
else:
    print("  ✗ role 集合不存在")

# 检查租户库
print("\n📦 租户库:")
all_dbs = client.list_database_names()
tenant_dbs = [db for db in all_dbs if db.startswith('mule_')]

if not tenant_dbs:
    print("  ⚠️  没有找到租户数据库 (mule_*)")
else:
    for db_name in tenant_dbs:
        db = client[db_name]
        tenant_id = db_name.replace('mule_', '')
        print(f"\n  数据库: {db_name} (tenant_id: {tenant_id})")
        
        if 'role' in db.list_collection_names():
            role_count = db['role'].count_documents({})
            print(f"    ✓ role 集合存在，共 {role_count} 条记录")
            if role_count > 0:
                roles = list(db['role'].find({}, {'_id': 1, 'name': 1, 'code': 1}).limit(3))
                for role in roles:
                    print(f"      - {role.get('name')} ({role.get('code')})")
        else:
            print(f"    ✗ role 集合不存在")

print("\n" + "="*60)
print("💡 结论:")
print("="*60)
print("如果 role 数据在系统库:")
print("  → RoleRepository.getCollection 应该返回系统库")
print("  → db := r.dbManager.GetDatabase(\"\")  // 空字符串")
print()
print("如果 role 数据在租户库:")
print("  → RoleRepository.getCollection 应该切换租户库")
print("  → db := r.dbManager.GetDatabase(tenantID)")
print("="*60)

client.close()

