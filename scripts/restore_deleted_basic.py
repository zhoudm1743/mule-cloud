#!/usr/bin/env python3
"""
恢复或清理已软删除的 basic 数据（用于测试）
"""

from pymongo import MongoClient
from bson import ObjectId

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def restore_deleted_basic():
    """恢复所有已软删除的 basic 数据"""
    
    print("\n" + "="*60)
    print("🔧 恢复已软删除的 basic 数据")
    print("="*60)
    
    system_db = client['tenant_system']
    
    # 获取所有租户
    tenants = list(system_db['tenant'].find({'is_deleted': 0}))
    
    total_restored = 0
    
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_db_name = f"mule_{tenant_id}"
        
        if tenant_db_name not in client.list_database_names():
            continue
        
        tenant_db = client[tenant_db_name]
        
        if 'basic' not in tenant_db.list_collection_names():
            continue
        
        # 查找已删除的记录
        deleted = list(tenant_db['basic'].find({'is_deleted': 1}))
        
        if not deleted:
            continue
        
        print(f"\n📦 租户: {tenant_name}")
        print(f"   数据库: {tenant_db_name}")
        print(f"   已删除记录数: {len(deleted)}")
        
        # 恢复记录
        result = tenant_db['basic'].update_many(
            {'is_deleted': 1},
            {'$set': {'is_deleted': 0, 'deleted_at': 0}}
        )
        
        print(f"   ✅ 已恢复: {result.modified_count} 条记录")
        total_restored += result.modified_count
    
    print("\n" + "="*60)
    print(f"✅ 总共恢复了 {total_restored} 条记录")
    print("="*60)

if __name__ == '__main__':
    try:
        restore_deleted_basic()
    except Exception as e:
        print(f"\n❌ 恢复失败: {e}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()

