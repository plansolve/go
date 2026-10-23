package solver

import (
	"encoding/json"
	"testing"
)

func TestSolverStatusUnmarshalPublicValues(t *testing.T) {
	cases := map[string]SolverStatus{
		"queued":    SolverStatusQueued,
		"solving":   SolverStatusSolving,
		"stopping":  SolverStatusStopping,
		"completed": SolverStatusCompleted,
	}
	for wire, want := range cases {
		var r SolverStatusResponse
		if err := json.Unmarshal([]byte(`{"status":"`+wire+`"}`), &r); err != nil {
			t.Fatalf("%s: unexpected error: %v", wire, err)
		}
		if r.SolverStatus != want {
			t.Errorf("%s: expected %q, got %q", wire, want, r.SolverStatus)
		}
	}
}

func TestSolverStatusUnmarshalLegacyValues(t *testing.T) {
	cases := map[string]SolverStatus{
		"SOLVING_SCHEDULED": SolverStatusQueued,
		"SOLVING_ACTIVE":    SolverStatusSolving,
		"TERMINATING_EARLY": SolverStatusStopping,
		"NOT_SOLVING":       SolverStatusCompleted,
	}
	for wire, want := range cases {
		var r SolverStatusResponse
		if err := json.Unmarshal([]byte(`{"solverStatus":"`+wire+`"}`), &r); err != nil {
			t.Fatalf("%s: unexpected error: %v", wire, err)
		}
		if r.SolverStatus != want {
			t.Errorf("%s: expected %q, got %q", wire, want, r.SolverStatus)
		}
	}
}

func TestSolverStatusResponsePrefersStatusOverLegacyField(t *testing.T) {
	cases := []struct {
		name string
		json string
		want SolverStatus
	}{
		{"both agree", `{"status":"completed","solverStatus":"NOT_SOLVING"}`, SolverStatusCompleted},
		{"both disagree, status wins", `{"status":"solving","solverStatus":"NOT_SOLVING"}`, SolverStatusSolving},
		{"only status", `{"status":"queued"}`, SolverStatusQueued},
		{"only legacy solverStatus", `{"solverStatus":"TERMINATING_EARLY"}`, SolverStatusStopping},
		{"null status falls back", `{"status":null,"solverStatus":"NOT_SOLVING"}`, SolverStatusCompleted},
		{"neither", `{}`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var r SolverStatusResponse
			if err := json.Unmarshal([]byte(tc.json), &r); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if r.SolverStatus != tc.want {
				t.Errorf("expected %q, got %q", tc.want, r.SolverStatus)
			}
		})
	}
}

func TestSolverStatusUnmarshalKeepsUnknownValue(t *testing.T) {
	var r SolverStatusResponse
	if err := json.Unmarshal([]byte(`{"status":"something_new"}`), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.SolverStatus != "something_new" {
		t.Errorf("expected unknown value kept as-is, got %q", r.SolverStatus)
	}
}

func TestSolverStatusMarshalEmitsPublicValue(t *testing.T) {
	data, err := json.Marshal(SolverStatusCompleted)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `"completed"` {
		t.Errorf("expected \"completed\", got %s", data)
	}
}

func TestSolverStatusResponseUnmarshalStringScore(t *testing.T) {
	var r SolverStatusResponse
	raw := `{"jobId":"j","status":"completed","solverStatus":"NOT_SOLVING","score":"0hard/0medium/-5soft","feasible":true,"solving":false}`
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Score == nil || r.Score.Soft != -5 {
		t.Errorf("expected soft -5, got %+v", r.Score)
	}
}

func TestDeprecatedSolverStatusNamesEqualPublicValues(t *testing.T) {
	cases := map[SolverStatus]SolverStatus{
		SolverStatusNotSolving:       SolverStatusCompleted,
		SolverStatusSolvingActive:    SolverStatusSolving,
		SolverStatusTerminatingEarly: SolverStatusStopping,
		SolverStatusSolvingScheduled: SolverStatusQueued,
	}
	for old, want := range cases {
		if old != want {
			t.Errorf("expected %q, got %q", want, old)
		}
	}
}

func TestParseSolverStatus(t *testing.T) {
	if got, ok := ParseSolverStatus("stopping"); !ok || got != SolverStatusStopping {
		t.Errorf("public value: expected stopping, got %q (ok=%t)", got, ok)
	}
	if got, ok := ParseSolverStatus("SOLVING_SCHEDULED"); !ok || got != SolverStatusQueued {
		t.Errorf("legacy value: expected queued, got %q (ok=%t)", got, ok)
	}
	if got, ok := ParseSolverStatus("something_new"); ok || got != "" {
		t.Errorf("unknown value: expected zero value, got %q (ok=%t)", got, ok)
	}
}
