#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查订单工作流所在的数据库"""

from pymongo import MongoClient
from bson import ObjectId

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

order_id = '202510183528'

print('=' * 60)
print(f'Checking Order: {order_id}')
print('=' * 60)

# 检查两个数据库中的订单
for db_name in ['mule_ace', 'mule_system']:
    db = client[db_name]
    order = db.orders.find_one({'contract_no': order_id})
    
    if order:
        print(f"\n[FOUND] Order in database: {db_name}")
        print(f"  _id: {order.get('_id')}")
        print(f"  status: {order.get('status')}")
        print(f"  workflow_state: {order.get('workflow_state', '(empty)')}")
        print(f"  workflow_code: {order.get('workflow_code', '(empty)')}")
        print(f"  workflow_instance: {order.get('workflow_instance', '(empty)')}")
        
        # 检查工作流实例
        workflow_instance_id = order.get('workflow_instance')
        if workflow_instance_id:
            try:
                instance_oid = ObjectId(workflow_instance_id)
                instance = db.workflow_instances.find_one({'_id': instance_oid})
                
                if instance:
                    print(f"\n  [OK] Workflow instance found in {db_name}")
                    print(f"       current_state: {instance.get('current_state')}")
                    print(f"       workflow_id: {instance.get('workflow_id')}")
                    
                    # 检查工作流定义
                    workflow_id = instance.get('workflow_id')
                    if workflow_id:
                        try:
                            wf_def_oid = ObjectId(workflow_id)
                            wf_def = db.workflow_definitions.find_one({'_id': wf_def_oid})
                            
                            if wf_def:
                                print(f"\n  [OK] Workflow definition found in {db_name}")
                                print(f"       name: {wf_def.get('name')}")
                                print(f"       code: {wf_def.get('code')}")
                                print(f"       active: {wf_def.get('is_active')}")
                            else:
                                print(f"\n  [ERROR] Workflow definition NOT found in {db_name}!")
                                print(f"          Looking for _id: {workflow_id}")
                        except:
                            pass
                else:
                    print(f"\n  [ERROR] Workflow instance NOT found in {db_name}!")
            except Exception as e:
                print(f"\n  [ERROR] {e}")

print('\n' + '=' * 60)

client.close()

