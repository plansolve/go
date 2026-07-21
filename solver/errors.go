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
//   - ErrorResponse (402/422):        {"error": "<string>"}
//   - ProblemDetails (401/403/502):   {"detail": "...", "title": "..."}
//
// Priority mirrors the other SDKs: errors(validation) > error > detail > title
// > raw body. If the body is not a JSON object, the raw body string is returned.
func ExtractErrorMessage(body []byte) string {
	raw := string(body)

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(body, &obj); err != nil {
		return raw
	}

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
				return strings.Join(messages, "; ")
			}
		}
	}

	// ErrorResponse.
	if s := stringField(obj, "error"); s != "" {
		return s
	}

	// ProblemDetails: detail, then title.
	if s := stringField(obj, "detail"); s != "" {
		return s
	}
	if s := stringField(obj, "title"); s != "" {
		return s
	}

	return raw
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
