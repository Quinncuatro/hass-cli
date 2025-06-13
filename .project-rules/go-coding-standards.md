Go Coding Standards
Code Formatting
Automatic Formatting
Use gofmt for all Go source files
Use goimports to manage imports automatically
Run formatting on save in development environment
Enforce formatting in CI/CD pipeline
Line Length
Prefer lines under 100 characters
Break long function calls across multiple lines
Use meaningful variable names even if longer
Naming Conventions
Packages
Use short, lowercase names
Avoid underscores or mixed caps
Use singular nouns (e.g., config, client, not configs, clients)
Package name should match directory name
Variables and Functions
Use camelCase for private members
Use PascalCase for public members
Use descriptive names over abbreviations
Prefer userID over uid, configuration over cfg
Constants
go
// Good
const (
    DefaultTimeout = 10 * time.Second
    MaxRetries     = 3
)

// Avoid
const DEFAULT_TIMEOUT = 10
Interface Names
Use -er suffix for single-method interfaces: Reader, Writer, Resolver
For multi-method interfaces, use descriptive nouns: EntityService, ConfigManager
Variable Declaration
Declaration Style
go
// Prefer short declaration when possible
entities := make([]Entity, 0, len(rawEntities))

// Use var for zero values
var (
    config Config
    client *http.Client
)

// Use explicit types when clarity is important
var timeout time.Duration = 30 * time.Second
Grouping Declarations
go
// Group related variables
var (
    ErrEntityNotFound     = errors.New("entity not found")
    ErrAmbiguousEntity    = errors.New("ambiguous entity name")
    ErrConnectionFailed   = errors.New("connection failed")
)
Function Design
Function Length
Keep functions under 50 lines when possible
Extract complex logic into helper functions
Use early returns to reduce nesting
Parameter Order
Context first: func Process(ctx context.Context, data Data) error
Required parameters before optional
Group related parameters into structs when there are many
Return Values
go
// Good: Clear error handling
func ResolveEntity(name string) (Entity, error) {
    // implementation
}

// Good: Multiple return values with clear meaning
func ParseCommand(args []string) (command string, params []string, err error) {
    // implementation
}

// Avoid: Too many return values
func ComplexOperation() (string, int, bool, []string, error) // Too many
Error Handling
Error Creation
go
// Use errors.New for simple errors
var ErrInvalidToken = errors.New("invalid authentication token")

// Use fmt.Errorf for formatted errors with context
return fmt.Errorf("failed to resolve entity %q: %w", name, err)

// Create custom error types for complex cases
type ValidationError struct {
    Field string
    Value interface{}
    Reason string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Reason)
}
Error Handling Patterns
go
// Check errors immediately
result, err := someOperation()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Use early returns for error cases
func ProcessEntities(entities []Entity) error {
    if len(entities) == 0 {
        return errors.New("no entities provided")
    }

    for _, entity := range entities {
        if err := processEntity(entity); err != nil {
            return fmt.Errorf("failed to process entity %s: %w", entity.ID, err)
        }
    }
    return nil
}

// Don't ignore errors
_ = conn.Close() // Bad
if err := conn.Close(); err != nil {
    log.Printf("failed to close connection: %v", err)
}
Comments and Documentation
Package Documentation
go
// Package client provides a Home Assistant API client with caching
// and automatic retry capabilities.
//
// The client supports both WebSocket and REST API connections,
// with intelligent fallback between connection types.
package client
Function Documentation
go
// ResolveEntity attempts to find Home Assistant entities matching the given criteria.
// It performs fuzzy matching against entity names, friendly names, and area assignments.
//
// The search process follows this priority:
//   1. Exact entity ID match
//   2. Exact friendly name match
//   3. Fuzzy area + domain matching
//   4. Partial name matching with similarity scoring
//
// Returns a slice of matching entities, or an error if the search fails.
// An empty slice indicates no matches were found.
func ResolveEntity(ctx context.Context, name, area, domain string) ([]Entity, error) {
    // implementation
}
Inline Comments
go
// Only comment complex or non-obvious logic
func calculateSimilarity(a, b string) float64 {
    // Use Levenshtein distance normalized by string length
    distance := levenshteinDistance(a, b)
    maxLen := max(len(a), len(b))
    if maxLen == 0 {
        return 1.0 // Both strings are empty
    }
    return 1.0 - float64(distance)/float64(maxLen)
}
Struct Design
Struct Definition
go
// Use clear, descriptive field names
type Entity struct {
    ID           string                 `json:"entity_id"`
    State        string                 `json:"state"`
    FriendlyName string                 `json:"friendly_name"`
    Attributes   map[string]interface{} `json:"attributes"`
    LastChanged  time.Time              `json:"last_changed"`
    LastUpdated  time.Time              `json:"last_updated"`
}

