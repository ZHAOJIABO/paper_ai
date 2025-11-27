# Paper AI 前端接口集成文档

## 概述

本文档提供Paper AI后端API的完整集成指南，包括用户认证、论文润色等功能的接口说明和前端实现示例。

## 基础配置

### API基础信息

- **Base URL (开发环境)**: `http://localhost:8080`
- **API版本**: `v1`
- **认证方式**: JWT Bearer Token
- **请求格式**: `application/json`
- **响应格式**: `application/json`

### 统一响应格式

所有接口都遵循统一的响应格式：

```typescript
interface ApiResponse<T = any> {
  code: number;        // 错误码，0表示成功
  message: string;     // 响应消息
  data: T;            // 响应数据
  trace_id: string;   // 请求追踪ID
}
```

**成功响应示例**：
```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "trace_id": "abc123"
}
```

**错误响应示例**：
```json
{
  "code": 20002,
  "message": "密码错误",
  "data": null,
  "trace_id": "xyz789"
}
```

## 错误码对照表

| 错误码 | 说明 | 前端处理建议 |
|-------|------|-------------|
| 0 | 成功 | - |
| 10001 | 参数错误 | 提示用户检查输入 |
| 20001 | 用户已存在 | 提示用户换用户名或邮箱 |
| 20002 | 密码错误 | 提示密码错误 |
| 20003 | 用户不存在 | 提示用户不存在 |
| 20004 | Token无效 | 清除本地token，跳转登录页 |
| 20005 | Token过期 | 尝试刷新token |
| 20006 | 密码强度不够 | 提示密码要求 |
| 20007 | 账号已被封禁 | 提示联系管理员 |
| 20008 | 未授权 | 跳转登录页 |
| 20009 | 禁止访问 | 提示权限不足 |

---

## 一、用户认证接口

### 1.1 用户注册

**接口**: `POST /api/v1/auth/register`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```typescript
interface RegisterRequest {
  username: string;         // 用户名，3-50位，只能包含字母、数字、下划线
  email: string;           // 邮箱地址
  password: string;        // 密码，至少8位，包含字母和数字
  confirm_password: string; // 确认密码，必须与password一致
  nickname?: string;       // 昵称（可选），最多50位
}
```

**请求示例**:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "Password123",
  "confirm_password": "Password123",
  "nickname": "John"
}
```

**响应数据**:
```typescript
interface UserInfo {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar_url: string;
  status: "active" | "inactive" | "banned";
  email_verified: boolean;
  last_login_at?: string;  // ISO 8601格式
  created_at: string;      // ISO 8601格式
}
```

**响应示例**:
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
  "trace_id": "abc123"
}
```

**前端校验规则**:
- username: 3-50位，只能包含字母、数字、下划线，正则：`^[a-zA-Z0-9_]{3,50}$`
- email: 标准邮箱格式
- password: 至少8位，必须包含字母和数字
- confirm_password: 必须与password完全一致
- nickname: 最多50位

---

### 1.2 用户登录

**接口**: `POST /api/v1/auth/login`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```typescript
interface LoginRequest {
  username: string;  // 用户名或邮箱
  password: string;  // 密码
}
```

**请求示例**:
```json
{
  "username": "john_doe",
  "password": "Password123"
}
```

**响应数据**:
```typescript
interface LoginResponse {
  access_token: string;   // 访问令牌
  refresh_token: string;  // 刷新令牌
  expires_in: number;     // access_token过期时间（秒）
  token_type: string;     // 令牌类型，固定为"Bearer"
  user: UserInfo;         // 用户信息
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
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
      "last_login_at": "2024-01-01T10:00:00Z",
      "created_at": "2024-01-01T09:00:00Z"
    }
  },
  "trace_id": "abc123"
}
```

**重要提示**:
- 登录成功后，需要将 `access_token` 和 `refresh_token` 保存到本地存储
- 后续需要认证的接口请求都需要在请求头中携带 `Authorization: Bearer {access_token}`
- `access_token` 有效期2小时，过期后需要使用 `refresh_token` 刷新

---

### 1.3 获取当前用户信息

**接口**: `GET /api/v1/auth/me`

**请求头**:
```
Authorization: Bearer {access_token}
```

**无请求参数**

