#!/bin/bash

# 服务更新脚本

set -e

echo "======================================"
echo "   Paper AI 更新脚本"
echo "======================================"

# 1. 备份数据库
echo "[1/5] 备份数据库..."
./scripts/backup.sh

# 2. 拉取最新代码
echo "[2/5] 拉取最新代码..."
git pull

# 3. 停止服务
echo "[3/5] 停止服务..."
docker-compose down

# 4. 重新构建并启动
echo "[4/5] 重新构建并启动服务..."
docker-compose up -d --build

# 5. 等待服务启动
echo "[5/5] 等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "======================================"
echo "服务状态："
docker-compose ps

echo ""
echo "查看日志："
docker-compose logs -f app

echo ""
echo "======================================"
echo "更新完成！"
echo "======================================"
