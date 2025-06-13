package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HomeAssistant HomeAssistantConfig `yaml:"homeassistant"`
	Aliases       map[string]string   `yaml:"aliases"`
	Preferences   PreferencesConfig   `yaml:"preferences"`
	Output        OutputConfig        `yaml:"output"`
	Discovery     DiscoveryConfig     `yaml:"discovery"`
	Security      SecurityConfig      `yaml:"security"`
}

type HomeAssistantConfig struct {
	URL            string        `yaml:"url"`
	Token          string        `yaml:"token"`
	Timeout        time.Duration `yaml:"timeout"`
	SkipTLSVerify  bool          `yaml:"skip_tls_verify"`
	CacheTimeout   time.Duration `yaml:"cache_timeout"`
}

type PreferencesConfig struct {
	FuzzyThreshold float64 `yaml:"fuzzy_threshold"`
	DefaultArea    string  `yaml:"default_area"`
	AutoDiscovery  bool    `yaml:"auto_discovery"`
}

type OutputConfig struct {
	Format    string `yaml:"format"`
	Color     bool   `yaml:"color"`
	Verbosity int    `yaml:"verbosity"`
}

type DiscoveryConfig struct {
	Enabled       bool          `yaml:"enabled"`
	Timeout       time.Duration `yaml:"timeout"`
	CustomPorts   []int         `yaml:"custom_ports"`
	IPRanges      []string      `yaml:"ip_ranges"`
	MDNSEnabled   bool          `yaml:"mdns_enabled"`
}

type SecurityConfig struct {
	UseKeyring       bool `yaml:"use_keyring"`
	KeyringService   string `yaml:"keyring_service"`
	ConfigFilePerms  os.FileMode `yaml:"config_file_perms"`
}

func DefaultConfig() *Config {
	return &Config{
		HomeAssistant: HomeAssistantConfig{
			Timeout:       10 * time.Second,
			SkipTLSVerify: false,
			CacheTimeout:  5 * time.Minute,
		},
		Aliases: make(map[string]string),
		Preferences: PreferencesConfig{
			FuzzyThreshold: 0.5,
			AutoDiscovery:  true,
		},
		Output: OutputConfig{
			Format:    "text",
			Color:     true,
			Verbosity: 1,
		},
		Discovery: DiscoveryConfig{
			Enabled:     true,
			Timeout:     30 * time.Second,
			CustomPorts: []int{8123, 8124},
			MDNSEnabled: true,
		},
		Security: SecurityConfig{
			UseKeyring:      true,
			KeyringService:  "hass-cli",
			ConfigFilePerms: 0600,
		},
	}
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	cfg := DefaultConfig()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, c.Security.ConfigFilePerms); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "hass", "config.yaml"), nil
}