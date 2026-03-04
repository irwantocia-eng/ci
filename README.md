# Go REST API Server

A production-ready REST API built with Go's standard library and SQLite.

## Features

- RESTful endpoints for users and products
- SQLite database with migrations
- HTTP middleware (logging, CORS)
- Comprehensive linting with golangci-lint (12 linters)
- Pre-commit hook for automatic code quality enforcement
- SonarCloud integration for continuous code analysis
- Unit tests with coverage reporting

## Quick Start

### Prerequisites

- **Go 1.26 or higher** - Download from [go.dev](https://go.dev/)
- Make (optional, for Makefile targets)

### Installation

```bash
# Download dependencies
go mod tidy

# Install linter (one-time)
make lint-install
```

### Running the Server

```bash
# Development mode
make run

# Build binary
make build

# Start server
./server
```

Server starts on `http://localhost:8089` (configurable via `PORT` environment variable).

### Database

SQLite database file `koban.db` is created automatically in the project root.

**Environment Variables:**
- `PORT` - Server port (default: 8089)
- `DB_PATH` - Database file path (default: koban.db)

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make run` | Start server in development mode |
| `make build` | Build binary |
| `make lint` | Run linter |
| `make lint-install` | Install golangci-lint |
| `make test` | Run lint + tests |
| `make coverage` | Generate test coverage report |
| `make coverage-html` | Open coverage in browser |
| `make sonar-scan` | Run SonarCloud analysis |
| `make clean` | Remove binary, database, and reports |

## Project Structure

```
.
├── main.go                 # Application entry point
├── main_test.go            # Main package tests
├── go.mod                  # Go module definition
├── Makefile                # Build commands
├── .golangci.yml           # Linter configuration
├── sonar-project.properties # SonarCloud configuration
├── .github/
│   └── hooks/
│       └── pre-commit     # Pre-commit linting hook
├── docs/
│   ├── LINTING.md         # Linting guide
│   └── SONAR.md           # SonarCloud integration guide
├── db/
│   ├── sqlite.go          # Database connection
│   ├── user_queries.go    # User operations
│   ├── user_queries_test.go
│   ├── product_queries.go  # Product operations
│   ├── product_queries_test.go
│   ├── config_queries.go  # Config operations
│   └── config_queries_test.go
├── handlers/
│   ├── users.go           # User HTTP handlers
│   ├── users_test.go
│   ├── products.go        # Product HTTP handlers
│   ├── products_test.go
│   ├── config.go          # Config HTTP handlers
│   └── config_test.go
└── models/
    ├── user.go            # User model
    ├── product.go         # Product model
    └── config.go          # Config model
```

## Code Quality

This project uses [golangci-lint](https://golangci-lint.run/) for automated code review with **12 enabled linters** covering security, reliability, and style.

### Pre-commit Hook

Linting is enforced at **commit time** via pre-commit hook:

```bash
# Install the hook (one-time)
cp .github/hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

Every `git commit` automatically runs `make lint`. If lint fails, the commit is rejected.

### SonarCloud Analysis

Continuous code quality monitoring via [SonarCloud](https://sonarcloud.io):

```bash
# Run locally (requires SONAR_TOKEN environment variable)
make sonar-scan

# Or view results online
# https://sonarcloud.io/dashboard?id=YOUR_PROJECT_KEY
```

**Quality Gates:**
- Coverage > 80%
- 0 Bugs
- 0 Vulnerabilities
- Code Smells < 50
- Duplication < 3%

See [docs/SONAR.md](docs/SONAR.md) for complete setup guide including:
- SonarCloud account setup (5 minutes)
- GitHub secrets configuration
- Understanding results and metrics
- Quality gates customization

See [docs/LINTING.md](docs/LINTING.md) for the complete linting guide including:
- All 12 linters explained
- Common issues and fixes
- IDE integration
- Troubleshooting

## License

MIT
