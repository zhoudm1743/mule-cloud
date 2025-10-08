# Mule-Cloud 微服务启动脚本
$Host.UI.RawUI.WindowTitle = "Mule-Cloud Services Manager"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Mule-Cloud Service Startup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 检查是否在项目根目录
if (!(Test-Path "go.mod")) {
    Write-Host "ERROR: Please run this script in project root directory!" -ForegroundColor Red
    Read-Host "Press Enter to exit"
    exit 1
}

# 服务定义
$services = @{
    "gateway" = @{ Name = "Gateway"; Path = "cmd/gateway/main.go"; Port = 8080; Color = "Green"; Process = $null }
    "auth"    = @{ Name = "Auth";    Path = "cmd/auth/main.go";    Port = 8081; Color = "Yellow"; Process = $null }
    "basic"   = @{ Name = "Basic";   Path = "cmd/basic/main.go";   Port = 8082; Color = "Blue"; Process = $null }
    "common"  = @{ Name = "Common";  Path = "cmd/common/main.go";  Port = 8083; Color = "Magenta"; Process = $null }
    "perms"   = @{ Name = "Perms";   Path = "cmd/perms/main.go";   Port = 8084; Color = "Cyan"; Process = $null }
    "order"   = @{ Name = "Order";   Path = "cmd/order/main.go";   Port = 8085; Color = "White"; Process = $null }
}

# 启动服务函数
function Start-Service {
    param([string]$serviceName)
    
    if (!$services.ContainsKey($serviceName)) {
        Write-Host "ERROR: Unknown service '$serviceName'" -ForegroundColor Red
        return
    }
    
    $svc = $services[$serviceName]
    
    # 如果已经在运行，先停止
    if ($svc.Process -and !$svc.Process.HasExited) {
        Write-Host "Service $($svc.Name) is already running, stopping first..." -ForegroundColor Yellow
        Stop-Service $serviceName
        Start-Sleep -Milliseconds 500
    }
    
    Write-Host "Starting $($svc.Name) Service (Port: $($svc.Port))..." -ForegroundColor $svc.Color
    try {
        $svc.Process = Start-Process -FilePath "go" -ArgumentList "run", $svc.Path -PassThru -WindowStyle Minimized
        Start-Sleep -Milliseconds 1500
        
        # 刷新进程状态
        try {
            $svc.Process.Refresh()
        } catch {
            # 忽略刷新错误
        }
        
        if ($svc.Process.HasExited) {
            Write-Host "  $($svc.Name) failed to start! (Exit code: $($svc.Process.ExitCode))" -ForegroundColor Red
            $svc.Process = $null
        } else {
            Write-Host "  $($svc.Name) started successfully! (PID: $($svc.Process.Id))" -ForegroundColor Green
        }
    } catch {
        Write-Host "  Failed to start $($svc.Name): $_" -ForegroundColor Red
        $svc.Process = $null
    }
}

# 停止服务函数
function Stop-Service {
    param([string]$serviceName)
    
    if (!$services.ContainsKey($serviceName)) {
        Write-Host "ERROR: Unknown service '$serviceName'" -ForegroundColor Red
        return
    }
    
    $svc = $services[$serviceName]
    
    if (!$svc.Process -or $svc.Process.HasExited) {
        Write-Host "$($svc.Name) is not running" -ForegroundColor Yellow
        return
    }
    
    try {
        $processId = $svc.Process.Id
        Write-Host "Stopping $($svc.Name) (PID: $processId)..." -ForegroundColor Gray
        
        # 使用 taskkill 停止进程树（包括所有子进程）
        taskkill /F /T /PID $processId 2>&1 | Out-Null
        Start-Sleep -Milliseconds 500
        
        $svc.Process = $null
        Write-Host "  $($svc.Name) stopped successfully!" -ForegroundColor Green
    } catch {
        Write-Host "  Failed to stop $($svc.Name): $_" -ForegroundColor Red
    }
}

# 显示服务状态
function Show-Status {
    Write-Host ""
    Write-Host "Service Status:" -ForegroundColor Yellow
    Write-Host "----------------------------------------" -ForegroundColor Gray
    foreach ($key in $services.Keys | Sort-Object) {
        $svc = $services[$key]
        $status = "STOPPED"
        $color = "Red"
        $processId = "N/A"
        
        if ($svc.Process -and !$svc.Process.HasExited) {
            $status = "RUNNING"
            $color = "Green"
            $processId = $svc.Process.Id
        }
        
        Write-Host ("  {0,-8} | {1,-8} | Port: {2} | PID: {3}" -f $svc.Name, $status, $svc.Port, $processId) -ForegroundColor $color
    }
    Write-Host "----------------------------------------" -ForegroundColor Gray
    Write-Host ""
}