**响应数据**:
```typescript
interface UserInfo {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar_url: string;
  status: "active" | "inactive" | "banned";
  email_verified: boolean;
  last_login_at: string;
  created_at: string;
}
```

**响应示例**:
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
  "trace_id": "abc123"
}
```

**使用场景**:
- 页面刷新后恢复用户登录状态
- 验证token是否有效
- 获取最新的用户信息

---

### 1.4 刷新访问令牌

**接口**: `POST /api/v1/auth/refresh`

**请求头**:
```
Content-Type: application/json
```

**请求参数**:
```typescript
interface RefreshTokenRequest {
  refresh_token: string;  // 刷新令牌
}
```

**请求示例**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应数据**:
```typescript
interface RefreshTokenResponse {
  access_token: string;  // 新的访问令牌
  expires_in: number;    // 过期时间（秒）
  token_type: string;    // 令牌类型
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 7200,
    "token_type": "Bearer"
  },
  "trace_id": "abc123"
}
```

**使用场景**:
- 当 `access_token` 过期（收到401错误码20005）时自动调用
- 建议在请求拦截器中实现自动刷新逻辑

---

### 1.5 用户登出

**接口**: `POST /api/v1/auth/logout`

**请求头**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**请求参数**:
```typescript
interface LogoutRequest {
  refresh_token: string;  // 刷新令牌
}
```

**请求示例**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应数据**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "登出成功"
  },
  "trace_id": "abc123"
}
```

**前端处理流程**:
1. 调用登出接口
2. 清除本地存储的 `access_token` 和 `refresh_token`
3. 清除用户状态
4. 跳转到登录页

---

## 二、论文润色接口

### 2.1 段落润色

**接口**: `POST /api/v1/polish`

**请求头**:
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**请求参数**:
```typescript
interface PolishRequest {
  content: string;     // 需要润色的文本内容，最多10000字符
  style?: string;      // 润色风格：academic(学术)、formal(正式)、concise(简洁)，默认academic
  language?: string;   // 目标语言：en(英文)、zh(中文)，默认en
  provider?: string;   // AI提供商：claude、doubao，可选，不指定则使用默认
}
```

**请求示例**:
```json
{
  "content": "这篇文章讨论了机器学习在软件开发中的作用。",
  "style": "academic",
  "language": "zh",
  "provider": "doubao"
}
```

**响应数据**:
```typescript
interface PolishResponse {
  polished_content: string;   // 润色后的文本
  original_length: number;    // 原始文本长度
  polished_length: number;    // 润色后文本长度
  suggestions: string[];      // 改进建议
  provider_used: string;      // 实际使用的AI提供商
  model_used: string;         // 实际使用的模型
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "polished_content": "本文深入探讨了机器学习技术在现代软件开发流程中的关键作用与实际应用。",
    "original_length": 25,
    "polished_length": 38,
    "suggestions": [
      "语言更加学术化",
      "增加了专业术语",
      "句式更加正式"
    ],
    "provider_used": "doubao",
    "model_used": "ep-m-20251124144251-5nxkx"
  },
  "trace_id": "abc123"
}
```

---

### 2.2 查询润色记录列表

**接口**: `GET /api/v1/polish/records`

**请求头**:
```
Authorization: Bearer {access_token}
```

**查询参数**:
```typescript
interface RecordsQueryParams {
  page?: number;          // 页码，默认1
  page_size?: number;     // 每页大小，默认20，最大100
  provider?: string;      // 按提供商过滤
  status?: string;        // 按状态过滤：success、failed
  language?: string;      // 按语言过滤：en、zh
  style?: string;         // 按风格过滤：academic、formal、concise
  exclude_text?: boolean; // 是否排除大文本字段，默认false
  start_time?: string;    // 开始时间，RFC3339格式
  end_time?: string;      // 结束时间，RFC3339格式
}
```

**请求示例**:
```
GET /api/v1/polish/records?page=1&page_size=20&language=zh&status=success
```

**响应数据**:
```typescript
interface RecordsResponse {
  records: PolishRecord[];
  total: number;
  page: number;
  page_size: number;
}

interface PolishRecord {
  id: number;
  user_id: number;
  trace_id: string;
  original_content?: string;    // exclude_text=true时为空
  polished_content?: string;    // exclude_text=true时为空
  style: string;
  language: string;
  provider: string;
  model: string;
  status: string;
  error_message?: string;
  process_time_ms: number;
  created_at: string;
}
```

