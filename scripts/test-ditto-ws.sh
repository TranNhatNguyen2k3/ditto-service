#!/bin/bash

# Test kết nối WebSocket tới Ditto
echo "Testing WebSocket connection to Ditto..."
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Host: localhost:8080" \
     -H "Origin: http://localhost:8080" \
     -H "Authorization: Basic ZGl0dG86ZGl0dG8=" \
     http://localhost:8080/ws/2 &

# Đợi kết nối WebSocket được thiết lập
sleep 2

# Gửi test message tới Ditto
echo "Sending test message to Ditto..."
curl -X PUT \
     -H "Content-Type: application/json" \
     -H "Authorization: Basic ZGl0dG86ZGl0dG8=" \
     -d '{"temperature": 60, "humidity": 45}' \
     http://localhost:8080/api/2/things/org.eclipse.ditto:device-1/features/temperature/properties

# Đợi message được xử lý
sleep 2

# Kiểm tra logs của Ditto
echo "Checking Ditto logs..."
docker logs ditto-ditto-1 | grep "Received message"

# Dừng WebSocket connection
pkill -f "curl.*websocket" 