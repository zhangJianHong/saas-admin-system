#!/bin/bash

# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "Token: $TOKEN"
echo ""

# 测试组织列表 API
echo "Testing GET /api/v1/organizations..."
curl -s -X GET "http://localhost:8080/api/v1/organizations?page=1&page_size=1" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

echo ""
echo "Checking for new fields:"
curl -s -X GET "http://localhost:8080/api/v1/organizations?page=1&page_size=1" \
  -H "Authorization: Bearer $TOKEN" | grep -E "subscription_status|active_subscription_count|subscription_end_date|days_until_expiration"
