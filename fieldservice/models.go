package fieldservice

// FieldServiceRequest is the request model for starting a field service optimization.
type FieldServiceRequest struct {
	Vehicles []Vehicle      `json:"vehicles"`
	Visits   []Visit        `json:"visits"`
	// Weights maps a constraint name to a score-notation string, e.g.
	// {"minimizeTravelTime": "0hard/0medium/1soft"}. The server models this as an
	// open map (additionalProperties: string), so any constraint key and score
	// string passes through unchanged — do not narrow it to a fixed struct.
	Weights map[string]string `json:"weights,omitempty"`
	Options *SolverOptions    `json:"options,omitempty"`
}

// Vehicle represents a vehicle in the field service optimization.
type Vehicle struct {
	ID            string     `json:"id"`
	Name          *string    `json:"name,omitempty"`
	Location      [2]float64 `json:"location"`
	Shifts        []Shift    `json:"shifts"`
	Skills        []string   `json:"skills"`
	DepartureTime *string    `json:"departureTime,omitempty"`
	Visits        []string   `json:"visits,omitempty"`
}

// Visit represents a visit (job/task) to be scheduled.
type Visit struct {
	ID                                       string       `json:"id"`
	Name                                     string       `json:"name"`
	Location                                 [2]float64   `json:"location"`
	TimeWindows                              []TimeWindow `json:"timeWindows"`
	ServiceDuration                          string       `json:"serviceDuration"`
	Priority                                 string       `json:"priority"`
	RequiredSkills                           []string     `json:"requiredSkills"`
	Pinned                                   *bool        `json:"pinned,omitempty"`
	Vehicle                                  *string      `json:"vehicle,omitempty"`
	ArrivalTime                              *string      `json:"arrivalTime,omitempty"`
	DepartureTime                            *string      `json:"departureTime,omitempty"`
	StartServiceTime                         *string      `json:"startServiceTime,omitempty"`
	DrivingTimeSecondsFromPreviousStandstill *int         `json:"drivingTimeSecondsFromPreviousStandstill,omitempty"`
}

// Shift represents a time window during which a vehicle is available.
type Shift struct {
	ID           string `json:"id"`
	MinStartTime string `json:"minStartTime"`
	MaxEndTime   string `json:"maxEndTime"`
}

// TimeWindow represents a time window for a visit.
type TimeWindow struct {
	MinStartTime string `json:"minStartTime"`
	MaxEndTime   string `json:"maxEndTime"`
}

// SolverOptions contains solver configuration options.
type SolverOptions struct {
	// Maximum time the solver can spend (ISO-8601 duration format, e.g., "PT5M" for 5 minutes).
	SpentLimit *string `json:"spentLimit,omitempty"`
	// Maximum time the solver can spend without improvement (ISO-8601 duration format).
	UnimprovedSpentLimit *string `json:"unimprovedSpentLimit,omitempty"`
}

// FieldServiceStartResponse is the response from starting a field service
// optimization. It mirrors the API's SolverJobResponse shape.
type FieldServiceStartResponse struct {
	JobID       string  `json:"jobId"`
	SolverJobID *string `json:"solverJobId,omitempty"`
	Result      *string `json:"result,omitempty"`
	Error       *string `json:"error,omitempty"`
}

// Spec-name aliases. The API contract names the result model FieldServiceSolution
// and its sub-types ExtendedVehicle / ExtendedVisit; the SDK keeps its idiomatic
// names but exposes the spec names as aliases for discoverability.
type (
	FieldServiceSolution = FieldServiceResultResponse
	ExtendedVehicle      = ScheduledVehicle
	ExtendedVisit        = ScheduledVisit
)

// FieldServiceResultResponse is the response from getting field service optimization results.
type FieldServiceResultResponse struct {
	JobID                  *string               `json:"jobId,omitempty"`
	Vehicles               []ScheduledVehicle    `json:"vehicles"`
	Visits                 []ScheduledVisit      `json:"visits"`
	Score                  *string               `json:"score,omitempty"`
	TotalDrivingTimeSeconds int64                `json:"totalDrivingTimeSeconds"`
	Weights                map[string]string     `json:"weights,omitempty"`
}

// ScheduledVehicle represents a vehicle in the optimization result. It
// corresponds to the API's ExtendedVehicle model.
type ScheduledVehicle struct {
	ID                      string            `json:"id"`
	Name                    *string           `json:"name,omitempty"`
	Location                [2]float64        `json:"location"`
	Skills                  []string          `json:"skills"`
	Shifts                  []Shift           `json:"shifts"`
	DepartureTime           *string           `json:"departureTime,omitempty"`
	Visits                  []string          `json:"visits"`
	DailyReturnTimes        map[string]string `json:"dailyReturnTimes,omitempty"`
	ArrivalTime             *string           `json:"arrivalTime,omitempty"`
	TotalDrivingTimeSeconds int64             `json:"totalDrivingTimeSeconds"`
}

// ScheduledVisit represents a visit in the optimization result. It corresponds
// to the API's ExtendedVisit model.
type ScheduledVisit struct {
	ID                                       string       `json:"id"`
	Name                                     string       `json:"name"`
	Location                                 [2]float64   `json:"location"`
	Pinned                                   *bool        `json:"pinned,omitempty"`
	TimeWindows                              []TimeWindow `json:"timeWindows"`
	// ServiceDuration is the service duration in seconds. The field-service solver
	// serializes a Java Duration in the RESULT as a number of seconds (e.g. 3600.0),
	// so this is a float64 — unlike the request Visit.ServiceDuration which is an
	// ISO-8601 string input like "PT1H".
	ServiceDuration                          float64      `json:"serviceDuration"`
	Priority                                 string       `json:"priority"`
	RequiredSkills                           []string     `json:"requiredSkills"`
	Vehicle                                  string       `json:"vehicle"`
	PreviousVisit                            *string      `json:"previousVisit"`
	PreviousVisitSameDay                     *string      `json:"previousVisitSameDay,omitempty"`
	ArrivalTime                              string       `json:"arrivalTime"`
	MinStartTime                             *string      `json:"minStartTime,omitempty"`
	MaxEndTime                               *string      `json:"maxEndTime,omitempty"`
	StartServiceTime                         string       `json:"startServiceTime"`
	DepartureTime                            string       `json:"departureTime"`
	IsDayHead                                *bool        `json:"isDayHead,omitempty"`
	DrivingTimeSecondsFromPreviousStandstill int          `json:"drivingTimeSecondsFromPreviousStandstill"`
}
