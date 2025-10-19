# PowerShell脚本：批量更新API路径前缀

$orderFile = "frontend/src/service/api/order.ts"

# 读取文件内容
$content = Get-Content $orderFile -Raw

# 替换 /order/ 为 /admin/order/
$content = $content -replace "'/order/", "'/admin/order/"
$content = $content -replace '`/order/', '`/admin/order/'

# 写回文件
$content | Set-Content $orderFile -NoNewline

Write-Host "✅ 已更新 $orderFile 中的所有 API 路径"

