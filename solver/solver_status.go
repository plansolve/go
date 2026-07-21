// Package solver holds types shared between the root plansolve package and the
// per-service subpackages (fieldservice, professionalservices, shift). It lives
// in its own leaf package so the subpackages can use these types without
// importing the root package — which would create an import cycle, since the
// root package imports the subpackages to assemble the aggregate Client.
package solver

// SolverStatus represents the state of the solver.
type SolverStatus string

const (
	SolverStatusNotSolving       SolverStatus = "NOT_SOLVING"
	SolverStatusSolvingActive    SolverStatus = "SOLVING_ACTIVE"
	SolverStatusTerminatingEarly SolverStatus = "TERMINATING_EARLY"
	SolverStatusSolvingScheduled SolverStatus = "SOLVING_SCHEDULED"
)

// SolverStatusResponse represents the status of a solver job.
type SolverStatusResponse struct {
	JobID        string       `json:"jobId"`
	SolverStatus SolverStatus `json:"solverStatus"`
	Score        string       `json:"score"`
	Feasible     bool         `json:"feasible"`
	Solving      bool         `json:"solving"`
	Message      string       `json:"message"`
}
