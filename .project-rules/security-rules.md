Security Requirements
Credential Management
Access Token Security
Storage Priority: OS keyring > config file > environment variable
File Permissions: Config files must have 600 permissions (owner read/write only)
Memory Handling: Clear sensitive data from memory after use
Logging: Never log tokens, passwords, or other credentials
OS Keyring Integration
go
// Use keyring library for secure credential storage
import "github.com/zalando/go-keyring"

const (
    KeyringService = "hass-cli"
    KeyringUser    = "default"
)

// Store token securely
func StoreToken(token string) error {
    return keyring.Set(KeyringService, KeyringUser, token)
}

// Retrieve token securely
func GetToken() (string, error) {
    return keyring.Get(KeyringService, KeyringUser)
}
Token Validation
Validate token format before use
Test token with API call on startup
Handle token expiration gracefully
Prompt for new token when invalid
Input Validation
Command Line Arguments
Validate all user inputs before processing
Sanitize inputs that will be used in API calls
Implement length limits on user inputs
Reject obviously malicious inputs
Entity Name Validation
go
func ValidateEntityName(name string) error {
    // Prevent injection attacks
    if strings.ContainsAny(name, "<>\"'&;") {
        return errors.New("invalid characters in entity name")
    }

    // Reasonable length limits
    if len(name) > 255 {
        return errors.New("entity name too long")
    }

    return nil
}
Configuration Validation
Validate URLs before making requests
Check file paths for directory traversal
Validate timeout values are reasonable
Ensure numeric inputs are within expected ranges
Network Security
TLS/SSL Configuration
go
// Default to secure TLS configuration
func createHTTPClient(config Config) *http.Client {
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS12,
        // Only allow insecure connections if explicitly configured
        InsecureSkipVerify: config.AllowInsecureSSL,
    }

    transport := &http.Transport{
        TLSClientConfig: tlsConfig,
        // Additional security settings
        DisableKeepAlives: false,
        IdleConnTimeout:   30 * time.Second,
    }

    return &http.Client{
        Transport: transport,
        Timeout:   config.Timeout,
    }
}
Certificate Validation
Verify SSL certificates by default
Allow custom CA certificates when specified
Warn users when using insecure connections
Provide clear error messages for certificate issues
Request Security
Set appropriate timeouts for all requests
Use secure HTTP headers
Implement request rate limiting
Validate response content types
API Security
Authentication Headers
go
func addAuthHeaders(req *http.Request, token string) {
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("User-Agent", "hass-cli/"+version)
    req.Header.Set("Content-Type", "application/json")
}
Request Validation
Validate entity IDs before API calls
Sanitize parameters passed to Home Assistant
Check response status codes and handle errors
Validate JSON responses before parsing
Rate Limiting
Implement client-side rate limiting
Respect Home Assistant API limits
Use exponential backoff for retries
Monitor and log excessive API usage
Data Protection
Sensitive Data Handling
Clear sensitive variables after use
Avoid storing credentials in temporary files
Use secure memory allocation for credentials when possible
Implement proper cleanup in error paths
Logging Security
go
// Safe logging that excludes sensitive data
func logAPIRequest(method, url string, headers http.Header) {
    // Create safe headers map without Authorization
    safeHeaders := make(map[string]string)
    for k, v := range headers {
        if strings.ToLower(k) != "authorization" {
            safeHeaders[k] = strings.Join(v, ", ")
        }
    }

    log.Printf("API %s %s headers=%v", method, url, safeHeaders)
}
Cache Security
Encrypt cached credentials if stored on disk
Set appropriate file permissions for cache files
Implement cache expiration for sensitive data
Clear cache on application exit
Error Handling Security
Information Disclosure
Don't expose internal system information in error messages
Sanitize error messages before displaying to users
Log detailed errors internally, show generic messages to users
Avoid stack traces in production error output
Error Message Examples
go
// Good: Generic user-facing error
return errors.New("authentication failed")

// Bad: Exposes internal details
return fmt.Errorf("HTTP 401 from https://internal.homeassistant.local:8123/api/states: invalid token eyJ0eXAi...")

// Good: Detailed internal logging, generic user message
log.Printf("Authentication failed: HTTP %d from %s: %v", resp.StatusCode, req.URL, err)
return errors.New("authentication failed - check your access token")
File System Security
Configuration Files
Create config files with 600 permissions
Validate file paths to prevent directory traversal
Use secure temporary files when needed
Clean up temporary files in all exit paths
Directory Permissions
go
func ensureConfigDir(path string) error {
    // Create directory with secure permissions
    if err := os.MkdirAll(path, 0700); err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }

    // Verify permissions are correct
    info, err := os.Stat(path)
    if err != nil {
        return fmt.Errorf("failed to check directory permissions: %w", err)
    }

    if info.Mode().Perm() != 0700 {
        return fmt.Errorf("config directory has incorrect permissions: %o", info.Mode().Perm())
    }

    return nil
}
Runtime Security
Process Security
Drop unnecessary privileges when possible
Handle signals gracefully for clean shutdown
Implement timeout mechanisms to prevent hangs
Validate environment variables
Memory Security
Use defer statements for resource cleanup
Implement proper context cancellation
Avoid memory leaks in long-running operations
Clear sensitive data structures on exit
Security Testing
Test Requirements
Test with invalid/malicious inputs
Verify credential storage security
Test SSL/TLS configuration
Validate error message content
Security Test Examples
go
func TestCredentialSecurity(t *testing.T) {
    // Test that tokens are not logged
    var logOutput bytes.Buffer
    log.SetOutput(&logOutput)

    client := NewClient(Config{Token: "secret-token"})
    client.makeRequest(context.Background(), "GET", "/api/states")

    logContent := logOutput.String()
    if strings.Contains(logContent, "secret-token") {
        t.Error("Token found in log output - security violation")
    }
}
Compliance Considerations
Data Handling
Don't store more data than necessary
Implement data retention policies for logs
Consider privacy implications of cached data
Provide clear data usage documentation
Audit Trail
Log authentication events
Log configuration changes
Track API usage patterns
Implement log rotation and retention
