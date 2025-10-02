// 权限体系改造：迁移现有角色数据
// 
// 新的角色体系：
// - super: 系统级超管（tenant_id 为空）
// - tenant_admin: 租户级超管（有 tenant_id）
// - user: 普通用户

db = db.getSiblingDB('mule-cloud');

print('===== 开始迁移角色体系 =====\n');

// 1. 统计当前数据
print('【统计当前数据】');
const totalAdmins = db.admin.countDocuments({ is_deleted: 0 });
const superAdmins = db.admin.countDocuments({ is_deleted: 0, role: 'super' });
const tenantAdmins = db.admin.countDocuments({ is_deleted: 0, tenant_id: { $ne: null, $ne: '' }, role: 'super' });
const systemUsers = db.admin.countDocuments({ is_deleted: 0, tenant_id: { $eq: '' } });
const tenantUsers = db.admin.countDocuments({ is_deleted: 0, tenant_id: { $ne: null, $ne: '' } });

print(`总管理员数: ${totalAdmins}`);
print(`当前 super 角色数: ${superAdmins}`);
print(`其中有租户的 super: ${tenantAdmins}`);
print(`系统级用户数: ${systemUsers}`);
print(`租户级用户数: ${tenantUsers}\n`);

// 2. 迁移租户级超管（将有 tenant_id 的 super 改为 tenant_admin）
print('【迁移租户级超管】');
const tenantSuperResult = db.admin.updateMany(
    {
        is_deleted: 0,
        tenant_id: { $ne: null, $ne: '' },
        role: 'super'
    },
    {
        $set: { role: ['tenant_admin'] }
    }
);
print(`✅ 已将 ${tenantSuperResult.modifiedCount} 个租户级 super 改为 tenant_admin\n`);

// 3. 确保系统级超管的 tenant_id 为空字符串
print('【规范化系统级超管】');
const systemSuperResult = db.admin.updateMany(
    {
        is_deleted: 0,
        $or: [
            { tenant_id: null },
            { tenant_id: '' }
        ],
        role: 'super'
    },
    {
        $set: { 
            tenant_id: '',
            role: ['super']
        }
    }
);
print(`✅ 已规范化 ${systemSuperResult.modifiedCount} 个系统级超管\n`);

// 4. 处理没有角色字段的用户（设置为普通用户）
print('【处理无角色用户】');
const noRoleResult = db.admin.updateMany(
    {
        is_deleted: 0,
        $or: [
            { role: { $exists: false } },
            { role: null },
            { role: [] }
        ]
    },
    {
        $set: { role: ['user'] }
    }
);
print(`✅ 已为 ${noRoleResult.modifiedCount} 个无角色用户设置默认角色 'user'\n`);

// 5. 确保所有 role 字段都是数组格式
print('【规范化角色字段为数组】');
const stringRoleResult = db.admin.updateMany(
    {
        is_deleted: 0,
        role: { $type: 'string' }
    },
    [
        {
            $set: {
                role: { $cond: { if: { $isArray: '$role' }, then: '$role', else: ['$role'] } }
            }
        }
    ]
);
print(`✅ 已规范化角色字段格式\n`);

// 6. 显示迁移后的统计
print('【迁移后统计】');
const afterSystemSuper = db.admin.countDocuments({ is_deleted: 0, tenant_id: '', role: 'super' });
const afterTenantAdmin = db.admin.countDocuments({ is_deleted: 0, tenant_id: { $ne: '' }, role: 'tenant_admin' });
const afterUsers = db.admin.countDocuments({ is_deleted: 0, role: 'user' });

print(`系统级超管（super，无租户）: ${afterSystemSuper}`);
print(`租户级超管（tenant_admin）: ${afterTenantAdmin}`);
print(`普通用户（user）: ${afterUsers}\n`);

// 7. 列出所有超管
print('【系统级超管列表】');
db.admin.find(
    { is_deleted: 0, tenant_id: '', role: 'super' },
    { phone: 1, nickname: 1, role: 1, tenant_id: 1 }
).forEach(admin => {
    print(`  - ${admin.nickname} (${admin.phone}) - 角色: ${admin.role}`);
});

print('\n【租户级超管列表】');
db.admin.find(
    { is_deleted: 0, tenant_id: { $ne: '' }, role: 'tenant_admin' },
    { phone: 1, nickname: 1, role: 1, tenant_id: 1 }
).forEach(admin => {
    print(`  - ${admin.nickname} (${admin.phone}) - 租户: ${admin.tenant_id} - 角色: ${admin.role}`);
});

print('\n===== 迁移完成 =====');
print('\n新的角色体系：');
print('  super: 系统级超管（tenant_id 为空，跨租户权限）');
print('  tenant_admin: 租户级超管（有 tenant_id，本租户所有权限）');
print('  user: 普通用户（通过角色分配权限）');

