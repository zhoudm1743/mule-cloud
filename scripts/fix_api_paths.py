#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""批量修复前端API路径前缀"""

import os
import re
from pathlib import Path

# API路径映射规则
PATH_RULES = {
    # Production服务使用 /api 前缀
    r"'/production/": "'/api/production/",
    r'`/production/': '`/api/production/',
    
    # 其他服务使用 /admin 前缀
    r"'/order/": "'/admin/order/",
    r'`/order/': '`/admin/order/',
    
    r"'/workflow/": "'/admin/workflow/",
    r'`/workflow/': '`/admin/workflow/',
    
    r"'/basic/": "'/admin/basic/",
    r'`/basic/': '`/admin/basic/',
    
    r"'/tenant/": "'/admin/tenant/",
    r'`/tenant/': '`/admin/tenant/',
    
    r"'/admin-user/": "'/admin/admin-user/",
    r'`/admin-user/': '`/admin/admin-user/',
    
    r"'/role/": "'/admin/role/",
    r'`/role/': '`/admin/role/',
    
    r"'/menu/": "'/admin/menu/",
    r'`/menu/': '`/admin/menu/',
    
    r"'/department/": "'/admin/department/",
    r'`/department/': '`/admin/department/',
    
    r"'/post/": "'/admin/post/",
    r'`/post/': '`/admin/post/',
    
    r"'/system/": "'/admin/system/",
    r'`/system/': '`/admin/system/',
}

# 排除规则（已经有前缀的不再添加）
EXCLUDE_PATTERNS = [
    r"'/admin/",
    r'`/admin/',
    r"'/api/",
    r'`/api/',
]

def should_skip(content, pattern):
    """检查是否应该跳过（已经有前缀）"""
    for exclude in EXCLUDE_PATTERNS:
        # 检查替换后的结果是否已经存在
        if re.search(exclude, content):
            return True
    return False

def fix_api_file(file_path):
    """修复单个API文件"""
    print(f"\n处理文件: {file_path}")
    
    try:
        # 读取文件
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        changes = []
        
        # 应用所有规则
        for pattern, replacement in PATH_RULES.items():
            # 统计匹配数量
            matches = re.findall(pattern, content)
            if matches:
                # 执行替换
                content = re.sub(pattern, replacement, content)
                changes.append(f"  - {pattern} → {replacement} ({len(matches)}处)")
        
        # 如果有更改，写回文件
        if content != original_content:
            with open(file_path, 'w', encoding='utf-8', newline='\n') as f:
                f.write(content)
            print(f"  [OK] 已更新，修改了以下内容：")
            for change in changes:
                print(change)
            return True
        else:
            print("  [SKIP] 无需修改")
            return False
            
    except Exception as e:
        print(f"  [ERROR] 处理失败: {e}")
        return False

def main():
    """主函数"""
    api_dir = Path("frontend/src/service/api")
    
    if not api_dir.exists():
        print(f"❌ 目录不存在: {api_dir}")
        return
    
    print("=" * 60)
    print("开始批量修复API路径前缀")
    print("=" * 60)
    
    # 获取所有.ts文件
    ts_files = list(api_dir.glob("*.ts"))
    
    if not ts_files:
        print("❌ 未找到任何.ts文件")
        return
    
    print(f"\n找到 {len(ts_files)} 个API文件")
    
    updated_count = 0
    skipped_count = 0
    
    # 处理每个文件
    for ts_file in ts_files:
        if fix_api_file(ts_file):
            updated_count += 1
        else:
            skipped_count += 1
    
    print("\n" + "=" * 60)
    print("处理完成！")
    print("=" * 60)
    print(f"  已更新: {updated_count} 个文件")
    print(f"  已跳过: {skipped_count} 个文件")
    print(f"  总计:   {len(ts_files)} 个文件")
    print()

if __name__ == "__main__":
    main()

