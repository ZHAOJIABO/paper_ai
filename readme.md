# Paper AI - 科研AI服务平台

一个面向科研人员的AI服务平台，提供论文润色、代码生成、数据分析等功能。目前已实现段落润色功能，支持多种AI模型接入。

## 🌟 功能特性

- ✅ **段落润色**: 支持学术（academic）、正式（formal）、简洁（concise）三种风格
- ✅ **多语言支持**: 支持英文（en）和中文（zh）润色
- 🔌 **高扩展性**: 基于接口设计，轻松接入多种AI模型（Claude、豆包、OpenAI、Gemini等）
- 🏗️ **低耦合架构**: 采用Clean Architecture，各层职责清晰，易于维护
- 🚀 **高性能**: 基于Gin框架，支持高并发请求
- 📝 **结构化日志**: 使用Zap实现高性能日志记录，便于问题排查
- 🔧 **配置驱动**: 支持YAML配置，灵活管理多个AI提供商
- 🔍 **请求追踪**: 每个请求生成唯一TraceID，便于追踪问题
- 🛡️ **完善的错误处理**: 统一错误码体系，返回友好错误信息

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **日志**: [Zap](https://github.com/uber-go/zap)
- **配置管理**: [Viper](https://github.com/spf13/viper)
- **AI模型**: Claude 3.5 Sonnet、豆包大模型 (可扩展至OpenAI、Gemini等)

## 📁 项目结构

```
paper_ai/
├── cmd/
│   └── server/
│       └── main.go                       # 程序入口
├── internal/
│   ├── api/                              # API层（HTTP处理）
│   │   ├── handler/
│   │   │   └── polish.go                 # 润色请求处理器
│   │   ├── middleware/
│   │   │   ├── logger.go                 # 日志中间件
│   │   │   ├── recovery.go              # Panic恢复中间件
│   │   │   └── cors.go                   # CORS跨域中间件
│   │   └── router/
│   │       └── router.go                 # 路由配置
│   ├── service/
│   │   └── polish.go                     # 业务逻辑层
│   ├── domain/
│   │   └── model/
│   │       └── polish.go                 # 领域模型（参数验证）
│   ├── infrastructure/                   # 基础设施层
│   │   └── ai/
│   │       ├── provider.go               # AI提供商接口定义
│   │       ├── types/
│   │       │   └── types.go              # 类型定义（避免循环依赖）
│   │       ├── factory.go                # 工厂模式（创建provider）
│   │       ├── claude/
│   │       │   └── client.go             # Claude客户端实现
│   │       └── doubao/
│   │           └── client.go             # 豆包客户端实现
│   └── config/
│       └── config.go                     # 配置管理
├── pkg/                                  # 公共包（可被外部引用）
│   ├── errors/
│   │   └── errors.go                     # 自定义错误类型
│   ├── logger/
│   │   └── logger.go                     # 日志工具
│   └── response/
│       └── response.go                   # 统一响应格式
├── config/
│   ├── config.yaml                       # 配置文件（需配置API Key）
│   └── config.example.yaml               # 配置示例
├── Makefile                              # 构建工具
├── test.sh                               # 测试脚本
├── QUICKSTART.md                         # 快速开始指南
├── .gitignore                            # Git忽略配置
└── readme.md                             # 本文档
```

## 🚀 快速开始

### 1. 克隆项目（如果是从Git获取）

```bash
cd /path/to/paper_ai
```

### 2. 安装依赖

```bash
go mod tidy
```

或使用Makefile：

```bash
make deps
```

### 3. 配置API Key

复制示例配置文件并编辑：

```bash
cp config/config.example.yaml config/config.yaml
vim config/config.yaml
```

修改配置文件，填入你的AI提供商API Key。

#### 使用Claude（推荐用于英文）

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

ai:
  default_provider: claude
  providers:
    claude:
      api_key: "sk-ant-你的API-Key-在这里"  # 替换为你的Claude API Key
      base_url: "https://api.anthropic.com"
      model: "claude-3-5-sonnet-20241022"
      timeout: 60s
```

> **获取Claude API Key**: 访问 [Anthropic Console](https://console.anthropic.com/) 注册账号并创建API Key

#### 使用豆包（推荐用于中文）

```yaml
server:
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

ai:
  default_provider: doubao
  providers:
    doubao:
      api_key: "your-doubao-api-key"         # 替换为你的豆包API Key
      base_url: "https://ark.cn-beijing.volces.com/api/v3"  # 豆包API地址
      model: "ep-xxxxx-xxxxx"                # 替换为你的模型endpoint ID
      timeout: 60s
```

> **获取豆包API Key**: 访问 [火山引擎-豆包大模型](https://console.volcengine.com/ark) 注册并创建推理接入点

#### 同时配置多个提供商

```yaml
ai:
  default_provider: doubao  # 默认使用豆包
  providers:
    claude:
      api_key: "sk-ant-xxx"
      base_url: "https://api.anthropic.com"
      model: "claude-3-5-sonnet-20241022"
      timeout: 60s
    doubao:
      api_key: "your-doubao-api-key"
      base_url: "https://ark.cn-beijing.volces.com/api/v3"
      model: "ep-xxxxx-xxxxx"
      timeout: 60s
```

### 4. 运行服务

**方式一：使用Makefile（推荐）**

```bash
make run
```

**方式二：直接运行**

```bash
go run cmd/server/main.go
```

**方式三：编译后运行**

```bash
make build
./paper_ai
```

服务将在 `http://localhost:8080` 启动。

### 5. 测试接口

**健康检查：**

```bash
curl http://localhost:8080/health
```

**段落润色（英文 - 使用Claude）：**

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This paper discuss the important of machine learning in modern software development.",
    "style": "academic",
    "language": "en",
    "provider": "claude"
  }'
```

**段落润色（中文 - 使用豆包）：**

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这篇文章讨论了机器学习在软件开发中的作用。",
    "style": "academic",
    "language": "zh",
    "provider": "doubao"
  }'
```

**使用默认提供商（不指定provider参数）：**

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一段需要润色的文本。",
    "style": "academic",
    "language": "zh"
  }'
