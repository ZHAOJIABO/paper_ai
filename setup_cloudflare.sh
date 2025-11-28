#!/bin/bash

# Paper AI - Cloudflare 域名配置脚本
# 专门为使用 Cloudflare CDN 优化的配置流程

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo "======================================"
echo "   Paper AI - Cloudflare 配置脚本"
echo "======================================"
echo ""

# 检查是否为 root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}请使用 sudo 运行此脚本${NC}"
    exit 1
fi

# 获取域名
echo -e "${BLUE}请输入你的域名（例如：example.com 或 api.example.com）：${NC}"
read -r DOMAIN

if [ -z "$DOMAIN" ]; then
    echo -e "${RED}域名不能为空${NC}"
    exit 1
fi

echo ""
echo -e "${CYAN}==================== 重要提示 ====================${NC}"
echo -e "${YELLOW}使用 Cloudflare CDN 时，请确保：${NC}"
echo ""
echo "1. ✅ 已在 Cloudflare 配置 DNS A 记录"
echo "2. ✅ DNS 记录的 Proxy status 为【橙色云朵】(Proxied)"
echo "3. ✅ SSL/TLS 加密模式设置为【Full (strict)】"
echo "4. ⏰ DNS 可能需要 2-10 分钟生效"
echo ""
echo -e "${CYAN}=================================================${NC}"
echo ""
echo -e "${BLUE}是否已完成上述 Cloudflare 配置？(y/n)${NC}"
read -r CLOUDFLARE_READY

if [ "$CLOUDFLARE_READY" != "y" ] && [ "$CLOUDFLARE_READY" != "Y" ]; then
    echo ""
    echo -e "${YELLOW}请先完成 Cloudflare 配置，参考：CLOUDFLARE_SETUP.md${NC}"
    echo ""
    echo "快速步骤："
    echo "1. 登录 Cloudflare 控制台"
    echo "2. 进入你的域名 > DNS"
    echo "3. 添加 A 记录："
    echo "   - Name: @"
    echo "   - Content: $(curl -s ifconfig.me)"
    echo "   - Proxy status: 🟠 Proxied"
    echo "4. SSL/TLS > 选择 Full (strict) 模式"
    echo ""
    exit 0
fi

# 询问是否需要 www 子域名
echo ""
echo -e "${BLUE}是否同时配置 www.$DOMAIN？(y/n)${NC}"
read -r ADD_WWW

# 显示配置摘要
echo ""
echo -e "${GREEN}==================== 配置摘要 ====================${NC}"
echo "域名: $DOMAIN"
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    echo "WWW: www.$DOMAIN"
fi
echo "服务器 IP: $(curl -s ifconfig.me)"
echo "SSL: Let's Encrypt"
echo "CDN: Cloudflare"
echo -e "${GREEN}================================================${NC}"
echo ""
echo -e "${BLUE}确认开始配置？(y/n)${NC}"
read -r CONFIRM

if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "配置已取消"
    exit 0
fi

# 1. 检查 Certbot
echo ""
echo -e "${GREEN}[1/7] 检查 Certbot...${NC}"
if ! command -v certbot &> /dev/null; then
    echo -e "${YELLOW}Certbot 未安装，正在安装...${NC}"
    apt-get update -qq
    apt-get install -y certbot
else
    echo "✓ Certbot 已安装: $(certbot --version | head -n1)"
fi

# 2. 创建必要目录
echo -e "${GREEN}[2/7] 创建目录...${NC}"
mkdir -p ssl
mkdir -p logs
echo "✓ 目录创建完成"

# 3. 备份现有配置
echo -e "${GREEN}[3/7] 备份配置文件...${NC}"
if [ -f "nginx.conf" ]; then
    cp nginx.conf "nginx.conf.backup.$(date +%Y%m%d_%H%M%S)"
    echo "✓ 已备份 nginx.conf"
fi

