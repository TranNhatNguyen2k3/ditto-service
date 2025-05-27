#!/bin/bash

# Tạo policy
echo "Creating policy..."
curl -X PUT 'http://localhost:8080/api/2/policies/org.eclipse.ditto:device-1' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGl0dG86ZGl0dG8=' \
  -d @../config/policy.json

# Tạo thing
echo "Creating thing..."
curl -X PUT 'http://localhost:8080/api/2/things/org.eclipse.ditto:device-1' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGl0dG86ZGl0dG8=' \
  -d @../config/thing.json

# Tạo connection
echo "Creating connection..."
curl -X POST 'http://localhost:8080/api/2/connections' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Basic ZGl0dG86ZGl0dG8=' \
  -d @../config/connection.json

echo "Setup completed successfully" 