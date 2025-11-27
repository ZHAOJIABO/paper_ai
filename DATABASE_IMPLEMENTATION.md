# PostgreSQL 数据库持久化方案实施文档

## ✅ 实施完成

已成功实现 PostgreSQL + TEXT字段 的高扩展性、低耦合持久化方案。

## 📊 架构总览

```
领域层（Domain Layer）
  ├── Entity（实体）- 纯业务模型
  ├── Repository Interface（仓储接口）- 定义契约
  └── Query Options（查询选项）- Options模式

服务层（Service Layer）
  └── PolishService - 依赖Repository接口

基础设施层（Infrastructure Layer）
  ├── Database（数据库管理）
  └── Persistence（仓储实现）
      ├── PO模型（包含GORM标签）
      └── Repository实现
```

## 📁 新增/修改的文件

### 新增文件

**领域层：**
- `internal/domain/entity/polish_record.go` - 润色记录实体
- `internal/domain/repository/polish_repository.go` - 仓储接口定义
- `internal/domain/repository/query_options.go` - 查询选项（Options模式）

**基础设施层：**
- `internal/infrastructure/database/database.go` - 数据库连接管理
- `internal/infrastructure/persistence/models.go` - 持久化对象（PO）
- `internal/infrastructure/persistence/polish_repository_impl.go` - 仓储实现
- `internal/infrastructure/persistence/polish_repository_stats.go` - 统计功能实现

**API层：**
- `internal/api/handler/polish_query_handler.go` - 查询处理器

### 修改文件

- `internal/config/config.go` - 添加数据库配置
- `internal/service/polish.go` - 集成数据库记录功能
- `internal/api/router/router.go` - 添加查询路由
- `cmd/server/main.go` - 依赖注入
- `config/config.yaml` - 添加数据库配置
- `go.mod` - 添加GORM依赖

## 🗄️ 数据库表结构

```sql
CREATE TABLE polish_records (
    id BIGSERIAL PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL UNIQUE,

    -- 输入信息
    original_content TEXT NOT NULL,
    style VARCHAR(20) NOT NULL,
    language VARCHAR(10) NOT NULL,

    -- 输出信息
    polished_content TEXT NOT NULL,
    original_length INT NOT NULL,
    polished_length INT NOT NULL,

    -- AI信息
    provider VARCHAR(50) NOT NULL,
    model VARCHAR(100) NOT NULL,

    -- 性能指标
    process_time_ms INT DEFAULT 0,

    -- 状态信息
    status VARCHAR(20) NOT NULL DEFAULT 'success',
    error_message TEXT,

    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 索引
CREATE UNIQUE INDEX idx_trace_id ON polish_records(trace_id);
CREATE INDEX idx_provider ON polish_records(provider);
CREATE INDEX idx_status ON polish_records(status);
CREATE INDEX idx_language ON polish_records(language);
CREATE INDEX idx_style ON polish_records(style);
CREATE INDEX idx_created_at ON polish_records(created_at);
CREATE INDEX idx_process_time ON polish_records(process_time_ms);
```

## ⚙️ 配置说明

### 数据库配置（config/config.yaml）

```yaml
database:
  type: postgres                 # 数据库类型
  host: localhost                # 数据库地址
  port: 5432                     # 数据库端口
  user: postgres                 # 用户名
  password: your_password        # 密码（请修改）
  dbname: paper_ai              # 数据库名
  max_idle_conns: 10            # 最大空闲连接数
  max_open_conns: 100           # 最大打开连接数
  conn_max_lifetime: 3600       # 连接最大生命周期（秒）
  auto_migrate: true            # 自动迁移表结构
  log_mode: info                # 日志级别
```

## 🚀 使用步骤

### 1. 安装PostgreSQL

**macOS:**
```bash
brew install postgresql@14
brew services start postgresql@14
```

**Ubuntu/Debian:**
```bash
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

### 2. 创建数据库

```bash
# 进入PostgreSQL
psql postgres

# 创建数据库
CREATE DATABASE paper_ai;

# 创建用户（可选）
CREATE USER paper_ai_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE paper_ai TO paper_ai_user;

# 退出
\q
```

### 3. 配置项目

修改 `config/config.yaml`：
```yaml
database:
  type: postgres
  host: localhost
  port: 5432
  user: postgres                # 或 paper_ai_user
  password: your_password       # 修改为实际密码
  dbname: paper_ai
  auto_migrate: true            # 首次运行自动创建表
```

### 4. 编译并运行

```bash
# 编译
go build -o paper_ai cmd/server/main.go

# 运行
./paper_ai
```

首次运行时，`auto_migrate: true` 会自动创建表结构。

## 📡 新增API接口

### 1. 查询记录列表

```bash
GET /api/v1/polish/records

# 参数：
# - page: 页码（默认1）
# - page_size: 每页大小（默认20，最大100）
# - provider: 按提供商过滤（可选）
# - status: 按状态过滤（success/failed，可选）
# - language: 按语言过滤（en/zh，可选）
# - style: 按风格过滤（academic/formal/concise，可选）
# - exclude_text: 是否排除大文本字段（true/false，可选）
# - start_time: 开始时间（RFC3339格式，可选）
# - end_time: 结束时间（RFC3339格式，可选）

# 示例：
curl "http://localhost:8080/api/v1/polish/records?page=1&page_size=20&provider=doubao&exclude_text=true"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "records": [
      {
        "id": 1,
        "trace_id": "550e8400-e29b-41d4-a716-446655440000",
        "style": "academic",
        "language": "zh",
        "original_length": 50,
        "polished_length": 60,
        "provider": "doubao",
        "model": "ep-m-20251124144251-5nxkx",
        "process_time_ms": 2500,
        "status": "success",
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  },
  "trace_id": "..."
}
```

### 2. 根据TraceID查询记录

```bash
GET /api/v1/polish/records/:trace_id

