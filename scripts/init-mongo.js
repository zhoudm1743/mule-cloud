// MongoDB初始化脚本

// 切换到mule_cloud数据库
db = db.getSiblingDB('mule_cloud');

// 创建应用用户
db.createUser({
  user: 'mule_cloud_user',
  pwd: 'mule_cloud_password',
  roles: [
    {
      role: 'readWrite',
      db: 'mule_cloud'
    }
  ]
});

// 创建用户集合并设置索引
db.createCollection('users');
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "status": 1 });
db.users.createIndex({ "created_at": -1 });

// 创建角色集合并设置索引
db.createCollection('roles');
db.roles.createIndex({ "code": 1 }, { unique: true });
db.roles.createIndex({ "status": 1 });

// 创建权限集合并设置索引
db.createCollection('permissions');
db.permissions.createIndex({ "code": 1 }, { unique: true });
db.permissions.createIndex({ "module": 1 });

// 创建客户集合并设置索引
db.createCollection('customers');
db.customers.createIndex({ "customer_no": 1 }, { unique: true });
db.customers.createIndex({ "company_name": 1 });
db.customers.createIndex({ "customer_type": 1 });
db.customers.createIndex({ "status": 1 });

// 创建业务员集合并设置索引
db.createCollection('salespersons');
db.salespersons.createIndex({ "user_id": 1 }, { unique: true });
db.salespersons.createIndex({ "status": 1 });

// 创建款式集合并设置索引
db.createCollection('styles');
db.styles.createIndex({ "style_no": 1 }, { unique: true });
db.styles.createIndex({ "category": 1 });
db.styles.createIndex({ "season": 1, "year": 1 });
db.styles.createIndex({ "status": 1 });

// 创建订单集合并设置索引
db.createCollection('orders');
db.orders.createIndex({ "order_no": 1 }, { unique: true });
db.orders.createIndex({ "customer_id": 1, "created_at": -1 });
db.orders.createIndex({ "status": 1, "created_at": -1 });
db.orders.createIndex({ "salesperson_id": 1 });
db.orders.createIndex({ "delivery_date": 1 });

// 创建生产计划集合并设置索引
db.createCollection('production_plans');
db.production_plans.createIndex({ "plan_no": 1 }, { unique: true });
db.production_plans.createIndex({ "order_id": 1 });
db.production_plans.createIndex({ "status": 1 });
db.production_plans.createIndex({ "start_date": 1, "end_date": 1 });

// 创建裁剪任务集合并设置索引
db.createCollection('cutting_tasks');
db.cutting_tasks.createIndex({ "task_no": 1 }, { unique: true });
db.cutting_tasks.createIndex({ "order_id": 1 });
db.cutting_tasks.createIndex({ "assigned_to": 1 });
db.cutting_tasks.createIndex({ "status": 1 });
db.cutting_tasks.createIndex({ "scheduled_date": 1 });

// 创建生产进度集合并设置索引
db.createCollection('production_progress');
db.production_progress.createIndex({ "order_id": 1, "process_id": 1 }, { unique: true });
db.production_progress.createIndex({ "order_id": 1, "updated_at": -1 });
db.production_progress.createIndex({ "status": 1 });

// 创建工作报告集合并设置索引
db.createCollection('work_reports');
db.work_reports.createIndex({ "worker_id": 1, "date": -1 });
db.work_reports.createIndex({ "order_id": 1, "process_id": 1 });
db.work_reports.createIndex({ "status": 1 });
db.work_reports.createIndex({ "reported_by": 1 });

// 创建工人集合并设置索引
db.createCollection('workers');
db.workers.createIndex({ "worker_no": 1 }, { unique: true });
db.workers.createIndex({ "user_id": 1 }, { unique: true });
db.workers.createIndex({ "status": 1 });
db.workers.createIndex({ "department": 1 });

// 创建工时记录集合并设置索引
db.createCollection('timesheets');
db.timesheets.createIndex({ "worker_id": 1, "date": -1 });
db.timesheets.createIndex({ "order_id": 1 });

// 创建工资单集合并设置索引
db.createCollection('payrolls');
db.payrolls.createIndex({ "worker_id": 1, "pay_period": -1 });
db.payrolls.createIndex({ "pay_period": 1 });

// 创建工序集合并设置索引
db.createCollection('processes');
db.processes.createIndex({ "process_no": 1 }, { unique: true });
db.processes.createIndex({ "category": 1 });
db.processes.createIndex({ "status": 1 });

// 创建尺码集合并设置索引
db.createCollection('sizes');
db.sizes.createIndex({ "size_no": 1 }, { unique: true });
db.sizes.createIndex({ "category": 1 });
db.sizes.createIndex({ "sort_order": 1 });