```

**使用测试脚本：**

```bash
chmod +x test.sh
./test.sh
```

## 📖 API文档

### 健康检查

**接口**: `GET /health`

**响应示例**:
```json
{
  "status": "ok"
}
```

---

### 段落润色

**接口**: `POST /api/v1/polish`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| content | string | 是 | 需要润色的文本内容（最大10000字符） | "This is a test paragraph." |
| style | string | 否 | 润色风格（默认：academic） | academic/formal/concise |
| language | string | 否 | 目标语言（默认：en） | en/zh |
| provider | string | 否 | AI提供商（默认使用配置的默认提供商） | claude/doubao |

**style 参数说明**:
- `academic`: 学术风格 - 适用于学术论文，更加正式、精确
- `formal`: 正式风格 - 适用于正式文档，更加专业
- `concise`: 简洁风格 - 去除冗余，保持简洁清晰

**请求示例**:

```json
{
  "content": "This paper discuss the important of machine learning.",
  "style": "academic",
  "language": "en",
  "provider": "claude"
}
```

**响应格式**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "polished_content": "This paper discusses the importance of machine learning.",
    "original_length": 52,
    "polished_length": 58,
    "suggestions": [],
    "provider_used": "claude",
    "model_used": "claude-3-5-sonnet-20241022"
  },
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| code | int | 错误码（0表示成功） |
| message | string | 响应消息 |
| data.polished_content | string | 润色后的文本 |
| data.original_length | int | 原始文本长度 |
| data.polished_length | int | 润色后文本长度 |
| data.suggestions | []string | 改进建议（预留字段） |
| data.provider_used | string | 实际使用的AI提供商 |
| data.model_used | string | 实际使用的模型 |
| trace_id | string | 请求追踪ID（用于问题排查） |

## 🔧 配置说明

### 服务器配置

```yaml
server:
  port: 8080              # 服务监听端口
  read_timeout: 30s       # HTTP读取超时时间
  write_timeout: 30s      # HTTP写入超时时间
