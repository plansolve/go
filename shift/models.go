package shift

// ShiftRequest is the request model for starting a shift optimization. It
// corresponds to the API's ShiftAssignmentRequest.
type ShiftRequest struct {
	// ID is the optional caller-supplied identifier for the plan.
	ID *string `json:"id,omitempty"`
	// Name is the plan display name.
	Name *string `json:"name,omitempty"`
	// Description is the free-text plan description.
	Description *string `json:"description,omitempty"`
	// Contracts available in this plan. Employees carry their own contract; this list is a
	// convenience catalogue for callers that reuse the same rules across employees.
	Contracts []Contract `json:"contracts,omitempty"`
	// Shifts that need to be staffed.
	Shifts []ShiftAssignment `json:"shifts,omitempty"`
	// Employees available to be assigned to shifts.
	Employees []ShiftEmployee `json:"employees,omitempty"`
	// TimeOffRequests are employee requests to be given a period off. Soft - honored when possible.
	TimeOffRequests []ShiftTimeOffRequest `json:"timeOffRequests,omitempty"`
	// Unavailabilities are periods employees cannot work. Hard - never scheduled over.
	Unavailabilities []ShiftUnavailability `json:"unavailabilities,omitempty"`
	// ShiftOffRequests are employee requests to avoid specific shifts.
	ShiftOffRequests []ShiftOffRequest `json:"shiftOffRequests,omitempty"`
	// Fairness holds the fairness configuration.
	Fairness *Fairness `json:"fairness,omitempty"`
	// Hook is an optional callback URL invoked when solving completes.
	Hook *string `json:"hook,omitempty"`
	// Weights holds per-constraint weight overrides keyed by constraint name, in
	// the form Xhard/Ymedium/Zsoft, e.g. {"shiftOffRequest": "0hard/0medium/4soft"}.
	// The server models this as an open map (additionalProperties: string), so any
	// constraint key and score string passes through unchanged.
	Weights map[string]string `json:"weights,omitempty"`
	// Options holds optional termination controls for the solve.
	Options *SolverOptions `json:"options,omitempty"`
}

// ShiftAssignment represents a shift that needs to be staffed, and - on the
// solution side - the employee the solver assigned to it.
type ShiftAssignment struct {
	// Name is the shift name, used as its identifier.
	Name *string `json:"name,omitempty"`
	// From is the shift start, as an ISO-8601 timestamp.
	From *string `json:"from,omitempty"`
	// To is the shift end, as an ISO-8601 timestamp.
	To *string `json:"to,omitempty"`
	// Skills an employee must have to be assigned this shift.
	Skills []string `json:"skills,omitempty"`
	// DesiredSkills are skills that are preferred but not required for this shift.
	DesiredSkills []string `json:"desiredSkills,omitempty"`
	// Tags are arbitrary tags used for grouping and affinity constraints.
	Tags []string `json:"tags,omitempty"`
	// CostFactor is the multiplier applied to the assigned employee's hourly cost for this shift.
	// Omit for the solver default (1.0).
	CostFactor *float64 `json:"costFactor,omitempty"`
	// Value is the number of employees this shift needs. The solver expands a shift with
	// value > 1 into one entity per required employee, named "{name}#1", "{name}#2", ... and
	// carrying OriginalName. Omit for the solver default (1).
	Value *int `json:"value,omitempty"`
	// Priority is the relative priority of staffing this shift. Omit for the solver default (5).
	Priority *int `json:"priority,omitempty"`
	// OriginalName is the name of the shift this one was expanded from when a shift needs
	// several employees; the solver splits it into one instance per required employee.
	OriginalName *string `json:"originalName,omitempty"`
	// PinnedByUser, when true, locks the existing assignment and leaves it unchanged by the solver.
	PinnedByUser bool `json:"pinnedByUser"`
	// AssignedEmployee (solution) is the id of the employee the solver assigned to this shift,
	// or nil if unstaffed. The wire format is the bare employee id, not a nested object.
	AssignedEmployee *string `json:"assignedEmployee,omitempty"`
}

