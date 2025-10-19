#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""为缺失工作流实例的订单创建实例"""

from pymongo import MongoClient
from bson import ObjectId
import time

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

# 🔥 系统数据库：存储工作流定义（全局共享）
system_db = client['mule_system']
# 🔥 租户数据库：存储订单和工作流实例（租户隔离）
tenant_db = client['mule_ace']

# 状态映射
status_to_state = {
    0: 'draft',
    1: 'ordered',
    2: 'production',
    3: 'completed',
    4: 'cancelled'
}

print('=' * 60)
print('Checking and Creating Missing Workflow Instances')
print('=' * 60)

# 🔥 重要：从系统数据库获取工作流定义
workflow_def = system_db.workflow_definitions.find_one({'code': 'basic_order', 'is_active': True})

if not workflow_def:
    print('[ERROR] No active basic_order workflow definition found!')
    client.close()
    exit(1)

workflow_id = str(workflow_def['_id'])
print(f"\n[OK] Found workflow definition: {workflow_def['name']}")
print(f"     Workflow ID: {workflow_id}")

# 查找所有订单（从租户数据库）
orders = list(tenant_db.orders.find({'is_deleted': 0}))
print(f"\nFound {len(orders)} orders")

created_count = 0
fixed_count = 0

for order in orders:
    order_id = str(order['_id'])
    contract_no = order.get('contract_no', 'Unknown')
    status = order.get('status', 0)
    workflow_state = order.get('workflow_state', '')
    workflow_instance_id = order.get('workflow_instance', '')
    
    # 检查是否需要创建工作流实例
    needs_instance = False
    
    if not workflow_instance_id:
        print(f"\n[CREATE] Order {contract_no} has no workflow_instance")
        needs_instance = True
    else:
        # 检查实例是否存在（从租户数据库）
        try:
            instance_oid = ObjectId(workflow_instance_id)
            instance = tenant_db.workflow_instances.find_one({'_id': instance_oid})
            if not instance:
                print(f"\n[FIX] Order {contract_no} has invalid workflow_instance: {workflow_instance_id}")
                needs_instance = True
        except:
            print(f"\n[FIX] Order {contract_no} has invalid workflow_instance format: {workflow_instance_id}")
            needs_instance = True
    
    if needs_instance:
        # 确定当前状态
        current_state = workflow_state if workflow_state else status_to_state.get(status, 'draft')
        
        # 创建新的工作流实例
        new_instance_id = ObjectId()
        now = int(time.time())
        
        instance = {
            '_id': new_instance_id,
            'workflow_id': workflow_id,
            'entity_type': 'order',
            'entity_id': order_id,
            'current_state': current_state,
            'variables': {},
            'history': [{
                'from_state': '',
                'to_state': current_state,
                'event': 'init',
                'operator': 'system',
                'reason': 'Auto-initialize workflow instance',
                'timestamp': now,
                'metadata': {}
            }],
            'created_at': now,
            'updated_at': now
        }
        
        # 插入工作流实例（到租户数据库）
        tenant_db.workflow_instances.insert_one(instance)
        
        # 更新订单（在租户数据库）
        tenant_db.orders.update_one(
            {'_id': order['_id']},
            {'$set': {
                'workflow_instance': str(new_instance_id),
                'workflow_code': 'basic_order',
                'workflow_state': current_state,
                'updated_at': now
            }}
        )
        
        print(f"  [OK] Created workflow instance: {new_instance_id}")
        print(f"       State: {current_state}")
        
        if workflow_instance_id:
            fixed_count += 1
        else:
            created_count += 1

print('\n' + '=' * 60)
print('Complete!')
print('=' * 60)
print(f"  Created: {created_count} instances")
print(f"  Fixed:   {fixed_count} instances")
print(f"\nPlease restart Order Service and refresh the page!")

client.close()

