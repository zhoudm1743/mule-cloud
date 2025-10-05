#!/usr/bin/env python3
"""
清理系统库中已迁移的租户数据
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
print("  清理系统库中已迁移的租户数据")
print("="*60)
print()

system_db = client['mule']

# 检查系统库中的租户数据
print("🔍 检查系统库中的租户数据...")
print()

admin_coll = system_db['admin']
role_coll = system_db['role']

# 查找有tenant_id的admin
admins_with_tenant = list(admin_coll.find({'tenant_id': {'$exists': True, '$ne': ''}}))
print(f"📊 Admin集合中有tenant_id的记录: {len(admins_with_tenant)}")
for admin in admins_with_tenant:
    print(f"  - ID: {admin.get('_id')}, TenantID: {admin.get('tenant_id')}, Phone: {admin.get('phone')}")

# 查找有tenant_id的role
roles_with_tenant = list(role_coll.find({'tenant_id': {'$exists': True, '$ne': ''}}))
print(f"📊 Role集合中有tenant_id的记录: {len(roles_with_tenant)}")
for role in roles_with_tenant:
    print(f"  - ID: {role.get('_id')}, TenantID: {role.get('tenant_id')}, Name: {role.get('name')}")

print()

if len(admins_with_tenant) == 0 and len(roles_with_tenant) == 0:
    print("✅ 系统库中无租户数据，无需清理")
    client.close()
    exit(0)

# 确认清理
print("⚠️  即将删除系统库中的租户数据（这些数据已迁移到租户数据库）")
print()
confirm = input("确认删除? (输入 yes 继续): ")

if confirm.lower() != 'yes':
    print("❌ 取消清理")
    client.close()
    exit(0)

print()
print("🗑️  开始清理...")
print()

# 删除有tenant_id的admin
if len(admins_with_tenant) > 0:
    result = admin_coll.delete_many({'tenant_id': {'$exists': True, '$ne': ''}})
    print(f"✅ Admin: 删除了 {result.deleted_count} 条租户数据")

# 删除有tenant_id的role
if len(roles_with_tenant) > 0:
    result = role_coll.delete_many({'tenant_id': {'$exists': True, '$ne': ''}})
    print(f"✅ Role: 删除了 {result.deleted_count} 条租户数据")

print()
print("="*60)
print("  ✅ 清理完成！")
print("="*60)
print()

# 验证结果
print("📊 清理后统计:")
system_admin_count = admin_coll.count_documents({})
system_role_count = role_coll.count_documents({})
print(f"  系统Admin: {system_admin_count} 条")
print(f"  系统Role: {system_role_count} 条")

print()
print("💡 建议重新运行验证脚本:")
print("  py scripts/verify_data_detail.py")

client.close()

