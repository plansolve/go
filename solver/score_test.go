package solver

import (
	"encoding/json"
	"testing"
)

// scoreHolder uses the same alias pattern the response models use, so these
// tests exercise the exact mechanism that ships. encoding/json allocates a
// *Score before calling UnmarshalJSON on it, so a method on *Score cannot
// report "no score" by nilling itself - the decision has to be made by the
// containing struct, which is what ParseScoreJSON is for.
type scoreHolder struct {
	Score *Score `json:"score"`
}

func (h *scoreHolder) UnmarshalJSON(data []byte) error {
	type alias scoreHolder
	aux := &struct {
		Score json.RawMessage `json:"score"`
		*alias
	}{alias: (*alias)(h)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	score, err := ParseScoreJSON(aux.Score)
	if err != nil {
		return err
	}
	h.Score = score
	return nil
}

func TestScoreUnmarshalFromString(t *testing.T) {
	var h scoreHolder
	if err := json.Unmarshal([]byte(`{"score":"-1hard/0medium/25soft"}`), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if h.Score == nil {
		t.Fatal("expected a score")
	}
	if h.Score.Hard != -1 || h.Score.Medium != 0 || h.Score.Soft != 25 {
		t.Errorf("got %+v", *h.Score)
	}
}

func TestScoreUnmarshalFromSolverObject(t *testing.T) {
	// The shape both the shift and professionalservices solvers emit.
	raw := `{"score":{"hardScore":-2,"mediumScore":0,"softScore":-400,` +
		`"feasible":false,"zero":false,"solutionInitialized":true,"initScore":0}}`
	var h scoreHolder
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if h.Score == nil {
		t.Fatal("expected a score")
	}
	if h.Score.Hard != -2 || h.Score.Medium != 0 || h.Score.Soft != -400 {
		t.Errorf("got %+v", *h.Score)
	}
}

func TestScoreUnmarshalUnsolvedYieldsNoScore(t *testing.T) {
	// Review Focus 1: an unscored job sends zeros with solutionInitialized
	// false. Reporting Score{0,0,0} would claim a perfect feasible score.
	raw := `{"score":{"hardScore":0,"mediumScore":0,"softScore":0,` +
		`"feasible":false,"zero":true,"solutionInitialized":false,"initScore":0}}`
	var h scoreHolder
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if h.Score != nil {
		t.Errorf("expected no score, got %+v", *h.Score)
	}
}

func TestScoreUnmarshalPartialPlanYieldsRealScore(t *testing.T) {
	// A time limit stopped the solver before it considered every unit of
	// work: solutionInitialized is false, but the levels are real and must
	// not be discarded.
	raw := `{"score":{"hardScore":0,"mediumScore":0,"softScore":-400,` +
		`"feasible":true,"solutionInitialized":false}}`
	var h scoreHolder
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if h.Score == nil {
		t.Fatal("expected a score")
	}
	if h.Score.Hard != 0 || h.Score.Medium != 0 || h.Score.Soft != -400 {
		t.Errorf("got %+v", *h.Score)
	}
}

func TestScoreUnmarshalEmptyAndNull(t *testing.T) {
	// Review Focus 2: field service sends "" while solving. Polling hits the
	// result endpoint repeatedly, so this must not error.
	for _, raw := range []string{`{"score":""}`, `{"score":null}`, `{}`} {
		var h scoreHolder
		if err := json.Unmarshal([]byte(raw), &h); err != nil {
			t.Fatalf("unmarshal %s: %v", raw, err)
		}
		if h.Score != nil {
			t.Errorf("%s: expected no score, got %+v", raw, *h.Score)
		}
	}
}

func TestScoreUnmarshalBeyondInt32(t *testing.T) {
	// Review Focus 3: field service and PS are HardMediumSoftLongScore.
	var h scoreHolder
	if err := json.Unmarshal([]byte(`{"score":"0hard/0medium/-9000000000soft"}`), &h); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if h.Score.Soft != -9000000000 {
		t.Errorf("got soft %d, want -9000000000", h.Score.Soft)
	}

	var o scoreHolder
	raw := `{"score":{"hardScore":9000000000,"mediumScore":0,"softScore":0,"solutionInitialized":true}}`
	if err := json.Unmarshal([]byte(raw), &o); err != nil {
		t.Fatalf("unmarshal object: %v", err)
	}
	if o.Score.Hard != 9000000000 {
		t.Errorf("got hard %d, want 9000000000", o.Score.Hard)
	}
}

func TestScoreUnmarshalRejectsGarbage(t *testing.T) {
	// Review Focus 5: silently yielding zeros would look like a real score.
	var h scoreHolder
	if err := json.Unmarshal([]byte(`{"score":"not-a-score"}`), &h); err == nil {
		t.Fatal("expected an error")
	}
}

func TestParseScoreRejectsOutOfRangeInt64(t *testing.T) {
	// Review finding: the input is regex-valid but overflows int64. Silently
	// returning Score{MaxInt64, 0, 0} with a nil error would be exactly the
	// silent-wrong-answer class this whole change exists to eliminate.
	_, err := ParseScore("99999999999999999999hard/0medium/0soft")
	if err == nil {
		t.Fatal("expected an error for a hard level that overflows int64")
	}
}

func TestScoreMarshalJSONUsesCanonicalString(t *testing.T) {
	score := Score{Hard: 0, Medium: 0, Soft: -5}
	data, err := json.Marshal(score)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `"0hard/0medium/-5soft"` {
		t.Errorf("got %s, want canonical string", data)
	}
}

func TestScoreMarshalJSONRoundTrips(t *testing.T) {
	// A Score marshalled on its own (e.g. cached, logged) must be readable by
	// the same ParseScoreJSON-based decoding path used by every response type.
	original := scoreHolder{Score: &Score{Hard: -1, Medium: 2, Soft: -400}}

	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var roundTripped scoreHolder
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("unmarshal round-tripped JSON %s: %v", data, err)
	}

	if roundTripped.Score == nil || *roundTripped.Score != *original.Score {
		t.Errorf("score did not survive round trip: got %+v, want %+v", roundTripped.Score, original.Score)
	}
}
