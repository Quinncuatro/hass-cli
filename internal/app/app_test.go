package app

import (
	"testing"
)

func TestRun_WithNoArgs_ShowsHelp(t *testing.T) {
	// Test that running with no arguments shows help
	err := Run([]string{})
	if err != nil {
		t.Fatalf("Expected no error for help, got: %v", err)
	}
}

func TestRun_WithHelpFlag_ShowsHelp(t *testing.T) {
	// Test help flag
	err := Run([]string{"--help"})
	if err != nil {
		t.Fatalf("Expected no error for help flag, got: %v", err)
	}
}

func TestRun_WithVersionFlag_ShowsVersion(t *testing.T) {
	// Test version flag
	err := Run([]string{"--version"})
	if err != nil {
		t.Fatalf("Expected no error for version flag, got: %v", err)
	}
}

func TestVersionInfo(t *testing.T) {
	info := VersionInfo()
	if info == "" {
		t.Fatal("Expected version info to not be empty")
	}
	
	if !contains(info, "hass-cli") {
		t.Fatal("Expected version info to contain 'hass-cli'")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}