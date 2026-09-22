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

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
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
		"id": "job1",
		"name": "Q1 Planning",
		"employees": [
			{
				"id": "emp1",
				"skills": ["C#", "JavaScript"],
				"shifts": []
			}
		],
		"tasks": [
			{
				"id": "task1",
				"name": "Develop API",
				"deadline": "2024-01-20T17:00:00",
				"duration": "PT8H",
				"priority": "HIGH",
				"requiredSkills": ["C#", "JavaScript"]
			}
		],
		"solverStatus": "NOT_SOLVING",
		"feasible": true,
		"scoreString": "0hard/-10soft",
		"score": {"hardScore": 0, "softScore": -10},
		"assignedTasks": ["task1"],
		"unassignedTasks": []
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
	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible to be true")
	}
	if result.ScoreString == nil || *result.ScoreString != "0hard/-10soft" {
		t.Errorf("expected scoreString '0hard/-10soft', got %v", result.ScoreString)
	}
	if result.Score["hardScore"] != float64(0) {
		t.Errorf("expected score.hardScore 0, got %v", result.Score["hardScore"])
	}
	if len(result.AssignedTasks) != 1 || result.AssignedTasks[0] != "task1" {
		t.Errorf("expected assignedTasks ['task1'], got %v", result.AssignedTasks)
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

func TestCanConstructResultResponseEcho(t *testing.T) {
	result := ProfessionalServicesResultResponse{
		ID:   strPtr("job1"),
		Name: strPtr("Q1 Planning"),
		Employees: []Employee{
			{ID: "emp1", Skills: []string{"C#"}, Shifts: []Shift{}},
		},
		Tasks: []Task{
			{ID: "task1", Name: "Code Review", Duration: "PT8H", Priority: "HIGH", RequiredSkills: []string{"C#"}},
		},
		Feasible:        boolPtr(true),
		ScoreString:     strPtr("0hard/-10soft"),
		AssignedTasks:   []string{"task1"},
		UnassignedTasks: []string{},
	}

	if *result.ID != "job1" {
		t.Errorf("expected ID 'job1', got '%s'", *result.ID)
	}
	if len(result.Employees) != 1 {
		t.Errorf("expected 1 employee, got %d", len(result.Employees))
	}
	if len(result.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(result.Tasks))
	}
	if len(result.AssignedTasks) != 1 || result.AssignedTasks[0] != "task1" {
		t.Errorf("expected assignedTasks ['task1'], got %v", result.AssignedTasks)
	}
	if result.Feasible == nil || !*result.Feasible {
		t.Errorf("expected feasible to be true")
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
		FreezeUntil:       strPtr("2024-01-16T00:00:00"),
		SkillValidity:     []SkillValidity{{Skill: "Audit", ValidFrom: strPtr("2024-01-01T00:00:00"), ValidTo: strPtr("2025-01-01T00:00:00")}},
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
	if *employee.FreezeUntil != "2024-01-16T00:00:00" {
		t.Errorf("expected freezeUntil '2024-01-16T00:00:00', got '%s'", *employee.FreezeUntil)
	}
	if len(employee.SkillValidity) != 1 || employee.SkillValidity[0].Skill != "Audit" {
		t.Errorf("expected one skill validity window for 'Audit', got %v", employee.SkillValidity)
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
		PreferredSkills: []string{"Docker", "Kubernetes"},
		DependsOn: []TaskDependency{
			{TaskID: "task0"},
			{TaskID: "review", Anchor: strPtr("END"), MinOffset: strPtr("PT12H"), MaxOffset: strPtr("P3D"), RequiresSameEmployee: true, NoIntermediateTasks: true},
		},
		PreferredEmployees:  []string{"emp1", "emp2"},
		ProhibitedEmployees: []string{"emp3"},
		AllowedEmployees:    []string{"emp1", "emp2", "emp4"},
		ClientID:            strPtr("client1"),
		ProjectID:           strPtr("project1"),
		SLA:                 strPtr("2024-01-22T17:00:00"),
		EarliestStart:       strPtr("2024-01-15T09:00:00"),
		TaskType:            strPtr("deep-work"),
		PriorityValue:       intPtr(120),
		DurationByEmployee:  map[string]string{"emp1": "PT8H"},
		Pinned:              true,
		ParentTaskID:        strPtr("feature"),
		SegmentIndex:        intPtr(1),
		TotalSegments:       intPtr(2),
		InitialEmployeeID:   strPtr("emp1"),
		InitialStartTime:    strPtr("2024-01-15T09:00:00"),
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
	if len(task.DependsOn) != 2 {
		t.Fatalf("expected 2 dependencies, got %d", len(task.DependsOn))
	}
	if task.DependsOn[0].TaskID != "task0" || task.DependsOn[0].Anchor != nil {
		t.Errorf("expected a plain dependency on 'task0', got %+v", task.DependsOn[0])
	}
	if *task.DependsOn[1].MinOffset != "PT12H" || !task.DependsOn[1].RequiresSameEmployee {
		t.Errorf("expected a rich dependency on 'review', got %+v", task.DependsOn[1])
	}
	if len(task.AllowedEmployees) != 3 {
		t.Errorf("expected 3 allowed employees, got %d", len(task.AllowedEmployees))
	}
	if *task.SLA != "2024-01-22T17:00:00" || *task.EarliestStart != "2024-01-15T09:00:00" || *task.TaskType != "deep-work" {
		t.Errorf("expected sla/earliestStart/taskType to round-trip, got %+v", task)
	}
	if *task.PriorityValue != 120 || task.DurationByEmployee["emp1"] != "PT8H" {
		t.Errorf("expected priorityValue 120 and a duration override for emp1, got %+v", task)
	}
	if !task.Pinned || *task.ParentTaskID != "feature" || *task.SegmentIndex != 1 || *task.TotalSegments != 2 {
		t.Errorf("expected pinning and split-task fields to round-trip, got %+v", task)
	}
	if *task.InitialEmployeeID != "emp1" || *task.InitialStartTime != "2024-01-15T09:00:00" {
		t.Errorf("expected re-planning fields to round-trip, got %+v", task)
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
			},
			{
				ID:              strPtr("part-time"),
				Name:            strPtr("Part Time"),
				MaxHoursPerDay:  strPtr("PT4H"),
				MaxHoursPerWeek: strPtr("PT20H"),
			},
		},
		TaskTypes: []string{"deep-work", "meeting"},
		TaskTypeTransitions: []TaskTypeTransition{
			{FromTaskType: "deep-work", ToTaskType: "meeting", SetupDuration: strPtr("PT15M")},
			{FromTaskType: "*", ToTaskType: "deep-work", Forbidden: true},
		},
		FreezeUntil: strPtr("2024-01-10T00:00:00"),
		Weights:     map[string]string{"meetDeadlines": "0hard/0medium/5soft"},
	}

	if len(request.TaskTypes) != 2 {
		t.Errorf("expected 2 task types, got %d", len(request.TaskTypes))
	}
	if *request.TaskTypeTransitions[0].SetupDuration != "PT15M" || !request.TaskTypeTransitions[1].Forbidden {
		t.Errorf("expected task type transitions to round-trip, got %+v", request.TaskTypeTransitions)
	}
	if *request.FreezeUntil != "2024-01-10T00:00:00" {
		t.Errorf("expected freezeUntil to round-trip, got %v", request.FreezeUntil)
	}
	if request.Weights["meetDeadlines"] != "0hard/0medium/5soft" {
		t.Errorf("expected weights to round-trip, got %v", request.Weights)
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
				ContractID:    strPtr("full-time"),
				FreezeUntil:   strPtr("2024-01-16T00:00:00"),
				SkillValidity: []SkillValidity{{Skill: "Audit", ValidTo: strPtr("2025-01-01T00:00:00")}},
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
				DependsOn:           []TaskDependency{{TaskID: "task0"}},
				PreferredEmployees:  []string{"emp1"},
				ProhibitedEmployees: []string{"emp2"},
				AllowedEmployees:    []string{"emp1"},
				SLA:                 strPtr("2024-01-20T17:00:00"),
				EarliestStart:       strPtr("2024-01-15T09:00:00"),
				TaskType:            strPtr("deep-work"),
				PriorityValue:       intPtr(80),
				DurationByEmployee:  map[string]string{"emp1": "PT6H"},
				Pinned:              true,
				InitialEmployeeID:   strPtr("emp1"),
				InitialStartTime:    strPtr("2024-01-15T09:00:00"),
			},
		},
		Contracts: []Contract{
			{ID: strPtr("full-time"), Name: strPtr("Full Time"), MaxHoursPerDay: strPtr("PT8H")},
		},
		TaskTypeTransitions: []TaskTypeTransition{{FromTaskType: "*", ToTaskType: "deep-work", SetupDuration: strPtr("PT30M")}},
		TaskTypes:           []string{"deep-work"},
		FreezeUntil:         strPtr("2024-01-10T00:00:00"),
		Weights:             map[string]string{"minimizeCost": "0hard/0medium/2soft"},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	jsonStr := string(data)
	for _, expected := range []string{
		"hourlyRate", "contractId", `"freezeUntil":"2024-01-16T00:00:00"`, `"skillValidity":[{"skill":"Audit"`,
		"availabilityTimeSpans", "preferredSkills",
		// Dependencies go over the wire in the solver's object form.
		`"dependsOn":[{"taskId":"task0"}]`,
		"preferredEmployees", "prohibitedEmployees", `"allowedEmployees":["emp1"]`,
		`"sla":"2024-01-20T17:00:00"`, `"earliestStart":`, `"taskType":"deep-work"`, `"priorityValue":80`,
		`"durationByEmployee":{"emp1":"PT6H"}`, `"pinned":true`, `"initialEmployeeId":"emp1"`,
		"contracts", "maxHoursPerDay",
		`"taskTypeTransitions":[{"fromTaskType":"*"`, `"taskTypes":["deep-work"]`, `"freezeUntil":"2024-01-10T00:00:00"`,
		`"weights":{"minimizeCost":"0hard/0medium/2soft"}`,
	} {
		if !strings.Contains(jsonStr, expected) {
			t.Errorf("expected JSON to contain '%s'", expected)
		}
	}
	for _, unexpected := range []string{"timeZoneId", "dedicatedClientId", "targetUtilization"} {
		if strings.Contains(jsonStr, unexpected) {
			t.Errorf("expected JSON not to contain the retired field '%s'", unexpected)
		}
	}
	// A plain dependency omits its unset options.
	if strings.Contains(jsonStr, `"anchor"`) || strings.Contains(jsonStr, `"requiresSameEmployee"`) {
		t.Errorf("expected a plain dependency to omit unset options, got %s", jsonStr)
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
				"freezeUntil": "2024-01-16T00:00:00",
				"skillValidity": [{"skill": "Audit", "validFrom": "2024-01-01T00:00:00", "validTo": "2025-01-01T00:00:00"}],
				"availabilityTimeSpans": [
					{"id": "a1", "start": "2024-01-15T08:00:00", "end": "2024-01-15T17:00:00", "type": "AVAILABLE"}
				],
				"tasks": ["task1"],
				"duration": "PT8H"
			}
		],
		"tasks": [
			{
				"id": "task1",
				"name": "API Work",
				"duration": "PT8H",
				"priority": "HIGH",
				"priorityValue": 80,
				"requiredSkills": ["C#"],
				"dependsOn": [{"taskId": "task0", "anchor": "END", "minOffset": "PT1H", "requiresSameEmployee": true}],
				"allowedEmployees": ["emp1"],
				"sla": "2024-01-20T17:00:00",
				"earliestStart": "2024-01-15T09:00:00",
				"taskType": "deep-work",
				"durationByEmployee": {"emp1": "PT6H"},
				"pinned": true,
				"parentTaskId": "feature",
				"segmentIndex": 0,
				"totalSegments": 2,
				"initialEmployeeId": "emp1",
				"initialStartTime": "2024-01-15T09:00:00",
				"employee": "emp1",
				"previousTask": null,
				"startTime": "2024-01-15T09:00:00",
				"endTime": "2024-01-15T15:00:00"
			}
		],
		"taskTypeTransitions": [{"fromTaskType": "deep-work", "toTaskType": "meeting", "setupDuration": "PT15M", "forbidden": false}],
		"taskTypes": ["deep-work", "meeting"],
		"freezeUntil": "2024-01-10T00:00:00"
	}`

	var result ProfessionalServicesResultResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(result.Employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(result.Employees))
	}
	employee := result.Employees[0]
	if *employee.HourlyRate != 100.0 {
		t.Errorf("expected hourlyRate 100.0, got %f", *employee.HourlyRate)
	}
	if *employee.ContractID != "full-time" {
		t.Errorf("expected contractId 'full-time', got '%s'", *employee.ContractID)
	}
	if *employee.FreezeUntil != "2024-01-16T00:00:00" {
		t.Errorf("expected freezeUntil to parse, got %v", employee.FreezeUntil)
	}
	if len(employee.SkillValidity) != 1 || employee.SkillValidity[0].Skill != "Audit" || *employee.SkillValidity[0].ValidTo != "2025-01-01T00:00:00" {
		t.Errorf("expected skill validity to parse, got %+v", employee.SkillValidity)
	}
	if len(employee.AvailabilityTimeSpans) != 1 {
		t.Fatalf("expected 1 availability time span, got %d", len(employee.AvailabilityTimeSpans))
	}
	if *employee.AvailabilityTimeSpans[0].Type != "AVAILABLE" {
		t.Errorf("expected type 'AVAILABLE', got '%s'", *employee.AvailabilityTimeSpans[0].Type)
	}
	if len(employee.Tasks) != 1 || employee.Tasks[0] != "task1" || *employee.Duration != "PT8H" {
		t.Errorf("expected the solver's assignment echo on the employee, got %+v", employee)
	}
	if len(result.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(result.Tasks))
	}
	task := result.Tasks[0]
	if task.ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", task.ID)
	}
	if *task.PriorityValue != 80 {
		t.Errorf("expected priorityValue 80, got %v", task.PriorityValue)
	}
	if len(task.DependsOn) != 1 || task.DependsOn[0].TaskID != "task0" || *task.DependsOn[0].Anchor != "END" || *task.DependsOn[0].MinOffset != "PT1H" || !task.DependsOn[0].RequiresSameEmployee || task.DependsOn[0].NoIntermediateTasks {
		t.Errorf("expected the object-form dependency to parse, got %+v", task.DependsOn)
	}
	if len(task.AllowedEmployees) != 1 || *task.SLA != "2024-01-20T17:00:00" || *task.EarliestStart != "2024-01-15T09:00:00" || *task.TaskType != "deep-work" {
		t.Errorf("expected whitelist/sla/earliestStart/taskType to parse, got %+v", task)
	}
	if task.DurationByEmployee["emp1"] != "PT6H" || !task.Pinned || *task.ParentTaskID != "feature" || *task.SegmentIndex != 0 || *task.TotalSegments != 2 {
		t.Errorf("expected duration overrides, pinning and split-task fields to parse, got %+v", task)
	}
	if *task.InitialEmployeeID != "emp1" || *task.Employee != "emp1" || task.PreviousTask != nil || *task.StartTime != "2024-01-15T09:00:00" || *task.EndTime != "2024-01-15T15:00:00" {
		t.Errorf("expected the solver's assignment fields to parse, got %+v", task)
	}
	if len(result.TaskTypeTransitions) != 1 || result.TaskTypeTransitions[0].FromTaskType != "deep-work" || *result.TaskTypeTransitions[0].SetupDuration != "PT15M" {
		t.Errorf("expected task type transitions to parse, got %+v", result.TaskTypeTransitions)
	}
	if len(result.TaskTypes) != 2 || *result.FreezeUntil != "2024-01-10T00:00:00" {
		t.Errorf("expected plan-level fields to parse, got %+v", result)
	}
}

func TestCanSerializeOptionsToJSON(t *testing.T) {
	request := ProfessionalServicesRequest{
		Employees: []Employee{{ID: "emp1", Skills: []string{"C#"}}},
		Tasks:     []Task{{ID: "task1", Name: "API Work", Duration: "PT8H", Priority: "HIGH", RequiredSkills: []string{"C#"}}},
		Options:   &SolverOptions{SpentLimit: strPtr("PT30S"), UnimprovedSpentLimit: strPtr("PT5S")},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}
	if !strings.Contains(string(data), `"options":{"spentLimit":"PT30S","unimprovedSpentLimit":"PT5S"}`) {
		t.Errorf("expected options in the request body, got %s", data)
	}

	// Unset option fields are omitted, and a nil Options is omitted entirely.
	request.Options = &SolverOptions{SpentLimit: strPtr("PT1M")}
	data, _ = json.Marshal(request)
	if !strings.Contains(string(data), `"options":{"spentLimit":"PT1M"}`) {
		t.Errorf("expected unset unimprovedSpentLimit to be omitted, got %s", data)
	}
	request.Options = nil
	data, _ = json.Marshal(request)
	if strings.Contains(string(data), `"options"`) {
		t.Errorf("expected nil options to be omitted, got %s", data)
	}
}
