package shift

// ShiftRequest is the request model for starting a shift optimization.
type ShiftRequest struct {
	// ID is the optional caller-supplied identifier for the plan.
	ID *string `json:"id,omitempty"`
	// Name is the plan display name.
	Name *string `json:"name,omitempty"`
	// Description is the free-text plan description.
	Description *string `json:"description,omitempty"`
	// Contracts define the working-time rules referenced by employees.
	Contracts []Contract `json:"contracts,omitempty"`
	// Shifts that need to be staffed.
	Shifts []ShiftAssignment `json:"shifts,omitempty"`
	// Employees available to be assigned to shifts.
	Employees []ShiftEmployee `json:"employees,omitempty"`
	// Options holds optional controls governing how the shift solver runs.
	Options *Options `json:"options,omitempty"`
	// Weights are the per-constraint penalty weights that steer the solver's objective.
	Weights *Weights `json:"weights,omitempty"`
	// Fairness holds the fairness configuration.
	Fairness *Fairness `json:"fairness,omitempty"`
	// DayOffRequests are employee requests to have specific days off.
	DayOffRequests []DayOffRequest `json:"dayOffRequests,omitempty"`
	// ShiftOffRequests are employee requests to avoid specific shifts.
	ShiftOffRequests []ShiftOffRequest `json:"shiftOffRequests,omitempty"`
	// Hook is an optional callback URL invoked when solving completes.
	Hook *string `json:"hook,omitempty"`
	// ConstraintWeightOverrides holds advanced knobs for the shift constraint set.
	ConstraintWeightOverrides *ConstraintWeightOverrides `json:"constraintWeightOverrides,omitempty"`
}

// Shift represents a time window for shift scheduling.
type Shift struct {
	// ID is the caller-supplied unique identifier for the shift.
	ID string `json:"id"`
	// MinStartTime is the earliest the shift may start.
	MinStartTime string `json:"minStartTime"`
	// MaxEndTime is the latest the shift may end.
	MaxEndTime string `json:"maxEndTime"`
}

// Employee represents an employee in shift scheduling.
type Employee struct {
	// ID is the caller-supplied unique identifier for the employee.
	ID string `json:"id"`
	// Shifts assigned to or available for the employee.
	Shifts []Shift `json:"shifts"`
	// Skills the employee possesses.
	Skills []string `json:"skills"`
}

// Task represents a task in shift scheduling.
type Task struct {
	// ID is the caller-supplied unique identifier for the task.
	ID string `json:"id"`
	// Name is the task display name.
	Name string `json:"name"`
	// Deadline is the latest the task may be completed, as an ISO-8601 timestamp.
	Deadline *string `json:"deadline,omitempty"`
	// Duration is the effort required to complete the task, as an ISO-8601 duration.
	Duration string `json:"duration"`
	// Priority is the scheduling priority (e.g. LOW, MEDIUM, HIGH).
	Priority string `json:"priority"`
	// RequiredSkills are the skills an employee must have to be assigned the task.
	RequiredSkills []string `json:"requiredSkills"`
}

// ShiftAssignment represents a shift assignment to be scheduled.
type ShiftAssignment struct {
	// Name is the shift name, used as its identifier.
	Name *string `json:"name,omitempty"`
	// From is the shift start, as an ISO-8601 timestamp.
	From *string `json:"from,omitempty"`
	// To is the shift end, as an ISO-8601 timestamp.
	To *string `json:"to,omitempty"`
	// Skills an employee must have to be assigned this shift.
	Skills []string `json:"skills,omitempty"`
	// Cost of staffing this shift.
	Cost float64 `json:"cost"`
	// Value is the business value of staffing this shift.
	Value int `json:"value"`
	// Priority is the relative priority of staffing this shift.
	Priority int `json:"priority"`
	// DesiredSkills are skills that are preferred but not required.
	DesiredSkills []string `json:"desiredSkills,omitempty"`
	// Tags are arbitrary tags used for grouping and affinity constraints.
	Tags []string `json:"tags,omitempty"`
	// PinnedByUser, when true, locks the existing assignment and leaves it unchanged by the solver.
	PinnedByUser bool `json:"pinnedByUser"`
}

