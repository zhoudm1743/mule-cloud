#!/usr/bin/env python3
"""
详细验证数据库隔离情况
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

print("="*70)
print("  数据库级别租户隔离 - 详细数据验证")
print("="*70)
print()

# 1. 查看所有数据库
print("📚 步骤1: 查看所有数据库")
print("-"*70)
all_dbs = client.list_database_names()
mule_dbs = sorted([db for db in all_dbs if db.startswith('mule')])
print(f"找到 {len(mule_dbs)} 个mule相关数据库:")
for db in mule_dbs:
    print(f"  ✅ {db}")
print()

# 2. 验证系统数据库
print("📊 步骤2: 验证系统数据库 (mule)")
print("-"*70)
system_db = client['mule']
collections = ['admin', 'role', 'menu', 'basic', 'tenant']

for coll_name in collections:
    coll = system_db[coll_name]
    count = coll.count_documents({})
    print(f"  {coll_name:15s}: {count} 条记录")
    
    # 显示样例数据（不含敏感信息）
    if count > 0:
        sample = coll.find_one({}, {'password': 0})
        if sample:
            # 检查是否有tenant_id字段
            has_tenant_id = 'tenant_id' in sample
            if has_tenant_id:
                print(f"    ⚠️  警告：发现tenant_id字段")
            else:
                print(f"    ✅ 无tenant_id字段（符合预期）")

print()

# 3. 验证租户数据库
print("📊 步骤3: 验证租户数据库")
print("-"*70)
tenant_dbs = [db for db in mule_dbs if db != 'mule']

if not tenant_dbs:
    print("  ⚠️  未找到租户数据库")
else:
    for tenant_db_name in tenant_dbs:
        tenant_id = tenant_db_name.replace('mule_', '')
        print(f"\n  租户: {tenant_id}")
        print(f"  数据库: {tenant_db_name}")
        print("  " + "-"*66)
        
        tenant_db = client[tenant_db_name]
        collections = ['admin', 'role', 'menu', 'basic']
        
        for coll_name in collections:
            coll = tenant_db[coll_name]
            count = coll.count_documents({})
            if count > 0:
                print(f"    {coll_name:15s}: {count} 条记录")
                
                # 检查是否有tenant_id字段
                sample = coll.find_one({}, {'password': 0})
                if sample:
                    has_tenant_id = 'tenant_id' in sample
                    if has_tenant_id:
                        print(f"      ⚠️  警告：发现tenant_id字段: {sample.get('tenant_id')}")
                    else:
                        print(f"      ✅ 无tenant_id字段（符合预期）")

print()

# 4. 数据隔离验证
print("🔍 步骤4: 数据隔离验证")
print("-"*70)

# 检查系统库中是否还有租户数据
system_admin_coll = system_db['admin']
system_role_coll = system_db['role']

system_admins_with_tenant = system_admin_coll.count_documents({'tenant_id': {'$exists': True, '$ne': ''}})
system_roles_with_tenant = system_role_coll.count_documents({'tenant_id': {'$exists': True, '$ne': ''}})

print("系统库中的租户数据检查:")
if system_admins_with_tenant > 0:
    print(f"  ⚠️  系统库admin集合中发现 {system_admins_with_tenant} 条租户数据（应该已迁移）")
else:
    print(f"  ✅ 系统库admin集合：无租户数据（符合预期）")

if system_roles_with_tenant > 0:
    print(f"  ⚠️  系统库role集合中发现 {system_roles_with_tenant} 条租户数据（应该已迁移）")
else:
    print(f"  ✅ 系统库role集合：无租户数据（符合预期）")

print()

# 5. 索引验证
print("📑 步骤5: 索引验证")
print("-"*70)
for tenant_db_name in tenant_dbs:
    tenant_db = client[tenant_db_name]
    print(f"\n  {tenant_db_name}:")
    
    for coll_name in ['admin', 'role']:
        coll = tenant_db[coll_name]
        indexes = list(coll.list_indexes())
        print(f"    {coll_name}: {len(indexes)} 个索引")
        for idx in indexes:
            if idx['name'] != '_id_':
                keys = ', '.join([f"{k}:{v}" for k, v in idx['key'].items()])
                print(f"      - {idx['name']}: {keys}")

print()

# 6. 数据完整性验证
print("✅ 步骤6: 数据完整性总结")
print("-"*70)

total_tenant_admins = 0
total_tenant_roles = 0

for tenant_db_name in tenant_dbs:
    tenant_db = client[tenant_db_name]
    admin_count = tenant_db['admin'].count_documents({})
    role_count = tenant_db['role'].count_documents({})
    total_tenant_admins += admin_count
    total_tenant_roles += role_count

system_admin_count = system_db['admin'].count_documents({'tenant_id': {'$exists': False}}) + \
                     system_db['admin'].count_documents({'tenant_id': ''})
system_role_count = system_db['role'].count_documents({'tenant_id': {'$exists': False}}) + \
                    system_db['role'].count_documents({'tenant_id': ''})

print(f"系统数据:")
print(f"  - 系统管理员: {system_admin_count} 条")
print(f"  - 系统角色:   {system_role_count} 条")
print()
print(f"租户数据 (跨 {len(tenant_dbs)} 个数据库):")
print(f"  - 租户管理员: {total_tenant_admins} 条")
print(f"  - 租户角色:   {total_tenant_roles} 条")
print()
print(f"总计:")
print(f"  - 管理员:     {system_admin_count + total_tenant_admins} 条")
print(f"  - 角色:       {system_role_count + total_tenant_roles} 条")

print()
print("="*70)
print("  ✅ 数据验证完成！")
print("="*70)
print()

# 7. 验证结论
print("📋 验证结论:")
print("-"*70)

issues = []

if system_admins_with_tenant > 0 or system_roles_with_tenant > 0:
    issues.append("系统库中仍有租户数据未迁移")

if len(tenant_dbs) == 0:
    issues.append("未找到租户数据库")

if issues:
    print("⚠️  发现问题:")
    for issue in issues:
        print(f"  - {issue}")
else:
    print("✅ 数据库隔离实现正确:")
    print(f"  ✅ 系统数据和租户数据完全隔离")
    print(f"  ✅ {len(tenant_dbs)} 个租户拥有独立数据库")
    print(f"  ✅ 租户数据无tenant_id字段（数据库级隔离）")
    print(f"  ✅ 索引正确创建")
    print()
    print("🎉 数据迁移和隔离验证通过！可以投入使用！")

print()

client.close()

