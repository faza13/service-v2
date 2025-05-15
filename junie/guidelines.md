# Developer Guidelines

## Project Overview
This service follows the **Clean Architecture** pattern, organizing code into layers with dependencies pointing inward:
- **Controllers**: Handle external requests (HTTP, gRPC, Kafka)
- **Usecases**: Implement business logic
- **Repositories**: Provide data access abstraction
- **Infrastructure**: Provide technical capabilities (database, cache, messaging, etc.)

## Directory Structure
```
service/
├── app/                    # Application code
│   ├── controllers/        # Handle external requests
│   │   ├── broker/         # Message broker handlers
│   │   ├── grpc/           # gRPC handlers
│   │   └── restapi/        # REST API handlers
│   ├── middlewares/        # HTTP middleware
│   ├── models/             # Domain models
│   ├── repositories/       # Data access layer
│   └── usecases/           # Business logic
├── config/                 # Configuration
├── container/              # Dependency injection
├── pkg/                    # Shared packages
│   ├── cache/              # Cache implementation
│   ├── datastore/          # Database implementations
│   ├── logger/             # Logging utilities
│   ├── message_broker/     # Message broker implementation
│   ├── otel/               # OpenTelemetry integration
│   ├── server/             # HTTP/gRPC server
│   ├── setting/            # Application settings
│   └── utilities/          # Utility functions
└── routes/                 # API route definitions
```

## Setup and Installation

### Prerequisites
- Go 1.24+
- Docker and Docker Compose

### Local Development
1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Start required services:
   ```
   docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
   docker run -p 9092:9092 --rm -ti apache/kafka:latest
   ```
4. Run the application:
   ```
   go run main.go
   ```

### Docker Deployment
To run the entire stack with Docker Compose:
```
docker-compose -f docker-composer.yml up
```

## Running Tests
Currently, the project does not have automated tests. When adding new features, consider adding tests to ensure code quality and prevent regressions.

## Best Practices

### Code Organization
1. **Follow Clean Architecture principles**:
   - Keep dependencies pointing inward
   - Don't let inner layers know about outer layers
   - Use interfaces for dependency inversion

2. **Package Structure**:
   - Group related functionality in packages
   - Keep packages focused on a single responsibility
   - Use meaningful names that reflect the package's purpose

### Coding Standards
1. **Go Conventions**:
   - Follow standard Go formatting (use `gofmt` or `go fmt`)
   - Use meaningful variable and function names
   - Write comments for exported functions and types

2. **Error Handling**:
   - Always check errors and handle them appropriately
   - Use custom error types for domain-specific errors
   - Avoid using panic in production code

3. **Dependency Injection**:
   - Use the container package for wiring dependencies
   - Avoid global state and singletons
   - Design for testability

### Adding New Features
1. Identify the appropriate layer for your code
2. Create interfaces for dependencies
3. Implement the feature following Clean Architecture principles
4. Wire it up in the container package
5. Update routes if necessary

## Monitoring and Observability
The service uses OpenTelemetry for distributed tracing and monitoring. Traces are sent to Jaeger, which can be accessed at http://localhost:16686 when running with Docker Compose.