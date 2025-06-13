package config

import (
	"os"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.HomeAssistant.Timeout != 10*time.Second {
		t.Errorf("expected timeout to be 10s, got %v", cfg.HomeAssistant.Timeout)
	}

	if cfg.HomeAssistant.SkipTLSVerify != false {
		t.Errorf("expected SkipTLSVerify to be false, got %v", cfg.HomeAssistant.SkipTLSVerify)
	}

	if cfg.Preferences.FuzzyThreshold != 0.6 {
		t.Errorf("expected FuzzyThreshold to be 0.6, got %v", cfg.Preferences.FuzzyThreshold)
	}

	if cfg.Output.Verbosity != 1 {
		t.Errorf("expected Verbosity to be 1, got %v", cfg.Output.Verbosity)
	}

	if cfg.Security.UseKeyring != true {
		t.Errorf("expected UseKeyring to be true, got %v", cfg.Security.UseKeyring)
	}
}

func TestConfigLoadNonExistent(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	
	_ = os.Setenv("HOME", tempDir)
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error loading non-existent config, got %v", err)
	}

	if cfg.HomeAssistant.Timeout != 10*time.Second {
		t.Errorf("expected default timeout, got %v", cfg.HomeAssistant.Timeout)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	
	_ = os.Setenv("HOME", tempDir)
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	cfg := DefaultConfig()
	cfg.HomeAssistant.URL = "http://test.local:8123"
	cfg.HomeAssistant.Token = "test-token"
	cfg.Preferences.FuzzyThreshold = 0.8

	err := cfg.Save()
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loadedCfg.HomeAssistant.URL != "http://test.local:8123" {
		t.Errorf("expected URL to be http://test.local:8123, got %s", loadedCfg.HomeAssistant.URL)
	}

	if loadedCfg.HomeAssistant.Token != "test-token" {
		t.Errorf("expected token to be test-token, got %s", loadedCfg.HomeAssistant.Token)
	}

	if loadedCfg.Preferences.FuzzyThreshold != 0.8 {
		t.Errorf("expected FuzzyThreshold to be 0.8, got %v", loadedCfg.Preferences.FuzzyThreshold)
	}
}