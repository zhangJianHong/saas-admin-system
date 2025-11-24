#!/bin/bash

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}' | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")

echo "====== 用户管理API完整测试 ======"
echo

echo "1. 获取用户列表"
curl -s -X GET "http://localhost:8080/api/v1/users?page=1&page_size=3" -H "Authorization: Bearer $TOKEN" | python3 -m json.tool | head -30
echo
echo

echo "2. 获取用户详情"
curl -s -X GET "http://localhost:8080/api/v1/users/302834b9-669a-474f-ab08-1110274f2e2b" -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo
echo

echo "3. 获取用户组织列表"
curl -s -X GET "http://localhost:8080/api/v1/users/302834b9-669a-474f-ab08-1110274f2e2b/organizations?page=1&page_size=3" -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo
echo

echo "4. 获取用户工作空间列表"
curl -s -X GET "http://localhost:8080/api/v1/users/302834b9-669a-474f-ab08-1110274f2e2b/workspaces?page=1&page_size=3" -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo
echo

echo "5. 获取用户订阅列表"
curl -s -X GET "http://localhost:8080/api/v1/users/302834b9-669a-474f-ab08-1110274f2e2b/subscriptions?page=1&page_size=3" -H "Authorization: Bearer $TOKEN" | python3 -m json.tool
echo

echo "====== 测试完成 ======"
