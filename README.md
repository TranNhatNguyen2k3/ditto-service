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
│   ├── grafana/         # Grafana dashboard configurations
│   └── config.go        # Configuration management
├── external/             # External dependencies and integrations
├── internal/             # Private application code
│   ├── app/             # Application core
│   │   └── module.go    # Application module definition
│   ├── config/          # Configuration management
│   ├── ditto/           # Ditto client implementation
│   │   ├── client.go    # Ditto API client implementation
│   │   ├── service.go   # Ditto service layer
│   │   └── module.go    # Ditto module definition
│   ├── http/            # HTTP server and handlers
│   │   ├── dto/         # Data Transfer Objects
│   │   ├── handler/     # HTTP request handlers
│   │   └── router/      # Route definitions
│   ├── influxdb/        # InfluxDB integration
│   │   ├── client.go    # InfluxDB client implementation
│   │   └── module.go    # InfluxDB module definition
│   ├── middleware/      # HTTP middleware
│   │   ├── cors.go      # CORS middleware
│   │   ├── error_handler.go # Error handling middleware
│   │   ├── jwt_auth.go  # JWT authentication middleware
│   │   ├── logging.go   # Request logging middleware
│   │   └── recover.go   # Panic recovery middleware
│   ├── model/           # Domain models
│   │   ├── entity/      # Domain entities
│   │   └── thing.go     # Thing model definition
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

The following environment variables need to be configured in `config/.env`:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Ditto Platform
DITTO_API_URL=https://ditto.example.com
DITTO_API_KEY=your-api-key
DITTO_TENANT=your-tenant

# InfluxDB
INFLUXDB_URL=http://influxdb:8086
INFLUXDB_TOKEN=your-token
INFLUXDB_ORG=your-org
INFLUXDB_BUCKET=your-bucket

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Application Configuration

The application configuration is defined in `config/config.yaml`:

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 5s
  write_timeout: 10s

ditto:
  api_url: "https://ditto.example.com"
  timeout: 30s
  retry_attempts: 3

influxdb:
  url: "http://influxdb:8086"
  timeout: 5s
  batch_size: 1000
```

## Setup

1. Install development tools:
```bash
make install-tools
```

2. Set up environment variables:
```bash
cp config/.env.example config/.env
# Edit config/.env with your configuration
```

3. Run the service:
```bash
make run
```

## Development

### Build Commands

- Build: `make build`
- Run: `make run`
- Test: `make test`
- Lint: `make lint`
- Format: `make fmt`
- Generate mocks: `make generate-mocks`

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

Run integration tests:
```bash
go test ./test/integration/...
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

### Docker Compose Services

- `ditto-service`: Main application service
- `influxdb`: Time-series database
- `grafana`: Metrics visualization (optional)

## API Documentation

API documentation is available at `/swagger/index.html` when the service is running.

### Example API Endpoints

- `GET /api/v1/devices`: List all devices
- `POST /api/v1/devices`: Create a new device
- `GET /api/v1/devices/{id}`: Get device details
- `GET /api/v1/telemetry`: Query telemetry data
- `POST /api/v1/telemetry`: Submit telemetry data

## Monitoring

The service exposes the following monitoring endpoints:

- `/metrics`: Prometheus metrics
- `/health`: Health check endpoint
- `/ready`: Readiness probe endpoint

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