---

### 2.3 根据TraceID查询记录

**接口**: `GET /api/v1/polish/records/{trace_id}`

**请求头**:
```
Authorization: Bearer {access_token}
```

**路径参数**:
- `trace_id`: 请求追踪ID

**请求示例**:
```
GET /api/v1/polish/records/abc123
```

**响应数据**: 同PolishRecord

---

### 2.4 获取统计信息

**接口**: `GET /api/v1/polish/statistics`

**请求头**:
```
Authorization: Bearer {access_token}
```

**查询参数**:
```typescript
interface StatisticsQueryParams {
  start_time?: string;  // 开始时间，RFC3339格式
  end_time?: string;    // 结束时间，RFC3339格式
}
```

**响应数据**:
```typescript
interface StatisticsResponse {
  total_count: number;           // 总请求数
  success_count: number;         // 成功数
  failed_count: number;          // 失败数
  success_rate: number;          // 成功率
  avg_process_time_ms: number;   // 平均处理时间（毫秒）
  provider_stats: Record<string, number>;  // 各提供商使用次数
  language_stats: Record<string, number>;  // 各语言使用次数
  style_stats: Record<string, number>;     // 各风格使用次数
}
```

---

## 三、前端实现示例

### 3.1 axios封装（推荐）

#### 创建axios实例

```typescript
// src/utils/request.ts
import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { message } from 'antd'; // 或其他UI库的提示组件

// 创建axios实例
const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Token管理
const TokenManager = {
  getAccessToken: () => localStorage.getItem('access_token'),
  getRefreshToken: () => localStorage.getItem('refresh_token'),
  setTokens: (accessToken: string, refreshToken: string) => {
    localStorage.setItem('access_token', accessToken);
    localStorage.setItem('refresh_token', refreshToken);
  },
  clearTokens: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  },
};

// 是否正在刷新token
let isRefreshing = false;
// 待重试的请求队列
let requestQueue: Array<(token: string) => void> = [];

// 请求拦截器
service.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 添加access token到请求头
    const token = TokenManager.getAccessToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response) => {
    const res = response.data;

    // code为0表示成功
    if (res.code === 0) {
      return res;
    }

    // 处理业务错误
    message.error(res.message || '请求失败');
    return Promise.reject(new Error(res.message || '请求失败'));
  },
  async (error: AxiosError) => {
    if (!error.response) {
      message.error('网络错误，请检查网络连接');
      return Promise.reject(error);
    }

    const { data, config } = error.response as any;
    const originalRequest = config;

    // Token过期，尝试刷新
    if (data?.code === 20005) {
      if (!isRefreshing) {
        isRefreshing = true;
        const refreshToken = TokenManager.getRefreshToken();

        if (!refreshToken) {
          // 没有refresh token，跳转登录
          TokenManager.clearTokens();
          window.location.href = '/login';
          return Promise.reject(error);
        }

        try {
          // 调用刷新token接口
          const response = await axios.post(
            `${service.defaults.baseURL}/api/v1/auth/refresh`,
            { refresh_token: refreshToken }
          );

          if (response.data.code === 0) {
            const { access_token } = response.data.data;
            TokenManager.setTokens(access_token, refreshToken);

            // 重试队列中的请求
            requestQueue.forEach(cb => cb(access_token));
            requestQueue = [];

            // 重试原请求
            originalRequest.headers.Authorization = `Bearer ${access_token}`;
            return service(originalRequest);
          }
        } catch (refreshError) {
          // 刷新失败，清除token并跳转登录
          TokenManager.clearTokens();
          window.location.href = '/login';
          return Promise.reject(refreshError);
        } finally {
          isRefreshing = false;
        }
      } else {
        // 正在刷新token，将请求加入队列
        return new Promise((resolve) => {
          requestQueue.push((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(service(originalRequest));
          });
        });
      }
    }

    // 未授权，跳转登录
    if (data?.code === 20008 || data?.code === 20004) {
      TokenManager.clearTokens();
      window.location.href = '/login';
      return Promise.reject(error);
    }

    // 其他错误
    message.error(data?.message || '请求失败');
    return Promise.reject(error);
  }
);

export default service;
export { TokenManager };
```

