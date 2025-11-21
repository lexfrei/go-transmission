# Testing

This project uses two types of tests: unit tests and end-to-end (E2E) tests.

## Unit Tests

Unit tests use HTTP mocking to test the client without a real Transmission daemon.

### Running Unit Tests

```bash
# Run all unit tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Test Structure

- `api/transmission/client_test.go` - Client method tests
- `api/transmission/mock_test.go` - HTTP RoundTripper mock using [mok](https://github.com/ymz-ncnk/mok)
- `api/transmission/testdata/responses/` - Real JSON responses captured from Transmission

### Adding Unit Tests

1. Capture a real response from Transmission and save it to `testdata/responses/`
2. Embed it in `client_test.go` using `//go:embed`
3. Create test using `NewRoundTripperMock()` to mock HTTP responses

Example:

```go
//go:embed testdata/responses/your-response.json
var yourResponse string

func TestYourMethod(t *testing.T) {
    t.Parallel()

    mock := NewRoundTripperMock()
    mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
        return sessionIDResponse(), nil
    })
    mock.RegisterRoundTrip(func(_ *http.Request) (*http.Response, error) {
        return jsonResponse(yourResponse), nil
    })

    client, err := newTestClient(mock)
    require.NoError(t, err)

    // Test your method
}
```

## E2E Tests

E2E tests run against a real Transmission daemon in a Docker container.

### Requirements

- Docker or Podman
- Network access to pull container images

### Running E2E Tests

```bash
# Run E2E tests
go test -tags=e2e ./api/transmission/...

# Verbose output
go test -v -tags=e2e ./api/transmission/...

# With timeout (recommended)
go test -v -tags=e2e -timeout=10m ./api/transmission/...
```

### E2E Test Structure

- `api/transmission/e2e_test.go` - All E2E tests
- `api/transmission/testdata/torrents/` - Test torrent files (Ubuntu ISO, Rocky Linux ISO)

### How E2E Tests Work

1. Tests start a Transmission container using [testcontainers-go](https://github.com/testcontainers/testcontainers-go)
2. Container uses `lscr.io/linuxserver/transmission:latest` image
3. Tests wait for Transmission to be ready on port 9091
4. Each test creates a fresh client and runs operations
5. Container is cleaned up after all tests complete

### E2E Tests in CI

E2E tests only run when a PR has the `ready-for-e2e` label. This prevents running slow tests on every commit.

## Code Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

Coverage is automatically uploaded to Codecov on every PR.
