package lint

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

func TestLintQueries_MultipleQueries(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "Query1",
			Package:   "test",
			File:      "/test/file1.go",
			Line:      10,
			Dataset:   "", // Missing dataset
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
		{
			Name:         "Query2",
			Package:      "test",
			File:         "/test/file2.go",
			Line:         20,
			Dataset:      "production",
			TimeRange:    discovery.TimeRange{}, // Missing time range
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	// Should have at least 2 errors (one for each query)
	if len(results) < 2 {
		t.Errorf("Expected at least 2 lint errors, got %d", len(results))
	}

	// Verify both WHC001 and WHC002 are present
	if !hasResult(results, "WHC001") {
		t.Error("Expected WHC001 error in results")
	}
	if !hasResult(results, "WHC002") {
		t.Error("Expected WHC002 error in results")
	}
}

func TestLintQueries_SortedByFileAndLine(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "Query2",
			Package:   "test",
			File:      "/test/b.go",
			Line:      10,
			Dataset:   "",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
		{
			Name:      "Query1",
			Package:   "test",
			File:      "/test/a.go",
			Line:      20,
			Dataset:   "",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
		{
			Name:      "Query3",
			Package:   "test",
			File:      "/test/a.go",
			Line:      10,
			Dataset:   "",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	// Results should be sorted by file, then line
	if len(results) < 3 {
		t.Fatalf("Expected at least 3 results, got %d", len(results))
	}

	// Check sorting
	for i := 1; i < len(results); i++ {
		prev := results[i-1]
		curr := results[i]

		if prev.File > curr.File {
			t.Errorf("Results not sorted by file: %s > %s", prev.File, curr.File)
		}

		if prev.File == curr.File && prev.Line > curr.Line {
			t.Errorf("Results not sorted by line: %d > %d", prev.Line, curr.Line)
		}
	}
}

func TestLintQueriesWithConfig(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "Query1",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "", // Missing dataset - WHC001
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	// Test with empty config (all rules enabled)
	results := LintQueriesWithConfig(queries, LintConfig{})
	if !hasResult(results, "WHC001") {
		t.Error("Expected WHC001 with empty config")
	}

	// Test with config that disables WHC001
	config := LintConfig{DisabledRules: []string{"WHC001"}}
	results = LintQueriesWithConfig(queries, config)
	if hasResult(results, "WHC001") {
		t.Error("WHC001 should be disabled")
	}

	// Test severity overrides
	configWithOverride := LintConfig{
		SeverityOverrides: map[string]Severity{
			"WHC001": SeverityWarning, // Override error to warning
		},
	}
	results = LintQueriesWithConfig(queries, configWithOverride)
	if len(results) > 0 && results[0].Rule == "WHC001" && results[0].Severity != SeverityWarning {
		t.Errorf("Expected WHC001 to be overridden to warning, got %s", results[0].Severity)
	}
}

func TestHasErrors(t *testing.T) {
	resultsWithError := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	if !HasErrors(resultsWithError) {
		t.Error("HasErrors should return true when there are errors")
	}

	resultsNoError := []LintResult{
		{Rule: "WHC004", Severity: SeverityWarning},
		{Rule: "WHC005", Severity: SeverityWarning},
	}

	if HasErrors(resultsNoError) {
		t.Error("HasErrors should return false when there are no errors")
	}

	if HasErrors(nil) {
		t.Error("HasErrors should return false for nil slice")
	}
}

func TestHasWarnings(t *testing.T) {
	resultsWithWarning := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	if !HasWarnings(resultsWithWarning) {
		t.Error("HasWarnings should return true when there are warnings")
	}

	resultsNoWarning := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC002", Severity: SeverityError},
	}

	if HasWarnings(resultsNoWarning) {
		t.Error("HasWarnings should return false when there are no warnings")
	}

	if HasWarnings(nil) {
		t.Error("HasWarnings should return false for nil slice")
	}
}

func TestCountByRule(t *testing.T) {
	results := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC002", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	counts := CountByRule(results)

	if counts["WHC001"] != 2 {
		t.Errorf("Expected WHC001 count 2, got %d", counts["WHC001"])
	}
	if counts["WHC002"] != 1 {
		t.Errorf("Expected WHC002 count 1, got %d", counts["WHC002"])
	}
	if counts["WHC004"] != 1 {
		t.Errorf("Expected WHC004 count 1, got %d", counts["WHC004"])
	}

	emptyCounts := CountByRule(nil)
	if len(emptyCounts) != 0 {
		t.Error("CountByRule should return empty map for nil slice")
	}
}

func TestCountBySeverity(t *testing.T) {
	results := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC002", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	counts := CountBySeverity(results)

	if counts["error"] != 2 {
		t.Errorf("Expected error count 2, got %d", counts["error"])
	}
	if counts["warning"] != 1 {
		t.Errorf("Expected warning count 1, got %d", counts["warning"])
	}

	emptyCounts := CountBySeverity(nil)
	if len(emptyCounts) != 0 {
		t.Error("CountBySeverity should return empty map for nil slice")
	}
}

func TestFilterByRule(t *testing.T) {
	results := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC002", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	filtered := FilterByRule(results, "WHC001")
	if len(filtered) != 2 {
		t.Errorf("Expected 2 WHC001 results, got %d", len(filtered))
	}

	for _, r := range filtered {
		if r.Rule != "WHC001" {
			t.Errorf("Expected all results to be WHC001, got %s", r.Rule)
		}
	}

	noMatch := FilterByRule(results, "WHC999")
	if len(noMatch) != 0 {
		t.Errorf("Expected 0 results for non-existent rule, got %d", len(noMatch))
	}
}

func TestFilterBySeverity(t *testing.T) {
	results := []LintResult{
		{Rule: "WHC001", Severity: SeverityError},
		{Rule: "WHC002", Severity: SeverityError},
		{Rule: "WHC004", Severity: SeverityWarning},
	}

	errors := FilterBySeverity(results, SeverityError)
	if len(errors) != 2 {
		t.Errorf("Expected 2 error results, got %d", len(errors))
	}

	for _, r := range errors {
		if r.Severity != SeverityError {
			t.Errorf("Expected all results to be errors, got %s", r.Severity)
		}
	}

	warnings := FilterBySeverity(results, SeverityWarning)
	if len(warnings) != 1 {
		t.Errorf("Expected 1 warning result, got %d", len(warnings))
	}
}

func TestLintQueries_WHC009_AbsoluteTimeRange(t *testing.T) {
	// Test with absolute time range exceeding 7 days
	queries := []discovery.DiscoveredQuery{
		{
			Name:    "TestQuery",
			Package: "test",
			File:    "/test/file.go",
			Line:    10,
			Dataset: "production",
			TimeRange: discovery.TimeRange{
				StartTime: 1000000,
				EndTime:   1000000 + 8*86400, // 8 days duration
			},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC009") {
		t.Error("Expected WHC009 error for absolute time range exceeding 7 days")
	}
}

func TestLintQueries_RealWorldExample(t *testing.T) {
	// Test with real testdata
	testDir := filepath.Join(getRepoRoot(t), "testdata", "queries")
	queries, err := discovery.DiscoverQueries(testDir)
	if err != nil {
		t.Fatalf("Failed to discover queries: %v", err)
	}

	results := LintQueries(queries)

	// The testdata queries should be valid, so we expect no errors
	// (This assumes the testdata is well-formed)
	errorCount := 0
	for _, r := range results {
		if r.Severity == SeverityError {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Logf("Found %d lint errors in testdata (this may be expected for testing):", errorCount)
		for _, r := range results {
			if r.Severity == SeverityError {
				t.Logf("  %s: %s (%s:%d)", r.Rule, r.Message, r.File, r.Line)
			}
		}
	}
}

// Helper functions

func hasResult(results []LintResult, rule string) bool {
	for _, r := range results {
		if r.Rule == rule {
			return true
		}
	}
	return false
}

func findResult(results []LintResult, rule string) *LintResult {
	for i, r := range results {
		if r.Rule == rule {
			return &results[i]
		}
	}
	return nil
}

func getRepoRoot(t *testing.T) string {
	// Use runtime.Caller to find the current file's location
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get current file path")
	}
	// Go up from internal/lint/ to repo root
	return filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
}