// ShiftEmployee represents an employee in the shift optimization.
type ShiftEmployee struct {
	// Name is the employee name, used as the identifier referenced by requests and assignments.
	Name *string `json:"name,omitempty"`
	// Contract is the name of the Contract governing this employee's working-time rules.
	Contract *string `json:"contract,omitempty"`
	// Skills the employee possesses.
	Skills []string `json:"skills,omitempty"`
	// LastRestDate is the date the employee last had rest.
	LastRestDate *string `json:"lastRestDate,omitempty"`
	// Availability lists the dates/periods the employee is available to work.
	Availability []string `json:"availability,omitempty"`
	// Preference lists the shifts/periods the employee prefers to work.
	Preference []string `json:"preference,omitempty"`
	// PeriodRules are the period-specific working-time rules for this employee.
	PeriodRules []PeriodRule `json:"periodRules,omitempty"`
	// UnavailableDates lists the dates the employee cannot work.
	UnavailableDates []string `json:"unavailableDates,omitempty"`
	// Tags are arbitrary tags used for grouping and affinity constraints.
	Tags []string `json:"tags,omitempty"`
	// MaximumMinutesPerWeek is the cap on weekly working minutes.
	MaximumMinutesPerWeek *int `json:"maximumMinutesPerWeek,omitempty"`
	// Shifts are the solution shifts assigned to this employee.
	Shifts []ShiftAssignment `json:"shifts,omitempty"`
}

// Contract represents a shift employee contract.
type Contract struct {
	// Name is the contract name, referenced by an employee's contract field.
	Name *string `json:"name,omitempty"`
	// Max is the maximum total working time, as an ISO-8601 duration.
	Max *string `json:"max,omitempty"`
	// Min is the minimum total working time, as an ISO-8601 duration.
	Min *string `json:"min,omitempty"`
	// MaxConsecutiveWorkDays is the maximum consecutive days that may be worked.
	MaxConsecutiveWorkDays int `json:"maxConsecutiveWorkDays"`
	// MaxShiftsDay is the maximum shifts allowed on a single day.
	MaxShiftsDay int `json:"maxShiftsDay"`
	// MinRestBetweenShiftsSameDay is the minimum rest between two shifts on the same day.
	MinRestBetweenShiftsSameDay *string `json:"minRestBetweenShiftsSameDay,omitempty"`
	// MaxWorkingDays is the maximum working days in the planning period.
	MaxWorkingDays int `json:"maxWorkingDays"`
	// LatestShiftStart is the latest permitted shift start time.
	LatestShiftStart *string `json:"latestShiftStart,omitempty"`
	// EarliestShiftStart is the earliest permitted shift start time.
	EarliestShiftStart *string `json:"earliestShiftStart,omitempty"`
	// MinimumConsecutiveDaysOff is the minimum consecutive days off between working stretches.
	MinimumConsecutiveDaysOff int `json:"minimumConsecutiveDaysOff"`
	// MinimumHoursOffBetweenShifts is the minimum hours off between consecutive shifts.
	MinimumHoursOffBetweenShifts int `json:"minimumHoursOffBetweenShifts"`
}

