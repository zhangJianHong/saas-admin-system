#!/bin/bash

# 获取token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

echo "====== 用户管理API测试 ======"
echo "Token获取成功: ${TOKEN:0:30}..."
echo

# 1. 测试获取用户列表
echo "1. 测试获取用户列表"
curl -s -X GET "http://localhost:8080/api/v1/users?page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN" | jq -c '{total, page, page_size, data_count: (.data | length)}'
echo

# 2. 测试不支持的操作
echo "2. 测试创建用户（应返回405）"
RESULT=$(curl -s -w "\n%{http_code}" -X POST "http://localhost:8080/api/v1/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username":"test"}')
echo "HTTP状态码: $(echo "$RESULT" | tail -1)"
echo "响应: $(echo "$RESULT" | head -1)"
echo

echo "====== 测试完成 ======"
