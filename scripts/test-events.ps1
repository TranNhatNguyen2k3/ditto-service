# Test script to send events to Ditto
Write-Host "Sending test events to Ditto..."

# Send temperature events
for ($i = 1; $i -le 10; $i++) {
    $temp = 55 + ($i * 2) # Temperature from 22 to 40
    Write-Host "Sending temperature event: $temp"
    $body = @{
        thingId = "org.eclipse.ditto:test-device"
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
        -Uri "http://localhost:8080/api/2/things/org.eclipse.ditto:test-device" `
        -Headers @{
            "Content-Type" = "application/json"
            "Authorization" = "Basic ZGl0dG86ZGl0dG8="
        } `
        -Body $body

    Start-Sleep -Seconds 2
}

# Send humidity events
for ($i = 1; $i -le 10; $i++) {
    $humidity = 50 + ($i * 3) # Humidity from 43 to 70
    Write-Host "Sending humidity event: $humidity"
    $body = @{
        thingId = "org.eclipse.ditto:test-device"
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
        -Uri "http://localhost:8080/api/2/things/org.eclipse.ditto:test-device" `
        -Headers @{
            "Content-Type" = "application/json"
            "Authorization" = "Basic ZGl0dG86ZGl0dG8="
        } `
        -Body $body

    Start-Sleep -Seconds 2
}

Write-Host "Test events sent successfully!" 