// Weights represents the constraint weights for the shift optimizer.
type Weights struct {
	// RequiredSkills is the penalty for assigning an employee who lacks a shift's required skills.
	RequiredSkills int `json:"requiredSkills"`
	// ShiftCapacity is the penalty for exceeding a shift's staffing capacity.
	ShiftCapacity int `json:"shiftCapacity"`
	// MinimumStaffing is the penalty for staffing a shift below its minimum required headcount.
	MinimumStaffing int `json:"minimumStaffing"`
	// NoDoubleBooking is the penalty for assigning an employee to overlapping shifts.
	NoDoubleBooking int `json:"noDoubleBooking"`
	// RestBetweenShifts is the penalty for violating the required rest between consecutive shifts.
	RestBetweenShifts int `json:"restBetweenShifts"`
	// EmployeeAvailability is the penalty for assigning an employee outside their declared availability.
	EmployeeAvailability int `json:"employeeAvailability"`
	// ShiftPreferences is the reward/penalty weight for honoring employees' shift preferences.
	ShiftPreferences int `json:"shiftPreferences"`
	// CostMinimization is the weight for minimizing total staffing cost.
	CostMinimization int `json:"costMinimization"`
	// WorkloadBalance is the weight for balancing workload evenly across employees.
	WorkloadBalance int `json:"workloadBalance"`
	// Fairness is the weight for the fairness objective across fairness buckets.
	Fairness int `json:"fairness"`
	// MaxConsecutiveWorkDays is the penalty for exceeding the maximum consecutive working days.
	MaxConsecutiveWorkDays int `json:"maxConsecutiveWorkDays"`
	// MaxShiftsPerDay is the penalty for exceeding the maximum shifts per day.
	MaxShiftsPerDay int `json:"maxShiftsPerDay"`
	// MaxWorkingDaysPerWeek is the penalty for exceeding the maximum working days per week.
	MaxWorkingDaysPerWeek int `json:"maxWorkingDaysPerWeek"`
	// ContractRestBetweenShifts is the penalty for violating the contract's minimum rest between shifts.
	ContractRestBetweenShifts int `json:"contractRestBetweenShifts"`
	// EarliestShiftStart is the penalty for starting a shift before the contract's earliest allowed start.
	EarliestShiftStart int `json:"earliestShiftStart"`
	// LatestShiftStart is the penalty for starting a shift after the contract's latest allowed start.
	LatestShiftStart int `json:"latestShiftStart"`
	// MinimumConsecutiveDaysOff is the penalty for providing fewer than the minimum consecutive days off.
	MinimumConsecutiveDaysOff int `json:"minimumConsecutiveDaysOff"`
	// PeriodRuleViolation is the penalty for violating an employee's period-specific rule.
	PeriodRuleViolation int `json:"periodRuleViolation"`
	// DesiredSkills is the weight for preferring employees who hold a shift's desired skills.
	DesiredSkills int `json:"desiredSkills"`
	// DesiredDayOff is the weight for honoring employees' desired day-off requests.
	DesiredDayOff int `json:"desiredDayOff"`
	// ShiftOffRequest is the weight for honoring employees' shift-off requests.
	ShiftOffRequest int `json:"shiftOffRequest"`
	// BalanceTimeWorked is the weight for balancing total time worked across employees.
	BalanceTimeWorked int `json:"balanceTimeWorked"`
	// EmployeeAffinity is the weight for keeping tagged employees working together.
	EmployeeAffinity int `json:"employeeAffinity"`
	// AvoidShiftCloseToDayOff is the weight for avoiding shifts scheduled immediately adjacent to a day off.
	AvoidShiftCloseToDayOff int `json:"avoidShiftCloseToDayOff"`
}

// Options represents solver options for the shift optimizer.
type Options struct {
	// PartialPlanning, when true, keeps existing assignments and only fills the gaps.
	PartialPlanning *bool `json:"partialPlanning,omitempty"`
	// MaxIterations is the maximum number of solver iterations before stopping.
	MaxIterations *int `json:"maxIterations,omitempty"`
	// TimeLimit is the maximum solve time in seconds.
	TimeLimit *int `json:"timeLimit,omitempty"`
}

// Fairness represents fairness configuration for the shift optimizer.
type Fairness struct {
	// FairnessBuckets define the groups the fairness objective balances over.
	FairnessBuckets []FairnessBucket `json:"fairnessBuckets,omitempty"`
}

// FairnessBucket represents a fairness bucket grouping employees and shifts.
type FairnessBucket struct {
	// Name is the name identifier for the fairness bucket.
	Name *string `json:"name,omitempty"`
	// Employees is the list of employee names in this bucket.
	Employees []string `json:"employees,omitempty"`
	// Shifts is the list of shift names in this bucket.
	Shifts []string `json:"shifts,omitempty"`
	// Period is the time period for this bucket (e.g. 2024-01-01/2024-01-08).
	Period *string `json:"period,omitempty"`
}

// DayOffRequest represents an employee's request for a day off.
type DayOffRequest struct {
	// ID is the optional caller-supplied identifier for the request.
	ID *string `json:"id,omitempty"`
	// EmployeeName is the name of the employee making the request.
	EmployeeName *string `json:"employeeName,omitempty"`
	// Date is the requested day off, as an ISO-8601 date.
	Date *string `json:"date,omitempty"`
	// Weight is the strength of the preference; higher values make honoring it more important.
	Weight int `json:"weight"`
}

