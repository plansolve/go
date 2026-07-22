package shift

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

func boolPtr(b bool) *bool {
	return &b
}

func TestCanConstructShiftRequest(t *testing.T) {
	request := ShiftRequest{
		ID:          strPtr("req1"),
		Name:        strPtr("Weekly Schedule"),
		Description: strPtr("Week 3 shift planning"),
		Contracts: []Contract{
			{
				Name:                         strPtr("full-time"),
				Max:                          strPtr("PT40H"),
				Min:                          strPtr("PT32H"),
				MaxConsecutiveWorkDays:       5,
				MaxShiftsDay:                 2,
				MaxWorkingDays:               5,
				MinimumConsecutiveDaysOff:    2,
				MinimumHoursOffBetweenShifts: 11,
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
				Cost:          150.0,
				Value:         10,
				Priority:      1,
				PinnedByUser:  false,
			},
		},
		Employees: []ShiftEmployee{
			{
				Name:         strPtr("Alice"),
				Contract:     strPtr("full-time"),
				Skills:       []string{"Nursing", "Pediatrics"},
				LastRestDate: strPtr("2024-01-13"),
				Availability: []string{"2024-01-15", "2024-01-16", "2024-01-17"},
				Preference:   []string{"Morning"},
				PeriodRules: []PeriodRule{
					{
						Period:             &PlanningPeriod{From: strPtr("2024-01-15"), To: strPtr("2024-01-21")},
						MaxWorkingDays:     5,
						MinWorkingDays:     3,
						MinWorkingDuration: strPtr("PT24H"),
						MaxWorkingDuration: strPtr("PT40H"),
						MinRestDuration:    strPtr("PT11H"),
					},
				},
				UnavailableDates:      []string{"2024-01-18"},
				Tags:                  []string{"senior"},
				MaximumMinutesPerWeek: intPtr(2400),
				Shifts: []ShiftAssignment{
					{Name: strPtr("Morning"), From: strPtr("2024-01-14T06:00:00"), To: strPtr("2024-01-14T14:00:00")},
				},
			},
		},
		DayOffRequests: []DayOffRequest{
			{ID: strPtr("dor1"), EmployeeName: strPtr("Alice"), Date: strPtr("2024-01-18"), Weight: 5},
		},
		ShiftOffRequests: []ShiftOffRequest{
			{ID: strPtr("sor1"), EmployeeName: strPtr("Alice"), ShiftName: strPtr("Night"), Weight: 3},
		},
		Hook: strPtr("https://example.com/webhook"),
		ConstraintWeightOverrides: &ConstraintWeightOverrides{
			KnownConstraintNames: []string{"requiredSkills", "noDoubleBooking"},
		},
		Options: &Options{PartialPlanning: boolPtr(true), MaxIterations: intPtr(1000), TimeLimit: intPtr(60)},
		Weights: &Weights{
			RequiredSkills:            10,
			ShiftCapacity:             8,
			MinimumStaffing:           9,
			NoDoubleBooking:           10,
			RestBetweenShifts:         7,
			EmployeeAvailability:      8,
			MaxConsecutiveWorkDays:    6,
			MaxShiftsPerDay:           10,
			MaxWorkingDaysPerWeek:     5,
			ContractRestBetweenShifts: 7,
			EarliestShiftStart:        3,
			LatestShiftStart:          3,
			MinimumConsecutiveDaysOff: 6,
			PeriodRuleViolation:       5,
			ShiftPreferences:          2,
			CostMinimization:          4,
			WorkloadBalance:           3,
			Fairness:                  3,
			DesiredSkills:             2,
			DesiredDayOff:             4,
			ShiftOffRequest:           4,
			BalanceTimeWorked:         3,
			EmployeeAffinity:         2,
			AvoidShiftCloseToDayOff:   3,
		},
		Fairness: &Fairness{
			FairnessBuckets: []FairnessBucket{
				{Name: strPtr("weekend"), Employees: []string{"Alice", "Bob"}, Shifts: []string{"Saturday", "Sunday"}, Period: strPtr("WEEKLY")},
			},
		},
	}

	if *request.ID != "req1" {
		t.Errorf("expected ID 'req1', got '%s'", *request.ID)
	}
	if len(request.Contracts) != 1 {
		t.Errorf("expected 1 contract, got %d", len(request.Contracts))
	}
	if len(request.Shifts) != 1 {
		t.Errorf("expected 1 shift, got %d", len(request.Shifts))
	}
	if len(request.Employees) != 1 {
		t.Errorf("expected 1 employee, got %d", len(request.Employees))
	}
	if len(request.DayOffRequests) != 1 {
		t.Errorf("expected 1 day off request, got %d", len(request.DayOffRequests))
	}
	if len(request.ShiftOffRequests) != 1 {
		t.Errorf("expected 1 shift off request, got %d", len(request.ShiftOffRequests))
	}
	if *request.Hook != "https://example.com/webhook" {
		t.Errorf("expected hook 'https://example.com/webhook', got '%s'", *request.Hook)
	}
	if len(request.ConstraintWeightOverrides.KnownConstraintNames) != 2 {
		t.Errorf("expected 2 known constraint names, got %d", len(request.ConstraintWeightOverrides.KnownConstraintNames))
	}
}

