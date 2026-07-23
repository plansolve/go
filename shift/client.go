package shift

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/plansolve/go/solver"
)

const routeShiftSolve = "/api/v1/shift"

// Client is the shift API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new shift API client.
func NewClient(httpClient *http.Client, baseURL, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// Start starts a new shift optimization job.
func (c *Client) Start(ctx context.Context, request ShiftRequest) (*ShiftStartResponse, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.baseURL + routeShiftSolve
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
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

	var result ShiftStartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetStatus gets the status of a shift optimization job.
func (c *Client) GetStatus(ctx context.Context, jobID string) (*solver.SolverStatusResponse, error) {
	url := fmt.Sprintf("%s%s/%s/status", c.baseURL, routeShiftSolve, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var result solver.SolverStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetResult gets the result of a completed shift optimization job.
func (c *Client) GetResult(ctx context.Context, jobID string) (*ShiftResultResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, routeShiftSolve, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var result ShiftResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result.JobID = &jobID
	return &result, nil
}

// Analyze gets analysis for a completed shift optimization job.
func (c *Client) Analyze(ctx context.Context, jobID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s/%s/analyze", c.baseURL, routeShiftSolve, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// Stop stops an in-progress shift solve and returns the best solution found so far.
func (c *Client) Stop(ctx context.Context, jobID string) (*ShiftResultResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.baseURL, routeShiftSolve, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
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

	var result ShiftResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	result.JobID = &jobID
	return &result, nil
}

// StartAndWaitForCompletion starts a job and polls until it completes.
func (c *Client) StartAndWaitForCompletion(ctx context.Context, request ShiftRequest, pollIntervalMs int, maxAttempts int) (*ShiftResultResponse, error) {
	if pollIntervalMs <= 0 {
		pollIntervalMs = 5000
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	startResp, err := c.Start(ctx, request)
	if err != nil {
		return nil, err
	}

	return c.WaitForCompletion(ctx, startResp.JobID, pollIntervalMs, maxAttempts)
}

// WaitForCompletion polls for job status until the solver finishes or max attempts is reached.
func (c *Client) WaitForCompletion(ctx context.Context, jobID string, pollIntervalMs int, maxAttempts int) (*ShiftResultResponse, error) {
	if jobID == "" {
		return nil, fmt.Errorf("jobId was not provided")
	}
	if pollIntervalMs <= 0 {
		pollIntervalMs = 5000
	}
	if maxAttempts <= 0 {
		maxAttempts = 60
	}

	pollInterval := time.Duration(pollIntervalMs) * time.Millisecond
	attempts := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(pollInterval):
		}

		status, err := c.GetStatus(ctx, jobID)
		if err != nil {
			return nil, err
		}
		attempts++

		if !isStillSolving(status) || attempts >= maxAttempts {
			if isStillSolving(status) {
				return nil, fmt.Errorf("solver still running after the client's poll budget elapsed (maxAttempts*pollInterval); increase maxAttempts/pollInterval, or set options.spentLimit in the request so the solver stops on its own")
			}
			return c.GetResult(ctx, jobID)
		}
	}
}

func isStillSolving(status *solver.SolverStatusResponse) bool {
	return status.Solving ||
		status.SolverStatus == solver.SolverStatusSolvingScheduled ||
		status.SolverStatus != solver.SolverStatusNotSolving ||
		status.Score == ""
}
