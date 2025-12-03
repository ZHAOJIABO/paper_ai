.PHONY: help build run clean test deps migrate-up migrate-down migrate-create migrate-force migrate-version migrate-drop

help: ## 显示帮助信息
	@echo "Paper AI - 科研AI服务平台"
	@echo ""
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go mod tidy

build: ## 编译项目
	go build -o paper_ai cmd/server/main.go
	@echo "编译完成: ./paper_ai"

run: ## 运行服务
	go run cmd/server/main.go

clean: ## 清理编译产物
	rm -f paper_ai
	@echo "清理完成"

test: ## 运行测试脚本（需要服务运行中）
	./test.sh

dev: deps build ## 开发环境准备
	@echo "开发环境准备完成"
	@echo "请先配置 config/config.yaml 中的 API Key"
	@echo "然后运行: make run"

# 数据库迁移相关命令
# 需要设置 DATABASE_URL 环境变量，格式：postgres://user:password@host:port/dbname?sslmode=disable
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/paper_ai?sslmode=disable
MIGRATIONS_DIR = migrations

migrate-up: ## 执行所有待执行的迁移
	@echo "执行数据库迁移..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
	@echo "迁移完成"

migrate-down: ## 回滚最后一次迁移
	@echo "回滚数据库迁移..."
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1
	@echo "回滚完成"

migrate-drop: ## 删除所有表（危险操作！）
	@echo "警告: 这将删除所有表！"
	@read -p "确认删除所有表? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" drop -f; \
		echo "所有表已删除"; \
	else \
		echo "操作已取消"; \
	fi

migrate-create: ## 创建新的迁移文件 (使用: make migrate-create name=add_users_table)
	@if [ -z "$(name)" ]; then \
		echo "错误: 请指定迁移文件名称"; \
		echo "用法: make migrate-create name=add_users_table"; \
		exit 1; \
	fi
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	@echo "迁移文件已创建在 $(MIGRATIONS_DIR)/"

migrate-force: ## 强制设置迁移版本 (使用: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "错误: 请指定版本号"; \
		echo "用法: make migrate-force version=1"; \
		exit 1; \
	fi
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)
	@echo "已强制设置为版本 $(version)"

migrate-version: ## 查看当前迁移版本
	@go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

