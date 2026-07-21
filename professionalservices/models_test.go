package professionalservices

import (
	"encoding/json"
	"strings"
	"testing"
)

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestCanConstructProfessionalServicesRequest(t *testing.T) {
	request := ProfessionalServicesRequest{
		Employees: []Employee{
			{
				ID:     "emp1",
				Skills: []string{"C#", "JavaScript"},
				Shifts: []Shift{
					{ID: "shift1", MinStartTime: "2024-01-15T08:00:00", MaxEndTime: "2024-01-15T18:00:00"},
				},
			},
		},
		Tasks: []Task{
			{
				ID:             "task1",
				Name:           "Develop API",
				Deadline:       strPtr("2024-01-20T17:00:00"),
				Duration:       "PT8H",
				Priority:       "HIGH",
				RequiredSkills: []string{"C#", "JavaScript"},
			},
		},
	}

	if len(request.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(request.Employees))
	}
	if len(request.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(request.Tasks))
	}
	if request.Employees[0].ID != "emp1" {
		t.Errorf("expected employee ID 'emp1', got '%s'", request.Employees[0].ID)
	}
	if request.Tasks[0].ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", request.Tasks[0].ID)
	}
}

func TestCanSerializeToJSON(t *testing.T) {
	request := ProfessionalServicesRequest{
		Employees: []Employee{},
		Tasks:     []Task{},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, "employees") {
		t.Errorf("expected JSON to contain 'employees'")
	}
	if !strings.Contains(jsonStr, "tasks") {
		t.Errorf("expected JSON to contain 'tasks'")
	}
}

func TestCanDeserializeResultResponse(t *testing.T) {
	jsonStr := `{
		"employees": [
			{
				"id": "emp1",
				"skills": ["C#", "JavaScript"],
				"shifts": [],
				"tasks": ["task1", "task2"]
			}
		],
		"tasks": [
			{
				"id": "task1",
				"name": "Develop API",
				"deadline": "2024-01-20T17:00:00",
				"duration": 28800.0,
				"priority": "HIGH",
				"requiredSkills": ["C#", "JavaScript"],
				"startTime": "2024-01-15T14:00:00",
				"endTime": "2024-01-15T18:00:00"
			}
		]
	}`

	var result ProfessionalServicesResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(result.Employees))
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(result.Tasks))
	}
	if result.Employees[0].ID != "emp1" {
		t.Errorf("expected employee ID 'emp1', got '%s'", result.Employees[0].ID)
	}
	if result.Tasks[0].ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", result.Tasks[0].ID)
	}
	if len(result.Employees[0].Tasks) != 2 {
		t.Errorf("expected 2 tasks on employee, got %d", len(result.Employees[0].Tasks))
	}
	if result.Tasks[0].StartTime != "2024-01-15T14:00:00" {
		t.Errorf("expected startTime '2024-01-15T14:00:00', got '%s'", result.Tasks[0].StartTime)
	}
	if result.Tasks[0].EndTime != "2024-01-15T18:00:00" {
		t.Errorf("expected endTime '2024-01-15T18:00:00', got '%s'", result.Tasks[0].EndTime)
	}
}

func TestCanConstructEmployeeWithShifts(t *testing.T) {
	employee := Employee{
		ID:     "emp1",
		Skills: []string{"C#", "JavaScript", "SQL"},
		Shifts: []Shift{
			{ID: "morning", MinStartTime: "2024-01-15T08:00:00", MaxEndTime: "2024-01-15T12:00:00"},
			{ID: "afternoon", MinStartTime: "2024-01-15T13:00:00", MaxEndTime: "2024-01-15T17:00:00"},
		},
	}

	if employee.ID != "emp1" {
		t.Errorf("expected ID 'emp1', got '%s'", employee.ID)
	}
	if len(employee.Skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(employee.Skills))
	}
	if len(employee.Shifts) != 2 {
		t.Errorf("expected 2 shifts, got %d", len(employee.Shifts))
	}
	if employee.Shifts[0].ID != "morning" {
		t.Errorf("expected first shift ID 'morning', got '%s'", employee.Shifts[0].ID)
	}
	if employee.Shifts[1].ID != "afternoon" {
		t.Errorf("expected second shift ID 'afternoon', got '%s'", employee.Shifts[1].ID)
	}
}

