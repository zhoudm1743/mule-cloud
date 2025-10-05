#!/usr/bin/env python3
"""
修复Repository编译错误
"""

import re

def fix_basic_repository():
    """修复basic.go的编译错误"""
    file_path = r"K:\Git\mule-cloud\internal\repository\basic.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除所有引用TenantID的代码
    # 查找并修复IsOwnedByTenant等方法
    content = re.sub(
        r'basic\.TenantID',
        'true  // 数据库隔离后，当前库的数据都属于当前租户',
        content
    )
    
    # 删除声明但未使用的collection
    # 找到那些只声明了collection但没有使用的方法，删除声明行
    lines = content.split('\n')
    new_lines = []
    i = 0
    while i < len(lines):
        line = lines[i]
        # 如果这一行是collection声明，检查后面的代码是否使用了它
        if '\tcollection := r.getCollection(ctx)' in line:
            # 查看接下来的20行是否使用了collection
            next_20_lines = '\n'.join(lines[i+1:min(i+21, len(lines))])
            if 'collection.' not in next_20_lines and 'collection)' not in next_20_lines:
                # 未使用，跳过这一行
                i += 1
                continue
        new_lines.append(line)
        i += 1
    
    content = '\n'.join(new_lines)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ basic.go 编译错误已修复")

def fix_menu_repository():
    """修复menu.go的编译错误"""
    file_path = r"K:\Git\mule-cloud\internal\repository\menu.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        lines = f.readlines()
    
    new_lines = []
    for i, line in enumerate(lines):
        new_lines.append(line)
        
        # 在函数定义后添加collection声明（如果还没有）
        if 'func (r *menuRepository)' in line and 'ctx context.Context' in line and '{' in line:
            # 检查下一行是否已经有collection声明
            if i + 1 < len(lines):
                next_line = lines[i + 1]
                if 'collection := r.getCollection(ctx)' not in next_line:
                    # 检查是否是getCollection方法本身
                    if 'getCollection' not in line:
                        # 检查接下来几行是否使用了collection
                        next_10_lines = ''.join(lines[i+1:min(i+11, len(lines))])
                        if 'collection.' in next_10_lines or 'collection)' in next_10_lines:
                            new_lines.append('\tcollection := r.getCollection(ctx)\n')
    
    content = ''.join(new_lines)
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ menu.go 编译错误已修复")

def main():
    print("修复Repository编译错误...")
    print()
    
    try:
        fix_basic_repository()
        fix_menu_repository()
        
        print()
        print("=" * 60)
        print("  ✅ 编译错误修复完成！")
        print("=" * 60)
        print()
        print("请重新验证编译：go build ./internal/repository/...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

