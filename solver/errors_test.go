package solver

import (
	"strings"
	"testing"
)

func TestExtractErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "validation problem details",
			body: `{"type":"t","title":"One or more validation errors occurred.","status":400,"errors":{"Name":["Name is required"],"Age":["Age must be positive"]}}`,
			// messages sorted ascending, joined with "; "
			want: "Age: Age must be positive; Name: Name is required",
		},
		{
			name: "error response",
			body: `{"error":"insufficient credits"}`,
			want: "insufficient credits",
		},
		{
			name: "problem details detail",
			body: `{"type":"t","title":"Unauthorized","status":401,"detail":"API key is invalid"}`,
			want: "API key is invalid",
		},
		{
			name: "problem details title only",
			body: `{"type":"t","title":"Bad Gateway","status":502}`,
			want: "Bad Gateway",
		},
		{
			name: "non-json fallback",
			body: `<html>500 Internal Server Error</html>`,
			want: `<html>500 Internal Server Error</html>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractErrorMessage([]byte(tt.body))
			if got != tt.want {
				t.Errorf("ExtractErrorMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestExtractErrorMessageValidationFieldForm asserts both field messages appear
// in the "Field: msg" form regardless of ordering concerns.
func TestExtractErrorMessageValidationFieldForm(t *testing.T) {
	body := `{"errors":{"Name":["Name is required","Name too long"],"Age":["Age must be positive"]}}`
	got := ExtractErrorMessage([]byte(body))

	for _, want := range []string{"Name: Name is required", "Name: Name too long", "Age: Age must be positive"} {
		if !strings.Contains(got, want) {
			t.Errorf("ExtractErrorMessage() = %q, missing %q", got, want)
		}
	}
}
