package jobs

import (
	"encoding/json"
	"testing"
)

func TestCanDeserializeCancelledJob(t *testing.T) {
	jsonStr := `{"id": 7, "jobId": "job-b", "subscriptionId": 1, "solverType": 1, "status": 4,
		"createdAt": "2026-10-13T12:00:00Z", "hasResponse": false, "billedItems": 0, "estimatedItems": 5}`

	var job JobDto
	if err := json.Unmarshal([]byte(jsonStr), &job); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if job.Status != JobStatusCancelled {
		t.Errorf("expected JobStatusCancelled, got %d", job.Status)
	}
	if job.SolverType != SolverTypeShift {
		t.Errorf("expected SolverTypeShift, got %d", job.SolverType)
	}
	if job.EstimatedItems == nil || *job.EstimatedItems != 5 {
		t.Errorf("expected estimatedItems 5, got %v", job.EstimatedItems)
	}
}

func TestUnknownJobStatusIsKept(t *testing.T) {
	var job JobDto
	if err := json.Unmarshal([]byte(`{"status": 9}`), &job); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if job.Status != 9 {
		t.Errorf("expected status 9, got %d", job.Status)
	}
}
