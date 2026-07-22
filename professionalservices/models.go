package professionalservices

// ProfessionalServicesRequest is the request model for starting a professional services optimization.
type ProfessionalServicesRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartDate   *string    `json:"startDate,omitempty"`
	EndDate     *string    `json:"endDate,omitempty"`
	Employees   []Employee `json:"employees"`
	Tasks       []Task     `json:"tasks"`
	Contracts   []Contract `json:"contracts,omitempty"`
}

// Employee represents an employee in the professional services optimization.
type Employee struct {
	ID                    string                 `json:"id"`
	Name                  *string                `json:"name,omitempty"`
	Email                 *string                `json:"email,omitempty"`
	Shifts                []Shift                `json:"shifts"`
	Skills                []string               `json:"skills"`
	HourlyRate            *float64               `json:"hourlyRate,omitempty"`
	ContractID            *string                `json:"contractId,omitempty"`
	TimeZoneID            *string                `json:"timeZoneId,omitempty"`
	DedicatedClientID     *string                `json:"dedicatedClientId,omitempty"`
	AvailabilityTimeSpans []AvailabilityTimeSpan `json:"availabilityTimeSpans,omitempty"`
}

// Shift represents a time window during which an employee is available.
type Shift struct {
	ID           string `json:"id"`
	MinStartTime string `json:"minStartTime"`
	MaxEndTime   string `json:"maxEndTime"`
}

// Task represents a task to be scheduled in the professional services optimization.
type Task struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Description        *string  `json:"description,omitempty"`
	Deadline           *string  `json:"deadline,omitempty"`
	Duration           string   `json:"duration"`
	Priority           string   `json:"priority"`
	RequiredSkills     []string `json:"requiredSkills"`
	PreferredSkills    []string `json:"preferredSkills,omitempty"`
	ClientID           *string  `json:"clientId,omitempty"`
	ProjectID          *string  `json:"projectId,omitempty"`
	TimeZoneID         *string  `json:"timeZoneId,omitempty"`
	DependsOn          []string `json:"dependsOn,omitempty"`
	PreferredEmployees []string `json:"preferredEmployees,omitempty"`
	ProhibitedEmployees []string `json:"prohibitedEmployees,omitempty"`
}

// Contract represents an employee contract with working hour constraints.
type Contract struct {
	ID                   *string  `json:"id,omitempty"`
	Name                 *string  `json:"name,omitempty"`
	MaxHoursPerDay       *string  `json:"maxHoursPerDay,omitempty"`
	MaxHoursPerWeek      *string  `json:"maxHoursPerWeek,omitempty"`
	MinRestBetweenShifts *string  `json:"minRestBetweenShifts,omitempty"`
	TargetUtilization    *float64 `json:"targetUtilization,omitempty"`
}

// AvailabilityTimeSpan represents a time span of employee availability.
type AvailabilityTimeSpan struct {
	ID    *string `json:"id,omitempty"`
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
	Type  *string `json:"type,omitempty"`
}

// ProfessionalServicesStartResponse is the response from starting a professional services optimization.
type ProfessionalServicesStartResponse struct {
	JobID       string  `json:"jobId"`
	SolverJobID *string `json:"solverJobId,omitempty"`
	Result      *string `json:"result,omitempty"`
	Error       *string `json:"error,omitempty"`
}

// ProfessionalServicesResultResponse is the response from getting professional
// services optimization results. It corresponds to the API's
// ProfessionalServicesResponse: a full echo of the request PLUS the result
// fields (solverStatus, feasible, scoreString, score, unassignedTasks,
// assignedTasks).
type ProfessionalServicesResultResponse struct {
	JobID           *string                `json:"jobId,omitempty"`
	ID              *string                `json:"id,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Description     *string                `json:"description,omitempty"`
	StartDate       *string                `json:"startDate,omitempty"`
	EndDate         *string                `json:"endDate,omitempty"`
	Employees       []Employee             `json:"employees,omitempty"`
	Tasks           []Task                 `json:"tasks,omitempty"`
	Contracts       []Contract             `json:"contracts,omitempty"`
	SolverStatus    *string                `json:"solverStatus,omitempty"`
	Feasible        *bool                  `json:"feasible,omitempty"`
	ScoreString     *string                `json:"scoreString,omitempty"`
	Score           map[string]interface{} `json:"score,omitempty"`
	UnassignedTasks []string               `json:"unassignedTasks,omitempty"`
	AssignedTasks   []string               `json:"assignedTasks,omitempty"`
}
