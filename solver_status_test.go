package plansolve

import "testing"

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
