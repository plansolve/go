# PlanSolve for Go

Official Go client for the [PlanSolve](https://getplansolve.com) optimization API. One typed client covers three solvers (field service routing, professional-services task assignment, and shift scheduling) with context support, built-in polling, and clean errors.

## Installation

```bash
go get github.com/plansolve/go
```

Requires Go 1.21+.

## Quick start

The module's import path ends in `go`, so alias it as `plansolve`:

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

	request := fieldservice.FieldServiceRequest{
		Vehicles: []fieldservice.Vehicle{{
			ID:       "tech1",
			Location: [2]float64{40.7128, -74.0060},
			Skills:   []string{"repair"},
			Shifts: []fieldservice.Shift{{
				ID:           "morning",
				MinStartTime: "2026-04-02T08:00:00",
				MaxEndTime:   "2026-04-02T17:00:00",
			}},
		}},
		Visits: []fieldservice.Visit{{
			ID:              "visit1",
			Name:            "AC Repair - Downtown Office",
			Location:        [2]float64{40.7589, -73.9851},
			ServiceDuration: "PT60M",
			Priority:        "HIGH",
			RequiredSkills:  []string{"repair"},
			TimeWindows: []fieldservice.TimeWindow{{
				MinStartTime: "2026-04-02T09:00:00",
				MaxEndTime:   "2026-04-02T17:00:00",
			}},
		}},
	}

	// Submit and block until the optimized plan is ready
	result, err := client.FieldService.StartAndWaitForCompletion(context.Background(), request, 5000, 150)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range result.Vehicles {
		fmt.Printf("Vehicle %s: %d visits\n", v.ID, len(v.Visits))
	}
}
```

## Solvers

One client, three solvers. All share the same submit, poll, result workflow:

| Solver | Field | Use for |
|--------|-------|---------|
| Field Service | `client.FieldService` | Vehicle routing with travel time, time windows, and skills |
| Professional Services | `client.ProfessionalServices` | Task assignment by skill, availability, priority, and deadlines |
| Shift | `client.Shift` | Shift scheduling across contracts, availability, and fairness |

Each field exposes the same methods. Every method takes a `context.Context` first:

| Method | What it does |
|--------|--------------|
| `Start(ctx, request)` | Submits a solve and returns the job ID |
| `GetStatus(ctx, jobID)` | Returns the current `SolverStatusResponse` |
| `GetResult(ctx, jobID)` | Returns the solver's result type |
| `Analyze(ctx, jobID)` | Returns the constraint analysis as `map[string]interface{}` (`GET /api/v1/{solver}/{jobId}/analyze`) |
| `Stop(ctx, jobID)` | Stops an in-progress solve and returns the best solution so far, same type as `GetResult` (`DELETE /api/v1/{solver}/{jobId}`) |
| `WaitForCompletion(ctx, jobID, pollIntervalMs, maxAttempts)` | Polls until the solve finishes, then returns the result |
| `StartAndWaitForCompletion(ctx, request, pollIntervalMs, maxAttempts)` | `Start` followed by `WaitForCompletion` |

## Polling

The wait methods take `pollIntervalMs, maxAttempts` (both `int`). When `pollIntervalMs <= 0` it defaults to `5000`, and when `maxAttempts <= 0` it defaults to `150`, so the default budget is 5 s x 150 = 12.5 minutes. That leaves headroom over the server's 10-minute cap on a single solve. `WaitForCompletion` waits one interval before its first status check.

A job counts as finished when the status reports `solving: false` and `solverStatus` is `NOT_SOLVING`, the same rule the server uses. A score is not required.

## Configuration

| Variable | Purpose |
|----------|---------|
| `PLANSOLVE_API_KEY` | API key, sent as `X-API-KEY`. Used by `plansolve.NewClientFromEnv()`. |
| `PLANSOLVE_BASE_URL` | Override the API base URL (default `https://plansolve.app`). |

Or pass the key directly: `plansolve.NewClient("...")`.

## Error handling

Methods return a plain `error` whose message is already readable. The library never prints to stdout or stderr; everything is reported through the returned error.

| Situation | Error |
|-----------|-------|
| Non-2xx HTTP response | `API error: status <code>: <message>`. A `400` lists the field validation messages, a `402` or `422` carries the body's `error`, a `401`, `403`, or `502` carries the problem `detail` (or `title`). |
| Solve failed | The status endpoint answers `422`, so the wait methods return `API error: status 422: <message>`. |
| Poll budget exhausted | `solver still running after N polls; raise maxAttempts or lower options.spentLimit` |
| Empty job ID passed to `WaitForCompletion` | `jobId was not provided` |
| Context cancelled or timed out | `ctx.Err()` (`context.Canceled` or `context.DeadlineExceeded`) |
| Transport or decode failure | `failed to send request: ...` / `failed to decode response: ...` (wrapped, usable with `errors.Is` / `errors.As`) |
| `NewClientFromEnv` without `PLANSOLVE_API_KEY` | The sentinel `plansolve.ErrMissingAPIKey` |

```go
result, err := client.FieldService.StartAndWaitForCompletion(context.Background(), request, 5000, 150)
if err != nil {
	log.Fatal(err)
}
```

## Documentation

Full guides, per-solver data models, and parameter reference live on the docs site:

- Field Service: https://getplansolve.com/docs/fieldservice/sdk/go
- Professional Services: https://getplansolve.com/docs/professionalservices/sdk/go
- Shift: https://getplansolve.com/docs/shiftsolver/sdk/go

Package: [pkg.go.dev](https://pkg.go.dev/github.com/plansolve/go)

## License

Apache-2.0
