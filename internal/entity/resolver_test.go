package entity

import (
	"testing"
	"time"

	"github.com/quinncuatro/hass-cli/internal/client"
	"github.com/quinncuatro/hass-cli/internal/config"
)


func TestParseEntityType(t *testing.T) {
	tests := []struct {
		input    string
		expected EntityType
	}{
		{"light", EntityTypeLight},
		{"lights", EntityTypeLight},
		{"lamp", EntityTypeLight},
		{"switch", EntityTypeSwitch},
		{"switches", EntityTypeSwitch},
		{"fan", EntityTypeFan},
		{"fans", EntityTypeFan},
		{"climate", EntityTypeClimate},
		{"thermostat", EntityTypeClimate},
		{"cover", EntityTypeCover},
		{"blinds", EntityTypeCover},
		{"sensor", EntityTypeSensor},
		{"unknown", EntityTypeUnknown},
	}

	for _, test := range tests {
		result := ParseEntityType(test.input)
		if result != test.expected {
			t.Errorf("ParseEntityType(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseAction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"on", "turn_on"},
		{"turn_on", "turn_on"},
		{"enable", "turn_on"},
		{"off", "turn_off"},
		{"turn_off", "turn_off"},
		{"disable", "turn_off"},
		{"toggle", "toggle"},
		{"switch", "toggle"},
		{"custom", "custom"},
	}

	for _, test := range tests {
		result := ParseAction(test.input)
		if result != test.expected {
			t.Errorf("ParseAction(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestParseNumericValue(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"123", 123.0, false},
		{"123.45", 123.45, false},
		{"0", 0.0, false},
		{"255", 255.0, false},
		{"abc", 0.0, true},
		{"", 0.0, true},
	}

	for _, test := range tests {
		result, err := ParseNumericValue(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("ParseNumericValue(%s) expected error but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseNumericValue(%s) unexpected error: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("ParseNumericValue(%s) = %f, expected %f", test.input, result, test.expected)
			}
		}
	}
}

func TestResolverScoreDomain(t *testing.T) {
	cfg := config.DefaultConfig()
	resolver := &Resolver{config: cfg}

	tests := []struct {
		domain     string
		entityType string
		expected   float64
	}{
		{"light", "light", 1.0},
		{"light", "lights", 1.0},
		{"light", "lamp", 1.0},
		{"switch", "switch", 1.0},
		{"fan", "fan", 1.0},
		{"climate", "thermostat", 1.0},
		{"cover", "blinds", 1.0},
		{"sensor", "sensor", 1.0},
		{"light", "switch", 0.0},
		{"unknown", "light", 0.0},
	}

	for _, test := range tests {
		result := resolver.scoreDomain(test.domain, test.entityType)
		if result != test.expected {
			t.Errorf("scoreDomain(%s, %s) = %f, expected %f", test.domain, test.entityType, result, test.expected)
		}
	}
}

func TestResolverFuzzyMatch(t *testing.T) {
	cfg := config.DefaultConfig()
	resolver := &Resolver{config: cfg}

	tests := []struct {
		s1       string
		s2       string
		expected bool // whether score should be > 0.6
	}{
		{"living room", "living", true},
		{"living room", "room", true},
		{"kitchen light", "kitchen", true},
		{"bedroom fan", "bedroom", true},
		{"completely different", "nothing", false},
		{"exact match", "exact match", true},
	}

	for _, test := range tests {
		score := resolver.fuzzyMatch(test.s1, test.s2)
		isAboveThreshold := score > 0.6
		if isAboveThreshold != test.expected {
			t.Errorf("fuzzyMatch(%s, %s) = %f (above threshold: %t), expected above threshold: %t", 
				test.s1, test.s2, score, isAboveThreshold, test.expected)
		}
	}
}

func TestResolverScoreEntity(t *testing.T) {
	cfg := config.DefaultConfig()
	resolver := &Resolver{config: cfg}

	state := client.EntityState{
		EntityID: "light.living_room_lamp",
		State:    "on",
		Attributes: map[string]interface{}{
			"friendly_name": "Living Room Lamp",
			"area_id":      "living_room",
		},
		LastChanged: time.Now(),
		LastUpdated: time.Now(),
	}

	match := resolver.scoreEntity(state, "living", "light", "lamp")

	if match.EntityID != "light.living_room_lamp" {
		t.Errorf("expected EntityID to be light.living_room_lamp, got %s", match.EntityID)
	}

	if match.FriendlyName != "Living Room Lamp" {
		t.Errorf("expected FriendlyName to be Living Room Lamp, got %s", match.FriendlyName)
	}

	if match.Domain != "light" {
		t.Errorf("expected Domain to be light, got %s", match.Domain)
	}

	if match.Area != "living_room" {
		t.Errorf("expected Area to be living_room, got %s", match.Area)
	}

	if match.Score <= 0.6 {
		t.Errorf("expected Score to be above threshold, got %f", match.Score)
	}
}