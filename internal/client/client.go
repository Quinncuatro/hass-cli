package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/quinncuatro/hass-cli/internal/config"
)

type HomeAssistantClient struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
	token      string
}

type EntityState struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged time.Time              `json:"last_changed"`
	LastUpdated time.Time              `json:"last_updated"`
}

type ServiceCallRequest struct {
	Domain      string                 `json:"domain"`
	Service     string                 `json:"service"`
	Target      map[string]interface{} `json:"target,omitempty"`
	ServiceData map[string]interface{} `json:"service_data,omitempty"`
}

type ServiceCallResponse struct {
	Context struct {
		ID       string `json:"id"`
		ParentID string `json:"parent_id"`
		UserID   string `json:"user_id"`
	} `json:"context"`
}

type UnitSystem struct {
	Length      string `json:"length"`
	Mass        string `json:"mass"`
	Temperature string `json:"temperature"`
	Volume      string `json:"volume"`
}

type SystemStatus struct {
	Message            string     `json:"message"`
	ConfigDir          string     `json:"config_dir"`
	Version            string     `json:"version"`
	Timezone           string     `json:"timezone"`
	SafeMode           bool       `json:"safe_mode"`
	State              string     `json:"state"`
	ExternalURL        string     `json:"external_url"`
	InternalURL        string     `json:"internal_url"`
	LocationName       string     `json:"location_name"`
	UnitSystem         UnitSystem `json:"unit_system"`
	ConfigSource       string     `json:"config_source"`
	RecoveryMode       bool       `json:"recovery_mode"`
	SupportsStatistics bool       `json:"supports_statistics"`
}

func New(cfg *config.Config) *HomeAssistantClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: cfg.HomeAssistant.SkipTLSVerify,
		},
	}

	httpClient := &http.Client{
		Timeout:   cfg.HomeAssistant.Timeout,
		Transport: transport,
	}

	return &HomeAssistantClient{
		config:     cfg,
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(cfg.HomeAssistant.URL, "/"),
		token:      cfg.HomeAssistant.Token,
	}
}

func (c *HomeAssistantClient) TestConnection(ctx context.Context) error {
	if c.baseURL == "" || c.token == "" {
		return fmt.Errorf("home Assistant URL and token must be configured")
	}

	_, err := c.GetSystemStatus(ctx)
	return err
}

func (c *HomeAssistantClient) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	var status SystemStatus
	err := c.makeRequest(ctx, "GET", "/api/config", nil, &status)
	if err != nil {
		return nil, fmt.Errorf("failed to get system status: %w", err)
	}
	return &status, nil
}

func (c *HomeAssistantClient) GetStates(ctx context.Context) ([]EntityState, error) {
	var states []EntityState
	err := c.makeRequest(ctx, "GET", "/api/states", nil, &states)
	if err != nil {
		return nil, fmt.Errorf("failed to get states: %w", err)
	}
	return states, nil
}

func (c *HomeAssistantClient) GetState(ctx context.Context, entityID string) (*EntityState, error) {
	var state EntityState
	path := fmt.Sprintf("/api/states/%s", entityID)
	err := c.makeRequest(ctx, "GET", path, nil, &state)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for %s: %w", entityID, err)
	}
	return &state, nil
}

func (c *HomeAssistantClient) CallService(ctx context.Context, domain, service string, target map[string]interface{}, serviceData map[string]interface{}) (*ServiceCallResponse, error) {
	request := ServiceCallRequest{
		Domain:      domain,
		Service:     service,
		Target:      target,
		ServiceData: serviceData,
	}

	var response ServiceCallResponse
	path := fmt.Sprintf("/api/services/%s/%s", domain, service)
	err := c.makeRequest(ctx, "POST", path, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to call service %s.%s: %w", domain, service, err)
	}
	return &response, nil
}

func (c *HomeAssistantClient) TurnOnEntity(ctx context.Context, entityID string) error {
	domain := strings.Split(entityID, ".")[0]
	
	// Create direct service request payload
	payload := map[string]interface{}{
		"entity_id": entityID,
	}
	
	path := fmt.Sprintf("/api/services/%s/turn_on", domain)
	err := c.makeRequest(ctx, "POST", path, payload, nil)
	return err
}

func (c *HomeAssistantClient) TurnOffEntity(ctx context.Context, entityID string) error {
	domain := strings.Split(entityID, ".")[0]
	
	payload := map[string]interface{}{
		"entity_id": entityID,
	}
	
	path := fmt.Sprintf("/api/services/%s/turn_off", domain)
	err := c.makeRequest(ctx, "POST", path, payload, nil)
	return err
}

func (c *HomeAssistantClient) ToggleEntity(ctx context.Context, entityID string) error {
	domain := strings.Split(entityID, ".")[0]
	
	payload := map[string]interface{}{
		"entity_id": entityID,
	}
	
	path := fmt.Sprintf("/api/services/%s/toggle", domain)
	err := c.makeRequest(ctx, "POST", path, payload, nil)
	return err
}

func (c *HomeAssistantClient) makeRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}