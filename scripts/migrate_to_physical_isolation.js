// MongoDB数据迁移脚本 - 从逻辑隔离迁移到物理隔离
// 用法：mongo --host 127.0.0.1 --port 27015 -u root -p bgg8384495 --authenticationDatabase admin mule scripts/migrate_to_physical_isolation.js

const sourceDB = db.getSiblingDB("mule");
const systemDB = db.getSiblingDB("tenant_system");

print("╔════════════════════════════════════════════════════════════════╗");
print("║        MongoDB 数据迁移 - 逻辑隔离 → 物理隔离                  ║");
print("╚════════════════════════════════════════════════════════════════╝");
print("");
print("源数据库: mule");
print("目标系统库: tenant_system");
print("开始时间: " + new Date().toISOString());
print("");

// 统计信息
let stats = {
    tenantCount: 0,
    superAdminCount: 0,
    tenantAdminCount: 0,
    regularUserCount: 0,
    mappingCount: 0,
    errors: []
};

try {
    // ============================================================
    // 1. 迁移租户数据到系统库
    // ============================================================
    print("[1/5] ════ 迁移租户数据到系统库 ════");
    const tenants = sourceDB.tenant.find({ is_deleted: 0 }).toArray();
    stats.tenantCount = tenants.length;
    print(`  找到 ${tenants.length} 个活跃租户`);

    if (tenants.length > 0) {
        systemDB.tenant.insertMany(tenants);
        print(`  ✅ 租户数据迁移完成`);
    } else {
        print(`  ⚠️  未找到租户数据`);
    }

    // ============================================================
    // 2. 迁移系统超管到系统库
    // ============================================================
    print("\n[2/5] ════ 迁移系统超管到系统库 ════");
    const superAdmins = sourceDB.admin.find({ 
        $or: [
            { tenant_id: { $eq: "" } },
            { tenant_id: { $exists: false } },
            { tenant_id: null }
        ],
        role: { $in: ["super", ["super"]] },
        is_deleted: 0 
    }).toArray();

    stats.superAdminCount = superAdmins.length;
    print(`  找到 ${superAdmins.length} 个系统超管`);

    // 删除 tenant_id 字段
    superAdmins.forEach(admin => {
        delete admin.tenant_id;
        // 确保 role 是数组格式
        if (typeof admin.role === 'string') {
            admin.role = [admin.role];
        }
    });

    if (superAdmins.length > 0) {
        systemDB.admin.insertMany(superAdmins);
        print(`  ✅ 系统超管迁移完成`);
    } else {
        print(`  ⚠️  未找到系统超管`);
    }

    // ============================================================
    // 3. 为每个租户创建独立数据库
    // ============================================================
    print("\n[3/5] ════ 为每个租户创建独立数据库 ════");
    print(`  共需处理 ${tenants.length} 个租户\n`);

    tenants.forEach((tenant, index) => {
        const tenantID = tenant._id.toString();
        const tenantDBName = `tenant_${tenantID}`;
        const tenantDB = db.getSiblingDB(tenantDBName);
        
        print(`  ┌─ [${index + 1}/${tenants.length}] ${tenant.name} (${tenant.code})`);
        print(`  │  数据库: ${tenantDBName}`);
        
        // 需要迁移的集合
        const collections = [
            "admin", "role", "menu", "basic", 
            "color", "customer", "order_type", 
            "procedure", "salesman", "size"
        ];
        
        let tenantStats = {
            total: 0,
            details: {}
        };

        collections.forEach(collName => {
            try {
                const filter = { tenant_id: tenantID, is_deleted: 0 };
                const docs = sourceDB[collName].find(filter).toArray();
                
                if (docs.length > 0) {
                    // 删除 tenant_id 字段
                    docs.forEach(doc => {
                        delete doc.tenant_id;
                        
                        // 特殊处理：确保 admin 的 role 是数组
                        if (collName === 'admin' && doc.role) {
                            if (typeof doc.role === 'string') {
                                doc.role = [doc.role];
                            }
                        }
                    });
                    
                    tenantDB[collName].insertMany(docs);
                    tenantStats.details[collName] = docs.length;
                    tenantStats.total += docs.length;
                    
                    print(`  │  ✓ ${collName}: ${docs.length} 条`);
                    
                    // 统计租户管理员和普通用户
                    if (collName === 'admin') {
                        docs.forEach(admin => {
                            if (admin.role && (admin.role.includes('tenant_admin') || admin.role === 'tenant_admin')) {
                                stats.tenantAdminCount++;
                            } else {
                                stats.regularUserCount++;
                            }
                        });
                    }
                }
            } catch (e) {
                print(`  │  ✗ ${collName}: 失败 - ${e.message}`);
                stats.errors.push(`租户 ${tenant.name} - ${collName}: ${e.message}`);
            }
        });
        
        // 创建索引
        try {
            tenantDB.admin.createIndex({ phone: 1 }, { unique: true, sparse: true });
            tenantDB.admin.createIndex({ email: 1 }, { sparse: true });
            tenantDB.admin.createIndex({ is_deleted: 1 });
            tenantDB.role.createIndex({ code: 1 }, { unique: true, sparse: true });
            tenantDB.role.createIndex({ is_deleted: 1 });
            tenantDB.menu.createIndex({ name: 1 });
            tenantDB.menu.createIndex({ is_deleted: 1 });
            print(`  │  ✓ 索引创建完成`);
        } catch (e) {
            print(`  │  ⚠️  索引创建失败: ${e.message}`);
        }
        
        print(`  └─ 完成，共迁移 ${tenantStats.total} 条记录\n`);
    });

    // ============================================================
    // 4. 创建手机号映射表（性能优化）
    // ============================================================
    print("[4/5] ════ 创建手机号映射表 ════");
    const allAdmins = sourceDB.admin.find({ 
        tenant_id: { $ne: "", $exists: true, $ne: null },
        phone: { $exists: true },
        is_deleted: 0 
    }).toArray();

    print(`  找到 ${allAdmins.length} 个租户用户`);

    if (allAdmins.length > 0) {
        const mappings = [];
        const phoneSet = new Set();
        
        allAdmins.forEach(admin => {
            if (admin.phone && !phoneSet.has(admin.phone)) {
                mappings.push({
                    _id: admin.phone,
                    tenant_id: admin.tenant_id,
                    created_at: Date.now()
                });
                phoneSet.add(admin.phone);
            }
        });
        
        stats.mappingCount = mappings.length;
        
        if (mappings.length > 0) {
            try {
                systemDB.phone_to_tenant.insertMany(mappings, { ordered: false });
                print(`  ✅ 创建 ${mappings.length} 条映射记录`);
            } catch (e) {
                // 忽略重复键错误
                if (!e.message.includes('duplicate key')) {
                    print(`  ⚠️  部分映射创建失败: ${e.message}`);
                }
            }
        }
    } else {
        print(`  ⚠️  未找到租户用户数据`);
    }

    // ============================================================
    // 5. 验证迁移结果
    // ============================================================
    print("\n[5/5] ════ 验证迁移结果 ════");
    let allSuccess = true;

    // 5.1 验证租户数量
    const sourceTenantCount = sourceDB.tenant.countDocuments({ is_deleted: 0 });
    const systemTenantCount = systemDB.tenant.countDocuments({ is_deleted: 0 });
    print(`\n  租户数量验证:`);
    print(`    原库: ${sourceTenantCount}`);
    print(`    系统库: ${systemTenantCount}`);
    if (sourceTenantCount !== systemTenantCount) {
        print("    ❌ 不匹配！");
        allSuccess = false;
    } else {
        print("    ✅ 匹配");
    }

    // 5.2 验证系统超管
    const sourceSuperCount = sourceDB.admin.countDocuments({ 
        $or: [
            { tenant_id: "" },
            { tenant_id: { $exists: false } },
            { tenant_id: null }
        ],
        role: { $in: ["super", ["super"]] },
        is_deleted: 0 
    });
    const systemSuperCount = systemDB.admin.countDocuments({ 
        role: { $in: ["super", ["super"]] },
        is_deleted: 0 
    });
    print(`\n  系统超管验证:`);
    print(`    原库: ${sourceSuperCount}`);
    print(`    系统库: ${systemSuperCount}`);
    if (sourceSuperCount !== systemSuperCount) {
        print("    ❌ 不匹配！");
        allSuccess = false;
    } else {
        print("    ✅ 匹配");
    }

    // 5.3 验证每个租户的数据
    print(`\n  租户数据验证:`);
    let tenantValidationErrors = 0;
    
    tenants.forEach(tenant => {
        const tenantID = tenant._id.toString();
        const tenantDBName = `tenant_${tenantID}`;
        const tenantDB = db.getSiblingDB(tenantDBName);
        
        const collections = ["admin", "role", "menu", "basic"];
        let hasError = false;
        
        collections.forEach(collName => {
            const sourceCount = sourceDB[collName].countDocuments({ 
                tenant_id: tenantID, 
                is_deleted: 0 
            });
            const tenantCount = tenantDB[collName].countDocuments({ is_deleted: 0 });
            
            if (sourceCount !== tenantCount) {
                if (!hasError) {
                    print(`    ❌ ${tenant.name}:`);
                    hasError = true;
                }
                print(`       ${collName}: 原库 ${sourceCount} vs 租户库 ${tenantCount}`);
                tenantValidationErrors++;
                allSuccess = false;
            }
        });
    });
    
    if (tenantValidationErrors === 0) {
        print(`    ✅ 所有租户数据验证通过`);
    }

    // ============================================================
    // 迁移总结
    // ============================================================
    print("\n╔════════════════════════════════════════════════════════════════╗");
    print("║                        迁移总结                                ║");
    print("╚════════════════════════════════════════════════════════════════╝");
    print("");
    print(`  租户数量:          ${stats.tenantCount}`);
    print(`  系统超管:          ${stats.superAdminCount}`);
    print(`  租户管理员:        ${stats.tenantAdminCount}`);
    print(`  普通用户:          ${stats.regularUserCount}`);
    print(`  手机号映射:        ${stats.mappingCount}`);
    print(`  错误数量:          ${stats.errors.length}`);
    print("");
    
    if (stats.errors.length > 0) {
        print("  错误详情:");
        stats.errors.forEach(err => {
            print(`    • ${err}`);
        });
        print("");
    }

    if (allSuccess && stats.errors.length === 0) {
        print("╔════════════════════════════════════════════════════════════════╗");
        print("║  ✅ 迁移成功！所有数据验证通过！                               ║");
        print("╚════════════════════════════════════════════════════════════════╝");
        print("");
        print("  ⚠️  下一步操作:");
        print("  1. 启动新版本服务，确保功能正常");
        print("  2. 测试各租户登录和数据访问");
        print("  3. 确认新系统稳定运行至少1周");
        print("  4. 备份原数据库后再删除");
        print("");
        print("  删除原数据库命令:");
        print("    use mule");
        print("    db.dropDatabase()");
    } else {
        print("╔════════════════════════════════════════════════════════════════╗");
        print("║  ❌ 迁移过程中发现问题，请检查！                               ║");
        print("╚════════════════════════════════════════════════════════════════╝");
        print("");
        print("  建议:");
        print("  1. 检查上述错误详情");
        print("  2. 修复问题后重新执行迁移");
        print("  3. 或手动修复不匹配的数据");
    }

} catch (e) {
    print("\n╔════════════════════════════════════════════════════════════════╗");
    print("║  ❌ 迁移失败！                                                 ║");
    print("╚════════════════════════════════════════════════════════════════╝");
    print("");
    print(`错误: ${e.message}`);
    print("");
    print(`堆栈: ${e.stack}`);
}

print("");
print("完成时间: " + new Date().toISOString());
print("");

