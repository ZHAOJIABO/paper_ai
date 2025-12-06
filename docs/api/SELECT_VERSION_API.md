# 版本选择接口文档

## 接口说明

当用户完成多版本润色并查看3个版本（conservative、balanced、aggressive）后，需要选择其中一个版本作为后续对比和修改的基础。此接口用于将用户选择的版本内容复制到主记录中。

**重要**：这是多版本润色工作流程中的**必需步骤**。只有调用此接口后，后续的对比和修改操作才能正确更新 `final_content`。

---

## 接口详情

### 基本信息

- **接口**: `POST /api/v1/polish/select-version/:trace_id`
- **认证**: 需要 JWT Token
- **功能**: 选择多版本润色中的一个版本

### 请求头

```
Authorization: Bearer {access_token}
Content-Type: application/json
```

### 路径参数

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| trace_id | string | 是 | 多版本润色返回的追踪ID |

### 查询参数

| 参数名 | 类型 | 必填 | 说明 | 可选值 |
|--------|------|------|------|--------|
| version | string | 是 | 要选择的版本类型 | `conservative` / `balanced` / `aggressive` |

---

## 请求示例

### cURL

```bash
curl -X POST 'http://localhost:8080/api/v1/polish/select-version/1764839051100?version=balanced' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json'
```

### JavaScript (Fetch)

```javascript
const traceId = '1764839051100';
const versionType = 'balanced';

fetch(`http://localhost:8080/api/v1/polish/select-version/${traceId}?version=${versionType}`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
    'Content-Type': 'application/json'
  }
})
  .then(response => response.json())
  .then(data => {
    if (data.code === 0) {
      console.log('版本选择成功:', data.data);
    } else {
      console.error('版本选择失败:', data.message);
    }
  });
```

### TypeScript (Axios)

```typescript
import axios from 'axios';

interface SelectVersionResponse {
  message: string;
  trace_id: string;
  selected_version: string;
}

