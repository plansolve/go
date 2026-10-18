package jobs

// JobStatus is the lifecycle state of a solver run as recorded in the Jobs table.
type JobStatus int

const (
	// JobStatusPending indicates the job has not yet reached a terminal state.
	JobStatusPending JobStatus = 0
	// JobStatusSucceeded indicates the job completed successfully.
	JobStatusSucceeded JobStatus = 1
	// JobStatusFailed indicates the job failed.
	JobStatusFailed JobStatus = 2
	// JobStatusLost indicates the job's outcome could not be determined.
	JobStatusLost JobStatus = 3
	// JobStatusCancelled indicates the job was stopped while still waiting in the
	// queue, before solving started; no plan was produced and nothing was billed.
	JobStatusCancelled JobStatus = 4
)

// Solver families reported in JobDto.SolverType, untyped so the int field stays unchanged.
const (
	// SolverTypeFieldService is a Field Service (vehicle routing) job.
	SolverTypeFieldService = 0
	// SolverTypeShift is a Shift scheduling job.
	SolverTypeShift = 1
	// SolverTypeProfessionalServices is a Professional Services (task assignment) job.
	SolverTypeProfessionalServices = 2
)

// CreateJobRequest is the payload for creating a new solver job via the public Jobs API.
// The job is billed against the tenant's subscription for the named solver.
type CreateJobRequest struct {
	// SolverID identifies the solver/model to run (e.g. field_service,
	// professional_services, shift).
	SolverID string `json:"solverId"`
	// Request is the optimization problem definition, as a raw JSON string matching
	// the chosen solver's request schema.
	Request string `json:"request"`
}

// JobDto is a solver job as exposed by the public Jobs API.
type JobDto struct {
	// ID is the internal numeric primary key. Not used to address the job over the
	// API - use JobID.
	ID int64 `json:"id"`
	// JobID is the public, opaque job identifier used in API routes.
	JobID *string `json:"jobId,omitempty"`
	// TenantID identifies the tenant that owns the job.
	TenantID *string `json:"tenantId,omitempty"`
	// SubscriptionID is the subscription the job was charged against.
	SubscriptionID int `json:"subscriptionId"`
	// SolverType is the solver family that produced the job (0 = Field Service,
	// 1 = Shift, 2 = Professional Services); compare with the SolverType* constants.
	SolverType int `json:"solverType"`
	// Status is the lifecycle state of the job.
	Status JobStatus `json:"status"`
	// CreatedAt is the UTC timestamp the job was created, as an ISO-8601 string.
	CreatedAt string `json:"createdAt"`
	// CompletedAt is the UTC timestamp the job reached a terminal state, or nil
	// while still in flight.
	CompletedAt *string `json:"completedAt,omitempty"`
	// Request is the optimization problem definition that was submitted, as raw JSON.
	Request *string `json:"request,omitempty"`
	// Response is the solver's result as raw JSON, or nil until the job completes.
	Response *string `json:"response,omitempty"`
	// HasResponse reports whether a solver response is available for the job.
	HasResponse bool `json:"hasResponse"`
	// BilledItems is the number of schedulable items billed for this job.
	BilledItems int `json:"billedItems"`
	// EstimatedItems is the number of schedulable items estimated for the job when
	// it was submitted, or nil when no estimate was recorded.
	EstimatedItems *int `json:"estimatedItems,omitempty"`
	// UserID identifies the user who submitted the job, or nil for API-key/system runs.
	UserID *int `json:"userId,omitempty"`
	// UserName is the display name of the user who submitted the job.
	UserName *string `json:"userName,omitempty"`
}

// PagedResultOfJobDto is a single page of jobs plus the pagination metadata needed
// to navigate the set.
type PagedResultOfJobDto struct {
	// Items are the jobs on the current page.
	Items []JobDto `json:"items,omitempty"`
	// TotalCount is the total number of items across all pages.
	TotalCount int `json:"totalCount"`
	// Page is the 1-based index of the current page.
	Page int `json:"page"`
	// PageSize is the maximum number of items per page.
	PageSize int `json:"pageSize"`
	// TotalPages is the total number of pages.
	TotalPages int `json:"totalPages"`
	// HasNextPage reports whether a page after the current one exists.
	HasNextPage bool `json:"hasNextPage"`
	// HasPreviousPage reports whether a page before the current one exists.
	HasPreviousPage bool `json:"hasPreviousPage"`
}