// ShiftEmployee represents an employee available for shift assignment.
type ShiftEmployee struct {
	// ID is the employee identifier, referenced by time-off, unavailability and shift-off requests.
	ID *string `json:"id,omitempty"`
	// Name is the employee display name.
	Name *string `json:"name,omitempty"`
	// Skills the employee possesses, matched against each shift's required skills.
	Skills []string `json:"skills,omitempty"`
	// LastShiftEnd is the end of the employee's last shift before the planning window, as an
	// ISO-8601 timestamp. Seeds the rest-between-shifts constraint across the window boundary.
	LastShiftEnd *string `json:"lastShiftEnd,omitempty"`
	// PreferredShifts are the names of shifts the employee prefers to work.
	PreferredShifts []string `json:"preferredShifts,omitempty"`
	// Tags are arbitrary tags used for grouping and affinity constraints.
	Tags []string `json:"tags,omitempty"`
	// CostPerHour is the hourly cost of this employee, used by cost-minimization.
	CostPerHour float64 `json:"costPerHour"`
	// Contract is the name of the Contract governing this employee, resolved against the plan's
	// contracts. The wire format is the bare name (e.g. "FULL_TIME"), not a nested object.
	Contract *string `json:"contract,omitempty"`
}

// Contract holds working-time rules that constrain how the employee holding it can be
// scheduled. It corresponds to the API's ShiftContract.
type Contract struct {
	// Name is the contract name, referenced by an employee's contract field.
	Name *string `json:"name,omitempty"`
	// MinWorkDurationPerWeek is the minimum working time per week, as an ISO-8601 duration.
	MinWorkDurationPerWeek *string `json:"minWorkDurationPerWeek,omitempty"`
	// MaxWorkDurationPerWeek is the maximum working time per week, as an ISO-8601 duration.
	MaxWorkDurationPerWeek *string `json:"maxWorkDurationPerWeek,omitempty"`
	// MaxConsecutiveWorkDays is the maximum number of consecutive days that may be worked.
	// Omit for the solver default (5).
	MaxConsecutiveWorkDays *int `json:"maxConsecutiveWorkDays,omitempty"`
	// MaxShiftsDay is the maximum number of shifts allowed on a single day. Omit for the solver default (2).
	MaxShiftsDay *int `json:"maxShiftsDay,omitempty"`
	// MinRestBetweenShifts is the minimum rest between two consecutive shifts, as an ISO-8601
	// duration. Omit for the solver default (PT11H).
	MinRestBetweenShifts *string `json:"minRestBetweenShifts,omitempty"`
	// MaxWorkingDaysPerWeek is the maximum number of working days per week. Omit for the solver default (6).
	MaxWorkingDaysPerWeek *int `json:"maxWorkingDaysPerWeek,omitempty"`
	// LatestShiftEnd is the latest permitted shift end time (time of day, e.g. 22:00:00).
	LatestShiftEnd *string `json:"latestShiftEnd,omitempty"`
	// EarliestShiftStart is the earliest permitted shift start time (time of day, e.g. 06:00:00).
	EarliestShiftStart *string `json:"earliestShiftStart,omitempty"`
	// MinConsecutiveDaysOff is the minimum number of consecutive days off between working
	// stretches. Omit for the solver default (1).
	MinConsecutiveDaysOff *int `json:"minConsecutiveDaysOff,omitempty"`
}

// SolverOptions holds optional termination controls for a shift solve, expressed
// as ISO-8601 durations. It corresponds to the API's ShiftSolverOptions.
type SolverOptions struct {
	// SpentLimit is the maximum total time the solver may run (e.g. PT30S), as an ISO-8601 duration.
	SpentLimit *string `json:"spentLimit,omitempty"`
	// UnimprovedSpentLimit stops the solver early after this much time passes with no score improvement (e.g. PT5S).
	UnimprovedSpentLimit *string `json:"unimprovedSpentLimit,omitempty"`
}

// Fairness represents fairness configuration for the shift optimizer.
type Fairness struct {
	// FairnessBuckets define the groups the fairness objective balances over.
	FairnessBuckets []FairnessBucket `json:"fairnessBuckets,omitempty"`
}

// FairnessBucket represents a fairness bucket for shift distribution.
type FairnessBucket struct {
	// Name is the name identifier for the fairness bucket.
	Name *string `json:"name,omitempty"`
	// EmployeeIDs are the ids of the employees in this bucket.
	EmployeeIDs []string `json:"employeeIds,omitempty"`
	// Shifts is the list of shift names in this bucket.
	Shifts []string `json:"shifts,omitempty"`
	// Period is the time period for this bucket (e.g. 2024-01-01/2024-01-08).
	Period *string `json:"period,omitempty"`
}

// ShiftTimeOffRequest is an employee's request to be given a period off. Soft:
// the solver honors it when it can.
type ShiftTimeOffRequest struct {
	// ID is the caller-supplied identifier for the request.
	ID *string `json:"id,omitempty"`
	// EmployeeID is the id of the employee making the request.
	EmployeeID *string `json:"employeeId,omitempty"`
	// From is the start of the requested period off, as an ISO-8601 timestamp.
	From *string `json:"from,omitempty"`
	// To is the end of the requested period off, as an ISO-8601 timestamp.
	To *string `json:"to,omitempty"`
}

