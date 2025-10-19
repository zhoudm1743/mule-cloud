#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查工作流定义和实例的 ID 类型"""

from pymongo import MongoClient
from bson import ObjectId

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

system_db = client['mule_system']
tenant_db = client['mule_ace']

print('=' * 60)
print('Checking Workflow Definition ID Type')
print('=' * 60)

# 查看系统库中的 workflow_definition
wf_def = system_db.workflow_definitions.find_one({'code': 'basic_order'})
if wf_def:
    print(f"\n[System DB] Workflow Definition:")
    print(f"  _id: {wf_def['_id']}")
    print(f"  _id type: {type(wf_def['_id'])}")
    print(f"  code: {wf_def.get('code')}")
    
    wf_id = wf_def['_id']
    print(f"\n  _id as string: {str(wf_id)}")
    
    # 尝试不同方式查询
    print(f"\n  Testing queries:")
    
    # 1. 使用原始 ObjectId
    test1 = system_db.workflow_definitions.find_one({'_id': wf_id})
    print(f"    Query with ObjectId: {'SUCCESS' if test1 else 'FAILED'}")
    
    # 2. 使用字符串
    test2 = system_db.workflow_definitions.find_one({'_id': str(wf_id)})
    print(f"    Query with string: {'SUCCESS' if test2 else 'FAILED'}")
    
    # 3. 使用 ObjectId(str)
    try:
        test3 = system_db.workflow_definitions.find_one({'_id': ObjectId(str(wf_id))})
        print(f"    Query with ObjectId(string): {'SUCCESS' if test3 else 'FAILED'}")
    except:
        print(f"    Query with ObjectId(string): FAILED (exception)")

print('\n' + '=' * 60)
print('Checking Workflow Instance ID Type')
print('=' * 60)

# 查看租户库中的 workflow_instance
wf_instance = tenant_db.workflow_instances.find_one()
if wf_instance:
    print(f"\n[Tenant DB] Workflow Instance:")
    print(f"  _id: {wf_instance['_id']}")
    print(f"  _id type: {type(wf_instance['_id'])}")
    print(f"  workflow_id: {wf_instance.get('workflow_id')}")
    print(f"  workflow_id type: {type(wf_instance.get('workflow_id'))}")
    
    wf_id_in_instance = wf_instance.get('workflow_id')
    print(f"\n  Testing if workflow_id matches definition:")
    
    # 尝试用实例中的 workflow_id 查询定义
    test1 = system_db.workflow_definitions.find_one({'_id': wf_id_in_instance})
    print(f"    Direct match: {'SUCCESS' if test1 else 'FAILED'}")
    
    # 如果是字符串，尝试转换为 ObjectId
    if isinstance(wf_id_in_instance, str):
        try:
            test2 = system_db.workflow_definitions.find_one({'_id': ObjectId(wf_id_in_instance)})
            print(f"    With ObjectId conversion: {'SUCCESS' if test2 else 'FAILED'}")
        except:
            print(f"    With ObjectId conversion: FAILED (invalid format)")

print('\n' + '=' * 60)

client.close()

