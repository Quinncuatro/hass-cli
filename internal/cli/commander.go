package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/quinncuatro/hass-cli/internal/config"
	"github.com/quinncuatro/hass-cli/internal/client"
	"github.com/quinncuatro/hass-cli/internal/entity"
	"github.com/quinncuatro/hass-cli/internal/tui"
)

type Commander struct {
	config   *config.Config
	client   *client.HomeAssistantClient
	resolver *entity.Resolver
}

func NewCommander(cfg *config.Config) *Commander {
	haClient := client.New(cfg)
	return &Commander{
		config:   cfg,
		client:   haClient,
		resolver: entity.NewResolver(cfg, haClient),
	}
}

func (c *Commander) Execute(args []string) error {
	if len(args) == 0 {
		return c.showHelp()
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "config":
		return c.handleConfigCommand(commandArgs)
	case "status":
		return c.handleStatusCommand(commandArgs)
	case "tui":
		return c.handleTUICommand(commandArgs)
	case "discover":
		return c.handleDiscoverCommand(commandArgs)
	case "automation":
		return c.handleAutomationCommand(commandArgs)
	case "scene":
		return c.handleSceneCommand(commandArgs)
	case "debug":
		return c.handleDebugCommand(commandArgs)
	case "help", "--help", "-h":
		return c.showHelp()
	case "version", "--version", "-v":
		return c.showVersion()
	default:
		return c.handleEntityCommand(args)
	}
}

func (c *Commander) handleConfigCommand(args []string) error {
	if len(args) == 0 {
		return c.showConfigHelp()
	}

	switch args[0] {
	case "init":
		return c.initConfig()
	case "show":
		return c.showConfig()
	case "test":
		return c.testConfig()
	default:
		return fmt.Errorf("unknown config command: %s", args[0])
	}
}

func (c *Commander) handleStatusCommand(args []string) error {
	if len(args) == 0 {
		return c.showSystemStatus()
	}
	
	query := strings.Join(args, " ")
	return c.showEntityStatus(query)
}

func (c *Commander) handleTUICommand(args []string) error {
	if c.config.HomeAssistant.URL == "" || c.config.HomeAssistant.Token == "" {
		return fmt.Errorf("home Assistant not configured. Run 'hass config init' to set up")
	}
	
	tuiApp := tui.NewApp(c.config, c.client)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	return tuiApp.Run(ctx)
}

func (c *Commander) handleDiscoverCommand(args []string) error {
	fmt.Println("Discovering Home Assistant instances...")
	fmt.Println("Checking common URLs:")
	fmt.Println("  http://homeassistant.local:8123")
	fmt.Println("  http://hassio.local:8123")
	fmt.Println("  http://192.168.1.100:8123")
	fmt.Println("  http://localhost:8123")
	fmt.Println()
	fmt.Println("Note: Full network discovery requires additional implementation")
	return nil
}

func (c *Commander) handleAutomationCommand(args []string) error {
	if len(args) == 0 {
		return c.listAutomations()
	}
	
	automationName := strings.Join(args, " ")
	return c.triggerAutomation(automationName)
}

func (c *Commander) handleSceneCommand(args []string) error {
	if len(args) == 0 {
		return c.listScenes()
	}
	
	sceneName := strings.Join(args, " ")
	return c.activateScene(sceneName)
}

func (c *Commander) handleDebugCommand(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	if len(args) == 0 {
		return c.showAllEntities(ctx)
	}

	switch args[0] {
	case "lights":
		return c.showLightEntities(ctx)
	case "match":
		if len(args) < 3 {
			return fmt.Errorf("usage: debug match <area> <entity-type>")
		}
		return c.debugEntityMatching(ctx, args[1], args[2])
	case "threshold":
		if len(args) < 2 {
			return fmt.Errorf("usage: debug threshold <value>")
		}
		return c.setFuzzyThreshold(args[1])
	default:
		return fmt.Errorf("unknown debug command: %s", args[0])
	}
}

