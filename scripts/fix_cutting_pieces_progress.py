#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
修复裁片监控进度
1. 重置所有裁片监控的progress为0
2. 根据工序上报记录重新统计（同时匹配床号和扎号）
"""

import sys
from pymongo import MongoClient
from collections import defaultdict

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

def fix_progress():
    """修复裁片监控进度"""
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
        reports_collection = db["procedure_reports"]
        batches_collection = db["cutting_batches"]
        
        # 步骤1: 重置所有裁片监控的进度为0
        result = pieces_collection.update_many(
            {},
            {"$set": {"progress": 0}}
        )
        print(f"[重置] 已重置 {result.modified_count} 条裁片监控记录的进度为0")
        
        # 步骤2: 获取所有上报记录
        reports = list(reports_collection.find({
            "is_deleted": 0,
            "batch_id": {"$ne": "", "$exists": True}
        }))
        
        print(f"[统计] 共找到 {len(reports)} 条有batch_id的上报记录")
        
        # 步骤3: 统计每个床号+扎号的工序进度
        # key: (bed_no, bundle_no), value: set(procedure_seq)
        progress_map = defaultdict(set)
        
        for report in reports:
            batch_id = report.get("batch_id")
            bundle_no = report.get("bundle_no")
            procedure_seq = report.get("procedure_seq")
            
            if not batch_id or not bundle_no or not procedure_seq:
                continue
            
            # 从批次获取床号
            batch = batches_collection.find_one({"_id": batch_id})
            if not batch:
                print(f"[警告] 找不到批次: batch_id={batch_id}")
                continue
            
            bed_no = batch.get("bed_no")
            if not bed_no:
                continue
            
            # 记录该床号+扎号完成的工序
            progress_map[(bed_no, bundle_no)].add(procedure_seq)
        
        print(f"[统计] 共有 {len(progress_map)} 个床号+扎号组合有上报记录")
        
        # 步骤4: 更新裁片监控进度
        updated_count = 0
        for (bed_no, bundle_no), procedures in progress_map.items():
            progress = len(procedures)  # 完成的工序数
            
            # 更新该床号+扎号的所有裁片监控记录
            result = pieces_collection.update_many(
                {
                    "bed_no": bed_no,
                    "bundle_no": bundle_no
                },
                {"$set": {"progress": progress}}
            )
            
            if result.modified_count > 0:
                updated_count += result.modified_count
                print(f"[更新] 床号={bed_no}, 扎号={bundle_no}, 进度={progress}道工序, 更新了{result.modified_count}条记录")
        
        print(f"\n[结果] 修复完成: 更新了 {updated_count} 条裁片监控记录")
            
    except Exception as e:
        print(f"[错误] 修复失败: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()
        print("\n[完成] 已关闭MongoDB连接")

if __name__ == "__main__":
    print("=" * 60)
    print("开始修复裁片监控进度")
    print("=" * 60)
    fix_progress()
    print("=" * 60)
    print("修复任务完成")
    print("=" * 60)

