# Service Architecture

## Docker Dependency

* docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
* docker run -p 9092:9092 --rm -ti apache/kafka:latest

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

## Architecture Diagram

```
+--------------------------------------------------------------------------------------------------+
|                                           CONTAINER                                              |
+--------------------------------------------------------------------------------------------------+
                |                    |                    |                    |
                v                    v                    v                    v
+---------------------------+  +-------------+  +------------------+  +------------------+
|        CONTROLLERS        |  |   USECASES  |  |   REPOSITORIES   |  |  INFRASTRUCTURE  |
|---------------------------|  |-------------|  |------------------|  |------------------|
| - REST API                |  | - User      |  | - UserDB         |  | - Config         |
|   * UserHandler           |  |   * List    |  |   * MariaDB/MySQL|  | - ORM (DB)       |
|                           |  |   * Register|  |   * PostgreSQL   |  |   * Transactor   |
| - gRPC                    |  |             |  | - UserElastic    |  | - Cache          |
|                           |  |             |  | - UserMongo      |  |   * Redis        |
| - Message Broker          |  |             |  |                  |  | - Elasticsearch  |
|   * Kafka                 |  |             |  |                  |  | - MongoDB        |
+------------+--------------+  +------^------+  +--------^---------+  | - Kafka          |
             |                        |                  |            |   * Publisher     |
             |                        |                  |            |   * Subscriber    |
             |                        |                  |            | - OpenTelemetry   |
             +------------------------+------------------+            | - HTTP/gRPC Server|
                                                                     +------------------+
```

## Detailed Architectural Dependency Map

```
                                  +----------------+
                                  |     main.go    |
                                  +-------+--------+
                                          |
                                          v
+--------------------------------------------------------------------------------------------------+
|                                      CONTAINER                                                   |
|                                                                                                  |
|  +----------------+  +----------------+  +----------------+  +----------------+  +---------------+|
|  |    Config      |  |  Infrastructure|  |  Repositories  |  |    Usecases    |  |  Controllers  ||
|  |----------------|  |----------------|  |----------------|  |----------------|  |---------------||
|  |- App           |  |- ORM           |  |- UserDB        |  |- UserUsecase   |  |- REST API     ||
|  |- Database      |  |  * MySQL       |  |  * List        |  |  * List        |  |  * UserHandler||
|  |- Cache         |  |  * PostgreSQL  |  |  * Create      |  |  * Register    |  |- gRPC         ||
|  |- Elastic       |  |- Cache         |  |- UserElastic   |  |                |  |- Broker       ||
|  |- MongoDB       |  |  * Redis       |  |- UserMongo     |  |                |  |  * UserHandler||
|  |- Kafka         |  |- Elastic       |  |                |  |                |  |               ||
|  |- GRPC          |  |- MongoDB       |  |                |  |                |  |               ||
|  |- REST          |  |- Kafka         |  |                |  |                |  |               ||
|  |- OTEL          |  |- OTEL          |  |                |  |                |  |               ||
|  +----------------+  +----------------+  +----------------+  +----------------+  +---------------+|
|                        ^      ^   ^             ^                   ^                  ^          |
|                        |      |   |             |                   |                  |          |
+------------------------|------|---|-------------|-------------------|------------------|----------+
                         |      |   |             |                   |                  |
                         |      |   |             |                   |                  |
+----------------------+ |      |   |             | +----------------+|          +------|-----------+
|      pkg/            | |      |   |             | |    app/        ||          |     app/        |
|                      | |      |   |             | |                ||          |                  |
| +------------------+ | |      |   |             | | +------------+ ||          | +-------------+  |
| |   datastore/     | | |      |   |             | | |repositories| ||          | | controllers/|  |
| |                  | | |      |   |             | | |            | ||          | |             |  |
| | +-------------+  | | |      |   |             | | | +---------+| ||          | | +---------+ |  |
| | |    orm/     |<-|-+ |      |   |             | | | |  user/  || ||          | | | restapi/ | |  |
| | +-------------+  | |        |   |             | | | +---------+| ||          | | +---------+ |  |
| |                  | |        |   |             | | +------------+ ||          | |             |  |
| | +-------------+  | |        |   |             | |                ||          | | +---------+ |  |
| | |   elastic/  |<-|-------+  |   |             | | +------------+ ||          | | |  grpc/  | |  |
| | +-------------+  | |     |  |   |             | | |  usecases/ | ||          | | +---------+ |  |
| |                  | |     |  |   |             | | |            | ||          | |             |  |
| | +-------------+  | |     |  |   |             | | | +---------+| ||          | | +---------+ |  |
| | |  mongodb/   |<-|--+    |  |   |             | | | |  user/  |<-+----------+ | | broker/ | |  |
| | +-------------+  | | |   |  |   |             | | | +---------+| |           | | +---------+ |  |
| +------------------+ | |   |  |   |             | | +------------+ |           | +-------------+  |
|                      | |   |  |   |             | |                |           |                  |
| +------------------+ | |   |  |   |             | | +------------+ |           | +-------------+  |
| |     cache/       |<|-+   |  |   |             | | |   models/  | |           | | middlewares/|  |
| +------------------+ |     |  |   |             | | +------------+ |           | +-------------+  |
|                      |     |  |   |             | +----------------+           |                  |
| +------------------+ |     |  |   |             |                              +------------------+
| | message_broker/  | |     |  |   |             |
| |                  | |     |  |   |             |                              +------------------+
| | +-------------+  | |     |  |   |             |                              |     routes/      |
| | |   kafka/    |<-|-----+    |   |             |                              |                  |
| | +-------------+  | |   |    |   |             |                              | +-------------+  |
| +------------------+ |   |    |   |             |                              | |    api/     |  |
|                      |   |    |   |             |                              | |             |  |
| +------------------+ |   |    |   |             |                              | | +---------+ |  |
| |     server/      |<|---+    |   |             |                              | | |  user/  | |  |
| +------------------+ |        |   |             |                              | | +---------+ |  |
|                      |        |   |             |                              | |             |  |
| +------------------+ |        |   |             |                              | | +---------+ |  |
| |     otel/        |<|--------+   |             |                              | | |permission| |  |
| +------------------+ |            |             |                              | | +---------+ |  |
|                      |            |             |                              | +-------------+  |
| +------------------+ |            |             |                              |                  |
| |    setting/      |<|------------+             |                              | +-------------+  |
| +------------------+ |                          |                              | |   broker/   |  |
|                      |                          |                              | |             |  |
| +------------------+ |                          |                              | | +---------+ |  |
| |    utilities/    | |                          |                              | | |  user/  | |  |
| +------------------+ |                          |                              | | +---------+ |  |
+----------------------+                          |                              | +-------------+  |
                                                  |                              +------------------+
+------------------------------------------+     |
|               config/                    |     |
|                                          |     |
| +----------------+ +------------------+  |     |
| |    config.go   | |    database.go   |<-------+
| +----------------+ +------------------+  |
|                                          |
| +----------------+ +------------------+  |
| |     app.go     | |     cache.go     |  |
| +----------------+ +------------------+  |
|                                          |
| +----------------+ +------------------+  |
| |    grpc.go     | |     rest.go      |  |
| +----------------+ +------------------+  |
|                                          |
| +----------------+ +------------------+  |
| |    kafka.go    | |    Elastic.go    |  |
| +----------------+ +------------------+  |
|                                          |
| +----------------+ +------------------+  |
| |   mongodb.go   | |     otel.go      |  |
| +----------------+ +------------------+  |
|                                          |
| +----------------+                       |
| |   setting.go   |                       |
| +----------------+                       |
+------------------------------------------+
```

