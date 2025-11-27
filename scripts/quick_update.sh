#!/bin/bash

# 快速更新脚本（不备份数据库，适合小改动）

set -e

echo "开始更新 Paper AI..."

# 进入项目目录
cd /opt/paper_ai

# 拉取最新代码
echo "[1/3] 拉取最新代码..."
git pull

# 重新构建并启动
echo "[2/3] 重新构建并启动..."
docker-compose up -d --build

# 等待服务启动
echo "[3/3] 等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "服务状态："
docker-compose ps

echo ""
echo "✅ 更新完成！"
echo ""
echo "查看日志："
echo "  docker-compose logs -f app"
