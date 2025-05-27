# Test WebSocket connection to Ditto
Write-Host "Testing WebSocket connection to Ditto..."

# Send test messages to Ditto
Write-Host "Sending test messages to Ditto..."

# Gửi 10 temperature events, value tăng dần
for ($i = 1; $i -le 10; $i++) {
    $temp = 40 + ($i * 3) # Giá trị tăng dần: 43, 46, ...
    Write-Host "Sending temperature event: $temp"
    $body = @{
        thingId = "org.eclipse.ditto:device-1"
        features = @{
            temperature = @{
                properties = @{
                    value = $temp
                    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ss.fffZ")
                }
            }
        }
    } | ConvertTo-Json -Depth 10

    Invoke-RestMethod -Method Put `
        -Uri "http://localhost:3001/api/v1/things/org.eclipse.ditto:device-1" `
        -Headers @{
            "Content-Type" = "application/json"
            "Authorization" = "Basic ZGl0dG86ZGl0dG8="
        } `
        -Body $body

    Start-Sleep -Seconds 5
}

# Gửi 10 humidity events, value tăng dần
for ($i = 1; $i -le 10; $i++) {
    $humidity = 30 + ($i * 4) # Giá trị tăng dần: 34, 38, ...
    Write-Host "Sending humidity event: $humidity"
    $body = @{
        thingId = "org.eclipse.ditto:device-1"
        features = @{
            humidity = @{
                properties = @{
                    value = $humidity
                    timestamp = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ss.fffZ")
                }
            }
        }
    } | ConvertTo-Json -Depth 10

    Invoke-RestMethod -Method Put `
        -Uri "http://localhost:3001/api/v1/things/org.eclipse.ditto:device-1" `
        -Headers @{
            "Content-Type" = "application/json"
            "Authorization" = "Basic ZGl0dG86ZGl0dG8="
        } `
        -Body $body

    Start-Sleep -Seconds 5
} 