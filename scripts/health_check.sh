#!/bin/bash

# 检测服务器 IP 是否可访问

SERVER_IP="${SERVER_IP:-your_server_ip}"
CHECK_URL="${CHECK_URL:-http://$SERVER_IP}"
WEBHOOK="${WEBHOOK:-}"  # 钉钉/企业微信等通知 webhook

echo "检测服务可用性: $CHECK_URL"

# 检测 HTTP 服务
if curl -s -o /dev/null -w "%{http_code}" --connect-timeout 10 $CHECK_URL | grep -q "200\|301\|302"; then
    echo "[$(date)] ✓ 服务正常"
    exit 0
else
    echo "[$(date)] ✗ 服务异常！"

    # 发送通知（如果配置了 webhook）
    if [ -n "$WEBHOOK" ]; then
        curl -X POST $WEBHOOK \
            -H 'Content-Type: application/json' \
            -d "{\"msg\": \"Paper AI 服务异常！时间: $(date)\"}"
    fi

    exit 1
fi
