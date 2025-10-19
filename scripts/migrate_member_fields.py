#!/usr/bin/env python3
"""
员工档案字段迁移脚本

用途：为现有的member集合添加新字段（47个新字段）
用法：py migrate_member_fields.py

注意：
1. 请确保已安装pymongo: pip install pymongo
2. 请先备份数据！
3. 本脚本会对所有租户数据库执行迁移
"""

from pymongo import MongoClient, ASCENDING
from datetime import datetime
import sys

# MongoDB连接配置
MONGO_URI = "mongodb://root:bgg8384495@127.0.0.1:27015/admin"

def connect_mongodb():
    """连接MongoDB"""
    try:
        client = MongoClient(MONGO_URI)
        # 测试连接
        client.admin.command('ping')
        print("✅ MongoDB连接成功")
        return client
    except Exception as e:
        print(f"❌ MongoDB连接失败: {e}")
        sys.exit(1)

def get_tenant_databases(client):
    """获取所有租户数据库"""
    databases = client.list_database_names()
    tenant_dbs = [db for db in databases if db.startswith('tenant_') and db != 'tenant_system']
    print(f"📊 找到 {len(tenant_dbs)} 个租户数据库")
    return tenant_dbs

def migrate_member_collection(db):
    """迁移member集合"""
    collection = db['member']
    
    # 统计现有记录
    total_count = collection.count_documents({'is_deleted': 0})
    print(f"  📝 找到 {total_count} 条员工记录")
    
    if total_count == 0:
        print("  ⚠️  没有需要迁移的记录")
        return True
    
    # 构建更新数据（新字段默认值）
    update_data = {
        '$set': {
            # 个人基本信息
            'name_pinyin': '',
            'id_card_type': 'idcard',
            'id_card_no': '',
            'birthday': 0,
            'age': 0,
            'nation': '汉族',
            'native_place': '',
            'marital_status': 'single',
            'political': 'masses',
            'education': '',
            'photo': '',
            
            # 联系信息
            'email': '',
            'address': '',
            'emergency_contact': '',
            'emergency_phone': '',
            'emergency_relation': '',
            
            # 企业信息扩展
            'department_id': '',
            'position_id': '',
            'workshop': '',
            'workshop_id': '',
            'team': '',
            'team_id': '',
            'team_leader': '',
            
            # 工作相关
            'regular_at': 0,
            'contract_type': 'fulltime',  # 默认全职
            'contract_start_at': 0,
            'contract_end_at': 0,
            'work_years': 0,
            'work_months': 0,
            
            # 技能与资质
            'skills': [],
            'certificates': [],
            
            # 薪资信息
            'salary_type': 'piece',  # 服装厂默认计件
            'base_salary': 0,
            'hourly_rate': 0,
            'piece_rate': 0,
            'bank_name': '',
            'bank_account': '',
            'bank_account_name': '',
            
            # 状态扩展
            'left_reason': '',
            
            # 备注
            'remark': ''
        }
    }
    
    try:
        # 批量更新所有记录
        result = collection.update_many(
            {'is_deleted': 0},
            update_data
        )
        print(f"  ✅ 更新完成: {result.modified_count} 条记录")
        return True
    except Exception as e:
        print(f"  ❌ 更新失败: {e}")
        return False

def create_indexes(db):
    """创建索引"""
    collection = db['member']
    
    try:
        # 工号唯一索引
        collection.create_index([('job_number', ASCENDING)], unique=True, sparse=True, background=True)
        print("  ✅ 创建job_number唯一索引")
        
        # 手机号索引
        collection.create_index([('phone', ASCENDING)], background=True)
        print("  ✅ 创建phone索引")
        
        # 身份证号唯一索引
        collection.create_index([('id_card_no', ASCENDING)], unique=True, sparse=True, background=True)
        print("  ✅ 创建id_card_no唯一索引")
        
        # 部门+状态复合索引
        collection.create_index([('department', ASCENDING), ('status', ASCENDING)], background=True)
        print("  ✅ 创建department+status复合索引")
        
        # 技能工序索引
        collection.create_index([('skills.process_ids', ASCENDING)], background=True)
        print("  ✅ 创建skills.process_ids索引")
        
        return True
    except Exception as e:
        print(f"  ⚠️  创建索引部分失败: {e}")
        return False

def verify_migration(db):
    """验证迁移结果"""
    collection = db['member']
    
    # 随机抽取一条记录检查
    sample = collection.find_one({'is_deleted': 0})
    if not sample:
        print("  ⚠️  没有记录可验证")
        return True
    
    # 检查新字段是否存在
    new_fields = ['name_pinyin', 'id_card_no', 'nation', 'education', 'skills', 'salary_type']
    missing_fields = [field for field in new_fields if field not in sample]
    
    if missing_fields:
        print(f"  ❌ 缺少字段: {', '.join(missing_fields)}")
        return False
    else:
        print("  ✅ 数据验证通过")
        return True

def main():
    """主函数"""
    print("=" * 60)
    print("🚀 员工档案字段迁移脚本")
    print("=" * 60)
    print()
    
    # 连接MongoDB
    client = connect_mongodb()
    
    # 获取租户数据库列表
    tenant_dbs = get_tenant_databases(client)
    
    if not tenant_dbs:
        print("⚠️  没有找到租户数据库，脚本退出")
        return
    
    print()
    print("⚠️  警告：请确保已备份数据！")
    print("⚠️  本脚本将修改以下数据库的member集合：")
    for db_name in tenant_dbs:
        print(f"   - {db_name}")
    print()
    
    # 确认执行
    confirm = input("是否继续执行迁移？(yes/no): ").strip().lower()
    if confirm != 'yes':
        print("❌ 用户取消迁移")
        return
    
    print()
    print("=" * 60)
    print("开始迁移...")
    print("=" * 60)
    print()
    
    success_count = 0
    failed_count = 0
    
    # 对每个租户数据库执行迁移
    for db_name in tenant_dbs:
        print(f"📦 处理数据库: {db_name}")
        db = client[db_name]
        
        # 1. 迁移字段
        if not migrate_member_collection(db):
            failed_count += 1
            continue
        
        # 2. 创建索引
        create_indexes(db)
        
        # 3. 验证迁移
        if verify_migration(db):
            success_count += 1
        else:
            failed_count += 1
        
        print()
    
    # 汇总结果
    print("=" * 60)
    print("迁移完成！")
    print("=" * 60)
    print(f"✅ 成功: {success_count} 个数据库")
    print(f"❌ 失败: {failed_count} 个数据库")
    print()
    
    if failed_count == 0:
        print("🎉 所有数据库迁移成功！")
        print()
        print("后续步骤：")
        print("1. 在小程序测试以下功能：")
        print("   - 查看个人档案")
        print("   - 编辑基本信息")
        print("   - 编辑联系信息")
        print("2. 如果有问题，可以从备份恢复")
    else:
        print("⚠️  部分数据库迁移失败，请检查错误日志")
    
    # 关闭连接
    client.close()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n❌ 用户中断执行")
        sys.exit(1)
    except Exception as e:
        print(f"\n\n❌ 发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