func TestCanConstructContractWithNewFields(t *testing.T) {
	contract := Contract{
		Name:                         strPtr("full-time"),
		Max:                          strPtr("PT40H"),
		Min:                          strPtr("PT32H"),
		MaxConsecutiveWorkDays:       5,
		MaxShiftsDay:                 2,
		MinRestBetweenShiftsSameDay:  strPtr("PT4H"),
		MaxWorkingDays:               5,
		LatestShiftStart:             strPtr("22:00"),
		EarliestShiftStart:           strPtr("06:00"),
		MinimumConsecutiveDaysOff:    2,
		MinimumHoursOffBetweenShifts: 11,
	}

	if contract.MinimumConsecutiveDaysOff != 2 {
		t.Errorf("expected minimumConsecutiveDaysOff 2, got %d", contract.MinimumConsecutiveDaysOff)
	}
	if contract.MinimumHoursOffBetweenShifts != 11 {
		t.Errorf("expected minimumHoursOffBetweenShifts 11, got %d", contract.MinimumHoursOffBetweenShifts)
	}
	if *contract.Name != "full-time" {
		t.Errorf("expected name 'full-time', got '%s'", *contract.Name)
	}
}

func TestCanConstructShiftAssignmentWithNewFields(t *testing.T) {
	sa := ShiftAssignment{
		Name:          strPtr("Morning"),
		From:          strPtr("2024-01-15T06:00:00"),
		To:            strPtr("2024-01-15T14:00:00"),
		Skills:        []string{"Nursing"},
		DesiredSkills: []string{"Pediatrics", "ICU"},
		Tags:          []string{"early", "weekday"},
		Cost:          200.0,
		Value:         10,
		Priority:      1,
		PinnedByUser:  true,
	}

	if len(sa.DesiredSkills) != 2 {
		t.Errorf("expected 2 desired skills, got %d", len(sa.DesiredSkills))
	}
	found := false
	for _, s := range sa.DesiredSkills {
		if s == "Pediatrics" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected desired skills to contain 'Pediatrics'")
	}
	if len(sa.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(sa.Tags))
	}
	if !sa.PinnedByUser {
		t.Errorf("expected pinnedByUser to be true")
	}
}