const selectVersion = async (traceId: string, versionType: 'conservative' | 'balanced' | 'aggressive') => {
  try {
    const response = await axios.post<{
      code: number;
      message: string;
      data: SelectVersionResponse;
      trace_id: string;
    }>(
      `/api/v1/polish/select-version/${traceId}`,
      null,
      {
        params: { version: versionType },
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`
        }
      }
    );

    if (response.data.code === 0) {
      console.log('版本选择成功:', response.data.data);
      return response.data.data;
    } else {
      throw new Error(response.data.message);
    }
  } catch (error) {
    console.error('版本选择失败:', error);
    throw error;
  }
};

// 使用示例
selectVersion('1764839051100', 'balanced');
```

---

## 响应数据

### 成功响应

**HTTP Status Code**: 200

**响应结构**:

```typescript
interface ApiResponse {
  code: number;              // 错误码，0表示成功
  message: string;           // 响应消息
  data: {
    message: string;         // 操作结果消息
    trace_id: string;        // 追踪ID
    selected_version: string; // 已选择的版本类型
  };
  trace_id: string;          // 请求追踪ID
}
```

**响应示例**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "版本选择成功",
    "trace_id": "1764839051100",
    "selected_version": "balanced"
  },
  "trace_id": "abc123"
}
```

### 错误响应

#### 1. 参数错误 (400)

```json
{
  "code": 10001,
  "message": "version 参数不能为空",
  "data": null,
  "trace_id": "xyz789"
}
```

```json
{
  "code": 10001,
  "message": "无效的版本类型: invalid_type",
  "data": null,
  "trace_id": "xyz789"
}
```

#### 2. 未授权 (401)

```json
{
  "code": 20008,
  "message": "请先登录",
  "data": null,
  "trace_id": "xyz789"
}
```

#### 3. 无权访问 (403)

```json
{
  "code": 20009,
  "message": "无权访问该记录",
  "data": null,
  "trace_id": "xyz789"
}
```

#### 4. 记录不存在 (404)

```json
{
  "code": 20003,
  "message": "润色记录不存在",
  "data": null,
  "trace_id": "xyz789"
}
```

```json
{
  "code": 20003,
  "message": "版本 balanced 不存在",
  "data": null,
  "trace_id": "xyz789"
}
```

#### 5. 业务逻辑错误 (400)

```json
{
  "code": 10001,
  "message": "该记录不是多版本润色",
  "data": null,
  "trace_id": "xyz789"
}
```

```json
{
  "code": 10001,
  "message": "版本 balanced 生成失败: AI call timeout",
  "data": null,
  "trace_id": "xyz789"
}
```

---

## 前端实现示例

### React + TypeScript 完整示例

```typescript
// src/api/polish.ts
import request from '@/utils/request';
import type { ApiResponse } from '@/types';

export interface SelectVersionResponse {
  message: string;
  trace_id: string;
  selected_version: string;
}

export const polishApi = {
  // 选择版本
  selectVersion: (traceId: string, version: 'conservative' | 'balanced' | 'aggressive') => {
    return request.post<ApiResponse<SelectVersionResponse>>(
      `/api/v1/polish/select-version/${traceId}`,
      null,
      { params: { version } }
    );
  },
};
```

```typescript
// src/pages/MultiVersionPolish.tsx
import React, { useState } from 'react';
import { Button, Card, message, Spin } from 'antd';
import { useNavigate } from 'react-router-dom';
import { polishApi } from '@/api/polish';

interface VersionData {
  polished_content: string;
  polished_length: number;
  suggestions: string[];
  status: string;
  error_message?: string;
}

const MultiVersionPolish: React.FC = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [traceId, setTraceId] = useState('');
  const [versions, setVersions] = useState<Record<string, VersionData> | null>(null);
  const [originalLength, setOriginalLength] = useState(0);
  const [selectingVersion, setSelectingVersion] = useState<string | null>(null);

  // 调用多版本润色
  const handlePolish = async (content: string) => {
    setLoading(true);
    try {
      const response = await polishApi.multiVersionPolish({
        content,
        style: 'academic',
        language: 'en',
        provider: 'claude'
      });

      setTraceId(response.data.trace_id);
      setVersions(response.data.versions);
      setOriginalLength(response.data.original_length);
      message.success('多版本润色完成，请选择一个版本');
    } catch (error: any) {
      message.error(error.message || '润色失败');
    } finally {
      setLoading(false);
    }
  };

  // 选择版本
  const handleSelectVersion = async (versionType: string) => {
    const versionData = versions![versionType];

    // 1. 检查版本状态
    if (versionData.status !== 'success') {
      message.error(`该版本生成失败：${versionData.error_message || '未知错误'}`);
      return;
    }

    setSelectingVersion(versionType);

    try {
      // 2. 调用版本选择接口
      await polishApi.selectVersion(traceId, versionType as any);

      message.success(`已选择${getVersionName(versionType)}，正在跳转到对比页面...`);

      // 3. 跳转到对比页面（不需要传 version 参数，因为已经选择了）
      setTimeout(() => {
        navigate(`/comparison/${traceId}`);
      }, 500);
    } catch (error: any) {
      message.error(error.message || '选择版本失败');
    } finally {
      setSelectingVersion(null);
    }
  };

  const getVersionName = (type: string) => {
    const names: Record<string, string> = {
      conservative: '保守版本',
      balanced: '均衡版本',
      aggressive: '激进版本'
    };
    return names[type] || type;
  };

  const getVersionDescription = (type: string) => {
    const descriptions: Record<string, string> = {
      conservative: '轻微润色，保持原意',
      balanced: '适度优化，提升表达',
      aggressive: '大幅改写，提升质量'
    };
    return descriptions[type] || '';
  };

  return (
    <div style={{ maxWidth: 1200, margin: '20px auto', padding: '0 20px' }}>
      <h1>多版本论文润色</h1>

      {/* 润色表单 - 省略... */}

      {/* 版本选择卡片 */}
      {versions && (
        <div style={{ marginTop: 40 }}>
          <h2>请选择一个版本进行后续编辑</h2>
          <p style={{ color: '#666', marginBottom: 20 }}>
            系统已生成3个不同风格的版本，请选择最符合您需求的版本。选择后可以进行详细的对比和修改。
          </p>

          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(350px, 1fr))',
            gap: 16
          }}>
            {Object.entries(versions).map(([versionType, versionData]) => (
              <Card
                key={versionType}
                title={
                  <div>
                    <div style={{ fontSize: 18, fontWeight: 'bold' }}>
                      {getVersionName(versionType)}
                    </div>
                    <div style={{ fontSize: 14, fontWeight: 'normal', color: '#666' }}>
                      {getVersionDescription(versionType)}
                    </div>
                  </div>
                }
                style={{
                  border: versionData.status === 'success' ? '1px solid #d9d9d9' : '1px solid #ff4d4f',
                  height: '100%'
                }}
              >
                {/* 状态指示 */}
                {versionData.status === 'success' ? (
                  <div style={{ color: '#52c41a', marginBottom: 12 }}>
                    ✓ 生成成功
                  </div>
                ) : (
                  <div style={{ color: '#ff4d4f', marginBottom: 12 }}>
                    ✗ 生成失败: {versionData.error_message}
                  </div>
                )}

                {/* 字数变化 */}
                {versionData.status === 'success' && (
                  <>
                    <div style={{ marginBottom: 12 }}>
                      <strong>字数变化：</strong>
                      {versionData.polished_length} 字
                      <span style={{
                        color: versionData.polished_length > originalLength ? '#52c41a' : '#1890ff',
                        marginLeft: 8
                      }}>
                        ({versionData.polished_length - originalLength > 0 ? '+' : ''}
                        {versionData.polished_length - originalLength})
                      </span>
                    </div>

                    {/* 内容预览 */}
                    <div style={{ marginBottom: 12 }}>
                      <strong>内容预览：</strong>
                      <div style={{
                        padding: 12,
                        background: '#f5f5f5',
                        borderRadius: 4,
                        marginTop: 8,
                        maxHeight: 150,
                        overflow: 'hidden',
                        position: 'relative'
                      }}>
                        <div style={{
                          lineHeight: 1.6,
                          fontSize: 14,
                          color: '#333'
                        }}>
                          {versionData.polished_content.substring(0, 200)}
                          {versionData.polished_content.length > 200 && '...'}
                        </div>
                      </div>
                    </div>

                    {/* 改进建议 */}
                    {versionData.suggestions && versionData.suggestions.length > 0 && (
                      <div style={{ marginBottom: 12 }}>
                        <strong>改进建议：</strong>
                        <ul style={{ marginTop: 8, paddingLeft: 20 }}>
                          {versionData.suggestions.slice(0, 3).map((suggestion, idx) => (
                            <li key={idx} style={{ fontSize: 14, color: '#666' }}>
                              {suggestion}
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}

                    {/* 选择按钮 */}
                    <Button
                      type="primary"
                      size="large"
                      block
                      onClick={() => handleSelectVersion(versionType)}
                      loading={selectingVersion === versionType}
                      disabled={selectingVersion !== null && selectingVersion !== versionType}
                      style={{ marginTop: 12 }}
                    >
                      {selectingVersion === versionType ? '正在选择...' : '选择此版本'}
                    </Button>
                  </>
                )}
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default MultiVersionPolish;
```

---

## 工作流程集成

### 完整的多版本润色流程

```typescript
/**
 * 多版本润色完整流程
 */

// Step 1: 调用多版本润色
const polishResponse = await polishApi.multiVersionPolish({
  content: "原文内容",
  style: "academic",
  language: "en",
  provider: "claude"
});

const traceId = polishResponse.data.trace_id;
const versions = polishResponse.data.versions;

// Step 2: 用户查看3个版本并选择一个
// 前端展示 conservative、balanced、aggressive 三个版本供用户选择

// Step 3: 用户点击选择某个版本（例如 balanced）
await polishApi.selectVersion(traceId, 'balanced');
// ⚠️ 重要：此步骤将 balanced 版本的内容复制到主记录

// Step 4: 跳转到对比页面（不需要传 version 参数）
navigate(`/comparison/${traceId}`);
// 因为已经选择了版本，主记录的 polished_content 已经是选中版本的内容

// Step 5: 用户在对比页面接受/拒绝修改
await comparisonApi.applyAction(traceId, changeId, 'accept');
// 修改会更新主记录的 final_content
```

---

## 重要说明

### 1. 必须先选择版本

在多版本润色流程中，**必须**先调用 `/polish/select-version/:trace_id` 接口选择版本，然后才能进行对比和修改操作。

**错误流程** ❌:
```
多版本润色 → 直接查看对比（传 version 参数） → 应用修改
```
这样会导致 `final_content` 没有正确更新。

**正确流程** ✅:
```
多版本润色 → 选择版本 → 查看对比 → 应用修改
```

### 2. 数据更新机制

选择版本时，系统会执行以下操作：

1. 将选中版本的 `polished_content` 复制到主记录的 `polished_content`
2. 将选中版本的内容同时复制到 `final_content`（初始值）
3. 记录 `selected_version` 字段为选中的版本类型
4. 更新相关元数据（长度、模型等）

### 3. 版本选择是一次性操作

一旦选择了某个版本，就无法"撤销"回到未选择状态。如果需要切换版本：

**方案 A**: 重新生成多版本润色
```typescript
// 重新调用多版本润色接口，生成新的 trace_id
const newResponse = await polishApi.multiVersionPolish({ ... });
```

**方案 B**: 再次选择其他版本
```typescript
// 可以再次调用选择接口，选择其他版本
await polishApi.selectVersion(traceId, 'aggressive');
// 这会覆盖之前的选择
```

### 4. 错误处理建议

```typescript
const handleSelectVersion = async (traceId: string, versionType: string) => {
  try {
    // 选择版本前先检查版本状态
    const versionData = versions[versionType];

    if (versionData.status !== 'success') {
      message.error(`该版本生成失败，无法选择`);
      return;
    }

    // 调用选择接口
    await polishApi.selectVersion(traceId, versionType);
    message.success('版本选择成功');

    // 跳转到对比页面
    navigate(`/comparison/${traceId}`);

  } catch (error: any) {
    // 根据错误码进行不同处理
    if (error.code === 20003) {
      message.error('记录不存在，请重新进行润色');
    } else if (error.code === 10001) {
      message.error('参数错误，请联系技术支持');
    } else if (error.code === 20009) {
      message.error('无权访问该记录');
    } else {
      message.error('选择版本失败，请重试');
    }
  }
};
```

---

## 常见问题

### Q1: 是否可以在对比页面直接传 `version` 参数而不调用选择接口？

**A**: 不推荐这样做。虽然对比接口支持 `version` 参数用于临时查看某个版本的对比，但这不会更新主记录。如果直接在对比页面应用修改，会更新错误的记录。

### Q2: 选择版本后，原来的3个版本数据还存在吗？

**A**: 是的。选择版本只是将选中版本的内容复制到主记录，`polish_versions` 表中的3个版本记录依然存在。如果需要，用户可以重新选择其他版本。

### Q3: 如果用户关闭页面没有选择版本，会怎样？

**A**: 主记录的 `polished_content` 会保持默认值（通常是第一个成功生成的版本）。用户可以随时回来选择版本。建议在前端添加提示，引导用户完成版本选择。

### Q4: 选择版本是否需要很长时间？

**A**: 不需要。选择版本只是数据库的更新操作，通常在几十毫秒内完成。不涉及AI调用。

---

## 测试清单

在集成此接口前，请确保测试以下场景：

- [ ] 正常流程：多版本润色 → 选择版本 → 对比 → 修改
- [ ] 选择成功状态的版本
- [ ] 尝试选择失败状态的版本（应返回错误）
- [ ] 使用无效的 version 参数（应返回400错误）
- [ ] 使用不存在的 trace_id（应返回404错误）
- [ ] 使用其他用户的 trace_id（应返回403错误）
- [ ] 重复选择同一版本（应成功）
- [ ] 切换选择不同版本（应成功）
- [ ] Token过期场景（应返回401错误）
- [ ] 网络错误处理

---

## 更新日志

- **2024-12-05**: 初始版本，新增版本选择接口文档
