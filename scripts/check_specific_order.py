#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查特定订单的工作流状态"""

from pymongo import MongoClient
from bson import ObjectId

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

system_db = client['mule_system']
tenant_db = client['mule_ace']

order_id = '68f3973b075a19ada180fcbe'

print('=' * 60)
print(f'Checking Order: {order_id}')
print('=' * 60)

# 查找订单
try:
    order_oid = ObjectId(order_id)
    order = tenant_db.orders.find_one({'_id': order_oid})
    
    if not order:
        print('[ERROR] Order not found!')
        client.close()
        exit(1)
    
    print(f"\n[OK] Order found:")
    print(f"  contract_no: {order.get('contract_no')}")
    print(f"  status: {order.get('status')}")
    print(f"  workflow_state: {order.get('workflow_state', '(empty)')}")
    print(f"  workflow_code: {order.get('workflow_code', '(empty)')}")
    print(f"  workflow_instance: {order.get('workflow_instance', '(empty)')}")
    
    # 检查工作流实例
    workflow_instance_id = order.get('workflow_instance')
    if workflow_instance_id:
        try:
            instance_oid = ObjectId(workflow_instance_id)
            instance = tenant_db.workflow_instances.find_one({'_id': instance_oid})
            
            if instance:
                print(f"\n[OK] Workflow instance found:")
                print(f"  _id: {instance['_id']}")
                print(f"  current_state: {instance.get('current_state')}")
                print(f"  workflow_id: {instance.get('workflow_id')}")
                
                # 检查工作流定义
                workflow_id = instance.get('workflow_id')
                if workflow_id:
                    # 先在租户库找
                    wf_def = tenant_db.workflow_definitions.find_one({'_id': workflow_id})
                    if wf_def:
                        print(f"\n[WARN] Workflow definition found in TENANT db (should be in SYSTEM db):")
                        print(f"  _id: {wf_def['_id']}")
                        print(f"  code: {wf_def.get('code')}")
                        print(f"  name: {wf_def.get('name')}")
                    else:
                        # 再在系统库找
                        wf_def = system_db.workflow_definitions.find_one({'_id': workflow_id})
                        if wf_def:
                            print(f"\n[OK] Workflow definition found in SYSTEM db:")
                            print(f"  _id: {wf_def['_id']}")
                            print(f"  code: {wf_def.get('code')}")
                            print(f"  name: {wf_def.get('name')}")
                        else:
                            print(f"\n[ERROR] Workflow definition NOT found in either db!")
                            print(f"  Looking for _id: {workflow_id}")
            else:
                print(f"\n[ERROR] Workflow instance NOT found!")
                print(f"  Looking for _id: {workflow_instance_id}")
        except Exception as e:
            print(f"\n[ERROR] Invalid workflow_instance format: {e}")
    else:
        print(f"\n[ERROR] Order has no workflow_instance field!")

except Exception as e:
    print(f"[ERROR] {e}")

print('\n' + '=' * 60)

client.close()

