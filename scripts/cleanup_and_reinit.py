#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""清理错误数据并重新初始化工作流"""

from pymongo import MongoClient
from bson import ObjectId
import time

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

system_db = client['mule_system']
tenant_db = client['mule_ace']

print('=' * 60)
print('Step 1: Cleanup Invalid Data')
print('=' * 60)

# 1. 删除租户数据库中的 workflow_definitions（不应该在租户库）
result = tenant_db.workflow_definitions.delete_many({})
print(f"\n[OK] Deleted {result.deleted_count} workflow_definitions from mule_ace")

# 2. 删除所有现有的 workflow_instances（将重新创建）
result = tenant_db.workflow_instances.delete_many({})
print(f"[OK] Deleted {result.deleted_count} workflow_instances from mule_ace")

# 3. 清空订单中的 workflow_instance 字段
result = tenant_db.orders.update_many(
    {},
    {'$set': {'workflow_instance': ''}}
)
print(f"[OK] Cleared workflow_instance from {result.modified_count} orders")

print('\n' + '=' * 60)
print('Step 2: Reinitialize Workflow Instances')
print('=' * 60)

# 获取系统数据库中的工作流定义
workflow_def = system_db.workflow_definitions.find_one({'code': 'basic_order', 'is_active': True})

if not workflow_def:
    print('\n[ERROR] No active basic_order workflow definition found in mule_system!')
    client.close()
    exit(1)

workflow_id = str(workflow_def['_id'])
print(f"\n[OK] Using workflow definition from mule_system")
print(f"     Workflow ID: {workflow_id}")
print(f"     Name: {workflow_def['name']}")

# 状态映射
status_to_state = {
    0: 'draft',
    1: 'ordered',
    2: 'production',
    3: 'completed',
    4: 'cancelled'
}

# 查找所有订单
orders = list(tenant_db.orders.find({'is_deleted': 0}))
print(f"\nFound {len(orders)} orders to reinitialize")

created_count = 0

for order in orders:
    order_id = str(order['_id'])
    contract_no = order.get('contract_no', 'Unknown')
    status = order.get('status', 0)
    workflow_state = order.get('workflow_state', '')
    
    # 确定当前状态
    current_state = workflow_state if workflow_state else status_to_state.get(status, 'draft')
    
    # 创建新的工作流实例
    new_instance_id = ObjectId()
    now = int(time.time())
    
    instance = {
        '_id': new_instance_id,
        'workflow_id': workflow_id,  # 🔥 使用系统数据库中的 workflow_id
        'entity_type': 'order',
        'entity_id': order_id,
        'current_state': current_state,
        'variables': {},
        'history': [{
            'from_state': '',
            'to_state': current_state,
            'event': 'init',
            'operator': 'system',
            'reason': 'Reinitialize with correct workflow_id from mule_system',
            'timestamp': now,
            'metadata': {}
        }],
        'created_at': now,
        'updated_at': now
    }
    
    # 插入工作流实例
    tenant_db.workflow_instances.insert_one(instance)
    
    # 更新订单
    tenant_db.orders.update_one(
        {'_id': order['_id']},
        {'$set': {
            'workflow_instance': str(new_instance_id),
            'workflow_code': 'basic_order',
            'workflow_state': current_state,
            'updated_at': now
        }}
    )
    
    created_count += 1
    print(f"\n[{created_count}] Order: {contract_no}")
    print(f"     State: {current_state}")
    print(f"     Instance ID: {new_instance_id}")

print('\n' + '=' * 60)
print('Complete!')
print('=' * 60)
print(f"  Created {created_count} workflow instances")
print(f"  All instances now use workflow_id from mule_system: {workflow_id}")
print(f"\n[NEXT] Please restart Order Service!")

client.close()

