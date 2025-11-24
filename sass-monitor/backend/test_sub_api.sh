#!/bin/bash

# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

ORG_ID="4808d051-3d69-4a82-82b2-28c134ea6572"

echo "Testing GET /api/v1/organizations/${ORG_ID}/subscriptions..."
curl -s -X GET "http://localhost:8080/api/v1/organizations/${ORG_ID}/subscriptions?page=1&page_size=2" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

echo ""
echo "Checking for subscription user fields:"
curl -s -X GET "http://localhost:8080/api/v1/organizations/${ORG_ID}/subscriptions?page=1&page_size=2" \
  -H "Authorization: Bearer $TOKEN" | grep -E "username|user_email|plan_pricing|days_until_expiry"
