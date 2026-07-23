package jobs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/plansolve/go/solver"
)

const routeJobs = "/api/v1/jobs"

// Client is the public Jobs API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new Jobs API client.
func NewClient(httpClient *http.Client, baseURL, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// GetJobsParams holds the optional filters for listing jobs.
type GetJobsParams struct {
	// Page is the 1-based page index (default 1).
	Page int
	// PageSize is the maximum number of items per page (default 10).
	PageSize int
	// SubscriptionID, when set, only returns jobs billed against this subscription.
	SubscriptionID *int
	// Start, when set, only returns jobs created at or after this timestamp (ISO-8601).
	Start *string
	// End, when set, only returns jobs created at or before this timestamp (ISO-8601).
	End *string
}

// GetJobs lists solver jobs for the authenticated tenant.
func (c *Client) GetJobs(ctx context.Context, params GetJobsParams) (*PagedResultOfJobDto, error) {
	query := url.Values{}
	if params.Page > 0 {
		query.Set("page", strconv.Itoa(params.Page))
	}
	if params.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(params.PageSize))
	}
	if params.SubscriptionID != nil {
		query.Set("subscriptionId", strconv.Itoa(*params.SubscriptionID))
	}
	if params.Start != nil {
		query.Set("start", *params.Start)
	}
	if params.End != nil {
		query.Set("end", *params.End)
	}

	requestURL := c.baseURL + routeJobs
	if encoded := query.Encode(); encoded != "" {
		requestURL += "?" + encoded
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, solver.ExtractErrorMessage(respBody))
	}

	var result PagedResultOfJobDto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// CreateJob creates a new solver job.
func (c *Client) CreateJob(ctx context.Context, request CreateJobRequest) (*JobDto, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestURL := c.baseURL + routeJobs
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, solver.ExtractErrorMessage(respBody))
	}

	var result JobDto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetJob gets a single job by its public job id.
func (c *Client) GetJob(ctx context.Context, jobID string) (*JobDto, error) {
	requestURL := fmt.Sprintf("%s%s/%s", c.baseURL, routeJobs, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if c.apiKey != "" {
		req.Header.Set("X-API-KEY", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d: %s", resp.StatusCode, solver.ExtractErrorMessage(respBody))
	}

	var result JobDto
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
