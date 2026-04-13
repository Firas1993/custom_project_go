# AGENTS.md - AI Agent Guidelines

This document provides guidance for AI coding agents working with this repository.

## Project Overview

**irp-app-from-template**: IRP scanner API

- **Type**: application
- **Organization**: mime-rona
- **Module**: `github.com/mime-rona/irp-app-from-template`

## Project Structure

```
irp-app-from-template/
├── main.go                    # Application entry point
├── server.go                  # HTTP server setup (Gin framework)
├── go.mod                     # Go module definition
├── Makefile                   # Development commands
├── Dockerfile                 # Multi-stage Docker build (applications only)
├── .golangci.yml              # Linter configuration
├── .mockery.yaml              # Mock generation config
├── .pre-commit-config.yaml    # Pre-commit hooks
├── .github/
│   └── workflows/
│       ├── ci-cd.yml          # Main CI/CD pipeline
│       ├── promote.yml        # Manual promotion workflow
│       └── cruft-update.yml   # Template auto-update
├── .makefiles/
│   ├── Makefile.docker        # Docker targets
│   ├── Makefile.migration     # Atlas migration targets (optional)
│   └── Makefile.gql           # GraphQL generation targets (optional)
├── cmd/                       # CLI binaries (applications only)
├── internal/                  # Private application code
│   └── testing/mocks/         # Generated mocks
├── pkg/                       # Public reusable packages
└── migrations/                # Atlas database migrations (optional)
```

## Key Technologies

| Technology    | Purpose               |
| ------------- | --------------------- |
| Go            | Language              |
| Gin-Gonic     | HTTP web framework    |
| Testify       | Testing assertions    |
| Mockery       | Mock generation       |
| golangci-lint | Linting (24+ linters) |
| Atlas         | Database migrations   |
| Pre-commit    | Git hooks             |

## Development Commands

### Code Quality

| Command                    | Description                         |
| -------------------------- | ----------------------------------- |
| `make format`              | Run code formatter via pre-commit   |
| `make lint`                | Run all linters via pre-commit      |
| `make test`                | Run tests with coverage enforcement |
| `make test package=./path` | Run tests for specific package      |
| `make mod`                 | Tidy go.mod and go.sum              |

### Mock Generation

| Command         | Description                  |
| --------------- | ---------------------------- |
| `make mock.gen` | Generate mocks using Mockery |

Mocks location:

- Applications: `internal/testing/mocks/`
- Libraries: `testing/mocks/`

### Docker (Applications Only)

| Command             | Description                         |
| ------------------- | ----------------------------------- |
| `make docker.build` | Build Docker image                  |
| `make docker.run`   | Build and run container (port 8080) |
| `make docker.clean` | Clean containers for this project   |

### Database Migrations (When atlas.hcl Exists)

| Command                             | Description                       |
| ----------------------------------- | --------------------------------- |
| `make migration desc="description"` | Create migration from schema diff |
| `make migration.new`                | Create blank migration template   |
| `make migration.recreate`           | Delete latest and regenerate      |
| `make migration.delete_latest`      | Remove latest migration file      |

### Dependency Management

| Command        | Description                          |
| -------------- | ------------------------------------ |
| `make update`  | Update all dependencies to latest    |
| `make upgrade` | Update cruft template + dependencies |
| `make clean`   | Run go clean                         |

## Code Conventions

### Logging Rules

We use zerolog for all our logging. The logger should always be taken from the context, like such: `log := zerolog.Ctx(ctx)`, except in `main()` functions.
Note: **Print statements are forbidden.** The `forbidigo` linter prevents:

- `fmt.Print()`, `fmt.Printf()`, `fmt.Println()`
- `print()`, `println()`

### Error Handling

- Always check and handle errors explicitly
- Either log the error or return it, but not both
- Enforced linters:
  - `errcheck` - Detects unchecked errors
  - `noctx` - Requires context.Context usage
  - `sqlclosecheck` - Ensures database connections close
  - `rowserrcheck` - Checks Row.Err() calls

### Code Complexity Limits

| Metric                | Limit              |
| --------------------- | ------------------ |
| Function statements   | 50 max             |
| Cyclomatic complexity | 15 max             |
| Line length           | Handled by golines |

### Commit Messages

Conventional commits are required (enforced by commitlint):

```
type(scope): description

Examples:
feat: Add user authentication
fix(api): Handle expired tokens correctly
docs: Update README with new commands
refactor(handlers): Extract validation logic
```

## Testing Patterns

### Test Framework

Tests use **Testify** for assertions:

```go
func TestExample(t *testing.T) {
    router := CreateRouter()
    w := httptest.NewRecorder()
    req, err := http.NewRequestWithContext(t.Context(), "GET", "/", nil)
    require.NoError(t, err)
	router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, `{"message":"Hello World"}`, w.Body.String())
}
```

### Mock Generation

Mocks are generated with Mockery using testify template:

```yaml
# .mockery.yaml
template: testify
mockname: "Mock{{.InterfaceName}}"
```

Run `make mock.gen` to regenerate all mocks.

### Coverage Requirements

- Test coverage is mandatory
- Tests fail if any file has less than 100% coverage
- Coverage report: `make cov.render`
- Generated files (`*.gen.go`) are excluded from coverage

## CI/CD Pipeline

### Workflow: ci-cd.yml

**Triggers:** Push to main, Pull requests

**Pipeline Steps:**

1. Setup Go, Python, SSH agent
2. Install Atlas migration tool
3. Cache golangci-lint
4. Run pre-commit linters
5. GCP authentication
6. Build Docker image (applications)
7. Run tests with coverage
8. Validate migrations (if enabled, PRs only)
9. Deploy to dev (main branch, applications only)

### Workflow: promote.yml

Manual workflow for promoting to environments (staging, production).

### Workflow: cruft-update.yml

Runs weekly (Mondays) to:

- Update template via cruft
- Update Go dependencies
- Create PR with changes

## Configuration Files

### .golangci.yml

Comprehensive linting with 24+ linters enabled:

- Code quality: `govet`, `gosec`, `gocritic`, `staticcheck`
- Error handling: `errcheck`, `noctx`, `sqlclosecheck`
- Code patterns: `gocyclo`, `funlen`, `dupl`
- Formatters: `gci` (imports), `golines` (line length)

### .mockery.yaml

Mock generation configuration:

- Template: testify
- Recursive interface discovery
- Auto-naming: `Mock{InterfaceName}`

### .pre-commit-config.yaml

Git hooks for:

- File validation (size, encoding, YAML/TOML)
- Code formatting (Prettier, golangci-lint)
- Go tools (mod tidy, no testing.T in non-test files)
- Commit message validation (conventional commits)
- Migration validation (when atlas.hcl exists)

## Important Notes for AI Agents

1. **Sensitive Files - DO NOT READ**: Never read files matching `*.env*` pattern (e.g., `.env`, `.env.local`, `.env.production`). These files contain secrets and sensitive configuration.

2. **Private Dependencies**: This project uses private GitHub repos via SSH. GOPRIVATE is set to `github.com/mime-rona/*`.

3. **Format code**: Always run `make format` before committing to ensure code follow guidelines

4. **Pre-commit Hooks**: Always run `make lint` before committing to ensure code passes all checks.

5. **Test Before PR**: Run `make test` to ensure all tests pass with required coverage.

6. **Migration Safety**: When modifying database schemas, use `make migration` to generate proper migration files.

7. **Docker Builds**: Optionally, use `make docker.build` to verify that the application builds correctly in a container, but be aware that `make test` will also build the container as part of testing
