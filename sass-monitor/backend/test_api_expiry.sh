#!/bin/bash

# 获取token
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' > /tmp/login_result.json

TOKEN=$(cat /tmp/login_result.json | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "Token obtained successfully"
echo ""

# 测试组织列表API
echo "====== Testing Organizations API ======"
curl -s -X GET "http://localhost:8080/api/v1/organizations?page=1&page_size=2" \
  -H "Authorization: Bearer $TOKEN" > /tmp/orgs.json

echo "Response:"
cat /tmp/orgs.json | python3 -m json.tool

# 获取第一个组织ID
ORG_ID=$(cat /tmp/orgs.json | python3 -c "import sys, json; d=json.load(sys.stdin); print(d['data'][0]['id'] if d.get('data') else '')")

if [ -n "$ORG_ID" ]; then
  echo ""
  echo "====== Testing Organization Detail (ID: $ORG_ID) ======"
  curl -s -X GET "http://localhost:8080/api/v1/organizations/$ORG_ID" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

  echo ""
  echo "====== Testing Organization Subscriptions ======"
  curl -s -X GET "http://localhost:8080/api/v1/organizations/$ORG_ID/subscriptions?page=1&page_size=5" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool

  echo ""
  echo "====== Testing Send Expiry Reminder ======"
  curl -s -X POST "http://localhost:8080/api/v1/organizations/$ORG_ID/send-expiry-reminder" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
fi
