package professionalservices

// ProfessionalServicesRequest is the request model for starting a professional services optimization.
type ProfessionalServicesRequest struct {
	// Name of the plan/project being scheduled.
	Name *string `json:"name,omitempty"`
	// Description is a free-text description of the plan.
	Description *string `json:"description,omitempty"`
	// StartDate is the start of the planning window, as an ISO-8601 date/timestamp.
	StartDate *string `json:"startDate,omitempty"`
	// EndDate is the end of the planning window, as an ISO-8601 date/timestamp.
	EndDate *string `json:"endDate,omitempty"`
	// Employees available to perform the tasks.
	Employees []Employee `json:"employees"`
	// Tasks to be assigned and scheduled.
	Tasks []Task `json:"tasks"`
	// Contracts referenced by employees, defining working-time limits.
	Contracts []Contract `json:"contracts,omitempty"`
}

// Employee represents an employee in the professional services optimization.
type Employee struct {
	// ID is the caller-supplied unique identifier for the employee.
	ID string `json:"id"`
	// Name is the employee display name.
	Name *string `json:"name,omitempty"`
	// Email is the employee email address.
	Email *string `json:"email,omitempty"`
	// Shifts are the solution shifts assigned to the employee.
	Shifts []Shift `json:"shifts"`
	// Skills the employee possesses.
	Skills []string `json:"skills"`
	// HourlyRate is the cost per hour, used by cost-minimization constraints.
	HourlyRate *float64 `json:"hourlyRate,omitempty"`
	// ContractID is the id of the Contract governing this employee's working-time limits.
	ContractID *string `json:"contractId,omitempty"`
	// TimeZoneID is the IANA time-zone id the employee works in.
	TimeZoneID *string `json:"timeZoneId,omitempty"`
	// DedicatedClientID, when set, restricts the employee to tasks for this client only.
	DedicatedClientID *string `json:"dedicatedClientId,omitempty"`
	// AvailabilityTimeSpans are the windows during which the employee is available (or explicitly unavailable).
	AvailabilityTimeSpans []AvailabilityTimeSpan `json:"availabilityTimeSpans,omitempty"`
}

// Shift represents a time window during which an employee is available.
type Shift struct {
	// ID is the caller-supplied unique identifier for the shift.
	ID string `json:"id"`
	// MinStartTime is the earliest the shift may start.
	MinStartTime string `json:"minStartTime"`
	// MaxEndTime is the latest the shift may end.
	MaxEndTime string `json:"maxEndTime"`
}

// Task represents a task to be scheduled in the professional services optimization.
type Task struct {
	// ID is the caller-supplied unique identifier for the task.
	ID string `json:"id"`
	// Name is the task display name.
	Name string `json:"name"`
	// Description is a free-text task description.
	Description *string `json:"description,omitempty"`
	// Deadline is the latest the task may be completed, as an ISO-8601 timestamp.
	Deadline *string `json:"deadline,omitempty"`
	// Duration is the effort required to complete the task, as an ISO-8601 duration.
	Duration string `json:"duration"`
	// Priority is the scheduling priority (e.g. LOW, MEDIUM, HIGH).
	Priority string `json:"priority"`
	// RequiredSkills are the skills an employee must have to be assigned the task.
	RequiredSkills []string `json:"requiredSkills"`
	// PreferredSkills are skills that are preferred but not mandatory.
	PreferredSkills []string `json:"preferredSkills,omitempty"`
	// ClientID is the identifier of the client the task belongs to.
	ClientID *string `json:"clientId,omitempty"`
	// ProjectID is the identifier of the project the task belongs to.
	ProjectID *string `json:"projectId,omitempty"`
	// TimeZoneID is the IANA time-zone id the task is scheduled in.
	TimeZoneID *string `json:"timeZoneId,omitempty"`
	// DependsOn are the ids of tasks that must be completed before this one can start.
	DependsOn []string `json:"dependsOn,omitempty"`
	// PreferredEmployees are the ids of employees preferred for this task.
	PreferredEmployees []string `json:"preferredEmployees,omitempty"`
	// ProhibitedEmployees are the ids of employees that must not be assigned this task.
	ProhibitedEmployees []string `json:"prohibitedEmployees,omitempty"`
}

