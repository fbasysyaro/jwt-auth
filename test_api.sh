#!/bin/bash

# API Testing Script
BASE_URL="http://localhost:8080"

echo "=== JWT Auth API Testing ==="
echo

# 1. Health Check
echo "1. Testing Health Check..."
curl -s -X GET $BASE_URL/health | jq .
echo

# 2. Register User
echo "2. Testing User Registration..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')
echo $REGISTER_RESPONSE | jq .
echo

# 3. Login User
echo "3. Testing User Login..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }')
echo $LOGIN_RESPONSE | jq .

# Extract tokens
ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.access_token')
REFRESH_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.refresh_token')
echo "Access Token: $ACCESS_TOKEN"
echo

# 4. Test Protected Route - Profile
echo "4. Testing Protected Route - Profile..."
curl -s -X GET $BASE_URL/api/v1/profile \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo

# 5. Test Protected Route - Dashboard
echo "5. Testing Protected Route - Dashboard..."
curl -s -X GET $BASE_URL/api/v1/dashboard \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo

# 6. Test Token Refresh
echo "6. Testing Token Refresh..."
REFRESH_RESPONSE=$(curl -s -X POST $BASE_URL/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}")
echo $REFRESH_RESPONSE | jq .
echo

# 7. Test Logout
echo "7. Testing Logout..."
curl -s -X POST $BASE_URL/api/v1/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo

# 8. Test Access After Logout (Should Fail)
echo "8. Testing Access After Logout (Should Fail)..."
curl -s -X GET $BASE_URL/api/v1/profile \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .
echo

echo "=== Testing Complete ==="