Architecture

This project follows a pragmatic layered architecture focused on simplicity, performance, and maintainability. The structure is intentionally lightweight and avoids unnecessary abstractions that add complexity without clear benefits.

Rather than strictly applying Clean Architecture, the design favors Go idioms: clear responsibilities, concrete types, and minimal indirection. This keeps the codebase easy to navigate and efficient to run.

The application is organized by responsibility:

```
/cmd/api            Application entry point
/internal/http      HTTP handlers and routing
/internal/service   Business logic and orchestration
/internal/repo      Database access and persistence
/internal/model     Shared domain data structures
/internal/db        Database initialization
```

Each layer has a single purpose. HTTP handlers deal only with request parsing and response formatting. Services contain business rules and validation. Repositories handle persistence and database queries. Models represent core application data and are shared across layers.

A single model is often used for database access, domain logic, and HTTP responses. When data structures are identical across layers, duplicating them provides no value and only adds mapping overhead. Separate DTOs are introduced only when fields differ, data must be hidden or transformed, or external contracts require it.

This approach reduces boilerplate, avoids unnecessary data copying, and keeps request handling fast and easy to reason about. It also makes the codebase easier to refactor as requirements evolve.

If the domain becomes more complex over time, this structure can naturally evolve toward more advanced architectural styles (such as vertical slices or hexagonal architecture) without requiring a rewrite.

The guiding principle is simple: keep things straightforward until additional structure is clearly needed.