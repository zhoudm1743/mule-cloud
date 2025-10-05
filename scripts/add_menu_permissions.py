#!/usr/bin/env python3
"""
为菜单添加 available_permissions（如果没有的话）
"""

from pymongo import MongoClient

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def add_menu_permissions():
    """为菜单添加权限配置"""
    
    system_db = client['tenant_system']
    
    # 标准CRUD权限
    crud_permissions = [
        {"action": "read", "label": "查看", "is_basic": True},
        {"action": "create", "label": "新增", "is_basic": True},
        {"action": "update", "label": "编辑", "is_basic": True},
        {"action": "delete", "label": "删除", "is_basic": True},
    ]
    
    # 获取所有菜单
    menus = list(system_db['menu'].find({'is_deleted': 0}))
    
    print("\n" + "="*60)
    print("🔧 为菜单添加权限配置")
    print("="*60)
    
    for menu in menus:
        menu_id = menu['_id']
        menu_name = menu.get('name')
        menu_title = menu.get('title')
        menu_type = menu.get('menuType')
        available_perms = menu.get('available_permissions', [])
        
        # 只为 page 类型的菜单添加权限
        if menu_type == 'page' and not available_perms:
            print(f"\n✨ 为 {menu_name} ({menu_title}) 添加权限...")
            
            system_db['menu'].update_one(
                {'_id': menu_id},
                {'$set': {'available_permissions': crud_permissions}}
            )
            
            print(f"   ✅ 已添加: read, create, update, delete")
        elif menu_type == 'page':
            print(f"\n✓ {menu_name} ({menu_title}) 已有权限配置")
        else:
            print(f"\n- {menu_name} ({menu_title}) [目录，跳过]")
    
    print("\n" + "="*60)
    print("✅ 完成！")
    print("="*60)

if __name__ == '__main__':
    try:
        add_menu_permissions()
    finally:
        client.close()

