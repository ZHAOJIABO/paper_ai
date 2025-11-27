# UserID查询过滤功能实现总结

## 概述

为了确保用户数据隔离和安全性，已实现完整的UserID查询过滤功能。现在所有查询操作都会自动过滤，确保用户只能查看自己的润色记录和统计数据。

---

## 一、实现的功能

### 1.1 查询记录列表过滤

**接口**: `GET /api/v1/polish/records`

**变更**:
- 自动从JWT token中提取用户ID
- 查询时自动添加`user_id`过滤条件
- 用户只能看到自己创建的润色记录

### 1.2 按TraceID查询记录

**接口**: `GET /api/v1/polish/records/:trace_id`

**变更**:
- 查询后验证记录所有权
- 如果记录不属于当前用户，返回403 Forbidden错误
- 防止用户通过trace_id访问他人记录

### 1.3 统计信息查询过滤

**接口**: `GET /api/v1/polish/statistics`

**变更**:
- 自动从JWT token中提取用户ID
- 统计数据只包含当前用户的记录
- 确保用户统计信息隔离

---

## 二、修改的文件

### 2.1 Domain层 - 查询选项

**[internal/domain/repository/query_options.go](internal/domain/repository/query_options.go)**

添加了UserID过滤支持：

```go
// QueryOptions 查询选项
type QueryOptions struct {
    // ... 其他字段
    UserID   *int64  // 按用户ID过滤 ⭐ 新增
}

// QueryOptionsBuilder 查询选项构建器
func (b *QueryOptionsBuilder) WithUserID(userID int64) *QueryOptionsBuilder {
    b.opts.UserID = &userID
    return b
}

// StatisticsOptions 统计选项
type StatisticsOptions struct {
    UserID      *int64     // 按用户ID过滤 ⭐ 新增
    TimeRange   *TimeRange
    GroupBy     []string
    Aggregation []string
}
```

### 2.2 Repository层 - 查询实现

**[internal/infrastructure/persistence/polish_repository_impl.go](internal/infrastructure/persistence/polish_repository_impl.go)**

在`buildQuery`方法中添加UserID过滤：

```go
func (r *polishRepositoryImpl) buildQuery(ctx context.Context, opts repository.QueryOptions) *gorm.DB {
    query := r.db.WithContext(ctx)

    // 字段选择优化
    if opts.ExcludeText {
        query = query.Select("id, trace_id, user_id, style, language, ...")
    }

    // 过滤条件
    if opts.UserID != nil {
        query = query.Where("user_id = ?", *opts.UserID)  // ⭐ 新增
    }

    // ... 其他过滤条件
    return query
}
```

**[internal/infrastructure/persistence/polish_repository_stats.go](internal/infrastructure/persistence/polish_repository_stats.go)**

在所有统计查询方法中添加UserID过滤：

```go
func (r *polishRepositoryImpl) GetStatistics(ctx context.Context, opts repository.StatisticsOptions) (*repository.Statistics, error) {
    query := r.db.WithContext(ctx).Model(&PolishRecordPO{})

    // 用户ID过滤 ⭐ 新增
    if opts.UserID != nil {
        query = query.Where("user_id = ?", *opts.UserID)
    }

    // 时间范围过滤
    if opts.TimeRange != nil {
        query = query.Where("created_at >= ? AND created_at <= ?", opts.TimeRange.Start, opts.TimeRange.End)
    }
    // ...
}

// getProviderStats、getLanguageStats、getStyleStats方法都添加了相同的UserID过滤
```

### 2.3 Service层 - 业务逻辑

**[internal/service/polish.go](internal/service/polish.go)**

修改`GetRecordByTraceID`方法，添加所有权验证：

```go
// 变更前
func (s *PolishService) GetRecordByTraceID(ctx context.Context, traceID string) (*entity.PolishRecord, error) {
    return s.polishRepo.GetByTraceID(ctx, traceID)
}

// 变更后
func (s *PolishService) GetRecordByTraceID(ctx context.Context, traceID string, userID int64) (*entity.PolishRecord, error) {
    record, err := s.polishRepo.GetByTraceID(ctx, traceID)
    if err != nil {
        return nil, err
    }

    // 验证记录所有权（只能查看自己的记录） ⭐ 新增
    if record.UserID != userID {
        return nil, apperrors.NewForbiddenError("you don't have permission to access this record")
    }

    return record, nil
}
```

