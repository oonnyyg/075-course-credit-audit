package audit

import (
	"reflect"
	"testing"
)

func baseTask() *Task {
	return &Task{
		Courses: []Course{
			{ID: "A", Credits: 3, Category: "cs"},
			{ID: "B", Credits: 3, Category: "cs", Prerequisites: []string{"A"}},
			{ID: "C", Credits: 2, Category: "math"},
		},
		Requirements: Requirements{
			Required:    []string{"A", "B"},
			CategoryMin: map[string]int{"cs": 6, "math": 2},
			TotalMin:    8,
		},
	}
}

func TestRetakeEarnsOnce(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{
		ID: "S1",
		Records: []Record{
			{Course: "A", Term: 1, Passed: false},
			{Course: "A", Term: 2, Passed: true},
			{Course: "A", Term: 3, Passed: true},
			{Course: "A", Term: 4, Passed: false},
			{Course: "B", Term: 3, Passed: true},
			{Course: "C", Term: 1, Passed: true},
		},
	}}
	if err := Validate(task); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(task).Students[0]
	if rep.TotalCredits != 8 {
		t.Fatalf("total = %d, want 8", rep.TotalCredits)
	}
	if rep.CategoryCredits["cs"] != 6 || rep.CategoryCredits["math"] != 2 {
		t.Fatalf("category credits = %v", rep.CategoryCredits)
	}
	if !rep.Graduated {
		t.Fatalf("expected graduated, report: %+v", rep)
	}
}

func TestSameTermPrerequisiteViolation(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{
		ID: "S1",
		Records: []Record{
			{Course: "A", Term: 1, Passed: true},
			{Course: "B", Term: 1, Passed: true},
		},
	}}
	if err := Validate(task); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(task).Students[0]
	want := []PrerequisiteViolation{{Course: "B", Term: 1, Missing: []string{"A"}}}
	if !reflect.DeepEqual(rep.PrerequisiteViolations, want) {
		t.Fatalf("violations = %+v, want %+v", rep.PrerequisiteViolations, want)
	}
	if rep.Graduated {
		t.Fatal("violation must block graduation")
	}
	if rep.TotalCredits != 6 {
		t.Fatalf("violation must not revoke credits, total = %d", rep.TotalCredits)
	}
}

func TestCategoryMetTotalShort(t *testing.T) {
	task := baseTask()
	task.Requirements.TotalMin = 9
	task.Students = []Student{{
		ID: "S1",
		Records: []Record{
			{Course: "A", Term: 1, Passed: true},
			{Course: "B", Term: 2, Passed: true},
			{Course: "C", Term: 1, Passed: true},
		},
	}}
	if err := Validate(task); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(task).Students[0]
	if len(rep.CategoryDeficits) != 0 {
		t.Fatalf("category deficits = %v, want none", rep.CategoryDeficits)
	}
	if rep.TotalDeficit != 1 {
		t.Fatalf("total deficit = %d, want 1", rep.TotalDeficit)
	}
	if rep.Graduated {
		t.Fatal("total deficit must block graduation")
	}
}

func TestEmptyRecordsReportDeficits(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{ID: "S1"}}
	if err := Validate(task); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(task).Students[0]
	if rep.TotalCredits != 0 || rep.TotalDeficit != 8 {
		t.Fatalf("got total=%d deficit=%d", rep.TotalCredits, rep.TotalDeficit)
	}
	if !reflect.DeepEqual(rep.MissingRequired, []string{"A", "B"}) {
		t.Fatalf("missing required = %v", rep.MissingRequired)
	}
	if rep.CategoryDeficits["cs"] != 6 || rep.CategoryDeficits["math"] != 2 {
		t.Fatalf("category deficits = %v", rep.CategoryDeficits)
	}
	if rep.Graduated {
		t.Fatal("empty records cannot graduate")
	}
}

func TestStableOrdering(t *testing.T) {
	task := baseTask()
	task.Students = []Student{
		{ID: "S2", Records: []Record{{Course: "A", Term: 1, Passed: true}}},
		{ID: "S1"},
	}
	if err := Validate(task); err != nil {
		t.Fatal(err)
	}
	rep := Evaluate(task)
	if rep.Students[0].ID != "S1" || rep.Students[1].ID != "S2" {
		t.Fatalf("students not sorted: %v", rep.Students)
	}
}
