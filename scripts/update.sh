#!/bin/bash

# 服务更新脚本

set -e

echo "======================================"
echo "   Paper AI 更新脚本"
echo "======================================"

# 1. 备份数据库
echo "[1/6] 备份数据库..."
./scripts/backup.sh

# 2. 拉取最新代码
echo "[2/6] 拉取最新代码..."
git pull

# 3. 检查 migrations 目录
echo "[3/6] 检查迁移文件..."
if [ ! -d "migrations" ]; then
  echo "❌ 错误: migrations 目录不存在"
  exit 1
fi
MIGRATION_COUNT=$(ls -1 migrations/*.sql 2>/dev/null | wc -l)
echo "✓ 找到 $MIGRATION_COUNT 个迁移文件"

# 4. 停止服务
echo "[4/6] 停止服务..."
docker-compose down

# 5. 重新构建并启动
echo "[5/6] 重新构建并启动服务..."
echo "注意: 数据库迁移将在应用启动时自动执行"
docker-compose up -d --build

# 6. 等待服务启动并检查迁移
echo "[6/6] 等待服务启动..."
sleep 10

# 检查迁移日志
echo ""
echo "======================================"
echo "检查迁移状态："
docker-compose logs app | grep -i "migration" || echo "未找到迁移日志"

echo ""
echo "服务状态："
docker-compose ps

echo ""
echo "======================================"
echo "更新完成！"
echo "如需查看完整日志，运行: docker-compose logs -f app"
echo "======================================"
