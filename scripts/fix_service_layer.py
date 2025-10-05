#!/usr/bin/env python3
"""
批量修复Service层的TenantID相关代码
"""

import re
import os

def fix_auth_service():
    """修复 app/auth/services/auth.go"""
    file_path = r"K:\Git\mule-cloud\app\auth\services\auth.go"
    
    if not os.path.exists(file_path):
        print("⚠️  auth.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除所有 admin.TenantID 的引用
    # 1. 删除JWT中的TenantID赋值
    content = re.sub(
        r',\s*TenantID:\s*admin\.TenantID',
        '',
        content
    )
    
    # 2. 删除结构体中的TenantID字段赋值
    content = re.sub(
        r'TenantID:\s*admin\.TenantID,?\s*\n',
        '',
        content
    )
    
    # 3. 删除单独的admin.TenantID引用
    content = re.sub(
        r'admin\.TenantID',
        '"" // 数据库隔离后不再需要TenantID',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ auth.go 修复完成")

def fix_admin_service():
    """修复 app/system/services/admin.go"""
    file_path = r"K:\Git\mule-cloud\app\system\services\admin.go"
    
    if not os.path.exists(file_path):
        print("⚠️  admin.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除创建Admin时的TenantID字段
    content = re.sub(
        r'TenantID:\s*[^,\n]+,\s*\n',
        '',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ admin.go 修复完成")

def fix_role_service():
    """修复 app/system/services/role.go"""
    file_path = r"K:\Git\mule-cloud\app\system\services\role.go"
    
    if not os.path.exists(file_path):
        print("⚠️  role.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 修复 GetByCode 调用（删除第三个参数）
    content = re.sub(
        r'\.GetByCode\(ctx,\s*([^,]+),\s*[^)]+\)',
        r'.GetByCode(ctx, \1)',
        content
    )
    
    # 2. 修复 GetByName 调用（删除第三个参数）
    content = re.sub(
        r'\.GetByName\(ctx,\s*([^,]+),\s*[^)]+\)',
        r'.GetByName(ctx, \1)',
        content
    )
    
    # 3. 删除创建Role时的TenantID字段
    content = re.sub(
        r'TenantID:\s*[^,\n]+,\s*\n',
        '',
        content
    )
    
    # 4. 删除所有 role.TenantID 的引用
    content = re.sub(
        r'role\.TenantID',
        '"" // 数据库隔离后不再需要TenantID',
        content
    )
    
    # 5. 修改 GetRolesByTenant 为 GetAllRoles
    content = re.sub(
        r'\.GetRolesByTenant\(ctx,\s*[^)]+\)',
        '.GetAllRoles(ctx)',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 修复完成")

def fix_menu_service():
    """修复 app/system/services/menu.go（如果有TenantID引用）"""
    file_path = r"K:\Git\mule-cloud\app\system\services\menu.go"
    
    if not os.path.exists(file_path):
        print("⚠️  menu.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除TenantID相关代码
    content = re.sub(
        r'TenantID:\s*[^,\n]+,?\s*\n',
        '',
        content
    )
    
    content = re.sub(
        r'menu\.TenantID',
        '"" // 数据库隔离后不再需要TenantID',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ menu.go 修复完成")

def fix_tenant_service():
    """修复 app/system/services/tenant.go（特殊：需要创建数据库）"""
    file_path = r"K:\Git\mule-cloud\app\system\services\tenant.go"
    
    if not os.path.exists(file_path):
        print("⚠️  tenant.go 不存在，跳过")
        return
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除TenantID相关代码
    content = re.sub(
        r'tenant\.TenantID',
        'tenant.ID',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ tenant.go 修复完成")

def fix_basic_services():
    """修复 app/basic/services/*.go"""
    base_path = r"K:\Git\mule-cloud\app\basic\services"
    
    if not os.path.exists(base_path):
        print("⚠️  basic/services 不存在，跳过")
        return
    
    for filename in os.listdir(base_path):
        if not filename.endswith('.go'):
            continue
        
        file_path = os.path.join(base_path, filename)
        
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # 删除TenantID字段赋值
        content = re.sub(
            r'TenantID:\s*[^,\n]+,?\s*\n',
            '',
            content
        )
        
        # 删除.TenantID引用
        content = re.sub(
            r'\.TenantID\s*==\s*tenantID',
            'true // 数据库隔离后不再需要TenantID比较',
            content
        )
        
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content)
    
    print(f"✅ basic/services 所有文件修复完成")

def main():
    print("=" * 60)
    print("  批量修复Service层")
    print("=" * 60)
    print()
    
    try:
        print("[1/6] 修复 auth/services/auth.go...")
        fix_auth_service()
        print()
        
        print("[2/6] 修复 system/services/admin.go...")
        fix_admin_service()
        print()
        
        print("[3/6] 修复 system/services/role.go...")
        fix_role_service()
        print()
        
        print("[4/6] 修复 system/services/menu.go...")
        fix_menu_service()
        print()
        
        print("[5/6] 修复 system/services/tenant.go...")
        fix_tenant_service()
        print()
        
        print("[6/6] 修复 basic/services/*.go...")
        fix_basic_services()
        print()
        
        print("=" * 60)
        print("  ✅ Service层修复完成！")
        print("=" * 60)
        print()
        print("验证编译：go build ./app/...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

