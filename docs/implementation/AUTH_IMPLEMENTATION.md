# 用户认证功能实现文档

## 概述

已成功实现完整的用户注册和登录功能，包括：
- ✅ 用户注册
- ✅ 用户登录
- ✅ JWT访问令牌（Access Token）
- ✅ JWT刷新令牌（Refresh Token）
- ✅ 用户登出
- ✅ 获取当前用户信息
- ✅ JWT认证中间件
- ✅ PostgreSQL持久化存储

## 技术栈

- **JWT**: github.com/golang-jwt/jwt/v5
- **密码加密**: bcrypt (golang.org/x/crypto/bcrypt)
- **数据库**: PostgreSQL + GORM
- **架构**: Clean Architecture（领域驱动设计）

## 数据库表结构

### 1. users 表
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(50),
    login_count INT DEFAULT 0,
    failed_login_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### 2. refresh_tokens 表
```sql
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    device_id VARCHAR(100),
    user_agent VARCHAR(500),
    ip_address VARCHAR(50),
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

## API接口

### 1. 用户注册

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "Password123",
    "confirm_password": "Password123",
    "nickname": "John"
}
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "nickname": "John",
        "avatar_url": "",
        "status": "active",
        "email_verified": false,
        "created_at": "2024-01-01T10:00:00Z"
    },
    "trace_id": "xxx"
}
```

### 2. 用户登录

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "john_doe",
    "password": "Password123"
}
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "eyJhbGc...",
        "refresh_token": "eyJhbGc...",
        "expires_in": 7200,
        "token_type": "Bearer",
        "user": {
            "id": 1,
            "username": "john_doe",
            "email": "john@example.com",
            "nickname": "John",
            "avatar_url": "",
            "status": "active",
            "email_verified": false,
            "created_at": "2024-01-01T10:00:00Z"
        }
    },
    "trace_id": "xxx"
}
```

### 3. 获取当前用户信息

```bash
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "nickname": "John",
        "avatar_url": "",
        "status": "active",
        "email_verified": false,
        "last_login_at": "2024-01-01T10:00:00Z",
        "created_at": "2024-01-01T09:00:00Z"
    },
    "trace_id": "xxx"
}
```

### 4. 刷新访问令牌

```bash
POST /api/v1/auth/refresh
Content-Type: application/json

{
    "refresh_token": "eyJhbGc..."
}
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "access_token": "eyJhbGc...",
        "expires_in": 7200,
        "token_type": "Bearer"
    },
    "trace_id": "xxx"
}
```

### 5. 用户登出

```bash
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
Content-Type: application/json

{
    "refresh_token": "eyJhbGc..."
}
```

**响应示例：**
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "message": "登出成功"
    },
    "trace_id": "xxx"
}
```

## 认证中间件

所有需要认证的接口都已添加JWT认证中间件保护：

- `POST /api/v1/polish` - 论文润色（需要认证）
- `GET /api/v1/polish/records` - 查询记录（需要认证）
- `GET /api/v1/polish/records/:trace_id` - 查询单条记录（需要认证）
- `GET /api/v1/polish/statistics` - 统计信息（需要认证）
- `GET /api/v1/auth/me` - 获取当前用户信息（需要认证）
- `POST /api/v1/auth/logout` - 登出（需要认证）

## 配置说明

在 `config/config.yaml` 中添加了JWT配置：

```yaml
jwt:
  secret_key: "your-secret-key-change-in-production-please-use-random-string"
  access_token_expiry: 7200       # Access Token过期时间（秒），2小时
  refresh_token_expiry: 604800    # Refresh Token过期时间（秒），7天
```

**⚠️ 重要：生产环境必须修改secret_key为强随机字符串！**

## 安全特性

1. **密码强度验证**
   - 最少8位
   - 必须包含字母和数字

2. **密码加密**
   - 使用bcrypt算法
   - 自动加盐

3. **Token机制**
   - Access Token短期有效（2小时）
   - Refresh Token长期有效（7天）
   - Token可撤销

4. **账号状态管理**
   - 支持封禁账号
   - 记录登录失败次数

5. **软删除**
   - 用户数据不会真正删除
   - 支持数据恢复

## 错误码

| 错误码 | 说明 | HTTP状态码 |
|-------|------|-----------|
| 20001 | 用户已存在 | 400 |
| 20002 | 密码错误 | 401 |
| 20003 | 用户不存在 | 404 |
| 20004 | Token无效 | 401 |
| 20005 | Token过期 | 401 |
| 20006 | 密码强度不够 | 400 |
| 20007 | 账号已被封禁 | 403 |
| 20008 | 未授权 | 401 |
| 20009 | 禁止访问 | 403 |

## 测试方法

### 方法一：使用测试脚本

```bash
./test_auth.sh
```

该脚本会自动测试所有认证相关功能。

### 方法二：手动测试

1. 启动服务：
```bash
./paper_ai
```

2. 注册用户：
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234",
    "confirm_password": "Test1234",
    "nickname": "测试用户"
  }'
```

3. 登录：
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test1234"
  }'
```

4. 使用返回的access_token访问需要认证的接口：
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

## 架构说明

采用Clean Architecture分层架构：

```
├── domain/                 # 领域层（核心业务逻辑）
│   ├── entity/            # 实体
│   ├── repository/        # 仓储接口
│   └── model/             # 请求/响应模型
├── service/               # 服务层（业务逻辑）
│   └── auth_service.go    # 认证服务
├── infrastructure/        # 基础设施层
│   ├── security/          # 安全工具（JWT、密码加密）
│   └── persistence/       # 数据持久化
└── api/                   # API层（HTTP处理）
    ├── handler/           # 处理器
    └── middleware/        # 中间件
```

## 后续扩展建议

1. **邮箱验证功能**
   - 注册后发送验证邮件
   - 验证邮箱后激活账号

2. **忘记密码功能**
   - 发送重置密码邮件
   - 验证重置令牌

3. **双因素认证（2FA）**
   - TOTP认证
   - 短信验证码

4. **社交登录**
   - Google OAuth
   - GitHub OAuth

5. **设备管理**
   - 查看登录设备列表
   - 踢出其他设备

6. **登录日志**
   - 记录登录历史
   - 异常登录提醒

## 注意事项

1. **首次运行**：确保PostgreSQL数据库已启动，并正确配置了数据库连接信息
2. **自动迁移**：首次运行时会自动创建数据库表（auto_migrate: true）
3. **生产环境**：务必修改JWT密钥为强随机字符串
4. **HTTPS**：生产环境必须启用HTTPS，保护Token传输安全
5. **限流**：建议添加限流中间件，防止暴力破解

## 完成状态

✅ 所有功能已实现并测试通过
✅ 代码已编译成功
✅ 数据库表已自动迁移
✅ 测试脚本已提供

现在可以启动服务并测试所有认证功能！
