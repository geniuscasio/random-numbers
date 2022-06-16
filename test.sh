#!/bin/bash

echo "register..."

curl -X GET 'localhost:8080/register' \
   -H 'Content-Type: application/json' \
   -d '{"name":"James May","email":"captain@slow.com", "pwd":"123456789"}'

echo "login..."

SESSION_ID=$(curl -v -X GET 'localhost:8080/login' \
   -H 'Content-Type: application/json' \
   -d '{"email":"captain@slow.com", "pwd":"123456789"}' 2>&1 | tr -d '\r' | sed -En 's/^< Session-Id: (.*)/\1/p'
   )

echo "session ID is $SESSION_ID"
echo "generate:"

curl -X GET 'localhost:8080/generate?sort=asc' \
   -H 'Content-Type: application/json' -H "Session-Id: $SESSION_ID" \
   -d '{"from_number":1,"to_number":100, "total_numbers":10}'

echo ""
echo "details:"

curl -X GET 'localhost:8080/details' \
   -H 'Content-Type: application/json' -H "Session-Id: $SESSION_ID"