# Course Credit Audit

A standard-library-only Go tool that audits graduation eligibility from a
course catalog and per-student enrollment records. It reads a JSON task file
and writes a JSON report; it never connects to any external system.

## Usage

```
go run ./cmd/degree <task.json> <report.json>
```

Exit code is non-zero and no report is written when the input is invalid.

## Task format

```json
{
  "catalog": [
    {"id": "CS100", "credits": 3, "category": "cs", "prerequisites": []},
    {"id": "CS101", "credits": 3, "category": "cs", "prerequisites": ["CS100"]}
  ],
  "rules": {
    "required": ["CS101"],
    "category_minima": {"cs": 6},
    "total_min_credits": 10
  },
  "students": [
    {"id": "S1", "records": [
      {"course": "CS100", "term": 1, "passed": true}
    ]}
  ]
}
```

- Course IDs are unique; `credits` is a positive integer; every course has
  one `category` and a list of direct `prerequisites`.
- Prerequisite references must exist; self-dependencies and cycles are
  rejected.
- `required`, `category_minima` and `total_min_credits` define graduation
  rules. Thresholds are non-negative; unknown courses/categories are rejected.
- Records carry a course ID, a positive integer `term`, and `passed`. The
  same student cannot have two records for the same course in the same term.

## Semantics

- A course earns credits exactly once, at its first passing record; later
  failures do not revoke credits and repeated passes do not add more.
- For every enrollment, each direct prerequisite must have been passed in a
  strictly earlier term; passing in the same term does not count. Every
  missing prerequisite is listed with the enrollment term. Future results are
  never used to justify past enrollments.
- Prerequisite violations do not revoke earned credits, but any violation
  blocks graduation.
- A course counts toward both its category and the total. A student is
  `graduation_eligible` only when all required courses are passed, every
  category minimum and the total threshold are met, and there are no
  prerequisite violations. Empty records simply report the full gaps.

## Report format

Per student (sorted by student ID, courses sorted by ID):

- `total_credits`, `category_credits`: earned credits.
- `missing_required`: unpassed required course IDs.
- `category_gaps`, `total_gap`: non-negative shortfalls.
- `prereq_violations`: `{course, term, missing}` entries.
- `graduation_eligible`: final verdict.

See `examples/task.json` and `examples/report.json` for a full run.

## Tests

```
go test ./...
```

Covers retakes, same-term prerequisites, category-met-but-total-short,
empty records, input-order independence, and rejection of invalid tasks.