func (c *Commander) showAllEntities(ctx context.Context) error {
	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	fmt.Printf("Total entities: %d\n\n", len(states))
	
	domainCounts := make(map[string]int)
	for _, state := range states {
		domain := strings.Split(state.EntityID, ".")[0]
		domainCounts[domain]++
	}

	fmt.Println("Entities by domain:")
	for domain, count := range domainCounts {
		fmt.Printf("  %s: %d\n", domain, count)
	}

	return nil
}

func (c *Commander) showLightEntities(ctx context.Context) error {
	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	fmt.Println("Light entities:")
	for _, state := range states {
		if strings.HasPrefix(state.EntityID, "light.") {
			friendlyName := state.EntityID
			if name, ok := state.Attributes["friendly_name"].(string); ok {
				friendlyName = name
			}
			
			fmt.Printf("  %s\n", state.EntityID)
			fmt.Printf("    friendly_name: %s\n", friendlyName)
			fmt.Printf("    state: %s\n", state.State)
			
			if areaID, ok := state.Attributes["area_id"].(string); ok {
				fmt.Printf("    area_id: %s\n", areaID)
			} else {
				fmt.Printf("    area_id: (not set)\n")
			}
			
			fmt.Println()
		}
	}

	return nil
}

func (c *Commander) debugEntityMatching(ctx context.Context, area, entityType string) error {
	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	fmt.Printf("Debug matching: area='%s', entityType='%s'\n\n", area, entityType)

	matches := c.resolver.DebugFindMatches(states, area, entityType, "")
	
	if len(matches) == 0 {
		fmt.Println("No matches found")
		return nil
	}

	fmt.Printf("Found %d matches:\n", len(matches))
	for i, match := range matches {
		fmt.Printf("%d. %s (score: %.3f)\n", i+1, match.FriendlyName, match.Score)
		fmt.Printf("   entity_id: %s\n", match.EntityID)
		fmt.Printf("   domain: %s\n", match.Domain)
		fmt.Printf("   area: %s\n", match.Area)
		fmt.Println()
	}

	return nil
}

func (c *Commander) setFuzzyThreshold(thresholdStr string) error {
	threshold, err := strconv.ParseFloat(thresholdStr, 64)
	if err != nil {
		return fmt.Errorf("invalid threshold value: %s", thresholdStr)
	}
	
	if threshold < 0.0 || threshold > 1.0 {
		return fmt.Errorf("threshold must be between 0.0 and 1.0")
	}
	
	c.config.Preferences.FuzzyThreshold = threshold
	fmt.Printf("Fuzzy threshold set to %.2f\n", threshold)
	return nil
}

