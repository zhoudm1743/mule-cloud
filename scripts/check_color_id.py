#!/usr/bin/env python3
"""
检查 color 数据的 _id 类型
"""

from pymongo import MongoClient
from bson import ObjectId

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def check_color_id():
    """检查 color 数据的 _id 类型"""
    
    print("\n" + "="*60)
    print("🔍 检查 color 数据的 _id 类型")
    print("="*60)
    
    system_db = client['tenant_system']
    
    # 获取所有租户
    tenants = list(system_db['tenant'].find({'is_deleted': 0}))
    
    for tenant in tenants:
        tenant_id = str(tenant['_id'])
        tenant_name = tenant.get('name', 'Unknown')
        tenant_db_name = f"mule_{tenant_id}"
        
        print(f"\n📦 租户: {tenant_name}")
        print(f"   数据库: {tenant_db_name}")
        
        if tenant_db_name not in client.list_database_names():
            print("   ⚠️  数据库不存在")
            continue
        
        tenant_db = client[tenant_db_name]
        
        # 检查 basic 集合
        if 'basic' not in tenant_db.list_collection_names():
            print("   ⚠️  basic 集合不存在")
            continue
        
        basics = list(tenant_db['basic'].find({'is_deleted': 0}).limit(5))
        
        print(f"\n   basic 集合中的数据（前5条）:")
        for basic in basics:
            id_value = basic.get('_id')
            id_type = type(id_value).__name__
            type_name = basic.get('type', 'Unknown')
            name = basic.get('name', 'Unknown')
            
            print(f"      - ID: {id_value}")
            print(f"        类型: {id_type}")
            print(f"        Type: {type_name}")
            print(f"        Name: {name}")
            
            # 测试是否可以用字符串查询
            str_id = str(id_value)
            result_str = tenant_db['basic'].find_one({'_id': str_id})
            result_oid = None
            
            if isinstance(id_value, str):
                try:
                    oid = ObjectId(id_value)
                    result_oid = tenant_db['basic'].find_one({'_id': oid})
                except:
                    pass
            
            print(f"        字符串查询: {'✓' if result_str else '✗'}")
            print(f"        ObjectID查询: {'✓' if result_oid else '✗'}")
            print()
    
    print("="*60)

if __name__ == '__main__':
    try:
        check_color_id()
    finally:
        client.close()

