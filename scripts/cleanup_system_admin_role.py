#!/usr/bin/env python3
"""
清理系统库中的 admin 和 role 数据（迁移后）
"""

from pymongo import MongoClient

# 连接 MongoDB
client = MongoClient('mongodb://root:bgg8384495@localhost:27015/', authSource='admin')

def cleanup_system_db():
    """清理系统库中的 admin 和 role"""
    
    print("\n" + "="*60)
    print("🧹 清理系统库中的 admin 和 role")
    print("="*60)
    
    system_db = client['tenant_system']
    
    # 统计当前数量
    admin_count = system_db['admin'].count_documents({})
    role_count = system_db['role'].count_documents({})
    
    print(f"\n当前系统库数据:")
    print(f"  - admin: {admin_count} 条")
    print(f"  - role: {role_count} 条")
    
    if admin_count == 0 and role_count == 0:
        print("\n✅ 系统库已经是干净的，无需清理")
        return
    
    # 确认
    print("\n⚠️  警告：即将删除系统库中的所有 admin 和 role 数据")
    print("   请确保已经迁移到租户库并测试通过")
    
    confirm = input("\n是否继续？(yes/no): ").strip().lower()
    
    if confirm != 'yes':
        print("\n❌ 已取消清理")
        return
    
    # 删除数据
    print("\n开始清理...")
    
    admin_result = system_db['admin'].delete_many({})
    print(f"✅ 删除 admin: {admin_result.deleted_count} 条")
    
    role_result = system_db['role'].delete_many({})
    print(f"✅ 删除 role: {role_result.deleted_count} 条")
    
    print("\n" + "="*60)
    print("✅ 清理完成！")
    print("="*60)
    
    # 验证
    print("\n📊 验证结果:")
    admin_count = system_db['admin'].count_documents({})
    role_count = system_db['role'].count_documents({})
    print(f"  - admin: {admin_count} 条")
    print(f"  - role: {role_count} 条")

if __name__ == '__main__':
    try:
        cleanup_system_db()
    except Exception as e:
        print(f"\n❌ 清理失败: {e}")
        import traceback
        traceback.print_exc()
    finally:
        client.close()

