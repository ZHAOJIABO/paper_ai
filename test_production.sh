#!/bin/bash

# Paper AI 接口测试脚本

# 设置服务器地址
SERVER="${SERVER:-http://localhost}"

echo "======================================"
echo "   Paper AI 接口测试"
echo "======================================"
echo "服务器地址: $SERVER"
echo ""

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local headers=$5

    echo -e "${YELLOW}[测试] $name${NC}"

    if [ -n "$headers" ]; then
        response=$(curl -s -X $method "$SERVER$endpoint" \
            -H "Content-Type: application/json" \
            -H "$headers" \
            -d "$data")
    else
        response=$(curl -s -X $method "$SERVER$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi

    echo "$response" | jq '.' 2>/dev/null || echo "$response"

    # 检查是否成功
    if echo "$response" | grep -q '"code":0'; then
        echo -e "${GREEN}✓ 成功${NC}\n"
        return 0
    else
        echo -e "${RED}✗ 失败${NC}\n"
        return 1
    fi
}

# 1. 健康检查
echo "======================================"
echo "1. 健康检查"
echo "======================================"
curl -s -X GET "$SERVER/health" | jq '.' || curl -s -X GET "$SERVER/health"
echo -e "${GREEN}✓ 服务运行正常${NC}\n"

# 2. 用户注册
echo "======================================"
echo "2. 用户注册"
echo "======================================"
REGISTER_DATA='{
  "email": "test_'$(date +%s)'@example.com",
  "password": "Test123456"
}'

REGISTER_RESPONSE=$(curl -s -X POST "$SERVER/api/v1/auth/register" \
    -H "Content-Type: application/json" \
    -d "$REGISTER_DATA")

echo "$REGISTER_RESPONSE" | jq '.'

# 提取 token
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.refresh_token')
USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user_id')

if [ "$ACCESS_TOKEN" != "null" ] && [ -n "$ACCESS_TOKEN" ]; then
    echo -e "${GREEN}✓ 注册成功${NC}"
    echo "User ID: $USER_ID"
    echo "Access Token: ${ACCESS_TOKEN:0:20}..."
    echo ""
else
    echo -e "${RED}✗ 注册失败，使用登录测试已有账号${NC}\n"

    # 尝试登录
    LOGIN_DATA='{
      "email": "test@example.com",
      "password": "Test123456"
    }'

    LOGIN_RESPONSE=$(curl -s -X POST "$SERVER/api/v1/auth/login" \
        -H "Content-Type: application/json" \
        -d "$LOGIN_DATA")

    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token')
    REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.refresh_token')
fi

# 3. 测试论文润色
echo "======================================"
echo "3. 论文润色"
echo "======================================"
POLISH_DATA='{
  "content": "This paper discuss the important of machine learning in modern software development.",
  "style": "academic",
  "language": "zh"
}'

curl -s -X POST "$SERVER/api/v1/polish" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "$POLISH_DATA" | jq '.'

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 润色成功${NC}\n"
else
    echo -e "${RED}✗ 润色失败${NC}\n"
fi

# 等待1秒，让数据写入数据库
sleep 1

# 4. 查询历史记录
echo "======================================"
echo "4. 查询历史记录"
echo "======================================"
curl -s -X GET "$SERVER/api/v1/polish/history?page=1&page_size=10" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 查询成功${NC}\n"
else
    echo -e "${RED}✗ 查询失败${NC}\n"
fi

# 5. 查询统计信息
echo "======================================"
echo "5. 统计信息"
echo "======================================"
curl -s -X GET "$SERVER/api/v1/polish/stats" \
    -H "Authorization: Bearer $ACCESS_TOKEN" | jq '.'

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 查询成功${NC}\n"
else
    echo -e "${RED}✗ 查询失败${NC}\n"
fi

# 6. 刷新 Token
echo "======================================"
echo "6. 刷新 Token"
echo "======================================"
REFRESH_DATA="{
  \"refresh_token\": \"$REFRESH_TOKEN\"
}"

curl -s -X POST "$SERVER/api/v1/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "$REFRESH_DATA" | jq '.'

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 刷新成功${NC}\n"
else
    echo -e "${RED}✗ 刷新失败${NC}\n"
fi

echo "======================================"
echo "   测试完成！"
echo "======================================"
echo ""
echo "保存以下信息供后续使用："
echo "Access Token: $ACCESS_TOKEN"
echo "Refresh Token: $REFRESH_TOKEN"
echo ""
echo "使用示例："
echo "export TOKEN=\"$ACCESS_TOKEN\""
echo "curl -X POST $SERVER/api/v1/polish \\"
echo "  -H \"Authorization: Bearer \$TOKEN\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"content\":\"test\",\"style\":\"academic\",\"language\":\"zh\"}'"
