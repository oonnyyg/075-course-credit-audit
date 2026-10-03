package audit

import "fmt"

// Validate checks the task for structural errors: duplicate or
// unknown references, non-positive numbers, self dependencies and
// prerequisite cycles.
func Validate(t *Task) error {
	courses := make(map[string]Course, len(t.Courses))
	categories := make(map[string]bool)
	for _, c := range t.Courses {
		if c.ID == "" {
			return fmt.Errorf("course with empty id")
		}
		if _, dup := courses[c.ID]; dup {
			return fmt.Errorf("duplicate course id %q", c.ID)
		}
		if c.Credits <= 0 {
			return fmt.Errorf("course %q: credits must be positive", c.ID)
		}
		if c.Category == "" {
			return fmt.Errorf("course %q: empty category", c.ID)
		}
		courses[c.ID] = c
		categories[c.Category] = true
	}
	for _, c := range t.Courses {
		seen := make(map[string]bool)
		for _, p := range c.Prerequisites {
			if p == c.ID {
				return fmt.Errorf("course %q: self dependency", c.ID)
			}
			if seen[p] {
				return fmt.Errorf("course %q: duplicate prerequisite %q", c.ID, p)
			}
			seen[p] = true
			if _, ok := courses[p]; !ok {
				return fmt.Errorf("course %q: unknown prerequisite %q", c.ID, p)
			}
		}
	}
	if err := checkCycles(courses); err != nil {
		return err
	}
	req := t.Requirements
	if req.TotalMin < 0 {
		return fmt.Errorf("totalMin must be non-negative")
	}
	seenReq := make(map[string]bool)
	for _, id := range req.Required {
		if seenReq[id] {
			return fmt.Errorf("duplicate required course %q", id)
		}
		seenReq[id] = true
		if _, ok := courses[id]; !ok {
			return fmt.Errorf("unknown required course %q", id)
		}
	}
	for cat, min := range req.CategoryMin {
		if !categories[cat] {
			return fmt.Errorf("unknown category %q in requirements", cat)
		}
		if min < 0 {
			return fmt.Errorf("category %q: negative minimum", cat)
		}
	}
	seenStudent := make(map[string]bool)
	for _, s := range t.Students {
		if s.ID == "" {
			return fmt.Errorf("student with empty id")
		}
		if seenStudent[s.ID] {
			return fmt.Errorf("duplicate student id %q", s.ID)
		}
		seenStudent[s.ID] = true
		type key struct {
			course string
			term   int
		}
		seenRec := make(map[key]bool)
		for _, r := range s.Records {
			if _, ok := courses[r.Course]; !ok {
				return fmt.Errorf("student %q: unknown course %q", s.ID, r.Course)
			}
			if r.Term <= 0 {
				return fmt.Errorf("student %q: course %q: term must be positive", s.ID, r.Course)
			}
			k := key{r.Course, r.Term}
			if seenRec[k] {
				return fmt.Errorf("student %q: duplicate record for %q in term %d", s.ID, r.Course, r.Term)
			}
			seenRec[k] = true
		}
	}
	return nil
}

func checkCycles(courses map[string]Course) error {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := make(map[string]int, len(courses))
	var visit func(id string) error
	visit = func(id string) error {
		color[id] = gray
		for _, p := range courses[id].Prerequisites {
			switch color[p] {
			case gray:
				return fmt.Errorf("prerequisite cycle involving %q", p)
			case white:
				if err := visit(p); err != nil {
					return err
				}
			}
		}
		color[id] = black
		return nil
	}
	for id := range courses {
		if color[id] == white {
			if err := visit(id); err != nil {
				return err
			}
		}
	}
	return nil
}
