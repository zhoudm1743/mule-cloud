#!/usr/bin/env python3
"""
员工档案数据查询脚本

用途：查询和查看员工档案数据
用法：py query_member_data.py

功能：
1. 查看所有租户数据库
2. 查看指定租户的员工列表
3. 查看指定员工的完整档案
4. 统计员工数据
"""

from pymongo import MongoClient
import sys
import json
from datetime import datetime

# MongoDB连接配置
MONGO_URI = "mongodb://root:bgg8384495@127.0.0.1:27015/admin"

def connect_mongodb():
    """连接MongoDB"""
    try:
        client = MongoClient(MONGO_URI)
        client.admin.command('ping')
        return client
    except Exception as e:
        print(f"❌ MongoDB连接失败: {e}")
        sys.exit(1)

def list_tenant_databases(client):
    """列出所有租户数据库"""
    databases = client.list_database_names()
    tenant_dbs = [db for db in databases if db.startswith('tenant_') and db != 'tenant_system']
    
    print("\n" + "=" * 60)
    print(f"📊 租户数据库列表 (共 {len(tenant_dbs)} 个)")
    print("=" * 60)
    
    for i, db_name in enumerate(tenant_dbs, 1):
        db = client[db_name]
        member_count = db['member'].count_documents({'is_deleted': 0})
        print(f"{i}. {db_name} - {member_count} 名员工")
    
    return tenant_dbs

def list_members(db, limit=10):
    """列出员工列表"""
    collection = db['member']
    
    members = collection.find(
        {'is_deleted': 0},
        {
            '_id': 1,
            'name': 1,
            'job_number': 1,
            'department': 1,
            'position': 1,
            'phone': 1,
            'status': 1
        }
    ).limit(limit)
    
    print("\n" + "=" * 60)
    print(f"员工列表 (最多显示{limit}条)")
    print("=" * 60)
    print(f"{'序号':<4} {'工号':<10} {'姓名':<10} {'部门':<12} {'岗位':<12} {'状态':<8}")
    print("-" * 60)
    
    for i, member in enumerate(members, 1):
        name = member.get('name', '未设置')
        job_number = member.get('job_number', '未设置')
        department = member.get('department', '未设置')
        position = member.get('position', '未设置')
        status = member.get('status', 'active')
        
        print(f"{i:<4} {job_number:<10} {name:<10} {department:<12} {position:<12} {status:<8}")

def view_member_detail(db, job_number):
    """查看员工详细信息"""
    collection = db['member']
    
    member = collection.find_one({
        'job_number': job_number,
        'is_deleted': 0
    })
    
    if not member:
        print(f"\n❌ 未找到工号为 {job_number} 的员工")
        return
    
    print("\n" + "=" * 60)
    print(f"员工档案详情 - {member.get('name', '未设置')}")
    print("=" * 60)
    
    # 基本信息
    print("\n【基本信息】")
    print(f"  工号: {member.get('job_number', '未设置')}")
    print(f"  姓名: {member.get('name', '未设置')}")
    print(f"  性别: {get_gender_text(member.get('gender', 0))}")
    print(f"  身份证号: {member.get('id_card_no', '未填写')}")
    if member.get('birthday', 0) > 0:
        birthday_str = datetime.fromtimestamp(member['birthday']).strftime('%Y-%m-%d')
        print(f"  出生日期: {birthday_str}")
        print(f"  年龄: {member.get('age', 0)}岁")
    print(f"  民族: {member.get('nation', '未设置')}")
    print(f"  籍贯: {member.get('native_place', '未设置')}")
    print(f"  婚姻状况: {get_marital_status_text(member.get('marital_status', 'single'))}")
    print(f"  政治面貌: {get_political_text(member.get('political', 'masses'))}")
    print(f"  学历: {get_education_text(member.get('education', ''))}")
    
    # 联系信息
    print("\n【联系信息】")
    print(f"  手机号: {member.get('phone', '未设置')}")
    print(f"  邮箱: {member.get('email', '未设置')}")
    print(f"  住址: {member.get('address', '未设置')}")
    print(f"  紧急联系人: {member.get('emergency_contact', '未设置')}")
    print(f"  紧急电话: {member.get('emergency_phone', '未设置')}")
    print(f"  关系: {member.get('emergency_relation', '未设置')}")
    
    # 企业信息
    print("\n【企业信息】")
    print(f"  部门: {member.get('department', '未设置')}")
    print(f"  岗位: {member.get('position', '未设置')}")
    print(f"  车间: {member.get('workshop', '未设置')}")
    print(f"  班组: {member.get('team', '未设置')}")
    print(f"  班组长: {member.get('team_leader', '未设置')}")
    
    # 工作相关
    print("\n【工作相关】")
    if member.get('employed_at', 0) > 0:
        employed_str = datetime.fromtimestamp(member['employed_at']).strftime('%Y-%m-%d')
        print(f"  入职日期: {employed_str}")
        print(f"  工龄: {member.get('work_years', 0)}年{member.get('work_months', 0)}个月")
    print(f"  合同类型: {get_contract_type_text(member.get('contract_type', 'fulltime'))}")
    print(f"  状态: {get_status_text(member.get('status', 'active'))}")
    
    # 技能
    skills = member.get('skills', [])
    if skills:
        print("\n【技能列表】")
        for skill in skills:
            print(f"  - {skill.get('name')} ({get_skill_level_text(skill.get('level', 'beginner'))})")
    
    # 薪资信息
    print("\n【薪资信息】")
    print(f"  薪资类型: {get_salary_type_text(member.get('salary_type', 'piece'))}")
    if member.get('base_salary', 0) > 0:
        print(f"  基本工资: {member.get('base_salary', 0)}元/月")
    if member.get('piece_rate', 0) > 0:
        print(f"  计件单价: {member.get('piece_rate', 0)}元/件")

