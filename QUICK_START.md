# Cloudflare 配置快速参考

## 🚀 5 分钟快速上手

### 你需要准备的：
- ✅ 一个邮箱（注册 Cloudflare）
- ✅ 信用卡或 PayPal（购买域名）
- ✅ Vultr 日本服务器 IP 地址
- ✅ SSH 访问权限

### 配置流程：

```
1️⃣ Cloudflare 注册 + 购买域名
   ↓ 10分钟
2️⃣ 配置 DNS 解析（A 记录）
   ↓ 2分钟
3️⃣ 配置 SSL 模式（Full strict）
   ↓ 1分钟
4️⃣ 服务器运行配置脚本
   ↓ 5分钟
5️⃣ 重新开启 Cloudflare CDN
   ↓ 1分钟
✅ 完成！
```

---

## 📝 详细步骤

### 步骤 1: Cloudflare 注册并购买域名

1. 访问 https://www.cloudflare.com
2. 点击 **Sign Up** 注册账号
3. 登录后，点击 **Domain Registration** 或 **注册域名**
4. 搜索并购买你想要的域名（推荐 .com）
5. 完成支付

**时间：10 分钟**

---

### 步骤 2: 配置 DNS

1. 进入你的域名管理页面
2. 点击 **DNS** 选项卡
3. 添加两条 A 记录：

**记录 1 - 根域名:**
```
Type: A
Name: @
IPv4 address: [你的服务器IP]
Proxy status: ⚪ DNS only（灰色云朵 - 临时）
```

**记录 2 - www:**
```
Type: A
Name: www
IPv4 address: [你的服务器IP]
Proxy status: ⚪ DNS only（灰色云朵 - 临时）
```

⚠️ **重要：** 首次配置时请使用**灰色云朵**（DNS only），方便申请 SSL 证书

**时间：2 分钟**

---

### 步骤 3: 配置 SSL 模式

1. 点击 **SSL/TLS** 选项卡
2. 加密模式选择：**Full (strict)**

**时间：1 分钟**

---

### 步骤 4: 服务器配置

SSH 连接到你的 Vultr 服务器：

```bash
# 连接服务器
ssh root@你的服务器IP

# 进入项目目录
cd /path/to/paper_ai

# 如果项目还没部署，先部署
git clone [你的项目仓库]
cd paper_ai

# 运行 Cloudflare 专用配置脚本
sudo ./setup_cloudflare.sh
```

按照提示输入：
- 你的域名
- 是否配置 www（输入 y）

**等待脚本自动完成配置**

**时间：5 分钟**

---

### 步骤 5: 开启 Cloudflare CDN

SSL 证书申请成功后，回到 Cloudflare：

1. 进入 **DNS** 设置
2. 点击两条 A 记录的**灰色云朵**，变成**橙色云朵**🟠
3. 等待 1-2 分钟生效

**时间：1 分钟**

---

## ✅ 验证配置

在本地电脑执行：

```bash
# 测试访问
curl -I https://your-domain.com/health

# 应该返回 200 OK
```

在浏览器访问：
- https://your-domain.com
- https://www.your-domain.com

检查：
- ✅ 显示 🔒 锁图标
- ✅ 证书有效
- ✅ 页面正常加载

---

## 🔧 常见问题快速解决

### ❌ 证书申请失败

**原因：** Cloudflare 代理开启导致验证失败

**解决：**
1. Cloudflare DNS 设置中，点击橙色云朵变成灰色
2. 等待 2-3 分钟
3. 重新运行 `./setup_cloudflare.sh`
4. 成功后再开启橙色云朵

### ❌ 522 错误（Cloudflare 无法连接）

**解决：**
```bash
# 检查服务状态
docker-compose ps

# 查看日志
docker-compose logs nginx

# 重启服务
docker-compose restart
```

### ❌ DNS 解析不生效

**解决：**
```bash
# 检查 DNS
nslookup your-domain.com

# 清除本地 DNS 缓存（Mac）
sudo dscacheutil -flushcache

# 等待 5-10 分钟后重试
```

---

## 📊 配置检查清单

在完成配置后，逐项检查：

- [ ] Cloudflare 账号已注册
- [ ] 域名已购买
- [ ] DNS A 记录已添加（@ 和 www）
- [ ] DNS Proxy status 为 🟠 Proxied（橙色云朵）
- [ ] SSL/TLS 模式为 Full (strict)
- [ ] 服务器已运行 setup_cloudflare.sh
- [ ] SSL 证书申请成功
- [ ] https://域名 可以访问
- [ ] 显示绿色锁图标
- [ ] 健康检查接口返回 200

---

## 💰 费用

- **域名（.com）**: ~$10/年
- **Cloudflare**: 免费
- **SSL 证书**: 免费（Let's Encrypt）
- **服务器**: $5-10/月

**总计**: 约 $70-130/年

---

## 📞 获取帮助

**查看详细文档：**
- [CLOUDFLARE_SETUP.md](CLOUDFLARE_SETUP.md) - 完整详细指南

**查看日志：**
```bash
docker-compose logs -f
```

**测试配置：**
```bash
# 测试 HTTPS
curl -I https://your-domain.com

# 检查 SSL
openssl s_client -connect your-domain.com:443
```

---

## 🎯 下一步

配置完成后，你可以：

1. **配置 API 客户端**
   - 将 API 地址改为你的域名

2. **监控服务**
   - 查看 Cloudflare 分析面板
   - 监控流量和安全威胁

3. **优化性能**
   - 配置缓存规则
   - 开启额外的优化选项

4. **备份数据**
   ```bash
   ./scripts/backup.sh
   ```

祝你使用愉快！🎉
