// 修复菜单的 componentPath，将 /system/ 改为 /setting/
db = db.getSiblingDB('mule');

print('========== 修复菜单组件路径 ==========');

// 更新所有包含 /system/ 的 componentPath
const result = db.menu.updateMany(
    { 
        componentPath: { $regex: '^/system/' },
        is_deleted: 0 
    },
    [
        { 
            $set: { 
                componentPath: { 
                    $replaceOne: { 
                        input: "$componentPath", 
                        find: "/system/", 
                        replacement: "/setting/" 
                    } 
                },
                updated_at: Math.floor(Date.now() / 1000)
            } 
        }
    ]
);

print('修改了 ' + result.modifiedCount + ' 条记录');

print('\n========== 更新后的菜单 ==========');
db.menu.find(
    { is_deleted: 0 },
    { title: 1, path: 1, componentPath: 1 }
).sort({ order: 1 }).forEach(menu => {
    print(menu.title + ' -> ' + menu.componentPath);
});

