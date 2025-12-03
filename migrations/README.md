# 数据库迁移说明

本项目使用 [golang-migrate](https://github.com/golang-migrate/migrate) 进行数据库版本管理。

## 快速开始

### 1. 设置数据库连接

在环境变量或命令行中设置 `DATABASE_URL`：

```bash
export DATABASE_URL="postgres://user:password@localhost:5432/paper_ai?sslmode=disable"
```

或在 Makefile 中修改默认值（第34行）。

### 2. 执行迁移

```bash
# 执行所有待执行的迁移
make migrate-up

# 回滚最后一次迁移
make migrate-down

# 查看当前版本
make migrate-version
```

## 迁移文件说明

迁移文件采用版本化命名格式：`{version}_{description}.{up|down}.sql`

### 当前迁移文件

1. **000000_initial_tables.sql** - 初始表结构
   - 创建 `users` 表
   - 创建 `polish_records` 表
   - 创建 `refresh_tokens` 表
   - 添加索引和触发器

2. **000001_initial_schema.sql** - 多版本润色功能
   - 扩展 `polish_records` 表（添加 `mode` 字段）
   - 创建 `polish_versions` 表（润色版本详情）
   - 创建 `polish_prompts` 表（Prompt 管理）
   - 扩展 `users` 表（多版本功能权限）
   - 插入初始 Prompt 数据

## 常用命令

### 查看帮助

```bash
make help
```

### 创建新迁移

```bash
# 创建新的迁移文件
make migrate-create name=add_new_feature

# 会生成两个文件：
# migrations/000002_add_new_feature.up.sql
# migrations/000002_add_new_feature.down.sql
```

### 迁移管理

```bash
# 执行所有待执行的迁移
make migrate-up

# 回滚最后一次迁移
make migrate-down

# 查看当前迁移版本
make migrate-version

# 强制设置迁移版本（修复损坏状态）
make migrate-force version=1

# 删除所有表（危险！）
make migrate-drop
```

## 应用启动时的自动迁移

应用在启动时会自动执行迁移（如果配置中 `database.auto_migrate: true`）。

代码位置：[internal/infrastructure/database/database.go:85-91](../internal/infrastructure/database/database.go#L85-L91)

## 迁移最佳实践

### 1. 编写迁移文件

**UP 文件（向前迁移）：**
```sql
-- 添加新表
CREATE TABLE new_table (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- 添加索引
CREATE INDEX idx_name ON new_table(name);
```

**DOWN 文件（回滚）：**
```sql
-- 删除表
DROP TABLE IF EXISTS new_table CASCADE;
```

### 2. 注意事项

- ✅ **DO**: 使用 `IF EXISTS` 和 `IF NOT EXISTS` 确保幂等性
- ✅ **DO**: 在 DOWN 文件中完全撤销 UP 文件的操作
- ✅ **DO**: 先在开发环境测试迁移
- ✅ **DO**: 备份生产数据库后再执行迁移
- ❌ **DON'T**: 修改已经执行过的迁移文件
- ❌ **DON'T**: 在生产环境直接使用 `migrate-drop`

### 3. 添加新字段（示例）

**UP:**
```sql
ALTER TABLE users
ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

CREATE INDEX IF NOT EXISTS idx_phone ON users(phone);

COMMENT ON COLUMN users.phone IS '用户手机号';
```

**DOWN:**
```sql
DROP INDEX IF EXISTS idx_phone;
ALTER TABLE users DROP COLUMN IF EXISTS phone;
```

### 4. 数据迁移（示例）

**UP:**
```sql
-- 1. 添加新列
ALTER TABLE users ADD COLUMN full_name VARCHAR(100);

-- 2. 迁移数据
UPDATE users SET full_name = CONCAT(first_name, ' ', last_name);

-- 3. 设置为 NOT NULL（在数据填充后）
ALTER TABLE users ALTER COLUMN full_name SET NOT NULL;
```

**DOWN:**
```sql
ALTER TABLE users DROP COLUMN IF EXISTS full_name;
```

## 故障排查

### 迁移失败（Dirty State）

如果迁移执行到一半失败，数据库会进入 "dirty" 状态：

```bash
# 1. 检查当前版本和状态
make migrate-version

# 2. 手动修复数据库到正确状态

# 3. 强制设置版本
make migrate-force version=1
```

### 查看迁移历史

```bash
# 连接数据库查看迁移表
psql $DATABASE_URL -c "SELECT * FROM schema_migrations;"
```

### 重置数据库

```bash
# 删除所有表并重新迁移
make migrate-drop
make migrate-up
```

## 参考资料

- [golang-migrate 官方文档](https://github.com/golang-migrate/migrate)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)
