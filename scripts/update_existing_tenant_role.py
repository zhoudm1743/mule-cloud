#!/usr/bin/env python3
"""
更新现有租户的默认角色，添加完整的菜单和权限
"""

from pymongo import MongoClient
from bson import ObjectId

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def update_existing_tenant_role():
    """更新现有租户的角色权限"""
    
    print("\n" + "="*60)
    print("🔧 更新现有租户的角色权限")
    print("="*60)
    
    system_db = client['tenant_system']
    
    # 获取所有租户
    tenants = list(system_db['tenant'].find({'is_deleted': 0}))
    
    if not tenants:
        print("\n⚠️  没有找到租户")
        return
    
    # 默认菜单和权限
    default_menus = [
        "dashboard",
        "system",
        "system_admin",
        "system_role",
        "system_menu",
        "system_tenant",
    ]
    
    default_menu_permissions = {
        "system_admin": ["read", "create", "update", "delete"],
        "system_role": ["read", "create", "update", "delete", "menus"],
        "system_menu": ["read", "create", "delete"],
        "system_tenant": ["read", "create", "update", "delete", "menus"],
    }
    
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_db_name = f"mule_{tenant_id}"
        
        print(f"\n📦 处理租户: {tenant_name}")
        print(f"   数据库: {tenant_db_name}")
        
        tenant_db = client[tenant_db_name]
        
        # 查找所有角色
        roles = list(tenant_db['role'].find({'is_deleted': 0}))
        
        if not roles:
            print("   ⚠️  没有找到角色")
            continue
        
        for role in roles:
            role_id = role['_id']
            role_name = role.get('name', 'Unknown')
            role_code = role.get('code', 'unknown')
            current_menus = role.get('menus', [])
            current_perms = role.get('menu_permissions', {})
            
            print(f"\n   角色: {role_name} ({role_code})")
            print(f"      当前菜单数: {len(current_menus)}")
            print(f"      当前权限数: {len(current_perms)}")
            
            # 更新菜单和权限
            tenant_db['role'].update_one(
                {'_id': role_id},
                {
                    '$set': {
                        'menus': default_menus,
                        'menu_permissions': default_menu_permissions,
                    }
                }
            )
            
            print(f"      ✅ 已更新为:")
            print(f"         菜单数: {len(default_menus)}")
            print(f"         权限数: {len(default_menu_permissions)}")
    
    print("\n" + "="*60)
    print("✅ 更新完成！")
    print("="*60)
    
    # 验证
    print("\n📊 验证结果:")
    print("-"*60)
    
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_db_name = f"mule_{tenant_id}"
        tenant_db = client[tenant_db_name]
        
        roles = list(tenant_db['role'].find({'is_deleted': 0}))
        
        print(f"\n{tenant_name} ({tenant_db_name}):")
        for role in roles:
            role_name = role.get('name', 'Unknown')
            menus = role.get('menus', [])
            perms = role.get('menu_permissions', {})
            print(f"  - {role_name}: {len(menus)} 个菜单, {len(perms)} 个权限映射")

if __name__ == '__main__':
    try:
        update_existing_tenant_role()
    except Exception as e:
        print(f"\n❌ 更新失败: {e}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()

