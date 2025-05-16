# Developer Guidelines

## Project Overview
This is a Go service built with a clean architecture approach, providing both REST API and gRPC endpoints, with Kafka for message broker functionality. The service includes OpenTelemetry for tracing and supports multiple data stores (MySQL/MariaDB, Redis, Elasticsearch, MongoDB).

## Directory Structure
```
├── app/                    # Application core
│   ├── controllers/        # API handlers (REST, gRPC, message broker)
│   ├── middlewares/        # HTTP middlewares
│   ├── models/             # Domain models
│   ├── repositories/       # Data access layer
│   └── usecases/           # Business logic
├── config/                 # Configuration definitions
├── container/              # Dependency injection
├── docs/                   # API documentation
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

## Setup Instructions
1. **Prerequisites**:
   - Go 1.16+
   - Docker and Docker Compose (for local development)
   - MariaDB/MySQL
   - Redis
   - Elasticsearch
   - Kafka
   - OpenTelemetry collector

2. **Installation**:
   ```bash
   # Clone the repository
   git clone <repository-url>
   
   # Install dependencies
   go mod download
   
   # Start required services using Docker Compose
   docker-compose up -d
   
   # Run the application
   go run main.go
   ```

## Development Workflow
1. **Making Changes**:
   - Follow the clean architecture principles
   - Add new functionality in the appropriate layer:
     - Models: Define domain entities
     - Repositories: Implement data access
     - Usecases: Implement business logic
     - Controllers: Handle API requests
     - Routes: Define API endpoints

2. **Adding New Features**:
   - Create interfaces in the appropriate layer
   - Implement the interfaces
   - Register in the container (container/container.go)
   - Add routes if needed

3. **Configuration**:
   - Default configuration is in config/config.go
   - Override with environment variables in production

## Testing
Currently, the project doesn't have automated tests. When implementing tests:
- Use Go's standard testing package
- Create *_test.go files alongside the code being tested
- Follow table-driven testing approach for comprehensive test cases
- Mock external dependencies using interfaces

## Best Practices
1. **Code Organization**:
   - Follow clean architecture principles
   - Keep layers separated with clear dependencies
   - Use interfaces for dependency injection

2. **Error Handling**:
   - Return errors rather than panicking
   - Use meaningful error messages
   - Log errors appropriately

3. **Logging**:
   - Use the provided logger package
   - Include context in log messages
   - Use appropriate log levels

4. **Performance**:
   - Use connection pooling for databases
   - Implement caching where appropriate
   - Consider using goroutines for concurrent operations

5. **Security**:
   - Validate all user input
   - Use prepared statements for database queries
   - Implement proper authentication and authorization

## Running in Production
1. **Configuration**:
   - Override default configuration with environment variables
   - Use secure credentials

2. **Monitoring**:
   - Use OpenTelemetry for tracing
   - Implement health checks
   - Set up alerts for critical errors

3. **Deployment**:
   - Use containerization (Docker)
   - Consider orchestration with Kubernetes
   - Implement CI/CD pipelines