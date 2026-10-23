// Package solver holds types shared between the root plansolve package and the
// per-service subpackages (fieldservice, professionalservices, shift). It lives
// in its own leaf package so the subpackages can use these types without
// importing the root package - which would create an import cycle, since the
// root package imports the subpackages to assemble the aggregate Client.
package solver

import "encoding/json"

// SolverStatus represents the state of the solver.
type SolverStatus string

const (
	// SolverStatusQueued indicates the job is queued and has not started solving yet.
	SolverStatusQueued SolverStatus = "queued"
	// SolverStatusSolving indicates the solver is actively working on the job.
	SolverStatusSolving SolverStatus = "solving"
	// SolverStatusStopping indicates the solver is stopping before its limits are reached.
	SolverStatusStopping SolverStatus = "stopping"
	// SolverStatusCompleted indicates the solver has finished working on the job.
	SolverStatusCompleted SolverStatus = "completed"
)

// The 0.34 names, kept so existing code compiles. They equal the public values.
const (
	// SolverStatusNotSolving is the 0.34 name for SolverStatusCompleted.
	//
	// Deprecated: use SolverStatusCompleted.
	SolverStatusNotSolving = SolverStatusCompleted
	// SolverStatusSolvingActive is the 0.34 name for SolverStatusSolving.
	//
	// Deprecated: use SolverStatusSolving.
	SolverStatusSolvingActive = SolverStatusSolving
	// SolverStatusTerminatingEarly is the 0.34 name for SolverStatusStopping.
	//
	// Deprecated: use SolverStatusStopping.
	SolverStatusTerminatingEarly = SolverStatusStopping
	// SolverStatusSolvingScheduled is the 0.34 name for SolverStatusQueued.
	//
	// Deprecated: use SolverStatusQueued.
	SolverStatusSolvingScheduled = SolverStatusQueued
)

// legacySolverStatuses maps the raw engine values older servers sent to the public values.
// Transitional: delete once servers no longer send the legacy solverStatus field (after 2026-10-30).
var legacySolverStatuses = map[string]SolverStatus{
	"SOLVING_SCHEDULED": SolverStatusQueued,
	"SOLVING_ACTIVE":    SolverStatusSolving,
	"TERMINATING_EARLY": SolverStatusStopping,
	"NOT_SOLVING":       SolverStatusCompleted,
}

// ParseSolverStatus maps a wire value to a SolverStatus, accepting both the
// public values and the legacy engine values. ok is false for any other value,
// in which case the returned status is the zero value.
func ParseSolverStatus(raw string) (SolverStatus, bool) {
	switch status := SolverStatus(raw); status {
	case SolverStatusQueued, SolverStatusSolving, SolverStatusStopping, SolverStatusCompleted:
		return status, true
	}
	status, ok := legacySolverStatuses[raw]
	return status, ok
}

// ResolveSolverStatus picks the wire value from the public status field when it
// is present, otherwise from the legacy solverStatus field, and parses it with
// ParseSolverStatus. An unrecognized value is returned as-is with ok false; when
// neither field is present the result is the zero value with ok false.
func ResolveSolverStatus(status, legacy *string) (SolverStatus, bool) {
	raw := status
	if raw == nil {
		raw = legacy
	}
	if raw == nil {
		return "", false
	}
	if parsed, ok := ParseSolverStatus(*raw); ok {
		return parsed, true
	}
	return SolverStatus(*raw), false
}

// UnmarshalJSON decodes the status, mapping legacy engine values to their
// public equivalents. Any other value is kept as-is.
func (s *SolverStatus) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if status, ok := ParseSolverStatus(raw); ok {
		*s = status
		return nil
	}
	*s = SolverStatus(raw)
	return nil
}

// SolverStatusResponse represents the status of a solver job.
type SolverStatusResponse struct {
	// JobID is the public job identifier the status refers to.
	JobID string `json:"jobId"`
	// SolverStatus is the public solver status (queued, solving, stopping or completed),
	// read from the wire's status field, or from the legacy solverStatus field when
	// an older server sends only that.
	SolverStatus SolverStatus `json:"solverStatus"`
	// Score is the current best score, when available.
	Score *Score `json:"score,omitempty"`
	// Feasible reports whether the current best solution satisfies all hard constraints.
	Feasible bool `json:"feasible"`
	// Solving reports whether the solver is still actively working on the job.
	Solving bool `json:"solving"`
}

// UnmarshalJSON decodes the response, routing the score field through
// ParseScoreJSON so both wire shapes - and the absence of a score - are
// handled in one place, and the status through ResolveSolverStatus so the
// public status field wins over the legacy solverStatus field.
func (r *SolverStatusResponse) UnmarshalJSON(data []byte) error {
	type alias SolverStatusResponse
	aux := &struct {
		Score        json.RawMessage `json:"score"`
		Status       *string         `json:"status"`
		SolverStatus *string         `json:"solverStatus"`
		*alias
	}{alias: (*alias)(r)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	score, err := ParseScoreJSON(aux.Score)
	if err != nil {
		return err
	}
	r.Score = score
	r.SolverStatus, _ = ResolveSolverStatus(aux.Status, aux.SolverStatus)
	return nil
}
