package solver

// SolverJobResponse is the answer of a solver's create operation (Start) and of
// the Jobs API's create operation: the new job's id plus, when the solve ran
// synchronously, its inline result or error. It mirrors the API's
// SolverJobResponse.
type SolverJobResponse struct {
	// JobID is the public PlanSolve job identifier - use it to poll status and fetch the solution.
	JobID string `json:"jobId"`
	// SolverJobID is the underlying solver engine's job identifier, when exposed.
	SolverJobID *string `json:"solverJobId,omitempty"`
	// Result is the inline solver result as raw JSON, when available synchronously.
	Result *string `json:"result,omitempty"`
	// Error is the error message when the solve request could not be accepted or run.
	Error *string `json:"error,omitempty"`
}
