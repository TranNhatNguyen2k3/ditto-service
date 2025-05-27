#!/bin/bash

# Tạo connection trong Ditto
curl -X POST 'http://localhost:8080/api/2/connections' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGl0dG86ZGl0dG8=' \
  -d @../config/ditto-connection.json

echo "Connection created successfully" 