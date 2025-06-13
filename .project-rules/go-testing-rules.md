Go Testing Requirements
Test Coverage Requirements
Coverage Targets
Minimum Overall Coverage: 80%
Critical Packages: 90% coverage required
internal/client (API client)
internal/entity (entity resolution)
internal/config (configuration management)
UI Packages: 60% coverage minimum
internal/cli (CLI commands)
internal/tui (terminal UI)
Coverage Validation
bash
# Generate and check coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Fail if coverage below threshold
go test -coverprofile=coverage.out ./... && \
  go tool cover -func=coverage.out | grep "total:" | \
  awk '{if ($3+0 < 80) exit 1}'
Test Organization
Test File Structure
internal/
├── client/
│   ├── client.go
│   ├── client_test.go          # Unit tests
│   ├── client_integration_test.go  # Integration tests
│   └── testdata/              # Test fixtures
│       ├── valid_response.json
│       └── error_response.json
Test Categories
Unit Tests: Test individual functions and methods
Integration Tests: Test package interactions
End-to-End Tests: Test complete workflows
Benchmark Tests: Performance validation
Example Tests: Executable documentation
Unit Testing Standards
Test Function Naming
go
func TestFunctionName_Condition_ExpectedResult(t *testing.T) {
    // Test implementation
}

// Examples:
func TestResolveEntity_ExactMatch_ReturnsEntity(t *testing.T) {}
func TestResolveEntity_NoMatch_ReturnsEmptySlice(t *testing.T) {}
func TestResolveEntity_InvalidInput_ReturnsError(t *testing.T) {}
Table-Driven Tests
go
func TestEntityMatching(t *testing.T) {
    tests := []struct {
        name        string
        input       EntityQuery
        entities    []Entity
        expected    []Entity
        expectedErr error
    }{
        {
            name: "exact entity ID match",
            input: EntityQuery{
                Name: "light.living_room_main",
            },
            entities: []Entity{
                {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
                {ID: "light.kitchen_main", FriendlyName: "Kitchen Main"},
            },
            expected: []Entity{
                {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
            },
            expectedErr: nil,
        },
        {
            name: "fuzzy area match",
            input: EntityQuery{
                Name: "living",
                Domain: "light",
            },
            entities: []Entity{
                {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
                {ID: "light.living_room_accent", FriendlyName: "Living Room Accent"},
                {ID: "light.bedroom_main", FriendlyName: "Bedroom Main"},
            },
            expected: []Entity{
                {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
                {ID: "light.living_room_accent", FriendlyName: "Living Room Accent"},
            },
            expectedErr: nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resolver := NewEntityResolver(tt.entities)
            result, err := resolver.ResolveEntity(context.Background(), tt.input)

            // Error checking
            if tt.expectedErr != nil {
                assert.Error(t, err)
                assert.IsType(t, tt.expectedErr, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
Test Helpers
go
// Create test helpers for common setup
func createTestEntities() []Entity {
    return []Entity{
        {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
        {ID: "light.kitchen_main", FriendlyName: "Kitchen Main"},
        {ID: "switch.coffee_maker", FriendlyName: "Coffee Maker"},
    }
}

func createTestConfig() Config {
    return Config{
        HAURL:     "http://localhost:8123",
        HAToken:   "test-token",
        HATimeout: 5 * time.Second,
    }
}
Integration Testing
Test Environment Setup
go
func TestMain(m *testing.M) {
    // Setup test environment
    if err := setupTestEnvironment(); err != nil {
        fmt.Printf("Failed to setup test environment: %v\n", err)
        os.Exit(1)
    }

    // Run tests
    code := m.Run()

    // Cleanup
    cleanupTestEnvironment()

    os.Exit(code)
}
HTTP Client Testing
go
func TestClient_GetEntities_Success(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request
        assert.Equal(t, "GET", r.Method)
        assert.Equal(t, "/api/states", r.URL.Path)
        assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

        // Return test response
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `[{"entity_id": "light.test", "state": "on"}]`)
    }))
    defer server.Close()

    // Test client
    client := NewClient(Config{
        HAURL:   server.URL,
        HAToken: "test-token",
    })

    entities, err := client.GetEntities(context.Background())
    assert.NoError(t, err)
    assert.Len(t, entities, 1)
    assert.Equal(t, "light.test", entities[0].ID)
}
Mocking Strategy
Interface Mocking
go
//go:generate mockery --name=EntityRepository --output=./mocks
type EntityRepository interface {
    GetEntities(ctx context.Context) ([]Entity, error)
    GetEntity(ctx context.Context, id string) (Entity, error)
}

// In tests
func TestEntityService_ResolveEntity_Success(t *testing.T) {
    mockRepo := &mocks.EntityRepository{}
    mockRepo.On("GetEntities", mock.Anything).Return(createTestEntities(), nil)

    service := NewEntityService(mockRepo)
    result, err := service.ResolveEntity(context.Background(), "living lights")

    assert.NoError(t, err)
    assert.NotEmpty(t, result)
    mockRepo.AssertExpectations(t)
}
HTTP Mocking
go
func TestClient_HandleAPIError(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusUnauthorized)
        fmt.Fprint(w, `{"message": "Invalid authentication"}`)
    }))
    defer server.Close()

    client := NewClient(Config{HAURL: server.URL, HAToken: "invalid"})
    _, err := client.GetEntities(context.Background())

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "authentication")
}
Benchmark Testing
Performance Benchmarks
go
func BenchmarkEntityResolver_ResolveEntity(b *testing.B) {
    entities := createLargeEntitySet(1000) // 1000 test entities
    resolver := NewEntityResolver(entities)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := resolver.ResolveEntity(context.Background(), EntityQuery{
            Name: "living lights",
        })
        if err != nil {
            b.Fatalf("Unexpected error: %v", err)
        }
    }
}