func TestCanConstructShiftEmployeeWithNewFields(t *testing.T) {
	employee := ShiftEmployee{
		Name:         strPtr("Bob"),
		Contract:     strPtr("part-time"),
		Skills:       []string{"Nursing"},
		LastRestDate: strPtr("2024-01-13"),
		Availability: []string{"2024-01-15"},
		Preference:   []string{"Morning"},
		PeriodRules: []PeriodRule{
			{
				Period:                              &PlanningPeriod{From: strPtr("2024-01-15"), To: strPtr("2024-01-21")},
				MaxWorkingDays:                      3,
				MinWorkingDays:                      1,
				MinWorkingDuration:                  strPtr("PT8H"),
				MaxWorkingDuration:                  strPtr("PT24H"),
				MinRestDurationBetweenShiftsSameDay: strPtr("PT4H"),
				MinRestDuration:                     strPtr("PT11H"),
			},
		},
		UnavailableDates:      []string{"2024-01-17", "2024-01-18"},
		Tags:                  []string{"junior", "trainee"},
		MaximumMinutesPerWeek: intPtr(1200),
		Shifts: []ShiftAssignment{
			{Name: strPtr("Morning"), From: strPtr("2024-01-14T06:00:00"), To: strPtr("2024-01-14T14:00:00")},
		},
	}

	if len(employee.PeriodRules) != 1 {
		t.Fatalf("expected 1 period rule, got %d", len(employee.PeriodRules))
	}
	if *employee.PeriodRules[0].Period.From != "2024-01-15" {
		t.Errorf("expected period from '2024-01-15', got '%s'", *employee.PeriodRules[0].Period.From)
	}
	if *employee.PeriodRules[0].Period.To != "2024-01-21" {
		t.Errorf("expected period to '2024-01-21', got '%s'", *employee.PeriodRules[0].Period.To)
	}
	if employee.PeriodRules[0].MaxWorkingDays != 3 {
		t.Errorf("expected maxWorkingDays 3, got %d", employee.PeriodRules[0].MaxWorkingDays)
	}
	if *employee.PeriodRules[0].MinRestDurationBetweenShiftsSameDay != "PT4H" {
		t.Errorf("expected minRestDurationBetweenShiftsSameDay 'PT4H', got '%s'", *employee.PeriodRules[0].MinRestDurationBetweenShiftsSameDay)
	}
	if len(employee.UnavailableDates) != 2 {
		t.Errorf("expected 2 unavailable dates, got %d", len(employee.UnavailableDates))
	}
	if len(employee.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(employee.Tags))
	}
	if *employee.MaximumMinutesPerWeek != 1200 {
		t.Errorf("expected maximumMinutesPerWeek 1200, got %d", *employee.MaximumMinutesPerWeek)
	}
	if len(employee.Shifts) != 1 {
		t.Errorf("expected 1 shift, got %d", len(employee.Shifts))
	}
}

func TestCanConstructWeightsWithAllFields(t *testing.T) {
	weights := Weights{
		RequiredSkills:            10,
		ShiftCapacity:             8,
		MinimumStaffing:           9,
		NoDoubleBooking:           10,
		RestBetweenShifts:         7,
		EmployeeAvailability:      8,
		MaxConsecutiveWorkDays:    6,
		MaxShiftsPerDay:           10,
		MaxWorkingDaysPerWeek:     5,
		ContractRestBetweenShifts: 7,
		EarliestShiftStart:        3,
		LatestShiftStart:          3,
		MinimumConsecutiveDaysOff: 6,
		PeriodRuleViolation:       5,
		ShiftPreferences:          2,
		CostMinimization:          4,
		WorkloadBalance:           3,
		Fairness:                  3,
		DesiredSkills:             2,
		DesiredDayOff:             4,
		ShiftOffRequest:           4,
		BalanceTimeWorked:         3,
		EmployeeAffinity:         2,
		AvoidShiftCloseToDayOff:   3,
	}

	if weights.MaxConsecutiveWorkDays != 6 {
		t.Errorf("expected maxConsecutiveWorkDays 6, got %d", weights.MaxConsecutiveWorkDays)
	}
	if weights.MaxShiftsPerDay != 10 {
		t.Errorf("expected maxShiftsPerDay 10, got %d", weights.MaxShiftsPerDay)
	}
	if weights.MaxWorkingDaysPerWeek != 5 {
		t.Errorf("expected maxWorkingDaysPerWeek 5, got %d", weights.MaxWorkingDaysPerWeek)
	}
	if weights.ContractRestBetweenShifts != 7 {
		t.Errorf("expected contractRestBetweenShifts 7, got %d", weights.ContractRestBetweenShifts)
	}
	if weights.EarliestShiftStart != 3 {
		t.Errorf("expected earliestShiftStart 3, got %d", weights.EarliestShiftStart)
	}
	if weights.LatestShiftStart != 3 {
		t.Errorf("expected latestShiftStart 3, got %d", weights.LatestShiftStart)
	}
	if weights.MinimumConsecutiveDaysOff != 6 {
		t.Errorf("expected minimumConsecutiveDaysOff 6, got %d", weights.MinimumConsecutiveDaysOff)
	}
	if weights.PeriodRuleViolation != 5 {
		t.Errorf("expected periodRuleViolation 5, got %d", weights.PeriodRuleViolation)
	}
	if weights.DesiredSkills != 2 {
		t.Errorf("expected desiredSkills 2, got %d", weights.DesiredSkills)
	}
	if weights.DesiredDayOff != 4 {
		t.Errorf("expected desiredDayOff 4, got %d", weights.DesiredDayOff)
	}
	if weights.ShiftOffRequest != 4 {
		t.Errorf("expected shiftOffRequest 4, got %d", weights.ShiftOffRequest)
	}
	if weights.BalanceTimeWorked != 3 {
		t.Errorf("expected balanceTimeWorked 3, got %d", weights.BalanceTimeWorked)
	}
	if weights.EmployeeAffinity != 2 {
		t.Errorf("expected employeeAffinity 2, got %d", weights.EmployeeAffinity)
	}
	if weights.AvoidShiftCloseToDayOff != 3 {
		t.Errorf("expected avoidShiftCloseToDayOff 3, got %d", weights.AvoidShiftCloseToDayOff)
	}
}

