package lint

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

func TestLintQueries_WHC011_CircularDependency_FilterSelfReference(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "UserMetrics",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Filters: []discovery.Filter{
				{Column: "usermetrics_id", Op: "=", Value: "123"}, // Self-reference in filter
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC011") {
		t.Error("Expected WHC011 warning for circular dependency in filter")
	}

	result := findResult(results, "WHC011")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC011_CircularDependency_CalculationSelfReference(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "ErrorRate",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "AVG", Column: "errorrate_value"}, // Self-reference in calculation
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC011") {
		t.Error("Expected WHC011 warning for circular dependency in calculation")
	}

	result := findResult(results, "WHC011")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC011_CircularDependency_NoSelfReference(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "LatencyQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "P99", Column: "duration_ms"},
			},
			Filters: []discovery.Filter{
				{Column: "service", Op: "=", Value: "api"},
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC011") {
		t.Error("Did not expect WHC011 warning for query without self-reference")
	}
}

func TestLintQueries_WHC011_CircularDependency_ShortQueryName(t *testing.T) {
	// Query names shorter than 3 characters should be skipped
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "QQ", // Short name, should be skipped
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
			Filters: []discovery.Filter{
				{Column: "qq_field", Op: "=", Value: "test"},
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC011") {
		t.Error("Did not expect WHC011 warning for short query name")
	}
}

func TestLintQueries_WHC011_CircularDependency_EmptyColumn(t *testing.T) {
	// Calculations with empty columns should be skipped
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT", Column: ""}, // Empty column, should not trigger
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC011") {
		t.Error("Did not expect WHC011 warning for calculation with empty column")
	}
}
