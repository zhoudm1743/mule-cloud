#!/usr/bin/env python3
"""
智能修复Service层的语法错误
"""

import re

def fix_auth_service():
    """修复 auth.go 的语法错误"""
    file_path = r"K:\Git\mule-cloud\app\auth\services\auth.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复GenerateToken调用 - 删除TenantID参数及其注释
    # admin.ID, admin.Nickname, "" // 注释, admin.Roles -> admin.ID, admin.Nickname, admin.Roles
    content = re.sub(
        r'GenerateToken\(([^,]+),\s*([^,]+),\s*""\s*//[^,\n]+,\s*([^)]+)\)',
        r'GenerateToken(\1, \2, \3)',
        content
    )
    
    # 修复getUserMenuPermissions调用
    content = re.sub(
        r'getUserMenuPermissions\(ctx,\s*([^,]+),\s*""\s*//[^\)]+\)',
        r'getUserMenuPermissions(ctx, \1)',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ auth.go 语法修复完成")

def fix_role_service():
    """修复 role.go 的语法错误"""
    file_path = r"K:\Git\mule-cloud\app\system\services\role.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 修复if条件中的注释
    # if "" // 注释 != "" { -> 直接删除这个if块或改为false
    content = re.sub(
        r'if\s+""\s*//[^\n]+\s*!=\s*""\s*\{[^}]*\}',
        '// 数据库隔离后不再需要租户验证',
        content,
        flags=re.DOTALL
    )
    
    # 修复Get调用中的注释参数
    content = re.sub(
        r'\.Get\(ctx,\s*""\s*//[^\)]+\)',
        '.Get(ctx, roleID)',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 语法修复完成")

def main():
    print("智能修复Service层语法错误...")
    print()
    
    try:
        fix_auth_service()
        fix_role_service()
        
        print()
        print("✅ 语法修复完成！")
        print()
        print("验证编译：go build ./app/...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