# 示例：
curl "http://localhost:8080/api/v1/polish/records/550e8400-e29b-41d4-a716-446655440000"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "trace_id": "550e8400-e29b-41d4-a716-446655440000",
    "original_content": "这是原始内容",
    "polished_content": "这是润色后的内容",
    "style": "academic",
    "language": "zh",
    "original_length": 50,
    "polished_length": 60,
    "provider": "doubao",
    "model": "ep-m-20251124144251-5nxkx",
    "process_time_ms": 2500,
    "status": "success",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
  },
  "trace_id": "..."
}
```

### 3. 获取统计信息

```bash
GET /api/v1/polish/statistics

# 参数：
# - start_time: 开始时间（RFC3339格式，可选）
# - end_time: 结束时间（RFC3339格式，可选）

# 示例：
curl "http://localhost:8080/api/v1/polish/statistics?start_time=2024-01-01T00:00:00Z&end_time=2024-12-31T23:59:59Z"
```

**响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total_count": 1000,
    "success_count": 950,
    "failed_count": 50,
    "success_rate": 95.0,
    "avg_process_time_ms": 2500.5,
    "provider_stats": {
      "doubao": {
        "count": 600,
        "success_count": 580,
        "failed_count": 20,
        "success_rate": 96.67,
        "avg_process_time_ms": 2400.0
      },
      "claude": {
        "count": 400,
        "success_count": 370,
        "failed_count": 30,
        "success_rate": 92.5,
        "avg_process_time_ms": 2650.0
      }
    },
    "language_stats": {
      "zh": {
        "count": 700,
        "success_count": 680,
        "success_rate": 97.14,
        "avg_process_time_ms": 2300.0
      },
      "en": {
        "count": 300,
        "success_count": 270,
        "success_rate": 90.0,
        "avg_process_time_ms": 2900.0
      }
    },
    "style_stats": {
      "academic": {
        "count": 500,
        "success_count": 480,
        "success_rate": 96.0,
        "avg_process_time_ms": 2600.0
      },
      "formal": {
        "count": 300,
        "success_count": 285,
        "success_rate": 95.0,
        "avg_process_time_ms": 2400.0
      },
      "concise": {
        "count": 200,
        "success_count": 185,
        "success_rate": 92.5,
        "avg_process_time_ms": 2500.0
      }
    }
  },
  "trace_id": "..."
}
```

## 🎯 架构优势

### 1. 高扩展性
- ✅ Repository接口可轻松切换实现（PostgreSQL → MySQL → MongoDB）
- ✅ Options模式支持灵活的查询条件组合
- ✅ 分层清晰，便于添加新功能

### 2. 低耦合
- ✅ 领域层不依赖任何外部框架（GORM、Gin等）
- ✅ Service层依赖Repository接口，不依赖具体实现
- ✅ PO和Entity分离，ORM标签不污染领域模型

### 3. 易测试
- ✅ 每层都可独立测试
- ✅ 可以Mock Repository接口进行单元测试

### 4. 性能优化
- ✅ 支持 `exclude_text` 参数，列表查询时排除大文本字段
- ✅ PostgreSQL的TEXT字段支持最大1GB
- ✅ 完善的索引设计，查询性能优异

## 🔍 查询优化建议

### 1. 列表查询优化
```bash
# 不需要查看内容时，排除大文本字段
curl "http://localhost:8080/api/v1/polish/records?page=1&page_size=20&exclude_text=true"
```

### 2. 时间范围查询
```bash
# 查询最近7天的记录
curl "http://localhost:8080/api/v1/polish/records?start_time=2024-01-20T00:00:00Z&end_time=2024-01-27T23:59:59Z"
```

### 3. 组合过滤
```bash
# 查询doubao提供商的成功记录
curl "http://localhost:8080/api/v1/polish/records?provider=doubao&status=success&exclude_text=true"
```

## 🛠️ 维护建议

### 1. 定期清理旧数据

```sql
-- 删除3个月前的失败记录
DELETE FROM polish_records
WHERE status = 'failed'
  AND created_at < NOW() - INTERVAL '3 months';

-- 归档6个月前的数据到历史表
INSERT INTO polish_records_archive
SELECT * FROM polish_records
WHERE created_at < NOW() - INTERVAL '6 months';
```

### 2. 监控数据库性能

```sql
-- 查看表大小
SELECT pg_size_pretty(pg_total_relation_size('polish_records'));

-- 查看索引使用情况
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
WHERE tablename = 'polish_records';
```

### 3. 优化建议

- 当数据量超过100万时，考虑分表（按月份或年份）
- 可以将历史数据归档到对象存储（OSS）
- 使用数据库连接池，调整 `max_open_conns` 参数

## ✨ 后续扩展

### 可扩展的功能
1. 添加用户认证，关联用户ID
2. 支持标签功能，便于分类管理
3. 添加收藏功能
4. 导出功能（导出为Excel、PDF等）
5. 数据可视化（图表展示统计信息）

### 可优化的点
1. 添加Redis缓存热点数据
2. 实现读写分离
3. 添加全文搜索（Elasticsearch）
4. 实现异步记录保存（消息队列）

## 🎉 总结

已成功实现：
- ✅ PostgreSQL持久化方案
- ✅ Clean Architecture分层架构
- ✅ Repository模式实现依赖倒置
- ✅ Options模式实现灵活查询
- ✅ 完整的CRUD和统计功能
- ✅ 性能优化（字段选择、索引优化）
- ✅ 自动表结构迁移
- ✅ 完善的错误处理

所有代码已编译通过，可以直接运行使用！
