package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/quinncuatro/hass-cli/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := &config.Config{
		HomeAssistant: config.HomeAssistantConfig{
			URL:           "http://test.local:8123",
			Token:         "test-token",
			Timeout:       5 * time.Second,
			SkipTLSVerify: true,
		},
	}

	client := New(cfg)

	if client.baseURL != "http://test.local:8123" {
		t.Errorf("expected baseURL to be http://test.local:8123, got %s", client.baseURL)
	}

	if client.token != "test-token" {
		t.Errorf("expected token to be test-token, got %s", client.token)
	}

	if client.httpClient.Timeout != 5*time.Second {
		t.Errorf("expected timeout to be 5s, got %v", client.httpClient.Timeout)
	}
}

func TestGetSystemStatus(t *testing.T) {
	mockStatus := SystemStatus{
		Message:      "API running.",
		Version:      "2024.1.0",
		State:        "RUNNING",
		Timezone:     "America/New_York",
		LocationName: "Home",
		UnitSystem: UnitSystem{
			Length:      "km",
			Mass:        "g",
			Temperature: "°C",
			Volume:      "L",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/config" {
			t.Errorf("expected path /api/config, got %s", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-token" {
			t.Errorf("expected Authorization header 'Bearer test-token', got %s", authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockStatus); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		HomeAssistant: config.HomeAssistantConfig{
			URL:     server.URL,
			Token:   "test-token",
			Timeout: 5 * time.Second,
		},
	}

	client := New(cfg)
	status, err := client.GetSystemStatus(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Version != "2024.1.0" {
		t.Errorf("expected version 2024.1.0, got %s", status.Version)
	}

	if status.State != "RUNNING" {
		t.Errorf("expected state RUNNING, got %s", status.State)
	}
}

func TestGetStates(t *testing.T) {
	mockStates := []EntityState{
		{
			EntityID: "light.living_room",
			State:    "on",
			Attributes: map[string]interface{}{
				"friendly_name": "Living Room Light",
				"brightness":    255,
			},
		},
		{
			EntityID: "switch.kitchen_coffee",
			State:    "off",
			Attributes: map[string]interface{}{
				"friendly_name": "Kitchen Coffee Maker",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/states" {
			t.Errorf("expected path /api/states, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mockStates); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		HomeAssistant: config.HomeAssistantConfig{
			URL:     server.URL,
			Token:   "test-token",
			Timeout: 5 * time.Second,
		},
	}

	client := New(cfg)
	states, err := client.GetStates(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(states) != 2 {
		t.Errorf("expected 2 states, got %d", len(states))
	}

	if states[0].EntityID != "light.living_room" {
		t.Errorf("expected first entity ID to be light.living_room, got %s", states[0].EntityID)
	}

	if states[1].EntityID != "switch.kitchen_coffee" {
		t.Errorf("expected second entity ID to be switch.kitchen_coffee, got %s", states[1].EntityID)
	}
}

func TestCallService(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/services/light/turn_on" {
			t.Errorf("expected path /api/services/light/turn_on, got %s", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		var request ServiceCallRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if request.Domain != "light" {
			t.Errorf("expected domain light, got %s", request.Domain)
		}

		if request.Service != "turn_on" {
			t.Errorf("expected service turn_on, got %s", request.Service)
		}

		response := ServiceCallResponse{
			Context: struct {
				ID       string `json:"id"`
				ParentID string `json:"parent_id"`
				UserID   string `json:"user_id"`
			}{
				ID:     "test-context-id",
				UserID: "test-user-id",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		HomeAssistant: config.HomeAssistantConfig{
			URL:     server.URL,
			Token:   "test-token",
			Timeout: 5 * time.Second,
		},
	}

	client := New(cfg)
	target := map[string]interface{}{
		"entity_id": "light.living_room",
	}

	response, err := client.CallService(context.Background(), "light", "turn_on", target, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Context.ID != "test-context-id" {
		t.Errorf("expected context ID test-context-id, got %s", response.Context.ID)
	}
}

func TestTurnOnEntity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/services/light/turn_on" {
			t.Errorf("expected path /api/services/light/turn_on, got %s", r.URL.Path)
		}

		var request ServiceCallRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		entityID, ok := request.Target["entity_id"].(string)
		if !ok || entityID != "light.living_room" {
			t.Errorf("expected entity_id light.living_room in target, got %v", request.Target)
		}

		response := ServiceCallResponse{}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := &config.Config{
		HomeAssistant: config.HomeAssistantConfig{
			URL:     server.URL,
			Token:   "test-token",
			Timeout: 5 * time.Second,
		},
	}

	client := New(cfg)
	err := client.TurnOnEntity(context.Background(), "light.living_room")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitSystemUnmarshal(t *testing.T) {
	jsonData := `{
		"message": "API running.",
		"version": "2024.1.0",
		"state": "RUNNING",
		"timezone": "America/New_York",
		"location_name": "Home",
		"unit_system": {
			"length": "km",
			"mass": "g",
			"temperature": "°C",
			"volume": "L"
		},
		"config_dir": "/config",
		"external_url": "http://homeassistant.local:8123",
		"internal_url": "http://homeassistant.local:8123",
		"safe_mode": false,
		"config_source": "yaml",
		"recovery_mode": false,
		"supports_statistics": true
	}`

	var status SystemStatus
	err := json.Unmarshal([]byte(jsonData), &status)
	if err != nil {
		t.Fatalf("failed to unmarshal SystemStatus: %v", err)
	}

	if status.UnitSystem.Length != "km" {
		t.Errorf("expected unit system length to be km, got %s", status.UnitSystem.Length)
	}

	if status.UnitSystem.Mass != "g" {
		t.Errorf("expected unit system mass to be g, got %s", status.UnitSystem.Mass)
	}

	if status.UnitSystem.Temperature != "°C" {
		t.Errorf("expected unit system temperature to be °C, got %s", status.UnitSystem.Temperature)
	}

	if status.UnitSystem.Volume != "L" {
		t.Errorf("expected unit system volume to be L, got %s", status.UnitSystem.Volume)
	}

	if status.Version != "2024.1.0" {
		t.Errorf("expected version 2024.1.0, got %s", status.Version)
	}
}