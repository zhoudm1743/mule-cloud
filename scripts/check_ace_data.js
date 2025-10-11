// 检查 ace 租户的数据
const { MongoClient } = require('mongodb');

async function main() {
  const uri = 'mongodb://admin:password@localhost:27017';
  const client = new MongoClient(uri);

  try {
    await client.connect();
    console.log('✅ MongoDB 连接成功\n');

    // 检查租户信息
    const systemDB = client.db('tenant_system');
    const tenant = await systemDB.collection('tenant').findOne({ code: 'ace' });
    
    if (tenant) {
      console.log('✅ 租户信息:');
      console.log(`   ID: ${tenant._id}`);
      console.log(`   Code: ${tenant.code}`);
      console.log(`   Name: ${tenant.name}\n`);
    } else {
      console.log('❌ 租户 ace 不存在\n');
      return;
    }

    // 检查 ace 租户的数据库
    const aceDB = client.db('mule_ace');
    
    // 检查 department 集合
    const deptCount = await aceDB.collection('department').countDocuments();
    console.log(`📊 mule_ace.department 记录数: ${deptCount}`);
    
    if (deptCount > 0) {
      const departments = await aceDB.collection('department').find({}).limit(5).toArray();
      console.log('\n前 5 条记录:');
      departments.forEach((dept, index) => {
        console.log(`  ${index + 1}. ${dept.name} (${dept._id})`);
      });
    }

    // 检查 admin 集合
    const adminCount = await aceDB.collection('admin').countDocuments();
    console.log(`\n📊 mule_ace.admin 记录数: ${adminCount}`);

    // 检查 role 集合
    const roleCount = await aceDB.collection('role').countDocuments();
    console.log(`📊 mule_ace.role 记录数: ${roleCount}`);

  } catch (error) {
    console.error('❌ 错误:', error.message);
  } finally {
    await client.close();
  }
}

main();


