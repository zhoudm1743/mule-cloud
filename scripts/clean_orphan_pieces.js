// MongoDB清理脚本：删除孤立的裁片监控记录
// 使用方法：
// 1. 连接到MongoDB
// 2. 切换到对应的租户数据库：use mule_ace（或你的租户代码）
// 3. 运行此脚本：load('scripts/clean_orphan_pieces.js')

// 获取所有裁片监控记录
const pieces = db.cutting_pieces.find({}).toArray();

print(`共找到 ${pieces.length} 条裁片监控记录`);

let orphanCount = 0;
let deletedIds = [];

// 检查每条裁片监控记录
pieces.forEach(piece => {
    // 查找对应的批次记录
    const batch = db.cutting_batches.findOne({
        bed_no: piece.bed_no,
        bundle_no: piece.bundle_no,
        is_deleted: 0  // 未删除的批次
    });
    
    // 如果找不到对应的批次，或者批次已被删除
    if (!batch) {
        orphanCount++;
        deletedIds.push(piece._id);
        print(`发现孤立记录: 合同号=${piece.contract_no}, 床号=${piece.bed_no}, 扎号=${piece.bundle_no}`);
    }
});

print(`\n共发现 ${orphanCount} 条孤立的裁片监控记录`);

if (orphanCount > 0) {
    print('\n是否删除这些孤立记录？(y/n)');
    // 在实际执行时，注释掉上面这行，直接执行删除
    const result = db.cutting_pieces.deleteMany({
        _id: { $in: deletedIds }
    });
    
    print(`\n删除完成！删除了 ${result.deletedCount} 条记录`);
} else {
    print('\n没有发现孤立记录，无需清理');
}

