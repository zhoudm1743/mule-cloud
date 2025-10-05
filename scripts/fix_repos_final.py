#!/usr/bin/env python3
"""
最终修复Repository的所有问题
"""

import re

def fix_repository_file(file_path, repo_name):
    """修复单个repository文件"""
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 确保每个使用collection的方法都有声明
    # 找到所有的函数定义
    func_pattern = r'(func \(r \*' + repo_name + r'Repository\) (\w+)\([^)]*\)[^{]*\{)'
    
    def ensure_collection(match):
        func_full = match.group(0)
        func_name = match.group(2)
        
        # 跳过getCollection方法本身
        if func_name == 'getCollection':
            return func_full
        
        # 查找这个函数到下一个函数之间的内容
        start_pos = match.end()
        # 简单查找：找到函数体中是否使用了collection
        next_func = re.search(r'\nfunc \(', content[start_pos:])
        if next_func:
            func_body = content[start_pos:start_pos + next_func.start()]
        else:
            func_body = content[start_pos:]
        
        # 检查是否使用了collection
        uses_collection = ('collection.' in func_body or 
                          'collection)' in func_body or
                          'collection,' in func_body)
        
        # 检查是否已经有collection声明
        has_declaration = 'collection := r.getCollection' in func_body[:200]
        
        if uses_collection and not has_declaration:
            # 需要添加声明
            if repo_name == 'tenant':
                # tenant特殊：不需要ctx
                return func_full + '\n\tcollection := r.getCollection()'
            else:
                return func_full + '\n\tcollection := r.getCollection(ctx)'
        else:
            return func_full
    
    content = re.sub(func_pattern, ensure_collection, content)
    
    # 2. 修复basic.go中的TenantID引用
    if repo_name == 'basic':
        # 删除IsOwnedByTenant等方法中对TenantID的引用
        content = re.sub(
            r'basic\.TenantID == tenantID',
            'true  // 数据库隔离后，所有数据都属于当前租户',
            content
        )
        
        # 修复方法签名
        content = re.sub(
            r'func \(r \*basicRepository\) IsOwnedByTenant\(ctx context\.Context, basic \*models\.Basic, tenantID string\) bool \{',
            'func (r *basicRepository) IsOwnedByTenant(ctx context.Context, basic *models.Basic, tenantID string) (bool, error) {',
            content
        )
        
        # 修复返回值
        content = re.sub(
            r'return true  // 数据库隔离后，所有数据都属于当前租户\n\}',
            'return true, nil  // 数据库隔离后，所有数据都属于当前租户\n}',
            content
        )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✅ {repo_name}.go 修复完成")

def main():
    print("=" * 60)
    print("  最终修复所有Repository")
    print("=" * 60)
    print()
    
    try:
        fix_repository_file(r"K:\Git\mule-cloud\internal\repository\basic.go", "basic")
        fix_repository_file(r"K:\Git\mule-cloud\internal\repository\menu.go", "menu")
        fix_repository_file(r"K:\Git\mule-cloud\internal\repository\role.go", "role")
        fix_repository_file(r"K:\Git\mule-cloud\internal\repository\tenant.go", "tenant")
        
        print()
        print("=" * 60)
        print("  ✅ 所有Repository修复完成！")
        print("=" * 60)
        print()
        print("验证编译：go build ./internal/repository/...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

