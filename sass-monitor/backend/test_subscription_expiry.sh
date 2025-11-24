#!/bin/bash

# 获取token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | \
  python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "====== 1. 测试组织列表（带订阅到期信息） ======"
echo
curl -s -X GET "http://localhost:8080/api/v1/light-admin/organizations?page=1&page_size=2" \
  -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo
echo

# 获取第一个组织ID进行详细测试
ORG_ID=$(curl -s -X GET "http://localhost:8080/api/v1/light-admin/organizations?page=1&page_size=1" \
  -H "Authorization: Bearer $TOKEN" | \
  python3 -c "import sys, json; d=json.load(sys.stdin); print(d['data'][0]['id'] if d['data'] else '')")

if [ -n "$ORG_ID" ]; then
  echo "====== 2. 测试组织详情（带订阅到期状态） ======"
  echo "组织ID: $ORG_ID"
  echo
  curl -s -X GET "http://localhost:8080/api/v1/light-admin/organizations/$ORG_ID" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
  echo
  echo

  echo "====== 3. 测试组织订阅列表（含用户信息和到期天数） ======"
  echo
  curl -s -X GET "http://localhost:8080/api/v1/light-admin/organizations/$ORG_ID/subscriptions?page=1&page_size=5" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
  echo
  echo

  echo "====== 4. 测试发送到期提醒邮件（预留功能） ======"
  echo
  curl -s -X POST "http://localhost:8080/api/v1/light-admin/organizations/$ORG_ID/send-expiry-reminder" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
  echo
fi

echo "====== 测试完成 ======"
