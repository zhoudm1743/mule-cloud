// 初始化菜单权限
// 为所有现有菜单添加默认的基础 CRUD 权限

db = db.getSiblingDB('mule-cloud'); // 切换到目标数据库

// 基础 CRUD 权限（适用于大部分菜单）
const basicPermissions = [
  { action: 'read', label: '查看', is_basic: true },
  { action: 'create', label: '创建', is_basic: true },
  { action: 'update', label: '修改', is_basic: true },
  { action: 'delete', label: '删除', is_basic: true }
];

// 只读权限（适用于仪表盘等展示类页面）
const readOnlyPermissions = [
  { action: 'read', label: '查看', is_basic: true }
];

// 菜单类型为 'page' 的添加权限，'dir' 目录不添加权限
const result = db.menus.updateMany(
  {
    menuType: 'page',
    available_permissions: { $exists: false }
  },
  {
    $set: { 
      available_permissions: basicPermissions 
    }
  }
);

print(`✅ 更新了 ${result.modifiedCount} 个菜单的权限配置`);

// 特殊处理：仪表盘只需要查看权限
db.menus.updateOne(
  { name: 'dashboard' },
  {
    $set: {
      available_permissions: readOnlyPermissions
    }
  }
);

print('✅ 仪表盘权限已设置为只读');

// 示例：如果有财务管理菜单，添加自定义业务权限
// db.menus.updateOne(
//   { name: 'finance' },
//   {
//     $set: {
//       available_permissions: [
//         { action: 'read', label: '查看', is_basic: true },
//         { action: 'create', label: '创建', is_basic: true },
//         { action: 'update', label: '修改', is_basic: true },
//         { action: 'delete', label: '删除', is_basic: true },
//         { action: 'pending', label: '挂账', is_basic: false, description: '将订单挂账延期支付' },
//         { action: 'verify', label: '核销', is_basic: false, description: '核销已挂账订单' },
//         { action: 'reverse', label: '冲账', is_basic: false, description: '冲销错误账目' },
//         { action: 'audit', label: '审核', is_basic: false, description: '财务审核' },
//         { action: 'export', label: '导出', is_basic: false }
//       ]
//     }
//   }
// );

print('✅ 菜单权限初始化完成！');
print('');
print('📝 提示：');
print('  - 所有 page 类型的菜单已添加基础 CRUD 权限');
print('  - dashboard 设置为只读权限');
print('  - 如需自定义业务权限，请手动修改对应菜单的 available_permissions 字段');

