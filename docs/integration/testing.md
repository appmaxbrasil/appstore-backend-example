# Testing Reference

Test organization, running tests, mock patterns, and how to add new tests.

---

## Running Tests

```bash
# In Docker (recommended — matches CI environment)
make test

# Locally (requires Go 1.25+ and database access for integration tests)
go test ./...

# Specific package
go test ./tests/unit/services/...

# Verbose output
go test -v ./...

# With coverage
go test -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

Inside Docker, `make test` runs:
```bash
docker compose exec -T app sh -lc 'go test ./...'
```

---

## Test Organization

```
tests/
├── unit/                           24 test files
│   ├── bootstrap/                  Config loading, module construction
│   │   ├── appmax_config_test.go
│   │   └── modules_test.go
│   ├── controllers/                Controller constructors, HTTP behavior
│   │   ├── constructors_test.go
│   │   ├── constructors_success_test.go
│   │   ├── install_controller_test.go
│   │   ├── upstream_error_message_test.go
│   │   ├── fake_http_context_test.go     (test helper)
│   │   └── setup_test.go                (test helper)
│   ├── gateway/appmax/             HTTP client, retry, error parsing
│   │   ├── client_test.go
│   │   ├── contracts_test.go
│   │   ├── helpers_test.go
│   │   └── setup_test.go                (test helper)
│   ├── models/                     Model behavior
│   │   └── webhook_event_test.go
│   ├── mocks/                      Hand-written mocks
│   │   ├── installation_repository_mock.go
│   │   ├── order_repository_mock.go
│   │   └── webhook_event_repository_mock.go
│   ├── repositories/               Repository constructor validation
│   │   └── constructors_test.go
│   ├── requests/                   Request parsing, payload detection
│   │   ├── webhook_data_test.go
│   │   └── webhook_envelope_payloads_test.go
│   └── services/                   Business logic
│       ├── appmax_service_constructor_test.go
│       ├── appmax_service_test.go
│       ├── checkout_service_test.go
│       ├── install_service_test.go
│       ├── token_manager_test.go
│       ├── webhook_service_test.go
│       ├── webhook_service_payloads_test.go
│       └── setup_test.go                (test helper)
└── integration/                    4 test files
    ├── db_test.go                  DB setup + truncateTables helper
    ├── install_test.go             Installation CRUD, unique constraints
    ├── checkout_test.go            Order persistence
    └── webhook_test.go             Webhook event storage, deduplication
```

28 test files total.

---

## Test Patterns

### Unit Tests

Unit tests use **hand-written function-based mocks**. Each mock struct has function
fields that tests configure per scenario:

```go
repo := &mocks.MockInstallationRepository{
    FindByExternalKeyFunc: func(_ context.Context, _ string) (*models.Installation, error) {
        return nil, nil // simulate "not found"
    },
    CreateFunc: func(_ context.Context, inst *models.Installation) error {
        inst.ID = 10 // simulate DB-assigned ID
        return nil
    },
}

svc, err := services.NewInstallService(repo)
require.NoError(t, err)

inst, wasCreated, err := svc.Upsert(context.Background(), input)
require.NoError(t, err)
assert.True(t, wasCreated)
assert.Equal(t, int64(10), inst.ID)
```

Advantages of this approach:
- No code generation tools required
- Each test explicitly declares the behavior it expects
- Easy to verify which repository methods were called

### Integration Tests

Integration tests run against a **real PostgreSQL database**. The `db_test.go` file
handles connection setup and provides a `truncateTables(t)` helper that clears all
tables between tests:

```go
func TestInstallation_CreateAndRetrieve(t *testing.T) {
    truncateTables(t)

    _, err := testDB.Exec(`INSERT INTO installations ...`)
    require.NoError(t, err)

    var externalKey string
    err = testDB.QueryRow(`SELECT external_key FROM installations WHERE ...`).Scan(&externalKey)
    require.NoError(t, err)
    assert.Equal(t, "key-abc", externalKey)
}
```

Integration tests verify:
- Database constraints (unique keys, foreign keys)
- Migration correctness (columns exist, types match)
- Default values and auto-generated fields (UUID, timestamps)

### Assertions

All tests use `testify`:
- `require.NoError(t, err)` — stops the test immediately on error
- `assert.Equal(t, expected, actual)` — reports failure but continues
- `assert.ErrorContains(t, err, "substring")` — verifies error messages
- `assert.ErrorIs(t, err, services.ErrNilDependency)` — verifies error types

### Constructor Tests

Every service, controller, and repository has constructor tests that verify:

```go
// Nil dependency rejection
func TestInstallServiceConstructor_RejectsNilDependency(t *testing.T) {
    svc, err := services.NewInstallService(nil)
    require.Error(t, err)
    assert.Nil(t, svc)
    assert.ErrorIs(t, err, services.ErrNilDependency)
}

