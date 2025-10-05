#!/usr/bin/env python3
"""
获取所有菜单及其权限，用于创建租户时自动分配
"""

from pymongo import MongoClient
import json

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def get_all_menus_and_permissions():
    """获取所有菜单及其权限配置"""
    
    system_db = client['tenant_system']
    
    menus = list(system_db['menu'].find({'is_deleted': 0}))
    
    print("\n" + "="*60)
    print("📋 所有菜单及权限配置")
    print("="*60)
    
    menu_names = []
    menu_permissions = {}
    
    for menu in menus:
        menu_name = menu.get('name')
        menu_title = menu.get('title')
        available_perms = menu.get('available_permissions', [])
        
        menu_names.append(menu_name)
        
        print(f"\n菜单: {menu_name} ({menu_title})")
        
        if available_perms:
            perm_list = [p.get('action') for p in available_perms]
            menu_permissions[menu_name] = perm_list
            print(f"  权限: {', '.join(perm_list)}")
        else:
            print("  权限: 无")
    
    print("\n" + "="*60)
    print("💡 用于创建租户时的配置:")
    print("="*60)
    
    print("\n// 所有菜单名称")
    print(f"menus := []string{{{', '.join([f'\\"{name}\\"' for name in menu_names])}}}")
    
    print("\n// 菜单权限映射")
    print("menuPermissions := map[string][]string{")
    for menu_name, perms in menu_permissions.items():
        perm_str = ', '.join([f'"{p}"' for p in perms])
        print(f'    "{menu_name}": {{{perm_str}}},')
    print("}")
    
    print("\n" + "="*60)
    
    return menu_names, menu_permissions

if __name__ == '__main__':
    try:
        get_all_menus_and_permissions()
    finally:
        client.close()

