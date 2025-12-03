# 多版本润色功能实现文档

## 功能概述

本项目已成功实现了多版本润色功能，该功能允许用户一次请求生成3个不同强度的润色版本：
- **Conservative（保守版本）**：仅修正语法错误，保持原文结构
- **Balanced（平衡版本）**：适度优化语法和结构
- **Aggressive（激进版本）**：大幅提升写作质量和学术水平

## 架构设计

### 核心特性

1. **并发多版本生成**：使用 Goroutine 并发调用AI，3个版本同时生成，总耗时约等于单个版本
2. **主从表设计**：扩展性强，易于维护
   - 主表：`polish_records` - 存储公共信息
   - 从表：`polish_versions` - 存储各版本详情
3. **Prompt数据库管理**：Prompt存储在数据库中，支持热更新、A/B测试
4. **权限控制**：三级开关（全局/用户/请求）精确控制功能开放范围
5. **LRU缓存**：Prompt缓存机制，减少数据库查询

### 技术栈

- **语言**：Go 1.21+
- **Web框架**：Gin
- **ORM**：GORM
- **数据库**：PostgreSQL
- **ID生成**：Snowflake算法

## 数据库迁移

### 1. 执行迁移脚本

```bash
# 进入项目根目录
cd paper_ai

# 连接到PostgreSQL数据库
psql -U postgres -d paper_ai -f migrations/001_multi_version_polish.sql
```

### 2. 验证表创建

```sql
-- 检查表是否创建成功
\dt polish_*

-- 查看Prompt初始数据
SELECT id, name, version_type, language, style, is_active FROM polish_prompts;
```

### 3. 回滚（如需要）

```bash
psql -U postgres -d paper_ai -f migrations/001_multi_version_polish_rollback.sql
```

## 配置说明

### config.yaml 配置示例

```yaml
# 功能开关配置
features:
  multi_version_polish:
    enabled: true           # 全局开关：是否启用多版本功能
    default_mode: "single"  # 默认模式：single（单版本）或 multi（多版本）
    max_concurrent: 3       # 最大并发数
```

## API 使用说明

### 1. 多版本润色接口

**端点**：`POST /api/v1/polish/multi-version`

**请求头**：
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**请求体**：
```json
{
  "content": "This paper discuss the important of machine learning.",
  "style": "academic",
  "language": "en",
  "provider": "claude",
  "versions": ["balanced", "aggressive"]  // 可选，不传则生成全部3个版本
}
```

**响应示例**：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "trace_id": "1234567890123456",
    "original_length": 56,
    "provider_used": "claude",
    "versions": {
      "balanced": {
        "polished_content": "This paper examines the significance of machine learning...",
        "polished_length": 65,
        "suggestions": [
          "Changed 'discuss' to 'examines'",
          "Changed 'important' to 'significance'"
        ],
        "process_time_ms": 1350,
        "model_used": "claude-3-5-sonnet",
        "status": "success"
      },
      "aggressive": {
        "polished_content": "This manuscript critically analyzes the pivotal role...",
        "polished_length": 78,
        "suggestions": [
          "Rewrote with sophisticated academic language",
          "Enhanced logical structure"
        ],
        "process_time_ms": 1500,
        "model_used": "claude-3-5-sonnet",
        "status": "success"
      }
    }
  }
}
```

### 2. 管理接口

#### 2.1 Prompt管理

**列出所有Prompts**：
```
GET /api/v1/admin/prompts?version_type=conservative&is_active=true
```

**创建Prompt**：
```
POST /api/v1/admin/prompts
Content-Type: application/json

{
  "name": "New Conservative Prompt",
  "version_type": "conservative",
  "language": "en",
  "style": "academic",
  "system_prompt": "You are an academic writing assistant...",
  "user_prompt_template": "Polish: {{content}}",
  "is_active": true
}
```

**激活/停用Prompt**：
```
POST /api/v1/admin/prompts/:id/activate
POST /api/v1/admin/prompts/:id/deactivate
```

#### 2.2 用户权限管理

**为用户开通多版本功能**：
```
POST /api/v1/admin/users/:user_id/multi-version/enable
Content-Type: application/json

