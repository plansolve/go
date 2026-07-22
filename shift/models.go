package shift

// ShiftRequest is the request model for starting a shift optimization.
type ShiftRequest struct {
	ID                      *string                  `json:"id,omitempty"`
	Name                    *string                  `json:"name,omitempty"`
	Description             *string                  `json:"description,omitempty"`
	Contracts               []Contract               `json:"contracts,omitempty"`
	Shifts                  []ShiftAssignment         `json:"shifts,omitempty"`
	Employees               []ShiftEmployee           `json:"employees,omitempty"`
	Options                 *Options                  `json:"options,omitempty"`
	Weights                 *Weights                  `json:"weights,omitempty"`
	Fairness                *Fairness                 `json:"fairness,omitempty"`
	DayOffRequests          []DayOffRequest           `json:"dayOffRequests,omitempty"`
	ShiftOffRequests        []ShiftOffRequest         `json:"shiftOffRequests,omitempty"`
	Hook                    *string                  `json:"hook,omitempty"`
	ConstraintWeightOverrides *ConstraintWeightOverrides `json:"constraintWeightOverrides,omitempty"`
}

// Shift represents a time window for shift scheduling.
type Shift struct {
	ID           string `json:"id"`
	MinStartTime string `json:"minStartTime"`
	MaxEndTime   string `json:"maxEndTime"`
}

// Employee represents an employee in shift scheduling.
type Employee struct {
	ID     string   `json:"id"`
	Shifts []Shift  `json:"shifts"`
	Skills []string `json:"skills"`
}

// Task represents a task in shift scheduling.
type Task struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Deadline       *string  `json:"deadline,omitempty"`
	Duration       string   `json:"duration"`
	Priority       string   `json:"priority"`
	RequiredSkills []string `json:"requiredSkills"`
}

