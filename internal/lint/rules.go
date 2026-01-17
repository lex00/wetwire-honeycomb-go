package lint

import (
	"fmt"
	"strings"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// Rule represents a lint rule that can be applied to queries.
type Rule struct {
	Code     string
	Severity Severity // SeverityError, SeverityWarning, or SeverityInfo
	Message  string
	Check    func(query discovery.DiscoveredQuery) []LintResult
}

// AllRules returns all available lint rules.
func AllRules() []Rule {
	return []Rule{
		WHC001MissingDataset(),
		WHC002MissingTimeRange(),
		WHC003EmptyCalculations(),
		WHC004BreakdownWithoutOrder(),
		WHC005HighCardinalityBreakdown(),
		WHC006InvalidCalculationForColumnType(),
		WHC007InvalidFilterOperator(),
		WHC008MissingLimitWithBreakdowns(),
		WHC009TimeRangeExceeds7Days(),
		WHC010ExcessiveFilterCount(),
		WHC011CircularDependency(),
		WHC012SecretInFilter(),
		WHC013SensitiveColumnExposure(),
		WHC014HardcodedCredentials(),
		WHC020InlineCalculationDefinition(),
		WHC021InlineFilterDefinition(),
		WHC022RawMapLiteral(),
		WHC023DeeplyNestedConfiguration(),
	}
}

// WHC001MissingDataset checks if a query is missing a dataset.
func WHC001MissingDataset() Rule {
	return Rule{
		Code:     "WHC001",
		Severity: SeverityError,
		Message:  "Query is missing dataset",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			if query.Dataset == "" {
				return []LintResult{
					{
						Rule:     "WHC001",
						Severity: SeverityError,
						Message:  "Query is missing dataset",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC002MissingTimeRange checks if a query is missing a time range.
func WHC002MissingTimeRange() Rule {
	return Rule{
		Code:     "WHC002",
		Severity: SeverityError,
		Message:  "Query is missing time_range",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// Time range is missing if all time fields are zero
			if query.TimeRange.TimeRange == 0 && query.TimeRange.StartTime == 0 && query.TimeRange.EndTime == 0 {
				return []LintResult{
					{
						Rule:     "WHC002",
						Severity: SeverityError,
						Message:  "Query is missing time_range",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC003EmptyCalculations checks if a query has no calculations.
func WHC003EmptyCalculations() Rule {
	return Rule{
		Code:     "WHC003",
		Severity: SeverityError,
		Message:  "Query has empty calculations",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			if len(query.Calculations) == 0 {
				return []LintResult{
					{
						Rule:     "WHC003",
						Severity: SeverityError,
						Message:  "Query has empty calculations",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC004BreakdownWithoutOrder checks if a query has breakdowns but no order specified.
func WHC004BreakdownWithoutOrder() Rule {
	return Rule{
		Code:     "WHC004",
		Severity: SeverityWarning,
		Message:  "Query has breakdowns but no order specified",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// For now, we check if there are breakdowns present
			// In the real implementation, we'd check if Orders field is empty
			// Since DiscoveredQuery doesn't have Orders field, we assume missing
			if len(query.Breakdowns) > 0 {
				// This is a simplified check - in reality you'd need to verify
				// that the query struct doesn't have orders
				return []LintResult{
					{
						Rule:     "WHC004",
						Severity: SeverityWarning,
						Message:  "Query has breakdowns but no order specified - results may be unpredictable",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC005HighCardinalityBreakdown checks if a query has high cardinality breakdowns.
func WHC005HighCardinalityBreakdown() Rule {
	return Rule{
		Code:     "WHC005",
		Severity: SeverityWarning,
		Message:  "Query has high cardinality breakdown (>100 groups)",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// High cardinality is indicated by a limit > 100
			if query.Limit > 100 {
				return []LintResult{
					{
						Rule:     "WHC005",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Query has high cardinality breakdown (limit=%d > 100 groups)", query.Limit),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC006InvalidCalculationForColumnType checks if calculations are appropriate for column types.
func WHC006InvalidCalculationForColumnType() Rule {
	return Rule{
		Code:     "WHC006",
		Severity: SeverityError,
		Message:  "Invalid calculation for column type",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			// Numeric operations that shouldn't be used on string columns
			numericOps := map[string]bool{
				"P50": true, "P75": true, "P90": true, "P95": true,
				"P99": true, "P999": true, "SUM": true, "AVG": true,
				"MIN": true, "MAX": true, "HEATMAP": true,
			}

			// Common string field patterns (heuristic-based detection)
			stringPatterns := []string{
				"name", "message", "error", "status", "endpoint",
				"path", "url", "type", "service", "env", "environment",
			}

			for _, calc := range query.Calculations {
				if calc.Column == "" {
					continue
				}

				if numericOps[calc.Op] {
					// Check if column name suggests it's a string field
					columnLower := strings.ToLower(calc.Column)
					for _, pattern := range stringPatterns {
						if strings.Contains(columnLower, pattern) && !strings.Contains(columnLower, "_ms") && !strings.Contains(columnLower, "_bytes") && !strings.Contains(columnLower, "_count") {
							results = append(results, LintResult{
								Rule:     "WHC006",
								Severity: SeverityError,
								Message:  fmt.Sprintf("Calculation %s should not be used on likely string column '%s'", calc.Op, calc.Column),
								File:     query.File,
								Line:     query.Line,
								Query:    query.Name,
							})
							break
						}
					}
				}
			}

			return results
		},
	}
}

// WHC007InvalidFilterOperator checks if filter operators are valid.
func WHC007InvalidFilterOperator() Rule {
	return Rule{
		Code:     "WHC007",
		Severity: SeverityError,
		Message:  "Invalid filter operator",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			validOps := map[string]bool{
				"=":                true,
				"!=":               true,
				">":                true,
				">=":               true,
				"<":                true,
				"<=":               true,
				"contains":         true,
				"does-not-contain": true,
				"exists":           true,
				"does-not-exist":   true,
				"starts-with":      true,
				"in":               true,
				"not-in":           true,
			}

			for _, filter := range query.Filters {
				if !validOps[filter.Op] {
					results = append(results, LintResult{
						Rule:     "WHC007",
						Severity: SeverityError,
						Message:  fmt.Sprintf("Invalid filter operator '%s' on column '%s'", filter.Op, filter.Column),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					})
				}
			}

			return results
		},
	}
}

// WHC008MissingLimitWithBreakdowns checks if a query with breakdowns has no limit.
func WHC008MissingLimitWithBreakdowns() Rule {
	return Rule{
		Code:     "WHC008",
		Severity: SeverityWarning,
		Message:  "Query has breakdowns but no limit specified",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			if len(query.Breakdowns) > 0 && query.Limit == 0 {
				return []LintResult{
					{
						Rule:     "WHC008",
						Severity: SeverityWarning,
						Message:  "Query has breakdowns but no limit specified - may return too many results",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC009TimeRangeExceeds7Days checks if time range exceeds 7 days.
func WHC009TimeRangeExceeds7Days() Rule {
	return Rule{
		Code:     "WHC009",
		Severity: SeverityError,
		Message:  "Time range exceeds 7 days",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			const sevenDays = 7 * 86400 // 7 days in seconds

			if query.TimeRange.TimeRange > sevenDays {
				days := query.TimeRange.TimeRange / 86400
				return []LintResult{
					{
						Rule:     "WHC009",
						Severity: SeverityError,
						Message:  fmt.Sprintf("Time range exceeds 7 days (current: %d days)", days),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}

			// Also check absolute time ranges
			if query.TimeRange.StartTime > 0 && query.TimeRange.EndTime > 0 {
				diff := query.TimeRange.EndTime - query.TimeRange.StartTime
				if diff > sevenDays {
					days := diff / 86400
					return []LintResult{
						{
							Rule:     "WHC009",
							Severity: SeverityError,
							Message:  fmt.Sprintf("Time range exceeds 7 days (current: %d days)", days),
							File:     query.File,
							Line:     query.Line,
							Query:    query.Name,
						},
					}
				}
			}

			return nil
		},
	}
}

// WHC010ExcessiveFilterCount checks if a query has too many filters.
func WHC010ExcessiveFilterCount() Rule {
	return Rule{
		Code:     "WHC010",
		Severity: SeverityWarning,
		Message:  "Query has excessive filter count (>50)",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			const maxFilters = 50

			if len(query.Filters) > maxFilters {
				return []LintResult{
					{
						Rule:     "WHC010",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Query has excessive filter count (%d > %d)", len(query.Filters), maxFilters),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}

			return nil
		},
	}
}

// WHC011CircularDependency checks for potential circular dependencies in queries.
// In Honeycomb queries, circular dependencies can occur when:
// - A query references itself through derived columns or query composition
// - Multiple queries reference each other creating a logical loop
//
// Since individual DiscoveredQuery objects don't contain explicit references to other
// queries, this rule currently serves as a placeholder that will be extended when
// cross-query analysis is implemented. For now, it checks for self-referential patterns
// where a query's name appears in its own filter or calculation column names.
func WHC011CircularDependency() Rule {
	return Rule{
		Code:     "WHC011",
		Severity: SeverityWarning,
		Message:  "Potential circular dependency detected",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			// Check for self-referential patterns where the query name
			// appears in filter columns or calculation columns
			queryNameLower := strings.ToLower(query.Name)

			// Skip if query name is empty or too short to be meaningful
			if len(queryNameLower) < 3 {
				return nil
			}

			// Check filters for self-references
			for _, filter := range query.Filters {
				columnLower := strings.ToLower(filter.Column)
				if strings.Contains(columnLower, queryNameLower) {
					results = append(results, LintResult{
						Rule:     "WHC011",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Potential circular dependency: filter column '%s' references query name '%s'", filter.Column, query.Name),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					})
				}
			}

			// Check calculations for self-references
			for _, calc := range query.Calculations {
				if calc.Column == "" {
					continue
				}
				columnLower := strings.ToLower(calc.Column)
				if strings.Contains(columnLower, queryNameLower) {
					results = append(results, LintResult{
						Rule:     "WHC011",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Potential circular dependency: calculation column '%s' references query name '%s'", calc.Column, query.Name),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					})
				}
			}

			return results
		},
	}
}

// WHC012SecretInFilter detects potential secrets or tokens in filter values.
// This rule helps prevent accidental exposure of sensitive credentials in queries.
func WHC012SecretInFilter() Rule {
	return Rule{
		Code:     "WHC012",
		Severity: SeverityError,
		Message:  "Potential secret detected in filter value",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			// Patterns that indicate secret/token values
			secretPatterns := []string{
				"sk-",      // OpenAI, Stripe keys
				"sk_live_", // Stripe live keys
				"sk_test_", // Stripe test keys
				"api_key",
				"apikey",
				"api-key",
				"bearer",
				"token",
				"secret",
				"password",
				"passwd",
				"credential",
				"private_key",
				"private-key",
				"access_key",
				"access-key",
				"auth_token",
				"auth-token",
			}

			for _, filter := range query.Filters {
				// Check the filter value if it's a string
				valueStr, ok := filter.Value.(string)
				if !ok {
					continue
				}

				valueLower := strings.ToLower(valueStr)
				for _, pattern := range secretPatterns {
					if strings.Contains(valueLower, pattern) {
						results = append(results, LintResult{
							Rule:     "WHC012",
							Severity: SeverityError,
							Message:  fmt.Sprintf("Potential secret detected in filter value for column '%s' (pattern: %s)", filter.Column, pattern),
							File:     query.File,
							Line:     query.Line,
							Query:    query.Name,
						})
						break // Only report once per filter
					}
				}
			}

			return results
		},
	}
}

// WHC013SensitiveColumnExposure warns when querying columns that might contain PII.
// This rule helps identify potential privacy concerns in query breakdowns.
func WHC013SensitiveColumnExposure() Rule {
	return Rule{
		Code:     "WHC013",
		Severity: SeverityWarning,
		Message:  "Query may expose sensitive/PII column data",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			// Patterns that indicate PII or sensitive data columns
			sensitivePatterns := []string{
				"password",
				"passwd",
				"ssn",
				"social_security",
				"social-security",
				"socialsecurity",
				"credit_card",
				"credit-card",
				"creditcard",
				"card_number",
				"card-number",
				"cardnumber",
				"cvv",
				"pin",
				"secret",
				"private_key",
				"private-key",
				"privatekey",
				"auth_token",
				"auth-token",
				"authtoken",
				"api_key",
				"api-key",
				"apikey",
				"access_token",
				"access-token",
				"accesstoken",
			}

			// Check breakdown columns for sensitive data patterns
			for _, breakdown := range query.Breakdowns {
				breakdownLower := strings.ToLower(breakdown)
				for _, pattern := range sensitivePatterns {
					if strings.Contains(breakdownLower, pattern) {
						results = append(results, LintResult{
							Rule:     "WHC013",
							Severity: SeverityWarning,
							Message:  fmt.Sprintf("Breakdown column '%s' may expose sensitive/PII data (pattern: %s)", breakdown, pattern),
							File:     query.File,
							Line:     query.Line,
							Query:    query.Name,
						})
						break // Only report once per column
					}
				}
			}

			return results
		},
	}
}

// WHC014HardcodedCredentials detects hardcoded credentials in dataset names.
// This rule catches cases where credentials might be accidentally embedded in dataset identifiers.
func WHC014HardcodedCredentials() Rule {
	return Rule{
		Code:     "WHC014",
		Severity: SeverityError,
		Message:  "Hardcoded credentials detected in dataset name",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			var results []LintResult

			// Patterns that indicate hardcoded credentials in dataset names
			credentialPatterns := []string{
				"password=",
				"passwd=",
				"key=",
				"token=",
				"secret=",
				"apikey=",
				"api_key=",
				"api-key=",
				"access_key=",
				"access-key=",
				"auth=",
				"credential=",
			}

			datasetLower := strings.ToLower(query.Dataset)
			for _, pattern := range credentialPatterns {
				if strings.Contains(datasetLower, pattern) {
					results = append(results, LintResult{
						Rule:     "WHC014",
						Severity: SeverityError,
						Message:  fmt.Sprintf("Hardcoded credentials detected in dataset name (pattern: %s)", pattern),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					})
					break // Only report once per dataset
				}
			}

			return results
		},
	}
}

// WHC020InlineCalculationDefinition detects inline calculation definitions that should be
// extracted to named variables for better readability and reusability.
func WHC020InlineCalculationDefinition() Rule {
	return Rule{
		Code:     "WHC020",
		Severity: SeverityWarning,
		Message:  "Extract inline Calculation definitions to named variables",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// Threshold: more than 3 inline calculations suggest extraction
			const threshold = 3

			if query.Style.InlineCalculationCount > threshold {
				return []LintResult{
					{
						Rule:     "WHC020",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Query has %d inline Calculation definitions - consider extracting to named variables for reusability", query.Style.InlineCalculationCount),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC021InlineFilterDefinition detects inline filter definitions that should be
// extracted to named variables for better readability and reusability.
func WHC021InlineFilterDefinition() Rule {
	return Rule{
		Code:     "WHC021",
		Severity: SeverityWarning,
		Message:  "Extract inline Filter definitions to named variables",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// Threshold: more than 3 inline filters suggest extraction
			const threshold = 3

			if query.Style.InlineFilterCount > threshold {
				return []LintResult{
					{
						Rule:     "WHC021",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Query has %d inline Filter definitions - consider extracting to named variables for reusability", query.Style.InlineFilterCount),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC022RawMapLiteral detects raw map literals used instead of typed query builders.
// Using typed builders provides better type safety and IDE support.
func WHC022RawMapLiteral() Rule {
	return Rule{
		Code:     "WHC022",
		Severity: SeverityWarning,
		Message:  "Use typed query builders instead of raw map literals",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			if query.Style.HasRawMapLiteral {
				return []LintResult{
					{
						Rule:     "WHC022",
						Severity: SeverityWarning,
						Message:  "Query contains raw map literals - use typed query builders for better type safety",
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}

// WHC023DeeplyNestedConfiguration detects deeply nested query configurations
// that should be flattened for better readability.
func WHC023DeeplyNestedConfiguration() Rule {
	return Rule{
		Code:     "WHC023",
		Severity: SeverityWarning,
		Message:  "Prevent deeply nested query configurations",
		Check: func(query discovery.DiscoveredQuery) []LintResult {
			// Threshold: nesting deeper than 4 levels is hard to read
			const maxDepth = 4

			if query.Style.MaxNestingDepth > maxDepth {
				return []LintResult{
					{
						Rule:     "WHC023",
						Severity: SeverityWarning,
						Message:  fmt.Sprintf("Query has nesting depth of %d - consider flattening configuration (max recommended: %d)", query.Style.MaxNestingDepth, maxDepth),
						File:     query.File,
						Line:     query.Line,
						Query:    query.Name,
					},
				}
			}
			return nil
		},
	}
}
