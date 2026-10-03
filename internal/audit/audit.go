package audit

import "sort"

// Audit evaluates every student against the catalog and rules.
// The task must have passed ValidateTask.
func Audit(t *Task) Report {
	courses := map[string]Course{}
	for _, c := range t.Catalog {
		courses[c.ID] = c
	}
	students := make([]Student, len(t.Students))
	copy(students, t.Students)
	sort.Slice(students, func(i, j int) bool { return students[i].ID < students[j].ID })

	report := Report{Students: []StudentReport{}}
	for _, s := range students {
		report.Students = append(report.Students, auditStudent(s, courses, t.Rules))
	}
	return report
}

func auditStudent(s Student, courses map[string]Course, rules Rules) StudentReport {
	records := make([]Record, len(s.Records))
	copy(records, s.Records)
	sort.Slice(records, func(i, j int) bool {
		if records[i].Term != records[j].Term {
			return records[i].Term < records[j].Term
		}
		return records[i].Course < records[j].Course
	})

	passedTerm := map[string]int{} // course -> earliest passed term
	earned := map[string]bool{}
	violations := []PrereqViolation{}
	for _, r := range records {
		c := courses[r.Course]
		if r.Passed {
			for _, p := range c.Prerequisites {
				pt, ok := passedTerm[p]
				if !ok || pt >= r.Term {
					violations = append(violations, PrereqViolation{Course: r.Course, Term: r.Term, Missing: p})
				}
			}
			if prev, ok := passedTerm[r.Course]; !ok || r.Term < prev {
				passedTerm[r.Course] = r.Term
			}
			earned[r.Course] = true
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].Course != violations[j].Course {
			return violations[i].Course < violations[j].Course
		}
		if violations[i].Term != violations[j].Term {
			return violations[i].Term < violations[j].Term
		}
		return violations[i].Missing < violations[j].Missing
	})

	total := 0
	catCredits := map[string]int{}
	for id := range earned {
		c := courses[id]
		total += c.Credits
		catCredits[c.Category] += c.Credits
	}

	missing := []string{}
	for _, id := range rules.Required {
		if !earned[id] {
			missing = append(missing, id)
		}
	}
	sort.Strings(missing)

	gaps := map[string]int{}
	for cat, min := range rules.CategoryMinima {
		gap := min - catCredits[cat]
		if gap < 0 {
			gap = 0
		}
		gaps[cat] = gap
	}
	totalGap := rules.TotalMinCredits - total
	if totalGap < 0 {
		totalGap = 0
	}

	eligible := len(missing) == 0 && len(violations) == 0 && totalGap == 0
	for _, g := range gaps {
		if g > 0 {
			eligible = false
			break
		}
	}

	return StudentReport{
		ID:                 s.ID,
		TotalCredits:       total,
		CategoryCredits:    catCredits,
		MissingRequired:    missing,
		CategoryGaps:       gaps,
		TotalGap:           totalGap,
		PrereqViolations:   violations,
		GraduationEligible: eligible,
	}
}
