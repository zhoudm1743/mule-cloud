#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
清理孤立的裁片监控记录
当批次被删除时，对应的裁片监控记录应该被删除
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

def clean_orphan_pieces():
    """清理孤立的裁片监控记录"""
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
        batches_collection = db["cutting_batches"]
        
        # 获取所有裁片监控记录
        pieces = list(pieces_collection.find({}))
        print(f"[统计] 共找到 {len(pieces)} 条裁片监控记录")
        
        orphan_count = 0
        deleted_ids = []
        
        # 检查每条裁片监控记录
        for piece in pieces:
            bed_no = piece.get("bed_no")
            bundle_no = piece.get("bundle_no")
            contract_no = piece.get("contract_no")
            
            # 查找对应的批次记录（未删除的）
            batch = batches_collection.find_one({
                "bed_no": bed_no,
                "bundle_no": bundle_no,
                "is_deleted": 0
            })
            
            # 如果找不到对应的批次，说明是孤立记录
            if not batch:
                orphan_count += 1
                deleted_ids.append(piece["_id"])
                print(f"[发现] 孤立记录: 合同号={contract_no}, 床号={bed_no}, 扎号={bundle_no}, 尺码={piece.get('size')}")
        
        print(f"\n[结果] 共发现 {orphan_count} 条孤立的裁片监控记录")
        
        if orphan_count > 0:
            # 删除孤立记录
            result = pieces_collection.delete_many({
                "_id": {"$in": deleted_ids}
            })
            print(f"[成功] 删除完成！成功删除了 {result.deleted_count} 条孤立记录")
        else:
            print("[信息] 没有发现孤立记录，无需清理")
            
    except Exception as e:
        print(f"[错误] 清理失败: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()
        print("\n[完成] 已关闭MongoDB连接")

if __name__ == "__main__":
    print("=" * 60)
    print("开始清理孤立的裁片监控记录")
    print("=" * 60)
    clean_orphan_pieces()
    print("=" * 60)
    print("清理任务完成")
    print("=" * 60)

