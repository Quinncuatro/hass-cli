package entity

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/quinncuatro/hass-cli/internal/client"
	"github.com/quinncuatro/hass-cli/internal/config"
)

type Resolver struct {
	config *config.Config
	client *client.HomeAssistantClient
}

type EntityMatch struct {
	EntityID   string
	FriendlyName string
	Domain     string
	Area       string
	Score      float64
}

type EntityType int

const (
	EntityTypeLight EntityType = iota
	EntityTypeSwitch
	EntityTypeFan
	EntityTypeClimate
	EntityTypeCover
	EntityTypeSensor
	EntityTypeUnknown
)

func (et EntityType) String() string {
	switch et {
	case EntityTypeLight:
		return "light"
	case EntityTypeSwitch:
		return "switch"
	case EntityTypeFan:
		return "fan"
	case EntityTypeClimate:
		return "climate"
	case EntityTypeCover:
		return "cover"
	case EntityTypeSensor:
		return "sensor"
	default:
		return "unknown"
	}
}

func (et EntityType) Domain() string {
	return et.String()
}

func NewResolver(cfg *config.Config, client *client.HomeAssistantClient) *Resolver {
	return &Resolver{
		config: cfg,
		client: client,
	}
}

func (r *Resolver) ResolveEntity(ctx context.Context, area, entityType, entityName string) (*EntityMatch, error) {
	states, err := r.client.GetStates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entities: %w", err)
	}

	matches := r.findMatches(states, area, entityType, entityName)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no entities found matching criteria")
	}

	if len(matches) == 1 {
		return &matches[0], nil
	}

	sort.Slice(matches, func(i, j int) bool {
		// First sort by score (higher scores first)
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		// For ties, prefer entities that are not "unavailable"
		// This requires getting state info, but we have it from the search
		return false // Keep original order for same scores
	})

	if matches[0].Score > matches[1].Score+0.05 {
		return &matches[0], nil
	}

	// If scores are very close, pick the first match (sorted by score)
	// This handles cases where multiple entities have the same score
	return &matches[0], nil
}

func (r *Resolver) findMatches(states []client.EntityState, area, entityType, entityName string) []EntityMatch {
	var matches []EntityMatch

	for _, state := range states {
		match := r.scoreEntity(state, area, entityType, entityName)
		if match.Score > r.config.Preferences.FuzzyThreshold {
			matches = append(matches, match)
		}
	}

	return matches
}

func (r *Resolver) DebugFindMatches(states []client.EntityState, area, entityType, entityName string) []EntityMatch {
	var matches []EntityMatch

	for _, state := range states {
		match := r.scoreEntity(state, area, entityType, entityName)
		matches = append(matches, match)
	}

	return matches
}

func (r *Resolver) scoreEntity(state client.EntityState, area, entityType, entityName string) EntityMatch {
	match := EntityMatch{
		EntityID: state.EntityID,
		Domain:   strings.Split(state.EntityID, ".")[0],
	}

	if friendlyName, ok := state.Attributes["friendly_name"].(string); ok {
		match.FriendlyName = friendlyName
	} else {
		match.FriendlyName = state.EntityID
	}

	if areaName, ok := state.Attributes["area_id"].(string); ok {
		match.Area = areaName
	}

	var score float64

	if entityType != "" {
		domainScore := r.scoreDomain(match.Domain, entityType)
		if domainScore == 0 {
			return match
		}
		score += domainScore * 0.4
	}

	if area != "" {
		areaScore := r.scoreArea(match, area)
		score += areaScore * 0.3
	}

	if entityName != "" {
		nameScore := r.scoreName(match.FriendlyName, entityName)
		score += nameScore * 0.3
	}

	// Penalize unavailable entities slightly
	if state.State == "unavailable" {
		score *= 0.95
	}

	// Slight bonus for more descriptive names (more words = more specific)
	wordCount := len(strings.Fields(match.FriendlyName))
	if wordCount > 2 {
		score += 0.01 * float64(wordCount-2) // Small bonus for specificity
	}

	match.Score = score
	return match
}

func (r *Resolver) scoreDomain(domain, entityType string) float64 {
	normalizedType := strings.ToLower(entityType)
	
	domainMappings := map[string][]string{
		"light":   {"light", "lights", "lamp", "lamps", "bulb", "bulbs"},
		"switch":  {"switch", "switches", "outlet", "outlets", "plug", "plugs"},
		"fan":     {"fan", "fans", "ceiling", "exhaust"},
		"climate": {"climate", "thermostat", "ac", "heat", "temp", "temperature", "hvac"},
		"cover":   {"cover", "covers", "blind", "blinds", "curtain", "curtains", "shade", "shades", "garage", "door", "doors"},
		"sensor":  {"sensor", "sensors", "temperature", "humidity", "motion", "occupancy"},
	}

	for targetDomain, keywords := range domainMappings {
		for _, keyword := range keywords {
			if strings.Contains(normalizedType, keyword) {
				if domain == targetDomain {
					return 1.0
				} else {
					return 0.0
				}
			}
		}
	}

	return 0.0
}

