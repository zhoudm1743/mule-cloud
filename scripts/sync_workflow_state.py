#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""同步现有订单的 status 和 workflow_state"""

from pymongo import MongoClient
import time

# 连接MongoDB
try:
    client = MongoClient('mongodb://root:bgg8384495@127.0.0.1:27015/')
    db = client['mule_ace']
    
    # 状态映射
    status_to_workflow_state = {
        0: 'draft',
        1: 'ordered',
        2: 'production',
        3: 'completed',
        4: 'cancelled'
    }
    
    print('=' * 60)
    print('同步订单状态字段')
    print('=' * 60)
    
    # 查找所有 status 和 workflow_state 不一致的订单
    orders = list(db.orders.find({'is_deleted': 0}))
    
    updated_count = 0
    skip_count = 0
    
    for order in orders:
        status = order.get('status', 0)
        current_workflow_state = order.get('workflow_state', '')
        expected_workflow_state = status_to_workflow_state.get(status, 'draft')
        
        if current_workflow_state != expected_workflow_state:
            print(f"\n订单: {order.get('contract_no', '未知')}")
            print(f"  当前: status={status}, workflow_state={current_workflow_state}")
            print(f"  应该: workflow_state={expected_workflow_state}")
            
            # 更新
            result = db.orders.update_one(
                {'_id': order['_id']},
                {'$set': {
                    'workflow_state': expected_workflow_state,
                    'updated_at': int(time.time())
                }}
            )
            
            if result.modified_count > 0:
                print(f"  [OK] 已同步")
                updated_count += 1
            else:
                print(f"  [SKIP] 无需更新")
                skip_count += 1
        else:
            skip_count += 1
    
    print('\n' + '=' * 60)
    print('同步完成！')
    print('=' * 60)
    print(f"  已更新: {updated_count} 个订单")
    print(f"  已跳过: {skip_count} 个订单")
    print(f"  总计:   {len(orders)} 个订单")
    
    client.close()
    
except Exception as e:
    print(f"错误: {e}")
    print("\n提示：如果MongoDB连接失败，请检查：")
    print("  1. MongoDB服务是否启动")
    print("  2. 连接地址和端口是否正确")
    print("  3. 用户名密码是否正确")

