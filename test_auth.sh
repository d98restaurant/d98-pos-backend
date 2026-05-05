#!/bin/bash

BASE_URL="https://pos-backend-h2c3.onrender.com"

echo "=== Testing Authentication ==="

# Register a test user
echo -e "\n1. Registering test user..."
REGISTER_RESPONSE=$(curl -s -X POST $BASE_URL/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demouser",
    "email": "demo@pos.com",
    "password": "demo123",
    "role": "admin"
  }')
echo $REGISTER_RESPONSE

# Login with the test user
echo -e "\n2. Logging in with test user..."
LOGIN_RESPONSE=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "demouser", "password": "demo123"}')
echo $LOGIN_RESPONSE

# Extract token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -n "$TOKEN" ]; then
  echo -e "\n✅ Login successful! Token: ${TOKEN:0:50}..."
  
  # Test authenticated endpoint
  echo -e "\n3. Testing authenticated endpoint..."
  curl -s -X GET $BASE_URL/api/menu \
    -H "Authorization: Bearer $TOKEN"
else
  echo -e "\n❌ Login failed"
fi
