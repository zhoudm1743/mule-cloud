#!/usr/bin/env python3
"""
数据库级别租户隔离 - 数据迁移脚本（Python版本）
将现有的单库多租户数据迁移到多库架构

使用方法：
    py migrate_data.py [--host localhost] [--port 27017] [--db mule] [--dry-run]

参数说明：
    --host: MongoDB主机地址，默认 localhost
    --port: MongoDB端口，默认 27017
    --db: 源数据库名称，默认 mule
    --username: MongoDB用户名（可选）
    --password: MongoDB密码（可选）
    --dry-run: 只查看迁移计划，不实际执行
"""

import argparse
from pymongo import MongoClient
from collections import defaultdict
import sys

class DataMigrator:
    def __init__(self, host='localhost', port=27017, db_name='mule', 
                 username=None, password=None, dry_run=False):
        self.host = host
        self.port = port
        self.db_name = db_name
        self.username = username
        self.password = password
        self.dry_run = dry_run
        self.client = None
        self.source_db = None
        
    def connect(self):
        """连接MongoDB"""
        try:
            if self.username and self.password:
                connection_string = f"mongodb://{self.username}:{self.password}@{self.host}:{self.port}/"
            else:
                connection_string = f"mongodb://{self.host}:{self.port}/"
            
            self.client = MongoClient(connection_string, serverSelectionTimeoutMS=5000)
            # 测试连接
            self.client.admin.command('ping')
            self.source_db = self.client[self.db_name]
            print(f"✅ 成功连接到 MongoDB: {self.host}:{self.port}/{self.db_name}")
            return True
        except Exception as e:
            print(f"❌ 连接MongoDB失败: {e}")
            return False
    
    def analyze_data(self):
        """分析现有数据"""
        print("\n" + "="*60)
        print("  数据分析")
        print("="*60 + "\n")
        
        # 需要迁移的集合
        collections_to_migrate = ['admin', 'role', 'menu', 'basic']
        
        stats = {}
        for coll_name in collections_to_migrate:
            collection = self.source_db[coll_name]
            
            # 统计总数
            total_count = collection.count_documents({})
            
            # 按租户统计
            pipeline = [
                {"$group": {"_id": "$tenant_id", "count": {"$sum": 1}}}
            ]
            tenant_stats = list(collection.aggregate(pipeline))
            
            stats[coll_name] = {
                'total': total_count,
                'tenants': {str(t['_id']): t['count'] for t in tenant_stats}
            }
            
            print(f"📊 {coll_name} 集合:")
            print(f"   总记录数: {total_count}")
            if tenant_stats:
                print(f"   租户数量: {len(tenant_stats)}")
                for tenant in tenant_stats:
                    tenant_id = tenant['_id'] if tenant['_id'] else '(空/系统)'
                    print(f"     - 租户 {tenant_id}: {tenant['count']} 条记录")
            print()
        
        return stats
    
    def get_all_tenant_ids(self):
        """获取所有租户ID"""
        tenant_ids = set()
        
        # 从admin集合获取租户ID
        for doc in self.source_db['admin'].find({'tenant_id': {'$ne': None, '$ne': ''}}):
            if 'tenant_id' in doc and doc['tenant_id']:
                tenant_ids.add(doc['tenant_id'])
        
        # 从role集合获取租户ID
        for doc in self.source_db['role'].find({'tenant_id': {'$ne': None, '$ne': ''}}):
            if 'tenant_id' in doc and doc['tenant_id']:
                tenant_ids.add(doc['tenant_id'])
        
        # 从basic集合获取租户ID
        for doc in self.source_db['basic'].find({'tenant_id': {'$ne': None, '$ne': ''}}):
            if 'tenant_id' in doc and doc['tenant_id']:
                tenant_ids.add(doc['tenant_id'])
        
        return list(tenant_ids)
    
    def migrate_collection(self, collection_name, tenant_id, target_db_name):
        """迁移单个集合的数据"""
        source_collection = self.source_db[collection_name]
        target_db = self.client[target_db_name]
        target_collection = target_db[collection_name]
        
        # 查询该租户的数据
        query = {'tenant_id': tenant_id} if tenant_id else {'$or': [{'tenant_id': None}, {'tenant_id': ''}]}
        
        documents = list(source_collection.find(query))
        
        if not documents:
            return 0
        
        # 删除tenant_id字段
        for doc in documents:
            if 'tenant_id' in doc:
                del doc['tenant_id']
        
        # 插入到目标数据库
        if not self.dry_run:
            try:
                result = target_collection.insert_many(documents, ordered=False)
                return len(result.inserted_ids)
            except Exception as e:
                print(f"      ⚠️  警告: 部分数据插入失败 ({e})")
                return len(documents)
        else:
            return len(documents)
    
    def migrate_tenant_data(self, tenant_id):
        """迁移单个租户的所有数据"""
        target_db_name = f"{self.db_name}_{tenant_id}"
        
        print(f"\n  迁移租户: {tenant_id}")
        print(f"  目标数据库: {target_db_name}")
        print("  " + "-"*50)
        
        collections = ['admin', 'role', 'menu', 'basic']
        total_migrated = 0
        
        for coll_name in collections:
            count = self.migrate_collection(coll_name, tenant_id, target_db_name)
            if count > 0:
                status = "✅" if not self.dry_run else "📋"
                action = "迁移" if not self.dry_run else "将迁移"
                print(f"    {status} {coll_name}: {action} {count} 条记录")
                total_migrated += count
        
        return total_migrated
    
    def migrate_system_data(self):
        """迁移系统数据（无租户ID或租户ID为空）"""
        print(f"\n  迁移系统数据")
        print(f"  目标数据库: {self.db_name}（保持在系统库）")
        print("  " + "-"*50)
        
        collections = ['admin', 'role', 'menu']
        total_migrated = 0
        
        for coll_name in collections:
            source_collection = self.source_db[coll_name]
            
            # 查询系统数据（无tenant_id或为空）
            query = {'$or': [{'tenant_id': None}, {'tenant_id': ''}]}
            count = source_collection.count_documents(query)
            
            if count > 0:
                status = "📋"
                print(f"    {status} {coll_name}: {count} 条系统记录（保持原位）")
                total_migrated += count
        
        # tenant集合始终在系统库
        tenant_count = self.source_db['tenant'].count_documents({})
        print(f"    📋 tenant: {tenant_count} 条记录（保持原位）")
        total_migrated += tenant_count
        
        return total_migrated
    
    def create_indexes(self, db_name):
        """为数据库创建索引"""
        if self.dry_run:
            return
        
        db = self.client[db_name]
        
        # admin集合索引
        try:
            db['admin'].create_index('phone', unique=True)
            db['admin'].create_index([('is_deleted', 1), ('status', 1)])
        except:
            pass
        
        # role集合索引
        try:
            db['role'].create_index('code', unique=True)
            db['role'].create_index([('is_deleted', 1), ('status', 1)])
        except:
            pass
        
        # menu集合索引
        try:
            db['menu'].create_index('name', unique=True)
            db['menu'].create_index([('is_deleted', 1), ('status', 1)])
        except:
            pass
        
        # basic集合索引
        try:
            db['basic'].create_index([('type', 1), ('is_deleted', 1)])
        except:
            pass
    
    def migrate(self):
        """执行完整的数据迁移"""
        if not self.connect():
            return False
        
        try:
            # 1. 分析数据
            stats = self.analyze_data()
            
            # 2. 获取所有租户ID
            tenant_ids = self.get_all_tenant_ids()
            
            if not tenant_ids:
                print("⚠️  未找到任何租户数据，只有系统数据")
                print()
            else:
                print(f"📋 发现 {len(tenant_ids)} 个租户需要迁移\n")
            
            # 确认
            if not self.dry_run:
                print("="*60)
                print("⚠️  即将开始数据迁移！")
                print("="*60)
                response = input("\n确认开始迁移? (输入 yes 继续): ")
                if response.lower() != 'yes':
                    print("❌ 迁移已取消")
                    return False
            else:
                print("="*60)
                print("📋 模拟运行模式 (Dry Run) - 不会实际修改数据")
                print("="*60)
            
            print("\n" + "="*60)
            print("  开始数据迁移")
            print("="*60)
            
            total_migrated = 0
            
            # 3. 迁移每个租户的数据
            for tenant_id in tenant_ids:
                count = self.migrate_tenant_data(tenant_id)
                total_migrated += count
                
                # 创建索引
                if not self.dry_run:
                    target_db_name = f"{self.db_name}_{tenant_id}"
                    print(f"    🔧 创建索引...")
                    self.create_indexes(target_db_name)
            
            # 4. 处理系统数据
            system_count = self.migrate_system_data()
            
            # 5. 完成
            print("\n" + "="*60)
            if not self.dry_run:
                print("  ✅ 数据迁移完成！")
            else:
                print("  📋 迁移计划分析完成（模拟运行）")
            print("="*60)
            print(f"\n  租户数量: {len(tenant_ids)}")
            print(f"  迁移记录: {total_migrated} 条")
            print(f"  系统记录: {system_count} 条")
            
            if not self.dry_run:
                print(f"\n  新建数据库: {len(tenant_ids)} 个")
                for tenant_id in tenant_ids:
                    print(f"    - {self.db_name}_{tenant_id}")
                
                print("\n  📚 下一步:")
                print("    1. 验证数据迁移结果")
                print("    2. 启动应用服务测试")
                print("    3. 确认无误后删除旧数据中的tenant_id字段（可选）")
            else:
                print("\n  💡 提示:")
                print("    移除 --dry-run 参数执行实际迁移")
            
            print()
            return True
            
        except Exception as e:
            print(f"\n❌ 迁移过程出错: {e}")
            import traceback
            traceback.print_exc()
            return False
        finally:
            if self.client:
                self.client.close()