func (c *Commander) handleEntityCommand(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	if len(args) < 2 {
		return fmt.Errorf("insufficient arguments for entity command")
	}

	var area, entityType, action, value string
	
	if len(args) == 2 {
		entityType = args[0]
		action = args[1]
	} else if len(args) == 3 {
		area = args[0]
		entityType = args[1]
		action = args[2]
	} else if len(args) >= 4 {
		area = args[0]
		entityType = args[1]
		action = args[2]
		value = strings.Join(args[3:], " ")
	}

	match, err := c.resolver.ResolveEntity(ctx, area, entityType, "")
	if err != nil {
		return fmt.Errorf("failed to resolve entity: %w", err)
	}

	fmt.Printf("🎯 Matched: %s (%s)\n", match.FriendlyName, match.EntityID)

	switch entity.ParseAction(action) {
	case "turn_on":
		err = c.client.TurnOnEntity(ctx, match.EntityID)
	case "turn_off":
		err = c.client.TurnOffEntity(ctx, match.EntityID)
	case "toggle":
		err = c.client.ToggleEntity(ctx, match.EntityID)
	default:
		if value != "" {
			err = c.handleEntityWithValue(ctx, match, action, value)
		} else {
			return fmt.Errorf("unsupported action: %s", action)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if c.config.Output.Verbosity > 0 {
		fmt.Printf("Successfully executed %s on %s (%s)\n", action, match.FriendlyName, match.EntityID)
	} else {
		fmt.Printf("✓ Turned %s %s (%s)\n", match.FriendlyName, action, match.EntityID)
	}

	return nil
}

func (c *Commander) handleEntityWithValue(ctx context.Context, match *entity.EntityMatch, action, value string) error {
	domain := match.Domain
	
	switch domain {
	case "light":
		return c.handleLightWithValue(ctx, match.EntityID, action, value)
	case "fan":
		return c.handleFanWithValue(ctx, match.EntityID, action, value)
	case "climate":
		return c.handleClimateWithValue(ctx, match.EntityID, action, value)
	case "cover":
		return c.handleCoverWithValue(ctx, match.EntityID, action, value)
	default:
		return fmt.Errorf("value-based actions not supported for domain: %s", domain)
	}
}

func (c *Commander) handleLightWithValue(ctx context.Context, entityID, action, value string) error {
	var serviceData map[string]interface{}

	switch strings.ToLower(action) {
	case "brightness", "bright", "dim":
		brightness, err := entity.ParseNumericValue(value)
		if err != nil {
			return fmt.Errorf("invalid brightness value: %s", value)
		}
		if brightness < 0 || brightness > 255 {
			return fmt.Errorf("brightness must be between 0 and 255")
		}
		serviceData = map[string]interface{}{
			"brightness": int(brightness),
		}
	case "color":
		serviceData = map[string]interface{}{
			"color_name": value,
		}
	default:
		return fmt.Errorf("unsupported light action: %s", action)
	}

	target := map[string]interface{}{
		"entity_id": entityID,
	}

	_, err := c.client.CallService(ctx, "light", "turn_on", target, serviceData)
	return err
}

func (c *Commander) handleFanWithValue(ctx context.Context, entityID, action, value string) error {
	var serviceData map[string]interface{}

	switch strings.ToLower(action) {
	case "speed", "percentage":
		speed, err := entity.ParseNumericValue(value)
		if err != nil {
			return fmt.Errorf("invalid speed value: %s", value)
		}
		if speed < 0 || speed > 100 {
			return fmt.Errorf("speed must be between 0 and 100")
		}
		serviceData = map[string]interface{}{
			"percentage": int(speed),
		}
	default:
		return fmt.Errorf("unsupported fan action: %s", action)
	}

	target := map[string]interface{}{
		"entity_id": entityID,
	}

	_, err := c.client.CallService(ctx, "fan", "set_percentage", target, serviceData)
	return err
}

func (c *Commander) handleClimateWithValue(ctx context.Context, entityID, action, value string) error {
	var serviceData map[string]interface{}

	switch strings.ToLower(action) {
	case "temp", "temperature":
		temp, err := entity.ParseNumericValue(value)
		if err != nil {
			return fmt.Errorf("invalid temperature value: %s", value)
		}
		serviceData = map[string]interface{}{
			"temperature": temp,
		}
	case "mode":
		serviceData = map[string]interface{}{
			"hvac_mode": value,
		}
	default:
		return fmt.Errorf("unsupported climate action: %s", action)
	}

	target := map[string]interface{}{
		"entity_id": entityID,
	}

	serviceName := "set_temperature"
	if action == "mode" {
		serviceName = "set_hvac_mode"
	}

	_, err := c.client.CallService(ctx, "climate", serviceName, target, serviceData)
	return err
}

func (c *Commander) handleCoverWithValue(ctx context.Context, entityID, action, value string) error {
	var serviceData map[string]interface{}

	switch strings.ToLower(action) {
	case "position", "pos":
		position, err := entity.ParseNumericValue(value)
		if err != nil {
			return fmt.Errorf("invalid position value: %s", value)
		}
		if position < 0 || position > 100 {
			return fmt.Errorf("position must be between 0 and 100")
		}
		serviceData = map[string]interface{}{
			"position": int(position),
		}
	default:
		return fmt.Errorf("unsupported cover action: %s", action)
	}

	target := map[string]interface{}{
		"entity_id": entityID,
	}

	_, err := c.client.CallService(ctx, "cover", "set_cover_position", target, serviceData)
	return err
}

func (c *Commander) showHelp() error {
	help := `Home Assistant CLI Tool

Usage:
  hass [command] [options]
  hass [area] [entity-type] [action] [value]

Commands:
  config      Configuration management
  status      Show entity or system status
  tui         Interactive terminal interface
  discover    Discover Home Assistant instances
  automation  Trigger automations
  scene       Activate scenes
  help        Show this help message
  version     Show version information

Entity Control Examples:
  hass living lights on              Turn on living room lights
  hass kitchen fan speed 75          Set kitchen fan to 75% speed
  hass bedroom climate temp 72       Set bedroom temperature to 72°F

For more information, visit: https://github.com/quinncuatro/hass-cli`

	fmt.Println(help)
	return nil
}

func (c *Commander) showVersion() error {
	fmt.Println("hass-cli v0.1.0")
	fmt.Println("A command-line interface for Home Assistant")
	fmt.Println("https://github.com/quinncuatro/hass-cli")
	return nil
}

func (c *Commander) showConfigHelp() error {
	help := `Config Commands:
  init    Initialize configuration with setup wizard
  show    Display current configuration
  test    Test connection to Home Assistant`

	fmt.Println(help)
	return nil
}

func (c *Commander) initConfig() error {
	fmt.Println("Home Assistant CLI Configuration Setup")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("This wizard will help you configure hass-cli to connect to your Home Assistant instance.")
	fmt.Println()
	fmt.Println("You'll need:")
	fmt.Println("1. Your Home Assistant URL (e.g., http://homeassistant.local:8123)")
	fmt.Println("2. A long-lived access token from Home Assistant")
	fmt.Println()
	fmt.Println("To create a token:")
	fmt.Println("1. Open Home Assistant in your browser")
	fmt.Println("2. Go to Profile → Security → Long-lived access tokens")
	fmt.Println("3. Click 'Create Token'")
	fmt.Println("4. Enter a name like 'CLI Tool'")
	fmt.Println("5. Copy the token")
	fmt.Println()
	fmt.Println("Example configuration file will be created at:")
	
	configDir, _ := os.UserConfigDir()
	configPath := filepath.Join(configDir, "hass", "config.yaml")
	fmt.Printf("  %s\n", configPath)
	fmt.Println()
	fmt.Println("Note: Interactive configuration setup requires additional implementation")
	fmt.Println("For now, manually create the config file with:")
	fmt.Println()
	fmt.Println("homeassistant:")
	fmt.Println("  url: \"http://homeassistant.local:8123\"")
	fmt.Println("  token: \"your-long-lived-access-token\"")
	fmt.Println("  timeout: \"10s\"")
	fmt.Println()
	fmt.Println("aliases:")
	fmt.Println("  lr: \"living room\"")
	fmt.Println("  br: \"bedroom\"")
	fmt.Println()
	fmt.Println("preferences:")
	fmt.Println("  fuzzy_threshold: 0.6")
	
	return nil
}

func (c *Commander) showConfig() error {
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")
	fmt.Printf("Home Assistant URL: %s\n", c.config.HomeAssistant.URL)
	
	if c.config.HomeAssistant.Token != "" {
		maskedToken := c.config.HomeAssistant.Token[:8] + "***MASKED***"
		fmt.Printf("Token: %s\n", maskedToken)
	} else {
		fmt.Println("Token: Not set")
	}
	
	fmt.Printf("Timeout: %s\n", c.config.HomeAssistant.Timeout)
	fmt.Printf("Skip TLS Verify: %t\n", c.config.HomeAssistant.SkipTLSVerify)
	fmt.Printf("Cache Timeout: %s\n", c.config.HomeAssistant.CacheTimeout)
	fmt.Println()
	
	if len(c.config.Aliases) > 0 {
		fmt.Println("Aliases:")
		for alias, target := range c.config.Aliases {
			fmt.Printf("  %s → %s\n", alias, target)
		}
		fmt.Println()
	}
	
	fmt.Printf("Output Format: %s\n", c.config.Output.Format)
	fmt.Printf("Color Output: %t\n", c.config.Output.Color)
	fmt.Printf("Verbosity: %d\n", c.config.Output.Verbosity)
	fmt.Printf("Fuzzy Threshold: %.2f\n", c.config.Preferences.FuzzyThreshold)
	
	return nil
}

func (c *Commander) testConfig() error {
	fmt.Println("Testing Home Assistant Connection...")
	fmt.Println("===================================")
	
	if c.config.HomeAssistant.URL == "" {
		return fmt.Errorf("❌ Home Assistant URL not configured")
	}
	
	if c.config.HomeAssistant.Token == "" {
		return fmt.Errorf("❌ Home Assistant token not configured")
	}
	
	fmt.Printf("✓ Configuration file found\n")
	fmt.Printf("✓ URL configured: %s\n", c.config.HomeAssistant.URL)
	fmt.Printf("✓ Token configured\n")
	
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()
	
	fmt.Println("Testing connection...")
	
	status, err := c.client.GetSystemStatus(ctx)
	if err != nil {
		return fmt.Errorf("❌ Connection failed: %w", err)
	}
	
	fmt.Printf("✓ Connection successful\n")
	fmt.Printf("✓ Authentication valid\n")
	fmt.Printf("✓ Home Assistant version: %s\n", status.Version)
	fmt.Printf("✓ Location: %s\n", status.LocationName)
	
	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("❌ Failed to fetch entities: %w", err)
	}
	
	fmt.Printf("✓ Found %d entities\n", len(states))
	
	domainCounts := make(map[string]int)
	for _, state := range states {
		domain := strings.Split(state.EntityID, ".")[0]
		domainCounts[domain]++
	}
	
	fmt.Println("\nEntity breakdown:")
	for domain, count := range domainCounts {
		fmt.Printf("  %s: %d\n", domain, count)
	}
	
	fmt.Println("\n✅ All tests passed! Your configuration is working correctly.")
	
	return nil
}

func (c *Commander) showSystemStatus() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	if c.config.HomeAssistant.URL == "" || c.config.HomeAssistant.Token == "" {
		return fmt.Errorf("home Assistant not configured. Run 'hass config init' to set up")
	}

	status, err := c.client.GetSystemStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system status: %w", err)
	}

	fmt.Printf("Home Assistant Status:\n")
	fmt.Printf("  Version: %s\n", status.Version)
	fmt.Printf("  State: %s\n", status.State)
	fmt.Printf("  Location: %s\n", status.LocationName)
	fmt.Printf("  Timezone: %s\n", status.Timezone)
	fmt.Printf("  Unit System: Temperature: %s, Length: %s, Mass: %s, Volume: %s\n", 
		status.UnitSystem.Temperature, status.UnitSystem.Length, status.UnitSystem.Mass, status.UnitSystem.Volume)
	fmt.Printf("  External URL: %s\n", status.ExternalURL)
	fmt.Printf("  Internal URL: %s\n", status.InternalURL)
	fmt.Printf("  Safe Mode: %t\n", status.SafeMode)
	fmt.Printf("  Recovery Mode: %t\n", status.RecoveryMode)

	return nil
}