# 4. 使用 Cloudflare 优化的配置
echo -e "${GREEN}[4/7] 配置 Nginx（Cloudflare 优化）...${NC}"
if [ -f "nginx.cloudflare.conf" ]; then
    cp nginx.cloudflare.conf nginx.conf
    echo "✓ 已应用 Cloudflare 优化配置"
else
    echo -e "${YELLOW}未找到 nginx.cloudflare.conf，使用默认配置${NC}"
fi

# 5. 停止 nginx 容器
echo -e "${GREEN}[5/7] 停止 Nginx 容器...${NC}"
docker-compose stop nginx 2>/dev/null || true
sleep 2
echo "✓ Nginx 已停止"

# 6. 申请 SSL 证书
echo -e "${GREEN}[6/7] 申请 SSL 证书...${NC}"
echo "正在向 Let's Encrypt 申请证书，请稍候..."

if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    certbot certonly --standalone \
        -d "$DOMAIN" \
        -d "www.$DOMAIN" \
        --agree-tos \
        --non-interactive \
        --email "admin@$DOMAIN" \
        --force-renewal 2>&1 | tee /tmp/certbot.log || {
            echo ""
            echo -e "${RED}❌ 证书申请失败${NC}"
            echo ""
            echo -e "${YELLOW}可能的原因：${NC}"
            echo "1. DNS 记录还未生效（请等待 5-10 分钟后重试）"
            echo "2. Cloudflare Proxy 开启导致 Let's Encrypt 无法验证"
            echo ""
            echo -e "${CYAN}解决方案（推荐）：${NC}"
            echo "临时关闭 Cloudflare 代理来申请证书："
            echo "1. 进入 Cloudflare > DNS 设置"
            echo "2. 点击 $DOMAIN 记录的橙色云朵，变成灰色（DNS only）"
            if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
                echo "3. 同样点击 www.$DOMAIN 的橙色云朵，变成灰色"
            fi
            echo "4. 等待 2-3 分钟"
            echo "5. 重新运行此脚本"
            echo "6. 证书申请成功后，重新开启橙色云朵（Proxied）"
            echo ""
            echo "查看详细错误："
            echo "  cat /tmp/certbot.log"
            echo ""
            exit 1
        }
else
    certbot certonly --standalone \
        -d "$DOMAIN" \
        --agree-tos \
        --non-interactive \
        --email "admin@$DOMAIN" \
        --force-renewal 2>&1 | tee /tmp/certbot.log || {
            echo ""
            echo -e "${RED}❌ 证书申请失败${NC}"
            echo ""
            echo -e "${YELLOW}可能的原因：${NC}"
            echo "1. DNS 记录还未生效（请等待 5-10 分钟后重试）"
            echo "2. Cloudflare Proxy 开启导致 Let's Encrypt 无法验证"
            echo ""
            echo -e "${CYAN}解决方案（推荐）：${NC}"
            echo "临时关闭 Cloudflare 代理来申请证书："
            echo "1. 进入 Cloudflare > DNS 设置"
            echo "2. 点击 $DOMAIN 记录的橙色云朵，变成灰色（DNS only）"
            echo "3. 等待 2-3 分钟"
            echo "4. 重新运行此脚本"
            echo "5. 证书申请成功后，重新开启橙色云朵（Proxied）"
            echo ""
            echo "查看详细错误："
            echo "  cat /tmp/certbot.log"
            echo ""
            exit 1
        }
fi

echo "✓ SSL 证书申请成功"

# 复制证书到项目目录
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ssl/
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ssl/
chmod 644 ssl/fullchain.pem
chmod 644 ssl/privkey.pem
echo "✓ 证书已复制到 ssl/ 目录"

# 7. 更新 Nginx 配置
echo -e "${GREEN}[7/7] 更新 Nginx 配置...${NC}"

# 构建 server_name
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    SERVER_NAME="$DOMAIN www.$DOMAIN"
else
    SERVER_NAME="$DOMAIN"
fi

# 更新配置文件中的域名
sed -i "s/server_name _;/server_name $SERVER_NAME;/g" nginx.conf

