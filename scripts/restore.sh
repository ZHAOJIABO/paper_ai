#!/bin/bash

# 数据库恢复脚本

set -e

if [ -z "$1" ]; then
    echo "用法: $0 <备份文件>"
    echo ""
    echo "可用备份文件："
    ls -1 ./backup/*.sql.gz 2>/dev/null || echo "暂无备份文件"
    exit 1
fi

BACKUP_FILE=$1
CONTAINER_NAME="paper_ai_db"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "错误: 文件不存在 - $BACKUP_FILE"
    exit 1
fi

echo "警告: 此操作将覆盖当前数据库！"
read -p "确认恢复数据库? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "操作已取消"
    exit 0
fi

echo "开始恢复数据库..."

# 先备份当前数据库
echo "正在备份当前数据库..."
./scripts/backup.sh

# 恢复数据库
echo "正在恢复数据: $BACKUP_FILE"
gunzip -c $BACKUP_FILE | docker exec -i $CONTAINER_NAME psql -U paperai paper_ai

echo "数据库恢复完成！"
