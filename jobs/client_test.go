package jobs

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateJobPostsRequestAndReturnsSolverJobResponse(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &gotBody)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jobId": "job-c", "solverJobId": "solver-1", "result": null, "error": null}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	created, err := client.CreateJob(context.Background(), CreateJobRequest{SolverID: "field_service", Request: `{"visits":[]}`})
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected POST, got %s", gotMethod)
	}
	if gotPath != "/api/v1/jobs" {
		t.Errorf("expected path '/api/v1/jobs', got '%s'", gotPath)
	}
	if _, ok := gotBody["subscriptionId"]; ok {
		t.Errorf("request must not carry subscriptionId, got %v", gotBody)
	}
	if gotBody["solverId"] != "field_service" || gotBody["request"] != `{"visits":[]}` {
		t.Errorf("unexpected request body %v", gotBody)
	}
	if created.JobID != "job-c" {
		t.Errorf("expected jobId 'job-c', got '%s'", created.JobID)
	}
	if created.SolverJobID == nil || *created.SolverJobID != "solver-1" {
		t.Errorf("expected solverJobId 'solver-1', got %v", created.SolverJobID)
	}
	if created.Result != nil || created.Error != nil {
		t.Errorf("expected nil result and error, got %v / %v", created.Result, created.Error)
	}
}

func TestGetJobsKeepsSubscriptionFilter(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items": [], "totalCount": 0, "page": 1, "pageSize": 10, "totalPages": 0}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")
	subscriptionID := 3
	if _, err := client.GetJobs(context.Background(), GetJobsParams{SubscriptionID: &subscriptionID}); err != nil {
		t.Fatalf("GetJobs failed: %v", err)
	}
	if gotQuery != "subscriptionId=3" {
		t.Errorf("expected query 'subscriptionId=3', got '%s'", gotQuery)
	}
}