# 取消注释 HTTPS 配置
sed -i 's/# server {/server {/g' nginx.conf
sed -i 's/#     listen/    listen/g' nginx.conf
sed -i 's/#     server_name/    server_name/g' nginx.conf
sed -i 's/#     access_log/    access_log/g' nginx.conf
sed -i 's/#     ssl_/    ssl_/g' nginx.conf
sed -i 's/#     add_header/    add_header/g' nginx.conf
sed -i 's/#     client_max_body_size/    client_max_body_size/g' nginx.conf
sed -i 's/#     location/    location/g' nginx.conf
sed -i 's/#         proxy/        proxy/g' nginx.conf
sed -i 's/#         access_log/        access_log/g' nginx.conf
sed -i 's/# }/}/g' nginx.conf

# 再次替换 HTTPS 部分的域名
sed -i "s/server_name your-domain.com www.your-domain.com;/server_name $SERVER_NAME;/g" nginx.conf

echo "✓ Nginx 配置已更新"

# 重启所有服务
echo ""
echo -e "${GREEN}重启服务...${NC}"
docker-compose up -d

# 等待服务启动
echo "等待服务启动..."
sleep 8

# 检查服务状态
echo ""
echo "======================================"
echo "服务状态："
docker-compose ps

echo ""
echo "======================================"
echo -e "${GREEN}✨ 配置完成！${NC}"
echo ""
echo -e "${CYAN}访问地址：${NC}"
echo -e "  ${GREEN}✓${NC} https://$DOMAIN"
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    echo -e "  ${GREEN}✓${NC} https://www.$DOMAIN"
fi
echo ""
echo -e "${CYAN}重要提示：${NC}"
echo -e "${YELLOW}1. 如果申请证书时临时关闭了 Cloudflare 代理（灰色云朵）${NC}"
echo -e "${YELLOW}   现在可以重新开启橙色云朵（Proxied）以启用 CDN${NC}"
echo ""
echo -e "2. 在 Cloudflare 控制台验证配置："
echo "   - SSL/TLS 模式：Full (strict) ✓"
echo "   - DNS 代理：🟠 Proxied ✓"
echo "   - Always Use HTTPS：ON ✓"
echo ""
echo -e "${CYAN}测试命令：${NC}"
echo "  # 测试 HTTPS"
echo "  curl -I https://$DOMAIN/health"
echo ""
echo "  # 查看日志"
echo "  docker-compose logs -f nginx"
echo ""
echo -e "${CYAN}证书信息：${NC}"
echo "  位置: /etc/letsencrypt/live/$DOMAIN/"
echo "  有效期: 90 天"
echo "  自动续期: 已配置"
echo ""
echo "======================================"

# 配置证书自动续期
echo ""
echo -e "${GREEN}配置证书自动续期...${NC}"
PROJECT_DIR=$(pwd)
CRON_JOB="0 0 1 * * certbot renew --quiet --post-hook 'cd $PROJECT_DIR && cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ssl/ && cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ssl/ && docker-compose restart nginx' 2>&1 | logger -t certbot-renew"

if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo "✓ 已添加证书自动续期任务（每月1号执行）"
else
    echo "✓ 证书自动续期任务已存在"
fi

# 验证配置
echo ""
echo -e "${CYAN}正在验证配置...${NC}"
sleep 3

if curl -s -o /dev/null -w "%{http_code}" https://$DOMAIN/health | grep -q "200"; then
    echo -e "${GREEN}✓ 健康检查通过！服务运行正常${NC}"
else
    echo -e "${YELLOW}⚠ 无法访问健康检查接口，请检查：${NC}"
    echo "  1. 服务是否正常运行：docker-compose ps"
    echo "  2. 查看日志：docker-compose logs"
    echo "  3. DNS 是否生效：nslookup $DOMAIN"
fi

echo ""
echo -e "${BLUE}配置完成！享受你的 Paper AI 服务吧！🎉${NC}"
echo ""
