# Paper AI - 科研AI服务平台

一个面向科研人员的AI服务平台，提供论文润色、用户认证等功能。支持多种AI模型接入。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)

## 🌟 功能特性

- ✅ **段落润色** - 支持学术、正式、简洁三种风格
- ✅ **多语言支持** - 英文和中文润色
- ✅ **用户认证** - JWT认证，支持注册、登录、Token刷新
- ✅ **数据持久化** - PostgreSQL存储用户和润色记录
- 🔌 **高扩展性** - 轻松接入多种AI模型（Claude、豆包等）
- 🏗️ **低耦合架构** - Clean Architecture，易于维护
- 🚀 **高性能** - 基于Gin框架，支持高并发

## 🚀 快速开始

### 本地开发

```bash
# 1. 克隆项目
git clone <your-repo-url>
cd paper_ai

# 2. 安装依赖
go mod tidy

# 3. 配置文件
cp config/config.example.yaml config/config.yaml
vim config/config.yaml  # 填入 Claude API Key 等配置

# 4. 运行服务
make run
```

**详细说明**：查看 [docs/QUICKSTART.md](docs/QUICKSTART.md)

### Docker 部署（推荐）

```bash
# 1. 配置文件
cp .env.example .env
cp config/config.example.yaml config/config.yaml
vim config/config.yaml  # 修改配置

# 2. 启动服务
docker-compose up -d

# 3. 查看状态
docker-compose ps
```

**部署指南**：查看 [docs/deployment/部署指南.md](docs/deployment/部署指南.md)

## 📖 文档

| 文档 | 说明 |
|------|------|
| [📚 文档中心](docs/README.md) | 所有文档的索引 |
| [🚀 快速开始](docs/QUICKSTART.md) | 5分钟快速上手 |
| [🔧 部署指南](docs/deployment/部署指南.md) | 生产环境部署 |
| [🔌 API文档](docs/api/openapi.yaml) | OpenAPI规范 |
| [💻 功能实现](docs/implementation/) | 各功能的实现文档 |

## 🛠️ 技术栈

- **语言**: Go 1.24+
- **Web框架**: [Gin](https://github.com/gin-gonic/gin)
- **数据库**: PostgreSQL + GORM
- **认证**: JWT
- **日志**: [Zap](https://github.com/uber-go/zap)
- **配置**: [Viper](https://github.com/spf13/viper)
- **AI模型**: Claude 3.5 Sonnet

## 📁 项目结构

```
paper_ai/
├── cmd/server/              # 程序入口
├── internal/
│   ├── api/                 # API层（HTTP处理）
│   ├── service/             # 业务逻辑层
│   ├── domain/              # 领域模型
│   ├── infrastructure/      # 基础设施层
│   └── config/              # 配置管理
├── pkg/                     # 公共包
├── config/                  # 配置文件
├── docs/                    # 📚 文档目录
├── scripts/                 # 部署和运维脚本
├── docker-compose.yml       # Docker编排
├── Dockerfile               # Docker镜像
└── Makefile                 # 构建工具
```

## 🔧 常用命令

```bash
# 开发
make run          # 运行服务
make build        # 编译
make test         # 运行测试

# Docker
docker-compose up -d           # 启动服务
docker-compose logs -f app     # 查看日志
docker-compose restart app     # 重启服务

# 部署
bash scripts/backup.sh         # 备份数据库
bash scripts/update.sh         # 更新服务
```

## 📊 API 示例

### 段落润色

```bash
curl -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "content": "This paper discuss the important of machine learning.",
    "style": "academic",
    "language": "en"
  }'
```

### 用户注册

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

**完整API文档**：[docs/api/openapi.yaml](docs/api/openapi.yaml)

## 🏗️ 架构设计

采用 **Clean Architecture** 设计：

```
HTTP Request → Middleware → Handler → Service → Repository → Database
                                    ↓
                              AI Provider → Claude API
```

**详细说明**：查看各功能的实现文档 [docs/implementation/](docs/implementation/)

## 📝 开发计划

### ✅ 已完成
- 段落润色功能
- 用户认证系统
- 数据持久化
- Claude AI集成
- Docker部署支持

### 🚧 进行中
- API限流功能
- 请求缓存

### 📋 计划中
- 更多AI提供商（OpenAI、Gemini）
- AI代码生成
- 文件批量处理
- 管理后台

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

**代码规范**：遵循 Go 官方规范，使用 `go fmt` 格式化

## 📄 许可证

MIT License

## 📧 联系方式

- 提交 Issue: [GitHub Issues](https://github.com/yourusername/paper_ai/issues)
- 查看文档: [docs/](docs/)

---

**⭐ 如果这个项目对你有帮助，请给个 Star 支持一下！**
