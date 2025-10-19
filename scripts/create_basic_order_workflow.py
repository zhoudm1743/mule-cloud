#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""创建基础订单工作流定义"""

from pymongo import MongoClient
from bson import ObjectId
import time

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

print('=' * 60)
print('Creating basic_order Workflow Definition')
print('=' * 60)

# 删除旧的定义（如果存在）
db.workflow_definitions.delete_many({'code': 'basic_order'})

# 创建工作流定义
now = int(time.time())
workflow_id = ObjectId()

workflow_def = {
    '_id': workflow_id,
    'name': '基础订单工作流',
    'code': 'basic_order',
    'description': '订单基础流程：草稿 → 已下单 → 生产中 → 已完成',
    'version': 1,
    'is_active': True,
    'states': [
        {
            'id': 'draft',
            'code': 'draft',
            'name': '草稿',
            'type': 'start',
            'color': '#909399',
            'description': '订单草稿状态',
            'position': {'x': 100, 'y': 200}
        },
        {
            'id': 'ordered',
            'code': 'ordered',
            'name': '已下单',
            'type': 'normal',
            'color': '#409EFF',
            'description': '订单已提交',
            'position': {'x': 300, 'y': 200}
        },
        {
            'id': 'production',
            'code': 'production',
            'name': '生产中',
            'type': 'normal',
            'color': '#E6A23C',
            'description': '订单正在生产',
            'position': {'x': 500, 'y': 200}
        },
        {
            'id': 'completed',
            'code': 'completed',
            'name': '已完成',
            'type': 'end',
            'color': '#67C23A',
            'description': '订单已完成',
            'position': {'x': 700, 'y': 200}
        },
        {
            'id': 'cancelled',
            'code': 'cancelled',
            'name': '已取消',
            'type': 'end',
            'color': '#F56C6C',
            'description': '订单已取消',
            'position': {'x': 500, 'y': 350}
        }
    ],
    'transitions': [
        {
            'id': 't1',
            'name': '提交订单',
            'from_state': 'draft',
            'to_state': 'ordered',
            'event': 'submit_order',
            'description': '从草稿提交为正式订单',
            'conditions': [],
            'actions': [
                {
                    'type': 'update_field',
                    'field': 'status',
                    'value': 1,
                    'description': '更新订单状态为已下单'
                }
            ],
            'require_role': ''
        },
        {
            'id': 't2',
            'name': '开始生产',
            'from_state': 'ordered',
            'to_state': 'production',
            'event': 'start_production',
            'description': '开始生产',
            'conditions': [],
            'actions': [
                {
                    'type': 'update_field',
                    'field': 'status',
                    'value': 2,
                    'description': '更新订单状态为生产中'
                }
            ],
            'require_role': ''
        },
        {
            'id': 't3',
            'name': '完成订单',
            'from_state': 'production',
            'to_state': 'completed',
            'event': 'complete',
            'description': '完成订单',
            'conditions': [],
            'actions': [
                {
                    'type': 'update_field',
                    'field': 'status',
                    'value': 3,
                    'description': '更新订单状态为已完成'
                }
            ],
            'require_role': ''
        },
        {
            'id': 't4',
            'name': '取消订单（从草稿）',
            'from_state': 'draft',
            'to_state': 'cancelled',
            'event': 'cancel',
            'description': '取消草稿订单',
            'conditions': [],
            'actions': [
                {
                    'type': 'update_field',
                    'field': 'status',
                    'value': 4,
                    'description': '更新订单状态为已取消'
                }
            ],
            'require_role': 'admin'
        },
        {
            'id': 't5',
            'name': '取消订单（从已下单）',
            'from_state': 'ordered',
            'to_state': 'cancelled',
            'event': 'cancel',
            'description': '取消已下单订单',
            'conditions': [],
            'actions': [
                {
                    'type': 'update_field',
                    'field': 'status',
                    'value': 4,
                    'description': '更新订单状态为已取消'
                }
            ],
            'require_role': 'admin'
        }
    ],
    'metadata': {
        'entity_type': 'order',
        'description': 'Auto-generated basic order workflow'
    },
    'created_at': now,
    'updated_at': now,
    'created_by': 'system',
    'updated_by': 'system'
}

# 插入数据库
result = db.workflow_definitions.insert_one(workflow_def)

print(f"\n[OK] Workflow definition created!")
print(f"     _id: {workflow_id}")
print(f"     code: {workflow_def['code']}")
print(f"     name: {workflow_def['name']}")
print(f"\nStates:")
for state in workflow_def['states']:
    print(f"  - {state['name']} ({state['code']})")

print(f"\nTransitions:")
for trans in workflow_def['transitions']:
    from_s = next(s['name'] for s in workflow_def['states'] if s['code'] == trans['from_state'])
    to_s = next(s['name'] for s in workflow_def['states'] if s['code'] == trans['to_state'])
    print(f"  - {trans['name']}: {from_s} -> {to_s}")

print('\n' + '=' * 60)
print('Done! Now running workflow instance initialization...')
print('=' * 60)

client.close()

