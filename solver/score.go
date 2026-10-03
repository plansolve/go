package solver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
)

var scorePattern = regexp.MustCompile(`^(-?\d+)hard/(-?\d+)medium/(-?\d+)soft$`)

// Score represents an optimization score with three constraint levels.
//
// The solvers express a score two different ways on the wire: field service
// sends Timefold's canonical string notation ("0hard/0medium/-5soft"), while
// shift and professionalservices send an object keyed hardScore/mediumScore/
// softScore. UnmarshalJSON accepts both, so callers see one type either way.
// MarshalJSON always emits the canonical string form, so a Score round-trips
// through any SDK that only reads the string (every reader accepts it).
type Score struct {
	// Hard is the hard-constraint level; negative means constraints are violated.
	Hard int64
	// Medium is the medium-constraint level.
	Medium int64
	// Soft is the soft-constraint level, used to rank feasible solutions.
	Soft int64
}

// ParseScore parses a score string in the format "Xhard/Ymedium/Zsoft".
func ParseScore(value string) (Score, error) {
	if value == "" {
		return Score{}, fmt.Errorf("score value cannot be empty")
	}

	match := scorePattern.FindStringSubmatch(value)
	if match == nil {
		return Score{}, fmt.Errorf("invalid score format: '%s', expected format: 'Xhard/Ymedium/Zsoft'", value)
	}

	hard, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		return Score{}, fmt.Errorf("invalid score format: '%s', hard level out of range for int64", value)
	}
	medium, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		return Score{}, fmt.Errorf("invalid score format: '%s', medium level out of range for int64", value)
	}
	soft, err := strconv.ParseInt(match[3], 10, 64)
	if err != nil {
		return Score{}, fmt.Errorf("invalid score format: '%s', soft level out of range for int64", value)
	}

	return Score{Hard: hard, Medium: medium, Soft: soft}, nil
}

// String returns the score in "Xhard/Ymedium/Zsoft" format.
func (s Score) String() string {
	return fmt.Sprintf("%dhard/%dmedium/%dsoft", s.Hard, s.Medium, s.Soft)
}

// MarshalJSON emits the score in the canonical "Xhard/Ymedium/Zsoft" string
// form - the same wire shape field service sends, and the one every SDK's
// reader (including ParseScoreJSON) accepts. Without this, Score's zero-value
// struct tags would round-trip as {"hard":...,"medium":...,"soft":...}, which
// ParseScoreJSON rejects: a result cached to disk or logged as JSON could
// never be read back.
func (s Score) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// scoreObject is the shape the shift and professionalservices solvers emit.
// When those solvers are simplified to emit the canonical string like field
// service does, this type and the object branch below can be deleted without
// changing any public signature.
type scoreObject struct {
	HardScore           *int64 `json:"hardScore"`
	MediumScore         *int64 `json:"mediumScore"`
	SoftScore           *int64 `json:"softScore"`
	SolutionInitialized *bool  `json:"solutionInitialized"`
}

// ParseScoreJSON reads a raw JSON score field in either shape the solvers
// emit and reports nil when there is no score yet.
//
// This is a function rather than an UnmarshalJSON method because
// encoding/json allocates the *Score before calling UnmarshalJSON, so a
// method could never report absence - it would hand back a non-nil pointer to
// a zero Score, which is indistinguishable from a real feasible score of
// 0hard/0medium/0soft. Response models call this from their own
// UnmarshalJSON.
//
// nil is returned for: absent or null input, the empty string field service
// sends while solving, and the server's all-zero placeholder object (all
// three levels zero and flagged solutionInitialized false).
func ParseScoreJSON(data []byte) (*Score, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		if asString == "" {
			return nil, nil
		}
		parsed, err := ParseScore(asString)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}

	var obj scoreObject
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("score must be a score string or a score object: %w", err)
	}
	if obj.HardScore == nil || obj.MediumScore == nil || obj.SoftScore == nil {
		return nil, fmt.Errorf("score object must have hardScore, mediumScore and softScore")
	}
	// An unscored solution sends zeros with solutionInitialized false. A
	// partial plan (time limit stopped the solver early) also sends
	// solutionInitialized false, but with real, non-zero levels - that is a
	// meaningful score and must not be discarded.
	allZero := *obj.HardScore == 0 && *obj.MediumScore == 0 && *obj.SoftScore == 0
	if obj.SolutionInitialized != nil && !*obj.SolutionInitialized && allZero {
		return nil, nil
	}
	return &Score{Hard: *obj.HardScore, Medium: *obj.MediumScore, Soft: *obj.SoftScore}, nil
}