func (c *Commander) showEntityStatus(query string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	if query == "" {
		states, err := c.client.GetStates(ctx)
		if err != nil {
			return fmt.Errorf("failed to get entity states: %w", err)
		}

		fmt.Printf("Total entities: %d\n\n", len(states))
		
		domainCounts := make(map[string]int)
		for _, state := range states {
			domain := strings.Split(state.EntityID, ".")[0]
			domainCounts[domain]++
		}

		fmt.Println("Entities by domain:")
		for domain, count := range domainCounts {
			fmt.Printf("  %s: %d\n", domain, count)
		}

		return nil
	}

	match, err := c.resolver.ResolveEntity(ctx, "", "", query)
	if err != nil {
		return fmt.Errorf("failed to resolve entity: %w", err)
	}

	state, err := c.client.GetState(ctx, match.EntityID)
	if err != nil {
		return fmt.Errorf("failed to get entity state: %w", err)
	}

	fmt.Printf("Entity: %s (%s)\n", match.FriendlyName, match.EntityID)
	fmt.Printf("State: %s\n", state.State)
	fmt.Printf("Domain: %s\n", match.Domain)
	if match.Area != "" {
		fmt.Printf("Area: %s\n", match.Area)
	}
	fmt.Printf("Last Changed: %s\n", state.LastChanged.Format(time.RFC3339))
	fmt.Printf("Last Updated: %s\n", state.LastUpdated.Format(time.RFC3339))

	if c.config.Output.Verbosity > 1 {
		fmt.Println("\nAttributes:")
		for key, value := range state.Attributes {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	return nil
}

func (c *Commander) listAutomations() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	fmt.Println("Available Automations:")
	found := false
	for _, state := range states {
		if strings.HasPrefix(state.EntityID, "automation.") {
			friendlyName := state.EntityID
			if name, ok := state.Attributes["friendly_name"].(string); ok {
				friendlyName = name
			}
			status := "enabled"
			if state.State == "off" {
				status = "disabled"
			}
			fmt.Printf("  %s (%s) - %s\n", friendlyName, state.EntityID, status)
			found = true
		}
	}

	if !found {
		fmt.Println("  No automations found")
	}

	return nil
}

func (c *Commander) triggerAutomation(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	var automationID string
	normalizedName := strings.ToLower(name)

	for _, state := range states {
		if strings.HasPrefix(state.EntityID, "automation.") {
			friendlyName := state.EntityID
			if fn, ok := state.Attributes["friendly_name"].(string); ok {
				friendlyName = fn
			}

			if strings.ToLower(friendlyName) == normalizedName ||
				strings.Contains(strings.ToLower(friendlyName), normalizedName) ||
				state.EntityID == name {
				automationID = state.EntityID
				break
			}
		}
	}

	if automationID == "" {
		return fmt.Errorf("automation '%s' not found", name)
	}

	target := map[string]interface{}{
		"entity_id": automationID,
	}

	_, err = c.client.CallService(ctx, "automation", "trigger", target, nil)
	if err != nil {
		return fmt.Errorf("failed to trigger automation: %w", err)
	}

	if c.config.Output.Verbosity > 0 {
		fmt.Printf("Successfully triggered automation: %s\n", automationID)
	}

	return nil
}

func (c *Commander) listScenes() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	fmt.Println("Available Scenes:")
	found := false
	for _, state := range states {
		if strings.HasPrefix(state.EntityID, "scene.") {
			friendlyName := state.EntityID
			if name, ok := state.Attributes["friendly_name"].(string); ok {
				friendlyName = name
			}
			fmt.Printf("  %s (%s)\n", friendlyName, state.EntityID)
			found = true
		}
	}

	if !found {
		fmt.Println("  No scenes found")
	}

	return nil
}

func (c *Commander) activateScene(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.HomeAssistant.Timeout)
	defer cancel()

	states, err := c.client.GetStates(ctx)
	if err != nil {
		return fmt.Errorf("failed to get states: %w", err)
	}

	var sceneID string
	normalizedName := strings.ToLower(name)

	for _, state := range states {
		if strings.HasPrefix(state.EntityID, "scene.") {
			friendlyName := state.EntityID
			if fn, ok := state.Attributes["friendly_name"].(string); ok {
				friendlyName = fn
			}

			if strings.ToLower(friendlyName) == normalizedName ||
				strings.Contains(strings.ToLower(friendlyName), normalizedName) ||
				state.EntityID == name {
				sceneID = state.EntityID
				break
			}
		}
	}

	if sceneID == "" {
		return fmt.Errorf("scene '%s' not found", name)
	}

	target := map[string]interface{}{
		"entity_id": sceneID,
	}

	_, err = c.client.CallService(ctx, "scene", "turn_on", target, nil)
	if err != nil {
		return fmt.Errorf("failed to activate scene: %w", err)
	}

	if c.config.Output.Verbosity > 0 {
		fmt.Printf("Successfully activated scene: %s\n", sceneID)
	}

	return nil
}