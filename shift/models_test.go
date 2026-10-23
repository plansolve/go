package shift

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestCanConstructShiftRequest(t *testing.T) {
	request := ShiftRequest{
		ID:          strPtr("req1"),
		Name:        strPtr("Weekly Schedule"),
		Description: strPtr("Week 3 shift planning"),
		Contracts: []Contract{
			{
				Name:                   strPtr("FULL_TIME"),
				MinWorkDurationPerWeek: strPtr("PT32H"),
				MaxWorkDurationPerWeek: strPtr("PT40H"),
				MaxConsecutiveWorkDays: intPtr(5),
				MaxShiftsDay:           intPtr(1),
				MinRestBetweenShifts:   strPtr("PT11H"),
				MaxWorkingDaysPerWeek:  intPtr(5),
				EarliestShiftStart:     strPtr("06:00:00"),
				LatestShiftEnd:         strPtr("22:00:00"),
				MinConsecutiveDaysOff:  intPtr(2),
			},
		},
		Shifts: []ShiftAssignment{
			{
				Name:          strPtr("Morning"),
				From:          strPtr("2024-01-15T06:00:00"),
				To:            strPtr("2024-01-15T14:00:00"),
				Skills:        []string{"Nursing"},
				DesiredSkills: []string{"Pediatrics"},
				Tags:          []string{"early"},
				CostFactor:    float64Ptr(1.5),
				Value:         intPtr(2),
				Priority:      intPtr(1),
				PinnedByUser:  false,
			},
		},
		Employees: []ShiftEmployee{
			{
				ID:              strPtr("alice"),
				Name:            strPtr("Alice"),
				Contract:        strPtr("FULL_TIME"),
				Skills:          []string{"Nursing", "Pediatrics"},
				LastShiftEnd:    strPtr("2024-01-14T22:00:00"),
				PreferredShifts: []string{"Morning"},
				Tags:            []string{"senior"},
				CostPerHour:     42.5,
			},
		},
		TimeOffRequests: []ShiftTimeOffRequest{
			{ID: strPtr("to1"), EmployeeID: strPtr("alice"), From: strPtr("2024-01-18T00:00:00"), To: strPtr("2024-01-19T00:00:00")},
		},
		Unavailabilities: []ShiftUnavailability{
			{ID: strPtr("u1"), EmployeeID: strPtr("alice"), From: strPtr("2024-01-20T00:00:00"), To: strPtr("2024-01-21T00:00:00")},
		},
		ShiftOffRequests: []ShiftOffRequest{
			{ID: strPtr("sor1"), EmployeeID: strPtr("alice"), ShiftName: strPtr("Night"), Weight: intPtr(3)},
		},
		Fairness: &Fairness{
			FairnessBuckets: []FairnessBucket{
				{Name: strPtr("weekend"), EmployeeIDs: []string{"alice", "bob"}, Shifts: []string{"Saturday", "Sunday"}, Period: strPtr("2024-01-15/2024-01-22")},
			},
		},
		Hook:    strPtr("https://example.com/webhook"),
		Weights: map[string]string{"shiftOffRequest": "0hard/0medium/4soft"},
		Options: &SolverOptions{SpentLimit: strPtr("PT30S"), UnimprovedSpentLimit: strPtr("PT5S")},
	}

	if *request.ID != "req1" {
		t.Errorf("expected ID 'req1', got '%s'", *request.ID)
	}
	if len(request.Contracts) != 1 || *request.Contracts[0].MaxConsecutiveWorkDays != 5 {
		t.Errorf("expected 1 contract with maxConsecutiveWorkDays 5, got %v", request.Contracts)
	}
	if len(request.Shifts) != 1 || *request.Shifts[0].Value != 2 {
		t.Errorf("expected 1 shift with value 2, got %v", request.Shifts)
	}
	if len(request.Employees) != 1 || request.Employees[0].CostPerHour != 42.5 {
		t.Errorf("expected 1 employee with costPerHour 42.5, got %v", request.Employees)
	}
	if len(request.TimeOffRequests) != 1 {
		t.Errorf("expected 1 time-off request, got %d", len(request.TimeOffRequests))
	}
	if len(request.Unavailabilities) != 1 {
		t.Errorf("expected 1 unavailability, got %d", len(request.Unavailabilities))
	}
	if len(request.ShiftOffRequests) != 1 || *request.ShiftOffRequests[0].EmployeeID != "alice" {
		t.Errorf("expected 1 shift-off request for 'alice', got %v", request.ShiftOffRequests)
	}
	if request.Weights["shiftOffRequest"] != "0hard/0medium/4soft" {
		t.Errorf("expected weights to round-trip, got %v", request.Weights)
	}
	if *request.Options.SpentLimit != "PT30S" {
		t.Errorf("expected options.spentLimit 'PT30S', got '%s'", *request.Options.SpentLimit)
	}
}

