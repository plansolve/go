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
	// TaskTypeTransitions is the setup cost or prohibition for switching between kinds of work on one employee's list.
	TaskTypeTransitions []TaskTypeTransition `json:"taskTypeTransitions,omitempty"`
	// TaskTypes is an optional vocabulary of task types, so wildcard transitions resolve to types no task carries yet.
	TaskTypes []string `json:"taskTypes,omitempty"`
	// FreezeUntil is the plan-wide freeze horizon, as an ISO-8601 timestamp: work already planned to start before it is pinned.
	FreezeUntil *string `json:"freezeUntil,omitempty"`
	// Weights are per-constraint weight overrides keyed by the solver's constraint identifier, with a
	// "Xhard/Ymedium/Zsoft" value. An identifier the solver does not recognise fails the job.
	Weights map[string]string `json:"weights,omitempty"`
	// Options holds optional termination controls for the solve.
	Options *SolverOptions `json:"options,omitempty"`
}

// SolverOptions holds optional termination controls for a Professional Services
// solve, expressed as ISO-8601 durations. It corresponds to the API's
// ProfessionalServicesSolverOptions.
type SolverOptions struct {
	// SpentLimit is the maximum total time the solver may run (e.g. PT30S), as an ISO-8601 duration.
	SpentLimit *string `json:"spentLimit,omitempty"`
	// UnimprovedSpentLimit stops the solver early after this much time passes with no score improvement (e.g. PT5S).
	UnimprovedSpentLimit *string `json:"unimprovedSpentLimit,omitempty"`
}

// Employee represents an employee in the professional services optimization.
type Employee struct {
	// ID is the caller-supplied unique identifier for the employee.
	ID string `json:"id"`
	// Name is the employee display name.
	Name *string `json:"name,omitempty"`
	// Email is the employee email address.
	Email *string `json:"email,omitempty"`
	// Shifts are the working windows during which tasks may be scheduled.
	Shifts []Shift `json:"shifts"`
	// Skills the employee possesses.
	Skills []string `json:"skills"`
	// SkillValidity holds validity windows for skills not held indefinitely; checked against each task's start time.
	SkillValidity []SkillValidity `json:"skillValidity,omitempty"`
	// HourlyRate is the cost per hour, used by cost-minimization constraints.
	HourlyRate *float64 `json:"hourlyRate,omitempty"`
	// ContractID is the id of the Contract governing this employee's working-time limits.
	ContractID *string `json:"contractId,omitempty"`
	// FreezeUntil is a freeze horizon for this employee alone: work already planned to start before it is pinned.
	// Combines with the plan-level freeze horizon, whichever reaches further.
	FreezeUntil *string `json:"freezeUntil,omitempty"`
	// AvailabilityTimeSpans are the windows during which the employee is explicitly unavailable (type UNAVAILABLE);
	// AVAILABLE is the inert default.
	AvailabilityTimeSpans []AvailabilityTimeSpan `json:"availabilityTimeSpans,omitempty"`
	// Tasks (result only) are the ids of the tasks the solver assigned to this employee, in schedule order.
	Tasks []string `json:"tasks,omitempty"`
	// Duration (result only) is the total assigned work, as an ISO-8601 duration.
	Duration *string `json:"duration,omitempty"`
}

