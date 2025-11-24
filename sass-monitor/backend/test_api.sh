#!/bin/bash

# API测试脚本
API_BASE="http://localhost:8080/api/v1"
TOKEN=""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TESTS_PASSED=0
TESTS_FAILED=0

# 测试函数
test_api() {
    local method=$1
    local url=$2
    local data=$3
    local expected_code=$4
    local description=$5

    echo -e "\n${YELLOW}测试: $description${NC}"
    echo "请求: $method $url"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "%{http_code}" -X GET "$url" -H "Authorization: Bearer $TOKEN")
    elif [ "$method" = "POST" ]; then
        response=$(curl -s -w "%{http_code}" -X POST "$url" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$data")
    elif [ "$method" = "PUT" ]; then
        response=$(curl -s -w "%{http_code}" -X PUT "$url" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$data")
    elif [ "$method" = "DELETE" ]; then
        response=$(curl -s -w "%{http_code}" -X DELETE "$url" -H "Authorization: Bearer $TOKEN")
    fi

    # 分离HTTP状态码和响应体
    http_code="${response: -3}"
    response_body="${response%???}"

    echo "状态码: $http_code (期望: $expected_code)"

    if [ "$http_code" -eq "$expected_code" ]; then
        echo -e "${GREEN}✓ 通过${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ 失败${NC}"
        echo "响应: $response_body"
        ((TESTS_FAILED++))
    fi
}

# 首先登录获取token
echo -e "\n${YELLOW}=== 步骤1: 用户认证 ===${NC}"
auth_response=$(curl -s -X POST "$API_BASE/auth/login" -H "Content-Type: application/json" -d '{
    "username": "admin",
    "password": "admin123"
}')

echo "认证响应: $auth_response"

# 提取token (这里需要根据实际的认证响应格式调整)
TOKEN=$(echo $auth_response | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ 认证失败，无法获取token${NC}"
    exit 1
fi

echo -e "${GREEN}✓ 认证成功，获得token${NC}"

# 测试组织管理API（只读模式）
echo -e "\n${YELLOW}=== 步骤2: 组织管理API测试（只读模式） ===${NC}"

# 测试获取组织列表
test_api "GET" "$API_BASE/organizations" "" "200" "获取组织列表"

# 测试获取组织详情（使用随机UUID）
test_api "GET" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001" "" "200" "获取组织详情"

# 测试获取组织用户列表
test_api "GET" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001/users" "" "200" "获取组织用户列表"

# 测试获取组织工作空间
test_api "GET" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001/workspaces" "" "200" "获取组织工作空间"

# 测试获取组织订阅信息
test_api "GET" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001/subscriptions" "" "200" "获取组织订阅信息"

# 测试不支持的操作（应该返回405 Method Not Allowed）
echo -e "\n${YELLOW}=== 步骤3: 测试不支持的操作（只读模式验证） ===${NC}"

test_api "POST" "$API_BASE/organizations" '{"name":"Test Org","owner_id":"00000000-0000-0000-0000-000000000001"}' "405" "创建组织（不支持）"

test_api "PUT" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001" '{"name":"Updated Org"}' "405" "更新组织（不支持）"

test_api "DELETE" "$API_BASE/organizations/00000000-0000-0000-0000-000000000001" "" "405" "删除组织（不支持）"

# 测试订阅计划管理API（完整CRUD）
echo -e "\n${YELLOW}=== 步骤4: 订阅计划管理API测试（完整CRUD） ===${NC}"

# 创建测试订阅计划
create_plan_data='{
    "tier_name": "测试计划-专业版",
    "pricing_monthly": 99.99,
    "pricing_quarterly": 269.99,
    "pricing_yearly": 999.99,
    "limits": "{\"users\": 100, \"storage\": \"100GB\"}",
    "features": "[\"无限API调用\", \"优先支持\", \"自定义集成\"]",
    "target_users": "专业团队和企业用户",
    "is_active": true
}'

create_response=$(curl -s -w "%{http_code}" -X POST "$API_BASE/subscription-plans" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$create_plan_data")

create_http_code="${create_response: -3}"
create_response_body="${create_response%???}"

if [ "$create_http_code" -eq "201" ]; then
    echo -e "${GREEN}✓ 创建订阅计划成功${NC}"
    ((TESTS_PASSED++))

    # 提取创建的计划ID
    PLAN_ID=$(echo $create_response_body | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "创建的计划ID: $PLAN_ID"

    if [ -n "$PLAN_ID" ]; then
        # 测试获取订阅计划详情
        test_api "GET" "$API_BASE/subscription-plans/$PLAN_ID" "" "200" "获取订阅计划详情"

        # 测试更新订阅计划
        update_plan_data='{
            "tier_name": "测试计划-专业版(更新)",
            "pricing_monthly": 89.99,
            "is_active": false
        }'
        test_api "PUT" "$API_BASE/subscription-plans/$PLAN_ID" "$update_plan_data" "200" "更新订阅计划"

        # 测试获取活跃的订阅计划
        test_api "GET" "$API_BASE/subscription-plans/active" "" "200" "获取活跃订阅计划"

        # 测试按价格范围搜索
        test_api "GET" "$API_BASE/subscription-plans/search?min_price=50&max_price=150" "" "200" "按价格范围搜索订阅计划"

        # 最后删除测试计划
        test_api "DELETE" "$API_BASE/subscription-plans/$PLAN_ID" "" "200" "删除订阅计划"
    fi
else
    echo -e "${RED}✗ 创建订阅计划失败${NC}"
    echo "响应: $create_response_body"
    ((TESTS_FAILED++))
fi

# 测试获取所有订阅计划
test_api "GET" "$API_BASE/subscription-plans" "" "200" "获取所有订阅计划"

# 测试仪表板API
echo -e "\n${YELLOW}=== 步骤5: 仪表板API测试 ===${NC}"

test_api "GET" "$API_BASE/dashboard/overview" "" "200" "获取仪表板概览"

test_api "GET" "$API_BASE/dashboard/organizations" "" "200" "获取仪表板组织列表"

test_api "GET" "$API_BASE/dashboard/database-status" "" "200" "获取数据库状态"

# 显示测试结果
echo -e "\n${YELLOW}=== 测试结果汇总 ===${NC}"
echo -e "通过: ${GREEN}$TESTS_PASSED${NC}"
echo -e "失败: ${RED}$TESTS_FAILED${NC}"
echo -e "总计: $((TESTS_PASSED + TESTS_FAILED))"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "\n${RED}❌ 有 $TESTS_FAILED 个测试失败${NC}"
    exit 1
fi