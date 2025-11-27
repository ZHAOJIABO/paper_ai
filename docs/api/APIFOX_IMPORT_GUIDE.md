# 如何将API导入到Apifox

## 方法一：导入OpenAPI文档（推荐）

### 步骤：

1. **打开Apifox**
   - 启动Apifox应用
   - 选择或创建一个项目

2. **导入OpenAPI文档**
   - 点击顶部菜单 **"导入"** 按钮
   - 选择 **"OpenAPI / Swagger"**
   - 选择 **"从文件导入"**
   - 选择项目根目录的 `openapi.yaml` 文件
   - 点击 **"确定"**

3. **配置环境变量**
   - 在Apifox左侧菜单找到 **"环境管理"**
   - 添加一个新环境，例如 "本地开发"
   - 添加环境变量：
     ```
     baseUrl: http://localhost:8080
     accessToken: （登录后填入）
     ```

4. **测试接口**
   - 导入完成后，所有接口会出现在左侧接口列表中
   - 按照下面的测试流程进行测试

## 方法二：手动创建接口

如果导入失败，可以手动创建接口：

### 1. 创建用户注册接口

- **接口名称**: 用户注册
- **请求方法**: POST
- **请求URL**: `{{baseUrl}}/api/v1/auth/register`
- **请求头**:
  ```
  Content-Type: application/json
  ```
- **请求体**（Body - JSON）:
  ```json
  {
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234",
    "confirm_password": "Test1234",
    "nickname": "测试用户"
  }
  ```

### 2. 创建用户登录接口

- **接口名称**: 用户登录
- **请求方法**: POST
- **请求URL**: `{{baseUrl}}/api/v1/auth/login`
- **请求头**:
  ```
  Content-Type: application/json
  ```
- **请求体**（Body - JSON）:
  ```json
  {
    "username": "testuser",
    "password": "Test1234"
  }
  ```
- **后置操作**（提取token）:
  在"后置操作"标签页添加脚本：
  ```javascript
  // 提取access_token并保存到环境变量
  const response = pm.response.json();
  if (response.code === 0 && response.data.access_token) {
    pm.environment.set("accessToken", response.data.access_token);
    console.log("Token已保存到环境变量");
  }
  ```

### 3. 创建需要认证的接口

以"获取当前用户信息"为例：

- **接口名称**: 获取当前用户信息
- **请求方法**: GET
- **请求URL**: `{{baseUrl}}/api/v1/auth/me`
- **请求头**:
  ```
  Authorization: Bearer {{accessToken}}
  ```

## 测试流程

### 第一步：启动服务

```bash
cd /Users/zhaojiabo/Documents/trae_projects/paper_ai
./paper_ai
```

确保服务在 `http://localhost:8080` 运行。

### 第二步：测试认证流程

1. **测试用户注册**
   - 发送 `POST /api/v1/auth/register` 请求
   - 确认返回 `code: 0` 表示注册成功

2. **测试用户登录**
   - 发送 `POST /api/v1/auth/login` 请求
   - 复制返回的 `access_token`
   - 如果设置了后置脚本，token会自动保存到环境变量

3. **测试获取用户信息**
   - 发送 `GET /api/v1/auth/me` 请求
   - 确认请求头包含 `Authorization: Bearer <token>`
   - 确认返回当前用户信息

4. **测试论文润色（需要认证）**
   - 发送 `POST /api/v1/polish` 请求
   - 请求体：
     ```json
     {
       "content": "这是一段测试文本。",
       "style": "academic",
       "language": "zh",
       "provider": "doubao"
     }
     ```
   - 确认返回润色结果

5. **测试刷新Token**
   - 发送 `POST /api/v1/auth/refresh` 请求
   - 使用登录返回的 `refresh_token`

6. **测试登出**
   - 发送 `POST /api/v1/auth/logout` 请求
   - 使用 `refresh_token`

## Apifox环境变量配置

在Apifox中配置以下环境变量，方便测试：

| 变量名 | 初始值 | 说明 |
|-------|--------|------|
| baseUrl | http://localhost:8080 | 基础URL |
| accessToken | （空） | 访问令牌（登录后自动填充） |
| refreshToken | （空） | 刷新令牌（登录后手动填充） |
| testUsername | testuser | 测试用户名 |
| testEmail | test@example.com | 测试邮箱 |
| testPassword | Test1234 | 测试密码 |

## 常见问题

### Q1: 导入OpenAPI文档后，接口没有显示？

**解决方法**：
- 确认导入时选择了正确的项目
- 检查文件格式是否正确
- 尝试重新导入

### Q2: 401 Unauthorized 错误？

**解决方法**：
- 确认已登录并获取token
- 检查请求头是否包含 `Authorization: Bearer <token>`
- 确认token格式正确（注意Bearer后有空格）
- 检查token是否过期（默认2小时过期）

### Q3: 如何自动刷新过期的token？

**解决方法**：
在Apifox的"前置操作"中添加脚本：
```javascript
// 检查token是否即将过期，自动刷新
const tokenExpiry = pm.environment.get("tokenExpiry");
if (tokenExpiry && Date.now() > tokenExpiry - 300000) { // 提前5分钟刷新
  // 调用refresh接口刷新token
  // 这里需要实现refresh逻辑
}
```

### Q4: 如何快速切换测试用户？

**解决方法**：
- 在Apifox中创建多个环境（用户A、用户B等）
- 每个环境配置不同的用户名、密码和token
- 通过切换环境来切换测试用户

## 接口测试顺序建议

按以下顺序测试接口最合理：

```
1. POST /api/v1/auth/register     (注册用户)
   ↓
2. POST /api/v1/auth/login        (登录获取token)
   ↓
3. GET /api/v1/auth/me            (验证token有效性)
   ↓
4. POST /api/v1/polish            (测试业务接口)
   ↓
5. GET /api/v1/polish/records     (查询记录)
   ↓
6. GET /api/v1/polish/statistics  (查看统计)
   ↓
7. POST /api/v1/auth/refresh      (刷新token)
   ↓
8. POST /api/v1/auth/logout       (登出)
```

## 批量测试脚本

也可以使用项目根目录的 `test_auth.sh` 脚本进行批量测试：

```bash
chmod +x test_auth.sh
./test_auth.sh
```

该脚本会自动测试所有认证相关接口。

## 参考资料

- OpenAPI文档位置: `/Users/zhaojiabo/Documents/trae_projects/paper_ai/openapi.yaml`
- 完整实现文档: `AUTH_IMPLEMENTATION.md`
- 测试脚本: `test_auth.sh`

---

**提示**: 建议在Apifox中创建一个接口测试套件（Test Suite），将所有接口按流程组织起来，可以一键运行所有测试。
