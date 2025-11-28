#!/bin/bash

# Paper AI 域名配置脚本
# 自动配置域名和 SSL 证书

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "======================================"
echo "   Paper AI 域名配置脚本"
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
echo -e "${GREEN}将为以下域名配置 SSL：${NC}"
echo "  域名: $DOMAIN"
echo ""

# 询问是否需要 www 子域名
echo -e "${BLUE}是否同时配置 www.$DOMAIN？(y/n)${NC}"
read -r ADD_WWW

# 1. 检查 Certbot
echo -e "${GREEN}[1/6] 检查 Certbot...${NC}"
if ! command -v certbot &> /dev/null; then
    echo -e "${YELLOW}Certbot 未安装，正在安装...${NC}"
    apt-get update
    apt-get install -y certbot
else
    echo "Certbot 已安装: $(certbot --version)"
fi

# 2. 创建 SSL 目录
echo -e "${GREEN}[2/6] 创建 SSL 目录...${NC}"
mkdir -p ssl

# 3. 停止 nginx 容器以便 certbot 使用 80 端口
echo -e "${GREEN}[3/6] 停止 nginx 容器...${NC}"
docker-compose stop nginx || true

# 4. 申请 SSL 证书
echo -e "${GREEN}[4/6] 申请 SSL 证书...${NC}"
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    certbot certonly --standalone \
        -d "$DOMAIN" \
        -d "www.$DOMAIN" \
        --agree-tos \
        --non-interactive \
        --email "admin@$DOMAIN" \
        --force-renewal || {
            echo -e "${RED}证书申请失败，请检查：${NC}"
            echo "1. 域名是否已正确解析到此服务器"
            echo "2. 防火墙是否允许 80 和 443 端口"
            echo "3. 服务器 IP: $(curl -s ifconfig.me)"
            exit 1
        }
else
    certbot certonly --standalone \
        -d "$DOMAIN" \
        --agree-tos \
        --non-interactive \
        --email "admin@$DOMAIN" \
        --force-renewal || {
            echo -e "${RED}证书申请失败，请检查：${NC}"
            echo "1. 域名是否已正确解析到此服务器"
            echo "2. 防火墙是否允许 80 和 443 端口"
            echo "3. 服务器 IP: $(curl -s ifconfig.me)"
            exit 1
        }
fi

# 5. 复制证书到项目目录
echo -e "${GREEN}[5/6] 配置证书...${NC}"
cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ssl/
cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ssl/
chmod 644 ssl/fullchain.pem
chmod 644 ssl/privkey.pem

# 6. 更新 nginx 配置
echo -e "${GREEN}[6/6] 更新 Nginx 配置...${NC}"

# 构建 server_name
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    SERVER_NAME="$DOMAIN www.$DOMAIN"
else
    SERVER_NAME="$DOMAIN"
fi

# 创建新的 nginx 配置
cat > nginx.conf << EOF
events {
    worker_connections 1024;
}

http {
    upstream paper_ai_backend {
        server app:8080;
    }

    # 启用 gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # HTTP 服务器 - 重定向到 HTTPS
    server {
        listen 80;
        server_name $SERVER_NAME;

        # Let's Encrypt 验证
        location /.well-known/acme-challenge/ {
            root /var/www/certbot;
        }

        # 重定向所有其他请求到 HTTPS
        location / {
            return 301 https://\$host\$request_uri;
        }
    }

    # HTTPS 服务器
    server {
        listen 443 ssl http2;
        server_name $SERVER_NAME;

        # SSL 证书配置
        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;
        ssl_prefer_server_ciphers on;

        # 安全头
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # 请求体大小限制
        client_max_body_size 10M;

        # 代理到后端
        location / {
            proxy_pass http://paper_ai_backend;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;

            # 超时配置
            proxy_connect_timeout 60s;
            proxy_send_timeout 120s;
            proxy_read_timeout 120s;
        }

        # 健康检查
        location /health {
            proxy_pass http://paper_ai_backend/health;
            access_log off;
        }
    }
}
EOF

echo -e "${GREEN}Nginx 配置已更新${NC}"

# 重启服务
echo ""
echo -e "${GREEN}重启服务...${NC}"
docker-compose up -d

# 等待服务启动
sleep 5

# 检查服务状态
echo ""
echo "======================================"
echo "服务状态："
docker-compose ps

echo ""
echo "======================================"
echo -e "${GREEN}域名配置完成！${NC}"
echo ""
echo "访问地址："
echo -e "  ${GREEN}✓${NC} https://$DOMAIN"
if [ "$ADD_WWW" = "y" ] || [ "$ADD_WWW" = "Y" ]; then
    echo -e "  ${GREEN}✓${NC} https://www.$DOMAIN"
fi
echo ""
echo "证书信息："
echo "  - 证书位置: /etc/letsencrypt/live/$DOMAIN/"
echo "  - 有效期: 90 天"
echo "  - 自动续期: 已配置 cron job"
echo ""
echo "======================================"

# 设置自动续期
echo -e "${GREEN}配置证书自动续期...${NC}"
CRON_JOB="0 0 1 * * certbot renew --quiet --post-hook 'cd $(pwd) && cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem ssl/ && cp /etc/letsencrypt/live/$DOMAIN/privkey.pem ssl/ && docker-compose restart nginx'"

# 检查 cron job 是否已存在
if ! crontab -l 2>/dev/null | grep -q "certbot renew"; then
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo -e "${GREEN}已添加证书自动续期任务${NC}"
else
    echo -e "${YELLOW}证书自动续期任务已存在${NC}"
fi

echo ""
echo -e "${BLUE}测试访问：${NC}"
echo "  curl -I https://$DOMAIN/health"
echo ""
