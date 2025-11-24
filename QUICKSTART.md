# 快速开始指南

## 第一步：配置 API Key

编辑 [config/config.yaml](config/config.yaml) 文件，将你的 Claude API Key 填入：

```yaml
ai:
  default_provider: claude
  providers:
    claude:
      api_key: "sk-ant-你的API-Key"  # 替换这里
      base_url: "https://api.anthropic.com"
      model: "claude-3-5-sonnet-20241022"
      timeout: 60s
```

## 第二步：运行服务

### 方式一：使用 Makefile（推荐）

```bash
# 安装依赖并编译
make dev

# 运行服务
make run
```

### 方式二：直接运行

```bash
# 安装依赖
go mod tidy

# 运行
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## 第三步：测试接口

### 1. 健康检查

```bash
curl http://localhost:8080/health
```

### 2. 段落润色（英文）

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This paper discuss the important of machine learning.",
    "style": "academic",
    "language": "en"
  }'
```

### 3. 段落润色（中文）

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这篇文章讨论了机器学习的重要性。",
    "style": "academic",
    "language": "zh"
  }'
```

### 4. 使用测试脚本

```bash
# 在另一个终端窗口运行测试脚本
./test.sh
```

## API 参数说明

### 请求参数

- `content` (必填): 需要润色的文本
- `style` (可选): 润色风格
  - `academic`: 学术风格（默认）
  - `formal`: 正式风格
  - `concise`: 简洁风格
- `language` (可选): 语言
  - `en`: 英文（默认）
  - `zh`: 中文
- `provider` (可选): AI 提供商名称，默认使用配置中的默认提供商

### 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "polished_content": "润色后的文本",
    "original_length": 45,
    "polished_length": 52,
    "suggestions": [],
    "provider_used": "claude",
    "model_used": "claude-3-5-sonnet-20241022"
  },
  "trace_id": "uuid-xxx"
}
```

## 常见问题

### 1. API Key 在哪里获取？

访问 [Anthropic Console](https://console.anthropic.com/) 注册账号并创建 API Key。

### 2. 如何更改端口？

修改 [config/config.yaml](config/config.yaml) 中的 `server.port` 配置。

### 3. 如何添加其他 AI 模型？

参考 [readme.md](readme.md) 中的"扩展指南"部分。

## 下一步

- 查看完整文档: [readme.md](readme.md)
- 了解项目架构: [readme.md#项目结构](readme.md#项目结构)
- 扩展新功能: [readme.md#扩展指南](readme.md#扩展指南)