// Valid construction
func TestInstallServiceConstructor_Success(t *testing.T) {
    svc, err := services.NewInstallService(repo)
    require.NoError(t, err)
    assert.NotNil(t, svc)
}
```

This pattern ensures that nil dependencies fail fast at startup rather than causing
nil pointer panics at runtime.

---

## Mock Strategy

Mocks live in `tests/unit/mocks/` — one file per repository interface:

| Mock | Interface | Methods |
|------|-----------|---------|
| `MockInstallationRepository` | `InstallationRepository` | `FindByExternalKey`, `Create`, `Save` |
| `MockOrderRepository` | `OrderRepository` | `FindByAppmaxOrderID`, `FindByAppmaxOrderIDAndInstallation`, `Create`, `Save` |
| `MockWebhookEventRepository` | `WebhookEventRepository` | `Create`, `Save`, `FindProcessedDuplicate` |

Each mock uses the function-field pattern:

```go
type MockInstallationRepository struct {
    FindByExternalKeyFunc func(ctx context.Context, key string) (*models.Installation, error)
    CreateFunc            func(ctx context.Context, inst *models.Installation) error
    SaveFunc              func(ctx context.Context, inst *models.Installation) error
}

func (m *MockInstallationRepository) FindByExternalKey(ctx context.Context, key string) (*models.Installation, error) {
    return m.FindByExternalKeyFunc(ctx, key)
}
```

No mock generation tools (mockgen, mockery) are used.

---

## What Each Test Suite Covers

| Suite | Key Scenarios |
|-------|---------------|
| `services/install_service_test.go` | Upsert create new, update existing, find error, create error, save error, nil dependency |
| `services/checkout_service_test.go` | Payment flows, customer+order auto-creation, best-effort persistence |
| `services/webhook_service_test.go` | Event handling, status mapping, deduplication, no-op events |
| `services/token_manager_test.go` | Token caching, expiry buffer, Redis fallback |
| `services/appmax_service_test.go` | Appmax API orchestration, token injection |
| `controllers/install_controller_test.go` | Start flow, callback with/without token, health check POST validation |
| `controllers/constructors_test.go` | All controller constructors reject nil deps and invalid config |
| `gateway/appmax/client_test.go` | HTTP client retry on 502/503/504, timeout, error parsing |
| `requests/webhook_data_test.go` | 5 payload model detection, order_id extraction from each model |
| `integration/install_test.go` | Create+retrieve, unique external_key constraint, credential update |
| `integration/checkout_test.go` | Order persistence, unique appmax_order_id constraint |
| `integration/webhook_test.go` | Webhook event storage, processed flag, deduplication |

---

## Adding a New Test

### Unit Test

1. Create `tests/unit/<layer>/<name>_test.go`
2. Use the `_test` package suffix (e.g., `package services_test`)
3. Import mocks from `tests/unit/mocks/` if needed
4. Follow the function-field mock pattern for new dependencies
5. Test both success and error paths

### New Mock

If your test needs a mock for a new interface:

1. Create `tests/unit/mocks/<name>_mock.go`
2. Define a struct with function fields matching the interface methods
3. Implement the interface by delegating to the function fields

### Integration Test

1. Add to `tests/integration/<name>_test.go`
2. Call `truncateTables(t)` at the start of each test
3. Use raw SQL via `testDB` (the shared `*sql.DB` connection)
4. Test constraints, defaults, and auto-generated values

### Run

```bash
make test
```
