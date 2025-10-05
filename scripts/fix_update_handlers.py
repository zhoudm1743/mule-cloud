#!/usr/bin/env python3
"""
修复所有 basic transport 的 Update handler 参数绑定顺序
"""

import os
import re

transport_files = [
    'app/basic/transport/color.go',
    'app/basic/transport/customer.go',
    'app/basic/transport/order_type.go',
    'app/basic/transport/salesman.go',
    'app/basic/transport/procedure.go',
]

for file_path in transport_files:
    print(f"\n🔧 修复 {file_path}...")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 查找 Update handler 中的参数绑定
    # 将 ShouldBindUri 和 ShouldBindJSON 的顺序交换
    pattern = r'(func Update\w+Handler.*?\{[\s\S]*?var req dto\.\w+UpdateRequest\n)'
    pattern += r'(\s+if err := c\.ShouldBindUri\(&req\); err != nil \{[\s\S]*?return\n\s+\}\n)'
    pattern += r'(\s+if err := c\.ShouldBindJSON\(&req\); err != nil \{[\s\S]*?return\n\s+\})'
    
    def swap_bindings(match):
        prefix = match.group(1)
        uri_binding = match.group(2)
        json_binding = match.group(3)
        
        # 添加注释并交换顺序
        json_with_comment = json_binding.replace(
            'if err := c.ShouldBindJSON(&req)',
            '// 先绑定 JSON body（包含 required 字段）\n\t\tif err := c.ShouldBindJSON(&req)'
        )
        uri_with_comment = uri_binding.replace(
            'if err := c.ShouldBindUri(&req)',
            '// 再绑定 URI 参数（ID）\n\t\tif err := c.ShouldBindUri(&req)'
        )
        
        return prefix + json_with_comment + '\n' + uri_with_comment
    
    content_new = re.sub(pattern, swap_bindings, content, flags=re.MULTILINE)
    
    if content_new != content:
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(content_new)
        print(f"   ✅ 已修复")
    else:
        print(f"   ⚠️  未找到需要修复的代码")

print("\n✅ 修复完成！")





