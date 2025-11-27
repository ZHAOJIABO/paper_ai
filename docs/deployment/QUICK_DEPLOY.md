# Paper AI 快速部署指南

## 一键部署（推荐）

### 1. 购买 Vultr 日本服务器

推荐配置：
- **区域**：Tokyo（东京）
- **系统**：Ubuntu 22.04 LTS
- **配置**：1 vCPU, 2GB RAM（$12/月起）

### 2. SSH 连接服务器

```bash
ssh root@your-server-ip
```

### 3. 克隆项目并部署

```bash
# 安装 git
apt update && apt install -y git

# 克隆项目
git clone <your-repo-url> /opt/paper_ai
cd /opt/paper_ai

# 运行一键部署脚本
sudo bash deploy.sh
```

### 4. 配置说明

部署脚本会提示你编辑配置文件，主要需要配置：

#### `config/config.yaml` - 应用配置

```yaml
ai:
  providers:
    claude:
      api_key: "sk-ant-你的Claude API密钥"  # 必填

jwt:
  secret_key: "生产环境强随机字符串"  # 必填，使用 openssl rand -base64 32 生成

database:
  password: "与 .env 中的 DB_PASSWORD 一致"  # 自动生成
```

### 5. 验证部署

```bash
# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 测试 API
curl http://localhost:8080/health
```

---

## 国内访问问题解答

### ✅ 能访问吗？

**可以！** Vultr 日本服务器可以从国内正常访问。

### ⚠️ 会被墙吗？

小概率事件，但建议：

1. **使用 Cloudflare CDN**（免费）
   - 隐藏真实 IP
   - 加速访问
   - 防 DDoS

2. **定期监控**
```bash
# 设置健康检查定时任务
crontab -e

# 每小时检查一次
0 * * * * /opt/paper_ai/scripts/health_check.sh
```

### 🚀 推荐配置域名 + Cloudflare

1. **购买域名**（如 Cloudflare、阿里云）
2. **添加 DNS 记录**
   ```
   A    api    your-server-ip
   ```
3. **开启 Cloudflare 代理**（橙色云朵图标）
4. **配置 SSL**（Cloudflare 自动提供）

这样国内访问：`https://api.yourdomain.com` 非常稳定！

---

## 常用命令

### 服务管理

```bash
# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f app

# 重启服务
docker-compose restart app

# 停止服务
docker-compose down

# 更新服务
./scripts/update.sh
```

### 数据库管理

```bash
# 备份数据库
./scripts/backup.sh

# 恢复数据库
./scripts/restore.sh backup/paper_ai_20250101_120000.sql.gz

# 进入数据库
docker exec -it paper_ai_db psql -U paperai -d paper_ai
```

### 日志查看

```bash
# 应用日志
docker-compose logs -f app

# Nginx 日志
docker-compose logs -f nginx

# 数据库日志
docker-compose logs -f postgres
```

---

## SSL 证书配置（可选但推荐）

### 使用 Let's Encrypt

```bash
# 安装 certbot
apt install -y certbot

# 停止 nginx 容器
docker-compose stop nginx

# 获取证书
certbot certonly --standalone -d api.yourdomain.com

# 复制证书到项目目录
mkdir -p ssl
cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem ssl/
cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem ssl/

# 编辑 nginx.conf，取消 HTTPS 部分注释
vim nginx.conf

# 重启服务
docker-compose up -d
```

### 证书自动续期

```bash
# 创建续期脚本
cat > /opt/scripts/renew_cert.sh << 'EOF'
#!/bin/bash
docker-compose stop nginx
certbot renew
cp /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem /opt/paper_ai/ssl/
cp /etc/letsencrypt/live/api.yourdomain.com/privkey.pem /opt/paper_ai/ssl/
docker-compose start nginx
EOF

chmod +x /opt/scripts/renew_cert.sh

# 添加定时任务（每月 1 号凌晨 2 点）
crontab -e
0 2 1 * * /opt/scripts/renew_cert.sh
```

---

## 监控与维护

### 资源监控

