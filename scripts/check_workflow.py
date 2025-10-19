#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""检查MongoDB中的工作流定义"""

from pymongo import MongoClient

# 连接MongoDB
client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
db = client['mule_ace']

# 查询工作流定义
workflows = list(db.workflow_definitions.find({}, {
    '_id': 1,
    'name': 1,
    'code': 1,
    'is_active': 1
}))

print('=' * 60)
print('工作流定义列表 (mule_ace 数据库)')
print('=' * 60)

if workflows:
    for wf in workflows:
        status = '✅ 激活' if wf.get('is_active', False) else '❌ 未激活'
        print(f"  - {wf['name']}")
        print(f"    编码: {wf['code']}")
        print(f"    状态: {status}")
        print(f"    ID: {wf['_id']}")
        print()
    print(f'总共 {len(workflows)} 个工作流定义')
else:
    print('  ❌ 没有找到任何工作流定义')

print('=' * 60)

# 检查是否有 order_basic
order_basic = db.workflow_definitions.find_one({'code': 'order_basic'})
basic_order = db.workflow_definitions.find_one({'code': 'basic_order'})

print('\n特定检查:')
print(f"  - order_basic (代码中使用): {'✅ 存在' if order_basic else '❌ 不存在'}")
print(f"  - basic_order (截图显示): {'✅ 存在' if basic_order else '❌ 不存在'}")

if basic_order and not order_basic:
    print('\n⚠️  问题发现:')
    print('  代码中使用 "order_basic"，但数据库中是 "basic_order"')
    print('  需要修改代码或重命名数据库中的工作流编码')

client.close()