// ShiftOffRequest represents an employee's request to be off a specific shift.
type ShiftOffRequest struct {
	// ID is the optional caller-supplied identifier.
	ID *string `json:"id,omitempty"`
	// EmployeeName is the name of the employee making the request.
	EmployeeName *string `json:"employeeName,omitempty"`
	// ShiftName is the name of the shift the employee wants to avoid.
	ShiftName *string `json:"shiftName,omitempty"`
	// Weight is the strength of the preference; higher values make honoring it more important.
	Weight int `json:"weight"`
}

// PeriodRule represents a period-based working constraint.
type PeriodRule struct {
	// Period is the planning period this rule applies to.
	Period *PlanningPeriod `json:"period,omitempty"`
	// MaxWorkingDays is the maximum working days allowed within the period.
	MaxWorkingDays int `json:"maxWorkingDays"`
	// MinWorkingDays is the minimum working days required within the period.
	MinWorkingDays int `json:"minWorkingDays"`
	// MinWorkingDuration is the minimum total working time within the period.
	MinWorkingDuration *string `json:"minWorkingDuration,omitempty"`
	// MaxWorkingDuration is the maximum total working time within the period.
	MaxWorkingDuration *string `json:"maxWorkingDuration,omitempty"`
	// MinRestDurationBetweenShiftsSameDay is the minimum rest between shifts on the same day.
	MinRestDurationBetweenShiftsSameDay *string `json:"minRestDurationBetweenShiftsSameDay,omitempty"`
	// MinRestDuration is the minimum rest between consecutive shifts.
	MinRestDuration *string `json:"minRestDuration,omitempty"`
}

// PlanningPeriod represents a time period for planning rules.
type PlanningPeriod struct {
	// From is the start of the period.
	From *string `json:"from,omitempty"`
	// To is the end of the period.
	To *string `json:"to,omitempty"`
}

// ConstraintWeightOverrides represents overrides for constraint weights.
type ConstraintWeightOverrides struct {
	// KnownConstraintNames are the names of constraints the solver is expected to recognize and apply.
	KnownConstraintNames []string `json:"knownConstraintNames,omitempty"`
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
// It corresponds to the API's ShiftAssignmentResponse: a full echo of the
// request PLUS the result fields (feasible, scoreString, score,
// unassignedShifts, assignedShifts). It carries no jobId.
type ShiftResultResponse struct {
	// JobID is the public job identifier this result was fetched with (stamped client-side).
	JobID *string `json:"jobId,omitempty"`
	// ID is the echoed plan identifier.
	ID *string `json:"id,omitempty"`
	// Name is the echoed plan display name.
	Name *string `json:"name,omitempty"`
	// Description is the echoed free-text plan description.
	Description *string `json:"description,omitempty"`
	// Contracts are the echoed contracts defining the working-time rules referenced by employees.
	Contracts []Contract `json:"contracts,omitempty"`
	// Shifts are the echoed shifts that need to be staffed.
	Shifts []ShiftAssignment `json:"shifts,omitempty"`
	// Employees are the employees with their assigned shifts populated by the solver.
	Employees []ShiftEmployee `json:"employees,omitempty"`
	// DayOffRequests are the echoed employee requests to have specific days off.
	DayOffRequests []DayOffRequest `json:"dayOffRequests,omitempty"`
	// ShiftOffRequests are the echoed employee requests to avoid specific shifts.
	ShiftOffRequests []ShiftOffRequest `json:"shiftOffRequests,omitempty"`
	// Options are the echoed controls governing how the shift solver runs.
	Options *Options `json:"options,omitempty"`
	// Weights are the echoed per-constraint penalty weights that steer the solver's objective.
	Weights *Weights `json:"weights,omitempty"`
	// Fairness is the echoed fairness configuration.
	Fairness *Fairness `json:"fairness,omitempty"`
	// Hook is the echoed callback URL invoked when solving completes.
	Hook *string `json:"hook,omitempty"`
	// ConstraintWeightOverrides are the echoed advanced knobs for the shift constraint set.
	ConstraintWeightOverrides *ConstraintWeightOverrides `json:"constraintWeightOverrides,omitempty"`
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