---

#### API接口封装

```typescript
// src/api/auth.ts
import request from '@/utils/request';
import type { ApiResponse } from '@/types';

// 类型定义
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  confirm_password: string;
  nickname?: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface UserInfo {
  id: number;
  username: string;
  email: string;
  nickname: string;
  avatar_url: string;
  status: string;
  email_verified: boolean;
  last_login_at?: string;
  created_at: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  token_type: string;
  user: UserInfo;
}

// API方法
export const authApi = {
  // 用户注册
  register: (data: RegisterRequest) => {
    return request.post<ApiResponse<UserInfo>>('/api/v1/auth/register', data);
  },

  // 用户登录
  login: (data: LoginRequest) => {
    return request.post<ApiResponse<LoginResponse>>('/api/v1/auth/login', data);
  },

  // 获取当前用户信息
  getCurrentUser: () => {
    return request.get<ApiResponse<UserInfo>>('/api/v1/auth/me');
  },

  // 刷新token
  refreshToken: (refreshToken: string) => {
    return request.post<ApiResponse<{ access_token: string; expires_in: number }>>
      ('/api/v1/auth/refresh', { refresh_token: refreshToken });
  },

  // 登出
  logout: (refreshToken: string) => {
    return request.post<ApiResponse<{ message: string }>>
      ('/api/v1/auth/logout', { refresh_token: refreshToken });
  },
};
```

```typescript
// src/api/polish.ts
import request from '@/utils/request';
import type { ApiResponse } from '@/types';

export interface PolishRequest {
  content: string;
  style?: 'academic' | 'formal' | 'concise';
  language?: 'en' | 'zh';
  provider?: 'claude' | 'doubao';
}

export interface PolishResponse {
  polished_content: string;
  original_length: number;
  polished_length: number;
  suggestions: string[];
  provider_used: string;
  model_used: string;
}

export const polishApi = {
  // 段落润色
  polish: (data: PolishRequest) => {
    return request.post<ApiResponse<PolishResponse>>('/api/v1/polish', data);
  },

  // 查询记录列表
  getRecords: (params?: any) => {
    return request.get<ApiResponse<any>>('/api/v1/polish/records', { params });
  },

  // 根据TraceID查询记录
  getRecordByTraceId: (traceId: string) => {
    return request.get<ApiResponse<any>>(`/api/v1/polish/records/${traceId}`);
  },

  // 获取统计信息
  getStatistics: (params?: any) => {
    return request.get<ApiResponse<any>>('/api/v1/polish/statistics', { params });
  },
};
```

---

### 3.2 React示例代码

#### 用户上下文

```typescript
// src/contexts/AuthContext.tsx
import React, { createContext, useContext, useState, useEffect } from 'react';
import { authApi } from '@/api/auth';
import { TokenManager } from '@/utils/request';
import type { UserInfo, LoginRequest, RegisterRequest } from '@/api/auth';

interface AuthContextType {
  user: UserInfo | null;
  loading: boolean;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [loading, setLoading] = useState(true);

  // 初始化时检查登录状态
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    const token = TokenManager.getAccessToken();
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const response = await authApi.getCurrentUser();
      setUser(response.data);
    } catch (error) {
      TokenManager.clearTokens();
    } finally {
      setLoading(false);
    }
  };

  const login = async (data: LoginRequest) => {
    const response = await authApi.login(data);
    const { access_token, refresh_token, user: userInfo } = response.data;

    TokenManager.setTokens(access_token, refresh_token);
    setUser(userInfo);
  };

  const register = async (data: RegisterRequest) => {
    await authApi.register(data);
    // 注册成功后自动登录
    await login({ username: data.username, password: data.password });
  };

  const logout = async () => {
    const refreshToken = TokenManager.getRefreshToken();
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch (error) {
        // 忽略登出错误
      }
    }
    TokenManager.clearTokens();
    setUser(null);
  };

  const refreshUser = async () => {
    const response = await authApi.getCurrentUser();
    setUser(response.data);
  };

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};
```

---

#### 登录页面

