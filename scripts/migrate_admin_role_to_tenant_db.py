#!/usr/bin/env python3
"""
将 admin 和 role 从系统库迁移到租户库
"""

from pymongo import MongoClient
from bson import ObjectId
import sys

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def migrate_to_tenant_databases():
    """将 admin 和 role 迁移到租户库"""
    
    print("\n" + "="*60)
    print("🚀 开始迁移 admin 和 role 到租户库")
    print("="*60)
    
    system_db = client['tenant_system']
    
    # 1. 获取所有租户
    tenants = list(system_db['tenant'].find({'is_deleted': 0}))
    
    if not tenants:
        print("\n⚠️  没有找到租户，无法迁移")
        return
    
    print(f"\n📋 找到 {len(tenants)} 个租户:")
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        print(f"  - {tenant.get('name')} (ID: {tenant_id}, Code: {tenant.get('code')})")
    
    # 2. 迁移每个租户的数据
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_code = tenant.get('code', 'unknown')
        tenant_db_name = f"mule_{tenant_id}"
        
        print(f"\n" + "-"*60)
        print(f"📦 处理租户: {tenant_name} ({tenant_code})")
        print(f"   目标数据库: {tenant_db_name}")
        print("-"*60)
        
        tenant_db = client[tenant_db_name]
        
        # 迁移 admin
        print("\n1️⃣  迁移 admin 数据...")
        admins_in_system = list(system_db['admin'].find({
            'is_deleted': 0
        }))
        
        if admins_in_system:
            print(f"   找到 {len(admins_in_system)} 个管理员")
            
            # 清空目标库的 admin 集合
            tenant_db['admin'].delete_many({})
            
            # 插入数据
            if admins_in_system:
                tenant_db['admin'].insert_many(admins_in_system)
                print(f"   ✅ 已迁移 {len(admins_in_system)} 个管理员到 {tenant_db_name}")
                
                for admin in admins_in_system:
                    print(f"      - {admin.get('nickname')} ({admin.get('phone')})")
        else:
            print(f"   ⚠️  系统库中没有管理员数据")
        
        # 迁移 role
        print("\n2️⃣  迁移 role 数据...")
        roles_in_system = list(system_db['role'].find({
            'is_deleted': 0
        }))
        
        if roles_in_system:
            print(f"   找到 {len(roles_in_system)} 个角色")
            
            # 清空目标库的 role 集合
            tenant_db['role'].delete_many({})
            
            # 插入数据
            if roles_in_system:
                tenant_db['role'].insert_many(roles_in_system)
                print(f"   ✅ 已迁移 {len(roles_in_system)} 个角色到 {tenant_db_name}")
                
                for role in roles_in_system:
                    print(f"      - {role.get('name')} ({role.get('code')})")
        else:
            print(f"   ⚠️  系统库中没有角色数据")
    
    print("\n" + "="*60)
    print("✅ 迁移完成！")
    print("="*60)
    
    # 3. 验证迁移结果
    print("\n📊 迁移验证:")
    print("-"*60)
    
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_db_name = f"mule_{tenant_id}"
        tenant_db = client[tenant_db_name]
        
        admin_count = tenant_db['admin'].count_documents({})
        role_count = tenant_db['role'].count_documents({})
        
        print(f"\n{tenant_name} ({tenant_db_name}):")
        print(f"  - admin: {admin_count} 条")
        print(f"  - role: {role_count} 条")
    
    # 4. 询问是否清理系统库
    print("\n" + "="*60)
    print("⚠️  注意：数据已迁移到租户库")
    print("="*60)
    print("\n系统库中的 admin 和 role 数据:")
    system_admin_count = system_db['admin'].count_documents({})
    system_role_count = system_db['role'].count_documents({})
    print(f"  - admin: {system_admin_count} 条")
    print(f"  - role: {system_role_count} 条")
    
    print("\n💡 建议:")
    print("  1. 先测试租户登录和角色功能")
    print("  2. 确认功能正常后，可以清理系统库中的数据")
    print("  3. 运行: py scripts/cleanup_system_admin_role.py")

if __name__ == '__main__':
    try:
        migrate_to_tenant_databases()
    except Exception as e:
        print(f"\n❌ 迁移失败: {e}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()

