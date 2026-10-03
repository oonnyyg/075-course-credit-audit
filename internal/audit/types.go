package audit

// Task is the JSON input describing the audit job.
type Task struct {
	Catalog  []Course  `json:"catalog"`
	Rules    Rules     `json:"rules"`
	Students []Student `json:"students"`
}

// Course describes a single catalog entry.
type Course struct {
	ID            string   `json:"id"`
	Credits       int      `json:"credits"`
	Category      string   `json:"category"`
	Prerequisites []string `json:"prerequisites"`
}

// Rules describe graduation requirements.
type Rules struct {
	Required        []string       `json:"required"`
	CategoryMinima  map[string]int `json:"category_minima"`
	TotalMinCredits int            `json:"total_min_credits"`
}

// Student groups enrollment records by student ID.
type Student struct {
	ID      string   `json:"id"`
	Records []Record `json:"records"`
}

// Record is one enrollment attempt.
type Record struct {
	Course string `json:"course"`
	Term   int    `json:"term"`
	Passed bool   `json:"passed"`
}

// Report is the JSON output.
type Report struct {
	Students []StudentReport `json:"students"`
}

// StudentReport is the audit result for one student.
type StudentReport struct {
	ID                 string            `json:"id"`
	TotalCredits       int               `json:"total_credits"`
	CategoryCredits    map[string]int    `json:"category_credits"`
	MissingRequired    []string          `json:"missing_required"`
	CategoryGaps       map[string]int    `json:"category_gaps"`
	TotalGap           int               `json:"total_gap"`
	PrereqViolations   []PrereqViolation `json:"prereq_violations"`
	GraduationEligible bool              `json:"graduation_eligible"`
}

// PrereqViolation describes one missing prerequisite at enrollment time.
type PrereqViolation struct {
	Course  string `json:"course"`
	Term    int    `json:"term"`
	Missing string `json:"missing"`
}
