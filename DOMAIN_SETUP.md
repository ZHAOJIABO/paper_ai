# 域名配置指南

本文档介绍如何为 Paper AI 配置域名访问和 HTTPS。

## 快速开始

### 方法一：自动配置（推荐）

```bash
# 在服务器上运行
sudo chmod +x setup_domain.sh
sudo ./setup_domain.sh
```

脚本会自动完成：
- 安装 Certbot
- 申请 SSL 证书
- 配置 Nginx
- 设置证书自动续期

### 方法二：手动配置

如果你想手动配置或自动脚本遇到问题，请按以下步骤操作。

## 手动配置步骤

### 1. DNS 解析配置

在你的域名服务商（阿里云/腾讯云/Cloudflare 等）添加 DNS 记录：

| 类型 | 主机记录 | 记录值 | TTL |
|------|---------|--------|-----|
| A    | @       | 你的服务器IP | 600 |
| A    | www     | 你的服务器IP | 600 |

**验证 DNS 是否生效：**
```bash
# 查询域名解析
nslookup your-domain.com

# 或使用 dig
dig your-domain.com +short
```

### 2. 安装 Certbot

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y certbot

# CentOS/RHEL
sudo yum install -y certbot
```

### 3. 停止 Nginx 容器

```bash
cd /path/to/paper_ai
docker-compose stop nginx
```

### 4. 申请 SSL 证书

**单个域名：**
```bash
sudo certbot certonly --standalone \
  -d your-domain.com \
  --agree-tos \
  --non-interactive \
  --email your-email@example.com
```

**包含 www 子域名：**
```bash
sudo certbot certonly --standalone \
  -d your-domain.com \
  -d www.your-domain.com \
  --agree-tos \
  --non-interactive \
  --email your-email@example.com
```

如果失败，请检查：
1. DNS 解析是否生效（等待 5-10 分钟）
2. 防火墙是否开放 80 和 443 端口
3. 80 端口是否被其他程序占用

### 5. 复制证书到项目目录

```bash
# 创建 ssl 目录
mkdir -p ssl

# 复制证书（替换 your-domain.com）
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ssl/

# 修改权限
sudo chmod 644 ssl/fullchain.pem
sudo chmod 644 ssl/privkey.pem
```

### 6. 更新 Nginx 配置

编辑 `nginx.conf` 文件，将第 19 行的 `server_name _;` 改为你的域名：

```nginx
server_name your-domain.com www.your-domain.com;
```

然后取消注释 HTTPS 配置部分（第 45-67 行），并修改域名：

```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    # ... 其他配置
}
```

**完整示例配置：**

```nginx
events {
    worker_connections 1024;
}

http {
    upstream paper_ai_backend {
        server app:8080;
    }

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml;

    # HTTP 重定向到 HTTPS
    server {
        listen 80;
        server_name your-domain.com www.your-domain.com;

        location / {
            return 301 https://$host$request_uri;
        }
    }

    # HTTPS 服务器
    server {
        listen 443 ssl http2;
        server_name your-domain.com www.your-domain.com;

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers HIGH:!aNULL:!MD5;

        client_max_body_size 10M;

        location / {
            proxy_pass http://paper_ai_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            proxy_connect_timeout 60s;
            proxy_send_timeout 120s;
            proxy_read_timeout 120s;
        }

        location /health {
            proxy_pass http://paper_ai_backend/health;
            access_log off;
        }
    }
}
```

### 7. 重启服务

```bash
docker-compose up -d
```

### 8. 验证配置

```bash
# 检查服务状态
docker-compose ps

# 测试 HTTP
curl -I http://your-domain.com

# 测试 HTTPS
curl -I https://your-domain.com

# 查看 Nginx 日志
docker-compose logs nginx
```

### 9. 配置证书自动续期

Let's Encrypt 证书有效期为 90 天，需要定期续期。

```bash
# 添加 cron 任务
sudo crontab -e

# 添加以下内容（替换路径）
0 0 1 * * certbot renew --quiet --post-hook 'cd /path/to/paper_ai && cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ssl/ && cp /etc/letsencrypt/live/your-domain.com/privkey.pem ssl/ && docker-compose restart nginx'
```

## 常见问题

### Q1: DNS 解析不生效

**解决方案：**
- 等待 5-30 分钟（DNS 传播需要时间）
- 检查 DNS 配置是否正确
- 使用 `nslookup` 或 `dig` 验证

### Q2: SSL 证书申请失败

**错误：Connection refused**
- 检查 80 端口是否开放：`sudo netstat -tlnp | grep :80`
- 确保 nginx 容器已停止：`docker-compose stop nginx`
- 检查防火墙：`sudo ufw status`

**错误：Timeout**
- DNS 可能未生效，等待后重试
- 检查服务器防火墙规则

### Q3: HTTPS 访问显示证书错误

**解决方案：**
- 检查证书文件是否正确复制到 `ssl/` 目录
- 验证 nginx 配置中的域名是否正确
- 查看 nginx 日志：`docker-compose logs nginx`

### Q4: HTTP 可以访问但 HTTPS 不行

**解决方案：**
- 检查 443 端口是否开放
- 验证证书文件权限：`ls -l ssl/`
- 检查 docker-compose 是否正确映射了 ssl 目录

### Q5: 如何强制使用 HTTPS

在 HTTP server 块中添加重定向：
```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$host$request_uri;
}
```

## 测试清单

配置完成后，请验证以下项目：

- [ ] DNS 解析正确指向服务器 IP
- [ ] HTTP 访问正常（http://your-domain.com）
- [ ] HTTPS 访问正常（https://your-domain.com）
- [ ] HTTP 自动重定向到 HTTPS
- [ ] SSL 证书有效（浏览器无警告）
- [ ] API 接口可以正常调用
- [ ] 证书自动续期已配置

## 安全建议

1. **定期更新系统和软件包**
   ```bash
   sudo apt-get update && sudo apt-get upgrade -y
   ```

2. **配置防火墙**
   ```bash
   sudo ufw allow 22/tcp   # SSH
   sudo ufw allow 80/tcp   # HTTP
   sudo ufw allow 443/tcp  # HTTPS
   sudo ufw enable
   ```

3. **更改默认密码**
   - 修改 `.env` 中的数据库密码
   - 修改 `config/config.yaml` 中的 JWT secret

4. **启用日志监控**
   ```bash
   # 实时查看日志
   docker-compose logs -f

   # 查看特定服务日志
   docker-compose logs nginx
   docker-compose logs app
   ```

5. **设置备份**
   ```bash
   # 使用项目提供的备份脚本
   ./scripts/backup.sh
   ```

## 监控和维护

### 查看证书过期时间
```bash
sudo certbot certificates
```

### 手动续期证书
```bash
sudo certbot renew
cd /path/to/paper_ai
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ssl/
docker-compose restart nginx
```

### 查看服务状态
```bash
docker-compose ps
docker-compose logs --tail=100
```

## 获取帮助

如果遇到问题：
1. 查看日志：`docker-compose logs`
2. 检查服务状态：`docker-compose ps`
3. 验证配置：`docker-compose config`
4. 测试 nginx 配置：`docker exec paper_ai_nginx nginx -t`
