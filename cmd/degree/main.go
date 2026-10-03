package main

import (
	"encoding/json"
	"fmt"
	"os"

	"example.com/course-credit-audit/internal/audit"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: degree <task.json> <report.json>")
		os.Exit(2)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	var task audit.Task
	if err := json.Unmarshal(data, &task); err != nil {
		fmt.Fprintln(os.Stderr, "error: invalid task JSON:", err)
		os.Exit(1)
	}
	if err := audit.ValidateTask(&task); err != nil {
		fmt.Fprintln(os.Stderr, "error: invalid task:", err)
		os.Exit(1)
	}
	report := audit.Audit(&task)
	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	out = append(out, '\n')
	if err := os.WriteFile(os.Args[2], out, 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
