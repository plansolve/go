package plansolve

import "errors"

// ErrMissingAPIKey is returned when the PLANSOLVE_API_KEY environment variable is not set.
var ErrMissingAPIKey = errors.New("PLANSOLVE_API_KEY environment variable is not set")
