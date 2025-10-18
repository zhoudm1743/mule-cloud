"""
修复扎号格式，将个位数补0
例如：1 -> 01, 2 -> 02
"""
import pymongo

# 连接配置
MONGO_HOST = "localhost"
MONGO_PORT = 27017
# 可能需要根据实际情况修改数据库名
DATABASES = ["ace", "system"]  # 添加你的租户数据库名

def fix_bundle_no_format():
    """修复扎号格式"""
    client = pymongo.MongoClient(f"mongodb://{MONGO_HOST}:{MONGO_PORT}/")
    
    for db_name in DATABASES:
        print(f"\n处理数据库: {db_name}")
        db = client[db_name]
        
        # 修复裁剪批次 (cutting_batches)
        print(f"  修复 cutting_batches...")
        batches_collection = db["cutting_batches"]
        batches = batches_collection.find({"is_deleted": 0})
        
        batches_updated = 0
        for batch in batches:
            bundle_no = batch.get("bundle_no", "")
            try:
                bundle_int = int(bundle_no)
                if bundle_int < 100:
                    formatted_bundle_no = f"{bundle_int:02d}"
                    if formatted_bundle_no != bundle_no:
                        batches_collection.update_one(
                            {"_id": batch["_id"]},
                            {"$set": {"bundle_no": formatted_bundle_no}}
                        )
                        batches_updated += 1
                        print(f"    批次 {batch['_id']}: {bundle_no} -> {formatted_bundle_no}")
            except ValueError:
                # 如果不是数字，跳过
                pass
        
        print(f"  cutting_batches 更新了 {batches_updated} 条记录")
        
        # 修复裁片监控 (cutting_pieces)
        print(f"  修复 cutting_pieces...")
        pieces_collection = db["cutting_pieces"]
        pieces = pieces_collection.find({})
        
        pieces_updated = 0
        for piece in pieces:
            bundle_no = piece.get("bundle_no", "")
            try:
                bundle_int = int(bundle_no)
                if bundle_int < 100:
                    formatted_bundle_no = f"{bundle_int:02d}"
                    if formatted_bundle_no != bundle_no:
                        pieces_collection.update_one(
                            {"_id": piece["_id"]},
                            {"$set": {"bundle_no": formatted_bundle_no}}
                        )
                        pieces_updated += 1
                        print(f"    裁片 {piece['_id']}: {bundle_no} -> {formatted_bundle_no}")
            except ValueError:
                # 如果不是数字，跳过
                pass
        
        print(f"  cutting_pieces 更新了 {pieces_updated} 条记录")
    
    client.close()
    print("\n✅ 扎号格式修复完成！")

if __name__ == "__main__":
    print("🔧 开始修复扎号格式...")
    print("=" * 60)
    fix_bundle_no_format()