func TestCanSerializeShiftRequestToJSON(t *testing.T) {
	request := ShiftRequest{
		ID:   strPtr("req1"),
		Name: strPtr("Test Schedule"),
		Employees: []ShiftEmployee{
			{
				Name:     strPtr("Alice"),
				Contract: strPtr("full-time"),
				Skills:   []string{"Nursing"},
				PeriodRules: []PeriodRule{
					{
						Period:         &PlanningPeriod{From: strPtr("2024-01-15"), To: strPtr("2024-01-21")},
						MaxWorkingDays: 5,
						MinWorkingDays: 3,
					},
				},
				UnavailableDates:      []string{"2024-01-18"},
				Tags:                  []string{"senior"},
				MaximumMinutesPerWeek: intPtr(2400),
			},
		},
		Shifts: []ShiftAssignment{
			{
				Name:          strPtr("Morning"),
				DesiredSkills: []string{"Pediatrics"},
				Tags:          []string{"early"},
				PinnedByUser:  true,
			},
		},
		DayOffRequests: []DayOffRequest{
			{ID: strPtr("dor1"), EmployeeName: strPtr("Alice"), Date: strPtr("2024-01-18"), Weight: 5},
		},
		ShiftOffRequests: []ShiftOffRequest{
			{ID: strPtr("sor1"), EmployeeName: strPtr("Alice"), ShiftName: strPtr("Night"), Weight: 3},
		},
		Hook: strPtr("https://example.com/webhook"),
		ConstraintWeightOverrides: &ConstraintWeightOverrides{
			KnownConstraintNames: []string{"requiredSkills"},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	for _, expected := range []string{
		"periodRules", "unavailableDates", "maximumMinutesPerWeek",
		"desiredSkills", "pinnedByUser", "dayOffRequests", "shiftOffRequests",
		"hook", "constraintWeightOverrides", "knownConstraintNames",
	} {
		if !strings.Contains(jsonStr, expected) {
			t.Errorf("expected JSON to contain '%s'", expected)
		}
	}
}

func TestCanDeserializeShiftRequestFromJSON(t *testing.T) {
	jsonStr := `{
		"id": "req1",
		"name": "Test Schedule",
		"employees": [
			{
				"name": "Alice",
				"contract": "full-time",
				"skills": ["Nursing"],
				"periodRules": [
					{
						"period": {"from": "2024-01-15", "to": "2024-01-21"},
						"maxWorkingDays": 5,
						"minWorkingDays": 3,
						"minRestDuration": "PT11H"
					}
				],
				"unavailableDates": ["2024-01-18"],
				"tags": ["senior"],
				"maximumMinutesPerWeek": 2400,
				"shifts": [
					{"name": "Morning", "from": "2024-01-14T06:00:00", "to": "2024-01-14T14:00:00"}
				]
			}
		],
		"shifts": [
			{
				"name": "Morning",
				"from": "2024-01-15T06:00:00",
				"to": "2024-01-15T14:00:00",
				"desiredSkills": ["Pediatrics"],
				"tags": ["early"],
				"pinnedByUser": true
			}
		],
		"contracts": [
			{
				"name": "full-time",
				"minimumConsecutiveDaysOff": 2,
				"minimumHoursOffBetweenShifts": 11
			}
		],
		"dayOffRequests": [
			{"id": "dor1", "employeeName": "Alice", "date": "2024-01-18", "weight": 5}
		],
		"shiftOffRequests": [
			{"id": "sor1", "employeeName": "Alice", "shiftName": "Night", "weight": 3}
		],
		"hook": "https://example.com/webhook",
		"constraintWeightOverrides": {
			"knownConstraintNames": ["requiredSkills", "noDoubleBooking"]
		},
		"weights": {
			"requiredSkills": 10,
			"maxConsecutiveWorkDays": 6,
			"desiredSkills": 2,
			"avoidShiftCloseToDayOff": 3
		}
	}`

	var result ShiftRequest
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if *result.ID != "req1" {
		t.Errorf("expected ID 'req1', got '%s'", *result.ID)
	}

	// Employee new fields
	if len(result.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(result.Employees))
	}
	if len(result.Employees[0].PeriodRules) != 1 {
		t.Fatalf("expected 1 period rule, got %d", len(result.Employees[0].PeriodRules))
	}
	if *result.Employees[0].PeriodRules[0].Period.From != "2024-01-15" {
		t.Errorf("expected period from '2024-01-15', got '%s'", *result.Employees[0].PeriodRules[0].Period.From)
	}
	if result.Employees[0].PeriodRules[0].MaxWorkingDays != 5 {
		t.Errorf("expected maxWorkingDays 5, got %d", result.Employees[0].PeriodRules[0].MaxWorkingDays)
	}
	if *result.Employees[0].PeriodRules[0].MinRestDuration != "PT11H" {
		t.Errorf("expected minRestDuration 'PT11H', got '%s'", *result.Employees[0].PeriodRules[0].MinRestDuration)
	}
	if len(result.Employees[0].UnavailableDates) != 1 {
		t.Errorf("expected 1 unavailable date, got %d", len(result.Employees[0].UnavailableDates))
	}
	if result.Employees[0].Tags[0] != "senior" {
		t.Errorf("expected tag 'senior', got '%s'", result.Employees[0].Tags[0])
	}
	if *result.Employees[0].MaximumMinutesPerWeek != 2400 {
		t.Errorf("expected maximumMinutesPerWeek 2400, got %d", *result.Employees[0].MaximumMinutesPerWeek)
	}
	if len(result.Employees[0].Shifts) != 1 {
		t.Errorf("expected 1 employee shift, got %d", len(result.Employees[0].Shifts))
	}

	// ShiftAssignment new fields
	if len(result.Shifts[0].DesiredSkills) != 1 {
		t.Fatalf("expected 1 desired skill, got %d", len(result.Shifts[0].DesiredSkills))
	}
	if result.Shifts[0].DesiredSkills[0] != "Pediatrics" {
		t.Errorf("expected desired skill 'Pediatrics', got '%s'", result.Shifts[0].DesiredSkills[0])
	}
	if len(result.Shifts[0].Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result.Shifts[0].Tags))
	}
	if !result.Shifts[0].PinnedByUser {
		t.Errorf("expected pinnedByUser to be true")
	}

	// Contract new fields
	if result.Contracts[0].MinimumConsecutiveDaysOff != 2 {
		t.Errorf("expected minimumConsecutiveDaysOff 2, got %d", result.Contracts[0].MinimumConsecutiveDaysOff)
	}
	if result.Contracts[0].MinimumHoursOffBetweenShifts != 11 {
		t.Errorf("expected minimumHoursOffBetweenShifts 11, got %d", result.Contracts[0].MinimumHoursOffBetweenShifts)
	}

	// DayOffRequests
	if len(result.DayOffRequests) != 1 {
		t.Fatalf("expected 1 day off request, got %d", len(result.DayOffRequests))
	}
	if *result.DayOffRequests[0].ID != "dor1" {
		t.Errorf("expected day off request ID 'dor1', got '%s'", *result.DayOffRequests[0].ID)
	}
	if *result.DayOffRequests[0].EmployeeName != "Alice" {
		t.Errorf("expected employee name 'Alice', got '%s'", *result.DayOffRequests[0].EmployeeName)
	}
	if result.DayOffRequests[0].Weight != 5 {
		t.Errorf("expected weight 5, got %d", result.DayOffRequests[0].Weight)
	}

	// ShiftOffRequests
	if len(result.ShiftOffRequests) != 1 {
		t.Fatalf("expected 1 shift off request, got %d", len(result.ShiftOffRequests))
	}
	if *result.ShiftOffRequests[0].ID != "sor1" {
		t.Errorf("expected shift off request ID 'sor1', got '%s'", *result.ShiftOffRequests[0].ID)
	}
	if *result.ShiftOffRequests[0].ShiftName != "Night" {
		t.Errorf("expected shift name 'Night', got '%s'", *result.ShiftOffRequests[0].ShiftName)
	}

	// Hook and overrides
	if *result.Hook != "https://example.com/webhook" {
		t.Errorf("expected hook 'https://example.com/webhook', got '%s'", *result.Hook)
	}
	if len(result.ConstraintWeightOverrides.KnownConstraintNames) != 2 {
		t.Errorf("expected 2 known constraint names, got %d", len(result.ConstraintWeightOverrides.KnownConstraintNames))
	}

	// Weights new fields
	if result.Weights.RequiredSkills != 10 {
		t.Errorf("expected requiredSkills 10, got %d", result.Weights.RequiredSkills)
	}
	if result.Weights.MaxConsecutiveWorkDays != 6 {
		t.Errorf("expected maxConsecutiveWorkDays 6, got %d", result.Weights.MaxConsecutiveWorkDays)
	}
	if result.Weights.DesiredSkills != 2 {
		t.Errorf("expected desiredSkills 2, got %d", result.Weights.DesiredSkills)
	}
	if result.Weights.AvoidShiftCloseToDayOff != 3 {
		t.Errorf("expected avoidShiftCloseToDayOff 3, got %d", result.Weights.AvoidShiftCloseToDayOff)
	}
}

