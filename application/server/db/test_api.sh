#!/bin/bash

# API 测试脚本
# 用于测试注册和登录接口

BASE_URL="http://localhost:8888/api/v2"

echo "========== 医疗信息管理系统 API 测试 =========="
echo ""

# 测试 1: 注册新用户
echo "测试 1: 注册新用户"
echo "请求: POST $BASE_URL/register"
curl -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account_name": "医生李四",
    "username": "doctor_li",
    "password": "password123",
    "role": "医生",
    "operator": "admin"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

# 等待 1 秒
sleep 1

# 测试 2: 使用正确的用户名和密码登录
echo "测试 2: 使用正确的用户名和密码登录"
echo "请求: POST $BASE_URL/loginWithPassword"
curl -X POST "$BASE_URL/loginWithPassword" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "doctor_li",
    "password": "password123"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

# 等待 1 秒
sleep 1

# 测试 3: 使用错误的密码登录
echo "测试 3: 使用错误的密码登录（应该失败）"
echo "请求: POST $BASE_URL/loginWithPassword"
curl -X POST "$BASE_URL/loginWithPassword" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "doctor_li",
    "password": "wrongpassword"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

# 等待 1 秒
sleep 1

# 测试 4: 注册重复的用户名（应该失败）
echo "测试 4: 注册重复的用户名（应该失败）"
echo "请求: POST $BASE_URL/register"
curl -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d '{
    "account_name": "医生李四2",
    "username": "doctor_li",
    "password": "password456",
    "role": "医生",
    "operator": "admin"
  }' \
  -w "\nHTTP Status: %{http_code}\n\n"

echo "========== 测试完成 =========="
