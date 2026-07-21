package fieldservice

import (
	"encoding/json"
	"strings"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestCanConstructFieldServiceRequest(t *testing.T) {
	request := FieldServiceRequest{
		Vehicles: []Vehicle{
			{
				ID:       "1",
				Location: [2]float64{51.52, -0.1},
				Shifts: []Shift{
					{ID: "full-day", MinStartTime: "2024-01-15T08:00:00", MaxEndTime: "2024-01-15T18:00:00"},
				},
				Skills: []string{"English"},
			},
		},
		Visits: []Visit{
			{
				ID:       "1",
				Name:     "Task 1",
				Location: [2]float64{51.51, -0.12},
				TimeWindows: []TimeWindow{
					{MinStartTime: "2024-01-15T09:00:00", MaxEndTime: "2024-01-15T12:00:00"},
				},
				ServiceDuration: "PT30M",
				Priority:        "HIGH",
				RequiredSkills:  []string{"English"},
			},
		},
	}

	if len(request.Vehicles) != 1 {
		t.Fatalf("expected 1 vehicle, got %d", len(request.Vehicles))
	}
	if len(request.Visits) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(request.Visits))
	}
}

func TestCanSerializeFieldServiceRequestToJSON(t *testing.T) {
	request := FieldServiceRequest{
		Vehicles: []Vehicle{},
		Visits:   []Visit{},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, "vehicles") {
		t.Errorf("expected JSON to contain 'vehicles'")
	}
	if !strings.Contains(jsonStr, "visits") {
		t.Errorf("expected JSON to contain 'visits'")
	}
}

func TestCanDeserializeFieldServiceResultResponse(t *testing.T) {
	jsonStr := `{
		"vehicles": [
			{"id": "1", "location": [51.52, -0.1], "shifts": [], "skills": ["English"], "visits": ["1"]}
		],
		"visits": [
			{
				"id": "1",
				"name": "Task 1",
				"location": [51.51, -0.12],
				"timeWindows": [],
				"serviceDuration": "PT30M",
				"priority": "HIGH",
				"requiredSkills": ["English"],
				"vehicle": "1",
				"previousVisit": null,
				"arrivalTime": "2024-01-15T10:15:33",
				"departureTime": "2024-01-15T10:45:33",
				"startServiceTime": "2024-01-15T10:15:33",
				"drivingTimeSecondsFromPreviousStandstill": 933
			}
		]
	}`

	var result FieldServiceResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.Vehicles) != 1 {
		t.Fatalf("expected 1 vehicle, got %d", len(result.Vehicles))
	}
	if len(result.Visits) != 1 {
		t.Fatalf("expected 1 visit, got %d", len(result.Visits))
	}
	if result.Visits[0].Vehicle != "1" {
		t.Errorf("expected vehicle '1', got '%s'", result.Visits[0].Vehicle)
	}
	if result.Visits[0].ArrivalTime != "2024-01-15T10:15:33" {
		t.Errorf("expected arrivalTime '2024-01-15T10:15:33', got '%s'", result.Visits[0].ArrivalTime)
	}
	if result.Visits[0].DepartureTime != "2024-01-15T10:45:33" {
		t.Errorf("expected departureTime '2024-01-15T10:45:33', got '%s'", result.Visits[0].DepartureTime)
	}
	if result.Visits[0].StartServiceTime != "2024-01-15T10:15:33" {
		t.Errorf("expected startServiceTime '2024-01-15T10:15:33', got '%s'", result.Visits[0].StartServiceTime)
	}
	if result.Visits[0].DrivingTimeSecondsFromPreviousStandstill != 933 {
		t.Errorf("expected drivingTimeSecondsFromPreviousStandstill 933, got %d", result.Visits[0].DrivingTimeSecondsFromPreviousStandstill)
	}
}
