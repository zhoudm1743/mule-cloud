#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复裁片监控记录的工序数量
将 TotalProcess=0 的记录更新为订单的实际工序数量
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

def fix_total_process():
    """修复裁片监控记录的工序数量"""
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
        print(f"[统计] 共找到 {len(pieces)} 条需要修复的记录 (total_process=0)")
        
        if len(pieces) == 0:
            print("[信息] 没有需要修复的记录")
            return
        
        fixed_count = 0
        failed_count = 0
        order_cache = {}  # 缓存订单数据，避免重复查询
        
        # 修复每条记录
        for piece in pieces:
            order_id = piece.get("order_id")
            contract_no = piece.get("contract_no")
            bed_no = piece.get("bed_no")
            bundle_no = piece.get("bundle_no")
            
            try:
                # 从缓存或数据库获取订单信息
                if order_id not in order_cache:
                    order = orders_collection.find_one({"_id": order_id})
                    if order:
                        order_cache[order_id] = order
                    else:
                        print(f"[警告] 找不到订单: order_id={order_id}, contract_no={contract_no}")
                        failed_count += 1
                        continue
                else:
                    order = order_cache[order_id]
                
                # 获取工序数量
                procedures = order.get("procedures", [])
                total_process = len(procedures)
                
                if total_process > 0:
                    # 更新裁片监控记录
                    result = pieces_collection.update_one(
                        {"_id": piece["_id"]},
                        {"$set": {"total_process": total_process}}
                    )
                    
                    if result.modified_count > 0:
                        fixed_count += 1
                        print(f"[修复] 合同号={contract_no}, 床号={bed_no}, 扎号={bundle_no}, 尺码={piece.get('size')} -> {total_process}道工序")
                    else:
                        print(f"[跳过] 合同号={contract_no}, 床号={bed_no}, 扎号={bundle_no} (已是最新)")
                else:
                    print(f"[警告] 订单没有工序: order_id={order_id}, contract_no={contract_no}")
                    failed_count += 1
                    
            except Exception as e:
                print(f"[错误] 修复失败: contract_no={contract_no}, 错误={str(e)}")
                failed_count += 1
        
        print(f"\n[结果] 修复完成: 成功={fixed_count}, 失败={failed_count}, 总计={len(pieces)}")
            
    except Exception as e:
        print(f"[错误] 修复失败: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()
        print("\n[完成] 已关闭MongoDB连接")

if __name__ == "__main__":
    print("=" * 60)
    print("开始修复裁片监控记录的工序数量")
    print("=" * 60)
    fix_total_process()
    print("=" * 60)
    print("修复任务完成")
    print("=" * 60)

