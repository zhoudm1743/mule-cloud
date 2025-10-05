#!/usr/bin/env python3
"""
修复role.go最后的collection声明问题
"""

file_path = r"K:\Git\mule-cloud\internal\repository\role.go"

with open(file_path, 'r', encoding='utf-8') as f:
    lines = f.readlines()

new_lines = []
skip_next = False
for i, line in enumerate(lines):
    if skip_next:
        skip_next = False
        continue
    
    # 如果这一行是声明collection，并且下一行是filter声明，而且这个方法调用了r.Find
    if '\tcollection := r.getCollection(ctx)' in line:
        # 检查后面10行是否有r.Find
        has_find = False
        for j in range(i+1, min(i+10, len(lines))):
            if 'r.Find(ctx, filter)' in lines[j]:
                has_find = True
                break
        
        # 如果使用了r.Find，说明不需要collection声明（Find会自己调用getCollection）
        if has_find:
            continue  # 跳过这一行
    
    new_lines.append(line)

with open(file_path, 'w', encoding='utf-8') as f:
    f.writelines(new_lines)

print("✅ role.go 最后的问题已修复")

