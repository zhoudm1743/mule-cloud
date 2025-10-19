// MongoDB Shell 脚本：同步订单的 status 和 workflow_state
// 使用方法：在MongoDB客户端中执行
//   use mule_ace
//   load('scripts/sync_workflow_state.js')

// 切换到数据库
db = db.getSiblingDB('mule_ace');

print('=' + '='.repeat(59));
print('同步订单状态字段');
print('=' + '='.repeat(59));

// 状态映射
const statusToWorkflowState = {
    0: 'draft',
    1: 'ordered',
    2: 'production',
    3: 'completed',
    4: 'cancelled'
};

// 查找所有订单
const orders = db.orders.find({ is_deleted: 0 }).toArray();

let updatedCount = 0;
let skipCount = 0;

print(`\n找到 ${orders.length} 个订单\n`);

orders.forEach(order => {
    const status = order.status || 0;
    const currentWorkflowState = order.workflow_state || '';
    const expectedWorkflowState = statusToWorkflowState[status] || 'draft';
    
    if (currentWorkflowState !== expectedWorkflowState) {
        print(`订单: ${order.contract_no || '未知'}`);
        print(`  当前: status=${status}, workflow_state=${currentWorkflowState}`);
        print(`  应该: workflow_state=${expectedWorkflowState}`);
        
        // 更新
        const result = db.orders.updateOne(
            { _id: order._id },
            {
                $set: {
                    workflow_state: expectedWorkflowState,
                    updated_at: Math.floor(Date.now() / 1000)
                }
            }
        );
        
        if (result.modifiedCount > 0) {
            print(`  [OK] 已同步\n`);
            updatedCount++;
        } else {
            print(`  [SKIP] 无需更新\n`);
            skipCount++;
        }
    } else {
        skipCount++;
    }
});

print('=' + '='.repeat(59));
print('同步完成！');
print('=' + '='.repeat(59));
print(`  已更新: ${updatedCount} 个订单`);
print(`  已跳过: ${skipCount} 个订单`);
print(`  总计:   ${orders.length} 个订单`);

