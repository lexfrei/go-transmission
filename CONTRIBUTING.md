# Contributing to go-transmission

Thank you for your interest in contributing to go-transmission!

## Requirements

- Go 1.24 or later
- Docker or Podman (for E2E tests)
- golangci-lint

## Getting Started

1. Fork the repository

2. Clone your fork:

    ```bash
    git clone https://github.com/YOUR_USERNAME/go-transmission.git
    cd go-transmission
    ```

3. Create a feature branch:

    ```bash
    git checkout -b feat/your-feature
    ```

## Development Workflow

### Running Tests

```bash
# Unit tests
go test ./...

# Unit tests with race detection
go test -race ./...

# E2E tests (requires Docker)
go test -tags=e2e ./api/transmission/...
```

### Linting

```bash
golangci-lint run
```

All code must pass linting before merge.

### Code Style

- Follow standard Go conventions
- Run `gofmt` and `goimports` on your code
- Keep functions focused and under 60 lines where possible
- Add tests for new functionality

## Commit Messages

Use semantic commit messages:

```text
type(scope): brief description

Optional longer explanation.
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `test`: Tests
- `refactor`: Code refactoring
- `chore`: Maintenance

Examples:

```text
feat(api): add TorrentRename method
fix(transport): handle CSRF token refresh
docs: update API examples
test: add unit tests for session methods
```

## Pull Request Process

1. Ensure all tests pass locally
2. Update documentation if needed
3. Create a PR against `master` branch
4. Wait for CI to pass
5. Request review

### E2E Tests in CI

E2E tests run only when a PR has the `ready-for-e2e` label. Add this label when your PR is ready for full testing.

## Questions?

Open an issue for questions or discussions.
