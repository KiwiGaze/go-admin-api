# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-admin is a backend admin panel framework built with Gin + GORM, providing RBAC access control via Casbin, JWT authentication, code generation, and a multi-command CLI (Cobra). It pairs with a separate Vue frontend (go-admin-ui). The project is Chinese-authored; comments, error messages, and config comments are in Chinese.

## Commands

```bash
# Build
make build                    # CGO_ENABLED=0 production binary
make build-sqlite             # Build with SQLite3 support (requires -tags sqlite3)

# Run (development)
go run main.go server -c config/settings.yml

# Run (Docker)
make run                      # docker-compose up
make stop                     # docker-compose down

# Database migration
go run main.go migrate -c config/settings.yml

# Generate migration file (local/custom)
go run main.go migrate -c config/settings.yml -g

# Generate migration file (go-admin core)
go run main.go migrate -c config/settings.yml -g -a

# Scaffold a new app module
go run main.go app -n <app-name>

# Swagger docs generation
go generate ./...
# or directly:
swag init --parseDependency --parseDepth=6 --instanceName admin -o ./docs/admin

# API check mode (validates registered routes against sys_api table)
go run main.go server -c config/settings.yml -a
```

## Architecture

### CLI Commands (cmd/)

Cobra-based multi-command structure. Entry: `main.go` -> `cmd.Execute()`.

| Command | Purpose |
|---------|---------|
| `server` | Start HTTP API server |
| `migrate` | Run database migrations / generate migration files |
| `version` | Print version |
| `config` | Print loaded config |
| `app` | Scaffold a new app module |

### App Modules (app/)

Each module is self-contained with a consistent internal structure:

```
app/<module>/
  apis/        # HTTP handlers (embed api.Api from go-admin-core)
  models/      # GORM models
  router/      # Gin route registration
  service/     # Business logic
  service/dto/ # Request/response DTOs with search tags
```

Three built-in modules:
- **admin** - Core RBAC: users, roles, menus, depts, posts, dicts, configs, logs
- **jobs** - Cron job management (robfig/cron)
- **other** - Code generation tools, file upload, server monitoring

### Router Registration Pattern

Routes register via `init()` functions that append to package-level slices. In `cmd/api/`, each module's init file appends its `router.InitRouter` to `AppRouters`. Within a module's router package:

- `routerNoCheckRole` - public routes (no auth)
- `routerCheckRole` - authenticated routes (JWT + Casbin)

Each resource router file uses `init()` to append its registration function.

### Request Flow (for authenticated endpoints)

```
Gin Engine
  -> Sentinel (rate limiting) -> RequestId -> Logger
  -> DemoEvn -> WithContextDb -> LoggerToFile -> CustomError -> CORS -> Secure
  -> JWT MiddlewareFunc -> AuthCheckRole (Casbin) -> PermissionAction (data scope)
  -> Handler (apis/) -> Service (service/) -> GORM (models/)
```

### Handler Pattern

Handlers embed `api.Api` from go-admin-core and use a fluent builder:
```go
err := e.MakeContext(c).MakeOrm().Bind(&req).MakeService(&s.Service).Errors
```

### DTO Search Tags

DTOs use struct tags for automatic query building:
```go
Username string `form:"username" search:"type:contains;column:username;table:sys_user"`
```
`common/dto.MakeCondition()` reflects on these tags to build GORM scopes. Supported search types: `exact`, `contains`, `order`, `left` (join).

### Data Permission (Data Scope)

Casbin handles API-level RBAC. Data-level scoping is a separate system controlled by `DataPermission` in `common/actions/permission.go`. Five scopes: all, custom dept, own dept, dept+children, self only. Toggle via `settings.yml` -> `application.enabledp`.

### Database

GORM with support for MySQL, PostgreSQL, SQLite, SQL Server. Driver configured in `config/settings.yml`. Migrations live in `cmd/migrate/migration/version/` (core) and `version-local/` (custom). Migration files are timestamped and tracked in `sys_migration` table.

### Configuration

Primary config: `config/settings.yml`. Extended config struct: `config/extend.go` (`ExtConfig`). The config is loaded by go-admin-core's `config.Setup()` which also triggers database and storage initialization callbacks.

### Code Generation

The `other` module provides code generation from database tables. It reads table metadata, applies Go templates from `template/v4/`, and generates model, API handler, service, DTO, router, and Vue frontend files. Templates use the `no_actions` variant (current default).

### Key Dependencies

- **go-admin-core** (`github.com/go-admin-team/go-admin-core`) - Framework core providing runtime, config, JWT, logger, SDK base types
- **Gin** - HTTP framework
- **GORM** - ORM
- **Casbin** - RBAC policy engine
- **Cobra** - CLI framework
- **robfig/cron** - Scheduled jobs
- **swaggo/swag** - Swagger doc generation

### Common Package (common/)

- `actions/` - Permission middleware and CRUD action helpers
- `dto/` - Base pagination, search condition builder, generic delete/get DTOs
- `middleware/` - All middleware: auth, Casbin, CORS, logging, Sentinel, demo-env guard
- `middleware/handler/` - JWT callback implementations (login, payload, identity)
- `models/` - Base model types (ControlBy, ModelTime, ActiveRecord interface)
- `database/` - DB connection initialization
- `storage/` - Cache/queue storage initialization
- `file_store/` - File upload adapters (OSS, OBS, Kodo)
