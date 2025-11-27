# Paper AI 部署指南

## 目录
- [服务器选择与网络说明](#服务器选择与网络说明)
- [准备工作](#准备工作)
- [方式一：Docker 部署（推荐）](#方式一docker-部署推荐)
- [方式二：systemd 服务部署](#方式二systemd-服务部署)
- [国内访问优化方案](#国内访问优化方案)
- [监控与维护](#监控与维护)

---

## 服务器选择与网络说明

### Vultr 日本服务器
- **优点**：延迟低（ping 约 50-100ms），价格合理
- **推荐配置**：1核2GB起步，根据并发需求调整

### 国内访问情况说明

#### 🟢 可以访问的情况
1. **普通 HTTP/HTTPS 服务**：Vultr 日本服务器可以被国内正常访问
2. **API 服务**：你的 Paper AI 后端 API 完全可以从国内访问
3. **访问速度**：日本节点对国内友好，延迟一般在 50-150ms

#### 🔴 需要注意的问题
1. **IP 被墙风险**：
   - 小概率事件，但如果同 IP 段有违规内容可能被连带
   - **解决方案**：定期检测 IP 可达性，必要时更换 IP（Vultr 支持删除重建）

2. **Claude API 访问**：
   - Vultr 日本服务器**可以**正常访问 Claude API（api.anthropic.com）
   - 无需额外配置代理

3. **部分地区网络限制**：
   - 个别地区（如学校、企业内网）可能有限制
   - **解决方案**：使用 CDN 加速（如 Cloudflare）

4. **DNS 污染**：
   - 如果使用域名，建议使用国外 DNS 解析服务
   - **推荐**：Cloudflare DNS、Google DNS

---

## 准备工作

### 1. 服务器初始化

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git vim

# 配置时区
sudo timedatectl set-timezone Asia/Shanghai

# 配置防火墙
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 8080/tcp  # API（可选，建议用反向代理）
sudo ufw enable
```

### 2. 域名配置（推荐）

建议购买域名并配置：
```
api.yourdomain.com  ->  Vultr 服务器 IP
```

---

## 方式一：Docker 部署（推荐）

### 1. 安装 Docker

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 启动 Docker
sudo systemctl enable docker
sudo systemctl start docker

# 安装 Docker Compose
sudo apt install -y docker-compose

# 验证安装
docker --version
docker-compose --version
```

### 2. 创建部署文件

在项目根目录创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# 安装 git（某些 Go 依赖需要）
RUN apk add --no-cache git

# 复制 go mod 文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o paper_ai ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/paper_ai .
COPY --from=builder /app/config ./config

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./paper_ai"]
```

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    container_name: paper_ai_db
    environment:
      POSTGRES_DB: paper_ai
      POSTGRES_USER: paperai
      POSTGRES_PASSWORD: ${DB_PASSWORD:-change_me_in_production}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U paperai"]
      interval: 10s
      timeout: 5s
      retries: 5

  app:
    build: .
    container_name: paper_ai_app
    environment:
      CONFIG_PATH: /root/config/config.yaml
    volumes:
      - ./config:/root/config
      - ./logs:/root/logs
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    container_name: paper_ai_nginx
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro  # SSL 证书目录（可选）
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
```

创建 `nginx.conf`：

```nginx
events {
    worker_connections 1024;
}

http {
    upstream paper_ai_backend {
        server app:8080;
    }

    # HTTP 服务器
    server {
        listen 80;
        server_name _;  # 替换为你的域名

        # 请求体大小限制
        client_max_body_size 10M;

        location / {
            proxy_pass http://paper_ai_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

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

    # HTTPS 配置（取消注释并配置 SSL 证书）
    # server {
    #     listen 443 ssl http2;
    #     server_name your-domain.com;
    #
    #     ssl_certificate /etc/nginx/ssl/fullchain.pem;
    #     ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    #     ssl_protocols TLSv1.2 TLSv1.3;
    #     ssl_ciphers HIGH:!aNULL:!MD5;
    #
    #     location / {
    #         proxy_pass http://paper_ai_backend;
    #         proxy_set_header Host $host;
    #         proxy_set_header X-Real-IP $remote_addr;
    #         proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    #         proxy_set_header X-Forwarded-Proto $scheme;
    #     }
    # }
}
```

创建 `.env` 文件：

```bash
# 数据库密码
DB_PASSWORD=your_secure_password_here
```

### 3. 配置生产环境

复制配置文件并修改：

```bash
cp config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml`：

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 120s  # 增加超时时间

ai:
  default_provider: claude
  providers:
    claude:
      api_key: "sk-ant-你的实际API密钥"
      base_url: "https://api.anthropic.com"
      model: "claude-3-5-sonnet-20241022"
      timeout: 120s

database:
  type: postgres
  host: postgres  # Docker Compose 服务名
  port: 5432
  user: paperai
  password: "your_secure_password_here"  # 与 .env 中一致
  dbname: paper_ai
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
  auto_migrate: true
  log_mode: info

jwt:
  secret_key: "生产环境请使用强随机字符串"  # 使用 openssl rand -base64 32 生成
  access_token_expiry: 7200      # 2小时
  refresh_token_expiry: 604800   # 7天

idgen:
  worker_id: 1  # 如果多实例部署，每个实例设置不同ID
```

### 4. 部署

```bash
# 构建并启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 查看状态
docker-compose ps

# 重启服务
docker-compose restart app

# 停止服务
docker-compose down
```

### 5. 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build

# 查看日志确认启动成功
docker-compose logs -f app
```

---

## 方式二：systemd 服务部署

### 1. 安装 PostgreSQL

```bash
# 安装 PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# 启动服务
sudo systemctl enable postgresql
sudo systemctl start postgresql

# 创建数据库和用户
sudo -u postgres psql <<EOF
CREATE DATABASE paper_ai;
CREATE USER paperai WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE paper_ai TO paperai;
\q
EOF
```

### 2. 安装 Go（如果需要在服务器编译）

```bash
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### 3. 部署应用

```bash
# 创建应用目录
sudo mkdir -p /opt/paper_ai
sudo chown $USER:$USER /opt/paper_ai

# 克隆代码（或上传编译好的二进制）
cd /opt/paper_ai
git clone <your-repo-url> .

# 编译
go build -o paper_ai ./cmd/server

# 配置文件
cp config/config.example.yaml config/config.yaml
vim config/config.yaml  # 修改配置
```

### 4. 创建 systemd 服务

创建 `/etc/systemd/system/paper_ai.service`：

```ini
[Unit]
Description=Paper AI Service
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/paper_ai
ExecStart=/opt/paper_ai/paper_ai
Environment="CONFIG_PATH=/opt/paper_ai/config/config.yaml"
Restart=always
RestartSec=10

# 安全配置
NoNewPrivileges=true
PrivateTmp=true

# 日志
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl enable paper_ai
sudo systemctl start paper_ai

# 查看状态
sudo systemctl status paper_ai

# 查看日志
sudo journalctl -u paper_ai -f
```

### 5. 安装 Nginx 反向代理

```bash
sudo apt install -y nginx

# 创建配置
sudo vim /etc/nginx/sites-available/paper_ai
```

配置内容：

```nginx
server {
    listen 80;
    server_name your-domain.com;  # 替换为你的域名

    client_max_body_size 10M;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_connect_timeout 60s;
        proxy_send_timeout 120s;
        proxy_read_timeout 120s;
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/paper_ai /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 国内访问优化方案

### 1. 使用 Cloudflare CDN

**步骤**：
1. 注册 Cloudflare 账号
2. 添加你的域名
3. 修改域名 NS 记录到 Cloudflare
4. 开启代理（橙色云朵图标）

**优势**：
- 免费 CDN 加速
- 隐藏真实 IP（防止被墙）
- 自动 HTTPS
- DDoS 防护

### 2. SSL 证书配置

使用 Let's Encrypt 免费证书：

```bash
# 安装 certbot
sudo apt install -y certbot python3-certbot-nginx

# 自动配置 HTTPS
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo systemctl enable certbot.timer
```

### 3. IP 可用性监控

创建监控脚本 `/opt/scripts/check_ip.sh`：

```bash
#!/bin/bash

# 检测 IP 是否可从国内访问
SERVER_IP="your_server_ip"
WEBHOOK="your_notification_webhook"  # 如钉钉/企业微信

if ! ping -c 3 $SERVER_IP > /dev/null 2>&1; then
    curl -X POST $WEBHOOK -d "{\"msg\": \"服务器 IP 可能被墙，请检查！\"}"
fi
```

设置定时任务：

```bash
crontab -e
# 每小时检测一次
0 * * * * /opt/scripts/check_ip.sh
```

### 4. 访问速度优化

**Nginx 配置优化**：

```nginx
http {
    # 启用 gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # 启用缓存
    proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m max_size=100m inactive=60m;

    server {
        location / {
            # 启用缓存（针对 GET 请求）
            proxy_cache api_cache;
            proxy_cache_valid 200 10m;
            proxy_cache_methods GET HEAD;
            proxy_cache_key "$scheme$request_method$host$request_uri";

            proxy_pass http://localhost:8080;
        }
    }
}
```

---

## 监控与维护

### 1. 日志管理

```bash
# 查看 Docker 日志
docker-compose logs -f --tail=100 app

# 查看 systemd 日志
sudo journalctl -u paper_ai -f --since "1 hour ago"

# 日志轮转配置
sudo vim /etc/logrotate.d/paper_ai
```

### 2. 性能监控

安装监控工具：

```bash
# 安装 htop
sudo apt install -y htop

# 安装 netdata（可选，Web 界面监控）
bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```

### 3. 数据库备份

创建备份脚本 `/opt/scripts/backup_db.sh`：

```bash
#!/bin/bash

BACKUP_DIR="/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR

# Docker 方式备份
docker exec paper_ai_db pg_dump -U paperai paper_ai | gzip > $BACKUP_DIR/paper_ai_$DATE.sql.gz

# systemd 方式备份
# sudo -u postgres pg_dump paper_ai | gzip > $BACKUP_DIR/paper_ai_$DATE.sql.gz

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: paper_ai_$DATE.sql.gz"
```

设置自动备份：

```bash
chmod +x /opt/scripts/backup_db.sh
crontab -e
# 每天凌晨 2 点备份
0 2 * * * /opt/scripts/backup_db.sh >> /var/log/db_backup.log 2>&1
```

### 4. 健康检查

在应用中添加健康检查端点（如果还没有）：

```bash
# 测试健康检查
curl http://localhost:8080/health
```

### 5. 更新流程

```bash
# 1. 备份数据库
/opt/scripts/backup_db.sh

# 2. 拉取最新代码
cd /opt/paper_ai && git pull

# 3. Docker 方式更新
docker-compose up -d --build

# 4. systemd 方式更新
go build -o paper_ai ./cmd/server
sudo systemctl restart paper_ai

# 5. 验证服务
curl http://localhost:8080/health
```

---

## 故障排查

### 常见问题

1. **无法连接数据库**
```bash
# 检查数据库状态
docker-compose ps postgres
sudo systemctl status postgresql

# 测试连接
psql -h localhost -U paperai -d paper_ai
```

2. **API 响应慢**
```bash
# 检查 Claude API 连接
curl -I https://api.anthropic.com

# 检查资源使用
docker stats
htop
```

3. **国内无法访问**
```bash
# 从国内服务器测试
ping your-server-ip
curl -I http://your-domain.com

# 检查防火墙
sudo ufw status
```

---

## 安全建议

1. **修改 SSH 端口**
2. **禁用 root 登录**
3. **配置 fail2ban 防暴力破解**
4. **定期更新系统和依赖**
5. **使用强密码和 SSH 密钥**
6. **配置 HTTPS（必须）**
7. **定期备份数据**
8. **监控服务器资源和日志**

---

## 总结

### 推荐方案
- **小项目/个人**：Docker + Docker Compose（简单快速）
- **生产环境**：Docker + Nginx + Cloudflare CDN + 自动备份
- **高可用需求**：Kubernetes + 负载均衡

### 国内访问结论
✅ **Vultr 日本服务器完全可以从国内访问**
✅ **服务器可以正常访问 Claude API**
⚠️ **建议使用 Cloudflare CDN 增加稳定性**
⚠️ **定期监控 IP 可用性**

如有疑问，请参考文档或提交 Issue。
