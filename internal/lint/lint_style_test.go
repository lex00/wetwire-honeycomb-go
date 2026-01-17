package lint

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// WHC020 Inline Calculation Definition Tests

func TestLintQueries_WHC020_InlineCalculationDefinition_TooMany(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				InlineCalculationCount: 5, // More than threshold of 3
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC020") {
		t.Error("Expected WHC020 warning for too many inline calculation definitions")
	}

	result := findResult(results, "WHC020")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC020_InlineCalculationDefinition_OK(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				InlineCalculationCount: 2, // Within threshold
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC020") {
		t.Error("Did not expect WHC020 warning for acceptable inline calculation count")
	}
}

// WHC021 Inline Filter Definition Tests

func TestLintQueries_WHC021_InlineFilterDefinition_TooMany(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				InlineFilterCount: 4, // More than threshold of 3
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC021") {
		t.Error("Expected WHC021 warning for too many inline filter definitions")
	}

	result := findResult(results, "WHC021")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC021_InlineFilterDefinition_OK(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				InlineFilterCount: 3, // At threshold, should be OK
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC021") {
		t.Error("Did not expect WHC021 warning for acceptable inline filter count")
	}
}

// WHC022 Raw Map Literal Tests

func TestLintQueries_WHC022_RawMapLiteral_Present(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				HasRawMapLiteral: true,
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC022") {
		t.Error("Expected WHC022 warning for raw map literal")
	}

	result := findResult(results, "WHC022")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC022_RawMapLiteral_NotPresent(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				HasRawMapLiteral: false,
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC022") {
		t.Error("Did not expect WHC022 warning when no raw map literals")
	}
}

// WHC023 Deeply Nested Configuration Tests

func TestLintQueries_WHC023_DeeplyNested_TooDeep(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				MaxNestingDepth: 5, // More than threshold of 4
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC023") {
		t.Error("Expected WHC023 warning for deeply nested configuration")
	}

	result := findResult(results, "WHC023")
	if result.Severity != SeverityWarning {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC023_DeeplyNested_OK(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				MaxNestingDepth: 4, // At threshold, should be OK
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC023") {
		t.Error("Did not expect WHC023 warning for acceptable nesting depth")
	}
}

func TestLintQueries_WHC023_DeeplyNested_VeryDeep(t *testing.T) {
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
			Style: discovery.StyleMetadata{
				MaxNestingDepth: 10, // Very deep nesting
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC023") {
		t.Error("Expected WHC023 warning for very deep nesting")
	}

	result := findResult(results, "WHC023")
	if result.Message == "" {
		t.Error("Expected WHC023 result to have a message")
	}
}

// Test AllRules includes style enforcement rules
func TestAllRules_IncludesStyleEnforcementRules(t *testing.T) {
	rules := AllRules()

	ruleMap := make(map[string]bool)
	for _, r := range rules {
		ruleMap[r.Code] = true
	}

	expectedRules := []string{"WHC020", "WHC021", "WHC022", "WHC023"}
	for _, code := range expectedRules {
		if !ruleMap[code] {
			t.Errorf("Expected AllRules to include %s", code)
		}
	}
}