## Module Dependencies

```
main.go
  └── container/
       ├── config/
       ├── pkg/datastore/orm/
       ├── pkg/cache/
       ├── pkg/datastore/elastic/
       ├── pkg/datastore/mongodb/
       ├── app/repositories/
       │    └── app/repositories/user/
       ├── app/usecases/
       │    └── app/usecases/user/
       ├── app/controllers/restapi/
       │    └── app/controllers/restapi/user/
       ├── app/controllers/grpc/
       ├── app/controllers/broker/
       ├── pkg/message_broker/kafka/
       ├── pkg/otel/
       └── pkg/server/
```

## Flow of Dependencies

1. **Controllers** depend on **Usecases**
2. **Usecases** depend on **Repositories**
3. **Repositories** depend on **Infrastructure** components
4. **Container** wires everything together

## Component Descriptions

- **Controllers**: Handle external requests (HTTP, gRPC, message broker)
- **Usecases**: Implement business logic
- **Repositories**: Provide data access abstraction
- **Infrastructure**: Provide technical capabilities (database, cache, messaging, etc.)

## Architecture Pattern

This service follows the **Clean Architecture** pattern with the following characteristics:

1. **Dependency Rule**: Dependencies only point inward. Inner layers don't know about outer layers.
   - Controllers → Usecases → Repositories → Infrastructure

2. **Separation of Concerns**:
   - **Controllers**: Handle input/output and protocol-specific concerns
   - **Usecases**: Contain application-specific business rules
   - **Repositories**: Abstract data access
   - **Infrastructure**: Implement technical details

3. **Dependency Injection**: The Container package wires everything together, injecting dependencies from outer layers into inner layers.

4. **Testability**: Each layer can be tested independently by mocking its dependencies.

5. **Framework Independence**: The core business logic (Usecases) doesn't depend on any external frameworks.

```
                  Dependency Direction
                         │
                         │
                         ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │     │             │
│ Controllers │────▶│  Usecases   │────▶│Repositories │────▶│Infrastructure│
│             │     │             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
     │                                                             ▲
     │                                                             │
     │                                                             │
     └─────────────────────────────────────────────────────────────┘
                  Dependency Inversion (via interfaces)

┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│                           Container                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
                  Wires all dependencies together
```

## Checklist

* ~~gin router~~
* ~~mariadb~~
* ~~kafka (watermill)~~
* sqs
* config(via aws param)
* ~~mongodb~~
* ~~cache (redis)~~
* ~~elastic~~
