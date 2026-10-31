// Package go-sdk provides the Helix observability platform SDK for Go
// This SDK allows Go applications to interact with the Helix platform
package helix // import helix.io/sdk/go

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client represents a Helix platform client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new Helix client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// Event represents an observability event
type Event struct {
	ID        string                 `json:"id"`
	Timestamp int64                  `json:"timestamp"`
	Source    string                 `json:"source"`
	Metrics   map[string]float64     `json:"metrics"`
	Tags      map[string]string      `json:"tags"`
	Tenant    string                 `json:"tenant"`
}

// QueryOptions represents options for a query
type QueryOptions struct {
	Expression string
	Tenant     string
	TimeRange  []int64
}

// QueryResult represents the result of a query
type QueryResult struct {
	Metric     map[string]string `json:"metric"`
	Value      float64           `json:"value"`
	Timestamp  int64             `json:"timestamp"`
}

// IngestEvent inserts an event into the Helix platform
func (c *Client) IngestEvent(ctx context.Context, event *Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/events", 
		strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	return nil
}

// Query executes a query against the Helix platform
func (c *Client) Query(ctx context.Context, options *QueryOptions) ([]QueryResult, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query options: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/query",
		strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Results []QueryResult `json:"results"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Results, nil
}

// TenantInfo represents tenant information
type TenantInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	RateLimit int    `json:"rate_limit"`
	Quota     int64  `json:"quota"`
}

// ListTenants lists all tenants
func (c *Client) ListTenants(ctx context.Context) ([]TenantInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/v1/tenants", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Tenant []TenantInfo `json:"tenants"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Tenant, nil
}

// HealthCheck checks the health of the Helix platform
func (c *Client) HealthCheck(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	return true, nil
}

// Metrics returns platform metrics
func (c *Client) Metrics(ctx context.Context) (map[string]int64, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/metrics", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	var result struct {
		EventsReceived   int64 `json:"events_received"`
		EventsProcessed  int64 `json:"events_processed"`
		ActiveTenants    int   `json:"active_tenants"`
		MemoryUsageBytes int64 `json:"memory_usage_bytes"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return map[string]int64{
		"events_received":   result.EventsReceived,
		"events_processed":  result.EventsProcessed,
		"active_tenants":    int(result.ActiveTenants),
		"memory_usage_bytes": result.MemoryUsageBytes,
	}, nil
}

// SearchEvents searches for events with optional filters
func (c *Client) SearchEvents(ctx context.Context, query string, filters map[string]string) ([]Event, error) {
	// Build URL with query parameters
	url := c.baseURL + "/v1/events/search?query=" + query
	
	for key, value := range filters {
		url += "&" + key + "=" + value
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	
	var result struct {
		Events []Event `json:"events"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return result.Events, nil
}