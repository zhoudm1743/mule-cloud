// MongoDB脚本：初始化基础订单工作流定义
// 使用方法：在MongoDB客户端中运行此脚本
//   mongo mule_ace scripts/init_order_workflow.js
// 或者：
//   mongosh mule_ace < scripts/init_order_workflow.js

// 切换到对应的数据库（根据您的租户代码修改）
db = db.getSiblingDB('mule_ace');

// 删除旧的 order_basic 工作流定义（如果存在）
db.workflow_definitions.deleteMany({ code: "order_basic" });

// 创建基础订单工作流定义
const workflowDef = {
    "_id": ObjectId(),
    "name": "基础订单工作流",
    "code": "order_basic",
    "description": "订单基础流程：草稿 → 已下单 → 生产中 → 已完成",
    "version": 1,
    "is_active": true,
    "states": [
        {
            "id": "draft",
            "code": "draft",
            "name": "草稿",
            "type": "start",
            "color": "#909399",
            "description": "订单草稿状态",
            "position": { "x": 100, "y": 200 }
        },
        {
            "id": "ordered",
            "code": "ordered",
            "name": "已下单",
            "type": "normal",
            "color": "#409EFF",
            "description": "订单已提交",
            "position": { "x": 300, "y": 200 }
        },
        {
            "id": "production",
            "code": "production",
            "name": "生产中",
            "type": "normal",
            "color": "#E6A23C",
            "description": "订单正在生产",
            "position": { "x": 500, "y": 200 }
        },
        {
            "id": "completed",
            "code": "completed",
            "name": "已完成",
            "type": "end",
            "color": "#67C23A",
            "description": "订单已完成",
            "position": { "x": 700, "y": 200 }
        },
        {
            "id": "cancelled",
            "code": "cancelled",
            "name": "已取消",
            "type": "end",
            "color": "#F56C6C",
            "description": "订单已取消",
            "position": { "x": 500, "y": 350 }
        }
    ],
    "transitions": [
        {
            "id": "t1",
            "name": "提交订单",
            "from_state": "draft",
            "to_state": "ordered",
            "event": "submit_order",
            "description": "从草稿提交为正式订单",
            "conditions": [],
            "actions": [
                {
                    "type": "update_field",
                    "field": "status",
                    "value": 1,
                    "description": "更新订单状态为已下单"
                }
            ],
            "require_role": ""
        },
        {
            "id": "t2",
            "name": "开始生产",
            "from_state": "ordered",
            "to_state": "production",
            "event": "start_production",
            "description": "开始生产",
            "conditions": [],
            "actions": [
                {
                    "type": "update_field",
                    "field": "status",
                    "value": 2,
                    "description": "更新订单状态为生产中"
                }
            ],
            "require_role": ""
        },
        {
            "id": "t3",
            "name": "完成订单",
            "from_state": "production",
            "to_state": "completed",
            "event": "complete",
            "description": "完成订单",
            "conditions": [
                {
                    "type": "field",
                    "field": "progress",
                    "operator": "gte",
                    "value": 1.0,
                    "description": "进度必须达到100%"
                }
            ],
            "actions": [
                {
                    "type": "update_field",
                    "field": "status",
                    "value": 3,
                    "description": "更新订单状态为已完成"
                }
            ],
            "require_role": ""
        },
        {
            "id": "t4",
            "name": "取消订单（从草稿）",
            "from_state": "draft",
            "to_state": "cancelled",
            "event": "cancel",
            "description": "取消草稿订单",
            "conditions": [],
            "actions": [
                {
                    "type": "update_field",
                    "field": "status",
                    "value": 4,
                    "description": "更新订单状态为已取消"
                }
            ],
            "require_role": "admin"
        },
        {
            "id": "t5",
            "name": "取消订单（从已下单）",
            "from_state": "ordered",
            "to_state": "cancelled",
            "event": "cancel",
            "description": "取消已下单订单",
            "conditions": [],
            "actions": [
                {
                    "type": "update_field",
                    "field": "status",
                    "value": 4,
                    "description": "更新订单状态为已取消"
                }
            ],
            "require_role": "admin"
        }
    ],
    "metadata": {
        "entity_type": "order",
        "description": "自动生成的基础订单工作流定义"
    },
    "created_at": NumberLong(Date.now() / 1000),
    "updated_at": NumberLong(Date.now() / 1000),
    "created_by": "system",
    "updated_by": "system"
};

// 插入工作流定义
const result = db.workflow_definitions.insertOne(workflowDef);

print("✅ 工作流定义创建成功！");
print("   工作流ID:", result.insertedId);
print("   工作流编码:", workflowDef.code);
print("   工作流名称:", workflowDef.name);
print("\n📊 工作流状态：");
workflowDef.states.forEach(state => {
    print("   -", state.name, "(" + state.code + ")");
});
print("\n🔄 工作流转换：");
workflowDef.transitions.forEach(trans => {
    const fromState = workflowDef.states.find(s => s.code === trans.from_state);
    const toState = workflowDef.states.find(s => s.code === trans.to_state);
    print("   -", trans.name + ":", fromState.name, "→", toState.name);
});

print("\n✨ 现在可以创建订单了，系统会自动使用此工作流！");

