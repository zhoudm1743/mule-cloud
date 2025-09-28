.PHONY: help build start stop restart logs test clean

# 默认目标
.DEFAULT_GOAL := help

# 颜色定义
RED = \033[31m
GREEN = \033[32m
YELLOW = \033[33m
BLUE = \033[34m
PURPLE = \033[35m
CYAN = \033[36m
WHITE = \033[37m
RESET = \033[0m

## help: 显示帮助信息
help:
	@echo "$(CYAN)信芙云服装生产管理系统 - 开发工具$(RESET)"
	@echo "$(CYAN)=====================================$(RESET)"
	@echo ""
	@echo "$(GREEN)可用命令:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(YELLOW)%-15s$(RESET) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

## build: 构建所有Docker镜像
build:
	@echo "$(BLUE)构建Docker镜像...$(RESET)"
	docker-compose build

## start: 启动所有服务
start:
	@echo "$(GREEN)启动所有服务...$(RESET)"
	docker-compose up -d
	@echo "$(GREEN)等待服务启动...$(RESET)"
	@sleep 10
	@echo "$(GREEN)服务启动完成！$(RESET)"
	@echo ""
	@echo "$(CYAN)服务访问地址：$(RESET)"
	@echo "  API网关:      http://localhost:8080"
	@echo "  用户服务:     http://localhost:8001"
	@echo "  Consul UI:    http://localhost:8500"
	@echo "  Grafana:      http://localhost:3000"

## stop: 停止所有服务
stop:
	@echo "$(YELLOW)停止所有服务...$(RESET)"
	docker-compose down

## restart: 重启所有服务
restart: stop start

## logs: 查看所有服务日志
logs:
	docker-compose logs -f

## logs-user: 查看用户服务日志
logs-user:
	docker-compose logs -f user-service

## logs-gateway: 查看网关日志
logs-gateway:
	docker-compose logs -f gateway

## status: 查看服务状态
status:
	@echo "$(BLUE)服务状态:$(RESET)"
	docker-compose ps

## test: 运行API测试
test:
	@echo "$(PURPLE)运行API测试...$(RESET)"
	@if command -v bash > /dev/null; then \
		bash scripts/test-api.sh; \
	else \
		echo "$(RED)错误: 需要bash环境来运行测试$(RESET)"; \
	fi

## health: 检查服务健康状态
health:
	@echo "$(BLUE)检查服务健康状态...$(RESET)"
	@echo -n "API网关: "
	@if curl -s http://localhost:8080/health > /dev/null 2>&1; then \
		echo "$(GREEN)✅ 正常$(RESET)"; \
	else \
		echo "$(RED)❌ 异常$(RESET)"; \
	fi
	@echo -n "用户服务: "
	@if curl -s http://localhost:8001/health > /dev/null 2>&1; then \
		echo "$(GREEN)✅ 正常$(RESET)"; \
	else \
		echo "$(RED)❌ 异常$(RESET)"; \
	fi

## dev: 启动开发环境（仅基础设施）
dev:
	@echo "$(BLUE)启动开发环境...$(RESET)"
	docker-compose up -d mongodb redis consul nats prometheus grafana
	@echo "$(GREEN)开发环境启动完成！$(RESET)"
	@echo ""
	@echo "$(CYAN)现在可以本地运行服务:$(RESET)"
	@echo "  cd cmd/user-service && go run main.go"

## clean: 清理所有容器和数据卷
clean:
	@echo "$(YELLOW)清理所有容器和数据卷...$(RESET)"
	docker-compose down -v
	docker system prune -f

## reset: 重置系统（清理并重新启动）
reset: clean build start

## init: 初始化数据库数据
init:
	@echo "$(BLUE)初始化数据库数据...$(RESET)"
	docker-compose exec mongodb mongosh mule_cloud /docker-entrypoint-initdb.d/init-mongo.js

## backup: 备份数据库
backup:
	@echo "$(BLUE)备份数据库...$(RESET)"
	mkdir -p backups
	docker-compose exec mongodb mongodump --out /tmp/backup
	docker cp $$(docker-compose ps -q mongodb):/tmp/backup ./backups/backup-$$(date +%Y%m%d-%H%M%S)

## restore: 恢复数据库（需要指定备份目录）
restore:
	@if [ -z "$(BACKUP_DIR)" ]; then \
		echo "$(RED)错误: 请指定备份目录，例如: make restore BACKUP_DIR=backups/backup-20231201-120000$(RESET)"; \
		exit 1; \
	fi
	@echo "$(BLUE)恢复数据库...$(RESET)"
	docker cp $(BACKUP_DIR) $$(docker-compose ps -q mongodb):/tmp/restore
	docker-compose exec mongodb mongorestore /tmp/restore

## build-user: 只构建用户服务
build-user:
	@echo "$(BLUE)构建用户服务...$(RESET)"
	docker-compose build user-service

## build-gateway: 只构建网关服务
build-gateway:
	@echo "$(BLUE)构建网关服务...$(RESET)"
	docker-compose build gateway

## fmt: 格式化Go代码
fmt:
	@echo "$(BLUE)格式化Go代码...$(RESET)"
	go fmt ./...

## lint: 运行代码检查
lint:
	@echo "$(BLUE)运行代码检查...$(RESET)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)警告: golangci-lint未安装，跳过代码检查$(RESET)"; \
	fi

## deps: 更新Go依赖
deps:
	@echo "$(BLUE)更新Go依赖...$(RESET)"
	go mod tidy
	go mod download

## quick-start: 快速启动（构建并启动）
quick-start: build start health

## prod: 生产环境部署
prod:
	@echo "$(RED)警告: 这是生产环境部署命令$(RESET)"
	@echo "$(YELLOW)请确保已修改默认密码和配置$(RESET)"
	@read -p "确认继续？(y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	docker-compose -f docker-compose.prod.yaml up -d
