package audit

import "sort"

// Evaluate runs the audit for a validated task and returns a report
// with students sorted by ID and all lists deterministically ordered.
func Evaluate(t *Task) Report {
	courses := make(map[string]Course, len(t.Courses))
	for _, c := range t.Courses {
		courses[c.ID] = c
	}
	students := make([]Student, len(t.Students))
	copy(students, t.Students)
	sort.Slice(students, func(i, j int) bool { return students[i].ID < students[j].ID })

	rep := Report{Students: make([]StudentReport, 0, len(students))}
	for _, s := range students {
		rep.Students = append(rep.Students, evaluateStudent(s, courses, t.Requirements))
	}
	return rep
}

func evaluateStudent(s Student, courses map[string]Course, req Requirements) StudentReport {
	records := make([]Record, len(s.Records))
	copy(records, s.Records)
	sort.Slice(records, func(i, j int) bool {
		if records[i].Term != records[j].Term {
			return records[i].Term < records[j].Term
		}
		return records[i].Course < records[j].Course
	})

	firstPass := make(map[string]int)
	earned := make(map[string]bool)
	var violations []PrerequisiteViolation
	for _, r := range records {
		var missing []string
		for _, p := range courses[r.Course].Prerequisites {
			term, ok := firstPass[p]
			if !ok || term >= r.Term {
				missing = append(missing, p)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			violations = append(violations, PrerequisiteViolation{
				Course:  r.Course,
				Term:    r.Term,
				Missing: missing,
			})
		}
		if r.Passed {
			if _, ok := firstPass[r.Course]; !ok {
				firstPass[r.Course] = r.Term
				earned[r.Course] = true
			}
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].Course != violations[j].Course {
			return violations[i].Course < violations[j].Course
		}
		return violations[i].Term < violations[j].Term
	})

	total := 0
	catCredits := make(map[string]int)
	for id := range earned {
		c := courses[id]
		total += c.Credits
		catCredits[c.Category] += c.Credits
	}

	missingReq := make([]string, 0)
	for _, id := range req.Required {
		if !earned[id] {
			missingReq = append(missingReq, id)
		}
	}
	sort.Strings(missingReq)

	catDeficits := make(map[string]int)
	for cat, min := range req.CategoryMin {
		if d := min - catCredits[cat]; d > 0 {
			catDeficits[cat] = d
		}
	}
	totalDeficit := req.TotalMin - total
	if totalDeficit < 0 {
		totalDeficit = 0
	}

	graduated := len(missingReq) == 0 && len(catDeficits) == 0 &&
		totalDeficit == 0 && len(violations) == 0
	if violations == nil {
		violations = []PrerequisiteViolation{}
	}
	return StudentReport{
		ID:                     s.ID,
		TotalCredits:           total,
		CategoryCredits:        catCredits,
		MissingRequired:        missingReq,
		CategoryDeficits:       catDeficits,
		TotalDeficit:           totalDeficit,
		PrerequisiteViolations: violations,
		Graduated:              graduated,
	}
}
