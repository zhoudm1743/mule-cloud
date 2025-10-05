#!/usr/bin/env python3
"""
批量修复role.go repository中的方法
在每个方法开始处添加 collection := r.getCollection(ctx)
并将所有 r.collection 替换为 collection
"""

import re

def fix_role_repository():
    file_path = r"K:\Git\mule-cloud\internal\repository\role.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 替换所有 r.collection. 为 collection.
    content = content.replace('r.collection.', 'collection.')
    
    # 2. 在每个方法开始处添加 collection := r.getCollection(ctx)
    # 匹配模式：func (r *roleRepository) MethodName(ctx context.Context...) ... {
    # 但排除已经有 collection := 的方法和 getCollection 方法本身
    
    def add_collection_line(match):
        func_sig = match.group(1)
        body_start = match.group(2)
        
        # 跳过 getCollection 方法本身
        if 'getCollection' in func_sig:
            return match.group(0)
        
        # 检查是否已经有 collection := r.getCollection(ctx)
        if 'collection := r.getCollection(ctx)' in body_start[:100]:
            return match.group(0)
        
        # 添加 collection 获取语句
        return func_sig + '\n\tcollection := r.getCollection(ctx)' + body_start
    
    # 正则：匹配函数签名到函数体开始
    pattern = r'(func \(r \*roleRepository\) \w+\(ctx context\.Context[^\)]*\)[^\{]*\{)(\n)'
    content = re.sub(pattern, add_collection_line, content)
    
    # 3. 删除 GetRolesByTenant 方法，替换为 GetAllRoles
    # 查找并替换方法名
    content = re.sub(
        r'func \(r \*roleRepository\) GetRolesByTenant\(ctx context\.Context, tenantID string\)',
        'func (r *roleRepository) GetAllRoles(ctx context.Context)',
        content
    )
    
    # 删除方法体中的 tenant_id 过滤
    def fix_get_all_roles(match):
        before = match.group(1)
        after = match.group(2)
        # 移除 tenant_id 过滤
        new_filter = 'bson.M{"is_deleted": 0}'
        return before + new_filter + after
    
    pattern = r'(filter := )bson\.M\{[^}]*"tenant_id":\s*tenantID[^}]*\}(\s*cursor)'
    content = re.sub(pattern, fix_get_all_roles, content)
    
    # 4. 修复 GetCollection 方法的返回值
    content = re.sub(
        r'func \(r \*roleRepository\) GetCollection\(\) \*mongo\.Collection \{\s*return r\.collection\s*\}',
        '''func (r *roleRepository) GetCollection() *mongo.Collection {
\t// 返回系统库的collection作为默认值（向下兼容）
\treturn r.dbManager.GetSystemDatabase().Collection("role")
}''',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 修复完成")
    print("- 已替换所有 r.collection 为 collection")
    print("- 已在所有方法添加 collection := r.getCollection(ctx)")
    print("- 已修复 GetRolesByTenant -> GetAllRoles")
    print("- 已修复 GetCollection 方法")

if __name__ == '__main__':
    fix_role_repository()

