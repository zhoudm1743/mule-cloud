#!/usr/bin/env python3
"""
修复System服务中的tenant_id过滤
数据库隔离后不需要tenant_id过滤
"""

import re

def fix_admin_service():
    file_path = r"K:\Git\mule-cloud\app\system\services\admin.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除 tenant_id 过滤
    content = re.sub(
        r'\s*if req\.TenantID != "" \{\s*filter\["tenant_id"\] = req\.TenantID // 租户过滤\s*\}',
        '\n\t// 数据库隔离后不需要tenant_id过滤',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ admin.go 修复完成")

def fix_role_service():
    file_path = r"K:\Git\mule-cloud\app\system\services\role.go"
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除 tenant_id 过滤
    content = re.sub(
        r'\s*if req\.TenantID != "" \{\s*filter\["tenant_id"\] = req\.TenantID\s*\}',
        '\n\t// 数据库隔离后不需要tenant_id过滤',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print("✅ role.go 修复完成")

def main():
    print("修复System服务中的tenant_id过滤...")
    print()
    
    try:
        fix_admin_service()
        fix_role_service()
        
        print()
        print("✅ 修复完成！")
        print()
        print("📝 注意：DTO中的TenantID字段保留是为了：")
        print("   1. 前端显示（只读）")
        print("   2. 向后兼容")
        print("   3. 创建时会被忽略（从Context获取）")
        print()
        print("验证编译：go build ./...")
        
    except Exception as e:
        print(f"\n❌ 错误：{e}")
        import traceback
        traceback.print_exc()

if __name__ == '__main__':
    main()

