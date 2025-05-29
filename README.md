# Ditto Service

A backend service for managing Ditto IoT platform integration. This service provides a robust API for managing IoT devices, collecting telemetry data, and integrating with the Ditto platform.

## Features

- RESTful API for device management
- Real-time telemetry data collection
- Integration with Ditto IoT platform
- InfluxDB time-series data storage
- Swagger API documentation
- Docker containerization
- Comprehensive test coverage

## Project Structure

```
.
├── bin/                   # Compiled binaries
├── cmd/                   # Application entry points
│   └── app/              # Main application
│       └── main.go       # Application entry point
├── config/               # Configuration files
│   └── config.go        # Configuration management
├── external/             # External dependencies and integrations
├── internal/             # Private application code
│   ├── app/             # Application core
│   │   └── module.go    # Application module definition
│   ├── ditto/           # Ditto client implementation
│   │   ├── client.go    # Ditto API client implementation
│   │   ├── service.go   # Ditto service layer
│   │   └── module.go    # Ditto module definition
│   ├── http/            # HTTP server and handlers
│   │   ├── handler/     # HTTP request handlers
│   │   └── router/      # Route definitions
│   ├── influxdb/        # InfluxDB integration
│   │   ├── client.go    # InfluxDB client implementation
│   │   └── module.go    # InfluxDB module definition
│   ├── middleware/      # HTTP middleware
│   │   ├── auth.go      # Authentication middleware
│   │   ├── cors.go      # CORS middleware
│   │   ├── error_handler.go # Error handling middleware
│   │   ├── jwt_auth.go  # JWT authentication middleware
│   │   ├── logging.go   # Request logging middleware
│   │   └── recover.go   # Panic recovery middleware
│   ├── model/           # Domain models
│   ├── repository/      # Repository implementations
│   │   ├── thing_repository.go      # Thing repository interface
│   │   ├── thing_repository_ditto.go # Ditto implementation
│   │   └── module.go    # Repository module definition
│   └── service/         # Business logic services
│       ├── thing_service.go # Thing service implementation
│       └── module.go    # Service module definition
├── pkg/                  # Public libraries
│   ├── constant/        # Constants and enums
│   ├── database/        # Database utilities
│   ├── errors/          # Error handling utilities
│   ├── graceful/        # Graceful shutdown utilities
│   ├── logger/          # Logging utilities
│   ├── swagger/         # Swagger documentation
│   ├── util/            # Common utilities
│   └── wrapper/         # Wrapper utilities
├── scripts/             # Build and deployment scripts
├── test/                # Additional test files
├── .gitignore          # Git ignore file
├── .mockery.yaml       # Mockery configuration
├── Dockerfile          # Docker build configuration
├── docker-compose.yml  # Docker compose configuration
├── go.mod             # Go module file
├── go.sum             # Go module checksum
└── Makefile           # Build automation
```

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for containerized deployment)
- Make (for build automation)
- InfluxDB (for time-series data storage)
- Ditto platform access credentials

## Configuration

### Environment Variables

The following environment variables need to be configured in `.env`:

```bash
# Server Configuration
PORT=3001
SERVER_URL=localhost
ENVIRONMENT=development
GIN_MODE=debug
PRODUCTION=false

# Ditto Platform
DITTO_URL=http://localhost:8080/api/2
DITTO_USERNAME=ditto
DITTO_PASSWORD=ditto
DITTO_WS_URL=ws://localhost:8080/ws/2

# Proxy Configuration
PROXY_AUTH_USERNAME=nguyen
PROXY_AUTH_PASSWORD=nguyen
PROXY_TARGET_URL=http://localhost:8080
PROXY_WS_URL=ws://localhost:8080

# InfluxDB
INFLUXDB_URL=http://influxdb:8086
INFLUXDB_TOKEN=your-token
INFLUXDB_ORG=your-org
INFLUXDB_BUCKET=your-bucket

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
SSL_MODE=disable

# JWT
JWT_SECRET=your-secret
JWT_EXPIRATION_TIME=24h
JWT_REFRESH_SECRET=your-refresh-secret
JWT_REFRESH_EXPIRATION_TIME=168h
```

## Setup

1. Install development tools:
```bash
make install-tools
```

2. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Run the service:
```bash
make run
```

## Development

### Build Commands

- Build: `make build`
- Run: `make run`


### Development Workflow

1. Create a new feature branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes and run tests:
```bash
make test
make lint
```

3. Generate mocks if needed:
```bash
make generate-mocks
```

4. Commit your changes:
```bash
git commit -m "feat: your feature description"
```

5. Push and create a pull request

### Testing

Run the test suite:
```bash
make test
```

Run specific test packages:
```bash
go test ./internal/...
```

## Docker Deployment

### Build and Run

Build and run with Docker:
```bash
docker-compose up -d
```

To stop the service:
```bash
docker-compose down
```

## API Documentation

### API Endpoints

#### Health Check
- `GET /health` - Health check endpoint (no authentication required)

#### Device Management
- `GET /api/devices` - List all devices with optional filtering
- `PUT /api/devices/:thingId` - Create or update a device
- `GET /api/devices/:thingId/state` - Get device state
- `PUT /api/devices/:thingId/features/:feature/command` - Send command to device feature
- `POST /api/devices/:thingId/features/:feature/command` - Send command to device feature

#### Ditto Integration
- `ANY /api/things/*path` - Proxy requests to Ditto API

### Authentication
All API endpoints (except `/health`) require Basic Authentication:
- Username: Configured via `PROXY_AUTH_USERNAME` environment variable
- Password: Configured via `PROXY_AUTH_PASSWORD` environment variable

### Example Usage

```bash
# Health check (no auth required)
curl http://localhost:3001/health

# List devices (auth required)
curl -u username:password http://localhost:3001/api/devices

# Create/Update device
curl -u username:password -X PUT http://localhost:3001/api/devices/device1 \
  -H "Content-Type: application/json" \
  -d '{"attributes": {"location": "room1"}}'

# Get device state
curl -u username:password http://localhost:3001/api/devices/device1/state

# Send command to device
curl -u username:password -X POST http://localhost:3001/api/devices/device1/features/temperature/command \
  -H "Content-Type: application/json" \
  -d '{"value": 25}'
```

## Contributing

1. Create a new branch for your feature
2. Make your changes
3. Run tests and ensure they pass
4. Update documentation if needed
5. Submit a pull request

### Code Style

- Follow Go best practices and idioms
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions small and focused
- Write unit tests for new features

## License

This project is proprietary and confidential. Unauthorized copying, distribution, or use is strictly prohibited. 