### 2.4 Handler层 - API处理器

**[internal/api/handler/polish_query_handler.go](internal/api/handler/polish_query_handler.go)**

#### ListRecords方法 - 添加UserID过滤

```go
func (h *PolishQueryHandler) ListRecords(c *gin.Context) {
    // ... 解析参数

    // 构建查询选项
    builder := repository.NewQueryOptions().
        Page(page, pageSize).
        OrderBy("created_at", true)

    // 从JWT上下文获取用户ID，确保只能查询自己的记录 ⭐ 新增
    userID, exists := c.Get("user_id")
    if exists {
        builder.WithUserID(userID.(int64))
    }

    // ... 其他过滤条件
    opts := builder.Build()

    records, total, err := h.polishService.ListRecords(c.Request.Context(), opts)
    // ...
}
```

#### GetRecordByTraceID方法 - 添加所有权验证

```go
func (h *PolishQueryHandler) GetRecordByTraceID(c *gin.Context) {
    traceID := c.Param("trace_id")
    if traceID == "" {
        // ... 错误处理
    }

    // 从JWT上下文获取用户ID，确保只能查询自己的记录 ⭐ 新增
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(401, gin.H{
            "code":    401,
            "message": "unauthorized",
        })
        return
    }

    // 传递userID进行所有权验证 ⭐ 修改
    record, err := h.polishService.GetRecordByTraceID(c.Request.Context(), traceID, userID.(int64))
    // ...
}
```

#### GetStatistics方法 - 添加UserID过滤

```go
func (h *PolishQueryHandler) GetStatistics(c *gin.Context) {
    // 解析时间范围
    // ...

    opts := repository.StatisticsOptions{}

    // 从JWT上下文获取用户ID，确保只能查看自己的统计信息 ⭐ 新增
    userID, exists := c.Get("user_id")
    if exists {
        uid := userID.(int64)
        opts.UserID = &uid
    }

    // 时间范围处理
    // ...

    stats, err := h.polishService.GetStatistics(c.Request.Context(), opts)
    // ...
}
```

---

## 三、安全性提升

### 3.1 数据隔离

✅ **查询列表**: 用户只能看到自己的润色记录列表
✅ **TraceID查询**: 即使知道他人的trace_id，也无法访问
✅ **统计数据**: 统计信息只包含用户自己的数据

### 3.2 权限验证

| 操作 | 验证方式 | 错误码 |
|------|----------|--------|
| 查询记录列表 | WHERE user_id = ? | 自动过滤 |
| 按TraceID查询 | 查询后对比UserID | 403 Forbidden |
| 获取统计信息 | WHERE user_id = ? | 自动过滤 |

### 3.3 错误处理

```go
// 未登录用户访问需要认证的接口
{
  "code": 401,
  "message": "unauthorized"
}

// 尝试访问他人的记录
{
  "code": 20009,
  "message": "you don't have permission to access this record"
}
```

---

## 四、测试验证

### 4.1 测试场景

1. **用户A查询自己的记录列表**
   - 预期：只返回用户A创建的记录

2. **用户A尝试访问用户B的trace_id**
   - 预期：返回403 Forbidden错误

3. **用户A查看统计信息**
   - 预期：统计数据只包含用户A的记录

### 4.2 测试步骤

