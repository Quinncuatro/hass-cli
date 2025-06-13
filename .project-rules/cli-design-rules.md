CLI Design Principles
Command Structure Philosophy
Natural Language Design
Commands should read like natural English sentences
hass living lights on instead of hass --area=living --type=light --action=on
Support common abbreviations and shortcuts
Implement intelligent defaults to reduce verbosity
Progressive Disclosure
Simple commands for common tasks
Advanced options available but not overwhelming
Help system that scales with user expertise
Context-sensitive suggestions
Command Parsing Strategy
Argument Order Priority
Area/Location: living, bedroom, kitchen
Entity Type: lights, fan, thermostat
Action: on, off, toggle, set
Value: 50, red, 72°F
Flexible Parsing Rules
go
// Support multiple valid formats
"hass living lights on"           // Standard format
"hass lights living on"           // Flexible order
"hass lr lights on"              // Area aliases
"hass living room main lights on" // Specific entity targeting
Disambiguation Strategy
go
type CommandContext struct {
    Area        string   // Resolved area name
    EntityType  string   // Resolved entity domain
    Action      string   // Resolved action
    Value       string   // Optional value parameter
    Matches     []Entity // Potential entity matches
    Confidence  float64  // Matching confidence score
}
User Experience Standards
Output Verbosity Levels
go
const (
    VerbosityQuiet  = 0 // Minimal output, errors only
    VerbosityNormal = 1 // Standard confirmations (default)
    VerbosityVerbose = 2 // Detailed information
)

// Example outputs for same action:
// Quiet (0):   "✓"
// Normal (1):  "✓ Living room lights turned on"
// Verbose (2): "✓ Living room lights turned on (3 entities, 234ms)"
Success Feedback
Immediate confirmation for actions
Clear indication of what changed
Number of affected entities when relevant
Response time for verbose mode
Error Messages
Clear problem description
Actionable suggestions when possible
Context about what was attempted
Examples of correct usage
Help System Design
Contextual Help
bash
# Command-specific help
hass lights --help          # Help for light commands
hass config --help          # Help for configuration commands
hass --help                 # General help

# Context-aware suggestions
$ hass livng lights on
No exact match for "livng". Did you mean:
  living → hass living lights on
Interactive Elements
Use prompts for ambiguous commands
Numbered choices for multiple matches
Confirmation for destructive actions
Progress indicators for slow operations
Error Handling Standards
Error Categories
User Input Errors: Typos, invalid syntax
Configuration Errors: Missing config, invalid tokens
Network Errors: Connection issues, timeouts
API Errors: Home Assistant API failures
System Errors: File permissions, OS issues
Error Message Format
✗ Problem description
  Suggestion or next step
  Example: hass bedroom lights on
Performance Requirements
Response Time Targets
Command parsing: < 10ms
Entity resolution: < 50ms
API calls: < 500ms (local network)
Cache lookups: < 10ms
Cold start: < 100ms
Optimization Strategies
Cache entity data locally
Parallel entity resolution
Connection pooling for API calls
Efficient string matching algorithms
Accessibility Standards
Screen Reader Support
Descriptive output messages
Structured information presentation
Clear success/failure indicators
Logical reading order
Color Usage
Default to no color in non-TTY environments
Respect NO_COLOR environment variable
Provide --no-color flag
Use color as enhancement, not requirement
Configuration Integration
Config File Interaction
bash
hass config init           # Interactive setup
hass config show          # Display current config
hass config set key value # Update configuration
hass config alias add lr "living room"
Environment Variable Support
HASS_URL, HASS_TOKEN for quick overrides
HASS_CONFIG for custom config file path
Standard CLI environment variables (NO_COLOR, etc.)
Command Completion
Shell Integration
Bash completion support
Zsh completion support
Fish completion support
PowerShell completion support
Dynamic Completion
Complete area names from Home Assistant
Complete entity names and types
Complete available actions
Complete configuration keys
Testing Standards
CLI Testing Requirements
Test all command parsing scenarios
Validate error message clarity
Test help system completeness
Benchmark command response times
User Experience Testing
go
func TestCommandParsing(t *testing.T) {
    tests := []struct {
        name     string
        args     []string
        expected CommandContext
    }{
        {
            name: "natural language command",
            args: []string{"living", "lights", "on"},
            expected: CommandContext{
                Area:   "living room",
                Type:   "light",
                Action: "turn_on",
            },
        },
        {
            name: "abbreviated command",
            args: []string{"lr", "lights", "on"},
            expected: CommandContext{
                Area:   "living room",
                Type:   "light",
                Action: "turn_on",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseCommand(tt.args)
            assert.Equal(t, tt.expected, result)
        })
    }
}
Validation Requirements
Command Validation
Validate all user inputs before processing
Provide clear error messages for invalid commands
Suggest corrections for typos and near-misses
Handle edge cases gracefully
Output Validation
Test output formatting at all verbosity levels
Ensure consistent message formatting
Validate color usage and accessibility
Test non-TTY output (pipes, redirects)