def statistics(db):
    """统计信息"""
    collection = db['member']
    
    print("\n" + "=" * 60)
    print("统计信息")
    print("=" * 60)
    
    # 总人数
    total = collection.count_documents({'is_deleted': 0})
    print(f"\n总员工数: {total}")
    
    # 在职人数
    active_count = collection.count_documents({'status': 'active', 'is_deleted': 0})
    print(f"在职人数: {active_count}")
    
    # 离职人数
    inactive_count = collection.count_documents({'status': 'inactive', 'is_deleted': 0})
    print(f"离职人数: {inactive_count}")
    
    # 按部门统计
    pipeline = [
        {'$match': {'is_deleted': 0, 'status': 'active'}},
        {'$group': {'_id': '$department', 'count': {'$sum': 1}}},
        {'$sort': {'count': -1}}
    ]
    dept_stats = list(collection.aggregate(pipeline))
    
    if dept_stats:
        print("\n按部门统计：")
        for stat in dept_stats:
            dept = stat['_id'] or '未设置'
            count = stat['count']
            print(f"  {dept}: {count}人")
    
    # 按学历统计
    pipeline = [
        {'$match': {'is_deleted': 0, 'status': 'active'}},
        {'$group': {'_id': '$education', 'count': {'$sum': 1}}},
        {'$sort': {'count': -1}}
    ]
    edu_stats = list(collection.aggregate(pipeline))
    
    if edu_stats:
        print("\n按学历统计：")
        for stat in edu_stats:
            edu = stat['_id'] or '未填写'
            count = stat['count']
            edu_text = get_education_text(edu) if edu else '未填写'
            print(f"  {edu_text}: {count}人")

# 辅助函数
def get_gender_text(gender):
    return {0: '未知', 1: '男', 2: '女'}.get(gender, '未知')

def get_marital_status_text(status):
    return {
        'single': '未婚',
        'married': '已婚',
        'divorced': '离异'
    }.get(status, status)

def get_political_text(political):
    return {
        'party': '党员',
        'league': '团员',
        'masses': '群众'
    }.get(political, political)

def get_education_text(education):
    return {
        'primary': '小学',
        'middle': '初中',
        'high': '高中',
        'college': '大专',
        'bachelor': '本科',
        'master': '硕士',
        'doctor': '博士'
    }.get(education, education or '未填写')

def get_contract_type_text(contract_type):
    return {
        'fulltime': '全职',
        'parttime': '兼职',
        'intern': '实习',
        'dispatch': '劳务派遣'
    }.get(contract_type, contract_type)

def get_status_text(status):
    return {
        'active': '在职',
        'probation': '试用期',
        'inactive': '离职',
        'suspended': '停职'
    }.get(status, status)

def get_skill_level_text(level):
    return {
        'beginner': '初级',
        'intermediate': '中级',
        'advanced': '高级',
        'expert': '专家'
    }.get(level, level)

def get_salary_type_text(salary_type):
    return {
        'hourly': '计时',
        'piece': '计件',
        'monthly': '月薪',
        'mixed': '混合'
    }.get(salary_type, salary_type)

def main():
    """主函数"""
    print("=" * 60)
    print("🔍 员工档案数据查询工具")
    print("=" * 60)
    
    client = connect_mongodb()
    
    while True:
        print("\n请选择操作：")
        print("1. 查看所有租户数据库")
        print("2. 查看指定租户的员工列表")
        print("3. 查看指定员工的详细信息")
        print("4. 统计信息")
        print("0. 退出")
        
        choice = input("\n请输入选项: ").strip()
        
        if choice == '0':
            print("👋 再见！")
            break
        elif choice == '1':
            list_tenant_databases(client)
        elif choice == '2':
            db_name = input("请输入数据库名称 (如 tenant_ace): ").strip()
            if db_name in client.list_database_names():
                db = client[db_name]
                list_members(db, limit=20)
            else:
                print(f"❌ 数据库 {db_name} 不存在")
        elif choice == '3':
            db_name = input("请输入数据库名称 (如 tenant_ace): ").strip()
            if db_name in client.list_database_names():
                db = client[db_name]
                job_number = input("请输入员工工号: ").strip()
                view_member_detail(db, job_number)
            else:
                print(f"❌ 数据库 {db_name} 不存在")
        elif choice == '4':
            db_name = input("请输入数据库名称 (如 tenant_ace): ").strip()
            if db_name in client.list_database_names():
                db = client[db_name]
                statistics(db)
            else:
                print(f"❌ 数据库 {db_name} 不存在")
        else:
            print("❌ 无效的选项")
    
    client.close()

if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n👋 再见！")
        sys.exit(0)
    except Exception as e:
        print(f"\n\n❌ 发生错误: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(1)

