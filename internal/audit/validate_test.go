package audit

import "testing"

func invalidTasks() map[string]func(*Task) {
	return map[string]func(*Task){
		"duplicate course": func(t *Task) { t.Catalog = append(t.Catalog, t.Catalog[0]) },
		"zero credits":     func(t *Task) { t.Catalog[0].Credits = 0 },
		"empty category":   func(t *Task) { t.Catalog[0].Category = "" },
		"self prereq":      func(t *Task) { t.Catalog[0].Prerequisites = []string{"CS100"} },
		"unknown prereq":   func(t *Task) { t.Catalog[0].Prerequisites = []string{"NOPE"} },
		"cycle": func(t *Task) {
			t.Catalog[0].Prerequisites = []string{"CS101"}
		},
		"negative total":    func(t *Task) { t.Rules.TotalMinCredits = -1 },
		"unknown required":  func(t *Task) { t.Rules.Required = []string{"NOPE"} },
		"unknown category":  func(t *Task) { t.Rules.CategoryMinima = map[string]int{"art": 1} },
		"negative cat min":  func(t *Task) { t.Rules.CategoryMinima = map[string]int{"cs": -1} },
		"duplicate student": func(t *Task) { t.Students = append(t.Students, t.Students[0]) },
		"unknown record course": func(t *Task) {
			t.Students[0].Records[0].Course = "NOPE"
		},
		"non-positive term": func(t *Task) { t.Students[0].Records[0].Term = 0 },
		"duplicate record": func(t *Task) {
			t.Students[0].Records = append(t.Students[0].Records, t.Students[0].Records[0])
		},
	}
}

func TestValidateRejectsInvalidTasks(t *testing.T) {
	for name, mutate := range invalidTasks() {
		t.Run(name, func(t *testing.T) {
			task := baseTask()
			task.Students = []Student{{ID: "S1", Records: []Record{
				{Course: "CS100", Term: 1, Passed: true},
			}}}
			mutate(task)
			if err := ValidateTask(task); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidateAcceptsBaseTask(t *testing.T) {
	if err := ValidateTask(baseTask()); err != nil {
		t.Fatal(err)
	}
}
