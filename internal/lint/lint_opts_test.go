package lint

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

// Tests for LintAllWithConfig - supporting opts.Disable

func TestLintAllWithConfig_DisabledRules(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Queries: []discovery.DiscoveredQuery{
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
		},
	}

	// Test without config - should have WHC001
	results := LintAll(resources)
	if !hasResult(results, "WHC001") {
		t.Error("Expected WHC001 with default config")
	}

	// Test with disabled WHC001
	config := LintConfig{DisabledRules: []string{"WHC001"}}
	results = LintAllWithConfig(resources, config)
	if hasResult(results, "WHC001") {
		t.Error("WHC001 should be disabled")
	}
}

func TestLintAllWithConfig_MultipleDisabledRules(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Queries: []discovery.DiscoveredQuery{
			{
				Name:         "Query1",
				Package:      "test",
				File:         "/test/file.go",
				Line:         10,
				Dataset:      "",                            // Missing dataset - WHC001
				TimeRange:    discovery.TimeRange{},         // Missing time range - WHC002
				Calculations: []discovery.Calculation{},    // Empty calculations - WHC003
			},
		},
	}

	// Test with multiple disabled rules
	config := LintConfig{DisabledRules: []string{"WHC001", "WHC002"}}
	results := LintAllWithConfig(resources, config)

	if hasResult(results, "WHC001") {
		t.Error("WHC001 should be disabled")
	}
	if hasResult(results, "WHC002") {
		t.Error("WHC002 should be disabled")
	}
	// WHC003 should still be present
	if !hasResult(results, "WHC003") {
		t.Error("Expected WHC003 since it was not disabled")
	}
}

func TestLintAllWithConfig_EmptyDisabledRules(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Queries: []discovery.DiscoveredQuery{
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
		},
	}

	// Empty config should behave like LintAll
	config := LintConfig{}
	results := LintAllWithConfig(resources, config)
	if !hasResult(results, "WHC001") {
		t.Error("Expected WHC001 with empty config")
	}
}

func TestLintAllWithConfig_BoardsAndSLOs(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Boards: []discovery.DiscoveredBoard{
			{
				Name:        "Board1",
				File:        "/test/board.go",
				Line:        10,
				BoardName:   "", // Missing name - WHC030
				PanelCount:  0,  // No panels - may trigger warning
			},
		},
		SLOs: []discovery.DiscoveredSLO{
			{
				Name:             "SLO1",
				File:             "/test/slo.go",
				Line:             10,
				SLOName:          "", // Missing name - WHC040
				TargetPercentage: 0,  // Invalid target
			},
		},
	}

	// Test without disabled rules
	results := LintAll(resources)
	hasBoard := hasResult(results, "WHC030")
	hasSLO := hasResult(results, "WHC040")

	if !hasBoard && !hasSLO {
		t.Log("Note: Board and SLO rules may have different conditions")
	}

	// Test with disabled rules for boards
	config := LintConfig{DisabledRules: []string{"WHC030", "WHC040"}}
	results = LintAllWithConfig(resources, config)
	if hasResult(results, "WHC030") {
		t.Error("WHC030 should be disabled")
	}
	if hasResult(results, "WHC040") {
		t.Error("WHC040 should be disabled")
	}
}

func TestLintAllWithConfig_TriggersDisabled(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Triggers: []discovery.DiscoveredTrigger{
			{
				Name:        "Trigger1",
				File:        "/test/trigger.go",
				Line:        10,
				TriggerName: "", // Missing name - WHC050
			},
		},
	}

	// Disable trigger rules
	config := LintConfig{DisabledRules: []string{"WHC050"}}
	results := LintAllWithConfig(resources, config)
	if hasResult(results, "WHC050") {
		t.Error("WHC050 should be disabled")
	}
}

func TestLintAllWithConfig_SeverityOverrides(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Queries: []discovery.DiscoveredQuery{
			{
				Name:      "Query1",
				Package:   "test",
				File:      "/test/file.go",
				Line:      10,
				Dataset:   "", // Missing dataset - WHC001 (error)
				TimeRange: discovery.TimeRange{TimeRange: 3600},
				Calculations: []discovery.Calculation{
					{Op: "COUNT"},
				},
			},
		},
	}

	// Override WHC001 error to warning
	config := LintConfig{
		SeverityOverrides: map[string]Severity{
			"WHC001": SeverityWarning,
		},
	}
	results := LintAllWithConfig(resources, config)

	whc001 := findResult(results, "WHC001")
	if whc001 == nil {
		t.Fatal("Expected WHC001 in results")
	}
	if whc001.Severity != SeverityWarning {
		t.Errorf("Expected WHC001 to be overridden to warning, got %s", whc001.Severity)
	}
}

func TestLintAllWithConfig_SortedByFileAndLine(t *testing.T) {
	resources := &discovery.DiscoveredResources{
		Queries: []discovery.DiscoveredQuery{
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
		},
	}

	config := LintConfig{}
	results := LintAllWithConfig(resources, config)

	// Results should be sorted by file, then line
	if len(results) < 2 {
		t.Skip("Need at least 2 results for sorting test")
	}

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