```typescript
// src/pages/Login.tsx
import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import type { LoginRequest } from '@/api/auth';

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const onFinish = async (values: LoginRequest) => {
    setLoading(true);
    try {
      await login(values);
      message.success('登录成功');
      navigate('/dashboard');
    } catch (error: any) {
      message.error(error.message || '登录失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '100px auto' }}>
      <h1>用户登录</h1>
      <Form
        name="login"
        onFinish={onFinish}
        autoComplete="off"
        layout="vertical"
      >
        <Form.Item
          label="用户名"
          name="username"
          rules={[{ required: true, message: '请输入用户名' }]}
        >
          <Input placeholder="请输入用户名或邮箱" />
        </Form.Item>

        <Form.Item
          label="密码"
          name="password"
          rules={[{ required: true, message: '请输入密码' }]}
        >
          <Input.Password placeholder="请输入密码" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            登录
          </Button>
        </Form.Item>

        <div style={{ textAlign: 'center' }}>
          还没有账号？<a href="/register">立即注册</a>
        </div>
      </Form>
    </div>
  );
};

export default Login;
```

---

#### 注册页面

```typescript
// src/pages/Register.tsx
import React, { useState } from 'react';
import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import type { RegisterRequest } from '@/api/auth';

const Register: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const onFinish = async (values: RegisterRequest) => {
    setLoading(true);
    try {
      await register(values);
      message.success('注册成功');
      navigate('/dashboard');
    } catch (error: any) {
      message.error(error.message || '注册失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 400, margin: '100px auto' }}>
      <h1>用户注册</h1>
      <Form
        name="register"
        onFinish={onFinish}
        autoComplete="off"
        layout="vertical"
      >
        <Form.Item
          label="用户名"
          name="username"
          rules={[
            { required: true, message: '请输入用户名' },
            { pattern: /^[a-zA-Z0-9_]{3,50}$/, message: '3-50位，只能包含字母、数字、下划线' },
          ]}
        >
          <Input placeholder="请输入用户名" />
        </Form.Item>

        <Form.Item
          label="邮箱"
          name="email"
          rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '请输入有效的邮箱地址' },
          ]}
        >
          <Input placeholder="请输入邮箱" />
        </Form.Item>

        <Form.Item
          label="昵称"
          name="nickname"
          rules={[{ max: 50, message: '昵称最多50位' }]}
        >
          <Input placeholder="请输入昵称（可选）" />
        </Form.Item>

        <Form.Item
          label="密码"
          name="password"
          rules={[
            { required: true, message: '请输入密码' },
            { min: 8, message: '密码至少8位' },
            { pattern: /^(?=.*[A-Za-z])(?=.*\d)/, message: '密码必须包含字母和数字' },
          ]}
        >
          <Input.Password placeholder="至少8位，包含字母和数字" />
        </Form.Item>

        <Form.Item
          label="确认密码"
          name="confirm_password"
          dependencies={['password']}
          rules={[
            { required: true, message: '请确认密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('两次输入的密码不一致'));
              },
            }),
          ]}
        >
          <Input.Password placeholder="请再次输入密码" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} block>
            注册
          </Button>
        </Form.Item>

        <div style={{ textAlign: 'center' }}>
          已有账号？<a href="/login">立即登录</a>
        </div>
      </Form>
    </div>
  );
};

export default Register;
```

---

#### 论文润色页面