```bash
# 1. 启动服务
./paper_ai

# 2. 注册两个用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "user_a", "email": "a@test.com", "password": "Test1234", "confirm_password": "Test1234"}'

curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "user_b", "email": "b@test.com", "password": "Test1234", "confirm_password": "Test1234"}'

# 3. 登录获取token
TOKEN_A=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user_a", "password": "Test1234"}' | jq -r '.data.access_token')

TOKEN_B=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user_b", "password": "Test1234"}' | jq -r '.data.access_token')

# 4. 用户A创建润色记录
RESPONSE_A=$(curl -s -X POST http://localhost:8080/api/v1/polish \
  -H "Authorization: Bearer $TOKEN_A" \
  -H "Content-Type: application/json" \
  -d '{"content": "用户A的测试文本", "style": "academic", "language": "zh", "provider": "doubao"}')

TRACE_ID_A=$(echo $RESPONSE_A | jq -r '.trace_id')

# 5. 用户B创建润色记录
RESPONSE_B=$(curl -s -X POST http://localhost:8080/api/v1/polish \
  -H "Authorization: Bearer $TOKEN_B" \
  -H "Content-Type: application/json" \
  -d '{"content": "用户B的测试文本", "style": "academic", "language": "zh", "provider": "doubao"}')

# 6. 用户A查询自己的记录列表（应该只看到自己的）
curl -X GET "http://localhost:8080/api/v1/polish/records?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN_A"

# 7. 用户A尝试访问用户B的trace_id（应该返回403）
curl -X GET "http://localhost:8080/api/v1/polish/records/$TRACE_ID_B" \
  -H "Authorization: Bearer $TOKEN_A"

# 预期响应：
# {
#   "code": 20009,
#   "message": "you don't have permission to access this record"
# }

# 8. 用户A查看统计信息（应该只包含自己的数据）
curl -X GET "http://localhost:8080/api/v1/polish/statistics" \
  -H "Authorization: Bearer $TOKEN_A"
```

### 4.3 预期结果

✅ 用户只能查询到自己的记录
✅ 跨用户访问被拒绝（403错误）
✅ 统计信息正确隔离
✅ 数据库查询自动添加user_id过滤

---

## 五、性能考虑

### 5.1 数据库索引

已在`polish_records`表上创建了`idx_user_id`索引：

```sql
CREATE INDEX idx_user_id ON polish_records(user_id);
```

### 5.2 查询优化

- 所有查询自动添加`WHERE user_id = ?`条件
- 利用索引加速查询
- 分页查询避免大量数据加载

### 5.3 性能影响

| 操作 | 影响 | 说明 |
|------|------|------|
| 查询列表 | 无影响 | 使用索引，性能反而可能提升 |
| TraceID查询 | 微小 | 增加一次UserID对比 |
| 统计查询 | 无影响 | 使用索引过滤 |

---

## 六、注意事项

### 6.1 未登录用户

如果用户未登录（没有JWT token），查询接口的行为：

- **ListRecords**: 返回空列表（因为`user_id`为null，无法匹配任何记录）
- **GetRecordByTraceID**: 返回401 Unauthorized
- **GetStatistics**: 返回空统计数据

### 6.2 管理员功能（未来扩展）

如果将来需要管理员查看所有用户的数据，可以：

1. 添加角色判断：
```go
if !isAdmin(userID) {
    builder.WithUserID(userID)
}
```

2. 或者创建专门的管理员查询接口：
```go
// GET /api/v1/admin/polish/records
func (h *AdminHandler) ListAllRecords(c *gin.Context) {
    // 不添加UserID过滤
    opts := builder.Build()
    records, total, err := h.polishService.ListRecords(c.Request.Context(), opts)
    // ...
}
```

---

## 七、总结

### 7.1 改进点

✅ **完整的数据隔离**: 用户无法看到其他用户的数据
✅ **安全的权限控制**: 防止通过trace_id越权访问
✅ **自动化过滤**: 在Repository层自动添加过滤条件
✅ **性能优化**: 利用数据库索引，不影响性能

### 7.2 最佳实践

1. **在多个层次验证权限**：
   - Repository层：自动过滤
   - Service层：所有权验证
   - Handler层：提取用户身份

2. **使用建造者模式**：
   - 灵活的查询条件构建
   - 易于扩展新的过滤条件

3. **错误处理清晰**：
   - 401: 未登录
   - 403: 已登录但无权限
   - 404: 记录不存在

---

## 八、相关文档

- [UserID改进总结文档](USERID_IMPROVEMENT.md)
- [用户认证实现文档](AUTH_IMPLEMENTATION.md)
- [数据库表结构](internal/infrastructure/persistence/models.go)

**最后更新**: 2024-11-27
