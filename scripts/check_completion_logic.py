#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查工序完成逻辑"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
tenant_db = client['mule_ace']

order_id = '68f3973b075a19ada180fcbe'

print('=' * 60)
print('Checking Procedure Completion Logic')
print('=' * 60)

# 查找订单的工序进度
procedures = list(tenant_db.order_procedure_progress.find({'order_id': order_id}))

print(f"\nAnalyzing {len(procedures)} procedures:\n")

for proc in procedures:
    name = proc.get('procedure_name', 'Unknown')
    total = proc.get('total_qty', 0)
    reported = proc.get('reported_qty', 0)
    progress_pct = proc.get('progress', 0)
    is_completed = proc.get('is_completed', False)
    
    should_be_completed = (reported >= total) if total > 0 else False
    
    status = 'OK' if is_completed else 'NO'
    
    print(f"  [{status}] {name}")
    print(f"      Reported: {reported} / {total} = {progress_pct:.2f}%")
    print(f"      is_completed: {is_completed}")
    
    if should_be_completed and not is_completed:
        print(f"      [WARN] Should be marked as completed!")
    elif not should_be_completed and reported > 0:
        print(f"      [INFO] In progress")
    elif reported == 0:
        print(f"      [INFO] Not started")
    
    print()

print('=' * 60)
print('Summary:')
print('=' * 60)

should_be_completed_count = sum(1 for p in procedures if p.get('reported_qty', 0) >= p.get('total_qty', 0))
actually_completed_count = sum(1 for p in procedures if p.get('is_completed', False))

print(f"  Procedures that should be completed: {should_be_completed_count}")
print(f"  Procedures marked as completed: {actually_completed_count}")

if should_be_completed_count != actually_completed_count:
    print(f"\n  [WARN] Mismatch detected! Some procedures are not properly marked.")

client.close()

