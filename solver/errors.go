package solver

import (
	"encoding/json"
	"sort"
	"strings"
)

// ExtractErrorMessage turns an HTTP error response body into a clean, human
// readable message. It understands the three error body shapes the PlanSolve
// API documents:
//
//   - ValidationProblemDetails (400): {"errors": {"<field>": ["<msg>", ...]}}
//   - ErrorResponse (402/422):        {"error": "<string>", "traceId": "<string>"}
//   - ProblemDetails (401/403/502):   {"detail": "...", "title": "..."}
//
// Priority mirrors the other SDKs: errors(validation) > error > detail > title
// > raw body. If the body is not a JSON object, the raw body string is returned.
//
// When the body carries a non-empty traceId, it is appended as
// " (traceId: <id>)" - except on the raw body fallback, where the id is already
// part of the text returned.
func ExtractErrorMessage(body []byte) string {
	raw := string(body)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return raw
	}

	if message, ok := extractMessage(obj); ok {
		return withTraceID(message, stringField(obj, "traceId"))
	}

	return raw
}

// extractMessage returns the best message in obj, and whether it held one.
func extractMessage(obj map[string]json.RawMessage) (string, bool) {
	// ValidationProblemDetails: errors is an object of field -> array of strings.
	if errsRaw, ok := obj["errors"]; ok {
		var fields map[string][]string
		if err := json.Unmarshal(errsRaw, &fields); err == nil {
			var messages []string
			for field, msgs := range fields {
				for _, msg := range msgs {
					messages = append(messages, field+": "+msg)
				}
			}
			if len(messages) > 0 {
				sort.Strings(messages)
				return strings.Join(messages, "; "), true
			}
		}
	}

	// ErrorResponse, then ProblemDetails: detail, then title.
	for _, key := range []string{"error", "detail", "title"} {
		if s := stringField(obj, key); s != "" {
			return s, true
		}
	}

	return "", false
}

// withTraceID appends traceID to message when it is non-empty.
func withTraceID(message, traceID string) string {
	if traceID == "" {
		return message
	}
	return message + " (traceId: " + traceID + ")"
}

// stringField returns the value of key if it is a non-empty JSON string.
func stringField(obj map[string]json.RawMessage, key string) string {
	rawVal, ok := obj[key]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(rawVal, &s); err != nil {
		return ""
	}
	return s
}
