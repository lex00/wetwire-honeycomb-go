// Package lint provides query validation and linting capabilities.
package lint

import (
	"sort"

	corelint "github.com/lex00/wetwire-core-go/lint"
	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

// Severity type alias and constants from wetwire-core-go/lint.
type Severity = corelint.Severity

const (
	SeverityError   = corelint.SeverityError
	SeverityWarning = corelint.SeverityWarning
	SeverityInfo    = corelint.SeverityInfo
)

// Issue is a type alias to the shared Issue type from wetwire-core-go/lint.
type Issue = corelint.Issue

// LintQueries runs all lint rules against the provided queries.
// Results are sorted by file and line number.
func LintQueries(queries []discovery.DiscoveredQuery) []Issue {
	return LintQueriesWithRules(queries, AllRules())
}

// LintQueriesWithRules runs specific lint rules against the provided queries.
// Results are sorted by file and line number.
func LintQueriesWithRules(queries []discovery.DiscoveredQuery, rules []Rule) []Issue {
	var results []Issue

	for _, query := range queries {
		for _, rule := range rules {
			ruleResults := rule.Check(query)
			results = append(results, ruleResults...)
		}
	}

	// Sort results by file, then line
	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}

// LintConfig holds configuration options for linting.
type LintConfig struct {
	// DisabledRules is a list of rule codes to skip
	DisabledRules []string

	// SeverityOverrides maps rule codes to custom severity levels
	SeverityOverrides map[string]Severity
}

// LintQueriesWithConfig runs lint rules with the specified configuration.
func LintQueriesWithConfig(queries []discovery.DiscoveredQuery, config LintConfig) []Issue {
	// Get all rules
	rules := AllRules()

	// Filter out disabled rules
	disabledSet := make(map[string]bool)
	for _, code := range config.DisabledRules {
		disabledSet[code] = true
	}

	var enabledRules []Rule
	for _, rule := range rules {
		if !disabledSet[rule.Code] {
			enabledRules = append(enabledRules, rule)
		}
	}

	// Run lint with enabled rules
	results := LintQueriesWithRules(queries, enabledRules)

	// Apply severity overrides
	for i := range results {
		if newSeverity, ok := config.SeverityOverrides[results[i].Rule]; ok {
			results[i].Severity = newSeverity
		}
	}

	return results
}

// HasErrors returns true if any lint results are errors.
func HasErrors(results []Issue) bool {
	for _, r := range results {
		if r.Severity == SeverityError {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any lint results are warnings.
func HasWarnings(results []Issue) bool {
	for _, r := range results {
		if r.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

// CountByRule counts lint results grouped by rule code.
func CountByRule(results []Issue) map[string]int {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Rule]++
	}
	return counts
}

// CountBySeverity counts lint results grouped by severity.
func CountBySeverity(results []Issue) map[string]int {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Severity.String()]++
	}
	return counts
}

// FilterByRule filters lint results by rule code.
func FilterByRule(results []Issue, rule string) []Issue {
	var filtered []Issue
	for _, r := range results {
		if r.Rule == rule {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// FilterBySeverity filters lint results by severity.
func FilterBySeverity(results []Issue, severity Severity) []Issue {
	var filtered []Issue
	for _, r := range results {
		if r.Severity == severity {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// LintBoards runs all board lint rules against the provided boards.
// Results are sorted by file and line number.
func LintBoards(boards []discovery.DiscoveredBoard) []Issue {
	return LintBoardsWithRules(boards, AllBoardRules())
}

// LintBoardsWithRules runs specific board lint rules against the provided boards.
// Results are sorted by file and line number.
func LintBoardsWithRules(boards []discovery.DiscoveredBoard, rules []BoardRule) []Issue {
	var results []Issue

	for _, board := range boards {
		for _, rule := range rules {
			ruleResults := rule.Check(board)
			results = append(results, ruleResults...)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}

// LintSLOs runs all SLO lint rules against the provided SLOs.
// Results are sorted by file and line number.
func LintSLOs(slos []discovery.DiscoveredSLO) []Issue {
	return LintSLOsWithRules(slos, AllSLORules())
}

// LintSLOsWithRules runs specific SLO lint rules against the provided SLOs.
// Results are sorted by file and line number.
func LintSLOsWithRules(slos []discovery.DiscoveredSLO, rules []SLORule) []Issue {
	var results []Issue

	for _, slo := range slos {
		for _, rule := range rules {
			ruleResults := rule.Check(slo)
			results = append(results, ruleResults...)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}

// LintTriggers runs all trigger lint rules against the provided triggers.
// Results are sorted by file and line number.
func LintTriggers(triggers []discovery.DiscoveredTrigger) []Issue {
	return LintTriggersWithRules(triggers, AllTriggerRules())
}

// LintTriggersWithRules runs specific trigger lint rules against the provided triggers.
// Results are sorted by file and line number.
func LintTriggersWithRules(triggers []discovery.DiscoveredTrigger, rules []TriggerRule) []Issue {
	var results []Issue

	for _, trigger := range triggers {
		for _, rule := range rules {
			ruleResults := rule.Check(trigger)
			results = append(results, ruleResults...)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}

// LintAll runs all lint rules against all discovered resources.
// Results are sorted by file and line number.
func LintAll(resources *discovery.DiscoveredResources) []Issue {
	return LintAllWithConfig(resources, LintConfig{})
}

// LintAllWithConfig runs all lint rules against all discovered resources with configuration.
// It respects DisabledRules and SeverityOverrides from the config.
// Results are sorted by file and line number.
func LintAllWithConfig(resources *discovery.DiscoveredResources, config LintConfig) []Issue {
	var results []Issue

	// Build disabled rules set
	disabledSet := make(map[string]bool)
	for _, code := range config.DisabledRules {
		disabledSet[code] = true
	}

	// Filter query rules
	queryRules := AllRules()
	var enabledQueryRules []Rule
	for _, rule := range queryRules {
		if !disabledSet[rule.Code] {
			enabledQueryRules = append(enabledQueryRules, rule)
		}
	}
	results = append(results, LintQueriesWithRules(resources.Queries, enabledQueryRules)...)

	// Filter board rules
	boardRules := AllBoardRules()
	var enabledBoardRules []BoardRule
	for _, rule := range boardRules {
		if !disabledSet[rule.Code] {
			enabledBoardRules = append(enabledBoardRules, rule)
		}
	}
	results = append(results, LintBoardsWithRules(resources.Boards, enabledBoardRules)...)

	// Filter SLO rules
	sloRules := AllSLORules()
	var enabledSLORules []SLORule
	for _, rule := range sloRules {
		if !disabledSet[rule.Code] {
			enabledSLORules = append(enabledSLORules, rule)
		}
	}
	results = append(results, LintSLOsWithRules(resources.SLOs, enabledSLORules)...)

	// Filter trigger rules
	triggerRules := AllTriggerRules()
	var enabledTriggerRules []TriggerRule
	for _, rule := range triggerRules {
		if !disabledSet[rule.Code] {
			enabledTriggerRules = append(enabledTriggerRules, rule)
		}
	}
	results = append(results, LintTriggersWithRules(resources.Triggers, enabledTriggerRules)...)

	// Apply severity overrides
	for i := range results {
		if newSeverity, ok := config.SeverityOverrides[results[i].Rule]; ok {
			results[i].Severity = newSeverity
		}
	}

	// Sort all results together by file, then line
	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}
