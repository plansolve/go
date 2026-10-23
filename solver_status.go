package plansolve

import "github.com/plansolve/go/solver"

// The solver status types live in the leaf `solver` package to avoid an import
// cycle (the subpackages need them too, and they can't import this root package
// because this package imports them). Re-exported here as aliases so the public
// API - plansolve.SolverStatusResponse, plansolve.SolverStatusCompleted, etc. -
// stays available from the root package.

// SolverStatus represents the state of the solver.
type SolverStatus = solver.SolverStatus

// SolverStatusResponse represents the status of a solver job.
type SolverStatusResponse = solver.SolverStatusResponse

const (
	// SolverStatusQueued indicates the job is queued and has not started solving yet.
	SolverStatusQueued = solver.SolverStatusQueued
	// SolverStatusSolving indicates the solver is actively working on the job.
	SolverStatusSolving = solver.SolverStatusSolving
	// SolverStatusStopping indicates the solver is stopping before its limits are reached.
	SolverStatusStopping = solver.SolverStatusStopping
	// SolverStatusCompleted indicates the solver has finished working on the job.
	SolverStatusCompleted = solver.SolverStatusCompleted
)

// The 0.34 names, kept so existing code compiles. They equal the public values.
const (
	// SolverStatusNotSolving is the 0.34 name for SolverStatusCompleted.
	//
	// Deprecated: use SolverStatusCompleted.
	SolverStatusNotSolving = solver.SolverStatusCompleted
	// SolverStatusSolvingActive is the 0.34 name for SolverStatusSolving.
	//
	// Deprecated: use SolverStatusSolving.
	SolverStatusSolvingActive = solver.SolverStatusSolving
	// SolverStatusTerminatingEarly is the 0.34 name for SolverStatusStopping.
	//
	// Deprecated: use SolverStatusStopping.
	SolverStatusTerminatingEarly = solver.SolverStatusStopping
	// SolverStatusSolvingScheduled is the 0.34 name for SolverStatusQueued.
	//
	// Deprecated: use SolverStatusQueued.
	SolverStatusSolvingScheduled = solver.SolverStatusQueued
)