```

### AI提供商配置

```yaml
ai:
  default_provider: claude  # 默认使用的AI提供商名称
  providers:                # 提供商配置列表
    claude:                 # Claude提供商（适合英文润色）
      api_key: "xxx"        # API密钥
      base_url: "xxx"       # API基础URL
      model: "xxx"          # 模型名称
      timeout: 60s          # 请求超时时间
    doubao:                 # 豆包提供商（适合中文润色）
      api_key: "xxx"        # API密钥
      base_url: "xxx"       # API基础URL
      model: "xxx"          # 模型endpoint ID
      timeout: 60s          # 请求超时时间
    # 可以配置更多提供商
    # openai:
    #   api_key: "sk-xxx"
    #   base_url: "https://api.openai.com"
    #   model: "gpt-4"
    #   timeout: 60s
```

### 环境变量

- `CONFIG_PATH`: 配置文件路径（默认：`./config/config.yaml`）

使用方式：
```bash
CONFIG_PATH=/path/to/config.yaml ./paper_ai
```

## 🎯 错误码说明

| 错误码 | 说明 | HTTP状态码 | 解决方案 |
|-------|------|-----------|---------|
| 0 | 成功 | 200 | - |
| 10001 | 参数错误 | 400 | 检查请求参数是否正确 |
| 10002 | AI服务错误 | 500 | 检查AI服务是否正常，API Key是否正确 |
| 10003 | 限流错误 | 429 | 降低请求频率 |
| 10004 | 超时错误 | 504 | 增加timeout配置或稍后重试 |
| 10005 | 内部错误 | 500 | 查看服务器日志 |
| 10006 | AI提供商不存在 | 400 | 检查provider参数是否正确 |
| 10007 | 配置错误 | 500 | 检查配置文件是否正确 |

**错误响应示例**:

```json
{
  "code": 10001,
  "message": "content cannot be empty",
  "trace_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 🔌 扩展指南

### 支持的AI提供商

目前已支持以下AI提供商：

| 提供商 | 适用场景 | API文档 |
|-------|---------|---------|
| **Claude** | 英文润色，学术写作 | [Anthropic API](https://docs.anthropic.com/) |
| **豆包（Doubao）** | 中文润色，本土化需求 | [火山引擎豆包](https://www.volcengine.com/docs/82379) |

### 添加新的AI提供商（以OpenAI为例）

#### 步骤1: 创建客户端实现

在 `internal/infrastructure/ai/openai/` 目录下创建 `client.go`：

```go
package openai

import (
    "context"
    "paper_ai/internal/infrastructure/ai/types"
)

type Client struct {
    apiKey  string
    baseURL string
    model   string
}

func NewClient(apiKey, baseURL, model string, timeout time.Duration) *Client {
    return &Client{
        apiKey:  apiKey,
        baseURL: baseURL,
        model:   model,
    }
}

// 实现 AIProvider 接口
func (c *Client) Polish(ctx context.Context, req *types.PolishRequest) (*types.PolishResponse, error) {
    // 实现OpenAI的调用逻辑
    // ...
}
```

#### 步骤2: 在工厂中注册

编辑 `internal/infrastructure/ai/factory.go`，在 `InitProviders` 方法中添加：

```go
case "openai":
    client := openai.NewClient(
        providerCfg.APIKey,
        providerCfg.BaseURL,
        providerCfg.Model,
        providerCfg.Timeout,
    )
    f.providers[name] = client
case "doubao":
    client := doubao.NewClient(
        providerCfg.APIKey,
        providerCfg.BaseURL,
        providerCfg.Model,
        providerCfg.Timeout,
    )
    f.providers[name] = client
```

#### 步骤3: 添加配置

在 `config/config.yaml` 中添加OpenAI配置：

```yaml
ai:
  default_provider: doubao  # 可选择使用豆包或其他
  providers:
    claude:
      # ... Claude配置
    doubao:
      # ... 豆包配置
    openai:
      api_key: "sk-xxx"
      base_url: "https://api.openai.com"
      model: "gpt-4"
      timeout: 60s
```

完成！现在可以通过 `"provider": "openai"` 参数使用OpenAI了。

### 添加新功能（以代码生成为例）

#### 步骤1: 在接口中添加新方法

编辑 `internal/infrastructure/ai/types/types.go` 添加新类型：

```go
// CodeGenRequest 代码生成请求
type CodeGenRequest struct {
    Description string `json:"description"`
    Language    string `json:"language"`
}

// CodeGenResponse 代码生成响应
type CodeGenResponse struct {
    Code         string `json:"code"`
    Explanation  string `json:"explanation"`
    ProviderUsed string `json:"provider_used"`
    ModelUsed    string `json:"model_used"`
}
```

编辑 `internal/infrastructure/ai/provider.go`：

```go
type AIProvider interface {
    Polish(ctx context.Context, req *types.PolishRequest) (*types.PolishResponse, error)
    GenerateCode(ctx context.Context, req *types.CodeGenRequest) (*types.CodeGenResponse, error)
}
```

#### 步骤2: 实现各提供商的方法

在 `internal/infrastructure/ai/claude/client.go` 中实现：

```go
func (c *Client) GenerateCode(ctx context.Context, req *types.CodeGenRequest) (*types.CodeGenResponse, error) {
    // 实现代码生成逻辑
}
```

#### 步骤3: 创建Service层

创建 `internal/service/codegen.go`：

```go
package service

type CodeGenService struct {
    providerFactory *ai.ProviderFactory
}

func NewCodeGenService(factory *ai.ProviderFactory) *CodeGenService {
    return &CodeGenService{providerFactory: factory}
}

func (s *CodeGenService) GenerateCode(ctx context.Context, req *model.CodeGenRequest) (*types.CodeGenResponse, error) {
    // 业务逻辑
}
```

#### 步骤4: 创建Handler和路由

创建 `internal/api/handler/codegen.go` 并在 `router.go` 中注册路由。

## 🏗️ 架构设计

### 核心设计原则

1. **依赖倒置原则（DIP）**
   - 通过接口抽象AI提供商，上层不依赖具体实现
   - 便于测试和替换实现

2. **开闭原则（OCP）**
   - 对扩展开放：轻松添加新的AI提供商
   - 对修改关闭：添加新功能不影响现有代码

3. **单一职责原则（SRP）**
   - 每层只负责自己的职责
   - API层：HTTP处理
   - Service层：业务逻辑
   - Infrastructure层：外部服务集成

### 依赖注入流程

```
main.go
  ↓ 创建
Factory
  ↓ 注入
Service
  ↓ 注入
Handler
  ↓ 注册
Router
```

### 请求处理流程

```
HTTP Request
  ↓
Middleware (Logger, CORS, Recovery)
  ↓
Router → Handler
  ↓
Service (业务逻辑 + 参数验证)
  ↓
AIProvider Interface
  ↓
Concrete Provider (Claude/OpenAI/...)
  ↓
AI API
  ↓
Response → Client
```

## 📊 性能优化建议

### 1. 添加缓存

对于相同的输入内容，可以缓存结果：

```go
// 使用Redis缓存
func (s *PolishService) Polish(ctx context.Context, req *model.PolishRequest) (*types.PolishResponse, error) {
    // 生成缓存key
    cacheKey := generateCacheKey(req)

    // 尝试从缓存获取
    if cached := s.cache.Get(cacheKey); cached != nil {
        return cached, nil
    }

    // 调用AI服务
    resp, err := s.provider.Polish(ctx, aiReq)

    // 缓存结果
    s.cache.Set(cacheKey, resp, 24*time.Hour)

    return resp, nil
}
```

### 2. 添加限流

使用中间件限制请求频率：

```go
// internal/api/middleware/ratelimit.go
func RateLimit() gin.HandlerFunc {
    limiter := rate.NewLimiter(10, 20) // 每秒10个请求，桶容量20
    return func(c *gin.Context) {
        if !limiter.Allow() {
            response.Error(c, apperrors.NewRateLimitError("too many requests"))
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. 异步处理

对于长文本润色，使用异步任务：

```go
// 返回任务ID
taskID := uuid.New().String()

// 异步处理
go func() {
    result := service.Polish(ctx, req)
    cache.Set(taskID, result)
}()

// 返回任务ID
return gin.H{"task_id": taskID}
```

## 🧪 测试

### 单元测试示例

```go
// internal/service/polish_test.go
func TestPolishService_Polish(t *testing.T) {
    // 创建mock provider
    mockProvider := &MockAIProvider{}
    factory := &MockFactory{provider: mockProvider}

    service := NewPolishService(factory)

    req := &model.PolishRequest{
        Content: "test content",
        Style:   "academic",
    }

    resp, err := service.Polish(context.Background(), req)

    assert.NoError(t, err)
    assert.NotEmpty(t, resp.PolishedContent)
}
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行指定包测试
go test ./internal/service

# 带覆盖率
go test -cover ./...
```

## 🚀 部署

### Docker部署（推荐）

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o paper_ai cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/paper_ai .
COPY config/config.example.yaml ./config/config.yaml

EXPOSE 8080
CMD ["./paper_ai"]
```

构建和运行：

```bash
# 构建镜像
docker build -t paper_ai:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/config:/app/config \
  --name paper_ai \
  paper_ai:latest
```

### 二进制部署

```bash
# 编译
make build

# 上传到服务器
scp paper_ai config/config.yaml user@server:/opt/paper_ai/

# 在服务器上运行
./paper_ai
```

### 使用systemd管理

创建 `/etc/systemd/system/paper_ai.service`：

```ini
[Unit]
Description=Paper AI Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/paper_ai
ExecStart=/opt/paper_ai/paper_ai
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

管理服务：

```bash
sudo systemctl enable paper_ai
sudo systemctl start paper_ai
sudo systemctl status paper_ai
```

## 📝 开发计划

### 已完成
- [x] 段落润色功能
- [x] Claude AI集成
- [x] 豆包AI集成
- [x] 多风格支持（academic/formal/concise）
- [x] 多语言支持（en/zh）
- [x] 多提供商支持与切换
- [x] 统一错误处理
- [x] 结构化日志
- [x] 请求追踪（TraceID）
- [x] CORS支持
- [x] 健康检查接口
- [x] 优雅关闭

### 待开发
- [ ] 用户认证系统（JWT）
- [ ] API限流功能
- [ ] 请求缓存（Redis）
- [ ] 支持更多AI提供商（OpenAI、Gemini、文心一言、通义千问）
- [ ] AI代码生成功能
- [ ] AI论文段落生成功能
- [ ] AI数据分析功能
- [ ] 文件上传支持（批量处理）
- [ ] 异步任务队列
- [ ] 监控和指标采集（Prometheus）
- [ ] 链路追踪（OpenTelemetry）
- [ ] 单元测试和集成测试
- [ ] API文档（Swagger）
- [ ] 管理后台界面

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

### 贡献步骤

1. Fork本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

### 代码规范

- 遵循Go官方代码规范
- 使用 `go fmt` 格式化代码
- 使用 `golangci-lint` 进行代码检查
- 添加必要的注释
- 编写单元测试

## 📄 许可证

MIT License

Copyright (c) 2024

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

## 📧 联系方式

如有问题或建议，欢迎通过以下方式联系：

- 提交Issue: [GitHub Issues](https://github.com/yourusername/paper_ai/issues)
- 邮箱: your.email@example.com

## 🙏 致谢

感谢以下开源项目：

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Zap Logger](https://github.com/uber-go/zap)
- [Viper](https://github.com/spf13/viper)
- [Anthropic Claude](https://www.anthropic.com/)
- [火山引擎豆包大模型](https://www.volcengine.com/product/doubao)

---

**⭐ 如果这个项目对你有帮助，请给个Star支持一下！**
