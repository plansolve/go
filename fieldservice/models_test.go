package fieldservice

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/plansolve/go/solver"
)

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
				"serviceDuration": 3600.0,
				"priority": "HIGH",
				"requiredSkills": ["English"],
				"vehicle": "1",
				"previousVisit": null,
				"arrivalTime": "2024-01-15T10:15:33",
				"minStartTime": "2024-01-15T09:00:00",
				"maxEndTime": "2024-01-15T12:00:00",
				"departureTime": "2024-01-15T10:45:33",
				"startServiceTime": "2024-01-15T10:15:33",
				"isDayHead": true,
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
	if result.Visits[0].ServiceDuration != 3600.0 {
		t.Errorf("expected serviceDuration 3600.0, got %v", result.Visits[0].ServiceDuration)
	}
	if result.Visits[0].IsDayHead == nil || *result.Visits[0].IsDayHead != true {
		t.Errorf("expected isDayHead true, got %v", result.Visits[0].IsDayHead)
	}
	if result.Visits[0].MinStartTime == nil || *result.Visits[0].MinStartTime != "2024-01-15T09:00:00" {
		t.Errorf("expected minStartTime '2024-01-15T09:00:00', got %v", result.Visits[0].MinStartTime)
	}
	if result.Visits[0].MaxEndTime == nil || *result.Visits[0].MaxEndTime != "2024-01-15T12:00:00" {
		t.Errorf("expected maxEndTime '2024-01-15T12:00:00', got %v", result.Visits[0].MaxEndTime)
	}
}

func TestFieldServiceResultResponseScoreRoundTrips(t *testing.T) {
	// A caller may cache a result to disk or log it as JSON, then read it back
	// with the same SDK. Score.MarshalJSON must emit the canonical string form
	// so that round trip succeeds - a bare {"hard":...,"medium":...,"soft":...}
	// (the struct's former json tags) is rejected by ParseScoreJSON.
	original := FieldServiceResultResponse{
		Vehicles: []ScheduledVehicle{},
		Visits:   []ScheduledVisit{},
		Score:    &solver.Score{Hard: 0, Medium: 0, Soft: -5},
	}

	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var roundTripped FieldServiceResultResponse
	if err := json.Unmarshal(data, &roundTripped); err != nil {
		t.Fatalf("unmarshal round-tripped JSON %s: %v", data, err)
	}

	if roundTripped.Score == nil || *roundTripped.Score != *original.Score {
		t.Errorf("score did not survive round trip: got %+v, want %+v", roundTripped.Score, original.Score)
	}
}

func TestCanDeserializeExtendedVehicleAndVisitFields(t *testing.T) {
	jsonStr := `{
		"score": "0hard/0medium/-30soft",
		"totalDrivingTimeSeconds": 5400,
		"vehicles": [
			{
				"id": "v1",
				"name": "Van 1",
				"location": [51.52, -0.1],
				"skills": ["English"],
				"shifts": [],
				"visits": ["1"],
				"dailyReturnTimes": {"2024-01-15": "2024-01-15T17:30:00"},
				"arrivalTime": "2024-01-15T17:30:00",
				"totalDrivingTimeSeconds": 5400
			}
		],
		"visits": [
			{
				"id": "1",
				"name": "Task 1",
				"location": [51.51, -0.12],
				"timeWindows": [],
				"serviceDuration": 1800.0,
				"priority": "HIGH",
				"requiredSkills": ["English"],
				"pinned": true,
				"vehicle": "v1",
				"arrivalTime": "2024-01-15T10:15:33",
				"departureTime": "2024-01-15T10:45:33",
				"startServiceTime": "2024-01-15T10:15:33",
				"drivingTimeSecondsFromPreviousStandstill": 933
			}
		]
	}`

	// FieldServiceSolution is a spec-name alias for FieldServiceResultResponse.
	var result FieldServiceSolution
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.Score == nil || result.Score.String() != "0hard/0medium/-30soft" {
		t.Errorf("expected score, got %v", result.Score)
	}
	if result.TotalDrivingTimeSeconds != 5400 {
		t.Errorf("expected totalDrivingTimeSeconds 5400, got %d", result.TotalDrivingTimeSeconds)
	}

	v := result.Vehicles[0]
	if v.Name == nil || *v.Name != "Van 1" {
		t.Errorf("expected vehicle name 'Van 1', got %v", v.Name)
	}
	if v.DailyReturnTimes["2024-01-15"] != "2024-01-15T17:30:00" {
		t.Errorf("expected dailyReturnTimes entry, got %v", v.DailyReturnTimes)
	}
	if v.TotalDrivingTimeSeconds != 5400 {
		t.Errorf("expected vehicle totalDrivingTimeSeconds 5400, got %d", v.TotalDrivingTimeSeconds)
	}

	vis := result.Visits[0]
	if vis.Pinned == nil || !*vis.Pinned {
		t.Errorf("expected visit pinned true, got %v", vis.Pinned)
	}
}
