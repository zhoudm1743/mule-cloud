// MongoDB 初始化认证服务测试用户脚本
// 使用方法: mongosh mongodb://root:bgg8384495@localhost:27015/mule?authSource=admin < scripts/init_auth_users.js

// 切换到 mule 数据库
db = db.getSiblingDB('mule');

// 删除已存在的测试用户（可选）
print("清理旧数据...");
db.admins.deleteMany({ phone: { $in: ["13800138000", "13900139000", "13700137000"] } });

// 创建索引
print("创建索引...");
db.admins.createIndex({ "phone": 1 }, { unique: true });
db.admins.createIndex({ "email": 1 });
db.admins.createIndex({ "status": 1 });

// 获取当前时间戳（秒）
var now = Math.floor(Date.now() / 1000);

// 插入测试用户
print("插入测试用户...");

// 1. 普通用户
db.admins.insertOne({
  phone: "13800138000",
  password: "e10adc3949ba59abbe56e057f20f883e",  // 123456 的 MD5
  nickname: "测试用户",
  email: "test@example.com",
  status: 1,
  role: ["user"],
  avatar: "https://avatar.example.com/test.jpg",
  created_at: NumberLong(now),
  updated_at: NumberLong(now)
});
print("✅ 创建普通用户: 13800138000 / 123456");

// 2. 管理员用户
db.admins.insertOne({
  phone: "13900139000",
  password: "e10adc3949ba59abbe56e057f20f883e",  // 123456 的 MD5
  nickname: "管理员",
  email: "admin@example.com",
  status: 1,
  role: ["admin", "user"],
  avatar: "https://avatar.example.com/admin.jpg",
  created_at: NumberLong(now),
  updated_at: NumberLong(now)
});
print("✅ 创建管理员用户: 13900139000 / 123456");

// 3. 编辑用户
db.admins.insertOne({
  phone: "13700137000",
  password: "e10adc3949ba59abbe56e057f20f883e",  // 123456 的 MD5
  nickname: "编辑员",
  email: "editor@example.com",
  status: 1,
  role: ["editor", "user"],
  avatar: "https://avatar.example.com/editor.jpg",
  created_at: NumberLong(now),
  updated_at: NumberLong(now)
});
print("✅ 创建编辑用户: 13700137000 / 123456");

// 验证插入结果
var count = db.admins.countDocuments({});
print("\n总共创建了 " + count + " 个用户");

// 显示所有用户
print("\n用户列表：");
db.admins.find({}, { password: 0 }).forEach(function(doc) {
  print("- " + doc.phone + " (" + doc.nickname + ") - 角色: " + doc.role.join(", "));
});

print("\n✅ 初始化完成！");
print("默认密码均为: 123456");

