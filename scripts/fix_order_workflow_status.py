#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""修复现有订单的工作流状态"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

# 查询所有有进度但状态不对的订单
orders = list(db.orders.find({
    'is_deleted': 0,
    'progress': {'$gt': 0},  # 有进度
    '$or': [
        {'status': 0},  # 草稿
        {'status': 1},  # 已下单
    ]
}))

print('=' * 60)
print('查找需要修复的订单')
print('=' * 60)

if not orders:
    print('没有找到需要修复的订单')
else:
    print(f'找到 {len(orders)} 个需要修复的订单:\n')
    
    for order in orders:
        contract_no = order.get('contract_no', '未知')
        progress = order.get('progress', 0)
        status = order.get('status', 0)
        status_name = ['草稿', '已下单', '生产中', '已完成', '已取消'][status]
        
        print(f"订单: {contract_no}")
        print(f"  进度: {progress * 100:.2f}%")
        print(f"  当前状态: {status} ({status_name})")
        print(f"  应该状态: 2 (生产中)")
        
        # 更新订单状态为生产中
        result = db.orders.update_one(
            {'_id': order['_id']},
            {'$set': {
                'status': 2,  # 生产中
                'updated_at': int(__import__('time').time())
            }}
        )
        
        if result.modified_count > 0:
            print(f"  [OK] 已更新\n")
        else:
            print(f"  [SKIP] 无需更新\n")

print('=' * 60)
print('修复完成！')
print('=' * 60)

client.close()

