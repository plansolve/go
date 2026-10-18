package solver

import (
	"encoding/json"
	"testing"
)

func TestScoreAnalysisResponseUnmarshal(t *testing.T) {
	body := `{
		"score": "0hard/0medium/-44soft",
		"initialized": true,
		"constraints": [{
			"name": "Minimize travel time",
			"package": "ai.timefold.solver.fieldservice",
			"weight": "0hard/0medium/-1soft",
			"score": "0hard/0medium/-44soft",
			"matchCount": 1,
			"matches": [
				{"score": "0hard/0medium/-44soft", "justification": {"vehicleId": "vehicle-1"}},
				{"score": null, "justification": null}
			]
		}]
	}`

	var analysis ScoreAnalysisResponse
	if err := json.Unmarshal([]byte(body), &analysis); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if analysis.Score == nil || *analysis.Score != (Score{Hard: 0, Medium: 0, Soft: -44}) {
		t.Errorf("unexpected total score %v", analysis.Score)
	}
	if !analysis.Initialized {
		t.Error("expected initialized true")
	}
	if len(analysis.Constraints) != 1 {
		t.Fatalf("expected 1 constraint, got %d", len(analysis.Constraints))
	}

	c := analysis.Constraints[0]
	if c.Name == nil || *c.Name != "Minimize travel time" {
		t.Errorf("unexpected name %v", c.Name)
	}
	if c.Package == nil || *c.Package != "ai.timefold.solver.fieldservice" {
		t.Errorf("unexpected package %v", c.Package)
	}
	if c.Weight == nil || c.Weight.Soft != -1 {
		t.Errorf("unexpected weight %v", c.Weight)
	}
	if c.Score == nil || c.Score.Soft != -44 {
		t.Errorf("unexpected constraint score %v", c.Score)
	}
	if c.MatchCount != 1 || len(c.Matches) != 2 {
		t.Fatalf("unexpected matches: count %d, len %d", c.MatchCount, len(c.Matches))
	}

	m := c.Matches[0]
	if m.Score == nil || m.Score.Soft != -44 {
		t.Errorf("unexpected match score %v", m.Score)
	}
	var justification map[string]string
	if err := json.Unmarshal(m.Justification, &justification); err != nil {
		t.Fatalf("justification is not raw JSON: %v", err)
	}
	if justification["vehicleId"] != "vehicle-1" {
		t.Errorf("unexpected justification %s", m.Justification)
	}

	if c.Matches[1].Score != nil || c.Matches[1].Justification != nil {
		t.Errorf("expected null score and justification to read as nil, got %v / %s", c.Matches[1].Score, c.Matches[1].Justification)
	}
}

func TestScoreAnalysisResponseNullFields(t *testing.T) {
	var analysis ScoreAnalysisResponse
	if err := json.Unmarshal([]byte(`{"score": null, "initialized": false, "constraints": null}`), &analysis); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if analysis.Score != nil || analysis.Initialized || analysis.Constraints != nil {
		t.Errorf("unexpected analysis %+v", analysis)
	}
}
