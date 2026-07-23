package fieldservice

// FieldServiceRequest is the request model for starting a field service optimization.
type FieldServiceRequest struct {
	// Vehicles (technicians/resources) available to service the visits.
	Vehicles []Vehicle `json:"vehicles"`
	// Visits (jobs/stops) to be scheduled across the vehicles.
	Visits []Visit `json:"visits"`
	// Weights holds per-constraint weight overrides keyed by constraint name, in
	// the form Xhard/Ymedium/Zsoft, e.g. {"minimizeTravelTime": "0hard/0medium/1soft"}.
	// The server models this as an open map (additionalProperties: string), so any
	// constraint key and score string passes through unchanged — do not narrow it
	// to a fixed struct.
	Weights map[string]string `json:"weights,omitempty"`
	// Options holds optional termination controls for the solve.
	Options *SolverOptions `json:"options,omitempty"`
}

// Vehicle represents a vehicle in the field service optimization.
type Vehicle struct {
	// ID is the caller-supplied unique identifier for the vehicle.
	ID string `json:"id"`
	// Name is the human-readable vehicle/resource name.
	Name *string `json:"name,omitempty"`
	// Location is the home/depot coordinates as [latitude, longitude].
	Location [2]float64 `json:"location"`
	// Shifts are the working shifts (time ranges) during which the vehicle is available.
	Shifts []Shift `json:"shifts"`
	// Skills the vehicle provides, matched against each visit's required skills.
	Skills []string `json:"skills"`
	// DepartureTime is the time the vehicle departs its depot, as an ISO-8601 timestamp.
	DepartureTime *string `json:"departureTime,omitempty"`
	// Visits are the ordered ids of the visits assigned to this vehicle (populated in the solution).
	Visits []string `json:"visits,omitempty"`
}

// Visit represents a visit (job/task) to be scheduled.
type Visit struct {
	// ID is the caller-supplied unique identifier for the visit.
	ID string `json:"id"`
	// Name is the human-readable visit name.
	Name string `json:"name"`
	// Location is the visit coordinates as [latitude, longitude].
	Location [2]float64 `json:"location"`
	// TimeWindows are the allowed time windows during which service may start.
	TimeWindows []TimeWindow `json:"timeWindows"`
	// ServiceDuration is the on-site service time, as an ISO-8601 duration or seconds.
	ServiceDuration string `json:"serviceDuration"`
	// Priority is the relative importance of servicing the visit (e.g. LOW, MEDIUM, HIGH).
	Priority string `json:"priority"`
	// RequiredSkills are the skills a vehicle must have to service this visit.
	RequiredSkills []string `json:"requiredSkills"`
	// Pinned, when true, locks the existing assignment and leaves it unchanged by the solver.
	Pinned *bool `json:"pinned,omitempty"`
	// Vehicle is the solution id of the vehicle assigned to the visit.
	Vehicle *string `json:"vehicle,omitempty"`
	// ArrivalTime is the solution planned arrival time at the visit.
	ArrivalTime *string `json:"arrivalTime,omitempty"`
	// DepartureTime is the solution planned departure time from the visit.
	DepartureTime *string `json:"departureTime,omitempty"`
	// StartServiceTime is the solution time service is planned to start.
	StartServiceTime *string `json:"startServiceTime,omitempty"`
	// DrivingTimeSecondsFromPreviousStandstill is the solution driving time in seconds from the previous stop.
	DrivingTimeSecondsFromPreviousStandstill *int `json:"drivingTimeSecondsFromPreviousStandstill,omitempty"`
}

// Shift represents a time window during which a vehicle is available.
type Shift struct {
	// ID is the caller-supplied unique identifier for the shift.
	ID string `json:"id"`
	// MinStartTime is the earliest the shift may start.
	MinStartTime string `json:"minStartTime"`
	// MaxEndTime is the latest the shift may end.
	MaxEndTime string `json:"maxEndTime"`
}

// TimeWindow represents a time window for a visit.
type TimeWindow struct {
	// MinStartTime is the earliest permitted service start, as an ISO-8601 timestamp.
	MinStartTime string `json:"minStartTime"`
	// MaxEndTime is the latest permitted service end, as an ISO-8601 timestamp.
	MaxEndTime string `json:"maxEndTime"`
}

// SolverOptions contains solver configuration options.
type SolverOptions struct {
	// SpentLimit is the maximum total time the solver may run (e.g. PT30S), as an ISO-8601 duration.
	SpentLimit *string `json:"spentLimit,omitempty"`
	// UnimprovedSpentLimit stops the solver early after this much time passes with no score improvement (e.g. PT5S).
	UnimprovedSpentLimit *string `json:"unimprovedSpentLimit,omitempty"`
}

