package plansolve

import "errors"

// ErrMissingAPIKey is returned when the PLANSOLVE_API_KEY environment variable is not set.
var ErrMissingAPIKey = errors.New("PLANSOLVE_API_KEY environment variable is not set")

// ErrorResponse is the API's error body shape: {"error": "...", "traceId": "..."}.
// It mirrors the contract's ErrorResponse model. The clients surface error bodies
// as strings via the solver package's ExtractErrorMessage; this type is provided
// for callers that want to decode the raw body themselves. The API answers 402,
// 409 and 422 with this shape.
type ErrorResponse struct {
	// Error is a human-readable description of what went wrong.
	Error *string `json:"error,omitempty"`
	// TraceID correlates this failure, and is quotable back to PlanSolve support.
	TraceID *string `json:"traceId,omitempty"`
}

// ApiError is the API's error envelope for 415 Unsupported Media Type:
// {"code": "...", "message": "...", "traceId": "..."}. ExtractErrorMessage
// surfaces its message; this type is provided for callers that want to decode
// the raw body themselves.
type ApiError struct {
	// Code is a machine-readable error code.
	Code *string `json:"code,omitempty"`
	// Message is a human-readable description of what went wrong.
	Message *string `json:"message,omitempty"`
	// TraceID correlates this failure, and is quotable back to PlanSolve support.
	TraceID *string `json:"traceId,omitempty"`
}
