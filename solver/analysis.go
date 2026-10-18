package solver

import "encoding/json"

// ScoreAnalysisResponse is the answer of the /{jobId}/analyze operations: how
// each constraint contributed to the plan's score. It mirrors the API's
// ScoreAnalysisResponse and is shared by all three solver clients.
type ScoreAnalysisResponse struct {
	// Score is the total score of the plan, or nil when the analysis carries none.
	Score *Score `json:"score,omitempty"`
	// Initialized reports whether every planning entity had been assigned when
	// the plan was scored.
	Initialized bool `json:"initialized"`
	// Constraints holds one entry per constraint of the solver, with its weight,
	// its share of the score and the matches that caused it.
	Constraints []ConstraintAnalysisResponse `json:"constraints,omitempty"`
}

// ConstraintAnalysisResponse is one constraint's contribution to a plan's score.
type ConstraintAnalysisResponse struct {
	// Name is the constraint name, the same name the request's weights accept.
	Name *string `json:"name,omitempty"`
	// Package is the package the constraint is defined in; identifies the solver.
	Package *string `json:"package,omitempty"`
	// Weight is the score one match carries, or nil when absent.
	Weight *Score `json:"weight,omitempty"`
	// Score is this constraint's share of the total score, or nil when absent.
	Score *Score `json:"score,omitempty"`
	// MatchCount is the number of matches; equals len(Matches) when they are included.
	MatchCount int32 `json:"matchCount"`
	// Matches are the individual matches, each with its own score and the entities involved.
	Matches []ConstraintMatchResponse `json:"matches,omitempty"`
}

// ConstraintMatchResponse is one match of a constraint: the score it
// contributed and the entities that caused it.
type ConstraintMatchResponse struct {
	// Score is the score of this single match, or nil when absent.
	Score *Score `json:"score,omitempty"`
	// Justification is the solver-specific description of the entities involved
	// (ids of the visit, vehicle, shift, employee or task), kept as raw JSON.
	Justification json.RawMessage `json:"justification,omitempty"`
}

// UnmarshalJSON decodes the analysis, routing the score field through
// ParseScoreJSON like SolverStatusResponse does.
func (r *ScoreAnalysisResponse) UnmarshalJSON(data []byte) error {
	type alias ScoreAnalysisResponse
	aux := &struct {
		Score json.RawMessage `json:"score"`
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
	return nil
}

// UnmarshalJSON decodes the constraint, routing the weight and score fields
// through ParseScoreJSON.
func (c *ConstraintAnalysisResponse) UnmarshalJSON(data []byte) error {
	type alias ConstraintAnalysisResponse
	aux := &struct {
		Weight json.RawMessage `json:"weight"`
		Score  json.RawMessage `json:"score"`
		*alias
	}{alias: (*alias)(c)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	weight, err := ParseScoreJSON(aux.Weight)
	if err != nil {
		return err
	}
	score, err := ParseScoreJSON(aux.Score)
	if err != nil {
		return err
	}
	c.Weight = weight
	c.Score = score
	return nil
}

// UnmarshalJSON decodes the match, routing the score field through
// ParseScoreJSON. A JSON null justification is kept as nil.
func (m *ConstraintMatchResponse) UnmarshalJSON(data []byte) error {
	type alias ConstraintMatchResponse
	aux := &struct {
		Score json.RawMessage `json:"score"`
		*alias
	}{alias: (*alias)(m)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	score, err := ParseScoreJSON(aux.Score)
	if err != nil {
		return err
	}
	m.Score = score
	if string(m.Justification) == "null" {
		m.Justification = nil
	}
	return nil
}