// FieldServiceStartResponse is the response from starting a field service
// optimization. It mirrors the API's SolverJobResponse shape.
type FieldServiceStartResponse struct {
	// JobID is the public PlanSolve job identifier - use it to poll status and fetch the solution.
	JobID string `json:"jobId"`
	// SolverJobID is the underlying solver engine's job identifier, when exposed.
	SolverJobID *string `json:"solverJobId,omitempty"`
	// Result is the inline solver result as raw JSON, when available synchronously.
	Result *string `json:"result,omitempty"`
	// Error is the error message when the solve request could not be accepted or run.
	Error *string `json:"error,omitempty"`
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
	// JobID is the public job identifier this result was fetched with (stamped client-side).
	JobID *string `json:"jobId,omitempty"`
	// Vehicles are the extended vehicles with additional metadata.
	Vehicles []ScheduledVehicle `json:"vehicles"`
	// Visits are the extended visits with additional metadata.
	Visits []ScheduledVisit `json:"visits"`
	// Score is the optimization score from the solver.
	Score *string `json:"score,omitempty"`
	// TotalDrivingTimeSeconds is the total driving time in seconds across all vehicles.
	TotalDrivingTimeSeconds int64 `json:"totalDrivingTimeSeconds"`
	// Weights are the weights in format Xhard/Ymedium/Zsoft.
	Weights map[string]string `json:"weights,omitempty"`
}

// ScheduledVehicle represents a vehicle in the optimization result. It
// corresponds to the API's ExtendedVehicle model.
type ScheduledVehicle struct {
	// ID is the caller-supplied unique identifier for the vehicle.
	ID string `json:"id"`
	// Name is the human-readable vehicle/resource name.
	Name *string `json:"name,omitempty"`
	// Location is the home/depot coordinates as [latitude, longitude].
	Location [2]float64 `json:"location"`
	// Skills the vehicle provides, matched against each visit's required skills.
	Skills []string `json:"skills"`
	// Shifts are the working shifts (time ranges) during which the vehicle is available.
	Shifts []Shift `json:"shifts"`
	// DepartureTime is the time the vehicle departs its depot, as an ISO-8601 timestamp.
	DepartureTime *string `json:"departureTime,omitempty"`
	// Visits are the ordered ids of the visits assigned to this vehicle.
	Visits []string `json:"visits"`
	// DailyReturnTimes is the return time for each day, keyed by ISO-8601 date.
	DailyReturnTimes map[string]string `json:"dailyReturnTimes,omitempty"`
	// ArrivalTime is the actual arrival time back at the depot.
	ArrivalTime *string `json:"arrivalTime,omitempty"`
	// TotalDrivingTimeSeconds is the total driving time in seconds for this vehicle.
	TotalDrivingTimeSeconds int64 `json:"totalDrivingTimeSeconds"`
}

// ScheduledVisit represents a visit in the optimization result. It corresponds
// to the API's ExtendedVisit model.
type ScheduledVisit struct {
	// ID is the caller-supplied unique identifier for the visit.
	ID string `json:"id"`
	// Name is the human-readable visit name.
	Name string `json:"name"`
	// Location is the visit coordinates as [latitude, longitude].
	Location [2]float64 `json:"location"`
	// Pinned, when true, locks the existing assignment and leaves it unchanged by the solver.
	Pinned *bool `json:"pinned,omitempty"`
	// TimeWindows are the allowed time windows during which service may start.
	TimeWindows []TimeWindow `json:"timeWindows"`
	// ServiceDuration is the service duration in seconds. The field-service solver
	// serializes a Java Duration in the RESULT as a number of seconds (e.g. 3600.0),
	// so this is a float64 — unlike the request Visit.ServiceDuration which is an
	// ISO-8601 string input like "PT1H".
	ServiceDuration float64 `json:"serviceDuration"`
	// Priority is the relative importance of servicing the visit (e.g. LOW, MEDIUM, HIGH).
	Priority string `json:"priority"`
	// RequiredSkills are the skills a vehicle must have to service this visit.
	RequiredSkills []string `json:"requiredSkills"`
	// Vehicle is the solution id of the vehicle assigned to the visit.
	Vehicle string `json:"vehicle"`
	// PreviousVisit is the id of the previous visit in the route.
	PreviousVisit *string `json:"previousVisit"`
	// PreviousVisitSameDay is the id of the previous visit on the same day.
	PreviousVisitSameDay *string `json:"previousVisitSameDay,omitempty"`
	// ArrivalTime is the solution planned arrival time at the visit.
	ArrivalTime string `json:"arrivalTime"`
	// MinStartTime is the minimum start time for the visit.
	MinStartTime *string `json:"minStartTime,omitempty"`
	// MaxEndTime is the maximum end time for the visit.
	MaxEndTime *string `json:"maxEndTime,omitempty"`
	// StartServiceTime is the solution time service is planned to start.
	StartServiceTime string `json:"startServiceTime"`
	// DepartureTime is the solution planned departure time from the visit.
	DepartureTime string `json:"departureTime"`
	// IsDayHead reports whether this is the first visit of the day.
	IsDayHead *bool `json:"isDayHead,omitempty"`
	// DrivingTimeSecondsFromPreviousStandstill is the solution driving time in seconds from the previous stop.
	DrivingTimeSecondsFromPreviousStandstill int `json:"drivingTimeSecondsFromPreviousStandstill"`
}
