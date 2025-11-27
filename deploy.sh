#!/bin/bash

# Paper AI 部署脚本
# 用于快速部署到服务器

set -e

echo "======================================"
echo "   Paper AI 部署脚本"
echo "======================================"
echo ""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查是否为 root
if [ "$EUID" -ne 0 ]; then
    echo -e "${YELLOW}建议使用 sudo 运行此脚本${NC}"
fi

# 1. 检查 Docker
echo -e "${GREEN}[1/8] 检查 Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}Docker 未安装，正在安装...${NC}"
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
else
    echo "Docker 已安装: $(docker --version)"
fi

# 2. 检查 Docker Compose
echo -e "${GREEN}[2/8] 检查 Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    echo -e "${YELLOW}Docker Compose 未安装，正在安装...${NC}"
    apt-get update
    apt-get install -y docker-compose
else
    echo "Docker Compose 已安装: $(docker-compose --version)"
fi

# 3. 创建必要的目录
echo -e "${GREEN}[3/8] 创建目录结构...${NC}"
mkdir -p logs ssl backup

# 4. 配置环境变量
echo -e "${GREEN}[4/8] 配置环境变量...${NC}"
if [ ! -f .env ]; then
    cp .env.example .env
    # 生成随机密码
    DB_PASSWORD=$(openssl rand -base64 32)
    sed -i "s/your_secure_password_here/$DB_PASSWORD/g" .env
    echo -e "${YELLOW}已生成随机数据库密码，请查看 .env 文件${NC}"
else
    echo ".env 文件已存在"
fi

# 5. 配置应用配置文件
echo -e "${GREEN}[5/8] 配置应用...${NC}"
if [ ! -f config/config.yaml ]; then
    cp config/config.example.yaml config/config.yaml
    echo -e "${YELLOW}请编辑 config/config.yaml 文件，配置：${NC}"
    echo "  - AI API Key"
    echo "  - JWT Secret"
    echo "  - 数据库密码（与 .env 一致）"
    read -p "按 Enter 继续编辑配置文件..."
    ${EDITOR:-vim} config/config.yaml
else
    echo "config.yaml 已存在"
fi

# 6. 配置防火墙
echo -e "${GREEN}[6/8] 配置防火墙...${NC}"
if command -v ufw &> /dev/null; then
    ufw allow 22/tcp
    ufw allow 80/tcp
    ufw allow 443/tcp
    echo "防火墙规则已添加"
else
    echo -e "${YELLOW}未检测到 ufw，请手动配置防火墙${NC}"
fi

# 7. 构建并启动服务
echo -e "${GREEN}[7/8] 构建并启动服务...${NC}"
docker-compose down
docker-compose up -d --build

# 8. 等待服务启动
echo -e "${GREEN}[8/8] 等待服务启动...${NC}"
sleep 10

# 检查服务状态
echo ""
echo "======================================"
echo "服务状态："
docker-compose ps

echo ""
echo "======================================"
echo -e "${GREEN}部署完成！${NC}"
echo ""
echo "服务地址："
echo "  - HTTP: http://$(curl -s ifconfig.me)"
echo "  - 本地: http://localhost:8080"
echo ""
echo "查看日志："
echo "  docker-compose logs -f"
echo ""
echo "重启服务："
echo "  docker-compose restart"
echo ""
echo "停止服务："
echo "  docker-compose down"
echo ""
echo -e "${YELLOW}重要提示：${NC}"
echo "1. 请配置域名解析到此服务器 IP"
echo "2. 建议配置 SSL 证书（使用 Let's Encrypt）"
echo "3. 定期备份数据库"
echo "======================================"