func TestCanDeserializeShiftResultResponseFromJSON(t *testing.T) {
	jsonStr := `{
		"feasible": true,
		"scoreString": "0hard/-120soft",
		"score": {"hardScore": 0, "softScore": -120},
		"assignedShifts": [
			{
				"name": "Morning",
				"from": "2024-01-15T06:00:00",
				"to": "2024-01-15T14:00:00",
				"skills": ["Nursing"],
				"desiredSkills": ["Pediatrics"],
				"tags": ["early"],
				"cost": 150.0,
				"value": 10,
				"priority": 1,
				"pinnedByUser": false
			}
		],
		"unassignedShifts": [
			{
				"name": "Night",
				"from": "2024-01-15T22:00:00",
				"to": "2024-01-16T06:00:00",
				"cost": 200.0,
				"value": 5,
				"priority": 2,
				"pinnedByUser": false
			}
		],
		"employees": [
			{
				"name": "Alice",
				"contract": "full-time",
				"skills": ["Nursing", "Pediatrics"],
				"shifts": [
					{"name": "Morning", "from": "2024-01-15T06:00:00", "to": "2024-01-15T14:00:00"}
				]
			}
		]
	}`

	var result ShiftResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible to be true")
	}
	if result.ScoreString == nil || *result.ScoreString != "0hard/-120soft" {
		t.Errorf("expected scoreString '0hard/-120soft', got '%v'", result.ScoreString)
	}
	if result.Score["hardScore"] != float64(0) {
		t.Errorf("expected score.hardScore 0, got %v", result.Score["hardScore"])
	}

	if len(result.AssignedShifts) != 1 {
		t.Fatalf("expected 1 assigned shift, got %d", len(result.AssignedShifts))
	}
	if *result.AssignedShifts[0].Name != "Morning" {
		t.Errorf("expected assigned shift name 'Morning', got '%s'", *result.AssignedShifts[0].Name)
	}
	if *result.AssignedShifts[0].From != "2024-01-15T06:00:00" {
		t.Errorf("expected assigned shift from '2024-01-15T06:00:00', got '%s'", *result.AssignedShifts[0].From)
	}
	if result.AssignedShifts[0].Cost != 150.0 {
		t.Errorf("expected assigned shift cost 150.0, got %f", result.AssignedShifts[0].Cost)
	}

	if len(result.UnassignedShifts) != 1 {
		t.Fatalf("expected 1 unassigned shift, got %d", len(result.UnassignedShifts))
	}
	if *result.UnassignedShifts[0].Name != "Night" {
		t.Errorf("expected unassigned shift name 'Night', got '%s'", *result.UnassignedShifts[0].Name)
	}

	if len(result.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(result.Employees))
	}
	if *result.Employees[0].Name != "Alice" {
		t.Errorf("expected employee name 'Alice', got '%s'", *result.Employees[0].Name)
	}
	if len(result.Employees[0].Shifts) != 1 {
		t.Errorf("expected 1 employee shift, got %d", len(result.Employees[0].Shifts))
	}
}

