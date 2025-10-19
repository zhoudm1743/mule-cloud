#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查订单的工作流状态"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

# 查询特定订单
order_id = '202510183528'
order = db.orders.find_one({'contract_no': order_id})

if order:
    print('=' * 60)
    print(f'订单详情: {order_id}')
    print('=' * 60)
    print(f"  _id: {order.get('_id')}")
    print(f"  status: {order.get('status')} (0=草稿,1=已下单,2=生产中,3=已完成,4=已取消)")
    print(f"  workflow_state: {order.get('workflow_state', '(空)')}")
    print(f"  workflow_code: {order.get('workflow_code', '(空)')}")
    print(f"  workflow_instance: {order.get('workflow_instance', '(空)')}")
    print(f"  progress: {order.get('progress', 0)}")
    
    # 检查工作流实例
    workflow_instance_id = order.get('workflow_instance')
    if workflow_instance_id:
        instance = db.workflow_instances.find_one({'_id': workflow_instance_id})
        if instance:
            print(f"\n✅ 工作流实例存在:")
            print(f"  当前状态: {instance.get('current_state')}")
            print(f"  工作流ID: {instance.get('workflow_id')}")
            print(f"  历史记录数: {len(instance.get('history', []))}")
        else:
            print(f"\n❌ 工作流实例不存在！(ID: {workflow_instance_id})")
    else:
        print(f"\n⚠️ 订单没有关联工作流实例！")
        print(f"   需要为订单初始化工作流实例")
        
        # 检查是否有工作流定义
        workflow_def = db.workflow_definitions.find_one({'code': 'basic_order', 'is_active': True})
        if workflow_def:
            print(f"\n✅ 找到工作流定义: {workflow_def.get('name')}")
            print(f"   工作流ID: {workflow_def.get('_id')}")
        else:
            print(f"\n❌ 找不到激活的 basic_order 工作流定义！")
else:
    print(f'❌ 订单不存在: {order_id}')

client.close()

