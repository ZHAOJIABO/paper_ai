#!/bin/bash

# Paper AI 测试脚本

echo "=== Paper AI 测试脚本 ==="
echo ""

# 检查服务是否运行
echo "1. 测试健康检查接口..."
curl -s http://localhost:8080/health | jq .
echo ""

# 测试段落润色接口（英文）
echo "2. 测试英文段落润色（academic风格）..."
curl -s -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This paper discuss the important of machine learning in modern software development.",
    "style": "academic",
    "language": "en"
  }' | jq .
echo ""

# 测试段落润色接口（中文）
echo "3. 测试中文段落润色（academic风格）..."
curl -s -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "content": "这篇文章讨论了机器学习在软件开发中的作用。",
    "style": "academic",
    "language": "zh"
  }' | jq .
echo ""

# 测试参数错误
echo "4. 测试参数错误（缺少必填字段）..."
curl -s -X POST http://localhost:8080/api/v1/polish \
  -H "Content-Type: application/json" \
  -d '{
    "style": "academic"
  }' | jq .
echo ""

echo "=== 测试完成 ==="
