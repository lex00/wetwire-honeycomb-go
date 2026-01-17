package lint

import (
	"testing"

	"github.com/lex00/wetwire-honeycomb-go/internal/discovery"
)

// WHC013 Sensitive Column Exposure Tests

func TestLintQueries_WHC013_SensitiveColumn_Password(t *testing.T) {
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
			Breakdowns: []string{"user_password"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for password column in breakdown")
	}

	result := findResult(results, "WHC013")
	if result.Severity != "warning" {
		t.Errorf("Expected warning severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC013_SensitiveColumn_SSN(t *testing.T) {
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
			Breakdowns: []string{"user_ssn"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for SSN column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_CreditCard(t *testing.T) {
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
			Breakdowns: []string{"credit_card_number"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for credit card column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_ApiKey(t *testing.T) {
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
			Breakdowns: []string{"user_api_key"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for api_key column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_SocialSecurity(t *testing.T) {
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
			Breakdowns: []string{"social_security_number"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for social_security column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_NoSensitive(t *testing.T) {
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
			Breakdowns: []string{"endpoint", "service", "status_code"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC013") {
		t.Error("Did not expect WHC013 warning for non-sensitive columns")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_CVV(t *testing.T) {
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
			Breakdowns: []string{"card_cvv"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for CVV column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_AccessToken(t *testing.T) {
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
			Breakdowns: []string{"user_access_token"},
			Limit:      10,
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC013") {
		t.Error("Expected WHC013 warning for access_token column in breakdown")
	}
}

func TestLintQueries_WHC013_SensitiveColumn_NoBreakdowns(t *testing.T) {
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
			Breakdowns: []string{}, // No breakdowns
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC013") {
		t.Error("Did not expect WHC013 warning when there are no breakdowns")
	}
}

// WHC014 Hardcoded Credentials Tests

func TestLintQueries_WHC014_HardcodedCredentials_PasswordInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "mydb?password=secret123",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for password in dataset name")
	}

	result := findResult(results, "WHC014")
	if result.Severity != "error" {
		t.Errorf("Expected error severity, got %s", result.Severity)
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_TokenInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "api-service?token=abc123xyz",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for token in dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_KeyInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production?key=myapikey123",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for key in dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_SecretInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "database?secret=supersecret",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for secret in dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_ApiKeyInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "service?api_key=abcdef12345",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for api_key in dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_CleanDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "production-api-events",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC014") {
		t.Error("Did not expect WHC014 error for clean dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_AuthInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "service?auth=mytoken",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for auth in dataset name")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_EmptyDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "", // Empty dataset (will trigger WHC001, but not WHC014)
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if hasResult(results, "WHC014") {
		t.Error("Did not expect WHC014 error for empty dataset")
	}
}

func TestLintQueries_WHC014_HardcodedCredentials_AccessKeyInDataset(t *testing.T) {
	queries := []discovery.DiscoveredQuery{
		{
			Name:      "TestQuery",
			Package:   "test",
			File:      "/test/file.go",
			Line:      10,
			Dataset:   "s3-bucket?access_key=AKIAIOSFODNN7EXAMPLE",
			TimeRange: discovery.TimeRange{TimeRange: 3600},
			Calculations: []discovery.Calculation{
				{Op: "COUNT"},
			},
		},
	}

	results := LintQueries(queries)

	if !hasResult(results, "WHC014") {
		t.Error("Expected WHC014 error for access_key in dataset name")
	}
}

// Test AllRules includes the new rules
func TestAllRules_IncludesSecretDetectionRules(t *testing.T) {
	rules := AllRules()

	ruleMap := make(map[string]bool)
	for _, r := range rules {
		ruleMap[r.Code] = true
	}

	expectedRules := []string{"WHC012", "WHC013", "WHC014"}
	for _, code := range expectedRules {
		if !ruleMap[code] {
			t.Errorf("Expected AllRules to include %s", code)
		}
	}
}

func TestAllRules_Count(t *testing.T) {
	rules := AllRules()
	// Should have 18 rules now (WHC001-WHC014, WHC020-WHC023)
	if len(rules) != 18 {
		t.Errorf("Expected 18 rules, got %d", len(rules))
	}
}
