package plansolve

import "github.com/plansolve/go/solver"

// The analysis and job types live in the leaf `solver` package for the same
// reason as the status types (see solver_status.go); re-exported here as aliases.

// ScoreAnalysisResponse is the answer of the solvers' Analyze operation.
type ScoreAnalysisResponse = solver.ScoreAnalysisResponse

// ConstraintAnalysisResponse is one constraint's contribution to a plan's score.
type ConstraintAnalysisResponse = solver.ConstraintAnalysisResponse

// ConstraintMatchResponse is one match of a constraint.
type ConstraintMatchResponse = solver.ConstraintMatchResponse

// SolverJobResponse is the answer of a solver's Start and of Jobs.CreateJob.
type SolverJobResponse = solver.SolverJobResponse