def main():
    parser = argparse.ArgumentParser(
        description='数据库级别租户隔离 - 数据迁移工具',
        formatter_class=argparse.RawDescriptionHelpFormatter
    )
    
    parser.add_argument('--host', default='localhost', help='MongoDB主机地址 (默认: localhost)')
    parser.add_argument('--port', type=int, default=27017, help='MongoDB端口 (默认: 27017)')
    parser.add_argument('--db', default='mule', help='源数据库名称 (默认: mule)')
    parser.add_argument('--username', help='MongoDB用户名（可选）')
    parser.add_argument('--password', help='MongoDB密码（可选）')
    parser.add_argument('--dry-run', action='store_true', help='只查看迁移计划，不实际执行')
    
    args = parser.parse_args()
    
    print("="*60)
    print("  数据库级别租户隔离 - 数据迁移工具")
    print("="*60)
    print(f"\n  MongoDB: {args.host}:{args.port}")
    print(f"  数据库: {args.db}")
    if args.dry_run:
        print(f"  模式: 模拟运行 (Dry Run)")
    else:
        print(f"  模式: 实际迁移")
    print()
    
    migrator = DataMigrator(
        host=args.host,
        port=args.port,
        db_name=args.db,
        username=args.username,
        password=args.password,
        dry_run=args.dry_run
    )
    
    success = migrator.migrate()
    sys.exit(0 if success else 1)

if __name__ == '__main__':
    main()