{
  "quota": 0  // 0表示无限配额
}
```

**关闭用户多版本功能**：
```
POST /api/v1/admin/users/:user_id/multi-version/disable
```

**查询用户状态**：
```
GET /api/v1/admin/users/:user_id/multi-version/status
```

## 代码集成

### 主程序初始化示例

在 `cmd/server/main.go` 中添加以下初始化代码：

```go
// 1. 初始化 Repository
polishRepo := persistence.NewPolishRepository(db)
versionRepo := persistence.NewPolishVersionRepository(db)
promptRepo := persistence.NewPolishPromptRepository(db)
userRepo := persistence.NewUserRepository(db)

// 2. 初始化 Service
promptService := service.NewPromptService(promptRepo)

featureConfig := &service.FeatureConfig{
    MultiVersionEnabled: config.Get().Features.MultiVersionPolish.Enabled,
    DefaultMode:         config.Get().Features.MultiVersionPolish.DefaultMode,
    MaxConcurrent:       config.Get().Features.MultiVersionPolish.MaxConcurrent,
}
featureService := service.NewFeatureService(userRepo, featureConfig)

// 单版本润色服务（保留原有）
polishService := service.NewPolishService(providerFactory, polishRepo)

// 多版本润色服务（新增）
multiVersionService := service.NewPolishMultiVersionService(
    providerFactory,
    polishRepo,
    versionRepo,
    promptService,
    featureService,
)

// 3. 初始化 Handler
polishHandler := handler.NewPolishHandler(polishService)
multiVersionHandler := handler.NewPolishMultiVersionHandler(multiVersionService)

// Admin handlers
promptAdminHandler := admin.NewPromptAdminHandler(promptRepo)
featureAdminHandler := admin.NewFeatureAdminHandler(userRepo)

// 4. 设置路由（需要修改router.Setup签名）
r := router.Setup(
    polishHandler,
    multiVersionHandler,
    queryHandler,
    comparisonHandler,
    authHandler,
    promptAdminHandler,
    featureAdminHandler,
    jwtManager,
)
```

## 扩展性设计

### 添加新的版本类型

1. 在 `internal/domain/entity/polish_version.go` 中添加新的版本类型常量
2. 在数据库中插入新的Prompt模板
3. 无需修改代码，系统自动支持

### 添加新的语言或风格

1. 在数据库 `polish_prompts` 表中插入新的Prompt
2. 使用Prompt管理接口创建即可

## 性能优化

### 短期优化建议

1. **智能版本生成**：先生成Balanced版本，用户按需生成其他版本
2. **Redis缓存**：缓存相同内容的润色结果
3. **异步队列**：使用消息队列处理多版本请求

### 中期优化建议

1. **A/B测试平台**：系统化的Prompt A/B测试
2. **版本质量评分**：为每个版本生成质量评分
3. **实时推送**：使用SSE或WebSocket实时返回结果

## 监控指标

建议监控以下指标：

1. **使用率**：多版本功能的使用比例
2. **成功率**：各版本的成功率
3. **耗时**：并发调用的总耗时
4. **成本**：AI API调用成本
5. **版本选择率**：用户最常选择哪个版本

## 故障排查

### 常见问题

**问题1：无权限使用多版本功能**
- 检查全局开关：`features.multi_version_polish.enabled`
- 检查用户权限：`users.enable_multi_version`
- 检查配额：`users.multi_version_quota`

**问题2：部分版本失败**
- 查看 `polish_versions` 表的 `error_message` 字段
- 检查Prompt是否正确
- 检查AI Provider配置

**问题3：响应时间过长**
- 检查并发调用是否正常工作
- 检查AI Provider的响应时间
- 考虑增加超时时间

## 安全建议

1. **管理接口权限**：只允许管理员访问Prompt管理接口
2. **Prompt审核**：建立Prompt修改审批流程
3. **配额控制**：合理设置用户配额，控制成本
4. **数据备份**：定期备份Prompt数据

## 版本历史

- **v1.0.0** (2025-12-03)
  - 初始版本，实现多版本润色核心功能
  - 支持3种版本类型：Conservative、Balanced、Aggressive
  - 实现Prompt数据库管理
  - 实现权限控制和配额管理

## 联系方式

如有问题，请提交Issue或联系开发团队。
