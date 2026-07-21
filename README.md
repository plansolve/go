# PlanSolve Go SDK

A Go client library for the PlanSolve optimization API. Solve complex field service routing, professional services task assignment, and shift scheduling problems with a simple, async API.

## Installation

```bash
go get github.com/plansolve/go
```

## Requirements

- Go 1.21+
- PlanSolve API key (optional but recommended for production use)

## API Workflow

All PlanSolve optimization APIs follow the same workflow:

1. **Construct Request** - Build your optimization request with the required data
2. **Start Solver** - Submit the request to start the optimization process
3. **Poll Status** - Check the status until the solver completes
4. **Get Results** - Retrieve the optimized solution

## Quick Start

### Field Service Optimization

To instantiate the client and perform a basic field service optimization:

```go
package main

import (
	"context"
	"fmt"
	"log"

	plansolve "github.com/plansolve/go"
	"github.com/plansolve/go/fieldservice"
)

func main() {
	client := plansolve.NewClient("YOUR_API_KEY")
	ctx := context.Background()

	// Create vehicles (field technicians)
	vehicles := []fieldservice.Vehicle{
		{
			ID:       "tech1",
			Location: []float64{40.7128, -74.0060},
			Shifts: []fieldservice.Shift{
				{
					ID:           "morning-shift",
					MinStartTime: "2024-01-15T08:00:00",
					MaxEndTime:   "2024-01-15T17:00:00",
				},
			},
			Skills: []string{"repair", "installation"},
		},
		{
			ID:       "tech2",
			Location: []float64{40.7128, -74.0060},
			Shifts: []fieldservice.Shift{
				{
					ID:           "afternoon-shift",
					MinStartTime: "2024-01-15T10:00:00",
					MaxEndTime:   "2024-01-15T19:00:00",
				},
			},
			Skills: []string{"repair"},
		},
	}

	// Create visits (customer appointments)
	visitName1 := "Times Square Repair"
	visitName2 := "Penn Station Installation"

	visits := []fieldservice.Visit{
		{
			ID:       "visit1",
			Name:     &visitName1,
			Location: []float64{40.7589, -73.9851},
			TimeWindows: []fieldservice.TimeWindow{
				{
					MinStartTime: "2024-01-15T09:00:00",
					MaxEndTime:   "2024-01-15T17:00:00",
				},
			},
			ServiceDuration: "PT30M",
			Priority:        "HIGH",
			RequiredSkills:  []string{"repair"},
		},
		{
			ID:       "visit2",
			Name:     &visitName2,
			Location: []float64{40.7505, -73.9934},
			TimeWindows: []fieldservice.TimeWindow{
				{
					MinStartTime: "2024-01-15T10:00:00",
					MaxEndTime:   "2024-01-15T16:00:00",
				},
			},
			ServiceDuration: "PT45M",
			Priority:        "MEDIUM",
			RequiredSkills:  []string{"installation"},
		},
	}

	request := fieldservice.FieldServiceRequest{
		Vehicles: vehicles,
		Visits:   visits,
	}

	// Option 1: Start and wait for completion
	result, err := client.FieldService.StartAndWaitForCompletion(ctx, request, 5000, 1000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Optimization completed with %d vehicles\n", len(result.Vehicles))

	// Option 2: Start and poll manually
	resp, err := client.FieldService.Start(ctx, request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Job ID: %s\n", resp.JobID)

	status, err := client.FieldService.GetStatus(ctx, resp.JobID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Status: %s\n", status.SolverStatus)
}
```

### Professional Services Optimization

For professional services task assignment optimization:

