package lint

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discover"
)

// WHC012 Secret Detection Tests

func TestLintQueries_WHC012_SecretInFilter_OpenAIKey(t *testing.T) {
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
				{Column: "api_token", Op: "=", Value: "sk-1234567890abcdef"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for OpenAI-style key in filter")
	}

	result := findResult(results, "WHC012")
	if result.Severity != SeverityError {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC012_SecretInFilter_BearerToken(t *testing.T) {
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
				{Column: "auth_header", Op: "=", Value: "Bearer abc123xyz"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for bearer token in filter")
	}
}

func TestLintQueries_WHC012_SecretInFilter_Password(t *testing.T) {
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
				{Column: "user_input", Op: "=", Value: "password123secret"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for password in filter")
	}
}

func TestLintQueries_WHC012_SecretInFilter_ApiKey(t *testing.T) {
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
				{Column: "header", Op: "=", Value: "api_key=xyz789"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for api_key in filter")
	}
}

func TestLintQueries_WHC012_SecretInFilter_NoSecret(t *testing.T) {
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
				{Column: "status", Op: "=", Value: "error"},
				{Column: "service", Op: "=", Value: "api-gateway"},
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC012") {
		t.Error("Did not expect WHC012 error for non-secret filter values")
	}
}

func TestLintQueries_WHC012_SecretInFilter_NonStringValue(t *testing.T) {
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
				{Column: "count", Op: ">", Value: 100}, // Non-string value
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC012") {
		t.Error("Did not expect WHC012 error for non-string filter values")
	}
}

func TestLintQueries_WHC012_SecretInFilter_StripeLiveKey(t *testing.T) {
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
				{Column: "payment_key", Op: "=", Value: "sk_live_abcdefghijklmnop"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for Stripe live key in filter")
	}
}

func TestLintQueries_WHC012_SecretInFilter_AccessToken(t *testing.T) {
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
				{Column: "oauth", Op: "=", Value: "auth_token_abc123"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC012") {
		t.Error("Expected WHC012 error for auth_token in filter")
	}
}
