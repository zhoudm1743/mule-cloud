#!/usr/bin/env python3
"""
批量完成所有Repository的改造
"""

import re
import os

def fix_role_remaining():
    """完成role.go剩余方法的改造"""
    file_path = r"K:\Git\mule-cloud\internal\repository\role.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 在每个方法开头添加 collection := r.getCollection(ctx)
    # 但要排除已经有的和getCollection方法本身
    def add_collection_if_needed(match):
        func_def = match.group(1)
        first_line = match.group(2)
        
        # 跳过getCollection方法本身
        if 'getCollection' in func_def:
            return match.group(0)
        
        # 如果第一行已经是collection声明，跳过
        if 'collection := r.getCollection(ctx)' in first_line:
            return match.group(0)
        
        # 添加collection声明
        indent = '\t'
        return func_def + '\n' + indent + 'collection := r.getCollection(ctx)' + first_line
    
    # 匹配函数定义和第一行
    pattern = r'(func \(r \*roleRepository\) \w+\(ctx context\.Context[^\)]*\)[^\{]*\{)(\n[^\n]*)'
    content = re.sub(pattern, add_collection_if_needed, content)
    
    # 修复GetRolesByTenant为GetAllRoles并删除tenant_id过滤
    content = re.sub(
        r'func \(r \*roleRepository\) GetRolesByTenant\(ctx context\.Context, tenantID string\)',
        'func (r *roleRepository) GetAllRoles(ctx context.Context)',
        content
    )
    
    # 删除GetAllRoles中的tenant_id过滤
    content = re.sub(
        r'(func \(r \*roleRepository\) GetAllRoles\(ctx context\.Context\)[^\{]*\{[^\}]*filter := )bson\.M\{[^}]*"tenant_id":\s*tenantID,\s*"is_deleted":\s*0[^}]*\}',
        r'\1bson.M{"is_deleted": 0}',
        content,
        flags=re.DOTALL
    )
    
    # 修复GetCollection方法
    pattern_getcoll = r'func \(r \*roleRepository\) GetCollection\(\) \*mongo\.Collection \{[^\}]*\}'
    replacement_getcoll = '''func (r *roleRepository) GetCollection() *mongo.Collection {
\treturn r.dbManager.GetSystemDatabase().Collection("role")
}'''
    content = re.sub(pattern_getcoll, replacement_getcoll, content)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 剩余方法改造完成")

