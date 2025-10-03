# TRex Architecture ASCII Diagram

```
╔══════════════════════════════════════════════════════════════════════════════════╗
║                            🦕 TRex Microservice Architecture                      ║
╚══════════════════════════════════════════════════════════════════════════════════╝

     HTTP Requests                           ┌─────────────────┐
         │                                   │   OpenAPI       │
         ▼                                   │   Generator     │
┌─────────────────┐                         │   & Docs        │
│   Router        │◄────────────────────────┤                 │
│   (Gorilla Mux) │                         └─────────────────┘
└─────────────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Middleware    │     │   Auth JWT      │     │   Metrics       │
│   Stack         │────►│   Validation    │────►│   & Logging     │
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │
         ▼
╔═════════════════════════════════════════════════════════════════════════════════╗
║                              API LAYER (pkg/handlers/)                          ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      ║
║  │  Dinosaur   │    │   Generic   │    │   Events    │    │  [Entity]   │      ║
║  │  Handlers   │    │  Handlers   │    │  Handlers   │    │  Handlers   │      ║
║  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘      ║
╚═════════════════════════════════════════════════════════════════════════════════╝
         │                    │                    │                    │
         ▼                    ▼                    ▼                    ▼
╔═════════════════════════════════════════════════════════════════════════════════╗
║                        SERVICE LAYER (pkg/services/)                            ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      ║
║  │  Dinosaur   │    │   Generic   │    │   Events    │    │  [Entity]   │      ║
║  │  Service    │    │  Service    │    │  Service    │    │  Service    │      ║
║  │             │    │             │    │             │    │             │      ║
║  │ • Business  │    │ • Auth      │    │ • Publish   │    │ • Custom    │      ║
║  │   Logic     │    │ • Generic   │    │ • Subscribe │    │   Logic     │      ║
║  │ • Events    │    │   CRUD      │    │ • Process   │    │ • Events    │      ║
║  │ • Locking   │    │             │    │             │    │             │      ║
║  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘      ║
╚═════════════════════════════════════════════════════════════════════════════════╝
         │                    │                    │                    │
         ▼                    ▼                    ▼                    ▼
╔═════════════════════════════════════════════════════════════════════════════════╗
║                          DAO LAYER (pkg/dao/)                                   ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐    ┌─────────────┐      ║
║  │  Dinosaur   │    │   Generic   │    │   Events    │    │  [Entity]   │      ║
║  │    DAO      │    │    DAO      │    │    DAO      │    │    DAO      │      ║
║  │             │    │             │    │             │    │             │      ║
║  │ • CRUD Ops  │    │ • Base CRUD │    │ • Event     │    │ • Entity    │      ║
║  │ • Queries   │    │ • Search    │    │   Storage   │    │   Queries   │      ║
║  │ • Joins     │    │ • Filters   │    │             │    │             │      ║
║  └─────────────┘    └─────────────┘    └─────────────┘    └─────────────┘      ║
╚═════════════════════════════════════════════════════════════════════════════════╝
         │                    │                    │                    │
         ▼                    ▼                    ▼                    ▼
╔═════════════════════════════════════════════════════════════════════════════════╗
║                       DATABASE LAYER (pkg/db/)                                  ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║           ┌─────────────────┐               ┌─────────────────┐                 ║
║           │   Session       │               │   Lock          │                 ║
║           │   Factory       │               │   Factory       │                 ║
║           │                 │               │                 │                 ║
║           │ • Connections   │               │ • Advisory      │                 ║
║           │ • Transactions  │               │   Locks         │                 ║
║           │ • Pooling       │               │ • Concurrency   │                 ║
║           └─────────────────┘               └─────────────────┘                 ║
╚═════════════════════════════════════════════════════════════════════════════════╝
                               │
                               ▼
                   ┌─────────────────────┐
                   │    PostgreSQL       │
                   │                     │
                   │ • Tables & Indexes  │
                   │ • Migrations        │
                   │ • NOTIFY/LISTEN     │
                   │ • Advisory Locks    │
                   └─────────────────────┘

╔═════════════════════════════════════════════════════════════════════════════════╗
║                      DEPENDENCY INJECTION FRAMEWORK                             ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║                        Environment (cmd/trex/environments/)                     ║
║                                                                                 ║
║    ┌─────────────────────────────────────────────────────────────────────┐     ║
║    │                        Service Locators                             │     ║
║    │                                                                     │     ║
║    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │     ║
║    │  │  Dinosaur   │  │   Generic   │  │   Events    │  │  [Entity]   │ │     ║
║    │  │  Locator    │  │  Locator    │  │  Locator    │  │  Locator    │ │     ║
║    │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │     ║
║    └─────────────────────────────────────────────────────────────────────┘     ║
║                                     │                                           ║
║    ┌─────────────────────────────────────────────────────────────────────┐     ║
║    │                    Infrastructure                                   │     ║
║    │                                                                     │     ║
║    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │     ║
║    │  │  Database   │  │  Handlers   │  │   Clients   │  │   Config    │ │     ║
║    │  │             │  │             │  │             │  │             │ │     ║
║    │  │ • Sessions  │  │ • Auth      │  │ • OCM       │  │ • Env Vars  │ │     ║
║    │  │ • Locks     │  │ • Middleware│  │ • External  │  │ • Secrets   │ │     ║
║    │  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘ │     ║
║    └─────────────────────────────────────────────────────────────────────┘     ║
╚═════════════════════════════════════════════════════════════════════════════════╝

╔═════════════════════════════════════════════════════════════════════════════════╗
║                           EVENT-DRIVEN ARCHITECTURE                             ║
╠═════════════════════════════════════════════════════════════════════════════════╣
║                                                                                 ║
║  Services ───► Event Bus ───► Controllers ───► Async Processing                 ║
║      │            │               │                    │                        ║
║      │            │               │                    ▼                        ║
║      ▼            ▼               ▼         ┌─────────────────┐                 ║
║  • Create      • Publish      • Listen      │   External      │                 ║
║  • Update      • Subscribe    • Process     │   Integrations  │                 ║
║  • Delete      • Route        • React      │                 │                 ║
║                                             │ • Notifications │                 ║
║                                             │ • Webhooks      │                 ║
║                                             │ • Auditing      │                 ║
║                                             └─────────────────┘                 ║
╚═════════════════════════════════════════════════════════════════════════════════╝

                        ┌─────────────────────────────────┐
                        │        CODE GENERATION          │
                        │                                 │
                        │  go run scripts/generator.go    │
                        │  --kind EntityName              │
                        │                                 │
                        │  Generates:                     │
                        │  • API Models                   │
                        │  • Services                     │
                        │  • DAOs                         │
                        │  • Handlers                     │
                        │  • OpenAPI Specs                │
                        │  • Service Locators             │
                        │  • Migrations                   │
                        │  • Tests                        │
                        └─────────────────────────────────┘
```

## Architecture Overview

TRex is a layered microservice architecture with clear separation of concerns:

- **API Layer**: HTTP handlers and routing with OpenAPI specification
- **Service Layer**: Business logic, event publishing, and transaction management  
- **DAO Layer**: Data access abstraction with GORM ORM
- **Database Layer**: PostgreSQL with migrations and advisory locking
- **Dependency Injection**: Service locator pattern for clean dependency management
- **Event-Driven**: Async processing via PostgreSQL NOTIFY/LISTEN
- **Code Generation**: Full CRUD scaffolding generator for rapid development

The architecture follows hexagonal/clean architecture principles with ports and adapters, enabling testability and maintainability at scale.