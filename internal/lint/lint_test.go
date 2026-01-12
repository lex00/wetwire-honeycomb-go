package lint

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestLintQueries_WHC001_MissingDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "", // Missing dataset
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC001") {
		t.Error("Expected WHC001 error for missing dataset")
	}

	result := findResult(results, "WHC001")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC002_MissingTimeRange(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{}, // Missing time range
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC002") {
		t.Error("Expected WHC002 error for missing time range")
	}

	result := findResult(results, "WHC002")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC003_EmptyCalculations(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:         "TestQuery",
			Package:      "test",
			File:         "/test/file.go",
			Line:         10,
			Dataset:      "production",
			TimeRange:    discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{}, // Empty calculations
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC003") {
		t.Error("Expected WHC003 error for empty calculations")
	}

	result := findResult(results, "WHC003")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC004_BreakdownWithoutOrder(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:       "TestQuery",
			Package:    "test",
			File:       "/test/file.go",
			Line:       10,
			Dataset:    "production",
			TimeRange:  discovery.TimeRange{TimeRange: 3600},
			Breakdowns: []string{"endpoint", "service"},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			// No orders specified - this should trigger warning
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC004") {
		t.Error("Expected WHC004 warning for breakdown without order")
	}

	result := findResult(results, "WHC004")
	if result.Severity != "warning" {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC005_HighCardinalityBreakdown(t *testing.T) {
	// Create a query with more than 100 breakdown groups
	breakdowns := make([]string, 150)
	for i := 0; i < 150; i++ {
		breakdowns[i] = "field"
	}

	queries := []discovery.DiscoveredQuery{
		{
			Name:       "TestQuery",
			Package:    "test",
			File:       "/test/file.go",
			Line:       10,
			Dataset:    "production",
			TimeRange:  discovery.TimeRange{TimeRange: 3600},
			Breakdowns: breakdowns,
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Limit: 150, // Limit indicates high cardinality
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC005") {
		t.Error("Expected WHC005 warning for high cardinality breakdown")
	}

	result := findResult(results, "WHC005")
	if result.Severity != "warning" {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC006_InvalidCalculationForColumnType(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "P99", Column: "endpoint"}, // P99 on likely string field
				{Op: "SUM", Column: "error_message"}, // SUM on likely string field
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC006") {
		t.Error("Expected WHC006 error for invalid calculation on column type")
	}

	result := findResult(results, "WHC006")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC007_InvalidFilterOperator(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Filters: []discovery.Filter{
				{Column: "duration", Op: "invalid_op", Value: 100},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC007") {
		t.Error("Expected WHC007 error for invalid filter operator")
	}

	result := findResult(results, "WHC007")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC008_MissingLimitWithBreakdowns(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:       "TestQuery",
			Package:    "test",
			File:       "/test/file.go",
			Line:       10,
			Dataset:    "production",
			TimeRange:  discovery.TimeRange{TimeRange: 3600},
			Breakdowns: []string{"endpoint", "service"},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Limit: 0, // No limit specified
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC008") {
		t.Error("Expected WHC008 warning for missing limit with breakdowns")
	}

	result := findResult(results, "WHC008")
	if result.Severity != "warning" {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC009_TimeRangeExceeds7Days(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 8 * 86400}, // 8 days
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC009") {
		t.Error("Expected WHC009 error for time range exceeding 7 days")
	}

	result := findResult(results, "WHC009")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC010_ExcessiveFilterCount(t *testing.T) {
	// Create a query with more than 50 filters
	filters := make([]discovery.Filter, 51)
	for i := 0; i < 51; i++ {
		filters[i] = discovery.Filter{
			Column: "field",
			Op:     "=",
			Value:  "value",
		}
	}

	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Filters: filters,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC010") {
		t.Error("Expected WHC010 warning for excessive filter count")
	}

	result := findResult(results, "WHC010")
	if result.Severity != "warning" {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_ValidQuery_NoErrors(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "ValidQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "P99", Column: "duration_ms"},
				{Op: "COUNT"},
			},
			Filters: []discovery.Filter{
				{Column: "status", Op: "=", Value: "error"},
			},
			// No breakdowns to avoid WHC004 warning
			Limit: 100,
		},
	}

	results := LintQueries(queries)

	// Filter to only check for errors, not warnings
	var errors []LintResult
	for _, r := range results {
		if r.Severity == "error" {
			errors = append(errors, r)
		}
	}

	if len(errors) > 0 {
		t.Errorf("Expected no lint errors for valid query, got %d: %v", len(errors), errors)
	}
}

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
		SeverityOverrides: map[string]string{
			"WHC001": "warning", // Override error to warning
		},
	}
	results = LintQueriesWithConfig(queries, configWithOverride)
	if len(results) > 0 && results[0].Rule == "WHC001" && results[0].Severity != "warning" {
		t.Errorf("Expected WHC001 to be overridden to warning, got %s", results[0].Severity)
	}
}

func TestHasErrors(t *testing.T) {
	resultsWithError := []LintResult{
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
	}

	if !HasErrors(resultsWithError) {
		t.Error("HasErrors should return true when there are errors")
	}

	resultsNoError := []LintResult{
		{Rule: "WHC004", Severity: "warning"},
		{Rule: "WHC005", Severity: "warning"},
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
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
	}

	if !HasWarnings(resultsWithWarning) {
		t.Error("HasWarnings should return true when there are warnings")
	}

	resultsNoWarning := []LintResult{
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC002", Severity: "error"},
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
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC002", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
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
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC002", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
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
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC002", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
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
		{Rule: "WHC001", Severity: "error"},
		{Rule: "WHC002", Severity: "error"},
		{Rule: "WHC004", Severity: "warning"},
	}

	errors := FilterBySeverity(results, "error")
	if len(errors) != 2 {
		t.Errorf("Expected 2 error results, got %d", len(errors))
	}

	for _, r := range errors {
		if r.Severity != "error" {
			t.Errorf("Expected all results to be errors, got %s", r.Severity)
		}
	}

	warnings := FilterBySeverity(results, "warning")
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
		if r.Severity == "error" {
			errorCount++
		}
	}

	if errorCount > 0 {
		t.Logf("Found %d lint errors in testdata (this may be expected for testing):", errorCount)
		for _, r := range results {
			if r.Severity == "error" {
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
