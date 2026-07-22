package plansolve

import "errors"

// ErrMissingAPIKey is returned when the PLANSOLVE_API_KEY environment variable is not set.
var ErrMissingAPIKey = errors.New("PLANSOLVE_API_KEY environment variable is not set")

// ErrorResponse is the API's error body shape: {"error": "..."}. It mirrors the
// contract's ErrorResponse model. The clients surface error bodies as strings
// via the solver package's ExtractErrorMessage; this type is provided for
// callers that want to decode the raw body themselves.
type ErrorResponse struct {
	Error *string `json:"error,omitempty"`
}
