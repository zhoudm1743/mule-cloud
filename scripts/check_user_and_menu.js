// 连接到 mule 数据库
db = db.getSiblingDB('mule');

print('\n========== 检查用户信息 ==========');
const user = db.admin.findOne({ phone: '13800138000' });
if (user) {
    print('用户ID: ' + user._id);
    print('手机号: ' + user.phone);
    print('昵称: ' + user.nickname);
    print('是否超级管理员: ' + user.is_super);
    print('角色: ' + JSON.stringify(user.roles));
    print('状态: ' + user.status);
} else {
    print('❌ 用户不存在');
}

print('\n========== 检查菜单数据 ==========');
const menuCount = db.menu.countDocuments({ is_deleted: 0 });
print('菜单总数: ' + menuCount);

if (menuCount === 0) {
    print('\n⚠️  数据库中没有菜单数据！');
    print('是否要创建测试菜单数据？(y/n)');
}

print('\n========== 菜单列表 ==========');
db.menu.find({ is_deleted: 0 }, { name: 1, title: 1, path: 1 }).forEach(menu => {
    print('- ' + menu.title + ' (' + menu.name + ') - ' + menu.path);
});