func TestCanSerializeShiftRequestToJSON(t *testing.T) {
	request := ShiftRequest{
		ID:   strPtr("req1"),
		Name: strPtr("Test Schedule"),
		Contracts: []Contract{
			{Name: strPtr("FULL_TIME"), MaxWorkDurationPerWeek: strPtr("PT40H"), MinConsecutiveDaysOff: intPtr(2)},
		},
		Employees: []ShiftEmployee{
			{
				ID:              strPtr("alice"),
				Name:            strPtr("Alice"),
				Contract:        strPtr("FULL_TIME"),
				Skills:          []string{"Nursing"},
				LastShiftEnd:    strPtr("2024-01-14T22:00:00"),
				PreferredShifts: []string{"Morning"},
				CostPerHour:     30,
			},
		},
		Shifts: []ShiftAssignment{
			{
				Name:          strPtr("Morning"),
				DesiredSkills: []string{"Pediatrics"},
				Tags:          []string{"early"},
				CostFactor:    float64Ptr(1.25),
				Value:         intPtr(2),
				PinnedByUser:  true,
			},
		},
		TimeOffRequests: []ShiftTimeOffRequest{
			{ID: strPtr("to1"), EmployeeID: strPtr("alice"), From: strPtr("2024-01-18T00:00:00"), To: strPtr("2024-01-19T00:00:00")},
		},
		Unavailabilities: []ShiftUnavailability{
			{ID: strPtr("u1"), EmployeeID: strPtr("alice"), From: strPtr("2024-01-20T00:00:00"), To: strPtr("2024-01-21T00:00:00")},
		},
		ShiftOffRequests: []ShiftOffRequest{
			{ID: strPtr("sor1"), EmployeeID: strPtr("alice"), ShiftName: strPtr("Night"), Weight: intPtr(3)},
		},
		Fairness: &Fairness{
			FairnessBuckets: []FairnessBucket{{Name: strPtr("all"), EmployeeIDs: []string{"alice"}}},
		},
		Hook:    strPtr("https://example.com/webhook"),
		Weights: map[string]string{"shiftOffRequest": "0hard/0medium/4soft"},
		Options: &SolverOptions{SpentLimit: strPtr("PT30S")},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	for _, expected := range []string{
		`"maxWorkDurationPerWeek":"PT40H"`, `"minConsecutiveDaysOff":2`,
		`"id":"alice"`, `"lastShiftEnd":"2024-01-14T22:00:00"`, `"preferredShifts":["Morning"]`, `"costPerHour":30`,
		`"contract":"FULL_TIME"`,
		`"desiredSkills":["Pediatrics"]`, `"costFactor":1.25`, `"value":2`, `"pinnedByUser":true`,
		`"timeOffRequests":[{"id":"to1","employeeId":"alice","from":"2024-01-18T00:00:00","to":"2024-01-19T00:00:00"}]`,
		`"unavailabilities":[{"id":"u1","employeeId":"alice","from":"2024-01-20T00:00:00","to":"2024-01-21T00:00:00"}]`,
		`"shiftOffRequests":[{"id":"sor1","employeeId":"alice","shiftName":"Night","weight":3}]`,
		`"fairnessBuckets":[{"name":"all","employeeIds":["alice"]}]`,
		`"hook":"https://example.com/webhook"`,
		`"weights":{"shiftOffRequest":"0hard/0medium/4soft"}`,
		`"options":{"spentLimit":"PT30S"}`,
	} {
		if !strings.Contains(jsonStr, expected) {
			t.Errorf("expected JSON to contain '%s', got %s", expected, jsonStr)
		}
	}

	// Unset nullable fields are omitted rather than sent as zero values.
	for _, unexpected := range []string{
		`"priority"`, `"maxShiftsDay"`, `"maxConsecutiveWorkDays"`, `"unimprovedSpentLimit"`, `"assignedEmployee"`,
		`"originalName"`,
		// Retired wire fields must not reappear.
		`"cost"`, `"dayOffRequests"`, `"constraintWeightOverrides"`, `"employeeName"`, `"periodRules"`,
		`"unavailableDates"`, `"maximumMinutesPerWeek"`, `"partialPlanning"`, `"timeLimit"`, `"employees":["`,
	} {
		if strings.Contains(jsonStr, unexpected) {
			t.Errorf("expected JSON not to contain '%s', got %s", unexpected, jsonStr)
		}
	}
}

func TestEmptySolverOptionsSerializesAsEmptyObject(t *testing.T) {
	data, err := json.Marshal(ShiftRequest{Options: &SolverOptions{}})
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(data) != `{"options":{}}` {
		t.Errorf("expected nil option fields to be omitted, got %s", data)
	}

	data, err = json.Marshal(ShiftRequest{})
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if string(data) != `{}` {
		t.Errorf("expected empty request to serialize as {}, got %s", data)
	}
}

func TestCanDeserializeShiftRequestFromJSON(t *testing.T) {
	jsonStr := `{
		"id": "req1",
		"name": "Test Schedule",
		"contracts": [
			{
				"name": "FULL_TIME",
				"minWorkDurationPerWeek": "PT32H",
				"maxWorkDurationPerWeek": "PT40H",
				"maxConsecutiveWorkDays": 5,
				"maxShiftsDay": 1,
				"minRestBetweenShifts": "PT12H",
				"maxWorkingDaysPerWeek": 5,
				"latestShiftEnd": "22:00:00",
				"earliestShiftStart": "06:00:00",
				"minConsecutiveDaysOff": 2
			}
		],
		"employees": [
			{
				"id": "alice",
				"name": "Alice",
				"contract": "FULL_TIME",
				"skills": ["Nursing"],
				"lastShiftEnd": "2024-01-14T22:00:00",
				"preferredShifts": ["Morning"],
				"tags": ["senior"],
				"costPerHour": 42.5
			}
		],
		"shifts": [
			{
				"name": "Morning",
				"from": "2024-01-15T06:00:00",
				"to": "2024-01-15T14:00:00",
				"desiredSkills": ["Pediatrics"],
				"tags": ["early"],
				"costFactor": 1.5,
				"value": 2,
				"priority": 1,
				"pinnedByUser": true
			}
		],
		"timeOffRequests": [
			{"id": "to1", "employeeId": "alice", "from": "2024-01-18T00:00:00", "to": "2024-01-19T00:00:00"}
		],
		"unavailabilities": [
			{"id": "u1", "employeeId": "alice", "from": "2024-01-20T00:00:00", "to": "2024-01-21T00:00:00"}
		],
		"shiftOffRequests": [
			{"id": "sor1", "employeeId": "alice", "shiftName": "Night", "weight": 3}
		],
		"fairness": {"fairnessBuckets": [{"name": "weekend", "employeeIds": ["alice"], "shifts": ["Saturday"]}]},
		"hook": "https://example.com/webhook",
		"weights": {"shiftOffRequest": "0hard/0medium/4soft"},
		"options": {"spentLimit": "PT30S", "unimprovedSpentLimit": "PT5S"}
	}`

	var result ShiftRequest
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	c := result.Contracts[0]
	if *c.MinWorkDurationPerWeek != "PT32H" || *c.MaxWorkDurationPerWeek != "PT40H" {
		t.Errorf("unexpected contract durations: %v / %v", *c.MinWorkDurationPerWeek, *c.MaxWorkDurationPerWeek)
	}
	if *c.MaxConsecutiveWorkDays != 5 || *c.MaxShiftsDay != 1 || *c.MaxWorkingDaysPerWeek != 5 || *c.MinConsecutiveDaysOff != 2 {
		t.Errorf("unexpected contract ints: %+v", c)
	}
	if *c.MinRestBetweenShifts != "PT12H" || *c.LatestShiftEnd != "22:00:00" || *c.EarliestShiftStart != "06:00:00" {
		t.Errorf("unexpected contract strings: %+v", c)
	}

	e := result.Employees[0]
	if *e.ID != "alice" || *e.Name != "Alice" || *e.Contract != "FULL_TIME" {
		t.Errorf("unexpected employee identity: %+v", e)
	}
	if *e.LastShiftEnd != "2024-01-14T22:00:00" || len(e.PreferredShifts) != 1 || e.CostPerHour != 42.5 {
		t.Errorf("unexpected employee fields: %+v", e)
	}

	s := result.Shifts[0]
	if *s.CostFactor != 1.5 || *s.Value != 2 || *s.Priority != 1 || !s.PinnedByUser {
		t.Errorf("unexpected shift fields: %+v", s)
	}

	if len(result.TimeOffRequests) != 1 || *result.TimeOffRequests[0].EmployeeID != "alice" || *result.TimeOffRequests[0].From != "2024-01-18T00:00:00" {
		t.Errorf("unexpected time-off requests: %v", result.TimeOffRequests)
	}
	if len(result.Unavailabilities) != 1 || *result.Unavailabilities[0].ID != "u1" || *result.Unavailabilities[0].To != "2024-01-21T00:00:00" {
		t.Errorf("unexpected unavailabilities: %v", result.Unavailabilities)
	}
	if *result.ShiftOffRequests[0].EmployeeID != "alice" || *result.ShiftOffRequests[0].Weight != 3 {
		t.Errorf("unexpected shift-off request: %+v", result.ShiftOffRequests[0])
	}
	if result.Fairness.FairnessBuckets[0].EmployeeIDs[0] != "alice" {
		t.Errorf("expected fairness bucket employeeIds ['alice'], got %v", result.Fairness.FairnessBuckets[0].EmployeeIDs)
	}
	if result.Weights["shiftOffRequest"] != "0hard/0medium/4soft" {
		t.Errorf("unexpected weights: %v", result.Weights)
	}
	if *result.Options.SpentLimit != "PT30S" || *result.Options.UnimprovedSpentLimit != "PT5S" {
		t.Errorf("unexpected options: %+v", result.Options)
	}
}

func TestCanDeserializeShiftResultResponseFromJSON(t *testing.T) {
	// testdata/shift_result.json is a Shift result as the API sends it: assignedShifts
	// and unassignedShifts hold the shift objects, each with its assignedEmployee
	// (plansolve/issues#173).
	data, err := os.ReadFile("testdata/shift_result.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	var result ShiftResultResponse
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible to be true")
	}
	if result.Score == nil || result.Score.Medium != -10000 || result.Score.Soft != -2220 {
		t.Errorf("expected score 0hard/-10000medium/-2220soft, got %+v", result.Score)
	}
	if len(result.Shifts) != 3 {
		t.Fatalf("expected 3 shifts, got %d", len(result.Shifts))
	}

	if len(result.AssignedShifts) != 2 {
		t.Fatalf("expected 2 assigned shifts, got %d", len(result.AssignedShifts))
	}
	morning := result.AssignedShifts[0]
	if *morning.Name != "morning_2026-04-06#1" || *morning.From != "2026-04-06T06:00:00" || *morning.To != "2026-04-06T14:00:00" {
		t.Errorf("unexpected assigned shift: %+v", morning)
	}
	if morning.AssignedEmployee == nil || *morning.AssignedEmployee != "alice" {
		t.Errorf("expected assignedEmployee 'alice', got %v", morning.AssignedEmployee)
	}
	if morning.OriginalName == nil || *morning.OriginalName != "morning_2026-04-06" {
		t.Errorf("expected originalName 'morning_2026-04-06', got %v", morning.OriginalName)
	}
	if morning.CostFactor == nil || *morning.CostFactor != 1.0 || morning.Priority == nil || *morning.Priority != 5 {
		t.Errorf("expected costFactor 1.0 and priority 5, got %+v", morning)
	}
	night := result.AssignedShifts[1]
	if *night.Name != "night_2026-04-06" || *night.AssignedEmployee != "carol" || night.OriginalName != nil {
		t.Errorf("unexpected night shift: %+v", night)
	}

	if len(result.UnassignedShifts) != 1 {
		t.Fatalf("expected 1 unassigned shift, got %d", len(result.UnassignedShifts))
	}
	unassigned := result.UnassignedShifts[0]
	if *unassigned.Name != "morning_2026-04-06#2" || *unassigned.From != "2026-04-06T06:00:00" {
		t.Errorf("unexpected unassigned shift: %+v", unassigned)
	}
	if unassigned.AssignedEmployee != nil {
		t.Errorf("expected unassigned shift to have nil assignedEmployee, got %v", *unassigned.AssignedEmployee)
	}

	if len(result.Employees) != 3 || *result.Employees[0].ID != "alice" {
		t.Errorf("unexpected employees: %v", result.Employees)
	}
}

func TestShiftResultResponseResolvesShiftNamesFromAnOlderAPI(t *testing.T) {
	// Before plansolve/issues#173 the API sent assignedShifts and unassignedShifts as
	// shift names. They resolve against shifts.
	jsonStr := `{
		"shifts": [
			{"name": "morning#1", "from": "2026-04-06T06:00:00", "to": "2026-04-06T14:00:00", "originalName": "morning", "assignedEmployee": "alice"},
			{"name": "morning#2", "from": "2026-04-06T06:00:00", "to": "2026-04-06T14:00:00", "originalName": "morning", "assignedEmployee": null}
		],
		"assignedShifts": ["morning#1", "ghost"],
		"unassignedShifts": ["morning#2"]
	}`

	var result ShiftResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.AssignedShifts) != 2 {
		t.Fatalf("expected 2 assigned shifts, got %d", len(result.AssignedShifts))
	}
	if *result.AssignedShifts[0].From != "2026-04-06T06:00:00" || *result.AssignedShifts[0].AssignedEmployee != "alice" {
		t.Errorf("expected morning#1 resolved from shifts, got %+v", result.AssignedShifts[0])
	}
	if *result.AssignedShifts[1].Name != "ghost" || result.AssignedShifts[1].From != nil {
		t.Errorf("expected an unknown name to keep only its name, got %+v", result.AssignedShifts[1])
	}
	if len(result.UnassignedShifts) != 1 || *result.UnassignedShifts[0].To != "2026-04-06T14:00:00" {
		t.Errorf("expected morning#2 resolved from shifts, got %+v", result.UnassignedShifts)
	}
}