```go
package main

import (
	"context"
	"fmt"
	"log"

	plansolve "github.com/plansolve/go"
	"github.com/plansolve/go/professionalservices"
)

func main() {
	client := plansolve.NewClient("YOUR_API_KEY")
	ctx := context.Background()

	employees := []professionalservices.Employee{
		{
			ID: "emp1",
			Shifts: []professionalservices.Shift{
				{
					ID:           "shift1",
					MinStartTime: "2024-01-15T08:00:00",
					MaxEndTime:   "2024-01-15T18:00:00",
				},
			},
			Skills: []string{"Java", "Spring", "Kotlin"},
		},
		{
			ID: "emp2",
			Shifts: []professionalservices.Shift{
				{
					ID:           "shift2",
					MinStartTime: "2024-01-15T09:00:00",
					MaxEndTime:   "2024-01-15T17:00:00",
				},
			},
			Skills: []string{"Python", "SQL", "React"},
		},
	}

	deadline1 := "2024-01-20T17:00:00"
	deadline2 := "2024-01-22T17:00:00"

	tasks := []professionalservices.Task{
		{
			ID:             "task1",
			Name:           "Develop REST API",
			Deadline:       &deadline1,
			Duration:       "PT16H",
			Priority:       "HIGH",
			RequiredSkills: []string{"Java", "Spring"},
		},
		{
			ID:             "task2",
			Name:           "Database Design",
			Deadline:       &deadline2,
			Duration:       "PT8H",
			Priority:       "MEDIUM",
			RequiredSkills: []string{"SQL"},
		},
	}

	request := professionalservices.ProfessionalServicesRequest{
		Employees: employees,
		Tasks:     tasks,
	}

	result, err := client.ProfessionalServices.StartAndWaitForCompletion(ctx, request, 5000, 1000)
	if err != nil {
		log.Fatal(err)
	}

	for _, emp := range result.Employees {
		fmt.Printf("Employee %s assigned tasks: %v\n", emp.ID, emp.Tasks)
	}
}
```

### Shift Assignment Optimization

For shift scheduling and assignment:

```go
package main

import (
	"context"
	"fmt"
	"log"

	plansolve "github.com/plansolve/go"
	"github.com/plansolve/go/shift"
)

func main() {
	client := plansolve.NewClient("YOUR_API_KEY")
	ctx := context.Background()

	contractName := "full-time"
	empName1 := "Alice"
	empName2 := "Bob"

	employees := []shift.ShiftEmployee{
		{
			Name:         &empName1,
			Contract:     &contractName,
			Skills:       []string{"cashier", "inventory"},
			Availability: []string{"2024-01-15", "2024-01-16"},
		},
		{
			Name:         &empName2,
			Contract:     &contractName,
			Skills:       []string{"manager", "cashier"},
			Availability: []string{"2024-01-15", "2024-01-16"},
		},
	}

	shiftName1 := "Morning Shift"
	shiftFrom1 := "2024-01-15T06:00:00"
	shiftTo1 := "2024-01-15T14:00:00"
	shiftName2 := "Evening Shift"
	shiftFrom2 := "2024-01-15T14:00:00"
	shiftTo2 := "2024-01-15T22:00:00"

	shifts := []shift.ShiftAssignment{
		{
			Name:     &shiftName1,
			From:     &shiftFrom1,
			To:       &shiftTo1,
			Skills:   []string{"cashier"},
			Cost:     0,
			Value:    1,
			Priority: 1,
		},
		{
			Name:     &shiftName2,
			From:     &shiftFrom2,
			To:       &shiftTo2,
			Skills:   []string{"manager"},
			Cost:     0,
			Value:    1,
			Priority: 1,
		},
	}

	request := shift.ShiftRequest{
		Employees: employees,
		Shifts:    shifts,
	}

	result, err := client.Shift.StartAndWaitForCompletion(ctx, request, 5000, 1000)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range result.Tasks {
		emp := ""
		if task.Employee != nil {
			emp = *task.Employee
		}
		fmt.Printf("Task %s assigned to %s\n", task.Name, emp)
	}
}
```

## API Reference

### Client

The main client struct for interacting with the PlanSolve API.

#### Constructors

```go
plansolve.NewClient(apiKey string) *Client          // With API key
plansolve.NewClientFromEnv() (*Client, error)       // Uses PLANSOLVE_API_KEY env var
```

#### Fields

- `FieldService` - Field service optimization client
- `ProfessionalServices` - Professional services optimization client
- `Shift` - Shift assignment optimization client

### FieldService Client

