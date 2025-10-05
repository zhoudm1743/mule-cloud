#!/usr/bin/env python3
"""
批量修改 basic 服务，添加 context 参数支持租户隔离
"""

import os
import re

# 服务列表（除了 color 已经修改完成）
services = ['size', 'customer', 'order_type', 'salesman', 'procedure']

# 服务名到类型名的映射
service_types = {
    'size': 'Size',
    'customer': 'Customer',
    'order_type': 'OrderType',
    'salesman': 'Salesman',
    'procedure': 'Procedure'
}

def fix_service_file(service):
    """修复 service 文件"""
    service_type = service_types[service]
    file_path = f'app/basic/services/{service}.go'
    
    print(f"\n🔧 修复 {file_path}...")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 1. 修改接口定义
    content = re.sub(
        rf'Get\(id string\)',
        r'Get(ctx context.Context, id string)',
        content
    )
    content = re.sub(
        rf'GetAll\(req dto\.{service_type}ListRequest\)',
        rf'GetAll(ctx context.Context, req dto.{service_type}ListRequest)',
        content
    )
    content = re.sub(
        rf'List\(req dto\.{service_type}ListRequest\)',
        rf'List(ctx context.Context, req dto.{service_type}ListRequest)',
        content
    )
    content = re.sub(
        rf'Create\(req dto\.{service_type}CreateRequest\)',
        rf'Create(ctx context.Context, req dto.{service_type}CreateRequest)',
        content
    )
    content = re.sub(
        rf'Update\(req dto\.{service_type}UpdateRequest\)',
        rf'Update(ctx context.Context, req dto.{service_type}UpdateRequest)',
        content
    )
    content = re.sub(
        r'Delete\(id string\)',
        r'Delete(ctx context.Context, id string)',
        content
    )
    
    # 2. 修改方法实现，删除 ctx := context.Background()
    content = re.sub(
        r'func \(s \*\w+Service\) Get\(id string\) \(\*models\.Basic, error\) \{\n\tctx := context\.Background\(\)\n\treturn',
        r'func (s *' + service_type + r'Service) Get(ctx context.Context, id string) (*models.Basic, error) {\n\treturn',
        content
    )
    
    content = re.sub(
        r'func \(s \*\w+Service\) GetAll\((ctx context\.Context, )?req dto\.\w+ListRequest\) \(\[\]\*?models\.Basic, error\) \{\n\tctx := context\.Background\(\)\n\n\t// 构建过滤条件',
        r'func (s *' + service_type + r'Service) GetAll(ctx context.Context, req dto.' + service_type + r'ListRequest) ([]' + ('*' if service != 'color' else '') + r'models.Basic, error) {\n\t// 构建过滤条件',
        content
    )
    
    content = re.sub(
        r'func \(s \*\w+Service\) List\((ctx context\.Context, )?req dto\.\w+ListRequest\) \(\[\]\*?models\.Basic, int64, error\) \{\n\tctx := context\.Background\(\)\n\n\t// 构建过滤条件',
        r'func (s *' + service_type + r'Service) List(ctx context.Context, req dto.' + service_type + r'ListRequest) ([]' + ('*' if service == 'size' else '') + r'models.Basic, int64, error) {\n\t// 构建过滤条件',
        content
    )
    
    content = re.sub(
        r'func \(s \*\w+Service\) Create\((ctx context\.Context, )?req dto\.\w+CreateRequest\) \(\*models\.Basic, error\) \{\n\tctx := context\.Background\(\)\n\tnow := time\.Now\(\)\.Unix\(\)',
        r'func (s *' + service_type + r'Service) Create(ctx context.Context, req dto.' + service_type + r'CreateRequest) (*models.Basic, error) {\n\tnow := time.Now().Unix()',
        content
    )
    
    content = re.sub(
        r'func \(s \*\w+Service\) Update\((ctx context\.Context, )?req dto\.\w+UpdateRequest\) \(\*models\.Basic, error\) \{\n\tctx := context\.Background\(\)\n\n\t// 更新字段',
        r'func (s *' + service_type + r'Service) Update(ctx context.Context, req dto.' + service_type + r'UpdateRequest) (*models.Basic, error) {\n\t// 更新字段',
        content
    )
    
    content = re.sub(
        r'func \(s \*\w+Service\) Delete\(id string\) error \{\n\tctx := context\.Background\(\)\n\treturn',
        r'func (s *' + service_type + r'Service) Delete(ctx context.Context, id string) error {\n\treturn',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"   ✅ 完成")

def fix_endpoint_file(service):
    """修复 endpoint 文件"""
    service_type = service_types[service]
    file_path = f'app/basic/endpoint/{service}.go'
    
    print(f"\n🔧 修复 {file_path}...")
    
    with open(file_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # 添加 ctx 参数到 service 调用
    content = re.sub(
        r'svc\.Get\(req\.ID\)',
        r'svc.Get(ctx, req.ID)',
        content
    )
    content = re.sub(
        r'svc\.GetAll\(req\)',
        r'svc.GetAll(ctx, req)',
        content
    )
    content = re.sub(
        r'svc\.List\(req\)',
        r'svc.List(ctx, req)',
        content
    )
    content = re.sub(
        r'svc\.Create\(req\)',
        r'svc.Create(ctx, req)',
        content
    )
    content = re.sub(
        r'svc\.Update\(req\)',
        r'svc.Update(ctx, req)',
        content
    )
    content = re.sub(
        r'svc\.Delete\(req\.ID\)',
        r'svc.Delete(ctx, req.ID)',
        content
    )
    
    with open(file_path, 'w', encoding='utf-8') as f:
        f.write(content)
    
    print(f"   ✅ 完成")

if __name__ == '__main__':
    print("="*60)
    print("🔧 批量修复 basic 服务的 context 参数")
    print("="*60)
    
    for service in services:
        try:
            fix_service_file(service)
            fix_endpoint_file(service)
        except Exception as e:
            print(f"   ❌ 失败: {e}")
    
    print("\n" + "="*60)
    print("✅ 修复完成！")
    print("="*60)