func (r *Resolver) scoreArea(match EntityMatch, area string) float64 {
	normalizedArea := strings.ToLower(area)
	
	if alias, exists := r.config.Aliases[normalizedArea]; exists {
		normalizedArea = strings.ToLower(alias)
	}

	// Check area_id first (if available)
	if match.Area != "" {
		areaScore := r.fuzzyMatch(strings.ToLower(match.Area), normalizedArea)
		if areaScore > 0.8 {
			return areaScore
		}
	}

	// Check if area appears in friendly name
	normalizedFriendly := strings.ToLower(match.FriendlyName)
	
	// Direct word match in friendly name gets higher score
	friendlyWords := strings.Fields(normalizedFriendly)
	for _, word := range friendlyWords {
		if word == normalizedArea {
			return 1.0 // Perfect score for exact word match
		}
		if strings.Contains(word, normalizedArea) {
			return 0.95 // High score for word containing area
		}
		if strings.Contains(normalizedArea, word) {
			return 0.90 // Slightly lower for area containing word
		}
	}
	
	// Additional scoring for exact substring matches in name
	if strings.Contains(normalizedFriendly, " "+normalizedArea+" ") {
		return 0.98 // High score for area as separate word with spaces
	}
	if strings.Contains(normalizedFriendly, normalizedArea+" ") || strings.Contains(normalizedFriendly, " "+normalizedArea) {
		return 0.95 // Good score for area at start/end with space
	}
	
	// Fallback to fuzzy matching
	nameScore := r.fuzzyMatch(normalizedFriendly, normalizedArea)
	if nameScore > 0.6 {
		return nameScore
	}

	return 0.0
}

func (r *Resolver) scoreName(friendlyName, entityName string) float64 {
	normalizedFriendly := strings.ToLower(friendlyName)
	normalizedEntity := strings.ToLower(entityName)

	if normalizedFriendly == normalizedEntity {
		return 1.0
	}

	if strings.Contains(normalizedFriendly, normalizedEntity) {
		return 0.9
	}

	return r.fuzzyMatch(normalizedFriendly, normalizedEntity)
}

func (r *Resolver) fuzzyMatch(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}

	s1 = r.normalizeString(s1)
	s2 = r.normalizeString(s2)

	if s1 == s2 {
		return 1.0
	}

	if strings.Contains(s1, s2) {
		return 0.9
	}
	if strings.Contains(s2, s1) {
		return 0.9
	}

	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)
	
	for _, word1 := range words1 {
		for _, word2 := range words2 {
			if word1 == word2 {
				return 0.8
			}
			if strings.Contains(word1, word2) || strings.Contains(word2, word1) {
				return 0.7
			}
		}
	}

	return r.levenshteinSimilarity(s1, s2)
}

func (r *Resolver) normalizeString(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			result.WriteRune(' ')
		}
	}
	return strings.TrimSpace(result.String())
}

func (r *Resolver) levenshteinSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	distance := r.levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

func (r *Resolver) levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b && a < c {
		return a
	}
	if b < c {
		return b
	}
	return c
}

func ParseEntityType(s string) EntityType {
	switch strings.ToLower(s) {
	case "light", "lights", "lamp", "lamps", "bulb", "bulbs":
		return EntityTypeLight
	case "switch", "switches", "outlet", "outlets", "plug", "plugs":
		return EntityTypeSwitch
	case "fan", "fans":
		return EntityTypeFan
	case "climate", "thermostat", "ac", "heat", "hvac":
		return EntityTypeClimate
	case "cover", "covers", "blind", "blinds", "curtain", "curtains", "shade", "shades", "garage", "door", "doors":
		return EntityTypeCover
	case "sensor", "sensors":
		return EntityTypeSensor
	default:
		return EntityTypeUnknown
	}
}

func ParseAction(s string) string {
	switch strings.ToLower(s) {
	case "on", "turn_on", "enable", "open":
		return "turn_on"
	case "off", "turn_off", "disable", "close":
		return "turn_off"
	case "toggle", "switch":
		return "toggle"
	default:
		return s
	}
}

func ParseNumericValue(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}