// Package solver holds types shared between the root plansolve package and the
// per-service subpackages (fieldservice, professionalservices, shift). It lives
// in its own leaf package so the subpackages can use these types without
// importing the root package — which would create an import cycle, since the
// root package imports the subpackages to assemble the aggregate Client.
package solver

// SolverStatus represents the state of the solver.
type SolverStatus string

const (
	// SolverStatusNotSolving indicates the solver is not currently working on the job.
	SolverStatusNotSolving SolverStatus = "NOT_SOLVING"
	// SolverStatusSolvingActive indicates the solver is actively working on the job.
	SolverStatusSolvingActive SolverStatus = "SOLVING_ACTIVE"
	// SolverStatusTerminatingEarly indicates the solver is stopping before its limits are reached.
	SolverStatusTerminatingEarly SolverStatus = "TERMINATING_EARLY"
	// SolverStatusSolvingScheduled indicates the job is queued and scheduled to start solving.
	SolverStatusSolvingScheduled SolverStatus = "SOLVING_SCHEDULED"
)

// SolverStatusResponse represents the status of a solver job.
type SolverStatusResponse struct {
	// JobID is the public job identifier the status refers to.
	JobID string `json:"jobId"`
	// SolverStatus is the raw solver engine status (e.g. SOLVING_ACTIVE, NOT_SOLVING).
	SolverStatus SolverStatus `json:"solverStatus"`
	// Score is the current best score as a solver score string, when available.
	Score string `json:"score"`
	// Feasible reports whether the current best solution satisfies all hard constraints.
	Feasible bool `json:"feasible"`
	// Solving reports whether the solver is still actively working on the job.
	Solving bool `json:"solving"`
}