def create_menu_repository():
    """改造menu.go"""
    file_path = r"K:\Git\mule-cloud\internal\repository\menu.go"
    
    if not os.path.exists(file_path):
        print("⚠️  menu.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 添加导入
    if 'tenantCtx "mule-cloud/core/context"' not in content:
        content = content.replace(
            'import (\n\t"context"',
            'import (\n\t"context"\n\ttenantCtx "mule-cloud/core/context"'
        )
    
    # 2. 修改结构体
    content = re.sub(
        r'type menuRepository struct \{[^\}]*collection \*mongo\.Collection[^\}]*\}',
        'type menuRepository struct {\n\tdbManager *database.DatabaseManager\n}',
        content
    )
    
    # 3. 修改构造函数
    content = re.sub(
        r'func NewMenuRepository\(\) MenuRepository \{[^\}]*collection := database\.MongoDB\.Collection\("menu"\)[^\}]*return &menuRepository\{[^\}]*collection: collection[^\}]*\}[^\}]*\}',
        '''func NewMenuRepository() MenuRepository {
\treturn &menuRepository{
\t\tdbManager: database.GetDatabaseManager(),
\t}
}

// getCollection 获取集合（自动根据Context中的租户ID切换数据库）
func (r *menuRepository) getCollection(ctx context.Context) *mongo.Collection {
\ttenantID := tenantCtx.GetTenantID(ctx)
\tdb := r.dbManager.GetDatabase(tenantID)
\treturn db.Collection("menu")
}''',
        content
    )
    
    # 4. 替换所有 r.collection. 为 collection.
    content = content.replace('r.collection.', 'collection.')
    
    # 5. 在每个方法开头添加 collection := r.getCollection(ctx)
    def add_collection(match):
        func_def = match.group(1)
        first_line = match.group(2)
        
        if 'getCollection' in func_def:
            return match.group(0)
        if 'collection := r.getCollection(ctx)' in first_line:
            return match.group(0)
        
        return func_def + '\n\tcollection := r.getCollection(ctx)' + first_line
    
    pattern = r'(func \(r \*menuRepository\) \w+\(ctx context\.Context[^\)]*\)[^\{]*\{)(\n[^\n]*)'
    content = re.sub(pattern, add_collection, content)
    
    # 6. 修复GetCollection方法
    content = re.sub(
        r'func \(r \*menuRepository\) GetCollection\(\) \*mongo\.Collection \{[^\}]*\}',
        '''func (r *menuRepository) GetCollection() *mongo.Collection {
\treturn r.dbManager.GetSystemDatabase().Collection("menu")
}''',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ menu.go 改造完成")

def create_basic_repository():
    """改造basic.go"""
    file_path = r"K:\Git\mule-cloud\internal\repository\basic.go"
    
    if not os.path.exists(file_path):
        print("⚠️  basic.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 添加导入
    if 'tenantCtx "mule-cloud/core/context"' not in content:
        content = content.replace(
            'import (\n\t"context"',
            'import (\n\t"context"\n\ttenantCtx "mule-cloud/core/context"'
        )
    
    # 2. 修改结构体
    content = re.sub(
        r'type basicRepository struct \{[^\}]*collection \*mongo\.Collection[^\}]*\}',
        'type basicRepository struct {\n\tdbManager *database.DatabaseManager\n}',
        content
    )
    
    # 3. 修改构造函数
    content = re.sub(
        r'func NewBasicRepository\(\) BasicRepository \{[^\}]*collection := database\.MongoDB\.Collection\("basic"\)[^\}]*return &basicRepository\{[^\}]*collection: collection[^\}]*\}[^\}]*\}',
        '''func NewBasicRepository() BasicRepository {
\treturn &basicRepository{
\t\tdbManager: database.GetDatabaseManager(),
\t}
}

// getCollection 获取集合（自动根据Context中的租户ID切换数据库）
func (r *basicRepository) getCollection(ctx context.Context) *mongo.Collection {
\ttenantID := tenantCtx.GetTenantID(ctx)
\tdb := r.dbManager.GetDatabase(tenantID)
\treturn db.Collection("basic")
}''',
        content
    )
    
    # 4. 替换所有 r.collection. 为 collection.
    content = content.replace('r.collection.', 'collection.')
    
    # 5. 在每个方法开头添加 collection := r.getCollection(ctx)
    def add_collection(match):
        func_def = match.group(1)
        first_line = match.group(2)
        
        if 'getCollection' in func_def:
            return match.group(0)
        if 'collection := r.getCollection(ctx)' in first_line:
            return match.group(0)
        
        return func_def + '\n\tcollection := r.getCollection(ctx)' + first_line
    
    pattern = r'(func \(r \*basicRepository\) \w+\(ctx context\.Context[^\)]*\)[^\{]*\{)(\n[^\n]*)'
    content = re.sub(pattern, add_collection, content)
    
    # 6. 修复GetCollection方法
    content = re.sub(
        r'func \(r \*basicRepository\) GetCollection\(\) \*mongo\.Collection \{[^\}]*\}',
        '''func (r *basicRepository) GetCollection() *mongo.Collection {
\treturn r.dbManager.GetSystemDatabase().Collection("basic")
}''',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ basic.go 改造完成")

def create_tenant_repository():
    """改造tenant.go - 特殊处理，固定使用系统库"""
    file_path = r"K:\Git\mule-cloud\internal\repository\tenant.go"
    
    if not os.path.exists(file_path):
        print("⚠️  tenant.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 修改结构体
    content = re.sub(
        r'type tenantRepository struct \{[^\}]*collection \*mongo\.Collection[^\}]*\}',
        'type tenantRepository struct {\n\tdbManager *database.DatabaseManager\n}',
        content
    )
    
    # 2. 修改构造函数
    content = re.sub(
        r'func NewTenantRepository\(\) TenantRepository \{[^\}]*collection := database\.MongoDB\.Collection\("tenant"\)[^\}]*return &tenantRepository\{[^\}]*collection: collection[^\}]*\}[^\}]*\}',
        '''func NewTenantRepository() TenantRepository {
\treturn &tenantRepository{
\t\tdbManager: database.GetDatabaseManager(),
\t}
}

// getCollection 获取集合（租户数据固定使用系统数据库）
func (r *tenantRepository) getCollection() *mongo.Collection {
\tdb := r.dbManager.GetSystemDatabase()
\treturn db.Collection("tenant")
}''',
        content
    )
    
    # 3. 替换所有 r.collection. 为 collection.
    content = content.replace('r.collection.', 'collection.')
    
    # 4. 在每个方法开头添加 collection := r.getCollection()（注意：不需要ctx）
    def add_collection(match):
        func_def = match.group(1)
        first_line = match.group(2)
        
        if 'getCollection' in func_def:
            return match.group(0)
        if 'collection := r.getCollection()' in first_line:
            return match.group(0)
        
        return func_def + '\n\tcollection := r.getCollection()' + first_line
    
    pattern = r'(func \(r \*tenantRepository\) \w+\([^\)]*\)[^\{]*\{)(\n[^\n]*)'
    content = re.sub(pattern, add_collection, content)
    
    # 5. 修复GetCollection方法
    content = re.sub(
        r'func \(r \*tenantRepository\) GetCollection\(\) \*mongo\.Collection \{[^\}]*\}',
        '''func (r *tenantRepository) GetCollection() *mongo.Collection {
\treturn r.dbManager.GetSystemDatabase().Collection("tenant")
}''',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ tenant.go 改造完成（特殊：固定使用系统库）")

def main():
    print("=" * 60)
    print("  批量完成所有Repository改造")
    print("=" * 60)
    print()
    
    try:
        print("[1/4] 完成 role.go 剩余方法...")
        fix_role_remaining()
        print()
        
        print("[2/4] 改造 menu.go...")
        create_menu_repository()
        print()
        
        print("[3/4] 改造 basic.go...")
        create_basic_repository()
        print()
        
        print("[4/4] 改造 tenant.go（特殊处理）...")
        create_tenant_repository()
        print()
        
        print("=" * 60)
        print("  ✅ 所有Repository改造完成！")
        print("=" * 60)
        print()
        print("下一步：")
        print("1. 验证编译：go build ./internal/repository/...")
        print("2. 继续改造Service层")
        print()
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

