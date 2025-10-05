#!/usr/bin/env python3
"""
查找特定 ID 的 color 数据
"""

from pymongo import MongoClient
from bson import ObjectId

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def find_color_by_id(target_id):
    """查找特定 ID 的 color 数据"""
    
    print("\n" + "="*60)
    print(f"🔍 查找 ID: {target_id}")
    print("="*60)
    
    # 获取所有数据库
    db_names = client.list_database_names()
    
    found = False
    
    for db_name in db_names:
        if db_name in ['admin', 'config', 'local']:
            continue
        
        db = client[db_name]
        
        # 检查 basic 集合
        if 'basic' not in db.list_collection_names():
            continue
        
        # 尝试字符串查询
        result_str = db['basic'].find_one({'_id': target_id})
        
        # 尝试 ObjectID 查询
        result_oid = None
        try:
            oid = ObjectId(target_id)
            result_oid = db['basic'].find_one({'_id': oid})
        except:
            pass
        
        if result_str or result_oid:
            found = True
            result = result_str or result_oid
            
            print(f"\n✅ 找到数据！")
            print(f"   数据库: {db_name}")
            print(f"   集合: basic")
            print(f"   _id 类型: {type(result['_id']).__name__}")
            print(f"   _id 值: {result['_id']}")
            print(f"   查询结果:")
            print(f"      - 字符串查询: {'✓' if result_str else '✗'}")
            print(f"      - ObjectID查询: {'✓' if result_oid else '✗'}")
            print(f"\n   完整数据:")
            for key, value in result.items():
                print(f"      {key}: {value}")
            
            # 检查 is_deleted 状态
            is_deleted = result.get('is_deleted', 0)
            print(f"\n   is_deleted: {is_deleted}")
            
            if is_deleted == 1:
                print("   ⚠️  该记录已被软删除！")
    
    if not found:
        print(f"\n❌ 在所有数据库中都没有找到 ID: {target_id}")
        print("\n提示：数据可能已被物理删除")
    
    print("\n" + "="*60)

if __name__ == '__main__':
    target_id = '68df0fe2aab0d60369d6d935'
    try:
        find_color_by_id(target_id)
    finally:
        client.close()