func TestShiftResultResponseEchoesRequestFields(t *testing.T) {
	jsonStr := `{
		"id": "req1",
		"name": "Weekly Schedule",
		"description": "Week 3",
		"contracts": [
			{"name": "FULL_TIME", "minConsecutiveDaysOff": 2, "minRestBetweenShifts": "PT11H"}
		],
		"shifts": [
			{"name": "Morning", "from": "2024-01-15T06:00:00", "to": "2024-01-15T14:00:00", "pinnedByUser": false, "assignedEmployee": "alice"}
		],
		"timeOffRequests": [
			{"id": "to1", "employeeId": "alice", "from": "2024-01-18T00:00:00", "to": "2024-01-19T00:00:00"}
		],
		"unavailabilities": [
			{"id": "u1", "employeeId": "bob", "from": "2024-01-20T00:00:00", "to": "2024-01-21T00:00:00"}
		],
		"shiftOffRequests": [
			{"id": "sor1", "employeeId": "alice", "shiftName": "Night"}
		],
		"fairness": {"fairnessBuckets": [{"name": "all", "employeeIds": ["alice", "bob"]}]},
		"hook": "https://example.com/webhook",
		"constraintWeightOverrides": {"shiftOffRequest": "0hard/0medium/4soft"},
		"feasible": true,
		"score": "0hard/0medium/-120soft",
		"solutionInitialized": true,
		"assignedShifts": [
			{"name": "Morning", "from": "2024-01-15T06:00:00", "to": "2024-01-15T14:00:00", "pinnedByUser": false, "assignedEmployee": "alice"}
		],
		"unassignedShifts": [],
		"employees": [
			{"id": "alice", "name": "Alice", "contract": "FULL_TIME", "skills": ["Nursing"], "costPerHour": 0}
		]
	}`

	var result ShiftResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Echoed request fields.
	if result.ID == nil || *result.ID != "req1" {
		t.Errorf("expected echoed id 'req1', got %v", result.ID)
	}
	if result.Name == nil || *result.Name != "Weekly Schedule" {
		t.Errorf("expected echoed name, got %v", result.Name)
	}
	if len(result.Contracts) != 1 || result.Contracts[0].MinConsecutiveDaysOff == nil || *result.Contracts[0].MinConsecutiveDaysOff != 2 {
		t.Errorf("expected echoed contracts, got %v", result.Contracts)
	}
	if len(result.Shifts) != 1 || *result.Shifts[0].AssignedEmployee != "alice" {
		t.Errorf("expected echoed shifts with assignedEmployee, got %v", result.Shifts)
	}
	if len(result.TimeOffRequests) != 1 || *result.TimeOffRequests[0].ID != "to1" {
		t.Errorf("expected echoed timeOffRequests, got %v", result.TimeOffRequests)
	}
	if len(result.Unavailabilities) != 1 || *result.Unavailabilities[0].EmployeeID != "bob" {
		t.Errorf("expected echoed unavailabilities, got %v", result.Unavailabilities)
	}
	if len(result.ShiftOffRequests) != 1 || result.ShiftOffRequests[0].Weight != nil {
		t.Errorf("expected echoed shiftOffRequests with nil weight, got %v", result.ShiftOffRequests)
	}
	if result.Fairness == nil || len(result.Fairness.FairnessBuckets[0].EmployeeIDs) != 2 {
		t.Errorf("expected echoed fairness, got %v", result.Fairness)
	}
	if result.ConstraintWeightOverrides["shiftOffRequest"] != "0hard/0medium/4soft" {
		t.Errorf("expected echoed constraintWeightOverrides, got %v", result.ConstraintWeightOverrides)
	}
	if result.Hook == nil || *result.Hook != "https://example.com/webhook" {
		t.Errorf("expected echoed hook, got %v", result.Hook)
	}

	// Result-only fields.
	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible true")
	}
	if result.Score == nil || result.Score.Soft != -120 {
		t.Errorf("expected soft -120, got %+v", result.Score)
	}
	if len(result.AssignedShifts) != 1 || *result.AssignedShifts[0].AssignedEmployee != "alice" {
		t.Errorf("expected 1 assigned shift worked by alice, got %+v", result.AssignedShifts)
	}
}

