# Ditto Service

A backend service for managing Ditto IoT platform integration.

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── app/               # Main application
│       └── main.go        # Application entry point
├── internal/              # Private application code
│   ├── api/              # API handlers and routes
│   │   ├── handler/      # HTTP handlers
│   │   ├── middleware/   # HTTP middleware
│   │   └── router/       # Route definitions
│   ├── config/           # Configuration
│   ├── domain/           # Business logic and domain models
│   │   ├── model/        # Domain models
│   │   ├── repository/   # Repository interfaces
│   │   └── service/      # Business logic
│   └── infrastructure/   # Infrastructure implementations
│       ├── ditto/        # Ditto client implementation
│       └── repository/   # Repository implementations
├── pkg/                   # Public libraries
│   ├── database/         # Database utilities
│   ├── logger/           # Logging utilities
│   └── utils/            # Common utilities
├── scripts/              # Build and deployment scripts
├── test/                 # Additional test files
├── .env                  # Environment variables
├── .gitignore           # Git ignore file
├── docker-compose.yml    # Docker compose configuration
├── go.mod               # Go module file
├── go.sum               # Go module checksum
└── Makefile             # Build automation
```

## Setup

1. Install dependencies:
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

- Build: `make build`
- Run: `make run`
- Test: `make test`
- Lint: `make lint`
- Format: `make fmt`

## API Documentation

API documentation is available at `/swagger/index.html` when the service is running.

## Docker

Build and run with Docker:
```bash
docker-compose up -d
``` 