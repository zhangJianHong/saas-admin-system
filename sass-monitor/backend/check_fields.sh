#!/bin/bash
curl -s -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}' > /tmp/tok.json
T=$(cat /tmp/tok.json | python3 -c "import sys, json; print(json.load(sys.stdin)['token'])")
curl -s -X GET "http://localhost:8080/api/v1/organizations?page=1&page_size=1" -H "Authorization: Bearer $T" > /tmp/check.json
echo "Raw response:"
cat /tmp/check.json
echo ""
echo "Checking for subscription_status field:"
grep "subscription_status" /tmp/check.json
echo "Checking for active_subscription_count field:"
grep "active_subscription_count" /tmp/check.json
