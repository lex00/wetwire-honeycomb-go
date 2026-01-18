package domain

import (
	"os"
	"testing"

	coredomain "github.com/lex00/wetwire-core-go/domain"
)

func TestHoneycombDomainImplementsInterface(t *testing.T) {
	// Compile-time check that HoneycombDomain implements Domain
	var _ coredomain.Domain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainImplementsListerDomain(t *testing.T) {
	// Compile-time check that HoneycombDomain implements ListerDomain
	var _ coredomain.ListerDomain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainImplementsGrapherDomain(t *testing.T) {
	// Compile-time check that HoneycombDomain implements GrapherDomain
	var _ coredomain.GrapherDomain = (*HoneycombDomain)(nil)
}

func TestHoneycombDomainName(t *testing.T) {
	d := &HoneycombDomain{}
	if d.Name() != "honeycomb" {
		t.Errorf("expected name 'honeycomb', got %q", d.Name())
	}
}

func TestHoneycombDomainVersion(t *testing.T) {
	d := &HoneycombDomain{}
	v := d.Version()
	if v == "" {
		t.Error("version should not be empty")
	}
}

func TestHoneycombDomainBuilder(t *testing.T) {
	d := &HoneycombDomain{}
	b := d.Builder()
	if b == nil {
		t.Error("builder should not be nil")
	}
}

func TestHoneycombDomainLinter(t *testing.T) {
	d := &HoneycombDomain{}
	l := d.Linter()
	if l == nil {
		t.Error("linter should not be nil")
	}
}

func TestHoneycombDomainInitializer(t *testing.T) {
	d := &HoneycombDomain{}
	i := d.Initializer()
	if i == nil {
		t.Error("initializer should not be nil")
	}
}

func TestHoneycombDomainValidator(t *testing.T) {
	d := &HoneycombDomain{}
	v := d.Validator()
	if v == nil {
		t.Error("validator should not be nil")
	}
}

func TestHoneycombDomainLister(t *testing.T) {
	d := &HoneycombDomain{}
	l := d.Lister()
	if l == nil {
		t.Error("lister should not be nil")
	}
}

func TestHoneycombDomainGrapher(t *testing.T) {
	d := &HoneycombDomain{}
	g := d.Grapher()
	if g == nil {
		t.Error("grapher should not be nil")
	}
}

func TestCreateRootCommand(t *testing.T) {
	cmd := CreateRootCommand(&HoneycombDomain{})
	if cmd == nil {
		t.Fatal("root command should not be nil")
	}
	if cmd.Use != "wetwire-honeycomb" {
		t.Errorf("expected Use 'wetwire-honeycomb', got %q", cmd.Use)
	}
}

// Tests for LintOpts.Fix and LintOpts.Disable support

func TestLinterLint_WithFixOption(t *testing.T) {
	d := &HoneycombDomain{}
	linter := d.Linter()

	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create a simple Go file with a query
	queryContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var TestQuery = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
}
`
	if err := os.WriteFile(tmpDir+"/queries.go", []byte(queryContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ctx := &coredomain.Context{}

	// Test with Fix=true
	result, err := linter.Lint(ctx, tmpDir, LintOpts{Fix: true})
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}

	// When Fix is requested but not implemented, the result should indicate this
	// The result message should contain information about Fix mode
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// If there are no lint issues, the result message should be "No lint issues found"
	// or if Fix mode was requested, it should indicate that
	if result.Message == "" {
		t.Error("Expected result to have a message")
	}
}

func TestLinterLint_WithDisableOption(t *testing.T) {
	d := &HoneycombDomain{}
	linter := d.Linter()

	// Create a temporary directory with test files that will trigger lint errors
	tmpDir := t.TempDir()

	// Create a Go file with a query missing dataset (triggers WHC001)
	queryContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var TestQuery = query.Query{
	Dataset:   "",  // Missing dataset - triggers WHC001
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
}
`
	if err := os.WriteFile(tmpDir+"/queries.go", []byte(queryContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ctx := &coredomain.Context{}

	// First, verify that WHC001 is triggered without Disable
	result, err := linter.Lint(ctx, tmpDir, LintOpts{})
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}

	// Check if WHC001 is in the errors
	hasWHC001 := false
	for _, e := range result.Errors {
		if e.Code == "WHC001" {
			hasWHC001 = true
			break
		}
	}
	if !hasWHC001 {
		t.Error("Expected WHC001 error for missing dataset")
	}

	// Now test with Disable=["WHC001"]
	result, err = linter.Lint(ctx, tmpDir, LintOpts{Disable: []string{"WHC001"}})
	if err != nil {
		t.Fatalf("Lint with Disable failed: %v", err)
	}

	// WHC001 should NOT be in the errors
	for _, e := range result.Errors {
		if e.Code == "WHC001" {
			t.Error("WHC001 should be disabled but was found in errors")
		}
	}
}

func TestLinterLint_DisableMultipleRules(t *testing.T) {
	d := &HoneycombDomain{}
	linter := d.Linter()

	tmpDir := t.TempDir()

	// Create a query with multiple lint issues
	queryContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var TestQuery = query.Query{
	Dataset:      "",  // Missing dataset - WHC001
	TimeRange:    query.TimeRange{},  // Missing time range - WHC002
	Calculations: []query.Calculation{},  // Empty calculations - WHC003
}
`
	if err := os.WriteFile(tmpDir+"/queries.go", []byte(queryContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ctx := &coredomain.Context{}

	// Disable WHC001 and WHC002
	result, err := linter.Lint(ctx, tmpDir, LintOpts{Disable: []string{"WHC001", "WHC002"}})
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}

	// Neither WHC001 nor WHC002 should be present
	for _, e := range result.Errors {
		if e.Code == "WHC001" {
			t.Error("WHC001 should be disabled")
		}
		if e.Code == "WHC002" {
			t.Error("WHC002 should be disabled")
		}
	}

	// WHC003 should still be present (not disabled)
	hasWHC003 := false
	for _, e := range result.Errors {
		if e.Code == "WHC003" {
			hasWHC003 = true
			break
		}
	}
	if !hasWHC003 {
		t.Error("Expected WHC003 to still be present since it wasn't disabled")
	}
}

func TestLinterLint_FixAndDisableTogether(t *testing.T) {
	d := &HoneycombDomain{}
	linter := d.Linter()

	tmpDir := t.TempDir()

	// Create a simple valid query
	queryContent := `package queries

import "github.com/lex00/wetwire-honeycomb-go/query"

var TestQuery = query.Query{
	Dataset:   "production",
	TimeRange: query.Hours(1),
	Calculations: []query.Calculation{
		query.Count(),
	},
}
`
	if err := os.WriteFile(tmpDir+"/queries.go", []byte(queryContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	ctx := &coredomain.Context{}

	// Use both Fix and Disable options
	result, err := linter.Lint(ctx, tmpDir, LintOpts{
		Fix:     true,
		Disable: []string{"WHC004"},
	})
	if err != nil {
		t.Fatalf("Lint failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}
}
