package audit

// Task describes one audit job: course catalog, graduation rules and
// the per-student enrollment records to evaluate.
type Task struct {
	Courses      []Course     `json:"courses"`
	Requirements Requirements `json:"requirements"`
	Students     []Student    `json:"students"`
}

// Course is a catalog entry. ID is unique, Credits is a positive
// integer and Prerequisites lists direct prerequisite course IDs.
type Course struct {
	ID            string   `json:"id"`
	Credits       int      `json:"credits"`
	Category      string   `json:"category"`
	Prerequisites []string `json:"prerequisites"`
}

// Requirements defines the graduation rules.
type Requirements struct {
	Required    []string       `json:"required"`
	CategoryMin map[string]int `json:"categoryMin"`
	TotalMin    int            `json:"totalMin"`
}

// Student holds one student's enrollment records.
type Student struct {
	ID      string   `json:"id"`
	Records []Record `json:"records"`
}

// Record is a single enrollment: course, term (positive) and result.
type Record struct {
	Course string `json:"course"`
	Term   int    `json:"term"`
	Passed bool   `json:"passed"`
}

// Report is the audit result for all students.
type Report struct {
	Students []StudentReport `json:"students"`
}

// StudentReport summarizes one student's standing.
type StudentReport struct {
	ID                     string                  `json:"id"`
	TotalCredits           int                     `json:"totalCredits"`
	CategoryCredits        map[string]int          `json:"categoryCredits"`
	MissingRequired        []string                `json:"missingRequired"`
	CategoryDeficits       map[string]int          `json:"categoryDeficits"`
	TotalDeficit           int                     `json:"totalDeficit"`
	PrerequisiteViolations []PrerequisiteViolation `json:"prerequisiteViolations"`
	Graduated              bool                    `json:"graduated"`
}

// PrerequisiteViolation records one enrollment whose direct
// prerequisites were not all passed in strictly earlier terms.
type PrerequisiteViolation struct {
	Course  string   `json:"course"`
	Term    int      `json:"term"`
	Missing []string `json:"missing"`
}
