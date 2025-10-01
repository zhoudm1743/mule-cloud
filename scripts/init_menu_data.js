// MongoDB初始化菜单数据脚本 - 适配Nova-admin前端
// 使用方法: mongosh mule --host 127.0.0.1:27015 -u root -p bgg8384495 --authenticationDatabase admin < scripts/init_menu_data.js

// 切换到mule数据库
db = db.getSiblingDB('mule');

// 清空现有菜单数据
db.menus.deleteMany({});

// 插入初始菜单数据（扁平结构）
const menus = [
  // 仪表盘（目录）
  {
    _id: ObjectId(),
    name: 'dashboard',
    path: '/dashboard',
    title: '仪表盘',
    requiresAuth: true,
    icon: 'icon-park-outline:analysis',
    menuType: 'dir',
    componentPath: null,
    pid: null,
    order: 1,
    status: 1,
    is_deleted: 0,
    created_at: NumberLong(Date.now() / 1000),
    updated_at: NumberLong(Date.now() / 1000),
  },
  // 工作台（页面）
  {
    _id: ObjectId(),
    name: 'dashboard_workbench',
    path: '/dashboard/workbench',
    title: '工作台',
    requiresAuth: true,
    icon: 'icon-park-outline:alarm',
    pinTab: true,
    menuType: 'page',
    componentPath: '/dashboard/workbench/index.vue',
    pid: null, // 这里需要引用上面dashboard的_id
    order: 1,
    status: 1,
    is_deleted: 0,
    created_at: NumberLong(Date.now() / 1000),
    updated_at: NumberLong(Date.now() / 1000),
  },
  // 监控页（页面）
  {
    _id: ObjectId(),
    name: 'dashboard_monitor',
    path: '/dashboard/monitor',
    title: '监控页',
    requiresAuth: true,
    icon: 'icon-park-outline:anchor',
    menuType: 'page',
    componentPath: '/dashboard/monitor/index.vue',
    pid: null, // 这里需要引用上面dashboard的_id
    order: 2,
    status: 1,
    is_deleted: 0,
    created_at: NumberLong(Date.now() / 1000),
    updated_at: NumberLong(Date.now() / 1000),
  },
];

// 先插入父级菜单，获取ID
const dashboardResult = db.menus.insertOne(menus[0]);
const dashboardId = dashboardResult.insertedId.toString();

// 更新子菜单的pid
menus[1].pid = dashboardId;
menus[2].pid = dashboardId;

// 插入子菜单
db.menus.insertMany([menus[1], menus[2]]);

print('✅ 菜单数据初始化完成！');
print('共插入 ' + menus.length + ' 条菜单记录');