// SkillValidity is a validity window for a skill an employee does not hold indefinitely (a certification, say).
// A skill named here counts even if absent from the employee's flat skills list, and several entries for one
// skill are alternatives.
type SkillValidity struct {
	// Skill the window applies to.
	Skill string `json:"skill"`
	// ValidFrom is the start of the window, as an ISO-8601 timestamp; open-ended when absent.
	ValidFrom *string `json:"validFrom,omitempty"`
	// ValidTo is the end of the window (exclusive), as an ISO-8601 timestamp; open-ended when absent.
	ValidTo *string `json:"validTo,omitempty"`
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
	// Deadline is the latest the task may be completed (soft), as an ISO-8601 timestamp.
	Deadline *string `json:"deadline,omitempty"`
	// SLA is the hard deadline that must not be violated (unlike Deadline, which is soft), as an ISO-8601 timestamp.
	SLA *string `json:"sla,omitempty"`
	// EarliestStart is the moment before which work may not begin (client kickoff, contract start), as an ISO-8601 timestamp.
	EarliestStart *string `json:"earliestStart,omitempty"`
	// Duration is the effort required to complete the task, as an ISO-8601 duration.
	Duration string `json:"duration"`
	// DurationByEmployee holds per-employee duration overrides keyed by employee id, falling back to Duration.
	DurationByEmployee map[string]string `json:"durationByEmployee,omitempty"`
	// Priority is the scheduling priority (e.g. LOW, MEDIUM, HIGH).
	Priority string `json:"priority"`
	// PriorityValue is the numeric urgency, higher being more important; it overrides Priority (worth HIGH 100,
	// MEDIUM 50, LOW 25) to rank tasks inside a tier.
	PriorityValue *int `json:"priorityValue,omitempty"`
	// TaskType is the kind of work (deep work, meeting, admin), driving the plan's task-type transitions.
	TaskType *string `json:"taskType,omitempty"`
	// RequiredSkills are the skills an employee must have to be assigned the task.
	RequiredSkills []string `json:"requiredSkills"`
	// PreferredSkills are skills that are preferred but not mandatory.
	PreferredSkills []string `json:"preferredSkills,omitempty"`
	// ClientID is the identifier of the client the task belongs to.
	ClientID *string `json:"clientId,omitempty"`
	// ProjectID is the identifier of the project the task belongs to.
	ProjectID *string `json:"projectId,omitempty"`
	// DependsOn are the predecessor tasks. A plain finish-to-start dependency is TaskDependency{TaskID: "..."}.
	DependsOn []TaskDependency `json:"dependsOn,omitempty"`
	// PreferredEmployees are the ids of employees preferred for this task.
	PreferredEmployees []string `json:"preferredEmployees,omitempty"`
	// ProhibitedEmployees are the ids of employees that must not be assigned this task.
	ProhibitedEmployees []string `json:"prohibitedEmployees,omitempty"`
	// AllowedEmployees, when set, are the only employee ids that may take the task (hard whitelist).
	AllowedEmployees []string `json:"allowedEmployees,omitempty"`
	// Pinned keeps the task at InitialEmployeeID/InitialStartTime; both are required when pinning.
	Pinned bool `json:"pinned,omitempty"`
	// ParentTaskID is the id of the parent task when this task is one segment of a split task.
	ParentTaskID *string `json:"parentTaskId,omitempty"`
	// SegmentIndex is the zero-based order of this segment within its split task.
	SegmentIndex *int `json:"segmentIndex,omitempty"`
	// TotalSegments is the number of segments the split task has in total.
	TotalSegments *int `json:"totalSegments,omitempty"`
	// InitialEmployeeID is the employee this task was assigned to before re-planning; the solver prefers to keep it.
	InitialEmployeeID *string `json:"initialEmployeeId,omitempty"`
	// InitialStartTime is the task's start time before re-planning, as an ISO-8601 timestamp.
	InitialStartTime *string `json:"initialStartTime,omitempty"`
	// Employee (result only) is the id of the employee the solver assigned the task to, or nil when unassigned.
	Employee *string `json:"employee,omitempty"`
	// PreviousTask (result only) is the id of the task scheduled immediately before this one on the same employee.
	PreviousTask *string `json:"previousTask,omitempty"`
	// StartTime (result only) is the scheduled start, as an ISO-8601 timestamp.
	StartTime *string `json:"startTime,omitempty"`
	// EndTime (result only) is the scheduled end, as an ISO-8601 timestamp.
	EndTime *string `json:"endTime,omitempty"`
}

// TaskDependency is a predecessor relation on a Task. A plain finish-to-start dependency only needs TaskID.
type TaskDependency struct {
	// TaskID is the id of the task that must precede this one.
	TaskID string `json:"taskId"`
	// Anchor is which end of the predecessor the offsets are measured from: START or END (the solver default).
	Anchor *string `json:"anchor,omitempty"`
	// MinOffset is the minimum lag after the anchor before this task may start, as an ISO-8601 duration.
	MinOffset *string `json:"minOffset,omitempty"`
	// MaxOffset is the maximum gap after the anchor within which this task must start, as an ISO-8601 duration.
	MaxOffset *string `json:"maxOffset,omitempty"`
	// RequiresSameEmployee sends both tasks to the same employee. Requires an END anchor.
	RequiresSameEmployee bool `json:"requiresSameEmployee,omitempty"`
	// NoIntermediateTasks makes the two tasks run back to back. Requires RequiresSameEmployee.
	NoIntermediateTasks bool `json:"noIntermediateTasks,omitempty"`
}

// TaskTypeTransition is the cost of switching between two kinds of work on one employee's list. "*" matches any
// type; the solver resolves wildcards to concrete pairs and lets a specific rule beat a broad one.
type TaskTypeTransition struct {
	// FromTaskType is the task type being switched from, or "*" for any.
	FromTaskType string `json:"fromTaskType"`
	// ToTaskType is the task type being switched to, or "*" for any.
	ToTaskType string `json:"toTaskType"`
	// SetupDuration is what the switch costs, as an ISO-8601 duration, penalised in the soft tier.
	SetupDuration *string `json:"setupDuration,omitempty"`
	// Forbidden means the switch may not happen at all (hard constraint).
	Forbidden bool `json:"forbidden,omitempty"`
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
}

// AvailabilityTimeSpan represents a time span of employee availability.
type AvailabilityTimeSpan struct {
	// ID is the optional caller-supplied identifier for the span.
	ID *string `json:"id,omitempty"`
	// Start of the span, as an ISO-8601 timestamp.
	Start *string `json:"start,omitempty"`
	// End of the span, as an ISO-8601 timestamp.
	End *string `json:"end,omitempty"`
	// Type is AVAILABLE or UNAVAILABLE. Only UNAVAILABLE has an effect; AVAILABLE is the inert default.
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
	// Employees with their assigned tasks populated by the solver.
	Employees []Employee `json:"employees,omitempty"`
	// Tasks with their assignment details populated by the solver.
	Tasks []Task `json:"tasks,omitempty"`
	// Contracts are the echoed contracts.
	Contracts []Contract `json:"contracts,omitempty"`
	// TaskTypeTransitions are the echoed work-type transitions, with wildcards resolved to concrete pairs.
	TaskTypeTransitions []TaskTypeTransition `json:"taskTypeTransitions,omitempty"`
	// TaskTypes is the echoed task-type vocabulary.
	TaskTypes []string `json:"taskTypes,omitempty"`
	// FreezeUntil is the echoed plan-wide freeze horizon.
	FreezeUntil *string `json:"freezeUntil,omitempty"`
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
