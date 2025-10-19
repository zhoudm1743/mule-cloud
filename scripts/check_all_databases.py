#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查所有数据库中的工作流定义"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')

print('=' * 60)
print('Checking All Databases for workflow_definitions')
print('=' * 60)

# 列出所有数据库
db_names = client.list_database_names()
print(f"\nFound databases: {db_names}")

for db_name in db_names:
    if db_name in ['admin', 'config', 'local']:
        continue
        
    db = client[db_name]
    
    # 检查是否有 workflow_definitions 集合
    if 'workflow_definitions' in db.list_collection_names():
        workflows = list(db.workflow_definitions.find({}, {'_id': 1, 'name': 1, 'code': 1, 'is_active': 1}))
        
        if workflows:
            print(f"\n{'=' * 60}")
            print(f"Database: {db_name}")
            print(f"{'=' * 60}")
            
            for wf in workflows:
                active = 'ACTIVE' if wf.get('is_active', False) else 'INACTIVE'
                print(f"\n  {wf['name']}")
                print(f"    code: {wf['code']}")
                print(f"    status: {active}")
                print(f"    _id: {wf['_id']}")

print('\n' + '=' * 60)
print('Check complete!')
print('=' * 60)

client.close()

