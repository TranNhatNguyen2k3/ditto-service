#!/bin/bash

# Test WebSocket connection to Ditto
echo "Testing WebSocket connection to Ditto..."
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Host: localhost:8080" \
     -H "Origin: http://localhost:8080" \
     -H "Authorization: Basic ZGl0dG86ZGl0dG8=" \
     http://localhost:8080/ws/2 &

# Wait for WebSocket connection
sleep 2

# Send test messages to Ditto
echo "Sending test messages to Ditto..."

# Send temperature events
for i in {1..5}; do
    temp=$((40 + RANDOM % 30))
    echo "Sending temperature event: $temp"
    curl -X PUT \
         -H "Content-Type: application/json" \
         -H "Authorization: Basic ZGl0dG86ZGl0dG8=" \
         -d "{\"temperature\": $temp, \"humidity\": 45}" \
         http://localhost:8080/api/2/things/org.eclipse.ditto:device-1/features/temperature/properties
    sleep 2
done

# Send humidity events
for i in {1..5}; do
    humidity=$((30 + RANDOM % 40))
    echo "Sending humidity event: $humidity"
    curl -X PUT \
         -H "Content-Type: application/json" \
         -H "Authorization: Basic ZGl0dG86ZGl0dG8=" \
         -d "{\"temperature\": 25, \"humidity\": $humidity}" \
         http://localhost:8080/api/2/things/org.eclipse.ditto:device-1/features/humidity/properties
    sleep 2
done

# Stop WebSocket connection
pkill -f "curl.*websocket" 