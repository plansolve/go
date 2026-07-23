package plansolve

import "github.com/plansolve/go/solver"

// The solver status types live in the leaf `solver` package to avoid an import
// cycle (the subpackages need them too, and they can't import this root package
// because this package imports them). Re-exported here as aliases so the public
// API — plansolve.SolverStatusResponse, plansolve.SolverStatusNotSolving, etc. —
// is unchanged.

// SolverStatus represents the state of the solver.
type SolverStatus = solver.SolverStatus

// SolverStatusResponse represents the status of a solver job.
type SolverStatusResponse = solver.SolverStatusResponse

const (
	// SolverStatusNotSolving indicates the solver is not currently working on the job.
	SolverStatusNotSolving = solver.SolverStatusNotSolving
	// SolverStatusSolvingActive indicates the solver is actively working on the job.
	SolverStatusSolvingActive = solver.SolverStatusSolvingActive
	// SolverStatusTerminatingEarly indicates the solver is stopping before its limits are reached.
	SolverStatusTerminatingEarly = solver.SolverStatusTerminatingEarly
	// SolverStatusSolvingScheduled indicates the job is queued and scheduled to start solving.
	SolverStatusSolvingScheduled = solver.SolverStatusSolvingScheduled
)