func TestCanConstructTaskWithRequiredSkills(t *testing.T) {
	task := Task{
		ID:             "task1",
		Name:           "Database Design",
		Deadline:       strPtr("2024-01-25T17:00:00"),
		Duration:       "PT16H",
		Priority:       "MEDIUM",
		RequiredSkills: []string{"SQL", "Database Design", "ERD"},
	}

	if task.ID != "task1" {
		t.Errorf("expected ID 'task1', got '%s'", task.ID)
	}
	if task.Name != "Database Design" {
		t.Errorf("expected name 'Database Design', got '%s'", task.Name)
	}
	if *task.Deadline != "2024-01-25T17:00:00" {
		t.Errorf("expected deadline '2024-01-25T17:00:00', got '%s'", *task.Deadline)
	}
	if task.Duration != "PT16H" {
		t.Errorf("expected duration 'PT16H', got '%s'", task.Duration)
	}
	if task.Priority != "MEDIUM" {
		t.Errorf("expected priority 'MEDIUM', got '%s'", task.Priority)
	}
	if len(task.RequiredSkills) != 3 {
		t.Errorf("expected 3 required skills, got %d", len(task.RequiredSkills))
	}
	found := false
	for _, s := range task.RequiredSkills {
		if s == "SQL" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected required skills to contain 'SQL'")
	}
}

func TestCanConstructScheduledEmployee(t *testing.T) {
	scheduledEmployee := ScheduledEmployee{
		ID:     "emp1",
		Skills: []string{"C#"},
		Shifts: []Shift{},
		Tasks:  []string{"task1", "task2", "task3"},
	}

	if scheduledEmployee.ID != "emp1" {
		t.Errorf("expected ID 'emp1', got '%s'", scheduledEmployee.ID)
	}
	if len(scheduledEmployee.Tasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(scheduledEmployee.Tasks))
	}
	for _, expected := range []string{"task1", "task2", "task3"} {
		found := false
		for _, task := range scheduledEmployee.Tasks {
			if task == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected tasks to contain '%s'", expected)
		}
	}
}

func TestCanConstructScheduledTask(t *testing.T) {
	scheduledTask := ScheduledTask{
		ID:             "task1",
		Name:           "Code Review",
		Deadline:       strPtr("2024-01-22T17:00:00"),
		Duration:       500.0,
		Priority:       "HIGH",
		RequiredSkills: []string{"C#", "Code Review"},
		StartTime:      "2024-01-15T14:00:00",
		EndTime:        "2024-01-15T18:00:00",
	}

	if scheduledTask.ID != "task1" {
		t.Errorf("expected ID 'task1', got '%s'", scheduledTask.ID)
	}
	if scheduledTask.Name != "Code Review" {
		t.Errorf("expected name 'Code Review', got '%s'", scheduledTask.Name)
	}
	if scheduledTask.StartTime != "2024-01-15T14:00:00" {
		t.Errorf("expected startTime '2024-01-15T14:00:00', got '%s'", scheduledTask.StartTime)
	}
	if scheduledTask.EndTime != "2024-01-15T18:00:00" {
		t.Errorf("expected endTime '2024-01-15T18:00:00', got '%s'", scheduledTask.EndTime)
	}
	if scheduledTask.Priority != "HIGH" {
		t.Errorf("expected priority 'HIGH', got '%s'", scheduledTask.Priority)
	}
}

