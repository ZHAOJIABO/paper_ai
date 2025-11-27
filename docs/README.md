# Paper AI 文档中心

## 📚 文档导航

### 🚀 快速开始
- **[快速开始](./QUICKSTART.md)** - 5分钟快速上手指南

### 🔧 部署文档
- **[部署指南](./deployment/部署指南.md)** - 简洁版部署指南（推荐）
- **[快速部署](./deployment/QUICK_DEPLOY.md)** - 详细的快速部署指南
- **[完整部署文档](./deployment/DEPLOYMENT.md)** - 最详细的部署文档

### 🔌 API 文档
- **[OpenAPI 规范](./api/openapi.yaml)** - API 接口定义
- **[Apifox 导入指南](./api/APIFOX_IMPORT_GUIDE.md)** - 如何导入 API 到 Apifox
- **[前端集成指南](./api/FRONTEND_INTEGRATION.md)** - 前端如何对接 API

### 💻 功能实现
- **[认证实现](./implementation/AUTH_IMPLEMENTATION.md)** - 用户认证系统实现
- **[数据库实现](./implementation/DATABASE_IMPLEMENTATION.md)** - 数据库设计与实现
- **[用户ID生成](./implementation/USERID_GENERATION.md)** - 用户ID生成机制
- **[用户ID改进](./implementation/USERID_IMPROVEMENT.md)** - ID生成优化方案
- **[用户ID查询过滤](./implementation/USERID_QUERY_FILTER.md)** - 查询过滤实现

---

## 🎯 推荐阅读顺序

### 新手入门
1. 先看 [快速开始](./QUICKSTART.md)
2. 然后看 [部署指南](./deployment/部署指南.md)
3. 最后看 [API 文档](./api/openapi.yaml)

### 开发者
1. [认证实现](./implementation/AUTH_IMPLEMENTATION.md)
2. [数据库实现](./implementation/DATABASE_IMPLEMENTATION.md)
3. [前端集成指南](./api/FRONTEND_INTEGRATION.md)

### 运维人员
1. [部署指南](./deployment/部署指南.md)
2. [完整部署文档](./deployment/DEPLOYMENT.md)（出问题时查看）

---

## 📂 文档结构

```
docs/
├── README.md                          # 本文件，文档索引
├── QUICKSTART.md                      # 快速开始
├── deployment/                        # 部署相关
│   ├── 部署指南.md                   # 简洁部署指南（推荐）
│   ├── QUICK_DEPLOY.md               # 快速部署
│   └── DEPLOYMENT.md                 # 完整部署文档
├── api/                               # API 相关
│   ├── openapi.yaml                  # API 规范
│   ├── APIFOX_IMPORT_GUIDE.md       # Apifox 导入
│   └── FRONTEND_INTEGRATION.md       # 前端集成
└── implementation/                    # 功能实现
    ├── AUTH_IMPLEMENTATION.md        # 认证实现
    ├── DATABASE_IMPLEMENTATION.md    # 数据库实现
    ├── USERID_GENERATION.md          # 用户ID生成
    ├── USERID_IMPROVEMENT.md         # ID优化
    └── USERID_QUERY_FILTER.md        # 查询过滤
```

---

## 🔍 快速查找

### 我想...
- **部署项目** → [部署指南](./deployment/部署指南.md)
- **更新服务** → [部署指南](./deployment/部署指南.md#更新服务)
- **对接 API** → [前端集成指南](./api/FRONTEND_INTEGRATION.md)
- **查看 API 定义** → [OpenAPI](./api/openapi.yaml)
- **了解认证流程** → [认证实现](./implementation/AUTH_IMPLEMENTATION.md)
- **了解数据库结构** → [数据库实现](./implementation/DATABASE_IMPLEMENTATION.md)

---

## 💡 贡献文档

如果你想补充或修改文档，请遵循以下规范：

1. 使用 Markdown 格式
2. 文件名使用英文，用下划线或横线分隔
3. 在本 README 中更新索引
4. 保持文档简洁清晰

---

## 📞 获取帮助

- **GitHub Issues** - 提交问题和建议
- **查看日志** - `docker-compose logs -f`
- **文档反馈** - 欢迎提 PR
