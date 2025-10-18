#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
清理无效的裁片监控记录
删除那些订单已不存在的裁片监控记录
"""

import sys
from pymongo import MongoClient

# 设置输出编码
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')

# MongoDB连接配置
MONGO_HOST = "127.0.0.1"
MONGO_PORT = 27015
MONGO_USER = "root"
MONGO_PASSWORD = "bgg8384495"
MONGO_AUTH_DB = "admin"
TENANT_DB = "mule_ace"  # 租户数据库名

def clean_invalid_pieces():
    """清理无效的裁片监控记录"""
    # 连接MongoDB
    client = MongoClient(
        host=MONGO_HOST,
        port=MONGO_PORT,
        username=MONGO_USER,
        password=MONGO_PASSWORD,
        authSource=MONGO_AUTH_DB
    )
    
    try:
        # 选择租户数据库
        db = client[TENANT_DB]
        
        pieces_collection = db["cutting_pieces"]
        orders_collection = db["orders"]
        
        # 获取所有 TotalProcess=0 的裁片监控记录
        pieces = list(pieces_collection.find({"total_process": 0}))
        print(f"[统计] 共找到 {len(pieces)} 条 total_process=0 的记录")
        
        if len(pieces) == 0:
            print("[信息] 没有需要处理的记录")
            return
        
        invalid_ids = []
        
        # 检查每条记录
        for piece in pieces:
            order_id = piece.get("order_id")
            
            # 检查订单是否存在
            order = orders_collection.find_one({"_id": order_id})
            
            if not order:
                # 订单不存在，标记为无效记录
                invalid_ids.append(piece["_id"])
                print(f"[无效] 合同号={piece.get('contract_no')}, 床号={piece.get('bed_no')}, 扎号={piece.get('bundle_no')} (订单不存在)")
        
        print(f"\n[统计] 共发现 {len(invalid_ids)} 条无效记录 (订单已不存在)")
        
        if len(invalid_ids) > 0:
            # 删除无效记录
            result = pieces_collection.delete_many({
                "_id": {"$in": invalid_ids}
            })
            print(f"[成功] 删除完成！成功删除了 {result.deleted_count} 条无效记录")
        else:
            print("[信息] 没有发现无效记录")
            
    except Exception as e:
        print(f"[错误] 清理失败: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()
        print("\n[完成] 已关闭MongoDB连接")

if __name__ == "__main__":
    print("=" * 60)
    print("开始清理无效的裁片监控记录")
    print("=" * 60)
    clean_invalid_pieces()
    print("=" * 60)
    print("清理任务完成")
    print("=" * 60)

