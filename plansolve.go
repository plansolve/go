package plansolve

import (
	"net/http"
	"os"
	"time"

	"github.com/plansolve/go/fieldservice"
	"github.com/plansolve/go/professionalservices"
	"github.com/plansolve/go/shift"
)

const (
	defaultBaseURL = "https://plansolve.app"
	apiKeyHeader   = "X-API-KEY"
)

// Client is the main PlanSolve API client.
type Client struct {
	FieldService         *fieldservice.Client
	ProfessionalServices *professionalservices.Client
	Shift                *shift.Client
}

// NewClient creates a new PlanSolve client with the given API key.
func NewClient(apiKey string) *Client {
	baseURL := os.Getenv("PLANSOLVE_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	return &Client{
		FieldService:         fieldservice.NewClient(httpClient, baseURL, apiKey),
		ProfessionalServices: professionalservices.NewClient(httpClient, baseURL, apiKey),
		Shift:                shift.NewClient(httpClient, baseURL, apiKey),
	}
}

// NewClientFromEnv creates a new PlanSolve client using the PLANSOLVE_API_KEY
// environment variable.
func NewClientFromEnv() (*Client, error) {
	apiKey := os.Getenv("PLANSOLVE_API_KEY")
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}
	return NewClient(apiKey), nil
}
