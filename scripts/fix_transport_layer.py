#!/usr/bin/env python3
"""
修复Transport层的TenantID引用
"""

import re

def fix_admin_transport():
    file_path = r"K:\Git\mule-cloud\app\system\transport\admin.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除租户权限检查（数据库隔离后不需要）
    content = re.sub(
        r'\s*if err := perm\.CheckTenantPermission\(targetAdmin\.TenantID, "[^"]+"\); err != nil \{[^}]+\}',
        '\n\t\t// 数据库隔离后不需要租户权限检查',
        content
    )
    
    content = re.sub(
        r'\s*if req\.TenantID != nil && \*req\.TenantID != targetAdmin\.TenantID \{[^}]+\}',
        '',
        content,
        flags=re.DOTALL
    )
    
    content = re.sub(
        r'\s*if err := perm\.CheckTenantPermission\(req\.TenantID, "[^"]+"\); err != nil \{[^}]+\}',
        '\n\t\t// 数据库隔离后不需要租户权限检查',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ admin.go 修复完成")

def fix_role_transport():
    file_path = r"K:\Git\mule-cloud\app\system\transport\role.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除租户权限检查
    content = re.sub(
        r'\s*if err := perm\.CheckTenantPermission\(targetRole\.TenantID, "[^"]+"\); err != nil \{[^}]+\}',
        '\n\t\t// 数据库隔离后不需要租户权限检查',
        content
    )
    
    content = re.sub(
        r'\s*if err := perm\.CheckTenantPermission\(req\.TenantID, "[^"]+"\); err != nil \{[^}]+\}',
        '\n\t\t// 数据库隔离后不需要租户权限检查',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 修复完成")

def main():
    print("修复Transport层TenantID引用...")
    print()
    
    try:
        fix_admin_transport()
        fix_role_transport()
        
        print()
        print("✅ Transport层修复完成！")
        print("验证编译：go build ./...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