func TestShiftResultResponseEchoesRequestFields(t *testing.T) {
	jsonStr := `{
		"id": "req1",
		"name": "Weekly Schedule",
		"description": "Week 3",
		"contracts": [
			{"name": "full-time", "minimumConsecutiveDaysOff": 2, "minimumHoursOffBetweenShifts": 11}
		],
		"shifts": [
			{"name": "Morning", "from": "2024-01-15T06:00:00", "to": "2024-01-15T14:00:00"}
		],
		"options": {"partialPlanning": true, "timeLimit": 60},
		"weights": {"requiredSkills": 10},
		"hook": "https://example.com/webhook",
		"feasible": true,
		"scoreString": "0hard/-120soft",
		"score": {"hardScore": 0, "softScore": -120},
		"assignedShifts": [
			{"name": "Morning", "from": "2024-01-15T06:00:00", "to": "2024-01-15T14:00:00"}
		],
		"unassignedShifts": [],
		"employees": [
			{"name": "Alice", "contract": "full-time", "skills": ["Nursing"], "shifts": []}
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
	if len(result.Contracts) != 1 || result.Contracts[0].MinimumConsecutiveDaysOff != 2 {
		t.Errorf("expected echoed contracts, got %v", result.Contracts)
	}
	if len(result.Shifts) != 1 {
		t.Errorf("expected echoed shifts, got %d", len(result.Shifts))
	}
	if result.Options == nil || result.Options.PartialPlanning == nil || !*result.Options.PartialPlanning {
		t.Errorf("expected echoed options.partialPlanning true, got %v", result.Options)
	}
	if result.Weights == nil || result.Weights.RequiredSkills != 10 {
		t.Errorf("expected echoed weights.requiredSkills 10, got %v", result.Weights)
	}
	if result.Hook == nil || *result.Hook != "https://example.com/webhook" {
		t.Errorf("expected echoed hook, got %v", result.Hook)
	}

	// Result-only fields.
	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible true")
	}
	if len(result.AssignedShifts) != 1 {
		t.Errorf("expected 1 assigned shift, got %d", len(result.AssignedShifts))
	}
}