// ShiftUnavailability is a period an employee cannot work. Hard: unlike a
// time-off request, which the solver only tries to honor, it will never assign
// a shift overlapping this period.
type ShiftUnavailability struct {
	// ID is the caller-supplied identifier for the unavailability.
	ID *string `json:"id,omitempty"`
	// EmployeeID is the id of the unavailable employee.
	EmployeeID *string `json:"employeeId,omitempty"`
	// From is the start of the unavailable period, as an ISO-8601 timestamp.
	From *string `json:"from,omitempty"`
	// To is the end of the unavailable period, as an ISO-8601 timestamp.
	To *string `json:"to,omitempty"`
}

// ShiftOffRequest is an employee's request to not be assigned a specific shift.
type ShiftOffRequest struct {
	// ID is the caller-supplied identifier for the request.
	ID *string `json:"id,omitempty"`
	// EmployeeID is the id of the employee making the request.
	EmployeeID *string `json:"employeeId,omitempty"`
	// ShiftName is the name of the shift the employee wants to avoid.
	ShiftName *string `json:"shiftName,omitempty"`
	// Weight is the strength of the preference; higher values make honoring it more important.
	// Omit for the solver default (1).
	Weight *int `json:"weight,omitempty"`
}

// ShiftStartResponse is the response from starting a shift optimization.
type ShiftStartResponse struct {
	// JobID is the public PlanSolve job identifier - use it to poll status and fetch the solution.
	JobID string `json:"jobId"`
	// SolverJobID is the underlying solver engine's job identifier, when exposed.
	SolverJobID *string `json:"solverJobId,omitempty"`
	// Result is the inline solver result as raw JSON, when available synchronously.
	Result *string `json:"result,omitempty"`
	// Error is the error message when the solve request could not be accepted or run.
	Error *string `json:"error,omitempty"`
}

// ShiftResultResponse is the response from getting shift optimization results.
// It corresponds to the API's ShiftAssignmentResponse: the echoed input PLUS
// the solution (feasible, scoreString, score, unassignedShifts, assignedShifts).
// Each shift carries its own AssignedEmployee; there is no employee-side
// back-reference. It carries no jobId on the wire.
type ShiftResultResponse struct {
	// JobID is the public job identifier this result was fetched with (stamped client-side).
	JobID *string `json:"jobId,omitempty"`
	// ID is the identifier of the solved plan.
	ID *string `json:"id,omitempty"`
	// Name is the echoed plan name.
	Name *string `json:"name,omitempty"`
	// Description is the echoed plan description.
	Description *string `json:"description,omitempty"`
	// Contracts are the echoed contracts from the request.
	Contracts []Contract `json:"contracts,omitempty"`
	// Employees are the echoed employees from the request.
	Employees []ShiftEmployee `json:"employees,omitempty"`
	// Shifts are all shifts with their assignments populated by the solver.
	Shifts []ShiftAssignment `json:"shifts,omitempty"`
	// TimeOffRequests are the echoed time-off requests from the request.
	TimeOffRequests []ShiftTimeOffRequest `json:"timeOffRequests,omitempty"`
	// Unavailabilities are the echoed unavailabilities from the request.
	Unavailabilities []ShiftUnavailability `json:"unavailabilities,omitempty"`
	// ShiftOffRequests are the echoed shift-off requests from the request.
	ShiftOffRequests []ShiftOffRequest `json:"shiftOffRequests,omitempty"`
	// Fairness is the echoed fairness configuration.
	Fairness *Fairness `json:"fairness,omitempty"`
	// Hook is the echoed completion callback URL.
	Hook *string `json:"hook,omitempty"`
	// ConstraintWeightOverrides are the echoed per-constraint weight overrides, keyed by
	// constraint name with a Xhard/Ymedium/Zsoft value.
	ConstraintWeightOverrides map[string]string `json:"constraintWeightOverrides,omitempty"`
	// Feasible reports whether the returned solution satisfies all hard constraints.
	Feasible *bool `json:"feasible,omitempty"`
	// ScoreString is the final score as a solver score string.
	ScoreString *string `json:"scoreString,omitempty"`
	// Score is the score broken down into its component levels.
	Score map[string]interface{} `json:"score,omitempty"`
	// UnassignedShifts are the shifts the solver could not staff.
	UnassignedShifts []ShiftAssignment `json:"unassignedShifts,omitempty"`
	// AssignedShifts are the shifts the solver successfully staffed.
	AssignedShifts []ShiftAssignment `json:"assignedShifts,omitempty"`
}
