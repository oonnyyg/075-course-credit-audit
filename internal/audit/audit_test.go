package audit

import (
	"reflect"
	"testing"
)

func baseTask() *Task {
	return &Task{
		Catalog: []Course{
			{ID: "CS100", Credits: 3, Category: "cs"},
			{ID: "CS101", Credits: 3, Category: "cs", Prerequisites: []string{"CS100"}},
			{ID: "MA100", Credits: 4, Category: "math"},
		},
		Rules: Rules{
			Required:        []string{"CS101"},
			CategoryMinima:  map[string]int{"cs": 6, "math": 4},
			TotalMinCredits: 10,
		},
	}
}

func TestRetakeEarnsOnceAndFailDoesNotRevoke(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{ID: "S1", Records: []Record{
		{Course: "CS100", Term: 1, Passed: false},
		{Course: "CS100", Term: 2, Passed: true},
		{Course: "CS100", Term: 3, Passed: true},
		{Course: "CS100", Term: 4, Passed: false},
		{Course: "CS101", Term: 3, Passed: true},
		{Course: "MA100", Term: 1, Passed: true},
	}}}
	if err := ValidateTask(task); err != nil {
		t.Fatal(err)
	}
	r := Audit(task).Students[0]
	if r.TotalCredits != 10 {
		t.Fatalf("total = %d, want 10", r.TotalCredits)
	}
	if r.CategoryCredits["cs"] != 6 || r.CategoryCredits["math"] != 4 {
		t.Fatalf("category credits = %v", r.CategoryCredits)
	}
	if len(r.PrereqViolations) != 0 {
		t.Fatalf("unexpected violations: %v", r.PrereqViolations)
	}
	if !r.GraduationEligible {
		t.Fatalf("expected eligible, report: %+v", r)
	}
}

func TestSameTermPrereqDoesNotSatisfy(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{ID: "S1", Records: []Record{
		{Course: "CS100", Term: 1, Passed: true},
		{Course: "CS101", Term: 1, Passed: true},
	}}}
	if err := ValidateTask(task); err != nil {
		t.Fatal(err)
	}
	r := Audit(task).Students[0]
	if len(r.PrereqViolations) != 1 {
		t.Fatalf("violations = %v, want 1", r.PrereqViolations)
	}
	v := r.PrereqViolations[0]
	if v.Course != "CS101" || v.Term != 1 || v.Missing != "CS100" {
		t.Fatalf("violation = %+v", v)
	}
	// Credits are kept despite the violation, but graduation is blocked.
	if r.TotalCredits != 6 {
		t.Fatalf("total = %d, want 6", r.TotalCredits)
	}
	if r.GraduationEligible {
		t.Fatal("must not be eligible with prereq violation")
	}
}

func TestCategoryMetButTotalShort(t *testing.T) {
	task := baseTask()
	task.Rules.CategoryMinima = map[string]int{"cs": 3}
	task.Rules.Required = nil
	task.Rules.TotalMinCredits = 5
	task.Students = []Student{{ID: "S1", Records: []Record{
		{Course: "CS100", Term: 1, Passed: true},
	}}}
	r := Audit(task).Students[0]
	if r.CategoryGaps["cs"] != 0 {
		t.Fatalf("cs gap = %d, want 0", r.CategoryGaps["cs"])
	}
	if r.TotalGap != 2 {
		t.Fatalf("total gap = %d, want 2", r.TotalGap)
	}
	if r.GraduationEligible {
		t.Fatal("must not be eligible with total gap")
	}
}

func TestEmptyRecordsReportGaps(t *testing.T) {
	task := baseTask()
	task.Students = []Student{{ID: "S1"}}
	r := Audit(task).Students[0]
	if r.TotalCredits != 0 || r.TotalGap != 10 {
		t.Fatalf("got total=%d gap=%d", r.TotalCredits, r.TotalGap)
	}
	if len(r.MissingRequired) != 1 || r.MissingRequired[0] != "CS101" {
		t.Fatalf("missing = %v", r.MissingRequired)
	}
	if r.CategoryGaps["cs"] != 6 || r.CategoryGaps["math"] != 4 {
		t.Fatalf("gaps = %v", r.CategoryGaps)
	}
	if r.GraduationEligible {
		t.Fatal("must not be eligible")
	}
}

func TestInputOrderDoesNotMatter(t *testing.T) {
	a := baseTask()
	a.Students = []Student{
		{ID: "S2", Records: []Record{{Course: "CS100", Term: 1, Passed: true}}},
		{ID: "S1", Records: []Record{
			{Course: "CS101", Term: 2, Passed: true},
			{Course: "CS100", Term: 1, Passed: true},
		}},
	}
	b := baseTask()
	b.Catalog = []Course{a.Catalog[2], a.Catalog[1], a.Catalog[0]}
	b.Students = []Student{a.Students[1], a.Students[0]}
	ra, rb := Audit(a), Audit(b)
	if len(ra.Students) != 2 || ra.Students[0].ID != "S1" || ra.Students[1].ID != "S2" {
		t.Fatalf("students not sorted: %+v", ra.Students)
	}
	for i := range ra.Students {
		if !reflect.DeepEqual(ra.Students[i], rb.Students[i]) {
			t.Fatalf("student %d differs: %+v vs %+v", i, ra.Students[i], rb.Students[i])
		}
	}
}
