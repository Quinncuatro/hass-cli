Go Architecture Standards
Project Structure
Directory Layout
hass-cli/
├── cmd/
│   └── hass/               # Main application entry point
│       └── main.go
├── internal/               # Private application code
│   ├── app/               # Application logic layer
│   ├── cli/               # CLI command handlers
│   ├── tui/               # Terminal UI components
│   ├── client/            # Home Assistant API client
│   ├── config/            # Configuration management
│   ├── entity/            # Entity resolution and matching
│   └── cache/             # Caching implementation
├── pkg/                   # Public library code (if any)
├── test/                  # Integration tests and test data
├── docs/                  # Documentation
├── .project-rules/        # Development rules (this directory)
├── specs/                 # Functional specifications
├── go.mod
├── go.sum
├── Makefile
└── README.md
Package Design Principles
1. Dependency Direction
cmd/ depends on internal/
internal/app/ coordinates between domain packages
Domain packages (cli/, client/, config/) should not depend on each other directly
Use interfaces for cross-package communication
No circular dependencies allowed
2. Domain-Driven Structure
Each domain package should:

Have a clear, single responsibility
Export minimal public API through interfaces
Contain its own types, logic, and tests
Be independently testable
3. Interface Design
Define interfaces in consuming packages, not providing packages
Keep interfaces small and focused (1-3 methods ideal)
Use interface{} sparingly and only when truly needed
Prefer composition over inheritance
Package Responsibilities
cmd/hass
Application bootstrapping
Dependency injection setup
Signal handling and graceful shutdown
Environment variable processing
internal/app
Application orchestration
Cross-domain coordination
Main business logic workflows
Error handling aggregation
internal/cli
Command parsing and validation
CLI-specific error formatting
Output formatting for different verbosity levels
Command routing to appropriate handlers
internal/tui
Terminal UI components and layouts
Input handling and navigation
Screen management and rendering
TUI-specific state management
internal/client
Home Assistant API communication
HTTP client configuration and middleware
API request/response handling
Connection management and retry logic
internal/config
Configuration file reading/writing
Environment variable processing
Credential management (keyring integration)
Configuration validation and migration
internal/entity
Entity name resolution and fuzzy matching
Area and domain mapping
Entity filtering and searching
Smart suggestions for typos/ambiguity
internal/cache
Entity state caching
Cache invalidation strategies
Persistent cache storage (optional)
Cache performance monitoring
Design Patterns
1. Repository Pattern
Use for data access (API, cache, config):

go
type EntityRepository interface {
    GetEntities(ctx context.Context) ([]Entity, error)
    GetEntity(ctx context.Context, id string) (Entity, error)
    UpdateEntity(ctx context.Context, id string, state EntityState) error
}
2. Service Pattern
Use for business logic:

go
type EntityService interface {
    ResolveEntity(name string, area string, domain string) ([]Entity, error)
    ControlEntity(ctx context.Context, id string, action Action) error
}
3. Command Pattern
Use for CLI commands:

go
type Command interface {
    Execute(ctx context.Context, args []string) error
    Usage() string
    Examples() []string
}
4. Observer Pattern
Use for TUI updates and cache invalidation:

go
type EventSubscriber interface {
    OnEntityStateChanged(entity Entity)
    OnConnectionStatusChanged(connected bool)
}
Error Handling Architecture
Error Types
Define domain-specific error types:

go
type EntityNotFoundError struct {
    Name string
    Suggestions []string
}

type AmbiguousEntityError struct {
    Name string
    Matches []Entity
}

type ConnectionError struct {
    URL string
    Cause error
}
Error Wrapping Strategy
Wrap errors at package boundaries
Include context at each layer
Use fmt.Errorf with %w verb for wrapping
Create typed errors for business logic failures
Error Handling Levels
Package Level: Handle technical errors, return domain errors
Application Level: Coordinate error handling, add context
CLI/TUI Level: Format errors for user display
Concurrency Patterns
1. Context Usage
Pass context.Context as first parameter to all operations
Use context for cancellation, timeouts, and values
Don't store context in structs
2. Goroutine Management
Use sync.WaitGroup for coordinating goroutines
Always handle goroutine cleanup
Use channels for communication, mutexes for shared state
3. Worker Pool Pattern
For API requests and cache operations:

go
type WorkerPool struct {
    workers int
    jobs    chan Job
    results chan Result
}
Configuration Management
Configuration Sources (Priority Order)
Command-line flags
Environment variables
Configuration file
Default values
Configuration Structure
go
type Config struct {
    HomeAssistant HAConfig     `yaml:"homeassistant"`
    Aliases       AliasConfig  `yaml:"aliases"`
    Preferences   PrefConfig   `yaml:"preferences"`
}
Testing Architecture
Test Organization
Unit tests alongside source files (*_test.go)
Integration tests in test/ directory
Benchmarks for performance-critical code
Example tests for public APIs
Test Patterns
Use table-driven tests for multiple scenarios
Create test helpers for common setup
Use testify for assertions and mocking
Test error paths explicitly
Mock Strategy
Generate mocks for interfaces using mockery
Mock external dependencies (HTTP, filesystem)
Use dependency injection to enable testing
Performance Considerations
Memory Management
Use object pools for frequently allocated objects
Avoid memory leaks with proper resource cleanup
Profile memory usage in performance tests
CPU Optimization
Use efficient algorithms for entity matching
Implement caching at appropriate levels
Profile CPU usage for hot paths
I/O Optimization
Use connection pooling for HTTP requests
Implement request batching where possible
Use buffered I/O for file operations
Security Architecture
Credential Management
Use OS keyring for secure token storage
Never log credentials or tokens
Implement secure config file permissions (600)
Input Validation
Validate all user inputs
Sanitize data before API calls
Use parameterized queries/requests
Network Security
Validate SSL certificates by default
Support custom CA certificates
Implement request timeouts
