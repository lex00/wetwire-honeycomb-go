// Package lint provides query validation and linting capabilities.
package lint

import (
	"sort"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// LintResult represents a single lint finding.
type LintResult struct {
	Rule     string // e.g., "WHC001"
	Severity string // "error" or "warning"
	Message  string
	File     string
	Line     int
	Query    string // Query name
}

// LintQueries runs all lint rules against the provided queries.
// Results are sorted by file and line number.
func LintQueries(queries []discovery.DiscoveredQuery) []LintResult {
	return LintQueriesWithRules(queries, AllRules())
}

// LintQueriesWithRules runs specific lint rules against the provided queries.
// Results are sorted by file and line number.
func LintQueriesWithRules(queries []discovery.DiscoveredQuery, rules []Rule) []LintResult {
	var results []LintResult

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
	SeverityOverrides map[string]string
}

// LintQueriesWithConfig runs lint rules with the specified configuration.
func LintQueriesWithConfig(queries []discovery.DiscoveredQuery, config LintConfig) []LintResult {
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
func HasErrors(results []LintResult) bool {
	for _, r := range results {
		if r.Severity == "error" {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any lint results are warnings.
func HasWarnings(results []LintResult) bool {
	for _, r := range results {
		if r.Severity == "warning" {
			return true
		}
	}
	return false
}

// CountByRule counts lint results grouped by rule code.
func CountByRule(results []LintResult) map[string]int {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Rule]++
	}
	return counts
}

// CountBySeverity counts lint results grouped by severity.
func CountBySeverity(results []LintResult) map[string]int {
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.Severity]++
	}
	return counts
}

// FilterByRule filters lint results by rule code.
func FilterByRule(results []LintResult, rule string) []LintResult {
	var filtered []LintResult
	for _, r := range results {
		if r.Rule == rule {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// FilterBySeverity filters lint results by severity.
func FilterBySeverity(results []LintResult, severity string) []LintResult {
	var filtered []LintResult
	for _, r := range results {
		if r.Severity == severity {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// LintBoards runs all board lint rules against the provided boards.
// Results are sorted by file and line number.
func LintBoards(boards []discovery.DiscoveredBoard) []LintResult {
	var results []LintResult

	rules := AllBoardRules()
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
func LintSLOs(slos []discovery.DiscoveredSLO) []LintResult {
	var results []LintResult

	rules := AllSLORules()
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
func LintTriggers(triggers []discovery.DiscoveredTrigger) []LintResult {
	var results []LintResult

	rules := AllTriggerRules()
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
func LintAll(resources *discovery.DiscoveredResources) []LintResult {
	var results []LintResult

	results = append(results, LintQueries(resources.Queries)...)
	results = append(results, LintBoards(resources.Boards)...)
	results = append(results, LintSLOs(resources.SLOs)...)
	results = append(results, LintTriggers(resources.Triggers)...)

	// Sort all results together by file, then line
	sort.Slice(results, func(i, j int) bool {
		if results[i].File != results[j].File {
			return results[i].File < results[j].File
		}
		return results[i].Line < results[j].Line
	})

	return results
}
