#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查批次工序进度数据"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
tenant_db = client['mule_ace']

order_id = '68f3973b075a19ada180fcbe'

print('=' * 60)
print(f'Checking Batch Procedure Progress: {order_id}')
print('=' * 60)

# 查找订单
order = tenant_db.orders.find_one({'_id': order_id})
if order:
    contract_no = order.get('contract_no', 'Unknown')
    total_qty = order.get('quantity', 0)
    print(f"\nOrder: {contract_no}")
    print(f"Total Quantity: {total_qty} pieces")

# 查找该订单的所有批次工序进度
batch_progress = list(tenant_db.batch_procedure_progress.find({'order_id': order_id}))

print(f"\nFound {len(batch_progress)} batch procedure progress records")

# 按工序分组统计
from collections import defaultdict

procedure_stats = defaultdict(lambda: {
    'total_qty': 0,
    'reported_qty': 0,
    'completed_batches': 0,
    'in_progress_batches': 0,
    'pending_batches': 0,
    'batch_count': 0
})

for record in batch_progress:
    proc_name = record.get('procedure_name', 'Unknown')
    batch_qty = record.get('batch_quantity', 0)
    reported_qty = record.get('reported_qty', 0)
    is_completed = record.get('is_completed', False)
    
    stats = procedure_stats[proc_name]
    stats['batch_count'] += 1
    stats['total_qty'] += batch_qty
    stats['reported_qty'] += reported_qty
    
    if is_completed:
        stats['completed_batches'] += 1
    elif reported_qty > 0:
        stats['in_progress_batches'] += 1
    else:
        stats['pending_batches'] += 1

print('\n' + '=' * 60)
print('Statistics by Procedure:')
print('=' * 60)

for proc_name, stats in sorted(procedure_stats.items()):
    print(f"\n{proc_name}:")
    print(f"  Total batches: {stats['batch_count']}")
    print(f"  Total quantity: {stats['total_qty']} pieces")
    print(f"  Reported quantity: {stats['reported_qty']} pieces")
    print(f"  Completed batches: {stats['completed_batches']}")
    print(f"  In progress batches: {stats['in_progress_batches']}")
    print(f"  Pending batches: {stats['pending_batches']}")
    
    if stats['total_qty'] > 0:
        progress_pct = (stats['reported_qty'] / stats['total_qty']) * 100
        print(f"  Progress: {progress_pct:.2f}%")

# 整体统计
print('\n' + '=' * 60)
print('Overall Statistics:')
print('=' * 60)

total_reported = sum(s['reported_qty'] for s in procedure_stats.values())
total_batches = sum(s['batch_count'] for s in procedure_stats.values())
total_completed_batches = sum(s['completed_batches'] for s in procedure_stats.values())
total_in_progress_batches = sum(s['in_progress_batches'] for s in procedure_stats.values())
total_pending_batches = sum(s['pending_batches'] for s in procedure_stats.values())

print(f"  Total batches: {total_batches}")
print(f"  Completed batches: {total_completed_batches}")
print(f"  In progress batches: {total_in_progress_batches}")
print(f"  Pending batches: {total_pending_batches}")

client.close()