func TestShiftResultSolutionInitialized(t *testing.T) {
	cases := []struct {
		name string
		json string
		want *bool
	}{
		{"present true", `{"solutionInitialized": true}`, boolPtr(true)},
		{"present false", `{"solutionInitialized": false}`, boolPtr(false)},
		{"absent", `{}`, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var result ShiftResultResponse
			if err := json.Unmarshal([]byte(tc.json), &result); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}
			switch {
			case tc.want == nil && result.SolutionInitialized != nil:
				t.Errorf("expected nil, got %v", *result.SolutionInitialized)
			case tc.want != nil && (result.SolutionInitialized == nil || *result.SolutionInitialized != *tc.want):
				t.Errorf("expected %v, got %v", *tc.want, result.SolutionInitialized)
			}
		})
	}
}

func TestShiftRequestDecodedFromResultDoesNotSendSolutionInitialized(t *testing.T) {
	// The server's Shift request rejects unknown members, so a result body reused as
	// a request must not carry solutionInitialized back.
	var request ShiftRequest
	if err := json.Unmarshal([]byte(`{"name": "Week", "solutionInitialized": true}`), &request); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if strings.Contains(string(data), "solutionInitialized") {
		t.Errorf("request must not carry solutionInitialized, got %s", data)
	}
}
