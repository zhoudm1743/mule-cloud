// 将用户设置为超级管理员
db = db.getSiblingDB('mule');

print('设置用户为超级管理员...');
const result = db.admin.updateOne(
    { phone: '13800138000' },
    { 
        $set: { 
            is_super: true,
            updated_at: Math.floor(Date.now() / 1000)
        } 
    }
);

if (result.modifiedCount > 0) {
    print('✅ 成功设置为超级管理员');
    const user = db.admin.findOne({ phone: '13800138000' });
    print('用户信息:');
    print('  - ID: ' + user._id);
    print('  - 手机号: ' + user.phone);
    print('  - 昵称: ' + user.nickname);
    print('  - 超级管理员: ' + user.is_super);
    print('  - 角色: ' + JSON.stringify(user.roles));
} else {
    print('❌ 未找到用户或设置失败');
}

