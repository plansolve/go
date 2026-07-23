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
)

// CreateJobRequest is the payload for creating a new solver job via the public Jobs API.
type CreateJobRequest struct {
	// SubscriptionID is the subscription the job should be billed against. Resolved
	// server-side from the authenticated tenant when omitted by the caller.
	SubscriptionID *int `json:"subscriptionId,omitempty"`
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
	// 1 = Professional Services, 2 = Shift).
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
