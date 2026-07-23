package plansolve

import (
	"fmt"
	"regexp"
	"strconv"
)

var scorePattern = regexp.MustCompile(`^(-?\d+)hard/(-?\d+)medium/(-?\d+)soft$`)

// Score represents an optimization score in the format "Xhard/Ymedium/Zsoft".
type Score struct {
	// Hard is the hard-constraint score level; non-zero means constraints are violated.
	Hard int `json:"hard"`
	// Medium is the medium-constraint score level.
	Medium int `json:"medium"`
	// Soft is the soft-constraint score level, used to rank feasible solutions.
	Soft int `json:"soft"`
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

	hard, _ := strconv.Atoi(match[1])
	medium, _ := strconv.Atoi(match[2])
	soft, _ := strconv.Atoi(match[3])

	return Score{Hard: hard, Medium: medium, Soft: soft}, nil
}

// String returns the score in "Xhard/Ymedium/Zsoft" format.
func (s Score) String() string {
	return fmt.Sprintf("%dhard/%dmedium/%dsoft", s.Hard, s.Medium, s.Soft)
}
