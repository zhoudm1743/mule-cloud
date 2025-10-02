// 检查 admin 集合中 role/roles 字段的实际情况
db = db.getSiblingDB('mule-cloud');

print('===== 检查 admin 数据 =====');
const admin = db.admin.findOne({ phone: '17858361617' });

if (admin) {
    print('\n找到用户:', admin.nickname);
    print('ID:', admin._id);
    print('Phone:', admin.phone);
    
    // 检查 role 字段（单数）
    if (admin.role !== undefined) {
        print('\n✅ 数据库中使用的是 role（单数）字段');
        print('role 值:', JSON.stringify(admin.role));
        print('role 类型:', typeof admin.role);
        print('是否数组:', Array.isArray(admin.role));
    }
    
    // 检查 roles 字段（复数）
    if (admin.roles !== undefined) {
        print('\n✅ 数据库中使用的是 roles（复数）字段');
        print('roles 值:', JSON.stringify(admin.roles));
        print('roles 类型:', typeof admin.roles);
        print('是否数组:', Array.isArray(admin.roles));
    }
    
    // 如果都没有
    if (admin.role === undefined && admin.roles === undefined) {
        print('\n⚠️  该用户没有 role 或 roles 字段');
    }
    
    // 打印完整对象（用于调试）
    print('\n===== 完整对象 =====');
    printjson(admin);
} else {
    print('❌ 未找到该用户');
}

