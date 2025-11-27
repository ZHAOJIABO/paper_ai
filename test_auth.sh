#!/bin/bash

BASE_URL="http://localhost:8080/api/v1"

echo "========================================="
echo "测试用户认证功能"
echo "========================================="
echo ""

# 1. 测试用户注册
echo "1. 测试用户注册"
echo "POST $BASE_URL/auth/register"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "testuser@example.com",
    "password": "Test1234",
    "confirm_password": "Test1234",
    "nickname": "测试用户"
  }')
echo "Response: $REGISTER_RESPONSE"
echo ""

# 2. 测试用户登录
echo "2. 测试用户登录"
echo "POST $BASE_URL/auth/login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test1234"
  }')
echo "Response: $LOGIN_RESPONSE"
echo ""

# 提取access_token和refresh_token
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "⚠️  登录失败，无法获取access_token"
    exit 1
fi

echo "✅ 登录成功，获得access_token"
echo ""

# 3. 测试获取当前用户信息（需要认证）
echo "3. 测试获取当前用户信息（需要认证）"
echo "GET $BASE_URL/auth/me"
ME_RESPONSE=$(curl -s -X GET "$BASE_URL/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "Response: $ME_RESPONSE"
echo ""

# 4. 测试刷新Token
echo "4. 测试刷新Token"
echo "POST $BASE_URL/auth/refresh"
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "Response: $REFRESH_RESPONSE"
echo ""

# 5. 测试论文润色（需要认证）
echo "5. 测试论文润色（需要认证）"
echo "POST $BASE_URL/polish"
POLISH_RESPONSE=$(curl -s -X POST "$BASE_URL/polish" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一段测试文本，用于测试论文润色功能。",
    "style": "academic",
    "language": "zh",
    "provider": "doubao"
  }')
echo "Response: $POLISH_RESPONSE"
echo ""

# 6. 测试登出
echo "6. 测试登出"
echo "POST $BASE_URL/auth/logout"
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")
echo "Response: $LOGOUT_RESPONSE"
echo ""

# 7. 测试未认证访问（应该返回401）
echo "7. 测试未认证访问（应该返回401）"
echo "GET $BASE_URL/polish/records"
UNAUTH_RESPONSE=$(curl -s -X GET "$BASE_URL/polish/records")
echo "Response: $UNAUTH_RESPONSE"
echo ""

echo "========================================="
echo "测试完成！"
echo "========================================="