// 创建颜色集合并设置索引
db.createCollection('colors');
db.colors.createIndex({ "color_no": 1 }, { unique: true });
db.colors.createIndex({ "sort_order": 1 });

// 插入初始数据

// 插入默认角色
db.roles.insertMany([
  {
    name: "系统管理员",
    code: "admin",
    description: "系统管理员，拥有所有权限",
    permissions: ["*"],
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "生产经理",
    code: "production_manager",
    description: "生产经理，管理生产相关功能",
    permissions: ["production.*", "order.view", "worker.*"],
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "业务员",
    code: "salesperson",
    description: "业务员，管理客户和订单",
    permissions: ["order.*", "customer.*", "style.view"],
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    name: "工人",
    code: "worker",
    description: "生产工人，上报工作进度",
    permissions: ["work_report.*", "timesheet.*"],
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

// 插入默认管理员用户
db.users.insertOne({
  username: "admin",
  email: "admin@mulecloud.com",
  password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
  real_name: "系统管理员",
  status: 1,
  role_ids: ["admin"],
  created_at: new Date(),
  updated_at: new Date(),
  version: 1
});

// 插入基础工序数据
db.processes.insertMany([
  {
    process_no: "P001",
    process_name: "裁剪",
    category: "裁剪",
    description: "面料裁剪工序",
    standard_rate: 0.5,
    rate_unit: "件",
    estimated_time: 0.1,
    difficulty_level: 2,
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    process_no: "P002",
    process_name: "缝制",
    category: "缝制",
    description: "服装缝制工序",
    standard_rate: 2.0,
    rate_unit: "件",
    estimated_time: 0.5,
    difficulty_level: 3,
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    process_no: "P003",
    process_name: "锁边",
    category: "缝制",
    description: "锁边工序",
    standard_rate: 0.3,
    rate_unit: "件",
    estimated_time: 0.08,
    difficulty_level: 1,
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    process_no: "P004",
    process_name: "整烫",
    category: "整理",
    description: "服装整烫工序",
    standard_rate: 0.8,
    rate_unit: "件",
    estimated_time: 0.15,
    difficulty_level: 2,
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    process_no: "P005",
    process_name: "包装",
    category: "整理",
    description: "成品包装工序",
    standard_rate: 0.2,
    rate_unit: "件",
    estimated_time: 0.05,
    difficulty_level: 1,
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
  }
]);

// 插入基础尺码数据
db.sizes.insertMany([
  { size_no: "XS", size_name: "XS码", category: "成人", sort_order: 1, status: 1, created_at: new Date(), updated_at: new Date() },
  { size_no: "S", size_name: "S码", category: "成人", sort_order: 2, status: 1, created_at: new Date(), updated_at: new Date() },
  { size_no: "M", size_name: "M码", category: "成人", sort_order: 3, status: 1, created_at: new Date(), updated_at: new Date() },
  { size_no: "L", size_name: "L码", category: "成人", sort_order: 4, status: 1, created_at: new Date(), updated_at: new Date() },
  { size_no: "XL", size_name: "XL码", category: "成人", sort_order: 5, status: 1, created_at: new Date(), updated_at: new Date() },
  { size_no: "XXL", size_name: "XXL码", category: "成人", sort_order: 6, status: 1, created_at: new Date(), updated_at: new Date() }
]);

// 插入基础颜色数据
db.colors.insertMany([
  { color_no: "BLK", color_name: "黑色", hex_code: "#000000", sort_order: 1, status: 1, created_at: new Date(), updated_at: new Date() },
  { color_no: "WHT", color_name: "白色", hex_code: "#FFFFFF", sort_order: 2, status: 1, created_at: new Date(), updated_at: new Date() },
  { color_no: "RED", color_name: "红色", hex_code: "#FF0000", sort_order: 3, status: 1, created_at: new Date(), updated_at: new Date() },
  { color_no: "BLU", color_name: "蓝色", hex_code: "#0000FF", sort_order: 4, status: 1, created_at: new Date(), updated_at: new Date() },
  { color_no: "GRN", color_name: "绿色", hex_code: "#00FF00", sort_order: 5, status: 1, created_at: new Date(), updated_at: new Date() },
  { color_no: "YEL", color_name: "黄色", hex_code: "#FFFF00", sort_order: 6, status: 1, created_at: new Date(), updated_at: new Date() }
]);

print("MongoDB初始化完成!");
print("默认管理员账号: admin / password");
print("数据库: mule_cloud");
print("应用用户: mule_cloud_user");