// Contract represents an employee contract with working hour constraints.
type Contract struct {
	// ID is the caller-supplied unique identifier for the contract.
	ID *string `json:"id,omitempty"`
	// Name is the human-readable contract name.
	Name *string `json:"name,omitempty"`
	// MaxHoursPerDay is the maximum working hours per day, as an ISO-8601 duration.
	MaxHoursPerDay *string `json:"maxHoursPerDay,omitempty"`
	// MaxHoursPerWeek is the maximum working hours per week, as an ISO-8601 duration.
	MaxHoursPerWeek *string `json:"maxHoursPerWeek,omitempty"`
	// MinRestBetweenShifts is the minimum rest between consecutive shifts, as an ISO-8601 duration.
	MinRestBetweenShifts *string `json:"minRestBetweenShifts,omitempty"`
	// TargetUtilization is the desired utilization ratio (0-1) the solver tries to hit.
	TargetUtilization *float64 `json:"targetUtilization,omitempty"`
}

// AvailabilityTimeSpan represents a time span of employee availability.
type AvailabilityTimeSpan struct {
	// ID is the optional caller-supplied identifier for the span.
	ID *string `json:"id,omitempty"`
	// Start of the span, as an ISO-8601 timestamp.
	Start *string `json:"start,omitempty"`
	// End of the span, as an ISO-8601 timestamp.
	End *string `json:"end,omitempty"`
	// Type marks whether the span is availability or unavailability (e.g. AVAILABLE, UNAVAILABLE).
	Type *string `json:"type,omitempty"`
}

// ProfessionalServicesStartResponse is the response from starting a professional services optimization.
type ProfessionalServicesStartResponse struct {
	// JobID is the public PlanSolve job identifier - use it to poll status and fetch the solution.
	JobID string `json:"jobId"`
	// SolverJobID is the underlying solver engine's job identifier, when exposed.
	SolverJobID *string `json:"solverJobId,omitempty"`
	// Result is the inline solver result as raw JSON, when available synchronously.
	Result *string `json:"result,omitempty"`
	// Error is the error message when the solve request could not be accepted or run.
	Error *string `json:"error,omitempty"`
}

// ProfessionalServicesResultResponse is the response from getting professional
// services optimization results. It corresponds to the API's
// ProfessionalServicesResponse: a full echo of the request PLUS the result
// fields (solverStatus, feasible, scoreString, score, unassignedTasks,
// assignedTasks).
type ProfessionalServicesResultResponse struct {
	// JobID is the public job identifier this result was fetched with (stamped client-side).
	JobID *string `json:"jobId,omitempty"`
	// ID is the identifier of the solved plan.
	ID *string `json:"id,omitempty"`
	// Name is the echoed plan name.
	Name *string `json:"name,omitempty"`
	// Description is the echoed plan description.
	Description *string `json:"description,omitempty"`
	// StartDate is the echoed planning-window start.
	StartDate *string `json:"startDate,omitempty"`
	// EndDate is the echoed planning-window end.
	EndDate *string `json:"endDate,omitempty"`
	// Employees with their assigned shifts populated by the solver.
	Employees []Employee `json:"employees,omitempty"`
	// Tasks with their assignment details populated by the solver.
	Tasks []Task `json:"tasks,omitempty"`
	// Contracts are the echoed contracts.
	Contracts []Contract `json:"contracts,omitempty"`
	// SolverStatus is the raw solver engine status.
	SolverStatus *string `json:"solverStatus,omitempty"`
	// Feasible reports whether the returned solution satisfies all hard constraints.
	Feasible *bool `json:"feasible,omitempty"`
	// ScoreString is the final score as a solver score string.
	ScoreString *string `json:"scoreString,omitempty"`
	// Score is the score broken down into its component levels.
	Score map[string]interface{} `json:"score,omitempty"`
	// UnassignedTasks are the ids of tasks the solver could not assign.
	UnassignedTasks []string `json:"unassignedTasks,omitempty"`
	// AssignedTasks are the ids of tasks the solver successfully assigned.
	AssignedTasks []string `json:"assignedTasks,omitempty"`
}
