package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestE2E_InitThenList(t *testing.T) {
	// End-to-end test: init project → list finds resources
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "e2e-test")

	// Step 1: Create directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	// Step 2: Create example query file (simulating init)
	exampleContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"endpoint", "service"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
		query.Count(),
	},
	Filters: []query.Filter{
		query.GT("duration_ms", 500),
	},
	Limit: 100,
}

var ErrorRate = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("status_code", 400),
	},
}
`
	if err := os.WriteFile(filepath.Join(projectPath, "queries.go"), []byte(exampleContent), 0644); err != nil {
		t.Fatalf("failed to write queries.go: %v", err)
	}

	// Step 3: List should find resources
	resources, err := discovery.DiscoverAll(projectPath)
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	// Should find SlowRequests and ErrorRate
	if len(resources.Queries) < 2 {
		t.Errorf("expected at least 2 queries, got %d", len(resources.Queries))
	}

	// Verify we found SlowRequests
	foundSlowRequests := false
	foundErrorRate := false
	for _, q := range resources.Queries {
		switch q.Name {
		case "SlowRequests":
			foundSlowRequests = true
			if q.Dataset != "production" {
				t.Errorf("SlowRequests dataset = %q, want production", q.Dataset)
			}
		case "ErrorRate":
			foundErrorRate = true
		}
	}

	if !foundSlowRequests {
		t.Error("list did not find SlowRequests query")
	}
	if !foundErrorRate {
		t.Error("list did not find ErrorRate query")
	}
}

func TestE2E_InitImportList(t *testing.T) {
	// Full e2e test: init → import → list
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "full-e2e")

	// Step 1: Create directory and initial query file (simulating init)
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	initContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var SlowRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(2),
	Breakdowns: []string{"endpoint"},
	Calculations: []query.Calculation{
		query.P99("duration_ms"),
	},
}
`
	if err := os.WriteFile(filepath.Join(projectPath, "queries.go"), []byte(initContent), 0644); err != nil {
		t.Fatalf("failed to write queries.go: %v", err)
	}

	// Step 2: Create Query JSON file and import it to a new Go file
	queryJSON := map[string]any{
		"time_range":   3600,
		"breakdowns":   []string{"service"},
		"calculations": []map[string]any{{"op": "COUNT"}},
		"filters": []map[string]any{
			{"column": "http.status_code", "op": ">=", "value": 500},
		},
	}
	jsonData, err := json.Marshal(queryJSON)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	jsonFile := filepath.Join(projectPath, "error_query.json")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		t.Fatalf("failed to write JSON file: %v", err)
	}

	// Simulate import by creating Go code from JSON
	importedContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var ErrorQuery = query.Query{
	TimeRange: query.Hours(1),
	Breakdowns: []string{"service"},
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.GTE("http.status_code", 500),
	},
}
`
	outputFile := filepath.Join(projectPath, "errors.go")
	if err := os.WriteFile(outputFile, []byte(importedContent), 0644); err != nil {
		t.Fatalf("failed to write imported file: %v", err)
	}

	// Step 3: List should find resources from both init and import
	resources, err := discovery.DiscoverAll(projectPath)
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	// Should find: SlowRequests (from init) + ErrorQuery (from import)
	if len(resources.Queries) < 2 {
		t.Errorf("expected at least 2 queries, got %d", len(resources.Queries))
	}

	// Verify we found resources from init
	foundSlowRequests := false
	// Verify we found resources from import
	foundErrorQuery := false

	for _, q := range resources.Queries {
		switch q.Name {
		case "SlowRequests":
			foundSlowRequests = true
		case "ErrorQuery":
			foundErrorQuery = true
		}
	}

	if !foundSlowRequests {
		t.Error("list did not find SlowRequests (from init)")
	}
	if !foundErrorQuery {
		t.Error("list did not find ErrorQuery (from import)")
	}
}

func TestE2E_AllResourceTypes(t *testing.T) {
	// E2E test with queries, SLOs, triggers, and boards
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, "all-types")

	if err := os.MkdirAll(projectPath, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	// Create file with all resource types
	content := `package monitoring

import (
	"github.com/lex00/wetwire-honeycomb-go/query"
	"github.com/lex00/wetwire-honeycomb-go/slo"
	"github.com/lex00/wetwire-honeycomb-go/trigger"
	"github.com/lex00/wetwire-honeycomb-go/board"
)

var SuccessfulRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
	Filters: []query.Filter{
		query.LT("http.status_code", 500),
	},
}

var AllRequests = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
}

var APIAvailability = slo.SLO{
	Name:        "API Availability",
	Description: "Track API availability",
	Dataset:     "production",
	SLI: slo.SLI{
		GoodEvents:  SuccessfulRequests,
		TotalEvents: AllRequests,
	},
	Target:     slo.Percentage(99.9),
	TimePeriod: slo.Days(30),
}

var HighLatencyAlert = trigger.Trigger{
	Name:        "High Latency Alert",
	Description: "Alert when P99 exceeds 1s",
	Dataset:     "production",
	Threshold:   trigger.GreaterThan(1000),
	Frequency:   trigger.Minutes(5),
	Recipients: []trigger.Recipient{
		trigger.SlackChannel("#alerts"),
	},
}

var PerformanceBoard = board.Board{
	Name:        "Performance Dashboard",
	Description: "Overview of system performance",
}
`
	if err := os.WriteFile(filepath.Join(projectPath, "monitoring.go"), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write monitoring.go: %v", err)
	}

	// Discover all resources
	resources, err := discovery.DiscoverAll(projectPath)
	if err != nil {
		t.Fatalf("DiscoverAll failed: %v", err)
	}

	// Verify counts
	if len(resources.Queries) < 2 {
		t.Errorf("expected at least 2 queries, got %d", len(resources.Queries))
	}
	if len(resources.SLOs) < 1 {
		t.Errorf("expected at least 1 SLO, got %d", len(resources.SLOs))
	}
	if len(resources.Triggers) < 1 {
		t.Errorf("expected at least 1 trigger, got %d", len(resources.Triggers))
	}
	if len(resources.Boards) < 1 {
		t.Errorf("expected at least 1 board, got %d", len(resources.Boards))
	}

	// Verify specific resources
	foundAPIAvailability := false
	for _, s := range resources.SLOs {
		if s.Name == "APIAvailability" {
			foundAPIAvailability = true
			if s.TargetPercentage != 99.9 {
				t.Errorf("APIAvailability target = %.1f, want 99.9", s.TargetPercentage)
			}
		}
	}
	if !foundAPIAvailability {
		t.Error("list did not find APIAvailability SLO")
	}

	foundHighLatencyAlert := false
	for _, tr := range resources.Triggers {
		if tr.Name == "HighLatencyAlert" {
			foundHighLatencyAlert = true
		}
	}
	if !foundHighLatencyAlert {
		t.Error("list did not find HighLatencyAlert trigger")
	}

	foundPerformanceBoard := false
	for _, b := range resources.Boards {
		if b.Name == "PerformanceBoard" {
			foundPerformanceBoard = true
		}
	}
	if !foundPerformanceBoard {
		t.Error("list did not find PerformanceBoard board")
	}
}
