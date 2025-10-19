#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from pymongo import MongoClient

client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

workflows = list(db.workflow_definitions.find({}, {'_id': 1, 'name': 1, 'code': 1, 'is_active': 1}))

print('=' * 60)
print('Workflow Definitions')
print('=' * 60)

if workflows:
    for wf in workflows:
        active_str = 'ACTIVE' if wf.get('is_active', False) else 'INACTIVE'
        print(f"\n{wf['name']}")
        print(f"  code: {wf['code']}")
        print(f"  status: {active_str}")
        print(f"  _id: {wf['_id']}")
else:
    print('No workflow definitions found!')
    print('\nYou need to create the basic_order workflow definition.')
    print('Please visit the workflow designer to create it.')

print('\n' + '=' * 60)

# 特别检查 basic_order
basic = db.workflow_definitions.find_one({'code': 'basic_order'})
if basic:
    print(f"basic_order exists: {basic.get('is_active', False)}")
else:
    print('basic_order NOT FOUND!')

client.close()