func BenchmarkFuzzyMatching(b *testing.B) {
    testCases := []struct {
        name   string
        input  string
        target string
    }{
        {"short", "lr", "living room"},
        {"medium", "living", "living room lights"},
        {"long", "living room main lights", "light.living_room_main_lights"},
    }

    for _, tc := range testCases {
        b.Run(tc.name, func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                calculateSimilarity(tc.input, tc.target)
            }
        })
    }
}
Error Testing
Error Path Testing
go
func TestClient_NetworkError_ReturnsWrappedError(t *testing.T) {
    // Test with invalid URL to trigger network error
    client := NewClient(Config{
        HAURL:   "http://invalid-host:1234",
        HAToken: "test-token",
    })

    _, err := client.GetEntities(context.Background())

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "connection")

    // Test error wrapping
    var netErr *net.OpError
    assert.True(t, errors.As(err, &netErr))
}

func TestEntityResolver_AmbiguousEntity_ReturnsProperError(t *testing.T) {
    entities := []Entity{
        {ID: "light.living_room_main", FriendlyName: "Living Room Main"},
        {ID: "light.living_room_accent", FriendlyName: "Living Room Accent"},
    }

    resolver := NewEntityResolver(entities)
    _, err := resolver.ResolveEntity(context.Background(), EntityQuery{
        Name: "living",
        Domain: "light",
        RequireUnique: true, // Force ambiguity error
    })

    var ambiguousErr *AmbiguousEntityError
    assert.True(t, errors.As(err, &ambiguousErr))
    assert.Len(t, ambiguousErr.Matches, 2)
}
Test Data Management
Test Fixtures
go
// Use embedded test data
//go:embed testdata/entities.json
var testEntitiesJSON []byte

func loadTestEntities(t *testing.T) []Entity {
    var entities []Entity
    err := json.Unmarshal(testEntitiesJSON, &entities)
    require.NoError(t, err)
    return entities
}
Test Database/Cache
go
func createTestCache(t *testing.T) *Cache {
    // Use in-memory cache for tests
    cache := NewCache(CacheConfig{
        TTL:     1 * time.Minute,
        MaxSize: 100,
        Storage: NewMemoryStorage(), // Instead of disk storage
    })

    t.Cleanup(func() {
        cache.Clear()
    })

    return cache
}
Test Utilities
Assertion Helpers
go
func assertEntityEqual(t *testing.T, expected, actual Entity) {
    t.Helper()
    assert.Equal(t, expected.ID, actual.ID)
    assert.Equal(t, expected.State, actual.State)
    assert.Equal(t, expected.FriendlyName, actual.FriendlyName)
}

func assertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
    t.Helper()
    if err != nil {
        t.Fatalf("Expected no error, got: %v %v", err, msgAndArgs)
    }
}
Context Helpers
go
func contextWithTimeout(t *testing.T, timeout time.Duration) context.Context {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    t.Cleanup(cancel)
    return ctx
}
Continuous Integration Testing
Test Pipeline Requirements
yaml
# .github/workflows/test.yml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'

      - name: Run tests
        run: |
          go test -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out

      - name: Check coverage
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$COVERAGE < 80" | bc -l) )); then
            echo "Coverage $COVERAGE% is below 80% threshold"
            exit 1
          fi
Race Condition Testing
Always run tests with -race flag
Test concurrent operations explicitly
Use proper synchronization in tests
Integration Test Environment
Use Docker containers for external dependencies
Implement proper cleanup between tests
Use separate test databases/caches
Mock external services appropriately
