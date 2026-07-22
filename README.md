# PlanSolve for Go

Official Go client for the [PlanSolve](https://getplansolve.com) optimization API. One typed client covers three solvers — field service routing, professional-services task assignment, and shift scheduling — with context support, built-in polling, and clean errors.

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
	result, err := client.FieldService.StartAndWaitForCompletion(context.Background(), request, 5000, 60)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range result.Vehicles {
		fmt.Printf("Vehicle %s: %d visits\n", v.ID, len(v.Visits))
	}
}
```

## Solvers

One client, three solvers — all share the same submit → poll → result workflow:

| Solver | Field | Use for |
|--------|-------|---------|
| Field Service | `client.FieldService` | Vehicle routing with travel time, time windows, and skills |
| Professional Services | `client.ProfessionalServices` | Task assignment by skill, availability, priority, and deadlines |
| Shift | `client.Shift` | Shift scheduling across contracts, availability, and fairness |

Each field exposes `Start`, `GetStatus`, `GetResult`, `Analyze`, `WaitForCompletion`, and `StartAndWaitForCompletion`. The wait methods take `pollIntervalMs, maxAttempts` (both `int`).

## Configuration

| Variable | Purpose |
|----------|---------|
| `PLANSOLVE_API_KEY` | API key, sent as `X-API-KEY`. Used by `plansolve.NewClientFromEnv()`. |
| `PLANSOLVE_BASE_URL` | Override the API base URL (default `https://plansolve.app`). |

Or pass the key directly: `plansolve.NewClient("...")`.

## Error handling

Methods return a plain `error` whose message is already readable (e.g. `API error: status 400: vehicles: At least one vehicle is required.`). `NewClientFromEnv` returns the sentinel `plansolve.ErrMissingAPIKey` when the key is unset:

```go
result, err := client.FieldService.StartAndWaitForCompletion(context.Background(), request, 5000, 60)
if err != nil {
	log.Fatal(err)
}
```

## Documentation

Full guides, per-solver data models, and parameter reference live on the docs site:

- Field Service — https://getplansolve.com/docs/fieldservice/sdk/go
- Professional Services — https://getplansolve.com/docs/professionalservices/sdk/go
- Shift — https://getplansolve.com/docs/shiftsolver/sdk/go

Package: [pkg.go.dev](https://pkg.go.dev/github.com/plansolve/go)

## License

Apache-2.0
