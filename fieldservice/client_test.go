package fieldservice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/plansolve/go/solver"
)

func TestStopSendsDeleteAndStampsJobID(t *testing.T) {
	var gotMethod, gotPath, gotKey string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotKey = r.Header.Get("X-API-KEY")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"vehicles": [], "visits": [], "score": "0hard/0medium/-100soft", "totalDrivingTimeSeconds": 100}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	result, err := client.Stop(context.Background(), "job-stop")
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", gotMethod)
	}
	if gotPath != "/api/v1/fieldservice/job-stop" {
		t.Errorf("expected path '/api/v1/fieldservice/job-stop', got '%s'", gotPath)
	}
	if gotKey != "test-key" {
		t.Errorf("expected X-API-KEY 'test-key', got '%s'", gotKey)
	}
	if result.JobID == nil || *result.JobID != "job-stop" {
		t.Errorf("expected JobID 'job-stop' to be stamped")
	}
}

func TestStopReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "job not found"}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	_, err := client.Stop(context.Background(), "missing")
	if err == nil {
		t.Fatalf("expected an error")
	}
	if !strings.Contains(err.Error(), "404") || !strings.Contains(err.Error(), "job not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestAnalyzeCallsAnalyzeEndpoint(t *testing.T) {
	var gotMethod, gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"score": "0hard/-5soft", "constraints": []}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	analysis, err := client.Analyze(context.Background(), "job-an")
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Errorf("expected GET, got %s", gotMethod)
	}
	if gotPath != "/api/v1/fieldservice/job-an/analyze" {
		t.Errorf("expected path '/api/v1/fieldservice/job-an/analyze', got '%s'", gotPath)
	}
	if analysis["score"] != "0hard/-5soft" {
		t.Errorf("expected score '0hard/-5soft', got %v", analysis["score"])
	}
}

func TestIsStillSolvingMatchesServerRule(t *testing.T) {
	cases := []struct {
		name   string
		status solver.SolverStatusResponse
		want   bool
	}{
		{"done without score", solver.SolverStatusResponse{Solving: false, SolverStatus: solver.SolverStatusNotSolving}, false},
		{"done with score", solver.SolverStatusResponse{Solving: false, SolverStatus: solver.SolverStatusNotSolving, Score: "0hard/0soft"}, false},
		{"solving flag set", solver.SolverStatusResponse{Solving: true, SolverStatus: solver.SolverStatusNotSolving}, true},
		{"active", solver.SolverStatusResponse{Solving: false, SolverStatus: solver.SolverStatusSolvingActive}, true},
		{"scheduled", solver.SolverStatusResponse{Solving: false, SolverStatus: solver.SolverStatusSolvingScheduled}, true},
		{"terminating", solver.SolverStatusResponse{Solving: false, SolverStatus: solver.SolverStatusTerminatingEarly}, true},
	}
	for _, tc := range cases {
		status := tc.status
		if got := isStillSolving(&status); got != tc.want {
			t.Errorf("%s: expected %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestWaitForCompletionFinishesWithoutScore(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/status") {
			_, _ = w.Write([]byte(`{"jobId": "job-w", "solverStatus": "NOT_SOLVING", "solving": false}`))
			return
		}
		_, _ = w.Write([]byte(`{"vehicles": [], "visits": [], "score": "0hard/0medium/-100soft", "totalDrivingTimeSeconds": 100}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	result, err := client.WaitForCompletion(context.Background(), "job-w", 1, 3)
	if err != nil {
		t.Fatalf("WaitForCompletion failed: %v", err)
	}
	if result.JobID == nil || *result.JobID != "job-w" {
		t.Errorf("expected JobID 'job-w' to be stamped")
	}
}

func TestWaitForCompletionTimesOut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jobId": "job-t", "solverStatus": "SOLVING_ACTIVE", "solving": true}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")

	_, err := client.WaitForCompletion(context.Background(), "job-t", 1, 2)
	if err == nil {
		t.Fatalf("expected a timeout error")
	}
	if !strings.Contains(err.Error(), "still running after 2 polls") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestWaitForCompletionRequiresJobID(t *testing.T) {
	client := NewClient(http.DefaultClient, "http://localhost", "test-key")

	_, err := client.WaitForCompletion(context.Background(), "", 1, 1)
	if err == nil || err.Error() != "jobId was not provided" {
		t.Errorf("expected 'jobId was not provided', got %v", err)
	}
}
