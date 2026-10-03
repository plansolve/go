package plansolve

import "github.com/plansolve/go/solver"

// Score represents an optimization score in the format "Xhard/Ymedium/Zsoft".
// It is an alias for solver.Score, which lives in the leaf package so the
// per-service subpackages can use it without an import cycle.
type Score = solver.Score

// ParseScore parses a score string in the format "Xhard/Ymedium/Zsoft".
func ParseScore(value string) (Score, error) {
	return solver.ParseScore(value)
}