```typescript
// src/pages/Polish.tsx
import React, { useState } from 'react';
import { Form, Input, Button, Select, Card, message } from 'antd';
import { polishApi } from '@/api/polish';
import type { PolishRequest, PolishResponse } from '@/api/polish';

const { TextArea } = Input;
const { Option } = Select;

const Polish: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<PolishResponse | null>(null);

  const onFinish = async (values: PolishRequest) => {
    if (!values.content || values.content.length === 0) {
      message.warning('请输入要润色的内容');
      return;
    }

    if (values.content.length > 10000) {
      message.error('内容长度不能超过10000字符');
      return;
    }

    setLoading(true);
    try {
      const response = await polishApi.polish(values);
      setResult(response.data);
      message.success('润色成功');
    } catch (error: any) {
      message.error(error.message || '润色失败');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 1200, margin: '20px auto', padding: '0 20px' }}>
      <h1>论文润色</h1>

      <Form
        name="polish"
        onFinish={onFinish}
        layout="vertical"
        initialValues={{
          style: 'academic',
          language: 'zh',
          provider: 'doubao',
        }}
      >
        <Form.Item
          label="待润色内容"
          name="content"
          rules={[
            { required: true, message: '请输入要润色的内容' },
            { max: 10000, message: '内容长度不能超过10000字符' },
          ]}
        >
          <TextArea
            rows={8}
            placeholder="请输入需要润色的段落内容..."
            showCount
            maxLength={10000}
          />
        </Form.Item>

        <Form.Item label="润色风格" name="style">
          <Select>
            <Option value="academic">学术风格</Option>
            <Option value="formal">正式风格</Option>
            <Option value="concise">简洁风格</Option>
          </Select>
        </Form.Item>

        <Form.Item label="目标语言" name="language">
          <Select>
            <Option value="en">英文</Option>
            <Option value="zh">中文</Option>
          </Select>
        </Form.Item>

        <Form.Item label="AI提供商" name="provider">
          <Select>
            <Option value="doubao">豆包</Option>
            <Option value="claude">Claude</Option>
          </Select>
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading} size="large">
            开始润色
          </Button>
        </Form.Item>
      </Form>

      {result && (
        <Card title="润色结果" style={{ marginTop: 20 }}>
          <div style={{ marginBottom: 16 }}>
            <strong>润色后内容：</strong>
            <div style={{
              padding: 16,
              background: '#f5f5f5',
              borderRadius: 4,
              marginTop: 8,
              whiteSpace: 'pre-wrap'
            }}>
              {result.polished_content}
            </div>
          </div>

          {result.suggestions && result.suggestions.length > 0 && (
            <div style={{ marginBottom: 16 }}>
              <strong>改进建议：</strong>
              <ul>
                {result.suggestions.map((suggestion, index) => (
                  <li key={index}>{suggestion}</li>
                ))}
              </ul>
            </div>
          )}

          <div style={{ color: '#666', fontSize: 12 }}>
            原文长度：{result.original_length} |
            润色后长度：{result.polished_length} |
            使用模型：{result.model_used}
          </div>
        </Card>
      )}
    </div>
  );
};

export default Polish;
```

---

#### 路由守卫

```typescript
// src/components/PrivateRoute.tsx
import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '@/contexts/AuthContext';
import { Spin } from 'antd';

interface PrivateRouteProps {
  children: React.ReactNode;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const { user, loading } = useAuth();

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh'
      }}>
        <Spin size="large" />
      </div>
    );
  }

  return user ? <>{children}</> : <Navigate to="/login" replace />;
};

export default PrivateRoute;
```

---

### 3.3 Vue 3示例代码

#### 用户Store (Pinia)

```typescript
// src/stores/auth.ts
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authApi } from '@/api/auth';
import { TokenManager } from '@/utils/request';
import type { UserInfo, LoginRequest, RegisterRequest } from '@/api/auth';

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null);
  const loading = ref(false);

  // 检查登录状态
  const checkAuth = async () => {
    const token = TokenManager.getAccessToken();
    if (!token) return;

    try {
      const response = await authApi.getCurrentUser();
      user.value = response.data;
    } catch (error) {
      TokenManager.clearTokens();
      user.value = null;
    }
  };

  // 登录
  const login = async (data: LoginRequest) => {
    const response = await authApi.login(data);
    const { access_token, refresh_token, user: userInfo } = response.data;

    TokenManager.setTokens(access_token, refresh_token);
    user.value = userInfo;
  };

  // 注册
  const register = async (data: RegisterRequest) => {
    await authApi.register(data);
    // 注册成功后自动登录
    await login({ username: data.username, password: data.password });
  };

  // 登出
  const logout = async () => {
    const refreshToken = TokenManager.getRefreshToken();
    if (refreshToken) {
      try {
        await authApi.logout(refreshToken);
      } catch (error) {
        // 忽略错误
      }
    }
    TokenManager.clearTokens();
    user.value = null;
  };

  // 刷新用户信息
  const refreshUser = async () => {
    const response = await authApi.getCurrentUser();
    user.value = response.data;
  };

  return {
    user,
    loading,
    checkAuth,
    login,
    register,
    logout,
    refreshUser,
  };
});
```

