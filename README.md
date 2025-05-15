# Service Architecture

## Docker Dependency

* docker run -p 3000:3000 -p 4317:4317 -p 4318:4318 --rm -ti grafana/otel-lgtm
* docker run -p 9092:9092 --rm -ti apache/kafka:latest

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

## Checklist

* ~~gin router~~
* ~~mariadb~~
* ~~kafka (watermill)~~
* sqs
* config(via aws param)
* ~~mongodb~~
* ~~cache (redis)~~
* ~~elastic~~
