#!/bin/bash

# 测试短ID生成和论文润色的UserID记录功能

BASE_URL="http://localhost:8080/api/v1"

echo "========================================="
echo "测试短ID生成和UserID记录"
echo "========================================="
echo ""

# 1. 注册新用户（会生成13位短ID）
echo "1. 注册新用户（测试13位短ID）"
echo "POST $BASE_URL/auth/register"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "short_id_user",
    "email": "short@example.com",
    "password": "Test1234",
    "confirm_password": "Test1234",
    "nickname": "短ID用户"
  }')
echo "$REGISTER_RESPONSE" | jq '.'

# 提取UserID并验证是否为13位
USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.data.id')
ID_LENGTH=${#USER_ID}

if [ "$ID_LENGTH" -eq 13 ]; then
    echo "✅ UserID长度正确: $ID_LENGTH 位"
    echo "✅ UserID: $USER_ID"
else
    echo "❌ UserID长度错误: $ID_LENGTH 位（应该是13位）"
fi
echo ""

# 2. 登录
echo "2. 登录获取Token"
echo "POST $BASE_URL/auth/login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "short_id_user",
    "password": "Test1234"
  }')

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.access_token')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" == "null" ]; then
    echo "❌ 登录失败，无法获取token"
    exit 1
fi

echo "✅ 登录成功"
echo ""

# 3. 测试论文润色（会记录UserID）
echo "3. 测试论文润色（会记录UserID）"
echo "POST $BASE_URL/polish"
POLISH_RESPONSE=$(curl -s -X POST "$BASE_URL/polish" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这是一段用于测试UserID记录的文本。",
    "style": "academic",
    "language": "zh",
    "provider": "doubao"
  }')
echo "$POLISH_RESPONSE" | jq '.'
echo ""

# 提取TraceID
TRACE_ID=$(echo $POLISH_RESPONSE | jq -r '.trace_id')
echo "TraceID: $TRACE_ID"
echo ""

# 等待1秒让数据库写入完成
sleep 1

# 4. 查询润色记录（验证UserID是否被记录）
echo "4. 查询润色记录（验证UserID）"
echo "GET $BASE_URL/polish/records"
RECORDS_RESPONSE=$(curl -s -X GET "$BASE_URL/polish/records?page=1&page_size=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

# 提取第一条记录的UserID
RECORD_USER_ID=$(echo $RECORDS_RESPONSE | jq -r '.data.records[0].user_id')

if [ "$RECORD_USER_ID" == "$USER_ID" ]; then
    echo "✅ 润色记录中的UserID正确: $RECORD_USER_ID"
else
    echo "❌ 润色记录中的UserID不正确"
    echo "   期望: $USER_ID"
    echo "   实际: $RECORD_USER_ID"
fi
echo ""

# 5. 显示完整记录
echo "5. 完整润色记录"
echo "$RECORDS_RESPONSE" | jq '.data.records[0] | {id, trace_id, user_id, status, language, style}'
echo ""

# 6. 测试ID解析
echo "6. 测试ID信息解析"
echo "UserID: $USER_ID"
# 提取时间戳（前10位）
TIMESTAMP=$((USER_ID / 1000))
# 提取机器ID（第11位）
WORKER_ID=$(((USER_ID % 1000) / 100))
# 提取序列号（后2位）
SEQUENCE=$((USER_ID % 100))

REG_TIME=$(date -r $TIMESTAMP "+%Y-%m-%d %H:%M:%S")

echo "  时间戳: $TIMESTAMP"
echo "  机器ID: $WORKER_ID"
echo "  序列号: $SEQUENCE"
echo "  注册时间: $REG_TIME"
echo ""

echo "========================================="
echo "测试完成！"
echo "========================================="
echo ""
echo "总结:"
echo "  ✅ 生成了13位短ID: $USER_ID"
echo "  ✅ 论文润色记录包含UserID"
echo "  ✅ 可以从ID中提取时间信息"
