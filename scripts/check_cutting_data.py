#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
检查裁片监控和上报记录的数据
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
TENANT_DB = "mule_ace"

def check_data():
    """检查数据"""
    client = MongoClient(
        host=MONGO_HOST,
        port=MONGO_PORT,
        username=MONGO_USER,
        password=MONGO_PASSWORD,
        authSource=MONGO_AUTH_DB
    )
    
    try:
        db = client[TENANT_DB]
        
        pieces_collection = db["cutting_pieces"]
        reports_collection = db["procedure_reports"]
        batches_collection = db["cutting_batches"]
        
        # 检查所有裁片监控记录
        print("=" * 60)
        print("裁片监控记录:")
        print("=" * 60)
        pieces = list(pieces_collection.find({}).sort("bed_no", 1).sort("bundle_no", 1))
        for piece in pieces:
            print(f"床号={piece.get('bed_no')}, 扎号={piece.get('bundle_no')}, 尺码={piece.get('size')}, 进度={piece.get('progress')}/{piece.get('total_process')}")
        
        # 检查所有上报记录
        print("\n" + "=" * 60)
        print("上报记录:")
        print("=" * 60)
        reports = list(reports_collection.find({"is_deleted": 0}))
        for report in reports:
            batch_id = report.get('batch_id')
            bundle_no = report.get('bundle_no')
            procedure_name = report.get('procedure_name')
            
            # 获取床号
            bed_no = "?"
            if batch_id:
                batch = batches_collection.find_one({"_id": batch_id})
                if batch:
                    bed_no = batch.get('bed_no', '?')
            
            print(f"batch_id={batch_id}, 床号={bed_no}, 扎号={bundle_no}, 工序={procedure_name}")
        
        # 检查批次记录
        print("\n" + "=" * 60)
        print("批次记录:")
        print("=" * 60)
        batches = list(batches_collection.find({"is_deleted": 0}))
        for batch in batches:
            print(f"batch_id={batch.get('_id')}, 床号={batch.get('bed_no')}, 扎号={batch.get('bundle_no')}")
            
    except Exception as e:
        print(f"[错误] {str(e)}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()

if __name__ == "__main__":
    check_data()