---

#### 登录页面

```vue
<!-- src/views/Login.vue -->
<template>
  <div class="login-container">
    <h1>用户登录</h1>
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="80px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="formData.username"
          placeholder="请输入用户名或邮箱"
        />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          placeholder="请输入密码"
          show-password
        />
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :loading="loading" @click="handleSubmit" block>
          登录
        </el-button>
      </el-form-item>
    </el-form>

    <div style="text-align: center">
      还没有账号？<router-link to="/register">立即注册</router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useAuthStore } from '@/stores/auth';
import type { FormInstance } from 'element-plus';

const router = useRouter();
const authStore = useAuthStore();

const formRef = ref<FormInstance>();
const loading = ref(false);

const formData = reactive({
  username: '',
  password: '',
});

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
};

const handleSubmit = async () => {
  if (!formRef.value) return;

  await formRef.value.validate(async (valid) => {
    if (!valid) return;

    loading.value = true;
    try {
      await authStore.login(formData);
      ElMessage.success('登录成功');
      router.push('/dashboard');
    } catch (error: any) {
      ElMessage.error(error.message || '登录失败');
    } finally {
      loading.value = false;
    }
  });
};
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 100px auto;
  padding: 20px;
}
</style>
```

---

## 四、常见问题

### 4.1 Token存储方案

**推荐方案**：localStorage

```typescript
// 存储
localStorage.setItem('access_token', accessToken);
localStorage.setItem('refresh_token', refreshToken);

// 读取
const accessToken = localStorage.getItem('access_token');
const refreshToken = localStorage.getItem('refresh_token');

// 清除
localStorage.removeItem('access_token');
localStorage.removeItem('refresh_token');
```

**不推荐**：cookie（跨域问题）

---

### 4.2 Token刷新策略

**方案一：被动刷新（推荐）**
- 在响应拦截器中检测到token过期（code=20005）时自动刷新
- 实现简单，不会产生不必要的请求

**方案二：主动刷新**
- 计算token过期时间，提前5-10分钟刷新
- 需要维护定时器，实现较复杂

---

### 4.3 多标签页同步

使用Storage事件监听token变化：

```typescript
window.addEventListener('storage', (e) => {
  if (e.key === 'access_token') {
    if (!e.newValue) {
      // Token被清除，跳转登录
      window.location.href = '/login';
    } else {
      // Token更新，刷新用户信息
      authStore.refreshUser();
    }
  }
});
```

---

### 4.4 CORS跨域配置

如果前后端分离部署，后端需要配置CORS：

```go
// 后端已配置CORS，允许以下来源：
// - 开发环境：http://localhost:3000, http://localhost:5173
// - 生产环境：需要在后端config中配置

// 前端无需额外配置
```

---

## 五、测试建议

### 5.1 接口测试流程

1. 使用Postman或Apifox导入OpenAPI文档（[openapi.yaml](openapi.yaml)）
2. 按照[APIFOX_IMPORT_GUIDE.md](APIFOX_IMPORT_GUIDE.md)配置环境变量
3. 按顺序测试：注册 → 登录 → 获取用户信息 → 论文润色 → 登出

### 5.2 前端开发建议

1. 先用Mock数据开发UI
2. 接口联调时使用开发环境API
3. 添加完善的错误处理和加载状态
4. 实现Token自动刷新机制
5. 添加请求/响应日志（开发环境）

---

## 六、部署说明

### 6.1 环境变量

```bash
# .env.development
VITE_API_BASE_URL=http://localhost:8080

# .env.production
VITE_API_BASE_URL=https://api.paperai.com
```

### 6.2 生产环境注意事项

1. 使用HTTPS协议
2. 配置正确的API域名
3. 启用请求加密
4. 添加请求频率限制
5. 配置CDN加速

---

## 七、联系方式

如有问题，请联系：

- 后端开发：[你的联系方式]
- API文档：[openapi.yaml](openapi.yaml)
- 接口测试指南：[APIFOX_IMPORT_GUIDE.md](APIFOX_IMPORT_GUIDE.md)
- 认证实现文档：[AUTH_IMPLEMENTATION.md](AUTH_IMPLEMENTATION.md)

---

**最后更新时间**: 2024-11-27
