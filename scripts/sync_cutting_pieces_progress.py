#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
同步裁片监控进度
根据工序上报记录，更新裁片监控的progress字段
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

def sync_progress():
    """同步裁片监控进度"""
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
        
        # 获取所有上报记录，按扎号分组统计
        reports = list(reports_collection.find({
            "is_deleted": 0,
            "bundle_no": {"$ne": "", "$exists": True}
        }))
        
        print(f"[统计] 共找到 {len(reports)} 条上报记录")
        
        # 统计每个扎号完成了多少个不同的工序
        bundle_progress = defaultdict(set)  # {bundle_no: set(procedure_seq)}
        
        for report in reports:
            bundle_no = report.get("bundle_no")
            procedure_seq = report.get("procedure_seq")
            
            if bundle_no and procedure_seq:
                bundle_progress[bundle_no].add(procedure_seq)
        
        print(f"[统计] 共有 {len(bundle_progress)} 个扎号有上报记录")
        
        # 更新每个扎号的裁片监控进度
        updated_count = 0
        for bundle_no, procedures in bundle_progress.items():
            progress = len(procedures)  # 完成的工序数
            
            # 更新该扎号的所有裁片监控记录
            result = pieces_collection.update_many(
                {"bundle_no": bundle_no},
                {"$set": {"progress": progress}}
            )
            
            if result.modified_count > 0:
                updated_count += result.modified_count
                print(f"[更新] 扎号={bundle_no}, 进度={progress}道工序, 更新了{result.modified_count}条记录")
        
        print(f"\n[结果] 同步完成: 更新了 {updated_count} 条裁片监控记录")
            
    except Exception as e:
        print(f"[错误] 同步失败: {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()
        print("\n[完成] 已关闭MongoDB连接")

if __name__ == "__main__":
    print("=" * 60)
    print("开始同步裁片监控进度")
    print("=" * 60)
    sync_progress()
    print("=" * 60)
    print("同步任务完成")
    print("=" * 60)

