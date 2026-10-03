// Command degree reads a JSON audit task and writes a JSON report.
//
// Usage: go run ./cmd/degree <task.json> <report.json>
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/course-credit-audit/audit"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: degree <task.json> <report.json>")
		os.Exit(2)
	}
	if err := run(os.Args[1], os.Args[2]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(taskPath, reportPath string) error {
	data, err := os.ReadFile(taskPath)
	if err != nil {
		return fmt.Errorf("read task: %w", err)
	}
	var task audit.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return fmt.Errorf("parse task: %w", err)
	}
	if err := audit.Validate(&task); err != nil {
		return fmt.Errorf("invalid task: %w", err)
	}
	report := audit.Evaluate(&task)
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode report: %w", err)
	}
	out = append(out, '\n')
	if err := os.WriteFile(reportPath, out, 0o644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}
