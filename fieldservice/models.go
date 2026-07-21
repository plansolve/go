package fieldservice

// FieldServiceRequest is the request model for starting a field service optimization.
type FieldServiceRequest struct {
	Vehicles []Vehicle      `json:"vehicles"`
	Visits   []Visit        `json:"visits"`
	Weights  *Weights       `json:"weights,omitempty"`
	Options  *SolverOptions `json:"options,omitempty"`
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

// Weights contains configuration weights for field service optimization constraints.
type Weights struct {
	PinnedVisitVehicleAssignment  *ConstraintWeight `json:"pinnedVisitVehicleAssignment,omitempty"`
	PinnedVisitServiceTime        *ConstraintWeight `json:"pinnedVisitServiceTime,omitempty"`
	NoMissingSkills               *ConstraintWeight `json:"noMissingSkills,omitempty"`
	ServiceTimeMissing            *ConstraintWeight `json:"serviceTimeMissing,omitempty"`
	VehicleUnassigned             *ConstraintWeight `json:"vehicleUnassigned,omitempty"`
	PreferHighPriority            *ConstraintWeight `json:"preferHighPriority,omitempty"`
	MinimizeTravelTime            *ConstraintWeight `json:"minimizeTravelTime,omitempty"`
	PreferUsingIdleVehicles       *ConstraintWeight `json:"preferUsingIdleVehicles,omitempty"`
	MinimizeEarlyArrivalWait      *ConstraintWeight `json:"minimizeEarlyArrivalWait,omitempty"`
	PreferInitialVehicleAssignment *ConstraintWeight `json:"preferInitialVehicleAssignment,omitempty"`
	PreferEarlierVisitDates       *ConstraintWeight `json:"PreferEarlierVisitDates,omitempty"`
}

// ConstraintWeight represents a constraint weight in the format "Xhard/Ymedium/Zsoft".
// Only one of Hard, Medium, or Soft can be non-zero at a time.
type ConstraintWeight struct {
	Hard   int `json:"hard"`
	Medium int `json:"medium"`
	Soft   int `json:"soft"`
}

// SolverOptions contains solver configuration options.
type SolverOptions struct {
	// Maximum time the solver can spend (ISO-8601 duration format, e.g., "PT5M" for 5 minutes).
	SpentLimit *string `json:"spentLimit,omitempty"`
	// Maximum time the solver can spend without improvement (ISO-8601 duration format).
	UnimprovedSpentLimit *string `json:"unimprovedSpentLimit,omitempty"`
}

// FieldServiceStartResponse is the response from starting a field service optimization.
type FieldServiceStartResponse struct {
	JobID      string  `json:"jobId"`
	Result     *string `json:"result,omitempty"`
	Error      *string `json:"error,omitempty"`
	StatusCode int     `json:"statusCode,omitempty"`
}

// FieldServiceResultResponse is the response from getting field service optimization results.
type FieldServiceResultResponse struct {
	JobID                  string                `json:"jobId,omitempty"`
	Vehicles               []ScheduledVehicle    `json:"vehicles"`
	Visits                 []ScheduledVisit      `json:"visits"`
	Score                  *string               `json:"score,omitempty"`
	TotalDrivingTimeSeconds int64                `json:"totalDrivingTimeSeconds"`
	Weights                map[string]string     `json:"weights,omitempty"`
}

// ScheduledVehicle represents a vehicle in the optimization result.
type ScheduledVehicle struct {
	ID                      string                 `json:"id"`
	Location                [2]float64             `json:"location"`
	Shifts                  []Shift                `json:"shifts"`
	Skills                  []string               `json:"skills"`
	Visits                  []string               `json:"visits"`
	DailyReturnTimes        map[string]interface{} `json:"dailyReturnTimes,omitempty"`
	TotalDrivingTimeSeconds *int64                 `json:"totalDrivingTimeSeconds,omitempty"`
	ArrivalTime             *string                `json:"arrivalTime,omitempty"`
	DepartureTime           *string                `json:"departureTime,omitempty"`
}

// ScheduledVisit represents a visit in the optimization result.
type ScheduledVisit struct {
	ID                                       string       `json:"id"`
	Name                                     string       `json:"name"`
	Location                                 [2]float64   `json:"location"`
	TimeWindows                              []TimeWindow `json:"timeWindows"`
	ServiceDuration                          string       `json:"serviceDuration"`
	Priority                                 string       `json:"priority"`
	RequiredSkills                           []string     `json:"requiredSkills"`
	Vehicle                                  string       `json:"vehicle"`
	PreviousVisit                            *string      `json:"previousVisit"`
	ArrivalTime                              string       `json:"arrivalTime"`
	DepartureTime                            string       `json:"departureTime"`
	StartServiceTime                         string       `json:"startServiceTime"`
	DrivingTimeSecondsFromPreviousStandstill int          `json:"drivingTimeSecondsFromPreviousStandstill"`
}