// Group related fields
type Config struct {
    // Home Assistant connection settings
    HAURL     string        `yaml:"url"`
    HAToken   string        `yaml:"token"`
    HATimeout time.Duration `yaml:"timeout"`

    // Application preferences
    CacheEnabled bool          `yaml:"cache_enabled"`
    CacheTTL     time.Duration `yaml:"cache_ttl"`
    Verbosity    int           `yaml:"verbosity"`
}
Constructor Functions
go
// Provide constructor functions for complex structs
func NewClient(config Config) (*Client, error) {
    if config.HAURL == "" {
        return nil, errors.New("Home Assistant URL is required")
    }

    client := &Client{
        baseURL:    strings.TrimSuffix(config.HAURL, "/"),
        token:      config.HAToken,
        timeout:    config.HATimeout,
        httpClient: &http.Client{Timeout: config.HATimeout},
    }

    return client, nil
}
Interface Usage
Interface Definition
go
// Keep interfaces small and focused
type EntityResolver interface {
    ResolveEntity(ctx context.Context, name string) ([]Entity, error)
}

type EntityController interface {
    TurnOn(ctx context.Context, entityID string) error
    TurnOff(ctx context.Context, entityID string) error
    SetState(ctx context.Context, entityID string, state EntityState) error
}

// Combine interfaces when needed
type EntityService interface {
    EntityResolver
    EntityController
}
Testing Standards
Test Function Names
go
func TestResolveEntity_ExactMatch_ReturnsEntity(t *testing.T) {
    // Test exact entity ID matching
}

func TestResolveEntity_NoMatch_ReturnsEmptySlice(t *testing.T) {
    // Test no matches scenario
}

func TestResolveEntity_InvalidInput_ReturnsError(t *testing.T) {
    // Test error conditions
}
Table-Driven Tests
go
func TestEntityMatching(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []string
        wantErr  bool
    }{
        {
            name:     "exact match",
            input:    "light.living_room_main",
            expected: []string{"light.living_room_main"},
            wantErr:  false,
        },
        {
            name:     "fuzzy match",
            input:    "living lights",
            expected: []string{"light.living_room_main", "light.living_room_accent"},
            wantErr:  false,
        },
        {
            name:     "no match",
            input:    "nonexistent",
            expected: []string{},
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := ResolveEntity(context.Background(), tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ResolveEntity() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Additional assertions...
        })
    }
}
Performance Guidelines
Memory Allocation
go
// Pre-allocate slices when size is known
entities := make([]Entity, 0, len(rawData))

// Use sync.Pool for frequently allocated objects
var entityPool = sync.Pool{
    New: func() interface{} {
        return &Entity{}
    },
}

// Avoid string concatenation in loops
var builder strings.Builder
for _, item := range items {
    builder.WriteString(item)
}
result := builder.String()
Efficient Patterns
go
// Use range with index when you need both
for i, entity := range entities {
    // Use i and entity
}

// Use early returns to avoid unnecessary work
func expensiveOperation(data []Data) error {
    if len(data) == 0 {
        return nil // Skip expensive work
    }

    // Do expensive work...
    return nil
}
Code Organization
Import Grouping
go
import (
    // Standard library
    "context"
    "fmt"
    "net/http"

    // Third-party packages
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"

    // Local packages
    "github.com/user/hass-cli/internal/client"
    "github.com/user/hass-cli/internal/config"
)
Constants and Variables
go
// Group related constants
const (
    DefaultTimeout = 10 * time.Second
    MaxRetries     = 3
    APIVersion     = "v1"
)

// Use typed constants for enums
type VerbosityLevel int

const (
    VerbosityQuiet VerbosityLevel = iota
    VerbosityNormal
    VerbosityVerbose
)
Validation Rules
Build Requirements
All code must pass go vet
All code must pass golangci-lint run
All code must pass go fmt -d . (no formatting changes)
All code must have >= 80% test coverage
Code Quality Checks
No unused variables or imports
No shadowed variables
Proper error handling (no ignored errors)
Consistent naming conventions
Clear function and variable names
