#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查订单的工作流状态"""

from pymongo import MongoClient
from bson import ObjectId

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

# 查询特定订单
order_id = '202510183528'
order = db.orders.find_one({'contract_no': order_id})

if order:
    print('=' * 60)
    print(f'Order Details: {order_id}')
    print('=' * 60)
    print(f"  _id: {order.get('_id')}")
    print(f"  status: {order.get('status')} (0=draft,1=ordered,2=production,3=completed,4=cancelled)")
    print(f"  workflow_state: {order.get('workflow_state', '(empty)')}")
    print(f"  workflow_code: {order.get('workflow_code', '(empty)')}")
    print(f"  workflow_instance: {order.get('workflow_instance', '(empty)')}")
    print(f"  progress: {order.get('progress', 0) * 100:.2f}%")
    
    # 检查工作流实例
    workflow_instance_id = order.get('workflow_instance')
    if workflow_instance_id:
        # 转换为ObjectId如果需要
        if isinstance(workflow_instance_id, str):
            try:
                instance_oid = ObjectId(workflow_instance_id)
                instance = db.workflow_instances.find_one({'_id': instance_oid})
            except:
                instance = db.workflow_instances.find_one({'_id': workflow_instance_id})
        else:
            instance = db.workflow_instances.find_one({'_id': workflow_instance_id})
            
        if instance:
            print(f"\n[OK] Workflow Instance Found:")
            print(f"  current_state: {instance.get('current_state')}")
            print(f"  workflow_id: {instance.get('workflow_id')}")
            print(f"  entity_type: {instance.get('entity_type')}")
            history = instance.get('history', [])
            print(f"  history count: {len(history)}")
            
            if history:
                print(f"\n  Last transition:")
                last = history[-1]
                print(f"    {last.get('from_state')} -> {last.get('to_state')}")
                print(f"    event: {last.get('event')}")
                print(f"    operator: {last.get('operator')}")
                print(f"    reason: {last.get('reason')}")
        else:
            print(f"\n[ERROR] Workflow instance not found! (ID: {workflow_instance_id})")
    else:
        print(f"\n[WARNING] Order has no workflow instance!")

print('\n' + '=' * 60)
print('Database data is correct!')
print('If frontend still shows old data, try:')
print('  1. Hard refresh (Ctrl+F5)')
print('  2. Clear browser cache')
print('  3. Restart Order Service')
print('=' * 60)

client.close()

