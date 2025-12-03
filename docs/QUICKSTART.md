# 多版本润色功能 - 快速启动指南

## 前置条件

1. Go 1.21+
2. PostgreSQL 12+
3. 已配置 AI Provider（Claude 或 Doubao）

## 步骤 1: 执行数据库迁移

```bash
# 连接到 PostgreSQL
psql -U postgres -d paper_ai

# 执行迁移脚本
\i migrations/001_multi_version_polish.sql

# 验证表创建
\dt polish_*

# 查看初始 Prompt 数据
SELECT id, name, version_type, language, style, is_active
FROM polish_prompts;

# 退出
\q
```

预期结果：
- ✅ `polish_records` 表新增 `mode` 字段
- ✅ `polish_versions` 表创建成功
- ✅ `polish_prompts` 表创建成功，包含 6 条初始数据
- ✅ `users` 表新增 `enable_multi_version` 和 `multi_version_quota` 字段

## 步骤 2: 配置文件

复制配置示例：

```bash
cp config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml`，确保包含以下配置：

```yaml
features:
  multi_version_polish:
    enabled: true           # 启用多版本功能
    default_mode: "single"  # 默认单版本
    max_concurrent: 3       # 最大并发数
```

## 步骤 3: 启动服务

```bash
# 设置环境变量（根据实际情况）
export CLAUDE_API_KEY="your_claude_api_key"
export DOUBAO_API_KEY="your_doubao_api_key"

# 启动服务
go run cmd/server/main.go
```

预期日志输出：

```
[INFO] starting paper_ai service...
[INFO] config loaded successfully
[INFO] database initialized successfully
[INFO] ID generator initialized worker_id=1
[INFO] AI providers initialized providers=[claude, doubao]
[INFO] Prompt service initialized with LRU cache
[INFO] Feature service initialized multi_version_enabled=true default_mode=single
[INFO] Multi-version polish service initialized
[INFO] Routes configured successfully
[INFO] server started port=8080
```

## 步骤 4: 为用户开通多版本功能

### 方法 1: 直接修改数据库（测试用）

```sql
-- 为用户 ID=1 开通多版本功能，无限配额
UPDATE users
SET enable_multi_version = true,
    multi_version_quota = 0
WHERE id = 1;
```

## 步骤 5: 测试多版本润色

### 5.1 登录获取 Token

```bash
# 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'

# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# 保存返回的 access_token
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### 5.2 调用多版本润色接口

```bash
curl -X POST http://localhost:8080/api/v1/polish/multi-version \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This paper discuss the important of machine learning.",
    "style": "academic",
    "language": "en",
    "provider": "claude"
  }'
```

## 故障排查

### 问题 1: 无权限错误

解决方案：为用户开通权限
```sql
UPDATE users SET enable_multi_version = true WHERE id = 1;
```

### 问题 2: Prompt 未找到

检查 Prompt 是否插入：
```sql
SELECT * FROM polish_prompts WHERE is_active = true;
```

---

🎉 恭喜！多版本润色功能已成功启动！
