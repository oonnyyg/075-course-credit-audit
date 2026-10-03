package audit

import (
	"fmt"
	"sort"
)

// ValidateTask checks the catalog, rules and student records for consistency.
func ValidateTask(t *Task) error {
	courses := map[string]Course{}
	for _, c := range t.Catalog {
		if c.ID == "" {
			return fmt.Errorf("course with empty id")
		}
		if _, dup := courses[c.ID]; dup {
			return fmt.Errorf("duplicate course id %q", c.ID)
		}
		if c.Credits <= 0 {
			return fmt.Errorf("course %q: credits must be a positive integer", c.ID)
		}
		if c.Category == "" {
			return fmt.Errorf("course %q: category must not be empty", c.ID)
		}
		courses[c.ID] = c
	}
	for _, c := range t.Catalog {
		for _, p := range c.Prerequisites {
			if p == c.ID {
				return fmt.Errorf("course %q: self prerequisite", c.ID)
			}
			if _, ok := courses[p]; !ok {
				return fmt.Errorf("course %q: unknown prerequisite %q", c.ID, p)
			}
		}
	}
	if err := checkAcyclic(courses); err != nil {
		return err
	}
	if t.Rules.TotalMinCredits < 0 {
		return fmt.Errorf("rules: total_min_credits must be non-negative")
	}
	for _, id := range t.Rules.Required {
		if _, ok := courses[id]; !ok {
			return fmt.Errorf("rules: unknown required course %q", id)
		}
	}
	knownCat := map[string]bool{}
	for _, c := range t.Catalog {
		knownCat[c.Category] = true
	}
	for cat, min := range t.Rules.CategoryMinima {
		if !knownCat[cat] {
			return fmt.Errorf("rules: unknown category %q", cat)
		}
		if min < 0 {
			return fmt.Errorf("rules: category %q minimum must be non-negative", cat)
		}
	}
	seenStudents := map[string]bool{}
	for _, s := range t.Students {
		if s.ID == "" {
			return fmt.Errorf("student with empty id")
		}
		if seenStudents[s.ID] {
			return fmt.Errorf("duplicate student id %q", s.ID)
		}
		seenStudents[s.ID] = true
		type key struct {
			course string
			term   int
		}
		seen := map[key]bool{}
		for _, r := range s.Records {
			if _, ok := courses[r.Course]; !ok {
				return fmt.Errorf("student %q: unknown course %q", s.ID, r.Course)
			}
			if r.Term <= 0 {
				return fmt.Errorf("student %q: term must be a positive integer for course %q", s.ID, r.Course)
			}
			k := key{r.Course, r.Term}
			if seen[k] {
				return fmt.Errorf("student %q: duplicate record for course %q in term %d", s.ID, r.Course, r.Term)
			}
			seen[k] = true
		}
	}
	return nil
}

func checkAcyclic(courses map[string]Course) error {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var visit func(id string) error
	visit = func(id string) error {
		color[id] = gray
		for _, p := range courses[id].Prerequisites {
			switch color[p] {
			case gray:
				return fmt.Errorf("prerequisite cycle involving course %q", p)
			case white:
				if err := visit(p); err != nil {
					return err
				}
			}
		}
		color[id] = black
		return nil
	}
	ids := make([]string, 0, len(courses))
	for id := range courses {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if color[id] == white {
			if err := visit(id); err != nil {
				return err
			}
		}
	}
	return nil
}
