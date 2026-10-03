package audit

import "testing"

func TestValidateRejects(t *testing.T) {
	cases := map[string]func(*Task){
		"duplicate course": func(t *Task) {
			t.Courses = append(t.Courses, t.Courses[0])
		},
		"zero credits": func(t *Task) {
			t.Courses[0].Credits = 0
		},
		"self dependency": func(t *Task) {
			t.Courses[0].Prerequisites = []string{"A"}
		},
		"cycle": func(t *Task) {
			t.Courses[0].Prerequisites = []string{"B"}
		},
		"unknown prerequisite": func(t *Task) {
			t.Courses[0].Prerequisites = []string{"ZZZ"}
		},
		"unknown required": func(t *Task) {
			t.Requirements.Required = []string{"ZZZ"}
		},
		"unknown category": func(t *Task) {
			t.Requirements.CategoryMin = map[string]int{"art": 1}
		},
		"negative total min": func(t *Task) {
			t.Requirements.TotalMin = -1
		},
		"negative category min": func(t *Task) {
			t.Requirements.CategoryMin = map[string]int{"cs": -1}
		},
		"duplicate student": func(t *Task) {
			t.Students = []Student{{ID: "S1"}, {ID: "S1"}}
		},
		"unknown record course": func(t *Task) {
			t.Students = []Student{{ID: "S1", Records: []Record{{Course: "ZZZ", Term: 1}}}}
		},
		"zero term": func(t *Task) {
			t.Students = []Student{{ID: "S1", Records: []Record{{Course: "A", Term: 0}}}}
		},
		"duplicate record": func(t *Task) {
			t.Students = []Student{{ID: "S1", Records: []Record{
				{Course: "A", Term: 1, Passed: true},
				{Course: "A", Term: 1, Passed: false},
			}}}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			task := baseTask()
			mutate(task)
			if err := Validate(task); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestValidateAcceptsBase(t *testing.T) {
	if err := Validate(baseTask()); err != nil {
		t.Fatal(err)
	}
}