// ShiftAssignment represents a shift assignment to be scheduled.
type ShiftAssignment struct {
	Name          *string  `json:"name,omitempty"`
	From          *string  `json:"from,omitempty"`
	To            *string  `json:"to,omitempty"`
	Skills        []string `json:"skills,omitempty"`
	Cost          float64  `json:"cost"`
	Value         int      `json:"value"`
	Priority      int      `json:"priority"`
	DesiredSkills []string `json:"desiredSkills,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	PinnedByUser  bool     `json:"pinnedByUser"`
}

// ShiftEmployee represents an employee in the shift optimization.
type ShiftEmployee struct {
	Name                 *string           `json:"name,omitempty"`
	Contract             *string           `json:"contract,omitempty"`
	Skills               []string          `json:"skills,omitempty"`
	LastRestDate         *string           `json:"lastRestDate,omitempty"`
	Availability         []string          `json:"availability,omitempty"`
	Preference           []string          `json:"preference,omitempty"`
	PeriodRules          []PeriodRule      `json:"periodRules,omitempty"`
	UnavailableDates     []string          `json:"unavailableDates,omitempty"`
	Tags                 []string          `json:"tags,omitempty"`
	MaximumMinutesPerWeek *int             `json:"maximumMinutesPerWeek,omitempty"`
	Shifts               []ShiftAssignment `json:"shifts,omitempty"`
}

// Contract represents a shift employee contract.
type Contract struct {
	Name                        *string `json:"name,omitempty"`
	Max                         *string `json:"max,omitempty"`
	Min                         *string `json:"min,omitempty"`
	MaxConsecutiveWorkDays      int     `json:"maxConsecutiveWorkDays"`
	MaxShiftsDay                int     `json:"maxShiftsDay"`
	MinRestBetweenShiftsSameDay *string `json:"minRestBetweenShiftsSameDay,omitempty"`
	MaxWorkingDays              int     `json:"maxWorkingDays"`
	LatestShiftStart            *string `json:"latestShiftStart,omitempty"`
	EarliestShiftStart          *string `json:"earliestShiftStart,omitempty"`
	MinimumConsecutiveDaysOff   int     `json:"minimumConsecutiveDaysOff"`
	MinimumHoursOffBetweenShifts int    `json:"minimumHoursOffBetweenShifts"`
}

// Weights represents the constraint weights for the shift optimizer.
type Weights struct {
	RequiredSkills            int `json:"requiredSkills"`
	ShiftCapacity             int `json:"shiftCapacity"`
	MinimumStaffing           int `json:"minimumStaffing"`
	NoDoubleBooking           int `json:"noDoubleBooking"`
	RestBetweenShifts         int `json:"restBetweenShifts"`
	EmployeeAvailability      int `json:"employeeAvailability"`
	ShiftPreferences          int `json:"shiftPreferences"`
	CostMinimization          int `json:"costMinimization"`
	WorkloadBalance           int `json:"workloadBalance"`
	Fairness                  int `json:"fairness"`
	MaxConsecutiveWorkDays    int `json:"maxConsecutiveWorkDays"`
	MaxShiftsPerDay           int `json:"maxShiftsPerDay"`
	MaxWorkingDaysPerWeek     int `json:"maxWorkingDaysPerWeek"`
	ContractRestBetweenShifts int `json:"contractRestBetweenShifts"`
	EarliestShiftStart        int `json:"earliestShiftStart"`
	LatestShiftStart          int `json:"latestShiftStart"`
	MinimumConsecutiveDaysOff int `json:"minimumConsecutiveDaysOff"`
	PeriodRuleViolation       int `json:"periodRuleViolation"`
	DesiredSkills             int `json:"desiredSkills"`
	DesiredDayOff             int `json:"desiredDayOff"`
	ShiftOffRequest           int `json:"shiftOffRequest"`
	BalanceTimeWorked         int `json:"balanceTimeWorked"`
	EmployeeAffinity          int `json:"employeeAffinity"`
	AvoidShiftCloseToDayOff   int `json:"avoidShiftCloseToDayOff"`
}

// Options represents solver options for the shift optimizer.
type Options struct {
	PartialPlanning *bool `json:"partialPlanning,omitempty"`
	MaxIterations   *int  `json:"maxIterations,omitempty"`
	TimeLimit       *int  `json:"timeLimit,omitempty"`
}

// Fairness represents fairness configuration for the shift optimizer.
type Fairness struct {
	FairnessBuckets []FairnessBucket `json:"fairnessBuckets,omitempty"`
}

// FairnessBucket represents a fairness bucket grouping employees and shifts.
type FairnessBucket struct {
	Name      *string  `json:"name,omitempty"`
	Employees []string `json:"employees,omitempty"`
	Shifts    []string `json:"shifts,omitempty"`
	Period    *string  `json:"period,omitempty"`
}

// DayOffRequest represents an employee's request for a day off.
type DayOffRequest struct {
	ID           *string `json:"id,omitempty"`
	EmployeeName *string `json:"employeeName,omitempty"`
	Date         *string `json:"date,omitempty"`
	Weight       int     `json:"weight"`
}

// ShiftOffRequest represents an employee's request to be off a specific shift.
type ShiftOffRequest struct {
	ID           *string `json:"id,omitempty"`
	EmployeeName *string `json:"employeeName,omitempty"`
	ShiftName    *string `json:"shiftName,omitempty"`
	Weight       int     `json:"weight"`
}

// PeriodRule represents a period-based working constraint.
type PeriodRule struct {
	Period                              *PlanningPeriod `json:"period,omitempty"`
	MaxWorkingDays                      int             `json:"maxWorkingDays"`
	MinWorkingDays                      int             `json:"minWorkingDays"`
	MinWorkingDuration                  *string         `json:"minWorkingDuration,omitempty"`
	MaxWorkingDuration                  *string         `json:"maxWorkingDuration,omitempty"`
	MinRestDurationBetweenShiftsSameDay *string         `json:"minRestDurationBetweenShiftsSameDay,omitempty"`
	MinRestDuration                     *string         `json:"minRestDuration,omitempty"`
}

// PlanningPeriod represents a time period for planning rules.
type PlanningPeriod struct {
	From *string `json:"from,omitempty"`
	To   *string `json:"to,omitempty"`
}

// ConstraintWeightOverrides represents overrides for constraint weights.
type ConstraintWeightOverrides struct {
	KnownConstraintNames []string `json:"knownConstraintNames,omitempty"`
}

// ShiftStartResponse is the response from starting a shift optimization.
type ShiftStartResponse struct {
	JobID       string  `json:"jobId"`
	SolverJobID *string `json:"solverJobId,omitempty"`
	Result      *string `json:"result,omitempty"`
	Error       *string `json:"error,omitempty"`
}

// ShiftResultResponse is the response from getting shift optimization results.
// It corresponds to the API's ShiftAssignmentResponse: a full echo of the
// request PLUS the result fields (feasible, scoreString, score,
// unassignedShifts, assignedShifts). It carries no jobId.
type ShiftResultResponse struct {
	JobID                     *string                    `json:"jobId,omitempty"`
	ID                        *string                    `json:"id,omitempty"`
	Name                      *string                    `json:"name,omitempty"`
	Description               *string                    `json:"description,omitempty"`
	Contracts                 []Contract                 `json:"contracts,omitempty"`
	Shifts                    []ShiftAssignment          `json:"shifts,omitempty"`
	Employees                 []ShiftEmployee            `json:"employees,omitempty"`
	DayOffRequests            []DayOffRequest            `json:"dayOffRequests,omitempty"`
	ShiftOffRequests          []ShiftOffRequest          `json:"shiftOffRequests,omitempty"`
	Options                   *Options                   `json:"options,omitempty"`
	Weights                   *Weights                   `json:"weights,omitempty"`
	Fairness                  *Fairness                  `json:"fairness,omitempty"`
	Hook                      *string                    `json:"hook,omitempty"`
	ConstraintWeightOverrides *ConstraintWeightOverrides `json:"constraintWeightOverrides,omitempty"`
	Feasible                  *bool                      `json:"feasible,omitempty"`
	ScoreString               *string                    `json:"scoreString,omitempty"`
	Score                     map[string]interface{}     `json:"score,omitempty"`
	UnassignedShifts          []ShiftAssignment          `json:"unassignedShifts,omitempty"`
	AssignedShifts            []ShiftAssignment          `json:"assignedShifts,omitempty"`
}
