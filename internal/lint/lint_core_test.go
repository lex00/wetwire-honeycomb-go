package lint

import (
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityWarning {
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
	if result.Severity != SeverityWarning {
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityWarning {
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
	if result.Severity != SeverityError {
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
	if result.Severity != SeverityWarning {
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
		if r.Severity == SeverityError {
			errors = append(errors, r)
		}
	}

	if len(errors) > 0 {
		t.Errorf("Expected no lint errors for valid query, got %d: %v", len(errors), errors)
	}
}
