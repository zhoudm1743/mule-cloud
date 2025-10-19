#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查订单工序进度数据"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
tenant_db = client['mule_ace']

order_id = '68f3973b075a19ada180fcbe'

print('=' * 60)
print(f'Checking Order Procedure Progress: {order_id}')
print('=' * 60)

# 查找订单的工序进度
procedures = list(tenant_db.order_procedure_progress.find({'order_id': order_id}))

print(f"\nFound {len(procedures)} procedure progress records:")

for proc in procedures:
    print(f"\n  Procedure: {proc.get('procedure_name')}")
    print(f"    seq: {proc.get('procedure_seq')}")
    print(f"    total_qty: {proc.get('total_qty')}")
    print(f"    reported_qty: {proc.get('reported_qty')}")
    print(f"    progress: {proc.get('progress')}")
    print(f"    is_completed: {proc.get('is_completed', False)}")

print('\n' + '=' * 60)
print('Statistics:')
print('=' * 60)

total = len(procedures)
completed = sum(1 for p in procedures if p.get('is_completed', False))
in_progress = sum(1 for p in procedures if not p.get('is_completed', False) and p.get('reported_qty', 0) > 0)
pending = sum(1 for p in procedures if p.get('reported_qty', 0) == 0)

print(f"  Total: {total}")
print(f"  Completed: {completed}")
print(f"  In Progress: {in_progress}")
print(f"  Pending: {pending}")

print('\n' + '=' * 60)

client.close()

