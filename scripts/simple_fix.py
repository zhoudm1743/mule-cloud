#!/usr/bin/env python3
"""
简单直接的修复：在每个需要collection的地方添加声明
"""

import re

def fix_menu():
    file_path = r"K:\Git\mule-cloud\internal\repository\menu.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    new_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]
        new_lines.append(line)
        
        # 如果是函数定义且包含ctx
        if ('func (r *' in line or 'func (r *MenuRepository)' in line) and 'ctx context.Context' in line and '{' in line:
            # 检查接下来3行是否有collection声明
            has_collection = False
            for j in range(i+1, min(i+4, len(lines))):
                if 'collection := r.getCollection(ctx)' in lines[j]:
                    has_collection = True
                    break
            
            # 检查后续是否使用collection
            uses_collection = False
            for j in range(i+1, min(i+30, len(lines))):
                if 'collection.' in lines[j] or 'collection)' in lines[j]:
                    uses_collection = True
                    break
            
            # 如果使用了collection但没有声明，添加
            if uses_collection and not has_collection and 'getCollection' not in line:
                new_lines.append('\tcollection := r.getCollection(ctx)\n')
        
        i += 1
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.writelines(new_lines)
    
    print("✅ menu.go 修复完成")

def fix_basic():
    file_path = r"K:\Git\mule-cloud\internal\repository\basic.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复返回值
    content = re.sub(
        r'func \(r \*basicRepository\) IsOwnedByTenant\(ctx context\.Context, basic \*models\.Basic, tenantID string\) \(bool, error\) \{[^}]*return true, nil  // 数据库隔离后，所有数据都属于当前租户',
        '''func (r *basicRepository) IsOwnedByTenant(ctx context.Context, basic *models.Basic, tenantID string) (bool, error) {
\treturn true, nil  // 数据库隔离后，所有数据都属于当前租户''',
        content
    )
    
    # 查找所有"return true"并改为"return true, nil"（如果函数签名要求返回(bool, error)）
    # 简单方法：全局替换
    content = content.replace(
        'return true  // 数据库隔离后，所有数据都属于当前租户\n}',
        'return true, nil  // 数据库隔离后，所有数据都属于当前租户\n}'
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ basic.go 修复完成")

def main():
    try:
        print("修复Repository编译错误...")
        fix_menu()
        fix_basic()
        print("\n✅ 修复完成！请验证编译。")
    except Exception as e:
        print(f"❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

