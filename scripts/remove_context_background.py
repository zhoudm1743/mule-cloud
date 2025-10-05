#!/usr/bin/env python3
"""
删除残留的 ctx := context.Background() 行
"""

import re

services = ['size', 'order_type', 'salesman', 'procedure']

for service in services:
    file_path = f'app/basic/services/{service}.go'
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 删除 Get 方法中的 ctx := context.Background()
    content = re.sub(
        r'func \(s \*\w+Service\) Get\(ctx context\.Context, id string\) \(\*models\.Basic, error\) \{\n\tctx := context\.Background\(\)\n\treturn',
        r'func (s *\\1Service) Get(ctx context.Context, id string) (*models.Basic, error) {\n\treturn',
        content
    )
    
    # 删除 Delete 方法中的 ctx := context.Background()
    content = re.sub(
        r'func \(s \*(\w+)Service\) Delete\(ctx context\.Context, id string\) error \{\n\tctx := context\.Background\(\)\n\treturn',
        r'func (s *\1Service) Delete(ctx context.Context, id string) error {\n\treturn',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"✅ 修复 {file_path}")

print("\n✅ 完成！")