```bash
# 查看 Docker 资源使用
docker stats

# 查看服务器资源
htop

# 查看磁盘使用
df -h

# 查看数据库大小
docker exec paper_ai_db psql -U paperai -d paper_ai -c "SELECT pg_size_pretty(pg_database_size('paper_ai'));"
```

### 日志轮转

Docker Compose 自动处理日志轮转，默认配置：
- 最大 10MB 每个日志文件
- 保留最近 3 个日志文件

如需调整，在 `docker-compose.yml` 中添加：

```yaml
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

---

## 故障排查

### 问题 1：无法访问服务

```bash
# 检查服务状态
docker-compose ps

# 检查防火墙
ufw status

# 检查端口监听
netstat -tlnp | grep 80

# 测试本地访问
curl http://localhost:8080/health
```

### 问题 2：数据库连接失败

```bash
# 查看数据库日志
docker-compose logs postgres

# 测试数据库连接
docker exec -it paper_ai_db psql -U paperai -d paper_ai

# 重启数据库
docker-compose restart postgres
```

### 问题 3：Claude API 调用失败

```bash
# 检查网络连接
docker exec paper_ai_app curl -I https://api.anthropic.com

# 查看应用日志
docker-compose logs -f app

# 验证 API Key
grep api_key config/config.yaml
```

### 问题 4：国内无法访问

```bash
# 从国内服务器测试
ping your-server-ip
curl -I http://your-server-ip

# 检查 IP 是否被墙（使用国内网络）
# 建议：切换到 Cloudflare CDN
```

---

## 性能优化

### 1. 数据库优化

```bash
# 进入数据库
docker exec -it paper_ai_db psql -U paperai -d paper_ai

# 创建索引（如果需要）
CREATE INDEX idx_polish_records_user_id ON polish_records(user_id);
CREATE INDEX idx_polish_records_created_at ON polish_records(created_at);
```

### 2. Nginx 缓存

已在 `nginx.conf` 中配置 gzip 压缩，如需缓存静态资源：

```nginx
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

### 3. 数据库连接池

在 `config/config.yaml` 中调整：

```yaml
database:
  max_idle_conns: 20    # 空闲连接数
  max_open_conns: 200   # 最大连接数
```

---

## 升级和更新

### 自动更新

```bash
# 使用更新脚本（会自动备份数据库）
./scripts/update.sh
```

### 手动更新

```bash
# 1. 备份
./scripts/backup.sh

# 2. 拉取代码
git pull

# 3. 重新构建
docker-compose down
docker-compose up -d --build

# 4. 验证
docker-compose logs -f app
```

---

## 安全加固

### 1. 修改 SSH 端口

```bash
vim /etc/ssh/sshd_config
# Port 22 改为 Port 2222
systemctl restart sshd

# 更新防火墙
ufw allow 2222/tcp
ufw delete allow 22/tcp
```

### 2. 安装 fail2ban

```bash
apt install -y fail2ban
systemctl enable fail2ban
systemctl start fail2ban
```

### 3. 禁用 root 登录

```bash
# 先创建普通用户
adduser deploy
usermod -aG sudo deploy

# 配置 SSH
vim /etc/ssh/sshd_config
# PermitRootLogin no

systemctl restart sshd
```

---

## 成本估算

### Vultr 日本服务器

| 配置 | 价格/月 | 适用场景 |
|------|---------|----------|
| 1核2GB | $12 | 测试/小型项目 |
| 2核4GB | $24 | 中型项目 |
| 4核8GB | $48 | 生产环境 |

### 额外成本

- **域名**：$10-15/年（可选）
- **Cloudflare**：免费
- **Let's Encrypt**：免费
- **Claude API**：按实际使用量计费

---

## 支持与帮助

- **查看完整文档**：[DEPLOYMENT.md](./DEPLOYMENT.md)
- **问题反馈**：提交 Issue
- **紧急支持**：查看日志 `docker-compose logs -f`

---

## 总结

✅ **Vultr 日本服务器完全适合部署此项目**
✅ **国内可以正常访问**
✅ **一键部署，5 分钟上线**
⚠️ **建议配置 Cloudflare CDN 提高稳定性**
⚠️ **定期备份数据库**

祝部署顺利！🚀
