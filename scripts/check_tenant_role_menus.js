const { MongoClient } = require('mongodb');

async function checkTenantRoleMenus() {
    const client = new MongoClient('mongodb://localhost:27017/');
    
    try {
        await client.connect();
        console.log('✅ 已连接到 MongoDB');
        
        // 查询租户数据库
        const tenantDB = 'mule_68dda6cd04ba0d6c8dda4b7a';
        const db = client.db(tenantDB);
        
        console.log(`\n📊 数据库: ${tenantDB}`);
        
        // 查询租户的默认角色
        const role = await db.collection('role').findOne({ 
            code: 'tenant_admin',
            is_deleted: 0 
        });
        
        if (role) {
            console.log('\n🔑 租户管理员角色:');
            console.log('  ID:', role._id);
            console.log('  名称:', role.name);
            console.log('  代码:', role.code);
            console.log('\n📋 分配的菜单:');
            if (role.menus && role.menus.length > 0) {
                role.menus.forEach((menu, index) => {
                    console.log(`  ${index + 1}. ${menu}`);
                });
            } else {
                console.log('  (无)');
            }
        } else {
            console.log('\n❌ 未找到租户管理员角色');
        }
        
    } catch (error) {
        console.error('❌ 错误:', error.message);
    } finally {
        await client.close();
    }
}

checkTenantRoleMenus();

