# Casbin 角色权限测试脚本 (PowerShell)

$BASE_URL = "http://localhost:8080/admin/system"

Write-Host "=== Casbin 角色权限系统测试 ===" -ForegroundColor Cyan
Write-Host ""

# 1. 创建租户
Write-Host "1. 创建测试租户..." -ForegroundColor Green
$tenantBody = @{
    name = "测试公司"
    code = "test_corp"
    contact = "测试联系人"
    phone = "13800138000"
    email = "test@example.com"
} | ConvertTo-Json

$tenantResponse = Invoke-RestMethod -Uri "$BASE_URL/tenants" -Method POST -Body $tenantBody -ContentType "application/json"
$tenantResponse | ConvertTo-Json -Depth 10
$TENANT_ID = $tenantResponse.data.id
Write-Host "租户ID: $TENANT_ID" -ForegroundColor Yellow
Write-Host ""

# 2. 创建角色
Write-Host "2. 创建测试角色..." -ForegroundColor Green
$roleBody = @{
    tenant_id = $TENANT_ID
    name = "系统管理员"
    code = "sys_admin"
    description = "系统管理员角色"
    menus = @()
} | ConvertTo-Json

$roleResponse = Invoke-RestMethod -Uri "$BASE_URL/roles" -Method POST -Body $roleBody -ContentType "application/json"
$roleResponse | ConvertTo-Json -Depth 10
$ROLE_ID = $roleResponse.data.id
Write-Host "角色ID: $ROLE_ID" -ForegroundColor Yellow
Write-Host ""

# 3. 创建菜单1
Write-Host "3. 创建测试菜单..." -ForegroundColor Green
$menu1Body = @{
    name = "dashboard"
    path = "/dashboard"
    title = "仪表盘"
    icon = "icon-park-outline:analysis"
    requiresAuth = $true
    menuType = "dir"
} | ConvertTo-Json

$menu1Response = Invoke-RestMethod -Uri "$BASE_URL/menus" -Method POST -Body $menu1Body -ContentType "application/json"
$menu1Response | ConvertTo-Json -Depth 10
$MENU1_ID = $menu1Response.data.id
Write-Host "菜单1 ID: $MENU1_ID" -ForegroundColor Yellow
Write-Host ""

# 创建菜单2
$menu2Body = @{
    pid = $MENU1_ID
    name = "dashboard_workbench"
    path = "/dashboard/workbench"
    title = "工作台"
    componentPath = "/dashboard/workbench/index.vue"
    icon = "icon-park-outline:alarm"
    requiresAuth = $true
    menuType = "page"
} | ConvertTo-Json

$menu2Response = Invoke-RestMethod -Uri "$BASE_URL/menus" -Method POST -Body $menu2Body -ContentType "application/json"
$menu2Response | ConvertTo-Json -Depth 10
$MENU2_ID = $menu2Response.data.id
Write-Host "菜单2 ID: $MENU2_ID" -ForegroundColor Yellow
Write-Host ""

# 4. 分配菜单给角色
Write-Host "4. 分配菜单权限给角色..." -ForegroundColor Green
$assignMenuBody = @{
    menus = @($MENU1_ID, $MENU2_ID)
} | ConvertTo-Json

$assignResponse = Invoke-RestMethod -Uri "$BASE_URL/roles/$ROLE_ID/menus" -Method POST -Body $assignMenuBody -ContentType "application/json"
$assignResponse | ConvertTo-Json -Depth 10
Write-Host ""

# 5. 获取角色的菜单
Write-Host "5. 查询角色的菜单权限..." -ForegroundColor Green
$roleMenus = Invoke-RestMethod -Uri "$BASE_URL/roles/$ROLE_ID/menus" -Method GET
$roleMenus | ConvertTo-Json -Depth 10
Write-Host ""

# 6. 创建管理员
Write-Host "6. 创建测试管理员..." -ForegroundColor Green
$adminBody = @{
    phone = "13900139000"
    password = "123456"
    nickname = "测试管理员"
    email = "admin@test.com"
    status = 1
} | ConvertTo-Json

$adminResponse = Invoke-RestMethod -Uri "$BASE_URL/admins" -Method POST -Body $adminBody -ContentType "application/json"
$adminResponse | ConvertTo-Json -Depth 10
$ADMIN_ID = $adminResponse.data.id
Write-Host "管理员ID: $ADMIN_ID" -ForegroundColor Yellow
Write-Host ""

# 7. 分配角色给管理员
Write-Host "7. 分配角色给管理员..." -ForegroundColor Green
$assignRoleBody = @{
    roles = @($ROLE_ID)
} | ConvertTo-Json

$assignRoleResponse = Invoke-RestMethod -Uri "$BASE_URL/admins/$ADMIN_ID/roles" -Method POST -Body $assignRoleBody -ContentType "application/json"
$assignRoleResponse | ConvertTo-Json -Depth 10
Write-Host ""

# 8. 获取管理员的角色
Write-Host "8. 查询管理员的角色..." -ForegroundColor Green
$adminRoles = Invoke-RestMethod -Uri "$BASE_URL/admins/$ADMIN_ID/roles" -Method GET
$adminRoles | ConvertTo-Json -Depth 10
Write-Host ""

# 9. 查询角色列表
Write-Host "9. 查询角色列表..." -ForegroundColor Green
$roleList = Invoke-RestMethod -Uri "$BASE_URL/roles?tenant_id=$TENANT_ID&page=1&page_size=10" -Method GET
$roleList | ConvertTo-Json -Depth 10
Write-Host ""

# 10. 查询菜单列表
Write-Host "10. 查询所有菜单..." -ForegroundColor Green
$menuList = Invoke-RestMethod -Uri "$BASE_URL/menus/all" -Method GET
$menuList | ConvertTo-Json -Depth 10
Write-Host ""

Write-Host "=== 测试完成 ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "创建的资源ID：" -ForegroundColor Cyan
Write-Host "  租户ID: $TENANT_ID"
Write-Host "  角色ID: $ROLE_ID"
Write-Host "  菜单1 ID: $MENU1_ID"
Write-Host "  菜单2 ID: $MENU2_ID"
Write-Host "  管理员ID: $ADMIN_ID"

