.PHONY: help build run clean test deps

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