- `Start(ctx, FieldServiceRequest) (*FieldServiceStartResponse, error)`
- `GetResult(ctx, jobID) (*FieldServiceResultResponse, error)`
- `GetStatus(ctx, jobID) (*SolverStatusResponse, error)`
- `StartAndWaitForCompletion(ctx, request, pollIntervalMs, maxAttempts) (*FieldServiceResultResponse, error)`
- `WaitForCompletion(ctx, jobID, pollIntervalMs, maxAttempts) (*FieldServiceResultResponse, error)`
- `Analyze(ctx, jobID) (map[string]interface{}, error)`

### ProfessionalServices Client

- `Start(ctx, ProfessionalServicesRequest) (*ProfessionalServicesStartResponse, error)`
- `GetResult(ctx, jobID) (*ProfessionalServicesResultResponse, error)`
- `GetStatus(ctx, jobID) (*SolverStatusResponse, error)`
- `StartAndWaitForCompletion(ctx, request, pollIntervalMs, maxAttempts) (*ProfessionalServicesResultResponse, error)`
- `WaitForCompletion(ctx, jobID, pollIntervalMs, maxAttempts) (*ProfessionalServicesResultResponse, error)`
- `Analyze(ctx, jobID) (map[string]interface{}, error)`

### Shift Client

- `Start(ctx, ShiftRequest) (*ShiftStartResponse, error)`
- `GetResult(ctx, jobID) (*ShiftResultResponse, error)`
- `GetStatus(ctx, jobID) (*SolverStatusResponse, error)`
- `StartAndWaitForCompletion(ctx, request, pollIntervalMs, maxAttempts) (*ShiftResultResponse, error)`
- `WaitForCompletion(ctx, jobID, pollIntervalMs, maxAttempts) (*ShiftResultResponse, error)`
- `Analyze(ctx, jobID) (map[string]interface{}, error)`

## Data Models

### Field Service Models

#### Vehicle

```go
type Vehicle struct {
    ID       string    `json:"id"`
    Location []float64 `json:"location"`   // [latitude, longitude]
    Shifts   []Shift   `json:"shifts"`
    Skills   []string  `json:"skills"`
}
```

#### Visit

```go
type Visit struct {
    ID              string       `json:"id"`
    Name            *string      `json:"name,omitempty"`
    Location        []float64    `json:"location"`
    TimeWindows     []TimeWindow `json:"timeWindows"`
    ServiceDuration string       `json:"serviceDuration"` // ISO 8601 duration
    Priority        string       `json:"priority"`
    RequiredSkills  []string     `json:"requiredSkills"`
}
```

### Professional Services Models

#### Employee

```go
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
```

### Shift Models

#### ShiftAssignment

```go
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
```

#### ShiftEmployee

```go
type ShiftEmployee struct {
    Name                  *string           `json:"name,omitempty"`
    Contract              *string           `json:"contract,omitempty"`
    Skills                []string          `json:"skills,omitempty"`
    LastRestDate          *string           `json:"lastRestDate,omitempty"`
    Availability          []string          `json:"availability,omitempty"`
    Preference            []string          `json:"preference,omitempty"`
    PeriodRules           []PeriodRule      `json:"periodRules,omitempty"`
    UnavailableDates      []string          `json:"unavailableDates,omitempty"`
    Tags                  []string          `json:"tags,omitempty"`
    MaximumMinutesPerWeek *int              `json:"maximumMinutesPerWeek,omitempty"`
    Shifts                []ShiftAssignment `json:"shifts,omitempty"`
}
```

## Advanced Usage

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PLANSOLVE_API_KEY` | API key for authentication | None |
| `PLANSOLVE_BASE_URL` | Base URL for the PlanSolve API | `https://plansolve.app` |

### Context and Cancellation

All methods accept a `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := client.FieldService.StartAndWaitForCompletion(ctx, request, 5000, 1000)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Optimization timed out")
    }
    log.Fatal(err)
}
```

### Error Handling

```go
result, err := client.FieldService.Start(ctx, request)
if err != nil {
    if strings.Contains(err.Error(), "status 401") {
        log.Println("Invalid API key")
    } else if strings.Contains(err.Error(), "status 400") {
        log.Println("Invalid request data")
    } else if strings.Contains(err.Error(), "status 429") {
        log.Println("Rate limit exceeded")
    } else {
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Support

For API support and questions, please refer to the main PlanSolve documentation or contact support.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