# 显示帮助
function Show-Help {
    Write-Host ""
    Write-Host "Available Commands:" -ForegroundColor Yellow
    Write-Host "  status                  - Show service status" -ForegroundColor Gray
    Write-Host "  restart <service>       - Restart a service" -ForegroundColor Gray
    Write-Host "  stop <service>          - Stop a service" -ForegroundColor Gray
    Write-Host "  start <service>         - Start a service" -ForegroundColor Gray
    Write-Host "  help                    - Show this help" -ForegroundColor Gray
    Write-Host "  exit                    - Stop all services and exit" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Available Services:" -ForegroundColor Yellow
    Write-Host "  gateway, auth, basic, common, perms, order" -ForegroundColor Gray
    Write-Host ""
}

Write-Host "Starting all services..." -ForegroundColor Yellow
Write-Host ""

try {
    # 启动所有服务
    $index = 1
    foreach ($key in @("gateway", "auth", "basic", "common", "perms", "order")) {
        Write-Host "[$index/6] " -NoNewline
        Start-Service $key
        $index++
    }

    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "All services started!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Cyan
    
    Show-Status
    Show-Help
    
    # 命令循环
    while ($true) {
        $input = Read-Host "Command"
        $parts = $input -split '\s+', 2
        $cmd = $parts[0].ToLower()
        $arg = if ($parts.Length -gt 1) { $parts[1] } else { "" }
        
        switch ($cmd) {
            "exit" {
                Write-Host "Exiting..." -ForegroundColor Yellow
                break
            }
            "status" {
                Show-Status
            }
            "restart" {
                if ($arg -eq "") {
                    Write-Host "Usage: restart <service>" -ForegroundColor Yellow
                } elseif ($services.ContainsKey($arg)) {
                    Write-Host ""
                    Stop-Service $arg
                    Start-Sleep -Milliseconds 500
                    Start-Service $arg
                    Write-Host ""
                } else {
                    Write-Host "Unknown service: $arg" -ForegroundColor Red
                    Write-Host "Available: gateway, auth, basic, common, perms, order" -ForegroundColor Gray
                }
            }
            "stop" {
                if ($arg -eq "") {
                    Write-Host "Usage: stop <service>" -ForegroundColor Yellow
                } elseif ($services.ContainsKey($arg)) {
                    Write-Host ""
                    Stop-Service $arg
                    Write-Host ""
                } else {
                    Write-Host "Unknown service: $arg" -ForegroundColor Red
                }
            }
            "start" {
                if ($arg -eq "") {
                    Write-Host "Usage: start <service>" -ForegroundColor Yellow
                } elseif ($services.ContainsKey($arg)) {
                    Write-Host ""
                    Start-Service $arg
                    Write-Host ""
                } else {
                    Write-Host "Unknown service: $arg" -ForegroundColor Red
                }
            }
            "help" {
                Show-Help
            }
            "" {
                # 空输入，忽略
            }
            default {
                Write-Host "Unknown command: $cmd" -ForegroundColor Red
                Write-Host "Type 'help' for available commands" -ForegroundColor Gray
            }
        }
        
        if ($cmd -eq "exit") {
            break
        }
    }

} finally {
    # 确保无论如何都会执行清理
    Write-Host ""
    Write-Host "Stopping all services..." -ForegroundColor Yellow
    
    $stoppedCount = 0
    foreach ($key in $services.Keys) {
        $svc = $services[$key]
        if ($svc.Process -and !$svc.Process.HasExited) {
            try {
                $processId = $svc.Process.Id
                Write-Host "  Stopping $($svc.Name) (PID: $processId)..." -ForegroundColor Gray
                # 使用 taskkill 停止进程树
                taskkill /F /T /PID $processId 2>&1 | Out-Null
                Start-Sleep -Milliseconds 300
                $stoppedCount++
            } catch {
                Write-Host "  Failed to stop $($svc.Name)" -ForegroundColor Red
            }
        }
    }
    
    # 额外清理：查找并结束所有相关的 go.exe 进程
    Write-Host "Cleaning up remaining processes..." -ForegroundColor Gray
    Get-Process | Where-Object { 
        $_.ProcessName -eq "go" -or $_.ProcessName -eq "main" 
    } | ForEach-Object {
        try {
            $_.Kill()
            Write-Host "  Killed process: $($_.ProcessName) (PID: $($_.Id))" -ForegroundColor Gray
        } catch {
            # 忽略错误
        }
    }
    
    Write-Host ""
    Write-Host "All services stopped! (Stopped: $stoppedCount)" -ForegroundColor Green
    Write-Host ""
    Start-Sleep -Seconds 2
}