func TestCanSerializeComplexRequest(t *testing.T) {
	request := ProfessionalServicesRequest{
		Employees: []Employee{
			{
				ID:     "emp1",
				Skills: []string{"C#", "JavaScript"},
				Shifts: []Shift{
					{ID: "shift1", MinStartTime: "2024-01-15T08:00:00", MaxEndTime: "2024-01-15T18:00:00"},
				},
			},
			{
				ID:     "emp2",
				Skills: []string{"Python", "SQL"},
				Shifts: []Shift{
					{ID: "shift2", MinStartTime: "2024-01-15T09:00:00", MaxEndTime: "2024-01-15T17:00:00"},
				},
			},
		},
		Tasks: []Task{
			{
				ID:             "task1",
				Name:           "Frontend Development",
				Deadline:       strPtr("2024-01-20T17:00:00"),
				Duration:       "PT16H",
				Priority:       "HIGH",
				RequiredSkills: []string{"JavaScript", "React"},
			},
			{
				ID:             "task2",
				Name:           "Backend API",
				Deadline:       strPtr("2024-01-22T17:00:00"),
				Duration:       "PT24H",
				Priority:       "MEDIUM",
				RequiredSkills: []string{"C#", "ASP.NET"},
			},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	for _, expected := range []string{"emp1", "emp2", "task1", "task2", "Frontend Development", "Backend API", "C#", "JavaScript", "Python", "SQL"} {
		if !strings.Contains(jsonStr, expected) {
			t.Errorf("expected JSON to contain '%s'", expected)
		}
	}
}

func TestCanConstructEmployeeWithNewFields(t *testing.T) {
	employee := Employee{
		ID:                "emp1",
		Name:              strPtr("Jane Smith"),
		Email:             strPtr("jane@example.com"),
		Skills:            []string{"C#", "JavaScript"},
		HourlyRate:        float64Ptr(95.50),
		ContractID:        strPtr("contract1"),
		TimeZoneID:        strPtr("Europe/Brussels"),
		DedicatedClientID: strPtr("client1"),
		Shifts: []Shift{
			{ID: "shift1", MinStartTime: "2024-01-15T08:00:00", MaxEndTime: "2024-01-15T18:00:00"},
		},
		AvailabilityTimeSpans: []AvailabilityTimeSpan{
			{ID: strPtr("avail1"), Start: strPtr("2024-01-15T08:00:00"), End: strPtr("2024-01-15T12:00:00"), Type: strPtr("AVAILABLE")},
			{ID: strPtr("avail2"), Start: strPtr("2024-01-15T13:00:00"), End: strPtr("2024-01-15T17:00:00"), Type: strPtr("UNAVAILABLE")},
		},
	}

	if employee.ID != "emp1" {
		t.Errorf("expected ID 'emp1', got '%s'", employee.ID)
	}
	if *employee.Name != "Jane Smith" {
		t.Errorf("expected name 'Jane Smith', got '%s'", *employee.Name)
	}
	if *employee.Email != "jane@example.com" {
		t.Errorf("expected email 'jane@example.com', got '%s'", *employee.Email)
	}
	if *employee.HourlyRate != 95.50 {
		t.Errorf("expected hourlyRate 95.50, got %f", *employee.HourlyRate)
	}
	if *employee.ContractID != "contract1" {
		t.Errorf("expected contractId 'contract1', got '%s'", *employee.ContractID)
	}
	if *employee.TimeZoneID != "Europe/Brussels" {
		t.Errorf("expected timeZoneId 'Europe/Brussels', got '%s'", *employee.TimeZoneID)
	}
	if *employee.DedicatedClientID != "client1" {
		t.Errorf("expected dedicatedClientId 'client1', got '%s'", *employee.DedicatedClientID)
	}
	if len(employee.AvailabilityTimeSpans) != 2 {
		t.Fatalf("expected 2 availability time spans, got %d", len(employee.AvailabilityTimeSpans))
	}
	if *employee.AvailabilityTimeSpans[0].ID != "avail1" {
		t.Errorf("expected first availability ID 'avail1', got '%s'", *employee.AvailabilityTimeSpans[0].ID)
	}
	if *employee.AvailabilityTimeSpans[0].Type != "AVAILABLE" {
		t.Errorf("expected first availability type 'AVAILABLE', got '%s'", *employee.AvailabilityTimeSpans[0].Type)
	}
}

func TestCanConstructTaskWithNewFields(t *testing.T) {
	task := Task{
		ID:                  "task1",
		Name:                "Full Stack Feature",
		Duration:            "PT16H",
		Priority:            "HIGH",
		RequiredSkills:      []string{"C#", "React"},
		PreferredSkills:     []string{"Docker", "Kubernetes"},
		TimeZoneID:          strPtr("America/New_York"),
		DependsOn:           []string{"task0"},
		PreferredEmployees:  []string{"emp1", "emp2"},
		ProhibitedEmployees: []string{"emp3"},
		ClientID:            strPtr("client1"),
		ProjectID:           strPtr("project1"),
	}

	if task.ID != "task1" {
		t.Errorf("expected ID 'task1', got '%s'", task.ID)
	}
	if len(task.PreferredSkills) != 2 {
		t.Errorf("expected 2 preferred skills, got %d", len(task.PreferredSkills))
	}
	found := false
	for _, s := range task.PreferredSkills {
		if s == "Docker" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected preferred skills to contain 'Docker'")
	}
	if *task.TimeZoneID != "America/New_York" {
		t.Errorf("expected timeZoneId 'America/New_York', got '%s'", *task.TimeZoneID)
	}
	if len(task.DependsOn) != 1 {
		t.Fatalf("expected 1 dependency, got %d", len(task.DependsOn))
	}
	if task.DependsOn[0] != "task0" {
		t.Errorf("expected dependency 'task0', got '%s'", task.DependsOn[0])
	}
	if len(task.PreferredEmployees) != 2 {
		t.Errorf("expected 2 preferred employees, got %d", len(task.PreferredEmployees))
	}
	if len(task.ProhibitedEmployees) != 1 {
		t.Fatalf("expected 1 prohibited employee, got %d", len(task.ProhibitedEmployees))
	}
	if task.ProhibitedEmployees[0] != "emp3" {
		t.Errorf("expected prohibited employee 'emp3', got '%s'", task.ProhibitedEmployees[0])
	}
}

func TestCanConstructRequestWithContracts(t *testing.T) {
	request := ProfessionalServicesRequest{
		Name:        strPtr("Q1 Planning"),
		Description: strPtr("First quarter resource planning"),
		StartDate:   strPtr("2024-01-01"),
		EndDate:     strPtr("2024-03-31"),
		Employees:   []Employee{},
		Tasks:       []Task{},
		Contracts: []Contract{
			{
				ID:                   strPtr("full-time"),
				Name:                 strPtr("Full Time"),
				MaxHoursPerDay:       strPtr("PT8H"),
				MaxHoursPerWeek:      strPtr("PT40H"),
				MinRestBetweenShifts: strPtr("PT12H"),
				TargetUtilization:    float64Ptr(0.85),
			},
			{
				ID:                strPtr("part-time"),
				Name:              strPtr("Part Time"),
				MaxHoursPerDay:    strPtr("PT4H"),
				MaxHoursPerWeek:   strPtr("PT20H"),
				TargetUtilization: float64Ptr(0.70),
			},
		},
	}

	if len(request.Contracts) != 2 {
		t.Fatalf("expected 2 contracts, got %d", len(request.Contracts))
	}
	if *request.Contracts[0].ID != "full-time" {
		t.Errorf("expected first contract ID 'full-time', got '%s'", *request.Contracts[0].ID)
	}
	if *request.Contracts[0].MaxHoursPerDay != "PT8H" {
		t.Errorf("expected maxHoursPerDay 'PT8H', got '%s'", *request.Contracts[0].MaxHoursPerDay)
	}
	if *request.Contracts[0].MaxHoursPerWeek != "PT40H" {
		t.Errorf("expected maxHoursPerWeek 'PT40H', got '%s'", *request.Contracts[0].MaxHoursPerWeek)
	}
	if *request.Contracts[0].MinRestBetweenShifts != "PT12H" {
		t.Errorf("expected minRestBetweenShifts 'PT12H', got '%s'", *request.Contracts[0].MinRestBetweenShifts)
	}
	if *request.Contracts[0].TargetUtilization != 0.85 {
		t.Errorf("expected targetUtilization 0.85, got %f", *request.Contracts[0].TargetUtilization)
	}
	if *request.Contracts[1].ID != "part-time" {
		t.Errorf("expected second contract ID 'part-time', got '%s'", *request.Contracts[1].ID)
	}
}

func TestCanSerializeNewFieldsToJSON(t *testing.T) {
	request := ProfessionalServicesRequest{
		Employees: []Employee{
			{
				ID:                "emp1",
				Skills:            []string{"C#"},
				HourlyRate:        float64Ptr(100.0),
				ContractID:        strPtr("full-time"),
				TimeZoneID:        strPtr("Europe/London"),
				DedicatedClientID: strPtr("client1"),
				AvailabilityTimeSpans: []AvailabilityTimeSpan{
					{ID: strPtr("a1"), Start: strPtr("2024-01-15T08:00:00"), End: strPtr("2024-01-15T17:00:00"), Type: strPtr("AVAILABLE")},
				},
			},
		},
		Tasks: []Task{
			{
				ID:                  "task1",
				Name:                "API Work",
				Duration:            "PT8H",
				Priority:            "HIGH",
				RequiredSkills:      []string{"C#"},
				PreferredSkills:     []string{"Azure"},
				TimeZoneID:          strPtr("Europe/London"),
				DependsOn:           []string{"task0"},
				PreferredEmployees:  []string{"emp1"},
				ProhibitedEmployees: []string{"emp2"},
			},
		},
		Contracts: []Contract{
			{ID: strPtr("full-time"), Name: strPtr("Full Time"), MaxHoursPerDay: strPtr("PT8H")},
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	for _, expected := range []string{
		"hourlyRate", "contractId", "timeZoneId", "dedicatedClientId",
		"availabilityTimeSpans", "preferredSkills", "dependsOn",
		"preferredEmployees", "prohibitedEmployees", "contracts", "maxHoursPerDay",
	} {
		if !strings.Contains(jsonStr, expected) {
			t.Errorf("expected JSON to contain '%s'", expected)
		}
	}
}

func TestCanDeserializeNewFieldsFromJSON(t *testing.T) {
	jsonStr := `{
		"employees": [
			{
				"id": "emp1",
				"skills": ["C#"],
				"shifts": [],
				"hourlyRate": 100.0,
				"contractId": "full-time",
				"timeZoneId": "Europe/London",
				"dedicatedClientId": "client1",
				"availabilityTimeSpans": [
					{"id": "a1", "start": "2024-01-15T08:00:00", "end": "2024-01-15T17:00:00", "type": "AVAILABLE"}
				]
			}
		],
		"tasks": [
			{
				"id": "task1",
				"name": "API Work",
				"duration": 28800.0,
				"priority": "HIGH",
				"requiredSkills": ["C#"],
				"startTime": "2024-01-15T14:00:00",
				"endTime": "2024-01-15T18:00:00"
			}
		]
	}`

	var result ProfessionalServicesResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(result.Employees))
	}
	if *result.Employees[0].HourlyRate != 100.0 {
		t.Errorf("expected hourlyRate 100.0, got %f", *result.Employees[0].HourlyRate)
	}
	if *result.Employees[0].ContractID != "full-time" {
		t.Errorf("expected contractId 'full-time', got '%s'", *result.Employees[0].ContractID)
	}
	if *result.Employees[0].TimeZoneID != "Europe/London" {
		t.Errorf("expected timeZoneId 'Europe/London', got '%s'", *result.Employees[0].TimeZoneID)
	}
	if *result.Employees[0].DedicatedClientID != "client1" {
		t.Errorf("expected dedicatedClientId 'client1', got '%s'", *result.Employees[0].DedicatedClientID)
	}
	if len(result.Employees[0].AvailabilityTimeSpans) != 1 {
		t.Fatalf("expected 1 availability time span, got %d", len(result.Employees[0].AvailabilityTimeSpans))
	}
	if *result.Employees[0].AvailabilityTimeSpans[0].Type != "AVAILABLE" {
		t.Errorf("expected type 'AVAILABLE', got '%s'", *result.Employees[0].AvailabilityTimeSpans[0].Type)
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(result.Tasks))
	}
	if result.Tasks[0].ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", result.Tasks[0].ID)
	}
}
