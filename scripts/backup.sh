#!/bin/bash

# 数据库备份脚本

set -e

BACKUP_DIR="./backup"
DATE=$(date +%Y%m%d_%H%M%S)
CONTAINER_NAME="paper_ai_db"

echo "开始备份数据库..."

# 创建备份目录
mkdir -p $BACKUP_DIR

# 执行备份
docker exec $CONTAINER_NAME pg_dump -U paperai paper_ai | gzip > $BACKUP_DIR/paper_ai_$DATE.sql.gz

# 检查备份文件大小
SIZE=$(du -h $BACKUP_DIR/paper_ai_$DATE.sql.gz | cut -f1)
echo "备份完成: paper_ai_$DATE.sql.gz ($SIZE)"

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
echo "已清理 7 天前的旧备份"

# 显示当前备份列表
echo ""
echo "当前备份列表："
ls -lh $BACKUP_DIR/*.sql.gz 2>/dev/null || echo "暂无备份文件"
