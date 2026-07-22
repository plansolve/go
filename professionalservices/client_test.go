package professionalservices

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResultStampsJobID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Response body deliberately carries no jobId; the client must stamp it.
		_, _ = w.Write([]byte(`{"id": "req1", "feasible": true, "scoreString": "0hard/-5soft", "employees": [], "tasks": [], "assignedTasks": [], "unassignedTasks": []}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	result, err := client.GetResult(context.Background(), "job-456")
	if err != nil {
		t.Fatalf("GetResult failed: %v", err)
	}

	if result.JobID == nil {
		t.Fatalf("expected JobID to be stamped, got nil")
	}
	if *result.JobID != "job-456" {
		t.Errorf("expected JobID 'job-456', got '%s'", *result.JobID)
	}
	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible true")
	}
}
